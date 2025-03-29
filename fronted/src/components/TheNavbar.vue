<template>
  <nav class="bg-primary-900 dark:bg-gray-900 text-white shadow-lg sticky top-0 z-50 transition-colors duration-300">
    <div class="container mx-auto">
      <div class="flex justify-between items-center py-3">
        <!-- Logo和标题 -->
        <div class="flex items-center space-x-4">
          <router-link to="/" class="text-xl font-bold flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 mr-2 text-white dark:text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            <span class="text-white">知识库</span>
          </router-link>
        </div>

        <!-- 桌面端导航 -->
        <div class="hidden md:flex items-center space-x-4">
          <router-link to="/" class="px-3 py-2 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white">首页</router-link>
          
          <!-- 已登录用户可见的导航 -->
          <template v-if="isAuthenticated">
            <router-link to="/notes" class="px-3 py-2 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white">我的笔记</router-link>
            <router-link to="/tags" class="px-3 py-2 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white">标签管理</router-link>
          </template>
          
          <div class="relative group">
            <input 
              type="text" 
              placeholder="搜索..." 
              class="bg-primary-800 rounded-full px-4 py-2 pl-10 focus:outline-none focus:ring-2 focus:ring-white text-white placeholder-gray-300 text-sm w-48 transition-all duration-300 group-hover:w-56"
            >
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              class="h-5 w-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-white" 
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          
          <!-- 主题设置按钮 -->
          <ThemeSettings />
          
          <!-- 用户菜单 -->
          <div class="relative ml-3" v-if="isAuthenticated">
            <button 
              @click="toggleUserMenu"
              class="flex items-center text-white hover:bg-primary-800 px-3 py-2 rounded-md focus:outline-none"
            >
              <div class="w-8 h-8 rounded-full bg-primary-700 flex items-center justify-center text-white font-medium mr-2">
                {{ userInitials }}
              </div>
              <span class="hidden sm:inline">{{ username }}</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            
            <div 
              v-show="userMenuOpen"
              class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg py-1 z-50 animate-fadeIn"
            >
              <router-link 
                to="/profile" 
                class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700"
                @click="userMenuOpen = false"
              >
                <div class="flex items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  个人资料
                </div>
              </router-link>
              <div class="border-t border-gray-200 dark:border-gray-700 my-1"></div>
              <button 
                @click="logout" 
                class="block w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                <div class="flex items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                  退出登录
                </div>
              </button>
            </div>
          </div>
          
          <!-- 未登录状态的按钮 -->
          <div v-else class="flex items-center space-x-2">
            <router-link 
              to="/login" 
              class="px-3 py-2 rounded-md text-white font-medium hover:bg-primary-700 transition-colors"
            >
              登录
            </router-link>
            <router-link 
              to="/register" 
              class="px-3 py-2 bg-white text-primary-700 hover:bg-gray-100 rounded-md font-medium transition-colors"
            >
              注册
            </router-link>
          </div>
        </div>

        <!-- 移动端菜单按钮和主题设置 -->
        <div class="md:hidden flex items-center">
          <ThemeSettings class="mr-2" />
          <button @click="toggleMobileMenu" class="text-white focus:outline-none p-2 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path v-if="mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7" />
            </svg>
          </button>
        </div>
      </div>

      <!-- 移动端导航菜单 -->
      <div v-if="mobileMenuOpen" class="md:hidden pb-4 space-y-2 animate-fadeIn">
        <router-link 
          to="/" 
          class="block px-4 py-3 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white"
          @click="mobileMenuOpen = false"
        >
          首页
        </router-link>
        
        <!-- 已登录用户可见的移动端导航 -->
        <template v-if="isAuthenticated">
          <router-link 
            to="/notes" 
            class="block px-4 py-3 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white"
            @click="mobileMenuOpen = false"
          >
            我的笔记
          </router-link>
          <router-link 
            to="/tags" 
            class="block px-4 py-3 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white"
            @click="mobileMenuOpen = false"
          >
            标签管理
          </router-link>
          <router-link 
            to="/profile" 
            class="block px-4 py-3 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white"
            @click="mobileMenuOpen = false"
          >
            个人资料
          </router-link>
          <button 
            @click="logout" 
            class="block w-full text-left px-4 py-3 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white"
          >
            退出登录
          </button>
        </template>
        
        <!-- 未登录状态的移动端导航 -->
        <template v-else>
          <router-link 
            to="/login" 
            class="block px-4 py-3 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white"
            @click="mobileMenuOpen = false"
          >
            登录
          </router-link>
          <router-link 
            to="/register" 
            class="block px-4 py-3 rounded-md hover:bg-primary-700 dark:hover:bg-primary-700 transition-colors font-medium text-white"
            @click="mobileMenuOpen = false"
          >
            注册
          </router-link>
        </template>
        
        <div class="mt-3 relative">
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            class="h-5 w-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-white" 
            fill="none" 
            viewBox="0 0 24 24" 
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input 
            type="text" 
            placeholder="搜索..." 
            class="bg-primary-800 rounded-full px-4 py-2 pl-10 focus:outline-none focus:ring-2 focus:ring-white text-white placeholder-gray-300 text-sm w-full"
          >
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import ThemeSettings from './ThemeSettings.vue';

const router = useRouter();
const mobileMenuOpen = ref(false);
const userMenuOpen = ref(false);

// 用户认证状态
const isAuthenticated = ref(false);
const username = ref('');
const userInitials = computed(() => {
  return username.value ? username.value.charAt(0).toUpperCase() : '';
});

// 切换移动菜单
const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value;
  
  // 关闭用户菜单
  if (mobileMenuOpen.value) {
    userMenuOpen.value = false;
  }
};

// 切换用户菜单
const toggleUserMenu = () => {
  userMenuOpen.value = !userMenuOpen.value;
};

// 退出登录
const logout = () => {
  localStorage.removeItem('user');
  isAuthenticated.value = false;
  userMenuOpen.value = false;
  router.push('/login');
};

// 检查用户登录状态
const checkAuth = () => {
  const user = localStorage.getItem('user');
  if (user) {
    const userData = JSON.parse(user);
    isAuthenticated.value = true;
    username.value = userData.username;
  } else {
    isAuthenticated.value = false;
    username.value = '';
  }
};

// 初始化时检查用户登录状态
onMounted(() => {
  checkAuth();
  
  // 点击外部关闭用户菜单
  document.addEventListener('click', (e) => {
    const userMenu = document.querySelector('.relative.ml-3');
    if (userMenu && !userMenu.contains(e.target) && userMenuOpen.value) {
      userMenuOpen.value = false;
    }
  });
});

// 监听路由变化，关闭菜单
watch(() => router.currentRoute.value, () => {
  mobileMenuOpen.value = false;
  checkAuth();
}, { deep: true });
</script>

<style scoped>
.animate-fadeIn {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { 
    opacity: 0;
    transform: translateY(-10px);
  }
  to { 
    opacity: 1;
    transform: translateY(0);
  }
}

.router-link-active {
  @apply bg-primary-700 bg-opacity-70 dark:bg-primary-700 dark:bg-opacity-40 font-semibold;
}
</style> 