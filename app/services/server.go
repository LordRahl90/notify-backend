package services

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lordrahl90/notify-backend/app/services/database"
	"log"
	"net/http"
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
	s.SetupWebsocketServer()
	s.setRoutes()
	return s
}

//Start - Starts the server
func (s *Server) Start(ctx context.Context, address string) {
	s.Router.Run(address)
}

func (s *Server) setRoutes() {
	s.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
}


func (s *Server) SetupWebsocketServer(){
	r:=s.Router
	//r.Use(middlewares.Auth())
	r.GET("/ws",func(c *gin.Context){
		fmt.Println("Socket connection started")
		socketHandler(c.Writer,c.Request)
	})
}

var wsupgrader=websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize:1024,
}

func socketHandler(w http.ResponseWriter,r *http.Request){

	wsupgrader.CheckOrigin=func(r *http.Request) bool {
		origin:=r.Host
		fmt.Println("Origin is: ",origin)
		if origin!=r.Host{
			return false
		}
		return true
	}
	conn,err:=wsupgrader.Upgrade(w,r,w.Header())
	if err!=nil{
		fmt.Println("Failed to upgrade the network")
		log.Fatal(err)
	}

	if err=conn.WriteMessage(1,[]byte("Connection Duly Established")); err!=nil{
		log.Println(err)
	}

	for{
		var t, msg, err = conn.ReadMessage()
		if err!=nil{
			break
		}

		fmt.Println(string(msg))

		if err=conn.WriteMessage(t,[]byte("Response: "+string(msg))); err!=nil{
			log.Println(err)
		}
	}

}