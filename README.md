# gofiber-validator-middleware

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

## License

Copyright (c) 2026 Rodolfo González González.

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.