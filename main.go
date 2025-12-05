package main

import (
	"log"

	"aegis.wlbt.nl/aegis-auth/database"
	"aegis.wlbt.nl/aegis-auth/features/auth"
	"aegis.wlbt.nl/aegis-auth/features/home"
	"aegis.wlbt.nl/aegis-auth/features/middleware"
	routes_cdn "aegis.wlbt.nl/aegis-auth/routes/cdn"
	"github.com/gofiber/fiber/v2"
)

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

	app.Use(middleware.Auth)

	// Static assets (CDN)
	app.Get("/aegis-htmx.js", routes_cdn.HTMXJS)
	app.Get("/aegis-alpine.js", routes_cdn.AlpineJS)
	app.Get("/aegis-styles.css", routes_cdn.TailwindCSS)

	// Home feature
	app.Get("/", home.IndexHandler)
	app.Get("/about", home.AboutHandler)
	app.Get("/message", home.MessageHandler)
	app.Get("/db/health", home.DBHealthHandler)
	app.Get("/api/navbar-user", home.NavbarUserHandler)

	// Auth feature
	app.Get("/login", auth.LoginHandler)
	app.Get("/register", auth.RegisterHandler)
	app.Post("/login", auth.LoginPostHandler)
	app.Post("/register", auth.RegisterPostHandler)
	app.Get("/logout", auth.LogoutHandler)

	log.Fatal(app.Listen(":3000"))
}
