# 架构与运行约束

项目使用 MyGo v0.8.3 的原生 UI 帧模型。开发流程和交付清单见 [贡献指南](../CONTRIBUTING.md)，名称定义见 [术语表](../CONTEXT.md)。

## 模块边界

| 模块 | 实际依赖与职责 |
|---|---|
| `theme` | 依赖 `mygo/ui`，定义 token、预设、密度、对比度及主题作用域 |
| `icons` | 依赖 `mygo/ui`，嵌入并解析 SVG，提供图标目录 |
| `ds` | 依赖标准库、`mygo/ui`、`theme`、`icons`，封装原生组件 |
| `registry` | 依赖 `mygo/ui`、`theme`、`icons`，保存组件文档、源文件名、指南与模板字符串 |
| `cmd/astryx` | 消费 registry，并从 `ds.SourceFS` 读取 swizzle 源码，提供 CLI 与 MCP |
| `cmd/gallery` | 消费组件和 registry，维护预览状态、交互示例与 snippet 生成 |

`registry.Components` 的文档和导出函数对应关系由测试核对。CLI 与 MCP 从同一个 registry 检索，gallery 使用元数据驱动目录，并通过自身的渲染分支展示组件。新增组件需分别维护这些接入点。

## 帧生命周期与状态

组件接收当前帧的 `*ui.Context`，构建并返回 `ui.Element`。`ui.*Parts` 提供本帧内的交互结果和子元素。持久状态由调用方的数据模型、状态指针或 MyGo 的局部状态机制持有。

跨帧保存业务数据、选择值和打开状态。帧句柄仅在当前构建路径使用。以 `Button` 的返回值读取 `.Clicked()`，在当前帧修改业务状态。以 `ModalCard` 的 `*bool` 参数保存打开状态，窗口每次构建时重新创建弹层。

重复、动态或有状态的控件在构建前使用 `c.Key(...)`，key 取稳定业务标识。异步 I/O 在 goroutine 执行，将普通结果数据通过 `Window.Update` 应用到 UI 状态，或使用 `Window.Invalidate` 请求重绘。共享数据的同步由调用方负责，goroutine 仅捕获业务数据及窗口等可长期持有的对象。

## 主题传播

`theme.WithTheme(c, th, fn)` 同时构建 `ui.Themed` 子树并压入帧内 token 栈，回调结束后弹出。栈通过 `ui.Local(c.Root(), ...)` 挂在当前 root 上。嵌套子树内使用 `theme.Of(c)` 读取最内层活动主题。

没有活动作用域时，`theme.Of` 根据 `c.Theme().Dark` 回退到 Neutral。组件从语义 token 读取视觉取值，并用 `th.Space(...)` 表达随密度变化的间距。修改作用域实现时运行主题隔离测试，检查嵌套子树与弹层的主题行为。

AppShell 组合导航与主区，并按窗口 DIP 宽度切换导航抽屉。NavState 和 AppShellState 保存普通状态，slots 使用当前帧的 `func()`。窗口根的主题容器使用 `.Fill()` 确定可用高度，内部滚动和五槽布局据此分配空间。接口和验收范围见 [P1a 交付记录](astryx-delivery-plan.md#p1a-本地交付与验收)。

## Registry、模板与发现

组件元数据包含签名、参数、变体、无障碍说明、示例和 `SourceFiles`。token、图标、预设与指南通过各自目录接入检索。CLI 的 JSON 输出与 MCP 的结构化结果读取这些现有表示。

`registry.BuildManifest()` 汇总既有资源和内嵌的 `parity.json`，由 `astryx manifest` 输出 schemaVersion 为 1 的确定性 JSON。图标只输出名称与 Go 符号。差距条目独立记录，已实现资源继续用于现有 CLI 与 MCP 检索。

上游目录快照位于 `registry/testdata/astryx-inventory.json`，记录固定 commit 和目录树 SHA。测试核对 354 个映射键、状态、资源引用、证据路径及 MyGo 原语符号。更新上游时显式更新快照、映射与版本约束。

模板以 `TemplateDoc.Files` 保存输出文件内容。`astryx template` 将其写入目标目录，端到端测试为每个模板建立临时 Go module，执行编译、Go vet、MyGo vet 和 Astryx 静态检查。block 与页面后续交付见 [路线图](astryx-parity-roadmap.md) 和 [交付计划](astryx-delivery-plan.md)。

## Swizzle 边界

`ds.SourceFS` 嵌入组件包的 Go 文件。CLI 根据组件声明的 `SourceFiles` 读取文件，改写包名，移除源码嵌入声明，再写到目标包。`--all` 选择全部非测试 Go 文件。

组件调用同包辅助函数时，`SourceFiles` 必须包含定义这些函数所需的文件。当前选择依据是显式文件清单。标准库和三个允许包之外的依赖需要先解决导出后的可用性，并调整相关约定与测试。

命令会解析、格式化导出的 Go 源码。导出后的类型检查和编译由使用方项目执行。仓库的 swizzle 端到端测试逐一导出 registry 中的全部组件，并覆盖 `--all`，检查构建、Go vet 和 MyGo vet。代表组件 ModalCard 和 AppShell 还会在独立临时应用中实际执行 `ui.Render`，构建 UI 树并核对图像尺寸，AppShell 覆盖宽窄窗口。

## 验证能力边界

`scripts/verify.sh` 是本仓库的完整验证入口。registry 测试核对导出函数覆盖、文档必填项、源文件存在性和导入约束。组件测试使用 MyGo 的 `ui.NewTester` 和 `ui.Render` 验证行为与渲染，gallery 与模板有独立测试。

渲染冒烟能发现构建异常和基础输出问题。键盘、焦点、无障碍与异步状态需要行为测试或明确的人工验证记录。原生文件对话框、系统服务和真实窗口行为按涉及平台补充验证，并在交付结果中说明覆盖范围。
