package god

import (
	"reflect"
	"strconv"
	"strings"
)

type BooleanSchema struct {
	BaseSchema
}

func Boolean() *BooleanSchema {
	return &BooleanSchema{
		BaseSchema: BaseSchema{isRequired: true},
	}
}

func Bool() *BooleanSchema {
	return Boolean()
}

func (s *BooleanSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *BooleanSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *BooleanSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *BooleanSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	b, ok := convertToBoolean(processedValue)
	if !ok {
		return ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "expected boolean", Code: "invalid_type", Value: value}},
		}
	}

	return ValidationResult{Valid: true, Value: b}
}

func convertToBoolean(value interface{}) (bool, bool) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool(), true
	case reflect.String:
		s := strings.ToLower(v.String())
		if b, err := strconv.ParseBool(s); err == nil {
			return b, true
		}
		switch s {
		case "yes", "y", "1":
			return true, true
		case "no", "n", "0":
			return false, true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := v.Int()
		if i == 0 {
			return false, true
		} else if i == 1 {
			return true, true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := v.Uint()
		if u == 0 {
			return false, true
		} else if u == 1 {
			return true, true
		}
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if f == 0.0 {
			return false, true
		} else if f == 1.0 {
			return true, true
		}
	}
	return false, false
}