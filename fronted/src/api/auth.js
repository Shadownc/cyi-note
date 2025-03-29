import { post, get, put, del } from '@/utils/request'

// 注册
export function register(data) {
  return post('/auth/register', data)
}

// 登录
export function login(data) {
  return post('/auth/login', data)
}

// 获取当前用户信息
export function getCurrentUser() {
  console.log('调用API获取当前用户信息');
  const result = get('/auth/user');
  console.log('API结果:', result);
  return result;
}

// 更新用户个人信息
export function updateProfile(data) {
  return put('/auth/user', data)
}

// 上传用户头像
export function updateAvatar(formData) {
  return post('/auth/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 注销账号
export function deleteAccount() {
  return del('/auth/user')
} 