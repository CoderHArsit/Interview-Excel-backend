package routes

import (
	"interviewexcel-backend-go/controllers"

	"github.com/gin-gonic/gin"
)

func GoogleLoginRoutes(router *gin.Engine) {
	router.POST("/auth/google/login", controllers.UserGoogleAuth)
	router.GET("/auth/google/callback", controllers.GoogleCallbackHandler)
}
