package models

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User        User        `json:"user"`
	Credentials AccessToken `json:"credentials"`
}
