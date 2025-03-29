<template>
  <div class="bg-card dark:bg-gray-800 rounded-lg shadow-sm hover:shadow p-4 border border-border dark:border-gray-700 transition duration-300 hover:border-primary-500 dark:hover:border-primary-400 cursor-pointer group dark:shadow-gray-900/20 mb-2" @click="navigateToNoteDetail">
    <div class="flex flex-col md:flex-row md:items-center">
      <div class="flex-grow">
        <div class="flex items-start justify-between">
          <h3 class="text-lg font-semibold text-text-primary dark:text-gray-50 group-hover:text-primary-600 dark:group-hover:text-primary-300 transition-colors mr-2">{{ note.title }}</h3>
          <div class="flex items-center shrink-0 space-x-2">
            <div class="flex flex-wrap gap-1">
              <span 
                v-for="tag in note.tags" 
                :key="tag.id" 
                class="px-2 py-0.5 bg-primary-100 dark:bg-primary-700 dark:bg-opacity-50 text-primary-800 dark:text-primary-100 text-xs rounded-full inline-flex items-center"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-primary-300 mr-1"></span>
                {{ tag.name }}
              </span>
            </div>
            <span class="text-xs text-text-secondary dark:text-gray-400 whitespace-nowrap">{{ formatDate(note.createdAt) }}</span>
          </div>
        </div>
        <p class="text-text-secondary dark:text-gray-300 text-sm mt-1 line-clamp-1">{{ note.summary || note.content.substring(0, 100) + '...' }}</p>
        
        <div class="flex justify-between items-center mt-2 text-xs">
          <div class="text-text-tertiary dark:text-gray-400">更新于: {{ formatDate(note.updatedAt, true) }}</div>
          <button class="text-primary-600 dark:text-primary-300 hover:text-primary-700 dark:hover:text-primary-200 transition-colors font-medium" @click.stop="navigateToNoteDetail">阅读更多</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
import { useRouter } from 'vue-router';

const props = defineProps({
  note: {
    type: Object,
    required: true
  }
});

const router = useRouter();

const navigateToNoteDetail = () => {
  router.push(`/notes/${props.note.id}`);
};

const formatDate = (dateString, showTime = false) => {
  const date = new Date(dateString);
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