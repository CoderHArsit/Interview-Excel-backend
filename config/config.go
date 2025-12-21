package config

import (
	"fmt"
	"interviewexcel-backend-go/models"
	"log"
	"os"

	"golang.org/x/oauth2/google"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/razorpay/razorpay-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var RazorpayClient *razorpay.Client

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func GoogleConfig() *oauth2.Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  "http://localhost:8080/google_callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	log.Println("This is correct")
	return &AppConfig.GoogleLoginConfig
}

var DB *gorm.DB
func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ensure logrus level allows SQL logs
	logrus.SetLevel(logrus.InfoLevel)

	dbURI := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: NewGormLogger(), 
	})
	if err != nil {
		log.Fatal("Error in connecting the DB: ", err)
	}

	// Auto migration
	if err := db.AutoMigrate(models.GetMigrationModel()...); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	DB = db
	log.Println("Database connected successfully")
}



func InitRazorpay() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file for Razorpay")
	}

	key := os.Getenv("RAZORPAY_KEY")
	secret := os.Getenv("RAZORPAY_SECRET")

	if key == "" || secret == "" {
		log.Fatal("Missing Razorpay credentials in .env")
	}

	RazorpayClient = razorpay.NewClient(key, secret)
	log.Println("Razorpay client initialized")
}
