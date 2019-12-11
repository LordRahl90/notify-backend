package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/lordrahl90/notify-backend/app/middlewares"
	"github.com/lordrahl90/notify-backend/app/services/database"
)

//Database - this is used to access the entire struvs
var Database *database.Database

//User - User Request struct
type User struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Fullname string `json:"fullname" form:"fullname"`
	Password string `json:"password" form:"password" binding:"required"`
	Token    string `json:"token" form:"token"`
}

type SocialUser struct {
	Email    string `json:"email" form:"email" binding:"required"`
	FullName string `json:"full_name" form:"full_name" binding:"required"`
	Token    string `json:"token" form:"token" binding:"required"`
	Avatar   string `json:"avatar" form:"avatar" binding:"required"`
}

//FriendRequest - Struct to manage the friendrequest request :-)
type FriendRequest struct {
	FriendID uint `json:"friend_id" form:"friend_id" binding:"required"`
}

//FriendRequestApproval - request format for friend request approval
type FriendRequestApproval struct {
	RequestKey string `json:"request_key" form:"request_key" binding:"required"`
	Response   bool   `json:"response" form:"response" binding:"required"`
}

//DeviceRequest - format for registering a user's device
type DeviceRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
	Name  string `json:"name" form:"name"`
}

//NewUserHandler - Returns a new Route for user handlers
func NewUserHandler(router *gin.Engine) {

	u := router.Group("/users")
	{
		u.POST("/authenticate", authenticate)
		u.POST("/authenticate/social", authenticateSocial)
		u.POST("/", newUser)
		u.Use(middlewares.Logger())
		u.Use(middlewares.Auth())
		u.GET("/", allUsers)
		u.POST("/search", searchUsers)
		u.GET("/me", singleUser)
		u.GET("/me/friends", getFriends)
		u.POST("/me/friend/request", sendFriendRequest)
		u.GET("/me/friend/requests", getFriendRequests)
		u.PUT("/me/friend/update", updateRequest)
		u.POST("/me/device", registerDevice)
	}
}

type sReq struct {
	Data string `json:"data" form:"data" binding:"required"`
}

func searchUsers(c *gin.Context) {
	var req sReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}

	users, err := Database.SearchForUser(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}
	c.JSON(200, returnFormat(true, "Matches found.", users))
}

func registerDevice(c *gin.Context) {
	var req DeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, returnFormat(false, err.Error(), nil))
		return
	}

	dReq := database.UserDevice{
		UserID:     userID,
		Token:      req.Token,
		DeviceName: req.Name,
	}

	if err = Database.NewUserDevice(&dReq); err != nil {
		c.JSON(http.StatusUnauthorized, returnFormat(false, err.Error(), nil))
		return
	}

	c.JSON(200, returnFormat(true, "User Device updated successfully.", nil))
}

func authenticateSocial(c *gin.Context) {
	var req SocialUser
	if err := c.ShouldBindJSON(&req); err != nil {
		returnResponse(c, 400, false, err.Error(), nil)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.MinCost)
	if err != nil {
		returnResponse(c, 500, false, err.Error(), nil)
		return
	}

	user, err := Database.NewUser(&database.User{
		Fullname: req.FullName,
		Email:    req.Email,
		Password: string(password),
	})
	if err != nil {
		returnResponse(c, 500, false, err.Error(), nil)
		return
	}

	//save device token as well
	err = Database.NewUserDevice(&database.UserDevice{
		UserID: user.ID,
		Token:  req.Token,
	})
	if err != nil {
		returnResponse(c, 500, false, err.Error(), nil)
		return
	}

	newToken, err := database.GenerateToken(user.ID)
	if err != nil {
		returnResponse(c, 500, false, err.Error(), nil)
		return
	}
	friends, err := Database.GetUserFriends(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, returnFormat(false, err.Error(), nil))
		return
	}
	user.Token = newToken
	resp := map[string]interface{}{
		"user":    user,
		"friends": friends,
	}

	returnResponse(c, 200, true, "Authentication Completed successfully.", resp)
}

func authenticate(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}

	fmt.Println(req.Email, " ", req.Password)

	user, err := Database.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}

	friends, err := Database.GetUserFriends(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, returnFormat(false, err.Error(), nil))
		return
	}

	resp := map[string]interface{}{
		"user":    user,
		"friends": friends,
	}

	c.JSON(200, returnFormat(true, "User Created", resp))
}

func newUser(c *gin.Context) {
	var userReq User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}

	user, err := Database.NewUser(&database.User{
		Email:    userReq.Email,
		Fullname: userReq.Fullname,
		Password: userReq.Password,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}

	if userReq.Token != "" {
		//add a user device token to the mix.
		u := database.UserDevice{
			UserID: user.ID,
			Token:  userReq.Token,
		}
		if err := Database.NewUserDevice(&u); err != nil {
			c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
			return
		}
	}

	user.Password = ""
	c.JSON(200, returnFormat(true, "Authentication Successful", user))
}

//Return all registered users
func allUsers(c *gin.Context) {
	users, err := Database.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, returnFormat(false, err.Error(), nil))
		return
	}
	c.JSON(200, returnFormat(true, "All users retrieved", users))
}

func singleUser(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, returnFormat(false, "Invalid User Record", nil))
	}
	user, err := Database.GetUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, returnFormat(false, err.Error(), nil))
		return
	}
	user.Password = ""

	c.JSON(200, returnFormat(true, "User profile retrieved", user))
}

func getFriends(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		returnResponse(c, http.StatusUnauthorized, false, "login to proceed.", nil)
		return
	}

	friends, err := Database.GetUserFriends(userID)
	if err != nil {
		returnResponse(c, http.StatusUnauthorized, false, err.Error(), nil)
		return
	}

	returnResponse(c, 200, false, "Hello World wide", friends)
}

func getFriendRequests(c *gin.Context) {
	userID, ok := c.Get("user_id")
	var reqs []database.Friend
	if !ok {
		c.JSON(http.StatusUnauthorized, returnFormat(false, "Invalid User Record", nil))
		return
	}

	reqType := c.Query("type")

	if reqType == "sent" {
		rr, err := Database.GetSentFriendRequest(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, returnFormat(false, "Invalid User Record", nil))
			return
		}
		reqs = rr

	} else {
		rr, err := Database.GetRecievedFriendRequest(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, returnFormat(false, "Invalid User Record", nil))
			return
		}
		reqs = rr
	}

	returnResponse(c, 200, true, "Friend Requests loaded successfully.", reqs)

}

func sendFriendRequest(c *gin.Context) {
	var req FriendRequest
	userID, err := getUserID(c)
	if err != nil {
		returnResponse(c, http.StatusUnauthorized, false, "Please login to proceed", nil)
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		returnResponse(c, http.StatusBadRequest, false, "provide a valid friend request...", nil)
		return
	}

	friendRequest := database.Friend{
		UserID:   userID,
		FriendID: req.FriendID,
		Status:   false,
	}

	err = Database.NewFriendRequest(friendRequest)
	if err != nil {
		returnResponse(c, http.StatusBadRequest, false, err.Error(), err)
		return
	}

	returnResponse(c, http.StatusOK, true, "Friend request sent successfully.", nil)
}

func updateRequest(c *gin.Context) {
	var f FriendRequestApproval
	_, err := getUserID(c)
	if err != nil {
		returnResponse(c, http.StatusUnauthorized, false, "login to proceed.", nil)
		return
	}

	if err = c.ShouldBindJSON(&f); err != nil {
		returnResponse(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	if err = Database.UpdateFriendRequest(f.RequestKey, f.Response); err != nil {
		returnResponse(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	returnResponse(c, 200, true, "Friend Request updated successfully.", nil)
}

func getUserID(c *gin.Context) (uint, error) {
	userID, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("unauthorised user account")
	}
	return userID.(uint), nil
}

func returnFormat(success bool, message string, data interface{}) gin.H {
	response := map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}
	return gin.H(response)
}

func returnResponse(c *gin.Context, code int, success bool, message string, data interface{}) {
	c.JSON(code, returnFormat(success, message, data))
}
