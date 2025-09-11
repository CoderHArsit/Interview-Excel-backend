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

type StudentProfile struct {
	UserID string `json:"user_uuid"`
	Role   string `json:"role"`

	// from UserRepo
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`

	Bio          string    `json:"bio,omitempty"`
	Sessions     string    `json:"sessions"`
	Points       string    `json:"points"`
	PreparingFor string    `json:"preparing_for"`
	DateOfBirth  time.Time `json:"dob"`
	City         string    `json:"city"`
	AboutMe      string    `json:"about_me"`
	Skills       []string  `gorm:"type:json" json:"skills"` // JSON column for skills
}

type ExpertProfile struct {
	// From User
	UserID   uint    `json:"id"`
	UserUUID string  `json:"user_uuid"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Picture  string  `json:"picture"`
	Phone    *string `json:"phone"`
	Role     string  `json:"role"`
	City     string  `json:"city"`

	// From Expert
	Bio                string    `json:"bio"`
	DOB                time.Time `json:"dob"`
	Languages          []string  `json:"languages"`
	Specializations    []string  `json:"specializations"`
	Expertise          string    `json:"expertise"`
	Education          string    `json:"education"`
	ExperienceYears    int       `json:"experience_years"`
	ProfilePictureUrl  string    `json:"profile_picture_url"`
	FeesPerSession     int       `json:"fees_per_session"`
	Rating             float64   `json:"rating"` // if you added
	TotalSessions      int       `json:"total_sessions"`
	VerificationStatus string    `json:"verification_status"`
	IsAvailable        bool      `json:"is_available"`
	StudentMentored    int64     `json:"student_mentored"`
}
