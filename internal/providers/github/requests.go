package github

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"go.breu.io/ctrlplane/internal/shared"
)

const (
	GithubEvent     = "github_event"
	UserAccessToken = "access_token"
	URL             = "url"
)

type GithubClient struct {
	client *Client
}

func (c *GithubClient) SetupAPIClient(url string) {
	var err error
	c.client, err = NewClient(url)
	if err != nil {
		panic(fmt.Sprintf("failed to create api client for github: %v", err))
	}
}

func AddGithubHeader(ctx context.Context, req *http.Request) error {

	req.Header.Set("X-GitHub-Event", ctx.Value(GithubEvent).(string))
	req.Header.Set("X-Hub-Signature-256", "sha256=e1d2476c46fba5a53d1937e5af702f6803d1766f27c849235faa6f22efda5deb")
	return nil
}

func (c *GithubClient) GithubWebHookPullRequest(ctx context.Context, installationId string) {
	url := ctx.Value(URL).(string) + "/providers/github/webhook"

	// var jsonStr = []byte(``)
	var jsonStr = []byte(GetPRBody(installationId))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-GitHub-Event", "pull_request")
	req.Header.Set("X-Hub-Signature-256", "sha256=006239c99fd6eba2004765ab97e06d18f9d18c21cf122474bd73fc89560c78de")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	shared.ValidateHttpResponse(client.Do(req))
	fmt.Printf("pull request webhook sent with installation id: %s\n", installationId)
}

func (c *GithubClient) GithubWebhookAppInstalled(ctx context.Context, installationId string) {
	url := ctx.Value(URL).(string) + "/providers/github/webhook"
	fmt.Println("URL:>", url)

	// var jsonStr = []byte(``)
	var jsonStr = []byte(GetWebhookInstallationBody(installationId))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-GitHub-Event", "installation")
	req.Header.Set("X-Hub-Signature-256", "sha256=006239c99fd6eba2004765ab97e06d18f9d18c21cf122474bd73fc89560c78de")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	shared.ValidateHttpResponse(client.Do(req))
	fmt.Printf("Installation webhook sent with installation id: %s\n", installationId)
}

func (c *GithubClient) CompleteInstallation(ctx context.Context, installationId string) {

	installationIdint, err := strconv.ParseInt(installationId, 10, 64)
	if err != nil {
		fmt.Printf("CompleteInstallation: cannot covert installation id: %s to integer\n", installationId)
	}
	completeInstallationBody := CompleteInstallationRequest{
		InstallationId: installationIdint,
		SetupAction:    SetupActionCreated,
	}

	response := shared.ValidateHttpResponse(c.client.GithubCompleteInstallation(ctx, completeInstallationBody, shared.AddAuthHeader))
	parsedResp, err := ParseGithubCompleteInstallationResponse(response)
	if err != nil {
		fmt.Printf("Error: Unable to parse register response: %v\n", err)
	}

	fmt.Printf("Run ID: %s, status: %s\n", parsedResp.JSON200.RunID, parsedResp.JSON200.Status)
}
