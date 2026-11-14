# MyGo Astryx 术语

本项目将 Astryx 的设计系统能力映射到 MyGo 原生 UI。以下名称用于代码、文档与能力盘点。

## Language

**Component（组件）** 是通过 `ds/` 导出、可在当前帧构建的原生 UI 单元。上游组件目录与本项目导出函数需要逐项映射后计算覆盖率。

**Block（组合示例）** 是由多个组件组成、表达一个具体使用场景的可复用 UI 片段。block 是路线图中的交付类型。

**Page template（页面模板）** 是表达完整页面结构与初始交互的脚手架。当前 registry 的 `starter` 和 `settings-shell` 是已有应用模板，固定上游版本的 58 个页面模板属于后续覆盖目标。

**Patterns** 是 block 与 page template 所表达的组合层能力。

**Registry（能力目录）** 是向开发者和 agent 描述已实现组件、token、预设、指南与模板的结构化目录。

**Capability manifest（能力清单）** 是提供给开发者和 agent 的版本化能力快照，包含可用资源及其上游覆盖情况。

**Parity（差距映射）** 是固定上游条目与本地能力的对应关系，记录覆盖状态、剩余差距、可复用原语和证据。

**Token（设计变量）** 是主题提供的颜色、排版、间距、圆角、阴影和动效等设计取值。

**ThemeSet（主题集合）** 是同一命名主题的浅色和深色配对。**Density（密度）** 表示主题的布局紧凑程度。

**Theme scope（主题作用域）** 是某一帧内由 `theme.WithTheme` 覆盖的 UI 子树及其活动 token 环境。

**Frame handle（帧句柄）** 是仅在当前构建帧有效的 `ui.Element`、`*ui.Context` 或 `ui.*Parts` 等 UI 引用。

**Swizzle（源码导出）** 是将组件及其声明的辅助源文件复制到使用方项目，并改写包名以便本地维护的操作。

**Gallery** 是原生交互组件目录，包含主题切换、组件属性预览和示例代码。
