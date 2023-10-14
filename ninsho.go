package ninsho

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var TOKEN_LEN = 32

type IdP[T any] struct {
	AuthURL   string
	TokenURL  string
	VerifyURL string
}

type Provider struct {
	ClientID     string
	ClientSecret string
	RedirectUri  string
}

type Ninsho[T any] struct {
	State    string
	Nonce    string
	Provider *Provider
	IdP      *IdP[T]
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func NewNinsho[T any](provider *Provider, idp *IdP[T]) (Ninsho[T], error) {
	var _ninsho Ninsho[T]
	var err error

	_ninsho.Nonce, err = secureRandom(TOKEN_LEN)

	if err != nil {
		return _ninsho, err
	}

	_ninsho.State, err = secureRandom(TOKEN_LEN)

	if err != nil {
		return _ninsho, err
	}

	_ninsho.Provider = provider
	_ninsho.IdP = idp

	return _ninsho, nil
}

func (_ninsho *Ninsho[T]) GetAuthURL() string {
	return fmt.Sprintf(_ninsho.IdP.AuthURL, _ninsho.Provider.ClientID, _ninsho.Provider.RedirectUri, _ninsho.State, _ninsho.Nonce)
}

func (_ninsho *Ninsho[T]) Auth(code string) (*T, error) {
	var jwt T
	var token Token

	provider := _ninsho.Provider

	values := url.Values{}

	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", provider.ClientID)
	values.Add("client_secret", provider.ClientSecret)
	values.Add("redirect_uri", provider.RedirectUri)

	req, err := http.NewRequest(
		"POST",
		_ninsho.IdP.TokenURL,
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &token)

	if err != nil {
		panic(err)
	}

	values = url.Values{}

	values.Add("id_token", token.IdToken)
	values.Add("client_id", provider.ClientID)

	req, err = http.NewRequest(
		"POST",
		_ninsho.IdP.VerifyURL,
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &jwt)

	if err != nil {
		panic(err)
	}

	return &jwt, nil
}
