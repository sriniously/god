package god

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type NumberSchema struct {
	BaseSchema
	min       *float64
	max       *float64
	int       bool
	positive  bool
	negative  bool
	nonNeg    bool
	nonPos    bool
	finite    bool
	safe      bool
	multipleOf *float64
}

func Number() *NumberSchema {
	return &NumberSchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func Int() *NumberSchema {
	return &NumberSchema{
		BaseSchema: BaseSchema{isRequired: true},
		int:        true,
	}
}

func Float() *NumberSchema {
	return &NumberSchema{
		BaseSchema: BaseSchema{isRequired: true},
		int:        false,
	}
}

func (s *NumberSchema) Min(value float64) *NumberSchema {
	s.min = &value
	return s
}

func (s *NumberSchema) Max(value float64) *NumberSchema {
	s.max = &value
	return s
}

func (s *NumberSchema) Positive() *NumberSchema {
	s.positive = true
	return s
}

func (s *NumberSchema) Negative() *NumberSchema {
	s.negative = true
	return s
}

func (s *NumberSchema) NonNegative() *NumberSchema {
	s.nonNeg = true
	return s
}

func (s *NumberSchema) NonPositive() *NumberSchema {
	s.nonPos = true
	return s
}

func (s *NumberSchema) Finite() *NumberSchema {
	s.finite = true
	return s
}

func (s *NumberSchema) Safe() *NumberSchema {
	s.safe = true
	return s
}

func (s *NumberSchema) MultipleOf(value float64) *NumberSchema {
	s.multipleOf = &value
	return s
}

func (s *NumberSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *NumberSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *NumberSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *NumberSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	num, ok := convertToFloat64(processedValue)
	if !ok {
		return ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "expected number", Code: "invalid_type", Value: value}},
		}
	}

	var errors []ValidationError

	if s.int && !isInteger(num) {
		errors = append(errors, ValidationError{
			Message: "expected integer",
			Code:    "invalid_type",
			Value:   num,
		})
	}

	if s.min != nil && num < *s.min {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("number must be greater than or equal to %g", *s.min),
			Code:    "too_small",
			Value:   num,
		})
	}

	if s.max != nil && num > *s.max {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("number must be less than or equal to %g", *s.max),
			Code:    "too_big",
			Value:   num,
		})
	}

	if s.positive && num <= 0 {
		errors = append(errors, ValidationError{
			Message: "number must be positive",
			Code:    "too_small",
			Value:   num,
		})
	}

	if s.negative && num >= 0 {
		errors = append(errors, ValidationError{
			Message: "number must be negative",
			Code:    "too_big",
			Value:   num,
		})
	}

	if s.nonNeg && num < 0 {
		errors = append(errors, ValidationError{
			Message: "number must be non-negative",
			Code:    "too_small",
			Value:   num,
		})
	}

	if s.nonPos && num > 0 {
		errors = append(errors, ValidationError{
			Message: "number must be non-positive",
			Code:    "too_big",
			Value:   num,
		})
	}

	if s.finite && (math.IsInf(num, 0) || math.IsNaN(num)) {
		errors = append(errors, ValidationError{
			Message: "number must be finite",
			Code:    "invalid_type",
			Value:   num,
		})
	}

	if s.safe && (num > 9007199254740991 || num < -9007199254740991) {
		errors = append(errors, ValidationError{
			Message: "number must be a safe integer",
			Code:    "too_big",
			Value:   num,
		})
	}

	if s.multipleOf != nil && math.Mod(num, *s.multipleOf) != 0 {
		errors = append(errors, ValidationError{
			Message: fmt.Sprintf("number must be a multiple of %g", *s.multipleOf),
			Code:    "invalid_type",
			Value:   num,
		})
	}

	if len(errors) > 0 {
		return ValidationResult{Valid: false, Errors: errors}
	}

	if s.int {
		return ValidationResult{Valid: true, Value: int64(num)}
	}

	return ValidationResult{Valid: true, Value: num}
}

func convertToFloat64(value interface{}) (float64, bool) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	case reflect.String:
		if f, err := parseFloat(v.String()); err == nil {
			return f, true
		}
	}
	return 0, false
}

func isInteger(num float64) bool {
	return num == math.Trunc(num)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}