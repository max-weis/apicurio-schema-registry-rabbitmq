package validator

import "context"

type Validator interface {
	Validate(ctx context.Context, obj map[string]any) error
}
