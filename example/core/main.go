package main

import (
	"bufio"
	"fmt"
	"line_login_core"
	"os"
	"strings"
)

func main() {
	var provider line_login_core.Provider

	provider.ClientID = os.Getenv("CLIENT_ID")
	provider.ClientSecret = os.Getenv("TOKEN")
	provider.RedirectURL = "http://localhost:3000/api/auth/callback/line"

	session, err := line_login_core.NewSession(&provider)
	if err != nil {
		panic(err)
	}

	authURL := session.AuthURL()

	fmt.Printf("Open this URL in your browser:\n%v\n\nEnter Code:", authURL)

	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	code = strings.Replace(code, "\n", "", -1)

	user, err := session.GetUser(code)

	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %v", user)
}
