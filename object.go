package god

import (
	"fmt"
	"reflect"
	"strings"
)

type ObjectSchema struct {
	BaseSchema
	fields        map[string]Schema
	strict        bool
	passthrough   bool
	catchall      Schema
	shape         map[string]Schema
	keyof         []string
	partial       bool
	deepPartial   bool
	required      []string
	pick          []string
	omit          []string
	extend        map[string]Schema
	merge         *ObjectSchema
}

func Object(fields map[string]Schema) *ObjectSchema {
	return &ObjectSchema{
		BaseSchema: BaseSchema{isRequired: true},
		fields:     fields,
		shape:      fields,
	}
}

func (s *ObjectSchema) Strict() *ObjectSchema {
	s.strict = true
	s.passthrough = false
	return s
}

func (s *ObjectSchema) Passthrough() *ObjectSchema {
	s.passthrough = true
	s.strict = false
	return s
}

func (s *ObjectSchema) Catchall(schema Schema) *ObjectSchema {
	s.catchall = schema
	return s
}

func (s *ObjectSchema) Partial() *ObjectSchema {
	s.partial = true
	return s
}

func (s *ObjectSchema) DeepPartial() *ObjectSchema {
	s.deepPartial = true
	return s
}

func (s *ObjectSchema) RequiredFields(fields ...string) *ObjectSchema {
	s.required = append(s.required, fields...)
	return s
}

func (s *ObjectSchema) Pick(fields ...string) *ObjectSchema {
	s.pick = fields
	return s
}

func (s *ObjectSchema) Omit(fields ...string) *ObjectSchema {
	s.omit = fields
	return s
}

func (s *ObjectSchema) Extend(fields map[string]Schema) *ObjectSchema {
	if s.extend == nil {
		s.extend = make(map[string]Schema)
	}
	for k, v := range fields {
		s.extend[k] = v
	}
	return s
}

func (s *ObjectSchema) Merge(other *ObjectSchema) *ObjectSchema {
	s.merge = other
	return s
}

func (s *ObjectSchema) Keyof() []string {
	var keys []string
	for key := range s.getEffectiveFields() {
		keys = append(keys, key)
	}
	return keys
}

func (s *ObjectSchema) Optional() Schema {
	s.BaseSchema.setOptional()
	return s
}

func (s *ObjectSchema) Required() Schema {
	s.BaseSchema.setRequired()
	return s
}

func (s *ObjectSchema) Default(value interface{}) Schema {
	s.BaseSchema.setDefault(value)
	return s
}

func (s *ObjectSchema) getEffectiveFields() map[string]Schema {
	fields := make(map[string]Schema)
	
	// Start with base fields
	for k, v := range s.fields {
		fields[k] = v
	}
	
	// Apply merge
	if s.merge != nil {
		for k, v := range s.merge.fields {
			fields[k] = v
		}
	}
	
	// Apply extend
	if s.extend != nil {
		for k, v := range s.extend {
			fields[k] = v
		}
	}
	
	// Apply pick
	if len(s.pick) > 0 {
		picked := make(map[string]Schema)
		for _, key := range s.pick {
			if schema, exists := fields[key]; exists {
				picked[key] = schema
			}
		}
		fields = picked
	}
	
	// Apply omit
	if len(s.omit) > 0 {
		for _, key := range s.omit {
			delete(fields, key)
		}
	}
	
	// Apply partial
	if s.partial || s.deepPartial {
		for k, v := range fields {
			fields[k] = v.Optional()
		}
	}
	
	// Apply required
	if len(s.required) > 0 {
		for _, key := range s.required {
			if schema, exists := fields[key]; exists {
				fields[key] = schema.Required()
			}
		}
	}
	
	return fields
}

func (s *ObjectSchema) Validate(value interface{}) ValidationResult {
	processedValue, shouldReturn, result := s.handleNil(value)
	if shouldReturn {
		return result
	}

	// Check if value is a map or struct
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
				Errors: []ValidationError{{Message: "expected object", Code: "invalid_type", Value: value}},
			}
		}
	case reflect.Struct:
		objMap = structToMap(v)
	default:
		return ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Message: "expected object", Code: "invalid_type", Value: value}},
		}
	}

	fields := s.getEffectiveFields()
	var errors []ValidationError
	validatedObj := make(map[string]interface{})

	// Validate known fields
	for fieldName, fieldSchema := range fields {
		fieldValue, exists := objMap[fieldName]
		if !exists {
			fieldValue = nil
		}

		result := fieldSchema.Validate(fieldValue)
		if !result.Valid {
			for _, err := range result.Errors {
				err.Field = fieldName
				errors = append(errors, err)
			}
		} else {
			if result.Value != nil {
				validatedObj[fieldName] = result.Value
			}
		}
	}

	// Handle unknown fields
	for fieldName, fieldValue := range objMap {
		if _, exists := fields[fieldName]; !exists {
			if s.strict {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: "unknown field",
					Code:    "unrecognized_keys",
					Value:   fieldValue,
				})
			} else if s.catchall != nil {
				result := s.catchall.Validate(fieldValue)
				if !result.Valid {
					for _, err := range result.Errors {
						err.Field = fieldName
						errors = append(errors, err)
					}
				} else {
					validatedObj[fieldName] = result.Value
				}
			} else if s.passthrough {
				validatedObj[fieldName] = fieldValue
			}
		}
	}

	if len(errors) > 0 {
		return ValidationResult{Valid: false, Errors: errors}
	}

	return ValidationResult{Valid: true, Value: validatedObj}
}

func convertMapToStringInterface(value interface{}) (map[string]interface{}, bool) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Map {
		return nil, false
	}

	result := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		result[keyStr] = v.MapIndex(key).Interface()
	}
	return result, true
}

func structToMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		fieldName := field.Name
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			if idx := strings.Index(tag, ","); idx != -1 {
				fieldName = tag[:idx]
			} else {
				fieldName = tag
			}
		}

		result[fieldName] = fieldValue.Interface()
	}

	return result
}