package ninsho

var LINE_LOGIN = IdP[LINE_USER]{
	AuthURL:   "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=profile openid&nonce=%s",
	TokenURL:  "https://api.line.me/oauth2/v2.1/token",
	VerifyURL: "https://api.line.me/oauth2/v2.1/verify",
}

type LINE_USER struct {
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
