import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMSystemEmailAuditList,
  sendGMSystemEmailAudit,
  deleteGMSystemEmailAudit,
  updateGMSystemEmailAudit,
  reviewGMSystemEmailAudit,
} from '@/api/gm_email_audit'
import { dateToTimestamp } from '@/utils/timestamp'
import { useEmailResource } from '@/composables/useEmailResource'


export const useGMSystemEmailAuditStore = defineStore('gmSystemEmailAudit', () => {
  // 状态
  const systemEmailAuditList = ref([])

  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)

  const searchInfo = ref({
    applicantId: null, // 申请人ID
    auditorId: null, // 审核人ID
    startTime: null, // 申请开始时间
    endTime: null, // 申请结束时间
    status: null, // 状态筛选
  })
  
  // 使用资源管理 composable
  const {
    resourceTypes,
    resourceList,
    resourceMap,
    fetchResourceTypes,
    fetchResourceList,
    preloadAllResources,
    formatAttachment
  } = useEmailResource()

  // 计算属性
  const hasItems = computed(() => systemEmailAuditList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 准备提交数据：转换时间戳格式
  const prepareSubmitData = (data) => {
    const submitData = { ...data }
    
    // 转换时间戳：将 Date 对象转换为时间戳（秒）
    submitData.startTime = dateToTimestamp(submitData.startTime)
    submitData.maxRegTime = dateToTimestamp(submitData.maxRegTime)

    return submitData
  }

  // 获取系统邮件审核申请列表
  const fetchSystemEmailAuditList = async (params = {}) => {
    loading.value = true
    try {
      const response = await getGMSystemEmailAuditList({
        page: page.value,
        pageSize: pageSize.value,
        applicantId: searchInfo.value.applicantId ? parseInt(searchInfo.value.applicantId) : null,
        auditorId: searchInfo.value.auditorId ? parseInt(searchInfo.value.auditorId) : null,
        startTime: searchInfo.value.startTime ? dateToTimestamp(searchInfo.value.startTime) : null,
        endTime: searchInfo.value.endTime ? dateToTimestamp(searchInfo.value.endTime) : null,
        status: searchInfo.value.status ? parseInt(searchInfo.value.status) : null,
        ...params
      })
      
      if (response.code === 0) {
        const list = response.data.list || []

        systemEmailAuditList.value = list
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
        
        // console.log('Pinia store - itemList.value 已更新:', itemList.value)
        // console.log('Pinia store - 数据长度:', itemList.value.length)
      } else {
        throw new Error(response.msg || '获取系统邮件审核申请列表失败')
      }
    } catch (error) {
      console.error('获取系统邮件审核申请列表失败:', error)
      systemEmailAuditList.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }


  // 发送系统邮件审核申请
  const sendSystemEmailAudit = async (data) => {
    try {
      const processedData = prepareSubmitData(data)
      const response = await sendGMSystemEmailAudit(processedData)
      if (response.code === 0) {
        await fetchSystemEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '发送系统邮件审核申请失败')
      }
    } catch (error) {
      console.error('发送系统邮件审核申请失败:', error)
      throw error
    }
  }

  // 更新系统邮件审核申请
  const updateSystemEmailAudit = async (data) => {
    try {
      const processedData = prepareSubmitData(data)
      const response = await updateGMSystemEmailAudit(processedData)
      if (response.code === 0) {
        await fetchSystemEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '更新系统邮件审核申请失败')
      }
    } catch (error) {
      console.error('更新系统邮件审核申请失败:', error)
      throw error
    }
  }

  // 撤回系统邮件审核申请
  const deleteSystemEmailAudit = async (email_id) => {
    try {
      const response = await deleteGMSystemEmailAudit(email_id)
      if (response.code === 0) {
        await fetchSystemEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '撤回系统邮件审核申请失败')
      }
    } catch (error) {
      console.error('撤回系统邮件审核申请失败:', error)
      throw error
    }
  }

  // 审核系统邮件申请
  const reviewSystemEmailAudit = async (data) => {
    try {
      const response = await reviewGMSystemEmailAudit(data)
      if (response.code === 0) {
        await fetchSystemEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '审核系统邮件申请失败')
      }
    } catch (error) {
      console.error('审核系统邮件申请失败:', error)
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
      applicantId: '',
      auditorId: '',
      startTime: '',
      endTime: '',
      status: '',
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

  return {
    // 状态
    systemEmailAuditList,
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
    fetchSystemEmailAuditList,
    sendSystemEmailAudit,
    fetchResourceTypes,
    fetchResourceList,
    preloadAllResources,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    deleteSystemEmailAudit,
    updateSystemEmailAudit,
    reviewSystemEmailAudit,
  }
})
