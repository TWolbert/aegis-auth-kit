package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html", true)
	})

	app.Get("/message", func(c *fiber.Ctx) error {
		return c.SendString("Test string hatsaa")
	})

	app.Get("/aegis-htmx.js", func(c *fiber.Ctx) error {
		return c.SendFile("./js/htmx.min.js", true)
	})

	app.Get("/aegis-alpine.js", func(c *fiber.Ctx) error {
		return c.SendFile("./js/alpine.min.js", true)
	})

	app.Get("/aegis-styles.css", func(c *fiber.Ctx) error {
		return c.SendFile("./css/output.css", true)
	})

	log.Fatal(app.Listen(":3000"))
}
