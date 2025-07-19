package routes

import (
	"interviewexcel-backend-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterExpertRoutes(router *gin.Engine) {
	router.POST("/expert/signup", controllers.ExpertSignUp)
	router.POST("/expert/signin", controllers.ExpertSignin)
	router.POST("/expert/generate-slots", controllers.GenerateWeeklyAvailability)
	router.POST("/expert/google-auth", controllers.ExpertGoogleAuth)
	router.POST("/book-slot", controllers.BookSlot)
}
