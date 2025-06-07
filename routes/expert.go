package routes

import (
	"interviewexcel-backend-go/controllers"

	"github.com/gin-gonic/gin"
)


func RegisterExpertRoutes(router *gin.Engine){
	router.POST("/expert/signup",controllers.ExpertSignUp)
	router.POST("/expert/sigin",controllers.ExpertSignin)
}