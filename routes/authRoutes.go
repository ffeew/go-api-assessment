package routes

import (
	"fiber-api/database"
	"fiber-api/models"
	"fiber-api/utils/types"
	"fiber-api/utils/types/requests"
	"fiber-api/utils/types/responses"
	"fiber-api/utils/validators"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golobby/dotenv"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {

	credentials := new(requests.LoginRequest)
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please provide the email and password as strings",
		})
	}

	errors := validators.IsValidParams(credentials)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: &errors})
	}

	// retrieve the user from the database
	teacher := models.Teacher{}
	if err := database.Database.Db.Where("email = ?", credentials.Email).First(&teacher).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: "No teacher found"})
	}

	// check if password matches
	err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(credentials.Password))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: "Wrong credentials provided"})
	}

	serverConfig := types.Config{}
	file, err := os.Open(".env")
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Message: "Something went wrong"})
		log.Fatal(err)
	}
	err = dotenv.NewDecoder(file).Decode(&serverConfig)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Message: "Something went wrong"})
		log.Fatal(err)
	}

	// Create the Claims
	accessTokenClaims := jwt.MapClaims{
		"email": teacher.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	refreshTokenClaims := jwt.MapClaims{
		"email": teacher.Email,
		"exp":   time.Now().Add(time.Hour * 168).Unix(),
	}

	// Create tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Generate encoded tokens
	accessT, err := accessToken.SignedString([]byte(serverConfig.App.AccessTokenSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Message: "Something went wrong"})
	}

	refreshT, err := refreshToken.SignedString([]byte(serverConfig.App.RefreshTokenSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Message: "Something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(responses.LoginTokenResponse{AccessToken: accessT, RefreshToken: refreshT})
}

// allow a teacher to register for an account
// the request body should contain the teacher's email, password, name and age
func RegisterTeacher(c *fiber.Ctx) error {

	credentials := new(requests.TeacherRegistrationRequest)
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide the email, password, name as strings and age as an integer",
		})
	}

	errors := validators.IsValidParams(credentials)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: &errors})
	} else if !validators.IsValidPassword(credentials.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: "Password must contain at least one number, one uppercase letter, one lowercase letter, and one special character"})
	}

	// check if teacher already exists in the db
	temp := models.Teacher{}
	if err := database.Database.Db.Where("email = ?", credentials.Email).First(&temp).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: "Account already exists"})
	}

	// hash the password
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), 8)
	hashedPassword := string(hashedBytes)

	// create a new user in the database
	teacher := models.Teacher{
		Email:    credentials.Email,
		Password: hashedPassword,
		Name:     credentials.Name,
		Age:      credentials.Age,
	}

	// save the user in the database
	if err := database.Database.Db.Create(&teacher).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: "Unable to create user"})
	}

	response := responses.RegisterTeacherResponse{Name: credentials.Name, Email: credentials.Email, Age: credentials.Age}
	return c.Status(fiber.StatusCreated).JSON(response)
}
