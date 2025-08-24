package main

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/routes"
	"time"

	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize core configs
	config.InitDB()
	config.GoogleConfig()
	config.InitRedis()
	config.InitRazorpay()

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	routes.RegisterExpertRoutes(r)
	routes.RegisterStudentRoutes(r)
	routes.GoogleLoginRoutes(r)

	// Banner
	banner := `

  __ __ _ ____ ____ ____ _  _ __ ____ _  _    ____ _  _ ___ ____ __      ____  __   ___ __ _ ____ __ _ ____ 
 (  (  ( (_  _(  __(  _ / )( (  (  __/ )( \  (  __( \/ / __(  __(  )    (  _ \/ _\ / __(  / (  __(  ( (    \
  )(/    / )(  ) _) )   \ \/ /)( ) _)\ /\ /   ) _) )  ( (__ ) _)/ (_/\   ) _ /    ( (__ )  ( ) _)/    /) D (
 (__\_)__)(__)(____(__\_)\__/(__(____(_/\_)  (____(_/\_\___(____\____/  (____\_/\_/\___(__\_(____\_)__(____/


		Welcome to InterviewExcel Backend - Powered by Go and PostgreSQL!
	`
	color.Green(banner)

	// Start server on :8080 (default)
	r.Run(":8080")
}
