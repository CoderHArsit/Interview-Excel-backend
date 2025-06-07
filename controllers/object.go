package controllers

type ExpertSignUpRequest struct {
	FullName   string `json:"full_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Phone      string `json:"phone" binding:"required"`
	Expertise  string `json:"expertise" binding:"required"`  // e.g., UPSC, IT, Banking
	Bio        string `json:"bio"`                            // optional
	Experience int    `json:"experience" binding:"required"` // in years
}
