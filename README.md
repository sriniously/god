# God - A Zod-like Validation Library for Go

God is a comprehensive validation library for Go that provides type-safe schema validation with a fluent API, inspired by Zod from the TypeScript ecosystem.

## Features

- **Type-safe validation** with comprehensive error messages
- **Fluent API** for building complex validation schemas
- **Primitive types**: String, Number, Boolean, Date validation
- **Complex types**: Object, Array, Tuple validation
- **Advanced features**: Union types, Discriminated unions, Enums, Literals
- **Flexible constraints**: Min/Max, Regex, Email, URL validation
- **Transformations**: String trimming, case conversion
- **Optional and default values** support
- **Nested validation** with detailed error paths
- **Lazy evaluation** for recursive schemas

## Installation

```bash
go get github.com/sriniously/god
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/sriniously/god"
)

func main() {
    // Define a user schema
    userSchema := god.Object(map[string]god.Schema{
        "name":  god.String().Min(2).Max(50),
        "email": god.String().Email(),
        "age":   god.Int().Min(0).Max(120),
    })

    // Validate data
    user := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   30,
    }

    result := userSchema.Validate(user)
    if result.Valid {
        fmt.Println("Validation passed!")
        // Access validated data
        validated := result.Value.(map[string]interface{})
        fmt.Printf("Name: %s\n", validated["name"])
    } else {
        fmt.Printf("Validation failed: %v\n", result.Error())
    }
}
```

## Core Types

### String Validation

```go
schema := god.String().Min(5).Max(100).Email()
schema = god.String().Regex(`^[a-zA-Z0-9]+$`)
schema = god.String().URL()
schema = god.String().UUID()

// Transformations
schema = god.String().Trim().ToLower()
schema = god.String().ToUpper()
```

### Number Validation

```go
schema := god.Number().Min(0).Max(100)
schema = god.Int().Positive()
schema = god.Number().Negative()
schema = god.Number().NonNegative()
schema = god.Number().MultipleOf(5)
```

### Boolean Validation

```go
schema := god.Boolean()
schema = god.Bool() // Alias
```

### Date Validation

```go
schema := god.Date()
schema = god.Date().Min(time.Now())
schema = god.Date().Max(time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC))
```

## Complex Types

### Object Validation

```go
userSchema := god.Object(map[string]god.Schema{
    "id":       god.Int().Positive(),
    "name":     god.String().Min(1).Max(100),
    "email":    god.String().Email(),
    "age":      god.Int().Min(0).Max(120).Optional(),
    "isActive": god.Boolean().Default(true),
})

// Object operations
schema = userSchema.Partial()        // Make all fields optional
schema = userSchema.RequiredFields("name", "email")  // Require specific fields
schema = userSchema.Pick("name", "email")            // Pick specific fields
schema = userSchema.Omit("id")                       // Omit specific fields
schema = userSchema.Strict()                         // Disallow unknown fields
schema = userSchema.Passthrough()                    // Allow unknown fields
```

### Array Validation

```go
schema := god.Array(god.String()).Min(1).Max(10)
schema = god.Array(god.Int()).Nonempty()
```

### Tuple Validation

```go
// Fixed-length tuple
coordinateSchema := god.Tuple(god.Number(), god.Number(), god.Number())

// Tuple with rest elements
csvSchema := god.Tuple(god.String(), god.String()).Rest(god.Union(god.String(), god.Number()))
```

## Advanced Features

### Union Types

```go
// Simple union
schema := god.Union(god.String(), god.Number())

// Discriminated union
shapeSchema := god.DiscriminatedUnion("type", map[string]god.Schema{
    "circle": god.Object(map[string]god.Schema{
        "type":   god.Literal("circle"),
        "radius": god.Number().Positive(),
    }),
    "rectangle": god.Object(map[string]god.Schema{
        "type":   god.Literal("rectangle"),
        "width":  god.Number().Positive(),
        "height": god.Number().Positive(),
    }),
})
```

### Enums and Literals

```go
// Enum validation
roleSchema := god.Enum("user", "admin", "moderator")

// Literal validation
typeSchema := god.Literal("success")
```

### Nullable Types

```go
schema := god.Nullable(god.String())
```

### Utility Types

```go
schema := god.Any()        // Accepts any value
schema = god.Unknown()     // Accepts any value (alias for Any)
schema = god.Void()        // Always returns nil
schema = god.Never()       // Always fails validation
```

### Lazy Evaluation

```go
// For recursive schemas
var nodeSchema god.Schema
nodeSchema = god.Lazy(func() god.Schema {
    return god.Object(map[string]god.Schema{
        "value":    god.String(),
        "children": god.Array(nodeSchema).Optional(),
    })
})
```

## Optional and Default Values

```go
schema := god.Object(map[string]god.Schema{
    "name":     god.String(),
    "email":    god.String().Email().Optional(),
    "isActive": god.Boolean().Default(true),
    "role":     god.Enum("user", "admin").Default("user"),
})
```

## Error Handling

```go
result := schema.Validate(data)
if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Field: %s, Error: %s, Code: %s\n", err.Field, err.Message, err.Code)
    }
    
    // Or get formatted error
    fmt.Printf("Error: %v\n", result.Error())
}
```

## Transformations

God supports data transformations during validation:

```go
userSchema := god.Object(map[string]god.Schema{
    "name":     god.String().Trim(),
    "email":    god.String().ToLower().Email(),
    "username": god.String().ToLower().Regex(`^[a-z0-9_]+$`),
})

// Input data with extra whitespace and mixed case
input := map[string]interface{}{
    "name":     "  John Doe  ",
    "email":    "  JOHN@EXAMPLE.COM  ",
    "username": "JOHN_DOE",
}

result := userSchema.Validate(input)
if result.Valid {
    validated := result.Value.(map[string]interface{})
    // validated["name"] = "John Doe"
    // validated["email"] = "john@example.com"
    // validated["username"] = "john_doe"
}
```

## Performance Considerations

- Schemas are reusable and thread-safe
- Compile schemas once and reuse them
- Use `Lazy()` for recursive schemas to avoid infinite recursion
- Consider using `Strict()` on objects when you don't need unknown fields

## Examples

See `example_test.go` for comprehensive examples including:
- User registration validation
- API request validation
- Complex nested object validation
- Discriminated union validation
- Array and tuple validation
- Data transformation examples

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.