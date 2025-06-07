package controllers

import (
	"interviewexcel-backend-go/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	// Generate the Google login URL
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	// Redirect the user to the Google login page
	c.Redirect(http.StatusSeeOther, url)
}
