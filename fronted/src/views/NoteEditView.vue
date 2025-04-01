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
        <form @submit.prevent="saveNote" class="space-y-4 sm:space-y-6" id="noteForm" @keydown.enter.prevent>
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
              @keydown.enter.prevent
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
                @generate-tags="generateTags"
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

          <!-- 附件管理 -->
          <div class="relative">
            <label class="block text-gray-700 dark:text-gray-200 font-medium mb-1 sm:mb-2 text-sm sm:text-base">附件</label>
            <div 
              class="p-3 sm:p-4 bg-gray-50 dark:bg-gray-750 rounded-md border border-gray-200 dark:border-gray-600" 
              @click.stop.prevent
              @submit.prevent
              @mousedown.stop
            >
              <AttachmentManager 
                v-if="!isNewNote" 
                :noteId="parseInt(route.params.id)" 
                @upload-success="handleAttachmentUploadSuccess"
                @delete-success="handleAttachmentDeleteSuccess"
              />
              <div v-else class="text-center py-4 text-gray-500 dark:text-gray-400">
                <p>请先保存笔记，然后才能添加附件</p>
              </div>
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-1 italic">
              上传图片、文档等文件作为笔记附件
            </div>
          </div>

          <!-- 提交按钮 -->
          <div class="flex flex-col-reverse sm:flex-row justify-end gap-3 sm:space-x-3 mt-6 sm:mt-8 pt-4 border-t border-gray-200 dark:border-gray-700">
            <button 
              type="button" 
              @click.prevent="$router.go(-1)" 
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
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import { useNotesStore } from '@/stores/notes';
import { useTagsStore } from '@/stores/tags';
import { storeToRefs } from 'pinia';
import TagManager from '@/components/TagManager.vue';
import AttachmentManager from '@/components/AttachmentManager.vue';

const route = useRoute();
const router = useRouter();
const notesStore = useNotesStore();
const tagsStore = useTagsStore();

// 从store获取状态
const { currentNote, isLoading } = storeToRefs(notesStore);
const { tags, isLoading: tagsLoading } = storeToRefs(tagsStore);

// 本地状态
const noteId = route.params.id ? parseInt(route.params.id) : null;
const isNewNote = !noteId;
const formTitle = computed(() => isNewNote ? '创建新笔记' : '编辑笔记');
const submitting = ref(false);
const isSaving = ref(false);
const loading = ref(false);
const showPreview = ref(false);

// 表单数据
const noteForm = reactive({
  title: '',
  content: '',
  tags: []
});

// 标签搜索
const tagSearchText = ref('');
const filteredTags = computed(() => {
  if (!tagSearchText.value) return tags.value;
  const search = tagSearchText.value.toLowerCase();
  return tags.value.filter(tag => tag.name.toLowerCase().includes(search));
});

// 添加标签到表单
const addTag = (tag) => {
  // 检查标签是否已存在
  if (!noteForm.tags.some(t => t.id === tag.id)) {
    noteForm.tags.push(tag);
  }
  tagSearchText.value = '';
};

// 创建新标签
const createTag = async () => {
  if (!tagSearchText.value.trim()) return;
  
  try {
    const newTag = await tagsStore.createTag({ name: tagSearchText.value.trim() });
    addTag(newTag);
  } catch (error) {
    console.error('创建标签失败:', error);
  }
};

// 从表单中移除标签
const removeTag = (tagId) => {
  noteForm.tags = noteForm.tags.filter(tag => tag.id !== tagId);
};

// 使用AI生成标签建议
const generateTags = async () => {
  if (!noteForm.content.trim()) return;
  
  try {
    const suggestions = await tagsStore.generateTagSuggestions(noteForm.content);
    // 添加建议的标签
    if (suggestions && suggestions.tags) {
      for (const tagName of suggestions.tags) {
        // 查找是否已存在该标签
        let tag = tags.value.find(t => t.name.toLowerCase() === tagName.toLowerCase());
        
        // 如果不存在，创建新标签
        if (!tag) {
          tag = await tagsStore.createTag({ name: tagName });
        }
        
        // 添加标签到表单
        addTag(tag);
      }
    }
  } catch (error) {
    console.error('生成标签建议失败:', error);
  }
};

// 处理附件上传成功
const handleAttachmentUploadSuccess = async (event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  console.log('附件上传成功，刷新笔记数据');
  // 刷新笔记数据以获取最新的附件列表，但不进行页面跳转
  if (noteId) {
    try {
      await notesStore.fetchNote(noteId);
      // 成功上传提示
      alert('附件上传成功！');
    } catch (error) {
      console.error('刷新笔记数据失败:', error);
    }
  }
}

// 处理附件删除成功
const handleAttachmentDeleteSuccess = async (event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  console.log('附件删除成功，刷新笔记数据');
  // 刷新笔记数据以获取最新的附件列表
  if (noteId) {
    try {
      await notesStore.fetchNote(noteId);
      // 删除成功提示
      alert('附件删除成功！');
    } catch (error) {
      console.error('刷新笔记数据失败:', error);
    }
  }
}

// 切换预览
const togglePreview = () => {
  showPreview.value = !showPreview.value;
};

// 保存笔记
const saveNote = async () => {
  if (!noteForm.title.trim() || !noteForm.content.trim()) {
    alert('标题和内容不能为空');
    return;
  }
  
  isSaving.value = true;
  
  try {
    // 准备提交的数据
    const noteData = {
      title: noteForm.title,
      content: noteForm.content,
      tags: noteForm.tags.map(tag => tag.name)
    };
    
    console.log('准备保存笔记数据:', noteData);
    
    if (isNewNote) {
      // 创建新笔记
      const response = await notesStore.createNote(noteData);
      console.log('创建笔记成功，完整响应数据:', response);
      
      // 获取笔记ID - 适配不同格式
      let newNoteId = null;
      
      if (typeof response === 'object') {
        console.log('分析响应格式...');
        
        // 检查各种可能的响应格式
        if (response.id) {
          // 直接是笔记对象
          console.log('找到ID在响应根级别');
          newNoteId = response.id;
        } else if (response.note && response.note.id) {
          // 包含在note字段中
          console.log('找到ID在response.note中');
          newNoteId = response.note.id;
        }
        
        console.log('提取的笔记ID:', newNoteId);
      }
      
      if (!newNoteId) {
        console.error('无法从响应中获取笔记ID，请检查API格式:', response);
        alert('笔记已保存，但无法跳转到详情页面。请刷新笔记列表查看。');
        router.push({ name: 'notes' });
      } else {
        console.log('跳转到笔记详情页:', newNoteId);
        router.push({ name: 'note-detail', params: { id: newNoteId } });
      }
    } else {
      // 更新笔记
      console.log('更新已有笔记:', noteId);
      const updatedNote = await notesStore.updateNote(noteId, noteData);
      console.log('更新笔记成功:', updatedNote);
      router.push({ name: 'note-detail', params: { id: noteId } });
    }
  } catch (error) {
    console.error('保存笔记失败:', error);
    alert('保存笔记失败: ' + (error.message || '请稍后重试'));
  } finally {
    isSaving.value = false;
  }
};

// 初始化数据
onMounted(async () => {
  loading.value = true;
  // 加载标签列表
  await tagsStore.fetchTags();
  
  // 如果是编辑现有笔记
  if (noteId) {
    console.log('正在加载笔记 ID:', noteId);
    await notesStore.fetchNote(noteId);
    
    if (currentNote.value) {
      console.log('成功加载笔记数据:', currentNote.value);
      // 填充表单数据
      noteForm.title = currentNote.value.title;
      noteForm.content = currentNote.value.content;
      noteForm.tags = currentNote.value.tags || [];
    } else {
      console.error('笔记加载失败或不存在，ID:', noteId);
      // 笔记不存在，重定向到笔记列表
      router.push({ name: 'notes' });
    }
  }
  loading.value = false;
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