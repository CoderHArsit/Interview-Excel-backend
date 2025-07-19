package controllers

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StudentSignUp(c *gin.Context) {
	var (
		request     StudentSignUpRequest
		studentRepo = models.InitStudentRepo(config.DB)
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	err = studentRepo.Create(&models.Student{
		FullName: request.FullName,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: hashedPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student registered successfully"})
}

func StudentSignIn(c *gin.Context) {
	var (
		request     StudentSignInRequest
		studentRepo = models.InitStudentRepo(config.DB)
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := studentRepo.GetByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := VerifyPassword(student.Password, request.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	accessToken, err := generateAccessToken(student.ID, "student")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	refreshToken, err := generateRefreshToken(student.ID, "student")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
