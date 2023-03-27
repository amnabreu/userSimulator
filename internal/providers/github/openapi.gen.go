// Package github provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen, a modified copy of github.com/deepmap/oapi-codegen.
// It was modified to add support for the following features:
//  - Support for custom templates by filename.
//  - Supporting x-breu-entity in the schema to generate a struct for the entity.
//
// DO NOT EDIT!!

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/scylladb/gocqlx/v2/table"
	externalRef1 "go.breu.io/ctrlplane/internal/shared"
)

const (
	APIKeyAuthScopes = "APIKeyAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

var (
	ErrInvalidSetupAction    = errors.New("invalid SetupAction value")
	ErrInvalidWorkflowStatus = errors.New("invalid WorkflowStatus value")
)

type (
	SetupActionMapType map[string]SetupAction // SetupActionMapType is a quick lookup map for SetupAction.
)

// Defines values for SetupAction.
const (
	SetupActionCreated SetupAction = "created"
	SetupActionDeleted SetupAction = "deleted"
	SetupActionUpdated SetupAction = "updated"
)

// SetupActionValues returns all known values for SetupAction.
var (
	SetupActionMap = SetupActionMapType{
		SetupActionCreated.String(): SetupActionCreated,
		SetupActionDeleted.String(): SetupActionDeleted,
		SetupActionUpdated.String(): SetupActionUpdated,
	}
)

/*
 * Helper methods for SetupAction for easy marshalling and unmarshalling.
 */
func (v SetupAction) String() string               { return string(v) }
func (v SetupAction) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *SetupAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := SetupActionMap[s]
	if !ok {
		return ErrInvalidSetupAction
	}

	*v = val

	return nil
}

type (
	WorkflowStatusMapType map[string]WorkflowStatus // WorkflowStatusMapType is a quick lookup map for WorkflowStatus.
)

// Defines values for WorkflowStatus.
const (
	WorkflowStatusFailure  WorkflowStatus = "failure"
	WorkflowStatusQueued   WorkflowStatus = "queued"
	WorkflowStatusSignaled WorkflowStatus = "signaled"
	WorkflowStatusSkipped  WorkflowStatus = "skipped"
	WorkflowStatusSuccess  WorkflowStatus = "success"
)

// WorkflowStatusValues returns all known values for WorkflowStatus.
var (
	WorkflowStatusMap = WorkflowStatusMapType{
		WorkflowStatusFailure.String():  WorkflowStatusFailure,
		WorkflowStatusQueued.String():   WorkflowStatusQueued,
		WorkflowStatusSignaled.String(): WorkflowStatusSignaled,
		WorkflowStatusSkipped.String():  WorkflowStatusSkipped,
		WorkflowStatusSuccess.String():  WorkflowStatusSuccess,
	}
)

/*
 * Helper methods for WorkflowStatus for easy marshalling and unmarshalling.
 */
func (v WorkflowStatus) String() string               { return string(v) }
func (v WorkflowStatus) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *WorkflowStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := WorkflowStatusMap[s]
	if !ok {
		return ErrInvalidWorkflowStatus
	}

	*v = val

	return nil
}

// CompleteInstallationRequest complete the installation given the installation_id & setup_action
type CompleteInstallationRequest struct {
	InstallationId int64       `json:"installation_id"`
	SetupAction    SetupAction `json:"setup_action"`
}

// Installation defines model for GithubInstallation.
type Installation struct {
	CreatedAt         time.Time  `cql:"created_at" json:"created_at"`
	ID                gocql.UUID `cql:"id" json:"id"`
	InstallationID    int64      `cql:"installation_id" json:"installation_id" validate:"required,db_unique"`
	InstallationLogin string     `cql:"installation_login" json:"installation_login"`
	InstallationType  string     `cql:"installation_type" json:"installation_type"`
	SenderID          int64      `cql:"sender_id" json:"sender_id"`
	SenderLogin       string     `cql:"sender_login" json:"sender_login"`
	Status            string     `cql:"status" json:"status"`
	TeamID            gocql.UUID `cql:"team_id" json:"team_id"`
	UpdatedAt         time.Time  `cql:"updated_at" json:"updated_at"`
}

var (
	githubinstallationColumns = []string{"created_at", "id", "installation_id", "installation_login", "installation_type", "sender_id", "sender_login", "status", "team_id", "updated_at"}

	githubinstallationMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_installations",
			Columns: githubinstallationColumns,
		},
	}

	githubinstallationTable = itable.New(*githubinstallationMeta.M)
)

func (githubinstallation *Installation) GetTable() itable.ITable {
	return githubinstallationTable
}

// Repo defines model for GithubRepo.
type Repo struct {
	CreatedAt      time.Time  `cql:"created_at" json:"created_at"`
	FullName       string     `cql:"full_name" json:"full_name"`
	GithubID       int64      `cql:"github_id" json:"github_id"`
	ID             gocql.UUID `cql:"id" json:"id"`
	InstallationID int64      `cql:"installation_id" json:"installation_id"`
	Name           string     `cql:"name" json:"name"`
	TeamID         gocql.UUID `cql:"team_id" json:"team_id"`
	UpdatedAt      time.Time  `cql:"updated_at" json:"updated_at"`
}

var (
	githubrepoColumns = []string{"created_at", "full_name", "github_id", "id", "installation_id", "name", "team_id", "updated_at"}

	githubrepoMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_repos",
			Columns: githubrepoColumns,
		},
	}

	githubrepoTable = itable.New(*githubrepoMeta.M)
)

func (githubrepo *Repo) GetTable() itable.ITable {
	return githubrepoTable
}

// SetupAction defines model for SetupAction.
type SetupAction string

// WorkflowResponse workflow status & run id
type WorkflowResponse struct {
	RunID string `json:"run_id"`

	// Status the workflow status
	Status WorkflowStatus `json:"status"`
}

// WorkflowStatus the workflow status
type WorkflowStatus string

// GithubCompleteInstallationJSONRequestBody defines body for GithubCompleteInstallation for application/json ContentType.
type GithubCompleteInstallationJSONRequestBody = CompleteInstallationRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GithubCompleteInstallation request with any body
	GithubCompleteInstallationWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	GithubCompleteInstallation(ctx context.Context, body GithubCompleteInstallationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GithubGetInstallations request
	GithubGetInstallations(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GithubGetRepos request
	GithubGetRepos(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GithubWebhook request
	GithubWebhook(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GithubCompleteInstallationWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGithubCompleteInstallationRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GithubCompleteInstallation(ctx context.Context, body GithubCompleteInstallationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGithubCompleteInstallationRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("authorization", "Token abc")
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GithubGetInstallations(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGithubGetInstallationsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GithubGetRepos(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGithubGetReposRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GithubWebhook(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGithubWebhookRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGithubCompleteInstallationRequest calls the generic GithubCompleteInstallation builder with application/json body
func NewGithubCompleteInstallationRequest(server string, body GithubCompleteInstallationJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewGithubCompleteInstallationRequestWithBody(server, "application/json", bodyReader)
}

// NewGithubCompleteInstallationRequestWithBody generates requests for GithubCompleteInstallation with any type of body
func NewGithubCompleteInstallationRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/providers/github/complete-installation")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGithubGetInstallationsRequest generates requests for GithubGetInstallations
func NewGithubGetInstallationsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/providers/github/installations")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGithubGetReposRequest generates requests for GithubGetRepos
func NewGithubGetReposRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/providers/github/repos")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGithubWebhookRequest generates requests for GithubWebhook
func NewGithubWebhookRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/providers/github/webhook")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GithubCompleteInstallation request with any body
	GithubCompleteInstallationWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GithubCompleteInstallationResponse, error)

	GithubCompleteInstallationWithResponse(ctx context.Context, body GithubCompleteInstallationJSONRequestBody, reqEditors ...RequestEditorFn) (*GithubCompleteInstallationResponse, error)

	// GithubGetInstallations request
	GithubGetInstallationsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GithubGetInstallationsResponse, error)

	// GithubGetRepos request
	GithubGetReposWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GithubGetReposResponse, error)

	// GithubWebhook request
	GithubWebhookWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GithubWebhookResponse, error)
}

type GithubCompleteInstallationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *WorkflowResponse
	JSON201      *WorkflowResponse
	JSON400      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r GithubCompleteInstallationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GithubCompleteInstallationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GithubGetInstallationsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Installation
	JSON400      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r GithubGetInstallationsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GithubGetInstallationsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GithubGetReposResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Repo
	JSON400      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r GithubGetReposResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GithubGetReposResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GithubWebhookResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *WorkflowResponse
	JSON201      *WorkflowResponse
	JSON400      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r GithubWebhookResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GithubWebhookResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GithubCompleteInstallationWithBodyWithResponse request with arbitrary body returning *GithubCompleteInstallationResponse
func (c *ClientWithResponses) GithubCompleteInstallationWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GithubCompleteInstallationResponse, error) {
	rsp, err := c.GithubCompleteInstallationWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGithubCompleteInstallationResponse(rsp)
}

func (c *ClientWithResponses) GithubCompleteInstallationWithResponse(ctx context.Context, body GithubCompleteInstallationJSONRequestBody, reqEditors ...RequestEditorFn) (*GithubCompleteInstallationResponse, error) {
	rsp, err := c.GithubCompleteInstallation(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGithubCompleteInstallationResponse(rsp)
}

// GithubGetInstallationsWithResponse request returning *GithubGetInstallationsResponse
func (c *ClientWithResponses) GithubGetInstallationsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GithubGetInstallationsResponse, error) {
	rsp, err := c.GithubGetInstallations(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGithubGetInstallationsResponse(rsp)
}

// GithubGetReposWithResponse request returning *GithubGetReposResponse
func (c *ClientWithResponses) GithubGetReposWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GithubGetReposResponse, error) {
	rsp, err := c.GithubGetRepos(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGithubGetReposResponse(rsp)
}

// GithubWebhookWithResponse request returning *GithubWebhookResponse
func (c *ClientWithResponses) GithubWebhookWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GithubWebhookResponse, error) {
	rsp, err := c.GithubWebhook(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGithubWebhookResponse(rsp)
}

// ParseGithubCompleteInstallationResponse parses an HTTP response from a GithubCompleteInstallationWithResponse call
func ParseGithubCompleteInstallationResponse(rsp *http.Response) (*GithubCompleteInstallationResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GithubCompleteInstallationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest WorkflowResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest WorkflowResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGithubGetInstallationsResponse parses an HTTP response from a GithubGetInstallationsWithResponse call
func ParseGithubGetInstallationsResponse(rsp *http.Response) (*GithubGetInstallationsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GithubGetInstallationsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Installation
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGithubGetReposResponse parses an HTTP response from a GithubGetReposWithResponse call
func ParseGithubGetReposResponse(rsp *http.Response) (*GithubGetReposResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GithubGetReposResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Repo
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGithubWebhookResponse parses an HTTP response from a GithubWebhookWithResponse call
func ParseGithubWebhookResponse(rsp *http.Response) (*GithubWebhookResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GithubWebhookResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest WorkflowResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest WorkflowResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Complete GitHub App installation
	// (POST /providers/github/complete-installation)
	GithubCompleteInstallation(ctx echo.Context) error
	// Get GitHub installations
	// (GET /providers/github/installations)
	GithubGetInstallations(ctx echo.Context) error
	// Get GitHub repositories
	// (GET /providers/github/repos)
	GithubGetRepos(ctx echo.Context) error
	// Webhook reciever for github
	// (POST /providers/github/webhook)
	GithubWebhook(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GithubCompleteInstallation converts echo context to params.
func (w *ServerInterfaceWrapper) GithubCompleteInstallation(ctx echo.Context) error {

	println("calling SIW GithubCompleteInstallation")
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GithubCompleteInstallation(ctx)
	return err
}

// GithubGetInstallations converts echo context to params.
func (w *ServerInterfaceWrapper) GithubGetInstallations(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GithubGetInstallations(ctx)
	return err
}

// GithubGetRepos converts echo context to params.
func (w *ServerInterfaceWrapper) GithubGetRepos(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GithubGetRepos(ctx)
	return err
}

// GithubWebhook converts echo context to params.
func (w *ServerInterfaceWrapper) GithubWebhook(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GithubWebhook(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/providers/github/complete-installation", wrapper.GithubCompleteInstallation)
	router.GET(baseURL+"/providers/github/installations", wrapper.GithubGetInstallations)
	router.GET(baseURL+"/providers/github/repos", wrapper.GithubGetRepos)
	router.POST(baseURL+"/providers/github/webhook", wrapper.GithubWebhook)

}
