package models

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(user User) error
}
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=16"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SuccessResponse представляет успешный ответ от API.
type SuccessResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

// TokenResponse представляет ответ с токеном.
type TokenResponse struct {
	Token string `json:"token"`
}

// ErrorResponse представляет ошибку в ответе от API.
type ErrorResponse struct {
	Error string `json:"error"`
}
