package cyberark

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Credentials provides an interface to store fields required for authentication
type Credentials struct {
	Username  string
	Password  string
	BaseURL   string
	AuthToken string
}

// PlatformAccountProperties type on an account in Cyberark
type PlatformAccountProperties map[string]interface{}

// SecretManagement data type for an account in Cyberark
type SecretManagement struct {
	AutomaticManagementEnabled bool   `json:"automaticManagementEnabled"`
	ManualManagementReason     string `json:"manualManagementReason"`
	Status                     string `json:"status"`
	LastModifiedTime           int64  `json:"lastModifiedTime"`
}

// AccountSummary is the response type for a particular account from a GetAccounts call to Cyberark
type AccountSummary struct {
	ID                        string                    `json:"id"`
	Name                      string                    `json:"name"`
	Address                   string                    `json:"address"`
	Username                  string                    `json:"userName"`
	PlatformID                string                    `json:"platformId"`
	SafeName                  string                    `json:"safeName"`
	SecretType                string                    `json:"secretType"`
	PlatformAccountProperties PlatformAccountProperties `json:"platformAccountProperties"`
	SecretManagement          SecretManagement          `json:"secretManagement"`
	CreatedTime               int64                     `json:"createdTime"`
}

// AccountsResponse is the response type from a GetAccounts call to Cyberark
type AccountsResponse struct {
	Value []AccountSummary `json:"value"`
	Count int64            `json:"count"`
}

// SafeSummary is the response type for a particular safe from a GetSafes call to Cyberark
type SafeSummary struct {
	Description               string      `json:"Description"`
	ManagingCPM               string      `json:"ManagingCPM"`
	NumberOfDaysRetention     int         `json:"NumberOfDaysRetention"`
	NumberOfVersionsRetention interface{} `json:"NumberOfVersionsRetention"`
	OLACEnabled               bool        `json:"OLACEnabled"`
	SafeName                  string      `json:"SafeName"`
}

// SafesResponse is the response type from a GetSafes call to Cyberark
type SafesResponse struct {
	GetSafesResult []SafeSummary `json:"GetSafesResult"`
}

// AccountsRequestParams provides an interface for fields needed to filter results returned by the GetAccounts call
type AccountsRequestParams struct {
	SearchBy      string `json:"search_by"`
	SortOn        string `json:"sort_on"`
	SortDirection string `json:"sort_direction"`
	Limit         string `json:"limit"`
	Offset        string `json:"offset"`
}

// CustomRequestParams provides an interface for fields needed for a generic call to Cyberark
type CustomRequestParams struct {
	Method   string
	Endpoint string
	Payload  map[string]interface{}
	Headers  map[string]string
}

// Authenticate authenticates a sessions and returns the session token
func Authenticate(creds Credentials) (string, error) {
	res, err := login(creds)
	return strings.Replace(string(res), "\"", "", -1), err
}

// GetAccounts gets a list of accounts based on the search criteria
func GetAccounts(creds Credentials, params AccountsRequestParams) (ar AccountsResponse, err error) {
	if creds.AuthToken == "" {
		creds.AuthToken, err = Authenticate(creds)
		if err != nil {
			return
		}
	}

	if params.SortOn != "" {
		if params.SortDirection != "" {
			params.SortOn = fmt.Sprintf("%s %s", params.SortOn, params.SortDirection)
		}
	}
	res, err := getAccounts(creds, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(res, &ar)
	return
}

// GetSafes gets a list of safes based on the search criteria
func GetSafes(creds Credentials) (sr SafesResponse, err error) {
	res, err := getSafes(creds)
	if err != nil {
		return
	}

	err = json.Unmarshal(res, &sr)
	return
}

// MakeCustomAPIRequest make a generic request to Cyberark
func MakeCustomAPIRequest(creds Credentials, params CustomRequestParams) (r map[string]interface{}, err error) {
	res, err := makeCustomAPIRequest(creds, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(res, &r)
	return
}
