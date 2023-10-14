package ninsho

import (
	"testing"
)

func sampleNinsho() Ninsho {
	var provider Provider

	provider.ClientID = "12345"
	provider.ClientSecret = ""
	provider.RedirectUri = "http://127.0.0.1:8080/callback"

	ninsho, err := NewNinsho(&provider)
	if err != nil {
		panic(err)
	}

	ninsho.State = "aaa"
	ninsho.Nonce = "bbb"

	return ninsho
}

func TestAuthURL(t *testing.T) {
	var expected = "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=12345&redirect_uri=http://127.0.0.1:8080/callback&state=aaa&scope=profile openid&nonce=bbb"

	ninsho := sampleNinsho()
	url := ninsho.AuthURL()

	if url != expected {
		t.Fatalf("Auth URL is not collect\nexpected: %v\nactual:   %v", url, expected)
	}
}

func TestGetUser(t *testing.T) {
	// todo
}
