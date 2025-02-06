package server

import (
	"html/template"
	"net/http"

	"github.com/ctf-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	handlers     *handlers.Handlers
	teamhandlers *handlers.TeamHandler
}

func NewServer(handlers *handlers.Handlers, teamhandlers *handlers.TeamHandler) *Server {
	router := gin.Default()
	server := &Server{router: router, handlers: handlers, teamhandlers: teamhandlers}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Add custom template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	s.router.SetFuncMap(funcMap)

	// Load templates and static files
	s.router.LoadHTMLGlob("../../../template/*.html")
	s.router.Static("/css", "../../../template/css")
	s.router.Static("/js", "../../../template/js")

	s.router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	s.router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	s.router.GET("/instructions", func(c *gin.Context) {
		c.HTML(http.StatusOK, "instructions.html", nil)
	})
	// s.router.GET("/hackerboard", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "hackerboard.html", nil)
	// })
	s.router.GET("/quests", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "quests.html", nil)
	})
	s.router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", nil)
	})
	// Define routes
	s.router.GET("/teams", s.handlers.GetTeams)
	s.router.GET("/challenges", s.handlers.GetChallenges)
	s.router.GET("/hackerboard", s.handlers.GetScores)
	s.router.GET("/ws", s.handlers.HandleWebSocket)

	// Team authentication routes
	s.router.POST("/register", s.teamhandlers.RegisterTeam)
	s.router.POST("/login", s.teamhandlers.LoginTeam)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
