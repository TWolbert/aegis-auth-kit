package utils

import (
	"aegis.wlbt.nl/aegis-auth/models"
	"github.com/gofiber/fiber/v2"
)

func GetUserFromContext(c *fiber.Ctx) *models.User {
	userLocal := c.Locals("user")
	if userLocal == nil {
		return nil
	}
	user, ok := userLocal.(*models.User)
	if !ok {
		return nil
	}
	return user
}
