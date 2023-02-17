package database

import (
	"fiber-api/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database! \n", err.Error())
	}

	log.Println("Connected to database!")
	db.Logger = logger.Default.LogMode(logger.Info)

	db.AutoMigrate(&models.Student{}, &models.Teacher{}, &models.TeacherStudent{})

	Database = DbInstance{Db: db}
}
