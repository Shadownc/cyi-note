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
        <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-6 text-center">用户注册</h2>
        
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
              placeholder="请输入用户名"
              required
              @focus="hasInteracted = true"
              @blur="fieldsInteracted.username = true"
            >
            <p v-if="!isUsernameValid && (fieldsInteracted.username || formSubmitted)" class="mt-1 text-sm text-red-600 dark:text-red-400">
              用户名至少需要3个字符
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
              @focus="hasInteracted = true"
              @blur="fieldsInteracted.email = true"
            >
            <p v-if="!isEmailValid && (fieldsInteracted.email || formSubmitted)" class="mt-1 text-sm text-red-600 dark:text-red-400">
              请输入有效的邮箱地址
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
              @focus="hasInteracted = true"
              @blur="fieldsInteracted.password = true"
            >
            <p v-if="!isPasswordValid && (fieldsInteracted.password || formSubmitted)" class="mt-1 text-sm text-red-600 dark:text-red-400">
              密码至少需要6个字符
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
              @focus="hasInteracted = true"
              @blur="fieldsInteracted.confirmPassword = true"
            >
            <p v-if="!isConfirmPasswordValid && (fieldsInteracted.confirmPassword || formSubmitted)" class="mt-1 text-sm text-red-600 dark:text-red-400">
              两次输入的密码不一致
            </p>
          </div>
          
          <div class="mb-4">
            <label for="agreement" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              <input 
                id="agreement" 
                v-model="agreement" 
                type="checkbox" 
                class="h-4 w-4 text-primary-600 border-gray-300 rounded"
                @click="fieldsInteracted.agreement = true"
              >
              <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">同意服务条款和隐私政策</span>
            </label>
            <p v-if="!agreement && (fieldsInteracted.agreement || formSubmitted)" class="mt-1 text-sm text-red-600 dark:text-red-400">
              请同意服务条款和隐私政策
            </p>
          </div>
          
          <div>
            <button 
              type="submit" 
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <svg v-if="isRegistering" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isRegistering ? '注册中...' : '注册' }}
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
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';

const router = useRouter();
const userStore = useUserStore();

// 表单数据
const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const agreement = ref(false);
const error = ref('');
const isRegistering = ref(false);

// 跟踪用户交互状态
const hasInteracted = ref(false);
const formSubmitted = ref(false);
const fieldsInteracted = ref({
  username: false,
  email: false,
  password: false,
  confirmPassword: false,
  agreement: false
});

// 表单验证
const isUsernameValid = computed(() => username.value.length >= 3 || username.value === '');
const isEmailValid = computed(() => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value) || email.value === '');
const isPasswordValid = computed(() => password.value.length >= 6 || password.value === '');
const isConfirmPasswordValid = computed(() => password.value === confirmPassword.value || confirmPassword.value === '');
const isFormValid = computed(() => 
  username.value.length >= 3 && 
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value) && 
  password.value.length >= 6 && 
  password.value === confirmPassword.value && 
  agreement.value
);

// 提交表单
const handleSubmit = async () => {
  formSubmitted.value = true;
  
  // 检查所有字段的交互状态，确保在提交时显示所有错误
  Object.keys(fieldsInteracted.value).forEach(key => {
    fieldsInteracted.value[key] = true;
  });
  
  if (!isFormValid.value) {
    if (!username.value) error.value = '请输入用户名';
    else if (!isUsernameValid.value) error.value = '用户名至少需要3个字符';
    else if (!email.value) error.value = '请输入邮箱';
    else if (!isEmailValid.value) error.value = '请输入有效的邮箱地址';
    else if (!password.value) error.value = '请输入密码';
    else if (!isPasswordValid.value) error.value = '密码至少需要6个字符';
    else if (!confirmPassword.value) error.value = '请确认密码';
    else if (!isConfirmPasswordValid.value) error.value = '两次输入的密码不一致';
    else if (!agreement.value) error.value = '请同意服务条款和隐私政策';
    return;
  }
  
  error.value = '';
  isRegistering.value = true;
  
  try {
    // 使用Pinia store注册
    const result = await userStore.register({
      username: username.value,
      email: email.value,
      password: password.value
    });
    
    if (result.success) {
      // 注册成功后自动登录
      await userStore.login({
        username: username.value,
        password: password.value
      });
    } else {
      error.value = result.message || '注册失败，请稍后重试';
    }
  } catch (err) {
    console.error('注册错误:', err);
    error.value = '注册失败，请稍后重试';
  } finally {
    isRegistering.value = false;
  }
};

// 跳转到登录页面
const goToLogin = () => {
  router.push({ name: 'login' });
};
</script>

<style scoped>
/* 确保页面占满整个视口且没有滚动条 */
:deep(body) {
  overflow: hidden;
}
</style> 