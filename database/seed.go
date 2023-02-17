package database

import (
	"fiber-api/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SeedDb(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database! \n", err.Error())
	}

	log.Println("Connected to database!")
	db.Logger = logger.Default.LogMode(logger.Info)

	db.AutoMigrate(&models.Student{}, &models.Teacher{}, &models.TeacherStudent{})

	// create a few default students for testing purposes
	students := []models.Student{
		{Email: "student1@gmail.com", Name: "Student 1", Age: 10},
		{Email: "student2@gmail.com", Name: "Student 2", Age: 10},
		{Email: "student3@gmail.com", Name: "Student 3", Age: 10},
		{Email: "student4@gmail.com", Name: "Student 4", Age: 10},
	}

	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("Password1!"), 8)
	hashedPassword := string(hashedBytes)

	teachers := []models.Teacher{
		{Email: "teacher1@gmail.com", Name: "Teacher 1", Age: 20, Password: hashedPassword},
		{Email: "teacher2@gmail.com", Name: "Teacher 2", Age: 20, Password: hashedPassword},
	}
	// use a transaction to guarantee the consistency of the data
	db.Transaction(func(tx *gorm.DB) error {
		// insert the students into the db
		for _, student := range students {
			if err := tx.Create(&student).Error; err != nil {
				return err
			}
		}

		// insert the teachers into the db
		for _, teacher := range teachers {
			if err := tx.Create(&teacher).Error; err != nil {
				return err
			}
		}
		// return nil will commit the whole transaction
		return nil
	})

}
