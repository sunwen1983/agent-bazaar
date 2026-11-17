# Open Agentic Web Challenge — 细化开发计划

> 状态：**M1 进行中**（2026-11-16 晚 Wen 授权最优路线后先行开工前几步：仓库已建、静态包已生成冒烟、mcpgate CORS PR #20 待合并；浏览器内真实跑通留待验证）。
> 计划版本：v1（2026-11-17）｜ 比赛：Open Agentic Web Challenge ｜ 参赛：Wen solo（Devpost 账号 wen-sun1217，已报名）
> 提交截止：**2027-01-11（周一）10:00 AM CT**（= 16:00 UTC）｜ 奖金：现金 $10,000（1st $5,000 / 2nd $3,000 / 3rd $2,000）+ Mozilla Otari credits

---

## 1. 产品一句话定位

**Agent Bazaar（暂名）——浏览器里的 agent 集市。**

一个公开的 agent 网站：来自不同网站的 agent 在你的浏览器里互相发现、验明身份、现场组队干活，全程看得见、每一步都经过你同意。没有中央服务器在幕后编排——浏览器就是编排器，mcpgate 是工具后端。

一句话卖点（写给评委的）：*Agents from different websites meet, verify each other, and team up — live in your browser, with no central orchestrator.*

---

## 2. 比赛事实（已查证，附来源）

查证时间：2026-11-17。以下全部来自官方页面原文，不是推测。

### 2.1 交什么（Devpost 官方页）

来源：https://cognizant-agenticweb.devpost.com/

参赛作品必须是一个**公开的、可工作的 agent 网站**，用标准开放 Web 技术构建，符合 open agentic web specification。提交时提供：

| # | 提交物 | 官方要求 |
|---|--------|----------|
| 1 | 公开网站 URL | 网站必须包含：(a) 公开的、符合规范的 `.json` 文件，描述 agent 网络及其能力；(b) 一个能在浏览器里运行该 agent 网络的 `runtime` 文件；(c) 允许用户或 agent 与该网络交互的界面；(d) 可选的、允许用户配置 LLM 与 key 的界面（BYOK）；(e) agent 网络要能与本站交互，并且可选地与其他 agent 网站交互、或被其他 agent 网站调用 |
| 2 | 文字说明 | 解释项目的功能与特性 |
| 3 | 演示视频 | < 2 分钟（评委不保证看超出的部分）；必须有项目在浏览器里跑的画面；必须有参赛者本人出镜讲解（不能 AI 生成）；必须上传到 YouTube 并公开，提交表单里给链接；不得含第三方商标、受版权保护的音乐/素材（除非有授权） |
| 4 | 截图 | 一张或多张网站截图 |
| 5 | （如果用了 Mozilla Otari） | 提供 Otari 账号邮箱，并用一两句话写使用 Otari 时遇到的问题及改进建议 |
| 6 | （可选但鼓励） | 开发过程中顺手做的、分享给社区的额外工具或服务，可在提交里描述并链接，计入 Potential Impact 评分 |

### 2.2 奖金（Devpost 官方页）

- 1st：$5,000 现金 + $1,500 Mozilla Otari credit + SWAG + 社交推广
- 2nd：$3,000 现金 + $750 Otari credit + SWAG + 社交推广
- 3rd：$2,000 现金 + $500 Otari credit + SWAG + 社交推广
- Employee prize：$250 Otari credit + SWAG + 社交推广
- 注：新闻稿说的 "$13,000 奖金池" = $10,000 现金 + $3,000 Otari credits 之和，两者一致。

### 2.3 评分标准（Devpost 官方页，四项）

1. **Technological Implementation**——做得好不好、多有效地用了 open agentic web specification
2. **Design**——是否是一个完整、连贯的产品体验，而不只是技术验证
3. **Potential Impact**——是否为真实受众解决真实问题、论证是否可信、演示是否支撑；是否包含分享给社区的工具/服务
4. **Quality of the Idea**——是否是创造性的、非显而易见的 agent 网站用法，是否真正理解问题空间

评委：Babak Hodjat（Cognizant 首席 AI 官）、Ramesh Raskar（MIT 教授、NANDA 创始人）、Risto Miikkulainen（UT Austin 计算机教授）。

### 2.4 两条赛道（worldwideagenticweb.org）

- Track 1：把已有网站 agent 化（agentify an existing website）
- Track 2：为 World Wide Agentic Web 造技术组件——agent 发现、身份与信任、互操作、交易、谈判、治理与政策工具、真实应用
- **我们的项目同时踩两条**：一个可用的 agent 集市网站（Track 1 形态）+ 发现/验证这一块做成可复用的机制（Track 2 形态）。（Track 1/2 在官网与 Devpost 上是主题方向，M0 确认提交表单是否要求单选。）

### 2.5 核心规范：MVP Profile 0.2（已查证）

来源：https://github.com/cognizant-ai-labs/agentic-web/blob/main/docs/worldwide_agentic_web_mvp_spec.md
（Minimum Viable Specification for a Worldwide Agentic Web，MVP Profile 0.2，2026-11-06，Cognizant AI Lab）

这是**整个计划的技术底座**，关键点如下（只写与我们相关的）：

**三个信条：**
1. **通用的多 agent 表示**：规范表示是 Neuro SAN HOCON 数据模型；因为 JSON 是合法的 HOCON 序列化，网站可以发布**完全解析好的 JSON 文档**供浏览器加载。网络定义必须保留 Neuro SAN 数据模型，不能依赖本地 includes、环境变量替换或发布者私有路径。
2. **浏览器执行与站点组合**：符合规范的 runtime 必须能在普通网页里用、必须用 JavaScript 在浏览器里执行已发布的图。参考实现是 **nsDeploy browser runtime**。
3. **可插拔推理**：网站可以发布默认 router 和模型，也可以提供 BYOK（访客自带 router、模型、key）。

**agent-manifest.json（必须有的文件）：**
- 网站发布的 agent 网络的公开机器可读索引
- 每个网络：名字、front agent（`front_man`）、已解析网络定义的位置（`hocon` 字段）
- 还可含：站点级推理/router 元数据、能力、策略、部署元数据、peer policy
- 必须公开 GET、允许跨域（CORS）

**网络定义（network definition）：**
- 解析好的、符合 Neuro SAN agent-network schema 的文档，以 JSON 序列化供浏览器加载
- 每个 agent：网内唯一名字、功能描述、用 LLM 时的 instructions、可调用的 agent/工具列表
- 可调用的 agent/工具要暴露 JSON Schema 兼容的输入；结构化输出也建议声明 schema
- 模型名与生成参数只是 hint，调用方的 router 和模型策略优先
- 惯例文件名：`<network>.network.json`（只是惯例，不是互操作要求；以 manifest 声明的位置为准）

**links.json（peer-network 链接表）：**
- 发布者拥有的默认 peer 绑定：把本地网络定义里的符号化外部网络引用（如 `transit`）映射到具体的 peer 网络 URL（如 `https://peer.example/agents/transit`）
- 公开 GET；调用方自己的覆盖存在调用方 origin 的浏览器存储里，不改发布者的文件

**加载 peer 站点的协议（§3.2，核心中的核心）：**
1. peer 网络用 URL 标识，URL 最后一段是网络名；runtime 把前面部分当作站点 base，从该 base 取 `agent-manifest.json`。例：`https://peer.example/agents/transit` 对应 `https://peer.example/agents/agent-manifest.json` 里的 `transit` 网络。HTTPS 必需（本地开发除外）。
2. 调用方在自己的浏览器里实例化 peer 网络；peer 用调用方的 router、模型策略和执行上限。用户 key 只发给 router，**绝不发给 peer 站点**。
3. **同意（Consent）**：交互式 runtime 必须提供同意机制；默认情况下，每次联系第三方 origin 之前，必须说明外部 origin 和联系目的，并取得用户同意。任何放宽都必须来自用户的明确选择。
4. 解析时名字必须**精确匹配**，没有精确匹配就必须失败，不能退而选另一个网络。
5. 每一跳 peer 或 MCP 都重复 origin 策略与同意检查；runtime 要拦住循环和超预算请求。
6. 符合 MVP 不要求下载或执行发布者提供的代码；跨域产物必须带合适的 CORS 头。

**MCP 是工具的基线接口：**
- 运行时必须支持 MCP 作为工具与数据的基线接口
- 发现或调用第三方 MCP 服务之前要先问用户批准
- 本地站点动作（site actions）：站点向 runtime 注册的能力，含名字、用途、JSON Schema、效果分类、handler；runtime 必须校验参数，写操作或不可逆操作之前必须向用户展示目标与参数

**信任边界（§5，演示里要看得见）：**
- peer 网络定义必须被当作**不可信输入**（它的 instructions 可能试图操纵调用方的模型）
- 私有上下文（sly_data）默认**不**跨网络传递，只有明确命名的字段才可以
- 执行失败要 fail closed：不能安全校验/解析/授权、或超出上限的调用，runtime 必须不开始、不继续、不谎报成功
- 用户拥有的写操作：确认归 runtime 不归模型；执行的动作、origin、参数、请求 ID、回执要绑定在已审核的提案上
- 可观察执行：实现**应该**暴露所选 router、联系过的 origin、agent 与工具调用、同意决定、错误、动作回执，且不暴露 secret

**BYOK（§4）：**
- 网站可以发布默认 router/模型，也可以提供 BYOK
- BYOK 时用户能选/填 router base URL、模型、key；key 只留在调用方浏览器里，只附加到发往所选 router 的模型请求上；**绝不能出现在 URL、manifest、网络定义、prompt、动作参数、trace、分析或 peer 请求里**
- reference runtime 把 BYOK profile（含 key）存在 origin 隔离的 localStorage；用户必须能删除
- 如果访客选的 router 不允许浏览器直接跨域访问，托管交互 runtime 的站点可以提供**同源转发 endpoint**（只转发访客自己的 Authorization 头，绝不能替换成发布者自己的 key）

**互操作自查清单（§6，M3 验收就按这个来）：**
- 从普通 HTTPS 主机发布一个双 agent HOCON 网络 + manifest
- 在浏览器里用网站发布的默认 router/模型加载并运行该网络；如果提供 BYOK，也演示用户自选路径
- 调用一个内部 agent 和至少一个经 MCP 或本地站点动作落地的能力
- 经用户批准后，从**第二个 origin**加载一个精确命名的网络，在调用方浏览器里本地运行
- 用户 key 激活时，证明第二个站点收不到 key、也改不了调用方的 router
- 证明私有上下文在 peer 处不存在，除非某个字段被明确允许
- 拦住无效网络、组合循环、超预算请求、被拒 origin，并显示可见错误

### 2.6 参考实现 nsDeploy（已查证到仓库与结构）

来源：https://github.com/cognizant-ai-lab/nsdeploy

- 浏览器 runtime 在 `nsdeploy/runtime/`：`engine.js`、`kernel.js`、`runner.js`、`mcp.js`、`byok.js`、`llm.js`、`hocon.js`、`sandbox-worker.js`、`trace.js`、`ui.js`、`action-bridge.js`、`actions.js`、`discovered-actions.js`、`form-actions.js`、`execution-stream.js`、`cruse.js`、`cruse-render.js`、`remote.js`、`web-crypto.js` 等
- Python 端在 `nsdeploy/`：`generate.py`、`linker.py`、`resolve.py`、`site_manifest.py`、`gateway.py`、`mcp_host.py`、`steward.py`、`studio.py`、`cailville.py` 等；另有 `commerce` 示例目录
- 参考实现目前支持 OpenAI 兼容 endpoint 的 LLM 客户端
- **待确认（M0 spike 解决）**：nsDeploy 的确切嵌入方式（script 引入？打包 bundle？有无 starter/示例站点的 README 步骤）、`commerce` / `cailville` 示例怎么跑起来

### 2.7 NANDA 发现与信任机制（已查证，见 2.7.1；未完全验证的见 2.7.2）

#### 2.7.1 已查证

- **NANDA Index**：去中心化的 agent 注册表，被称为"agent 的 DNS"。agent 在 index 注册后可被网上任何系统发现。底层是 federated 注册表组成的"quilt"架构，不同生态的注册表可以互通。
- **AgentFacts**：W3C Verifiable Credentials（已签名的 JSON-LD 文档），描述 agent 的能力、身份、endpoint、策略；用于证明"这个 agent 是它声称的那个、有它声称的能力"，防冒充。
- **身份**：DID（去中心化标识符）+ 公私钥对，每个 agent 有可验证身份；AgentAddr 约 120 字节。
- **Zero Trust Agentic Access (ZTAA)**：零信任访问层。
- **Web Agentica**：NANDA 的研究原型——已签名的、可移植的 agent 包，在用户控制的 runtime 里跑（OpenClaw 或浏览器）。
- **NandaIndex 注册路径（真实存在）**：https://my.nandaindex.org/ —— 流程是：先用 Inkbox 给 agent 一个专用邮箱地址，再用 AIR 标识符（Agent Internet Resource，如 `urn:air:host39.org:personal:...`）在 NandaIndex 注册并完成邮箱验证；注册后其他 agent 可以解析这个标识符。官方还发布了 SKILL.md 指导 AI 助手完成这个设置。
- **NANDA 不取代 MCP / A2A / HTTPS**：它是包在这些协议外面的发现与信任层。

#### 2.7.2 待确认（M0 spike 解决）

- 网站型 agent 网络（不是个人助理型 agent）能不能走同一套 NandaIndex 注册流程？注册后浏览器里的 runtime 能不能程序化地解析 AIR 标识符、取回 AgentFacts？
- Web Agentica 的签名包格式与我们的 HOCON/JSON 网络定义如何对应？是直接复用，还是只借鉴"签名 + 可验证"这个思路？
- **退路已定**：如果 NandaIndex 程序化接入走不通，信任层退到"本地实现"——manifest 带 Ed25519 签名 + 公钥、集市 UI 做四步可见验证（见 §5.3），并在文字说明里写明与 NANDA 机制的对应关系。比赛不强制用 NANDA（官方明确说这些贡献的使用不是参赛必需），所以退路不影响合规。

### 2.8 mcpgate 现状（已查证，本地仓库）

来源：`~/workspace/mcpgate/`（Wen 的 Go 开源项目，MIT 协议，已发布 v0.1.0，主线 8131c80）

- 定位：MCP 的生产层网关，坐在 MCP 服务器前面，加 auth、按 key 限流、响应缓存、熔断、结构化日志、Prometheus 指标，不动服务器本身
- HTTP 入口：`POST /mcp/{backend}`，JSON-RPC 2.0；支持 Streamable HTTP 与 SSE 流式（`Accept: text/event-stream` 时中继后端 SSE 事件）
- Auth：`X-API-Key` 或 `Authorization: Bearer`，另支持 OIDC；后端可以是 stdio 子进程或 HTTP MCP 服务器；支持多 endpoint 负载均衡
- 可观测性：`/healthz`、`/readyz`、`/metrics`；缓存命中/未命中、限流 429 等均有指标
- 已验证过的行为（2026-11-16 实录 demo）：401、缓存 MISS→HIT、`add(2,3)=5`、200×10 后 429×2、healthz/metrics
- **待确认（M0 spike 解决）**：mcpgate 是否默认发 CORS 头、浏览器直连跨域 MCP 是否畅通；若不通，方案是同源代理还是给 mcpgate 加 CORS 支持（后者是他自己的项目，改动小，还能算作"分享给社区的额外工具"）

---

## 3. 魔法瞬间：demo 前 10 秒的具体设计

FixLens 的教训：demo 前 10 秒必须有一个"不可能"时刻。我们的不可能时刻是：**三个互不认识的网站的 agent，在你的浏览器里现场相识、验明正身、组队干活——没有任何中央服务器在指挥。**

### 3.1 前 10 秒分镜（精确到秒）

| 时间 | 画面 | 观众看到什么 | 旁白（Wen 出镜口述，约 15 字以内一句） |
|------|------|-------------|----------------------------------------|
| 0:00–0:02 | 集市首页，一个输入框 | 用户输入一句话："Plan me a 3-day weekend in Austin for under $800."（输入框已预填好，点一下就发，避免打字浪费时间） | "One sentence. Three strangers." |
| 0:02–0:04 | 同意弹窗 | 一次列出三个外部 origin 与联系目的，点一下 Allow（规范要求联系第三方 origin 前须取得同意；单次弹窗列明三家是一次性拿全三家同意） | "Your browser asks first." |
| 0:04–0:07 | 三张 agent 卡片依次亮起 | 三张卡片从三个不同的域名"飞"进来（卡片上印着各自的域名：flights / stays / local-guide），每张卡片依次出现绿色对勾：✓ Manifest verified、✓ Signature verified、✓ NANDA registered | "They never met. They just verified each other." |
| 0:07–0:10 | 三张卡片拼成一个"小队" | 卡片连线，汇成一个 team 面板，出现一行字："Team formed in your browser — no central orchestrator."；trace 开始滚动，三个 agent 并行开工 | "No orchestrator. The browser is the orchestrator." |

> 注意：① 策略检查（§5.3 第④步）在验证面板里可见，不占魔法瞬间的三个勾；② 若走 B 档（本地信任），第三个勾显示为"本地信任"而非 NANDA 注册，文案相应调整；③ 第一条真实结果落在 0:10–0:45 段，不承诺 10 秒内出结果（LLM 首 token 延迟）。

**为什么这是"不可能"时刻**：今天的 agent 都被锁在各自的平台里；三个不同网站的 agent 在用户浏览器里互相发现、验证、组队，是观众没见过的画面。验证对勾和"no central orchestrator"这行字是点睛之笔——把信任机制做成看得见的东西。

### 3.2 完整 2 分钟视频结构（< 120 秒，留 10 秒余量）

| 段 | 时长 | 内容 |
|----|------|------|
| 魔法瞬间 | 0:00–0:10 | §3.1 |
| 组队干活 | 0:10–0:45 | 三个 agent 并行工作，trace 实时滚动；mcpgate 工具调用可见（查航班、查酒店、查活动各一次，带缓存 MISS→HIT 的小字）；中途一个 agent 提议、另一个接力（delegation 可见） |
| 成果 | 0:45–1:10 | 一份完整行程单生成：机票 + 酒店 + 每天活动 + 总价 $742（在预算内），每项标着是哪个 agent 给的、可点开溯源 |
| 信任收尾 | 1:10–1:30 | 镜头切到验证面板：点开一个 agent 的 manifest、签名、NANDA 注册信息；再演示一次"拒绝"——输入一个未注册/签名不对的 peer URL，runtime 当场拦下并显示原因（fail closed 看得见） |
| 收尾 | 1:30–1:50 | Wen 出镜一句话总结 + 屏幕显示开源仓库地址 + "Try it live" 的网站链接；结束画面留 5 秒静帧给链接 |

**录制要点**：Wen 本人出镜讲解（官方要求，不能 AI 生成）；浏览器画面必须是真实跑的项目；全程不用第三方商标与版权音乐；提前写好逐字稿并排练 2–3 遍再正式录；录两版备选，取更好的一版。

---

## 4. 核心功能清单

### 4.1 Must-have（没有这些就不叫完成）

| # | 功能 | 对应规范/评分 | 验收标准 |
|---|------|-------------|----------|
| M1 | 集市首页：agent 卡片墙 + 一句话输入框 + 实时 trace 面板 | Design / Quality of the Idea | 打开即懂；输入一句话就能启动一次组队 |
| M2 | `agent-manifest.json`：发布至少 2 个网络（concierge + 至少 1 个 peer 网络），含 front agent、hocon 位置、router 元数据、peer policy | Technological Implementation（合规硬要求） | 文件公开可 GET、带 CORS；名字精确匹配 |
| M3 | 网络定义 JSON：concierge + 3 个 peer 网络（flights / stays / local-guide），Neuro SAN 数据模型，含 agent instructions 与 JSON Schema 输入 | Technological Implementation | nsDeploy 能解析、校验通过；无效网络能被拦下 |
| M4 | 跨 origin 加载：peer 网络放在**不同 origin**上，经用户同意后在调用方浏览器里本地运行 | §3.2 / §6 自查 | 演示里真实发生，不是同 origin 模拟 |
| M5 | 验证面板（信任看得见）：每个 peer 加载前依次显示 manifest 校验、签名校验、NANDA 注册查询结果；失败时拦下并显示原因 | Quality of the Idea / §5 | 拒绝案例可演示（坏签名 / 未注册 / 名字不匹配） |
| M6 | mcpgate 工具后端：3 个 MCP 工具服务（航班、酒店、活动）经 mcpgate 暴露，带 API key auth、限流、缓存；浏览器 runtime 经 MCP 调用 | Technological Implementation / §3.1 | 工具调用真实发生；缓存 MISS→HIT 可见；无 key 时 401 |
| M7 | BYOK 配置界面：用户能填 router base URL、模型、key；key 只存浏览器 localStorage，可删除；默认提供 Wen 出资的便宜模型（见 D4） | §4 / 提交要求 (d) | 无 key 也能跑默认路径；有 key 时 peer 收不到 key |
| M8 | 同意机制：每次联系第三方 origin 前弹窗说明 origin 与目的，等用户点同意；记住的同意只限当前 session 与该 origin | §3.2 / §5 | 演示里真实弹窗一次 |
| M9 | 执行上限与 fail closed：组合深度、调用次数、耗时上限；超限/无效/被拒时可见错误，不谎报成功 | §3.1 / §6 | 自查清单逐项通过 |
| M10 | 公开网站 + 开源仓库：网站公开可访问；代码 MIT 开源；README 含本地运行步骤 | 提交硬要求 | 评委点开链接就能跑 |
| M11 | 演示视频（<2 分钟，Wen 出镜，YouTube 公开）+ 文字说明 + 截图 | 提交硬要求 | 按 §3.2 结构录制；提前 2 天完成 |

### 4.2 Nice-to-have（有时间再做，按优先级排序）

| # | 功能 | 价值 | 做的前提 |
|---|------|------|----------|
| N1 | 第四、第五个 peer agent（如 food / transport），让"集市"更像集市 | Design | M1–M10 全绿之后 |
| N2 | "Add your agent" 入口：访客贴一个 peer URL，集市当场加载并验证它（真正的开放集市） | Quality of the Idea / Potential Impact | M5 稳定之后；需防恶意 URL（只读 manifest、严格校验） |
| N3 | mcpgate 加 CORS 支持（或同源代理），作为分享给社区的改进 | Potential Impact（额外工具加分） | M0 确认缺口存在 |
| N4 | manifest 在线校验小工具（贴 URL → 逐项打勾） | Potential Impact；评委可能喜欢 | M2 之后，半天工作量 |
| N5 | CRUSE 式生成界面：agent 提议的表单/选择/报表由 runtime 渲染（规范 §7 扩展） | Design 亮点 | M1–M10 全绿且有整块空闲时间 |
| N6 | 多语言（中英）界面 | 无 | 最后有空再做；不列入主线 |
| N7 | 用 Mozilla Otari 做默认 router | 可能多拿 Otari 相关奖项/credits；但要多写邮箱+问题反馈 | D8 决定；如用，放在 M7 之后 |

### 4.3 明确不做

- 不做交易/支付（谈判、结算在规范 §7 里是扩展，不在 MVP 内；demo 里只做到"出方案"，不真下单）
- 不做账号系统、数据库、后端编排器（违背"浏览器就是编排器"的定位）
- 不做移动 App（规范 §7 也列为扩展）
- 不重造 NANDA 注册表（只用它的机制：DID + AgentFacts 概念 + NandaIndex 注册；见 §5.3）
- 不做实时协作/多人（单用户 demo 足够）

---

## 5. 技术方案

### 5.1 架构（文字图）

```
访客浏览器（唯一执行地）
┌─────────────────────────────────────────────────────────────────┐
│ 集市前端 (静态站, origin A)                                       │
│  ┌──────────┐  ┌──────────────┐  ┌────────────┐  ┌───────────┐  │
│  │ 卡片墙    │  │ 验证面板      │  │ Trace 面板  │  │ BYOK 设置  │  │
│  │ + 输入框  │  │ (manifest/   │  │ (实时执行  │  │ (local-   │  │
│  └────┬─────┘  │  签名/NANDA)  │  │  流)       │  │ Storage)  │  │
│       │        └──────┬───────┘  └─────┬──────┘  └─────┬─────┘  │
│       ▼               ▼                ▼               ▼        │
│  nsDeploy browser runtime (engine.js/kernel.js/runner.js/...)     │
│   ├─ 解析/校验 HOCON-JSON 网络定义 + manifest                     │
│   ├─ agent 循环：流式模型调用、工具调用、agent 间委派              │
│   ├─ peer 加载：fetch 对方 manifest → 精确匹配 → 取网络定义 → 本地运行│
│   ├─ 同意弹窗 / 执行上限 / fail closed / 私有上下文过滤            │
│   └─ MCP 客户端 (mcp.js)：调工具前问用户批准                       │
└───────┬──────────────────────────┬──────────────────┬───────────┘
        │ HTTPS + CORS             │ HTTPS + CORS     │ HTTPS
        ▼                          ▼                  ▼
 peer 站 B (静态站)          peer 站 C (静态站)    mcpgate (Go 服务)
 flights agent 网络          stays + guide 网络    POST /mcp/{backend}
 agent-manifest.json         agent-manifest.json   ├─ API key auth
 flights.network.json        stays.network.json    ├─ 限流 / 缓存 / 熔断
 links.json                  guide.network.json    ├─ /healthz /metrics
                             links.json            └─ 后端: 3 个 MCP 工具服务
                                                     (flights / stays / activities)
        ▲                                                       ▲
        │ 发现: links.json + manifest                            │ 信任: NandaIndex
        │      (+ 可选 .well-known/nsdeploy/)                   │   注册查询 + 签名验证
        └──────────────────────────────────────────────────────┘
              LLM 推理: 默认 router (DeepSeek 兼容 endpoint, Wen 出资, 见 D4)
                        或访客 BYOK (OpenAI / DeepSeek / 其他兼容 endpoint)
```

**关键设计选择（都是规范直接支持的）：**
- 浏览器是唯一执行地；没有后端编排器；peer 网络在调用方浏览器里跑，用调用方的 router 与上限。
- 三个 origin 是真实的三个静态托管（不是同站不同路径的模拟），这样跨 origin 加载、同意、CORS、私钥不外泄都是真实发生的。
- mcpgate 是唯一的"服务器味"组件，但它只做工具网关，不做 agent 编排——定位清晰，评委问起来也好讲。

### 5.2 前端

- 技术：纯静态站，HTML + CSS + 原生 JavaScript modules（不引入重型框架，减少出问题面；nsDeploy runtime 本身就是 JS 模块）。构建：能用 `vite` 打包就打包，不能就直接静态文件——以 M0 验证结果为准。
- 页面结构（单页）：
  1. 顶部：一句话定位 + "How it works" 三步小字（发现 → 验证 → 组队）
  2. 输入框 + 预设示例（一句话任务，点一下就发）
  3. Agent 卡片墙：每个 agent 一张卡（名字、来源域名、能力、验证状态徽章、NANDA 注册状态）
  4. Team 面板：组队成功后出现，显示小队成员与连线
  5. Trace 面板：实时执行流（agent 调用、工具调用、同意记录、错误），可展开看参数（不含 secret）
  6. 成果区：行程单（每项标来源 agent，可溯源）
  7. 验证面板：点开任意 agent 看 manifest 原文、签名、公钥、NANDA 查询结果
  8. BYOK 设置：router、模型、key 输入 + 删除 profile
  9. 页脚：开源仓库链接、规范链接、"Built for the Open Agentic Web Challenge"
- 视觉要求：卡片"飞入"与对勾动画是魔法瞬间的载体，要流畅；trace 要像机场出发屏一样实时滚动；整体风格干净，不堆装饰。

### 5.3 NANDA 怎么用（发现与信任，不重造）

分两层，M0 spike 后按验证结果二选一（见 D6）：

**A 档（目标）：真实 NandaIndex 注册 + AgentFacts 验证**
1. 每个 peer agent 在 https://my.nandaindex.org/ 注册，拿到 AIR 标识符（`urn:air:...`），完成邮箱验证。
2. manifest 里写上 agent 的 AIR 标识符与 DID/公钥。
3. 集市 UI 加载 peer 时做四步可见验证：
   - ① Manifest 可达：HTTPS + CORS 拉到 `agent-manifest.json`，名字精确匹配；
   - ② 签名验证：用 manifest 里声明的公钥（Ed25519，WebCrypto）验网络定义的签名；
   - ③ NANDA 注册查询：解析 AIR 标识符，查 NandaIndex 记录存在且能力描述与 manifest 一致；
   - ④ 策略检查：peer policy 允许被本站联系。
   四步全过才显示绿勾、才允许组队；任何一步失败都拦下并显示原因。
4. 文字说明里写明这套验证与 NANDA 机制的对应关系（Index = 发现、AgentFacts = 能力凭证、DID = 身份）。

**B 档（退路，若 A 档程序化接入走不通）：本地签名验证 + links.json**
1. ①②④ 同上（签名验证不依赖外部注册表，完全本地可做）。
2. ③ 换成"注册表状态：未接入 NandaIndex（本地信任）"，并在文字说明里诚实写明原因与后续接入路径。
3. 仍保留 AIR 标识符字段（占位），为赛后接入留接口。

**不做的事**：不自建注册表、不自发证书体系、不把 NandaIndex 当数据库用。

### 5.4 mcpgate 怎么接

- 部署：一台常驻小服务（见 D3），跑 mcpgate + 3 个 MCP 工具后端（每个后端一个 MCP server，工具分别是航班查询、酒店查询、活动查询）。
- 浏览器 → mcpgate：nsDeploy 的 `mcp.js` 发 JSON-RPC 2.0 到 `POST https://<mcpgate>/mcp/{backend}`，带 `X-API-Key: <demo key>`（demo key 公开在前端配置里，但限流极严 + 只读工具，见下）。
- 配置要点：
  - demo key：单独一个 key，限流如 30 请求/分钟、只读工具；比赛结束即轮换。
  - 缓存：工具结果 TTL 缓存开着，demo 里"第二次查同一家酒店直接 HIT"是看得见的卖点。
  - CORS：M0 确认 mcpgate 是否发 CORS 头；若没有，二选一：(a) 给 mcpgate 加 CORS 支持（小改动，顺手贡献回社区）；(b) 集市站同源代理转发（只转发、不代持 key）。
  - 健康检查：`/healthz` 挂在页面页脚小字或 trace 里，证明网关活着。
- 工具数据：默认 mock 确定性数据（见 D5），三个后端返回固定但真实感强的数据集（如 6 个航班、8 家酒店、10 个活动，价格总和能凑进预算）；工具本身是标准 MCP server，可换真实 API 而不用改网关。
- mcpgate 自己的改动（如果需要）单独提 PR 回 mcpgate 主线，不与比赛仓库混在一起。

### 5.5 网络定义（HOCON / JSON）怎么写

- 优先用 nsDeploy 自带流程（写 HOCON，nsdeploy 自动解析；不需要单独装 Studio，它会按需从 neuro-san-studio 仓库拉依赖）；Neuro SAN Studio（`ns init` / `ns run`）只在需要可视化设计时再考虑；以 M0 验证为准。
- 每个网络一个 JSON 文件（`<network>.network.json`），内容：
  - `metadata`：名字、描述
  - `tools`：agent 列表，每个 agent 有 `name`、`function.description`、`instructions`、可调用的 `tools`（内部 agent 名、peer 网络 URL、或 MCP descriptor `{url, tools:[...]}`）
  - 调用 schema：JSON Schema 输入；结构化输出声明 schema
  - 模型设置只写 hint，不写死
- 4 个网络分工：
  - `concierge`（主站）：front agent，理解一句话任务、拆解子任务、把子任务委派给 peer、汇总成行程单
  - `flights`（peer B）：查航班、比价、返回结构化航班选项
  - `stays`（peer C）：查酒店、按预算筛、返回结构化酒店选项
  - `local-guide`（peer C）：查活动/餐厅/天气、排每天日程
- **红线**：网络定义是公开文件，instructions 里不写任何 secret、key、内部 URL；所有敏感的东西只走 MCP 或 BYOK。

### 5.6 文件清单（最终网站上必须有的）

```
origin A (集市主站):
  /index.html                  集市单页
  /runtime/                    nsDeploy browser runtime (vendored 或构建产物)
  /agents/agent-manifest.json  主站 manifest (concierge 网络 + router 元数据 + peer policy)
  /agents/concierge.network.json
  /agents/links.json           符号名 -> peer URL (flights/stays/guide)
  /.well-known/nsdeploy/       可选发现元数据
  /byok.html 或站内设置面板     BYOK 配置界面
origin B (peer 站, 如 flights):
  /agents/agent-manifest.json
  /agents/flights.network.json
origin C (peer 站, 如 stays+guide):
  /agents/agent-manifest.json
  /agents/stays.network.json
  /agents/guide.network.json
mcpgate 服务:
  /mcp/flights  /mcp/stays  /mcp/activities   (三个 MCP 后端)
  /healthz  /metrics
```

### 5.7 部署与域名

- 静态站：三 origin 用同一家静态托管的不同项目/子域（见 D3：GitHub Pages / Cloudflare Pages / Vercel 三选一，要求：免费、支持自定义 CORS 头、HTTPS 默认开启）。
- mcpgate：小服务托管（见 D3：fly.io / Railway / 其他，要求：常驻、HTTPS、有固定域名、月费接近 $0）。
- 域名：先用托管商默认子域（如 `*.pages.dev`），如有空闲时间再考虑自定义域（不列入主线）。但若 D6 选 A 档（真实 NandaIndex 注册），域名锚定注册需要 DNS TXT 证明，默认子域做不了，自定义域就变成 must-have——M0 定 D6 时一并定。

---

## 6. 里程碑排期（2026-11-17 → 2027-01-11）

**总时长 8 周。硬约束**：Apple SDS 面试 **12/10（周四）9:00–9:45 AM CT 是头等大事**，12/8–10 面试周不碰黑客松；12/5（周六）米粒生日会；Phase 0 每日刷题照常；圣诞/新年假期节奏放缓。节奏按"工作日晚上 + 周末整块"估算，solo 开发，留足 buffer。

| 里程碑 | 时间 | 交付物 | 退出标准 |
|--------|------|--------|----------|
| **M0 探路** | 11/17（周二）– 11/23（周一） | nsDeploy + agentic-web clone 到本地；跑通官方示例（commerce / cailville）；确认嵌入方式；确认 NANDA 注册流程；确认 mcpgate CORS；确认托管方案；决策 D1–D6 定稿 | 能在浏览器里跑起一个官方示例网络；四个待确认项每个都有"是/否 + 退路"结论；决策清单签字 |
| **M1 最小闭环** | 11/24（周二）– 12/7（周一） | 集市壳子（首页+输入框+卡片墙）；manifest + concierge 网络；1 个 peer 网络在第二 origin 上；mcpgate 1 个工具后端通；第一次端到端"发现→验证→组队"在浏览器里跑通 | 输入一句话，两个 origin 的 agent 真能组队并返回一条结果；§6 自查清单过一半；若跨 origin 延期，退出标准降为单 origin 端到端 + mock peer，跨 origin 移到 M2，提交日不变 |
| **面试周停工** | 12/8（周二）– 12/10（周四） | ——（不碰黑客松，全力准备 SDS 面试） | —— |
| **M2 完整功能** | 12/11（周五）– 12/21（周一） | 3 个 peer 网络全齐；验证面板四步全实现；BYOK 界面；同意弹窗；执行上限；trace 面板；行程单成果区 | 完整魔法瞬间（§3.1）从头到尾可演示；坏签名/未注册拒绝案例可演示 |
| **M3 打磨与自查** | 12/22（周二）– 12/28（周一） | 按 §6 互操作清单逐项自查并修；边缘 case（超预算、循环、超时、断网）处理；截图素材；性能与首屏优化 | §6 清单全绿；连续 3 次演示零失误 |
| **M4 视频与文字** | 12/29（周二）– 1/4（周一） | 视频逐字稿 → 排练 → 正式录制（2 版备选）→ 剪辑 → YouTube 公开；文字说明定稿；截图定稿；这周含元旦，提前和家里确认录制时间，魔法瞬间备用镜头争取 M3 录好 | 视频 < 2 分钟、本人出镜、浏览器实跑画面；文字说明讲清四项评分点 |
| **M5 提交** | 1/5（周二）– 1/8（周五） | Devpost 表单填写、链接全检查、提交 | 表单提交成功；所有链接点开可用 |
| **Buffer** | 1/9（周六）– 1/10（周日） | 修提交后发现的问题；最后检查 | —— |
| **截止** | 1/11（周一）10:00 AM CT | —— | 目标 1/8 提交，最晚 1/10，不拖到截止当天 |

**每周节奏建议**：工作日晚上 1–1.5 小时（写代码/调 bug），周六或周日一个 3–4 小时整块（联调/录视频）；M0 与 M4 的周末块最重要。面试周（12/8–10）与 12/5 生日会当天不排任务。

---

## 7. 风险与待确认事项

| # | 风险 | 影响 | 缓解 / 退路 | 何时消除 |
|---|------|------|------------|----------|
| R1 | nsDeploy 嵌入方式未完全验证（无 README 步骤在手） | M1 开工受阻 | M0 spike：跑通官方示例；若跑不通，退路是自己写最小 runtime（只实现 manifest 解析 + peer 加载 + MCP 调用 + 同意弹窗），工作量约 1 周，M1 顺延 | 11/23 |
| R2 | NandaIndex 程序化接入未验证 | 验证面板第 ③ 步可能做不了真的 | 退路 B 档（本地签名验证）；信任故事仍完整；比赛不强制用 NANDA | 11/23 |
| R3 | mcpgate CORS 未验证 | 浏览器直连 MCP 可能被拦 | 同源代理 或 给 mcpgate 加 CORS（二选一，M0 定） | 11/23 |
| R4 | 默认推理 key 的成本与滥用 | 公开 demo key 被刷爆 | 便宜模型（DeepSeek）+ 硬性额度上限 + 演示 key 严格限流 + 比赛结束轮换；BYOK 分流 | M2 |
| R5 | 时间风险：面试 + 生日会 + 圣诞新年 + solo | M2/M3 挤压 | 面试周停工已排进计划；M5 之后留整整一个周末 buffer；nice-to-have 随时可砍 | 全程 |
| R6 | 三 origin 的 CORS 配置 | 跨域加载失败 demo 当场翻车 | M1 就用真实三 origin 联调，不用同站模拟；M3 自查清单覆盖 | M1 |
| R7 | 视频重录风险（口误、demo 翻车、超时） | M4 拖期 | 逐字稿 + 排练 2–3 遍 + 录两版备选；魔法瞬间单独录一条备用镜头 | M4 |
| R8 | 对规范理解偏差（如 manifest 字段名、peer URL 规则） | 返工 | M1 结束时对照 §6 清单自查一次；M3 再查一次；有疑问直接读 spec 原文，不猜 | M1/M3 |
| R9 | 评委点开链接时网站/网关挂了 | 直接丢分 | 提交前做"冷启动检查"（无痕浏览器 + 手机网络各跑一次完整流程）；mcpgate 加健康检查与自动重启 | M5 |
| R10 | 默认推理 key 保管方式未定：nsDeploy 的 served site 明确不持有发布者 key（BYOK guarantee），D4-A"默认 DeepSeek、Wen 出资"用现成 nsDeploy 做不出来 | 默认路径不存在，评委打开就得自己贴 key | M0 必须定方案：(a) 纯 BYOK（体验差）；(b) 自建 token broker/代理注入 key（额外工作量，且规范 §4 不鼓励发布者代持 key 的转发）；(c) 前端嵌演示 key + DeepSeek 账户余额封顶（如 $5）+ 赛后轮换（最简单，key 公开但损失有上限）；视频用 BYOK + Wen 自己的 key（镜头外输入），保证画面是真实模型行为 | 11/23 |

**待确认事项清单（计划外需要查的东西，M0 解决）：**
- [ ] nsDeploy 嵌入方式与示例运行步骤
- [ ] NandaIndex 网站型 agent 注册与程序化解析是否可行
- [ ] mcpgate 默认 CORS 行为
- [ ] 静态托管三家的 CORS 头配置方法（选定一家后确认）
- [ ] Otari 用还是不用（D8）
- [ ] 默认推理模型与月预算上限（D4）
- [ ] 默认推理模型的 key 保管方案（见 R10）
- [ ] Devpost 规则是否要求作品在比赛期内新建（mcpgate 是赛前已有项目，作后端组件是否合规、是否需披露）
- [ ] 提交表单是否要选赛道（Track 1/2 是主题方向还是必填项）
- [ ] NANDA A 档与自定义域名耦合（域名锚定注册需 DNS TXT 证明；个人 agent 路径 Inkbox + host39.org 是否适用网站型网络）
- [ ] nsDeploy pin 一个 commit（仓库 11 月 5 日才建、变化快）

---

## 8. 需要 Wen 拍板的决策点

**规则：以下每项由 Wen 逐项点头，不替他决定。点头后写入本计划"已定"一节。**

**已定（2026-11-16 晚，Wen 授权 Mila 按最优路线决定）**：D1-A Austin 三日游；D2-A nsDeploy（pin 4e4cb36）；D3 静态站 Cloudflare Pages，mcpgate 放 fly.io；D4-(b) 前端嵌演示 key + DeepSeek 余额封顶约 $5、赛后轮换，视频用 BYOK + Wen 自己的 key（镜头外输入）；D5-A mock 确定性数据；D6-B 档本地签名验证为主线，A 档不做；D7 `agent-bazaar`；D8 不用 Otari；D9 视频 M4 期间录（手机 + 录屏，具体工具到时定）；D10 工作日晚上 1–1.5h + 周末 3–4h 整块。

| # | 决策 | 选项 | 推荐 | 影响 |
|---|------|------|------|------|
| D1 | Demo 场景 | A. Austin 三日周末游（机票+酒店+本地向导） / B. 孩子生日派对策划（场地+餐饮+娱乐） / C. 搬到新城市安家（找房+搬家+水电网） | A | 决定 3 个 peer agent 是谁、mock 数据长什么样、视频台词 |
| D2 | Runtime 路线 | A. 直接用 nsDeploy 参考 runtime / B. 自己写轻量 runtime（只实现 MVP 最小集） | A | A 快但要吃透别人代码；B 可控但多约 1 周工作量 |
| D3 | 托管 | 静态站：GitHub Pages / Cloudflare Pages / Vercel；mcpgate：fly.io / Railway / 其他 | 静态站 Cloudflare Pages（CORS 头好配），mcpgate fly.io（常驻免费档） | 决定部署步骤与域名 |
| D4 | 推理与预算 | A. 默认 DeepSeek 便宜模型（Wen 出资，设月上限 $___）+ BYOK / B. 纯 BYOK（评委须自带 key） / C. 默认 OpenAI 便宜模型 + BYOK | A | B 会让评委体验变差，不推荐；A 的月上限金额由 Wen 定 |
| D5 | 工具数据 | A. Mock 确定性数据（演示稳定、零成本） / B. 接真实 API（更真、但要 key、钱、处理失败） | A | 决定 M1 工具后端怎么写 |
| D6 | NANDA 做到哪档 | A. 真实 NandaIndex 注册（M0 验证可行才选） / B. 本地签名验证 + 占位 AIR 字段 | M0 验证后定 | 决定验证面板第 ③ 步的实现 |
| D7 | 项目名与 repo 名 | 候选：`agent-bazaar` / `agentic-bazaar` / `agent-market` / 其他（Wen 定） | `agent-bazaar` | 决定仓库名、域名、视频结尾画面 |
| D8 | 用不用 Mozilla Otari | 用（默认 router 走 Otari；提交里要写邮箱 + 使用问题反馈）/ 不用 | 不用（除非 M0 发现 Otari 接入极顺且有额外奖项） | 影响 M7 与提交材料 |
| D9 | 视频 | 录制方式：手机 + 录屏软件（Wen 定具体工具）；录制时间：M4 期间（12/29–1/4）周末整块 | —— | 决定 M4 安排 |
| D10 | 每周投入 | 工作日晚上 1–1.5h + 周末 3–4h 整块（推荐）/ 其他节奏（Wen 定） | 推荐节奏 | 决定里程碑是否现实；如时间不够，优先砍 nice-to-have，不砍 must-have |

---

## 9. 查证记录（本计划的事实来源）

| 内容 | 来源 | 查证时间 |
|------|------|----------|
| 提交物、截止时间、奖金、评分标准、评委 | Devpost 官方页 https://cognizant-agenticweb.devpost.com/ | 2026-11-17 |
| MVP Profile 0.2 全文（manifest、peer 协议、BYOK、MCP、自查清单） | https://github.com/cognizant-ai-labs/agentic-web/blob/main/docs/worldwide_agentic_web_mvp_spec.md | 2026-11-17 |
| nsDeploy 仓库结构（runtime 文件清单、Python 工具） | https://github.com/cognizant-ai-lab/nsdeploy | 2026-11-17 |
| NANDA Index / AgentFacts / DID / Web Agentica 机制 | NANDA 论文 arXiv:2507.14263、knowyourmodel.ai NANDA 文档、NANDA 官方 LinkedIn（my.nandaindex.org 注册流程） | 2026-11-17 |
| Neuro SAN Studio 用法（ns init / ns run / 104 示例网络） | https://github.com/cognizant-ai-lab/neuro-san-studio | 2026-11-17 |
| mcpgate 能力（auth/限流/缓存/SSE/端点） | 本地仓库 ~/workspace/mcpgate/README.md + 2026-11-16 实录 demo | 2026-11-17 |
| 两条赛道、$13,000 口径、合作方名单 | worldwideagenticweb.org、PR Newswire 通稿（2026-11-12） | 2026-11-17 |

**本计划未编造任何规则**。凡是写了"待确认"的，都是查不到官方原文的，等 M0 spike 验证后更新本计划。

---

## 10. 审查记录（2026-11-17，Mila 复核，be critical）

**已核实无误**：nsDeploy 仓库真实存在（cognizant-ai-lab/nsdeploy，2026-11-05 建），`--static-bundle` 可出静态站、`python -m nsdeploy examples/travel_agent_system.hocon --port 8080` 可跑官方旅行示例；MVP Profile 0.2 全文与计划 §2.5 一致；截止 2027-01-11 16:00 UTC、奖金、评分标准与 Devpost 页一致；NandaIndex 注册路径真实（Inkbox + AIR），但分个人 agent 与域名锚定两档。

**本次已修正的问题**：
1. "Team formed in your browser — no server involved." 是错的——demo 里 mcpgate 和 LLM router 都是服务器，trace 面板还会把调用打出来；已改为 "no central orchestrator"，与旁白一致。
2. 分镜缺了同意弹窗：规范要求联系第三方 origin 前须取得同意，原 10 秒时间线塞不下；已加 0:02–0:04 同意弹窗（一次列明三家），卡片飞入顺延到 0:04–0:07。
3. 四步验证 vs 三个勾不一致：策略检查（第④步）只在验证面板可见，不占魔法瞬间的三个勾；B 档时第三个勾显示"本地信任"；已在分镜下注明。
4. "第一条结果 10 秒内落下来"不现实（LLM 首 token 延迟）：已改为 trace 开始滚动、结果落在 0:10–0:45 段。
5. M1 增加 fallback：跨 origin 延期时退出标准降为单 origin 端到端 + mock peer，提交日不变。
6. M4 备注：元旦周录制需提前确认家里时间，魔法瞬间备用镜头争取 M3 录好。
7. §5.5 改为优先 nsDeploy 自带流程（不需要单独装 Studio）。
8. §5.7 增加 NANDA A 档与自定义域名的耦合说明。

**仍待 M0 解决的关键问题**：R10 默认模型 key 保管（最大坑，见 §7）；Devpost 是否要求作品比赛期内新建；提交表单是否单选赛道；NANDA 域名锚定 vs 个人路径；nsDeploy pin commit。

**给 Wen 的建议**：D4-A（默认 DeepSeek、你出资）需要重新评估——现成 nsDeploy 不支持发布者出资默认路径；最省事的组合是：视频用 BYOK + 你自己的 key（镜头外输入），live 站给评委的默认路径用演示 key + 余额封顶 + 赛后轮换，或接受纯 BYOK。

---

## 11. M1 进度（2026-11-16 晚，先行前几步）

- 仓库已建：https://github.com/sunwen1983/agent-bazaar（公开，MIT），commit `2c7131d`，159 个文件。
- 四个 HOCON 网络骨架已写（concierge / flights / stays / local-guide），`llm_config` 指向 `deepseek-chat`；concierge 的三个 peer  fallback 靠 `demo_mode` 标记（nsdeploy 识别 link slot 的机制）才被正确链接，初版漏了标记已修。
- 三 origin 静态包已生成并冒烟：17 个关键资源全 200（`/`、manifest、network JSON、links.json、runtime JS、site.js）；manifest 的 `linking.imports` 三条全 `linked`，`links.json` 映射正确；`_headers` 含 `Access-Control-Allow-Origin: *`（文件内容已核对；普通 HTTP 服务器不执行 `_headers`，真实 CORS 头要等 Cloudflare Pages 部署后验证）。
- origin C（stays + guide 同站）是 nsdeploy 单 bundle 单网络的限制下用 `scripts/merge-origin-c.py` 合并的：manifest `networks` 与 `linking.exports` 各两条，`_headers` 补了 guide 的 CORS 规则。
- mcpgate CORS：PR [#20](https://github.com/sunwen1983/mcpgate/pull/20)（分支 `feature/cors-browser-mcp`，commit `d6aca6e`），只开 `/mcp/` 与 `/healthz`，`Allow-Origin: *`，OPTIONS 在 auth 之前回 200；4 个新测试全过，`go test ./...` 全绿，CI build success。**未合并，等 Wen 点头。**
- 未验证：浏览器内真实跑通（本环境无浏览器）；Pages 部署后的真实 CORS 头；MCP 后端（mock 数据服务）还没写；HOCON 里的 mcpgate URL（`agent-bazaar-mcpgate.fly.dev`）与 Pages 子域都是占位，部署时确认。
- 已知注意点：Git Data API 拒绝在空仓库建 blob（409），新仓库需先用 Contents API seed 一个文件（`/tmp/seed_repo.py`）；`gh` CLI 未登录、`git push` 无 credential，分支/PR 推送走 API 脚本（`/tmp/create_pr.py`）；nsdeploy 需用 `PYTHONPATH` 指到 clone 目录才能跑（venv 里没有装 nsdeploy 包本身）。
