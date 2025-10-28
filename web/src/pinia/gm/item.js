import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMResourceLogList,
  createGMItem,
  updateGMItem,
  deleteGMItem,
  getGMItem,
  batchDeleteGMItem,
  exportGMItem,
  cleanupGMItem,
  getGMItemStats,
  getGMItemTypes,
  getGMItemOperationTypes,
  getGMResourceTypeList,
  getGMResourceList
} from '@/api/gm_item'

export const useGMItemStore = defineStore('gmItem', () => {
  // 状态
  const itemList = ref([])
  const currentItem = ref(null)
  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const searchInfo = ref({
    player_id: '',
    res_type: '',
    res_id: '',
    operation_type: '',
    log_time_range: []
  })
  const itemStats = ref({
    totalRecords: 0,
    gainRecords: 0,
    consumeRecords: 0,
    tradeRecords: 0,
    systemRecords: 0
  })
  const itemTypes = ref([])
  const operationTypes = ref([])
  const resourceTypes = ref([])  // 资源类型列表
  const resourceList = ref([])    // 资源列表（根据类型动态获取）

  // 计算属性
  const hasItems = computed(() => itemList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 获取道具流水列表
  const fetchResourceLogList = async (params = {}) => {
    loading.value = true
    try {
      const response = await getGMResourceLogList({
        page: page.value,
        pageSize: pageSize.value,
        ...searchInfo.value,
        ...params
      })
      
      if (response.code === 0) {
        const list = response.data.list || []

        // 预处理数据，转换时间戳为日期时间对象
        list.forEach(item => {
          item.log_time_formatted = item.log_time ? new Date(item.log_time * 1000).toLocaleString() : '-'
        })

        itemList.value = list
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
        
        // console.log('Pinia store - itemList.value 已更新:', itemList.value)
        // console.log('Pinia store - 数据长度:', itemList.value.length)
      } else {
        throw new Error(response.msg || '获取道具流水列表失败')
      }
    } catch (error) {
      console.error('获取道具流水列表失败:', error)
      itemList.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }


  // 获取道具流水详情
  const fetchItem = async (id) => {
    try {
      const response = await getGMItem(id)
      if (response.code === 0) {
        currentItem.value = response.data
        return response.data
      } else {
        throw new Error(response.msg || '获取道具流水详情失败')
      }
    } catch (error) {
      console.error('获取道具流水详情失败:', error)
      throw error
    }
  }

  // 创建道具流水记录
  const createItem = async (itemData) => {
    try {
      const response = await createGMItem(itemData)
      if (response.code === 0) {
        // 刷新列表
        await fetchResourceLogList()
        return response.data
      } else {
        throw new Error(response.msg || '创建道具流水记录失败')
      }
    } catch (error) {
      console.error('创建道具流水记录失败:', error)
      throw error
    }
  }

  // 更新道具流水记录
  const updateItem = async (itemData) => {
    try {
      const response = await updateGMItem(itemData)
      if (response.code === 0) {
        // 刷新列表
        await fetchResourceLogList()
        return response.data
      } else {
        throw new Error(response.msg || '更新道具流水记录失败')
      }
    } catch (error) {
      console.error('更新道具流水记录失败:', error)
      throw error
    }
  }

  // 删除道具流水记录
  const removeItem = async (id) => {
    try {
      const response = await deleteGMItem({ id })
      if (response.code === 0) {
        // 刷新列表
        await fetchResourceLogList()
        return true
      } else {
        throw new Error(response.msg || '删除道具流水记录失败')
      }
    } catch (error) {
      console.error('删除道具流水记录失败:', error)
      throw error
    }
  }

  // 批量删除道具流水记录
  const batchDeleteItems = async (ids) => {
    try {
      const response = await batchDeleteGMItem({ ids })
      if (response.code === 0) {
        // 刷新列表
        await fetchResourceLogList()
        return true
      } else {
        throw new Error(response.msg || '批量删除道具流水记录失败')
      }
    } catch (error) {
      console.error('批量删除道具流水记录失败:', error)
      throw error
    }
  }

  // 导出道具流水数据
  const exportItems = async (params = {}) => {
    try {
      const response = await exportGMItem({
        ...searchInfo.value,
        ...params
      })
      return response
    } catch (error) {
      console.error('导出道具流水数据失败:', error)
      throw error
    }
  }

  // 清理旧数据
  const cleanupOldData = async (days = 30) => {
    try {
      const response = await cleanupGMItem({ days })
      if (response.code === 0) {
        // 刷新列表
        await fetchResourceLogList()
        return true
      } else {
        throw new Error(response.msg || '清理旧数据失败')
      }
    } catch (error) {
      console.error('清理旧数据失败:', error)
      throw error
    }
  }

  // 获取道具流水统计信息
  const fetchItemStats = async (params = {}) => {
    try {
      const response = await getGMItemStats({
        ...searchInfo.value,
        ...params
      })
      if (response.code === 0) {
        itemStats.value = response.data
        return response.data
      } else {
        throw new Error(response.msg || '获取道具流水统计失败')
      }
    } catch (error) {
      console.error('获取道具流水统计失败:', error)
      throw error
    }
  }

  // 获取道具类型列表
  const fetchItemTypes = async () => {
    try {
      const response = await getGMItemTypes()
      if (response.code === 0) {
        itemTypes.value = response.data
        return response.data
      } else {
        throw new Error(response.msg || '获取道具类型失败')
      }
    } catch (error) {
      console.error('获取道具类型失败:', error)
      throw error
    }
  }

  // 获取操作类型列表
  const fetchOperationTypes = async () => {
    try {
      const response = await getGMItemOperationTypes()
      if (response.code === 0) {
        operationTypes.value = response.data
        return response.data
      } else {
        throw new Error(response.msg || '获取操作类型失败')
      }
    } catch (error) {
      console.error('获取操作类型失败:', error)
      throw error
    }
  }

  // 获取资源类型列表
  const fetchResourceTypes = async () => {
    try {
      const response = await getGMResourceTypeList()
      if (response.code === 0) {
        resourceTypes.value = response.data.list || []
      } else {
        throw new Error(response.msg || '获取资源类型失败')
      }
    } catch (error) {
      console.error('获取资源类型失败:', error)
      throw error
    }
  }

  // 根据资源类型获取资源列表
  const fetchResourceList = async (resType) => {
    try {
      if (!resType) {
        resourceList.value = []
        return []
      }
      
      const response = await getGMResourceList(resType)
      if (response.code === 0) {
        resourceList.value = response.data.list || []
      } else {
        throw new Error(response.msg || '获取资源列表失败')
      }
    } catch (error) {
      console.error('获取资源列表失败:', error)
      resourceList.value = []
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
      player_id: '',
      res_type: '',
      res_id: '',
      operation_type: '',
      log_time_range: []
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
    itemList.value = []
    currentItem.value = null
    total.value = 0
    page.value = 1
    resetSearchInfo()
  }

  return {
    // 状态
    itemList,
    currentItem,
    loading,
    total,
    page,
    pageSize,
    searchInfo,
    itemStats,
    itemTypes,
    operationTypes,
    resourceTypes,
    resourceList,
    
    // 计算属性
    hasItems,
    totalPages,
    
    // 方法
    fetchResourceLogList,
    fetchItem,
    createItem,
    updateItem,
    removeItem,
    batchDeleteItems,
    exportItems,
    cleanupOldData,
    fetchItemStats,
    fetchItemTypes,
    fetchOperationTypes,
    fetchResourceTypes,
    fetchResourceList,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    clearState
  }
})
