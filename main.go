package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/lordrahl90/notify-backend/app/handlers"
	"github.com/lordrahl90/notify-backend/app/services"
	"github.com/lordrahl90/notify-backend/app/services/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/api/option"
)

func main() {
	sigs := make(chan os.Signal, 1)
	ctx := context.Background()
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
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go startMonitoringServer(ctx)
	fmt.Println("Starting server")
	go s.Start(ctx, "0.0.0.0:5500")

	select {
	case sig := <-sigs:
		fmt.Println("Done Completely", sig)
		// gracefulShutdown(s)

	}

	gracefulShutdown(s)
	fmt.Println("Hello World")

}

func gracefulShutdown(s *services.Server) {
	s.DB.DB.Close()
	fmt.Println("Shutting down the DB gracefully...")
	os.Exit(0)
}

func startMonitoringServer(ctx context.Context) {
	http.Handle("/metrics", promhttp.Handler())
	println("Monitoring server added successfully.")
	if err := http.ListenAndServe("0.0.0.0:5501", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server stopped")
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
