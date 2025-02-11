package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ctf/api/internal/handlers"
	"github.com/ctf/api/internal/middleware"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

// Server struct
type Server struct {
	router            *gin.Engine
	handlers          *handlers.Handlers
	teamhandlers      *handlers.TeamHandler
	challengehandlers *handlers.ChallengeHandler
	socketServer      *socketio.Server
}

// NewServer initializes the Gin server with WebSocket support
func NewServer(handlers *handlers.Handlers, teamhandlers *handlers.TeamHandler, challengehandlers *handlers.ChallengeHandler) *Server {
	router := gin.Default()
	socketServer := SetupSocketServer(router)

	server := &Server{
		router:            router,
		handlers:          handlers,
		teamhandlers:      teamhandlers,
		challengehandlers: challengehandlers,
		socketServer:      socketServer,
	}

	server.setupRoutes()
	return server
}

// setupRoutes defines all API routes
func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(CORSMiddleware())

	// Custom template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	s.router.SetFuncMap(funcMap)

	// Load templates and static files
	s.router.LoadHTMLGlob("../../../template/*.html")
	s.router.Static("/css", "../../../template/css")
	s.router.Static("/js", "../../../template/js")
	s.router.Static("/images", "../../../template/images")
	s.router.Static("/fonts", "../../../template/fonts")
	s.router.Static("/admin/css", "../../../template/css")

	// 📌 Public Routes (No Authentication Required)
	s.router.GET("/", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "index.html", nil) })
	s.router.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	s.router.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	s.router.GET("/instructions", func(c *gin.Context) { c.HTML(http.StatusOK, "instructions.html", nil) })
	s.router.GET("/quests", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "quests.html", nil) })
	s.router.GET("/about", func(c *gin.Context) { c.HTML(http.StatusOK, "about.html", nil) })

	// API Endpoints - Public
	s.router.POST("/register", s.teamhandlers.RegisterTeam)
	s.router.POST("/login", s.teamhandlers.LoginTeam)
	s.router.GET("/hackerboard", s.handlers.GetScores)
	s.router.GET("/getchallenges", s.challengehandlers.GetChallenges)

	// 📌 **Authenticated Team Routes (Requires Login)**
	teamRoutes := s.router.Group("/team")
	teamRoutes.Use(middleware.TeamAuth)
	{
		teamRoutes.POST("/verifyflag", s.challengehandlers.SubmitFlag)
	}

	// 📌 **Admin Routes (Requires Admin Role)**
	adminRoutes := s.router.Group("/admin")
	adminRoutes.Use(middleware.AdminAuth)
	{
		adminRoutes.GET("/", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "admin.html", nil) })
		adminRoutes.POST("/addchallenge", s.challengehandlers.CreateChallenge)
		adminRoutes.GET("/challenges", s.challengehandlers.GetChallenges)
		adminRoutes.DELETE("/deletechallenge/:id", s.challengehandlers.DeleteChallenge)
		adminRoutes.GET("/teams", s.handlers.GetTeams)
		adminRoutes.PUT("/blockteam/:id", s.teamhandlers.BlockTeam)
		adminRoutes.PUT("/unblockteam/:id", s.teamhandlers.UnblockTeam)
	}

	// 📌 WebSocket (Real-time updates)
	s.router.GET("/ws", s.handlers.HandleWebSocket)
	s.router.GET("/socket.io/*any", gin.WrapH(s.socketServer))
	s.router.POST("/socket.io/*any", gin.WrapH(s.socketServer))
}

// Run starts the server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// CORSMiddleware handles CORS policy
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// SetupSocketServer initializes WebSocket server
func SetupSocketServer(router *gin.Engine) *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("New client connected:", s.ID())
		return nil
	})

	// Timer controls
	server.OnEvent("/", "freeze-timer", func(s socketio.Conn) {
		log.Println("Timer frozen by admin")
		server.BroadcastToNamespace("/", "timer-frozen")
	})

	server.OnEvent("/", "resume-timer", func(s socketio.Conn) {
		log.Println("Timer resumed by admin")
		server.BroadcastToNamespace("/", "timer-resumed")
	})

	server.OnEvent("/", "start-timer", func(s socketio.Conn) {
		log.Println("Timer started")
		server.BroadcastToNamespace("/", "timer-started")
	})

	server.OnEvent("/", "reset-timer", func(s socketio.Conn) {
		log.Println("Timer reset")
		server.BroadcastToNamespace("/", "timer-reset")
	})

	// API route to update the challenge score in real-time
	router.POST("/update-score", func(c *gin.Context) {
		var data struct {
			ChallengeID int `json:"challenge_id"`
			NewScore    int `json:"new_score"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Emit score update to WebSocket clients
		server.BroadcastToNamespace("/", "update-score", data)

		c.JSON(http.StatusOK, gin.H{"message": "Score updated", "data": data})
	})

	// Start WebSocket server in a goroutine
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()

	return server
}
