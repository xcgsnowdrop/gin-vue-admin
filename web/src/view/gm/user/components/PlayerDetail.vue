<template>
  <el-dialog
    v-model="visible"
    title="🏠 玩家详情"
    width="85%"
    :before-close="handleClose"
    align-center
    destroy-on-close
  >
    <div v-loading="loading" class="player-detail">
      <!-- 基本信息 -->
      <div class="detail-section">
        <h3 class="section-title">
          <el-icon class="title-icon"><User /></el-icon>
          基本信息
        </h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="UserID">
            {{ detail.user_id }}
          </el-descriptions-item>
          <el-descriptions-item label="PlayerID">
            {{ detail.player_id }}
          </el-descriptions-item>
          <el-descriptions-item label="UniqueID">
            {{ detail.unique_id }}
          </el-descriptions-item>
          <el-descriptions-item label="昵称">
            {{ detail.nickname }}
          </el-descriptions-item>
          <el-descriptions-item label="等级">
            {{ detail.lv }}
          </el-descriptions-item>
          <el-descriptions-item label="战力">
            {{ detail.power }}
          </el-descriptions-item>
          <el-descriptions-item label="区服">
            {{ detail.area_id }}
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{ formatTimestamp(detail.register_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="最后登录时间" :span="2">
            {{ formatTimestamp(detail.login_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="玩家状态">
            <div class="status-container">
              <el-tag :type="detail.online ? 'success' : 'info'" size="small">
                {{ detail.online ? '在线' : '离线' }}
              </el-tag>
              <el-tag 
                v-if="detail.ban_chat" 
                type="warning" 
                size="small" 
                style="margin-left: 4px;"
              >
                已禁言
              </el-tag>
              <el-tag 
                v-if="detail.ban_login" 
                type="danger" 
                size="small" 
                style="margin-left: 4px;"
              >
                已封号
              </el-tag>
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 属性信息 -->
      <!-- <div class="detail-section" v-if="detail.attrs">
        <h3 class="section-title">属性信息</h3>
        <div class="json-display">
          <vue-json-pretty 
            :data="detail.attrs" 
            :deep="10"
            :showLength="true"
            :highlightMouseoverNode="true"
          />
        </div>
      </div> -->


      <!-- 全部信息 -->
      <div class="detail-section">
        <h3 class="section-title">
          <el-icon class="title-icon"><Document /></el-icon>
          全部信息
        </h3>
        <div class="json-display">
          <vue-json-pretty 
            :data="detail" 
            :collapsedNodeLength="2"
            :deep="1"
            :showLength="true"
            :highlightMouseoverNode="true"
            :showSelectController="true"
          />
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="handleClose">
          <el-icon><Check /></el-icon>
          关闭
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { getGMUser } from '@/api/gm_user'
import { ElMessage } from 'element-plus'
import { User, Document, Check } from '@element-plus/icons-vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'

defineOptions({
  name: 'GmPlayerDetail'
})

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  playerId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'close'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const detail = ref({})

// 监听弹窗打开，加载详情
watch(visible, (newVal) => {
  if (newVal && props.playerId) {
    loadDetail()
  }
})

// 加载玩家详情
const loadDetail = async () => {
  loading.value = true
  try {
    const response = await getGMUser(props.playerId)
    if (response.code === 0) {
      detail.value = response.data.player_data
    } else {
      ElMessage.error(response.msg || '获取玩家详情失败')
    }
  } catch (error) {
    console.error('获取玩家详情失败:', error)
    ElMessage.error('获取玩家详情失败')
  } finally {
    loading.value = false
  }
}

// 格式化时间戳
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

const handleClose = () => {
  visible.value = false
  emit('close')
}
</script>

<style lang="scss" scoped>
.player-detail {
  .detail-section {
    margin-bottom: 24px;
    
    &:last-child {
      margin-bottom: 0;
    }
  }

  .section-title {
    margin: 0 0 20px 0;
    font-size: 18px;
    font-weight: 600;
    color: #303133;
    padding: 12px 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    border-radius: 8px;
    display: flex;
    align-items: center;
    gap: 8px;
    
    .title-icon {
      font-size: 20px;
    }
  }

  .status-container {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 4px;
  }

  .json-display {
    background: #f8f9fa;
    border: 1px solid #e9ecef;
    border-radius: 8px;
    padding: 20px;
    overflow-x: auto;
    max-height: 600px;
    overflow-y: auto;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    
    &:hover {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    }
  }

  :deep(.vue-json-pretty) {
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.8;
    
    .vjs-tree-node {
      &:hover {
        background: #f0f9ff;
        border-radius: 4px;
      }
    }
    
    .vjs-tree-node__content {
      padding: 4px 8px;
    }
    
    .vjs-tree-node__toggle {
      color: #409eff;
      font-size: 14px;
      
      &:hover {
        color: #66b1ff;
      }
    }
    
    .vjs-key {
      color: #881391;
      font-weight: 600;
    }
    
    .vjs-value-string {
      color: #0b7500;
    }
    
    .vjs-value-number {
      color: #1c00cf;
      font-weight: 600;
    }
    
    .vjs-value-boolean {
      color: #aa5d00;
      font-weight: 600;
    }
  }
  
  .dialog-footer {
    text-align: center;
    
    .el-button {
      min-width: 100px;
      height: 40px;
      font-size: 14px;
      
      .el-icon {
        margin-right: 4px;
      }
    }
  }
}

/* 全局样式优化 */
:deep(.el-dialog) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  padding: 20px 24px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  margin-right: 0;
  
  .el-dialog__title {
    color: #fff;
    font-size: 18px;
    font-weight: 600;
  }
  
  .el-dialog__headerbtn {
    .el-dialog__close {
      color: #fff;
      font-size: 20px;
      
      &:hover {
        color: #f0f9ff;
      }
    }
  }
}

:deep(.el-dialog__body) {
  padding: 24px;
  background: #fafbfc;
}

:deep(.el-descriptions) {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}
</style>

