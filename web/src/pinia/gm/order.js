import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getGMOrderList, getGMRechargeList } from '@/api/gm_order'
import { dateToTimestamp } from '@/utils/timestamp'
import { stringToIntArray } from '@/utils/stringFun'

export const useGMOrderStore = defineStore('gmOrder', () => {
  // 状态
  const orderList = ref([])
  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const rechargeList = ref([]) // 充值商品列表

  // 初始化搜索信息结构
  const initSearchInfo = () => ({
    player_id: '',
    area_ids: '',
    recharge_ids: '',
    pay_start_time: null,
    pay_end_time: null,
  })

  const searchInfo = ref(initSearchInfo())

  // 计算属性
  const hasOrders = computed(() => orderList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))
  // 充值商品映射表，用于根据 id 快速查找商品名称
  const rechargeMap = computed(() => {
    const map = {}
    rechargeList.value.forEach(item => {
      map[item.id] = item.name
    })
    return map
  })


  // 方法
  const setSearchInfo = (info) => {
    searchInfo.value = { ...searchInfo.value, ...info }
  }

  const resetSearchInfo = () => {
    searchInfo.value = initSearchInfo()
  }

  const setPage = (newPage) => {
    page.value = newPage
  }

  const setPageSize = (newPageSize) => {
    pageSize.value = newPageSize
    page.value = 1
  }

  const fetchOrderList = async (params = {}) => {
    loading.value = true
    try {
      const response = await getGMOrderList({
        page: page.value,
        pageSize: pageSize.value,
        player_id: searchInfo.value.player_id,
        area_ids: stringToIntArray(searchInfo.value.area_ids),
        recharge_ids: stringToIntArray(searchInfo.value.recharge_ids),
        pay_start_time: dateToTimestamp(searchInfo.value.pay_start_time),
        pay_end_time: dateToTimestamp(searchInfo.value.pay_end_time),
        ...params
      })

      if (response.code === 0) {
        orderList.value = response.data.list || []
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
      } else {
        throw new Error(response.msg || '获取订单列表失败')
      }
    } catch (error) {
      console.error('获取订单列表失败:', error)
      orderList.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  // 获取充值商品列表
  const fetchRechargeList = async () => {
    try {
      const response = await getGMRechargeList()
      if (response.code === 0) {
        rechargeList.value = response.data.list || []
      } else {
        throw new Error(response.msg || '获取充值商品列表失败')
      }
    } catch (error) {
      console.error('获取充值商品列表失败:', error)
      rechargeList.value = []
    }
  }

  return {
    // 状态
    orderList,
    loading,
    total,
    page,
    pageSize,
    searchInfo,
    rechargeList,

    // 计算属性
    hasOrders,
    totalPages,
    rechargeMap,
    
    // 方法
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    fetchOrderList,
    fetchRechargeList,
  }
})