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

package gofibervalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Helper functions for common validation patterns

// ErrorResponse represents a validation error
type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
}

// ValidateStruct validates a struct using the validator from context
func ValidateStruct(c *fiber.Ctx, s interface{}) []ErrorResponse {
	v := GetValidator(c)
	if v == nil {
		return nil
	}

	var errors []ErrorResponse
	err := v.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Message = getErrorMessage(err)
			errors = append(errors, element)
		}
	}
	return errors
}

// getErrorMessage returns a user-friendly error message
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (min: " + fe.Param() + ")"
	case "max":
		return "Value is too long (max: " + fe.Param() + ")"
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	case "len":
		return "Value must have exactly " + fe.Param() + " characters"
	case "url":
		return "Invalid URL format"
	case "uuid":
		return "Invalid UUID format"
	default:
		return "Invalid value for " + fe.Field()
	}
}
