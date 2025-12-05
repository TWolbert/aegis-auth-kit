package utils

import (
	"net/url"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type UrlParams struct {
	Key     string
	Message string
}

func HTMXRedirect(c *fiber.Ctx, path string, params []UrlParams) error {
	query := "?"

	for _, param := range params {
		query += url.QueryEscape(param.Key) + "=" + url.QueryEscape(param.Message) + "&"
	}

	// remove if no params are supplied
	if query == "?" {
		query = ""
	}

	query = strings.TrimRight(query, "&")

	c.Set("HX-Redirect", path+query)
	return c.SendStatus(fiber.StatusOK)
}

func Redirect(c *fiber.Ctx, path string, params []UrlParams) error {
	query := "?"

	for _, param := range params {
		query += url.QueryEscape(param.Key) + "=" + url.QueryEscape(param.Message) + "&"
	}

	// remove if no params are supplied
	if query == "?" {
		query = ""
	}

	query = strings.TrimRight(query, "&")

	c.Set("HX-Redirect", path+query)
	return c.Redirect(path + query)
}

func RenderTempl(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
