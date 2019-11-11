package cyberark

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient http.Client

type serverError struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

func init() {
	// Modify transport protocol to skip cert validation
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{},
			InsecureSkipVerify: true,
		},
	}

	// Initialize HTTP client with the transport protocol modifications and custom http timeout
	httpClient = http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
}

func login(credentials Credentials) ([]byte, error) {
	url := fmt.Sprintf("%s/API/Auth/CyberArk/Logon", credentials.BaseURL)
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf("{\n	\"username\": \"%s\",\n	\"password\": \"%s\"\n}", credentials.Username, credentials.Password))

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	res, err := sendRequest(req, credentials)
	if err != nil {
		return nil, err
	}

	var errRes serverError
	if mErr := json.Unmarshal(res, &errRes); mErr == nil {
		return nil, fmt.Errorf("Server returned error code: %s. Error message: %s", errRes.ErrorCode, strings.Replace(errRes.ErrorMessage, "\n", " ", -1))
	}
	return res, nil
}

func getAccounts(credentials Credentials, params AccountsRequestParams) ([]byte, error) {
	qp := url.QueryEscape(fmt.Sprintf("search=%s&sort=%s&offset=%s&limit=%s", params.SearchBy, params.SortOn, params.Offset, params.Limit))
	endpoint := fmt.Sprintf("%s/api/Accounts?%s", credentials.BaseURL, qp)
	method := "GET"

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := sendRequest(req, credentials)
	if err != nil {
		return nil, err
	}

	var errRes serverError
	if mErr := json.Unmarshal(res, &errRes); mErr == nil && (errRes.ErrorCode != "" || errRes.ErrorMessage != "") {
		return nil, fmt.Errorf("Server returned error code: %s. Error message: %s", errRes.ErrorCode, strings.Replace(errRes.ErrorMessage, "\n", " ", -1))
	}
	return res, nil
}

func getSafes(credentials Credentials) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/WebServices/PIMServices.svc/Safes", credentials.BaseURL)
	method := "GET"

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := sendRequest(req, credentials)
	if err != nil {
		return nil, err
	}

	var errRes serverError
	if mErr := json.Unmarshal(res, &errRes); mErr == nil && (errRes.ErrorCode != "" || errRes.ErrorMessage != "") {
		return nil, fmt.Errorf("Server returned error code: %s. Error message: %s", errRes.ErrorCode, strings.Replace(errRes.ErrorMessage, "\n", " ", -1))
	}
	return res, nil
}

func makeCustomAPIRequest(credentials Credentials, params CustomRequestParams) ([]byte, error) {
	qp := strings.Split(params.Endpoint, "?")
	if len(qp) > 1 {
		params.Endpoint = qp[0] + url.QueryEscape(qp[1])
	}
	endpoint := fmt.Sprintf("%s/%s", credentials.BaseURL, params.Endpoint)
	endpoint = url.QueryEscape(endpoint)
	method := params.Method

	reqBody, err := json.Marshal(params.Payload)
	if err != nil {
		return nil, fmt.Errorf("Malformed request payload: %s", err.Error())
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	for k, v := range params.Headers {
		req.Header.Add(k, v)
	}

	res, err := sendRequest(req, credentials)
	if err != nil {
		return nil, err
	}

	var errRes serverError
	if mErr := json.Unmarshal(res, &errRes); mErr == nil && (errRes.ErrorCode != "" || errRes.ErrorMessage != "") {
		return nil, fmt.Errorf("Server returned error code: %s. Error message: %s", errRes.ErrorCode, strings.Replace(errRes.ErrorMessage, "\n", " ", -1))
	}
	return res, nil
}

// SendRequest preps and sends the appropriate http request to the splunk server and retrieves the data from the response
func sendRequest(req *http.Request, credentials Credentials) ([]byte, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "close")
	req.Header.Add("Authorization", credentials.AuthToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	return data, err
}
