package healthcheck

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockHandlerSuccessHiringLine(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func MockHandlerFailureHiringLine(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func TestSendReportToHiringLine(t *testing.T) {
	tests := map[string]struct{
		handler func(w http.ResponseWriter, r *http.Request)
		isErrorNil bool
	} {
		"MustBePassed": {MockHandlerSuccessHiringLine, true },
		"MustBeFailed": { MockHandlerFailureHiringLine, false },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(test.handler))
			defer ts.Close()

			err := SendReportToHiringLine(ts.URL,
				&oauth2.Token{AccessToken: "LINE_ACCESS_TOKEN"},
				Report{20, 7, 13},
				2000)

			assert.Equal(t, test.isErrorNil, err == nil)
		})
	}


}