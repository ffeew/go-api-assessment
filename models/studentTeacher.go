package models

import "time"

type TeacherStudent struct {
	StudentEmail string  `gorm:"primaryKey"`
	Student      Student `gorm:"foreignKey:StudentEmail"`
	TeacherEmail string  `gorm:"primaryKey"`
	Teacher      Teacher `gorm:"foreignKey:TeacherEmail"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
