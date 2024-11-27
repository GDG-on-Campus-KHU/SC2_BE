package config

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"log"
)

var FirebaseApp *firebase.App

// Firebase 초기화
func InitFirebase() {
	// 서비스 계정 JSON 파일 경로
	opt := option.WithCredentialsFile("C:\\Users\\user\\Downloads\\khu-gdg-sc2-firebase-adminsdk-d60ps-1ca08cefb6.json")

	// Firebase 앱 초기화
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase 초기화 실패: %v\n", err)
	}

	FirebaseApp = app
}
