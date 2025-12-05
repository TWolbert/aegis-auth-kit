package middleware

import (
	"aegis.wlbt.nl/aegis-auth/database"
	"aegis.wlbt.nl/aegis-auth/features/utils"
	"aegis.wlbt.nl/aegis-auth/models"
	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	cookie := c.Cookies("aegis-token")

	if cookie == "" {
		return c.Next()
	}

	user, token, err := models.GetUserByToken(c.Context(), database.DB, cookie)

	if err != nil {
		print("No user found")
		c.ClearCookie("aegis-token")
		return c.Next()
	}

	if token.IsExpired() {
		print("Expired token deleted")
		c.ClearCookie("aegis-token")
		token.Delete(c.Context(), database.DB)
		return c.Next()
	}

	if !token.OwnedByIp(c) {
		print("Token deleted for security concerns")
		c.ClearCookie("aegis-token")
		token.Delete(c.Context(), database.DB)
		return c.Next()
	}

	c.Locals("user", user)
	return c.Next()
}

// Pass after Auth() Middleware
func RequiresAuth(c *fiber.Ctx) error {
	user := c.Locals("user")

	if user == nil {
		return utils.Redirect(c, "/", []utils.UrlParams{
			{Key: "statusType", Message: "error"},
			{Key: "statusMessage", Message: "You have to be signed in to view this page."},
		})
	}

	return c.Next()
}
