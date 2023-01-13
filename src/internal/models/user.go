package models

import "gorm.io/gorm"

// @Description User account information
type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email" example:"username@domain.tld" validate:"required"`
	Password string `json:"password" example:"qwerty" validate:"required"`
} // @name User
