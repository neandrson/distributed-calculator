package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"not null;index:idx_email;unique"`
	Password string `gorm:"not null"`
}
type Expression struct {
	gorm.Model
	UserId     uint   `gorm:"index:idx_userid"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	EX         string `gorm:"not null"`
	SubEx      string `gorm:"type:text"`
	Inwork     bool   `gorm:"default:false"`
	Ready      bool   `gorm:"default:false"`
	Timestamp  time.Time
	TimeUpdate time.Time
	Error      bool `gorm:"default:false"`
}
type StatusClient struct {
	ID            string
	ActiveWorkers int
	AllWorkers    int
}
