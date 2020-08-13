package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique_index"`
	Password string `json:"-"`
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	if db.First(&user, id).RecordNotFound() {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func UserSignUp(db *gorm.DB, name, password string) (err error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			Name:     name,
			Password: fmt.Sprintf("%s", pass),
		}).Error

		if pgError, ok := err.(*pq.Error); ok {
			if pgError.Code == "23505" {
				return ErrUserConflict
			}
		}

		return err
	})

}

func UserLogin(db *gorm.DB, name, password string) (*User, error) {
	var user User
	if db.First(&user, "name=?", name).RecordNotFound() {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
