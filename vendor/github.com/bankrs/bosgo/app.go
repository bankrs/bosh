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

package bosgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AppClient is a client used for interacting with services in the context of
// a registered application and a valid user or developer session. It is safe
// for concurrent use by multiple goroutines.
type AppClient struct {
	// never modified once they have been set
	hc            *http.Client
	addr          string
	token         string // session token
	applicationID string
	ua            string

	Categories *CategoriesService
	Providers  *ProvidersService
	Users      *AppUsersService
}

// NewAppClient creates a new client that may be used to interact with
// services that require a specific application context.
func NewAppClient(client *http.Client, addr string, token string, applicationID string) *AppClient {
	ac := &AppClient{
		hc:            client,
		addr:          addr,
		token:         token,
		applicationID: applicationID,
	}

	ac.Categories = NewCategoriesService(ac)
	ac.Providers = NewProvidersService(ac)
	ac.Users = NewAppUsersService(ac)
	return ac
}

func (a *AppClient) newReq(path string) req {
	return req{
		hc:   a.hc,
		addr: a.addr,
		path: path,
		headers: headers{
			"User-Agent":       a.userAgent(),
			"x-token":          a.token,
			"x-application-id": a.applicationID,
		},
		par: params{},
	}
}

func (a *AppClient) userAgent() string {
	if a.ua == "" {
		return UserAgent
	}

	return UserAgent + " " + a.ua
}

// CategoriesService provides access to category related API services.
type CategoriesService struct {
	client *AppClient
}

func NewCategoriesService(a *AppClient) *CategoriesService { return &CategoriesService{client: a} }

// List returns a request that may be used to request a list of classification categories.
func (c *CategoriesService) List() *CategoriesReq {
	return &CategoriesReq{
		req: c.client.newReq(apiV1 + "/categories"),
	}
}

type CategoriesReq struct {
	req
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *CategoriesReq) Context(ctx context.Context) *CategoriesReq {
	r.req.ctx = ctx
	return r
}

// Send sends the request to list categories.
func (r *CategoriesReq) Send() (*CategoryList, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var list CategoryList
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, wrap(errContextInvalidServiceResponse, err)
	}

	return &list, nil
}

type CategoryList []Category

type Category struct {
	ID    int64             `json:"id"`
	Names map[string]string `json:"names"`
	Group string            `json:"group"` // spending or income, e.g.
}

// ProvidersService provides access to financial provider related API services.
type ProvidersService struct {
	client *AppClient
}

func NewProvidersService(a *AppClient) *ProvidersService { return &ProvidersService{client: a} }

// Search returns a request that may be used to search the list of financial providers.
func (c *ProvidersService) Search(query string) *ProvidersSearchReq {
	r := c.client.newReq(apiV1 + "/providers")
	r.par.Set("q", query)
	return &ProvidersSearchReq{
		req: r,
	}
}

type ProvidersSearchReq struct {
	req
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *ProvidersSearchReq) Context(ctx context.Context) *ProvidersSearchReq {
	r.req.ctx = ctx
	return r
}

// Send sends the request to search providers.
func (r *ProvidersSearchReq) Send() (*ProviderSearchResults, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var srch ProviderSearchResults
	if err := json.NewDecoder(res.Body).Decode(&srch); err != nil {
		return nil, wrap(errContextInvalidServiceResponse, err)
	}

	return &srch, nil
}

// Get returns a request that may be used to get the details of a single financial provider.
func (c *ProvidersService) Get(id string) *ProvidersGetReq {
	return &ProvidersGetReq{
		req: c.client.newReq(apiV1 + "/providers/" + id),
	}
}

type ProvidersGetReq struct {
	req
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *ProvidersGetReq) Context(ctx context.Context) *ProvidersGetReq {
	r.req.ctx = ctx
	return r
}

// Send sends the request to get a single financial provider.
func (r *ProvidersGetReq) Send() (*Provider, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var p Provider
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, wrap(errContextInvalidServiceResponse, err)
	}

	return &p, nil
}

type ProviderSearchResults []ProviderSearchResult

type ProviderSearchResult struct {
	Score    float64  `json:"score"`
	Provider Provider `json:"provider"`
}

type Provider struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Country     string          `json:"country"`
	URL         string          `json:"url"`
	Address     string          `json:"address"`
	PostalCode  string          `json:"postal_code"`
	Challenges  []ChallengeSpec `json:"challenges"`
}

type ChallengeSpec struct {
	ID          string            `json:"id"`
	Description string            `json:"desc"`
	Type        string            `json:"type"`
	Secure      bool              `json:"secure"`
	UnStoreable bool              `json:"unstoreable"`
	Options     map[string]string `json:"options,omitempty"`
}

// AppUsersService provides access to application user related API services.
type AppUsersService struct {
	client *AppClient
}

func NewAppUsersService(c *AppClient) *AppUsersService { return &AppUsersService{client: c} }

func (a *AppUsersService) List() *ListDevUsersReq {
	return &ListDevUsersReq{
		req: a.client.newReq(apiV1 + "/developers/users"),
	}
}

type ListDevUsersReq struct {
	req
	data PageParams
}

type PageParams struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

func (r *ListDevUsersReq) Context(ctx context.Context) *ListDevUsersReq {
	r.req.ctx = ctx
	return r
}

func (r *ListDevUsersReq) Cursor(cursor string) *ListDevUsersReq {
	r.data.Cursor = cursor
	return r
}

func (r *ListDevUsersReq) Limit(v int) *ListDevUsersReq {
	r.data.Limit = v
	return r
}

func (r *ListDevUsersReq) Send() (*UserListPage, error) {
	if r.data.Limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative")
	}

	var res *http.Response
	var cleanup func()
	var err error
	if r.data.Limit == 0 {
		res, cleanup, err = r.req.get()
	} else {
		res, cleanup, err = r.req.postJSON(r.data)
	}
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var list UserListPage
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, wrap(errContextInvalidServiceResponse, err)
	}
	return &list, nil
}

type UserListPage struct {
	Users      []string `json:"users,omitempty"`
	NextCursor string   `json:"next,omitempty"`
}

// Create returns a request that may be used to create a user with the given username and password.
func (a *AppUsersService) Create(username, password string) *UserCreateReq {
	return &UserCreateReq{
		req:    a.client.newReq(apiV1 + "/users"),
		client: a.client,
		data: UserCredentials{
			Username: username,
			Password: password,
		},
	}
}

// UserCreateReq is a request that may be used to create a user.
type UserCreateReq struct {
	req
	client *AppClient
	data   UserCredentials
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *UserCreateReq) Context(ctx context.Context) *UserCreateReq {
	r.req.ctx = ctx
	return r
}

// Send sends the request to create the user and returns a client that can be
// used to access services within the new users's session.
func (r *UserCreateReq) Send() (*UserClient, error) {
	res, cleanup, err := r.req.postJSON(r.data)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var t UserToken
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, wrap(errContextInvalidServiceResponse, err)
	}

	uc := NewUserClient(r.client.hc, r.client.addr, t.Token, r.client.applicationID)
	uc.ua = r.client.ua
	return uc, nil
}

type UserToken struct {
	ID    string `json:"id"`    // globally unique identifier for a user
	Token string `json:"token"` // session token
}

// Login returns a request that may be used to login a user with the given username and password.
func (a *AppUsersService) Login(username, password string) *UserLoginReq {
	return &UserLoginReq{
		req:    a.client.newReq(apiV1 + "/users/login"),
		client: a.client,
		data: UserCredentials{
			Username: username,
			Password: password,
		},
	}
}

type UserLoginReq struct {
	req
	client *AppClient
	data   UserCredentials
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *UserLoginReq) Context(ctx context.Context) *UserLoginReq {
	r.req.ctx = ctx
	return r
}

// Send sends the request to login the user and returns a client that can be
// used to access services within the new users's session.
func (r *UserLoginReq) Send() (*UserClient, error) {
	res, cleanup, err := r.req.postJSON(r.data)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var t UserToken
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, wrap(errContextInvalidServiceResponse, err)
	}

	uc := NewUserClient(r.client.hc, r.client.addr, t.Token, r.client.applicationID)
	uc.ua = r.client.ua
	return uc, nil
}

// ResetPassword prepares and returns a request to reset a user's password.
func (a *AppUsersService) ResetPassword(username, password string) *ResetUserPasswordReq {
	return &ResetUserPasswordReq{
		req: a.client.newReq(apiV1 + "/users/reset_password"),
		data: UserCredentials{
			Username: username,
			Password: password,
		},
	}
}

type ResetUserPasswordReq struct {
	req
	data UserCredentials
}

func (r *ResetUserPasswordReq) Context(ctx context.Context) *ResetUserPasswordReq {
	r.req.ctx = ctx
	return r
}

// Send sends the request to retrieve the developer's profile.
func (r *ResetUserPasswordReq) Send() error {
	_, cleanup, err := r.req.postJSON(r.data)
	defer cleanup()
	if err != nil {
		return err
	}
	return nil
}
