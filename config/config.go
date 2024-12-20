package config

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	BaseURL       string
	ServiceKey    string
	JSONFileRoute string
	AIModelURL    string
)

var FirebaseApp *firebase.App

// Firebase 초기화
func InitFirebase() {
	JSONFileRoute = os.Getenv("JSON_FILE_ROUTE")
	// 서비스 계정 JSON 파일 경로
	opt := option.WithCredentialsFile(JSONFileRoute)

	// Firebase 앱 초기화
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase 초기화 실패: %v\n", err)
	}

	FirebaseApp = app
}

// Initialize 환경 변수 로드 및 기본 설정 초기화
func InitEnv() {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// 환경 변수 로드
	BaseURL = os.Getenv("BASE_URL")
	ServiceKey = os.Getenv("SERVICE_KEY")
	JSONFileRoute = os.Getenv("JSON_FILE_ROUTE")
	AIModelURL = os.Getenv("AI_MODEL_URL")

	if BaseURL == "" {
		log.Fatalf("BASE_URL is not set in environment variables")
	}

	if ServiceKey == "" {
		log.Fatalf("SERVICE_KEY is not set in environment variables")
	}

	if JSONFileRoute == "" {
		log.Fatalf("JSON_FILE_ROUTE is not set in environment variables")
	}

	if AIModelURL == "" {
		log.Fatalf("AI_MODEL_URL is not set in environment variables")
	}
}
