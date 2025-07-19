package routes

import (
	"interviewexcel-backend-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(router *gin.Engine) {
	router.POST("/student/signup", controllers.StudentSignUp)
	router.POST("/student/signin", controllers.StudentSignIn)

}
