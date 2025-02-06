package main

import (
	"log"

	wire "github.com/ctf-api/internal/di"
)

func main() {
	// Initialize Database
	// db := database.InitDB()

	// Initialize Server with Dependency Injection
	server := wire.InitializeServer()

	// Start Server
	log.Println("🚀 Server running on port 3000")
	if err := server.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
