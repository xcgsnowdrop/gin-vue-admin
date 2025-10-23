import gmService from '@/utils/gmRequest'

// 获取游戏用户列表
export const getGMUserList = (data) => {
  return gmService({
    url: '/gm/user/list',
    method: 'post',
    data
  })
}

// 创建游戏用户
export const createGMUser = (data) => {
  return gmService({
    url: '/gm/user',
    method: 'post',
    data
  })
}

// 更新游戏用户信息
export const updateGMUser = (data) => {
  return gmService({
    url: '/gm/user',
    method: 'put',
    data
  })
}

// 删除游戏用户
export const deleteGMUser = (data) => {
  return gmService({
    url: '/gm/user',
    method: 'delete',
    data
  })
}

// 获取游戏用户详情
export const getGMUser = (id) => {
  return gmService({
    url: `/gm/user/${id}`,
    method: 'get'
  })
}

// 重置游戏用户密码
export const resetGMUserPassword = (data) => {
  return gmService({
    url: '/gm/user/resetPassword',
    method: 'post',
    data
  })
}

// 设置游戏用户状态
export const setGMUserStatus = (data) => {
  return gmService({
    url: '/gm/user/status',
    method: 'post',
    data
  })
}

// 切换禁言状态
export const toggleGMBanChat = (data) => {
  return gmService({
    url: '/gm/user/toggleBanChat',
    method: 'post',
    data
  })
}

// 切换封号状态
export const toggleGMBanLogin = (data) => {
  return gmService({
    url: '/gm/user/toggleBanLogin',
    method: 'post',
    data
  })
}

// 批量操作游戏用户
export const batchOperateGMUser = (data) => {
  return gmService({
    url: '/gm/user/batch',
    method: 'post',
    data
  })
}

// 导出游戏用户数据
export const exportGMUser = (data) => {
  return gmService({
    url: '/gm/user/export',
    method: 'post',
    data,
    responseType: 'blob'
  })
}

// 获取游戏用户统计信息
export const getGMUserStats = () => {
  return gmService({
    url: '/gm/user/stats',
    method: 'get'
  })
}
