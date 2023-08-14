package registry

import (
	"context"
	"github.com/Apicurio/apicurio-registry-client-sdk-go"
	"io"
)

type Client struct {
	apiClient *registryclient.APIClient
}

func NewClient(host string) Client {
	cfg := registryclient.NewConfiguration()
	cfg.Host = host
	cfg.Scheme = "http"

	cfg.OperationServers["ArtifactsApiService.GetLatestArtifact"] = registryclient.ServerConfigurations{
		registryclient.ServerConfiguration{
			URL: "/apis/registry/v2",
		},
	}

	return Client{
		apiClient: registryclient.NewAPIClient(cfg),
	}
}

func (c Client) GetLatestArtifact(ctx context.Context, group, artifactId string) (string, error) {
	_, resp, err := c.apiClient.ArtifactsApi.GetLatestArtifact(ctx, group, artifactId).Execute()
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
