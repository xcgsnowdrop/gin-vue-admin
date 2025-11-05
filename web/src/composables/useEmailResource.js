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
    loadResourcesForAttachments,
    getResourceListByType,
    getResourceTypeName,
    getResourceName,
    formatAttachment
  } = useResource()

  /**
   * 批量加载邮件列表中的附件资源
   * @param {Array} emailList - 邮件列表
   * @param {string} attachmentField - 附件字段名，默认为 'emailAttachments' 或 'attachments'
   */
  const loadAttachmentsFromEmailList = async (emailList, attachmentField = 'emailAttachments') => {
    if (!emailList || emailList.length === 0) return

    // 收集所有附件中的资源类型
    const allAttachments = []
    emailList.forEach(item => {
      const attachments = item[attachmentField] || item.attachments
      if (attachments && Array.isArray(attachments)) {
        allAttachments.push(...attachments)
      }
    })

    // 批量加载附件所需的资源信息
    if (allAttachments.length > 0) {
      await loadResourcesForAttachments(allAttachments)
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
    loadResourcesForAttachments,
    loadAttachmentsFromEmailList,
    getResourceListByType,
    getResourceTypeName,
    getResourceName,
    formatAttachment
  }
}

