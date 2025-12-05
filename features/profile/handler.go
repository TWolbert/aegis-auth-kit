package profile

import (
	"errors"

	"aegis.wlbt.nl/aegis-auth/database"
	"aegis.wlbt.nl/aegis-auth/features/home"
	"aegis.wlbt.nl/aegis-auth/features/utils"
	"aegis.wlbt.nl/aegis-auth/models"
	v "aegis.wlbt.nl/aegis-auth/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func IndexHandler(c *fiber.Ctx) error {
	statusType := c.Query("statusType", "")
	statusMessage := c.Query("statusMessage", "")
	return utils.RenderTempl(c, ProfilePage(c.Locals("user").(*models.User), home.StatusMessage{
		StatusType:    statusType,
		StatusMessage: statusMessage,
	}))
}

func PostUpdateProfile(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	currentPassword := c.FormValue("currentPassword")

	if username != "" {
		if err := v.Validate(username, v.IsNotEmpty("username"), v.IsMinLength("username", 3), v.IsntExisting("username", models.User{}, "username = ?", username, c.Context())); err != nil {
			return v.ErrorToHTML(c, err)
		}
	}

	if email != "" {
		if err := v.Validate(email, v.IsNotEmpty("email"), v.IsEmail("email"), v.IsntExisting("email", models.User{}, "email = ?", email, c.Context())); err != nil {
			return v.ErrorToHTML(c, err)
		}
	}

	if password != "" {
		if err := v.Validate(password, v.IsNotEmpty("password"), v.IsStrongPassword("password"), v.IsMinLength("password", 8)); err != nil {
			return v.ErrorToHTML(c, err)
		}
	}

	if err := v.Validate(currentPassword, v.IsNotEmpty("currentPassword")); err != nil {
		return v.ErrorToHTML(c, err)
	}
	currentUser := c.Locals("user").(*models.User)

	err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(currentPassword))

	if err != nil {
		return v.ErrorToHTML(c, fiber.NewError(401, "Current password is wrong"))
	}

	ok, err := currentUser.Update(c.Context(), database.DB, username, email, password)

	if err != nil || !ok {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return v.ErrorToHTML(c, fiber.NewError(500, "Username or email already in use!"))
		}

		return v.ErrorToHTML(c, fiber.NewError(500, "Failed to update profile"))
	}

	return utils.HTMXRedirect(c, "/profile", []utils.UrlParams{
		{
			Key: "statusType", Message: "success",
		},
		{
			Key: "statusMessage", Message: "Profile updated",
		},
	})
}
