package routes

import (
	"interviewexcel-backend-go/controllers"
	"interviewexcel-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(r *gin.Engine) {
	studentRoutes := r.Group("/student")
	studentRoutes.Use(middleware.AuthMiddleware()) 

	studentRoutes.GET("/profile", controllers.GetStudentProfile)
	studentRoutes.PUT("/profile", controllers.UpdateStudentProfile)
	studentRoutes.POST("/book-slot", controllers.BookAvailabilitySlotHandler)
	studentRoutes.GET("/experts", controllers.GetAllExpertsHandler)
	studentRoutes.GET("/expert/:id/slots", controllers.GetAvailableSlotsForExpertHandler)
	studentRoutes.GET("/bookings", controllers.GetStudentBookingsHandler)
	studentRoutes.POST("/preview-slot", controllers.PreviewSlotForPaymentHandler)

}
