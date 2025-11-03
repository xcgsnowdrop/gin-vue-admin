# Vue 时间戳处理最佳实践

## 问题背景

在处理时间戳数据时，经常会遇到以下场景：

1. **API 返回**：时间戳（秒级）
2. **表格显示**：格式化字符串（如 `2024-01-01 12:00:00`）
3. **表单编辑**：Date 对象（用于 `el-date-picker`）
4. **提交数据**：时间戳（秒级）

### 原有问题

在优化前，时间格式转换逻辑分散在多个地方：

```javascript
// ❌ 问题1：在 store 中手动转换格式化字符串
list.forEach(item => {
  item.start_time_formatted = item.startTime 
    ? new Date(item.startTime * 1000).toLocaleString() 
    : '-'
  // 重复代码...
})

// ❌ 问题2：在组件中手动转换 Date 对象
if (data.startTime && typeof data.startTime === 'number') {
  data.startTime = new Date(data.startTime * 1000)
} else if (data.startTime && typeof data.startTime === 'string') {
  const timestamp = Number(data.startTime)
  if (!isNaN(timestamp)) {
    data.startTime = new Date(timestamp * 1000)
  }
}

// ❌ 问题3：提交时手动转换时间戳
if (submitData.startTime instanceof Date) {
  submitData.startTime = Math.floor(submitData.startTime.getTime() / 1000)
} else if (submitData.startTime) {
  const date = new Date(submitData.startTime)
  if (!isNaN(date.getTime())) {
    submitData.startTime = Math.floor(date.getTime() / 1000)
  }
}
```

**存在的问题**：
- 代码重复，维护困难
- 转换逻辑分散，容易出错
- 没有统一的错误处理
- 类型检查不完善

## 优化方案

### 1. 创建统一的时间工具函数

创建 `web/src/utils/timestamp.js`，集中处理所有时间戳转换：

```javascript
// ✅ 统一的时间戳工具函数
export function timestampToDate(timestamp) {
  // 统一的转换逻辑，支持多种输入类型
}

export function dateToTimestamp(date) {
  // 统一的转换逻辑，支持 Date、字符串、时间戳
}

export function formatTimestamp(timestamp, format = 'localeString') {
  // 统一的格式化逻辑
}

export function convertTimestampsToDates(data, fields = []) {
  // 批量转换时间戳为 Date 对象
}

export function convertDatesToTimestamps(data, fields = []) {
  // 批量转换 Date 对象为时间戳
}
```

### 2. Store 层：数据预处理

**职责**：负责从 API 获取数据后的格式化处理

```javascript
// ✅ 在 store 中使用工具函数
import { formatTimestamp } from '@/utils/timestamp'

const fetchAnnouncementList = async () => {
  // ...
  list.forEach(item => {
    // 使用工具函数，代码简洁清晰
    item.start_time_formatted = formatTimestamp(item.startTime)
    item.end_time_formatted = formatTimestamp(item.endTime)
    item.create_time_formatted = formatTimestamp(item.createTime)
  })
  // ...
}
```

**优势**：
- 代码简洁，可读性强
- 统一的格式化逻辑
- 易于维护和修改

### 3. 组件层：编辑数据转换

**职责**：负责编辑时将时间戳转换为 Date 对象

```javascript
// ✅ 在组件中使用工具函数
import { convertTimestampsToDates } from '@/utils/timestamp'

const updateRow = async (row) => {
  let data = { ...row }
  
  // 批量转换时间字段，一行代码搞定
  data = convertTimestampsToDates(data, ['startTime', 'endTime'])
  
  formData.value = data
}
```

**优势**：
- 代码量大幅减少
- 支持批量转换
- 统一的转换逻辑

### 4. 提交数据转换

**职责**：负责提交时将 Date 对象转换为时间戳

```javascript
// ✅ 在 store 中使用工具函数
import { convertDatesToTimestamps } from '@/utils/timestamp'

const prepareSubmitData = (data) => {
  // 批量转换时间字段，一行代码搞定
  return convertDatesToTimestamps(data, ['startTime', 'endTime'])
}
```

**优势**：
- 统一的转换逻辑
- 减少重复代码
- 易于扩展和维护

## 最佳实践总结

### 1. 单一职责原则

每个层次的代码只负责自己的职责：

- **工具函数层** (`utils/timestamp.js`)：
  - 纯函数，只负责转换逻辑
  - 不依赖业务逻辑
  - 可复用性强

- **Store 层** (`pinia/*.js`)：
  - 数据预处理（添加格式化字段）
  - 提交数据转换（Date → Timestamp）
  - 不涉及 UI 逻辑

- **组件层** (`view/*.vue`)：
  - 编辑数据转换（Timestamp → Date）
  - 专注于 UI 交互
  - 不涉及复杂业务逻辑

### 2. DRY 原则（Don't Repeat Yourself）

**优化前**：转换逻辑重复 3-4 次
**优化后**：统一使用工具函数，一处定义，多处使用

### 3. 统一的数据流

```
API 返回 (Timestamp)
  ↓
Store: 添加格式化字段 (formatTimestamp)
  ↓
组件: 表格显示 (使用 *_formatted 字段)
  ↓
组件: 编辑时转换 (convertTimestampsToDates)
  ↓
组件: 用户编辑 (Date 对象)
  ↓
Store: 提交时转换 (convertDatesToTimestamps)
  ↓
API 提交 (Timestamp)
```

### 4. 类型安全

工具函数提供了完善的类型检查：

```javascript
// 支持多种输入类型
timestampToDate(1761926400)           // number
timestampToDate("1761926400")         // string
timestampToDate(null)                 // null
timestampToDate(new Date())           // Date（直接返回）
```

### 5. 错误处理

工具函数统一处理错误情况：

```javascript
timestampToDate(null)        // → null
timestampToDate('invalid')   // → null
formatTimestamp(null)        // → '-'
```

## 使用示例

### 示例1：在 Store 中预处理数据

```javascript
// ✅ 优化后
import { formatTimestamp } from '@/utils/timestamp'

list.forEach(item => {
  item.start_time_formatted = formatTimestamp(item.startTime)
  item.end_time_formatted = formatTimestamp(item.endTime)
})
```

### 示例2：在组件中编辑数据

```javascript
// ✅ 优化后
import { convertTimestampsToDates } from '@/utils/timestamp'

const updateRow = async (row) => {
  let data = { ...row }
  data = convertTimestampsToDates(data, ['startTime', 'endTime'])
  formData.value = data
}
```

### 示例3：提交数据

```javascript
// ✅ 优化后
import { convertDatesToTimestamps } from '@/utils/timestamp'

const prepareSubmitData = (data) => {
  return convertDatesToTimestamps(data, ['startTime', 'endTime'])
}
```

## 对比总结

### 代码量对比

| 场景 | 优化前 | 优化后 | 减少 |
|------|--------|--------|------|
| Store 格式化 | ~15 行 | ~3 行 | 80% |
| 组件编辑转换 | ~20 行 | ~2 行 | 90% |
| 提交数据转换 | ~15 行 | ~1 行 | 93% |

### 维护性对比

| 指标 | 优化前 | 优化后 |
|------|--------|--------|
| 代码重复 | 多处重复 | 统一函数 |
| 错误处理 | 分散处理 | 统一处理 |
| 类型检查 | 手动检查 | 统一检查 |
| 扩展性 | 难扩展 | 易扩展 |

## 其他优化建议

### 1. 使用 TypeScript

如果项目支持 TypeScript，可以为工具函数添加类型定义：

```typescript
export function timestampToDate(
  timestamp: number | string | null | undefined
): Date | null
```

### 2. 统一时间格式配置

可以在配置文件中统一管理时间格式：

```javascript
// config/timeFormats.js
export const TIME_FORMATS = {
  display: 'localeString',
  date: 'localeDateString',
  time: 'localeTimeString'
}
```

### 3. 考虑使用 dayjs 或 date-fns

如果项目中有更复杂的时间处理需求，可以考虑使用专业的时间库：

```javascript
import dayjs from 'dayjs'

export function formatTimestamp(timestamp, format = 'YYYY-MM-DD HH:mm:ss') {
  return dayjs.unix(timestamp).format(format)
}
```

## 总结

通过创建统一的时间工具函数，我们实现了：

1. ✅ **代码复用**：一处定义，多处使用
2. ✅ **职责分离**：每层代码只负责自己的职责
3. ✅ **易于维护**：修改转换逻辑只需改一处
4. ✅ **类型安全**：统一的类型检查和错误处理
5. ✅ **可扩展性**：易于添加新的转换功能

这种优化方式遵循了 Vue 和现代前端开发的最佳实践，使代码更加清晰、可维护和可扩展。

