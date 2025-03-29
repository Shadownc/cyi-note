import axios from 'axios'

// 创建axios实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api', // 从环境变量获取API基础URL
  timeout: 15000, // 请求超时时间
  withCredentials: false // 跨域请求不带凭证
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    console.log(`请求API: ${config.method.toUpperCase()} ${config.url}`, config.params || config.data)
    
    // 从localStorage获取token，如果存在则添加到请求头
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
      console.log('请求添加了Authorization令牌')
    }
    return config
  },
  error => {
    console.error('请求错误', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    console.log(`API响应: ${response.config.method.toUpperCase()} ${response.config.url}`, res)
    
    // 处理统一响应格式
    if (res.hasOwnProperty('success')) {
      // 新的统一格式
      if (res.success === false) {
        // 处理业务层面的错误
        const errorMsg = res.error || '请求失败'
        console.error('API业务错误:', errorMsg)
        return Promise.reject(new Error(errorMsg))
      }
      
      // 如果有data字段，直接返回data，否则返回整个响应
      const result = res.data !== undefined ? res.data : res
      console.log('统一格式下解析后的数据:', result)
      return result
    }
    
    // 兼容旧的API格式 - 直接返回响应数据
    console.log('非统一格式响应，直接返回:', res)
    return res
  },
  error => {
    // 处理错误响应
    if (error.response) {
      const { status, data } = error.response
      const method = error.config?.method?.toUpperCase() || 'UNKNOWN'
      const url = error.config?.url || 'UNKNOWN'
      const errorMsg = data?.error || data?.message || '服务器错误'
      
      console.error(`API错误: ${method} ${url} [${status}] - ${errorMsg}`)
      
      // 处理401未授权错误，重定向到登录页面
      if (status === 401) {
        console.log('认证失败，重定向到登录页')
        localStorage.removeItem('token')
        // 使用window.location进行重定向，而不是router
        window.location.href = '/#/login'
      }
    } else {
      console.error('网络错误', error.message)
    }
    
    return Promise.reject(error)
  }
)

// 封装GET请求
export function get(url, params, config = {}) {
  return service({
    url,
    method: 'get',
    params,
    ...config
  })
}

// 封装POST请求
export function post(url, data, config = {}) {
  return service({
    url,
    method: 'post',
    data,
    ...config
  })
}

// 封装PUT请求
export function put(url, data, config = {}) {
  return service({
    url,
    method: 'put',
    data,
    ...config
  })
}

// 封装DELETE请求
export function del(url, params, config = {}) {
  return service({
    url,
    method: 'delete',
    params,
    ...config
  })
}

export default service 