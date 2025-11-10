import service from '@/utils/request'

// 获取游戏订单列表
export const getGMOrderList = (data) => {
  return service({
    url: '/gm/order/list',
    method: 'post',
    data: data
  })
}
