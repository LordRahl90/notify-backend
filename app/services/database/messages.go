package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

//Message struct for holding messages
type Message struct {
	gorm.Model
	MessageKey string    `json:"message_key"`
	SenderID   int       `json:"sender_id"`
	Sender     User      `json:"sender" gorm:"foreignkey:SenderID"`
	RecieverID int       `json:"reciever_id"`
	Reciever   User      `json:"reciever" gorm:"foreignkey:RecieverID"`
	Content    string    `json:"content"`
	Media      string    `json:"media"`
	Read       bool      `json:"read" gorm:"default=false"`
	DateSent   time.Time `json:"date_sent"`
	DateRead   time.Time `json:"date_read" gorm:"default=null"`
}

//BeforeCreate - Add the message key by default.
func (m *Message) BeforeCreate(scope *gorm.Scope) {
	scope.SetColumn("message_key", uuid.New().String())
	scope.SetColumn("date_sent", time.Now())
	scope.SetColumn("date_read", time.Now())
}

//NewMessage - create a new message on the platform
func (d *Database) NewMessage(m *Message) error {
	db := d.DB
	err := m.Validate()
	if err != nil {
		return err
	}
	err = db.Create(&m).Error
	return err
}

//UserConversation - Return all the messages that belongs to a user.
func (d *Database) UserConversation(sender, recipient uint) ([]Message, error) {
	db := d.DB
	var messages []Message
	if err := db.Preload("Sender").Preload("Reciever").
		Where("(sender_id in (?,?) and reciever_id in(?,?))", sender, recipient, sender, recipient).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

//Validate - make sure that the message is validated
func (m *Message) Validate() error {
	if m.SenderID <= 0 {
		return errors.New("Invalid Sender ID")
	}

	if m.RecieverID <= 0 {
		return errors.New("Invalid Reciever ID")
	}

	if m.Content == "" {
		return errors.New("Empty messages cannot be processed")
	}

	return nil
}
