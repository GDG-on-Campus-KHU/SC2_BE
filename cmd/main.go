package main

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/routes"
	"github.com/joho/godotenv"
	"log"
	"fmt"
)

func main() {
	if err := godotenv.Load(); err != nil{
		log.Fatal("Error loading .env file")
	}
	r := routes.Routes()
	
	port := 8080

	log.Printf("Server is running on port: %d", port)
	r.Run(fmt.Sprintf(":%d", port))
}