package pocdemo

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

const (
	EMAIL    = "amna@breu.io"
	PASSWORD = "amna"
)

var (
	coreClient *core.CoreClient
	gitClient  *github.GithubClient
	authClient *auth.AuthClient
)

const (
	SNS_REPO_ID  = "53c7cecb-5923-4479-8e16-2147534ee97c"
	SQS_REPO_ID  = "87760bf4-88c0-4d4d-8537-05a6303002e8"
	TEST_REPO_ID = "2eb0688b-9c3b-4918-aeaf-87c6eef9a76b"
)

func Main_pocdemo() {

	ctx := context.Background()
	url := "http://localhost:8000"

	authClient = &auth.AuthClient{}
	gitClient = &github.GithubClient{}
	coreClient = &core.CoreClient{}

	authClient.SetupAPIClient(url)
	gitClient.SetupAPIClient(url)
	coreClient.SetupAPIClient(url)

	request := auth.LoginRequest{Email: EMAIL, Password: PASSWORD}
	ctx = context.WithValue(ctx, shared.URL, url)
	accessToken := authClient.Login(ctx, request)
	ctx = context.WithValue(ctx, shared.UserAccessToken, accessToken)

	// installationId := "38385949"
	// gitClient.GithubWebhookAppInstalled(ctx, installationId)
	// gitClient.CompleteInstallation(ctx, installationId)

	stackName := "Quantum POC"
	stackID, _ := coreClient.CreateStack(ctx, stackName)
	ctx = createRepos(ctx, stackID)
	ctx = createResources(ctx, stackID)
	createWorkloads(ctx, stackID)
	regions := core.BluePrintRegions{Aws: []string{"us-west1"}, Gcp: []string{"asia-southeast1"}}
	bp := core.BlueprintCreateRequest{Name: "Helloworld blueprint", StackID: stackID, RolloutBudget: "300", Regions: regions, ProviderConfig: `{"project": "breu-dev"}`}
	coreClient.CreateBlueprint(ctx, bp)
}

func createRepos(ctx context.Context, stackID gocql.UUID) context.Context {

	request := core.RepoCreateRequest{Name: "HelloWorld", ProviderID: "648084184", DefaultBranch: "main", Provider: "github", IsMonorepo: true, StackID: stackID}
	id := coreClient.CreateRepo(ctx, request)
	ctx = context.WithValue(ctx, "hello_world_id", id)

	return ctx
}

func createResources(ctx context.Context, stackID gocql.UUID) context.Context {

	request := core.ResourceCreateRequest{Name: "CloudRun_HelloWorld", Driver: core.DriverCloudrun.String(), Provider: "GCP", Immutable: true, StackID: stackID,
		Config: `{"properties":{"generation":"second-generation","cpu":"2000m","memory":"1024Mi"},"output":{"env":[{"url":"CloudRun_HelloWorld_URL"}]}}`}
	rsid := coreClient.CreateResource(ctx, request)

	fmt.Printf("resource created with id: %v\n", rsid)
	ctx = context.WithValue(ctx, "rsid", rsid)
	return ctx
}

func createWorkloads(ctx context.Context, stackID gocql.UUID) {

	repoID, _ := gocql.ParseUUID(SNS_REPO_ID)
	if ctx.Value("hello_world_id") != nil {
		repoID = ctx.Value("hello_world_id").(gocql.UUID)
	}

	var rsid gocql.UUID
	if ctx.Value("rsid") != nil {
		rsid = ctx.Value("rsid").(gocql.UUID)
	}

	fmt.Printf("rsid: %v\n", rsid)
	request := core.WorkloadCreateRequest{Name: "helloworld", Kind: "worker", RepoID: repoID, StackID: stackID, RepoPath: "https://github.com/amnabreu/HelloWorld",
		ResourceID: rsid, Container: `{"image": "asia-southeast1-docker.pkg.dev/breu-dev/ctrlplane/helloworld"}`}
	coreClient.CreateWorkload(ctx, request)
}
