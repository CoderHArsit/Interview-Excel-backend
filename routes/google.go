package routes

import (
	"interviewexcel-backend-go/controllers"

	"github.com/gin-gonic/gin"
)


func GoogleLoginRoutes(router *gin.Engine){
	router.GET("/google_login", controllers.GoogleLogin)
    //app.Post("/google_callback", controllers.GoogleCallback)
}