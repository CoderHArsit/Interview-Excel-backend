package routes

import (
	"interviewexcel-backend-go/controllers"
	"interviewexcel-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(r *gin.Engine) {
	studentRoutes := r.Group("/student")
	studentRoutes.Use(middleware.AuthMiddleware()) // already created earlier

	studentRoutes.GET("/profile", controllers.GetStudentProfile)
	studentRoutes.PUT("/profile", controllers.UpdateStudentProfile)
	studentRoutes.POST("/book-slot", controllers.BookAvailabilitySlotHandler)

}
