import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMResourceLogList,
  exportGMItem,
  getGMResourceTypeList,
  getGMResourceList
} from '@/api/gm_item'
import { dateToTimestamp } from '@/utils/timestamp'

export const useGMItemStore = defineStore('gmItem', () => {
  // 状态
  const itemList = ref([])
  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  // 获取当前月份，格式为YYYYMM
  const getCurrentMonth = () => {
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    return `${year}${month}`
  }

  const searchInfo = ref({
    player_id: '',
    res_type: '',
    res_id: '',
    month: getCurrentMonth(), // 默认当前月份
    log_time_range: []
  })

  const resourceTypes = ref([])  // 资源类型列表
  const resourceList = ref([])    // 资源列表（根据类型动态获取）

  // 计算属性
  const hasItems = computed(() => itemList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 获取道具流水列表
  const fetchResourceLogList = async (params = {}) => {
    loading.value = true
    try {
      // 处理时间范围转换
      const processedSearchInfo = { ...searchInfo.value }
      if (processedSearchInfo.log_time_range && processedSearchInfo.log_time_range.length === 2) {
        // 将日期时间字符串转换为时间戳（秒）
        processedSearchInfo.start_time = dateToTimestamp(processedSearchInfo.log_time_range[0])
        processedSearchInfo.end_time = dateToTimestamp(processedSearchInfo.log_time_range[1])
        
        // 删除原始的时间范围字段，避免传给后端
        delete processedSearchInfo.log_time_range
      }
      
      const response = await getGMResourceLogList({
        page: page.value,
        pageSize: pageSize.value,
        ...processedSearchInfo,
        ...params
      })
      
      if (response.code === 0) {
        const list = response.data.list || []

        // 预处理数据，转换时间戳为日期时间对象
        // list.forEach(item => {
        //   item.log_time = timestampToDate(item.log_time)
        // })

        itemList.value = list
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
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
      // operation_type: '',
      month: getCurrentMonth(), // 重置为当前月份
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
    total.value = 0
    page.value = 1
    resetSearchInfo()
  }

  return {
    // 状态
    itemList,
    loading,
    total,
    page,
    pageSize,
    searchInfo,
    resourceTypes,
    resourceList,
    
    // 计算属性
    hasItems,
    totalPages,
    
    // 方法
    fetchResourceLogList,
    exportItems,
    fetchResourceTypes,
    fetchResourceList,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    clearState
  }
})
