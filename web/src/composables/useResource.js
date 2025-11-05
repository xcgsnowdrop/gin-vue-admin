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

  /**
   * 预加载所有资源类型的资源列表
   * 优化：一次性加载所有资源类型，避免后续多次请求
   */
  const preloadAllResources = async () => {
    try {
      // 先确保资源类型列表已加载
      if (resourceTypes.value.length === 0) {
        await fetchResourceTypes()
      }

      // 如果没有资源类型，直接返回
      if (resourceTypes.value.length === 0) {
        return
      }

      // 收集所有资源类型的 type 值
      const allResourceTypes = resourceTypes.value.map(rt => rt.type).filter(type => type != null)

      // 如果没有资源类型，直接返回
      if (allResourceTypes.length === 0) {
        return
      }

      // 使用批量接口一次性加载所有资源类型
      const response = await getGMResourceListBatch({ resourceTypes: allResourceTypes })

      if (response.code === 0 && response.data) {
        const batchData = response.data

        // 遍历返回的数据，更新 resourceMap
        Object.keys(batchData).forEach(typeStr => {
          const type = parseInt(typeStr)
          const typeResourceList = batchData[typeStr] || []

          // 更新 resourceMap
          if (!resourceMap.value[type]) {
            resourceMap.value[type] = {}
          }
          typeResourceList.forEach(item => {
            if (item && item.id !== undefined && item.name) {
              resourceMap.value[type][item.id] = item.name
            }
          })
        })
      } else {
        throw new Error(response.msg || '批量获取资源列表失败')
      }
    } catch (error) {
      console.error('预加载所有资源失败:', error)
      throw error
    }
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
    preloadAllResources,
    getResourceListByType,
    getResourceTypeName,
    getResourceName,
    formatAttachment
  }
}

