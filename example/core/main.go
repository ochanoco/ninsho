package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ochanoco/ninsho"
)

func main() {
	var provider ninsho.Provider

	provider.ClientID = os.Getenv("CLIENT_ID")
	provider.ClientSecret = os.Getenv("TOKEN")
	provider.RedirectUri = "http://127.0.0.1:8080/callback"

	n, err := ninsho.NewNinsho(&provider, &ninsho.LINE_LOGIN)
	if err != nil {
		panic(err)
	}

	authURL := n.GetAuthURL()

	fmt.Printf("Open this URL in your browser:\n%v\n\nEnter Code:", authURL)

	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	code = strings.Replace(code, "\n", "", -1)

	jwt, err := n.Auth(code)

	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %v", jwt)
}
