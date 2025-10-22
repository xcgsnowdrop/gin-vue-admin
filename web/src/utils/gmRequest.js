import axios from 'axios'
import { useUserStore } from '@/pinia/modules/user'
import { ElMessage } from 'element-plus'
import { emitter } from '@/utils/bus'
import router from '@/router/index'

// 创建专门用于GM模块的service实例
const gmService = axios.create({
  baseURL: '', // 不设置baseURL，直接使用相对路径
  timeout: 99999
})

// 为GM service添加请求拦截器
gmService.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    config.headers = {
      'Content-Type': 'application/json',
      'x-token': userStore.token,
      'x-user-id': userStore.userInfo.ID,
      ...config.headers
    }
    return config
  },
  (error) => {
    emitter.emit('show-error', {
      code: 'request',
      message: error.message || '请求发送失败'
    })
    return error
  }
)

// 为GM service添加响应拦截器
gmService.interceptors.response.use(
  (response) => {
    const userStore = useUserStore()
    if (response.headers['new-token']) {
      userStore.setToken(response.headers['new-token'])
    }
    if (typeof response.data.status === 'undefined') {
      return response
    }
    if (response.data.status === 0 || response.headers.success === 'true') {
      if (response.headers.msg) {
        response.data.msg = decodeURI(response.headers.msg)
      }
      return response.data
    } else {
      ElMessage({
        showClose: true,
        message: response.data.msg || decodeURI(response.headers.msg),
        type: 'error'
      })
      return response.data.msg ? response.data : response
    }
  },
  (error) => {
    if (!error.response) {
      // 网络错误
      emitter.emit('show-error', {
        code: 'network',
        message: error.response?.data?.msg || '请求失败'
      })
      return Promise.reject(error)
    }

    // HTTP 状态码错误
    if (error.response.status === 401) {
      emitter.emit('show-error', {
        code: '401',
        message: error.response?.data?.msg || '请求失败',
        fn: () => {
          const userStore = useUserStore()
          userStore.ClearStorage()
          router.push({ name: 'Login', replace: true })
        }
      })
      return Promise.reject(error)
    }

    emitter.emit('show-error', {
      code: error.response.status,
      message: error.response?.data?.msg || '请求失败'
    })
    return Promise.reject(error)
  }
)

// 导出GM service实例
export default gmService
