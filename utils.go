package god

import (
	"fmt"
	"time"
)

type AnySchema struct {
	BaseSchema
}

func Any() *AnySchema {
	return &AnySchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func (s *AnySchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *AnySchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *AnySchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *AnySchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	return ValidationResult{Valid: true, Value: processedValue}
}

type UnknownSchema struct {
	BaseSchema
}

func Unknown() *UnknownSchema {
	return &UnknownSchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func (s *UnknownSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *UnknownSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *UnknownSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *UnknownSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	return ValidationResult{Valid: true, Value: processedValue}
}

type VoidSchema struct {
	BaseSchema
}

func Void() *VoidSchema {
	return &VoidSchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func (s *VoidSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *VoidSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *VoidSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *VoidSchema) Validate(value interface{}) ValidationResult {
	_, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	return ValidationResult{Valid: true, Value: nil}
}

type NeverSchema struct {
	BaseSchema
}

func Never() *NeverSchema {
	return &NeverSchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func (s *NeverSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *NeverSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *NeverSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *NeverSchema) Validate(value interface{}) ValidationResult {
	return ValidationResult{
		Valid: false,
		Errors: []ValidationError{{
			Message: "never type should never be used",
			Code:    "invalid_type",
			Value:   value,
		}},
	}
}

type DateSchema struct {
	BaseSchema
	min *time.Time
	max *time.Time
}

func Date() *DateSchema {
	return &DateSchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func (s *DateSchema) Min(date time.Time) *DateSchema {
	s.min = &date
	return s
}

func (s *DateSchema) Max(date time.Time) *DateSchema {
	s.max = &date
	return s
}

func (s *DateSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *DateSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *DateSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *DateSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	var date time.Time
	var ok bool

	switch v := processedValue.(type) {
	case time.Time:
		date = v
		ok = true
	case string:
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			date = parsed
			ok = true
		} else if parsed, err := time.Parse("2006-01-02", v); err == nil {
			date = parsed
			ok = true
		}
	}

	if !ok {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Message: "expected valid date",
				Code:    "invalid_date",
				Value:   value,
			}},
		}
	}

	var errors []ValidationError

	if s.min != nil && date.Before(*s.min) {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("date must be after %s", s.min.Format(time.RFC3339)),
			Code:    "too_small",
			Value:   date,
		})
	}

	if s.max != nil && date.After(*s.max) {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("date must be before %s", s.max.Format(time.RFC3339)),
			Code:    "too_big",
			Value:   date,
		})
	}

	if len(errors) > 0 {
		return ValidationResult{Valid: false, Errors: errors}
	}

	return ValidationResult{Valid: true, Value: date}
}

func Lazy(schemaFn func() Schema) Schema {
	return &LazySchema{
		BaseSchema: BaseSchema{isRequired: true},
		schemaFn:   schemaFn,
	}
}

type LazySchema struct {
	BaseSchema
	schemaFn func() Schema
	cached   Schema
}

func (s *LazySchema) getSchema() Schema {
	if s.cached == nil {
		s.cached = s.schemaFn()
	}
	return s.cached
}

func (s *LazySchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *LazySchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *LazySchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *LazySchema) Validate(value interface{}) ValidationResult {
	_, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	return s.getSchema().Validate(value)
}