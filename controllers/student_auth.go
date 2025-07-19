package controllers

import (
	"encoding/json"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"io"
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

func StudentGoogleAuth(c *gin.Context) {
	var req GoogleAuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Google auth code"})
		return
	}

	// Step 1: Exchange code for token
	token, err := config.AppConfig.GoogleLoginConfig.Exchange(c, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Step 2: Get user info from Google
	client := config.AppConfig.GoogleLoginConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info from Google"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Google response"})
		return
	}

	var gUser struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.Unmarshal(body, &gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Google user info"})
		return
	}

	if !gUser.VerifiedEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified by Google"})
		return
	}

	// Step 3: Check if student already exists
	var student models.Student
	err = config.DB.Where("email = ?", gUser.Email).First(&student).Error

	if err != nil {
		// New student → create
		student = models.Student{
			Email:    gUser.Email,
			FullName: gUser.Name,
			Picture:  gUser.Picture,
		}
		if err := config.DB.Create(&student).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
			return
		}
	}

	// Step 4: Generate tokens
	accessToken, err := generateAccessToken(student.ID, "student")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	refreshToken, err := generateRefreshToken(student.ID, "student")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	// Step 5: Respond
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"student": gin.H{
			"id":     student.ID,
			"name":   student.FullName,
			"email":  student.Email,
			"avatar": student.Picture,
		},
	})
}
