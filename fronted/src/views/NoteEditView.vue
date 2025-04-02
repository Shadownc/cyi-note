<template>
  <div class="note-edit-view bg-gray-50 dark:bg-gray-900 min-h-screen pb-10">
    <div class="container mx-auto px-4 py-6 max-w-7xl">
      <!-- 顶部操作栏 -->
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-2xl md:text-3xl font-bold text-gray-800 dark:text-white flex items-center gap-2">
          <router-link :to="isNewNote ? '/notes' : `/notes/${noteId}`" class="text-primary-500 hover:text-primary-600 dark:text-primary-400 dark:hover:text-primary-300 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
          </router-link>
          {{ isNewNote ? '创建新笔记' : '编辑笔记' }}
        </h1>
        <div class="flex space-x-3">
          <button 
            @click="cancelEdit" 
            class="px-5 py-2.5 rounded-lg bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors shadow-sm font-medium text-sm"
          >
            取消
          </button>
          <button 
            @click="saveNote" 
            class="px-6 py-2.5 rounded-lg bg-primary-600 hover:bg-primary-700 text-white shadow-md hover:shadow-lg transition-all flex items-center justify-center font-medium text-sm min-w-[100px]" 
            :disabled="isSaving"
          >
            <svg v-if="isSaving" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ isSaving ? '保存中...' : '保存笔记' }}
          </button>
        </div>
      </div>
      
      <!-- 笔记标题 -->
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-5 mb-6">
        <label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">笔记标题</label>
        <input 
          type="text" 
          id="title" 
          v-model="noteForm.title" 
          class="w-full px-4 py-3 text-lg font-medium border border-gray-200 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white transition-shadow"
          placeholder="输入一个引人注目的标题..."
        />
      </div>
      
      <!-- 编辑器和预览区域 -->
      <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
        <!-- 左侧：Markdown编辑器 (占3列) -->
        <div class="lg:col-span-3 space-y-6">
          <!-- Markdown编辑器 -->
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm overflow-hidden">
            <!-- 编辑器工具栏 -->
            <div class="bg-gray-50 dark:bg-gray-750 border-b border-gray-200 dark:border-gray-700 px-4 py-2 flex flex-wrap gap-1.5 sticky top-0 z-10">
              <button @click="insertMarkdown('# ', '', '标题')" title="标题" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                </svg>
              </button>
              <button @click="insertMarkdown('**', '**', '粗体')" title="粗体" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M6 4h8a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"></path>
                  <path d="M6 12h9a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"></path>
                </svg>
              </button>
              <button @click="insertMarkdown('*', '*', '斜体')" title="斜体" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="19" y1="4" x2="10" y2="4"></line>
                  <line x1="14" y1="20" x2="5" y2="20"></line>
                  <line x1="15" y1="4" x2="9" y2="20"></line>
                </svg>
              </button>
              <div class="h-5 border-r border-gray-300 dark:border-gray-600 mx-1"></div>
              <button @click="insertMarkdown('- ', '', '列表项')" title="无序列表" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="8" y1="6" x2="21" y2="6"></line>
                  <line x1="8" y1="12" x2="21" y2="12"></line>
                  <line x1="8" y1="18" x2="21" y2="18"></line>
                  <line x1="3" y1="6" x2="3.01" y2="6"></line>
                  <line x1="3" y1="12" x2="3.01" y2="12"></line>
                  <line x1="3" y1="18" x2="3.01" y2="18"></line>
                </svg>
              </button>
              <button @click="insertMarkdown('1. ', '', '列表项')" title="有序列表" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="10" y1="6" x2="21" y2="6"></line>
                  <line x1="10" y1="12" x2="21" y2="12"></line>
                  <line x1="10" y1="18" x2="21" y2="18"></line>
                  <path d="M4 6h1v4"></path>
                  <path d="M4 10h2"></path>
                  <path d="M6 18H4c0-1 2-2 2-3s-1-1.5-2-1"></path>
                </svg>
              </button>
              <div class="h-5 border-r border-gray-300 dark:border-gray-600 mx-1"></div>
              <button @click="insertMarkdown('[', '](链接URL)', '链接文本')" title="链接" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
                </svg>
              </button>
              <button @click="openImageUploader" title="添加图片/附件" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                  <circle cx="8.5" cy="8.5" r="1.5"></circle>
                  <polyline points="21 15 16 10 5 21"></polyline>
                </svg>
              </button>
              <!-- 隐藏的文件上传输入 -->
              <input 
                type="file" 
                ref="imageInput" 
                @change="handleImageUpload" 
                accept="image/*" 
                class="hidden" 
              />
              <button @click="insertMarkdown('```\n', '\n```', '代码块')" title="代码块" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="16 18 22 12 16 6"></polyline>
                  <polyline points="8 6 2 12 8 18"></polyline>
                </svg>
              </button>
              <button @click="insertMarkdown('> ', '', '引用内容')" title="引用" class="editor-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                </svg>
              </button>
            </div>
            
            <!-- Markdown编辑区域 -->
            <textarea 
              ref="contentTextarea"
              v-model="noteForm.content" 
              class="w-full h-[550px] px-5 py-4 resize-none focus:outline-none dark:bg-gray-800 dark:text-gray-100 bg-white text-gray-800 font-mono leading-relaxed"
              placeholder="开始撰写笔记内容，支持Markdown格式..."
            ></textarea>
          </div>
        </div>
        
        <!-- 右侧：预览和附件管理 (占2列) -->
        <div class="lg:col-span-2 space-y-6">
          <!-- 笔记设置卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm overflow-hidden">
            <div class="bg-gray-50 dark:bg-gray-750 border-b border-gray-200 dark:border-gray-700 px-4 py-3">
              <h3 class="font-medium text-gray-800 dark:text-gray-200 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                笔记设置
              </h3>
            </div>
            <div class="p-4 space-y-5">
              <!-- 笔记可见性 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">笔记可见性</label>
                <div class="flex items-center space-x-4">
                  <label class="inline-flex items-center cursor-pointer">
                    <input type="radio" v-model="noteForm.is_public" :value="false" class="form-radio text-primary-600 focus:ring-primary-500 h-4 w-4">
                    <span class="ml-2 text-gray-700 dark:text-gray-300">私密</span>
                  </label>
                  <label class="inline-flex items-center cursor-pointer">
                    <input type="radio" v-model="noteForm.is_public" :value="true" class="form-radio text-primary-600 focus:ring-primary-500 h-4 w-4">
                    <span class="ml-2 text-gray-700 dark:text-gray-300">公开</span>
                  </label>
                </div>
              </div>
              
              <!-- 笔记标签 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">标签</label>
                <div class="flex flex-wrap items-center gap-2 mb-2">
                  <!-- 已选标签 -->
                  <div 
                    v-for="tag in noteForm.tags" 
                    :key="tag" 
                    class="bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300 px-2.5 py-1 rounded-full text-sm flex items-center"
                  >
                    <span>{{ tag }}</span>
                    <button 
                      @click="removeTag(tag)" 
                      class="ml-1.5 text-primary-500 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-200"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                <!-- 添加标签输入框 -->
                <div class="flex">
                  <input 
                    type="text" 
                    v-model="newTag" 
                    @keydown.enter="addTag"
                    placeholder="添加标签..." 
                    class="flex-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-l-lg focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200"
                  />
                  <button 
                    @click="addTag" 
                    class="px-3 py-2 bg-primary-600 text-white text-sm rounded-r-lg hover:bg-primary-700 transition-colors"
                  >
                    添加
                  </button>
                </div>
                
                <!-- 生成标签按钮 -->
                <button 
                  @click="generateTags" 
                  class="mt-2 text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 flex items-center"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  自动生成标签
                </button>
                
                <!-- 标签建议 -->
                <div v-if="suggestedTags.length > 0" class="mt-3">
                  <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">建议标签:</p>
                  <div class="flex flex-wrap gap-1">
                    <button 
                      v-for="tag in suggestedTags" 
                      :key="tag"
                      @click="addSuggestedTag(tag)"
                      class="px-2 py-0.5 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                    >
                      + {{ tag }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 预览标签页 -->
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm overflow-hidden">
            <div class="bg-gray-50 dark:bg-gray-750 border-b border-gray-200 dark:border-gray-700 px-4 py-3">
              <h3 class="font-medium text-gray-800 dark:text-gray-200 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                预览
              </h3>
            </div>
            <div class="p-5 max-h-[300px] overflow-y-auto prose dark:prose-invert prose-sm max-w-none">
              <MarkdownRenderer :content="noteForm.content" />
            </div>
          </div>
          
          <!-- 附件管理 -->
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm overflow-hidden">
            <div class="bg-gray-50 dark:bg-gray-750 border-b border-gray-200 dark:border-gray-700 px-4 py-3">
              <h3 class="font-medium text-gray-800 dark:text-gray-200 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                </svg>
                附件管理
              </h3>
            </div>
            <div class="p-4">
              <AttachmentManager 
                :noteId="noteId" 
                :tempMode="isNewNote"
                @upload-success="handleAttachmentUploadSuccess"
                @delete-success="handleAttachmentDeleteSuccess"
                @insert-to-editor="insertImageToEditor"
              />
            </div>
          </div>
        </div>
      </div>
      
      <!-- 加载中显示 -->
      <div v-if="loading" class="fixed inset-0 bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm flex flex-col items-center justify-center z-50">
        <div class="relative w-20 h-20">
          <div class="absolute top-0 left-0 w-full h-full rounded-full border-4 border-primary-200 dark:border-primary-900"></div>
          <div class="absolute top-0 left-0 w-full h-full rounded-full border-t-4 border-primary-500 dark:border-primary-400 animate-spin"></div>
        </div>
        <p class="mt-4 text-primary-600 dark:text-primary-400 font-medium">正在加载笔记内容...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useNotesStore } from '@/stores/notes';
import { useTagsStore } from '@/stores/tags';
import { useAttachmentsStore } from '@/stores/attachments';
import { storeToRefs } from 'pinia';
import TagManager from '@/components/TagManager.vue';
import AttachmentManager from '@/components/AttachmentManager.vue';
import MarkdownRenderer from '@/components/MarkdownRenderer.vue';
import { buildAttachmentUrl, createTempImageUrl, uploadImage, loadAuthenticatedImage } from '@/utils/imageHelper';
import toast from '@/utils/toast';

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
const newTag = ref('');
const suggestedTags = ref([]);

// 表单数据
const noteForm = reactive({
  title: '',
  content: '',
  tags: [],
  is_public: false,
  tempAttachments: []
});

// 编辑器状态
const contentTextarea = ref(null);
const imageInput = ref(null);

// 从附件管理器插入图片到编辑器
const insertImageToEditor = (attachment) => {
  if (!contentTextarea.value || !attachment) return;
  
  const markdownImage = `![${attachment.filename}](${attachment.url || `/api/attachments/${attachment.id}`})`;
  
  const textarea = contentTextarea.value;
  const cursorPosition = textarea.selectionStart;
  
  // 在光标位置插入图片链接
  let textBeforeCursor = noteForm.content.substring(0, cursorPosition);
  let textAfterCursor = noteForm.content.substring(cursorPosition);
  
  // 确保图片前后有空行
  if (cursorPosition > 0 && textBeforeCursor.charAt(textBeforeCursor.length - 1) !== '\n') {
    textBeforeCursor += '\n\n';
  } else if (cursorPosition > 0 && !textBeforeCursor.endsWith('\n\n')) {
    textBeforeCursor += '\n';
  }
  
  // 检查光标后是否需要添加换行
  if (textAfterCursor.length > 0 && textAfterCursor.charAt(0) !== '\n') {
    textAfterCursor = '\n\n' + textAfterCursor;
  } else if (textAfterCursor.length > 0 && !textAfterCursor.startsWith('\n\n')) {
    textAfterCursor = '\n' + textAfterCursor;
  }
  
  // 更新内容，插入图片链接
  noteForm.content = textBeforeCursor + markdownImage + textAfterCursor;
  
  // 更新光标位置
  nextTick(() => {
    const newPosition = textBeforeCursor.length + markdownImage.length;
    textarea.selectionStart = newPosition;
    textarea.selectionEnd = newPosition;
    textarea.focus();
  });
  
  toast.success(`已插入图片 ${attachment.filename}`);
};

// 添加标签
const addTag = (tag) => {
  let tagToAdd;
  
  // 如果是字符串输入（从输入框添加）
  if (typeof tag === 'string' || !tag) {
    tagToAdd = newTag.value.trim();
    
    // 清空输入框
    newTag.value = '';
  } 
  // 如果是标签对象（从API或建议中添加）
  else if (typeof tag === 'object') {
    tagToAdd = tag.name;
  }
  
  // 验证并添加标签
  if (tagToAdd && !noteForm.tags.includes(tagToAdd)) {
    noteForm.tags.push(tagToAdd);
  }
};

// 从表单中移除标签
const removeTag = (tag) => {
  noteForm.tags = noteForm.tags.filter(t => t !== tag);
};

// 使用AI生成标签建议
const generateTags = async () => {
  if (!noteForm.content.trim()) return;
  
  try {
    const suggestions = await tagsStore.generateTagSuggestions(noteForm.content);
    // 添加建议的标签
    if (suggestions && suggestions.tags) {
      suggestedTags.value = suggestions.tags.filter(tag => 
        !noteForm.tags.includes(tag)
      );
    }
  } catch (error) {
    console.error('生成标签建议失败:', error);
  }
};

// 添加建议的标签
const addSuggestedTag = (tag) => {
  if (tag && !noteForm.tags.includes(tag)) {
    noteForm.tags.push(tag);
    // 从建议列表中移除已添加的标签
    suggestedTags.value = suggestedTags.value.filter(t => t !== tag);
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
async function saveNote() {
  isSaving.value = true;
  
  try {
    // 基本验证
    if (!noteForm.title) {
      toast.error('请输入笔记标题');
      isSaving.value = false;
      return;
    }
    
    console.log('准备保存笔记:', { ...noteForm });
    
    // 检查内容中的blob URL，替换为实际的API URL
    const content = replaceBlobUrlsWithApiUrls(noteForm.content);
    
    // 创建要保存的笔记对象
    const noteToSave = {
      title: noteForm.title,
      content: content,
      tags: noteForm.tags,
      is_public: noteForm.is_public
    };
    
    let savedNoteId;
    
    // 如果是新建笔记
    if (isNewNote) {
      console.log('创建新笔记');
      
      const createdNote = await notesStore.createNote(noteToSave);
      savedNoteId = createdNote.id;
      
      // 处理临时附件
      if (noteForm.tempAttachments && noteForm.tempAttachments.length > 0) {
        console.log(`处理 ${noteForm.tempAttachments.length} 个临时附件`);
        
        for (const tempAttachment of noteForm.tempAttachments) {
          try {
            if (tempAttachment.id) {
              // 调用关联API
              const result = await attachmentsStore.associateTempAttachment(tempAttachment.id, savedNoteId);
              
              if (result.success) {
                console.log(`附件 ${tempAttachment.file.name} 已关联到笔记`);
                
                // 如果服务器返回了新的URL，替换内容中的URL
                if (result.url && tempAttachment.serverUrl && result.url !== tempAttachment.serverUrl) {
                  noteToSave.content = noteToSave.content.replace(
                    new RegExp(tempAttachment.serverUrl.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, "\\$1"), 'g'),
                    result.url
                  );
                  console.log(`已更新附件URL: ${tempAttachment.serverUrl} -> ${result.url}`);
                }
              } else {
                console.error('关联临时附件失败:', result.message);
                toast.warning(`附件 ${tempAttachment.file.name} 关联失败`);
              }
            }
          } catch (error) {
            console.error('关联临时附件失败:', error);
            toast.warning(`附件 ${tempAttachment.file.name} 关联失败: ${error.message}`);
          }
        }
      }
      
      // 清空临时附件列表
      noteForm.tempAttachments = [];
      
      // 提示保存成功
      toast.success('笔记创建成功');
      
      // 导航到新笔记页面
      router.push({ name: 'note-detail', params: { id: savedNoteId } });
    } else {
      // 更新现有笔记
      console.log('更新现有笔记:', noteId);
      
      await notesStore.updateNote(noteId, noteToSave);
      savedNoteId = noteId;
      
      // 提示保存成功
      toast.success('笔记更新成功');
    }
    
    console.log('笔记保存成功, ID:', savedNoteId);
    return savedNoteId;
  } catch (error) {
    console.error('保存笔记失败:', error);
    toast.error(`保存失败: ${error.message || '未知错误'}`);
    return null;
  } finally {
    isSaving.value = false;
  }
}

// 替换内容中的blob URL为API URL
function replaceBlobUrlsWithApiUrls(content) {
  if (!content) return '';
  
  // 查找内容中的blob URL
  const blobUrlMatches = content.match(/!\[.*?\]\(blob:.*?\)/g);
  if (!blobUrlMatches || blobUrlMatches.length === 0) {
    return content;
  }
  
  console.log(`检测到 ${blobUrlMatches.length} 个blob URL，尝试替换`);
  let updatedContent = content;
  
  // 遍历所有blob URL
  for (const match of blobUrlMatches) {
    console.log('发现blob URL:', match);
    
    // 尝试在tempAttachments中查找对应的服务器URL
    const blobUrlMatch = match.match(/!\[.*?\]\((blob:.*?)\)/);
    if (blobUrlMatch && blobUrlMatch[1]) {
      const blobUrl = blobUrlMatch[1];
      const attachment = noteForm.tempAttachments.find(a => a.tempUrl === blobUrl);
      
      if (attachment && attachment.serverUrl) {
        // 替换blob URL为服务器URL
        updatedContent = updatedContent.replace(blobUrl, attachment.serverUrl);
        console.log(`已替换blob URL为服务器URL: ${blobUrl} -> ${attachment.serverUrl}`);
      } else {
        console.warn(`未找到blob URL ${blobUrl} 对应的服务器URL`);
      }
    }
  }
  
  return updatedContent;
}

// 插入Markdown格式文本
const insertMarkdown = (prefix, suffix, placeholder) => {
  if (!contentTextarea.value) return;
  
  const textarea = contentTextarea.value;
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const selectedText = noteForm.content.substring(start, end);
  
  const textToInsert = selectedText || placeholder;
  const newText = noteForm.content.substring(0, start) + prefix + textToInsert + suffix + noteForm.content.substring(end);
  
  noteForm.content = newText;
  
  // 设置新的光标位置
  nextTick(() => {
    if (selectedText) {
      textarea.selectionStart = start + prefix.length;
      textarea.selectionEnd = end + prefix.length;
    } else {
      textarea.selectionStart = start + prefix.length;
      textarea.selectionEnd = start + prefix.length + placeholder.length;
    }
    textarea.focus();
  });
};

// 打开图片上传器
const openImageUploader = () => {
  if (imageInput.value) {
    imageInput.value.click();
  }
};

// 处理图片上传
async function handleImageUpload(event) {
  console.log('开始上传图片');
  const file = event.target.files[0];
  if (!file) return;

  // 验证文件类型
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
  if (!allowedTypes.includes(file.type)) {
    toast.error('只支持 JPG, PNG, GIF 或 WebP 格式的图片');
    resetFileInput();
    return;
  }

  // 创建唯一占位符ID
  const placeholderId = `img_placeholder_${Date.now()}`;
  const placeholderText = `![图片上传中...](${placeholderId})`;
  
  // 获取光标位置并插入占位符
  const textarea = contentTextarea.value;
  if (!textarea) {
    // 如果无法获取文本区域，直接在内容末尾添加
    noteForm.content += `\n\n${placeholderText}\n\n`;
  } else {
    const cursorPosition = textarea.selectionStart;
    
    // 如果光标所在行开头不是空行，确保图片前后有空行
    let textBeforeCursor = noteForm.content.substring(0, cursorPosition);
    let textAfterCursor = noteForm.content.substring(cursorPosition);
    
    // 检查光标前是否需要添加换行
    if (cursorPosition > 0 && textBeforeCursor.charAt(textBeforeCursor.length - 1) !== '\n') {
      textBeforeCursor += '\n\n';
    } else if (cursorPosition > 0 && !textBeforeCursor.endsWith('\n\n')) {
      textBeforeCursor += '\n';
    }
    
    // 检查光标后是否需要添加换行
    if (textAfterCursor.length > 0 && textAfterCursor.charAt(0) !== '\n') {
      textAfterCursor = '\n\n' + textAfterCursor;
    } else if (textAfterCursor.length > 0 && !textAfterCursor.startsWith('\n\n')) {
      textAfterCursor = '\n' + textAfterCursor;
    }
    
    // 更新内容，插入占位符
    noteForm.content = textBeforeCursor + placeholderText + textAfterCursor;
  }
  
  isLoading.value = true;
  
  try {
    // 第1步：为图片创建临时本地URL用于即时显示
    const tempUrl = createTempImageUrl(file);
    
    // 第2步：立即替换占位符为临时URL，这样用户可以立刻看到图片
    const tempRegex = new RegExp(`!\\[图片上传中...\\]\\(${placeholderId}\\)`, 'g');
    noteForm.content = noteForm.content.replace(tempRegex, `![${file.name}(上传中...)](${tempUrl})`);
    console.log('已将占位符替换为临时预览URL:', tempUrl);
    
    // 第3步：上传图片到服务器
    let uploadResult;
    
    if (isNewNote) {
      // 新建笔记：使用临时附件上传方式
      console.log('新建笔记模式，使用临时附件API');
      try {
        uploadResult = await attachmentsStore.uploadTempAttachment(file);
        console.log('上传临时附件结果:', uploadResult);
        
        // 添加必要的字段，保持与uploadImage相同的返回格式
        if (!uploadResult.success) {
          uploadResult.success = true;
        }
        
        // 如果URL为空，但有ID，则构建一个URL
        if (!uploadResult.url && uploadResult.id) {
          uploadResult.url = `/api/attachments/${uploadResult.id}`;
          console.log('URL为空，使用ID构建URL:', uploadResult.url);
        }
      } catch (error) {
        console.error('临时附件上传失败:', error);
        throw new Error(error.message || '临时附件上传失败');
      }
    } else {
      // 已有笔记：使用常规上传方式
      console.log('现有笔记模式，使用常规上传API');
      uploadResult = await uploadImage(file, noteId);
    }
    
    if (!uploadResult || !uploadResult.success) {
      throw new Error(uploadResult?.message || '上传失败');
    }
    
    console.log('图片上传成功:', uploadResult);
    
    // 第4步：替换临时URL为服务器URL
    const finalRegex = new RegExp(`!\\[${file.name}\\(上传中...\\)\\]\\(${tempUrl.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, "\\$1")}\\)`, 'g');
    
    // 确保URL可以正确显示
    let finalUrl = uploadResult.url || '';
    console.log('服务器返回的URL:', finalUrl);
    
    // 检查URL是否需要调整
    if (finalUrl && finalUrl.startsWith('/api/')) {
      console.log('使用服务器返回的API URL:', finalUrl);
      
      // 如果是需要认证的URL，尝试预加载
      if (finalUrl.startsWith('/api/attachments/')) {
        try {
          // 预加载图片并获取blob URL
          const blobUrl = await loadAuthenticatedImage(finalUrl);
          if (blobUrl && blobUrl !== finalUrl) {
            console.log('已预加载图片并创建blob URL:', blobUrl);
            finalUrl = blobUrl;
          }
        } catch (error) {
          console.error('预加载图片失败:', error);
          // 继续使用原始URL
        }
      }
    } else if (finalUrl && !finalUrl.startsWith('http') && !finalUrl.startsWith('/')) {
      // 如果是相对路径但不以/开头，添加/api/前缀
      finalUrl = `/api/${finalUrl}`;
      console.log('将相对路径调整为API路径:', finalUrl);
    } else if (finalUrl && finalUrl.startsWith('/') && !finalUrl.startsWith('/api/')) {
      // 如果以/开头但不是/api/开头，添加/api前缀
      finalUrl = `/api${finalUrl}`;
      console.log('添加API前缀到路径:', finalUrl);
    } else if (!finalUrl && uploadResult.id) {
      // 如果URL为空但有ID，使用ID构建URL
      finalUrl = `/api/attachments/${uploadResult.id}`;
      console.log('URL为空，使用ID构建URL:', finalUrl);
      
      // 预加载需要认证的URL
      try {
        const blobUrl = await loadAuthenticatedImage(finalUrl);
        if (blobUrl && blobUrl !== finalUrl) {
          console.log('已预加载图片并创建blob URL:', blobUrl);
          finalUrl = blobUrl;
        }
      } catch (error) {
        console.error('预加载图片失败:', error);
        // 继续使用原始URL
      }
    } else if (!finalUrl) {
      // 如果URL为空，使用一个默认或占位符URL
      console.error('上传结果中缺少URL和ID');
      finalUrl = '/api/attachments/placeholder';
      console.log('使用占位符URL:', finalUrl);
    }
    
    noteForm.content = noteForm.content.replace(finalRegex, `![${file.name}](${finalUrl})`);
    console.log('已将临时URL替换为服务器URL:', finalUrl);
    
    // 第5步：如果是新建笔记，保存临时附件信息以便后续关联
    if (isNewNote && uploadResult.id) {
      noteForm.tempAttachments.push({
        id: uploadResult.id,
        file: file,
        tempUrl: tempUrl,
        serverUrl: uploadResult.url
      });
      console.log('已将临时附件信息添加到列表:', noteForm.tempAttachments);
    }
    
    // 释放临时URL
    URL.revokeObjectURL(tempUrl);
    
    toast.success('图片上传成功');
  } catch (error) {
    console.error('图片上传失败:', error);
    // 移除占位符或临时预览
    const regex = new RegExp(`!\\[.+?\\]\\((?:${placeholderId}|blob:.+?)\\)`, 'g');
    noteForm.content = noteForm.content.replace(regex, '');
    toast.error(`图片上传失败: ${error.message || '未知错误'}`);
  } finally {
    isLoading.value = false;
    resetFileInput();
  }
}

// 重置文件输入
const resetFileInput = () => {
  if (imageInput.value) {
    imageInput.value.value = '';
  }
};

// 取消编辑
const cancelEdit = () => {
  // 确认是否放弃编辑
  if (noteForm.title || noteForm.content || noteForm.tags.length > 0) {
    if (!confirm('确定要放弃当前编辑吗？所有更改将丢失。')) {
      return;
    }
  }
  
  // 返回上一页
  router.back();
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
      noteForm.is_public = currentNote.value.is_public;
    } else {
      // 笔记不存在，重定向到笔记列表
      router.push({ name: 'notes' });
    }
  }
  loading.value = false;
});
</script>

<style scoped>
/* 基础样式 */
.editor-btn {
  @apply p-1.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors;
}

.editor-btn:hover {
  @apply transform scale-110 text-primary-600 dark:text-primary-400;
}

/* 输入框与文本区域样式 */
textarea, input[type="text"] {
  @apply transition-all duration-200 ease-in-out;
}

textarea:focus, input[type="text"]:focus {
  @apply shadow-md;
}

textarea {
  @apply font-mono text-base leading-relaxed tracking-wide;
}

/* 卡片样式 */
.rounded-xl {
  @apply transition-shadow duration-200;
}

.rounded-xl:hover {
  @apply shadow-md;
}

/* 标签样式 */
.bg-primary-50 {
  @apply transition-colors duration-200;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  @apply w-2;
}

::-webkit-scrollbar-track {
  @apply bg-transparent;
}

::-webkit-scrollbar-thumb {
  @apply bg-gray-300 dark:bg-gray-700 rounded-full;
}

::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-400 dark:bg-gray-600;
}

/* 按钮交互效果 */
button {
  @apply transition-all duration-200 ease-in-out;
}

button:active:not(:disabled) {
  @apply transform scale-95;
}

/* 加载动画 */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.note-edit-view {
  animation: fadeIn 0.3s ease-out;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* 响应式调整 */
@media (max-width: 1023px) {
  .prose {
    @apply text-sm;
  }
}

/* 暗色模式特殊处理 */
.dark .bg-white {
  @apply bg-opacity-[0.03] backdrop-blur-sm;
}
</style> 