package routes_cdn

import "github.com/gofiber/fiber/v2"

// Only serves static files

func AlpineJS(c *fiber.Ctx) error {
	return c.SendFile("./js/alpine.min.js", true)
}

func HTMXJS(c *fiber.Ctx) error {
	return c.SendFile("./js/htmx.min.js", true)
}

func TailwindCSS(c *fiber.Ctx) error {
	return c.SendFile("./css/output/output.css", true)
}

func IndexHTML(c *fiber.Ctx) error {
	return c.SendFile("./templates/index.html", true)
}
