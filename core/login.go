package line_login_core

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
const TOKEN_URL = "https://api.line.me/oauth2/v2.1/token"

const TOKEN_LEN = 32

type Provider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type Session struct {
	State    string
	Nonce    string
	provider *Provider
}

type User struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func secureRandom(b int) (string, error) {
	k := make([]byte, b)
	if _, err := rand.Read(k); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", k), nil
}

func NewSession(provider *Provider) (Session, error) {
	var session Session
	var err error

	session.Nonce, err = secureRandom(TOKEN_LEN)

	if err != nil {
		return session, err
	}

	session.State, err = secureRandom(TOKEN_LEN)

	if err != nil {
		return session, err
	}

	session.provider = provider

	return session, nil
}

func (session *Session) AuthURL() string {
	return fmt.Sprintf(AUTH_URL, session.provider.ClientID, session.provider.RedirectURL, session.State, session.Nonce)
}

func (session *Session) GetUser(code string) (*User, error) {
	var user User

	provider := session.provider

	values := url.Values{}

	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", provider.ClientID)
	values.Add("client_secret", provider.ClientSecret)
	values.Add("redirect_uri", provider.RedirectURL)

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

	err = json.Unmarshal(b, &user)

	if err != nil {
		panic(err)
	}

	return &user, nil
}
