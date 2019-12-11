package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)


 var app *firebase.App

func init(){
	opt := option.WithCredentialsFile("../../../fire-messaging.json")
	fireApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("error Initializing app ", err)
	}

	app=fireApp
}

//func TestFireApp_SendMessage(t *testing.T) {
//	if app==nil{
//		log.Fatal("App is nkt initialized")
//	}
//	token:=`cjBSwhwGiQU:APA91bG5aBPcskrSZrv6qugkOOVc2PY7r2auN8MOgaVLjlP7blgIZBZE3vO9NTjwttLveOhHsNVp6oiY9RxcnhekgYSCoypkBRnu_oAUKyJfZ8r4t_fHc998XxS4X3fQyr-jKZuXJF5i`
//	fireApp,err:=New(app)
//	if err!=nil{
//		log.Fatal(err)
//	}
//
//	str,err:=fireApp.SendMessage(token)
//	if err!=nil{
//		log.Fatal(err)
//	}
//
//	log.Println(&str)
//}