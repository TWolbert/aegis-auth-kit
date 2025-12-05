package auth

import (
	"time"

	"aegis.wlbt.nl/aegis-auth/database"
	"aegis.wlbt.nl/aegis-auth/features/utils"
	"aegis.wlbt.nl/aegis-auth/models"
	v "aegis.wlbt.nl/aegis-auth/validation"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func renderTempl(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func LoginHandler(c *fiber.Ctx) error {
	return renderTempl(c, LoginPage())
}

func RegisterHandler(c *fiber.Ctx) error {
	return renderTempl(c, RegisterPage())
}

// TODO: Implement actual login logic
func LoginPostHandler(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := v.Validate(email, v.IsNotEmpty("email"), v.IsEmail("email")); err != nil {
		return v.ErrorToHTML(c, err)
	}

	if err := v.Validate(password, v.IsNotEmpty("password")); err != nil {
		return v.ErrorToHTML(c, err)
	}

	user, err := gorm.G[models.User](database.DB).Where("email = ?", email).First(c.Context())

	if err != nil || user.ID == 0 {
		return v.ErrorToHTML(c, fiber.NewError(404, "User not found!"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return v.ErrorToHTML(c, fiber.NewError(401, "Password incorrect!"))
	}

	token, err := models.CreateToken(c.Context(), database.DB, user, time.Now().AddDate(0, 1, 0), c.IP(), string(c.Context().UserAgent()))

	if err != nil {
		return v.ErrorToHTML(c, fiber.NewError(500, err.Error()))
	}

	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Name:     "aegis-token",
		Value:    token.Token,
		Expires:  token.ExpiresAt,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	// Redirect to home with success message using HX-Redirect header for HTMX
	c.Set("HX-Redirect", "/?statusType=success&statusMessage=Login+successful!+Welcome+back+ "+user.Username+".")
	return c.SendStatus(fiber.StatusOK)
}

func RegisterPostHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := v.Validate(username, v.IsNotEmpty("username"), v.IsMinLength("username", 3)); err != nil {
		return v.ErrorToHTML(c, err)
	}

	if err := v.Validate(email, v.IsNotEmpty("email"), v.IsEmail("email")); err != nil {
		return v.ErrorToHTML(c, err)
	}

	if err := v.Validate(password, v.IsNotEmpty("password"), v.IsStrongPassword("password"), v.IsMinLength("password", 8)); err != nil {
		return v.ErrorToHTML(c, err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return v.ErrorToHTML(c, fiber.NewError(500, "Hashing error"))
	}

	user := models.User{Username: username, Email: email, Password: string(hash)}

	err = gorm.G[models.User](database.DB).Create(c.Context(), &user)

	if err != nil {
		return v.ErrorToHTML(c, fiber.NewError(500, err.Error()))
	}

	token, err := models.CreateToken(c.Context(), database.DB, user, time.Now().AddDate(0, 1, 0), c.IP(), string(c.Context().UserAgent()))

	if err != nil {
		return v.ErrorToHTML(c, fiber.NewError(500, err.Error()))
	}

	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Name:     "aegis-token",
		Value:    token.Token,
		Expires:  token.ExpiresAt,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	// Redirect to home with success message using HX-Redirect header for HTMX
	c.Set("HX-Redirect", "/?statusType=success&statusMessage=Registration+successful!+Welcome+To+Aegis.")
	return c.SendStatus(fiber.StatusOK)
}

func LogoutHandler(c *fiber.Ctx) error {
	user := utils.GetUserFromContext(c)

	for _, token := range user.SessionTokens {
		token.Delete(c.Context(), database.DB)
	}

	c.ClearCookie("aegis-token")

	c.Set("HX-Redirect", "/?statusType=success&statusMessage=Logged+out+succesfully.")

	return c.SendStatus(fiber.StatusOK)
}
