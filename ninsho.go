package ninsho

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const AUTH_URL = "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=profile openid&nonce=%s"

var TOKEN_URL = "https://api.line.me/oauth2/v2.1/token"
var VERIFY_URL = "https://api.line.me/oauth2/v2.1/verify"

var TOKEN_LEN = 32

type Provider struct {
	ClientID     string
	ClientSecret string
	RedirectUri  string
}

type Ninsho struct {
	State    string
	Nonce    string
	Provider *Provider
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type JWT struct {
	Iss     string   `json:"iss"`
	Sub     string   `json:"sub"`
	Aud     string   `json:"aud"`
	Exp     int      `json:"exp"`
	Iat     int      `json:"iat"`
	Nonce   string   `json:"nonce"`
	Amr     []string `json:"amr"`
	Name    string   `json:"name"`
	Picture string   `json:"picture"`
	Email   string   `json:"email"`
}

func secureRandom(b int) (string, error) {
	k := make([]byte, b)
	if _, err := rand.Read(k); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", k), nil
}

func NewNinsho(provider *Provider) (Ninsho, error) {
	var _ninsho Ninsho
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

	return _ninsho, nil
}

func (_ninsho *Ninsho) AuthURL() string {
	return fmt.Sprintf(AUTH_URL, _ninsho.Provider.ClientID, _ninsho.Provider.RedirectUri, _ninsho.State, _ninsho.Nonce)
}

func (_ninsho *Ninsho) GetUser(code string) (*JWT, error) {
	var jwt JWT
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
		TOKEN_URL,
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
		VERIFY_URL,
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
