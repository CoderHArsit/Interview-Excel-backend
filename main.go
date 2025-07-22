package main

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/routes"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.InitDB()
	config.GoogleConfig()
	config.InitRedis()
	config.InitRazorpay()

	routes.RegisterExpertRoutes(r)
	banner := `
		

  __ __ _ ____ ____ ____ _  _ __ ____ _  _    ____ _  _ ___ ____ __      ____  __   ___ __ _ ____ __ _ ____ 
 (  (  ( (_  _(  __(  _ / )( (  (  __/ )( \  (  __( \/ / __(  __(  )    (  _ \/ _\ / __(  / (  __(  ( (    \
  )(/    / )(  ) _) )   \ \/ /)( ) _)\ /\ /   ) _) )  ( (__ ) _)/ (_/\   ) _ /    ( (__ )  ( ) _)/    /) D (
 (__\_)__)(__)(____(__\_)\__/(__(____(_/\_)  (____(_/\_\___(____\____/  (____\_/\_/\___(__\_(____\_)__(____/


		Welcome to InterviewExcel Backend - Powered by Go and PostgreSQL!
		`
	color.Green(banner)

	r.Run()

}
