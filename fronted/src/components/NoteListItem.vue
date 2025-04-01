<template>
  <div class="bg-card dark:bg-gray-800 rounded-lg shadow-sm hover:shadow p-4 border border-border dark:border-gray-700 transition duration-300 hover:border-primary-500 dark:hover:border-primary-400 cursor-pointer group dark:shadow-gray-900/20 mb-2" @click="navigateToNoteDetail">
    <div class="flex flex-col md:flex-row md:items-center">
      <div class="flex-grow">
        <div class="flex items-start justify-between">
          <h3 class="text-lg font-semibold text-text-primary dark:text-gray-50 group-hover:text-primary-600 dark:group-hover:text-primary-300 transition-colors mr-2 truncate">{{ note.title }}</h3>
          <div class="flex items-center shrink-0 space-x-2">
            <div class="flex flex-wrap gap-1">
              <span 
                v-for="tag in note.tags" 
                :key="tag.id" 
                class="px-2 py-0.5 bg-primary-100 dark:bg-primary-600 text-primary-800 dark:text-white text-xs rounded-full inline-flex items-center shadow-sm"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-white mr-1 inline-block"></span>
                {{ tag.name }}
              </span>
            </div>
            <span class="text-xs text-text-secondary dark:text-gray-400 whitespace-nowrap">{{ formatDate(note.created_at) }}</span>
          </div>
        </div>
        <p class="text-text-secondary dark:text-gray-300 text-sm mt-1 line-clamp-1">{{ notePreview }}</p>
        
        <div class="flex justify-between items-center mt-2 text-xs">
          <div class="text-text-tertiary dark:text-gray-400">更新于: {{ formatDate(note.updated_at, true) }}</div>
          <button class="text-primary-600 dark:text-primary-300 hover:text-primary-700 dark:hover:text-primary-200 transition-colors font-medium" @click.stop="navigateToNoteDetail">阅读更多</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
import { useRouter } from 'vue-router';
import { computed } from 'vue';
import { formatDate } from '@/utils/dateUtils';

const props = defineProps({
  note: {
    type: Object,
    required: true
  }
});

const router = useRouter();

// 计算属性：笔记摘要，优先使用summary字段，如无则转换content
const notePreview = computed(() => {
  if (props.note.summary) return props.note.summary;
  
  const plainText = stripMarkdown(props.note.content || '');
  if (plainText.length <= 100) return plainText;
  return plainText.substring(0, 100) + '...';
});

// 移除Markdown格式，转换为纯文本
const stripMarkdown = (markdownText) => {
  if (!markdownText) return '';
  
  // 替换代码块
  let plainText = markdownText.replace(/```[\s\S]*?```/g, '[代码块]');
  
  // 替换内联代码
  plainText = plainText.replace(/`([^`]+)`/g, '$1');
  
  // 替换标题
  plainText = plainText.replace(/#{1,6}\s+/g, '');
  
  // 替换链接 [text](url) 为 text
  plainText = plainText.replace(/\[([^\]]+)\]\([^)]+\)/g, '$1');
  
  // 替换图片 ![alt](url) 为 [图片]
  plainText = plainText.replace(/!\[([^\]]*)\]\([^)]+\)/g, '[图片]');
  
  // 替换加粗和斜体
  plainText = plainText.replace(/(\*\*|__)(.*?)\1/g, '$2');
  plainText = plainText.replace(/(\*|_)(.*?)\1/g, '$2');
  
  // 替换列表
  plainText = plainText.replace(/^\s*[\-*+]\s+/gm, '');
  plainText = plainText.replace(/^\s*\d+\.\s+/gm, '');
  
  return plainText;
};

const navigateToNoteDetail = () => {
  router.push(`/notes/${props.note.id}`);
};
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 暗色模式下列表项悬停效果增强 */
.dark .group:hover {
  @apply bg-gray-750 shadow-md;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2), 0 2px 4px -1px rgba(0, 0, 0, 0.1);
}
</style> 