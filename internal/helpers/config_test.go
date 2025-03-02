package helpers

import "testing"

func TestLoadConfig(t *testing.T) {
	c, err := LoadAppConfig(".", "test_config")
	if err != nil {
		t.Fatalf("config load failed: %v", err)
	}
	if c.OAuth.ClientID != "local_auth_manager" {
		t.Fatalf("failed to read client_id: expect: 'local_auth_manager', got: '%v'", c.OAuth.ClientID)
	}
	if c.OAuth.ClientSecret != "test" {
		t.Fatalf("failed to read client_id: expect: 'test', got: '%v'", c.OAuth.ClientSecret)
	}
	if c.OAuth.Authority != "test" {
		t.Fatalf("failed to read authority: expect: 'test', got: '%v'", c.OAuth.Authority)
	}
	if c.OAuth.CallbackPath != "/auth/login/callback" {
		t.Fatalf("failed to read callback_path: expect: '/auth/login/callback', got: '%v'", c.OAuth.CallbackPath)
	}
	if c.Session.SessionKey != "test" {
		t.Fatalf("failed to read session_key: expect: 'test', got: '%v'", c.Session.SessionKey)
	}
}
