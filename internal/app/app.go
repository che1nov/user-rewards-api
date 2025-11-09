package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"user-rewards-api/internal/adapters/postgresql"
	"user-rewards-api/internal/config"
	httpController "user-rewards-api/internal/controllers/http"
	"user-rewards-api/internal/database"
	authMiddleware "user-rewards-api/internal/middleware"
	"user-rewards-api/internal/usecases"
)

type App struct {
	config *config.Config
	db     *sql.DB
	router *gin.Engine
	server *http.Server
}

// NewApp создает новое приложение
func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка ping базы данных: %w", err)
	}

	slog.Info("Подключение к базе данных установлено")

	migrationsPath := "migrations"
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка выполнения миграций: %w", err)
	}

	slog.Info("Миграции выполнены успешно")

	sqlxDB := sqlx.NewDb(db, "postgres")

	postgresAdapter := postgresql.NewPostgreSQLAdapter(sqlxDB)

	createUserUC := usecases.NewCreateUserUseCase(postgresAdapter, cfg.JWTSecret)
	getUserStatusUC := usecases.NewGetUserStatusUseCase(postgresAdapter)
	getLeaderboardUC := usecases.NewGetLeaderboardUseCase(postgresAdapter)
	completeTaskUC := usecases.NewCompleteTaskUseCase(postgresAdapter)
	processReferralUC := usecases.NewProcessReferralUseCase(postgresAdapter)

	userController := httpController.NewUserController(
		createUserUC,
		getUserStatusUC,
		getLeaderboardUC,
		completeTaskUC,
		processReferralUC,
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.POST("/users", userController.CreateUser)

	protected := router.Group("")
	protected.Use(authMiddleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/users/leaderboard", userController.GetLeaderboard)
		protected.GET("/users/:id/status", userController.GetUserStatus)
		protected.POST("/users/:id/task/complete", userController.CompleteTask)
		protected.POST("/users/:id/referrer", userController.ProcessReferral)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		config: cfg,
		db:     db,
		router: router,
		server: server,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	go func() {
		slog.Info("Сервер запущен", "port", a.config.ServerPort)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка запуска сервера", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка при остановке сервера: %w", err)
	}

	slog.Info("Сервер остановлен")
	return nil
}

// Close закрывает ресурсы приложения
func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}
