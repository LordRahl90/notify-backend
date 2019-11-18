package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/lordrahl90/notify-backend/app/handlers"
	"github.com/lordrahl90/notify-backend/app/services"
	"github.com/lordrahl90/notify-backend/app/services/database"
	"google.golang.org/api/option"
)

func main() {

	err := godotenv.Load("./.envs/.app.env")
	if err != nil {
		log.Fatal(err)
	}

	dbase := os.Getenv("DATABASE")
	user := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	dsn := "root:@(localhost)/notifier?charset=utf8&parseTime=True&loc=Local"

	env := flag.String("env", "local", "The Current development environment")
	flag.Parse()

	if *env != "local" {
		dsn = fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbase)
	}

	s := services.NewServer()
	db, err := setupDB("mysql", dsn)
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
