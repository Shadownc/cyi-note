import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 初始化主题
function initializeTheme() {
  // 从localStorage读取主题设置
  const savedTheme = localStorage.getItem('theme') || 'blue';
  const savedDarkMode = localStorage.getItem('darkMode');
  
  // 添加主题类
  document.documentElement.classList.add(`theme-${savedTheme}`);
  
  // 应用暗色模式
  if (savedDarkMode === 'dark' || 
      (!savedDarkMode && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
}

// 在应用挂载前初始化主题
initializeTheme();

app.mount('#app')
