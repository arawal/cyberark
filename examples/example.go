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
	res2, err := cyberark.GetSafes(creds)
	fmt.Println(res2, err)
	res3, err := cyberark.MakeCustomAPIRequest(creds, cyberark.CustomRequestParams{
		Method:   "GET",
		Endpoint: "/api/accounts?limit=1",
		Headers:  map[string]string{},
		Payload:  nil,
	})
	fmt.Println(res3, err)
}
