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
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

//-----------------------------------------------------------------------------

// Config defines the config for middleware
type Config struct {
	// Next defines a function to skip this middleware when returned true
	Next func(c *fiber.Ctx) bool

	// Validator instance, if nil a new validator will be created
	Validator *validator.Validate

	// ContextKey is the key used to store the validator in context
	// Optional. Default: "validator"
	ContextKey string

	// CustomValidations allows you to register custom validation functions
	// Optional. Default: nil
	CustomValidations map[string]validator.Func
}

//-----------------------------------------------------------------------------

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:       nil,
	Validator:  nil,
	ContextKey: "validator",
}

//-----------------------------------------------------------------------------

// Singleton
var (
	validatorInstance *validator.Validate
	once              sync.Once
)

//-----------------------------------------------------------------------------

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.ContextKey == "" {
			cfg.ContextKey = ConfigDefault.ContextKey
		}
	}

	// Initialize validator instance (singleton)
	once.Do(func() {
		if cfg.Validator != nil {
			validatorInstance = cfg.Validator
		} else {
			validatorInstance = validator.New()
		}

		// Register custom validations if provided
		if cfg.CustomValidations != nil {
			for tag, fn := range cfg.CustomValidations {
				validatorInstance.RegisterValidation(tag, fn)
			}
		}
	})

	// Return middleware handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Store validator in context
		c.Locals(cfg.ContextKey, validatorInstance)

		// Continue stack
		return c.Next()
	}
}

//-----------------------------------------------------------------------------

// GetValidator retrieves the validator from context
func GetValidator(c *fiber.Ctx, key ...string) *validator.Validate {
	contextKey := "validator"
	if len(key) > 0 {
		contextKey = key[0]
	}

	if v := c.Locals(contextKey); v != nil {
		return v.(*validator.Validate)
	}
	return nil
}
