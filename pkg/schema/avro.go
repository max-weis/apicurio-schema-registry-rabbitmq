package schema

import (
	"context"
	"fmt"
	"github.com/linkedin/goavro"
)

type AvroValidator struct {
	schema string
}

func NewAvroValidator(schema string) Validator {
	return AvroValidator{
		schema: schema,
	}
}

func (av AvroValidator) Validate(ctx context.Context, obj map[string]any) error {
	// Create a Codec for the Avro schema
	codec, err := goavro.NewCodec(av.schema)
	if err != nil {
		return fmt.Errorf("failed to create codec: %w", err)
	}

	// Convert the map into Avro binary format
	binary, err := codec.BinaryFromNative(nil, obj)
	if err != nil {
		return fmt.Errorf("failed to encode to binary: %w", err)
	}

	// Convert back from binary to Go native data structure
	_, _, err = codec.NativeFromBinary(binary)
	if err != nil {
		return fmt.Errorf("failed to decode from binary: %w", err)
	}

	return nil
}
