import service from '@/utils/request'

// 获取游戏用户列表
export const getGMUserList = (data) => {
  return service({
    url: '/gm/user/list',
    method: 'post',
    data
  })
}

// 获取游戏用户详情
export const getGMUser = (id) => {
  return service({
    url: `/gm/user/detail/${id}`,
    method: 'get'
  })
}

// 切换禁言状态
export const toggleGMBanChat = (data) => {
  return service({
    url: '/gm/user/toggleBanChat',
    method: 'post',
    data
  })
}

// 切换封号状态
export const toggleGMBanLogin = (data) => {
  return service({
    url: '/gm/user/toggleBanLogin',
    method: 'post',
    data
  })
}

// 批量操作游戏用户
export const batchOperateGMUser = (data) => {
  return service({
    url: '/gm/user/batch',
    method: 'post',
    data
  })
}

// 导出游戏用户数据
export const exportGMUser = (data) => {
  return service({
    url: '/gm/user/export',
    method: 'post',
    data,
    responseType: 'blob'
  })
}
