package middleware

import (
	"aegis.wlbt.nl/aegis-auth/database"
	"aegis.wlbt.nl/aegis-auth/models"
	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	cookie := c.Cookies("aegis-token")

	if cookie == "" {
		print("No cookie")
		return c.Next()
	}

	println(cookie)

	user, err := models.GetUserByToken(c.Context(), database.DB, cookie)

	if err != nil {
		print("No user found")
		return c.Next()
	}

	println(user.ID)
	c.Locals("user", user)
	return c.Next()
}
