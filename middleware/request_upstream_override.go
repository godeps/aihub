package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
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
	allowlist := operation_setting.GetRequestUpstreamOverrideAllowlist()
	allowed, allowErr := isUpstreamAllowlisted(baseURLHeader, allowlist)
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
	return &upstreamOverride{
		BaseURL: baseURLHeader,
		Headers: headers,
	}, nil
}

func isUpstreamAllowlisted(baseURL string, allowlist []string) (bool, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return false, fmt.Errorf("request upstream base url invalid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false, fmt.Errorf("request upstream base url scheme not allowed")
	}
	host := strings.ToLower(parsed.Host)
	if host == "" {
		return false, fmt.Errorf("request upstream base url host missing")
	}
	for _, entry := range allowlist {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if entry == "*" {
			return true, nil
		}
		entryHost := entry
		if strings.Contains(entry, "://") {
			parsedEntry, err := url.Parse(entry)
			if err != nil {
				continue
			}
			entryHost = parsedEntry.Host
		}
		entryHost = strings.ToLower(strings.TrimSpace(entryHost))
		if entryHost != "" && host == entryHost {
			return true, nil
		}
	}
	return false, nil
}
