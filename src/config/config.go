package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Sender            string
	MaxSendPerDay     int
	Host              string
	Port              int
	Encryption        string
	ConnectionTimeout int
	SendTimeout       int
	Auth              bool
	Username          string
	Password          string
}

type DBConfig struct {
	Host       string
	Username   string
	Password   string
	DBName     string
	Collection *Collection
}

type Collection struct {
	AuthAttempt string
	User        string
	Location    string
	Menu        string
	Cart        string
	Email       string
	Chat        string
}

type Bedrock struct {
	Model   string
	Region  string
	Profile string
}

type Config struct {
	Host           string
	Port           string
	FrontendAPIKey string
	DB             *DBConfig
	SMTP           *SMTPConfig
	Bedrock        *Bedrock
}

func LoadEnv() *Config {
	if len(os.Args) < 2 {
		log.Println("load default env")
		godotenv.Load()
	} else {
		file := os.Args[1]
		log.Println("load", file)
		godotenv.Load(file)
	}

	cfg := Config{
		DB:   &DBConfig{},
		SMTP: &SMTPConfig{},
	}

	cfg.Host = os.Getenv("HOST")
	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}

	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	cfg.FrontendAPIKey = os.Getenv("FRONTEND_API_KEY")

	cfg.DB.Host = os.Getenv("MONGO_HOST")
	cfg.DB.Username = os.Getenv("MONGO_USERNAME")
	cfg.DB.Password = os.Getenv("MONGO_PASSWORD")
	cfg.DB.DBName = os.Getenv("MONGO_DATABASE")

	authAttemptCollection := os.Getenv("MONGO_COLLECTION_AUTH_ATTEMPT")
	if authAttemptCollection == "" {
		authAttemptCollection = "login_attempts"
	}
	userCollection := os.Getenv("MONGO_COLLECTION_USER")
	if userCollection == "" {
		userCollection = "users"
	}
	locationCollection := os.Getenv("MONGO_COLLECTION_LOCATION")
	if locationCollection == "" {
		locationCollection = "locations"
	}
	menuCollection := os.Getenv("MONGO_COLLECTION_MENU")
	if menuCollection == "" {
		menuCollection = "menus"
	}
	cartCollection := os.Getenv("MONGO_COLLECTION_CART")
	if cartCollection == "" {
		cartCollection = "carts"
	}
	emailCollection := os.Getenv("MONGO_COLLECTION_EMAIL")
	if emailCollection == "" {
		emailCollection = "emails"
	}
	chatCollection := os.Getenv("MONGO_COLLECTION_CHAT")
	if chatCollection == "" {
		chatCollection = "chats"
	}
	cfg.DB.Collection = &Collection{
		AuthAttempt: authAttemptCollection,
		User:        userCollection,
		Location:    locationCollection,
		Menu:        menuCollection,
		Cart:        cartCollection,
		Email:       emailCollection,
		Chat:        chatCollection,
	}

	cfg.SMTP.Sender = os.Getenv("SMTP_SENDER")
	if cfg.SMTP.Sender == "" {
		log.Fatal("sender is required")
	}

	maxSendPerDay := os.Getenv("SMTP_MAX_SEND_PER_DAY")
	if maxSendPerDay == "" {
		cfg.SMTP.MaxSendPerDay = 20
	} else {
		maxSendPerDayNum, err := strconv.Atoi(maxSendPerDay)
		if err != nil {
			log.Fatal("invalid max send per day value")
		}
		cfg.SMTP.MaxSendPerDay = maxSendPerDayNum
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	if cfg.SMTP.Host == "" {
		cfg.SMTP.Host = "localhost"
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		cfg.SMTP.Port = 1025
	} else {
		portNum, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal("invalid smtp port")
		}
		cfg.SMTP.Port = portNum
	}

	cfg.SMTP.Encryption = os.Getenv("SMTP_ENCRYPTION")

	if os.Getenv("SMTP_AUTH") == "true" {
		cfg.SMTP.Auth = true
		cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
		cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	} else {
		cfg.SMTP.Auth = false
	}

	ct := os.Getenv("SMTP_CONNECTION_TIMEOUT")
	if ct == "" {
		cfg.SMTP.ConnectionTimeout = 10
	} else {
		ctNum, err := strconv.Atoi(ct)
		if err != nil {
			log.Fatal("invalid connection timeout value")
		}
		cfg.SMTP.ConnectionTimeout = ctNum
	}

	st := os.Getenv("SMTP_SEND_TIMEOUT")
	if st == "" {
		cfg.SMTP.SendTimeout = 10
	} else {
		stNum, err := strconv.Atoi(ct)
		if err != nil {
			log.Fatal("invalid send timeout value")
		}
		cfg.SMTP.SendTimeout = stNum
	}

	cfg.Bedrock = &Bedrock{}
	cfg.Bedrock.Model = os.Getenv("BEDROCK_MODEL")
	cfg.Bedrock.Region = os.Getenv("BEDROCK_REGION")
	cfg.Bedrock.Profile = os.Getenv("BEDROCK_PROFILE")

	return &cfg
}
