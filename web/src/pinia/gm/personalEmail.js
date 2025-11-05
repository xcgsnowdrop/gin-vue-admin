import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMPersonalEmailList,
  sendGMPersonalEmail,
} from '@/api/gm_email'
import { dateToTimestamp } from '@/utils/timestamp'
import { useEmailResource } from '@/composables/useEmailResource'

export const useGMPersonalEmailStore = defineStore('gmPersonalEmail', () => {
  // 状态
  const personalEmailList = ref([])

  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)

  // 初始化搜索信息结构
  const initSearchInfo = () => ({
    playerId: null, // 玩家ID筛选
    tplId: null, // 模板ID筛选
    startTime: null, // 个人邮件创建开始时间
    endTime: null, // 个人邮件创建结束时间
  })

  const searchInfo = ref(initSearchInfo())
  
  // 使用资源管理 composable
  const {
    resourceTypes,
    resourceList,
    resourceMap,
    fetchResourceTypes,
    fetchResourceList,
    loadResourcesForAttachments,
    loadAttachmentsFromEmailList,
    preloadAllResources,
    formatAttachment
  } = useEmailResource()

  // 计算属性
  const hasItems = computed(() => personalEmailList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 准备提交数据：转换时间戳格式
  const prepareSubmitData = (data) => {
    const submitData = { ...data }
    
    // 转换时间戳：将 Date 对象转换为时间戳（秒）
    submitData.startTime = dateToTimestamp(submitData.startTime)
    
    return submitData
  }

  // 获取个人邮件列表
  const fetchPersonalEmailList = async (params = {}) => {
    loading.value = true
    try {
      const response = await getGMPersonalEmailList({
        page: page.value,
        pageSize: pageSize.value,
        playerId: searchInfo.value.playerId,
        tplId: searchInfo.value.tplId ? parseInt(searchInfo.value.tplId) : null,
        startTime: searchInfo.value.startTime ? dateToTimestamp(searchInfo.value.startTime) : null,
        endTime: searchInfo.value.endTime ? dateToTimestamp(searchInfo.value.endTime) : null,
        ...params
      })
      
      if (response.code === 0) {
        const list = response.data.list || []

        // 批量加载附件所需的资源信息
        await loadAttachmentsFromEmailList(list, 'attachments')

        personalEmailList.value = list
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
      personalEmailList.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }


  // 发送个人邮件
  const sendPersonalEmail = async (data) => {
    try {
      const processedData = prepareSubmitData(data)
      const response = await sendGMPersonalEmail(processedData)
      if (response.code === 0) {
        await fetchPersonalEmailList()
        return true
      } else {
        throw new Error(response.msg || '发送个人邮件失败')
      }
    } catch (error) {
      console.error('发送个人邮件失败:', error)
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
    personalEmailList,
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
    fetchPersonalEmailList,
    sendPersonalEmail,
    fetchResourceTypes,
    fetchResourceList,
    preloadAllResources,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
  }
})
