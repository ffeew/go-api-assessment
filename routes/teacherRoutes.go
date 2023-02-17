package routes

import (
	"fiber-api/database"
	"fiber-api/models"
	"fiber-api/utils/types/requests"
	"fiber-api/utils/types/responses"
	"fiber-api/utils/validators"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TeacherAssignment(c *fiber.Ctx) error {
	body := new(requests.StudentAssignmentRequest)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please provide the teacher's email and the students' emails in this format: {teacher: <email>, students: [<email>, <email>, ...]}",
		})
	}

	errors := validators.IsValidParams(body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: &errors})
	}

	// use a transaction to guarantee the consistency of the data
	err := database.Database.Db.Transaction(func(tx *gorm.DB) error {
		// insert into the StudentTeacher table the teacher's email and the students' emails
		for _, student := range body.StudentEmail {
			studentTeacher := models.TeacherStudent{StudentEmail: student, TeacherEmail: body.TeacherEmail}
			// insert the data if it doesn't exist
			if err := tx.Where(&studentTeacher).FirstOrCreate(&studentTeacher).Error; err != nil {
				return err
			}
		}
		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to complete the student assignment",
		})
	} else {
		return c.Status(204).JSON(fiber.Map{
			"message": "successfully assigned students to the teacher",
		})
	}
}

func CommonStudents(c *fiber.Ctx) error {
	q := string(c.Request().URI().QueryString())
	if q == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide >1 teachers' emails as a query string",
		})
	}
	qSlice := strings.Split(q, "&")
	for i, v := range qSlice {
		qSlice[i] = strings.TrimPrefix(v, "teacher=")
		qSlice[i] = strings.Replace(qSlice[i], "%40", "@", 1)
	}

	// split the query string into an array of strings
	query := requests.CommonStudentsRequest{TeacherEmail: qSlice}

	// validate the query parameters
	errors := validators.IsValidParams(query)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: &errors})
	}

	// retrieve the students that are registered to all the teachers
	students := []models.Student{}
	if err := database.Database.Db.Table("students").Joins("INNER JOIN teacher_students ON students.email = teacher_students.student_email").Where("teacher_students.teacher_email IN ?", query.TeacherEmail).Group("students.email").Having("COUNT(DISTINCT teacher_students.teacher_email) = ?", len(query.TeacherEmail)).Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve the common students",
		})
	}
	// return the students' emails
	studentEmails := []string{}
	for _, student := range students {
		studentEmails = append(studentEmails, student.Email)
	}
	return c.Status(fiber.StatusOK).JSON(responses.CommonStudentsResponse{Students: studentEmails})
}

func SuspendStudent(c *fiber.Ctx) error {
	body := new(requests.SuspendStudentRequest)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide the student's email in this format: {student: <email>}",
		})
	}

	errors := validators.IsValidParams(body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: &errors})
	}

	// check if the student exists
	student := models.Student{}
	if err := database.Database.Db.Where("email = ?", body.StudentEmail).First(&student).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: "No student found"})
	}

	// update the student's status to suspended
	if err := database.Database.Db.Model(&student).Update("is_suspended", true).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: "Unable to suspend student"})
	}
	return c.Status(204).JSON(fiber.Map{
		"message": "successfully suspended the student",
	})
}

func StudentNotificaton(c *fiber.Ctx) error {
	body := new(requests.StudentNotificationRequest)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide the teacher's email and the notification in this format: {teacher: <email>, notification: <string>}",
		})
	}

	errors := validators.IsValidParams(body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: &errors})
	}

	// retrieve the students that are registered to the teacher
	students := []models.Student{}
	if err := database.Database.Db.Table("students").Joins("INNER JOIN teacher_students ON students.email = teacher_students.student_email").Where("teacher_students.teacher_email = ?", body.TeacherEmail).Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve the students",
		})
	}

	// retrieve the students that are mentioned in the notification
	re := regexp.MustCompile(`\b[\w\.-]+@[\w\.-]+\.\w{2,}\b`)
	mentionedStudents := re.FindAllString(body.Notification, -1)

	// retrieve the students that are suspended
	suspendedStudents := []models.Student{}
	if err := database.Database.Db.Where("is_suspended = ?", true).Find(&suspendedStudents).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve the suspended students",
		})
	}

	// generate the list of students that are registered to the teacher or has been mentioned in the notification
	temp := []string{}
	for _, student := range students {
		temp = append(temp, student.Email)
	}
	for _, student := range mentionedStudents {
		temp = append(temp, student)
	}

	occurred := map[string]bool{}
	notificationStudents := []string{}
	for _, email := range temp {
		if !occurred[email] {
			occurred[email] = true
			notificationStudents = append(notificationStudents, email)
		}
	}

	// find the students that are fufil the criteria
	studentEmails := []string{}
	for _, student := range notificationStudents {
		suspended := false
		for _, suspendedStudent := range suspendedStudents {
			if student == suspendedStudent.Email {
				suspended = true
				break
			}
		}
		if !suspended {
			studentEmails = append(studentEmails, student)
		}
	}
	return c.Status(fiber.StatusOK).JSON(responses.StudentNotificationResponse{Recipients: studentEmails})
}
