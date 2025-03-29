<template>
  <div class="container mx-auto px-3 sm:px-4 py-4 sm:py-6 md:py-8 max-w-5xl">
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 dark:border-primary-300"></div>
    </div>

    <div v-else class="animate-fadeIn">
      <!-- 返回按钮 -->
      <div class="mb-4 sm:mb-6">
        <router-link 
          :to="isNewNote ? '/notes' : `/notes/${route.params.id}`" 
          class="flex items-center text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 transition-colors text-sm sm:text-base"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 sm:h-5 sm:w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          {{ isNewNote ? '返回笔记列表' : '返回笔记详情' }}
        </router-link>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-4 sm:p-5 md:p-6 lg:p-8">
        <h1 class="text-lg sm:text-xl md:text-2xl font-bold text-gray-800 dark:text-gray-50 mb-4 sm:mb-6 flex flex-wrap items-center">
          <span class="mr-2">{{ isNewNote ? '新建笔记' : '编辑笔记' }}</span>
          <span v-if="isSaving" class="inline-block animate-pulse text-primary-500 text-sm font-normal">保存中...</span>
        </h1>

        <!-- 表单 -->
        <form @submit.prevent="saveNote" class="space-y-4 sm:space-y-6">
          <!-- 标题 -->
          <div class="relative">
            <label for="title" class="block text-gray-700 dark:text-gray-200 font-medium mb-1 sm:mb-2 text-sm sm:text-base">标题</label>
            <input 
              id="title" 
              v-model="noteForm.title" 
              type="text" 
              class="w-full px-3 py-2 text-base border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-white transition-colors"
              placeholder="请输入标题" 
              required
              autocomplete="off"
            >
          </div>

          <!-- 标签 -->
          <div class="relative">
            <label class="block text-gray-700 dark:text-gray-200 font-medium mb-1 sm:mb-2 text-sm sm:text-base">标签</label>
            <div class="p-3 sm:p-4 bg-gray-50 dark:bg-gray-750 rounded-md border border-gray-200 dark:border-gray-600">
              <TagManager 
                :tags="noteForm.tags" 
                @add-tag="addTag" 
                @remove-tag="removeTag" 
              />
            </div>
          </div>

          <!-- Markdown编辑器 -->
          <div class="relative">
            <label class="block text-gray-700 dark:text-gray-200 font-medium mb-1 sm:mb-2 text-sm sm:text-base">内容</label>
            <div class="p-3 sm:p-4 bg-gray-50 dark:bg-gray-750 rounded-md border border-gray-200 dark:border-gray-600">
              <MarkdownEditor 
                :initial-value="noteForm.content" 
                @update:content="noteForm.content = $event" 
              />
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-1 italic">
              支持Markdown格式，可使用上方工具栏添加格式化内容
            </div>
          </div>

          <!-- 提交按钮 -->
          <div class="flex flex-col-reverse sm:flex-row justify-end gap-3 sm:space-x-3 mt-6 sm:mt-8 pt-4 border-t border-gray-200 dark:border-gray-700">
            <button 
              type="button" 
              @click="$router.go(-1)" 
              class="btn btn-secondary dark:text-gray-200 dark:border-gray-600 py-2 sm:py-2 text-sm sm:text-base w-full sm:w-auto"
            >
              取消
            </button>
            <button 
              type="submit" 
              class="btn btn-primary dark:bg-primary-600 dark:hover:bg-primary-700 py-2 sm:py-2 text-sm sm:text-base w-full sm:w-auto flex justify-center items-center"
              :disabled="isSaving"
            >
              <svg v-if="isSaving" class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isSaving ? '保存中...' : '保存笔记' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import TagManager from '@/components/TagManager.vue';

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const isSaving = ref(false);

// 判断是新建还是编辑
const isNewNote = computed(() => route.params.id === 'new');

// 笔记表单数据
const noteForm = ref({
  title: '',
  content: '',
  tags: []
});

// 获取所有标签（模拟数据）
const allTags = [
  { id: 1, name: "前端" },
  { id: 2, name: "Vue" },
  { id: 3, name: "JavaScript" },
  { id: 4, name: "后端" },
  { id: 5, name: "Docker" }
];

// 初始化笔记数据
const fetchNoteData = async () => {
  loading.value = true;
  
  if (isNewNote.value) {
    // 新建笔记
    noteForm.value = {
      title: '',
      content: '',
      tags: []
    };
  } else {
    // 编辑现有笔记，获取数据
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 模拟笔记数据
    const mockNotes = [
      {
        id: 1,
        title: "Vue.js 学习笔记",
        content: "# Vue.js 基础\n\nVue.js 是一套用于构建用户界面的**渐进式框架**。\n\n## 组件化开发\n\nVue 组件包含三个部分：\n\n- template（模板）\n- script（脚本）\n- style（样式）\n\n```js\nconst app = Vue.createApp({\n  data() {\n    return {\n      message: 'Hello Vue!'\n    }\n  }\n})\n```",
        tags: [
          { id: 1, name: "前端" },
          { id: 2, name: "Vue" }
        ]
      },
      {
        id: 2,
        title: "JavaScript 异步编程",
        content: "# JavaScript 异步编程\n\n异步编程是 JavaScript 中的重要概念。\n\n## Promise\n\nPromise 是一种处理异步操作的对象，代表了一个最终可能成功或失败的操作及其结果值。\n\n```js\nconst promise = new Promise((resolve, reject) => {\n  setTimeout(() => {\n    resolve('成功!');\n  }, 1000);\n});\n\npromise.then(value => {\n  console.log(value); // 成功!\n});\n```\n\n## Async/Await\n\n`async/await` 是基于 Promise 的语法糖，使异步代码更容易编写和理解。",
        tags: [
          { id: 1, name: "前端" },
          { id: 3, name: "JavaScript" }
        ]
      },
      {
        id: 3,
        title: "Docker 容器化部署",
        content: "# Docker 容器化部署\n\nDocker 是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中。\n\n## 基本概念\n\n- 镜像（Image）\n- 容器（Container）\n- 仓库（Repository）\n\n## 常用命令\n\n```bash\n# 构建镜像\ndocker build -t myapp .\n\n# 运行容器\ndocker run -p 8080:80 myapp\n```",
        tags: [
          { id: 4, name: "后端" },
          { id: 5, name: "Docker" }
        ]
      }
    ];
    
    const noteId = parseInt(route.params.id);
    const note = mockNotes.find(n => n.id === noteId);
    
    if (note) {
      noteForm.value = {
        title: note.title,
        content: note.content,
        tags: [...note.tags] // 复制标签数组
      };
    } else {
      // 笔记不存在，跳转回列表页
      router.push('/notes');
      return;
    }
  }
  
  loading.value = false;
};

// 添加标签
const addTag = (tag) => {
  // 检查是否已存在该标签名
  const existingTag = allTags.find(t => t.name.toLowerCase() === tag.name.toLowerCase());
  
  if (existingTag) {
    // 使用现有标签
    if (!noteForm.value.tags.some(t => t.id === existingTag.id)) {
      noteForm.value.tags.push(existingTag);
    }
  } else {
    // 创建新标签（在实际应用中，这里会调用API创建新标签）
    const newTag = {
      id: Date.now(), // 模拟生成ID
      name: tag.name
    };
    noteForm.value.tags.push(newTag);
  }
};

// 移除标签
const removeTag = (tagId) => {
  noteForm.value.tags = noteForm.value.tags.filter(tag => tag.id !== tagId);
};

// 保存笔记
const saveNote = async () => {
  if (!noteForm.value.title || !noteForm.value.content) {
    alert('请填写标题和内容');
    return;
  }
  
  isSaving.value = true;
  
  // 模拟API请求
  await new Promise(resolve => setTimeout(resolve, 800));
  
  // 显示成功消息
  alert(isNewNote.value ? '笔记创建成功！' : '笔记更新成功！');
  
  // 导航到笔记列表或详情页
  if (isNewNote.value) {
    router.push('/notes');
  } else {
    router.push(`/notes/${route.params.id}`);
  }
  
  isSaving.value = false;
};

onMounted(() => {
  fetchNoteData();
});
</script>

<style scoped>
.animate-fadeIn {
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 移动端优化 */
@media (max-width: 640px) {
  .container {
    padding-left: 0.75rem;
    padding-right: 0.75rem;
  }
  
  /* 确保表单控件在移动端有足够的点击面积 */
  input, button {
    min-height: 40px;
  }
  
  textarea {
    min-height: 150px;
  }
  
  /* 调整按钮和间距 */
  .btn {
    padding-top: 0.5rem;
    padding-bottom: 0.5rem;
  }
}

/* 确保safari中输入框有正确的样式 */
input {
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
}

/* 移除表单元素聚焦时的默认边框 */
*:focus {
  outline: none;
}
</style> 