import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { 
    getGMAnnouncementList, 
    addGMAnnouncement, 
    deleteGMAnnouncement, 
    updateGMAnnouncement, 
    toppingGMAnnouncement 
} from '@/api/gm_announcement'
import { timestampToDate, convertDatesToTimestamps } from '@/utils/timestamp'


export const useGMAnnouncementStore = defineStore('gmAnnouncement', () => {
  // 状态
  const announcementList = ref([])
  const currentAnnouncement = ref(null)
  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)

  const searchInfo = ref({
    type: '',
    startCreatedAt: null,
    endCreatedAt: null
  })

  // 计算属性
  const hasAnnouncements = computed(() => announcementList.value.length > 0)
  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))


  // 获取公告列表
  const fetchAnnouncementList = async (params = {}) => {
    loading.value = true
    try {
      const processedSearchInfo = convertDatesToTimestamps(searchInfo.value, ['startCreatedAt', 'endCreatedAt'])
      const response = await getGMAnnouncementList({
        page: page.value,
        pageSize: pageSize.value,
        ...processedSearchInfo,
        ...params
      })

      if (response.code === 0) {
        const list = response.data.announcementList || []

        // 预处理数据，添加格式化时间字段（使用工具函数）
        list.forEach(item => {
          item.startTime = timestampToDate(item.startTime)
          item.endTime = timestampToDate(item.endTime)
          item.createTime = timestampToDate(item.createTime)
        })

        announcementList.value = list
        total.value = response.data.total || 0
        page.value = response.data.page || 1
        pageSize.value = response.data.pageSize || 10
      } else {
        throw new Error(response.msg || '获取公告列表失败')
      }
    } catch (error) {
      console.error('获取公告列表失败:', error)
      announcementList.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  // 添加公告
  const addAnnouncement = async (data) => {
    try {
      const processedData = convertDatesToTimestamps(data, ['startTime', 'endTime'])
      const response = await addGMAnnouncement(processedData)
      if (response.code === 0) {
        // 添加成功后自动刷新列表
        await fetchAnnouncementList()
        return true
      } else {
        throw new Error(response.msg || '添加公告失败')
      }
    } catch (error) {
      console.error('添加公告失败:', error)
      throw error
    }
  }

  // 删除公告
  const deleteAnnouncement = async (data) => {
    try {
      const response = await deleteGMAnnouncement(data)
      if (response.code === 0) {
        // 删除成功后自动刷新列表
        await fetchAnnouncementList()
        return true
      } else {
        throw new Error(response.msg || '删除公告失败')
      }
    } catch (error) {
      console.error('删除公告失败:', error)
      throw error
    }
  }

  // 更新公告
  const updateAnnouncement = async (data) => {
    try {
      const processedData = convertDatesToTimestamps(data, ['startTime', 'endTime'])
      const response = await updateGMAnnouncement(processedData)
      if (response.code === 0) {
        // 更新成功后自动刷新列表
        await fetchAnnouncementList()
        return true
      } else {
        throw new Error(response.msg || '更新公告失败')
      }
    } catch (error) {
      console.error('更新公告失败:', error)
      throw error
    }
  }

  // 置顶公告
  const toppingAnnouncement = async (data) => {
    try {
      const response = await toppingGMAnnouncement(data)
      if (response.code === 0) {
        // 置顶成功后自动刷新列表
        await fetchAnnouncementList()
        return true
      } else {
        throw new Error(response.msg || '置顶公告失败')
      }
    } catch (error) {
      console.error('置顶公告失败:', error)
      throw error
    }
  }

    // 设置搜索条件
    const setSearchInfo = (info) => {
        searchInfo.value = { ...searchInfo.value, ...info }
    }

    // 重置搜索条件
    const resetSearchInfo = () => {
        searchInfo.value = {
            type: '',
            start_time: null,
            end_time: null
        }
    }

    // 设置分页
    const setPage = (newPage) => {
        page.value = newPage
    }

    const setPageSize = (newPageSize) => {
        pageSize.value = newPageSize
        page.value = 1
    }

    // 清空状态
    const clearState = () => {
        announcementList.value = []
        total.value = 0
        page.value = 1
        resetSearchInfo()
    }

  return {
    // 状态
    announcementList,
    currentAnnouncement,
    loading,
    total,
    page,
    pageSize,
    searchInfo,
    
    // 计算属性
    hasAnnouncements,
    totalPages,
    
    // 方法
    fetchAnnouncementList,
    addAnnouncement,
    deleteAnnouncement,
    updateAnnouncement,
    toppingAnnouncement,
    setSearchInfo,
    resetSearchInfo,
    setPage,
    setPageSize,
    clearState
  }
})
