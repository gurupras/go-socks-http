package sockshttp

import (
	"net/http"
	"os"
	"testing"
)

func TestNew_DefaultNoEnvSet(t *testing.T) {
	// Unset the default env var
	os.Unsetenv(DefaultEnvVar)

	client, err := New("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	// Should be default transport
	if client.Transport != nil {
		t.Errorf("expected default transport (nil), got %T", client.Transport)
	}
}

func TestNew_CustomEnvVar(t *testing.T) {
	const envVar = "CUSTOM_SOCKS_PROXY"
	os.Unsetenv(envVar)
	client, err := New(envVar)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.Transport != nil {
		t.Errorf("expected default transport (nil), got %T", client.Transport)
	}
}

func TestNew_InvalidProxyURL(t *testing.T) {
	const envVar = "SOCKS_PROXY_URL"
	os.Setenv(envVar, "%%%bad-url%%%")
	defer os.Unsetenv(envVar)

	client, err := New(envVar)
	if err == nil {
		t.Fatal("expected error for invalid proxy URL, got nil")
	}
	if client != nil {
		t.Errorf("expected nil client on error, got %#v", client)
	}
}

func TestNew_ValidProxyURL(t *testing.T) {
	const envVar = "SOCKS_PROXY_URL"
	// This is a syntactically valid address (localhost:1080), but we don't actually connect.
	os.Setenv(envVar, "socks5://127.0.0.1:1080")
	defer os.Unsetenv(envVar)

	client, err := New(envVar)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	// Should have a custom transport
	if client.Transport == nil {
		t.Errorf("expected custom transport, got nil")
	}
	_, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Errorf("expected *http.Transport, got %T", client.Transport)
	}
}
