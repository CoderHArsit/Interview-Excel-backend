package controllers

import (
	"github.com/gin-gonic/gin"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"
)

func ExpertSignUp(c *gin.Context) {
	var (
		request    ExpertSignUpRequest
		expertRepo = models.InitExpertRepo(config.DB)
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//Hash the password
	hashedPassword, err := HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = expertRepo.Create(&models.Expert{
		FullName:        request.FullName,
		Password:        hashedPassword,
		Email:           request.Email,
		Phone:           request.Phone,
		Expertise:       request.Expertise,
		Bio:             request.Bio,
		ExperienceYears: request.Experience,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expert registered SuccessFully"})
}

// Expert Signin handles expert Login
func ExpertSignin(c *gin.Context) {
	var (
		request   ExpertSignInRequest
		experRepo = models.InitExpertRepo(config.DB)
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expert, err := experRepo.GetWithTx(config.DB, &models.Expert{
		Email: request.Email,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = VerifyPassword(expert.Password, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	accessToken, err := generateAccessToken(expert.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	refreshToken, err := generateRefreshToken(expert.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
