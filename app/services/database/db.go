package database

import (
	"context"

	"github.com/jinzhu/gorm"
)

//Database struct
type Database struct {
	DB  *gorm.DB
	Ctx context.Context
}

//NewDatabase - returns a new instance of database connection after migrating the necessary models.
func NewDatabase(dialect, connString string) (*Database, error) {
	db, err := gorm.Open(dialect, connString)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{}, &UserDevice{}, &UserInfo{}, &Message{}, &Friend{})

	dBase := &Database{
		DB: db,
	}
	return dBase, nil
}
