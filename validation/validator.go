package v

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Rule func(input any) (bool, string, string)

// ValidationError holds validation error details
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Validate checks input against rules and returns a ValidationError if validation fails
func Validate(input any, rules ...Rule) error {
	for _, rule := range rules {
		if validated, msg, field := rule(input); msg != "" || !validated {
			return &ValidationError{
				Message: msg,
				Field:   field,
			}
		}
	}

	return nil
}

func IsNotEmpty(field string) Rule {
	return func(value any) (bool, string, string) {
		valueString, ok := value.(string)
		if !ok || strings.TrimSpace(valueString) == "" {
			return false, underscoresToSpaces(capitalizeFirst(field)) + " is not allowed to be empty!", field
		}

		return true, "", ""
	}
}

func IsMinLength(field string, minLength int) Rule {
	return func(value any) (bool, string, string) {
		valueString, ok := value.(string)

		if !ok || len(valueString) < minLength {
			return false, underscoresToSpaces(capitalizeFirst(field)) + " has to be at least " + strconv.Itoa(minLength) + " characters", field
		}

		return true, "", ""
	}
}

func IsMaxLength(field string, maxLength int) Rule {
	return func(value any) (bool, string, string) {
		valueString, ok := value.(string)

		if !ok || len(valueString) > maxLength {
			return false, underscoresToSpaces(capitalizeFirst(field)) + " cannot be longer than " + strconv.Itoa(maxLength) + " characters", field
		}

		return true, "", ""
	}
}

func IsEmail(field string) Rule {
	return func(value any) (bool, string, string) {
		valueString, ok := value.(string)
		regex := "(?:[a-z0-9!#$%&'*+\\x2f=?^_`\\x7b-\\x7d~\\x2d]+(?:\\.[a-z0-9!#$%&'*+\\x2f=?^_`\\x7b-\\x7d~\\x2d]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9\\x2d]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9\\x2d]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9\\x2d]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"

		r, _ := regexp.Compile(regex)

		if !ok || !r.MatchString(valueString) {
			return false, underscoresToSpaces(capitalizeFirst(field)) + " is not an email!", field
		}

		return true, "", ""
	}
}

func IsStrongPassword(field string) Rule {
	strongPasswordLower := regexp.MustCompile(`[a-z]`)
	strongPasswordUpper := regexp.MustCompile(`[A-Z]`)
	strongPasswordDigit := regexp.MustCompile(`[0-9]`)
	strongPasswordSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)

	return func(value any) (bool, string, string) {
		valueString, ok := value.(string)
		if !ok {
			return false, underscoresToSpaces(capitalizeFirst(field)) + " must be a string", field
		}

		var errors []string

		if !strongPasswordLower.MatchString(valueString) {
			errors = append(errors, "one lowercase letter")
		}

		if !strongPasswordUpper.MatchString(valueString) {
			errors = append(errors, "one uppercase letter")
		}

		if !strongPasswordDigit.MatchString(valueString) {
			errors = append(errors, "one digit")
		}

		if !strongPasswordSpecial.MatchString(valueString) {
			errors = append(errors, "one special character")
		}

		if len(errors) > 0 {
			msg := underscoresToSpaces(capitalizeFirst(field)) + " must contain " + strings.Join(errors, ", ")
			return false, msg, field
		}

		return true, "", ""
	}
}

func capitalizeFirst(input string) string {
	capitalized := ""

	for i, r := range input {
		r := string(r)
		if i == 0 {
			capitalized += strings.ToUpper(r)
		} else {
			capitalized += r
		}
	}

	return capitalized
}

func underscoresToSpaces(input string) string {
	return strings.ReplaceAll(input, "_", " ")
}

func ErrorToHTML(c *fiber.Ctx, err error) error {
	if validationErr, ok := err.(*ValidationError); ok {
		html := fmt.Sprintf(
			`<div class="error" data-field="%s">%s</div>`,
			validationErr.Field,
			validationErr.Message,
		)
		return c.Type("html").SendString(html)
	}

	html := fmt.Sprintf(`<div class="error">%s</div>`, err.Error())
	return c.Type("html").SendString(html)
}
