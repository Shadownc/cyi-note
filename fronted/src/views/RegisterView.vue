<template>
  <div class="flex min-h-screen bg-gray-100 dark:bg-gray-900 overflow-hidden">
    <div class="m-auto w-full max-w-md px-6">
      <div class="text-center mb-8">
        <router-link to="/" class="inline-flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
          <span class="ml-3 text-xl font-bold text-gray-800 dark:text-white">知识库</span>
        </router-link>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-8">
        <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-6 text-center">用户注册</h2>
        
        <form @submit.prevent="handleRegister">
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
              placeholder="请输入用户名"
              required
            >
            <p v-if="v$.username.$error" class="mt-1 text-sm text-red-600 dark:text-red-400">
              {{ v$.username.$errors[0].$message }}
            </p>
          </div>
          
          <div class="mb-4">
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">邮箱</label>
            <input 
              id="email" 
              v-model="email" 
              type="email" 
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white" 
              placeholder="请输入邮箱"
              required
            >
            <p v-if="v$.email.$error" class="mt-1 text-sm text-red-600 dark:text-red-400">
              {{ v$.email.$errors[0].$message }}
            </p>
          </div>
          
          <div class="mb-4">
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">密码</label>
            <input 
              id="password" 
              v-model="password" 
              type="password" 
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white" 
              placeholder="请输入密码"
              required
            >
            <p v-if="v$.password.$error" class="mt-1 text-sm text-red-600 dark:text-red-400">
              {{ v$.password.$errors[0].$message }}
            </p>
          </div>
          
          <div class="mb-6">
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">确认密码</label>
            <input 
              id="confirmPassword" 
              v-model="confirmPassword" 
              type="password" 
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white" 
              placeholder="请再次输入密码"
              required
            >
            <p v-if="v$.confirmPassword.$error" class="mt-1 text-sm text-red-600 dark:text-red-400">
              {{ v$.confirmPassword.$errors[0].$message }}
            </p>
          </div>
          
          <div>
            <button 
              type="submit" 
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              :disabled="loading"
            >
              <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ loading ? '注册中...' : '注册' }}
            </button>
          </div>
        </form>
        
        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600 dark:text-gray-400">
            已有账号? 
            <router-link to="/login" class="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400">立即登录</router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useVuelidate } from '@vuelidate/core';
import { required, email as emailValidator, minLength, sameAs, helpers } from '@vuelidate/validators';

const router = useRouter();
const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const loading = ref(false);
const error = ref('');

// 表单验证规则
const rules = {
  username: { 
    required: helpers.withMessage('用户名不能为空', required),
    minLength: helpers.withMessage('用户名长度至少为3个字符', minLength(3)) 
  },
  email: { 
    required: helpers.withMessage('邮箱不能为空', required),
    email: helpers.withMessage('邮箱格式不正确', emailValidator) 
  },
  password: { 
    required: helpers.withMessage('密码不能为空', required),
    minLength: helpers.withMessage('密码长度至少为6个字符', minLength(6)) 
  },
  confirmPassword: { 
    required: helpers.withMessage('请确认密码', required),
    sameAsPassword: helpers.withMessage('两次输入的密码不一致', sameAs(password)) 
  }
};

const v$ = useVuelidate(rules, { username, email, password, confirmPassword });

const handleRegister = async () => {
  const result = await v$.value.$validate();
  if (!result) return;
  
  loading.value = true;
  error.value = '';
  
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // TODO: 实现实际的注册API调用
    // const response = await axios.post('/api/auth/register', {
    //   username: username.value,
    //   email: email.value,
    //   password: password.value
    // });
    
    // 模拟注册成功，自动登录
    localStorage.setItem('user', JSON.stringify({ 
      username: username.value,
      email: email.value,
      role: 'user'
    }));
    
    // 注册成功后重定向到登录页面
    router.push('/login');
  } catch (err) {
    console.error('注册失败:', err);
    error.value = '注册失败，该用户名或邮箱可能已被使用';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* 确保页面占满整个视口且没有滚动条 */
:deep(body) {
  overflow: hidden;
}
</style> 