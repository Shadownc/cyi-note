<template>
  <div class="container mx-auto py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-50">我的笔记</h1>
      <router-link to="/notes/new" class="btn btn-primary flex items-center shadow-sm hover:shadow dark:bg-primary-600 dark:hover:bg-primary-700">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新建笔记
      </router-link>
    </div>

    <!-- 筛选栏 -->
    <div class="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 mb-6 shadow-sm dark:shadow-gray-900/30">
      <div class="flex flex-wrap gap-4">
        <div class="relative flex-grow">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="搜索笔记..." 
            class="w-full px-4 py-2 rounded border border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-white"
          >
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            class="h-5 w-5 absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 dark:text-gray-300" 
            fill="none" 
            viewBox="0 0 24 24" 
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <select 
          v-model="selectedTag" 
          class="px-4 py-2 rounded border border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-white"
        >
          <option value="">所有标签</option>
          <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
        </select>
        <select 
          v-model="sortBy" 
          class="px-4 py-2 rounded border border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-white"
        >
          <option value="updated_desc">最近更新</option>
          <option value="updated_asc">最早更新</option>
          <option value="created_desc">最近创建</option>
          <option value="created_asc">最早创建</option>
        </select>
      </div>
    </div>
    
    <!-- 视图切换和分页控制 -->
    <div class="flex flex-wrap items-center justify-between mb-4">
      <div class="flex items-center space-x-2 mb-2 sm:mb-0">
        <span class="text-sm text-gray-600 dark:text-gray-300">视图:</span>
        <button 
          @click="viewMode = 'card'" 
          class="p-2 rounded-md transition-colors" 
          :class="viewMode === 'card' ? 'bg-primary-100 text-primary-700 dark:bg-primary-600 dark:text-primary-50' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'"
          title="卡片视图"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </button>
        <button 
          @click="viewMode = 'list'" 
          class="p-2 rounded-md transition-colors" 
          :class="viewMode === 'list' ? 'bg-primary-100 text-primary-700 dark:bg-primary-600 dark:text-primary-50' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'"
          title="列表视图"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
        </button>
      </div>
      
      <div class="text-sm text-gray-600 dark:text-gray-300">
        显示 {{ pagedNotes.length > 0 ? (currentPage - 1) * pageSize + 1 : 0 }} - {{ Math.min(currentPage * pageSize, filteredNotes.length) }} 项，共 {{ filteredNotes.length }} 项
      </div>
    </div>

    <!-- 笔记列表 -->
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 dark:border-primary-300"></div>
    </div>

    <div v-else-if="filteredNotes.length === 0" class="text-center my-12">
      <div class="text-gray-500 dark:text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto mb-4 text-gray-400 dark:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-lg mb-4 dark:text-gray-100">暂无笔记</p>
        <router-link to="/notes/new" class="btn btn-primary dark:bg-primary-600 dark:hover:bg-primary-700 shadow hover:shadow-md">创建第一篇笔记</router-link>
      </div>
    </div>

    <!-- 卡片视图 -->
    <div v-else-if="viewMode === 'card'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="note in pagedNotes" :key="note.id" class="mb-4">
        <NoteCard :note="note" />
      </div>
    </div>

    <!-- 列表视图 -->
    <div v-else class="flex flex-col space-y-2">
      <NoteListItem v-for="note in pagedNotes" :key="note.id" :note="note" />
    </div>
    
    <!-- 分页控制 -->
    <div v-if="filteredNotes.length > pageSize" class="flex justify-center mt-8">
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
import { ref, computed, onMounted, watch } from 'vue';
import NoteCard from '@/components/NoteCard.vue';
import NoteListItem from '@/components/NoteListItem.vue';
import { useNotesStore } from '@/stores/notes';
import { useTagsStore } from '@/stores/tags';
import { storeToRefs } from 'pinia';

// 使用Pinia store
const notesStore = useNotesStore();
const tagsStore = useTagsStore();

// 从store中获取状态
const { notes, totalNotes, isLoading: loading, currentPage, pageSize, activeTag, searchKeyword } = storeToRefs(notesStore);
const { tags } = storeToRefs(tagsStore);

// 本地状态
const sortBy = ref('updated_desc');
const viewMode = ref('card'); // 'card' 或 'list'
const selectedTag = ref('');
const searchQuery = ref('');

// 根据筛选条件过滤笔记
const filteredNotes = computed(() => {
  return notes.value;
});

// 总页数
const totalPages = computed(() => {
  return Math.ceil(totalNotes.value / pageSize.value);
});

// 当前页的笔记
const pagedNotes = computed(() => {
  return notes.value;
});

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

// 页面跳转函数
const goToPage = (page) => {
  notesStore.fetchNotes(page, pageSize.value, selectedTag.value ? parseInt(selectedTag.value) : null);
  // 滚动回顶部
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

// 处理搜索
const handleSearch = () => {
  if (searchQuery.value.trim()) {
    notesStore.searchNotes(searchQuery.value);
  } else {
    notesStore.clearSearch();
  }
};

// 处理标签选择
const handleTagSelect = () => {
  if (selectedTag.value) {
    notesStore.filterByTag(parseInt(selectedTag.value));
  } else {
    notesStore.clearSearch();
  }
};

// 处理排序
const handleSortChange = () => {
  // 更新排序条件参数并重新加载笔记
  const params = {
    page: currentPage.value,
    pageSize: pageSize.value,
    sortBy: sortBy.value
  };
  
  if (selectedTag.value) {
    params.tagId = parseInt(selectedTag.value);
  }
  
  notesStore.fetchNotes(params.page, params.pageSize, params.tagId, params.sortBy);
};

// 监听筛选条件变化
watch([sortBy], () => {
  handleSortChange();
});

// 监听搜索输入变化
watch([searchQuery], () => {
  // 使用setTimeout实现简单的防抖
  clearTimeout(searchDebounce.value);
  searchDebounce.value = setTimeout(() => {
    handleSearch();
  }, 500); // 500ms防抖
});

// 监听标签选择变化
watch([selectedTag], () => {
  handleTagSelect();
});

// 搜索防抖定时器
const searchDebounce = ref(null);

// 初始化数据
onMounted(async () => {
  // 获取笔记列表
  await notesStore.fetchNotes();
  
  // 获取标签列表
  await tagsStore.fetchTags();
  
  // 同步本地状态和store状态
  if (activeTag.value) {
    selectedTag.value = activeTag.value.toString();
  }
  
  if (searchKeyword.value) {
    searchQuery.value = searchKeyword.value;
  }
});
</script>

<style scoped>
/* 禁用状态样式 */
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 悬浮效果 */
.hover-scale {
  transition: transform 0.2s;
}
.hover-scale:hover {
  transform: scale(1.02);
}
</style> 