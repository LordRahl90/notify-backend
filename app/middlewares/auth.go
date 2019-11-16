package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordrahl90/notify-backend/app/services/database"
)

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

		auth := c.Request.Header["Authorization"]
		if len(auth) <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization token",
			})
			c.Abort()
			return
		}

		authToken := strings.Split(auth[0], " ")
		if len(authToken) <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid token format provided",
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
