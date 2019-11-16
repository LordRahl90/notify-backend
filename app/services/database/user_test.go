package database

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *Database

func init() {
	dbase, err := NewDatabase("mysql", "root:@/notifier_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	db = dbase
}

func TestNewUser(t *testing.T) {
	user := User{
		Email:     "tolaabbey009@gmail.com",
		Password:  "secret",
		Fullname:  "Alugbin LordRahl",
		LastLogon: time.Now(),
	}

	n, err := db.NewUser(user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New User's ID is: ", n.ID)
}

func TestGetAllUsers(t *testing.T) {
	users, err := db.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", users)
}

func TestGetUser(t *testing.T) {
	id := uint(2)
	user, err := db.GetUser(id)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(user.Email)
}

func TestGetUserByEmail(t *testing.T) {
	email := "tolaabbey009@gmail.com"
	user, err := db.GetUserByEmail(email)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user.ID)
}

func TestAuthentication(t *testing.T) {
	email, password := "tolaabbey009@gmail.com", "secret"
	user, err := db.Authenticate(email, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user.Token)
}

func TestGenerateToken(t *testing.T) {
	token, err := generateToken(100)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)
}

func TestDecodeToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVfb24iOiIyMDIwLTAyLTE1VDAxOjA4OjM5LjI1MDk3NyswMTowMCIsInVzZXJfaWQiOjV9.aDCtPwOONHUFyicVdf-5JrdRqsFKa5Y4uoS97JyQ3TM"
	res, err := DecodeToken(tokenString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
