package controllers

import "time"

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

type StudentProfileResponse struct {
	UserID string `json:"user_uuid"`
	Role   string `json:"role"`

	// from UserRepo
	Email string `json:"email"`
	Phone string `json:"phone"`

	Bio          string    `json:"bio,omitempty"`
	Sessions     string    `json:"sessions"`
	Points       string    `json:"points"`
	PreparingFor string    `json:"preparing_for"`
	DateOfBirth  time.Time `json:"dob"`
	City         string    `json:"city"`
	AboutMe      string    `json:"about_me"`
	Skills       []string  `gorm:"type:json" json:"skills"` // JSON column for skills
}
