package models

import "time"

type Student struct {
	Name        string
	Email       string `gorm:"primaryKey;size:255"`
	Age         uint8
	IsSuspended bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
