package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/lordrahl90/notify-backend/app/handlers"
	"github.com/lordrahl90/notify-backend/app/services"
	"github.com/lordrahl90/notify-backend/app/services/database"
	"google.golang.org/api/option"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	s := services.NewServer()
	db, err := setupDB("mysql", "root:@/notifier?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("An error occured while setting up the database")
	}
	s.DB = db
	setupEndpoints(s.Router)
	handlers.Database = s.DB
	fmt.Println("Starting server")
	s.Start("0.0.0.0:5500")
}

func createFirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile("./fire-messaging.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("Error Initializing app %v", err)
	}

	return app, nil
}

func setupDB(dialect, conString string) (*database.Database, error) {
	dBase, err := database.NewDatabase(dialect, conString)
	if err != nil {
		return nil, err
	}

	return dBase, nil
}

func setupEndpoints(router *gin.Engine) {
	handlers.NewUserHandler(router)
	handlers.NewMessageHandler(router)
}
