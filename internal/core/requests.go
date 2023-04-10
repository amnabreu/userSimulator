package core

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/shared"
)

type CoreClient struct {
	client *Client
}

func (c *CoreClient) SetupAPIClient(url string) {
	var err error
	c.client, err = NewClient(url)
	if err != nil {
		panic(fmt.Sprintf("failed to create api client for core: %v", err))
	}
}

func (c *CoreClient) CreateStack(ctx context.Context, stackName string) (gocql.UUID, error) {

	csbody := StackCreateRequest{Name: stackName}
	response := shared.ValidateHttpResponse(c.client.CreateStack(ctx, csbody, shared.AddAuthHeader))
	parsedResp, err := ParseCreateStackResponse(response)
	if err != nil {
		fmt.Printf("Error: Unable to parse create stack response: %v\n", err)
		return gocql.UUID{}, err
	}

	fmt.Printf("Stack created with ID: %s", parsedResp.JSON201.ID.String())
	return parsedResp.JSON201.ID, nil
}

func (c *CoreClient) CreateRepo(ctx context.Context, stackID gocql.UUID, providerID string) {

	crbody := RepoCreateRequest{
		StackID:       stackID,
		Provider:      "github",
		IsMonorepo:    false,
		DefaultBranch: "main",
		ProviderID:    providerID,
	}

	shared.ValidateHttpResponse(c.client.CreateRepo(ctx, crbody, shared.AddAuthHeader))
	fmt.Printf("Repo created")
}
