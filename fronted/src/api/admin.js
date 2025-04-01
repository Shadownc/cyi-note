import { get, post, put, del } from '@/utils/request'

// 获取用户列表
export function getUsers(params) {
  return get('/admin/users', params)
}

// 获取用户详情
export function getUser(id) {
  return get(`/admin/users/${id}`)
}

// 创建用户
export function createUser(data) {
  return post('/admin/users', data)
}

// 更新用户角色
export function updateUserRole(id, role) {
  return put(`/admin/users/${id}/role`, { role })
}

// 删除用户
export function deleteUser(id) {
  return del(`/admin/users/${id}`)
} 