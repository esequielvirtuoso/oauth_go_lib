package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	restErrors "github.com/esequielvirtuoso/go_utils_lib/rest_errors"
	"github.com/mercadolibre/golang-restclient/rest"
)

const (
	// This X-Public header is not being set by the client. It will be set by the NGINX server
	// based on the condition if the request was generated from outside the network or inside the network.
	headerXPublic   = "X-Public"
	headerXClientID = "X-Client_Id"
	headerXCallerID = "X-Caller-Id"

	paramAccessToken = "access_token"
)

var (
	oauthRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:5002",
		Timeout: 200 * time.Millisecond,
	}
)

type accessToken struct {
	ID       string `json:"id"`
	UserID   int64  `json:"user_id"`
	ClientID int64  `json:"client_id"`
}

func IsPublic(request *http.Request) bool {
	if request != nil {
		// this means it is a public request
		return true
	}

	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	callerID, err := strconv.ParseInt(request.Header.Get(headerXCallerID), 10, 64)
	if err != nil {
		return 0
	}
	return callerID
}

func GetClientId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	clientID, err := strconv.ParseInt(request.Header.Get(headerXClientID), 10, 64)
	if err != nil {
		return 0
	}
	return clientID
}

func AuthenticateRequest(request *http.Request) restErrors.RestErr {
	if request == nil {
		// this means it is a public request
		return nil
	}

	cleanRequest(request)

	// http://localhost:8082/resource?access_token=abc13
	accessTokenID := strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
	if accessTokenID == "" {
		return nil
	}

	at, err := getAccessToken(accessTokenID)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil
		}
		return err
	}

	request.Header.Add(headerXCallerID, fmt.Sprintf("%v", at.UserID))
	request.Header.Add(headerXClientID, fmt.Sprintf("%v", at.ClientID))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}

	request.Header.Del(headerXClientID)
	request.Header.Del(headerXCallerID)
}

func getAccessToken(accessTokenID string) (*accessToken, restErrors.RestErr) {
	response := oauthRestClient.Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenID))
	if response == nil || response.Response == nil {
		return nil, restErrors.NewInternalServerError("invalid rest client response when trying to get access token", nil)
	}
	if response.StatusCode > 299 {
		var restErr restErrors.RestErr
		if err := json.Unmarshal(response.Bytes(), &restErr); err != nil {
			return nil, restErrors.NewInternalServerError("invalid error interface when trying to get access token", err)
		}

		return nil, restErr
	}

	var at accessToken
	if err := json.Unmarshal(response.Bytes(), &at); err != nil {
		return nil, restErrors.NewInternalServerError("error when trying to unmarshal access token response", err)
	}
	return &at, nil
}
