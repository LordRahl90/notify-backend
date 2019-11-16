package database

import (
	"errors"

	"github.com/jinzhu/gorm"
)

//UserDevice struct for managing user device and token
type UserDevice struct {
	gorm.Model
	UserID     uint
	User       User   `json:"user_id" gorm:"foreignkey:UserID"`
	DeviceName string `json:"device_name"`
	Token      string `json:"token"`
}

//NewUserDevice - Register a new user device onto the platform
//we need to check if the user's device has been registered before.
//This will inform us on how to proceed.
func (d *Database) NewUserDevice(u *UserDevice) error {
	db := d.DB
	if err := u.Validate(); err != nil {
		return err
	}

	err := db.FirstOrCreate(&u, &u).Error
	return err
}

//GetUserDevice - This function returns the user's device to us.
func (d *Database) GetUserDevice(userID uint) (*UserDevice, error) {
	db := d.DB
	var u *UserDevice
	err := db.Where("user_id=?", userID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

//Validate - Validate the entries provided from the controller.
func (u *UserDevice) Validate() error {
	if u.UserID <= 0 {
		return errors.New("Invalid User ID Provided")
	}
	if u.Token == "" {
		return errors.New("Invalid Device token provided")
	}

	return nil
}
