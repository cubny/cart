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

// AccessKey contains the key of a user
type AccessKey struct {
	ID     int64
	Key    string
	UserID int64
}

// New stubs an actual auth server and populates an list of valid
// access keys for 3 sample users
func New() *Client {
	accessKeys := []AccessKey{
		{
			ID:     1,
			Key:    "abcdef123456",
			UserID: 1,
		},
		{
			ID:     2,
			Key:    "bcdefg123456",
			UserID: 12,
		},
		{
			ID:     3,
			Key:    "cdefgh123456",
			UserID: 20,
		},
	}
	return &Client{accessKeys: accessKeys}
}

// KeyFromClientRequest gets the Key from the headers of a request in a standardized way. Empty if not found.
func KeyFromClientRequest(req *http.Request) string {
	parts := strings.Fields(req.Header.Get("Authorisation"))
	if len(parts) != 2 || parts[0] != "Key" {
		return ""
	}

	return parts[1]
}

//  VerifyKey verifies an Key.
func (c *Client) VerifyKey(ctx context.Context, key string) (*AccessKey, error) {
	// mocking the auth client
	// verifying key should be done against an actual auth service using http or grpc
	// here we only mock the service to keep things simple for the assignment
	for _, ak := range c.accessKeys {
		if ak.Key == key {
			return &ak, nil
		}
	}
	return nil, ErrNotFound
}

// AddKeyToRequest adds the Key to the headers of a request in a standardized way.
func AddKeyToRequest(req *http.Request, token string) {
	req.Header.Set("Authorisation", fmt.Sprintf("Key %v", token))
}
