# MES 前端界面逻辑设计

## 目标

这份设计只描述前端界面逻辑和页面组织，不直接进入视觉稿和代码实现。

核心方向：
- AI 助手保持现状，不重做交互形态。
- MES 主应用继续使用当前侧边栏工作台结构。
- Admin 页面也统一进入侧边栏结构，不再像孤立后台页面。
- 主体区域以列表为核心，详情、编辑、创建使用抽屉或弹窗承接。
- 扫码页面按现场操作设计：部分流程先扫一个入口码进入操作界面，再扫物料、单体或流转码执行动作。
- 入口码和业务码由前端生成，后端只按解析后的业务 id 和动作执行。

## 总体信息架构

### App Shell

保留当前 `MESLayout` 作为主壳：

```text
┌──────────────────────────────────────────────────────────────┐
│ 顶部全局栏                                                     │
├──────────────┬───────────────────────────────────────────────┤
│ MES 侧边栏     │ 顶部页面栏：当前模块 / 快捷状态 / 工单提醒        │
│              ├───────────────────────────────────────────────┤
│ 角色菜单       │ 主体：列表 / 扫码 / 表单 / 详情抽屉                │
│              │                                               │
└──────────────┴───────────────────────────────────────────────┘
```

侧边栏只负责模块导航，不承载复杂筛选。筛选和动作放到主体列表顶部。

### AI 助手

AI 助手不动：
- 继续保留当前 `FloatingAssistant`。
- `/mes/assistant` 仍作为独立助手页。
- 普通 MES 页面右下角浮动助手继续存在。
- AI 的权限、可见数据、可操作草稿逻辑不由前端额外扩大，只调用当前用户可用接口。

## Admin 侧边栏统一

当前 admin 页面是 `/admin/userManage`、`/admin/appManage`、`/admin/chatManage` 独立页面。后续应统一挂到侧边栏壳里。

推荐结构：

```text
/mes/admin/users       员工管理
/mes/admin/apps        应用管理
/mes/admin/chats       对话管理
/mes/admin/items       物料管理
/mes/admin/processes   工艺管理
/mes/admin/orders      工程单管理
/mes/admin/flows       流转单管理
```

侧边栏中 admin 可见项：
- 员工管理
- 物料管理
- 工艺管理
- 工程单管理
- 流转单管理
- 应用管理
- 对话管理

Admin 主体仍然以列表为主：

```text
┌──────────────────────────────────────────────────────────────┐
│ 员工管理                         搜索用户  角色筛选  新建员工 │
├──────────────────────────────────────────────────────────────┤
│ 表格：姓名 / 账号 / 角色 / 状态 / 更新时间 / 操作              │
│ 详情抽屉：基础信息 / 角色 / 最近活动 / 保存                   │
└──────────────────────────────────────────────────────────────┘
```

## 主体列表页面规范

所有业务主体页面默认采用：

```text
列表页 = 顶部工具条 + 状态 tabs + 列表/表格 + 右侧详情抽屉
```

### 顶部工具条

包含：
- 模块标题
- 搜索框
- 主要筛选
- 创建草稿按钮
- 扫码入口按钮

示例：

```text
工程单                         搜索工程单/物料
[全部] [草稿] [已提交] [已完成]          + 新建工程单草稿  扫码
```

### 列表行

列表行展示业务判断所需的最少信息，不把详情堆在列表里。

工程单行：
- 工程单号
- 物料
- 工艺
- 负责人
- 预计数量 / 已生产数量
- 状态
- 更新时间
- 操作：查看、编辑草稿、提交草稿

流转单行：
- 流转单号
- 发起人 / 接收人
- 类型
- 状态
- 物料摘要
- 更新时间
- 操作：查看、编辑草稿、提交、审核

单体行：
- 单体码
- 物料
- 工程单
- 库存状态
- 质检状态
- 更新时间
- 操作：查看、改质检状态、加入流转草稿

### 详情和编辑

列表点击行打开右侧抽屉：

```text
┌──────────────────────┐
│ 工程单 #1024          │
│ 状态：草稿            │
├──────────────────────┤
│ 基础信息              │
│ 物料 / 工艺 / 负责人   │
│ 数量 / 说明           │
├──────────────────────┤
│ 关联单体              │
│ UNIT-001 ...          │
├──────────────────────┤
│ 保存草稿  提交         │
└──────────────────────┘
```

规则：
- 草稿可编辑。
- 非草稿只读，除非角色有明确审核/状态操作权限。
- 创建不跳页面，优先使用抽屉。
- 大表单可用全屏抽屉。

## 扫码体系

扫码分两类码：

1. 业务对象码：指向某个已有对象。
2. 操作入口码：指向某个操作上下文，扫完后进入对应操作界面。

### 业务对象码

沿用当前 `MES:<KIND>:<ID>` 结构：

```text
MES:ITEM_UNIT:123
MES:FLOW:456
MES:ENGINEERING_ORDER:789
MES:WORK_ORDER:321
```

用途：
- 打开详情。
- 添加到当前操作。
- 作为某个操作的目标对象。

### 操作入口码

新增前端约定：

```text
MES:ACTION:<ACTION_NAME>:<TARGET_KIND>:<TARGET_ID>
```

示例：

```text
MES:ACTION:ENGINEERING_INBOUND:ENGINEERING_ORDER:1024
MES:ACTION:QUALITY_INSPECT:ENGINEERING_ORDER:1024
MES:ACTION:FLOW_RECEIVE:FLOW:2048
MES:ACTION:FLOW_OUTBOUND:FLOW:2049
MES:ACTION:ITEM_UNIT_DETAIL:ITEM_UNIT:889
```

说明：
- 入口码由前端根据当前页面、当前对象和当前动作生成。
- 扫入口码后，前端进入一个明确的操作模式。
- 进入操作模式后，再扫业务对象码执行具体动作。
- 后端不需要理解入口码字符串，前端解析后调用现有接口。

## 扫码页面结构

扫码页移动端优先，桌面端兼容扫码枪。

### 扫码入口页

默认扫码页面：

```text
┌────────────────────┐
│ 扫码操作      手输  │
├────────────────────┤
│                    │
│     相机取景区       │
│   ┌────────────┐   │
│   │            │   │
│   └────────────┘   │
│                    │
├────────────────────┤
│ 当前：等待入口码     │
│ 可扫：工程单入口码   │
│      流转单入口码    │
│      单体码          │
├────────────────────┤
│ 最近扫描             │
│ ...                 │
└────────────────────┘
```

如果扫到业务对象码：
- `ITEM_UNIT`：进入单体详情或提示选择动作。
- `FLOW`：进入流转单领取/审核/查看，具体取决于角色和状态。
- `ENGINEERING_ORDER`：进入工程单上下文，可选择入库、质检、查看单体。
- `WORK_ORDER`：进入工单详情或标记已读。

如果扫到操作入口码：
- 直接进入对应操作界面。

### 操作界面

操作界面由三块组成：

```text
┌────────────────────┐
│ 工程单入库      退出 │
├────────────────────┤
│ 相机取景区 / 扫码枪输入 │
├────────────────────┤
│ 上下文               │
│ 工程单 #1024         │
│ 物料：电机外壳        │
│ 已扫：18 / 200       │
├────────────────────┤
│ 最近扫描             │
│ ✓ UNIT-00821 已加入  │
│ ✕ UNIT-00819 非本单  │
├────────────────────┤
│ 保存草稿 / 完成       │
└────────────────────┘
```

操作界面不是详情页，它只回答四件事：
- 当前在做什么。
- 当前上下文是什么。
- 刚扫到了什么。
- 下一步能做什么。

## 典型扫码流程

### 流程一：工程单入库

入口码：

```text
MES:ACTION:ENGINEERING_INBOUND:ENGINEERING_ORDER:1024
```

状态机：

```text
等待入口码
  -> 扫工程单入库入口码
  -> 加载工程单详情
  -> 进入工程单入库模式
  -> 连续扫 ITEM_UNIT 码
  -> 校验单体是否属于目标物料/是否可入库
  -> 加入入库流转草稿
  -> 保存草稿或提交
```

前端动作：
- 解析入口码，得到 `engineeringOrderId = 1024`。
- 调 `getEngineeringOrder`。
- 后续每扫一个 `ITEM_UNIT`，调 `getItemUnit`。
- 合法则加入本地待提交列表。
- 保存时调 `createInventoryFlowDraft` 或 `updateInventoryFlowDraft`。
- 提交时调 `submitInventoryFlow`。

### 流程二：质检单体

入口码：

```text
MES:ACTION:QUALITY_INSPECT:ENGINEERING_ORDER:1024
```

状态机：

```text
等待入口码
  -> 扫质检入口码
  -> 进入质检模式
  -> 扫 ITEM_UNIT
  -> 展示单体信息
  -> 选择 合格 / 不合格 / 待检
  -> 调 updateItemUnitStatus
  -> 继续扫描
```

界面重点：
- 三个大按钮：合格、不合格、待检。
- 最近扫描列表保留成功和异常。
- 如果单体不属于该工程单，显示错误但不退出模式。

### 流程三：领取流转单

入口码：

```text
MES:ACTION:FLOW_RECEIVE:FLOW:2048
```

状态机：

```text
等待入口码
  -> 扫流转领取入口码
  -> 加载流转单
  -> 进入领取模式
  -> 扫 ITEM_UNIT 确认领取
  -> 本地记录已确认单体
  -> 全部完成后提交确认
```

界面重点：
- 显示流转单目标数量和已确认数量。
- 单体不在流转单内时提示异常。
- 已扫重复单体时提示重复，不重复加入。

### 流程四：出库流转

入口码：

```text
MES:ACTION:FLOW_OUTBOUND:FLOW:2049
```

状态机：

```text
等待入口码
  -> 扫出库入口码
  -> 进入出库模式
  -> 连续扫 ITEM_UNIT
  -> 加入出库流转草稿
  -> 保存草稿
  -> 提交流转单
```

如果没有已有流转单，也可以从列表页点击“新建出库草稿”，进入同一个操作界面，只是没有入口码预设上下文。

## 二维码生成位置

### 列表行生成

每个可被扫描的对象行提供二维码入口：
- 工程单：对象码 + 操作入口码。
- 流转单：对象码 + 操作入口码。
- 单体：对象码。
- 工单：对象码。

不要把二维码直接塞进表格行，表格只放一个二维码按钮。点击后打开小弹层：

```text
┌──────────────────────┐
│ 工程单 #1024          │
├──────────────────────┤
│ [对象码] [入库入口] [质检入口] │
│                      │
│       QR Code        │
│                      │
│ MES:ACTION:...       │
│ 复制 / 打印           │
└──────────────────────┘
```

### 详情抽屉生成

详情抽屉中提供完整二维码区域：
- 对象码。
- 当前角色可用的操作入口码。
- 复制文本。
- 打印标签。

### 扫码页生成

扫码页本身不生成入口码，扫码页负责消费入口码。

## 角色和页面

### 采购

侧边栏：
- 物品建档
- 扫描入库
- 流转单

主体：
- 物品建档：列表 + 新建物品/单体。
- 扫描入库：扫码操作页。
- 流转单：列表 + 详情抽屉。

### 普通员工

侧边栏：
- 新增单品
- 领取货物
- 检验单品

主体：
- 新增单品：工程单上下文 + 单体生成。
- 领取货物：先扫流转入口码，再扫单体确认。
- 检验单品：先扫质检入口码，再扫单体并更新质量。

### 组长

侧边栏：
- 工程单
- 申请物资
- 入库流转单
- 发工单

主体：
- 工程单：列表为主，草稿抽屉创建/编辑。
- 申请物资：创建流转草稿。
- 入库流转单：工程单入库扫码或列表。
- 发工单：工单列表和创建草稿。

### 仓库管理员

侧边栏：
- 审批流转单
- 物资情况
- 物资类型
- 发工单

主体：
- 审批流转单：列表 + 审核抽屉。
- 物资情况：物料列表 + 单体列表。
- 物资类型：物料管理。
- 发工单：工单列表。

### 销售

侧边栏：
- 领取物资
- 申请物资

主体：
- 领取物资：扫码领取或列表查看。
- 申请物资：创建流转草稿。

### Admin

侧边栏：
- 员工管理
- 物料管理
- 工艺管理
- 工程单管理
- 流转单管理
- 应用管理
- 对话管理

主体全部列表化。

## 路由建议

保留现有角色工作台路由：

```text
/mes/purchase?panel=add
/mes/purchase?panel=scan
/mes/purchase?panel=flows
/mes/worker?panel=add
/mes/worker?panel=receive
/mes/worker?panel=inspect
/mes/leader?panel=engineering
/mes/leader?panel=material
/mes/leader?panel=inbound
/mes/leader?panel=workorder
/mes/warehouse?panel=audit
/mes/warehouse?panel=inventory
/mes/warehouse?panel=item
/mes/warehouse?panel=workorder
/mes/sales?panel=receive
/mes/sales?panel=apply
```

新增统一扫码路由：

```text
/mes/scan
/mes/scan?action=ENGINEERING_INBOUND&engineeringOrderId=1024
/mes/scan?action=QUALITY_INSPECT&engineeringOrderId=1024
/mes/scan?action=FLOW_RECEIVE&flowId=2048
/mes/scan?action=FLOW_OUTBOUND&flowId=2049
```

Admin 迁移路由：

```text
/mes/admin/users
/mes/admin/apps
/mes/admin/chats
/mes/admin/items
/mes/admin/processes
/mes/admin/orders
/mes/admin/flows
```

旧 `/admin/*` 可以保留重定向，避免已有链接失效。

## 前端解析规则

### 解析业务对象码

当前已有：

```text
MES:FLOW:<id>
MES:ITEM_UNIT:<id>
MES:ENGINEERING_ORDER:<id>
MES:WORK_ORDER:<id>
```

解析结果：

```ts
{
  type: 'flow' | 'itemUnit' | 'engineeringOrder' | 'workOrder',
  id: number
}
```

### 解析操作入口码

新增解析：

```text
MES:ACTION:<action>:<targetKind>:<targetId>
```

解析结果：

```ts
{
  type: 'action',
  action: 'ENGINEERING_INBOUND' | 'QUALITY_INSPECT' | 'FLOW_RECEIVE' | 'FLOW_OUTBOUND' | 'ITEM_UNIT_DETAIL',
  targetKind: 'ENGINEERING_ORDER' | 'FLOW' | 'ITEM_UNIT',
  targetId: number
}
```

### 扫码分发

伪逻辑：

```text
scan(value):
  parsed = parseMesCode(value)

  if parsed.type == action:
    enterActionMode(parsed.action, parsed.targetKind, parsed.targetId)
    return

  if currentActionMode exists:
    handleObjectInCurrentMode(parsed)
    return

  openObjectDetailOrAskAction(parsed)
```

## 操作模式状态

统一状态结构：

```ts
type ScanModeState = {
  action: string
  targetKind: string
  targetId: number
  targetLabel: string
  loaded: boolean
  records: ScanRecord[]
  draftId?: number
}

type ScanRecord = {
  code: string
  objectType: string
  objectId: number
  status: 'success' | 'warning' | 'error'
  message: string
  time: string
}
```

所有扫码操作都往 `records` 写结果，页面底部展示最近记录。

## 页面优先级

第一阶段只做界面逻辑和结构统一：
1. Admin 合入 MES 侧边栏。
2. 抽出通用列表页结构。
3. 新增统一扫码页。
4. 扩展 `mesCode.ts` 支持操作入口码。
5. 在工程单/流转单/单体详情里生成二维码。

第二阶段做真实相机扫码：
1. 移动端调用摄像头。
2. 桌面端保持扫码枪输入框。
3. 增加连续扫码、防重复、错误保留。

第三阶段做体验增强：
1. 打印标签。
2. 批量生成入口码。
3. 异常记录导出。
4. 离线缓存最近扫码记录。

## 不做的事

当前阶段不做：
- 不重做 AI 助手。
- 不把扫码页做成营销式页面。
- 不把主列表改成大卡片网格。
- 不让前端绕过后端权限。
- 不让二维码承载敏感权限，只承载动作上下文和对象 id。

## 实现注意点

- 二维码由前端生成，但所有写操作仍调用后端接口。
- 扫入口码只是进入界面，不代表操作成功。
- 扫单体码后必须二次校验对象状态。
- 所有“提交、审核、删除”仍由按钮触发，连续扫码默认只加入草稿或本地待处理列表。
- 异常不要只用 toast，扫码页必须保留错误记录。
- 移动端扫码页应该减少表格，桌面端列表页可以保持表格密度。

## 后端复杂度约束

前端列表统一使用游标语义，不依赖传统页码的真实总数：
- 列表继续保留 `pageNum/pageSize` 兼容旧组件。
- 新页面优先使用 `hasMore + nextCursor*` 加载下一页。
- `totalRow/totalPage` 在 cursor 列表里只作为兼容字段，不作为真实总数展示。

库存、物品、工艺、工程单、流转单、工单列表都应满足：
- 主查询不执行 `COUNT(*)`。
- 主查询不使用 `OFFSET`。
- 按索引定位后读取 `pageSize + 1` 条记录。
- 复杂度为 `O(log n + m)`。

工单列表额外要求：
- 按 `updated_at DESC, id DESC` 排序。
- 发起人列表走 `(deleted_at, from_user_id, updated_at, id)`。
- 接收人列表走 `(deleted_at, to_user_id, updated_at, id)`。
- 未读筛选走 `(deleted_at, from_user_id, read_status, updated_at, id)` 或 `(deleted_at, to_user_id, read_status, updated_at, id)`。
- 旧的 `created_at` 列表索引不再作为主路径索引。
