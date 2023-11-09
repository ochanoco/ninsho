package ninsho

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var DEFAULT_PKCE_METHOD = "S256"

var TOKEN_LEN = 64

type User interface{}

type IdP[T User] struct {
	AuthURL   string
	TokenURL  string
	VerifyURL string
}

type Provider struct {
	ClientID     string
	ClientSecret string
	RedirectUri  string
	Scope        string
	UsePKCE      bool
}

type PKCEAuth struct {
	CodeChallenge       string
	CodeVerifier        string
	CodeChallengeMethod string
}

func NewPKCEAuth() *PKCEAuth {
	var pkce PKCEAuth
	return &pkce
}

func RandomPKCE() (*PKCEAuth, error) {
	pkce := NewPKCEAuth()

	codeVerifier, err := secureRandom(TOKEN_LEN)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write([]byte(codeVerifier))
	s := h.Sum(nil)

	codeChallenge := base64.RawURLEncoding.EncodeToString(s)

	pkce.CodeVerifier = codeVerifier
	pkce.CodeChallenge = codeChallenge
	pkce.CodeChallengeMethod = DEFAULT_PKCE_METHOD

	return pkce, nil
}

func InitPKCEAuth(usePKCE bool) (*PKCEAuth, error) {
	if usePKCE {
		return RandomPKCE()
	} else {
		return NewPKCEAuth(), nil
	}
}

type Session struct {
	State    string
	PkceAuth *PKCEAuth
}

type Ninsho[T any] struct {
	Provider *Provider
	IdP      *IdP[T]
	Session  *Session
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func RandomSession(usePKCE bool) (*Session, error) {
	var session Session
	var err error

	session.State, err = secureRandom(TOKEN_LEN)

	if err != nil {
		return nil, err
	}

	session.PkceAuth, err = InitPKCEAuth(usePKCE)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func NewNinsho[T User](provider *Provider, idp *IdP[T], session *Session) *Ninsho[T] {
	var _ninsho Ninsho[T]

	_ninsho.Provider = provider
	_ninsho.IdP = idp
	_ninsho.Session = session

	return &_ninsho
}

func InitNinsho[T User](provider *Provider, idp *IdP[T]) (*Ninsho[T], error) {
	session, err := RandomSession(provider.UsePKCE)

	if err != nil {
		return nil, err
	}

	_ninsho := NewNinsho[T](provider, idp, session)

	return _ninsho, nil
}

func (_ninsho *Ninsho[T]) MakeAuthURL() (string, string, error) {
	nonce, err := secureRandom(TOKEN_LEN)

	if err != nil {
		return "", "", err
	}

	values := url.Values{}
	values.Add("response_type", "code")
	values.Add("client_id", _ninsho.Provider.ClientID)
	values.Add("redirect_uri", _ninsho.Provider.RedirectUri)
	values.Add("nonce", nonce)
	values.Add("state", _ninsho.Session.State)
	values.Add("scope", _ninsho.Provider.Scope)

	if _ninsho.Provider.UsePKCE {
		values.Add("code_challenge", _ninsho.Session.PkceAuth.CodeChallenge)
		values.Add("code_challenge_method", _ninsho.Session.PkceAuth.CodeChallengeMethod)
	}

	url := _ninsho.IdP.AuthURL + "?" + values.Encode()
	return url, nonce, nil
}

func (_ninsho *Ninsho[T]) Auth(code string, state string) (*T, error) {
	var jwt T
	var token Token

	provider := _ninsho.Provider

	if state != _ninsho.Session.State {
		return nil, fmt.Errorf("state is wrong: %v", state)
	}

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", provider.ClientID)
	values.Add("client_secret", provider.ClientSecret)
	values.Add("redirect_uri", provider.RedirectUri)

	if _ninsho.Provider.UsePKCE {
		values.Add("code_verifier", _ninsho.Session.PkceAuth.CodeVerifier)
	}

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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to authentiction #1: (%v)", resp.StatusCode)
	}

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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to authentiction #1: (%v)", resp.StatusCode)
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &jwt)

	if err != nil {
		return nil, err
	}

	return &jwt, nil
}
