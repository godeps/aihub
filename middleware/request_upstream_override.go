package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

const (
	requestUpstreamBaseURLHeader = "X-Relay-Upstream-Base-URL"
	requestUpstreamHeadersHeader = "X-Relay-Upstream-Headers"
)

type upstreamOverride struct {
	BaseURL string
	Headers map[string]string
}

type upstreamOverrideError struct {
	StatusCode int
	Code       string
	Message    string
}

const (
	proxyNoOverrideValue       = "none"
	upstreamProxyCacheSizeLimit = 1024
)

var upstreamProxyCache sync.Map
var upstreamProxyCacheCount atomic.Int64
var upstreamProxyMapPtr atomic.Uintptr

func applyRequestUpstreamOverride(c *gin.Context) bool {
	override, overrideErr := parseUpstreamOverrideFromRequest(c)
	if overrideErr != nil {
		abortWithOpenAiMessage(c, overrideErr.StatusCode, overrideErr.Message, overrideErr.Code)
		return false
	}
	if override == nil {
		return true
	}
	common.SetContextKey(c, constant.ContextKeyChannelBaseUrl, override.BaseURL)
	if len(override.Headers) > 0 {
		headerOverride := make(map[string]interface{}, len(override.Headers))
		for key, value := range override.Headers {
			headerOverride[key] = value
		}
		common.SetContextKey(c, constant.ContextKeyChannelHeaderOverride, headerOverride)
	}
	logger.LogInfo(
		c.Request.Context(),
		fmt.Sprintf(
			"request upstream override applied: base_url=%s headers=%d user=%d",
			common.MaskSensitiveInfo(override.BaseURL),
			len(override.Headers),
			c.GetInt("id"),
		),
	)
	return true
}

func parseUpstreamOverrideFromRequest(c *gin.Context) (*upstreamOverride, *upstreamOverrideError) {
	baseURLHeader := strings.TrimSpace(c.GetHeader(requestUpstreamBaseURLHeader))
	headersHeader := strings.TrimSpace(c.GetHeader(requestUpstreamHeadersHeader))
	if baseURLHeader == "" && headersHeader == "" {
		return nil, nil
	}
	if !operation_setting.IsRequestUpstreamOverrideEnabled() {
		return nil, &upstreamOverrideError{
			StatusCode: http.StatusForbidden,
			Code:       string(types.ErrorCodeAccessDenied),
			Message:    "request upstream override is disabled",
		}
	}
	if baseURLHeader == "" {
		return nil, &upstreamOverrideError{
			StatusCode: http.StatusBadRequest,
			Code:       string(types.ErrorCodeInvalidRequest),
			Message:    "request upstream base url is required",
		}
	}
	upstreamHost, hostErr := parseUpstreamHost(baseURLHeader)
	if hostErr != nil {
		return nil, &upstreamOverrideError{
			StatusCode: http.StatusBadRequest,
			Code:       string(types.ErrorCodeInvalidRequest),
			Message:    hostErr.Error(),
		}
	}
	allowlist := operation_setting.GetRequestUpstreamOverrideAllowlist()
	allowed, allowErr := isUpstreamAllowlisted(upstreamHost, allowlist)
	if allowErr != nil {
		return nil, &upstreamOverrideError{
			StatusCode: http.StatusBadRequest,
			Code:       string(types.ErrorCodeInvalidRequest),
			Message:    allowErr.Error(),
		}
	}
	if !allowed {
		return nil, &upstreamOverrideError{
			StatusCode: http.StatusForbidden,
			Code:       string(types.ErrorCodeAccessDenied),
			Message:    "request upstream override not allowlisted",
		}
	}
	headers := map[string]string{}
	if headersHeader != "" {
		if err := json.Unmarshal([]byte(headersHeader), &headers); err != nil {
			return nil, &upstreamOverrideError{
				StatusCode: http.StatusBadRequest,
				Code:       string(types.ErrorCodeInvalidRequest),
				Message:    "request upstream headers invalid",
			}
		}
		for key, value := range headers {
			if strings.TrimSpace(key) == "" {
				return nil, &upstreamOverrideError{
					StatusCode: http.StatusBadRequest,
					Code:       string(types.ErrorCodeInvalidRequest),
					Message:    "request upstream headers include empty key",
				}
			}
			headers[key] = value
		}
	}

	if proxyURL, ok := resolveUpstreamProxy(upstreamHost); ok {
		channelSetting, ok := common.GetContextKeyType[dto.ChannelSettings](c, constant.ContextKeyChannelSetting)
		if ok {
			channelSetting.Proxy = proxyURL
			common.SetContextKey(c, constant.ContextKeyChannelSetting, channelSetting)
		}
	}

	return &upstreamOverride{
		BaseURL: baseURLHeader,
		Headers: headers,
	}, nil
}

func parseUpstreamHost(baseURL string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("request upstream base url invalid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("request upstream base url scheme not allowed")
	}
	host := strings.ToLower(parsed.Host)
	if host == "" {
		return "", fmt.Errorf("request upstream base url host missing")
	}
	return host, nil
}

func isUpstreamAllowlisted(host string, allowlist []string) (bool, error) {
	for _, entry := range allowlist {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if entry == "*" {
			return true, nil
		}
		entryHost := normalizeUpstreamHost(entry)
		if entryHost != "" && host == entryHost {
			return true, nil
		}
	}
	return false, nil
}

func resolveUpstreamProxy(host string) (string, bool) {
	proxyMap := operation_setting.GetRequestUpstreamProxyMap()
	if len(proxyMap) == 0 {
		return "", false
	}
	mapPtr := reflect.ValueOf(proxyMap).Pointer()
	if upstreamProxyMapPtr.Load() != mapPtr {
		upstreamProxyCache.Clear()
		upstreamProxyCacheCount.Store(0)
		upstreamProxyMapPtr.Store(mapPtr)
	}
	if cached, ok := upstreamProxyCache.Load(host); ok {
		return cached.(string), true
	}
	if upstreamProxyCacheCount.Load() > upstreamProxyCacheSizeLimit {
		upstreamProxyCache.Clear()
		upstreamProxyCacheCount.Store(0)
	}
	exactMatch := ""
	exactMatched := false
	wildcardMatch := ""
	wildcardLen := 0
	suffixMatch := ""
	suffixLen := 0

	for key, proxyURL := range proxyMap {
		entryHost := normalizeUpstreamHost(key)
		if entryHost == "" {
			continue
		}
		proxyURL = strings.TrimSpace(proxyURL)
		if strings.EqualFold(proxyURL, proxyNoOverrideValue) {
			proxyURL = ""
		}
		if entryHost == host {
			exactMatch = proxyURL
			exactMatched = true
			break
		}
		if strings.HasPrefix(entryHost, "*.") {
			suffix := strings.TrimPrefix(entryHost, "*.")
			if suffix != "" && strings.HasSuffix(host, "."+suffix) {
				if len(suffix) > wildcardLen {
					wildcardLen = len(suffix)
					wildcardMatch = proxyURL
				}
			}
			continue
		}
		if strings.HasPrefix(entryHost, ".") {
			suffix := strings.TrimPrefix(entryHost, ".")
			if suffix != "" && (host == suffix || strings.HasSuffix(host, "."+suffix)) {
				if len(suffix) > suffixLen {
					suffixLen = len(suffix)
					suffixMatch = proxyURL
				}
			}
		}
	}

	if exactMatched {
		cacheProxyResolution(host, exactMatch)
		return exactMatch, true
	}
	if wildcardMatch != "" {
		cacheProxyResolution(host, wildcardMatch)
		return wildcardMatch, true
	}
	if suffixMatch != "" {
		cacheProxyResolution(host, suffixMatch)
		return suffixMatch, true
	}
	return "", false
}

func cacheProxyResolution(host string, proxyURL string) {
	if _, loaded := upstreamProxyCache.LoadOrStore(host, proxyURL); !loaded {
		upstreamProxyCacheCount.Add(1)
	} else {
		upstreamProxyCache.Store(host, proxyURL)
	}
	if upstreamProxyCacheCount.Load() > upstreamProxyCacheSizeLimit {
		upstreamProxyCache.Clear()
		upstreamProxyCacheCount.Store(0)
	}
}

func resetUpstreamProxyCache() {
	upstreamProxyCache.Clear()
	upstreamProxyCacheCount.Store(0)
	upstreamProxyMapPtr.Store(0)
}

func normalizeUpstreamHost(entry string) string {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return ""
	}
	if strings.Contains(entry, "://") {
		parsed, err := url.Parse(entry)
		if err != nil {
			return ""
		}
		return strings.ToLower(strings.TrimSpace(parsed.Host))
	}
	return strings.ToLower(entry)
}
