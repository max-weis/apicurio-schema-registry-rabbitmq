package schema

import (
	"fmt"
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
