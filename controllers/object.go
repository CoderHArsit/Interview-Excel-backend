package controllers

type ExpertSignUpRequest struct {
	FullName   string `json:"full_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Phone      string `json:"phone" binding:"required"`
	Expertise  string `json:"expertise" binding:"required"`  // e.g., UPSC, IT, Banking
	Bio        string `json:"bio"`                           // optional
	Experience int    `json:"experience" binding:"required"` // in years
}

type ExpertSignInRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// For binding signup request
type StudentSignUpRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// For binding signin request
type StudentSignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
