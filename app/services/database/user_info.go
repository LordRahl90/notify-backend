package database

import "github.com/jinzhu/gorm"

//UserInfo - struct to hold a user information
type UserInfo struct {
	gorm.Model
	User      User `json:"user_id" gorm:"foreignkey:UserRefer"`
	UserRefer uint
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
	Phone     string `json:"phone"`
}
