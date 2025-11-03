import service from '@/utils/request'

// 获取个人邮件列表
export const getGMPersonalEmailAuditList = (data) => {
  return service({
    url: '/email/audit/personal/list',
    method: 'post',
    data
  })
}

// 发送个人邮件
export const sendGMPersonalEmailAudit = (data) => {
  return service({
    url: '/email/audit/personal/send',
    method: 'post',
    data
  })
}

// 获取系统邮件审核申请列表
export const getGMSystemEmailAuditList = (data) => {
  return service({
    url: '/email/audit/list',
    method: 'post',
    data
  })
}

// 发送系统邮件审核申请
export const sendGMSystemEmailAudit = (data) => {
  return service({
    url: '/email/audit/apply',
    method: 'post',
    data
  })
}

// 撤回系统邮件审核申请
export const deleteGMSystemEmailAudit = (id) => {
  return service({
    url: `/email/audit/${id}`,
    method: 'delete'
  })
}

// 更新系统邮件
export const updateGMSystemEmailAudit = (data) => {
  return service({
    url: '/email/audit/apply',
    method: 'put',
    data
  })
}

// 审核系统邮件申请
export const reviewGMSystemEmailAudit = (data) => {
  return service({
    url: '/email/audit/review',
    method: 'post',
    data
  })
}