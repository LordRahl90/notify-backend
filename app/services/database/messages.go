package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Message struct for holding messages
type Message struct {
	gorm.Model
	SenderID   int       `json:"sender_id"`
	Sender     User      `json:"sender" gorm:"foreignkey:SenderID"`
	RecieverID int       `json:"reciever_id"`
	Reciever   User      `json:"reciever" gorm:"foreignkey:RecieverID"`
	Content    string    `json:"content"`
	Media      string    `json:"media"`
	Read       bool      `json:"read"`
	DateSent   time.Time `json:"date_sent"`
	DateRead   time.Time `json:"date_read"`
}
