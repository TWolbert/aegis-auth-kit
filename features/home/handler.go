package home

import (
	"time"

	"aegis.wlbt.nl/aegis-auth/database"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func renderTempl(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func IndexHandler(c *fiber.Ctx) error {
	return renderTempl(c, IndexPage())
}

func AboutHandler(c *fiber.Ctx) error {
	return renderTempl(c, AboutPage())
}

func MessageHandler(c *fiber.Ctx) error {
	return c.SendString("Test string hatsaa " + time.Now().Format(time.RFC3339))
}

func DBHealthHandler(c *fiber.Ctx) error {
	sqlDB, err := database.DB.DB()
	if err != nil {
		return renderTempl(c, DBHealthError(err.Error()))
	}

	if err := sqlDB.Ping(); err != nil {
		return renderTempl(c, DBHealthError("Ping failed: "+err.Error()))
	}

	var count int64
	database.DB.Table("users").Count(&count)

	return renderTempl(c, DBHealthSuccess(count, time.Now()))
}
