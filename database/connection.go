package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBConnect() (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open("baleCorp.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
