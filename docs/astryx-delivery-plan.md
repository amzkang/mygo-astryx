# P1a 与 58 页交付计划

## 阶段 0b 交付依据

能力清单由 `astryx manifest` 离线输出，差距数据位于 [parity.json](../registry/parity.json)，目录键与树 SHA 由 [独立快照](../registry/testdata/astryx-inventory.json) 核对。上游固定为 [e364afc](https://github.com/meta/astryx/tree/e364afc958609c990d42514b8927a6ef54b86b00)，MyGo 为 v0.8.3。

| 上游清单 | missing | partial | deferred | notApplicable | 总计 |
|---|---|---|---|---|---|
| core 组件族 | 37 | 67 | 2 | 1 | 107 |
| core 支撑与内部目录 | 5 | 3 | 0 | 9 | 17 |
| block | 172 | 0 | 0 | 0 | 172 |
| page | 56 | 2 | 0 | 0 | 58 |

core 明确分类为 107 个组件族、12 个支撑目录和 5 个内部目录，category 字段用于机器核对。本次按上游完整组件族的可移植行为核对。已有组件的功能子集记为 partial，并逐项写明未覆盖的行为。P1a 后共有 71 个本地组件函数，其可用接口由 manifest.components 描述，逐组件 swizzle 构建与 vet 已验证。

## P1a 的接口与顺序

P1a 的 10 个组件、registry 元数据、gallery playground 与八个组合场景已完成本地实施，交付验证见本节末尾。原估算为 8–12 个单人有效开发日，后续批次根据实际复杂度重估。

| 顺序 | 已交付接口 | 复用能力 | 核心验收 |
|---|---|---|---|
| 1 | `Layout(c, slots LayoutSlots, opts ...LayoutOpts)` | Column、Row、Scroll | 五槽组合、fill/auto、嵌套、空槽与尺寸边界 |
| 2 | `TopNav(c, slots TopNavSlots, opts ...TopNavOpts)` | Row、Button、Menu | 标题及 start/center/end、可访问名称和键盘操作 |
| 3 | `SideNav(c, state *NavState, items []NavItem, opts ...NavOpts)` | Sidebar 族 | 层级、分区、选择、折叠与稳定 identity |
| 4 | `MobileNav(c, open *bool, fn func(), opts ...NavOpts)` | 现有 Drawer | Escape、外部关闭、初始焦点与关闭恢复 |
| 5 | `AppShell(c, state *AppShellState, slots AppShellSlots, opts ...AppShellOpts)` | 上述模块 | 宽窄导航切换、主内容区域语义、键盘跳转与状态保留 |
| 6 | `ScrollableArea(c, fn func(), opts ...ScrollOpts)`、`MoreMenu(c, open *bool, items []MenuItemSpec)` | Scroll 族、MenuBase | 溢出与键盘滚动、菜单定位和选择后关闭 |
| 7 | `Grid(c, fn func(), opts ...GridOpts)`、`AspectRatio(c, ratio float32, fn func())`、`Link(c, label, url string)` | Grid、Element.AspectRatio、Link | 网格收窄、媒体比例、链接名称与跳转 |

所有接口接收当前帧的 `*ui.Context` 并返回 `ui.Element`。slots 与子内容回调统一采用现有 `ds/` 的 `func()`，通过调用方当前帧的构建闭包使用 `c`，由容器接管父节点构建栈。slots 参数在当前帧创建并执行。`NavState` 持有 SelectedID、Collapsed 和按稳定 ID 保存的 OpenSections，`AppShellState` 持有 NavState 与 MobileNavOpen。持久状态仅包含普通数据。

LayoutSlots 提供 Header、Start、Content、End、Footer。TopNavSlots 提供 Heading、Start、Center、End。AppShellSlots 提供 TopNav、SideNav、Banner、Content，窄窗口复用 SideNav 内容构建导航抽屉。NavItem 提供稳定 ID、Label、Icon、Disabled 与 Children。空槽不占空间，初始空选择允许显示未选中导航，禁用项保持选择不变。

AppShell 用于窗口根布局，基于 `c.Size()` 的 DIP 内容宽度切换。默认导航宽度为 `th.Space(60)`、内容最小宽度为 640 DIP、区域间距为 `th.Space(6)`，阈值取三者之和，comfortable 密度为 904 DIP。AppShellOpts 允许覆盖这三项，零值使用默认值。阈值按展开导航计算，手动折叠保持同一响应式规则。

进入窄窗口后侧栏呈现为抽屉，选择项目关闭抽屉。恢复宽窗口时关闭抽屉，并保留选择、手动折叠和分区状态。测试覆盖阈值两侧、三种密度、RTL、键盘焦点及窗口连续缩放。导航使用 RoleNavigation，主内容使用带 Label 的 RoleRegion，F6 在导航与主内容之间移动焦点。

首批安排八个组合场景，形成每个新增组件至少两次的覆盖关系。

| 场景 | 覆盖组件 |
|---|---|
| app-shell-basic | AppShell、TopNav、SideNav、Layout、Grid、Link、ScrollableArea |
| app-shell-narrow | AppShell、MobileNav、SideNav、Layout |
| nav-hierarchy | SideNav、TopNav、MoreMenu、Link |
| nav-collapsed | AppShell、SideNav、TopNav |
| scroll-long-content | ScrollableArea、Layout、Grid |
| media-grid | Grid、AspectRatio、Link |
| action-overflow | MoreMenu、TopNav、Link |
| responsive-media-detail | AspectRatio、Layout、MobileNav |

## P1a 本地交付与验收

- [x] 交付 Layout、TopNav、SideNav、MobileNav、AppShell、ScrollableArea、MoreMenu、Grid、AspectRatio、Link，同步元数据和独立 swizzle 源文件清单。
- [x] gallery 的 Components 提供十个组件的实际 playground，Patterns 提供上表八个组合场景。app-shell-narrow 显式提高断点以展示抽屉，nav-collapsed 初始使用图标窄栏。
- [x] 行为测试覆盖五槽布局、auto 高度、空槽、嵌套、网格与比例、URL 指针和键盘激活、禁用项、分区键盘切换、稳定 ID、抽屉焦点约束与关闭恢复、重复选择关闭、F6 和三种密度下的断点缩放。
- [x] 七套主题、浅深模式、三种密度与 RTL 的组合通过 headless 布局和语义检查，新增场景通过渲染与交互测试。
- [x] 全部 71 个组件与 `--all` 独立导出后执行 build、Go vet 和 MyGo vet。AppShell 另在导出后的独立应用中实际渲染宽窄布局。
- [x] `./scripts/verify.sh` 全量验证通过。

AppShell 按窗口根布局使用，根主题容器写作 `theme.WithTheme(c, th, func() { ... }).Fill()`，使固定高度窗口的内部滚动和 footer 获得明确的可用高度。嵌套主题容器按自身布局选择尺寸。

SideNav 采用原生 Sidebar 处理目的地选择，分区使用可聚焦按钮。方向键、Home、End 选择可用目的地，Tab 进入分区按钮并通过 Enter 或 Space 展开。图标窄栏保留名称与 tooltip，标签前缀搜索仅在展开模式可用。AppShell 在抽屉中复用调用方的 SideNav 槽位，保留其手动折叠状态。自定义导航宽度时，同步设置 SideNav 的 NavOpts 与 AppShellOpts。

本批包含同步 UI 状态与导航，异步任务验收不适用。所有有状态入口与链接已纳入循环 key 检查。真实操作系统窗口、屏幕阅读器和宿主浏览器打开行为仍需平台人工验证，本地测试通过 MyGo headless 宿主验证 URL 请求。远端 CI 首次运行继续跟踪。

P1a 完成范围按本节接口和验收记录。上游的 mega menu、侧栏拖动调宽、自动列数及 masonry、滚动有效轴 tab 顺序等能力仍在差距映射中记为 partial，并记录剩余行为。

## 后续组件的排期输入

以下本地组件作为原生扩展能力单独记录，完整接口与源文件可通过 manifest.components 查询。

| 本地能力 | 当前用途 |
|---|---|
| `Meter` | 带范围、读数和阈值色的计量控件，可用于指标与容量展示 |
| `StatTile` | 组合指标名称、数值和变化徽标的卡片，可复用到 dashboard 页面 |
| `Stepper` | 数值加减输入，可复用到数量与参数编辑，流程步骤能力规划为 `Steps` |

| 批次 | 能力与预期接口方向 | 真实前置条件 |
|---|---|---|
| P1b 表单 | RadioList、CheckboxList、FieldStatus、InputGroup、FormLayout、Calendar、DateRangeInput、DateTimeInput、TimeInput、FileInput、ToggleButton、ButtonGroup、Typeahead、MultiSelector | 组级校验、受控状态、异步与过期结果，日期范围多月和约束需要组合现有单日期原语 |
| P1b 流程与筛选 | `Steps(c, current *int, items []StepItem)`、PowerSearch | 上游 Stepper 表示流程步骤，现有数值 Stepper 的调用方继续使用原接口，流程推进需统一状态与阻断校验 |
| P1c 集合与媒体 | List、Item、MetadataList、OverflowList、StatusDot、AvatarGroup、Thumbnail、Lightbox、Carousel、ClickableCard、SelectableCard | 集合布局、图片加载、溢出和焦点，Table 的排序、分组和选择属于已有组件的能力扩展 |
| P1c 富内容 | Markdown、CodeBlock、Outline、Timestamp、Timer、Chat | 富文本扩展包的主题和导出契约，文档 Outline 需要标题锚点与滚动联动，系统语音另行验证 |
| 页面图表最小集 | sparkline、折线、柱状、分布与漏斗所需封装 | 复用 ui/chart 和 Canvas，在需要该能力的页面批次前完成主题与无障碍适配 |

FileInput 的系统对话框通过调用方回调接入，在后台执行并经窗口发布结果。Markdown、editor、chart 的子包依赖由扩展包承担，组件目录中的核心导入约束用于验证 swizzle。

## 页面批次与逐页验收

以下估算单位为单人有效开发日，前提为相应组件已完成。58 页累计约 43–71 日，组件补齐与系统适配另计。每批完成后更新剩余估算。九类基础页面先进入交付，之后按各页已满足的依赖推进；组件和页面可在依赖允许时并行。

| 批次 | 页面数 | 初估 | 主要前置条件 |
|---|---|---|---|
| A 九类基础页面 | 9 | 6–10 日 | 布局、导航、表单、表格与富文本基础封装就绪 |
| B 表单与设置变体 | 13 | 8–13 日 | Steps、组级表单、日期时间、文件选择与校验 |
| C 分析与数据工作台 | 13 | 10–16 日 | 表格插件、结构化筛选、分组展开与图表最小集 |
| D 内容与媒体变体 | 12 | 7–12 日 | 媒体比例、条目集合、文档目录与主题化代码块 |
| E 复杂工作区与通信 | 11 | 12–20 日 | Chat、Markdown、绘图、重排与系统能力验证 |

每页复用原生应用壳、主题与交互规范。表中的依赖从固定提交的 page.tsx 导入和页面行为核对，Button、Text 等基础能力随组件批次统一完善。与原生系统相关的语音、文件和权限能力在对应批次记录平台验证结果。

### A 九类基础页面

| 页面 | 主要前置能力 | 首批验收场景 |
|---|---|---|
| [blank](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/blank/page.tsx) | Layout | 最小窗口、内容区与缩放 |
| [dashboard](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/dashboard/page.tsx) | Grid、Layout、Table、图表最小集 | 全局过滤联动指标、图形和表格 |
| [settings](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/settings/page.tsx) | Grid、Layout、List、Typeahead | 分区、字段状态与保存反馈 |
| [table](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table/page.tsx) | Layout、Table | 列与行、状态标签和页头操作入口 |
| [login](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/login/page.tsx) | 基础布局与控件 | 输入校验、提交与错误提示 |
| [documentation](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/documentation/page.tsx) | ClickableCard、Grid、Layout | 组件分类、目录卡片与跳转入口呈现 |
| [editor](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/editor/page.tsx) | Layout、List、Table | 块增删、属性编辑与视口预览 |
| [ide](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/ide/page.tsx) | CodeBlock、Layout、List、MetadataList、Resizable、TreeList | 树选择、代码展示、面板缩放与页签 |
| [library](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/library/page.tsx) | ButtonGroup、Grid、Layout、OverflowList、ToggleButton | 搜索、视图切换与卡片操作 |

### B 表单与设置变体

| 页面 | 主要前置能力 | 首批验收场景 |
|---|---|---|
| [checkout-wizard](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/checkout-wizard/page.tsx) | FieldStatus、FormLayout、Layout、RadioList、Steps、Thumbnail | 步骤推进、购物数据保留与提交校验 |
| [contact-form](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/contact-form/page.tsx) | Grid、RadioList | 必填、单选、复选和提交结果 |
| [form-two-column](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/form-two-column/page.tsx) | AspectRatio、Grid | 双列收窄、输入值与标签关联 |
| [form-wizard-dialog](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/form-wizard-dialog/page.tsx) | CheckboxList、DateInput、FieldStatus、FormLayout、Layout、MetadataList、MultiSelector、RadioList、Steps、TimeInput | 弹层内步骤、日期时间与焦点恢复 |
| [form-wizard-inline](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/form-wizard-inline/page.tsx) | CheckboxList、FieldStatus、FormLayout、Layout、MetadataList、RadioList、Steps | 内联步骤、返回保值与异步提交 |
| [form-wizard-vertical](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/form-wizard-vertical/page.tsx) | AspectRatio、CheckboxList、DateInput、FieldStatus、FileInput、FormLayout、Layout、List、MetadataList、MultiSelector、RadioList、Steps | 纵向步骤、文件校验与错误定位 |
| [form-wizard](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/form-wizard/page.tsx) | CheckboxList、FormLayout、Grid、Layout、MetadataList、RadioList、SelectableCard、Steps | 步骤导航、组级校验与完成状态 |
| [payment-form](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/payment-form/page.tsx) | Grid、Layout、RadioList、Thumbnail | 金额、选择项和提交状态 |
| [settings-dialog](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/settings-dialog/page.tsx) | AspectRatio、Grid、Layout、List、SelectableCard、SideNav | 设置导航、主题切换与关闭恢复焦点 |
| [settings-sidebar](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/settings-sidebar/page.tsx) | Layout、List | 侧栏页签、字段保留与保存 |
| [login-card](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/login-card/page.tsx) | 基础布局与控件 | 卡片布局、提交与键盘操作 |
| [login-split](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/login-split/page.tsx) | Grid | 宽窄布局、输入与验证 |
| [login-sso](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/login-sso/page.tsx) | 基础布局与控件 | 身份来源选择、反馈与导航动作 |

### C 分析与数据工作台

| 页面 | 主要前置能力 | 首批验收场景 |
|---|---|---|
| [dashboard-alert-rail](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/dashboard-alert-rail/page.tsx) | Grid、Layout、List、图表最小集 | 告警筛选、侧栏和状态时间 |
| [dashboard-cohort-funnel](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/dashboard-cohort-funnel/page.tsx) | Grid、Layout、Table、图表最小集 | 筛选联动漏斗和 cohort 数据 |
| [dashboard-comparison](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/dashboard-comparison/page.tsx) | AspectRatio、Grid、Layout、Table、图表最小集 | 区间比较、图表与表格联动 |
| [dashboard-composition](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/dashboard-composition/page.tsx) | Grid、Layout、List、Table、图表最小集 | 组合筛选、图形和资源行操作 |
| [dashboard-progress](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/dashboard-progress/page.tsx) | Grid、Layout、List、Table、图表最小集 | 项目状态、进度与悬停详情 |
| [dashboard-scorecard](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/dashboard-scorecard/page.tsx) | Grid、Layout、图表最小集 | 指标切换、趋势与侧面板 |
| [table-collapsible](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table-collapsible/page.tsx) | Layout、Table、图表最小集 | 行排序、展开与折叠状态 |
| [table-comparison](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table-comparison/page.tsx) | Layout、Table、图表最小集 | 分组与固定列、比较选择 |
| [table-filter](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table-filter/page.tsx) | AspectRatio、Layout、List、MetadataList、MultiSelector、OverflowList、PowerSearch、RadioList、Resizable、Table、ToggleButton | 结构化筛选、清空、结果数与详情 |
| [table-grouped](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table-grouped/page.tsx) | Layout、MetadataList、PowerSearch、RadioList、Resizable、Table | 分组、筛选、重排与面板尺寸 |
| [table-inbox](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table-inbox/page.tsx) | Layout、OverflowList、Resizable、Table、Thumbnail、Typeahead | 列表详情切换、标签输入与编辑状态 |
| [table-page](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table-page/page.tsx) | Layout、MetadataList、PowerSearch、Table | 搜索、排序、过滤和空态 |
| [table-tree](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/table-tree/page.tsx) | Layout、Markdown、Table | 树行展开、列尺寸与 Markdown 详情 |

### D 内容与媒体变体

| 页面 | 主要前置能力 | 首批验收场景 |
|---|---|---|
| [centered-hero](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/centered-hero/page.tsx) | AspectRatio、Layout | 首屏居中、比例媒体与主操作 |
| [classic-gallery](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/classic-gallery/page.tsx) | Grid、Layout | 网格、页签和窗口收窄 |
| [gallery-hero](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/gallery-hero/page.tsx) | AspectRatio、Grid、Layout | 媒体网格、RTL 与导航动作 |
| [mixed-gallery](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/mixed-gallery/page.tsx) | AspectRatio、Layout | 多比例媒体和布局切换 |
| [product-gallery](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/product-gallery/page.tsx) | AspectRatio、Grid、Layout | 图片比例、资源选择与操作 |
| [side-gallery](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/side-gallery/page.tsx) | AspectRatio、Grid、Layout | 侧栏收窄和媒体选择 |
| [course-detail](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/course-detail/page.tsx) | AspectRatio、Grid、Layout、RadioList | 内容展开、进度与选项选择 |
| [detail-page](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/detail-page/page.tsx) | Layout、List、MetadataList、Thumbnail | 页签、元数据、详情弹层与缩放 |
| [product-detail](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/product-detail/page.tsx) | AspectRatio、Grid、Layout、SelectableCard | 变体选择、数量和折叠说明 |
| [theme-showcase](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/theme-showcase/page.tsx) | AppShell、AspectRatio、Chat、Grid、Layout、RadioList、SelectableCard、Table、TopNav | 主题作用域、全部演示控件和导航 |
| [documentation-design](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/documentation-design/page.tsx) | CodeBlock、Layout、Outline、Table | 目录定位、代码、表格与弹层 |
| [documentation-technical](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/documentation-technical/page.tsx) | CodeBlock、Layout、List、Outline | 标题锚点、代码块和目录联动 |

### E 复杂工作区与通信

| 页面 | 主要前置能力 | 首批验收场景 |
|---|---|---|
| [ai-chat-landing](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/ai-chat-landing/page.tsx) | ButtonGroup、Chat、ClickableCard、Grid、Layout、ToggleButton、Typeahead、语音适配验证 | 建议选择、输入提交与语音权限适配 |
| [ai-chat](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/ai-chat/page.tsx) | Chat、ClickableCard、CodeBlock、Layout、Markdown、Resizable | 消息流、工具调用、代码和消息操作 |
| [canvas-editor](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/canvas-editor/page.tsx) | AspectRatio、Canvas / 拖动、Layout、List、Resizable、Thumbnail、TreeList | 选区、拖动、缩放、颜色与属性同步 |
| [file-explorer](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/file-explorer/page.tsx) | Layout、List、MetadataList | 目录与文件选择、视图和元数据 |
| [incident-console](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/incident-console/page.tsx) | Layout、List、MetadataList、PowerSearch、Resizable | 事件筛选、状态时间与详情面板 |
| [kanban-board](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/kanban-board/page.tsx) | Layout、集合重排 | 跨列拖放、移动结果与键盘替代动作 |
| [messaging-shell](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/messaging-shell/page.tsx) | Chat、Layout、List | 会话选择、消息内容和输入状态 |
| [shell-nav](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/shell-nav/page.tsx) | AppShell、Layout、SideNav、TopNav、TreeList、Typeahead | 命令面板、树导航与侧栏切换 |
| [shell-side-nav](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/shell-side-nav/page.tsx) | AppShell、Layout、SideNav | 分区导航、折叠与菜单动作 |
| [shell-top-nav](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/shell-top-nav/page.tsx) | AppShell、Grid、TopNav | 顶栏、菜单和窄窗口导航 |
| [work-item-detail](https://github.com/meta/astryx/blob/e364afc958609c990d42514b8927a6ef54b86b00/packages/cli/assets/templates/pages/work-item-detail/page.tsx) | AspectRatio、CheckboxList、ClickableCard、Layout、List、MetadataList、Typeahead | 负责人、标签、评论编辑与操作菜单 |

## 验收与追踪

每页生成到独立 module，通过 build、Go vet、MyGo vet 和 Astryx check。使用 ui.NewTester 验证状态变化与关键键盘路径，通过 gallery 检查七套主题、浅深模式、三种密度及 RTL。页面名称逐项对照固定的 58 项清单，按已完成、部分完成与阻塞项记录。

30–40 个 block 的选择以新增组件覆盖矩阵为准，首批八个场景计入总量。复杂场景复用同一数据模型和示例数据，并同步脚手架与 gallery 的可见行为。所有批次完成后再核定页面与组件族的完整覆盖。
