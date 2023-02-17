package models

import "time"

type Teacher struct {
	Name      string
	Email     string `gorm:"primaryKey;size:255"`
	Password  string
	Age       uint8
	CreatedAt time.Time
	UpdatedAt time.Time
}
