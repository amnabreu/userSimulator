// Package core provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen, a modified copy of github.com/deepmap/oapi-codegen.
// It was modified to add support for the following features:
//  - Support for custom templates by filename.
//  - Supporting x-breu-entity in the schema to generate a struct for the entity.
//
// DO NOT EDIT!!

package core

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
	"github.com/deepmap/oapi-codegen/pkg/runtime"
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
	ErrInvalidRepoProvider = errors.New("invalid RepoProvider value")
)

type (
	RepoProviderMapType map[string]RepoProvider // RepoProviderMapType is a quick lookup map for RepoProvider.
)

// Defines values for RepoProvider.
const (
	RepoProviderBitbucket RepoProvider = "bitbucket"
	RepoProviderGithub    RepoProvider = "github"
	RepoProviderGitlab    RepoProvider = "gitlab"
)

// RepoProviderValues returns all known values for RepoProvider.
var (
	RepoProviderMap = RepoProviderMapType{
		RepoProviderBitbucket.String(): RepoProviderBitbucket,
		RepoProviderGithub.String():    RepoProviderGithub,
		RepoProviderGitlab.String():    RepoProviderGitlab,
	}
)

/*
 * Helper methods for RepoProvider for easy marshalling and unmarshalling.
 */
func (v RepoProvider) String() string               { return string(v) }
func (v RepoProvider) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *RepoProvider) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := RepoProviderMap[s]
	if !ok {
		return ErrInvalidRepoProvider
	}

	*v = val

	return nil
}

// Repo defines model for Repo.
type Repo struct {
	CreatedAt     time.Time    `cql:"created_at" json:"created_at"`
	DefaultBranch string       `cql:"default_branch" json:"default_branch"`
	ID            gocql.UUID   `cql:"id" json:"id"`
	IsMonorepo    bool         `cql:"is_monorepo" json:"is_monorepo"`
	Provider      RepoProvider `cql:"provider" json:"provider"`
	ProviderID    string       `cql:"provider_id" json:"provider_id"`
	StackID       gocql.UUID   `cql:"stack_id" json:"stack_id"`
	UpdatedAt     time.Time    `cql:"updated_at" json:"updated_at"`
}

var (
	repoColumns = []string{"created_at", "default_branch", "id", "is_monorepo", "provider", "provider_id", "stack_id", "updated_at"}

	repoMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "repos",
			Columns: repoColumns,
		},
	}

	repoTable = itable.New(*repoMeta.M)
)

func (repo *Repo) GetTable() itable.ITable {
	return repoTable
}

// RepoCreateRequest defines model for RepoCreateRequest.
type RepoCreateRequest struct {
	DefaultBranch string       `json:"default_branch"`
	IsMonorepo    bool         `json:"is_monorepo"`
	Provider      RepoProvider `json:"provider"`
	ProviderID    string       `json:"provider_id"`
	StackID       gocql.UUID   `json:"stack_id"`
}

// RepoListResponse defines model for RepoListResponse.
type RepoListResponse = []Repo

// RepoProvider defines model for RepoProvider.
type RepoProvider string

// Stack defines model for Stack.
type Stack struct {
	Config    StackConfig `cql:"config" json:"config"`
	CreatedAt time.Time   `cql:"created_at" json:"created_at"`
	ID        gocql.UUID  `cql:"id" json:"id"`
	Name      string      `cql:"name" json:"name" validate:"required"`
	Slug      string      `cql:"slug" json:"slug"`
	TeamID    gocql.UUID  `cql:"team_id" json:"team_id"`
	UpdatedAt time.Time   `cql:"updated_at" json:"updated_at"`
}

var (
	stackColumns = []string{"config", "created_at", "id", "name", "slug", "team_id", "updated_at"}

	stackMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "stacks",
			Columns: stackColumns,
		},
	}

	stackTable = itable.New(*stackMeta.M)
)

func (stack *Stack) GetTable() itable.ITable {
	return stackTable
}

// StackConfig defines model for StackConfig.
type StackConfig map[string]interface{}

// StackCreateRequest defines model for StackCreateRequest.
type StackCreateRequest struct {
	Config StackConfig `json:"config"`
	Name   string      `json:"name"`
}

// StackListResponse defines model for StackListResponse.
type StackListResponse = []Stack

// CreateRepoJSONRequestBody defines body for CreateRepo for application/json ContentType.
type CreateRepoJSONRequestBody = RepoCreateRequest

// CreateStackJSONRequestBody defines body for CreateStack for application/json ContentType.
type CreateStackJSONRequestBody = StackCreateRequest

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
	// ListRepos request
	ListRepos(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateRepo request with any body
	CreateRepoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateRepo(ctx context.Context, body CreateRepoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetRepo request
	GetRepo(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListStacks request
	ListStacks(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateStack request with any body
	CreateStackWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateStack(ctx context.Context, body CreateStackJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetStack request
	GetStack(ctx context.Context, slug string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListRepos(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListReposRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateRepoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRepoRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateRepo(ctx context.Context, body CreateRepoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRepoRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetRepo(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetRepoRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListStacks(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListStacksRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateStackWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateStackRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateStack(ctx context.Context, body CreateStackJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateStackRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetStack(ctx context.Context, slug string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetStackRequest(c.Server, slug)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListReposRequest generates requests for ListRepos
func NewListReposRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/core/repos")
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

// NewCreateRepoRequest calls the generic CreateRepo builder with application/json body
func NewCreateRepoRequest(server string, body CreateRepoJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateRepoRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateRepoRequestWithBody generates requests for CreateRepo with any type of body
func NewCreateRepoRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/core/repos")
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

// NewGetRepoRequest generates requests for GetRepo
func NewGetRepoRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/core/repos/%s", pathParam0)
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

// NewListStacksRequest generates requests for ListStacks
func NewListStacksRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/core/stacks")
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

// NewCreateStackRequest calls the generic CreateStack builder with application/json body
func NewCreateStackRequest(server string, body CreateStackJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateStackRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateStackRequestWithBody generates requests for CreateStack with any type of body
func NewCreateStackRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/core/stacks")
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

// NewGetStackRequest generates requests for GetStack
func NewGetStackRequest(server string, slug string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "slug", runtime.ParamLocationPath, slug)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/core/stacks/%s", pathParam0)
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
	// ListRepos request
	ListReposWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListReposResponse, error)

	// CreateRepo request with any body
	CreateRepoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRepoResponse, error)

	CreateRepoWithResponse(ctx context.Context, body CreateRepoJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRepoResponse, error)

	// GetRepo request
	GetRepoWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetRepoResponse, error)

	// ListStacks request
	ListStacksWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListStacksResponse, error)

	// CreateStack request with any body
	CreateStackWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateStackResponse, error)

	CreateStackWithResponse(ctx context.Context, body CreateStackJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateStackResponse, error)

	// GetStack request
	GetStackWithResponse(ctx context.Context, slug string, reqEditors ...RequestEditorFn) (*GetStackResponse, error)
}

type ListReposResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *RepoListResponse
	JSON404      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r ListReposResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListReposResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateRepoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Repo
	JSON400      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r CreateRepoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateRepoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetRepoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Repo
	JSON404      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r GetRepoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetRepoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListStacksResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *StackListResponse
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r ListStacksResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListStacksResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateStackResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Stack
	JSON400      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r CreateStackResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateStackResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetStackResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Stack
	JSON404      *externalRef1.APIError
	JSON500      *externalRef1.APIError
}

// Status returns HTTPResponse.Status
func (r GetStackResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetStackResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListReposWithResponse request returning *ListReposResponse
func (c *ClientWithResponses) ListReposWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListReposResponse, error) {
	rsp, err := c.ListRepos(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListReposResponse(rsp)
}

// CreateRepoWithBodyWithResponse request with arbitrary body returning *CreateRepoResponse
func (c *ClientWithResponses) CreateRepoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRepoResponse, error) {
	rsp, err := c.CreateRepoWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRepoResponse(rsp)
}

func (c *ClientWithResponses) CreateRepoWithResponse(ctx context.Context, body CreateRepoJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRepoResponse, error) {
	rsp, err := c.CreateRepo(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRepoResponse(rsp)
}

// GetRepoWithResponse request returning *GetRepoResponse
func (c *ClientWithResponses) GetRepoWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetRepoResponse, error) {
	rsp, err := c.GetRepo(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetRepoResponse(rsp)
}

// ListStacksWithResponse request returning *ListStacksResponse
func (c *ClientWithResponses) ListStacksWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListStacksResponse, error) {
	rsp, err := c.ListStacks(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListStacksResponse(rsp)
}

// CreateStackWithBodyWithResponse request with arbitrary body returning *CreateStackResponse
func (c *ClientWithResponses) CreateStackWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateStackResponse, error) {
	rsp, err := c.CreateStackWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateStackResponse(rsp)
}

func (c *ClientWithResponses) CreateStackWithResponse(ctx context.Context, body CreateStackJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateStackResponse, error) {
	rsp, err := c.CreateStack(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateStackResponse(rsp)
}

// GetStackWithResponse request returning *GetStackResponse
func (c *ClientWithResponses) GetStackWithResponse(ctx context.Context, slug string, reqEditors ...RequestEditorFn) (*GetStackResponse, error) {
	rsp, err := c.GetStack(ctx, slug, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetStackResponse(rsp)
}

// ParseListReposResponse parses an HTTP response from a ListReposWithResponse call
func ParseListReposResponse(rsp *http.Response) (*ListReposResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListReposResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest RepoListResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseCreateRepoResponse parses an HTTP response from a CreateRepoWithResponse call
func ParseCreateRepoResponse(rsp *http.Response) (*CreateRepoResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateRepoResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Repo
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

// ParseGetRepoResponse parses an HTTP response from a GetRepoWithResponse call
func ParseGetRepoResponse(rsp *http.Response) (*GetRepoResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetRepoResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Repo
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseListStacksResponse parses an HTTP response from a ListStacksWithResponse call
func ParseListStacksResponse(rsp *http.Response) (*ListStacksResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListStacksResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest StackListResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseCreateStackResponse parses an HTTP response from a CreateStackWithResponse call
func ParseCreateStackResponse(rsp *http.Response) (*CreateStackResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateStackResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Stack
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

// ParseGetStackResponse parses an HTTP response from a GetStackWithResponse call
func ParseGetStackResponse(rsp *http.Response) (*GetStackResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetStackResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Stack
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest externalRef1.APIError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

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
	// List Repos
	// (GET /core/repos)
	ListRepos(ctx echo.Context) error
	// Create repo
	// (POST /core/repos)
	CreateRepo(ctx echo.Context) error
	// Get repo
	// (GET /core/repos/{id})
	GetRepo(ctx echo.Context, id string) error
	// List stacks
	// (GET /core/stacks)
	ListStacks(ctx echo.Context) error
	// Create stack
	// (POST /core/stacks)
	CreateStack(ctx echo.Context) error
	// Get stack
	// (GET /core/stacks/{slug})
	GetStack(ctx echo.Context, slug string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListRepos converts echo context to params.
func (w *ServerInterfaceWrapper) ListRepos(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListRepos(ctx)
	return err
}

// CreateRepo converts echo context to params.
func (w *ServerInterfaceWrapper) CreateRepo(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateRepo(ctx)
	return err
}

// GetRepo converts echo context to params.
func (w *ServerInterfaceWrapper) GetRepo(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetRepo(ctx, id)
	return err
}

// ListStacks converts echo context to params.
func (w *ServerInterfaceWrapper) ListStacks(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListStacks(ctx)
	return err
}

// CreateStack converts echo context to params.
func (w *ServerInterfaceWrapper) CreateStack(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateStack(ctx)
	return err
}

// GetStack converts echo context to params.
func (w *ServerInterfaceWrapper) GetStack(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "slug" -------------
	var slug string

	err = runtime.BindStyledParameterWithLocation("simple", false, "slug", runtime.ParamLocationPath, ctx.Param("slug"), &slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter slug: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStack(ctx, slug)
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

	router.GET(baseURL+"/core/repos", wrapper.ListRepos)
	router.POST(baseURL+"/core/repos", wrapper.CreateRepo)
	router.GET(baseURL+"/core/repos/:id", wrapper.GetRepo)
	router.GET(baseURL+"/core/stacks", wrapper.ListStacks)
	router.POST(baseURL+"/core/stacks", wrapper.CreateStack)
	router.GET(baseURL+"/core/stacks/:slug", wrapper.GetStack)

}
