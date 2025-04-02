<template>
  <div class="bg-card dark:bg-gray-800 rounded-lg shadow-md hover:shadow-lg p-5 border border-border dark:border-gray-700 transition duration-300 hover:border-primary-500 dark:hover:border-primary-400 cursor-pointer group dark:shadow-gray-900/30 h-[230px] flex flex-col" @click="navigateToNoteDetail">
    <div class="flex justify-between items-start mb-3">
      <h3 class="text-lg font-semibold text-text-primary dark:text-gray-50 group-hover:text-primary-600 dark:group-hover:text-primary-300 transition-colors truncate mr-2">
        {{ note.title }}
        <span v-if="note.is_public" class="inline-flex items-center text-green-600 dark:text-green-400 ml-2 text-xs">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 mr-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          公开
        </span>
      </h3>
      <span class="text-sm text-text-secondary dark:text-gray-400 ml-2 shrink-0">{{ formatDate(note.created_at) }}</span>
    </div>
    <p class="text-text-secondary dark:text-gray-300 text-sm mb-4 line-clamp-2 flex-grow">{{ notePreview }}</p>
    <div class="flex flex-wrap gap-2 mb-2">
      <span 
        v-for="tag in note.tags" 
        :key="tag.id" 
        class="px-2 py-1 bg-primary-100 dark:bg-primary-600 text-primary-800 dark:text-white text-xs rounded-full inline-flex items-center shadow-sm"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-white mr-1 inline-block"></span>
        {{ tag.name }}
      </span>
    </div>
    <div class="flex justify-between items-center mt-auto pt-3 border-t border-border dark:border-gray-700">
      <div class="text-xs text-text-tertiary dark:text-gray-400">更新于: {{ formatDate(note.updated_at, true) }}</div>
      <button class="text-primary-600 dark:text-primary-300 text-sm hover:text-primary-700 dark:hover:text-primary-200 transition-colors font-medium" @click.stop="navigateToNoteDetail">阅读更多</button>
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

const navigateToNoteDetail = () => {
  router.push(`/notes/${props.note.id}`);
};

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
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 暗色模式下卡片悬停效果增强 */
.dark .group:hover {
  @apply bg-gray-750 shadow-lg;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.2), 0 4px 6px -2px rgba(0, 0, 0, 0.1);
}
</style> 