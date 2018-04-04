// Copyright 2017 Bankrs AG.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package bosgo provides a Go client for accessing the Bankrs OS API.
package bosgo

import (
	"context"
	"encoding/json"
	"net/http"
)

const (

	// Version is the current version of the bosgo library.
	Version = "0.1"

	// DefaultUserAgent is the default user agent header used by the bosgo library.
	DefaultUserAgent = "bosgo-bankrs-os-client/" + Version

	apiV1 = "/v1"
)

const (
	SandboxAddr    = "api.sandbox.bankrs.com"
	ProductionAddr = "api.bankrs.com"
)

// Client is the base client used for interacting with services that do not
// require authentication. Use Login to initiate a developer session.  It is
// safe for concurrent use by multiple goroutines.
type Client struct {
	// never modified once they have been set
	hc          *http.Client
	addr        string
	ua          string
	environment string
}

type ClientOption func(*Client)

// New creates a new client that will use the supplied HTTP client and connect
// via the specified API host address.
func New(client *http.Client, addr string, opts ...ClientOption) *Client {
	c := &Client{
		hc:   client,
		addr: addr,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) newReq(path string) req {
	return req{
		hc:   c.hc,
		addr: c.addr,
		path: path,
		headers: headers{
			"User-Agent": c.userAgent(),
		},
		par:         params{},
		environment: c.environment,
	}
}

func (c *Client) userAgent() string {
	if c.ua == "" {
		return DefaultUserAgent
	}

	return DefaultUserAgent + " " + c.ua
}

func (c *Client) WithApplicationID(applicationID string) *AppClient {
	ac := NewAppClient(c.hc, c.addr, applicationID)
	ac.ua = c.ua
	ac.environment = c.environment
	return ac
}

// Login prepares and returns a request to log a developer into the Bankrs
// API. Sending a successful request will return a new client that allows
// access to services requiring a valid developer session.
func (c *Client) Login(email, password string) *DeveloperLoginReq {
	return &DeveloperLoginReq{
		client: c,
		req:    c.newReq(apiV1 + "/developers/login"),
		data: DeveloperCredentials{
			Email:    email,
			Password: password,
		},
	}
}

type DeveloperLoginReq struct {
	req
	client *Client
	data   DeveloperCredentials
}

type sessionToken struct {
	Token string `json:"token"`
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *DeveloperLoginReq) Context(ctx context.Context) *DeveloperLoginReq {
	r.req.ctx = ctx
	return r
}

// ClientID sets a client identifier that will be passed to the Bankrs API in
// the X-Client-Id header.
func (r *DeveloperLoginReq) ClientID(id string) *DeveloperLoginReq {
	r.req.clientID = id
	return r
}

// Send sends the login request and returns a client that can be used to
// access services within the developer's session.
func (r *DeveloperLoginReq) Send() (*DevClient, error) {
	res, cleanup, err := r.req.postJSON(&r.data)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var t sessionToken
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, decodeError(err, res)
	}

	dc := NewDevClient(r.client.hc, r.client.addr, t.Token)
	dc.ua = r.client.ua
	dc.environment = r.client.environment
	return dc, nil
}

// CreateDeveloper prepares and returns a request to create a developer account for the
// Bankrs API. Sending a successful request will return a new client that
// allows access to services requiring a valid developer session.
func (c *Client) CreateDeveloper(email, password string) *DeveloperCreateReq {
	return &DeveloperCreateReq{
		client: c,
		req:    c.newReq(apiV1 + "/developers"),
		data: DeveloperCredentials{
			Email:    email,
			Password: password,
		},
	}
}

type DeveloperCreateReq struct {
	req
	client *Client
	data   DeveloperCredentials
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *DeveloperCreateReq) Context(ctx context.Context) *DeveloperCreateReq {
	r.req.ctx = ctx
	return r
}

// ClientID sets a client identifier that will be passed to the Bankrs API in
// the X-Client-Id header.
func (r *DeveloperCreateReq) ClientID(id string) *DeveloperCreateReq {
	r.req.clientID = id
	return r
}

// Send sends the create request and returns a client that can be used to
// access services within the developer's session.
func (r *DeveloperCreateReq) Send() (*DevClient, error) {
	res, cleanup, err := r.req.postJSON(&r.data)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var t sessionToken
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, decodeError(err, res)
	}

	dc := NewDevClient(r.client.hc, r.client.addr, t.Token)
	dc.ua = r.client.ua
	dc.environment = r.client.environment
	return dc, nil

}

// LostPassword prepares and returns a request to start the lost password process.
func (c *Client) LostPassword(email string) *LostPasswordReq {
	return &LostPasswordReq{
		req: c.newReq(apiV1 + "/developers/lost_password"),
		data: developerEmail{
			Email: email,
		},
	}
}

type developerEmail struct {
	Email string `json:"email"`
}

type LostPasswordReq struct {
	req
	data developerEmail
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *LostPasswordReq) Context(ctx context.Context) *LostPasswordReq {
	r.req.ctx = ctx
	return r
}

// ClientID sets a client identifier that will be passed to the Bankrs API in
// the X-Client-Id header.
func (r *LostPasswordReq) ClientID(id string) *LostPasswordReq {
	r.req.clientID = id
	return r
}

// Send sends the lost password request.
func (r *LostPasswordReq) Send() error {
	_, cleanup, err := r.req.postJSON(&r.data)
	defer cleanup()
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword prepares and returns a request to reset a lost password.
func (c *Client) ResetPassword(password string, token string) *ResetPasswordReq {
	return &ResetPasswordReq{
		req: c.newReq(apiV1 + "/developers/reset_password"),
		data: developerPasswordReset{
			Password: password,
			Token:    token,
		},
	}
}

type developerPasswordReset struct {
	Password string `json:"email"`
	Token    string `json:"token"`
}

type ResetPasswordReq struct {
	req
	data developerPasswordReset
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *ResetPasswordReq) Context(ctx context.Context) *ResetPasswordReq {
	r.req.ctx = ctx
	return r
}

// ClientID sets a client identifier that will be passed to the Bankrs API in
// the X-Client-Id header.
func (r *ResetPasswordReq) ClientID(id string) *ResetPasswordReq {
	r.req.clientID = id
	return r
}

// Send sends the reset password request.
func (r *ResetPasswordReq) Send() error {
	_, cleanup, err := r.req.postJSON(&r.data)
	defer cleanup()
	if err != nil {
		return err
	}

	return nil
}

// UserAgent is a client option that may be used to add information to the user agent header used by
// the client.
func UserAgent(ua string) ClientOption {
	return func(c *Client) { c.ua = ua }
}

// Environment is a client option that may be used to set the X-Environment header used by
// the client.
func Environment(environment string) ClientOption {
	return func(c *Client) { c.environment = environment }
}
