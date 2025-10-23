<template>
  <div>
    <warning-bar title="注：GM管理 - 游戏用户管理" />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="UserID">
          <el-input v-model="searchInfo.userId" placeholder="用户UserID" />
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
        <el-table-column align="left" label="UserID" min-width="120" prop="user_id" />
        <el-table-column
          align="left"
          label="PlayerID"
          min-width="120"
          prop="player_id"
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
              :type="scope.row.ban_chat ? 'success' : 'warning'"
              link
              :icon="scope.row.ban_chat ? 'unlock' : 'lock'"
              @click="toggleBanChat(scope.row)"
            >
              {{ scope.row.ban_chat ? '解除禁言' : '禁言' }}
            </el-button>
            <el-button
              :type="scope.row.ban_login ? 'success' : 'danger'"
              link
              :icon="scope.row.ban_login ? 'unlock' : 'lock'"
              @click="toggleBanLogin(scope.row)"
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
  </div>
</template>

<script setup>
  import { useGMUserStore } from '@/pinia/gm/user'
  import WarningBar from '@/components/warningBar/warningBar.vue'

  import { watch, onMounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
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


  // 切换禁言状态
  const toggleBanChat = async (row) => {
    const action = row.ban_chat ? '解除禁言' : '禁言'
    const confirmText = row.ban_chat ? '确定要解除该玩家的禁言吗？' : '确定要禁言该玩家吗？'
    
    ElMessageBox.confirm(confirmText, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        // 调用 API 切换禁言状态
        await gmUserStore.setUserStatus(row.user_id, {
          ban_chat: !row.ban_chat
        })
        ElMessage.success(`${action}成功`)
        // 刷新列表
        await fetchUserList()
      } catch (error) {
        ElMessage.error(error.message || `${action}失败`)
      }
    })
  }

  // 切换封号状态
  const toggleBanLogin = async (row) => {
    const action = row.ban_login ? '解除封号' : '封号'
    const confirmText = row.ban_login ? '确定要解除该玩家的封号吗？' : '确定要封号该玩家吗？'
    
    ElMessageBox.confirm(confirmText, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        // 调用 API 切换封号状态
        await gmUserStore.setUserStatus(row.user_id, {
          ban_login: !row.ban_login
        })
        ElMessage.success(`${action}成功`)
        // 刷新列表
        await fetchUserList()
      } catch (error) {
        ElMessage.error(error.message || `${action}失败`)
      }
    })
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
