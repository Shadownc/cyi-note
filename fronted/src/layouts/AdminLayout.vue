<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';

const userStore = useUserStore();
const route = useRoute();
const router = useRouter();
const isSidebarOpen = ref(true);
const currentSection = ref('users');

// 响应式断点检测
const isSmallScreen = ref(false);
const checkScreenSize = () => {
  isSmallScreen.value = window.innerWidth < 768;
  
  // 在小屏幕上自动折叠侧边栏
  if (isSmallScreen.value) {
    isSidebarOpen.value = false;
  }
};

// 从路由路径获取当前section
const updateSectionFromRoute = () => {
  const path = route.path;
  if (path.includes('/admin/users')) {
    currentSection.value = 'users';
  } else if (path.includes('/admin/notes')) {
    currentSection.value = 'notes';
  } else if (path.includes('/admin/tags')) {
    currentSection.value = 'tags';
  } else if (path.includes('/admin/settings')) {
    currentSection.value = 'settings';
  }
};

// 切换侧边栏
const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

// 切换管理版块
const setCurrentSection = (section) => {
  currentSection.value = section;
  router.push(`/admin/${section}`);
  
  // 在移动端切换版块后自动折叠侧边栏
  if (isSmallScreen.value) {
    isSidebarOpen.value = false;
  }
};

// 退出登录
const logout = () => {
  userStore.logout();
  router.push('/login');
};

// 监听路由变化
watch(() => route.path, updateSectionFromRoute);

// 监听屏幕大小变化
onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
  updateSectionFromRoute();
});

// 组件卸载前移除事件监听器
onBeforeUnmount(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>

<template>
  <div class="flex h-screen bg-gray-100 dark:bg-gray-900 overflow-hidden">
    <!-- 侧边栏遮罩层 (移动端) -->
    <div 
      v-if="isSmallScreen && isSidebarOpen" 
      class="fixed inset-0 bg-black/50 z-20"
      @click="toggleSidebar"
    ></div>
    
    <!-- 侧边栏 -->
    <aside 
      :class="[
        'transition-all duration-300 ease-in-out',
        'bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700',
        'h-screen z-30 overflow-y-auto',
        isSidebarOpen ? 'w-56' : 'w-0 md:w-14',
        isSmallScreen && isSidebarOpen ? 'fixed inset-y-0 left-0 shadow-2xl' : '',
      ]"
    >
      <!-- 顶部标志区域 -->
      <div :class="[
        'flex items-center px-3 h-14 border-b border-gray-200 dark:border-gray-700',
        isSidebarOpen ? 'justify-between' : 'justify-center'
      ]">
        <div v-if="isSidebarOpen" class="flex items-center space-x-2">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-r from-purple-600 to-indigo-600 flex items-center justify-center text-white font-bold shadow-lg">
            A
          </div>
          <span class="font-bold text-gray-900 dark:text-white">管理后台</span>
        </div>
        <div v-else class="w-8 h-8 rounded-full bg-gradient-to-r from-purple-600 to-indigo-600 flex items-center justify-center text-white font-bold shadow-lg">
          A
        </div>
        
        <!-- 侧边栏折叠按钮 -->
        <button v-if="isSidebarOpen" @click="toggleSidebar" class="p-1 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
          </svg>
        </button>
      </div>
      
      <!-- 折叠状态下的展开按钮 -->
      <button v-if="!isSidebarOpen" @click="toggleSidebar" class="mt-2 w-full flex justify-center p-1 hover:bg-gray-100 dark:hover:bg-gray-700">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7" />
        </svg>
      </button>
      
      <!-- 用户信息 (移动到侧边栏顶部) -->
      <div v-if="isSidebarOpen" class="border-b border-gray-200 dark:border-gray-700 py-2 px-3 bg-white dark:bg-gray-800">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 rounded-full bg-purple-100 dark:bg-purple-900 flex items-center justify-center">
              <span class="text-purple-600 dark:text-purple-400 font-bold">
                {{ userStore.user?.username?.[0]?.toUpperCase() || 'A' }}
              </span>
            </div>
          </div>
          <div class="ml-2">
            <p class="text-xs font-medium text-gray-900 dark:text-white">{{ userStore.user?.username || '管理员' }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ userStore.user?.email || 'admin@example.com' }}</p>
          </div>
        </div>
      </div>
      
      <!-- 导航菜单 -->
      <nav class="p-2">
        <ul class="space-y-1">
          <!-- 返回首页 -->
          <li>
            <router-link 
              to="/" 
              :class="[
                'w-full flex items-center p-2 rounded-lg transition-all',
                'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50',
                !isSidebarOpen ? 'justify-center' : '',
              ]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
              </svg>
              <span v-if="isSidebarOpen" class="ml-2 text-sm">返回首页</span>
            </router-link>
          </li>
          
          <!-- 用户管理 -->
          <li>
            <button 
              @click="setCurrentSection('users')" 
              :class="[
                'w-full flex items-center p-2 rounded-lg transition-all',
                currentSection === 'users' 
                  ? 'bg-purple-50 dark:bg-purple-900/20 text-purple-600 dark:text-purple-400' 
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50',
                !isSidebarOpen ? 'justify-center' : '',
              ]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
              <span v-if="isSidebarOpen" class="ml-2 text-sm">用户管理</span>
            </button>
          </li>
          
          <!-- 笔记管理 -->
          <li>
            <button 
              @click="setCurrentSection('notes')" 
              :class="[
                'w-full flex items-center p-2 rounded-lg transition-all',
                currentSection === 'notes' 
                  ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400' 
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50',
                !isSidebarOpen ? 'justify-center' : '',
              ]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
              <span v-if="isSidebarOpen" class="ml-2 text-sm">笔记管理</span>
            </button>
          </li>
          
          <!-- 标签管理 -->
          <li>
            <button 
              @click="setCurrentSection('tags')" 
              :class="[
                'w-full flex items-center p-2 rounded-lg transition-all',
                currentSection === 'tags' 
                  ? 'bg-orange-50 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400' 
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50',
                !isSidebarOpen ? 'justify-center' : '',
              ]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
              </svg>
              <span v-if="isSidebarOpen" class="ml-2 text-sm">标签管理</span>
            </button>
          </li>
          
          <!-- 系统设置 -->
          <li>
            <button 
              @click="setCurrentSection('settings')" 
              :class="[
                'w-full flex items-center p-2 rounded-lg transition-all',
                currentSection === 'settings' 
                  ? 'bg-green-50 dark:bg-green-900/20 text-green-600 dark:text-green-400' 
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50',
                !isSidebarOpen ? 'justify-center' : '',
              ]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              </svg>
              <span v-if="isSidebarOpen" class="ml-2 text-sm">系统设置</span>
            </button>
          </li>
        </ul>
        
        <!-- 分隔线 -->
        <div v-if="isSidebarOpen" class="my-3 border-t border-gray-200 dark:border-gray-700"></div>
        
        <!-- 底部菜单 -->
        <ul v-if="isSidebarOpen" class="space-y-1">
          <!-- 退出登录 -->
          <li>
            <button 
              @click="logout"
              class="w-full flex items-center p-2 text-red-600 dark:text-red-400 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 text-sm"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              <span class="ml-2">退出登录</span>
            </button>
          </li>
        </ul>
      </nav>
    </aside>
    
    <!-- 主内容区 -->
    <main :class="[
      'flex-1 overflow-y-auto',
      'transition-all duration-300 ease-in-out',
      isSmallScreen && !isSidebarOpen ? 'w-full' : '',
      !isSmallScreen && !isSidebarOpen ? 'ml-14' : '',
    ]">
      <!-- 切换侧边栏按钮（仅在移动端显示） -->
      <div v-if="isSmallScreen" class="bg-white dark:bg-gray-800 p-3 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-10 flex items-center">
        <button @click="toggleSidebar" class="p-1.5 mr-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <h1 class="text-base font-medium text-gray-800 dark:text-white">
          {{ 
            currentSection === 'users' ? '用户管理' : 
            currentSection === 'notes' ? '笔记管理' : 
            currentSection === 'tags' ? '标签管理' : '系统设置' 
          }}
        </h1>
      </div>
      
      <!-- 内容区域 -->
      <div class="p-4 w-full">
        <router-view />
      </div>
    </main>
  </div>
</template>

<style scoped>
/* 滚动条样式 */
::-webkit-scrollbar {
  width: 3px;
  height: 3px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(193, 193, 193, 0.3);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(193, 193, 193, 0.5);
}

.dark ::-webkit-scrollbar-thumb {
  background: rgba(90, 90, 90, 0.3);
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(90, 90, 90, 0.5);
}
</style> 