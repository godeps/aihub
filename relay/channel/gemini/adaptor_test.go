package gemini

import (
	"net/http"
	"net/http/httptest"
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

func TestSetupRequestHeaderKeepsOverrideKey(t *testing.T) {
	adapt := &Adaptor{}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("POST", "/v1beta/models/gemini:generateContent", nil)
	c.Request = req

	headers := http.Header{}
	headers.Set("x-goog-api-key", "override-key")
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiKey: "default-key",
		},
	}

	if err := adapt.SetupRequestHeader(c, &headers, info); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := headers.Get("x-goog-api-key"); got != "override-key" {
		t.Fatalf("expected override key, got %s", got)
	}
}

func TestSetupRequestHeaderSkipsWhenAuthorizationPresent(t *testing.T) {
	adapt := &Adaptor{}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("POST", "/v1beta/models/gemini:generateContent", nil)
	c.Request = req

	headers := http.Header{}
	headers.Set("Authorization", "Bearer override")
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiKey: "default-key",
		},
	}

	if err := adapt.SetupRequestHeader(c, &headers, info); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := headers.Get("x-goog-api-key"); got != "" {
		t.Fatalf("expected no x-goog-api-key set, got %s", got)
	}
}
