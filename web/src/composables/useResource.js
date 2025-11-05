import { ref } from 'vue'
import { getGMResourceTypeList, getGMResourceList, getGMResourceListBatch } from '@/api/gm_item'

/**
 * 资源管理 Composable
 * 提供资源类型、资源列表的获取和附件格式化等功能
 */
export function useResource() {
  const resourceTypes = ref([])  // 资源类型列表
  const resourceList = ref([])    // 资源列表（根据类型动态获取）
  const resourceMap = ref({})      // 资源映射 { type: { id: name } }，用于快速查找资源名称

  // 附件资源列表缓存（每个类型对应一个资源列表）
  const attachmentResourceLists = ref({})

  /**
   * 获取资源类型列表
   */
  const fetchResourceTypes = async () => {
    try {
      const response = await getGMResourceTypeList()
      if (response.code === 0) {
        resourceTypes.value = response.data.list || response.data || []
      } else {
        throw new Error(response.msg || '获取资源类型列表失败')
      }
    } catch (error) {
      console.error('获取资源类型列表失败:', error)
      throw error
    }
  }

  /**
   * 获取指定类型的资源列表
   */
  const fetchResourceList = async (resType) => {
    try {
      const response = await getGMResourceList(resType)
      if (response.code === 0) {
        const list = response.data.list || response.data || []
        resourceList.value = list
        
        // 更新 resourceMap
        if (!resourceMap.value[resType]) {
          resourceMap.value[resType] = {}
        }
        list.forEach(item => {
          resourceMap.value[resType][item.id] = item.name
        })
        
        return list
      } else {
        throw new Error(response.msg || '获取资源列表失败')
      }
    } catch (error) {
      console.error(`获取资源类型 ${resType} 的列表失败:`, error)
      throw error
    }
  }

  /**
   * 批量加载附件所需的资源信息
   * 优化：使用批量接口，避免多次 HTTP 请求
   */
  const loadResourcesForAttachments = async (attachments) => {
    if (!attachments || attachments.length === 0) return

    // 收集所有唯一的资源类型
    const types = new Set()
    attachments.forEach(att => {
      if (att && att.type) {
        types.add(att.type)
      }
    })

    if (types.size === 0) return

    // 多个资源类型时，使用批量接口
    try {
      const resourceTypesArray = Array.from(types)
      const response = await getGMResourceListBatch({ resourceTypes: resourceTypesArray })
      
      if (response.code === 0 && response.data) {
        // 处理批量返回的数据，更新 resourceMap 和 resourceList
        const batchData = response.data
        
        // 遍历返回的数据，更新 resourceMap
        Object.keys(batchData).forEach(typeStr => {
          const type = parseInt(typeStr)
          const resourceList = batchData[typeStr] || []
          
          // 更新 resourceMap
          if (!resourceMap.value[type]) {
            resourceMap.value[type] = {}
          }
          resourceList.forEach(item => {
            if (item && item.id !== undefined && item.name) {
              resourceMap.value[type][item.id] = item.name
            }
          })
          
          // 更新缓存（如果当前 resourceList 正好是这个类型）
          // 注意：批量接口返回多个类型，我们只更新当前 resourceList 对应的类型（如果有的话）
          // 实际上批量加载时，我们主要更新 resourceMap，resourceList 会在需要时从 resourceMap 构建
        })
      } else {
        throw new Error(response.msg || '批量获取资源列表失败')
      }
    } catch (error) {
      console.error('批量加载资源列表失败:', error)
      throw error
      // 如果批量请求失败，回退到原来的并行单个请求方式
      // const promises = Array.from(types).map(type => fetchResourceList(type))
      // await Promise.all(promises)
    }
  }

  /**
   * 根据资源类型获取资源列表（从缓存或 resourceMap）
   */
  const getResourceListByType = (type, attachmentResourceLists = {}) => {
    if (!type) return []
    // 优先从缓存中获取
    if (attachmentResourceLists[type]) {
      return attachmentResourceLists[type]
    }
    // 如果缓存中没有，从 resourceMap 中构建
    if (resourceMap.value[type]) {
      return Object.keys(resourceMap.value[type]).map(id => ({
        id: parseInt(id),
        name: resourceMap.value[type][id]
      }))
    }
    return []
  }

  /**
   * 获取资源类型名称
   */
  const getResourceTypeName = (type) => {
    if (!type) return '未知类型'
    const resType = resourceTypes.value.find(t => t.type === type)
    return resType ? resType.name : `类型${type}`
  }

  /**
   * 获取资源名称
   */
  const getResourceName = (type, id) => {
    if (!type || !id) return '未知资源'
    if (resourceMap.value[type] && resourceMap.value[type][id]) {
      return resourceMap.value[type][id]
    }
    return `资源${id}`
  }

  /**
   * 格式化附件显示
   */
  const formatAttachment = (attachment) => {
    if (!attachment) return ''
    const typeName = getResourceTypeName(attachment.type)
    const resourceName = getResourceName(attachment.type, attachment.id)
    return `${typeName} - ${resourceName} × ${attachment.num}`
  }

  return {
    // 状态
    resourceTypes,
    resourceList,
    resourceMap,
    attachmentResourceLists,
    
    // 方法
    fetchResourceTypes,
    fetchResourceList,
    loadResourcesForAttachments,
    getResourceListByType,
    getResourceTypeName,
    getResourceName,
    formatAttachment
  }
}

