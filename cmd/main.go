package main

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/routes"
	"log"
	"fmt"
)

func main() {
	r := routes.Routes()
	
	port := 8080

	log.Printf("Server is running on port: %d", port)
	r.Run(fmt.Sprintf(":%d", port))
}