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

// 模拟更多数据 - 实际应用中会从API获取
const notes = ref([
  {
    id: 1,
    title: "Vue.js 学习笔记",
    content: "# Vue.js 基础\n\nVue.js 是一套用于构建用户界面的**渐进式框架**。\n\n## 组件化开发\n\nVue 组件包含三个部分：\n\n- template（模板）\n- script（脚本）\n- style（样式）\n\n```js\nconst app = Vue.createApp({\n  data() {\n    return {\n      message: 'Hello Vue!'\n    }\n  }\n})\n```",
    createdAt: "2023-01-15T10:30:00Z",
    updatedAt: "2023-01-16T14:20:00Z",
    tags: [
      { id: 1, name: "前端" },
      { id: 2, name: "Vue" }
    ]
  },
  {
    id: 2,
    title: "JavaScript 异步编程",
    content: "# JavaScript 异步编程\n\n异步编程是 JavaScript 中的重要概念。\n\n## Promise\n\nPromise 是一种处理异步操作的对象，代表了一个最终可能成功或失败的操作及其结果值。\n\n```js\nconst promise = new Promise((resolve, reject) => {\n  setTimeout(() => {\n    resolve('成功!');\n  }, 1000);\n});\n\npromise.then(value => {\n  console.log(value); // 成功!\n});\n```\n\n## Async/Await\n\n`async/await` 是基于 Promise 的语法糖，使异步代码更容易编写和理解。",
    createdAt: "2023-02-10T08:15:00Z",
    updatedAt: "2023-02-12T11:45:00Z",
    tags: [
      { id: 1, name: "前端" },
      { id: 3, name: "JavaScript" }
    ]
  },
  {
    id: 3,
    title: "Docker 容器化部署",
    content: "# Docker 容器化部署\n\nDocker 是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中。\n\n## 基本概念\n\n- 镜像（Image）\n- 容器（Container）\n- 仓库（Repository）\n\n## 常用命令\n\n```bash\n# 构建镜像\ndocker build -t myapp .\n\n# 运行容器\ndocker run -p 8080:80 myapp\n```",
    createdAt: "2023-03-05T15:20:00Z",
    updatedAt: "2023-03-07T09:30:00Z",
    tags: [
      { id: 4, name: "后端" },
      { id: 5, name: "Docker" }
    ]
  },
  {
    id: 4,
    title: "React Hooks 详解",
    content: "# React Hooks 详解\n\nReact Hooks 是 React 16.8 引入的新特性，它可以让你在不编写 class 的情况下使用 state 和其他 React 特性。\n\n## useState\n\n```jsx\nimport React, { useState } from 'react';\n\nfunction Example() {\n  const [count, setCount] = useState(0);\n\n  return (\n    <div>\n      <p>You clicked {count} times</p>\n      <button onClick={() => setCount(count + 1)}>\n        Click me\n      </button>\n    </div>\n  );\n}\n```",
    createdAt: "2023-04-05T10:20:00Z",
    updatedAt: "2023-04-07T14:30:00Z",
    tags: [
      { id: 1, name: "前端" },
      { id: 6, name: "React" }
    ]
  },
  {
    id: 5,
    title: "Git 工作流指南",
    content: "# Git 工作流指南\n\nGit 是一个分布式版本控制系统，用于跟踪文件的更改并协调多人的工作。\n\n## 常用命令\n\n```bash\n# 创建新分支\ngit checkout -b feature/new-feature\n\n# 添加修改\ngit add .\n\n# 提交修改\ngit commit -m \"Add new feature\"\n\n# 推送到远程\ngit push origin feature/new-feature\n```",
    createdAt: "2023-05-15T09:30:00Z",
    updatedAt: "2023-05-16T11:20:00Z",
    tags: [
      { id: 7, name: "工具" },
      { id: 8, name: "Git" }
    ]
  },
  {
    id: 6,
    title: "Python 数据分析入门",
    content: "# Python 数据分析入门\n\nPython 是数据分析领域的主流语言，与 Pandas, NumPy, Matplotlib 等库一起使用非常强大。\n\n## Pandas 基础\n\n```python\nimport pandas as pd\n\n# 创建 DataFrame\ndf = pd.DataFrame({\n    'A': [1, 2, 3],\n    'B': [4, 5, 6],\n    'C': [7, 8, 9]\n})\n\n# 数据操作\nprint(df.describe())\n```",
    createdAt: "2023-06-20T14:15:00Z",
    updatedAt: "2023-06-22T10:30:00Z",
    tags: [
      { id: 4, name: "后端" },
      { id: 9, name: "Python" },
      { id: 10, name: "数据分析" }
    ]
  }
]);

const tags = ref([
  { id: 1, name: "前端" },
  { id: 2, name: "Vue" },
  { id: 3, name: "JavaScript" },
  { id: 4, name: "后端" },
  { id: 5, name: "Docker" },
  { id: 6, name: "React" },
  { id: 7, name: "工具" },
  { id: 8, name: "Git" },
  { id: 9, name: "Python" },
  { id: 10, name: "数据分析" }
]);

const loading = ref(false);
const searchQuery = ref('');
const selectedTag = ref('');
const sortBy = ref('updated_desc');
const viewMode = ref('card'); // 'card' 或 'list'
const currentPage = ref(1);
const pageSize = ref(3); // 每页显示的笔记数量，从6改为3

// 根据筛选条件过滤笔记
const filteredNotes = computed(() => {
  let result = notes.value;
  
  // 根据搜索关键词过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(note => 
      note.title.toLowerCase().includes(query) || 
      note.content.toLowerCase().includes(query)
    );
  }
  
  // 根据标签过滤
  if (selectedTag.value) {
    result = result.filter(note => 
      note.tags.some(tag => tag.id === Number(selectedTag.value))
    );
  }
  
  // 排序
  result = [...result].sort((a, b) => {
    if (sortBy.value === 'updated_desc') {
      return new Date(b.updatedAt) - new Date(a.updatedAt);
    } else if (sortBy.value === 'updated_asc') {
      return new Date(a.updatedAt) - new Date(b.updatedAt);
    } else if (sortBy.value === 'created_desc') {
      return new Date(b.createdAt) - new Date(a.createdAt);
    } else if (sortBy.value === 'created_asc') {
      return new Date(a.createdAt) - new Date(b.createdAt);
    }
    return 0;
  });
  
  return result;
});

// 总页数
const totalPages = computed(() => {
  return Math.ceil(filteredNotes.value.length / pageSize.value);
});

// 当前页的笔记
const pagedNotes = computed(() => {
  const startIndex = (currentPage.value - 1) * pageSize.value;
  const endIndex = startIndex + pageSize.value;
  return filteredNotes.value.slice(startIndex, endIndex);
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
  currentPage.value = page;
  // 滚动回顶部
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

// 监听筛选条件变化，重置页码
watch([searchQuery, selectedTag, sortBy], () => {
  currentPage.value = 1;
});

// 模拟加载数据
onMounted(() => {
  loading.value = true;
  // 模拟API请求延迟
  setTimeout(() => {
    loading.value = false;
  }, 500);
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