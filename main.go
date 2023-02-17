package main

import (
	"fiber-api/database"
	"fiber-api/routes"
	"fiber-api/utils/types"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golobby/dotenv"
)

func main() {
	app := fiber.New()

	config := types.Config{}
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}

	err = dotenv.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	// connect to database
	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)
	database.ConnectDb(dsn)

	// uncomment the line below to seed the database with dummy data
	// database.SeedDb(dsn)

	// setup routes
	routes.Setup(app, config)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.App.Port)))
}
