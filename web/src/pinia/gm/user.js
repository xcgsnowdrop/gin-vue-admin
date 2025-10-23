import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMUserList,
  createGMUser,
  updateGMUser,
  deleteGMUser,
  getGMUser,
  setGMUserStatus,
  batchOperateGMUser,
  exportGMUser,
  getGMUserStats
} from '@/api/gm_user'

export const useGMUserStore = defineStore('gmUser', () => {
  // 状态
  const userList = ref([])
  const currentUser = ref(null)
  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const searchInfo = ref({
    userId: '',
    playerId: '',
    uniqueId: '',
    nickname: '',
  })
  const userStats = ref({
    totalUsers: 0,
    activeUsers: 0,
    newUsersToday: 0,
    onlineUsers: 0
  })

  // 计算属性
  const hasUsers = computed(() => userList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 获取用户列表
  const fetchUserList = async (params = {}) => {
    loading.value = true
    try {
      const response = await getGMUserList({
        page: page.value,
        pageSize: pageSize.value,
        ...searchInfo.value,
        ...params
      })
      
      // 调试信息 - 在Chrome开发工具中查看
      console.log('🔍 GM User API Response:', response)
      console.log('🔍 Response code:', response.code)
      console.log('🔍 Response data:', response.data)
      
      if (response.status === 0) {
        const playerList = response.data.player_list || response.data.list || []
        console.log('🔍 Player list data:', playerList)
        console.log('🔍 First player item:', playerList[0])
        
        // 预处理数据，转换时间戳为日期时间对象
        playerList.forEach(user => {
          user.register_time_formatted = user.register_time ? new Date(user.register_time * 1000).toLocaleString() : '-'
          user.login_time_formatted = user.login_time ? new Date(user.login_time * 1000).toLocaleString() : '-'
        })

        userList.value = playerList
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
        
        console.log('🔍 Updated userList.value:', userList.value)
      } else {
        throw new Error(response.msg || '获取用户列表失败')
      }
    } catch (error) {
      console.error('获取用户列表失败:', error)
      userList.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  // 获取用户详情
  const fetchUser = async (id) => {
    try {
      const response = await getGMUser(id)
      if (response.code === 0) {
        currentUser.value = response.data
        return response.data
      } else {
        throw new Error(response.msg || '获取用户详情失败')
      }
    } catch (error) {
      console.error('获取用户详情失败:', error)
      throw error
    }
  }

  // 创建用户
  const createUser = async (userData) => {
    try {
      const response = await createGMUser(userData)
      if (response.code === 0) {
        // 刷新列表
        await fetchUserList()
        return response.data
      } else {
        throw new Error(response.msg || '创建用户失败')
      }
    } catch (error) {
      console.error('创建用户失败:', error)
      throw error
    }
  }

  // 更新用户
  const updateUser = async (userData) => {
    try {
      const response = await updateGMUser(userData)
      if (response.code === 0) {
        // 刷新列表
        await fetchUserList()
        return response.data
      } else {
        throw new Error(response.msg || '更新用户失败')
      }
    } catch (error) {
      console.error('更新用户失败:', error)
      throw error
    }
  }

  // 删除用户
  const removeUser = async (id) => {
    try {
      const response = await deleteGMUser({ id })
      if (response.code === 0) {
        // 刷新列表
        await fetchUserList()
        return true
      } else {
        throw new Error(response.msg || '删除用户失败')
      }
    } catch (error) {
      console.error('删除用户失败:', error)
      throw error
    }
  }

  // 设置用户状态
  const setUserStatus = async (id, status) => {
    try {
      const response = await setGMUserStatus({
        id,
        status
      })
      if (response.code === 0) {
        // 刷新列表
        await fetchUserList()
        return true
      } else {
        throw new Error(response.msg || '设置用户状态失败')
      }
    } catch (error) {
      console.error('设置用户状态失败:', error)
      throw error
    }
  }

  // 批量操作
  const batchOperate = async (operation, userIds) => {
    try {
      const response = await batchOperateGMUser({
        operation,
        userIds
      })
      if (response.code === 0) {
        // 刷新列表
        await fetchUserList()
        return true
      } else {
        throw new Error(response.msg || '批量操作失败')
      }
    } catch (error) {
      console.error('批量操作失败:', error)
      throw error
    }
  }

  // 导出用户数据
  const exportUsers = async (params = {}) => {
    try {
      const response = await exportGMUser({
        ...searchInfo.value,
        ...params
      })
      return response
    } catch (error) {
      console.error('导出用户数据失败:', error)
      throw error
    }
  }

  // 获取用户统计信息
  const fetchUserStats = async () => {
    try {
      const response = await getGMUserStats()
      if (response.code === 0) {
        userStats.value = response.data
        return response.data
      } else {
        throw new Error(response.msg || '获取用户统计失败')
      }
    } catch (error) {
      console.error('获取用户统计失败:', error)
      throw error
    }
  }

  // 设置搜索条件
  const setSearchInfo = (info) => {
    searchInfo.value = { ...searchInfo.value, ...info }
  }

  // 重置搜索条件
  const resetSearchInfo = () => {
    searchInfo.value = {
      userId: '',
      playerId: '',
      uniqueId: '',
      nickname: '',
    }
  }

  // 设置分页
  const setPage = (newPage) => {
    page.value = newPage
  }

  const setPageSize = (newPageSize) => {
    pageSize.value = newPageSize
    page.value = 1
  }

  // 清空状态
  const clearState = () => {
    userList.value = []
    currentUser.value = null
    total.value = 0
    page.value = 1
    resetSearchInfo()
  }

  return {
    // 状态
    userList,
    currentUser,
    loading,
    total,
    page,
    pageSize,
    searchInfo,
    userStats,
    
    // 计算属性
    hasUsers,
    totalPages,
    
    // 方法
    fetchUserList,
    fetchUser,
    createUser,
    updateUser,
    removeUser,
    setUserStatus,
    batchOperate,
    exportUsers,
    fetchUserStats,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    clearState
  }
})
