package main

import (
	"log"

	wire "github.com/ctf/api/internal/di"
)

func main() {
	// Initialize Database
	// db := database.InitDB()
	// config, configErr := config.LoadConfig()
	// if configErr != nil {
	// 	log.Fatal("Error in configuration", configErr)
	// }

	// Initialize Server with Dependency Injection
	server, err := wire.InitializeServer()
	if err != nil {
		log.Fatalf("Failed to inject Dependencies :%v", err)
	}

	// Start Server
	log.Println("🚀 Server running on port 3000")
	if err := server.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
