# Astryx 复刻路线图（mygo-astryx ← meta/astryx）

> 上游基准固定为 2026-11-14 核对的 [meta/astryx 提交 e364afc958609c990d42514b8927a6ef54b86b00](https://github.com/meta/astryx/tree/e364afc958609c990d42514b8927a6ef54b86b00)。
> `packages/core/src` 的 124 个直接子目录已分类为 107 个组件族、12 个支撑目录与 5 个内部目录，分类记录于 parity 的 category 字段。
> `packages/cli/assets/templates/blocks/components` 有 172 个 block 目录，内部共 729 个 TSX 示例文件和 728 个文档文件；`packages/cli/assets/templates/pages` 有 58 个 page 目录。
> 本项目现状：71 组件（`ds/`）、2 模板、7 主题（与 Astryx 同名）、gallery、CLI（init/search/component/docs/theme/swizzle/template/check/mcp）、MCP server。

阶段 2 按全部 58 个 page 模板交付，按复杂度分批排期。目录数与本地导出函数采用各自口径，覆盖率通过逐项映射计算。

## 总原则

1. **复刻"系统能力"，不复刻"Web 特性"。** CSS 级联层、Tailwind 桥接、CDN 免构建、`className` 覆盖是 DOM 平台专属，原生 UI 的对应物是 token 引擎本身（已有）。Storybook/docsite 的对应物是 gallery + registry（已有雏形）。
2. **每个阶段末尾跑全量验证**（AGENTS.md）：`./scripts/verify.sh`（build + vet + mygo vet + astryx check --vet + test）。
3. **新增 `ds/` 组件遵守 swizzle 安全导入约束**：只 import stdlib、`mygo/ui`、本仓库 `theme`、`icons`。
4. **优先复用已核实的 MyGo v0.8.3 原语**。阶段 0b 记录公开 API、行为差距和依赖边界，真正缺失的交互再评估组合、自绘或上游需求。

## 已核实的 MyGo 原语与接入边界

| 能力 | v0.8.3 已有接口 | 本项目需要完成的工作 |
|---|---|---|
| 日期与时间 | `ui.Calendar`、`ui.DateInput`、`ui.TimeInput` | 复用区域格式、键盘和无障碍行为，补主题封装与测试，日期范围和日期时间组合另行定义状态与校验 |
| 切换与分组 | `ui.ToggleBase`、`ui.Toggle`、`ui.ToggleGroup` | 定义选择语义、选项与 token 映射 |
| 导航与大纲 | `ui.Sidebar`、`ui.SidebarSection`、`ui.SidebarItem`、`ui.Outline`、`ui.OutlineTable` | 组合应用壳，明确导航状态、折叠与焦点行为 |
| 滚动与图像 | `ui.Scroll`、`ui.ScrollBoth`、`ui.ScrollHorizontal`、`ui.Image` | 封装滚动视觉和媒体布局，另行核对轮播分页、手势与资源加载行为 |
| 文件选择 | `mygo.Dialog.Open`、`mygo.Dialog.Save` | `FileInput` 接收调用方的选择回调，系统对话框由调用方在后台执行，结果经窗口更新发布 |
| 富文本与代码 | `markdown.View`、`markdown.Stream`、`markdown.DefaultCodeBlock`、`editor.View`、`editor.DiffView` | 复用官方子包，确定扩展包位置、主题适配和源码导出支持 |
| 绘图与图表 | `ui.NewCanvas`、`ui/chart` 的比例尺、坐标轴、图形、选择与降采样 | 按需求封装主题化图表并验证交互、无障碍与渲染 |

源码核对覆盖 [公开 UI API](https://github.com/egoist/mygo/blob/v0.8.3/ui/api_gen.go)、[文件对话框](https://github.com/egoist/mygo/blob/v0.8.3/dialog.go)、[Markdown](https://github.com/egoist/mygo/blob/v0.8.3/ui/markdown/code.go)、[Canvas](https://github.com/egoist/mygo/blob/v0.8.3/ui/canvas.go) 和 [chart](https://github.com/egoist/mygo/blob/v0.8.3/ui/chart/chart.go)。

`mygo` 根包和 `ui/markdown`、`ui/editor`、`ui/chart` 超出现行 `ds/` 导入白名单。文件选择采用调用方适配，富文本与图表在扩展包设计时明确依赖和导出方式。阶段 0b 记录这些设计前置条件。

---

## 阶段 0 — 开发基础设施与能力基线

阶段 0a 与 0b 已完成本地交付和验证，远端 CI 首次运行继续单独跟踪。AppShell 和页面的后续接口、依赖与分批安排见 [交付计划](astryx-delivery-plan.md)。

### 阶段 0a — Agent 开发基础设施

- [x] `AGENTS.md` 提供统一入口，`CLAUDE.md` 指向相同约束，README 链接维护文档。
- [x] [术语表](../CONTEXT.md)、[架构说明](architecture.md) 和 [贡献指南](../CONTRIBUTING.md) 记录模块边界、状态与主题模型、开发流程及组件完成标准。
- [x] GitHub Actions 在 push、PR 和手动触发时通过 `go.mod` 配置 Go，复用 `./scripts/verify.sh`，支持依赖与构建缓存。
- [x] `astryx init` 根据目标项目选择验证脚本或完整 Go 命令，包含测试、UI 线程与稳定 identity 指引，并验证幂等更新和用户内容保留。
- [x] registry 的导入检查覆盖所有外部模块来源，swizzle 端到端测试验证代表组件与全部文件导出的 build、Go vet、MyGo vet。
- [x] 内置指南使用完整验证要求，明确 swizzle 源文件选择和导出后的验证责任。

验收以文档链接和源码一致、上述测试通过、`./scripts/verify.sh` 全量通过为准。本地验证已完成，远端 CI 结果在工作流推送并首次触发后记录。

- [ ] 记录首次远端 CI 的运行结果。

### 阶段 0b — 能力清单与上游映射

- [x] `registry.BuildManifest()` 与 `astryx manifest` 输出全部组件、模板、token、预设、图标、指南与 parity 的机器可读 JSON。
- [x] 新 manifest 使用 schemaVersion 1，已有 `astryx search --json`、`astryx component --json` 的输出结构与 CLI/MCP 契约通过测试。
- [x] 对固定提交的 124 个 core 目录、172 个 block 目录和 58 个 page 目录记录本地对应物、原语、差距、排除理由与证据，并通过独立快照核对全部 354 项。
- [x] 阶段 0b 逐一验证当时的 61 个组件的独立 swizzle 构建、Go vet、MyGo vet，修复显式辅助源文件清单。
- [x] P1a 与 58 个页面的依赖和分批估算已写入交付计划。

本地 `./scripts/verify.sh` 全量通过。阶段 0b 验收时，107 个组件族有 58 项部分覆盖、46 项缺失、2 项延期与 1 项不适用，支撑与内部目录单独分类。page 有 2 项部分覆盖与 56 项缺失，block 的 172 个目录待场景交付。具体依据见 [阶段 0b 实施记录](phase-0b-plan.md)。

## 阶段 1 — 高价值组件补齐（分批估算）

下表为优先候选清单。新增数量和排除范围由阶段 0b 的逐项映射核定，已有等价物需比较功能、状态与无障碍行为。按使用频率分三批。

**P1a — 应用骨架族（最高优先，解锁 Patterns 层）**

本地交付已完成。新增 10 个组件，registry 共 71 个组件函数，gallery 提供对应 playground 和八个组合场景，完整验证通过。当前 core 组件族为 67 项部分覆盖、37 项缺失、2 项延期、1 项不适用。上游剩余变体继续逐项记录，详见 [P1a 验收记录](astryx-delivery-plan.md#p1a-本地交付与验收)。

| 组件 | 说明 | 前置原语 |
|---|---|---|
| `AppShell` | 顶栏+侧栏+主区+面板槽位的应用框架 | 已有布局原语即可 |
| `TopNav` / `SideNav` / `MobileNav` | 导航壳（含折叠、section、item 状态） | 已有 `ui.Sidebar` 族及本地 Drawer |
| `Layout` 族 | Header/Content/Footer/Panel 槽位布局 | — |
| `Grid` / `AspectRatio` / `Link` | 页面需要的网格、媒体比例与导航链接 | 已有 Grid、Element.AspectRatio、Link |
| `ScrollableArea` | 统一滚动容器（含溢出阴影） | 已有 `ui.Scroll` 族 |
| `MoreMenu` | 溢出操作菜单（`…`） | 已有 MenuBase |

P1a 已确定调用方持有的 AppShell 状态、导航选择、展开与窄窗口抽屉状态，以及键盘和焦点契约。响应式使用 DIP 尺寸，根据导航和内容最小宽度制定统一断点，并验证断点两侧与窗口缩放。页面模板复用这一约定。

**P1b — 表单与选择族**

| 组件 | 说明 |
|---|---|
| `RadioList` / `CheckboxList` | 列表型单选/多选（已有单项控件，补列表形态） |
| `Calendar` / `DateRangeInput` / `DateTimeInput` / `TimeInput` | 复用已存在的 Calendar、DateInput、TimeInput，补范围与日期时间组合，核对区域格式与校验 |
| `FileInput` | 使用已有原生对话框，通过调用方选择回调接入，验证取消、错误、多选与异步结果 |
| `InputGroup` / `FieldStatus` | 组合输入、字段校验态展示 |
| `Typeahead` / `MultiSelector` | 异步建议、多选（Combobox 的强化形态） |
| `ToggleButton` / `ButtonGroup` | 复用已有 ToggleBase、Toggle、ToggleGroup，明确选择与分组语义 |
| `Steps` / `PowerSearch` | 流程进度与结构化筛选；上游 Stepper 的流程语义独立映射，现有数值 Stepper 仍按原接口使用 |

**P1c — 展示与媒体族**

| 组件 | 说明 |
|---|---|
| `AvatarGroup` / `StatusDot` | 头像组溢出、状态点 |
| `List` / `OverflowList` / `Outline` | 列表、溢出折叠列表，大纲复用已有 Outline / OutlineTable |
| `HoverCard` / `InfoTip` | 富悬停卡、增强 Tooltip |
| `Markdown` / `CodeBlock` | 复用官方 markdown / editor 子包，先明确扩展包边界与主题适配 |
| `MetadataList` | 键值元数据列表（KeyValueRow 泛化） |
| `Thumbnail` / `Lightbox` / `Carousel` | 复用已有图像原语，核对资源加载、分页、手势和焦点行为 |
| `Banner` | 与现有 `Alert` 合并或别名，补 dismissible/floating 变体 |

每批合入时同步 gallery 页面、registry 元数据、swizzle 源文件清单与行为测试，按 [组件完成标准](../CONTRIBUTING.md#新增或修改组件的完成标准) 验收主题、布局、交互和无障碍行为。

新增有状态控件时评估并更新 `statefulLoopControls`，增加循环内 identity 的正反例。每个新增或改变同包依赖的组件增加独立 swizzle 编译与 vet 用例，现有 ModalCard 与 `--all` 检查继续作为基础覆盖。自动推导源文件闭包在手工清单出现实际维护问题时再评估。

## 阶段 2 — Patterns 层（按 58 个页面分批估算）

固定基准有 58 个 page 模板和 172 个 block 示例组目录。差距映射按目录记录，block 交付按具体场景记录，选择 30–40 个场景。每个新增组件出现在至少两个 block 场景中，同一个组合 block 可覆盖多个组件，以覆盖矩阵验收。

页面按 9/13/13/12/11 项分为五批，单人有效开发日初估累计 43–71 日，组件补齐与系统适配另计。P1a 初估 8–12 日。逐页前置条件和验收场景见 [交付计划](astryx-delivery-plan.md)。

- [ ] page 模板 ×58，以固定提交的 pages 目录逐项验收。先交付 blank、dashboard、settings、table、login、documentation、editor、ide、library 九类基础页面，再完成其余场景与变体。`settings-shell` 通过能力核对确定可复用部分。
- [x] `astryx template --list` 已有基础实现。
- [ ] `astryx template --skeleton` 提供带空间标注的骨架输出。
- [ ] block 模板进 registry，`astryx search` 可检索，gallery 可交互预览
- [ ] 首批依赖阶段 1a 的 AppShell 族已交付，后续按每个页面的表单、媒体和富文本依赖分批执行，已满足依赖的页面与 P1b/P1c 并行。

模板数量增长时将脚手架资源整理为独立文件并通过 `go:embed` 接入 registry。沿用现有逐模板生成后独立 build、Go vet、MyGo vet 和 Astryx check 的端到端验证。gallery 的 block 使用实时函数渲染，脚手架通过统一场景和示例内容与之核对。二进制体积与启动时间按测量结果评估。

## 阶段 3：CLI agent 能力补齐（~1–2 周）

- [ ] **`astryx doctor`** 检查实际 Go 工具链、`go.mod` / `go list -m` 解析的 MyGo 版本、有效 GOMODCACHE 与离线依赖可用性、`mygo.json`、`.mcp.json` 和托管指引版本，提供 `--ci` 非零退出码门禁。离线检查使用关闭网络下载的 Go 查询，报告缺失模块与配置来源。
- [ ] **`astryx gap-report`**：对照阶段 0 清单输出缺失能力，供上报/排期
- [ ] **`astryx discover`** 复用 Go module 解析定位外部依赖，另行定义和版本化组件元数据格式，处理 replace 与本地模块。脚手架生成与 `mygo.json` 的打包配置关系在此阶段明确。
- [ ] `docs` 子命令扩充：`tokens`、`theming`、`styling`、`principles` 等主题文档（Astryx 有 30+ 篇，先做 8–10 篇高频的）
- [ ] MCP 工具从 2 个扩到 4 个：`search`、`get`、`list`（模板/图标枚举）、`doctor`

## 阶段 4：升级与迁移机制（~1–2 周，可延后）

- [ ] registry 记录组件版本与弃用标记
- [ ] `astryx upgrade --dry-run/--apply` 使用 registry 声明迁移规则，通过 `go/ast` 定位符号和字节范围，必要时结合 `go/types` 确认符号归属，使用 `go/format` 输出。首版支持可确认归属的符号和结构体字段重命名，歧义或不支持的调用生成手工提示。
- [ ] 迁移验证同名符号隔离、注释和字符串保留、重复执行幂等性及失败不落盘。在候选结果上执行项目完整验证，并包含 Astryx check 后提交修改。
- [ ] swizzle 回迁提示：`astryx swizzle` eject 过的文件在升级时给出 diff 警告

## 阶段 5（可选）：图表

- [ ] MyGo v0.8.3 已有 Canvas 与 chart 原语。阶段 0b 核对 58 个页面的图表需求，必要的最小封装随对应页面批次推进，完整图表包结合产品需求与 Astryx charts/vega 的版本状态排期。
- [ ] 先确定图表扩展包与主题适配，再按 sparkline、折线与柱状的顺序交付，并验证选择交互、无障碍与渲染性能。

## 明确不复刻

| Astryx 能力 | 原因 |
|---|---|
| CSS `@layer` 级联、reset.css | DOM 专属 |
| Tailwind 桥接（`@theme reference`） | DOM 专属 |
| CDN 免构建模式 | DOM 专属 |
| StyleX / build 包 | 原生渲染无需样式编译 |
| npm 多包分发、pnpm monorepo | Go 单 module 即可；主题已在仓库内 |
| Storybook / 独立 docsite / Sandbox | gallery + registry 已是等价物 |
| 完整 30 语言 locales 包 | 首版按需求提供语言覆盖，复用 MyGo i18n / 区域格式，新增组件文案与格式化保留注入能力 |
| React hooks（useAnimator 等） | 帧模型不同，不等价映射 |

## 里程碑顺序与并行性

```
阶段0a ──→ 阶段0b ──→ 阶段1a ──→ 阶段2 ──→ 阶段3 ──→ 阶段4
                         ├─→ 阶段1b ┘（可与阶段2并行）
                         └─→ 阶段1c ┘
页面所需图表最小集按依赖插入阶段1/2，完整阶段5可选
```

总量与工期在阶段 0b 完成逐项映射及 58 个页面的依赖分类后重估。阶段 0–2 的主体交付以映射范围和组件完成标准为准，受原语阻塞的能力保留未完成状态。阶段 3–4 继续补齐 agent 工具与迁移体验。
