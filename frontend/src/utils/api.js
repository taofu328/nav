import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

api.interceptors.request.use(
  config => {
    // 首先尝试获取admin_token
    let token = localStorage.getItem('admin_token')
    // 如果没有admin_token，尝试获取普通用户token
    if (!token) {
      token = localStorage.getItem('token')
    }
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      // 移除所有token
      localStorage.removeItem('admin_token')
      localStorage.removeItem('token')
      // 跳转到普通登录页面
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
