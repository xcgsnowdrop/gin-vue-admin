import service from '@/utils/request'

// 获取私人邮件审核申请列表
export const getGMPersonalEmailAuditList = (data) => {
  return service({
    url: '/email/audit/list',
    method: 'post',
    data: {
      ...data,
      isSystem: false // 只获取私人邮件
    }
  })
}

// 获取系统邮件审核申请列表
export const getGMSystemEmailAuditList = (data) => {
  return service({
    url: '/email/audit/list',
    method: 'post',
    data: {
      ...data,
      isSystem: true // 只获取系统邮件
    }
  })
}

// 创建邮件审核申请（系统/私人邮件通用）
export const sendGMSystemEmailAudit = (data) => {
  return service({
    url: '/email/audit/apply',
    method: 'post',
    data
  })
}

// 撤回邮件审核申请（系统/私人邮件通用）
export const deleteGMSystemEmailAudit = (id) => {
  return service({
    url: `/email/audit/${id}`,
    method: 'delete'
  })
}

// 更新邮件审核申请（系统/私人邮件通用）
export const updateGMSystemEmailAudit = (data) => {
  return service({
    url: '/email/audit/apply',
    method: 'put',
    data
  })
}

// 审核邮件申请（系统/私人邮件通用）
export const reviewGMSystemEmailAudit = (data) => {
  return service({
    url: '/email/audit/review',
    method: 'post',
    data
  })
}