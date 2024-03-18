package models

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	AccountID uint
	Token     string
}
