package routes

import (
	"fiber-api/utils/types"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, config types.Config) {
	// unauthenticated routes
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", Login)
	auth.Post("/register", RegisterTeacher)

	// // authenticated routes
	// api.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(config.App.AccessTokenSecret),
	// }))
	api.Post("/register", TeacherAssignment)
	api.Get("/commonstudents", CommonStudents)
	api.Post("/suspend", SuspendStudent)
	api.Post("/retrievefornotifications", StudentNotificaton)
}
