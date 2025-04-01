<template>
  <div class="bg-white dark:bg-gray-800 rounded-xl shadow-md overflow-hidden mb-6">
    <div class="p-5">
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4 flex items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        资源统计
      </h2>
      
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
        <!-- 总文件数 -->
        <div class="relative rounded-lg bg-gradient-to-br from-blue-50 to-blue-100 dark:from-blue-900/20 dark:to-blue-800/20 p-4 border border-blue-200 dark:border-blue-800 overflow-hidden">
          <div class="absolute top-0 right-0 mt-2 mr-2 opacity-10">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <div class="text-sm text-blue-800 dark:text-blue-300 font-medium mb-1">文件总数</div>
          <div class="text-2xl text-blue-900 dark:text-blue-100 font-bold">{{ totalFiles }}</div>
        </div>
        
        <!-- 照片数量 -->
        <div class="relative rounded-lg bg-gradient-to-br from-green-50 to-green-100 dark:from-green-900/20 dark:to-green-800/20 p-4 border border-green-200 dark:border-green-800 overflow-hidden">
          <div class="absolute top-0 right-0 mt-2 mr-2 opacity-10">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
          <div class="text-sm text-green-800 dark:text-green-300 font-medium mb-1">照片数量</div>
          <div class="text-2xl text-green-900 dark:text-green-100 font-bold">{{ imageCount }}</div>
        </div>
        
        <!-- 文档数量 -->
        <div class="relative rounded-lg bg-gradient-to-br from-red-50 to-red-100 dark:from-red-900/20 dark:to-red-800/20 p-4 border border-red-200 dark:border-red-800 overflow-hidden">
          <div class="absolute top-0 right-0 mt-2 mr-2 opacity-10">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-red-600 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
          </div>
          <div class="text-sm text-red-800 dark:text-red-300 font-medium mb-1">文档数量</div>
          <div class="text-2xl text-red-900 dark:text-red-100 font-bold">{{ documentCount }}</div>
        </div>
        
        <!-- 其他文件 -->
        <div class="relative rounded-lg bg-gradient-to-br from-purple-50 to-purple-100 dark:from-purple-900/20 dark:to-purple-800/20 p-4 border border-purple-200 dark:border-purple-800 overflow-hidden">
          <div class="absolute top-0 right-0 mt-2 mr-2 opacity-10">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-purple-600 dark:text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
          <div class="text-sm text-purple-800 dark:text-purple-300 font-medium mb-1">其他文件</div>
          <div class="text-2xl text-purple-900 dark:text-purple-100 font-bold">{{ otherCount }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  attachments: {
    type: Array,
    default: () => []
  }
});

// 总文件数
const totalFiles = computed(() => {
  return props.attachments.length;
});

// 照片数量
const imageCount = computed(() => {
  return props.attachments.filter(att => 
    att.filetype && att.filetype.startsWith('image/')
  ).length;
});

// 文档数量
const documentCount = computed(() => {
  return props.attachments.filter(att => {
    const type = att.filetype || '';
    return type.includes('pdf') || 
      type.includes('word') || 
      type.includes('document') || 
      type.includes('excel') || 
      type.includes('sheet') || 
      type.includes('powerpoint') || 
      type.includes('presentation') ||
      type.startsWith('text/');
  }).length;
});

// 其他文件数量
const otherCount = computed(() => {
  return totalFiles.value - imageCount.value - documentCount.value;
});
</script>

<style scoped>
/* 可以根据需要添加样式 */
</style> 