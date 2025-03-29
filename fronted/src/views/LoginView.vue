<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <router-link to="/" class="text-primary-600 text-2xl font-bold flex items-center justify-center">
          <span class="i-carbon-notebook mr-2 text-3xl"></span>
          知识库
        </router-link>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-8">
        <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-6 text-center">用户登录</h2>
        
        <form @submit.prevent="handleSubmit">
          <div v-if="error" class="mb-4 p-3 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-md text-sm">
            {{ error }}
          </div>
          
          <div class="mb-4">
            <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">用户名</label>
            <input 
              id="username" 
              v-model="username" 
              type="text" 
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white" 
              placeholder="请输入用户名或邮箱"
              required
              @focus="hasInteracted = true"
              @blur="fieldsInteracted.username = true"
            >
            <p v-if="!username && fieldsInteracted.username" class="mt-1 text-sm text-red-600 dark:text-red-400">
              请输入用户名或邮箱
            </p>
          </div>
          
          <div class="mb-6">
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">密码</label>
            <input 
              id="password" 
              v-model="password" 
              type="password" 
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white" 
              placeholder="请输入密码"
              required
              @focus="hasInteracted = true"
              @blur="fieldsInteracted.password = true"
            >
            <p v-if="!password && fieldsInteracted.password" class="mt-1 text-sm text-red-600 dark:text-red-400">
              请输入密码
            </p>
          </div>
          
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center">
              <input 
                id="remember" 
                v-model="remember" 
                type="checkbox" 
                class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              >
              <label for="remember" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">记住我</label>
            </div>
            <div class="text-sm">
              <a href="#" class="text-primary-600 hover:text-primary-500 dark:text-primary-400">忘记密码?</a>
            </div>
          </div>
          
          <div>
            <button 
              type="submit" 
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              :disabled="isLoggingIn || (!username || !password)"
            >
              <svg v-if="isLoggingIn" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isLoggingIn ? '登录中...' : '登录' }}
            </button>
          </div>
        </form>
        
        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600 dark:text-gray-400">
            还没有账号? 
            <router-link to="/register" class="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400">立即注册</router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';

const router = useRouter();
const userStore = useUserStore();

// 表单数据
const username = ref('');
const password = ref('');
const remember = ref(false);
const error = ref('');
const isLoggingIn = ref(false);

// 跟踪用户交互状态
const hasInteracted = ref(false);
const formSubmitted = ref(false);
const fieldsInteracted = ref({
  username: false,
  password: false
});

// 提交表单
const handleSubmit = async () => {
  formSubmitted.value = true;
  fieldsInteracted.username = true;
  fieldsInteracted.password = true;
  
  // 表单验证
  if (!username.value || !password.value) {
    error.value = '请输入用户名和密码';
    return;
  }
  
  error.value = '';
  isLoggingIn.value = true;
  
  try {
    // 使用Pinia store登录
    const result = await userStore.login({
      username: username.value,
      password: password.value
    });
    
    if (!result.success) {
      error.value = result.message || '登录失败，请检查用户名和密码';
    }
  } catch (err) {
    console.error('登录错误:', err);
    error.value = '登录失败，请稍后重试';
  } finally {
    isLoggingIn.value = false;
  }
};

// 跳转到注册页面
const goToRegister = () => {
  router.push({ name: 'register' });
};
</script>

<style scoped>
/* 确保页面占满整个视口且没有滚动条 */
:deep(body) {
  overflow: hidden;
}
</style> 