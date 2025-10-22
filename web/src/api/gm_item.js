import gmService from '@/utils/gmRequest'

// 获取游戏道具列表
export const getGMItemList = (data) => {
  return gmService({
    url: '/gm/item/list',
    method: 'get',
    params: data
  })
}

// 创建游戏道具
export const createGMItem = (data) => {
  return gmService({
    url: '/gm/item',
    method: 'post',
    data
  })
}

// 更新游戏道具信息
export const updateGMItem = (data) => {
  return gmService({
    url: '/gm/item',
    method: 'put',
    data
  })
}

// 删除游戏道具
export const deleteGMItem = (data) => {
  return gmService({
    url: '/gm/item',
    method: 'delete',
    data
  })
}

// 获取游戏道具详情
export const getGMItem = (id) => {
  return gmService({
    url: `/gm/item/${id}`,
    method: 'get'
  })
}

// 批量操作游戏道具
export const batchOperateGMItem = (data) => {
  return gmService({
    url: '/gm/item/batch',
    method: 'post',
    data
  })
}

// 导出游戏道具数据
export const exportGMItem = (data) => {
  return gmService({
    url: '/gm/item/export',
    method: 'post',
    data,
    responseType: 'blob'
  })
}

// 获取游戏道具统计信息
export const getGMItemStats = () => {
  return gmService({
    url: '/gm/item/stats',
    method: 'get'
  })
}