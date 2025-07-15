package god

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
	Code    string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
	Value  interface{}
}

func (r ValidationResult) Error() error {
	if r.Valid {
		return nil
	}
	var messages []string
	for _, err := range r.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
}

type Schema interface {
	Validate(value interface{}) ValidationResult
	Optional() Schema
	Required() Schema
	Default(value interface{}) Schema
}

type BaseSchema struct {
	isOptional   bool
	isRequired   bool
	defaultValue interface{}
	hasDefault   bool
}

func (s *BaseSchema) setOptional() {
	s.isOptional = true
	s.isRequired = false
}

func (s *BaseSchema) setRequired() {
	s.isRequired = true
	s.isOptional = false
}

func (s *BaseSchema) setDefault(value interface{}) {
	s.defaultValue = value
	s.hasDefault = true
}

func (s *BaseSchema) handleNil(value interface{}) (interface{}, bool, ValidationResult) {
	if value == nil {
		if s.hasDefault {
			return s.defaultValue, false, ValidationResult{Valid: true, Value: s.defaultValue}
		}
		if s.isOptional {
			return nil, true, ValidationResult{Valid: true, Value: nil}
		}
		if s.isRequired {
			return nil, true, ValidationResult{
				Valid:  false,
				Errors: []ValidationError{{Message: "field is required", Code: "required"}},
			}
		}
		return nil, true, ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "field is required", Code: "required"}},
		}
	}
	return value, false, ValidationResult{}
}