package model

import "github.com/jinzhu/gorm"

func AutoMigrate(db *gorm.DB) (err error) {
	err = db.AutoMigrate(&User{}).Error
	return
}
