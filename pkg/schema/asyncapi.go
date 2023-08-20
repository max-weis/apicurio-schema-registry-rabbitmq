package schema

import (
	"context"
	"fmt"
	"github.com/max-weis/go-asyncapi-validator/pkg/validator"
)

type AsyncAPIValidator struct {
	spec     string
	jsonPath string
}

func NewAsyncAPIValidator(spec, jsonPath string) *AsyncAPIValidator {
	return &AsyncAPIValidator{spec: spec, jsonPath: jsonPath}
}

func (v *AsyncAPIValidator) Validate(ctx context.Context, obj map[string]interface{}) error {
	spec, err := validator.LoadAsyncAPISpec(v.spec)
	if err != nil {
		return fmt.Errorf("failed to load AsyncAPI spec: %w", err)
	}

	schema, err := validator.ExtractSchemaWithJSONPath(spec, v.jsonPath)
	if err != nil {
		return fmt.Errorf("failed to extract schema: %w", err)
	}

	if err = validator.ValidateJSONAgainstSchema(obj, schema); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
