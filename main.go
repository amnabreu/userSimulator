package main

import (
	"context"
	"log"
	"net/http"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/providers/github"
)

func main() {

	ctx := context.Background()
	url := "http://localhost:8000"
	authclient, gitClient := SetupAPIClient(url)

	RegisterRequest(ctx, authclient)
	accessToken := LoginRequest(ctx, authclient)

	ctx = context.WithValue(ctx, "url", url)
	ctx = context.WithValue(ctx, "access_token", accessToken)
	CompleteInstallation(ctx, gitClient)

}

func CompleteInstallation(ctx context.Context, client *github.Client) {

	completeInstallationBody := github.CompleteInstallationRequest{}
	completeInstallationBody.InstallationId = 35046675
	completeInstallationBody.SetupAction = github.SetupActionCreated

	response := validateHttpResponse(client.GithubCompleteInstallation(ctx, completeInstallationBody, AddAuthHeader))
	parsedResp, err := github.ParseGithubCompleteInstallationResponse(response)
	if err != nil {
		log.Panicf("Error: Unable to parse register response: %v", err)
	}

	log.Printf("Run ID: %s, status: %s", parsedResp.JSON200.RunID, parsedResp.JSON200.Status)
}

func validateHttpResponse(response *http.Response, err error) *http.Response {
	if err != nil {
		log.Panicf("Request failed with error: %v", err)
	}

	if response.StatusCode != 200 {
		log.Panicf("Error: complete Installation requested failed with status: %d", response.StatusCode)
	}

	return response
}

func AddAuthHeader(ctx context.Context, req *http.Request) error {

	req.Header.Set("authorization", "Token "+ctx.Value("access_token").(string))
	return nil
}

func SetupAPIClient(url string) (*auth.Client, *github.Client) {
	authclient, err := auth.NewClient(url)
	if err != nil {
		println("failed to create api client for auth: %v", err)
	}

	gitclient, err := github.NewClient(url)
	if err != nil {
		println("failed to create api client for github: %v", err)
	}

	return authclient, gitclient
}

func RegisterRequest(ctx context.Context, client *auth.Client) {
	registerRequestBody := auth.RegisterationRequest{}
	registerRequestBody.FirstName = "Amna"
	registerRequestBody.LastName = "Tehreem"
	registerRequestBody.Email = "amna@breu.io"
	registerRequestBody.Password = "amna"
	registerRequestBody.ConfirmPassword = "amna"
	registerRequestBody.TeamName = "ctrlplane"

	response, err := client.Register(ctx, registerRequestBody)

	if err != nil {
		panic("Error: Unable to register user")
	}

	parsedResp, err := auth.ParseRegisterResponse(response)

	if err != nil {
		panic("Error: Unable to parse register response")
	}

	switch response.StatusCode {
	case 200:
		log.Println("new user added")
		return
	case 400:
		err, _ := parsedResp.JSON400.Errors.Get("email")
		if err == "already exists" {
			log.Println("User already exists")
		} else {
			panic("Unable to register user")
		}
	default:
		panic("Unable to register user")
	}
}

func LoginRequest(ctx context.Context, client *auth.Client) string {
	loginRequestBody := auth.LoginRequest{}
	loginRequestBody.Email = "amna@breu.io"
	loginRequestBody.Password = "amna"

	response := validateHttpResponse(client.Login(ctx, loginRequestBody))
	parsedResp, err := auth.ParseLoginResponse(response)

	if err != nil {
		panic("Error: Unable to parse Login response")
	}

	return parsedResp.JSON200.AccessToken

}
