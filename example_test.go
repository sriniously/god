package god

import (
	"fmt"
	"strings"
	"time"
)

// Example demonstrates a comprehensive user validation schema
func Example() {
	// Define a user schema with various validation rules
	userSchema := Object(map[string]Schema{
		"id":       Int().Positive(),
		"username": String().Min(3).Max(50).Regex(`^[a-zA-Z0-9_]+$`),
		"email":    String().Email(),
		"age":      Int().Min(13).Max(120).Optional(),
		"bio":      String().Max(500).Optional(),
		"website":  String().URL().Optional(),
		"isActive": Boolean().Default(true),
		"tags":     Array(String()).Min(1).Max(10),
		"role":     Enum("user", "admin", "moderator"),
		"profile": Object(map[string]Schema{
			"firstName": String().Min(1).Max(50),
			"lastName":  String().Min(1).Max(50),
			"avatar":    String().URL().Optional(),
			"birthDate": Date().Max(time.Now()),
		}),
		"settings": Object(map[string]Schema{
			"notifications": Boolean().Default(true),
			"theme":         Enum("light", "dark").Default("light"),
			"language":      String().Default("en"),
		}),
	})

	// Valid user data
	validUser := map[string]interface{}{
		"id":       123,
		"username": "john_doe",
		"email":    "john@example.com",
		"age":      30,
		"bio":      "Software developer",
		"isActive": true,
		"tags":     []interface{}{"developer", "golang"},
		"role":     "user",
		"profile": map[string]interface{}{
			"firstName": "John",
			"lastName":  "Doe",
			"avatar":    "https://example.com/avatar.jpg",
			"birthDate": time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"settings": map[string]interface{}{
			"notifications": true,
			"theme":         "dark",
			"language":      "en",
		},
	}

	result := userSchema.Validate(validUser)
	if result.Valid {
		fmt.Println("User validation passed!")
		// Access validated data
		validated := result.Value.(map[string]interface{})
		fmt.Printf("User ID: %v\n", validated["id"])
		fmt.Printf("Username: %v\n", validated["username"])
	} else {
		fmt.Printf("Validation failed: %v\n", result.Error())
	}

	// Output:
	// User validation passed!
	// User ID: 123
	// Username: john_doe
}

// Example_api demonstrates request validation
func Example_api() {
	// Define API request schema
	requestSchema := Object(map[string]Schema{
		"method": Enum("GET", "POST", "PUT", "DELETE"),
		"path":   String().Regex(`^/[a-zA-Z0-9/_-]*$`),
		"query":  Object(map[string]Schema{}).Passthrough().Optional(),
		"body":   Union(Object(map[string]Schema{}).Passthrough(), String(), Array(Any())).Optional(),
		"headers": Object(map[string]Schema{
			"Content-Type": String().Optional(),
			"Authorization": String().Regex(`^Bearer [a-zA-Z0-9._-]+$`).Optional(),
		}).Passthrough(),
	})

	// Valid API request
	apiRequest := map[string]interface{}{
		"method": "POST",
		"path":   "/api/users",
		"query": map[string]interface{}{
			"include": "profile",
		},
		"body": map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
		},
		"headers": map[string]interface{}{
			"Content-Type":  "application/json",
			"Authorization": "Bearer abc123.def456.ghi789",
			"User-Agent":    "MyApp/1.0",
		},
	}

	result := requestSchema.Validate(apiRequest)
	if result.Valid {
		fmt.Println("API request validation passed!")
	} else {
		fmt.Printf("Validation failed: %v\n", result.Error())
	}

	// Output:
	// API request validation passed!
}

// Example_union demonstrates discriminated union validation
func Example_union() {
	// Define shapes with discriminated unions
	shapeSchema := DiscriminatedUnion("type", map[string]Schema{
		"circle": Object(map[string]Schema{
			"type":   Literal("circle"),
			"radius": Number().Positive(),
		}),
		"rectangle": Object(map[string]Schema{
			"type":   Literal("rectangle"),
			"width":  Number().Positive(),
			"height": Number().Positive(),
		}),
		"triangle": Object(map[string]Schema{
			"type": Literal("triangle"),
			"base": Number().Positive(),
			"height": Number().Positive(),
		}),
	})

	// Valid circle
	circle := map[string]interface{}{
		"type":   "circle",
		"radius": 5.5,
	}

	result := shapeSchema.Validate(circle)
	if result.Valid {
		fmt.Println("Circle validation passed!")
	}

	// Valid rectangle
	rectangle := map[string]interface{}{
		"type":   "rectangle",
		"width":  10,
		"height": 20,
	}

	result = shapeSchema.Validate(rectangle)
	if result.Valid {
		fmt.Println("Rectangle validation passed!")
	}

	// Invalid shape (missing required field)
	invalid := map[string]interface{}{
		"type": "circle",
		// missing radius
	}

	result = shapeSchema.Validate(invalid)
	if !result.Valid {
		fmt.Printf("Invalid shape error: %v\n", result.Error())
	}

	// Output:
	// Circle validation passed!
	// Rectangle validation passed!
	// Invalid shape error: validation failed: radius: field is required
}

// Example_transform demonstrates string transformations
func Example_transform() {
	// Schema with transformations
	userSchema := Object(map[string]Schema{
		"name":     String().Trim().Min(1),
		"email":    String().Transform(func(s string) string {
			return strings.ToLower(strings.TrimSpace(s))
		}).Email(),
		"username": String().ToLower().Regex(`^[a-z0-9_]+$`),
		"bio":      String().Trim().Max(100).Optional(),
	})

	// Data with whitespace and mixed case
	userData := map[string]interface{}{
		"name":     "  John Doe  ",
		"email":    "  JOHN@EXAMPLE.COM  ",
		"username": "JOHN_DOE",
		"bio":      "  Software developer  ",
	}

	result := userSchema.Validate(userData)
	if result.Valid {
		validated := result.Value.(map[string]interface{})
		fmt.Printf("Transformed name: '%v'\n", validated["name"])
		fmt.Printf("Transformed email: '%v'\n", validated["email"])
		fmt.Printf("Transformed username: '%v'\n", validated["username"])
		fmt.Printf("Transformed bio: '%v'\n", validated["bio"])
	} else {
		fmt.Printf("Validation failed: %v\n", result.Error())
	}

	// Output:
	// Transformed name: 'John Doe'
	// Transformed email: 'john@example.com'
	// Transformed username: 'john_doe'
	// Transformed bio: 'Software developer'
}

// Example_array demonstrates array validation
func Example_array() {
	// Schema for a todo list
	todoSchema := Object(map[string]Schema{
		"title": String().Min(1).Max(100),
		"items": Array(Object(map[string]Schema{
			"id":        Int().Positive(),
			"text":      String().Min(1).Max(500),
			"completed": Boolean().Default(false),
			"priority":  Enum("low", "medium", "high").Default("medium"),
			"tags":      Array(String()).Max(5).Optional(),
		})).Min(1).Max(100),
		"created": Date(),
	})

	todoList := map[string]interface{}{
		"title": "My Todo List",
		"items": []interface{}{
			map[string]interface{}{
				"id":        1,
				"text":      "Complete Go project",
				"completed": false,
				"priority":  "high",
				"tags":      []interface{}{"work", "golang"},
			},
			map[string]interface{}{
				"id":        2,
				"text":      "Buy groceries",
				"completed": true,
				"priority":  "medium",
			},
		},
		"created": time.Now(),
	}

	result := todoSchema.Validate(todoList)
	if result.Valid {
		fmt.Println("Todo list validation passed!")
		validated := result.Value.(map[string]interface{})
		items := validated["items"].([]interface{})
		fmt.Printf("Number of items: %d\n", len(items))
	}

	// Output:
	// Todo list validation passed!
	// Number of items: 2
}

// Example_tuple demonstrates tuple validation
func Example_tuple() {
	// Define a coordinate tuple (x, y, z)
	coordinateSchema := Tuple(Number(), Number(), Number())

	coordinate := []interface{}{10.5, 20.3, 5.0}

	result := coordinateSchema.Validate(coordinate)
	if result.Valid {
		validated := result.Value.([]interface{})
		fmt.Printf("Coordinate: x=%v, y=%v, z=%v\n", validated[0], validated[1], validated[2])
	}

	// Tuple with rest elements
	csvRowSchema := Tuple(String(), String()).Rest(Union(String(), Number()))

	csvRow := []interface{}{"John", "Doe", 30, "Engineer", "Active"}

	result = csvRowSchema.Validate(csvRow)
	if result.Valid {
		fmt.Println("CSV row validation passed!")
		validated := result.Value.([]interface{})
		fmt.Printf("Name: %v %v\n", validated[0], validated[1])
	}

	// Output:
	// Coordinate: x=10.5, y=20.3, z=5
	// CSV row validation passed!
	// Name: John Doe
}

// Example_nested demonstrates complex nested validation
func Example_nested() {
	// Define a complex nested schema for a blog post
	blogPostSchema := Object(map[string]Schema{
		"title":   String().Min(1).Max(200),
		"content": String().Min(1),
		"author": Object(map[string]Schema{
			"id":       Int().Positive(),
			"name":     String().Min(1).Max(100),
			"email":    String().Email(),
			"profile":  Object(map[string]Schema{
				"bio":     String().Max(500).Optional(),
				"avatar":  String().URL().Optional(),
				"social":  Object(map[string]Schema{
					"twitter":  String().Regex(`^@[a-zA-Z0-9_]+$`).Optional(),
					"linkedin": String().URL().Optional(),
					"github":   String().Regex(`^[a-zA-Z0-9_-]+$`).Optional(),
				}).Optional(),
			}),
		}),
		"tags":      Array(String().Min(1).Max(50)).Max(10),
		"published": Boolean().Default(false),
		"metadata": Object(map[string]Schema{
			"created":    Date(),
			"updated":    Date().Optional(),
			"views":      Int().NonNegative().Default(0),
			"likes":      Int().NonNegative().Default(0),
			"comments":   Array(Object(map[string]Schema{
				"id":      Int().Positive(),
				"author":  String().Min(1).Max(100),
				"content": String().Min(1).Max(1000),
				"created": Date(),
			})).Optional(),
		}),
	})

	blogPost := map[string]interface{}{
		"title":   "Introduction to Go Validation",
		"content": "This is a comprehensive guide to validation in Go...",
		"author": map[string]interface{}{
			"id":    1,
			"name":  "John Doe",
			"email": "john@example.com",
			"profile": map[string]interface{}{
				"bio":    "Software developer passionate about Go",
				"avatar": "https://example.com/john.jpg",
				"social": map[string]interface{}{
					"twitter": "@johndoe",
					"github":  "johndoe",
				},
			},
		},
		"tags":      []interface{}{"go", "validation", "tutorial"},
		"published": true,
		"metadata": map[string]interface{}{
			"created": time.Now(),
			"views":   150,
			"likes":   25,
			"comments": []interface{}{
				map[string]interface{}{
					"id":      1,
					"author":  "Jane Smith",
					"content": "Great tutorial!",
					"created": time.Now(),
				},
			},
		},
	}

	result := blogPostSchema.Validate(blogPost)
	if result.Valid {
		fmt.Println("Blog post validation passed!")
		validated := result.Value.(map[string]interface{})
		fmt.Printf("Title: %v\n", validated["title"])
		
		author := validated["author"].(map[string]interface{})
		fmt.Printf("Author: %v\n", author["name"])
		
		metadata := validated["metadata"].(map[string]interface{})
		fmt.Printf("Views: %v\n", metadata["views"])
	} else {
		fmt.Printf("Validation failed: %v\n", result.Error())
	}

	// Output:
	// Blog post validation passed!
	// Title: Introduction to Go Validation
	// Author: John Doe
	// Views: 150
}