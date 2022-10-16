package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const AUTH_URL = "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=profile openid&nonce=%s"
const TOKEN_URL = "https://api.line.me/oauth2/v2.1/token"

// const CallBackURL = "http://%s/callback?code=%s&state=%s&friendship_status_changed=true"
// const CallBackErrURL = "https://%s/callback?error=access_denied&error_description=The+resource+owner+denied+the+request.&state=%s"

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

func secureRandom(b int) string {
	k := make([]byte, b)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", k)
}

func NewSession(provider *Provider) Session {
	var session Session

	session.Nonce = secureRandom(32)
	session.State = secureRandom(32)

	session.provider = provider

	return session
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

func main() {
	var provider Provider

	provider.ClientID = os.Getenv("CLIENT_ID")
	provider.ClientSecret = os.Getenv("TOKEN")
	provider.RedirectURL = "http://localhost:3000/api/auth/callback/line"

	session := NewSession(&provider)
	authURL := session.AuthURL()

	fmt.Printf("Open this URL in your browser:\n%v\n\nEnter Code:", authURL)

	reader := bufio.NewReader(os.Stdin)
	code, _ := reader.ReadString('\n')
	code = strings.Replace(code, "\n", "", -1)

	user, err := session.GetUser(code)

	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %v", user)
}
