<template>
  <div class="w-full">
    <!-- 笔记标题和操作按钮 -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4 gap-3">
      <h1 class="text-xl md:text-2xl font-bold text-text-primary dark:text-gray-50 break-words">{{ note.title }}</h1>
      <div class="flex items-center space-x-2 shrink-0">
        <button @click="$emit('edit')" class="btn btn-secondary flex items-center text-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
          编辑
        </button>
        <button @click="$emit('delete')" class="btn flex items-center bg-red-600 dark:bg-red-700 text-white hover:bg-red-700 dark:hover:bg-red-800 transition-colors text-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
          删除
        </button>
      </div>
    </div>

    <!-- 笔记元信息 -->
    <div class="flex flex-wrap items-center text-xs sm:text-sm text-text-tertiary dark:text-gray-400 mb-4">
      <span>创建于: {{ formatDate(note.createdAt) }}</span>
      <span class="mx-2">•</span>
      <span>更新于: {{ formatDate(note.updatedAt, true) }}</span>
    </div>

    <!-- 笔记标签 -->
    <div class="flex flex-wrap gap-2 mb-4 sm:mb-6">
      <span 
        v-for="tag in note.tags" 
        :key="tag.id" 
        class="px-2 py-1 bg-primary-100 dark:bg-primary-700 dark:bg-opacity-50 text-primary-800 dark:text-primary-100 text-xs rounded-full"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-primary-300 mr-1 inline-block"></span>
        {{ tag.name }}
      </span>
    </div>

    <!-- 笔记内容 - Markdown渲染 -->
    <div class="prose dark:prose-invert sm:prose-sm md:prose-base max-w-none overflow-hidden dark:bg-gray-800 dark:bg-opacity-30 dark:p-4 dark:rounded-lg" v-html="renderedContent"></div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, onMounted, watch } from 'vue';
import { marked } from 'marked';

const props = defineProps({
  note: {
    type: Object,
    required: true
  }
});

defineEmits(['edit', 'delete']);

const renderedContent = ref('');

// 格式化日期
const formatDate = (dateString, shortFormat = false) => {
  const date = new Date(dateString);
  
  // 在移动端使用更简短的日期格式
  if (shortFormat && window.innerWidth < 640) {
    return date.toLocaleDateString('zh-CN', { 
      month: 'numeric', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  
  return date.toLocaleDateString('zh-CN', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 渲染Markdown内容
const renderContent = () => {
  if (props.note && props.note.content) {
    renderedContent.value = marked(props.note.content);
  }
};

// 监听props变化重新渲染
watch(() => props.note, () => {
  renderContent();
}, { deep: true });

// 组件加载时渲染内容
onMounted(() => {
  renderContent();
});
</script>

<style scoped>
.prose {
  @apply text-text-primary transition-colors;
}

.dark .prose {
  @apply text-gray-100;
}

/* 增强Markdown内容在暗色模式下的可读性 */
:deep(.dark .prose) {
  --tw-prose-headings: theme('colors.gray.100');
  --tw-prose-body: theme('colors.gray.200');
  --tw-prose-bold: theme('colors.gray.50');
  --tw-prose-links: theme('colors.primary.300');
  --tw-prose-code: theme('colors.gray.100');
  --tw-prose-quotes: theme('colors.gray.300');
}

:deep(.prose) {
  --tw-prose-links: theme('colors.primary.600');
}

/* 增强暗色模式下的样式 */
:deep(.dark .prose blockquote) {
  @apply border-l-4 border-primary-700 pl-4 bg-gray-800 bg-opacity-50 py-1 rounded-r-md;
}

:deep(.dark .prose h1) {
  @apply text-primary-300 border-b pb-2 border-gray-700;
}

:deep(.dark .prose h2) {
  @apply text-primary-300 opacity-90 border-b pb-1 border-gray-700;
}

:deep(.dark .prose h3) {
  @apply text-primary-300 opacity-80;
}

:deep(.dark .prose strong) {
  @apply text-primary-200;
}

:deep(.dark .prose code) {
  @apply bg-gray-800 text-primary-200 border border-gray-700;
}

:deep(.dark .prose pre) {
  @apply bg-gray-900 border border-gray-700;
}

:deep(.dark .prose pre code) {
  @apply bg-transparent border-0 text-gray-200;
}

:deep(.dark .prose a) {
  @apply text-primary-300 hover:text-primary-200 underline;
}

:deep(.dark .prose table) {
  @apply border-collapse border border-gray-700;
}

:deep(.dark .prose th) {
  @apply bg-gray-800 text-gray-100 border border-gray-700 px-4 py-2;
}

:deep(.dark .prose td) {
  @apply border border-gray-700 px-4 py-2;
}

/* 移动端优化 Markdown 内容 */
@media (max-width: 640px) {
  :deep(.prose img) {
    @apply w-full;
  }
  
  :deep(.prose pre) {
    @apply text-sm;
  }
}
</style> 