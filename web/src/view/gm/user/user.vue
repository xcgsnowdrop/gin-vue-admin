<template>
  <div>
    <warning-bar title="注：GM管理 - 玩家管理" />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="UserID">
          <el-input v-model="searchInfo.userId" placeholder="玩家UserID" />
        </el-form-item>
        <el-form-item label="PlayerID">
          <el-input v-model="searchInfo.playerId" placeholder="玩家PlayerID" />
        </el-form-item>
        <el-form-item label="UniqueID">
          <el-input v-model="searchInfo.uniqueId" placeholder="玩家UniqueID" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" placeholder="玩家昵称" />
        </el-form-item>
        <el-form-item label="最后登录日期">
          <template #label>
            <span>
              最后登录日期
              <el-tooltip
                content="搜索范围是开始日期（包含）至结束日期（包含），可以只选择开始日期或只选择结束日期"
              >
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startLoginTime"
            type="datetime"
            placeholder="登录开始时间"
            :disabled-date="
              (time) =>
                searchInfo.endLoginTime
                  ? time.getTime() > searchInfo.endLoginTime.getTime()
                  : false
            "
            style="width: 200px"
          />
          <span style="margin: 0 8px; color: #606266;">—</span>
          <el-date-picker
            v-model="searchInfo.endLoginTime"
            type="datetime"
            placeholder="登录结束时间"
            :disabled-date="
              (time) =>
                searchInfo.startLoginTime
                  ? time.getTime() < searchInfo.startLoginTime.getTime()
                  : false
            "
            style="width: 200px"
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
      <el-table :data="tableData" row-key="user_id">
        <el-table-column
          align="left"
          label="PlayerID"
          min-width="130"
          prop="player_id"
        />
        <el-table-column
          align="left"
          label="UserID"
          min-width="120"
          prop="user_id"
        />
        <el-table-column
          align="left"
          label="UniqueID"
          min-width="120"
          prop="unique_id"
        />
        <el-table-column
          align="left"
          label="昵称"
          min-width="150"
          prop="nickname"
        />
        <el-table-column
          align="left"
          label="等级"
          min-width="80"
          prop="lv"
        />
        <el-table-column
          align="left"
          label="战力"
          min-width="180"
          prop="power"
        />
        <el-table-column
          align="left"
          label="区服"
          min-width="60"
          prop="area_id"
        />
        <el-table-column align="left" label="注册时间" min-width="180" prop="register_time_formatted" />
        <el-table-column align="left" label="最后登录时间" min-width="180" prop="login_time_formatted" />
        <el-table-column align="left" label="玩家状态" min-width="200">
          <template #default="scope">
            <div class="status-container">
              <el-tag :type="scope.row.online ? 'success' : 'info'" size="small">
                {{ scope.row.online ? '在线' : '离线' }}
              </el-tag>
              <el-tag 
                v-if="scope.row.ban_chat" 
                type="warning" 
                size="small" 
                style="margin-left: 4px;"
              >
                已禁言
              </el-tag>
              <el-tag 
                v-if="scope.row.ban_login" 
                type="danger" 
                size="small" 
                style="margin-left: 4px;"
              >
                已封号
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="view"
              @click="openPlayerDetail(scope.row.player_id)"
            >
              详情
            </el-button>
            <el-button
              :type="scope.row.ban_chat ? 'success' : 'warning'"
              link
              :icon="scope.row.ban_chat ? 'unlock' : 'lock'"
              @click="handleBanChat(scope.row)"
            >
              {{ scope.row.ban_chat ? '解除禁言' : '禁言' }}
            </el-button>
            <el-button
              :type="scope.row.ban_login ? 'success' : 'danger'"
              link
              :icon="scope.row.ban_login ? 'unlock' : 'lock'"
              @click="handleBanLogin(scope.row)"
            >
              {{ scope.row.ban_login ? '解除封号' : '封号' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    
    <!-- 玩家详情弹窗 -->
    <PlayerDetail 
      v-model="showPlayerDetail" 
      :player-id="selectedPlayerId"
      @close="handlePlayerDetailClose"
    />

    <!-- 禁言时间选择对话框 -->
    <el-dialog
      v-model="banChatDialogVisible"
      title="⏰ 设置禁言截止时间"
      width="500px"
      align-center
      destroy-on-close
    >
      <div style="padding: 20px 0;">
        <el-form :model="banChatData" label-width="120px">
          <el-form-item label="禁言截止时间" required>
            <el-date-picker
              v-model="banChatData.expireTime"
              type="datetime"
              placeholder="请选择禁言截止时间"
              style="width: 100%"
              :disabled-date="(time) => time.getTime() <= Date.now()"
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="banChatDialogVisible = false">取消</el-button>
          <el-button 
            type="warning" 
            @click="submitBanChat"
            :disabled="!banChatData.expireTime"
          >
            <el-icon><Lock /></el-icon>
            确定禁言
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 封号时间选择对话框 -->
    <el-dialog
      v-model="banLoginDialogVisible"
      title="🚫 设置封号截止时间"
      width="500px"
      align-center
      destroy-on-close
    >
      <div style="padding: 20px 0;">
        <el-form :model="banLoginData" label-width="120px">
          <el-form-item label="封号截止时间" required>
            <el-date-picker
              v-model="banLoginData.expireTime"
              type="datetime"
              placeholder="请选择封号截止时间"
              style="width: 100%"
              :disabled-date="(time) => time.getTime() <= Date.now()"
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="banLoginDialogVisible = false">取消</el-button>
          <el-button 
            type="danger" 
            @click="submitBanLogin"
            :disabled="!banLoginData.expireTime"
          >
            <el-icon><Lock /></el-icon>
            确定封号
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { useGMUserStore } from '@/pinia/gm/user'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import PlayerDetail from './components/PlayerDetail.vue'

  import { ref, watch, onMounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Lock } from '@element-plus/icons-vue'
  import { useAppStore } from "@/pinia";
  import { storeToRefs } from 'pinia'

  defineOptions({
    name: 'GMUser'
  })

  const appStore = useAppStore()
  const gmUserStore = useGMUserStore()

  // 状态和getters使用storeToRefs响应式解构
  const { 
    userList: tableData, 
    total, 
    page, 
    pageSize, 
    searchInfo,
  } = storeToRefs(gmUserStore)

  // 方法直接解构
  const {
    fetchUserList,
    resetSearchInfo,
    setPage,
    setPageSize
  } = gmUserStore


  // 玩家详情相关
  const showPlayerDetail = ref(false)
  const selectedPlayerId = ref('')
  
  // 打开玩家详情
  const openPlayerDetail = (playerId) => {
    selectedPlayerId.value = playerId
    showPlayerDetail.value = true
  }
  
  // 关闭玩家详情
  const handlePlayerDetailClose = () => {
    showPlayerDetail.value = false
    selectedPlayerId.value = ''
  }

  // 禁言相关
  const banChatDialogVisible = ref(false)
  const banChatData = ref({
    expireTime: null,
    currentRow: null  // 保存当前操作的行数据
  })

  // 封号相关
  const banLoginDialogVisible = ref(false)
  const banLoginData = ref({
    expireTime: null,
    currentRow: null  // 保存当前操作的行数据
  })

  const onSubmit = () => {
    setPage(1)
    fetchUserList()
  }

  const onReset = () => {
    resetSearchInfo()
    setPage(1)
    fetchUserList()
  }

  // 分页
  const handleSizeChange = (val) => {
    setPageSize(val)
    fetchUserList()
  }

  const handleCurrentChange = (val) => {
    setPage(val)
    fetchUserList()
  }

  watch(
    () => tableData.value,
    (newValue, oldValue) => {
      console.log('tableData 变化了')
      console.log('新值:', newValue)
      console.log('旧值:', oldValue)
    }
  )

  const initPage = async () => {
    await fetchUserList()
  }

  onMounted(() => {
    initPage()
  })


  // 处理禁言/解除禁言
  const handleBanChat = (row) => {
    if (row.ban_chat) {
      // 解除禁言：直接确认
      ElMessageBox.confirm(`确定要解除该玩家的禁言吗？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        executeBanChat(row, null)
      }).catch(() => {})
    } else {
      // 禁言：打开对话框选择时间
      banChatData.value = {
        expireTime: null,
        currentRow: row  // 保存完整的行数据
      }
      banChatDialogVisible.value = true
    }
  }

  // 执行禁言/解除禁言操作
  const executeBanChat = async (row, expireTime) => {
    const action = row.ban_chat ? '解除禁言' : '禁言'
    
    try {
      const params = {
        player_id: row.player_id,
        expire_time: expireTime
      }
      
      await gmUserStore.toggleBanChat(params)
      ElMessage.success(`${action}成功`)
      
      // 关闭对话框并重置
      banChatDialogVisible.value = false
      banChatData.value = { expireTime: null, currentRow: null }
      
      // 刷新列表
      await fetchUserList()
    } catch (error) {
      ElMessage.error(error.message || `${action}失败`)
    }
  }

  // 提交禁言（带时间选择）
  const submitBanChat = () => {
    if (!banChatData.value.expireTime) {
      ElMessage.warning('请选择禁言截止时间')
      return
    }
    
    ElMessageBox.confirm('确定要禁言该玩家吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      // 直接使用保存的行数据，无需构造 tempRow
      const expireTimestamp = Math.floor(new Date(banChatData.value.expireTime).getTime() / 1000)
      executeBanChat(banChatData.value.currentRow, expireTimestamp)
    }).catch(() => {})
  }

  // 处理封号/解除封号
  const handleBanLogin = (row) => {
    if (row.ban_login) {
      // 解除封号：直接确认
      ElMessageBox.confirm(`确定要解除该玩家的封号吗？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        executeBanLogin(row, null)
      }).catch(() => {})
    } else {
      // 封号：打开对话框选择时间
      banLoginData.value = {
        expireTime: null,
        currentRow: row  // 保存完整的行数据
      }
      banLoginDialogVisible.value = true
    }
  }

  // 执行封号/解除封号操作
  const executeBanLogin = async (row, expireTime) => {
    const action = row.ban_login ? '解除封号' : '封号'
    
    try {
      const params = {
        player_id: row.player_id,
        expire_time: expireTime
      }
      
      await gmUserStore.toggleBanLogin(params)
      ElMessage.success(`${action}成功`)
      
      // 关闭对话框并重置
      banLoginDialogVisible.value = false
      banLoginData.value = { expireTime: null, currentRow: null }
      
      // 刷新列表
      await fetchUserList()
    } catch (error) {
      ElMessage.error(error.message || `${action}失败`)
    }
  }

  // 提交封号（带时间选择）
  const submitBanLogin = () => {
    if (!banLoginData.value.expireTime) {
      ElMessage.warning('请选择封号截止时间')
      return
    }
    
    ElMessageBox.confirm('确定要封号该玩家吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      // 直接使用保存的行数据，无需构造 tempRow
      const expireTimestamp = Math.floor(new Date(banLoginData.value.expireTime).getTime() / 1000)
      executeBanLogin(banLoginData.value.currentRow, expireTimestamp)
    }).catch(() => {})
  }

</script>

<style lang="scss">
  .status-container {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 4px;
  }
  
  .status-container .el-tag {
    margin: 0;
  }
</style>
