<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="text-center mb-16">
      <h1 class="text-4xl font-bold text-gray-800 dark:text-gray-50 mb-4 sm:text-5xl">关于<span class="text-blue-600 dark:text-blue-400">知识库</span></h1>
      <p class="text-lg text-gray-600 dark:text-gray-300 max-w-2xl mx-auto mt-4">
        一个简单易用的个人知识管理工具，帮助您整理、分类和管理您的笔记和知识。
      </p>
      <div class="mt-8">
        <router-link 
          :to="isAuthenticated ? '/notes' : '/auth/login'" 
          class="btn btn-primary px-8 py-3 text-lg shadow-lg hover:shadow-xl"
        >
          {{ isAuthenticated ? '查看笔记' : '立即登录' }}
        </router-link>
        <router-link 
          v-if="!isAuthenticated"
          to="/auth/register" 
          class="btn btn-outline ml-4 px-8 py-3 text-lg"
        >
          免费注册
        </router-link>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
      <div class="card text-center p-8 transform transition duration-300 hover:scale-105 hover:border-blue-500">
        <div class="text-blue-600 dark:text-blue-400 mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-3 dark:text-gray-50">Markdown支持</h3>
        <p class="text-gray-600 dark:text-gray-300">
          使用Markdown语法编写笔记，支持标题、列表、代码块等格式，让您的笔记更加丰富多彩。
        </p>
      </div>

      <div class="card text-center p-8 transform transition duration-300 hover:scale-105 hover:border-blue-500">
        <div class="text-blue-600 dark:text-blue-400 mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-3 dark:text-gray-50">标签管理</h3>
        <p class="text-gray-600 dark:text-gray-300">
          通过标签对笔记进行分类，方便查找和整理您的知识，让信息管理更加条理清晰。
        </p>
      </div>

      <div class="card text-center p-8 transform transition duration-300 hover:scale-105 hover:border-blue-500">
        <div class="text-blue-600 dark:text-blue-400 mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-3 dark:text-gray-50">全文搜索</h3>
        <p class="text-gray-600 dark:text-gray-300">
          快速搜索您的笔记内容，找到您需要的信息，高效获取知识，提升工作学习效率。
        </p>
      </div>
    </div>

    <div class="bg-blue-50 dark:bg-blue-900/30 rounded-xl p-8 shadow-inner">
      <div class="max-w-4xl mx-auto">
        <h2 class="text-2xl font-bold text-center mb-6 text-gray-800 dark:text-gray-100">功能亮点</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div class="flex items-start">
            <div class="bg-blue-100 dark:bg-blue-800 rounded-full p-2 mr-4 text-blue-600 dark:text-blue-300">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900 dark:text-gray-100">响应式设计</h3>
              <p class="text-gray-600 dark:text-gray-300 mt-1">自适应PC和移动端，随时随地管理您的笔记</p>
            </div>
          </div>
          <div class="flex items-start">
            <div class="bg-blue-100 dark:bg-blue-800 rounded-full p-2 mr-4 text-blue-600 dark:text-blue-300">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900 dark:text-gray-100">实时预览</h3>
              <p class="text-gray-600 dark:text-gray-300 mt-1">编辑Markdown时实时查看渲染效果</p>
            </div>
          </div>
          <div class="flex items-start">
            <div class="bg-blue-100 dark:bg-blue-800 rounded-full p-2 mr-4 text-blue-600 dark:text-blue-300">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900 dark:text-gray-100">分类管理</h3>
              <p class="text-gray-600 dark:text-gray-300 mt-1">通过标签系统对笔记进行分类整理</p>
            </div>
          </div>
          <div class="flex items-start">
            <div class="bg-blue-100 dark:bg-blue-800 rounded-full p-2 mr-4 text-blue-600 dark:text-blue-300">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-gray-900 dark:text-gray-100">代码高亮</h3>
              <p class="text-gray-600 dark:text-gray-300 mt-1">支持多种编程语言的代码块语法高亮</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-16 text-center">
      <h2 class="text-2xl font-bold text-gray-800 dark:text-gray-100 mb-4">立即开始使用</h2>
      <p class="text-gray-600 dark:text-gray-300 mb-6 max-w-2xl mx-auto">
        不论是记录工作笔记、整理学习资料，还是管理个人知识体系，知识库都能满足您的需求。
      </p>
      <router-link 
        :to="isAuthenticated ? '/notes/new' : '/auth/register'" 
        class="btn btn-primary px-8 py-3 text-lg shadow-lg hover:shadow-xl"
      >
        {{ isAuthenticated ? '创建笔记' : '免费注册' }}
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useUserStore } from '@/stores/user';

// 获取用户认证状态
const userStore = useUserStore();
const isAuthenticated = computed(() => userStore.isAuthenticated);
</script>

<style scoped>
.card {
  @apply bg-white dark:bg-gray-800 rounded-lg shadow-md hover:shadow-lg border border-gray-100 dark:border-gray-700 transition duration-300;
}

.btn {
  @apply inline-block text-center rounded-md transition-colors duration-300 font-medium;
}

.btn-primary {
  @apply bg-primary-600 hover:bg-primary-700 text-white;
}

.btn-outline {
  @apply border border-primary-600 text-primary-600 hover:bg-primary-600 hover:bg-opacity-10;
}
</style>
