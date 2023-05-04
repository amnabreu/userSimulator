package auth

import (
	"context"
	"fmt"

	"go.breu.io/ctrlplane/internal/shared"
)

type AuthClient struct {
	client *Client
}

func (c *AuthClient) SetupAPIClient(url string) {
	var err error
	c.client, err = NewClient(url)
	if err != nil {
		panic(fmt.Sprintf("failed to create api client for auth: %v", err))
	}
}

func (c *AuthClient) Login(ctx context.Context, request LoginRequest) string {
	response := shared.ValidateHttpResponse(c.client.Login(ctx, request))
	parsedResp, err := ParseLoginResponse(response)

	if err != nil {
		panic("Error: Unable to parse Login response")
	}

	println("User logged in")

	return parsedResp.JSON200.AccessToken
}

func (c *AuthClient) RegisterRequest(ctx context.Context) {
	registerRequestBody := RegisterationRequest{
		FirstName:       "Amna",
		LastName:        "Tehreem",
		Email:           "amna@breu.io",
		Password:        "amna",
		ConfirmPassword: "amna",
		TeamName:        "ctrlplane",
	}

	response, err := c.client.Register(ctx, registerRequestBody)

	if err != nil {
		panic(fmt.Sprintf("failed to register user: %v", err))
	}

	parsedResp, err := ParseRegisterResponse(response)
	if err != nil {
		panic("Error: Unable to parse register response")
	}

	switch response.StatusCode {
	case 200:
		fmt.Println("new user added")
		return
	case 400:
		err, _ := parsedResp.JSON400.Errors.Get("email")
		if err == "already exists" {
			fmt.Println("User already exists")
		} else {
			panic("Unable to register user")
		}
	default:
		panic("Unable to register user")
	}
}
