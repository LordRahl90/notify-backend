package database

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

//Friend struct
type Friend struct {
	gorm.Model
	RequestKey string `json:"request_key"`
	UserID     uint   `json:"user" json:"-"`
	FriendID   uint   `json:"friend_id" json:"-"`
	Friend     User   `gorm:"foreignkey:FriendID"`
	User       User   `gorm:"foreignkey:UserID" json:"-"`
	Status     bool   `json:"status" gorm:"default:false"`
	Response   bool   `json:"res_status" gorm:"default:false"`
}

//BeforeCreate - Function to handle the beforecreate hook
func (f *Friend) BeforeCreate(scope *gorm.Scope) {
	scope.SetColumn("request_key", uuid.New().String())
}

//NewFriendRequest - Create a new friend record
func (d *Database) NewFriendRequest(f Friend) error {
	if f.FriendID == f.UserID {
		return errors.New("you cannot send a friend request to yourself")
	}
	db := d.DB
	if _, err := d.GetUser(f.UserID); err != nil {
		return err
	}

	err := db.Where(Friend{UserID: f.UserID, FriendID: f.FriendID}).FirstOrCreate(&f).Error
	if err != nil {
		return err
	}

	return err
}

//GetRecievedFriendRequest - Retrieve all the friend request of a user.
func (d *Database) GetRecievedFriendRequest(userID uint) ([]Friend, error) {
	var request []Friend
	db := d.DB
	err := db.Preload("User").Preload("Friend").Where("friend_id=? and response=false", userID).Find(&request).Error
	if err != nil {
		return nil, err
	}

	return request, nil
}

//GetSentFriendRequest - Retrieve all the friend request of a user.
func (d *Database) GetSentFriendRequest(userID uint) ([]Friend, error) {
	var request []Friend
	db := d.DB
	err := db.Preload("Friend").Where("user_id=? and response=false", userID).Find(&request).Error
	if err != nil {
		return nil, err
	}

	return request, nil
}

//UpdateFriendRequest -
func (d *Database) UpdateFriendRequest(requestKey string, response bool) error {
	db := d.DB
	var f Friend
	err := db.Where("request_key=?", requestKey).First(&f).Error
	if err != nil {
		return err
	}

	if response {
		err := db.Model(&f).Update("response", response).Error
		if err != nil {
			return err
		}

		//create a new friend for the new friend
		var nf = Friend{
			UserID:   f.FriendID,
			FriendID: f.UserID,
			Status:   true,
			Response: true,
		}
		err = d.NewFriendRequest(nf)
		if err != nil {
			return err
		}

	}
	err = db.Model(&f).Update("status", true).Error
	if err != nil {
		return err
	}

	return nil
}

//GetUserFriends - Retrieves the friends list of a user.
func (d *Database) GetUserFriends(userID uint) ([]Friend, error) {
	db := d.DB
	var friends []Friend
	if err := db.Preload("Friend").Where("user_id=? and response=?", userID, true).Find(&friends).Error; err != nil {
		return nil, err
	}

	return friends, nil
}
