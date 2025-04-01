<template>
  <div class="bg-white dark:bg-gray-800 rounded-xl shadow-md overflow-hidden transition-all duration-300">
    <div v-if="$slots.title || $slots.toolbar" class="border-b border-gray-100 dark:border-gray-700">
      <div class="p-4">
        <!-- 头部区域：包含标题和工具栏 -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <!-- 标题插槽 -->
          <div v-if="$slots.title" class="flex-shrink-0">
            <slot name="title"></slot>
          </div>
          
          <!-- 工具栏插槽 -->
          <div v-if="$slots.toolbar" class="flex-grow">
            <slot name="toolbar"></slot>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 内容区域 -->
    <div class="p-4">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
// 无需定义任何Props，使用插槽进行内容分发
</script>

<style scoped>
/* 卡片悬停效果 */
.card-hover:hover {
  @apply shadow-lg;
  transform: translateY(-2px);
}

/* 暗色模式下卡片效果 */
.dark .bg-gray-800 {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2), 0 2px 4px -1px rgba(0, 0, 0, 0.1);
}
</style> 