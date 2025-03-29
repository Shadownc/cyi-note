<script setup>
import { RouterView, useRoute } from 'vue-router'
import TheNavbar from './components/TheNavbar.vue'
import { computed } from 'vue';

const route = useRoute();

// 判断当前页面是否为登录或注册页面
const isAuthPage = computed(() => {
  return route.name === 'login' || route.name === 'register';
});
</script>

<template>
  <div class="flex flex-col min-h-screen bg-background text-text-primary dark:bg-dark-background dark:text-dark-text-primary transition-colors duration-300" :class="{ 'overflow-hidden': isAuthPage }">
    <!-- 非登录/注册页面才显示导航栏 -->
    <TheNavbar v-if="!isAuthPage" />
    <main :class="{'flex-grow py-6': !isAuthPage, 'min-h-screen': isAuthPage}">
      <RouterView />
    </main>
    <!-- 非登录/注册页面才显示页脚 -->
    <footer v-if="!isAuthPage" class="bg-neutral-100 text-neutral-700 dark:bg-neutral-800 dark:text-neutral-300 py-4 transition-colors duration-300">
      <div class="container mx-auto text-center text-sm">
        <p>© {{ new Date().getFullYear() }} 知识库 - 一个简单易用的笔记管理工具</p>
      </div>
    </footer>
  </div>
</template>

<style>
/* 全局样式 */
body {
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  margin: 0;
  padding: 0;
}

/* 确保深色模式在整个应用中正确应用 */
html {
  transition: color-scheme 0.5s ease;
}

html.dark {
  color-scheme: dark;
}

/* 登录和注册页面没有滚动条 */
.overflow-hidden {
  overflow: hidden;
}
</style>
