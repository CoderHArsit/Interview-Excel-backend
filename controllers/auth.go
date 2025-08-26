package controllers

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	logger "interviewexcel-backend-go/pkg/errors"
	utils "interviewexcel-backend-go/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func Signup(c *gin.Context) {
	var req SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// check if user already exists
	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		logger.Errorf("User already exists: %v\n", existing.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// hash password if provided
	var password *string
	if req.Password != "" {
		hash, _ := utils.HashPassword(req.Password)
		password = &hash
	}

	user := models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: password,
		Role:     req.Role, // "student" or "expert"
	}

	if err := config.DB.Create(&user).Error; err != nil {
		logger.Errorf("Error creating user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 🔹 Create respective profile
	switch req.Role {
	case "student":
		student := models.Student{UserID: user.ID}
		if err := config.DB.Create(&student).Error; err != nil {
			logger.Errorf("Error creating student profile: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student profile"})
			return
		}
	case "expert":
		expert := models.Expert{UserID: user.ID}
		if err := config.DB.Create(&expert).Error; err != nil {
			logger.Errorf("Error creating expert profile: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expert profile"})
			return
		}
	}

	// 🔹 Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		logger.Errorf("Error generating access token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		logger.Errorf("Error generating refresh token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie(
		"refresh_token", // cookie name
		refreshToken,    // value
		7*24*60*60,      // maxAge in seconds
		"/",             // path: available to all endpoints
		"localhost",     // domain: only hostname
		false,           // secure: false for http in dev
		true,            // httpOnly: true is fine
	)

	// Save refresh token in Redis
	err = config.RedisClient.Set(
		config.Ctx,
		refreshToken,   // key
		user.ID,        // value
		7*24*time.Hour, // expiration = 7 days
	).Err()

	if err != nil {
		logger.Errorf("Error saving refresh token to Redis: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"access_token": accessToken,
	})
}

func UserGoogleAuth(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Role != "student" && req.Role != "expert") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Google token or invalid role"})
		return
	}

	// Verify the Google token
	payload, err := idtoken.Validate(c, req.Token, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google token"})
		return
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)

	if !emailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified by Google"})
		return
	}

	// Initialize repos
	userRepo := models.InitUserRepo(config.DB)
	studentRepo := models.InitStudentRepo(config.DB)
	expertRepo := models.InitExpertRepo(config.DB)

	// Check if user exists
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		// New user
		user = &models.User{
			FullName: name,
			Email:    email,
			Picture:  picture,
			Role:     req.Role,
		}

		if err := userRepo.Create(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Create respective profile
		switch req.Role {
		case "student":
			if err := studentRepo.Create(&models.Student{UserID: user.ID}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student profile"})
				return
			}
		case "expert":
			if err := expertRepo.Create(&models.Expert{UserID: user.ID}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expert profile"})
				return
			}
		}
	}

	// Generate access and refresh tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie(
		"refresh_token", // cookie name
		refreshToken,    // value
		7*24*60*60,      // maxAge in seconds
		"/",             // path: available to all endpoints
		"localhost",     // domain: only hostname
		false,           // secure: false for http in dev
		true,            // httpOnly: true is fine
	)

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id":     user.ID,
			"name":   user.FullName,
			"email":  user.Email,
			"role":   user.Role,
			"avatar": user.Picture,
		},
	})
}

func UserSignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		logger.Errorf("User not found: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if user.Password == nil || utils.VerifyPassword(*user.Password, req.Password) != nil {
		logger.Errorf("Invalid password for user: %s\n", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		logger.Errorf("Error generating access token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		logger.Errorf("Error generating refresh token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie(
		"refresh_token", // cookie name
		refreshToken,    // value
		7*24*60*60,      // maxAge in seconds
		"/",             // path: available to all endpoints
		"localhost",     // domain: only hostname
		false,           // secure: false for http in dev
		true,            // httpOnly: true is fine
	)

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Error("refresh token error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}
	logger.Info("refresh toke", cookie)
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.FullName,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func RefreshSession(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Error("refresh token error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}

	claims, err := utils.VerifyRefreshToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Generate new access token
	accessToken, _ := utils.GenerateAccessToken(claims.UserID, claims.Role)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
