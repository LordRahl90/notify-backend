package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordrahl90/notify-backend/app/services/database"
	"github.com/lordrahl90/notify-backend/app/services/prometheus"
)

//BasicMonitor - Logs some basic Informaton about the request and response time
func BasicMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		host := c.Request.Host
		prometheus.IncrementRequestCount(host, path)

		//we might want to log the response time on here too....
		c.Next()
	}
}

//Logger - Test Middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		log.Print(latency)
		status := c.Writer.Status()
		log.Print(status)
	}
}

//Auth Middleware for authenticated routes
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		prometheus.IncrementResponseCount(c.Writer.Status(), c.Writer.Size())
		var auth []string

		fmt.Printf("%v\n", c.Request.Header)

		if len(c.Request.Header["Sec-Websocket-Protocol"]) > 0 && c.Request.Header["Sec-Websocket-Protocol"][0] != "" {
			authHeader := c.Request.Header["Sec-Websocket-Protocol"]
			auth = []string{"Bearer " + authHeader[0]}
		} else {
			auth = c.Request.Header["Authorization"]
		}

		if len(auth) <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "No Authorization token",
			})
			c.Abort()
			return
		}

		authToken := strings.Split(auth[0], " ")
		if len(authToken) <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid token format provided",
			})
			c.Abort()
			return
		}
		token := authToken[1]
		userID, err := database.DecodeToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
		fmt.Println("User ID is: ", userID)
		c.Set("user_id", userID)
		c.Next()
	}
}
