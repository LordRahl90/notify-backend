package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lordrahl90/notify-backend/app/middlewares"
	"github.com/lordrahl90/notify-backend/app/services/database"
)

//NewMessage request struct
type NewMessage struct {
	Sender   int    `json:"sender" form:"sender" binding:"required"`
	Reciever int    `json:"reciever" form:"reciever" binding:"required"`
	Content  string `json:"content" form:"content" binding:"required"`
	Media    string `json:"media" form:"media"`
}

//NewMessageHandler - function to initialize the message routes
func NewMessageHandler(router *gin.Engine) {
	m := router.Group("/messages")
	{
		m.Use(middlewares.Auth())
		m.POST("/", sendMessage)
		m.GET("/conversations", fetchConversations)
		m.GET("/conversations/:friendID", fetchIndividualConversation)
	}
}

func sendMessage(c *gin.Context) {
	var m NewMessage
	if err := c.ShouldBindJSON(&m); err != nil {
		returnResponse(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	msg := database.Message{
		SenderID:   m.Sender,
		RecieverID: m.Reciever,
		Content:    m.Content,
		Media:      m.Media,
	}

	if err := Database.NewMessage(&msg); err != nil {
		returnResponse(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	returnResponse(c, 200, true, "Message Sent successfully", msg)
}

func fetchIndividualConversation(c *gin.Context) {
	friend := c.Param("friendID")
	userID, err := getUserID(c)
	if err != nil {
		returnResponse(c, http.StatusUnauthorized, false, "Please login to proceed", nil)
		return
	}

	friendID, err := strconv.Atoi(friend)
	if err != nil {
		returnResponse(c, http.StatusBadRequest, false, "Invalid Friend ID provided", nil)
		return
	}

	msgs, err := Database.UserConversation(userID, uint(friendID))
	if err != nil {
		returnResponse(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	returnResponse(c, http.StatusOK, true, "Conversation loaded successfully!", msgs)
}

func fetchContacts(c *gin.Context) {
	returnResponse(c, 200, true, "Returning your contacts in bulk", nil)
}

func fetchConversations(c *gin.Context) {
	returnResponse(c, 200, true, "Sending your messages in good pack.", nil)
}
