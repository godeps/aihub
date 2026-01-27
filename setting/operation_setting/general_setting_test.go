package operation_setting

import "testing"

func TestGeneralSettingEnvOverrides(t *testing.T) {
	t.Setenv("REQUEST_UPSTREAM_OVERRIDE_ENABLED", "true")
	t.Setenv("REQUEST_UPSTREAM_OVERRIDE_ALLOWLIST", "[\"example.com\"]")
	t.Setenv("REQUEST_UPSTREAM_PROXY_MAP", "{\"example.com\":\"http://proxy:8080\"}")

	setting := GetGeneralSetting()
	if !setting.RequestUpstreamOverrideEnabled {
		t.Fatalf("expected override enabled")
	}
	if len(setting.RequestUpstreamOverrideAllowlist) != 1 || setting.RequestUpstreamOverrideAllowlist[0] != "example.com" {
		t.Fatalf("unexpected allowlist: %#v", setting.RequestUpstreamOverrideAllowlist)
	}
	if setting.RequestUpstreamProxyMap["example.com"] != "http://proxy:8080" {
		t.Fatalf("unexpected proxy map: %#v", setting.RequestUpstreamProxyMap)
	}
}

func TestGeneralSettingEnvOverridesInvalidJSON(t *testing.T) {
	t.Setenv("REQUEST_UPSTREAM_OVERRIDE_ALLOWLIST", "invalid")
	t.Setenv("REQUEST_UPSTREAM_PROXY_MAP", "invalid")

	setting := GetGeneralSetting()
	if setting == nil {
		t.Fatal("expected setting")
	}
}
