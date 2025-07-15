package god

import (
	"fmt"
	"reflect"
)

type ArraySchema struct {
	BaseSchema
	element   Schema
	minLength *int
	maxLength *int
	length    *int
	nonempty  bool
}

func Array(element Schema) *ArraySchema {
	return &ArraySchema{
		BaseSchema: BaseSchema{isRequired: true},
		element:    element,
	}
}

func (s *ArraySchema) Min(length int) *ArraySchema {
	s.minLength = &length
	return s
}

func (s *ArraySchema) Max(length int) *ArraySchema {
	s.maxLength = &length
	return s
}

func (s *ArraySchema) Length(length int) *ArraySchema {
	s.length = &length
	return s
}

func (s *ArraySchema) Nonempty() *ArraySchema {
	s.nonempty = true
	return s
}

func (s *ArraySchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *ArraySchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *ArraySchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *ArraySchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	v := reflect.ValueOf(processedValue)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "expected array", Code: "invalid_type", Value: value}},
		}
	}

	length := v.Len()
	var errors []ValidationError

	if s.length != nil && length != *s.length {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("array must have exactly %d elements", *s.length),
			Code:    "invalid_type",
			Value:   value,
		})
	}

	if s.minLength != nil && length < *s.minLength {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("array must have at least %d elements", *s.minLength),
			Code:    "too_small",
			Value:   value,
		})
	}

	if s.maxLength != nil && length > *s.maxLength {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("array must have at most %d elements", *s.maxLength),
			Code:    "too_big",
			Value:   value,
		})
	}

	if s.nonempty && length == 0 {
		errors = append(errors, ValidationError{
			Message: "array must not be empty",
			Code:    "too_small",
			Value:   value,
		})
	}

	validatedArray := make([]interface{}, length)
	for i := 0; i < length; i++ {
		elementValue := v.Index(i).Interface()
		result := s.element.Validate(elementValue)
		if !result.Valid {
			for _, err := range result.Errors {
				err.Field = fmt.Sprintf("[%d]", i)
				errors = append(errors, err)
			}
		} else {
			validatedArray[i] = result.Value
		}
	}

	if len(errors) > 0 {
		return ValidationResult{Valid: false, Errors: errors}
	}

	return ValidationResult{Valid: true, Value: validatedArray}
}

type TupleSchema struct {
	BaseSchema
	elements []Schema
	rest     Schema
}

func Tuple(elements ...Schema) *TupleSchema {
	return &TupleSchema{
		BaseSchema: BaseSchema{isRequired: true},
		elements:   elements,
	}
}

func (s *TupleSchema) Rest(schema Schema) *TupleSchema {
	s.rest = schema
	return s
}

func (s *TupleSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *TupleSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *TupleSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *TupleSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	v := reflect.ValueOf(processedValue)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "expected tuple", Code: "invalid_type", Value: value}},
		}
	}

	length := v.Len()
	var errors []ValidationError

	if s.rest == nil && length != len(s.elements) {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("tuple must have exactly %d elements", len(s.elements)),
			Code:    "invalid_type",
			Value:   value,
		})
	}

	if s.rest != nil && length < len(s.elements) {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("tuple must have at least %d elements", len(s.elements)),
			Code:    "too_small",
			Value:   value,
		})
	}

	validatedTuple := make([]interface{}, length)

	// Validate fixed elements
	for i, elementSchema := range s.elements {
		if i >= length {
			break
		}
		elementValue := v.Index(i).Interface()
		result := elementSchema.Validate(elementValue)
		if !result.Valid {
			for _, err := range result.Errors {
				err.Field = fmt.Sprintf("[%d]", i)
				errors = append(errors, err)
			}
		} else {
			validatedTuple[i] = result.Value
		}
	}

	// Validate rest elements
	if s.rest != nil {
		for i := len(s.elements); i < length; i++ {
			elementValue := v.Index(i).Interface()
			result := s.rest.Validate(elementValue)
			if !result.Valid {
				for _, err := range result.Errors {
					err.Field = fmt.Sprintf("[%d]", i)
					errors = append(errors, err)
				}
			} else {
				validatedTuple[i] = result.Value
			}
		}
	}

	if len(errors) > 0 {
		return ValidationResult{Valid: false, Errors: errors}
	}

	return ValidationResult{Valid: true, Value: validatedTuple}
}