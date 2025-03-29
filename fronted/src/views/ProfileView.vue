<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-2xl font-bold mb-6 text-gray-800 dark:text-white">个人资料</h1>
      
      <div v-if="!isDataReady && !error" class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 text-center py-12">
        <div class="animate-spin w-12 h-12 mx-auto mb-4 border-4 border-primary-500 border-t-transparent rounded-full"></div>
        <p class="text-gray-600 dark:text-gray-400">正在加载个人资料...</p>
      </div>
      
      <div v-else-if="error" class="bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 p-4 rounded-lg mb-6">
        <p class="font-medium">加载失败</p>
        <p>{{ error }}</p>
        <button @click="retryLoading" class="mt-3 px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600">重试</button>
      </div>
      
      <template v-else>
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
          <div class="flex flex-col md:flex-row items-start md:items-center mb-6">
            <div class="flex-shrink-0 mb-4 md:mb-0 md:mr-6">
              <div class="relative">
                <div v-if="userAvatar" class="w-24 h-24 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center overflow-hidden border-4 border-white dark:border-gray-700 shadow">
                  <img :src="userAvatar" alt="用户头像" class="w-full h-full object-cover" />
                </div>
                <div v-else class="w-24 h-24 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center text-primary-700 dark:text-primary-300 text-3xl font-bold border-4 border-white dark:border-gray-700 shadow">
                  {{ userInitials }}
                </div>
                <button 
                  class="absolute bottom-0 right-0 bg-primary-600 text-white p-1.5 rounded-full shadow-md hover:bg-primary-700"
                  title="更换头像"
                  @click="openAvatarUpload"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </button>
                <input 
                  type="file" 
                  ref="avatarInput" 
                  class="hidden" 
                  accept="image/*"
                  @change="handleAvatarChange" 
                />
              </div>
            </div>
            
            <div class="flex-grow">
              <h2 class="text-xl font-bold text-gray-800 dark:text-white mb-1">{{ displayName }}</h2>
              <p class="text-gray-600 dark:text-gray-400 mb-2">{{ displayEmail }}</p>
              <div class="flex items-center">
                <span class="px-2 py-1 text-xs rounded-full bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-300">
                  {{ displayRole }}
                </span>
                <span class="ml-3 text-sm text-gray-500 dark:text-gray-400">
                  注册于 {{ formatDate(displayCreatedAt) }}
                </span>
              </div>
            </div>
          </div>
          
          <div class="border-t dark:border-gray-700 pt-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <div class="text-sm text-gray-500 dark:text-gray-400 mb-1">笔记总数</div>
                <div class="text-lg font-medium text-gray-800 dark:text-white">
                  {{ isLoading ? '加载中...' : stats.notesCount }}
                </div>
              </div>
              <div>
                <div class="text-sm text-gray-500 dark:text-gray-400 mb-1">使用的标签</div>
                <div class="text-lg font-medium text-gray-800 dark:text-white">
                  {{ isLoading ? '加载中...' : stats.tagsCount }}
                </div>
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
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useUserStore } from '@/stores/user';
import { useNotesStore } from '@/stores/notes';
import { useTagsStore } from '@/stores/tags';
import { storeToRefs } from 'pinia';
import api from '@/api';

const userStore = useUserStore();
const notesStore = useNotesStore();
const tagsStore = useTagsStore();

// 从store获取状态
const { user } = storeToRefs(userStore);
const { notes, totalNotes } = storeToRefs(notesStore);
const { tags } = storeToRefs(tagsStore);

// 本地状态
const activeTab = ref('profile'); // 'profile', 'password', 'preferences'
const isLoading = ref(true);
const isDataReady = ref(false);
const error = ref('');
const isSaving = ref(false);
const successMessage = ref('');
const errorMessage = ref('');
const loading = ref(false);
const message = ref('');
const messageType = ref('');
const avatarInput = ref(null);

// 安全处理的用户数据显示
const userAvatar = computed(() => {
  return user.value?.avatar || null;
});

const displayName = computed(() => {
  return user.value?.username || '未知用户';
});

const displayEmail = computed(() => {
  return user.value?.email || '暂无邮箱';
});

const displayRole = computed(() => {
  return user.value?.role === 'admin' ? '管理员' : '普通用户';
});

const displayCreatedAt = computed(() => {
  return user.value?.created_at || null;
});

// 用户统计信息
const stats = computed(() => ({
  notesCount: totalNotes.value || 0,
  tagsCount: tags.value ? tags.value.length : 0
}));

// 表单数据
const form = ref({
  username: '',
  email: '',
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
  theme: 'blue',
  darkMode: 'auto'
});

// 主题设置
const themes = [
  { id: 'blue', name: '蓝色', color: 'bg-blue-500' },
  { id: 'green', name: '绿色', color: 'bg-green-500' },
  { id: 'purple', name: '紫色', color: 'bg-purple-500' },
  { id: 'indigo', name: '靛蓝', color: 'bg-indigo-500' },
  { id: 'red', name: '红色', color: 'bg-red-500' }
];

// 黑暗模式选项
const darkModeOptions = [
  { id: 'light', name: '浅色模式' },
  { id: 'dark', name: '深色模式' },
  { id: 'auto', name: '跟随系统' }
];

// 获取用户头像显示用字符
const userInitials = computed(() => {
  if (!user.value || !user.value.username) return '?';
  
  // 检查用户名是否包含中文字符
  if (/[\u4e00-\u9fa5]/.test(user.value.username)) {
    // 如果是中文，取第一个字
    return user.value.username.charAt(0);
  } else {
    // 如果是英文或其他，取首字母大写
    return user.value.username.charAt(0).toUpperCase();
  }
});

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '暂无数据';
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return '日期格式无效';
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  } catch (e) {
    console.error('日期格式化错误:', e);
    return '日期格式错误';
  }
};

// 打开头像上传文件选择器
const openAvatarUpload = () => {
  avatarInput.value.click();
};

// 处理头像选择
const handleAvatarChange = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  try {
    loading.value = true;
    message.value = '';
    
    // 创建FormData对象
    const formData = new FormData();
    formData.append('avatar', file);
    
    // 调用API上传头像
    const response = await api.auth.updateAvatar(formData);
    
    // 更新用户信息
    await userStore.fetchUserInfo();
    
    message.value = '头像更新成功';
    messageType.value = 'success';
  } catch (error) {
    console.error('头像上传失败:', error);
    message.value = '头像上传失败，请稍后重试';
    messageType.value = 'error';
  } finally {
    loading.value = false;
    // 清除文件输入，以便可以再次选择同一文件
    event.target.value = '';
  }
};

// 更新基本信息
const updateProfile = async () => {
  loading.value = true;
  message.value = '';
  
  try {
    const userData = {
      username: form.value.username,
      email: form.value.email
    };
    
    // 如果有填写密码，则添加到更新数据中
    if (form.value.currentPassword && form.value.newPassword) {
      userData.currentPassword = form.value.currentPassword;
      userData.newPassword = form.value.newPassword;
    }
    
    // 调用API更新用户信息
    const response = await api.auth.updateProfile(userData);
    
    // 检查是否修改了密码
    if (response.passwordChanged) {
      message.value = response.message || '密码已更新，将在3秒后重新登录';
      messageType.value = 'success';
      
      // 3秒后登出并重定向到登录页
      setTimeout(() => {
        userStore.logout();
      }, 3000);
    } else {
      // 常规更新，重新获取最新的用户信息
      await userStore.fetchUserInfo();
      
      // 清空密码字段
      form.value.currentPassword = '';
      form.value.newPassword = '';
      form.value.confirmPassword = '';
      
      message.value = response.message || '个人信息已更新';
      messageType.value = 'success';
    }
  } catch (error) {
    console.error('更新个人信息失败:', error);
    message.value = error.response?.data?.message || '更新信息失败，请稍后重试';
    messageType.value = 'error';
  } finally {
    loading.value = false;
  }
};

// 重试加载
const retryLoading = async () => {
  error.value = '';
  await loadUserData();
};

// 加载用户数据
const loadUserData = async () => {
  isLoading.value = true;
  isDataReady.value = false;
  error.value = '';
  try {
    console.log('开始加载用户数据...');
    
    // 先检查本地存储是否有用户数据
    const storedUser = JSON.parse(localStorage.getItem('user') || '{}');
    console.log('本地存储的用户数据:', storedUser);
    
    if (Object.keys(storedUser).length > 0) {
      // 先使用本地存储的数据快速显示
      console.log('使用本地存储的用户数据');
      user.value = storedUser;
    }
    
    // 然后从服务器获取最新数据
    console.log('从服务器获取最新用户数据');
    await userStore.fetchUserInfo();
    console.log('获取到的用户数据:', user.value);
    
    // 加载笔记和标签数据
    console.log('加载笔记和标签数据');
    await notesStore.fetchNotes();
    await tagsStore.fetchTags();
    
    // 填充表单数据
    console.log('填充表单数据');
    form.value.username = user.value.username || '';
    form.value.email = user.value.email || '';
    
    // 加载主题设置
    form.value.theme = localStorage.getItem('theme') || 'blue';
    form.value.darkMode = localStorage.getItem('darkMode') || 'auto';
    
    console.log('数据加载完成');
    isDataReady.value = true;
  } catch (err) {
    console.error('加载个人资料数据失败:', err);
    // 尝试使用本地存储的数据作为后备
    const backupUser = JSON.parse(localStorage.getItem('user') || '{}');
    if (Object.keys(backupUser).length > 0) {
      user.value = backupUser;
      form.value.username = backupUser.username || '';
      form.value.email = backupUser.email || '';
      isDataReady.value = true;
    } else {
      error.value = '无法加载用户数据，请检查您的网络连接并刷新页面';
    }
  } finally {
    isLoading.value = false;
  }
};

// 应用主题和深色模式
const applyTheme = (theme, darkMode) => {
  // 移除所有主题类
  document.documentElement.classList.forEach(className => {
    if (className.startsWith('theme-')) {
      document.documentElement.classList.remove(className);
    }
  });
  
  // 添加新主题类
  document.documentElement.classList.add(`theme-${theme}`);
  
  // 应用深色模式
  if (darkMode === 'dark') {
    document.documentElement.classList.add('dark');
  } else if (darkMode === 'light') {
    document.documentElement.classList.remove('dark');
  } else if (darkMode === 'auto') {
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }
};

// 初始化数据
onMounted(async () => {
  await loadUserData();
});

// 确认删除账号
const confirmDeleteAccount = () => {
  if (window.confirm('您确定要注销账号吗？此操作不可逆，所有数据将被永久删除。')) {
    loading.value = true;
    message.value = '';
    
    // 调用账号注销API
    api.auth.deleteAccount()
      .then(() => {
        // 成功后登出用户
        userStore.logout();
      })
      .catch(error => {
        console.error('账号注销失败:', error);
        message.value = '账号注销失败，请稍后重试';
        messageType.value = 'error';
        loading.value = false;
      });
  }
};
</script> 