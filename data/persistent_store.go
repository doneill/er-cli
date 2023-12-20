package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// ----------------------------------------------
// DB connection
// ----------------------------------------------

func OpenDB(file string) {
	database, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic("Failed to open database")
	}

	db = database
}

// ----------------------------------------------
// table data
// ----------------------------------------------

func GetTables() (tableList []string, err error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func GetTableRowCount(name string) int64 {
	var count int64
	db.Table(name).Count(&count)
	return count
}

// ----------------------------------------------
// queries
// ----------------------------------------------

func SelectUser() Accounts_User {
	var user Accounts_User
	db.First(&user)
	return user
}

func SelectUserProfileById(id int) User_Profile {
	var profile User_Profile
	db.Where("id = ?", id).Find(&profile)
	return profile
}

func SelectPendingSyncEvents() []Event {
	var events []Event
	db.Where("is_draft = 0").Where(db.Where("remote_id IS NULL").Or("remote_id = ?", "")).Find(&events)
	return events
}
