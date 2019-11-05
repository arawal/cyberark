package main

import (
	"fmt"
	"strings"

	"github.com/arawal/cyberark"
)

func main() {
	var err error
	creds := cyberark.Credentials{
		Username: "usernamehere",
		Password: "passwordhere",
		BaseURL:  "urlhere",
	}
	token, err := cyberark.Authenticate(creds)
	creds.AuthToken = strings.Replace(token, "\"", "", -1)
	fmt.Println(creds.AuthToken, err)
	res, err := cyberark.GetAccounts(creds)
	fmt.Println(res, err)
}
