<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <RouterView />
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

/* 自定义动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fadeIn {
  animation: fadeIn 0.4s ease forwards;
}

.animate-slideUp {
  animation: slideUp 0.5s ease forwards;
}
</style>
