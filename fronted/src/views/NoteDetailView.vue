<template>
  <div class="container mx-auto px-4 py-6 sm:py-8 max-w-5xl">
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 dark:border-primary-300"></div>
    </div>

    <div v-else-if="!note" class="text-center my-12 p-6 bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700">
      <p class="text-gray-600 dark:text-gray-300 mb-4">笔记不存在或已被删除。</p>
      <router-link to="/notes" class="btn btn-primary mt-4">返回笔记列表</router-link>
    </div>

    <div v-else class="animate-fadeIn">
      <!-- 返回按钮 -->
      <div class="mb-6">
        <router-link to="/notes" class="flex items-center text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          返回笔记列表
        </router-link>
      </div>

      <!-- 笔记标题 -->
      <div class="mb-6">
        <div class="flex items-start justify-between">
          <h1 class="text-2xl sm:text-3xl font-bold text-gray-800 dark:text-white break-words">
            {{ note.title }}
            <span v-if="note.is_public" class="inline-flex items-center text-green-600 dark:text-green-400 ml-2 text-sm">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              公开笔记
            </span>
          </h1>
          <div class="flex space-x-2">
            <!-- 右侧按钮区域 -->
          </div>
        </div>
      </div>

      <!-- 笔记详情 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-5 sm:p-6 md:p-8">
        <NoteDetail 
          :note="note" 
          @edit="navigateToEdit" 
          @delete="showDeleteConfirm = true" 
        />
      </div>

      <!-- 删除确认模态框 -->
      <div v-if="showDeleteConfirm" class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50 p-4">
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl p-5 sm:p-6 max-w-md w-full border dark:border-gray-700">
          <h3 class="text-xl font-bold mb-4 text-gray-800 dark:text-gray-100">确认删除</h3>
          <p class="mb-6 text-gray-600 dark:text-gray-300">您确定要删除笔记《{{ note.title }}》吗？此操作无法撤销。</p>
          <div class="flex justify-end space-x-3">
            <button @click="showDeleteConfirm = false" class="btn btn-secondary dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600">取消</button>
            <button @click="deleteNote" class="btn bg-red-600 dark:bg-red-700 text-white hover:bg-red-700 dark:hover:bg-red-800 transition-colors">确认删除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import NoteDetail from '@/components/NoteDetail.vue';
import { useNotesStore } from '@/stores/notes';
import { useTagsStore } from '@/stores/tags';
import { storeToRefs } from 'pinia';

const route = useRoute();
const router = useRouter();
const notesStore = useNotesStore();
const tagsStore = useTagsStore();

// 从store获取状态
const { currentNote, isLoading } = storeToRefs(notesStore);
const { tags } = storeToRefs(tagsStore);

// 将 currentNote 映射为 note（供模板使用）
const note = computed(() => currentNote.value);
const loading = computed(() => isLoading.value);
const showDeleteConfirm = ref(false);

// 本地状态
const noteId = parseInt(route.params.id);

// 格式化日期的函数
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 计算属性：格式化后的创建时间
const formattedCreatedAt = computed(() => {
  return currentNote.value ? formatDate(currentNote.value.created_at) : '';
});

// 计算属性：格式化后的更新时间
const formattedUpdatedAt = computed(() => {
  return currentNote.value ? formatDate(currentNote.value.updated_at) : '';
});

// 确认删除
const confirmDelete = () => {
  showDeleteConfirm.value = true;
};

// 取消删除
const cancelDelete = () => {
  showDeleteConfirm.value = false;
};

// 执行删除
const deleteNote = async () => {
  if (!noteId) return;
  
  try {
    await notesStore.deleteNote(noteId);
    router.push({ name: 'notes' });
  } catch (error) {
    console.error('删除笔记失败:', error);
    // 可以添加错误提示
  }
};

// 编辑笔记
const navigateToEdit = () => {
  router.push({ name: 'note-edit', params: { id: noteId } });
};

// 返回笔记列表
const goBack = () => {
  router.push({ name: 'notes' });
};

// 加载笔记数据
onMounted(async () => {
  if (noteId) {
    console.log('正在加载笔记详情页的笔记数据，ID:', noteId);
    await notesStore.fetchNote(noteId);
    
    console.log('笔记详情页加载结果:', currentNote.value);
    
    // 笔记不存在，重定向到笔记列表
    if (!currentNote.value) {
      console.error('笔记详情页找不到笔记，ID:', noteId);
      router.push({ name: 'notes' });
    }
  }
  
  // 加载标签列表
  await tagsStore.fetchTags();
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
}
</style> 