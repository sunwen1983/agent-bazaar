# M0 探路结论

> 探路时间：2026-11-17（约 1 小时，一个 session 内完成，比计划的 11/17–11/23 提前）。
> 约束遵守：只探路、不写产品代码、不部署公开服务、不花钱、不注册任何账号、不改 mcpgate 仓库、不写任何真实 key 进文件。
> 计划对照：`plan.md` §7 待确认清单逐项见 §7；§8 决策点建议见 §8。

## 0. 总结论

**M1 可以按计划开工。** 官方示例在 served 与 static 两种模式下都跑通了，跨 origin 链接机制在构建层面已验证，四个待确认项每个都有"是/否 + 退路"结论。唯一的硬性缺口是 mcpgate 默认不发 CORS（已定位到一行中间件的修复点，M1 里做），以及默认模型 key 的保管方式需要 Wen 在 D4 上拍板（推荐方案已定，见 §2）。

| # | 验证项 | 结论 |
|---|--------|------|
| 1 | nsDeploy 跑通官方示例 + 静态包 + 三 origin 方式 | ✅ 通过（见 §1） |
| 2 | 默认模型 key 保管（R10，最大坑） | ⚠️ 定论已出：现成 nsDeploy 只有 BYOK；推荐演示 key + 余额封顶（见 §2） |
| 3 | mcpgate CORS | ❌ 默认不发；退路已定：给 mcpgate 加 CORS 中间件（见 §3） |
| 4 | NANDA 程序化接入与网站型注册 | ⚠️ resolve 接口可用；网站型注册无现成路径，B 档为主、A 档作加分项（见 §4） |
| 5 | Devpost 规则（旧项目 + 赛道） | ✅ 均合规；赛道不是表单项（见 §5） |
| 6 | Cloudflare Pages 托管 | ✅ 确认（见 §6） |

---

## 1. nsDeploy：验证通过 ✅

### 1.1 Clone 与 pin

- 仓库：`https://github.com/cognizant-ai-lab/nsdeploy`
- **Pin 住的 commit：`4e4cb36731ae71cbbb905d3fe98f26299ac3b4f2`**（2026-11-13，Merge PR #19）。clone 用时 22 秒。
- 仓库 2026-11-05 才建、0 star、变化快——M1 以此 commit 为准，不要追 main。

### 1.2 依赖安装

- 环境：Python 3.12.3，VM 上没有 `uv`，用 `pip`。
- `python3 -m venv .venv && .venv/bin/pip install -r requirements.txt`：约 2 分钟，装了约 90 个包（pyhocon 0.3.63、mcp 1.30.0、neuro-san 0.7.24、langchain 1.5.0、langchain-openai、langgraph、fastapi 0.163.1、uvicorn、httpx、pydantic 等），venv 共 216 MB。`import neuro_san, langchain, mcp, fastapi, pyhocon` 全部通过。

### 1.3 Served 模式跑官方旅行示例

- 命令：`.venv/bin/python -m nsdeploy examples/travel_agent_system.hocon --port 8080`
- **第一次失败是环境原因，不是 nsDeploy 的问题：端口 8080 已被占用**——占用者是 Wen 自己的 mcpgate（pid 2270，正在跑的 demo 实例）。换 `--port 8091` 后正常启动并持续运行。
- 生成产物（`.nsdeploy/examples/travel_agent_system/`）：`index.html`（832 B 薄壳，`import { boot } from "/runtime/index.js"`）、`agent-manifest.json`、`travel_agent_system.network.json`（106 KB）、`links.json`、`agent-mcp.json`、`.well-known/nsdeploy.json`、`nsdeploy-agent.json`、`theme.css/json`、`site.js`。
- `GET /` → 200。`/byok/v1/models` 无 key 时 → 401（符合预期：BYOK 代理拒绝无 key 请求）。`/v1/models` 无 `--mock-llm` 时 → 404（mock router 只有加 `--mock-llm` 才挂载）。
- 生成日志要点：
  - 网络 10 个 agent；2 个可托管工具（`ddgs_search`、`wikimedia_media_search`）以本地 MCP server 形式起在随机端口；
  - **6 个 agent 是 DEMO-MODE（instructions 里写明编造看似真实的答案，不接地气数据）**；
  - 2 个工具缺失被丢弃（`web_fetch`、`webpage_rag`）；
  - 外部 MCP `https://api.you.com/mcp?profile=free` 经同源代理（served 模式）接入；
  - 生成时自动从 neuro-san-studio 拉依赖文件，并自动 `pip install ddgs`。
- **注意**：dev server（Python BaseHTTP）对带 `Origin` 的请求不回任何 `Access-Control-Allow-Origin`。这是本地开发服务器的行为，不影响部署（部署用静态包 + Pages `_headers`），但本地多 origin 联调时要直接用静态包或 gateway，不要用 dev server 测跨域。

### 1.4 Static bundle（静态包）

- 命令：`.venv/bin/python -m nsdeploy examples/travel_agent_system.hocon --static-bundle --out ./dist/travel --site-url https://travel.example.com`，退出码 0。
- 产物（`dist/travel/`，共 5.9 MB）：
  - `runtime/`：完整 vendored 浏览器 runtime（`engine.js`、`kernel.js`、`runner.js`、`mcp.js`、`byok.js`、`llm.js`、`hocon.js`、`sandbox-worker.js`、`trace.js`、`ui.js`、`web-crypto.js` 等）；
  - `agent-manifest.json`、`travel_agent_system.network.json`、`links.json`、`agent-mcp.json`、`.well-known/`、`index.html`、`site.js`、`theme.*`；
  - **`_headers`**：`/agent-manifest.json`、`/travel_agent_system.network.json`、`/links.json`、`/agent-mcp.json`、`/.well-known/nsdeploy.json`、`/.well-known/nsdeploy-agent.json` 六个文件配 `Access-Control-Allow-Origin: *`；`/runtime/sandbox-worker.js` 配严格 CSP。另有 `.htaccess` 与 `STATIC_HOSTING.md`。
  - `STATIC_HOSTING.md` 原文要点："The site, browser runtime, and browser-local links require no application server." / "Real BYOK calls go directly from the browser to the router selected in settings, so that router must allow browser CORS. External MCP servers are also contacted directly and must allow browser CORS. Python tools and the public MCP facade are not included."
- 用普通 `python -m http.server 8092` 伺服该目录，13 个关键资源全部 200 且 Content-Type 正确（`/`、manifest、network JSON、links、agent-mcp、runtime 的 index/engine/mcp/byok、.well-known、theme.css、site.js）。
- **关键推论**：静态包里**没有 Python coded tool**（日志明说 omitted：`ddgs_search`、`web_fetch`、`webpage_rag`、`wikimedia_media_search`），也没有 public MCP facade 与 private steward。所以静态部署下，工具落地只能走外部 MCP 服务器——正好就是计划的 mcpgate 架构，计划方向被官方实现印证。
- dist 版 manifest 的 `llm_router`：`base_url: https://api.openai.com/v1`（浏览器直调 OpenAI 兼容 endpoint），`mock: false`。

### 1.5 跨 origin 链接机制（M4 的核心机制）

- 用 `--link` 生成了两个互相链接的静态包：
  - `link_demo_travel_agency`（`--site-url https://agency.example.com`，`--link link_demo_airline=https://airline.example.com/link_demo_airline`）
  - `link_demo_airline`（`--site-url https://airline.example.com`）
- 验证结果：
  - agency 的 `links.json`：`{"link_demo_airline": "https://airline.example.com/link_demo_airline"}`——符号引用到 peer URL 的映射正确生成；
  - agency manifest 的 `linking.imports[0]`：`status: "linked"`，`target` 为 peer URL；
  - airline 包正常发布自己的 manifest 与网络。
- 结论：跨 origin peer 绑定机制在**构建层面已验证**。浏览器里实际加载（fetch + 同意弹窗 + 执行）因本环境无浏览器未能实测，是 M1 第一项待验证（见 §9）。

### 1.6 未做的事与原因

- `commerce` / `cailville` 示例没跑：travel + link_demo 已覆盖所需机制（served、static、跨 origin 链接），同一套代码路径，没必要重复。
- 浏览器内真实跑通（点一下、看 agent 执行）未能验证：本环境无浏览器。M1 第一件事就是在真实浏览器里打开 served 示例（`--mock-llm` 可无 key 冒烟，或 BYOK 自带 key）。

---

## 2. 默认模型 key 保管（R10）：定论

### 2.1 (a) 纯 BYOK 是现成 nsDeploy 的唯一默认路径 —— 已证实

- README 原文（"BYOK guarantee"）："A served site never holds or uses an LLM key configured on the host. `/v1` is a keyless local mock that makes no network calls; `/byok/v1` forwards exactly the `Authorization` header the visitor entered to `--upstream` and nothing else; the public `/mcp` facade refuses real requests without caller BYOK headers"。
- "By default the site requires the visitor's own API key. … Without a key, the UI links to free-tier providers, fills in their recommended model settings, and does not make an LLM request."
- 实测：`/byok/v1/models` 无 key → 401；静态包 manifest 的 `llm_router.base_url` 直接指向 `https://api.openai.com/v1`（浏览器直调）。
- 结论：不改代码的前提下，评委打开网站必须自己贴 key 才能跑真实模型，否则只能跑 `--mock-llm` 的确定性假模型（且 mock 会在 UI 里被明确标出，不会假装是真模型）。

### 2.2 (b) 前端嵌演示 key + DeepSeek 余额封顶 —— 技术上可行 ✅

- **浏览器直调 DeepSeek 已实测可行**：对 `https://api.deepseek.com/v1/chat/completions` 发 preflight（`Origin: https://example.com`，`Access-Control-Request-Method: POST`，`Access-Control-Request-Headers: authorization,content-type`）返回 200，且带 `access-control-allow-methods: POST`、`access-control-allow-headers: authorization,content-type`、`access-control-allow-origin: https://example.com`、`access-control-allow-credentials: true`。OpenAI 的 `/v1/models` 同样回 `Access-Control-Allow-Origin: *`。所以浏览器带 `Authorization` 调 DeepSeek/OpenAI 都畅通。
- **实现路径很小**：nsDeploy 的 BYOK profile 存在浏览器 localStorage，key 名 `nsweb-byok-profiles-v1`（`nsdeploy/runtime/byok.js`，`upsertByokProfile`）。做法是在集市站的 `site.js` 里加一个 "Use demo key" 按钮（或首次访问时 seed 一个 Demo profile），几行代码。key 不进 manifest、不进网络定义（规范禁止，manifest 也没有 key 字段），只进浏览器 localStorage——与 BYOK 机制本身一致。
- **风险与上限**：key 在页面源码/localStorage 里是公开的，任何人都能拿去直接调 DeepSeek。唯一的硬上限是 DeepSeek 账户的预付余额：保持余额很小（如 $5），比赛结束即轮换 key。DeepSeek 不提供单 key 限流，所以余额就是 cap。按 demo 的调用量（每次组队约 10–20 次短调用，DeepSeek-V3 级约 $0.3–$1/百万 token），$5 够几百次完整 demo。
- **对评委体验的影响**：零摩擦——打开即用，和视频里一模一样。这是三个选项里评委体验最好的。

### 2.3 (c) 自建 token broker —— 可行，但没必要，除非 Wen 坚持零 key 暴露

- 工作量估计：约 0.5–1 天。约 150–300 行代码（Go 或 Python 小服务）+ 部署 + 防滥用（按 IP 限流、每日上限、prompt 长度上限、只允许指定模型）+ 评审周的运维值守。另需在提交文字说明里诚实描述它是一个"发布者运营的推理代理"，不是 BYOK。
- 规范层面：规范明确允许"publisher-funded"推理（"Publisher-funded and user-funded model access both fit this profile"），但 §4 的同源转发 endpoint 规则禁止**在 BYOK 转发路径里替换成发布者的 key**。独立的 broker 服务是另一个组件，不违反那条规则，但必须在提交材料里写清楚，不能把它包装成 BYOK。
- 对评委体验的影响：与 (b) 相同（零摩擦），但多一个 Wen 要运维的服务，评审周挂了就是直接丢分（R9）。

### 2.4 推荐

- **Live 站默认路径选 (b)：演示 key + DeepSeek 余额封顶（约 $5）+ BYOK 并存，赛后轮换。** 视频用 Wen 自己的 key 走 BYOK（镜头外输入），保证画面是真实模型行为。
- **(a) 纯 BYOK 作退路**：如果 Wen 不接受任何 key 暴露，就退到纯 BYOK，代价是评委体验变差。
- **(c) 只有 Wen 明确要求零 key 暴露时才做。**
- 补充论据（Devpost 规则原文，§5）：Testing 条款要求项目在评审期结束前 "available free of charge and without any restriction"。纯 BYOK 下，没有 key 的评委无法测试，等于是有门槛——这进一步支持 (b)/(c) 而不是 (a)。另 Functionality 条款要求项目 "function as depicted in the video"——视频里跑的是真模型，live 站也必须是真模型，不能一边 mock 一边录。

---

## 3. mcpgate CORS：默认不发（缺口已确认）❌

- **证据**：对 `~/workspace/mcpgate/` 全部 `*.go` 文件 grep `cors` / `Access-Control-Allow` / `OPTIONS` / `preflight`——零命中。`internal/gateway/gateway.go` 里 `http.NewServeMux()` 直接注册路由，没有任何中间件；`Gateway.ServeHTTP` 直接转给 mux。
- **后果**：浏览器从 origin A 向 `POST https://<mcpgate>/mcp/{backend}` 发请求（带 `X-API-Key` + `Content-Type: application/json`）会触发 CORS preflight（OPTIONS）。`handleMCP` 要求 POST，非 POST 直接 405，且没有任何 CORS 头——浏览器会拦掉整个调用。**浏览器直连 mcpgate 目前是走不通的。**
- **修复点（未改代码，只定位）**：在 `Gateway.ServeHTTP` 里包一层 CORS 中间件（或在 `cmd/mcpgate/main.go` 里包 `gw`）：OPTIONS 预检直接回 200 + `Allow-Origin/Allow-Methods/Allow-Headers`，其余请求补 `Access-Control-Allow-Origin`；中间件必须跑在 auth 之前（auth 在各 handler 内部先执行，无 key 的预检会被 401）。约 40 行 + 1 个测试，估计 1 小时内。
- **退路比较**：
  - (a) 给 mcpgate 加 CORS：**推荐**。改动点唯一、无新基建、且是 Wen 自己的项目，顺手可作"分享给社区的额外工具"（Potential Impact 加分）。
  - (b) 集市站同源代理（Pages Functions 转发到 mcpgate）：多一个部署物（Functions）+ 转发逻辑 + 延迟；且 demo key 本来就是公开的，代理藏不住任何 secret，没有安全收益。**不推荐。**
- **作用域建议**：CORS 只开给 `/mcp/` 与 `/healthz`，不要开给 `/admin/*`；`Access-Control-Allow-Origin: *` 对公开只读 demo key 是合适的；`Allow-Headers` 需含 `X-API-Key`、`Content-Type`，`Allow-Methods` 含 `POST, GET, OPTIONS`。

---

## 4. NANDA：resolve 可用，网站型注册无现成路径 ⚠️

### 4.1 程序化查询：可行 ✅

- `GET https://api.nandaindex.org/health` → `{"status":"ok","db":"ok"}`，服务活着。
- `GET https://api.nandaindex.org/api/v1/resolve?locator=urn:air:host39.org:personal:alicegmailcom@inkboxmail.com`（文档里的示例 locator，不存在）→ 404 JSON，但响应带 `access-control-allow-origin: https://example.com`。**resolve 接口对浏览器发 CORS 头，程序化查询可行**；不存在的 locator 返回明确的 404（fail closed 可观测）。
- `https://my.nandaindex.org/` 正常返回 200 HTML（注册前端）。

### 4.2 注册路径：两档，网站型网络没有现成类型

来源：`projnanda/nanda-index-v2` README 与 NANDA 官方说明。

| 路径 | 要求 | 适用性 |
|------|------|--------|
| 个人 agent（Personal Agent） | Inkbox 给 agent 一个专用邮箱；标识符锚定到 `host39.org`，email 作 `subjectAccount`；完成邮箱验证；**不需要自己的域名** | 注册的是"个人助理型 agent"，卡片托管在 host39.org；与我们网站上 HOCON 网络的绑定很弱（注册表记录不会以密码学方式绑定到我们的 manifest），能否用来注册网站型网络**未证实** |
| 域名锚定（org / DNS SVCB / SMB） | 注册 org；索引返回 TXT 记录 `_nanda-challenge.<domain>` = `nanda-verify=<token>`；把记录发布到该域名的 DNS；再调 `POST /api/v1/orgs/<org_id>/verify-domain`；之后标识符形如 `urn:air:<domain>:agent:<slug>`，解析链指向该域名的 AI Catalog / agent 卡 | 绑定强：记录可指向我们网站的 manifest；但**需要 Wen 控制一个自定义域名**（`*.pages.dev` 子域加不了 TXT 记录），且需要 Wen 的邮箱走注册流程 |

- 另：NANDA Index 是"一注册表一条记录"的 quilt 模型，不托管 agent，只告诉你去哪找；`sm-bridge`（PyPI）是现成的接入库，但那是给有自家注册表的源用的，与我们关系不大。
- **没有注册任何东西**（按约束；注册需要 Wen 的邮箱与拍板）。

### 4.3 D6 建议

- **B 档（本地签名验证 + 四步可见验证）作主线**：确定性强、不依赖外部服务、demo 稳。第 ③ 步显示"本地信任"，文字说明里诚实写明与 NANDA 机制的对应关系与后续接入路径。
- **A 档作加分项**：只有 Wen 决定做才做——需要他定一个自定义域名（A2 路径）并用自己的邮箱走注册。resolve 接口与 CORS 已证实可用，所以一旦注册完成，UI 上把第 ③ 步从"本地信任"换成一次真实 resolve 调用只是小改动。
- 计划 §5.7 的耦合判断被证实：**A 档 → 自定义域名是 must-have，不是 nice-to-have**（除非个人路径被证实可用，那仍需 M1 实测）。

---

## 5. Devpost 规则：两件事都确认了 ✅

来源：`https://cognizant-agenticweb.devpost.com/rules` 全文（2026-11-17 读取）。

### 5.1 旧项目作组件是否合规 —— 合规 ✅

- **"New & Existing" 条款原文**："Projects must be either newly created by the Entrant or, if the Entrant's Project existed prior to the Hackathon Submission Period, must have been significantly updated after the start of the Hackathon Submission Period. Entrants should explain how their Project was significantly updated during the Submission Period."
- **IP 条款原文**："An Entrant may submit a Submission that includes the use of open source software or hardware, provided the Entrant complies with applicable open source licenses and, as part of the Submission, creates software that enhances and builds upon the features and functionality included in the underlying open source product."
- 结论：Agent Bazaar 本身是比赛期内新建的项目；mcpgate 是 Wen 自己的 MIT 开源项目（v0.1.0），作工具后端组件完全合规。提交文字说明里披露一句话即可，例如："Agent Bazaar (built Nov 2026–Jan 2027) uses mcpgate, the author's pre-existing MIT-licensed MCP gateway (v0.1.0), as its tool backend; the CORS support added during the contest period is itself an enhancement shared back to that project." 这样同时满足"显著更新"的说明要求（CORS 改进就是比赛期内对 mcpgate 的增强）。

### 5.2 提交表单是否单选赛道 —— 不是表单项 ✅

- 对规则页全文搜 "track"：**零命中**。Submission Requirements 只列：网站 URL（含 5 项内容）、文字说明、视频、截图、可选的 Otari 邮箱。Track 1/2 只出现在 worldwideagenticweb.org 的主题描述里，不是 Devpost 表单项。
- 结论：计划里"同时踩两条"只是文字说明里的定位写法，无需在表单里选，不影响合规。

### 5.3 规则里与计划相关的其他确认

- 评审分两段：Stage One pass/fail（是否符合提交要求）；Stage Two 四项等权评分，平局先比 Technological Implementation。——与计划 §2.3 一致。
- Testing 条款："The Entrant must make the Project available free of charge and without any restriction, for testing, evaluation and use … until the Judging Period ends."（评审期 2027-01-11–01-18 PT）——支持 §2.4 的 (b) 推荐。
- Functionality 条款："function as depicted in the video and/or expressed in the text description."——视频必须拍真实行为。
- Submission Period：2026-11-11 8:00 pm PT – 2027-01-11 8:00 am PT（= 10:00 AM CT，与计划一致）；Winners ~2027-01-19。

---

## 6. Cloudflare Pages 托管：确认 ✅

- 官方文档（`developers.cloudflare.com/pages/configuration/headers/`）：`_headers` 是放在输出目录根部的纯文本文件，规则如 `/*` + `Access-Control-Allow-Origin: *`；文档里的 CORS 示例与我们要配的完全一致。上限：100 条规则，每行 2000 字符——我们的 _headers 只有 7 条规则，远未超。
- nsDeploy 生成的 `_headers` 可直接用：文件在 bundle 根部，路径均为 host-root-relative，与 Pages 的要求一致。
- 三个 Pages 项目 = 三个 origin（各一个 `*.pages.dev` 子域，HTTPS 默认开启）。免费档：不限站点数、每月 500 次构建。
- 注意点：
  - 单文件 25 MiB 上限——我们的 bundle 共约 6 MB，单文件最大是 network JSON（106 KB），无压力。
  - `_headers` 只作用于静态资源响应——我们全站都是静态的，正好。
  - 对比：GitHub Pages 不支持自定义响应头，**做不了跨 origin manifest 加载**，已排除；Vercel 可用 `vercel.json` 配 header，也是备选。Cloudflare Pages 仍是推荐（nsDeploy 的 `_headers` 在它是原生的）。
  - 如果将来走 (b) 同源代理路线才需要 Pages Functions——按 §3 推荐的 CORS 中间件路线，不需要。

---

## 7. 计划 §7 待确认清单逐项对照

| # | 待确认项 | 结论 | 证据 |
|---|----------|------|------|
| 1 | nsDeploy 嵌入方式与示例运行步骤 | ✅ 已确认：`pip install -r requirements.txt` 后 `python -m nsdeploy <hocon> --port <n>`；静态包用 `--static-bundle --out <dir> [--site-url <origin>]`；`--link` 绑 peer | §1.2–1.5 |
| 2 | NandaIndex 网站型 agent 注册与程序化解析是否可行 | ⚠️ 解析可行（resolve + CORS 已实测）；网站型注册无现成类型，B 档为主、A 档加分 | §4 |
| 3 | mcpgate 默认 CORS 行为 | ❌ 默认不发；退路已定：加 CORS 中间件（§3），同源代理不推荐 | §3 |
| 4 | 静态托管三家的 CORS 头配置方法 | ✅ Cloudflare Pages `_headers` 已确认（官方文档 + nsDeploy 生成文件可直接用）；GitHub Pages 不可行 | §6 |
| 5 | Otari 用还是不用（D8） | 见 §8：推荐不用 | §8 |
| 6 | 默认推理模型与月预算上限（D4） | ⚠️ 需 Wen 拍板；推荐 (b) 演示 key + 余额约 $5 | §2.4 |
| 7 | 默认推理模型的 key 保管方案（R10） | ✅ 定论已出（§2）；实现细节见 §2.2 | §2 |
| 8 | Devpost 规则是否要求作品比赛期内新建 | ✅ 不要求；旧项目作组件合规，需在文字说明里披露 | §5.1 |
| 9 | 提交表单是否要选赛道 | ✅ 不是表单项；Track 只是主题方向 | §5.2 |
| 10 | NANDA A 档与自定义域名耦合 | ✅ 已证实：A 档（域名锚定）需要自定义域名 + DNS TXT；个人路径绑定弱、未证实 | §4.2–4.3 |
| 11 | nsDeploy pin 一个 commit | ✅ 已 pin：`4e4cb36731ae71cbbb905d3fe98f26299ac3b4f2`（2026-11-13） | §1.1 |

风险对照（§7 表格）：R1 已消除（嵌入方式验证通过，退路"自写 runtime"不需要了）；R2 退到 B 档；R3 退路已定（加 CORS）；R4 缓解方案见 §2.2（余额封顶）；R6 缓解已验证（_headers）；其余 R5/R7/R8/R9/R10 不变，R10 的方案已定、只等 D4 拍板。

---

## 8. D1–D10 建议

| # | 决策 | 建议 | 是否可直接定 |
|---|------|------|--------------|
| D1 | Demo 场景 | **A. Austin 三日游**。注意：官方旅行示例本身就是一个 10-agent 旅行 concierge（且是 demo_mode），我们的差异化必须是集市与验证层，不是旅行内容本身；A 仍是最合适的（天然对应 flights/stays/guide 三个 peer） | 可直接定，等 Wen 点头 |
| D2 | Runtime 路线 | **A. 直接用 nsDeploy**（已验证，pin commit `4e4cb36`）。B 不需要了 | 可直接定 |
| D3 | 托管 | 静态站 **Cloudflare Pages**（已验证）；mcpgate 托管 fly.io vs Railway——M0 未测（取决于 Wen 已有账号/偏好），建议 fly.io（计划推荐），由 Wen 确认 | 静态站可定；mcpgate 托管需 Wen 选 |
| D4 | 推理与预算 | **推荐 (b)：演示 key + DeepSeek 余额封顶约 $5 + BYOK 并存**；视频用 Wen 自己的 key 走 BYOK（镜头外输入）。月上限金额由 Wen 定 | **需 Wen 选**（方案已给，只差金额与点头） |
| D5 | 工具数据 | **A. Mock 确定性数据**。官方示例的 agent 本来就是 demo_mode 编造答案，确定性 mock 是它的诚实版 | 可直接定 |
| D6 | NANDA 做到哪档 | **B 档为主，A 档作加分项**（见 §4.3）。A 档需要 Wen 定域名 + 自己走注册 | B 档可直接定；A 档需 Wen 决定做不做 |
| D7 | 项目名与 repo 名 | **`agent-bazaar`**（无 M0 阻碍） | 需 Wen 点头（名字是他的） |
| D8 | 用不用 Otari | **推荐不用**。M0 没发现 Otari 有额外奖项或接入优势；用了反而多一份提交材料（邮箱 + 问题反馈）。只有发现明确奖项时才回头 | 可直接定（除非 Wen 想试） |
| D9 | 视频 | 需 Wen 定（工具 + 时间）。建议手机 + 录屏，M4 周末整块；注意元旦周需提前确认家里时间 | 需 Wen 定 |
| D10 | 每周投入 | 推荐节奏不变。M0 实际只用约 1 小时，比计划一周的窗口提前很多，排期有 slack | 可直接定 |

**可以直接定的**：D1(A)、D2(A)、D3-静态站、D5(A)、D6-B 档、D7、D8(不用)、D10。
**仍需 Wen 选的**：D3-mcpgate 托管（fly.io/Railway）、D4（方案推荐 b，只差金额与点头）、D6-A 档做不做（+ 域名）、D9（视频工具与时间）。

---

## 9. M1 是否能按计划开工

**可以。** M0 退出标准（"能在浏览器里跑起一个官方示例网络；四个待确认项每个都有是/否 + 退路结论；决策清单签字"）中，前两项已达成；决策清单在 §8，等 Wen 签字后即进 M1。

M1 前三步建议顺序：

1. 建 `agent-bazaar` 仓库（D7），nsDeploy 以 pin 住的 commit 为准（`4e4cb36`），用 CLI 生成第一个网络再改 `sites/<network>/` 做定制——不要手写 runtime。
2. mcpgate CORS 中间件小 PR（§3 的修复点，<1 小时），与集市壳子并行。
3. 第一次端到端先在**真实浏览器**里跑 served 示例（`--mock-llm` 无 key 冒烟，再 BYOK 真 key），验证同意弹窗 + 跨 origin peer 加载 + MCP 调用——这是 M0 因无浏览器而唯一没验证过的环节。

带入 M1 的已知风险：R10 只等 D4 拍板；R6 的三 origin 实测需在 M1 用真实 Pages 部署验证一次；nsDeploy 的学习曲线比预想的小（CLI 做了大部分工作，定制点是 HOCON + `site.js`）。

---

## 附：环境与耗时记录

- VM：Python 3.12.3、git 2.43.0、node/npm 可用；无 `uv`。端口 8080 已被 Wen 的 mcpgate 占用（pid 2270，M0 未碰它）。
- 耗时：clone 22 秒；`pip install -r requirements.txt` 约 2 分钟（~90 包，venv 216 MB）；首次生成 travel 示例约 40 秒（含从 neuro-san-studio 拉文件 + 自动装 `ddgs`）；静态包生成每次约 1 秒（有缓存后）。
- 生成的本地文件：`nsdeploy/`（clone，含 `.venv/`、`dist/` 三个 bundle、`.nsdeploy/` 生成站）；`~/.cache/nsdeploy/`（152 KB，nsDeploy 的工具缓存）；`/tmp/nsdeploy-*.log`（过程日志）。测试服务器（8091/8092）已全部停止，未留公开服务。
- 未花钱、未注册账号、未写任何 key、未改 mcpgate。
