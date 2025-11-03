# Vue Store 与 View 职责划分 - 快速参考

## 快速回答

### Q1: 第59-61行的转换应该在 store 层还是 view 层？

**答：应该在 Store 层，但推荐改为转换为 Date 对象**

- ✅ 当前做法（方案B）：在 Store 层添加格式化字段
- ⭐ **推荐做法（方案A）**：在 Store 层转换为 Date 对象

### Q2: 第83行的转换应该在 store 层还是 view 层？

**答：应该在 Store 层** ✅

- 提交数据转换是业务逻辑，属于 Store 层职责

### Q3: 第395-396行的转换应该在 view 层还是 store 层？

**答：取决于方案选择**

- **方案A（推荐）**：Store 层已转换，View 层无需转换
- **方案B（当前）**：View 层需要转换（合理但非最优）

### Q4: 哪种方式更符合最佳实践？

**答：方案A（Store 层转换为 Date 对象）** ⭐

**原因**：
- Store 统一管理数据格式
- 编辑时无需转换
- 数据格式一致性好
- 符合单一数据源原则

### Q5: JavaScript 命名规范：下划线还是驼峰？

**答：推荐使用驼峰命名法（camelCase）** ✅

**理由**：
- JavaScript 官方规范
- Vue 框架约定
- 代码一致性
- IDE 友好

---

## 推荐方案实施（方案A）

### Store 层修改

```javascript
// web/src/pinia/gm/announcement.js
import { timestampToDate, convertDatesToTimestamps } from '@/utils/timestamp'

const fetchAnnouncementList = async (params = {}) => {
  // ...
  if (response.code === 0) {
    const list = response.data.announcementList || []

    // ✅ 方案A：转换为 Date 对象（基础格式）
    list.forEach(item => {
      item.startTime = timestampToDate(item.startTime)
      item.endTime = timestampToDate(item.endTime)
      item.createTime = timestampToDate(item.createTime)
    })

    announcementList.value = list
  }
}

// ✅ 提交时转换（保持不变）
const prepareSubmitData = (data) => {
  return convertDatesToTimestamps(data, ['startTime', 'endTime'])
}
```

### View 层修改

```vue
<!-- web/src/view/gm/announcement/announcement.vue -->
<template>
  <!-- ✅ 表格显示时格式化 -->
  <el-table-column label="开始时间" width="180">
    <template #default="scope">
      {{ formatDate(scope.row.startTime) }}
    </template>
  </el-table-column>
</template>

<script setup>
import { formatTimestamp } from '@/utils/timestamp'

// ✅ 格式化函数（UI 展示逻辑）
const formatDate = (date) => {
  if (!date) return '-'
  return date.toLocaleString()
  // 或使用工具函数：return formatTimestamp(date)
}

// ✅ 编辑时无需转换
const updateRow = async (row) => {
  let data = { ...row }
  // 无需转换，直接使用 Date 对象
  formData.value = data
}
</script>
```

---

## 命名规范建议

### 当前命名（下划线）→ 推荐命名（驼峰）

```javascript
// ❌ 当前（下划线命名）
item.start_time_formatted
item.end_time_formatted
item.create_time_formatted

// ✅ 推荐（驼峰命名，如果保留格式化字段）
item.startTimeFormatted
item.endTimeFormatted
item.createTimeFormatted

// ⭐ 更推荐（方案A：使用 Date 对象）
item.startTime  // Date 对象，显示时格式化
item.endTime    // Date 对象，显示时格式化
```

### 命名规范总结

| 类型 | 规范 | 示例 |
|------|------|------|
| 变量/属性 | camelCase | `userName`, `startTime` |
| 函数/方法 | camelCase | `getUserInfo`, `formatDate` |
| 常量 | UPPER_SNAKE_CASE | `API_BASE_URL`, `MAX_COUNT` |
| 组件名 | PascalCase | `UserProfile`, `AnnouncementList` |
| API 字段 | 保持后端格式 | `start_time` (通常 snake_case) |

---

## 方案对比表

| 对比项 | 方案A（推荐）⭐ | 方案B（当前） |
|--------|--------------|--------------|
| Store 层处理 | 时间戳 → Date | 时间戳 → 格式化字符串 |
| View 层表格 | 需要格式化 | 直接使用格式化字段 |
| View 层编辑 | ✅ 无需转换 | ❌ 需要转换 |
| 数据一致性 | ✅ 高 | ⚠️ 中等 |
| 代码复杂度 | ✅ 低 | ⚠️ 中等 |
| 维护成本 | ✅ 低 | ⚠️ 中等 |

**结论**：推荐使用方案A，并采用 camelCase 命名规范。

