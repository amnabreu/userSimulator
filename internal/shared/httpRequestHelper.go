package shared

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GithubEvent     = "github_event"
	UserAccessToken = "access_token"
	URL             = "url"
)

// TODO: add array of valid responses
func ValidateHttpResponse(response *http.Response, err error) *http.Response {
	if err != nil {
		fmt.Printf("Request failed with error: %v", err)
	}

	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		fmt.Printf("Error: complete Installation requested failed with status: %d", response.StatusCode)
	}

	return response
}

func AddAuthHeader(ctx context.Context, req *http.Request) error {

	req.Header.Set("authorization", "Token "+ctx.Value(UserAccessToken).(string))
	return nil
}
