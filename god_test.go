package god

import (
	"testing"
	"time"
)

func TestStringSchema(t *testing.T) {
	schema := String()

	// Valid string
	result := schema.Validate("hello")
	if !result.Valid {
		t.Errorf("Expected valid result for string 'hello', got invalid")
	}

	// Invalid type
	result = schema.Validate(123)
	if result.Valid {
		t.Errorf("Expected invalid result for number, got valid")
	}

	// Test min length
	schema = String().Min(5)
	result = schema.Validate("hi")
	if result.Valid {
		t.Errorf("Expected invalid result for short string, got valid")
	}

	// Test max length
	schema = String().Max(3)
	result = schema.Validate("hello")
	if result.Valid {
		t.Errorf("Expected invalid result for long string, got valid")
	}

	// Test email
	schema = String().Email()
	result = schema.Validate("test@example.com")
	if !result.Valid {
		t.Errorf("Expected valid result for email, got invalid")
	}

	result = schema.Validate("invalid-email")
	if result.Valid {
		t.Errorf("Expected invalid result for invalid email, got valid")
	}

	// Test optional
	optionalSchema := String().Optional()
	result = optionalSchema.Validate(nil)
	if !result.Valid {
		t.Errorf("Expected valid result for nil on optional field, got invalid: %v", result.Errors)
	}
}

func TestNumberSchema(t *testing.T) {
	schema := Number()

	// Valid number
	result := schema.Validate(42.5)
	if !result.Valid {
		t.Errorf("Expected valid result for number 42.5, got invalid")
	}

	// Invalid type
	result = schema.Validate("hello")
	if result.Valid {
		t.Errorf("Expected invalid result for string, got valid")
	}

	// Test min
	schema = Number().Min(10)
	result = schema.Validate(5)
	if result.Valid {
		t.Errorf("Expected invalid result for number below min, got valid")
	}

	// Test max
	schema = Number().Max(100)
	result = schema.Validate(150)
	if result.Valid {
		t.Errorf("Expected invalid result for number above max, got valid")
	}

	// Test positive
	schema = Number().Positive()
	result = schema.Validate(-5)
	if result.Valid {
		t.Errorf("Expected invalid result for negative number, got valid")
	}

	// Test integer
	schema = Int()
	result = schema.Validate(42)
	if !result.Valid {
		t.Errorf("Expected valid result for integer, got invalid")
	}

	result = schema.Validate(42.5)
	if result.Valid {
		t.Errorf("Expected invalid result for float when expecting integer, got valid")
	}
}

func TestBooleanSchema(t *testing.T) {
	schema := Boolean()

	// Valid boolean
	result := schema.Validate(true)
	if !result.Valid {
		t.Errorf("Expected valid result for boolean true, got invalid")
	}

	// Invalid type
	result = schema.Validate("hello")
	if result.Valid {
		t.Errorf("Expected invalid result for string, got valid")
	}

	// Test string conversion
	result = schema.Validate("true")
	if !result.Valid {
		t.Errorf("Expected valid result for string 'true', got invalid")
	}

	result = schema.Validate("false")
	if !result.Valid {
		t.Errorf("Expected valid result for string 'false', got invalid")
	}
}

func TestObjectSchema(t *testing.T) {
	schema := Object(map[string]Schema{
		"name": String(),
		"age":  Number(),
	})

	// Valid object
	obj := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	result := schema.Validate(obj)
	if !result.Valid {
		t.Errorf("Expected valid result for valid object, got invalid: %v", result.Errors)
	}

	// Missing required field
	obj = map[string]interface{}{
		"name": "John",
	}
	result = schema.Validate(obj)
	if result.Valid {
		t.Errorf("Expected invalid result for object missing required field, got valid")
	}

	// Invalid field type
	obj = map[string]interface{}{
		"name": "John",
		"age":  "thirty",
	}
	result = schema.Validate(obj)
	if result.Valid {
		t.Errorf("Expected invalid result for object with invalid field type, got valid")
	}

	// Test optional fields
	schema = Object(map[string]Schema{
		"name":  String(),
		"email": String().Optional(),
	})

	obj = map[string]interface{}{
		"name": "John",
	}
	result = schema.Validate(obj)
	if !result.Valid {
		t.Errorf("Expected valid result for object with optional field missing, got invalid: %v", result.Errors)
	}
}

func TestArraySchema(t *testing.T) {
	schema := Array(String())

	// Valid array
	arr := []interface{}{"hello", "world"}
	result := schema.Validate(arr)
	if !result.Valid {
		t.Errorf("Expected valid result for valid array, got invalid: %v", result.Errors)
	}

	// Invalid element type
	arr = []interface{}{"hello", 123}
	result = schema.Validate(arr)
	if result.Valid {
		t.Errorf("Expected invalid result for array with invalid element type, got valid")
	}

	// Test min length
	schema = Array(String()).Min(3)
	arr = []interface{}{"hello"}
	result = schema.Validate(arr)
	if result.Valid {
		t.Errorf("Expected invalid result for array below min length, got valid")
	}

	// Test max length
	schema = Array(String()).Max(2)
	arr = []interface{}{"hello", "world", "test"}
	result = schema.Validate(arr)
	if result.Valid {
		t.Errorf("Expected invalid result for array above max length, got valid")
	}
}

func TestUnionSchema(t *testing.T) {
	schema := Union(String(), Number())

	// Valid string
	result := schema.Validate("hello")
	if !result.Valid {
		t.Errorf("Expected valid result for string in union, got invalid")
	}

	// Valid number
	result = schema.Validate(42)
	if !result.Valid {
		t.Errorf("Expected valid result for number in union, got invalid")
	}

	// Invalid type
	result = schema.Validate(true)
	if result.Valid {
		t.Errorf("Expected invalid result for boolean in string/number union, got valid")
	}
}

func TestLiteralSchema(t *testing.T) {
	schema := Literal("hello")

	// Valid literal
	result := schema.Validate("hello")
	if !result.Valid {
		t.Errorf("Expected valid result for matching literal, got invalid")
	}

	// Invalid literal
	result = schema.Validate("world")
	if result.Valid {
		t.Errorf("Expected invalid result for non-matching literal, got valid")
	}
}

func TestEnumSchema(t *testing.T) {
	schema := Enum("red", "green", "blue")

	// Valid enum value
	result := schema.Validate("red")
	if !result.Valid {
		t.Errorf("Expected valid result for valid enum value, got invalid")
	}

	// Invalid enum value
	result = schema.Validate("yellow")
	if result.Valid {
		t.Errorf("Expected invalid result for invalid enum value, got valid")
	}
}

func TestNullableSchema(t *testing.T) {
	schema := Nullable(String())

	// Valid string
	result := schema.Validate("hello")
	if !result.Valid {
		t.Errorf("Expected valid result for string in nullable, got invalid")
	}

	// Valid null
	result = schema.Validate(nil)
	if !result.Valid {
		t.Errorf("Expected valid result for nil in nullable, got invalid")
	}

	// Invalid type
	result = schema.Validate(123)
	if result.Valid {
		t.Errorf("Expected invalid result for number in nullable string, got valid")
	}
}

func TestDateSchema(t *testing.T) {
	schema := Date()

	// Valid time.Time
	now := time.Now()
	result := schema.Validate(now)
	if !result.Valid {
		t.Errorf("Expected valid result for time.Time, got invalid")
	}

	// Valid RFC3339 string
	result = schema.Validate("2023-01-01T00:00:00Z")
	if !result.Valid {
		t.Errorf("Expected valid result for RFC3339 string, got invalid")
	}

	// Valid date string
	result = schema.Validate("2023-01-01")
	if !result.Valid {
		t.Errorf("Expected valid result for date string, got invalid")
	}

	// Invalid date string
	result = schema.Validate("invalid-date")
	if result.Valid {
		t.Errorf("Expected invalid result for invalid date string, got valid")
	}

	// Test min date
	minDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	schema = Date().Min(minDate)
	testDate := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	result = schema.Validate(testDate)
	if result.Valid {
		t.Errorf("Expected invalid result for date before min, got valid")
	}
}

func TestComplexObjectValidation(t *testing.T) {
	// Define a complex nested schema
	userSchema := Object(map[string]Schema{
		"id":    Number().Positive(),
		"name":  String().Min(2),
		"email": String().Email(),
		"age":   Number().Min(0).Max(150).Optional(),
		"address": Object(map[string]Schema{
			"street": String(),
			"city":   String(),
			"zip":    String().Regex(`^\d{5}$`),
		}),
		"hobbies": Array(String()).Min(1),
		"status":  Enum("active", "inactive", "suspended"),
	})

	// Valid user
	user := map[string]interface{}{
		"id":    123,
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
			"zip":    "10001",
		},
		"hobbies": []interface{}{"reading", "swimming"},
		"status":  "active",
	}

	result := userSchema.Validate(user)
	if !result.Valid {
		t.Errorf("Expected valid result for complex valid user, got invalid: %v", result.Errors)
	}

	// Invalid user - bad email
	user["email"] = "invalid-email"
	result = userSchema.Validate(user)
	if result.Valid {
		t.Errorf("Expected invalid result for user with bad email, got valid")
	}

	// Reset email and test bad zip
	user["email"] = "john@example.com"
	address := user["address"].(map[string]interface{})
	address["zip"] = "invalid-zip"
	result = userSchema.Validate(user)
	if result.Valid {
		t.Errorf("Expected invalid result for user with bad zip, got valid")
	}
}

func TestDefaultValues(t *testing.T) {
	schema := Object(map[string]Schema{
		"name":   String(),
		"active": Boolean().Default(true),
		"count":  Number().Default(0),
	})

	obj := map[string]interface{}{
		"name": "Test",
	}

	result := schema.Validate(obj)
	if !result.Valid {
		t.Errorf("Expected valid result with defaults, got invalid: %v", result.Errors)
	}

	if result.Value != nil {
		validated := result.Value.(map[string]interface{})
		if validated["active"] != true {
			t.Errorf("Expected default value true for active, got %v", validated["active"])
		}

		if validated["count"] != 0.0 {
			t.Errorf("Expected default value 0 for count, got %v (type: %T)", validated["count"], validated["count"])
		}
	}
}

func TestTupleSchema(t *testing.T) {
	schema := Tuple(String(), Number(), Boolean())

	// Valid tuple
	tuple := []interface{}{"hello", 42, true}
	result := schema.Validate(tuple)
	if !result.Valid {
		t.Errorf("Expected valid result for valid tuple, got invalid: %v", result.Errors)
	}

	// Invalid tuple - wrong length
	tuple = []interface{}{"hello", 42}
	result = schema.Validate(tuple)
	if result.Valid {
		t.Errorf("Expected invalid result for tuple with wrong length, got valid")
	}

	// Invalid tuple - wrong type
	tuple = []interface{}{"hello", "world", true}
	result = schema.Validate(tuple)
	if result.Valid {
		t.Errorf("Expected invalid result for tuple with wrong type, got valid")
	}

	// Test tuple with rest
	schema = Tuple(String(), Number()).Rest(Boolean())
	tuple = []interface{}{"hello", 42, true, false, true}
	result = schema.Validate(tuple)
	if !result.Valid {
		t.Errorf("Expected valid result for tuple with rest, got invalid: %v", result.Errors)
	}
}

func TestDiscriminatedUnion(t *testing.T) {
	schema := DiscriminatedUnion("type", map[string]Schema{
		"user": Object(map[string]Schema{
			"type": Literal("user"),
			"name": String(),
		}),
		"admin": Object(map[string]Schema{
			"type":        Literal("admin"),
			"name":        String(),
			"permissions": Array(String()),
		}),
	})

	// Valid user
	user := map[string]interface{}{
		"type": "user",
		"name": "John",
	}
	result := schema.Validate(user)
	if !result.Valid {
		t.Errorf("Expected valid result for discriminated union user, got invalid: %v", result.Errors)
	}

	// Valid admin
	admin := map[string]interface{}{
		"type":        "admin",
		"name":        "Jane",
		"permissions": []interface{}{"read", "write"},
	}
	result = schema.Validate(admin)
	if !result.Valid {
		t.Errorf("Expected valid result for discriminated union admin, got invalid: %v", result.Errors)
	}

	// Invalid - missing discriminant
	invalid := map[string]interface{}{
		"name": "John",
	}
	result = schema.Validate(invalid)
	if result.Valid {
		t.Errorf("Expected invalid result for missing discriminant, got valid")
	}

	// Invalid - unknown discriminant
	invalid = map[string]interface{}{
		"type": "unknown",
		"name": "John",
	}
	result = schema.Validate(invalid)
	if result.Valid {
		t.Errorf("Expected invalid result for unknown discriminant, got valid")
	}
}