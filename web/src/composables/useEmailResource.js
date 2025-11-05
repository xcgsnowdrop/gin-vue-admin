import { useResource } from './useResource'

/**
 * 邮件资源管理 Composable
 * 专门用于邮件相关的资源管理，扩展了 useResource 的功能
 */
export function useEmailResource() {
  const {
    resourceTypes,
    resourceList,
    resourceMap,
    attachmentResourceLists,
    fetchResourceTypes,
    fetchResourceList,
    preloadAllResources,
    getResourceListByType,
    getResourceTypeName,
    getResourceName,
    formatAttachment
  } = useResource()

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

