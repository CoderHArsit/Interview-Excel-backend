package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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