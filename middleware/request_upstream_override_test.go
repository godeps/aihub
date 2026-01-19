package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/gin-gonic/gin"
)

func TestParseUpstreamOverrideEnabledAllowlisted(t *testing.T) {
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"proxy.example.com"}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://proxy.example.com")
	c.Request.Header.Set(requestUpstreamHeadersHeader, "{\"Authorization\":\"Bearer test\"}")

	override, err := parseUpstreamOverrideFromRequest(c)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	if override.BaseURL != "https://proxy.example.com" {
		t.Fatalf("unexpected base url: %s", override.BaseURL)
	}
	if override.Headers["Authorization"] != "Bearer test" {
		t.Fatalf("unexpected header value: %s", override.Headers["Authorization"])
	}
}

func TestParseUpstreamOverrideDisabled(t *testing.T) {
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	setting.RequestUpstreamOverrideEnabled = false
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://proxy.example.com")

	override, err := parseUpstreamOverrideFromRequest(c)
	if err == nil || err.StatusCode != 403 {
		t.Fatalf("expected 403 error, got override=%v err=%+v", override, err)
	}
}

func TestParseUpstreamOverrideNotAllowlisted(t *testing.T) {
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"allowed.example.com"}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://proxy.example.com")

	override, err := parseUpstreamOverrideFromRequest(c)
	if err == nil || err.StatusCode != 403 {
		t.Fatalf("expected 403 error, got override=%v err=%+v", override, err)
	}
}

func TestParseUpstreamOverrideInvalidHeaders(t *testing.T) {
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"proxy.example.com"}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://proxy.example.com")
	c.Request.Header.Set(requestUpstreamHeadersHeader, "{not-json}")

	override, err := parseUpstreamOverrideFromRequest(c)
	if err == nil || err.StatusCode != 400 {
		t.Fatalf("expected 400 error, got override=%v err=%+v", override, err)
	}
}
