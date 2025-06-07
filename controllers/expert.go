package controllers

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

// RefreshToken handles the creation of a new access token using a refresh token
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//parse and verify the refresh token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expertID := uint(claims["sub"].(float32))
		newAccessToken, err := generateAccessToken(expertID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": newAccessToken,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
	}
}

// function to generate an access token
func generateAccessToken(expertID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": expertID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Function to generate a refresh token
func generateRefreshToken(expertID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": expertID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7-day expiration
	})
	return token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//VerifyPassword Compares a hashed password with a plain text one

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func ExpertSignUp(c *gin.Context) {
	var expert models.Expert
	log.Info("gedrf")
	err := c.ShouldBindJSON(&expert)

	if err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//Hash the password
	hashedPassword, err := HashPassword(expert.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expert.Password = hashedPassword

	err = config.DB.Create(&expert).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register Expert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expert registered SuccessFully"})
}

//Expert Signin handles expert Login

func ExpertSignin(c *gin.Context) {

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
    log.Info("nvdkjnvkfv")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var expert models.Expert
	err = config.DB.Where("email =?", input.Email).First(&expert).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = VerifyPassword(expert.Password, input.Password)
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
