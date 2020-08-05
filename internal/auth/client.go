package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrNotFound = errors.New("accessKey not found")

// Client is the client for Auth service, usually we enquiry this service
// using http or grpc, but here we mock the service with a list for the
//  assignment and to keep things simple
type Client struct {
	accessKeys []AccessKey
}

type AccessKey struct {
	ID        int64
	AccessKey string
	UserID    int64
}

func New() *Client {
	accessKeys := []AccessKey{
		{
			ID:        1,
			AccessKey: "abcdef123456",
			UserID:    1,
		},
		{
			ID:        2,
			AccessKey: "bcdefg123456",
			UserID:    12,
		},
		{
			ID:        3,
			AccessKey: "cdefgh123456",
			UserID:    20,
		},
	}
	return &Client{accessKeys: accessKeys}
}

// AccessKeyFromClientRequest gets the AccessKey from the headers of a request in a standardized way. Empty if not found.
func AccessKeyFromClientRequest(req *http.Request) string {
	parts := strings.Fields(req.Header.Get("Authorization"))
	if len(parts) != 2 || parts[0] != "AccessKey" {
		return ""
	}

	return parts[1]
}

//  VerifyAccessKey verifies an AccessKey.
func (c *Client) VerifyAccessKey(ctx context.Context, accessKey string) (*AccessKey, error) {
	// mocking the auth client
	// verifying accessKey should be done against an actual auth service using http or grpc
	// here we only mock the service to keep things simple for the assignment
	for _, ak := range c.accessKeys {
		if ak.AccessKey == accessKey {
			return &ak, nil
		}
	}
	return nil, ErrNotFound
}

// AddAccessKeyToRequest adds the AccessKey to the headers of a request in a standardized way.
func AddAccessKeyToRequest(req *http.Request, accessKey string) {
	req.Header.Set("Authorization", fmt.Sprintf("AccessKey %v", accessKey))
}
