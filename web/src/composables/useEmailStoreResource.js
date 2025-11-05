import { useEmailResource } from './useEmailResource'

/**
 * 邮件 Store 资源管理 Composable
 * 用于在 Pinia store 中集成资源管理功能
 */
export function useEmailStoreResource() {
  const {
    resourceTypes,
    resourceList,
    resourceMap,
    fetchResourceTypes,
    fetchResourceList,
    getResourceTypeName,
    getResourceName,
    formatAttachment
  } = useEmailResource()

  return {
    // 状态
    resourceTypes,
    resourceList,
    resourceMap,
    
    // 方法
    fetchResourceTypes,
    fetchResourceList,
    getResourceTypeName,
    getResourceName,
    formatAttachment
  }
}

