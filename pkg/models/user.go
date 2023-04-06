package models

import "time"

// UserLoginInput - структура запроса в апи для логина
type UserLoginInput struct {
	Username string
	Password string
}

// UserLoginOutput - структура ответа апи для логина
type UserLoginOutput struct {
	Token string `json:"Token"`
}

// UserSignupInput - структура запроса в апи для регистрации
type UserSignupInput struct {
	Username string `json:"Username" validate:"required,max=60"`
	Password string `json:"Password" validate:"required,min=8,max=60"`
	Email    string `json:"Email" validate:"required,email,max=60"`
	Phone    string `json:"Phone"`
	Country  string `json:"Country"`
}

// UserUpdateInput - структура запроса в апи для обновления некоторых данных
type UserUpdateInput struct {
	Email   string `json:"Email" validate:"email"`
	Phone   string `json:"Phone"`
	Country string `json:"Country"`
}

// UpdateUserPasswordInput - структура запроса в апи для обновления пароля
type UpdateUserPasswordInput struct {
	NewPassword string `json:"NewPassword" validate:"required,min=8,max=60"`
}

// UserOutput - вью апи для получения данных о пользователе
type UserOutput struct {
	UserUID    string
	Username   string
	Hash       string `json:"-"`
	Email      string
	Phone      string
	CountryUID string
	Country    string
	HistoryUID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}
