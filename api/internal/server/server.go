package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ctf-api/internal/handlers"
	"github.com/ctf-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router            *gin.Engine
	handlers          *handlers.Handlers
	teamhandlers      *handlers.TeamHandler
	challengehandlers *handlers.ChallengeHandler
}

func NewServer(handlers *handlers.Handlers, teamhandlers *handlers.TeamHandler, challengehandlers *handlers.ChallengeHandler) *Server {
	router := gin.Default()
	server := &Server{
		router:            router,
		handlers:          handlers,
		teamhandlers:      teamhandlers,
		challengehandlers: challengehandlers,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

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
	s.router.Static("/admin/css","../../../template/css")

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
	teamRoutes.Use(middleware.TeamAuth) // Middleware for authentication
	{
		teamRoutes.POST("/verifyflag", s.challengehandlers.SubmitFlag)
		log.Println("server called") // List Challenges
	}

	// 📌 **Admin Routes (Requires Admin Role)**
	adminRoutes := s.router.Group("/admin")
	adminRoutes.Use(middleware.AdminAuth) // Middleware for admin authentication
	{
		adminRoutes.GET("/", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "admin.html", nil) })
		adminRoutes.POST("/addchallenge", s.challengehandlers.CreateChallenge)
		adminRoutes.GET("/challenges", s.challengehandlers.GetChallenges)
		adminRoutes.PUT("/updatechallenge/:id", s.challengehandlers.UpdateChallenge)
		adminRoutes.DELETE("/deletechallenge/:id", s.challengehandlers.DeleteChallenge)
		adminRoutes.GET("/teams", s.handlers.GetTeams)
	}
	SetupSocketServer(s.router)
	

	// 📌 WebSocket (Real-time updates)
	s.router.GET("/ws", s.handlers.HandleWebSocket)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
