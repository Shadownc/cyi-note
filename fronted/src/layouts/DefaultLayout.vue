<script setup>
import TheNavbar from '@/components/TheNavbar.vue';
</script>

<template>
  <div class="flex flex-col min-h-screen bg-background text-text-primary dark:bg-dark-background dark:text-dark-text-primary transition-colors duration-300">
    <!-- 头部导航栏 -->
    <TheNavbar />
    
    <!-- 主内容区 -->
    <main class="flex-grow py-6">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8">
        <!-- 确保内容插槽正确显示 -->
        <router-view />
      </div>
    </main>
    
    <!-- 页脚 -->
    <footer class="bg-neutral-100 text-neutral-700 dark:bg-neutral-800 dark:text-neutral-300 py-4 transition-colors duration-300">
      <div class="container mx-auto text-center text-sm px-4">
        <p>© {{ new Date().getFullYear() }} 知识库 - 一个简单易用的笔记管理工具</p>
      </div>
    </footer>
  </div>
</template> 