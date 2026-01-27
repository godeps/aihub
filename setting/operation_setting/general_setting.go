package operation_setting

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/config"
)

// 额度展示类型
const (
	QuotaDisplayTypeUSD    = "USD"
	QuotaDisplayTypeCNY    = "CNY"
	QuotaDisplayTypeTokens = "TOKENS"
	QuotaDisplayTypeCustom = "CUSTOM"
)

type GeneralSetting struct {
	DocsLink            string `json:"docs_link"`
	PingIntervalEnabled bool   `json:"ping_interval_enabled"`
	PingIntervalSeconds int    `json:"ping_interval_seconds"`
	// 当前站点额度展示类型：USD / CNY / TOKENS
	QuotaDisplayType string `json:"quota_display_type"`
	// 自定义货币符号，用于 CUSTOM 展示类型
	CustomCurrencySymbol string `json:"custom_currency_symbol"`
	// 自定义货币与美元汇率（1 USD = X Custom）
	CustomCurrencyExchangeRate float64 `json:"custom_currency_exchange_rate"`
	// Request-level upstream override for relay requests
	RequestUpstreamOverrideEnabled   bool              `json:"request_upstream_override_enabled"`
	RequestUpstreamOverrideAllowlist []string          `json:"request_upstream_override_allowlist"`
	RequestUpstreamProxyMap          map[string]string `json:"request_upstream_proxy_map"`
}

// 默认配置
var generalSetting = GeneralSetting{
	DocsLink:                   "https://docs.newapi.pro",
	PingIntervalEnabled:        false,
	PingIntervalSeconds:        60,
	QuotaDisplayType:           QuotaDisplayTypeUSD,
	CustomCurrencySymbol:       "¤",
	CustomCurrencyExchangeRate: 1.0,
	RequestUpstreamOverrideEnabled:   true,
	RequestUpstreamOverrideAllowlist: []string{"*"},
	RequestUpstreamProxyMap:          map[string]string{},
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("general_setting", &generalSetting)
}

func GetGeneralSetting() *GeneralSetting {
	if value := os.Getenv("REQUEST_UPSTREAM_OVERRIDE_ENABLED"); value != "" {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			common.SysError("failed to parse REQUEST_UPSTREAM_OVERRIDE_ENABLED: " + err.Error())
		} else {
			generalSetting.RequestUpstreamOverrideEnabled = enabled
		}
	}
	if value := os.Getenv("REQUEST_UPSTREAM_OVERRIDE_ALLOWLIST"); value != "" {
		var allowlist []string
		if err := json.Unmarshal([]byte(value), &allowlist); err != nil {
			common.SysError("failed to parse REQUEST_UPSTREAM_OVERRIDE_ALLOWLIST: " + err.Error())
		} else {
			generalSetting.RequestUpstreamOverrideAllowlist = allowlist
		}
	}
	if value := os.Getenv("REQUEST_UPSTREAM_PROXY_MAP"); value != "" {
		var proxyMap map[string]string
		if err := json.Unmarshal([]byte(value), &proxyMap); err != nil {
			common.SysError("failed to parse REQUEST_UPSTREAM_PROXY_MAP: " + err.Error())
		} else {
			generalSetting.RequestUpstreamProxyMap = proxyMap
		}
	}
	return &generalSetting
}

func IsRequestUpstreamOverrideEnabled() bool {
	return generalSetting.RequestUpstreamOverrideEnabled
}

func GetRequestUpstreamOverrideAllowlist() []string {
	return generalSetting.RequestUpstreamOverrideAllowlist
}

func GetRequestUpstreamProxyMap() map[string]string {
	return generalSetting.RequestUpstreamProxyMap
}

// IsCurrencyDisplay 是否以货币形式展示（美元或人民币）
func IsCurrencyDisplay() bool {
	return generalSetting.QuotaDisplayType != QuotaDisplayTypeTokens
}

// IsCNYDisplay 是否以人民币展示
func IsCNYDisplay() bool {
	return generalSetting.QuotaDisplayType == QuotaDisplayTypeCNY
}

// GetQuotaDisplayType 返回额度展示类型
func GetQuotaDisplayType() string {
	return generalSetting.QuotaDisplayType
}

// GetCurrencySymbol 返回当前展示类型对应符号
func GetCurrencySymbol() string {
	switch generalSetting.QuotaDisplayType {
	case QuotaDisplayTypeUSD:
		return "$"
	case QuotaDisplayTypeCNY:
		return "¥"
	case QuotaDisplayTypeCustom:
		if generalSetting.CustomCurrencySymbol != "" {
			return generalSetting.CustomCurrencySymbol
		}
		return "¤"
	default:
		return ""
	}
}

// GetUsdToCurrencyRate 返回 1 USD = X <currency> 的 X（TOKENS 不适用）
func GetUsdToCurrencyRate(usdToCny float64) float64 {
	switch generalSetting.QuotaDisplayType {
	case QuotaDisplayTypeUSD:
		return 1
	case QuotaDisplayTypeCNY:
		return usdToCny
	case QuotaDisplayTypeCustom:
		if generalSetting.CustomCurrencyExchangeRate > 0 {
			return generalSetting.CustomCurrencyExchangeRate
		}
		return 1
	default:
		return 1
	}
}
