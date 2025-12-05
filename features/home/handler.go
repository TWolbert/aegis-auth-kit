package home

import (
	"time"

	"aegis.wlbt.nl/aegis-auth/database"
	"aegis.wlbt.nl/aegis-auth/features/utils"
	"aegis.wlbt.nl/aegis-auth/templates"
	"github.com/gofiber/fiber/v2"
)

type StatusMessage struct {
	StatusType    string
	StatusMessage string
}

func IndexHandler(c *fiber.Ctx) error {
	statusType := c.Query("statusType", "")
	statusMessage := c.Query("statusMessage", "")

	return utils.RenderTempl(c, IndexPage(StatusMessage{
		StatusType:    statusType,
		StatusMessage: statusMessage,
	}))
}

func AboutHandler(c *fiber.Ctx) error {
	return utils.RenderTempl(c, AboutPage())
}

func MessageHandler(c *fiber.Ctx) error {
	return c.SendString("Test string hatsaa " + time.Now().Format(time.RFC3339))
}

func DBHealthHandler(c *fiber.Ctx) error {
	sqlDB, err := database.DB.DB()
	if err != nil {
		return utils.RenderTempl(c, DBHealthError(err.Error()))
	}

	if err := sqlDB.Ping(); err != nil {
		return utils.RenderTempl(c, DBHealthError("Ping failed: "+err.Error()))
	}

	var count int64
	database.DB.Table("users").Count(&count)

	return utils.RenderTempl(c, DBHealthSuccess(count, time.Now()))
}

func NavbarUserHandler(c *fiber.Ctx) error {
	user := utils.GetUserFromContext(c)
	return utils.RenderTempl(c, templates.NavbarUser(user))
}
