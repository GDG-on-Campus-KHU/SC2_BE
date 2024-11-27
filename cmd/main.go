package main

import (
	"context"
	"fmt"
	"github.com/GDG-on-Campus-KHU/SC2_BE/config"
	"github.com/GDG-on-Campus-KHU/SC2_BE/controllers"
	"github.com/GDG-on-Campus-KHU/SC2_BE/routes"
	"log"
)

func main() {

	// Firebase 초기화
	config.InitFirebase()

	go func() {
		log.Println("Starting disaster message polling...")
		controllers.GetDisasterMessagesHandler(context.TODO())
	}()

	//Gin 서버 설정
	r := routes.Routes()
	port := 8080

	log.Printf("Server is running on port: %d", port)
	r.Run(fmt.Sprintf(":%d", port))

}
