# Request-Level Upstream Override 使用指南

本指南说明如何使用 add-request-upstream-override 以请求级方式指定上游地址与鉴权信息。

## 1. 启用功能
在系统配置中设置以下 key：

- `general_setting.request_upstream_override_enabled = true`
- `general_setting.request_upstream_override_allowlist = ["proxy.example.com","https://proxy2.example.com"]`

说明：
- 白名单只比较 host，可以写纯域名或完整 URL。
- 不在白名单内的 host 将被拒绝（403）。
- `*` 表示允许所有上游（默认配置已启用）。

### API 设置示例（管理员）

```bash
curl -X PUT "http://localhost:19000/api/option/" \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"key":"general_setting.request_upstream_override_enabled","value":"true"}'

curl -X PUT "http://localhost:19000/api/option/" \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"key":"general_setting.request_upstream_override_allowlist","value":"[\"*\"]"}'
```

## 2. 请求头
在原有 Gemini 接口基础上新增两个请求头：

- `X-Relay-Upstream-Base-URL`
  - 真实上游地址，例如：`https://proxy.example.com`
- `X-Relay-Upstream-Headers`
  - JSON 字符串，指定要透传给上游的鉴权头，例如：
    `{"Authorization":"Bearer <token>","x-goog-api-key":"<key>"}`

## 3. 示例

```bash
curl -X POST "http://localhost:19000/v1beta/models/gemini-3-pro-image-preview:generateContent" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-local-api-key>" \
  -H "X-Relay-Upstream-Base-URL: https://proxy.example.com" \
  -H "X-Relay-Upstream-Headers: {\"Authorization\":\"Bearer upstream-token\"}" \
  -d '{
    "contents":[{"role":"user","parts":[{"text":"draw a cat"}]}],
    "generationConfig":{"responseModalities":["TEXT","IMAGE"]}
  }'
```

## 4. 行为说明
- 覆盖逻辑在选好渠道之后执行：不会影响本项目配额、计费、限流。
- 如果 `X-Relay-Upstream-Headers` 包含 `Authorization` 或 `x-goog-api-key`，Gemini 适配器不会再写默认 key。
- 功能关闭或 host 不在白名单时，返回 403。

## 5. 常见错误
- `403 request upstream override is disabled`：功能未开启。
- `403 request upstream override not allowlisted`：host 不在白名单。
- `400 request upstream headers invalid`：`X-Relay-Upstream-Headers` 不是合法 JSON。
