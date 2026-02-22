# gofiber-validator-middleware

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/gofiber-validator-middleware/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/gofiber-validator-middleware)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/gofiber-validator-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/gofiber-validator-middleware)](https://goreportcard.com/report/github.com/rgglez/gofiber-validator-middleware)
[![GitHub release](https://img.shields.io/github/release/rgglez/gofiber-validator-middleware.svg)](https://github.com/rgglez/gofiber-validator-middleware/releases/)
![GitHub stars](https://img.shields.io/github/stars/rgglez/gofiber-validator-middleware?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/gofiber-validator-middleware?style=social)

**gofiber-validator-middleware** is a [gofiber](https://docs.gofiber.io/category/-middleware/) [middleware](https://drstearns.github.io/tutorials/gomiddleware/) to generate a [go-playground/validator](go-playground/validator) [singleton](https://leangaurav.medium.com/golang-channels-vs-sync-once-for-one-time-execution-of-code-fafc81d2f54d) for use in validating DTOs in handlers
(or whatever you need to validate).

## Installation

```bash
go get github.com/rgglez/gofiber-validator-middleware
```

```go
import gofibervalidator "github.com/rgglez/gofiber-validator-middleware/gofibervalidator"
```

## Configuration

* ```Next``` defines a function to skip this middleware when returned true.
* ```Validator``` injected instance, if nil a new validator will be created.
* ```ContextKey``` is the key used to store the validator in context. Optional. Default: "**validator**"
* ```CustomValidations``` allows you to register custom validation functions. Optional. Default: nil

## Example

An example is included in the [example](example/) directory. To execute it:

1. Enter the example directory.
1. Run the example:
   ```bash
   go run .
   ```
1. You can try 3 use cases:
  * Skip validation (GET):
   [http://127.0.0.1:3000/health](http://127.0.0.1:3000/health)
  * Validate user "creation" (POST):
   [http://127.0.0.1:3000/users](http://127.0.0.1:3000/users)
  * Validate product "creation" (POST):
   [http://127.0.0.1:3000/products](http://127.0.0.1:3000/products)

You can use [resting](https://addons.mozilla.org/en-US/firefox/addon/resting/) or the plugin of your choice, or tools like [curl](https://curl.se/), to try the endpoints.

## Dependencies

* [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2)
* [go-playground/validator](https://github.com/go-playground/validator)

## Tests

Run the tests with:

```bash
go test ./gofibervalidator/... -v
```

### Test cases

| Test | What it verifies |
|---|---|
| `TestNew_DefaultConfig` | Middleware stores a `*validator.Validate` under the `"validator"` key by default |
| `TestNew_CustomContextKey` | Custom `ContextKey` is used; default key is absent |
| `TestNew_EmptyContextKeyFallsBackToDefault` | Empty `ContextKey` falls back to `"validator"` |
| `TestNew_NextSkipsMiddleware` | When `Next` returns `true`, middleware is skipped and no validator is stored |
| `TestNew_CustomValidatorInstance` | A pre-built `*validator.Validate` provided via config is used as-is |
| `TestNew_CustomValidations` | Custom tag functions registered via `CustomValidations` work correctly |
| `TestGetValidator_DefaultKey` | `GetValidator` retrieves the instance using the default key |
| `TestGetValidator_CustomKey` | `GetValidator` retrieves via a custom key; wrong key returns `nil` |
| `TestGetValidator_NotInContext` | `GetValidator` returns `nil` when no middleware is registered |
| `TestGetValidator_ValidatesStruct` | The retrieved validator correctly validates struct fields |

## License

Copyright (c) 2026 Rodolfo González González.

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.