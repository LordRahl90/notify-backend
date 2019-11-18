package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//User - User Object
type User struct {
	gorm.Model
	UserKey   string    `json:"user_key"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Token     string    `json:"token"`
	LastLogon time.Time `json:"last_logon"`
}

//BeforeCreate Add a default key to the user key stuff
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UserKey", uuid.New().String())
	scope.SetColumn("LastLogon", time.Now())
	return nil
}

//NewUser - Create a new user account
func (d *Database) NewUser(u *User) (*User, error) {
	db := d.DB
	_, cancel := context.WithCancel(d.Ctx)
	defer cancel()
	err := u.Validate()
	if err != nil {
		return nil, err
	}

	//check if record exists
	var oUser User
	err = db.First(&oUser, "email=?", u.Email).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if oUser.Email != "" {
		return nil, errors.New("User Exists already")
	}

	newPassword, err := generatePassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = newPassword

	err = db.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return u, err
}

//GetAllUsers - Return all the user information
func (d *Database) GetAllUsers() ([]*User, error) {
	db := d.DB
	var users []*User
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		for _, u := range users {
			u.Password = ""
		}
	}

	return users, err
}

//GetUser - Fetches One User's Record
func (d *Database) GetUser(id uint) (*User, error) {
	db := d.DB
	var user User
	// err := db.Where("id=?", id).First(&user).Error
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//GetUserByEmail - Retrieve the user information by email
func (d *Database) GetUserByEmail(email string) (*User, error) {
	fmt.Println(email)
	db := d.DB
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

//GetUserByKey - Find user by their userkey
func (d *Database) GetUserByKey(userKey int) (*User, error) {
	db := d.DB
	var user User
	err := db.Where("user_key=?", userKey).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//Authenticate a user by validating their email and password against what we have in the database
func (d *Database) Authenticate(email, password string) (*User, error) {
	db := d.DB
	if email == "" {
		return nil, errors.New("Please provide an email")
	}
	if password == "" {
		return nil, errors.New("Please provide a password")
	}

	user, err := d.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if chk := comparePassword(password, user.Password); !chk {
		return nil, errors.New("invalid password provided for this user")
	}

	token, err := generateToken(user.ID)
	if err != nil {
		return nil, err
	}
	user.Token = token
	user.LastLogon = time.Now()
	db.Save(&user)
	user.Password = ""

	return user, nil
}

//Validate the entry provided
func (u User) Validate() error {
	if u.Fullname == "" {
		return errors.New("Fullname should be provided")
	}

	if u.Email == "" {
		return errors.New("Email should be provided")
	}

	if u.Password == "" {
		return errors.New("Password should be provided")
	}

	if len(u.Password) < 6 {
		return errors.New("Password should be more than 6")
	}
	return nil
}

// Function to generate hashed strings from provided password.
func generatePassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed), err
}

//Function to compare password with the one provided in the database
func comparePassword(password, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}

//function to generate a JWT token, This will be returned with the user's information when they login
func generateToken(userID uint) (string, error) {
	expiryDate := time.Now().AddDate(0, 3, 0)
	hashSecret := "my hashing secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"expire_on": expiryDate, //fmt.Sprintf("%d-%d-%d 11:59:00", expiryDate.Year(), expiryDate.Month(), expiryDate.Day()),
	})
	tokenString, err := token.SignedString([]byte(hashSecret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

//DecodeToken -
func DecodeToken(tokenString string) (uint, error) {
	hashSecret := "my hashing secret"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hashSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiryDate := claims["expire_on"].(string)
		exp, err := time.Parse(time.RFC3339, expiryDate)
		if err != nil {
			return 0, err
		}
		if time.Now().After(exp) {
			return 0, errors.New("Expired Token")
		}

		return uint(claims["user_id"].(float64)), nil
	}
	return 0, errors.New("Invalid Token Provided")
}
