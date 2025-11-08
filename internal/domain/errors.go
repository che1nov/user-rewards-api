package domain

import "errors"

// Доменные ошибки
var (
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrUserExists        = errors.New("пользователь уже существует")
	ErrInvalidUsername   = errors.New("некорректный username")
	ErrInvalidEmail      = errors.New("некорректный email")
	ErrTaskNotFound      = errors.New("задание не найдено")
	ErrTaskAlreadyExists = errors.New("задание уже выполнено")
	ErrInvalidTaskType   = errors.New("неизвестный тип задания")
	ErrReferralExists    = errors.New("реферальный код уже использован")
	ErrSelfReferral      = errors.New("нельзя использовать свой собственный реферальный код")
	ErrReferrerNotFound  = errors.New("реферер не найден")
)

