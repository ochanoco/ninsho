package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ochanoco/ninsho"
	// "ninsho"
)

func main() {
	var provider ninsho.Provider

	provider.ClientID = os.Getenv("NINSHO_CLIENT_ID")
	provider.ClientSecret = os.Getenv("NINSHO_CLIENT_SECRET")
	provider.RedirectUri = "http://127.0.0.1:8080/callback"
	provider.Scope = "profile openid"
	provider.UsePKCE = true

	n, err := ninsho.InitNinsho(&provider, &ninsho.LINE_LOGIN)
	if err != nil {
		panic(err)
	}

	authURL, _, _ := n.MakeAuthURL()
	fmt.Printf("Open this URL in your browser:\n%v\n\n", authURL)

	code := input("code")
	state := input("state")

	jwt, err := n.Auth(code, state)

	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %v", jwt)
}

func input(name string) string {
	fmt.Printf("Enter %v...:", name)

	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	line := strings.Replace(code, "\n", "", -1)

	return line
}
