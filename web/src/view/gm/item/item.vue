<template>
  <div>
    <warning-bar title="注：GM管理 - 游戏资源流水日志" />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="PlayerId">
          <el-input v-model="searchInfo.player_id" placeholder="PlayerId" />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select 
            v-model="searchInfo.res_type" 
            placeholder="请选择资源类型" 
            clearable
            @change="handleResourceTypeChange"
          >
            <el-option 
              v-for="resType in resourceTypes" 
              :key="resType.type" 
              :label="resType.name" 
              :value="resType.type" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="资源ID">
          <el-select 
            v-model="searchInfo.res_id" 
            placeholder="请选择资源ID" 
            clearable
            :disabled="!searchInfo.res_type"
          >
            <el-option 
              v-for="resource in resourceList" 
              :key="resource.id" 
              :label="resource.name" 
              :value="resource.id" 
            />
          </el-select>
        </el-form-item>
        <!-- <el-form-item label="操作类型">
          <el-select v-model="searchInfo.operation_type" placeholder="请选择操作类型" clearable>
            <el-option label="获得" value="gain" />
            <el-option label="消耗" value="consume" />
            <el-option label="交易" value="trade" />
            <el-option label="系统" value="system" />
          </el-select>
        </el-form-item> -->
        <el-form-item label="查询月份">
          <template #label>
            <span>查询月份</span>
            <el-tooltip content="选择要查询的月份，资源流水按月分表存储，一次只能查询一个月的数据" placement="top">
              <el-icon style="margin-left: 4px; color: #909399; cursor: help;">
                <QuestionFilled />
              </el-icon>
            </el-tooltip>
          </template>
          <el-date-picker
            v-model="searchInfo.month"
            type="month"
            placeholder="请选择月份"
            format="YYYY年MM月"
            value-format="YYYYMM"
            :default-value="new Date()"
          />
        </el-form-item>
        <el-form-item label="时间范围">
          <template #label>
            <span>时间范围</span>
            <el-tooltip content="在选定月份内选择具体的时间范围进行筛选" placement="top">
              <el-icon style="margin-left: 4px; color: #909399; cursor: help;">
                <QuestionFilled />
              </el-icon>
            </el-tooltip>
          </template>
          <el-date-picker
            v-model="searchInfo.log_time_range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledDate"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <!-- <div class="gva-btn-list">
        <el-button type="success" icon="download" @click="exportData">
          导出数据
        </el-button>
        <el-button type="warning" icon="delete" @click="clearOldData">
          清理旧数据
        </el-button>
      </div> -->
      <el-table :data="tableData" row-key="ID" v-loading="loading">
        <el-table-column align="left" label="ID" min-width="120" prop="_id" />
        <el-table-column
          align="left"
          label="PlayerId"
          min-width="120"
          prop="player_id"
        />
        <el-table-column
          align="left"
          label="资源类型"
          min-width="100"
          prop="res_type_name"
        />
        <el-table-column
          align="left"
          label="资源ID"
          min-width="100"
          prop="res_id"
        />
        <el-table-column
          align="left"
          label="资源名称"
          min-width="150"
          prop="res_name"
        />
        <el-table-column
          align="left"
          label="操作类型"
          min-width="100"
          prop="operation_type"
        >
          <template #default="scope">
            <el-tag
              :type="getOperationTypeTag(scope.row.operation_type)"
              size="small"
            >
              {{ getOperationTypeText(scope.row.operation_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="数量变化"
          min-width="120"
          prop="change_num"
        >
          <template #default="scope">
            <span :class="getQuantityClass(scope.row.change_num)">
              {{ scope.row.change_num > 0 ? '+' : '' }}{{ scope.row.change_num }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="剩余数量"
          min-width="120"
          prop="new_num"
        />
        <el-table-column
          align="left"
          label="操作原因"
          min-width="150"
          prop="log_text"
          show-overflow-tooltip
        />
        <el-table-column
          align="left"
          label="操作时间"
          min-width="180"
          prop="log_time_formatted"
        />
        <!-- <el-table-column align="left" label="操作" min-width="200">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="view"
              @click="viewDetail(scope.row)"
            >
              查看详情
            </el-button>
            <el-button
              type="warning"
              link
              icon="edit"
              @click="editRecord(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              icon="delete"
              @click="deleteRecord(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column> -->
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

  </div>
</template>

<script setup>
import { onMounted, watch } from 'vue'
// import { ElMessage, ElMessageBox } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { useGMItemStore } from '@/pinia/gm/item'
import { storeToRefs } from 'pinia'

defineOptions({
  name: 'GmItem'
})

const gmItemStore = useGMItemStore()

// 使用store中的状态
const { 
  itemList: tableData, 
  loading, 
  total, 
  page, 
  pageSize, 
  searchInfo,
  resourceTypes,
  resourceList
} = storeToRefs(gmItemStore)

const {
  fetchResourceLogList,
  fetchResourceTypes,
  fetchResourceList,
  // exportItems,
  // cleanupOldData,
  resetSearchInfo,
  setPage,
  setPageSize
} = gmItemStore

// 限制时间范围选择器只能在选定月份内选择
const disabledDate = (time) => {
  if (!searchInfo.value.month) {
    return false // 如果没有选择月份，不限制
  }
  
  const selectedMonth = searchInfo.value.month
  const year = parseInt(selectedMonth.substring(0, 4))
  const month = parseInt(selectedMonth.substring(4, 6)) - 1 // JavaScript月份从0开始
  
  const currentDate = new Date(time)
  const currentYear = currentDate.getFullYear()
  const currentMonth = currentDate.getMonth()
  
  // 禁用不在选定年月的日期
  return currentYear !== year || currentMonth !== month
}

// 获取操作类型标签样式
const getOperationTypeTag = (type) => {
  const tagMap = {
    gain: 'success',
    consume: 'danger',
    trade: 'warning',
    system: 'info'
  }
  return tagMap[type] || 'info'
}

// 获取操作类型文本
const getOperationTypeText = (type) => {
  const textMap = {
    gain: '获得',
    consume: '消耗',
    trade: '交易',
    system: '系统'
  }
  return textMap[type] || type
}

// 获取数量变化样式
const getQuantityClass = (quantity) => {
  return quantity > 0 ? 'text-green-600' : 'text-red-600'
}

// 查询数据
const onSubmit = () => {
  setPage(1)
  fetchResourceLogList()
}

// 重置搜索
const onReset = () => {
  resetSearchInfo()
  setPage(1)
  fetchResourceLogList()
}

// 分页处理
const handleCurrentChange = (val) => {
  setPage(val)
  fetchResourceLogList()
}

const handleSizeChange = (val) => {
  setPageSize(val)
  fetchResourceLogList()
}

// // 导出数据
// const exportData = async () => {
//   try {
//     await exportItems()
//     ElMessage.success('导出成功')
//   } catch (error) {
//     ElMessage.error(error.message || '导出失败')
//   }
// }

// // 清理旧数据
// const clearOldData = async () => {
//   try {
//     await ElMessageBox.confirm('确定要清理30天前的旧数据吗？此操作不可恢复！', '警告', {
//       confirmButtonText: '确定',
//       cancelButtonText: '取消',
//       type: 'warning'
//     })
    
//     await cleanupOldData(30)
//     ElMessage.success('清理完成')
//   } catch (error) {
//     if (error !== 'cancel') {
//       ElMessage.error(error.message || '清理失败')
//     }
//   }
// }

// 监听月份变化，清空时间范围选择
watch(
  () => searchInfo.value.month,
  (newMonth, oldMonth) => {
    if (newMonth !== oldMonth && searchInfo.value.log_time_range && searchInfo.value.log_time_range.length > 0) {
      // 检查当前选择的时间是否还在新月份范围内
      const selectedMonth = newMonth
      const year = parseInt(selectedMonth.substring(0, 4))
      const month = parseInt(selectedMonth.substring(4, 6)) - 1
      
      const startTime = new Date(searchInfo.value.log_time_range[0])
      const endTime = new Date(searchInfo.value.log_time_range[1])
      
      const startYear = startTime.getFullYear()
      const startMonth = startTime.getMonth()
      const endYear = endTime.getFullYear()
      const endMonth = endTime.getMonth()
      
      // 如果时间范围不在新月份内，清空时间范围
      if (startYear !== year || startMonth !== month || endYear !== year || endMonth !== month) {
        searchInfo.value.log_time_range = []
      }
    }
  }
)

watch(
    () => tableData.value,
    (newValue, oldValue) => {
      console.log('tableData 变化了')
      console.log('新值:', newValue)
      console.log('旧值:', oldValue)
      console.log('新值长度:', newValue?.length)
    },
    { deep: true, immediate: true }
)

// 处理资源类型变化
const handleResourceTypeChange = async (resType) => {
  // 清空资源ID选择
  searchInfo.value.res_id = ''
  
  if (resType) {
    try {
      await fetchResourceList(resType)
    } catch (error) {
      console.error('获取资源列表失败:', error)
    }
  }
}

// 初始化
onMounted(async () => {
  try {
    // 并行获取资源类型和道具列表
    await Promise.all([
      fetchResourceTypes(),
      fetchResourceLogList()
    ])
  } catch (error) {
    console.error('初始化失败:', error)
  }
})
</script>

<style scoped>
.text-green-600 {
  color: #16a34a;
}

.text-red-600 {
  color: #dc2626;
}
</style>
