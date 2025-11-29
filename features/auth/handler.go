package auth

import (
	"aegis.wlbt.nl/aegis-auth/database"
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

	return c.SendString("Login successful - email: " + email)
}

// TODO: Implement actual registration logic
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

	// Redirect to home with success message using HX-Redirect header for HTMX
	c.Set("HX-Redirect", "/?statusType=success&statusMessage=Registration+successful!+Welcome+to+Aegis.")
	return c.SendStatus(fiber.StatusOK)
}
