package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func dbConnect(connectString string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(connectString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, db.AutoMigrate(&DeviceStatus{})
}

func insert(db *gorm.DB, ds DeviceStatus) error {
	var result = db.Create(ds)
	return result.Error
}

func getLatest(db *gorm.DB) (DeviceStatus, error) {
	var ds = new(DeviceStatus)
	var result = db.First(ds)
	return *ds, result.Error
}
