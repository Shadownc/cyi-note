<template>
  <div class="container mx-auto py-8">
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
    </div>

    <div v-else-if="!note" class="text-center my-12 text-gray-600">
      <p>笔记不存在或已被删除。</p>
      <router-link to="/notes" class="btn btn-primary mt-4">返回笔记列表</router-link>
    </div>

    <div v-else>
      <!-- 返回按钮 -->
      <div class="mb-6">
        <router-link to="/notes" class="flex items-center text-blue-600 hover:text-blue-800">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          返回笔记列表
        </router-link>
      </div>

      <!-- 笔记详情 -->
      <NoteDetail 
        :note="note" 
        @edit="navigateToEdit" 
        @delete="showDeleteConfirm = true" 
      />

      <!-- 删除确认模态框 -->
      <div v-if="showDeleteConfirm" class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50">
        <div class="bg-white rounded-lg shadow-xl p-6 max-w-md w-full">
          <h3 class="text-xl font-bold mb-4">确认删除</h3>
          <p class="mb-6">您确定要删除笔记《{{ note.title }}》吗？此操作无法撤销。</p>
          <div class="flex justify-end space-x-3">
            <button @click="showDeleteConfirm = false" class="btn btn-secondary">取消</button>
            <button @click="deleteNote" class="btn bg-red-600 text-white hover:bg-red-700">确认删除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import NoteDetail from '@/components/NoteDetail.vue';

const route = useRoute();
const router = useRouter();
const loading = ref(true);
const note = ref(null);
const showDeleteConfirm = ref(false);

// 模拟API调用，获取笔记详情
const fetchNote = async () => {
  loading.value = true;
  // 模拟API延迟
  await new Promise(resolve => setTimeout(resolve, 500));
  
  // 从本地模拟数据中查找笔记
  // 在真实应用中，这里会是一个API调用
  const noteId = parseInt(route.params.id);
  
  // 模拟数据，真实应用中会从API获取
  const mockNotes = [
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
    }
  ];
  
  note.value = mockNotes.find(n => n.id === noteId) || null;
  loading.value = false;
};

// 导航到编辑页面
const navigateToEdit = () => {
  router.push(`/notes/${route.params.id}/edit`);
};

// 删除笔记
const deleteNote = async () => {
  // 模拟删除API调用
  await new Promise(resolve => setTimeout(resolve, 300));
  
  // 显示成功消息 (在实际应用中可能使用toast组件)
  alert('笔记已成功删除');
  
  // 返回笔记列表
  router.push('/notes');
};

onMounted(() => {
  fetchNote();
});
</script> 