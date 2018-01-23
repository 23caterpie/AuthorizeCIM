package AuthorizeCIM

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

var api_endpoint string = "https://apitest.authorize.net/xml/v1/request.api"
var apiName *string
var apiKey *string
var testMode string
var showDebug bool = true
var connected bool = false
var apiLogger *log.Logger
var apihttpClient *http.Client

func SetAPIInfo(name string, key string, mode string, showDebugLogs bool, logger *log.Logger, httpClient *http.Client) {
	apiKey = &key
	apiName = &name
	apiLogger = logger
	apihttpClient = httpClient
	showDebug = showDebugLogs
	if mode == "live" {
		testMode = "liveMode"
		api_endpoint = "https://api.authorize.net/xml/v1/request.api"
	} else {
		testMode = "testMode"
		api_endpoint = "https://apitest.authorize.net/xml/v1/request.api"
	}
}

func IsConnected() (bool, error) {
	info, err := GetMerchantDetails()
	if err != nil {
		return false, err
	}
	if info.Ok() {
		return true, err
	}
	return false, err
}

func GetAuthentication() MerchantAuthentication {
	auth := MerchantAuthentication{
		Name:           apiName,
		TransactionKey: apiKey,
	}
	return auth
}

func SendRequest(input []byte) ([]byte, error) {
	if apihttpClient == nil {
		return nil, errors.New("http client must be set")
	}
	req, err := http.NewRequest("POST", api_endpoint, bytes.NewBuffer(input))
	req.Header.Set("Content-Type", "application/json")
	resp, err := apihttpClient.Do(req)
	if err != nil {
		logf("Endpoint: %s", api_endpoint)
		logf("Request: %s", input)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logf("Endpoint: %s", api_endpoint)
		logf("Request: %s", input)
		return nil, err
	}
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	debugLogf("Endpoint: %s", api_endpoint)
	debugLogf("Request Body: %s", input)
	debugLogf("Response Body: %s", body)
	return body, err
}

func logf(format string, args ...interface{}) {
	if apiLogger != nil {
		apiLogger.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func debugLogf(format string, args ...interface{}) {
	if !showDebug {
		return
	}
	if apiLogger != nil {
		apiLogger.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (r AVS) Text() string {
	var response string
	switch r.avsResultCode {
	case "E":
		response = "AVS data provided is invalid or AVS is not allowed for the card type that was used."
	case "R":
		response = "The AVS system was unavailable at the time of processing."
	case "G":
		response = "The card issuing bank is of non-U.S. origin and does not support AVS"
	case "U":
		response = "The address information for the cardholder is unavailable."
	case "S":
		response = "The U.S. card issuing bank does not support AVS."
	case "N":
		response = "Address: No Match ZIP Code: No Match"
	case "A":
		response = "Address: Match ZIP Code: No Match"
	case "Z":
		response = "Address: No Match ZIP Code: Match"
	case "W":
		response = "Address: No Match ZIP Code: Matched 9 digits"
	case "X":
		response = "Address: Match ZIP Code: Matched 9 digits"
	case "Y":
		response = "Address: Match ZIP: Matched first 5 digits"
	}
	return response
}
