package services

import (
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/lordrahl90/notify-backend/app/middlewares"
	"github.com/lordrahl90/notify-backend/app/services/database"
)

//Server - Struct for the server objects.
type Server struct {
	DB      *database.Database
	FireApp *firebase.App
	Router  *gin.Engine
}

//NewServer - Returns a New Server
func NewServer() *Server {
	r := gin.Default()
	s := &Server{
		Router: r,
	}
	s.setRoutes()
	return s
}

//Start - Starts the server
func (s *Server) Start(address string) {
	s.Router.Run(address)
}

func (s *Server) setRoutes() {
	s.Router.RouterGroup.Use(middlewares.BasicMonitor())
	s.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
}
