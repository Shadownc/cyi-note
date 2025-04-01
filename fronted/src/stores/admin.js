import { ref } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api'

export const useAdminStore = defineStore('admin', () => {
  // 状态
  const users = ref([])
  const selectedUser = ref(null)
  const pagination = ref({
    page: 1,
    pageSize: 10,
    total: 0,
    totalPage: 0
  })
  const isLoading = ref(false)
  const keyword = ref('')

  // 获取用户列表
  async function fetchUsers(page = 1, pageSize = 10, searchKeyword = '') {
    isLoading.value = true
    try {
      const params = {
        page: page,
        pageSize: pageSize
      }
      
      if (searchKeyword) {
        params.keyword = searchKeyword
        keyword.value = searchKeyword
      } else {
        keyword.value = ''
      }
      
      const response = await api.admin.getUsers(params)
      users.value = response.users
      pagination.value = response.pagination
      return response
    } catch (error) {
      console.error('获取用户列表失败:', error)
      return null
    } finally {
      isLoading.value = false
    }
  }

  // 获取用户详情
  async function fetchUserDetail(id) {
    isLoading.value = true
    try {
      const response = await api.admin.getUser(id)
      selectedUser.value = response.user
      return response
    } catch (error) {
      console.error(`获取用户 ${id} 详情失败:`, error)
      return null
    } finally {
      isLoading.value = false
    }
  }

  // 创建用户
  async function createUser(userData) {
    isLoading.value = true
    try {
      const response = await api.admin.createUser(userData)
      // 刷新用户列表
      await fetchUsers(pagination.value.page, pagination.value.pageSize, keyword.value)
      return { success: true, data: response }
    } catch (error) {
      console.error('创建用户失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.message || '创建用户失败' 
      }
    } finally {
      isLoading.value = false
    }
  }

  // 更新用户角色
  async function updateUserRole(id, role) {
    isLoading.value = true
    try {
      const response = await api.admin.updateUserRole(id, role)
      // 更新本地状态
      if (selectedUser.value && selectedUser.value.id === id) {
        selectedUser.value.role = role
      }
      // 更新用户列表中的用户
      const userIndex = users.value.findIndex(u => u.id === id)
      if (userIndex !== -1) {
        users.value[userIndex].role = role
      }
      return { success: true, data: response }
    } catch (error) {
      console.error(`更新用户 ${id} 角色失败:`, error)
      return { 
        success: false, 
        message: error.response?.data?.message || '更新用户角色失败' 
      }
    } finally {
      isLoading.value = false
    }
  }

  // 删除用户
  async function deleteUser(id) {
    isLoading.value = true
    try {
      await api.admin.deleteUser(id)
      // 从列表中移除
      users.value = users.value.filter(user => user.id !== id)
      // 如果当前选中的用户被删除，清空选中
      if (selectedUser.value && selectedUser.value.id === id) {
        selectedUser.value = null
      }
      // 更新分页信息
      pagination.value.total--
      return { success: true }
    } catch (error) {
      console.error(`删除用户 ${id} 失败:`, error)
      return { 
        success: false, 
        message: error.response?.data?.message || '删除用户失败' 
      }
    } finally {
      isLoading.value = false
    }
  }

  return {
    users,
    selectedUser,
    pagination,
    isLoading,
    keyword,
    fetchUsers,
    fetchUserDetail,
    createUser,
    updateUserRole,
    deleteUser
  }
}) 