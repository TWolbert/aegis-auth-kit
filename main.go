package main

import (
	"log"
	"time"

	"aegis.wlbt.nl/aegis-auth/database"
	"aegis.wlbt.nl/aegis-auth/templates"
	routes_cdn "aegis.wlbt.nl/aegis-auth/routes/cdn"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func renderTempl(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	// Initialize database (silent mode for child processes)
	database.Connect(fiber.IsChild())

	// Only run migrations in master process
	if !fiber.IsChild() {
		database.Migrate()
	}

	// static
	app.Get("/", routes_cdn.IndexHTML)
	app.Get("/aegis-htmx.js", routes_cdn.HTMXJS)
	app.Get("/aegis-alpine.js", routes_cdn.AlpineJS)
	app.Get("/aegis-styles.css", routes_cdn.TailwindCSS)

	app.Get("/message", func(c *fiber.Ctx) error {
		return c.SendString("Test string hatsaa " + time.Now().Format(time.RFC3339))
	})

	// Database health check
	app.Get("/db/health", func(c *fiber.Ctx) error {
		sqlDB, err := database.DB.DB()
		if err != nil {
			return renderTempl(c, templates.DBHealthError(err.Error()))
		}

		if err := sqlDB.Ping(); err != nil {
			return renderTempl(c, templates.DBHealthError("Ping failed: "+err.Error()))
		}

		var count int64
		database.DB.Table("users").Count(&count)

		return renderTempl(c, templates.DBHealthSuccess(count, time.Now()))
	})

	log.Fatal(app.Listen(":3000"))
}
