package auth

import (
	v "aegis.wlbt.nl/aegis-auth/validation"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
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

	if err := v.Validate(password, v.IsNotEmpty("password"), v.IsMinLength("password", 8), v.IsStrongPassword("password")); err != nil {
		return v.ErrorToHTML(c, err)
	}

	return c.SendString("Login successful - email: " + email)
}

// TODO: Implement actual registration logic
func RegisterPostHandler(c *fiber.Ctx) error {
	return c.SendString("Register endpoint - to be implemented")
}
