<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-2xl font-bold mb-6 text-gray-800 dark:text-white">个人资料</h1>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
        <div class="flex flex-col md:flex-row items-start md:items-center mb-6">
          <div class="flex-shrink-0 mb-4 md:mb-0 md:mr-6">
            <div class="relative">
              <div class="w-24 h-24 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center text-primary-700 dark:text-primary-300 text-3xl font-bold border-4 border-white dark:border-gray-700 shadow">
                {{ userInitials }}
              </div>
              <button 
                class="absolute bottom-0 right-0 bg-primary-600 text-white p-1.5 rounded-full shadow-md hover:bg-primary-700"
                title="更换头像"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </button>
            </div>
          </div>
          
          <div class="flex-grow">
            <h2 class="text-xl font-bold text-gray-800 dark:text-white mb-1">{{ user.username }}</h2>
            <p class="text-gray-600 dark:text-gray-400 mb-2">{{ user.email }}</p>
            <div class="flex items-center">
              <span class="px-2 py-1 text-xs rounded-full bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-300">
                {{ user.role === 'admin' ? '管理员' : '普通用户' }}
              </span>
              <span class="ml-3 text-sm text-gray-500 dark:text-gray-400">
                注册于 {{ formatDate(user.created_at) }}
              </span>
            </div>
          </div>
        </div>
        
        <div class="border-t dark:border-gray-700 pt-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400 mb-1">笔记总数</div>
              <div class="text-lg font-medium text-gray-800 dark:text-white">{{ stats.notesCount }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400 mb-1">使用的标签</div>
              <div class="text-lg font-medium text-gray-800 dark:text-white">{{ stats.tagsCount }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
        <h2 class="text-xl font-bold text-gray-800 dark:text-white mb-4">账号设置</h2>
        
        <form @submit.prevent="updateProfile">
          <div v-if="message" class="mb-4 p-3 rounded-md text-sm" :class="messageType === 'success' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'">
            {{ message }}
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">用户名</label>
              <input 
                id="username" 
                v-model="form.username" 
                type="text" 
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white" 
                required
              >
            </div>
            
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">邮箱</label>
              <input 
                id="email" 
                v-model="form.email" 
                type="email" 
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white" 
                required
              >
            </div>
          </div>
          
          <div class="mb-6">
            <h3 class="text-lg font-medium text-gray-800 dark:text-white mb-3">修改密码</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">如果不需要修改密码，请将以下字段留空</p>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label for="currentPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">当前密码</label>
                <input 
                  id="currentPassword" 
                  v-model="form.currentPassword" 
                  type="password" 
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
                >
              </div>
              
              <div>
                <label for="newPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">新密码</label>
                <input 
                  id="newPassword" 
                  v-model="form.newPassword" 
                  type="password" 
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
                >
              </div>
            </div>
          </div>
          
          <div class="flex justify-end">
            <button 
              type="submit" 
              class="btn btn-primary px-6"
              :disabled="loading"
            >
              <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              保存修改
            </button>
          </div>
        </form>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-bold text-gray-800 dark:text-white mb-4">注销账号</h2>
        <p class="text-gray-600 dark:text-gray-400 mb-4">注销账号后，您的所有数据将被永久删除且无法恢复。</p>
        <button 
          @click="confirmDeleteAccount" 
          class="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
        >
          注销账号
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

// 用户信息（通常从 API 或 vuex/pinia 获取）
const user = ref({
  username: '张三',
  email: 'zhangsan@example.com',
  role: 'user',
  created_at: new Date('2023-01-15').toISOString()
});

// 用户统计数据（模拟数据）
const stats = ref({
  notesCount: 15,
  tagsCount: 8
});

// 表单数据
const form = ref({
  username: '',
  email: '',
  currentPassword: '',
  newPassword: ''
});

const loading = ref(false);
const message = ref('');
const messageType = ref('success');

// 计算用户首字母头像
const userInitials = computed(() => {
  if (!user.value.username) return '';
  return user.value.username.charAt(0).toUpperCase();
});

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

// 初始化表单数据
onMounted(() => {
  form.value.username = user.value.username;
  form.value.email = user.value.email;
});

// 更新个人资料
const updateProfile = async () => {
  loading.value = true;
  message.value = '';
  
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // TODO: 实现实际的更新个人资料 API 调用
    // const response = await axios.put('/api/auth/user', {
    //   username: form.value.username,
    //   email: form.value.email,
    //   current_password: form.value.currentPassword,
    //   new_password: form.value.newPassword || undefined
    // });
    
    // 更新本地用户信息
    user.value.username = form.value.username;
    user.value.email = form.value.email;
    
    // 清空密码字段
    form.value.currentPassword = '';
    form.value.newPassword = '';
    
    message.value = '个人资料更新成功';
    messageType.value = 'success';
  } catch (err) {
    console.error('更新个人资料失败:', err);
    message.value = '更新个人资料失败，请稍后重试';
    messageType.value = 'error';
  } finally {
    loading.value = false;
  }
};

// 确认删除账号
const confirmDeleteAccount = () => {
  if (window.confirm('您确定要注销账号吗？此操作不可逆，所有数据将被永久删除。')) {
    // TODO: 实现账号注销的 API 调用
    alert('账号注销功能尚未实现');
  }
};
</script> 