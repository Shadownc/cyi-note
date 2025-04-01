<template>
  <div class="bg-gradient-to-br from-indigo-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900 min-h-screen pb-16">
    <!-- 顶部装饰 -->
    <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-primary-400 via-purple-500 to-pink-500"></div>
    
    <div class="container mx-auto px-4 sm:px-6 pt-6 max-w-5xl relative z-10">
      <!-- 返回与状态栏 -->
      <header class="flex items-center justify-between mb-5">
        <router-link 
          :to="isNewNote ? '/notes' : `/notes/${route.params.id}`" 
          class="group flex items-center space-x-2 text-gray-700 dark:text-gray-200 font-medium"
        >
          <div class="w-8 h-8 flex items-center justify-center rounded-full bg-white dark:bg-gray-800 shadow-sm group-hover:shadow-md transition-all transform group-hover:-translate-x-1">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
          </div>
          <span class="text-sm group-hover:underline">{{ isNewNote ? '返回笔记列表' : '返回笔记详情' }}</span>
        </router-link>
        
        <div v-if="loading" class="flex items-center space-x-2 text-primary-500 dark:text-primary-400">
          <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="text-sm font-medium">加载中...</span>
        </div>
        
        <div v-else-if="isSaving" class="flex items-center space-x-2 text-primary-500 dark:text-primary-400">
          <span class="inline-block w-2 h-2 rounded-full bg-primary-500 animate-ping"></span>
          <span class="text-sm font-medium">保存中...</span>
        </div>
      </header>

      <!-- 主要内容区 -->
      <div v-if="!loading" class="transition-all duration-500 ease-out" :class="{'opacity-40': isSaving}">
        <!-- 标题输入 -->
        <div class="relative mb-4">
          <div class="absolute inset-0 bg-gradient-to-r from-primary-500/5 to-purple-500/5 dark:from-primary-700/10 dark:to-purple-700/10 rounded-xl transform rotate-180 blur-xl"></div>
          <div class="relative bg-white/80 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-md overflow-hidden">
            <input 
              id="title" 
              v-model="noteForm.title" 
              type="text" 
              class="w-full text-xl sm:text-2xl font-bold bg-transparent px-5 py-4 border-0 focus:outline-none focus:ring-0 text-gray-800 dark:text-white placeholder-gray-400/50 dark:placeholder-gray-500/50"
              placeholder="给笔记一个标题..." 
              required
              autocomplete="off"
            >
          </div>
        </div>

        <!-- 主要布局：左侧内容，右侧工具 -->
        <div class="grid grid-cols-1 lg:grid-cols-7 gap-4">
          <!-- 左侧内容区 (占5/7) -->
          <div class="lg:col-span-5">
            <div class="relative h-full">
              <div class="absolute inset-0 bg-gradient-to-r from-primary-500/5 to-purple-500/5 dark:from-primary-700/10 dark:to-purple-700/10 rounded-xl blur-xl"></div>
              <div class="relative bg-white/80 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-md overflow-hidden h-full">
                <div class="flex space-x-1 bg-gray-100/80 dark:bg-gray-700/60 p-2 border-b border-gray-200 dark:border-gray-600/30">
                  <button 
                    @click="activeTab = 'edit'" 
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors"
                    :class="activeTab === 'edit' ? 'bg-white dark:bg-gray-800 text-primary-600 dark:text-primary-400 shadow-sm' : 'hover:bg-white/50 dark:hover:bg-gray-800/50 text-gray-600 dark:text-gray-300'"
                  >
                    编辑
                  </button>
                  <button 
                    @click="activeTab = 'preview'" 
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors"
                    :class="activeTab === 'preview' ? 'bg-white dark:bg-gray-800 text-primary-600 dark:text-primary-400 shadow-sm' : 'hover:bg-white/50 dark:hover:bg-gray-800/50 text-gray-600 dark:text-gray-300'"
                  >
                    预览
                  </button>
                </div>
                
                <div class="h-[calc(100%-44px)]">
                  <!-- 编辑模式 -->
                  <div v-show="activeTab === 'edit'">
                    <textarea 
                      id="content"
                      v-model="noteForm.content" 
                      class="w-full px-5 py-4 text-base bg-transparent border-0 min-h-[350px] sm:min-h-[450px] md:min-h-[500px] h-full resize-none focus:outline-none focus:ring-0 text-gray-700 dark:text-gray-200 overflow-auto"
                      placeholder="开始书写你的想法... 支持Markdown格式"
                    ></textarea>
                  </div>
                  
                  <!-- 预览模式 -->
                  <div v-show="activeTab === 'preview'" class="h-full overflow-auto px-5 py-4">
                    <MarkdownRenderer 
                      :content="noteForm.content" 
                      class="prose dark:prose-invert sm:prose-sm md:prose-base max-w-none"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 右侧工具区 (占2/7) -->
          <div class="lg:col-span-2 space-y-4">
            <!-- 标签区域 -->
            <div class="relative">
              <div class="absolute inset-0 bg-gradient-to-r from-primary-500/5 to-indigo-500/5 dark:from-primary-700/10 dark:to-indigo-700/10 rounded-xl blur-lg"></div>
              <div class="relative bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-md overflow-hidden">
                <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-700/50 flex items-center space-x-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-primary-500 dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                  </svg>
                  <h2 class="text-base font-semibold text-gray-800 dark:text-gray-200">标签</h2>
                </div>
                <div class="p-4">
                  <TagManager 
                    :tags="noteForm.tags" 
                    @add-tag="addTag" 
                    @remove-tag="removeTag"
                    @generate-tags="generateTags"
                  />
                </div>
              </div>
            </div>

            <!-- 附件区域 -->
            <div class="relative">
              <div class="absolute inset-0 bg-gradient-to-r from-indigo-500/5 to-purple-500/5 dark:from-indigo-700/10 dark:to-purple-700/10 rounded-xl blur-lg"></div>
              <div class="relative bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-md overflow-hidden">
                <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-700/50 flex items-center space-x-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-indigo-500 dark:text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                  </svg>
                  <h2 class="text-base font-semibold text-gray-800 dark:text-gray-200">附件</h2>
                </div>
                <div class="p-4" @click.stop.prevent @submit.prevent @mousedown.stop>
                  <AttachmentManager 
                    :noteId="noteId" 
                    :tempMode="isNewNote"
                    @upload-success="handleAttachmentUploadSuccess"
                    @delete-success="handleAttachmentDeleteSuccess"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 平板/手机视图下的工具栏：内容下方展示工具栏 -->
        <div class="lg:hidden grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
          <!-- 标签区域（移动视图） -->
          <div class="relative">
            <div class="absolute inset-0 bg-gradient-to-r from-primary-500/5 to-indigo-500/5 dark:from-primary-700/10 dark:to-indigo-700/10 rounded-xl blur-lg"></div>
            <div class="relative bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-md overflow-hidden">
              <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-700/50 flex items-center space-x-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-primary-500 dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
                <h2 class="text-base font-semibold text-gray-800 dark:text-gray-200">标签</h2>
              </div>
              <div class="p-4">
                <TagManager 
                  :tags="noteForm.tags" 
                  @add-tag="addTag" 
                  @remove-tag="removeTag"
                  @generate-tags="generateTags"
                />
              </div>
            </div>
          </div>

          <!-- 附件区域（移动视图） -->
          <div class="relative">
            <div class="absolute inset-0 bg-gradient-to-r from-indigo-500/5 to-purple-500/5 dark:from-indigo-700/10 dark:to-purple-700/10 rounded-xl blur-lg"></div>
            <div class="relative bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-md overflow-hidden">
              <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-700/50 flex items-center space-x-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-indigo-500 dark:text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                </svg>
                <h2 class="text-base font-semibold text-gray-800 dark:text-gray-200">附件</h2>
              </div>
              <div class="p-4" @click.stop.prevent @submit.prevent @mousedown.stop>
                <AttachmentManager 
                  :noteId="noteId" 
                  :tempMode="isNewNote"
                  @upload-success="handleAttachmentUploadSuccess"
                  @delete-success="handleAttachmentDeleteSuccess"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载中显示 -->
      <div v-if="loading" class="flex flex-col items-center justify-center h-[60vh]">
        <div class="relative w-20 h-20">
          <div class="absolute top-0 left-0 w-full h-full rounded-full border-4 border-primary-200 dark:border-primary-900"></div>
          <div class="absolute top-0 left-0 w-full h-full rounded-full border-t-4 border-primary-500 dark:border-primary-400 animate-spin"></div>
          <div class="absolute top-0 left-0 w-full h-full flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-primary-500 dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
        </div>
        <p class="mt-4 text-sm text-primary-600 dark:text-primary-400 font-medium">正在加载笔记内容...</p>
      </div>
    </div>

    <!-- 底部保存操作栏 -->
    <div v-if="!loading" class="fixed bottom-0 left-0 right-0 bg-white/70 dark:bg-gray-900/80 backdrop-blur-md shadow-lg border-t border-gray-200 dark:border-gray-800 py-3 z-20">
      <div class="container mx-auto px-4 sm:px-6 max-w-5xl flex justify-end space-x-3">
        <button 
          type="button" 
          @click.prevent="$router.go(-1)" 
          class="px-5 py-2 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 rounded-full border border-gray-300 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-750 shadow-sm transition-all duration-200 text-sm"
        >
          取消
        </button>
        <button 
          type="button" 
          @click.prevent="saveNote" 
          class="px-6 py-2 bg-gradient-to-r from-primary-500 to-indigo-600 hover:from-primary-600 hover:to-indigo-700 text-white rounded-full shadow-md hover:shadow-lg transition-all duration-200 min-w-[100px] flex items-center justify-center font-medium text-sm"
          :disabled="isSaving"
        >
          <svg v-if="isSaving" class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isSaving ? '保存中...' : '保存笔记' }}
        </button>
      </div>
    </div>

    <!-- 背景装饰元素 -->
    <div class="fixed top-0 right-0 w-1/3 h-1/3 bg-gradient-to-b from-primary-200/10 via-indigo-200/5 to-transparent dark:from-primary-900/10 dark:via-indigo-900/5 dark:to-transparent blur-3xl -z-10 rounded-full transform translate-x-1/3 -translate-y-1/4"></div>
    <div class="fixed bottom-0 left-0 w-1/2 h-1/2 bg-gradient-to-t from-purple-200/10 via-pink-200/5 to-transparent dark:from-purple-900/10 dark:via-pink-900/5 dark:to-transparent blur-3xl -z-10 rounded-full transform -translate-x-1/4 translate-y-1/4"></div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useNotesStore } from '@/stores/notes';
import { useTagsStore } from '@/stores/tags';
import { useAttachmentsStore } from '@/stores/attachments';
import { storeToRefs } from 'pinia';
import TagManager from '@/components/TagManager.vue';
import AttachmentManager from '@/components/AttachmentManager.vue';
import MarkdownRenderer from '@/components/MarkdownRenderer.vue';

const route = useRoute();
const router = useRouter();
const notesStore = useNotesStore();
const tagsStore = useTagsStore();
const attachmentsStore = useAttachmentsStore();

// 从store获取状态
const { currentNote, isLoading } = storeToRefs(notesStore);
const { tags } = storeToRefs(tagsStore);

// 本地状态
const noteId = route.params.id ? parseInt(route.params.id) : null;
const isNewNote = !noteId;
const isSaving = ref(false);
const loading = ref(false);

// 表单数据
const noteForm = reactive({
  title: '',
  content: '',
  tags: [],
  tempAttachments: []
});

// 编辑器状态
const activeTab = ref('edit');

// 添加标签到表单
const addTag = (tag) => {
  // 检查标签是否已存在
  if (!noteForm.tags.some(t => t.id === tag.id)) {
    noteForm.tags.push(tag);
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
const handleAttachmentUploadSuccess = async (event, fileData) => {
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  if (isNewNote && fileData) {
    // 新建笔记模式下，将附件添加到临时存储
    noteForm.tempAttachments.push(fileData);
  } else if (noteId) {
    // 刷新笔记数据以获取最新的附件列表
    await notesStore.fetchNote(noteId);
  }
};

// 处理附件删除成功
const handleAttachmentDeleteSuccess = async (event, deletedAttachment) => {
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  if (isNewNote && deletedAttachment) {
    // 临时模式下，从表单中移除删除的附件
    noteForm.tempAttachments = noteForm.tempAttachments.filter(
      attachment => attachment.id !== deletedAttachment.id
    );
  } else if (noteId) {
    // 刷新笔记数据以获取最新的附件列表
    await notesStore.fetchNote(noteId);
  }
};

// 保存笔记
const saveNote = async () => {
  if (!noteForm.title.trim()) {
    alert('标题不能为空');
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
    
    if (isNewNote) {
      // 创建新笔记
      const response = await notesStore.createNote(noteData);
      
      // 获取笔记ID
      let newNoteId = null;
      
      if (typeof response === 'object') {
        if (response.id) {
          newNoteId = response.id;
        } else if (response.note && response.note.id) {
          newNoteId = response.note.id;
        }
      }
      
      if (!newNoteId) {
        alert('笔记已保存，但无法跳转到详情页面。请刷新笔记列表查看。');
        router.push({ name: 'notes' });
      } else {
        // 上传所有暂存的附件
        if (noteForm.tempAttachments.length > 0) {
          try {
            for (const fileData of noteForm.tempAttachments) {
              await attachmentsStore.uploadAttachment(fileData.file, newNoteId);
            }
          } catch (attachmentError) {
            alert('笔记保存成功，但部分附件上传失败');
          }
        }
        
        router.push({ name: 'note-detail', params: { id: newNoteId } });
      }
    } else {
      // 更新笔记
      await notesStore.updateNote(noteId, noteData);
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
    await notesStore.fetchNote(noteId);
    
    if (currentNote.value) {
      // 填充表单数据
      noteForm.title = currentNote.value.title;
      noteForm.content = currentNote.value.content;
      noteForm.tags = currentNote.value.tags || [];
    } else {
      // 笔记不存在，重定向到笔记列表
      router.push({ name: 'notes' });
    }
  }
  loading.value = false;
});
</script>

<style scoped>
/* 页面进入动画 */
@keyframes fadeUpIn {
  from { opacity: 0; transform: translateY(15px); }
  to { opacity: 1; transform: translateY(0); }
}

.animate-fadeIn {
  animation: fadeUpIn 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

/* 自定义滚动条 */
textarea::-webkit-scrollbar {
  width: 3px;
}

textarea::-webkit-scrollbar-track {
  background: transparent;
}

textarea::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.3);
  border-radius: 6px;
}

.dark textarea::-webkit-scrollbar-thumb {
  background-color: rgba(75, 85, 99, 0.3);
}

/* 输入框样式 */
textarea {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  line-height: 1.7;
  letter-spacing: 0.2px;
}

/* 美化输入体验 */
textarea, input[type="text"] {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

textarea:focus, input[type="text"]:focus {
  transform: translateY(-1px);
}

/* 移动端优化 */
@media (max-width: 640px) {
  textarea, input {
    font-size: 16px; /* 防止iOS缩放 */
  }
  
  input[type="text"] {
    font-size: 20px;
  }
}

/* 渐变按钮效果 */
button {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

button:active {
  transform: scale(0.97);
}

/* 毛玻璃按钮效果 */
.backdrop-blur-md {
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

/* 加载动画 */
@keyframes ping {
  0% { transform: scale(1); opacity: 1; }
  75%, 100% { transform: scale(2); opacity: 0; }
}

.animate-ping {
  animation: ping 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

/* 复制成功提示 */
.copy-success-tooltip {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  background-color: #10b981;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  z-index: 50;
  animation: fadeInOut 2s forwards;
}

@keyframes fadeInOut {
  0%, 100% { opacity: 0; }
  10%, 90% { opacity: 1; }
}
</style> 