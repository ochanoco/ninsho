package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const AuthURL = "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=profile openid&nonce=%s"
const TokenURL = "https://api.line.me/oauth2/v2.1/token"

// const CallBackURL = "http://%s/callback?code=%s&state=%s&friendship_status_changed=true"
// const CallBackErrURL = "https://%s/callback?error=access_denied&error_description=The+resource+owner+denied+the+request.&state=%s"

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func main() {
	var config Config
	config.ClientID = os.Getenv("CLIENT_ID")
	config.ClientSecret = os.Getenv("TOKEN")
	config.RedirectURL = "http://localhost:3000/api/auth/callback/line"

	state := "12345"
	nonce := "09876"

	authUrl := fmt.Sprintf(AuthURL, config.ClientID, config.RedirectURL, state, nonce)
	fmt.Println(authUrl)

	reader := bufio.NewReader(os.Stdin)
	code, _ := reader.ReadString('\n')

	code = strings.Replace(code, "\n", "", -1)

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", config.ClientID)
	values.Add("client_secret", config.ClientSecret)
	values.Add("redirect_uri", config.RedirectURL)

	req, err := http.NewRequest(
		"POST",
		TokenURL,
		strings.NewReader(values.Encode()),
	)
	
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

}
