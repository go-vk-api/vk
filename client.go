package vk

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-vk-api/vk/httputil"
	"github.com/pkg/errors"
)

const (
	DefaultBaseURL = "https://api.vk.com/method"
	DefaultLang    = "en"
	DefaultVersion = "5.103"
)

// Client manages communication with VK API.
type Client struct {
	BaseURL string
	Lang    string
	Version string

	Token string

	HTTPClient httputil.RequestDoer
}

// CallMethod invokes the named method and stores the result in the value pointed to by response.
func (c *Client) CallMethod(method string, params RequestParams, response interface{}) error {
	queryParams, err := params.URLValues()
	if err != nil {
		return err
	}

	setIfEmpty := func(param, value string) {
		if queryParams.Get(param) == "" {
			queryParams.Set(param, value)
		}
	}

	setIfEmpty("v", c.Version)
	setIfEmpty("lang", c.Lang)
	if c.Token != "" {
		setIfEmpty("access_token", c.Token)
	}

	rawBody, err := httputil.Post(c.HTTPClient, c.BaseURL+"/"+method, queryParams)
	if err != nil {
		return err
	}

	var body struct {
		Response interface{}  `json:"response"`
		Error    *MethodError `json:"error"`
	}

	if response != nil {
		valueOfResponse := reflect.ValueOf(response)
		if valueOfResponse.Kind() != reflect.Ptr || valueOfResponse.IsNil() {
			return errors.New("response must be a valid pointer")
		}

		body.Response = response
	}

	if err = json.Unmarshal(rawBody, &body); err != nil {
		return err
	}

	if body.Error != nil {
		return body.Error
	}

	return nil
}

// NewClient initializes a new VK API client with default values.
func NewClient() (*Client, error) {
	return NewClientWithOptions()
}

// NewClientWithOptions initializes a new VK API client with default values. It takes functors
// to modify values when creating it, like `NewClientWithOptions(WithToken(…))`.
func NewClientWithOptions(options ...Option) (*Client, error) {
	client := &Client{
		BaseURL: DefaultBaseURL,
		Lang:    DefaultLang,
		Version: DefaultVersion,
	}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}

	if client.HTTPClient == nil {
		client.HTTPClient = http.DefaultClient
	}

	return client, nil
}
