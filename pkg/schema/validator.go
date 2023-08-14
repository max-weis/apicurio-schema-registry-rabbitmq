package schema

import "context"

type Validator interface {
	Validate(ctx context.Context, obj map[string]any) (bool, error)
}
