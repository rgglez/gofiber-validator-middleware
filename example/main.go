/*
Copyright 2025 Rodolfo González González

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	gofibervalidator "github.com/rgglez/gofiber-validator-middleware/gofibervalidator"
)

// User struct with validation tags
type User struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,gte=18,lte=100"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
}

// Product struct with validation tags
type Product struct {
	Name  string  `json:"name" validate:"required,min=2"`
	Price float64 `json:"price" validate:"required,gt=0"`
	SKU   string  `json:"sku" validate:"required,len=8,alphanum"`
}

// Custom validation function example
func isEvenNumber(fl validator.FieldLevel) bool {
	value := fl.Field().Int()
	return value%2 == 0
}

func main() {
	app := fiber.New()

	// Option 1: Use default configuration
	//app.Use(gofibervalidator.New())

	// Option 2: Use with custom configuration
	app.Use(gofibervalidator.New(gofibervalidator.Config{
		ContextKey: "validator",
		CustomValidations: map[string]validator.Func{
			"even": isEvenNumber,
		},
		Next: func(c *fiber.Ctx) bool {
			// Skip validation for specific routes
			return c.Path() == "/health"
		},
	}))

	// Routes
	app.Post("/users", createUser)
	app.Post("/products", createProduct)
	app.Get("/health", healthCheck)

	app.Listen(":3000")
}

// Handler using ValidateStruct helper
func createUser(c *fiber.Ctx) error {
	var user User

	// Parse request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// Validate using helper function
	if errors := gofibervalidator.ValidateStruct(c, user); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// Process valid user
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

// Handler using GetValidator directly
func createProduct(c *fiber.Ctx) error {
	var product Product

	// Parse request body
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// Get validator from context
	validate := gofibervalidator.GetValidator(c)
	if validate == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Validator not available",
		})
	}

	// Validate manually
	if err := validate.Struct(product); err != nil {
		var errors []gofibervalidator.ErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, gofibervalidator.ErrorResponse{
				Field:   err.Field(),
				Message: err.Tag(),
				Tag:     err.Tag(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// Process valid product
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": product,
	})
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}