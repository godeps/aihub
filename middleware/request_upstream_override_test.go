package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/gin-gonic/gin"
)

func TestParseUpstreamOverrideEnabledAllowlisted(t *testing.T) {
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"proxy.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"proxy.example.com": "http://corp-proxy:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://proxy.example.com")
	c.Request.Header.Set(requestUpstreamHeadersHeader, "{\"Authorization\":\"Bearer test\"}")
	commonSetting := dto.ChannelSettings{Proxy: "http://existing-proxy:8080"}
	c.Set(string(constant.ContextKeyChannelSetting), commonSetting)

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
	updatedSetting := c.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "http://corp-proxy:8080" {
		t.Fatalf("unexpected proxy: %s", updatedSetting.Proxy)
	}
}

func TestParseUpstreamOverrideDisabled(t *testing.T) {
	resetUpstreamProxyCache()
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
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"allowed.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"allowed.example.com": "http://corp-proxy:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
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
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"proxy.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"proxy.example.com": "http://corp-proxy:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
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

func TestParseUpstreamOverrideProxyMapNoMatch(t *testing.T) {
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"proxy.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"other.example.com": "http://corp-proxy:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://proxy.example.com")
	commonSetting := dto.ChannelSettings{Proxy: "http://existing-proxy:8080"}
	c.Set(string(constant.ContextKeyChannelSetting), commonSetting)

	override, err := parseUpstreamOverrideFromRequest(c)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	updatedSetting := c.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "http://existing-proxy:8080" {
		t.Fatalf("unexpected proxy: %s", updatedSetting.Proxy)
	}
}

func TestParseUpstreamOverrideProxyMapWildcard(t *testing.T) {
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"api.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"*.example.com": "http://corp-proxy:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://api.example.com")
	commonSetting := dto.ChannelSettings{}
	c.Set(string(constant.ContextKeyChannelSetting), commonSetting)

	override, err := parseUpstreamOverrideFromRequest(c)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	updatedSetting := c.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "http://corp-proxy:8080" {
		t.Fatalf("unexpected proxy: %s", updatedSetting.Proxy)
	}
}

func TestParseUpstreamOverrideProxyMapSuffix(t *testing.T) {
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"api.example.com", "example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		".example.com": "http://corp-proxy:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://api.example.com")
	commonSetting := dto.ChannelSettings{}
	c.Set(string(constant.ContextKeyChannelSetting), commonSetting)

	override, err := parseUpstreamOverrideFromRequest(c)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	updatedSetting := c.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "http://corp-proxy:8080" {
		t.Fatalf("unexpected proxy: %s", updatedSetting.Proxy)
	}
}

func TestParseUpstreamOverrideProxyMapExactPrecedence(t *testing.T) {
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"api.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"api.example.com": "http://exact-proxy:8080",
		"*.example.com":   "http://wildcard-proxy:8080",
		".example.com":    "http://suffix-proxy:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://api.example.com")
	commonSetting := dto.ChannelSettings{}
	c.Set(string(constant.ContextKeyChannelSetting), commonSetting)

	override, err := parseUpstreamOverrideFromRequest(c)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	updatedSetting := c.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "http://exact-proxy:8080" {
		t.Fatalf("unexpected proxy: %s", updatedSetting.Proxy)
	}
}

func TestParseUpstreamOverrideProxyMapNoProxyOverride(t *testing.T) {
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"api.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"api.example.com": "none",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://api.example.com")
	commonSetting := dto.ChannelSettings{Proxy: "http://existing-proxy:8080"}
	c.Set(string(constant.ContextKeyChannelSetting), commonSetting)

	override, err := parseUpstreamOverrideFromRequest(c)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	updatedSetting := c.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "" {
		t.Fatalf("expected proxy cleared, got %s", updatedSetting.Proxy)
	}
}

func TestParseUpstreamOverrideProxyMapCacheInvalidation(t *testing.T) {
	resetUpstreamProxyCache()
	setting := operation_setting.GetGeneralSetting()
	originalEnabled := setting.RequestUpstreamOverrideEnabled
	originalAllowlist := append([]string{}, setting.RequestUpstreamOverrideAllowlist...)
	originalProxyMap := map[string]string{}
	for key, value := range setting.RequestUpstreamProxyMap {
		originalProxyMap[key] = value
	}
	setting.RequestUpstreamOverrideEnabled = true
	setting.RequestUpstreamOverrideAllowlist = []string{"api.example.com"}
	setting.RequestUpstreamProxyMap = map[string]string{
		"api.example.com": "http://proxy-1:8080",
	}
	defer func() {
		setting.RequestUpstreamOverrideEnabled = originalEnabled
		setting.RequestUpstreamOverrideAllowlist = originalAllowlist
		setting.RequestUpstreamProxyMap = originalProxyMap
	}()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c.Request.Header.Set(requestUpstreamBaseURLHeader, "https://api.example.com")
	c.Set(string(constant.ContextKeyChannelSetting), dto.ChannelSettings{})

	override, err := parseUpstreamOverrideFromRequest(c)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	updatedSetting := c.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "http://proxy-1:8080" {
		t.Fatalf("unexpected proxy: %s", updatedSetting.Proxy)
	}

	setting.RequestUpstreamProxyMap = map[string]string{
		"api.example.com": "http://proxy-2:8080",
	}

	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	c2.Request = httptest.NewRequest("POST", "/v1beta/models/test:generateContent", nil)
	c2.Request.Header.Set(requestUpstreamBaseURLHeader, "https://api.example.com")
	c2.Set(string(constant.ContextKeyChannelSetting), dto.ChannelSettings{})

	override, err = parseUpstreamOverrideFromRequest(c2)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if override == nil {
		t.Fatal("expected override")
	}
	updatedSetting = c2.MustGet(string(constant.ContextKeyChannelSetting)).(dto.ChannelSettings)
	if updatedSetting.Proxy != "http://proxy-2:8080" {
		t.Fatalf("unexpected proxy: %s", updatedSetting.Proxy)
	}
}
