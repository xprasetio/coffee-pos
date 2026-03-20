package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator.Validate instance
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate validates a struct and returns error messages
// Returns nil if validation passes, or map[field]error_message if failed
func (v *Validator) Validate(i interface{}) map[string]string {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := v.getFieldName(i, err.Field())
		errors[field] = v.getErrorMessage(err)
	}

	return errors
}

// getFieldName extracts the JSON tag name from the struct field
func (v *Validator) getFieldName(i interface{}, fieldName string) string {
	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				// Handle json tag with options (e.g., "name,omitempty")
				parts := strings.Split(jsonTag, ",")
				return parts[0]
			}
			// Convert camelCase to snake_case if no json tag
			return toSnakeCase(fieldName)
		}
	}

	return toSnakeCase(fieldName)
}

// getErrorMessage returns human-readable error message in Indonesian
func (v *Validator) getErrorMessage(err validator.FieldError) string {
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return "wajib diisi"
	case "min":
		return fmt.Sprintf("minimal %s", param)
	case "max":
		return fmt.Sprintf("maksimal %s", param)
	case "email":
		return "format email tidak valid"
	case "oneof":
		return fmt.Sprintf("harus salah satu dari: %s", strings.ReplaceAll(param, " ", ", "))
	case "uuid", "uuid4":
		return "format UUID tidak valid"
	case "len":
		return fmt.Sprintf("harus tepat %s karakter", param)
	case "numeric":
		return "harus berupa angka"
	case "url":
		return "format URL tidak valid"
	case "alpha":
		return "hanya boleh mengandung huruf"
	case "alphanum":
		return "hanya boleh mengandung huruf dan angka"
	case "gt":
		return fmt.Sprintf("harus lebih besar dari %s", param)
	case "lt":
		return fmt.Sprintf("harus lebih kecil dari %s", param)
	case "gte":
		return fmt.Sprintf("harus lebih besar atau sama dengan %s", param)
	case "lte":
		return fmt.Sprintf("harus lebih kecil atau sama dengan %s", param)
	case "eq":
		return fmt.Sprintf("harus sama dengan %s", param)
	case "ne":
		return fmt.Sprintf("tidak boleh sama dengan %s", param)
	case "contains":
		return fmt.Sprintf("harus mengandung %s", param)
	case "startswith":
		return fmt.Sprintf("harus diawali dengan %s", param)
	case "endswith":
		return fmt.Sprintf("harus diakhiri dengan %s", param)
	default:
		return "tidak valid"
	}
}

// toSnakeCase converts camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
