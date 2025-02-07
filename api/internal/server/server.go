package server

import (
	"html/template"
	"net/http"

	"github.com/ctf-api/internal/handlers"
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
	s.router.Static("/images", "../../../template/images")
	s.router.Static("/fonts", "../../../template/fonts")

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
	s.router.GET("/admin",func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK,"admin.html",nil)
	})
	s.router.GET("/teams", s.handlers.GetTeams)
	s.router.GET("/challenges", s.handlers.GetChallenges)
	s.router.GET("/hackerboard", s.handlers.GetScores)
	s.router.GET("/ws", s.handlers.HandleWebSocket)

	// Team authentication routes
	s.router.POST("/register", s.teamhandlers.RegisterTeam)
	s.router.POST("/login", s.teamhandlers.LoginTeam)
	s.router.POST("/addchallenge", s.challengehandlers.CreateChallenge)
	s.router.GET("/viewchallenge",s.challengehandlers.GetChallenges)
	s.router.PUT("/updatechallenge",s.challengehandlers.UpdateChallenge)
	s.router.DELETE("/deletechallenge/:id", s.challengehandlers.DeleteChallenge)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
