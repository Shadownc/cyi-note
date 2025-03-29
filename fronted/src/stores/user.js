import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  // 用户状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
  const isAuthenticated = computed(() => !!token.value)

  // 登录方法
  async function login(credentials) {
    try {
      const response = await api.auth.login(credentials)
      
      // 保存令牌和用户信息
      token.value = response.token
      user.value = response.user
      
      // 存储到本地存储
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      
      // 跳转到主页
      router.push({ name: 'notes' })
      
      return { success: true }
    } catch (error) {
      console.error('登录失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.message || '登录失败，请检查您的凭据' 
      }
    }
  }

  // 注册方法
  async function register(userData) {
    try {
      const response = await api.auth.register(userData)
      return { success: true }
    } catch (error) {
      console.error('注册失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.message || '注册失败，请稍后再试' 
      }
    }
  }

  // 登出方法
  function logout() {
    // 清除状态
    token.value = ''
    user.value = {}
    
    // 清除本地存储
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    // 返回登录页
    router.push({ name: 'login' })
  }

  // 获取当前用户信息
  async function fetchUserInfo() {
    if (!token.value) {
      console.log('没有找到token，跳过获取用户信息');
      return;
    }
    
    console.log('开始获取用户信息...');
    try {
      const userData = await api.auth.getCurrentUser();
      console.log('成功获取到用户信息:', userData);
      
      // 检查返回的数据是否有效
      if (!userData || !userData.user) {
        console.error('获取的用户数据无效:', userData);
        throw new Error('获取到的用户数据无效');
      }
      
      // 更新用户状态
      user.value = userData.user || userData; // 兼容两种可能的API返回格式
      
      // 存储到本地
      localStorage.setItem('user', JSON.stringify(user.value));
      console.log('用户信息已更新和保存');
      
      return user.value;
    } catch (error) {
      console.error('获取用户信息失败:', error);
      if (error.response?.status === 401) {
        console.log('认证失败，执行登出');
        // 如果认证失败，登出用户
        logout();
      }
      throw error;
    }
  }

  return { 
    token, 
    user, 
    isAuthenticated, 
    login, 
    register, 
    logout, 
    fetchUserInfo 
  }
}) 