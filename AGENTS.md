# AGENTS.md

## 项目与开发入口

`mygo-astryx` is a native-UI design system, component library, multi-theme token engine, and agent-ready CLI/MCP toolkit built on [MyGo](https://github.com/egoist/mygo) v0.8.3 (inspired by [Meta Astryx](https://github.com/meta/astryx)).

- 开始修改前阅读 [CONTRIBUTING.md](CONTRIBUTING.md)，按对应修改类型核对关联交付物。
- 调整模块边界、状态管理、主题或 swizzle 时阅读 [架构说明](docs/architecture.md)。
- 确认 component、block、page template 等术语时阅读 [CONTEXT.md](CONTEXT.md)。
- 处理能力补齐与排期时阅读 [路线图](docs/astryx-parity-roadmap.md)，区分已实现能力与规划项。

## 模块职责

| 目录 | 职责 |
|---|---|
| `theme/` | token、主题定义、预设、帧内主题作用域与对比度 |
| `icons/` | 嵌入式 SVG 与图标检索 |
| `ds/` | 基于 MyGo 原语的组件及可导出的源文件 |
| `registry/` | 组件、token、预设、指南和模板的结构化元数据 |
| `cmd/astryx/` | CLI、MCP、初始化、生成与静态检查 |
| `cmd/gallery/` | 交互预览、属性 playground 和示例代码 |

## 工具链与验证

仓库内使用 `./bin/go`，版本由 `go.mod` 与 wrapper 配置。每批修改后运行统一入口，交付时报告结果及尚未验证的内容。

```sh
./scripts/verify.sh
```

脚本执行 build、Go vet、MyGo vet、Astryx check 和全部测试。新增或修改组件时同步 registry、gallery、行为测试及 swizzle 源文件清单，完成标准见贡献指南。

## MyGo Native UI Rules (`mygo vet`)

- **Elements do not outlive their frame.** Never store `ui.Element`, `*ui.Context`, or `*Parts` in structs, package variables, or closures sent to goroutines.
- **The UI thread never waits.** Run I/O or background work in goroutines and publish back via `Window.Update` / `Window.Invalidate`.
- **Theme scopes stay in-frame.** Use `theme.WithTheme(c, th, fn)` (which wraps `ui.Themed` and pushes onto a frame-local stack via `ui.Local`) and read tokens with `theme.Of(c)`.
- **Identity is stable.** Repeated or stateful controls must use `c.Key(...)` before construction.
- **Swizzle-safe component files.** Exported component files in `ds/` must keep imports limited to `github.com/egoist/mygo/ui`, `github.com/amzkang/mygo-astryx/theme`, `github.com/amzkang/mygo-astryx/icons`, and standard library packages so `astryx swizzle` can rewrite and eject them cleanly.
