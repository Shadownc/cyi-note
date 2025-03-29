<template>
  <div class="container mx-auto py-8">
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 dark:border-primary-300"></div>
    </div>

    <div v-else>
      <!-- 返回按钮 -->
      <div class="mb-6">
        <router-link 
          :to="isNewNote ? '/notes' : `/notes/${route.params.id}`" 
          class="flex items-center text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          {{ isNewNote ? '返回笔记列表' : '返回笔记详情' }}
        </router-link>
      </div>

      <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-50 mb-6">{{ isNewNote ? '新建笔记' : '编辑笔记' }}</h1>

      <!-- 表单 -->
      <form @submit.prevent="saveNote" class="space-y-6">
        <!-- 标题 -->
        <div>
          <label for="title" class="block text-gray-700 dark:text-gray-200 font-medium mb-2">标题</label>
          <input 
            id="title" 
            v-model="noteForm.title" 
            type="text" 
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-white"
            placeholder="请输入标题" 
            required
          >
        </div>

        <!-- 标签 -->
        <div>
          <label class="block text-gray-700 dark:text-gray-200 font-medium mb-2">标签</label>
          <TagManager 
            :tags="noteForm.tags" 
            @add-tag="addTag" 
            @remove-tag="removeTag" 
          />
        </div>

        <!-- Markdown编辑器 -->
        <div>
          <label class="block text-gray-700 dark:text-gray-200 font-medium mb-2">内容</label>
          <MarkdownEditor 
            :initial-value="noteForm.content" 
            @update:content="noteForm.content = $event" 
          />
        </div>

        <!-- 提交按钮 -->
        <div class="flex justify-end space-x-3">
          <button 
            type="button" 
            @click="$router.go(-1)" 
            class="btn btn-secondary dark:text-gray-200 dark:border-gray-600"
          >
            取消
          </button>
          <button 
            type="submit" 
            class="btn btn-primary dark:bg-primary-600 dark:hover:bg-primary-700"
            :disabled="isSaving"
          >
            {{ isSaving ? '保存中...' : '保存笔记' }}
          </button>
        </div>
      </form>
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