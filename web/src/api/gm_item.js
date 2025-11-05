import service from '@/utils/request'

// 获取游戏道具列表
export const getGMResourceLogList = (data) => {
  return service({
    url: '/gm/resource/log/list',
    method: 'post',
    data: data
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

// 清理旧数据
export const cleanupGMItem = (data) => {
  return service({
    url: '/gm/item/cleanup',
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

// 批量获取资源列表
// 参数：data: 资源类型列表，格式为：{"resourceTypes":[1,2,3]}
// 返回：资源列表，格式为：{"code":0,"data":{"1":[{"id":1,"name":"resourceName1"},{"id":2,"name":"resourceName2"}],"2":[{"id":3,"name":"resourceName3"},{"id":4,"name":"resourceName4"}],"3":[{"id":5,"name":"resourceName5"},{"id":6,"name":"resourceName6"}]}}
export const getGMResourceListBatch = (data) => {
  return service({
    url: '/gm/resource/batch/list',
    method: 'post',
    data
  })
}