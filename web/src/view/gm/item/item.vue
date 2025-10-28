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
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchInfo.log_time_range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
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
      <div class="gva-btn-list">
        <el-button type="success" icon="download" @click="exportData">
          导出数据
        </el-button>
        <el-button type="warning" icon="delete" @click="clearOldData">
          清理旧数据
        </el-button>
      </div>
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
        <el-table-column align="left" label="操作" min-width="200">
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
        </el-table-column>
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

    <!-- 详情对话框 -->
    <!-- <el-dialog
      v-model="detailDialogVisible"
      title="道具流水详情"
      width="600px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="流水ID">{{ detailData.ID }}</el-descriptions-item>
        <el-descriptions-item label="游戏用户ID">{{ detailData.userId }}</el-descriptions-item>
        <el-descriptions-item label="道具ID">{{ detailData.itemId }}</el-descriptions-item>
        <el-descriptions-item label="道具名称">{{ detailData.itemName }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getOperationTypeTag(detailData.operationType)">
            {{ getOperationTypeText(detailData.operationType) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="数量变化">
          <span :class="getQuantityClass(detailData.quantityChange)">
            {{ detailData.quantityChange > 0 ? '+' : '' }}{{ detailData.quantityChange }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="剩余数量">{{ detailData.remainingQuantity }}</el-descriptions-item>
        <el-descriptions-item label="操作原因" :span="2">{{ detailData.reason }}</el-descriptions-item>
        <el-descriptions-item label="操作时间" :span="2">{{ formatDateTime(detailData.createdAt) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog> -->

    <!-- 编辑对话框 -->
    <!-- <el-dialog
      v-model="editDialogVisible"
      title="编辑道具流水"
      width="500px"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-width="100px"
      >
        <el-form-item label="操作原因" prop="reason">
          <el-input
            v-model="editForm.reason"
            type="textarea"
            :rows="3"
            placeholder="请输入操作原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog> -->
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { useGMItemStore } from '@/pinia/gm/item'
import { storeToRefs } from 'pinia'

defineOptions({
  name: 'GMItem'
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
  // updateItem,
  removeItem,
  exportItems,
  cleanupOldData,
  resetSearchInfo,
  setPage,
  setPageSize
} = gmItemStore

// 详情对话框
const detailDialogVisible = ref(false)
const detailData = ref({})

// 编辑对话框
const editDialogVisible = ref(false)
const editForm = reactive({
  reason: ''
})
// const editFormRef = ref()
// const editRules = {
//   reason: [
//     { required: true, message: '请输入操作原因', trigger: 'blur' }
//   ]
// }

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

// // 格式化时间
// const formatDateTime = (dateTime) => {
//   if (!dateTime) return ''
//   return new Date(dateTime).toLocaleString('zh-CN')
// }

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

// 生成模拟数据
// const generateMockData = () => {
//   const data = []
//   const operationTypes = ['gain', 'consume', 'trade', 'system']
//   const itemNames = ['金币', '钻石', '经验药水', '装备强化石', '技能书', '宠物蛋']
  
//   for (let i = 1; i <= pageSize.value; i++) {
//     const operationType = operationTypes[Math.floor(Math.random() * operationTypes.length)]
//     const quantityChange = operationType === 'gain' 
//       ? Math.floor(Math.random() * 1000) + 1
//       : -(Math.floor(Math.random() * 100) + 1)
    
//     data.push({
//       ID: (page.value - 1) * pageSize.value + i,
//       userId: Math.floor(Math.random() * 10000) + 1000,
//       itemId: Math.floor(Math.random() * 1000) + 1,
//       itemName: itemNames[Math.floor(Math.random() * itemNames.length)],
//       operationType,
//       quantityChange,
//       remainingQuantity: Math.floor(Math.random() * 10000) + 1000,
//       reason: getRandomReason(operationType),
//       createdAt: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString()
//     })
//   }
  
//   return {
//     data,
//     total: 1000
//   }
// }

// 获取随机原因
// const getRandomReason = (type) => {
//   const reasons = {
//     gain: ['任务奖励', '活动奖励', '签到奖励', '充值获得', '系统补偿'],
//     consume: ['购买道具', '强化装备', '学习技能', '升级消耗', '交易消耗'],
//     trade: ['玩家交易', '拍卖行交易', '公会交易', '好友赠送'],
//     system: ['系统回收', '数据修正', '活动调整', '维护补偿']
//   }
//   const typeReasons = reasons[type] || ['系统操作']
//   return typeReasons[Math.floor(Math.random() * typeReasons.length)]
// }

// 分页处理
const handleCurrentChange = (val) => {
  setPage(val)
  fetchResourceLogList()
}

const handleSizeChange = (val) => {
  setPageSize(val)
  fetchResourceLogList()
}

// 查看详情
const viewDetail = (row) => {
  detailData.value = { ...row }
  detailDialogVisible.value = true
}

// 编辑记录
const editRecord = (row) => {
  editForm.reason = row.reason
  editForm.id = row.ID
  editDialogVisible.value = true
}

// // 保存编辑
// const saveEdit = async () => {
//   if (!editFormRef.value) return
  
//   try {
//     await editFormRef.value.validate()
//     await updateItem({
//       id: editForm.id,
//       reason: editForm.reason
//     })
//     ElMessage.success('保存成功')
//     editDialogVisible.value = false
//   } catch (error) {
//     ElMessage.error(error.message || '保存失败')
//   }
// }

// 删除记录
const deleteRecord = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await removeItem(row.ID)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 导出数据
const exportData = async () => {
  try {
    await exportItems()
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error(error.message || '导出失败')
  }
}

// 清理旧数据
const clearOldData = async () => {
  try {
    await ElMessageBox.confirm('确定要清理30天前的旧数据吗？此操作不可恢复！', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await cleanupOldData(30)
    ElMessage.success('清理完成')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '清理失败')
    }
  }
}

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
