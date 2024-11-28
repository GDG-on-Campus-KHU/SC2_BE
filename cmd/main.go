package main

import (
	"fmt"
	"github.com/GDG-on-Campus-KHU/SC2_BE/config"
	"github.com/GDG-on-Campus-KHU/SC2_BE/routes"
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"github.com/GDG-on-Campus-KHU/SC2_BE/db"
	"github.com/GDG-on-Campus-KHU/SC2_BE/controllers"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	// Firebase 초기화
	config.InitFirebase()

	// 고루틴 실행
	// 백그라운드에서 비동기적으로 작업을 처리하면서도 메인 프로그램의 실행을 방해하지 않기 위해서 사용한다.
	// 현재 서버가 실행된 상태에서 지속적으로 API를 호출하고 메시지를 확인하는 Polling 작업 수행.
	go service.PollForUpdates() // Polling 작업을 비동기로 실행한다.

	// Gin 서버 설정
	if err := godotenv.Load(); err != nil{
		log.Fatal("Error loading .env file")
	}

	// MongoDB 연결
    client, err := db.ConnectDB()
	controllers.SetMongoClient(client)
    if err != nil {
        log.Fatal(err)
    }
    defer db.DisconnectDB(client)

	r := routes.Routes()
	port := 8080

	log.Printf("Server is running on port: %d", port)
	// r.Run(...)은 블로킹 함수로, 서버가 종료되지 않는 한 계속 실행된다.
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
