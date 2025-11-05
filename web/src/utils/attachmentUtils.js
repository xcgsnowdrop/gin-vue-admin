/**
 * 附件工具函数
 */

/**
 * 过滤有效的附件（只保留 type、id、num 都有效的附件）
 */
export function filterValidAttachments(attachments) {
  if (!attachments || !Array.isArray(attachments)) {
    return []
  }
  return attachments.filter(
    att => att.type !== null && att.type !== undefined &&
           att.id !== null && att.id !== undefined &&
           att.num > 0
  )
}

/**
 * 初始化附件对象
 */
export function initAttachment() {
  return {
    id: null,
    type: null,
    num: 1
  }
}

