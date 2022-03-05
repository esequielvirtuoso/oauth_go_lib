package oauth

import (
	"net/http"
	"os"
	"testing"

	"github.com/esequielvirtuoso/go_utils_lib/logger"
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.Info("starting oauth lib unit tests...")
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestOauthConstants(t *testing.T) {
	assert.EqualValues(t, "X-Public", headerXPublic)
	assert.EqualValues(t, "X-Client-Id", headerXClientID)
	assert.EqualValues(t, "X-Caller-Id", headerXCallerID)
	assert.EqualValues(t, "access_token", paramAccessToken)
}

func TestIsPublicNilRequest(t *testing.T) {
	assert.True(t, IsPublic(nil))
}

func TestIsPublicNoError(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}
	assert.False(t, IsPublic(&request))

	request.Header.Add("X-Public", "true")
	assert.True(t, IsPublic(&request))
}

func TestGetCallerIdNilRequest(t *testing.T) {
	assert.Zero(t, GetCallerId(nil))
}

func TestGetCallerInvalidFormat(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}

	request.Header.Add("X-Caller-Id", "notNumber")
	assert.Zero(t, GetCallerId(&request))
}

func TestGetCallerNoError(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}

	request.Header.Add("X-Caller-Id", "1")
	assert.EqualValues(t, 1, GetCallerId(&request))
}

func TestClientCallerIdNilRequest(t *testing.T) {
	assert.Zero(t, GetClientId(nil))
}

func TestGetClientInvalidFormat(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}

	request.Header.Add("X-Caller-Id", "notNumber")
	assert.Zero(t, GetClientId(&request))
}

func TestGetClientNoError(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}

	request.Header.Add("X-Client-Id", "1")
	assert.EqualValues(t, 1, GetClientId(&request))
}

func TestGetAccessTokenInvalidRestClientResponse(t *testing.T) {
	rest.FlushMockups() // flush all the mockups we have
	// create a mock
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodGet,
		URL:          "http://localhost:5002/oauth/access_token/Abc123",
		ReqBody:      ``,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	accessToken, err := getAccessToken("Abc123")
	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid rest client response when trying to get access token", err.Message())
}

func TestGetAccessTokenInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups() // flush all the mockups we have
	// create a mock
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodGet,
		URL:          "http://localhost:5002/oauth/access_token/Abc123",
		ReqBody:      ``,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid error interface when trying to get access token"}`,
	})

	accessToken, err := getAccessToken("Abc123")
	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid error interface when trying to get access token", err.Message())
}

func TestGetAccessTokenInvalidJson(t *testing.T) {
	rest.FlushMockups() // flush all the mockups we have
	// create a mock
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodGet,
		URL:          "http://localhost:5002/oauth/access_token/Abc123",
		ReqBody:      ``,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid json"}`,
	})

	accessToken, err := getAccessToken("Abc123")
	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid json", err.Message())
}

func TestGetAccessTokenUnmarshalTokenError(t *testing.T) {
	rest.FlushMockups() // flush all the mockups we have
	// create a mock
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodGet,
		URL:          "http://localhost:5002/oauth/access_token/Abc123",
		ReqBody:      ``,
		RespHTTPCode: http.StatusInternalServerError,
		RespBody:     `{"message": "error when trying to unmarshal access token response"}`,
	})

	accessToken, err := getAccessToken("Abc123")
	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error when trying to unmarshal access token response", err.Message())
}

func TestGetAccessTokenSuccess(t *testing.T) {
	rest.FlushMockups() // flush all the mockups we have
	// create a mock
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodGet,
		URL:          "http://localhost:5002/oauth/access_token/Abc123",
		ReqBody:      ``,
		RespHTTPCode: http.StatusOK,
		RespBody: `{"id": "1",
						"user_id": 2,
						"client_id": 3}`,
	})

	accessToken, err := getAccessToken("Abc123")
	assert.NotNil(t, accessToken)
	assert.Nil(t, err)
	assert.EqualValues(t, "1", accessToken.ID)
	assert.EqualValues(t, 2, accessToken.UserID)
	assert.EqualValues(t, 3, accessToken.ClientID)
}
