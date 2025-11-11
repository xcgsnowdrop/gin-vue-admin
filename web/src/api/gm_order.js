import service from '@/utils/request'

// 获取游戏订单列表
export const getGMOrderList = (data) => {
  return service({
    url: '/gm/order/list',
    method: 'post',
    data: data
  })
}

// 获取充值商品名称列表
export const getGMRechargeList = () => {
  return service({
    url: '/gm/recharge/list',
    method: 'get',
  })
}
