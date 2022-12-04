package line_login_core

import (
	"testing"
)

func sampleSession() Session {
	var provider Provider

	provider.ClientID = "12345"
	provider.ClientSecret = ""
	provider.RedirectUri = "http://127.0.0.1:8080/callback"

	session, err := NewSession(&provider)
	if err != nil {
		panic(err)
	}

	session.State = "aaa"
	session.Nonce = "bbb"

	return session
}

func TestAuthURL(t *testing.T) {
	var expected = "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=12345&redirect_uri=http://127.0.0.1:8080/callback&state=aaa&scope=profile openid&nonce=bbb"

	session := sampleSession()
	url := session.AuthURL()

	if url != expected {
		t.Fatalf("Auth URL is not collect\nexpected: %v\nactual:   %v", url, expected)
	}
}

func TestGetUser(t *testing.T) {
	// todo
}
