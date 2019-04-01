package line

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

)

func TestAuthLineCallback(t *testing.T) {
	state := "state"
	sendCode := make(chan string, 1)
	handler := lineLoginHandler{state, sendCode}

	req, err := http.NewRequest("GET", "http://localhost:8080/auth/callback", nil)
	if err != nil {
		t.Error(err)
	}

	urlQuery := url.Values{}
	urlQuery.Add("code", "wflp321")
	urlQuery.Add("state", "state")
	urlQuery.Add("friendship_status_changed", "false")

	req.URL.RawQuery = urlQuery.Encode()

	rcd := httptest.NewRecorder()
	if err != nil {
		t.Error(err)
	}

	handler.authLineCallback(rcd, req)

	resp := rcd.Result()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	code := <-sendCode
	assert.NotEmpty(t, code)
}