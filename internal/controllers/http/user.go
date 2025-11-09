package http

import (
	"log/slog"
	"net/http"

	"user-rewards-api/internal/domain"
	"user-rewards-api/internal/dto"
	"user-rewards-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	createUserUC      *usecases.CreateUserUseCase
	getUserStatusUC   *usecases.GetUserStatusUseCase
	getLeaderboardUC  *usecases.GetLeaderboardUseCase
	completeTaskUC    *usecases.CompleteTaskUseCase
	processReferralUC *usecases.ProcessReferralUseCase
}

func NewUserController(
	createUserUC *usecases.CreateUserUseCase,
	getUserStatusUC *usecases.GetUserStatusUseCase,
	getLeaderboardUC *usecases.GetLeaderboardUseCase,
	completeTaskUC *usecases.CompleteTaskUseCase,
	processReferralUC *usecases.ProcessReferralUseCase,
) *UserController {
	return &UserController{
		createUserUC:      createUserUC,
		getUserStatusUC:   getUserStatusUC,
		getLeaderboardUC:  getLeaderboardUC,
		completeTaskUC:    completeTaskUC,
		processReferralUC: processReferralUC,
	}
}

// CreateUser создает нового пользователя
// POST /users
func (c *UserController) CreateUser(ctx *gin.Context) {
	var input dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		sendError(ctx, domain.ErrInvalidUsername)
		return
	}

	output, err := c.createUserUC.Execute(ctx.Request.Context(), input)
	if err != nil {
		handleError(ctx, err)
		return
	}

	slog.Info("Пользователь создан", "user_id", output.UserID, "username", output.Username)
	ctx.JSON(http.StatusCreated, output)
}

// GetUserStatus получает статус пользователя
// GET /users/:id/status
func (c *UserController) GetUserStatus(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	output, err := c.getUserStatusUC.Execute(ctx.Request.Context(), userIDStr)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetLeaderboard получает таблицу лидеров
// GET /users/leaderboard
func (c *UserController) GetLeaderboard(ctx *gin.Context) {
	limit := 100

	output, err := c.getLeaderboardUC.Execute(ctx.Request.Context(), limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// CompleteTask выполняет задание
// POST /users/:id/task/complete
func (c *UserController) CompleteTask(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	var input dto.CompleteTaskInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		sendError(ctx, domain.ErrInvalidTaskType)
		return
	}

	output, err := c.completeTaskUC.Execute(ctx.Request.Context(), userIDStr, input)
	if err != nil {
		handleError(ctx, err)
		return
	}

	slog.Info("Задание выполнено", "user_id", userIDStr, "task_type", input.TaskType, "points", output.Points, "new_balance", output.NewBalance)
	ctx.JSON(http.StatusOK, output)
}

// ProcessReferral обрабатывает реферальный код
// POST /users/:id/referrer
func (c *UserController) ProcessReferral(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	var input dto.ProcessReferralInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		sendError(ctx, domain.ErrReferrerNotFound)
		return
	}

	output, err := c.processReferralUC.Execute(ctx.Request.Context(), userIDStr, input)
	if err != nil {
		handleError(ctx, err)
		return
	}

	slog.Info("Реферальный код использован", "user_id", userIDStr, "referrer_id", input.ReferrerID, "bonus_points", output.BonusPoints, "new_balance", output.NewBalance)
	ctx.JSON(http.StatusOK, output)
}

// handleError обрабатывает ошибки и возвращает соответствующий HTTP статус
func handleError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	errStr := err.Error()
	switch {
	case err == domain.ErrUserNotFound:
		sendError(ctx, err, http.StatusNotFound)
	case err == domain.ErrUserExists:
		sendError(ctx, err, http.StatusConflict)
	case err == domain.ErrTaskAlreadyExists:
		sendError(ctx, err, http.StatusConflict)
	case err == domain.ErrReferralExists:
		sendError(ctx, err, http.StatusConflict)
	case err == domain.ErrSelfReferral:
		sendError(ctx, err, http.StatusBadRequest)
	case err == domain.ErrReferrerNotFound:
		sendError(ctx, err, http.StatusNotFound)
	case err == domain.ErrInvalidUsername || err == domain.ErrInvalidEmail || err == domain.ErrInvalidTaskType:
		sendError(ctx, err, http.StatusBadRequest)
	default:
		slog.Error("Внутренняя ошибка", "error", err, "error_string", errStr, "path", ctx.Request.URL.Path)
		sendError(ctx, err, http.StatusInternalServerError)
	}
}

// sendError отправляет ошибку в формате JSON
func sendError(ctx *gin.Context, err error, statusCode ...int) {
	code := http.StatusInternalServerError
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	response := map[string]string{
		"error": err.Error(),
	}
	ctx.JSON(code, response)
}
