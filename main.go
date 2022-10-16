package main

import (
	"fmt"
	"os"
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

	url := fmt.Sprintf(AuthURL, config.ClientID, config.RedirectURL, state, nonce)
	fmt.Println(url)
}
