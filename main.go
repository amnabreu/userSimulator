package main

import (
	"context"
	"fmt"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

func main() {

	ctx := context.Background()
	url := "http://localhost:8000"

	authClient := &auth.AuthClient{}
	gitClient := &github.GithubClient{}
	coreClient := &core.CoreClient{}

	authClient.SetupAPIClient(url)
	gitClient.SetupAPIClient(url)
	coreClient.SetupAPIClient(url)

	ctx = context.WithValue(ctx, shared.URL, url)
	var installationId string

	for {
		println(`please select the option 
		1. Register user
		2. Login user
		3. Installation webhook 
		4. Installation complete
		5. Pull request
		6. Create stack and test repo`,
		)

		var option int
		fmt.Scan(&option)

		switch option {
		case 1:
			authClient.RegisterRequest(ctx)
		case 2:
			accessToken := authClient.LoginRequest(ctx)
			ctx = context.WithValue(ctx, shared.UserAccessToken, accessToken)
		case 3:
			installationId = readInstallationId(installationId)
			gitClient.GithubWebhookAppInstalled(ctx, installationId)

		case 4:
			installationId = readInstallationId(installationId)
			gitClient.CompleteInstallation(ctx, installationId)
		case 5:
			installationId = readInstallationId(installationId)
			gitClient.GithubWebHookPullRequest(ctx, installationId)
		case 6:
			installationId = "123456789"
			stackName := "AWS stack"
			providerID := "611620220"
			stackID, err := coreClient.CreateStack(ctx, stackName)
			if err != nil {
				fmt.Printf("Unable to create stack, error:%v", err)
				break
			}

			coreClient.CreateRepo(ctx, stackID, providerID)

		default:
			println("Please select valid option")
		}
	}

}

func readInstallationId(installId string) string {
	var id string
	fmt.Print("Enter installation id, press enter to use previous installation id:")
	fmt.Scan(&id)
	if id == "" {
		return installId
	} else {
		return id
	}
}
