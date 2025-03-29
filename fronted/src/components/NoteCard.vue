<template>
  <div class="bg-card dark:bg-gray-800 rounded-lg shadow-md hover:shadow-lg p-5 border border-border dark:border-gray-700 transition duration-300 hover:border-primary-500 dark:hover:border-primary-400 cursor-pointer group dark:shadow-gray-900/30" @click="navigateToNoteDetail">
    <div class="flex justify-between items-start mb-3">
      <h3 class="text-lg font-semibold text-text-primary dark:text-gray-50 group-hover:text-primary-600 dark:group-hover:text-primary-300 transition-colors">{{ note.title }}</h3>
      <span class="text-sm text-text-secondary dark:text-gray-400 ml-2 shrink-0">{{ formatDate(note.created_at) }}</span>
    </div>
    <p class="text-text-secondary dark:text-gray-300 text-sm mb-4 line-clamp-2">{{ notePreview }}</p>
    <div class="flex flex-wrap gap-2">
      <span 
        v-for="tag in note.tags" 
        :key="tag.id" 
        class="px-2 py-1 bg-primary-100 dark:bg-primary-700 dark:bg-opacity-50 text-primary-800 dark:text-primary-100 text-xs rounded-full inline-flex items-center"
      >
        <span class="w-2 h-2 rounded-full bg-primary-600 dark:bg-primary-300 mr-1"></span>
        {{ tag.name }}
      </span>
    </div>
    <div class="flex justify-between items-center mt-4 pt-3 border-t border-border dark:border-gray-700">
      <div class="text-xs text-text-tertiary dark:text-gray-400">更新于: {{ formatDate(note.updated_at, true) }}</div>
      <button class="text-primary-600 dark:text-primary-300 text-sm hover:text-primary-700 dark:hover:text-primary-200 transition-colors font-medium" @click.stop="navigateToNoteDetail">阅读更多</button>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
import { useRouter } from 'vue-router';
import { computed } from 'vue';

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

const formatDate = (dateString, showTime = false) => {
  console.log('NoteCard 格式化日期，原始值:', dateString);
  
  if (!dateString) {
    console.warn('日期字符串为空或无效');
    return '无日期';
  }
  
  // 尝试修复日期格式
  let fixedDateString = dateString;
  if (typeof dateString === 'string' && dateString.includes('T')) {
    // 确保日期字符串有Z或时区
    if (!dateString.includes('Z') && !dateString.includes('+')) {
      fixedDateString = dateString + 'Z';
    }
  }
  
  const date = new Date(fixedDateString);
  console.log('NoteCard 转换后的日期对象:', date);
  
  if (isNaN(date.getTime())) {
    console.warn('日期解析失败:', dateString);
    return '日期无效';
  }
  
  try {
    if (showTime) {
      return date.toLocaleDateString('zh-CN', { 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    });
  } catch (e) {
    console.error('日期格式化错误:', e);
    
    // 降级方案：手动格式化
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    
    if (showTime) {
      return `${month}月${day}日 ${hours}:${minutes}`;
    }
    return `${year}年${month}月${day}日`;
  }
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