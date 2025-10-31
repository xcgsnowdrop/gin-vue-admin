import service from '@/utils/request'

// 获取个人邮件列表
export const getGMPersonalEmailList = (data) => {
  return service({
    url: '/gm/email/personal/list',
    method: 'post',
    data
  })
}

// 发送个人邮件
export const sendGMPersonalEmail = (data) => {
  return service({
    url: '/gm/email/personal/send',
    method: 'post',
    data
  })
}

// 获取系统邮件列表
export const getGMSystemEmailList = (data) => {
  return service({
    url: '/gm/email/system/list',
    method: 'post',
    data
  })
}

// 发送系统邮件
export const sendGMSystemEmail = (data) => {
  return service({
    url: '/gm/email/system/send',
    method: 'post',
    data
  })
}

// 删除系统邮件
export const deleteGMSystemEmail = (data) => {
  return service({
    url: '/gm/email/system/delete',
    method: 'post',
    data
  })
}

// 更新系统邮件
export const updateGMSystemEmail = (data) => {
  return service({
    url: '/gm/email/system/update',
    method: 'post',
    data
  })
}
