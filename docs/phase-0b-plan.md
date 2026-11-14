# 阶段 0b 实施计划

## 目标与固定基准

本阶段已完成本地实现与验证，交付离线可用的能力清单、兼容现有调用方的 JSON 契约，以及固定上游版本的差距映射。P1a 与全部 58 个页面模板的依赖、接口方向和分批安排见 [交付计划](astryx-delivery-plan.md)，完成状态见 [路线图](astryx-parity-roadmap.md)。

- [x] 三份上游目录快照与 354 条差距记录完成，固定目录树 SHA 已交叉核对。
- [x] `registry.BuildManifest()` 与 `astryx manifest` 实现 schemaVersion 1，支持离线和确定性输出。
- [x] CLI/MCP 契约、错误输出、映射完整性与原语引用检查通过，61 个组件和 `--all` 独立导出构建及两类 vet 通过。
- [x] P1a 接口、状态和响应式规则，以及全部 58 页的分批估算、依赖和验收场景已记录。

以下保留实施要求，作为后续更新固定基准时的验收依据。

上游固定为 `meta/astryx` 的提交 `e364afc958609c990d42514b8927a6ef54b86b00`，MyGo 固定为当前依赖 v0.8.3。此提交在 2026-11-14 通过 GitHub API 核对。

| 清单 | 计数口径 | 已核实数量 |
|---|---|---|
| core | `packages/core/src` 的直接子目录，包含需要分类的辅助目录 | 124 |
| block | `packages/cli/assets/templates/blocks/components` 的直接子目录 | 172 |
| page | `packages/cli/assets/templates/pages` 的直接子目录 | 58 |

core 已分类为 107 个组件族、12 个支撑目录与 5 个内部目录，原始目录全部保留可追溯记录。三个清单共 354 个上游条目。页面范围已确定为全部 58 个目录，先交付原定九类基础页面，再完成其余场景和变体。

172 个 block 目录包含 729 个 TSX 示例文件和 728 个文档文件。映射粒度为目录，后续 30–40 个 block 的交付粒度为具体场景，按场景到目录与组件的关系验收。

来源为 [固定提交](https://github.com/meta/astryx/tree/e364afc958609c990d42514b8927a6ef54b86b00)、[core 清单](https://github.com/meta/astryx/tree/e364afc958609c990d42514b8927a6ef54b86b00/packages/core/src)、[block 清单](https://github.com/meta/astryx/tree/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/blocks/components) 和 [page 清单](https://github.com/meta/astryx/tree/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages)。

## 工作包一 — 上游快照与差距映射

1. 保存三个目录清单的独立快照 fixture，记录 repository、commit、来源路径和排序后的条目。快照提交到仓库，常规构建和测试使用本地数据。
2. 在 registry 增加独立的 parity 数据文件与读取类型，逐项比较上游公开 API、文档、模板和本地代码。以 `(kind, sourcePath)` 标识条目，kind 使用 `core`、`block`、`page`。
3. 每条记录包含 status、localRefs、primitives、reason 和 evidence。core 额外使用 category 分类为 component、support 或 internal。localRefs 使用现有资源的 kind 与规范名称，primitives 记录 MyGo 包路径和公开符号，证据引用固定版本源码或本地实际实现。
4. 使用下表中的状态。辅助目录、平台专属项、组件等价合并都要给出理由。已有函数名称相近时继续核对状态、交互和无障碍能力。部分实现记录剩余差距。
5. 对日期、时间、切换、导航、大纲、文件选择、富文本与图表逐项引用已核实的 MyGo 原语，记录扩展包设计、异步系统对话框或组合行为等真实前置条件。

| status | 判定与所需信息 |
|---|---|
| `implemented` | 声明范围内的能力已具备，给出有效 localRefs 和实现证据，允许多个上游条目映射到同一等价实现 |
| `partial` | 有可复用实现，给出有效 localRefs、已覆盖行为和具体剩余差距 |
| `missing` | 尚无对应交付物，说明需要补齐的能力与已有可复用原语 |
| `deferred` | 已评估并延期，说明技术或范围原因以及恢复条件 |
| `notApplicable` | 辅助目录或不适用的平台能力，给出分类与排除理由 |

验收时三份快照与映射键集合完全一致，数量分别为 124、172、58。每个条目都有明确状态与理由，已实现和部分实现的 localRefs 能在本地目录中找到。规划中的新组件通过 parity 记录表达，registry 的已实现组件列表继续接受导出函数一致性检查。

## 工作包二 — Capability manifest

新增 `registry.CapabilityManifest` 和 `registry.BuildManifest()`，由 registry 汇总现有数据与 parity 记录。新增 CLI 命令 `astryx manifest`，默认向 stdout 输出 JSON，错误沿用现有 stderr 与非零退出码处理。

| 字段 | 首版契约 |
|---|---|
| `schemaVersion` | 整数 `1`，用于新 manifest 的兼容性演进 |
| `baseline` | repository、commit、MyGo 版本及三个上游清单的来源路径 |
| `components` | 已有 `ComponentDoc` 的完整记录 |
| `templates` | 已有 `TemplateDoc` 的完整记录，包含模板文件内容 |
| `tokens`、`presets`、`guides` | 复用对应的现有文档结构 |
| `icons` | 图标 name 与 symbol，使用可序列化元数据 |
| `parity` | 工作包一的全部差距记录 |

资源数组按规范名称排序，token 按路径排序，parity 按 kind 和 sourcePath 排序。空数组输出 `[]`，避免输出生成时间等随调用改变的字段。汇总排序使用副本，避免影响现有检索顺序。所有数据通过编译时资源携带，命令运行期间支持离线使用。

更新命令 help、README 用法与相关指南。MCP 继续从同一 registry 获取已实现资源，后续阶段再扩展工具集合。

## 工作包三 — JSON 契约与 swizzle 完整性

- 通过正常 CLI 入口验证 manifest 是完整、可解析、可重复的 JSON，stdout 包含命令结果，错误进入 stderr。
- 为现有 `search --json` 的数组结果和 `component --json` 的对象结果补契约测试，保留当前字段名、可选字段、空查询、无匹配和错误行为。`schemaVersion` 用于新增 manifest。
- 比较 CLI 与 MCP 的同一组件记录。search 比较相同查询及 kind 下的结果，并核对 MCP 现有最多 30 条的限制，保留结构化输出的现有包装。
- 测试映射缺项、重复键、未知状态、悬空 localRefs、缺少理由及固定版本不一致会导致检查失败。测试读取快照和 parity 数据时使用严格解码以发现拼写错误。
- 将现有 swizzle 端到端测试扩展为逐组件导出验证，遍历 registry 中所有组件，每个独立临时 module 执行 build、Go vet 和 MyGo vet，继续覆盖 `--all`。发现不完整的 `SourceFiles` 时修正显式清单。
- 新增或改变行为的有状态组件依照贡献指南同步循环 identity 检查。阶段 0b 记录当前检查表与候选新组件的接入要求。

每个工作包完成后运行 `./scripts/verify.sh`。本阶段的自动化测试使用提交的快照，线上上游更新通过显式更新 commit 和快照进入后续变更。

## 工作包四 — P1a 与 58 页的排期输入

根据已分类的差距记录输出下一阶段候选清单，逐项列出可复用原语、预期新增 API、依赖和验收场景。先确定 AppShell 的调用方状态、导航选择、折叠、抽屉、键盘与焦点行为，再按内容最小宽度制定 DIP 断点。

58 个页面使用固定目录名称作为验收标识，记录对导航壳、表单、媒体、富文本、图表和系统集成的依赖。按基础页面、同族变体和复杂工作区分批，逐批重估工期。原语或扩展包设计尚未完成的页面标记真实阻塞项。页面需要的最小图表封装纳入该页面批次，完整图表包继续独立排期。

```text
ai-chat-landing, ai-chat, blank, canvas-editor, centered-hero,
checkout-wizard, classic-gallery, contact-form, course-detail,
dashboard-alert-rail, dashboard-cohort-funnel, dashboard-comparison,
dashboard-composition, dashboard-progress, dashboard-scorecard, dashboard,
detail-page, documentation-design, documentation-technical, documentation,
editor, file-explorer, form-two-column, form-wizard-dialog,
form-wizard-inline, form-wizard-vertical, form-wizard, gallery-hero,
ide, incident-console, kanban-board, library, login-card, login-split,
login-sso, login, messaging-shell, mixed-gallery, payment-form,
product-detail, product-gallery, settings-dialog, settings-sidebar,
settings, shell-nav, shell-side-nav, shell-top-nav, side-gallery,
table-collapsible, table-comparison, table-filter, table-grouped,
table-inbox, table-page, table-tree, table, theme-showcase, work-item-detail
```

完成条件为 manifest 与兼容性测试通过、354 个上游条目核对完成、现有组件的独立 swizzle 验证通过，以及 P1a 和全部页面具备可追溯的依赖与分批依据。在路线图中记录验证结果，并分别跟踪远端 CI 首次运行与本阶段实现进度。
