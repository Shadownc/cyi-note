<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-white mb-4">公开笔记</h1>
      <p class="text-gray-600 dark:text-gray-300">社区分享的精选笔记</p>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 dark:border-primary-300"></div>
    </div>

    <!-- 无笔记状态 -->
    <div v-else-if="publicNotes.length === 0" class="text-center my-12">
      <div class="text-gray-500 dark:text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto mb-4 text-gray-400 dark:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-lg mb-4 dark:text-gray-100">暂无公开笔记</p>
      </div>
    </div>

    <!-- 笔记内容直接展示 -->
    <div v-else class="grid grid-cols-1 gap-10">
      <div v-for="note in publicNotes" :key="note.id" class="bg-white dark:bg-gray-800 rounded-xl shadow-md overflow-hidden border border-gray-200 dark:border-gray-700 transition-all duration-300 hover:shadow-lg">
        <!-- 笔记顶部信息栏 -->
        <div class="p-6 bg-gradient-to-r from-primary-50 to-indigo-50 dark:from-primary-900/30 dark:to-indigo-900/30 border-b border-gray-200 dark:border-gray-700">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <h2 class="text-2xl font-bold text-gray-800 dark:text-white">{{ note.title }}</h2>
            
            <!-- 作者与时间信息 -->
            <div class="flex items-center gap-3">
              <div class="flex items-center space-x-2">
                <div class="w-8 h-8 rounded-full bg-primary-500 dark:bg-primary-600 flex items-center justify-center text-white font-medium">
                  {{ getInitials(note.user?.name || '匿名') }}
                </div>
                <div>
                  <div class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ note.user?.name || '匿名用户' }}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">{{ formatDate(note.created_at) }}</div>
                </div>
              </div>
              
              <router-link :to="`/notes/${note.id}`" class="ml-auto bg-primary-500 hover:bg-primary-600 text-white px-4 py-2 rounded-full text-sm font-medium shadow-sm transition-colors">
                查看详情
              </router-link>
            </div>
          </div>
          
          <!-- 标签展示 -->
          <div class="flex flex-wrap gap-2 mt-4">
            <span 
              v-for="tag in note.tags" 
              :key="tag.id" 
              class="px-3 py-1 bg-primary-100 dark:bg-primary-600/70 text-primary-800 dark:text-white text-xs rounded-full inline-flex items-center shadow-sm"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-white mr-1.5 inline-block"></span>
              {{ typeof tag === 'string' ? tag : tag.name }}
            </span>
          </div>
        </div>
        
        <!-- 笔记内容预览 -->
        <div class="p-6">
          <!-- 笔记内容 -->
          <div class="prose dark:prose-invert max-w-none mb-4" :class="{'line-clamp-5': !expandedNotes[note.id]}">
            <MarkdownRenderer :content="note.content.slice(0, 1000) + (note.content.length > 1000 ? '...' : '')" />
          </div>
          
          <!-- 展开/折叠按钮 -->
          <div class="flex justify-center mt-4" v-if="note.content.length > 600">
            <button 
              @click="toggleExpand(note.id)" 
              class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 text-sm font-medium flex items-center"
            >
              <span v-if="expandedNotes[note.id]">
                收起
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                </svg>
              </span>
              <span v-else>
                展开阅读更多
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </span>
            </button>
          </div>
        </div>
        
        <!-- 笔记底部 -->
        <div class="px-6 py-4 bg-gray-50 dark:bg-gray-800/60 border-t border-gray-200 dark:border-gray-700 flex justify-between items-center">
          <div class="text-gray-500 dark:text-gray-400 text-sm">
            <span>更新于: {{ formatDate(note.updated_at, true) }}</span>
          </div>
          <div class="flex space-x-4">
            <router-link :to="`/notes/${note.id}`" class="text-primary-600 dark:text-primary-300 hover:text-primary-700 dark:hover:text-primary-200 text-sm font-medium">
              查看全文
            </router-link>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 分页控制 -->
    <div v-if="totalNotes > pageSize && !isLoading" class="flex justify-center mt-8">
      <div class="flex items-center space-x-2">
        <button 
          @click="goToPage(1)" 
          class="p-2 rounded-md border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="currentPage === 1"
          title="第一页"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
          </svg>
        </button>
        <button 
          @click="goToPage(currentPage - 1)" 
          class="p-2 rounded-md border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="currentPage === 1"
          title="上一页"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        
        <div class="flex items-center">
          <template v-for="page in pageButtons" :key="page">
            <button 
              v-if="page !== '...'" 
              @click="goToPage(page)" 
              class="px-3 py-1 mx-0.5 rounded-md transition-colors"
              :class="currentPage === page ? 'bg-primary-600 dark:bg-primary-500 text-white' : 'text-gray-600 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700'"
            >
              {{ page }}
            </button>
            <span v-else class="px-1 text-gray-600 dark:text-gray-300 mx-0.5">...</span>
          </template>
        </div>
        
        <button 
          @click="goToPage(currentPage + 1)" 
          class="p-2 rounded-md border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="currentPage === totalPages"
          title="下一页"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
        <button 
          @click="goToPage(totalPages)" 
          class="p-2 rounded-md border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="currentPage === totalPages"
          title="最后一页"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useNotesStore } from '@/stores/notes';
import { storeToRefs } from 'pinia';
import MarkdownRenderer from '@/components/MarkdownRenderer.vue';

// 使用store
const notesStore = useNotesStore();
const { publicNotes, totalPublicNotes, isLoading, publicNotesPage } = storeToRefs(notesStore);

// 本地状态
const pageSize = ref(5); // 减少每页显示的笔记数量，因为现在直接显示内容
const currentPage = computed(() => publicNotesPage.value);
const totalNotes = computed(() => totalPublicNotes.value);
const expandedNotes = reactive({}); // 跟踪哪些笔记是展开的

// 计算要显示的页码按钮
const pageButtons = computed(() => {
  const total = totalPages.value;
  const current = currentPage.value;
  
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }
  
  if (current <= 3) {
    return [1, 2, 3, 4, '...', total - 1, total];
  }
  
  if (current >= total - 2) {
    return [1, 2, '...', total - 3, total - 2, total - 1, total];
  }
  
  return [1, '...', current - 1, current, current + 1, '...', total];
});

// 总页数
const totalPages = computed(() => {
  return Math.ceil(totalNotes.value / pageSize.value);
});

// 获取用户名首字母作为头像
const getInitials = (name) => {
  if (!name) return '?';
  return name.charAt(0).toUpperCase();
};

// 格式化日期
const formatDate = (dateString, shortFormat = false) => {
  if (!dateString) return '未知时间';
  
  try {
    const date = new Date(dateString);
    
    // 显示相对时间（今天、昨天等）
    const now = new Date();
    const diffDays = Math.floor((now - date) / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) {
      // 今天，显示具体时间
      const hours = date.getHours().toString().padStart(2, '0');
      const minutes = date.getMinutes().toString().padStart(2, '0');
      return `今天 ${hours}:${minutes}`;
    } else if (diffDays === 1) {
      return '昨天';
    } else if (diffDays < 7) {
      return `${diffDays}天前`;
    } else {
      // 一周以上，显示完整日期
      const year = date.getFullYear();
      const month = (date.getMonth() + 1).toString().padStart(2, '0');
      const day = date.getDate().toString().padStart(2, '0');
      
      if (shortFormat) {
        return `${month}-${day}`;
      } else {
        return `${year}-${month}-${day}`;
      }
    }
  } catch (e) {
    console.error('日期格式化错误:', e);
    return '日期无效';
  }
};

// 切换笔记展开/折叠状态
const toggleExpand = (noteId) => {
  expandedNotes[noteId] = !expandedNotes[noteId];
};

// 页面跳转函数
const goToPage = (page) => {
  notesStore.fetchPublicNotes(page, pageSize.value);
  // 重置展开状态
  Object.keys(expandedNotes).forEach(key => {
    expandedNotes[key] = false;
  });
  // 滚动回顶部
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

// 初始化加载
onMounted(() => {
  notesStore.fetchPublicNotes(1, pageSize.value);
});
</script>

<style scoped>
.line-clamp-5 {
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
  overflow: hidden;
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

/* 动画效果 */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.animate-fadeIn {
  animation: fadeIn 0.5s ease-in-out;
}
</style> 