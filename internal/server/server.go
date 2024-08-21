package server

import (
	"os"

	"github.com/SatvikG7/Scalable-Notification-System/internal/handlers"
	"github.com/SatvikG7/Scalable-Notification-System/internal/handlers/users"

	"github.com/gofiber/fiber/v2"
)

func Init() error {
	app := fiber.New(
		fiber.Config{
			AppName: "Scalable Notification System",
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Scalable Notification System")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/users", users.GetUsers)
	app.Get("/users/get/:id", users.GetUser)
	app.Post("/users/create", users.CreateUser)
	app.Delete("/users/delete/:id", users.DeleteUser)

	app.Post("/notify", handlers.CreateNotification)

	app.Listen("127.0.0.1:" + os.Getenv("PORT"))

	return nil
}
