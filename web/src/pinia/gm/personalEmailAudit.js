import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMPersonalEmailAuditList,
  sendGMSystemEmailAudit,
  deleteGMSystemEmailAudit, // 使用相同的删除接口
  updateGMSystemEmailAudit,
  reviewGMSystemEmailAudit, // 使用相同的审核接口
} from '@/api/gm_email_audit'
import { getGMResourceTypeList, getGMResourceList } from '@/api/gm_item'
import { dateToTimestamp } from '@/utils/timestamp'


export const useGMPersonalEmailAuditStore = defineStore('gmPersonalEmailAudit', () => {
  // 状态
  const personalEmailAuditList = ref([])

  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)

  const searchInfo = ref({
    applicantId: null, // 申请人ID
    auditorId: null, // 审核人ID
    playerId: '', // 玩家ID
    startTime: null, // 申请开始时间
    endTime: null, // 申请结束时间
    status: null, // 状态筛选
  })
  
  const resourceTypes = ref([])  // 资源类型列表
  const resourceList = ref([])    // 资源列表（根据类型动态获取）
  const resourceMap = ref({})      // 资源映射 { type: { id: name } }，用于快速查找资源名称

  // 计算属性
  const hasItems = computed(() => personalEmailAuditList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 准备提交数据：转换时间戳格式
  const prepareSubmitData = (data) => {
    const submitData = { ...data }
    
    // 转换时间戳：将 Date 对象转换为时间戳（秒）
    submitData.startTime = dateToTimestamp(submitData.startTime)
    // 私人邮件不需要maxRegTime，需要playerId

    return submitData
  }

  // 获取私人邮件审核申请列表
  const fetchPersonalEmailAuditList = async (params = {}) => {
    loading.value = true
    try {
      const response = await getGMPersonalEmailAuditList({
        page: page.value,
        pageSize: pageSize.value,
        applicantId: searchInfo.value.applicantId ? parseInt(searchInfo.value.applicantId) : null,
        auditorId: searchInfo.value.auditorId ? parseInt(searchInfo.value.auditorId) : null,
        playerId: searchInfo.value.playerId || '',
        startTime: searchInfo.value.startTime ? dateToTimestamp(searchInfo.value.startTime) : null,
        endTime: searchInfo.value.endTime ? dateToTimestamp(searchInfo.value.endTime) : null,
        status: searchInfo.value.status ? parseInt(searchInfo.value.status) : null,
        ...params
      })
      
      if (response.code === 0) {
        const list = response.data.list || []

        // 收集所有附件中的资源类型
        const allAttachments = []
        list.forEach(item => {
          if (item.emailAttachments && Array.isArray(item.emailAttachments)) {
            allAttachments.push(...item.emailAttachments)
          }
        })

        // 批量加载附件所需的资源信息
        if (allAttachments.length > 0) {
          await loadResourcesForAttachments(allAttachments)
        }

        personalEmailAuditList.value = list
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
      } else {
        throw new Error(response.msg || '获取私人邮件审核申请列表失败')
      }
    } catch (error) {
      console.error('获取私人邮件审核申请列表失败:', error)
      personalEmailAuditList.value = []
      total.value = 0
    } finally {
      loading.value = false
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
        const list = response.data.list || []
        resourceList.value = list
        
        // 更新资源映射，方便快速查找
        if (!resourceMap.value[resType]) {
          resourceMap.value[resType] = {}
        }
        list.forEach(resource => {
          resourceMap.value[resType][resource.id] = resource.name
        })
        
        return list
      } else {
        throw new Error(response.msg || '获取资源列表失败')
      }
    } catch (error) {
      console.error('获取资源列表失败:', error)
      resourceList.value = []
      throw error
    }
  }

  // 批量加载资源列表（用于附件展示）
  const loadResourcesForAttachments = async (attachments) => {
    if (!attachments || !Array.isArray(attachments) || attachments.length === 0) {
      return
    }
    
    // 收集所有需要的资源类型
    const typesNeeded = [...new Set(attachments.map(att => att.type))]
    
    // 并行加载所有需要的资源类型
    const loadPromises = typesNeeded.map(type => {
      // 如果已经加载过，跳过
      if (resourceMap.value[type]) {
        return Promise.resolve()
      }
      return fetchResourceList(type).catch(err => {
        console.error(`加载资源类型 ${type} 失败:`, err)
      })
    })
    
    await Promise.all(loadPromises)
  }

  // 获取资源类型名称
  const getResourceTypeName = (type) => {
    const resourceType = resourceTypes.value.find(rt => rt.type === type)
    return resourceType ? resourceType.name : `类型${type}`
  }

  // 获取资源名称
  const getResourceName = (type, id) => {
    if (resourceMap.value[type] && resourceMap.value[type][id]) {
      return resourceMap.value[type][id]
    }
    return `资源${id}`
  }

  // 发送私人邮件审核申请
  const sendPersonalEmailAudit = async (data) => {
    try {
      const processedData = prepareSubmitData(data)
      const response = await sendGMSystemEmailAudit(processedData)
      if (response.code === 0) {
        await fetchPersonalEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '发送私人邮件审核申请失败')
      }
    } catch (error) {
      console.error('发送私人邮件审核申请失败:', error)
      throw error
    }
  }

  // 更新私人邮件审核申请
  const updatePersonalEmailAudit = async (data) => {
    try {
      const processedData = prepareSubmitData(data)
      const response = await updateGMSystemEmailAudit(processedData)
      if (response.code === 0) {
        await fetchPersonalEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '更新私人邮件审核申请失败')
      }
    } catch (error) {
      console.error('更新私人邮件审核申请失败:', error)
      throw error
    }
  }

  // 撤回私人邮件审核申请
  const deletePersonalEmailAudit = async (email_id) => {
    try {
      const response = await deleteGMSystemEmailAudit(email_id)
      if (response.code === 0) {
        await fetchPersonalEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '撤回私人邮件审核申请失败')
      }
    } catch (error) {
      console.error('撤回私人邮件审核申请失败:', error)
      throw error
    }
  }

  // 审核私人邮件申请
  const reviewPersonalEmailAudit = async (data) => {
    try {
      const response = await reviewGMSystemEmailAudit(data)
      if (response.code === 0) {
        await fetchPersonalEmailAuditList()
        return true
      } else {
        throw new Error(response.msg || '审核私人邮件申请失败')
      }
    } catch (error) {
      console.error('审核私人邮件申请失败:', error)
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
      applicantId: null,
      auditorId: null,
      playerId: '',
      startTime: null,
      endTime: null,
      status: null,
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
    personalEmailAuditList,
    loading,
    total,
    page,
    pageSize,
    searchInfo,
    resourceTypes,
    resourceList,
    resourceMap,
    
    // 计算属性
    hasItems,
    totalPages,
    
    // 方法
    fetchPersonalEmailAuditList,
    sendPersonalEmailAudit,
    fetchResourceTypes,
    fetchResourceList,
    loadResourcesForAttachments,
    getResourceTypeName,
    getResourceName,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    deletePersonalEmailAudit,
    updatePersonalEmailAudit,
    reviewPersonalEmailAudit,
  }
})

