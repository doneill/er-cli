package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ----------------------------------------------
// exported funtions
// ----------------------------------------------

func DbConnect(file string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
