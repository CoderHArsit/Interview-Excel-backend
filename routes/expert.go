package routes

import (
	"interviewexcel-backend-go/controllers"
	"interviewexcel-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterExpertRoutes(router *gin.Engine) {
	// Protected routes
	expertGroup := router.Group("/expert")
	expertGroup.Use(middleware.AuthMiddleware()) // ✅ Apply middleware here

	// authexpertGroup.POST("/generate-slots", controllers.GenerateWeeklyAvailability)
	expertGroup.GET("/profile", controllers.GetExpertProfile)
	expertGroup.GET("/my-slots", controllers.GetAvailableSlotsForExpertHandler)
	expertGroup.PUT("/profile", controllers.UpdateExpertProfile)
	expertGroup.POST("/generate-slots", controllers.GenerateWeeklyAvailability)
	expertGroup.GET("/all-slots", controllers.GetAllSlotsOfExpert)
	// Add more protected expert routes here
}
