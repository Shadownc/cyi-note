<template>
  <div class="relative">
    <!-- 主题设置按钮 -->
    <button 
      @click="toggleMenu"
      class="rounded-full p-2 text-white dark:text-white transition-colors focus:outline-none hover:bg-primary-800 hover:bg-opacity-50 hover:text-white"
      :title="t('主题设置')"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path v-if="!darkMode" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
        <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
      </svg>
    </button>

    <!-- 主题设置下拉菜单 -->
    <div 
      v-show="menuOpen"
      class="absolute right-0 mt-2 w-64 bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4 z-50 animate-fadeIn"
    >
      <div class="border-b dark:border-gray-700 pb-3 mb-3">
        <h3 class="text-gray-700 dark:text-gray-200 font-medium mb-2">{{ t('显示模式') }}</h3>
        <div class="flex items-center gap-4">
          <button
            @click="setDarkMode(false)"
            class="flex-1 px-3 py-2 rounded-md text-center text-sm transition-colors"
            :class="darkMode ? 'text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700' : 'bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200'"
          >
            <div class="flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
              {{ t('亮色') }}
            </div>
          </button>
          <button
            @click="setDarkMode(true)"
            class="flex-1 px-3 py-2 rounded-md text-center text-sm transition-colors"
            :class="darkMode ? 'bg-primary-800 text-primary-100 dark:bg-primary-700 dark:text-primary-50' : 'text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'"
          >
            <div class="flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
              {{ t('暗色') }}
            </div>
          </button>
        </div>
      </div>

      <div>
        <h3 class="text-gray-700 dark:text-gray-200 font-medium mb-2">{{ t('主题颜色') }}</h3>
        <div class="grid grid-cols-5 gap-2">
          <button
            v-for="(theme, index) in themes"
            :key="index"
            @click="setTheme(theme.value)"
            class="h-8 w-full rounded-md transition-transform hover:scale-110 focus:outline-none transform duration-200"
            :style="{backgroundColor: theme.color}"
            :class="{'ring-2 ring-offset-2 dark:ring-offset-gray-800': currentTheme === theme.value}"
            :title="theme.label"
          ></button>
        </div>
      </div>
    </div>

    <!-- 点击其他区域关闭菜单 -->
    <div 
      v-if="menuOpen" 
      class="fixed inset-0 z-40"
      @click="menuOpen = false"
    ></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';

// 支持的主题列表
const themes = [
  { value: 'blue', label: '蓝色', color: '#1a84ff' },
  { value: 'purple', label: '紫色', color: '#7c3aed' },
  { value: 'green', label: '绿色', color: '#059669' },
  { value: 'amber', label: '琥珀色', color: '#d97706' },
  { value: 'red', label: '红色', color: '#dc2626' },
];

const menuOpen = ref(false);
const darkMode = ref(false);
const currentTheme = ref('blue');

// 简单的i18n支持
const t = (text) => text; // 这里可以替换为实际的翻译函数

// 打开/关闭菜单
const toggleMenu = () => {
  menuOpen.value = !menuOpen.value;
};

// 设置暗色模式
const setDarkMode = (isDark) => {
  darkMode.value = isDark;
  localStorage.setItem('darkMode', isDark ? 'dark' : 'light');
  
  if (isDark) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
};

// 设置主题颜色
const setTheme = (theme) => {
  // 移除所有主题类
  document.documentElement.classList.remove(
    ...themes.map(t => `theme-${t.value}`)
  );
  
  // 添加新主题类
  document.documentElement.classList.add(`theme-${theme}`);
  
  // 保存到 localStorage
  localStorage.setItem('theme', theme);
  currentTheme.value = theme;
};

// 初始化主题
onMounted(() => {
  // 读取保存的深色模式设置
  const savedDarkMode = localStorage.getItem('darkMode');
  if (savedDarkMode === 'dark' || (!savedDarkMode && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    setDarkMode(true);
  }
  
  // 读取保存的主题
  const savedTheme = localStorage.getItem('theme') || 'blue';
  setTheme(savedTheme);
  
  // 监听系统主题变化
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (!localStorage.getItem('darkMode')) {
      setDarkMode(e.matches);
    }
  });
});

// 关闭菜单的ESC键监听
watch(menuOpen, (isOpen) => {
  if (isOpen) {
    const handleEsc = (e) => {
      if (e.key === 'Escape') {
        menuOpen.value = false;
      }
    };
    document.addEventListener('keydown', handleEsc);
    return () => {
      document.removeEventListener('keydown', handleEsc);
    };
  }
});
</script>

<style scoped>
.animate-fadeIn {
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { 
    opacity: 0;
    transform: translateY(-10px) scale(0.95);
  }
  to { 
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style> 