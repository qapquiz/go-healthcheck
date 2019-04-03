package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
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

func TestReport_IsCheckAnyWebsite(t *testing.T) {
	tests := map[string]struct{
		report Report
		expected bool
	} {
		"MustReturnTrue": { Report{totalWebsites: 20}, true },
		"MustReturnFalse": { Report{totalWebsites: 0}, false },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.report.IsCheckAnyWebsite())
		})
	}
}