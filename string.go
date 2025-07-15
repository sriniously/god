package god

import (
	"fmt"
	"regexp"
	"strings"
)

type StringSchema struct {
	BaseSchema
	minLength *int
	maxLength *int
	pattern   *regexp.Regexp
	email     bool
	url       bool
	uuid      bool
	transform func(string) string
}

func String() *StringSchema {
	return &StringSchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func (s *StringSchema) Min(length int) *StringSchema {
	s.minLength = &length
	return s
}

func (s *StringSchema) Max(length int) *StringSchema {
	s.maxLength = &length
	return s
}

func (s *StringSchema) Length(length int) *StringSchema {
	s.minLength = &length
	s.maxLength = &length
	return s
}

func (s *StringSchema) Regex(pattern string) *StringSchema {
	s.pattern = regexp.MustCompile(pattern)
	return s
}

func (s *StringSchema) Email() *StringSchema {
	s.email = true
	return s
}

func (s *StringSchema) URL() *StringSchema {
	s.url = true
	return s
}

func (s *StringSchema) UUID() *StringSchema {
	s.uuid = true
	return s
}

func (s *StringSchema) Transform(fn func(string) string) *StringSchema {
	s.transform = fn
	return s
}

func (s *StringSchema) Trim() *StringSchema {
	s.transform = strings.TrimSpace
	return s
}

func (s *StringSchema) ToLower() *StringSchema {
	s.transform = strings.ToLower
	return s
}

func (s *StringSchema) ToUpper() *StringSchema {
	s.transform = strings.ToUpper
	return s
}

func (s *StringSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *StringSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *StringSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *StringSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	str, ok := processedValue.(string)
	if !ok {
		return ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "expected string", Code: "invalid_type", Value: value}},
		}
	}

	if s.transform != nil {
		str = s.transform(str)
	}

	var errors []ValidationError

	if s.minLength != nil && len(str) < *s.minLength {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("string must be at least %d characters", *s.minLength),
			Code:    "too_small",
			Value:   str,
		})
	}

	if s.maxLength != nil && len(str) > *s.maxLength {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("string must be at most %d characters", *s.maxLength),
			Code:    "too_big",
			Value:   str,
		})
	}

	if s.pattern != nil && !s.pattern.MatchString(str) {
		errors = append(errors, ValidationError{
			Message: "string does not match required pattern",
			Code:    "invalid_string",
			Value:   str,
		})
	}

	if s.email && !isValidEmail(str) {
		errors = append(errors, ValidationError{
			Message: "invalid email format",
			Code:    "invalid_string",
			Value:   str,
		})
	}

	if s.url && !isValidURL(str) {
		errors = append(errors, ValidationError{
			Message: "invalid URL format",
			Code:    "invalid_string",
			Value:   str,
		})
	}

	if s.uuid && !isValidUUID(str) {
		errors = append(errors, ValidationError{
			Message: "invalid UUID format",
			Code:    "invalid_string",
			Value:   str,
		})
	}

	if len(errors) > 0 {
		return ValidationResult{Valid: false, Errors: errors}
	}

	return ValidationResult{Valid: true, Value: str}
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

func isValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(uuid))
}