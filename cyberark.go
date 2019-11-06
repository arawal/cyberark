package cyberark

import (
	"encoding/json"
)

// Credentials provides an interface to store fields required for authentication
type Credentials struct {
	Username  string
	Password  string
	BaseURL   string
	AuthToken string
}

type PlatformAccountProperties map[string]interface{}

type SecretManagement struct {
	AutomaticManagementEnabled bool   `json:"automaticManagementEnabled"`
	ManualManagementReason     string `json:"manualManagementReason"`
	Status                     string `json:"status"`
	LastModifiedTime           int64  `json:"lastModifiedTime"`
}

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

type AccountsResponse struct {
	Value []AccountSummary `json:"value"`
	Count int64            `json:"count"`
}

type SafeSummary struct {
	Description               string      `json:"Description"`
	ManagingCPM               string      `json:"ManagingCPM"`
	NumberOfDaysRetention     int         `json:"NumberOfDaysRetention"`
	NumberOfVersionsRetention interface{} `json:"NumberOfVersionsRetention"`
	OLACEnabled               bool        `json:"OLACEnabled"`
	SafeName                  string      `json:"SafeName"`
}

type SafesResponse struct {
	GetSafesResult []SafeSummary `json:"GetSafesResult"`
}

type CustomRequestParams struct {
	Method   string
	Endpoint string
	Payload  map[string]interface{}
	Headers  map[string]string
}

// Authenticate authenticates a sessions and returns the session token
func Authenticate(creds Credentials) (string, error) {
	res, err := login(creds)
	return string(res), err
}

func GetAccounts(creds Credentials) (ar AccountsResponse, err error) {
	res, err := getAccounts(creds)
	if err != nil {
		return ar, err
	}

	err = json.Unmarshal(res, &ar)
	return ar, err
}

func GetSafes(creds Credentials) (sr SafesResponse, err error) {
	res, err := getSafes(creds)
	if err != nil {
		return sr, err
	}

	err = json.Unmarshal(res, &sr)
	return sr, err
}

func MakeCustomAPIRequest(creds Credentials, params CustomRequestParams) (r map[string]interface{}, err error) {
	res, err := makeCustomAPIRequest(creds, params)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(res, &r)
	return r, err
}
