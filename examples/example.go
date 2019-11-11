package main

import (
	"fmt"

	"github.com/arawal/cyberark"
)

func main() {
	var err error
	creds := cyberark.Credentials{
		Username: "usernamehere",
		Password: "passwordhere",
		BaseURL:  "urlhere",
	}
	creds.AuthToken, err = cyberark.Authenticate(creds)
	fmt.Println(creds.AuthToken, err)
	res, err := cyberark.GetAccounts(creds, cyberark.AccountsRequestParams{
		SearchBy:      "",
		SortOn:        "",
		SortDirection: "",
		Offset:        "",
		Limit:         "",
	})
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
