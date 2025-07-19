package routes

import (
	"interviewexcel-backend-go/controllers"

	"github.com/gin-gonic/gin"
)


func GoogleLoginRoutes(router *gin.Engine){
router.GET("/auth/google/login", controllers.GoogleLoginHandler)
router.GET("/auth/google/callback", controllers.GoogleCallbackHandler)

}