package controllers

import (
	"encoding/json"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"interviewexcel-backend-go/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExpertSignUp(c *gin.Context) {
	var (
		request    ExpertSignUpRequest
		expertRepo = models.InitExpertRepo(config.DB)
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error binding request": err.Error()})
		return
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

	accessToken, err := generateAccessToken(expert.ID, "expert")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	refreshToken, err := generateRefreshToken(expert.ID, "expert")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

func GenerateWeeklyAvailability(c *gin.Context) {
	var input struct {
		ExpertID uint `json:"expert_id" binding:"required"` //TODO - remove this api, expert id should be extracted from the logged in expert token
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slots := utils.GenerateWeeklySlots(input.ExpertID)

	if err := config.DB.Create(&slots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create availability slots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Availability slots created successfully"})
}



func ExpertGoogleAuth(c *gin.Context) {
	var req GoogleAuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var gUser struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.Unmarshal(body, &gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	if !gUser.VerifiedEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Google account email not verified"})
		return
	}

	// Step 3: Check if expert exists
	var expert models.Expert
	err = config.DB.Where("email = ?", gUser.Email).First(&expert).Error

	if err != nil {
		// If not found, create a new expert
		expert = models.Expert{
			Email:    gUser.Email,
			FullName: gUser.Name,
			Picture:  gUser.Picture,
		}
		if err := config.DB.Create(&expert).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expert"})
			return
		}
	}

	// Step 4: Generate JWT tokens
	accessToken, err := generateAccessToken(expert.ID, "expert")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	refreshToken, err := generateRefreshToken(expert.ID, "expert")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expert": gin.H{
			"id":     expert.ID,
			"name":   expert.FullName,
			"email":  expert.Email,
			"avatar": expert.Picture,
		},
	})
}
