import service from '@/utils/request'

// 获取游戏道具列表
export const getGMResourceLogList = (data) => {
  return service({
    url: '/gm/resource/log/list',
    method: 'post',
    data: data
  })
}

// 创建游戏道具
export const createGMItem = (data) => {
  return service({
    url: '/gm/item',
    method: 'post',
    data
  })
}

// 更新游戏道具信息
export const updateGMItem = (data) => {
  return service({
    url: '/gm/item',
    method: 'put',
    data
  })
}

// 删除游戏道具
export const deleteGMItem = (data) => {
  return service({
    url: '/gm/item',
    method: 'delete',
    data
  })
}

// 获取游戏道具详情
export const getGMItem = (id) => {
  return service({
    url: `/gm/item/${id}`,
    method: 'get'
  })
}

// 批量操作游戏道具
export const batchOperateGMItem = (data) => {
  return service({
    url: '/gm/item/batch',
    method: 'post',
    data
  })
}

// 导出游戏道具数据
export const exportGMItem = (data) => {
  return service({
    url: '/gm/item/export',
    method: 'post',
    data,
    responseType: 'blob'
  })
}

// 获取游戏道具统计信息
export const getGMItemStats = () => {
  return service({
    url: '/gm/item/stats',
    method: 'get'
  })
}

// 清理旧数据
export const cleanupGMItem = (data) => {
  return service({
    url: '/gm/item/cleanup',
    method: 'post',
    data
  })
}

// 批量删除道具流水记录
export const batchDeleteGMItem = (data) => {
  return service({
    url: '/gm/item/batchDelete',
    method: 'post',
    data
  })
}

// 获取操作类型列表
export const getGMItemOperationTypes = () => {
  return service({
    url: '/gm/item/operationTypes',
    method: 'get'
  })
}

// 获取道具类型列表
export const getGMItemTypes = () => {
  return service({
    url: '/gm/item/types',
    method: 'get'
  })
}

// 获取资源类型列表
export const getGMResourceTypeList = () => {
  return service({
    url: '/gm/resource/type/list',
    method: 'get'
  })
}

// 根据资源类型获取资源列表
export const getGMResourceList = (resType) => {
  return service({
    url: `/gm/resource/list/${resType}`,
    method: 'get'
  })
}