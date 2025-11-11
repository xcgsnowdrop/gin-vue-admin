import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMUserList,
  getGMUser,
  toggleGMBanChat,
  toggleGMBanLogin,
  batchOperateGMUser,
  exportGMUser,
} from '@/api/gm_user'
import { dateToTimestamp } from '@/utils/timestamp'

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
    startLoginTime: null,  // 登录查询开始时间
    endLoginTime: null,    // 登录查询结束时间
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
      // 准备查询参数，将日期转换为时间戳
      const queryParams = {
        page: page.value,
        pageSize: pageSize.value,
        userId: searchInfo.value.userId,
        playerId: searchInfo.value.playerId,
        uniqueId: searchInfo.value.uniqueId,
        nickname: searchInfo.value.nickname,
        // 将 Date 对象转换为时间戳（秒）
        startLoginTime: searchInfo.value.startLoginTime ? dateToTimestamp(searchInfo.value.startLoginTime) : null,
        endLoginTime: searchInfo.value.endLoginTime ? dateToTimestamp(searchInfo.value.endLoginTime) : null,
        ...params
      }
      
      const response = await getGMUserList(queryParams)
      
      if (response.code === 0) {
        const playerList = response.data.player_list || response.data.list || []

        userList.value = playerList
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
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

  // 切换禁言状态
  const toggleBanChat = async (params) => {
    try {
      const response = await toggleGMBanChat(params)
      if (response.code === 0) {
        // 刷新列表
        await fetchUserList()
        return true
      } else {
        throw new Error(response.msg || '切换禁言状态失败')
      }
    } catch (error) {
      console.error('切换禁言状态失败:', error)
      throw error
    }
  }

  // 切换封号状态
  const toggleBanLogin = async (params) => {
    try {
      const response = await toggleGMBanLogin(params)
      if (response.code === 0) {
        // 刷新列表
        await fetchUserList()
        return true
      } else {
        throw new Error(response.msg || '切换封号状态失败')
      }
    } catch (error) {
      console.error('切换封号状态失败:', error)
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
      startLoginTime: null,
      endLoginTime: null,
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
    toggleBanChat,
    toggleBanLogin,
    batchOperate,
    exportUsers,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    clearState
  }
})
