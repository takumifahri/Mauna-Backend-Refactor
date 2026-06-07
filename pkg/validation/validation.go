package validation

import (
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
)

type FieldError struct {
	Field   string
	Message string
}

type ValidationError struct {
	Fields []FieldError
}

func (e ValidationError) Error() string {
	if len(e.Fields) == 0 {
		return "validation failed"
	}

	messages := make([]string, 0, len(e.Fields))
	for _, field := range e.Fields {
		messages = append(messages, field.Message)
	}
	return strings.Join(messages, "; ")
}

func Validate(input any) error {
	if input == nil {
		return ValidationError{Fields: []FieldError{{Field: "", Message: "request body is required"}}}
	}

	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return ValidationError{Fields: []FieldError{{Field: "", Message: "request body is required"}}}
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil
	}

	typ := value.Type()
	errors := make([]FieldError, 0)

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		rules := field.Tag.Get("validate")
		if rules == "" {
			continue
		}

		fieldName := jsonFieldName(field)
		fieldValue := value.Field(i)
		for _, rule := range strings.Split(rules, ",") {
			if rule == "" {
				continue
			}
			if err := validateRule(fieldName, fieldValue, rule); err != nil {
				errors = append(errors, *err)
				break
			}
		}
	}

	if len(errors) > 0 {
		return ValidationError{Fields: errors}
	}
	return nil
}

func validateRule(fieldName string, value reflect.Value, rule string) *FieldError {
	switch {
	case rule == "required":
		if isEmpty(value) {
			return &FieldError{Field: fieldName, Message: fmt.Sprintf("%s is required", fieldName)}
		}
	case rule == "email":
		email := strings.TrimSpace(stringValue(value))
		address, err := mail.ParseAddress(email)
		if err != nil || address.Address != email {
			return &FieldError{Field: fieldName, Message: fmt.Sprintf("%s must be a valid email address", fieldName)}
		}
	case strings.HasPrefix(rule, "min="):
		min, err := strconv.Atoi(strings.TrimPrefix(rule, "min="))
		if err == nil && len(stringValue(value)) < min {
			return &FieldError{Field: fieldName, Message: fmt.Sprintf("%s must be at least %d characters", fieldName, min)}
		}
	case strings.HasPrefix(rule, "max="):
		max, err := strconv.Atoi(strings.TrimPrefix(rule, "max="))
		if err == nil && len(stringValue(value)) > max {
			return &FieldError{Field: fieldName, Message: fmt.Sprintf("%s must be at most %d characters", fieldName, max)}
		}
	}

	return nil
}

func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return strings.TrimSpace(value.String()) == ""
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map:
		return value.IsNil()
	default:
		return value.IsZero()
	}
}

func stringValue(value reflect.Value) string {
	if value.Kind() == reflect.String {
		return strings.TrimSpace(value.String())
	}
	return strings.TrimSpace(fmt.Sprint(value.Interface()))
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}

	name := strings.Split(tag, ",")[0]
	if name == "" || name == "-" {
		return field.Name
	}
	return name
}
