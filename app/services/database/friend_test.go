package database

import (
	"log"
	"testing"

	"github.com/bxcodec/faker/v3"
)

func createUsers() {
	for i := 1; i <= 10; i++ {
		u := User{
			Fullname: faker.FirstName() + " " + faker.LastName(),
			Email:    faker.Email(),
			Password: "secret",
		}

		_, err := db.NewUser(&u)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func tearDown() {
	db.DB.Raw("TRUNCATE TABLE users")
}

func TestNewFriend(t *testing.T) {
	// createUsers()
	// defer tearDown()

	friends := []Friend{
		Friend{
			UserID:   2,
			FriendID: 3,
		},
		Friend{
			UserID:   2,
			FriendID: 4,
		},
		Friend{
			UserID:   2,
			FriendID: 5,
		},
	}

	for _, f := range friends {
		err := db.NewFriendRequest(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestLoadSentFriendRequest(t *testing.T) {
	userID := uint(2)
	req, err := db.GetSentFriendRequest(userID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(req)
}

func TestApproveFriendship(t *testing.T) {
	// friendID := 2
}
