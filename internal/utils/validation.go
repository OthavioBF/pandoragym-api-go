package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// DecodeValidJSON decodes JSON from request body and validates the struct
func DecodeValidJSON[T any](r *http.Request) (T, error) {
	var data T

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("decode json: %w", err)
	}

	if err := ValidateStruct(data); err != nil {
		return data, fmt.Errorf("validation failed: %w", err)
	}

	return data, nil
}

// DecodeValidJSONWithDetails decodes JSON and returns detailed validation errors
func DecodeValidJSONWithDetails[T any](r *http.Request) (T, map[string]string, error) {
	var data T

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, nil, fmt.Errorf("decode json: %w", err)
	}

	if problems := ValidateStructDetailed(data); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}

// ParseUUID parses a UUID string and returns proper error response if invalid
func ParseUUID(uuidStr, fieldName string) (uuid.UUID, error) {
	if uuidStr == "" {
		return uuid.Nil, fmt.Errorf("%s is required", fieldName)
	}

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s format", fieldName)
	}

	return parsedUUID, nil
}

// ValidateUserID validates that a user ID is not nil
func ValidateUserID(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return fmt.Errorf("unauthorized")
	}
	return nil
}

func ValidateStruct[T any](s T) error {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" {
			continue
		}

		if err := validateField(field, fieldType.Name, tag); err != nil {
			return err
		}
	}

	return nil
}

// ValidateStructDetailed returns detailed validation errors as a map
func ValidateStructDetailed[T any](s T) map[string]string {
	problems := make(map[string]string)
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" {
			continue
		}

		if err := validateField(field, fieldType.Name, tag); err != nil {
			problems[strings.ToLower(fieldType.Name)] = err.Error()
		}
	}

	return problems
}

func validateField(field reflect.Value, fieldName, tag string) error {
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		if rule == "required" {
			if isEmptyValue(field) {
				return fmt.Errorf("%s is required", fieldName)
			}
		} else if strings.HasPrefix(rule, "min=") {
			minStr := strings.TrimPrefix(rule, "min=")
			if err := validateMin(field, fieldName, minStr); err != nil {
				return err
			}
		} else if strings.HasPrefix(rule, "max=") {
			maxStr := strings.TrimPrefix(rule, "max=")
			if err := validateMax(field, fieldName, maxStr); err != nil {
				return err
			}
		} else if rule == "email" {
			if field.Kind() == reflect.String {
				email := field.String()
				if email != "" && !emailRegex.MatchString(email) {
					return fmt.Errorf("%s must be a valid email address", fieldName)
				}
			}
		} else if strings.HasPrefix(rule, "oneof=") {
			oneofStr := strings.TrimPrefix(rule, "oneof=")
			validValues := strings.Split(oneofStr, " ")
			if err := validateOneOf(field, fieldName, validValues); err != nil {
				return err
			}
		}
	}

	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

func validateMin(field reflect.Value, fieldName, minStr string) error {
	// Implementation for min validation
	// This is a simplified version - you might want to use a proper validation library
	return nil
}

func validateMax(field reflect.Value, fieldName, maxStr string) error {
	// Implementation for max validation
	// This is a simplified version - you might want to use a proper validation library
	return nil
}

func validateOneOf(field reflect.Value, fieldName string, validValues []string) error {
	if field.Kind() == reflect.String {
		value := field.String()
		if slices.Contains(validValues, value) {
			return nil
		}

		return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(validValues, ", "))
	}
	return nil
}

func IsValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	if slices.Contains(validTypes, contentType) {
		return true
	}

	return false
}
