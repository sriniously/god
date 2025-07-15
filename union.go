package god

import (
	"fmt"
	"reflect"
)

type UnionSchema struct {
	BaseSchema
	schemas []Schema
}

func Union(schemas ...Schema) *UnionSchema {
	return &UnionSchema{
		BaseSchema: BaseSchema{isRequired: true},
		schemas:    schemas,
	}
}

func (s *UnionSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *UnionSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *UnionSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *UnionSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	var allErrors []ValidationError

	for i, schema := range s.schemas {
		result := schema.Validate(processedValue)
		if result.Valid {
			return result
		}
		
		for _, err := range result.Errors {
			err.Field = fmt.Sprintf("union[%d]", i)
			allErrors = append(allErrors, err)
		}
	}

	return ValidationResult{
		Valid: false,
		Errors: []ValidationError{{
			Message: fmt.Sprintf("value does not match any of the union types (%d alternatives tried)", len(s.schemas)),
			Code:    "invalid_union",
			Value:   value,
		}},
	}
}

type DiscriminatedUnionSchema struct {
	BaseSchema
	discriminant string
	options      map[string]Schema
}

func DiscriminatedUnion(discriminant string, options map[string]Schema) *DiscriminatedUnionSchema {
	return &DiscriminatedUnionSchema{
		BaseSchema:   BaseSchema{isRequired: true},
		discriminant: discriminant,
		options:      options,
	}
}

func (s *DiscriminatedUnionSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *DiscriminatedUnionSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *DiscriminatedUnionSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *DiscriminatedUnionSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	// Must be an object
	v := reflect.ValueOf(processedValue)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var objMap map[string]interface{}
	var ok bool

	switch v.Kind() {
	case reflect.Map:
		objMap, ok = convertMapToStringInterface(processedValue)
		if !ok {
			return ValidationResult{
				Valid:  false,
				Errors: []ValidationError{{Message: "expected object for discriminated union", Code: "invalid_type", Value: value}},
			}
		}
	case reflect.Struct:
		objMap = structToMap(v)
	default:
		return ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "expected object for discriminated union", Code: "invalid_type", Value: value}},
		}
	}

	// Check for discriminant field
	discriminantValue, exists := objMap[s.discriminant]
	if !exists {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Message: fmt.Sprintf("missing discriminant field '%s'", s.discriminant),
				Code:    "invalid_union",
				Value:   value,
			}},
		}
	}

	discriminantStr := fmt.Sprintf("%v", discriminantValue)
	schema, exists := s.options[discriminantStr]
	if !exists {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Message: fmt.Sprintf("unknown discriminant value '%s'", discriminantStr),
				Code:    "invalid_union",
				Value:   discriminantValue,
			}},
		}
	}

	return schema.Validate(processedValue)
}

type LiteralSchema struct {
	BaseSchema
	value interface{}
}

func Literal(value interface{}) *LiteralSchema {
	return &LiteralSchema{
		BaseSchema: BaseSchema{isRequired: true},
		value:      value,
	}
}

func (s *LiteralSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *LiteralSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *LiteralSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *LiteralSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	if !reflect.DeepEqual(processedValue, s.value) {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Message: fmt.Sprintf("expected literal value %v", s.value),
				Code:    "invalid_literal",
				Value:   value,
			}},
		}
	}

	return ValidationResult{Valid: true, Value: processedValue}
}

type EnumSchema struct {
	BaseSchema
	values []interface{}
}

func Enum(values ...interface{}) *EnumSchema {
	return &EnumSchema{
		BaseSchema: BaseSchema{isRequired: true},
		values:     values,
	}
}

func (s *EnumSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *EnumSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *EnumSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *EnumSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	for _, enumValue := range s.values {
		if reflect.DeepEqual(processedValue, enumValue) {
			return ValidationResult{Valid: true, Value: processedValue}
		}
	}

	return ValidationResult{
		Valid: false,
		Errors: []ValidationError{{
			Message: fmt.Sprintf("expected one of %v", s.values),
			Code:    "invalid_enum_value",
			Value:   value,
		}},
	}
}

type NullableSchema struct {
	BaseSchema
	schema Schema
}

func Nullable(schema Schema) *NullableSchema {
	return &NullableSchema{
		BaseSchema: BaseSchema{isRequired: true},
		schema:     schema,
	}
}

func (s *NullableSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *NullableSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *NullableSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *NullableSchema) Validate(value interface{}) ValidationResult {
	if value == nil {
		return ValidationResult{Valid: true, Value: nil}
	}

	return s.schema.Validate(value)
}