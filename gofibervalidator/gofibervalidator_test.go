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
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// resetSingleton resets the package-level singleton so each test starts fresh.
func resetSingleton() {
	once = sync.Once{}
	validatorInstance = nil
}

//-----------------------------------------------------------------------------

func TestNew_DefaultConfig(t *testing.T) {
	resetSingleton()

	app := fiber.New()
	app.Use(New())
	app.Get("/", func(c *fiber.Ctx) error {
		v := c.Locals("validator")
		if v == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if _, ok := v.(*validator.Validate); !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestNew_CustomContextKey(t *testing.T) {
	resetSingleton()

	app := fiber.New()
	app.Use(New(Config{ContextKey: "myValidator"}))
	app.Get("/", func(c *fiber.Ctx) error {
		// Default key should be absent
		if c.Locals("validator") != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		// Custom key should be present
		v := c.Locals("myValidator")
		if v == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestNew_EmptyContextKeyFallsBackToDefault(t *testing.T) {
	resetSingleton()

	app := fiber.New()
	app.Use(New(Config{ContextKey: ""}))
	app.Get("/", func(c *fiber.Ctx) error {
		if c.Locals("validator") == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestNew_NextSkipsMiddleware(t *testing.T) {
	resetSingleton()

	app := fiber.New()
	app.Use(New(Config{
		Next: func(c *fiber.Ctx) bool { return true },
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		// Middleware was skipped, so the validator must not be in context
		if c.Locals("validator") != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestNew_CustomValidatorInstance(t *testing.T) {
	resetSingleton()

	customV := validator.New()
	app := fiber.New()
	app.Use(New(Config{Validator: customV}))
	app.Get("/", func(c *fiber.Ctx) error {
		v, ok := c.Locals("validator").(*validator.Validate)
		if !ok || v != customV {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestNew_CustomValidations(t *testing.T) {
	resetSingleton()

	type Payload struct {
		Name string `validate:"is_alice"`
	}

	app := fiber.New()
	app.Use(New(Config{
		CustomValidations: map[string]validator.Func{
			"is_alice": func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "alice"
			},
		},
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		v := GetValidator(c)
		if v == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		validPayload := Payload{Name: "alice"}
		if err := v.Struct(validPayload); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		invalidPayload := Payload{Name: "bob"}
		if err := v.Struct(invalidPayload); err == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

//-----------------------------------------------------------------------------

func TestGetValidator_DefaultKey(t *testing.T) {
	resetSingleton()

	app := fiber.New()
	app.Use(New())
	app.Get("/", func(c *fiber.Ctx) error {
		v := GetValidator(c)
		if v == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetValidator_CustomKey(t *testing.T) {
	resetSingleton()

	app := fiber.New()
	app.Use(New(Config{ContextKey: "myValidator"}))
	app.Get("/", func(c *fiber.Ctx) error {
		// Wrong key returns nil
		if GetValidator(c) != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		// Correct key returns the instance
		v := GetValidator(c, "myValidator")
		if v == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetValidator_NotInContext(t *testing.T) {
	// No middleware registered, so the validator is never stored.
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		if GetValidator(c) != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetValidator_ValidatesStruct(t *testing.T) {
	resetSingleton()

	type User struct {
		Email string `validate:"required,email"`
		Age   int    `validate:"gte=0,lte=130"`
	}

	app := fiber.New()
	app.Use(New())
	app.Get("/valid", func(c *fiber.Ctx) error {
		v := GetValidator(c)
		u := User{Email: "test@example.com", Age: 30}
		if err := v.Struct(u); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/invalid", func(c *fiber.Ctx) error {
		v := GetValidator(c)
		u := User{Email: "not-an-email", Age: 200}
		if err := v.Struct(u); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	for _, tc := range []struct {
		path string
		want int
	}{
		{"/valid", http.StatusOK},
		{"/invalid", http.StatusBadRequest},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != tc.want {
			t.Errorf("path %s: expected %d, got %d", tc.path, tc.want, resp.StatusCode)
		}
	}
}
