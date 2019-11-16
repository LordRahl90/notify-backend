package firebase

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

//FireApp - App That contains the firebase App
type FireApp struct {
	App *firebase.App
}

//New Create New Firebase App
func New() (*firebase.App, error) {
	return nil, nil
}

//SendMessage - Sends Message from firebase to devices
func (app FireApp) SendMessage(token string) (string, error) {
	ctx := context.Background()
	client, err := app.App.Messaging(ctx)
	if err != nil {
		log.Fatal("Error retrieving messaging client")
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "App Title",
			Body:  "Yup, Whatsup???",
		},
		Data: map[string]string{
			"fullname": "Adeleke Raheem",
			"email":    "I am d louvre",
		},
		Token: token,
	}

	resp, err := client.Send(ctx, message)
	if err != nil {
		return "", err
	}

	fmt.Println("Message Sent successfully!!!", resp)
	return resp, nil
}