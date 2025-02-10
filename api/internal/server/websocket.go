package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func SetupSocketServer(router *gin.Engine) {
	server := socketio.NewServer(nil)

	// Handle new WebSocket connections
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("New client connected:", s.ID())
		return nil
	})

	// Handle timer freeze event
	server.OnEvent("/", "freeze-timer", func(s socketio.Conn) {
		fmt.Println("Timer frozen by admin")
		server.BroadcastToNamespace("/", "timer-frozen") // Notify all clients
	})

	// Handle timer resume event
	server.OnEvent("/", "resume-timer", func(s socketio.Conn) {
		fmt.Println("Timer resumed by admin")
		server.BroadcastToNamespace("/", "timer-resumed") // Notify all clients
	})

	// Handle timer start event
	server.OnEvent("/", "start-timer", func(s socketio.Conn) {
		fmt.Println("Timer started")
		server.BroadcastToNamespace("/", "timer-started") // Notify all clients
	})

	// Start Socket.IO server in a separate goroutine
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()

	// Serve WebSocket at /socket.io/
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
}
