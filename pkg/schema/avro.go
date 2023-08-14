package schema

import (
	"context"
	"fmt"
	"github.com/linkedin/goavro"
	"io"
	"net/http"
)

type Client struct {
	host     string
	basePath string
}

func NewClient(host string) Client {
	return Client{
		host:     host,
		basePath: "apis/registry/v2",
	}
}

// GetSchemaByGlobalId retrieves the artifact with the given id.
func (c Client) GetSchemaByGlobalId(id string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/ids/globalIds/%s", c.host, c.basePath, id))
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

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
