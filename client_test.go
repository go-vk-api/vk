package vk

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	testDefaultClient(client, t)
}

func TestNewClientWithOptions(t *testing.T) {
	token := "TOKEN"

	client, err := NewClientWithOptions(
		WithToken(token),
		WithHTTPClient(http.DefaultClient),
	)
	if err != nil {
		t.Error(err)
	}

	testDefaultClient(client, t)

	if client.Token != token {
		t.Errorf("client.Token == %q, want %q", client.Token, token)
	}

	if client.HTTPClient != http.DefaultClient {
		t.Errorf("client.HTTPClient == %v, want %v (http.DefaultClient)", client.HTTPClient, http.DefaultClient)
	}
}

func testDefaultClient(client *Client, t *testing.T) {
	if client.Lang != DefaultLang {
		t.Errorf("client.Lang == %q, want %q", client.Lang, DefaultLang)
	}

	if client.Version != DefaultVersion {
		t.Errorf("client.Version == %q, want %q", client.Version, DefaultVersion)
	}

	if client.BaseURL != DefaultBaseURL {
		t.Errorf("client.BaseURL == %q, want %q", client.BaseURL, DefaultBaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("client.HTTPClient == nil")
	}
}
