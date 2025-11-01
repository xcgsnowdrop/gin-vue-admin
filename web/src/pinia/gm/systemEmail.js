import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getGMSystemEmailList,
  sendGMSystemEmail,
  deleteGMSystemEmail,
  updateGMSystemEmail,
} from '@/api/gm_email'
import { getGMResourceTypeList, getGMResourceList } from '@/api/gm_item'


export const useGMSystemEmailStore = defineStore('gmSystemEmail', () => {
  // 状态
  const systemEmailList = ref([])

  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)

  const searchInfo = ref({
    player_id: '',
  })
  
  const resourceTypes = ref([])  // 资源类型列表
  const resourceList = ref([])    // 资源列表（根据类型动态获取）
  const resourceMap = ref({})      // 资源映射 { type: { id: name } }，用于快速查找资源名称

  // 计算属性
  const hasItems = computed(() => systemEmailList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  // 准备提交数据：转换时间戳格式
  const prepareSubmitData = (data) => {
    const submitData = { ...data }
    
    // 转换时间戳：将 Date 对象转换为时间戳（秒）
    if (submitData.startTime instanceof Date) {
      submitData.startTime = Math.floor(submitData.startTime.getTime() / 1000)
    } else if (submitData.startTime) {
      const date = new Date(submitData.startTime)
      if (!isNaN(date.getTime())) {
        submitData.startTime = Math.floor(date.getTime() / 1000)
      }
    }

    if (submitData.maxRegTime instanceof Date) {
      submitData.maxRegTime = Math.floor(submitData.maxRegTime.getTime() / 1000)
    } else if (submitData.maxRegTime) {
      const date = new Date(submitData.maxRegTime)
      if (!isNaN(date.getTime())) {
        submitData.maxRegTime = Math.floor(date.getTime() / 1000)
      }
    }

    if (submitData.areaIds) {
      submitData.areaIds = submitData.areaIds.split(',').map(id => parseInt(id))
    }

    return submitData
  }

  // 获取个人邮件列表
  const fetchSystemEmailList = async (params = {}) => {
    loading.value = true
    try {
      const response = await getGMSystemEmailList({
        page: page.value,
        pageSize: pageSize.value,
        ...searchInfo.value,
        ...params
      })
      
      if (response.code === 0) {
        const list = response.data.list || []

        // 收集所有附件中的资源类型
        const allAttachments = []
        list.forEach(item => {
          if (item.attachments && Array.isArray(item.attachments)) {
            allAttachments.push(...item.attachments)
          }
        })

        // 批量加载附件所需的资源信息
        if (allAttachments.length > 0) {
          await loadResourcesForAttachments(allAttachments)
        }

        // 预处理数据，转换时间戳为日期时间对象
        list.forEach(item => {
          item.create_time_formatted = item.create_time ? new Date(item.create_time * 1000).toLocaleString() : '-'
          item.start_time_formatted = item.start_time ? new Date(item.start_time * 1000).toLocaleString() : '-'
          item.end_time_formatted = item.end_time ? new Date(item.end_time * 1000).toLocaleString() : '-'
          item.max_reg_time_formatted = item.max_reg_time ? new Date(item.max_reg_time * 1000).toLocaleString() : '-'
        })

        systemEmailList.value = list
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
      systemEmailList.value = []
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
    searchInfo.value = {
      player_id: '',
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
    systemEmailList,
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
    fetchSystemEmailList,
    sendSystemEmail,
    fetchResourceTypes,
    fetchResourceList,
    loadResourcesForAttachments,
    getResourceTypeName,
    getResourceName,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    deleteSystemEmail,
    updateSystemEmail,
  }
})
