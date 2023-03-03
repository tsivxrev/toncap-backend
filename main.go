package main

import (
	"log"
	"os"
	"toncap-backend/controller"
	"toncap-backend/database"
	"toncap-backend/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	//"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Toncap v3",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return controller.Error(c, fiber.StatusInternalServerError, err)
		},
		CaseSensitive: true,
		Prefork:       true,
	})

	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())
	//app.Use(limiter.New())

	router.Setup(app)

	log.Fatal(app.Listen(os.Getenv("ADDRESS")))
}
