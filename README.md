# Demo-Vue

> 基于 Go (Gin) + Vue 3 的全栈管理后台脚手架，前端静态资源通过 `go:embed` 嵌入编译为单二进制部署。

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)
[![Vue Version](https://img.shields.io/badge/Vue-3.5+-4FC08D?logo=vue.js)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9+-3178C6?logo=typescript)](https://www.typescriptlang.org)
[![Vite](https://img.shields.io/badge/Vite-7.1+-646CFF?logo=vite)](https://vitejs.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](./LICENSE)

[![GitHub](https://img.shields.io/badge/GitHub-skyle1995/Demo--Vue-181717?logo=github)](https://github.com/skyle1995/Demo-Vue)
[![Releases](https://img.shields.io/badge/Releases-GitHub%20Releases-181717?logo=github)](https://github.com/skyle1995/Demo-Vue/releases)

> **仓库地址（GitHub）**：<https://github.com/skyle1995/Demo-Vue> ｜ **版本下载**：<https://github.com/skyle1995/Demo-Vue/releases>

---

## 功能特性

| 分类 | 特性 |
|---|---|
| 认证与鉴权 | JWT 双 Token（Access + Refresh）、一次一换、重用检测、会话绝对过期、密码变更自动失效 |
| 用户管理 | 超级管理员 / 子账号、账号状态（启用/禁用）、角色分级 |
| 系统设置 | 数据库驱动的键值配置、缓存层、运行时热更新 |
| 门户导航 | 树形导航结构、分组/链接、首页动态标记、排序与显隐控制 |
| 安全措施 | 图形验证码、CSRF 防护、登录日志、操作审计日志 |
| 中间件 | CORS、限速（令牌桶）、维护模式、开发模式快捷入口 |
| 前端框架 | Vue 3 Composition API + TypeScript、Element Plus、Tailwind CSS、ECharts、Pinia |
| 构建与部署 | 跨平台交叉编译（Windows / Linux / macOS）、`go:embed` 前端嵌入、一键构建脚本 |

---

## 技术栈

| 层 | 技术 |
|---|---|
| 后端框架 | [Gin](https://github.com/gin-gonic/gin) v1.11 |
| ORM | [GORM](https://gorm.io) v1.30 + SQLite / MySQL |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| 配置 | [Viper](https://github.com/spf13/viper) |
| 认证 | [golang-jwt](https://github.com/golang-jwt/jwt) v5 |
| 验证码 | [base64Captcha](https://github.com/mojocn/base64Captcha) |
| 日志 | [Logrus](https://github.com/sirupsen/logrus) + [lumberjack](https://github.com/natefinch/lumberjack) |
| 前端脚手架 | [pure-admin-thin](https://github.com/pure-admin/pure-admin-thin) v6.2 |
| UI 组件 | [Element Plus](https://element-plus.org) v2.11 |
| 图表 | [ECharts](https://echarts.apache.org) v6 |
| 构建 | Vite v7 + pnpm |

---

## 项目结构

```
Demo-Vue/
├── main.go                  # 入口：go:embed 前端 + Cobra 启动
├── cmd/                     # Cobra 命令（server 子命令）
├── server/                  # 路由注册、前端静态资源分发
│   ├── routes.go            # API 路由 + SPA 兜底
│   ├── admin.go             # 后台路由
│   ├── default.go           # 公开路由
│   └── install.go           # 安装引导路由
├── controllers/             # 控制器（请求处理）
│   ├── admin/               # 后台接口
│   │   ├── auth.go          # 登录/登出/刷新 Token
│   │   ├── dashboard.go     # 仪表盘统计
│   │   ├── settings.go      # 系统设置 CRUD
│   │   ├── portal_navigation.go  # 门户导航管理
│   │   ├── log.go           # 日志查询
│   │   ├── profile.go       # 个人资料
│   │   ├── captcha.go       # 验证码
│   │   └── utils.go         # 工具接口（如 Excel 导出）
│   ├── default/             # 公开接口
│   └── install/             # 安装向导
├── models/                  # 数据模型（GORM）
│   ├── user.go              # 用户
│   ├── settings.go          # 系统设置
│   ├── portal_navigation.go # 门户导航
│   ├── refresh_token.go     # 刷新令牌持久化
│   ├── login_log.go         # 登录日志
│   └── operation_log.go     # 操作日志
├── services/                # 业务逻辑层
│   ├── settings.go          # 设置缓存服务
│   ├── portal_navigation.go # 导航规范化
│   ├── refresh_token.go     # Token 族管理
│   ├── operation_log.go     # 操作日志记录
│   └── log_cleanup.go       # 日志清理
├── database/                # 数据库初始化与迁移
├── middleware/               # 中间件（CORS、限速、维护模式等）
├── utils/                   # 工具集（加密、验证码、日志、Excel）
├── config/                  # Viper 配置加载
├── constants/               # 应用常量（版本号等）
├── docs/                    # 技术文档
├── frontend/                # 前端工程（Vue 3 + Vite）
│   ├── src/
│   │   ├── views/           # 页面组件
│   │   ├── router/          # 路由配置
│   │   ├── store/           # Pinia 状态管理
│   │   ├── api/             # 后端 API 封装
│   │   ├── layout/          # 布局组件
│   │   ├── components/      # 公共组件
│   │   ├── utils/           # 前端工具函数
│   │   └── style/           # 全局样式
│   └── vite.config.ts       # Vite 构建配置
├── build.sh                 # 一键构建脚本
├── config.json              # 默认配置文件（本地调试用，不入库）
├── CHANGELOG.md             # 变更日志
└── LICENSE
```

---

## 快速开始

### 环境要求

| 依赖 | 版本 |
|---|---|
| Go | >= 1.26 |
| Node.js | ^20.19 \|\| >= 22.13 |
| pnpm | >= 9 |

数据库默认使用 **SQLite**（零配置），无需额外安装。可选 MySQL 或 Redis（见配置说明）。

### 克隆并启动

```bash
# 克隆仓库（GitHub）
git clone https://github.com/skyle1995/Demo-Vue.git
cd Demo-Vue

# 也可不克隆，直接从 GitHub Releases 下载对应平台压缩包
# https://github.com/skyle1995/Demo-Vue/releases

# 安装前端依赖
cd frontend && pnpm install && cd ..

# 构建前端（产物到 frontend/dist）
cd frontend && pnpm build && cd ..

# 启动后端
go run main.go server
```

访问 `http://localhost:8080` 进入管理后台。

### 默认配置

服务启动前请在项目根目录下准备 `config.json`：

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "dist": "",
    "dev_mode": true,
    "access_log": true
  },
  "database": {
    "type": "sqlite",
    "sqlite": { "path": "database.db" }
  },
  "log": {
    "level": "info",
    "file": ""
  }
}
```

- `server.dist`：设为空字符串使用 `go:embed` 嵌入前端；设为 `http://localhost:5173` 可反向代理到前端开发服务器。
- `database.type`：支持 `sqlite` 或 `mysql`。
- 首次启动后会自动建表与初始化默认设置。

---

## 构建部署

### 使用构建脚本

```bash
chmod +x build.sh
./build.sh
```

菜单选项：

- **自动构建 + 生成更新包**：自动迭代版本号、编译前端、交叉编译后端、打包 zip
- **全部构建**：前端 + 所有平台后端
- **仅编译前端**：单独构建 Vue 前端

### 手动交叉编译

```bash
# 编译前端
cd frontend && pnpm build && cd ..

# 编译后端（示例：Linux amd64）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath --ldflags="-s -w" -o dist/linux_amd64/Demo-Vue

# 产物即为包含前端的单二进制文件
```

### 产物结构

```
dist/
├── linux_amd64/Demo-Vue
├── linux_arm64/Demo-Vue
├── darwin_amd64/Demo-Vue
├── darwin_arm64/Demo-Vue
├── windows_amd64/Demo-Vue.exe
└── packages/               # zip 更新包（含 SHA256 校验）
```

---

## 认证机制

> 详细说明见 [docs/JWT-AUTH.md](./docs/JWT-AUTH.md)

- Access Token 有效期 2 小时（可配置），Refresh Token 有效期 7 天（可配置）
- 刷新时执行一次一换：旧 Token 即行撤销，分配新 Token
- 重复使用已撤销 Token 将触发整族销毁（防泄露攻击）
- 密码变更后所有现存 Token 立即失效
- 会话绝对过期上限 30 天，超时强制重新登录

---

## 配置项

通过在 Web 后台"系统设置"页面或数据库 `settings` 表动态配置（修改后即时生效）：

| 配置项 | 默认值 | 说明 |
|---|---|---|
| `jwt_secret` | 随机生成 | JWT 签名密钥 |
| `jwt_expire` | 2h | Access Token 有效期 |
| `refresh_token_expire_days` | 7d | Refresh Token 有效期 |
| `session_absolute_expire_days` | 30d | 会话绝对过期上限 |
| `refresh_advance_seconds` | 60s | 前端提前刷新阈值 |

---

## 开发说明

### 前端热重载

```bash
# 启动前端开发服务器
cd frontend && pnpm dev

# 配置 config.json 中 server.dist = "http://localhost:5173"
# 后端会将非 /api 请求反向代理到 Vite 开发服务器
go run main.go server
```

### 项目约定

- 强制使用 Composition API（`<script setup>`）与 TypeScript
- 控制器层仅做参数校验与响应，业务逻辑下沉至 `services/`
- 测试文件统一放入 `./demo-test/` 目录
- 提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org)

---

## 后台登录安全加固

后台登录已内置多层防护（本次安全修复），上线前请逐项确认：

- **验证码与 `dev_mode` 解耦**：`server.dev_mode: true` 不再跳过登录验证码；本地调试确需跳过时，显式设置 `security.skip_captcha: true`（启动会打印严重安全警告），生产环境必须为 `false`。
- **客户端 IP 不可伪造**：未配置 `server.trusted_proxies` 时禁用 `X-Forwarded-For`/`X-Real-IP` 解析，客户端 IP 取 TCP 对端地址。部署在 Nginx/CDN 之后时必须把代理出口地址加入白名单（如 `["127.0.0.1/32"]`），并让代理覆盖客户端可伪造的转发头，否则 IP 级限流会退化为全局共享额度、登录日志 IP 失真。
- **验证码令牌绑定客户端与类型**：滑块/点击验证码通过后签发的一次性令牌绑定签发时的 IP、User-Agent 与验证码类型，5 分钟内有效且只能使用一次；登录只接受与当前配置类型一致的令牌，未启用的验证码接口会直接拒绝（修复「配置 click 后仍可调用 slide 接口刷令牌」的绕过）。
- **限流 fail-secure**：Redis 不可用时自动降级为进程内限流，不再静默放行。
- **登录失败锁定**：同一账号连续失败达到阈值后临时锁定，重复触发按指数退避（最高 6 小时）；同一 IP 的失败（含验证码错误）也会触发临时锁定。锁定状态保存在进程内存，重启服务即清除。
- **建议**：在 Nginx/WAF 层对 `/api/admin/login`、`/api/admin/captcha/` 再叠加一层基于 `$binary_remote_addr` 的 `limit_req`。

对应配置（`config.json`，均可省略、使用安全默认值）：

```json
{
  "server": {
    "dev_mode": false,
    "trusted_proxies": []
  },
  "security": {
    "skip_captcha": false,
    "login_max_failures": 10,
    "login_window_minutes": 15,
    "login_lock_minutes": 15,
    "login_ip_max_failures": 30,
    "login_ip_lock_minutes": 10
  }
}
```

---

## Git 版本发布（Tag）

本项目的发布通常以 Tag 为触发点：推送 `v*` Tag 后，由 GitHub Actions（[.github/workflows/act_runner.yml](./.github/workflows/act_runner.yml)）自动构建并在 **GitHub Releases** 生成对应版本。建议使用语义化版本号（SemVer），并统一以 `v` 前缀命名，例如 `v1.2.3`。

- 仓库：<https://github.com/skyle1995/Demo-Vue>
- 工作流运行记录：<https://github.com/skyle1995/Demo-Vue/actions>
- 版本发布页：<https://github.com/skyle1995/Demo-Vue/releases>

### 发布一个 Tag 版本

```bash
# 0) 同步远端 Tag（推荐）
#    目的：避免本地残留 Tag（远端已删除）被误推回去并触发工作流
#    Git 版本需支持 --prune-tags（不支持时看下方兼容方案）
git fetch origin --tags --prune --prune-tags

# 1) 确保代码已提交且工作区干净
git status

# 2) 确认远端不存在同名 Tag（避免误复用）
git ls-remote --tags origin v1.2.3

# 3) 创建“附注（annotated）”Tag（推荐）
git tag -a v1.2.3 -m "release: v1.2.3"

# 4) 推送 Tag 到远端（推荐：仅推送本次发布的 Tag）
git push origin v1.2.3
```

注意：不建议直接使用 `git push origin --tags`。如果远端曾删除某些 Tag，但本地仍保留，这条命令会把这些 Tag 重新推送回远端，从而可能再次触发 Release/工作流。

如你的 Git 不支持 `--prune-tags`，可用“清空本地 Tag → 重新从远端拉取 Tag”的方式确保对齐后再发布：

```bash
git tag -l | xargs -n 1 git tag -d
git fetch origin --tags
```

```powershell
git tag | ForEach-Object { git tag -d $_ }
git fetch origin --tags
```

### 重新发布同名 Tag 版本（修正打错的 Tag）

同名 Tag 一般不建议复用；如确需修正（例如 Tag 指向了错误提交），请先删除远端 Tag，再重新创建并推送。

```bash
# 假设要修正的 Tag 为 v1.2.3

# 1) 删除本地 Tag
git tag -d v1.2.3

# 2) 删除远端 Tag
git push origin --delete v1.2.3

# 3) 在正确的提交上重建 Tag（<commit> 可省略，默认当前 HEAD）
git tag -a v1.2.3 <commit> -m "release: v1.2.3"

# 4) 重新推送到远端
git push origin v1.2.3
```

如果远端已基于该 Tag 生成了 Release（例如 GitHub Release），通常也需要同步删除并重新创建对应 Release，避免附件与版本信息不一致。GitHub 上的操作入口：仓库 → **Releases** → 选中该 Tag → **Delete release**（如需删除 Tag 本身，在 **Tags** 页面操作），之后重新推送 Tag 即会由 GitHub Actions 重新生成 Release。

- Releases 页面：<https://github.com/skyle1995/Demo-Vue/releases>
- Tags 页面：<https://github.com/skyle1995/Demo-Vue/tags>

---

## CI/CD Secrets 配置

以下变量需在 GitHub 仓库 **Settings → Actions → Secrets** 中配置（需登录且具备仓库管理权限）。

### 必选（缺一不可）

| Secret 名称 | 说明 |
|---|---|
| `RELEASE_TOKEN` | GitHub 访问令牌（classic PAT），权限需包含 `repo`（公开仓库用 `public_repo` 亦可）。用于自动创建 Release 并上传附件；未配置时回退使用内置 `GITHUB_TOKEN` |

### 对象存储上传（可选，未配置则跳过）

以下凭据需在 **Settings → Actions → Secrets** 中配置，其余配置项在 **Settings → Actions → Variables** 中配置。

#### 腾讯云 COS

**Secrets（凭据）**

| Secret 名称 | 必填 | 说明 |
|---|---|---|
| `COS_SECRET_ID` | 是 | API 密钥 SecretId |
| `COS_SECRET_KEY` | 是 | API 密钥 SecretKey |

**Variables（配置）**

| Variable 名称 | 必填 | 说明 |
|---|---|---|
| `COS_BUCKET` | 是 | 存储桶名称（不含 AppId 后缀） |
| `COS_REGION` | 是 | 地域，如 `ap-guangzhou` |
| `COS_UPLOAD_PATH` | 否 | 上传路径前缀，如 `releases`；留空则放桶根目录 |
| `COS_CLEAR_BEFORE_UPLOAD` | 否 | 上传前清空目标路径；默认 `false`，设为 `true` 启用 |

#### 阿里云 OSS

**Secrets（凭据）**

| Secret 名称 | 必填 | 说明 |
|---|---|---|
| `OSS_ACCESS_KEY_ID` | 是 | RAM 用户 AccessKey ID |
| `OSS_ACCESS_KEY_SECRET` | 是 | RAM 用户 AccessKey Secret |

**Variables（配置）**

| Variable 名称 | 必填 | 说明 |
|---|---|---|
| `OSS_BUCKET` | 是 | Bucket 名称 |
| `OSS_ENDPOINT` | 是 | Endpoint，如 `oss-cn-hangzhou.aliyuncs.com` |
| `OSS_UPLOAD_PATH` | 否 | 上传路径前缀，同上 |
| `OSS_CLEAR_BEFORE_UPLOAD` | 否 | 上传前清空目标路径；默认 `false`，设为 `true` 启用 |

> **优先级**：COS > OSS。两类同时配置时仅使用 COS。全部未配置则跳过存储上传（Release 仍正常创建）。

---

## License

MIT © 2025 skyle1995
