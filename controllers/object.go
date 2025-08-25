package controllers

type SignUpRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone"`
	Password        string `json:"password" binding:"required"`
	ConfirrPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	Role            string `json:"role" binding:"required,oneof=student expert"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleAuthRequest struct {
	Role  string `json:"role" binding:"required"`
	Token string `json:"token" binding:"required"`
}
