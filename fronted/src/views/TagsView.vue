<template>
  <div class="container mx-auto py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-50">标签管理</h1>
    </div>

    <!-- 标签管理界面 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md dark:shadow-gray-900/30 p-6">
      <div class="mb-6">
        <label class="block text-gray-700 dark:text-gray-200 font-medium mb-2">新增标签</label>
        <div class="flex">
          <input 
            v-model="newTagName" 
            type="text" 
            placeholder="输入标签名..." 
            class="flex-grow px-4 py-2 rounded-l border border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-white"
            @keyup.enter="addTag"
          >
          <button 
            @click="addTag" 
            class="btn btn-primary rounded-l-none px-4 dark:bg-primary-600 dark:hover:bg-primary-700"
            :disabled="!newTagName.trim()"
          >
            添加
          </button>
        </div>
      </div>

      <div class="mb-4">
        <label class="block text-gray-700 dark:text-gray-200 font-medium mb-2">搜索标签</label>
        <div class="relative">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="输入关键词搜索..." 
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
      </div>

      <div v-if="loading" class="flex justify-center my-8">
        <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-primary-500 dark:border-primary-300"></div>
      </div>

      <div v-else-if="filteredTags.length === 0" class="text-center my-8 text-gray-500 dark:text-gray-400">
        <p v-if="searchQuery" class="dark:text-gray-300">没有找到匹配的标签</p>
        <p v-else class="dark:text-gray-300">暂无标签，请添加新标签</p>
      </div>

      <div v-else>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <div 
            v-for="tag in filteredTags" 
            :key="tag.id"
            class="border border-gray-200 dark:border-gray-600 rounded-lg p-3 flex justify-between items-center group hover:border-primary-400 dark:hover:border-primary-500 transition-colors bg-white dark:bg-gray-700"
          >
            <div class="flex items-center">
              <span 
                class="w-3 h-3 rounded-full mr-2"
                :style="{ backgroundColor: getTagColor(tag.name) }"
              ></span>
              <span class="text-gray-800 dark:text-gray-100">{{ tag.name }}</span>
            </div>
            <div class="flex items-center">
              <span class="text-xs text-gray-500 dark:text-gray-300 mr-2">{{ tag.noteCount || 0 }}篇笔记</span>
              <button 
                @click="deleteTag(tag.id)" 
                class="text-gray-400 hover:text-red-500 dark:text-gray-300 dark:hover:text-red-400 opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

const loading = ref(true);
const tags = ref([]);
const newTagName = ref('');
const searchQuery = ref('');
const showDeleteConfirm = ref(false);
const tagToDelete = ref(null);

// 根据标签名生成一致的颜色
const getTagColor = (name) => {
  const colors = [
    '#3B82F6', // blue-500
    '#10B981', // emerald-500
    '#F59E0B', // amber-500
    '#EF4444', // red-500
    '#8B5CF6', // violet-500
    '#EC4899', // pink-500
    '#06B6D4', // cyan-500
    '#F97316', // orange-500
  ];
  
  // 使用标签名的字符编码总和作为索引
  const sum = name.split('')
    .map(char => char.charCodeAt(0))
    .reduce((acc, val) => acc + val, 0);
  
  return colors[sum % colors.length];
};

// 过滤标签
const filteredTags = computed(() => {
  if (!searchQuery.value) return tags.value;
  
  const query = searchQuery.value.toLowerCase();
  return tags.value.filter(tag => 
    tag.name.toLowerCase().includes(query)
  );
});

// 获取标签数据
const fetchTags = async () => {
  loading.value = true;
  
  // 模拟API调用
  await new Promise(resolve => setTimeout(resolve, 500));
  
  // 模拟数据
  tags.value = [
    { id: 1, name: "前端", noteCount: 2 },
    { id: 2, name: "Vue", noteCount: 1 },
    { id: 3, name: "JavaScript", noteCount: 1 },
    { id: 4, name: "后端", noteCount: 1 },
    { id: 5, name: "Docker", noteCount: 1 },
    { id: 6, name: "数据库", noteCount: 0 },
    { id: 7, name: "Python", noteCount: 0 },
    { id: 8, name: "算法", noteCount: 0 }
  ];
  
  loading.value = false;
};

// 添加标签
const addTag = async () => {
  if (!newTagName.value.trim()) return;
  
  const tagName = newTagName.value.trim();
  
  // 检查标签是否已存在
  if (tags.value.some(tag => tag.name.toLowerCase() === tagName.toLowerCase())) {
    alert('标签已存在');
    return;
  }
  
  // 模拟API调用
  await new Promise(resolve => setTimeout(resolve, 300));
  
  // 添加新标签
  const newTag = {
    id: Date.now(), // 模拟生成ID
    name: tagName,
    noteCount: 0
  };
  
  tags.value.push(newTag);
  newTagName.value = '';
};

// 删除标签
const deleteTag = async (id) => {
  // 检查该标签是否有笔记使用
  const tag = tags.value.find(t => t.id === id);
  
  if (tag && tag.noteCount > 0) {
    if (!confirm(`标签"${tag.name}"正在被${tag.noteCount}篇笔记使用，确定要删除吗？`)) {
      return;
    }
  }
  
  // 模拟API调用
  await new Promise(resolve => setTimeout(resolve, 300));
  
  // 从列表中移除
  tags.value = tags.value.filter(tag => tag.id !== id);
};

onMounted(() => {
  fetchTags();
});
</script> 