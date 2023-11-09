package ninsho

import (
	"testing"
)

func sampleNinsho() Ninsho[LINE_USER] {
	var provider Provider

	provider.ClientID = "12345"
	provider.ClientSecret = ""
	provider.RedirectUri = "http://127.0.0.1:8080/callback"
	provider.Scope = "profile openid"
	provider.UsePKCE = false

	ninsho, err := NewNinsho(&provider, &LINE_LOGIN)
	if err != nil {
		panic(err)
	}

	ninsho.State = "aaa"
	ninsho.Nonce = "bbb"

	return ninsho
}

func TestAuthURL(t *testing.T) {
	expected := "https://access.line.me/oauth2/v2.1/authorize?client_id=12345&nonce=bbb&redirect_uri=http%3A%2F%2F127.0.0.1%3A8080%2Fcallback&response_type=code&scope=profile+openid&state=aaa"

	ninsho := sampleNinsho()
	url := ninsho.GetAuthURL()

	if url != expected {
		t.Fatalf("Auth URL is not collect\nexpected: %v\nactual:   %v", expected, url)
	}
}

func TestGetUser(t *testing.T) {
	// todo
}
