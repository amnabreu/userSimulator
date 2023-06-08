package main

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

func main() {

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

	var installationId string

	stackID, _ := gocql.ParseUUID("00cac1f7-6f09-4cf4-8221-f4b8caf1cce3")
	for {
		println(`please select the option 
		3. github app Installation
		4. Create stack and repos
		5. create resources
		6. create workload
		7. create blueprint
		8. create stack, repos, resources, workload, blueprint
		9. Pull request
		`,
		)

		var option int
		fmt.Scan(&option)

		installationId = "38385949"

		switch option {
		case 1:
			authClient.RegisterRequest(ctx)
		case 2:
			request := auth.LoginRequest{Email: EMAIL, Password: PASSWORD}
			accessToken := authClient.Login(ctx, request)
			ctx = context.WithValue(ctx, shared.UserAccessToken, accessToken)
		case 3:
			// installationId = readInstallationId(installationId)
			gitClient.GithubWebhookAppInstalled(ctx, installationId)
			gitClient.CompleteInstallation(ctx, installationId)
		case 4:
			stackName := "AWS stack"
			stackID, _ = coreClient.CreateStack(ctx, stackName)
			ctx = createRepos(ctx, stackID)
		case 5:
			createResources(ctx, stackID)
		case 6:
			createWorkloads(ctx, stackID)
		case 7:
			regions := core.BluePrintRegions{Aws: []string{"us-west1"}, Gcp: []string{"asia-southeast1"}}
			bp := core.BlueprintCreateRequest{Name: "pubsub blueprint", StackID: stackID, RolloutBudget: "300", Regions: regions}
			coreClient.CreateBlueprint(ctx, bp)

		case 8:
			stackName := "AWS stack"
			stackID, _ = coreClient.CreateStack(ctx, stackName)
			ctx = createRepos(ctx, stackID)
			createResources(ctx, stackID)
			createWorkloads(ctx, stackID)
			regions := core.BluePrintRegions{Aws: []string{"us-west1"}, Gcp: []string{"asia-southeast1"}}
			bp := core.BlueprintCreateRequest{Name: "pubsub blueprint", StackID: stackID, RolloutBudget: "300", Regions: regions}
			coreClient.CreateBlueprint(ctx, bp)
		case 9:
			// installationId = readInstallationId(installationId)
			gitClient.GithubWebHookPullRequest(ctx, installationId)

		default:
			println("Please select valid option")
		}
	}
}

func createRepos(ctx context.Context, stackID gocql.UUID) context.Context {
	// request := core.RepoCreateRequest{Name: "snspublisher", ProviderID: "620175899", DefaultBranch: "master", Provider: "github", IsMonorepo: true, StackID: stackID}
	// id := coreClient.CreateRepo(ctx, request)
	// ctx = context.WithValue(ctx, "sns_repo_id", id)

	// request = core.RepoCreateRequest{Name: "sqssubscriber", ProviderID: "620179406", DefaultBranch: "master", Provider: "github", IsMonorepo: true, StackID: stackID}
	// id = coreClient.CreateRepo(ctx, request)
	// ctx = context.WithValue(ctx, "sqs_repo_id", id)

	// request = core.RepoCreateRequest{Name: "test", ProviderID: "611620220", DefaultBranch: "main", Provider: "github", IsMonorepo: true, StackID: stackID}
	// id = coreClient.CreateRepo(ctx, request)
	// ctx = context.WithValue(ctx, "test_repo_id", id)

	request := core.RepoCreateRequest{Name: "HelloWorld", ProviderID: "648084184", DefaultBranch: "main", Provider: "github", IsMonorepo: true, StackID: stackID}
	id := coreClient.CreateRepo(ctx, request)
	ctx = context.WithValue(ctx, "hello_world_id", id)

	return ctx
}

func createResources(ctx context.Context, stackID gocql.UUID) {
	// request := core.ResourceCreateRequest{Name: "SNS_TOPIC", Driver: "SNS", Provider: "AWS", Immutable: true, StackID: stackID, Config: ""}
	// coreClient.CreateResource(ctx, request)

	// request = core.ResourceCreateRequest{Name: "SQS_QUEUE", Driver: "SQS", Provider: "AWS", Immutable: true, StackID: stackID, Config: ""}
	// coreClient.CreateResource(ctx, request)

	// request = core.ResourceCreateRequest{Name: "GKE_Cluster", Driver: "GKE", Provider: "GCP", Immutable: true, StackID: stackID, Config: ""}
	// coreClient.CreateResource(ctx, request)

	request := core.ResourceCreateRequest{Name: "CloudRun_HelloWorld", Driver: core.DriverCloudrun.String(), Provider: "GCP", Immutable: true, StackID: stackID, Config: `{
		"properties":{
		   "generation":"second-generation",
		   "cpu":"2000m",
		   "memory":"1024Mi"
		},
		"output":{
		   "env":[
			  {
				 "url":"CloudRun_HelloWorld_URL"
			  }
		   ]
		}
	 }`}
	coreClient.CreateResource(ctx, request)
}

func createWorkloads(ctx context.Context, stackID gocql.UUID) {

	repoID, _ := gocql.ParseUUID(SNS_REPO_ID)
	if ctx.Value("hello_world_id") != nil {
		repoID = ctx.Value("hello_world_id").(gocql.UUID)
	}
	request := core.WorkloadCreateRequest{Name: "hello_world", Kind: "worker", RepoID: repoID, StackID: stackID, RepoPath: "https://github.com/amnabreu/HelloWorld"}
	coreClient.CreateWorkload(ctx, request)

	// // sns workload
	// repoID, _ := gocql.ParseUUID(SNS_REPO_ID)
	// if ctx.Value("sns_repo_id") != nil {
	// 	repoID = ctx.Value("sns_repo_id").(gocql.UUID)
	// }
	// request := core.WorkloadCreateRequest{Name: "sns_publisher", Kind: "worker", RepoID: repoID, StackID: stackID, RepoPath: "https://github.com/amnabreu/snsPublisher"}
	// coreClient.CreateWorkload(ctx, request)

	// // sqs workload
	// repoID, _ = gocql.ParseUUID(SQS_REPO_ID)
	// if ctx.Value("sqs_repo_id") != nil {
	// 	repoID = ctx.Value("sqs_repo_id").(gocql.UUID)
	// }
	// request = core.WorkloadCreateRequest{Name: "sqs_publisher", Kind: "worker", RepoID: repoID, StackID: stackID, RepoPath: "https://github.com/amnabreu/sqsSubscriber"}
	// coreClient.CreateWorkload(ctx, request)

	// // test workload
	// repoID, _ = gocql.ParseUUID(TEST_REPO_ID)
	// if ctx.Value("test_repo_id") != nil {
	// 	repoID = ctx.Value("test_repo_id").(gocql.UUID)
	// }
	// request = core.WorkloadCreateRequest{Name: "test_repo", Kind: "worker", RepoID: repoID, StackID: stackID, RepoPath: "https://github.com/amnabreu/testGithubapp"}
	// coreClient.CreateWorkload(ctx, request)

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
