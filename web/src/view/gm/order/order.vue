<template>
    <div>
      <warning-bar title="注：GM管理 - 订单查询" />
      <div class="gva-search-box">
        <el-form ref="searchForm" :inline="true" :model="searchInfo">
          <el-form-item label="PlayerId">
            <el-input v-model="searchInfo.player_id" placeholder="PlayerId" />
          </el-form-item>
          <el-form-item label="AreaId">
            <el-input v-model="searchInfo.area_id" placeholder="AreaId" />
          </el-form-item>
          <el-form-item label="RechargeId">
            <el-input v-model="searchInfo.recharge_id" placeholder="RechargeId" />
          </el-form-item>
          
          <el-form-item label="下单开始时间">
            <el-date-picker
              v-model="searchInfo.startTime"
              type="datetime"
              placeholder="请选择邮件创建开始时间"
              style="width: 100%"
              :disabled-date="
                (time) =>
                  searchInfo.endTime
                    ? time.getTime() > searchInfo.endTime.getTime()
                    : false
              "
            />
          </el-form-item>
          <el-form-item label="下单结束时间">
            <el-date-picker
              v-model="searchInfo.endTime"
              type="datetime"
              placeholder="请选择邮件创建结束时间"
              style="width: 100%"
              :disabled-date="
                (time) =>
                  searchInfo.startTime
                    ? time.getTime() < searchInfo.startTime.getTime()
                    : false
              "
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="search" @click="onSubmit">
              查询
            </el-button>
            <el-button icon="refresh" @click="onReset"> 重置 </el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="gva-table-box">
        <!-- <div class="gva-btn-list">
          <el-button type="success" icon="download" @click="exportData">
            导出数据
          </el-button>
        </div> -->
        <el-table :data="tableData" row-key="_id" v-loading="loading">
            <el-table-column align="left" label="订单号" min-width="180" prop="_id" />
          <el-table-column align="left" label="平台订单号" min-width="180" prop="platform_order_id" />
          <el-table-column
            align="left"
            label="PlayerId"
            min-width="120"
            prop="player_id"
          />
          <el-table-column
            align="left"
            label="区服"
            min-width="100"
            prop="area_id"
          />
          <el-table-column
            align="left"
            label="商品ID"
            min-width="100"
            prop="recharge_id"
          />
          <el-table-column
            align="left"
            label="商品名称"
            min-width="150"
            prop="recharge_name"
          >
            <template #default="scope">
              {{ getRechargeName(scope.row.recharge_id) }}
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="价格"
            min-width="150"
            prop="price"
          />
          <el-table-column
            align="left"
            label="币种"
            min-width="150"
            prop="currency"
          />
          <el-table-column
            align="left"
            label="订单状态"
            min-width="100"
            prop="status"
          >
            <template #default="scope">
              <el-tag
                :type="getOrderStatusTag(scope.row.status)"
                size="small"
              >
                {{ getOrderStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column
            align="left"
            label="下单时间"
            min-width="180"
            prop="pay_time"
          >
            <template #default="scope">
              {{ formatTimestamp(scope.row.pay_time) }}
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>
  
    </div>
  </template>
  
  <script setup>
  import { onMounted, watch } from 'vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { useGMOrderStore } from '@/pinia/gm/order'
  import { storeToRefs } from 'pinia'
  import { formatTimestamp } from '@/utils/timestamp'
  
  defineOptions({
    name: 'GmOrder'
  })
  
  const gmOrderStore = useGMOrderStore()
  
  // 使用store中的状态
  const { 
    orderList: tableData, 
    loading, 
    total, 
    page, 
    pageSize, 
    searchInfo,
    rechargeMap
  } = storeToRefs(gmOrderStore)
  
  const {
    fetchOrderList,
    resetSearchInfo,
    setPage,
    setPageSize,
    fetchRechargeList
  } = gmOrderStore
  
  // 获取操作类型标签样式
  const getOrderStatusTag = (status) => {
    const tagMap = {
      3: 'success', // 已发货
      2: 'warning', // 待发货
      1: 'info' // 订单创建
    }
    return tagMap[status] || 'info'
  }
  
  // 获取操作类型文本
  const getOrderStatusText = (status) => {
    const textMap = {
      3: '已发货',
      2: '待发货',
      1: '订单创建'
    }
    return textMap[status] || status
  }
  
  // 根据充值商品ID获取商品名称
  const getRechargeName = (rechargeId) => {
    return rechargeMap.value[rechargeId] || '-'
  }
  
  // 查询数据
  const onSubmit = () => {
    setPage(1)
    fetchOrderList()
  }
  
  // 重置搜索
  const onReset = () => {
    resetSearchInfo()
    setPage(1)
    fetchOrderList()
  }
  
  // 分页处理
  const handleCurrentChange = (val) => {
    setPage(val)
    fetchOrderList()
  }
  
  const handleSizeChange = (val) => {
    setPageSize(val)
    fetchOrderList()
  }
  
  // // 导出数据
  // const exportData = async () => {
  //   try {
  //     await exportItems()
  //     ElMessage.success('导出成功')
  //   } catch (error) {
  //     ElMessage.error(error.message || '导出失败')
  //   }
  // }
  
  // 监听月份变化，清空时间范围选择
  watch(
    () => searchInfo.value.month,
    (newMonth, oldMonth) => {
      if (newMonth !== oldMonth && searchInfo.value.log_time_range && searchInfo.value.log_time_range.length > 0) {
        // 检查当前选择的时间是否还在新月份范围内
        const selectedMonth = newMonth
        const year = parseInt(selectedMonth.substring(0, 4))
        const month = parseInt(selectedMonth.substring(4, 6)) - 1
        
        const startTime = new Date(searchInfo.value.log_time_range[0])
        const endTime = new Date(searchInfo.value.log_time_range[1])
        
        const startYear = startTime.getFullYear()
        const startMonth = startTime.getMonth()
        const endYear = endTime.getFullYear()
        const endMonth = endTime.getMonth()
        
        // 如果时间范围不在新月份内，清空时间范围
        if (startYear !== year || startMonth !== month || endYear !== year || endMonth !== month) {
          searchInfo.value.log_time_range = []
        }
      }
    }
  )
  
  watch(
      () => tableData.value,
      (newValue, oldValue) => {
        console.log('tableData 变化了')
        console.log('新值:', newValue)
        console.log('旧值:', oldValue)
        console.log('新值长度:', newValue?.length)
      },
      { deep: true, immediate: true }
  )
  
  
  // 初始化
  onMounted(async () => {
    try {
      // 并行获取充值商品列表和订单列表
      await Promise.all([
        fetchRechargeList(),
        fetchOrderList()
      ])
    } catch (error) {
      console.error('初始化失败:', error)
    }
  })
  </script>
  
  <style scoped>
  .text-green-600 {
    color: #16a34a;
  }
  
  .text-red-600 {
    color: #dc2626;
  }
  </style>
  