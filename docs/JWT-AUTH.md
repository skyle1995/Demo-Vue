# Demo-Vue JWT 认证机制说明

> 适用范围：`Demo-Vue` 项目（admin 后台 + user 门户）  
> 认证模式：OAuth2 风格的双 Token（Access + Refresh），纯 `Authorization: Bearer` 传输

---

## 1. 总览

| 项 | 取值 |
|---|---|
| 签名算法 | `HMAC-SHA256` |
| 传输方式 | `Authorization: Bearer <accessToken>` |
| 存储位置 | 前端 `localStorage`（key：`admin-info` / `user-info`） |
| Token 形态 | `accessToken`（短期）+ `refreshToken`（长期，可轮换） |
| 续期方式 | 前端按 `expires` 主动调用 `/admin/refresh-token` |
| 安全特性 | jti 唯一、family 分组、一次一换、重用检测、绝对过期上限、密码修改自动失效 |

> 不使用 HttpOnly Cookie 存储认证 Token；`csrf_token`、`captcha_id`、`multiple-tabs` 等 Cookie 与认证无关。

---

## 2. 数据库可配置项

全部位于 `settings` 表，由 [`database/SeedDefaultSettings`](../database/settings.go) 初始化：

| 配置 key | 默认值 | 含义 |
|---|---|---|
| `jwt_secret` | 启动时随机生成 | HS256 签名密钥 |
| `jwt_expire` | `2`（小时） | accessToken 有效期，建议 0.5~2 小时 |
| `refresh_token_expire_days` | `7`（天） | refreshToken 滑动有效期，每次刷新重置 |
| `session_absolute_expire_days` | `30`（天） | 会话绝对上限，refresh 滑动续期不能突破 |
| `refresh_advance_seconds` | `60`（秒） | 前端提前多少秒触发刷新 |
| `cookie_secure` / `cookie_samesite` / `cookie_domain` / `cookie_max_age` | 见 settings | 仅作用于 `csrf_token` 等非认证 Cookie |

修改后无需重启，`SettingsService` 会即时读取。`<= 0` 时自动回退到上述默认值。

---

## 3. 数据模型

### `refresh_tokens` 表（[`models/refresh_token.go`](../models/refresh_token.go)）

| 字段 | 说明 |
|---|---|
| `jti` | 每个 refreshToken 唯一标识，UUID v4 |
| `family_id` | 同一登录会话的 token 族；登录时新建，刷新时继承 |
| `user_uuid` / `user_type` | 归属用户；`user_type` ∈ `admin` / `user` |
| `issued_at` / `expires_at` | 本条 token 签发与过期时间 |
| `absolute_expires_at` | 整个会话绝对过期时间（family 共享） |
| `revoked` | 是否已撤销（轮换 / 登出 / 重用检测会置为 true） |
| `replaced_by` | 被哪个新 jti 替换（轮换审计用） |
| `user_agent` / `ip` | 登录环境，用于审计 |

### JWT Claims（[`controllers/admin/auth.go`](../controllers/admin/auth.go) 中 `JWTClaims`）

```jsonc
{
  "username": "admin",
  "uuid": "xxxx-xxxx",
  "role": 0,
  "password_hash": "<sha256(密码哈希)>",
  "typ": "access",         // "access" 或 "refresh"
  "fid": "<family_id>",    // 仅 refresh token 携带
  "jti": "<uuid>",         // 仅 refresh token 携带
  "exp": 1714800000,
  "iat": 1714792800,
  "nbf": 1714792800,
  "iss": "Demo",
  "sub": "admin"
}
```

---

## 4. 接口列表

| 路由 | 方法 | 鉴权 | 用途 |
|---|---|---|---|
| `/admin/login` | POST | 否 | 登录，签发 access + refresh |
| `/admin/refresh-token` | POST | 否（用 refreshToken 自证） | 轮换出新的 access + refresh |
| `/admin/logout` | POST | 否 | 撤销该用户全部 refreshToken |
| `/admin/*`（其余） | * | `AdminAuthRequired` 中间件 | 需要有效 access token |

---

## 5. 登录流程

```
前端                                    后端
 │  POST /admin/csrf                     │
 │ ◀──── csrf_token ─────────────────── │
 │                                        │
 │  POST /admin/login                    │
 │  {username,password,captcha,csrf_token}│
 │ ─────────────────────────────────────▶│
 │                                        │  1. 校验 CSRF / 验证码
 │                                        │  2. 校验账户存在 + Status=1 + Role≤1
 │                                        │  3. VerifyPasswordWithSalt
 │                                        │  4. generateJWTTokenForAdmin → access(2h)
 │                                        │  5. NewFamilyID + NewJTI
 │                                        │  6. generateRefreshTokenForAdmin → refresh(7d)
 │                                        │  7. RefreshTokenService.Create 持久化
 │                                        │  8. 记录登录日志
 │ ◀─── {accessToken, refreshToken,      │
 │       expires, username, role, ...} ──│
 │                                        │
 │  setToken() → storageLocal('admin-info')│
```

> 注意：access token **不写入 cookie**；前端只把它放在 `localStorage` 与 `Authorization` Header 中。

---

## 6. 鉴权流程（每次受保护请求）

[`AdminAuthRequired()` 中间件](../controllers/admin/auth.go) → `GetCurrentAdminUserWithRefresh()`：

1. 从 `Authorization: Bearer <token>` 读取 access token
2. `parseJWTToken`：校验签名 / `exp` / `nbf`
3. 校验 `typ == "access"`（旧 token 无该字段时按 access 视）
4. `validateAdminPasswordHash`：
   - DB 取最新用户
   - 用户存在 / `Status=1` / `Role<=1`
   - JWT 中 `password_hash` 与当前数据库密码 SHA256 一致 → 否则强制下线
5. 通过后将 `admin_uuid / admin_username / admin_role` 注入 `gin.Context`

> 中间件**不再做滑动续期**，所有续期统一走 `/admin/refresh-token`。

---

## 7. 刷新流程（OAuth2 一次一换 + 重用检测）

[`RefreshTokenHandler`](../controllers/admin/auth.go)：

1. 解析请求体 `{refreshToken}` → `parseJWTToken`
2. `claims.typ` 必须为 `refresh`，否则 `401`
3. `RefreshTokenService.FindByJTI(jti)`：
   - **未找到** → `401 刷新令牌不存在或已失效`（防伪造 / 已被清理）
4. **重用检测**：`rec.Revoked == true` →
   - `RevokeFamily(family_id)` 整族撤销
   - 返回 `401 检测到刷新令牌重用，会话已强制失效`
5. `now > rec.ExpiresAt` → `RevokeByJTI` → `401 已过期`
6. `now > rec.AbsoluteExpiresAt` → `RevokeFamily` → `401 已达最长有效期`
7. 再次 `validateAdminPasswordHash` —— 密码改了立刻失效
8. 签发新 token：
   - 新 access：`generateJWTTokenForAdmin`（有效期 = `jwt_expire` 小时）
   - 新 refresh：新 `jti` + **复用旧 `family_id`** + 新 `expires_at = min(now + refresh_token_expire_days, rec.AbsoluteExpiresAt)`
9. `RefreshTokenService.Create(newJTI, sameFamily, …)`
10. `RefreshTokenService.Rotate(oldJTI → newJTI)`：旧记录置 `revoked=true` 并写 `replaced_by`
11. 响应 `{accessToken, refreshToken, expires}`

### 关键不变量

- 同一 family 任意时刻最多只有 1 条 `revoked=false` 的 refresh token
- 滑动续期的 `expires_at` 永远不超过 `absolute_expires_at`
- 任何旧 refresh token 被再次使用 → 整族立即失效（典型攻击：refresh token 泄露）

---

## 8. 登出流程

[`LogoutHandler`](../controllers/admin/auth.go)：

1. 解析当前 access token（如有）取 `claims.UUID`
2. `revokeAllRefreshOfUser(uuid)`：将该用户全部 `user_type=admin` 且 `revoked=false` 的 refreshToken 标记为 `revoked=true`
3. 返回 `{redirect: "/admin/login"}`

> 前端 `removeToken()` 同时清理 localStorage 中的 `admin-info` 与残留的 `admin-authorized-token` cookie。

---

## 9. 前端触发逻辑

### 9.1 存储

[`frontend/src/utils/auth.ts`](../frontend/src/utils/auth.ts)：

- `setToken()`：写入 `storageLocal('admin-info')`，**不再写入 cookie**
- `getToken()`：优先 `storageLocal`，cookie 仅作旧会话兜底
- `removeToken()`：清理 localStorage + 残留 cookie

### 9.2 请求拦截器

[`frontend/src/utils/http/index.ts`](../frontend/src/utils/http/index.ts)：

```text
请求 → 拦截器
  ├── 白名单 (/login, /refresh-token) → 直接放行
  └── 其他请求
        ├── data.expires - now <= 0
        │     ├── 是 → 调 useUserStoreHook().handRefreshToken({refreshToken})
        │     │       新 token 回填到 storageLocal + 当前请求 Authorization
        │     │       同时通过 PureHttp.requests 队列重放被挂起的并发请求
        │     └── 否 → Authorization = Bearer <accessToken>
```

并发安全：`PureHttp.isRefreshing` 标志位 + `requests` 队列，确保同一时刻只发起一次刷新。

### 9.3 提前刷新（可选优化）

后端已暴露 `refresh_advance_seconds`（默认 60s）。前端拦截器若改为 `data.expires - now <= advance * 1000`，即可在过期前 N 秒主动刷新，避免请求第一次因 401 被拒后再重试。

---

## 10. 安全特性矩阵

| 风险 | 缓解 |
|---|---|
| accessToken 泄露 | 短期 (2h)，过期即失效 |
| refreshToken 泄露 | 一次一换 + 重用检测 → 攻击者一旦使用过的旧 token 就触发整族撤销 |
| 永久驻留 | `session_absolute_expire_days` (30d) 强制重新登录 |
| 密码被修改后旧 token 仍可用 | `password_hash` 写入 claim，每次校验对比当前数据库 |
| 用户被禁用后旧 token 仍可用 | `Status != 1` 立即拒绝 |
| 越权访问 admin | `Role > 1` 拒绝 |
| CSRF | 登录接口走独立 `csrf_token` cookie + Header / body 双源校验；其余接口因使用 Bearer Header 天然免疫 |
| 会话固定 | 每次登录新建 `family_id` |
| 暴力穷举 jti | jti = UUID v4，且必须命中 DB |
| 跨用户重放 | 校验 `claims.UUID` 必须匹配 DB 中存在且启用的用户 |

---

## 11. 维护建议

- **定期清理**：可由定时任务调用 `RefreshTokenService.CleanupExpired(retentionDays)` 删除已过期 N 天的记录
- **监控**：`revoked + replaced_by` 链可重建会话轨迹；批量 `RevokeFamily` 事件应触发安全告警（可能存在 token 泄露）
- **修改默认有效期**：直接更新 `settings` 表对应行，无需改代码
- **轮换 jwt_secret**：更新 `settings.jwt_secret` 后所有现存 access/refresh 立即失效，相当于全员强制下线
