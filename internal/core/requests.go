package core

import (
	"context"
	"fmt"
	"io"

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

func (c *CoreClient) CreateRepo(ctx context.Context, request RepoCreateRequest) gocql.UUID {
	response := shared.ValidateHttpResponse(c.client.CreateRepo(ctx, request, shared.AddAuthHeader))
	fmt.Printf("The following repo is successfully created")
	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body[:]))

	parsedresp, err := ParseCreateRepoResponse(response)
	if err != nil {
		panic(fmt.Sprintf("unable to parse create repo response: %v", err))
	}

	return parsedresp.JSON201.ID

}

func (c *CoreClient) CreateResource(ctx context.Context, request ResourceCreateRequest) {
	response := shared.ValidateHttpResponse(c.client.CreateResource(ctx, request, shared.AddAuthHeader))
	fmt.Println("The following resource is successfully created")
	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body[:]))
}

func (c *CoreClient) CreateWorkload(ctx context.Context, request WorkloadCreateRequest) {
	response := shared.ValidateHttpResponse(c.client.CreateWorkload(ctx, request, shared.AddAuthHeader))
	fmt.Println("The following workload is successfully created")
	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body[:]))
}

func (c *CoreClient) CreateBlueprint(ctx context.Context, request BlueprintCreateRequest) {
	response := shared.ValidateHttpResponse(c.client.CreateBlueprint(ctx, request, shared.AddAuthHeader))
	fmt.Println("The following blueprint is successfully created")
	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body[:]))
}
