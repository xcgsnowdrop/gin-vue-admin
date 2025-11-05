import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMSystemEmailList,
  sendGMSystemEmail,
  deleteGMSystemEmail,
  updateGMSystemEmail,
} from '@/api/gm_email'
import { dateToTimestamp } from '@/utils/timestamp'
import { useResource } from '@/composables/useResource'


export const useGMSystemEmailStore = defineStore('gmSystemEmail', () => {
  // 状态
  const systemEmailList = ref([])

  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)

  // 初始化搜索信息结构
  const initSearchInfo = () => ({
    startTime: null, // 邮件创建开始时间
    endTime: null, // 邮件创建结束时间
  })

  const searchInfo = ref(initSearchInfo())
  
  // 使用资源管理 composable
  const {
    resourceTypes,
    resourceList,
    resourceMap,
    fetchResourceTypes,
    fetchResourceList,
    preloadAllResources,
    formatAttachment
  } = useResource()

  // 计算属性
  const hasItems = computed(() => systemEmailList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 准备提交数据：转换时间戳格式
  const prepareSubmitData = (data) => {
    const submitData = { ...data }
    
    // 转换时间戳：将 Date 对象转换为时间戳（秒）
    submitData.startTime = dateToTimestamp(submitData.startTime)
    submitData.maxRegTime = dateToTimestamp(submitData.maxRegTime)
    if (submitData.areaIds) {
      submitData.areaIds = submitData.areaIds.split(',').map(id => parseInt(id))
    }

    return submitData
  }

  // 获取个人邮件列表
  const fetchSystemEmailList = async (params = {}) => {
    loading.value = true
    try {
      // 构建查询参数，处理时间字段
      const queryParams = {
        page: page.value,
        pageSize: pageSize.value,
        ...params
      }
      
      // 处理时间字段（转换为时间戳）
      if (searchInfo.value.startTime) {
        queryParams.startTime = dateToTimestamp(searchInfo.value.startTime)
      }
      
      if (searchInfo.value.endTime) {
        queryParams.endTime = dateToTimestamp(searchInfo.value.endTime)
      }

      const response = await getGMSystemEmailList(queryParams)
      
      if (response.code === 0) {
        const list = response.data.list || []

        systemEmailList.value = list
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
      } else {
        throw new Error(response.msg || '获取道具流水列表失败')
      }
    } catch (error) {
      console.error('获取道具流水列表失败:', error)
      systemEmailList.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }


  // 发送个人邮件
  const sendSystemEmail = async (data) => {
    try {
      const processedData = prepareSubmitData(data)
      const response = await sendGMSystemEmail(processedData)
      if (response.code === 0) {
        await fetchSystemEmailList()
        return true
      } else {
        throw new Error(response.msg || '发送个人邮件失败')
      }
    } catch (error) {
      console.error('发送个人邮件失败:', error)
      throw error
    }
  }

  // 更新系统邮件
  const updateSystemEmail = async (data) => {
    try {
      const processedData = prepareSubmitData(data)
      const response = await updateGMSystemEmail(processedData)
      if (response.code === 0) {
        await fetchSystemEmailList()
        return true
      } else {
        throw new Error(response.msg || '更新系统邮件失败')
      }
    } catch (error) {
      console.error('更新系统邮件失败:', error)
      throw error
    }
  }

  // 删除系统邮件
  const deleteSystemEmail = async (email_id) => {
    try {
      const response = await deleteGMSystemEmail(email_id)
      if (response.code === 0) {
        await fetchSystemEmailList()
        return true
      }
    } catch (error) {
      console.error('删除系统邮件失败:', error)
      throw error
    }
  }

  // 设置搜索条件
  const setSearchInfo = (info) => {
    searchInfo.value = { ...searchInfo.value, ...info }
  }

  // 重置搜索条件
  const resetSearchInfo = () => {
    searchInfo.value = initSearchInfo()
  }

  // 设置分页
  const setPage = (newPage) => {
    page.value = newPage
  }

  const setPageSize = (newPageSize) => {
    pageSize.value = newPageSize
    page.value = 1
  }

  return {
    // 状态
    systemEmailList,
    loading,
    total,
    page,
    pageSize,
    searchInfo,
    resourceTypes,
    resourceList,
    resourceMap,
    formatAttachment,
    
    // 计算属性
    hasItems,
    totalPages,
    
    // 方法
    fetchSystemEmailList,
    sendSystemEmail,
    fetchResourceTypes,
    fetchResourceList,
    preloadAllResources,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    deleteSystemEmail,
    updateSystemEmail,
  }
})
