<template>
  <div class="w-full">
    <!-- 笔记标题和操作按钮 -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4 gap-3">
      <h1 class="text-xl md:text-2xl font-bold text-text-primary dark:text-gray-50 break-words">{{ note.title }}</h1>
      <div class="flex items-center space-x-2 shrink-0">
        <button @click="$emit('edit')" class="btn btn-secondary flex items-center text-sm sm:text-base">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
          编辑
        </button>
        <button @click="$emit('delete')" class="btn flex items-center bg-red-600 dark:bg-red-700 text-white hover:bg-red-700 dark:hover:bg-red-800 transition-colors text-sm sm:text-base">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
          删除
        </button>
      </div>
    </div>

    <!-- 笔记元信息 -->
    <div class="flex flex-wrap items-center text-xs sm:text-sm text-text-tertiary dark:text-gray-400 mb-4">
      <span>创建于: {{ formatDate(note.created_at) }}</span>
      <span class="mx-2">•</span>
      <span>更新于: {{ formatDate(note.updated_at, true) }}</span>
      <span class="mx-2">•</span>
      <span class="flex items-center">
        <span v-if="note.is_public" class="inline-flex items-center text-green-600 dark:text-green-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          公开笔记
        </span>
        <span v-else class="inline-flex items-center text-gray-600 dark:text-gray-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          私密笔记
        </span>
      </span>
    </div>

    <!-- 笔记标签 -->
    <div class="flex flex-wrap gap-2 mb-5 sm:mb-6">
      <span 
        v-for="tag in note.tags" 
        :key="tag.id" 
        class="px-2 py-1 bg-primary-100 dark:bg-primary-600 text-primary-800 dark:text-white text-xs rounded-full inline-flex items-center shadow-sm"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-white mr-1 inline-block"></span>
        {{ tag.name }}
      </span>
    </div>

    <!-- 笔记内容 - 使用MarkdownRenderer组件 -->
    <MarkdownRenderer 
      :content="note.content" 
      class="prose dark:prose-invert sm:prose-sm md:prose-base max-w-none overflow-hidden dark:bg-gray-800/30 dark:p-4 dark:rounded-lg animate-fadeIn"
    />

    <!-- 附件区域 -->
    <div v-if="note.attachments && note.attachments.length > 0" class="mt-6 border-t dark:border-gray-700 pt-4">
      <h3 class="text-lg font-semibold mb-3 text-gray-800 dark:text-gray-200">附件</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
        <div 
          v-for="attachment in note.attachments" 
          :key="attachment.id"
          class="attachment-item border dark:border-gray-700 rounded-lg overflow-hidden bg-white dark:bg-gray-800 shadow-sm hover:shadow-md transition-shadow"
        >
          <!-- 图片附件 -->
          <div v-if="isImage(attachment)" class="relative">
            <div class="w-full h-32 bg-gray-100 dark:bg-gray-700 flex items-center justify-center overflow-hidden">
              <img
                v-if="loading"
                class="w-6 h-6 animate-spin text-primary-500"
                src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJjdXJyZW50Q29sb3IiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIj48Y2lyY2xlIGN4PSIxMiIgY3k9IjEyIiByPSIxMCIgb3BhY2l0eT0iMC4yNSIvPjxwYXRoIGQ9Ik0xNi4yNCAxNy43OEE5IDkgMCAwIDEgNS43NiA2LjIyIiBvcGFjaXR5PSIwLjc1Ii8+PC9zdmc+"
                alt="加载中"
              />
              <img
                v-else
                :src="attachmentBlobUrls[attachment.id] || ''"
                :alt="attachment.filename"
                class="w-full h-full object-cover cursor-pointer"
                @click="openImageViewer(attachment)"
                loading="lazy"
                @error="handleImageError(attachment)"
                ref="imageRefs"
              />
              <div v-if="failedImages[attachment.id]" class="absolute inset-0 flex flex-col items-center justify-center bg-gray-200 dark:bg-gray-700">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-400 dark:text-gray-500 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <p class="text-sm text-gray-500 dark:text-gray-400">图片加载失败</p>
                <button
                  @click.stop="reloadImage(attachment)"
                  class="mt-2 text-xs text-white bg-primary-500 hover:bg-primary-600 px-2 py-1 rounded"
                >
                  重试
                </button>
              </div>
            </div>
            <div class="absolute top-1 right-1">
              <button
                @click="downloadAttachment(attachment)"
                class="bg-gray-800/70 hover:bg-gray-900/70 text-white rounded p-1 transition-colors"
                title="下载"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </button>
            </div>
          </div>

          <!-- 非图片附件 -->
          <div v-else class="flex justify-between p-3">
            <div class="flex items-center">
              <div class="file-icon mr-3 text-gray-500 dark:text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div class="flex flex-col">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300 truncate" :title="attachment.filename">
                  {{ attachment.filename }}
                </span>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  {{ formatFileSize(attachment.filesize) }}
                </span>
              </div>
            </div>
            <button
              @click="downloadAttachment(attachment)"
              class="bg-primary-500 hover:bg-primary-600 text-white rounded p-1 transition-colors"
              title="下载"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片查看器 - 新设计 -->
    <div v-if="showImageViewer" class="fixed inset-0 z-50" @click.self="closeImageViewer">
      <!-- 玻璃态背景 -->
      <div class="absolute inset-0 bg-black/90 backdrop-blur-sm"></div>
      
      <!-- 主内容区 -->
      <div class="relative z-10 w-full h-full flex flex-col">
        
        <!-- 顶部工具栏 -->
        <div class="flex items-center justify-between p-4 bg-gradient-to-b from-black/80 to-transparent">
          <h3 class="text-white text-lg font-medium truncate max-w-[60%] opacity-80" v-if="currentImageAttachment">
            {{ currentImageAttachment.filename }}
          </h3>
          
          <div class="flex space-x-3">
            <button 
              v-if="currentImageAttachment"
              @click="downloadAttachment(currentImageAttachment)"
              class="bg-black/25 hover:bg-blue-600/90 text-white p-2.5 rounded-full transition-all shadow-lg"
              title="下载图片"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
            </button>
            <button 
              @click="closeImageViewer"
              class="bg-black/25 hover:bg-red-600/90 text-white p-2.5 rounded-full transition-all shadow-lg"
              title="关闭"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
        
        <!-- 中央内容区 -->
        <div class="flex-grow flex items-center justify-center p-4">
          <!-- 加载中状态 -->
          <div v-if="imageLoading" class="flex flex-col items-center justify-center">
            <!-- 更现代的加载动画 -->
            <div class="relative w-20 h-20">
              <div class="absolute inset-0 rounded-full border-4 border-blue-200/30"></div>
              <div class="absolute inset-0 rounded-full border-4 border-transparent border-t-blue-500 animate-spin"></div>
              <div class="absolute inset-0 rounded-full border-4 border-transparent border-l-blue-400 animate-spin" style="animation-duration: 1.5s;"></div>
            </div>
            <p class="mt-6 text-blue-300 text-lg font-medium animate-pulse">载入图片中...</p>
          </div>
          
          <!-- 图片加载错误 -->
          <div v-else-if="imageError" class="text-center">
            <div class="inline-flex items-center justify-center w-24 h-24 rounded-full bg-red-500/20 mb-6 animate-pulse">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <p class="text-red-200 text-xl font-medium mb-5">图片载入失败</p>
            <button 
              @click="reloadImage(currentImageAttachment)"
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-medium rounded-full transition-colors shadow-lg shadow-blue-900/30"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              重新加载
            </button>
          </div>
          
          <!-- 图片显示 -->
          <div 
            v-else-if="currentImageUrl" 
            class="relative w-full h-full flex items-center justify-center"
            :class="{'cursor-grab': !isZoomed, 'cursor-grabbing': isZoomed}"
          >
            <div class="relative group">
              <img
                :src="currentImageUrl"
                :alt="currentImageAttachment ? currentImageAttachment.filename : ''"
                class="max-h-[85vh] max-w-[90vw] object-contain rounded-lg shadow-2xl transition-all duration-300"
                :class="{'scale-[1.5] rotate-0': isZoomed}"
                @click.stop="toggleZoom"
                @load="handleImageLoaded"
              />
              
            </div>
          
          </div>
        </div>
      </div>
      
      <!-- 移动端专用底部工具栏 - 移除重复的下载和关闭按钮，改为显示当前状态和图片名称 -->
      <div class="fixed bottom-8 left-0 right-0 flex justify-center z-30 sm:hidden">
        <div class="flex items-center px-5 py-2 rounded-full bg-gray-900/70 backdrop-blur-sm shadow-xl">
          <div class="text-white/80 text-sm truncate max-w-[60vw]" v-if="currentImageAttachment">
            {{ currentImageAttachment.filename }}
          </div>
        </div>
      </div>
    </div>
    
    <!-- 移动端返回顶部按钮 -->
    <button 
      v-if="showScrollToTop" 
      @click="scrollToTop" 
      class="fixed bottom-5 right-5 bg-primary-500 dark:bg-primary-600 text-white p-2 rounded-full shadow-lg z-10 sm:hidden"
      aria-label="返回顶部"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
      </svg>
    </button>
    
    <!-- 复制成功提示 -->
    <div 
      v-if="showCopySuccess" 
      class="fixed bottom-4 right-4 bg-primary-600 text-white px-4 py-2 rounded-md shadow-lg z-50 animate-fade-in-out"
    >
      已复制到剪贴板
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, onMounted, watch, onUnmounted, nextTick, computed } from 'vue';
import { useAttachmentsStore } from '@/stores/attachments';
import MarkdownRenderer from '@/components/MarkdownRenderer.vue';

const props = defineProps({
  note: {
    type: Object,
    required: true
  }
});

defineEmits(['edit', 'delete']);

const showScrollToTop = ref(false);
const showCopySuccess = ref(false);
const showImageViewer = ref(false);
const currentImageAttachment = ref(null);
const currentImageUrl = computed(() => {
  if (!currentImageAttachment.value || !currentImageAttachment.value.id) return '';
  return attachmentBlobUrls.value[currentImageAttachment.value.id] || '';
});
const imageLoading = ref(false);
const imageError = ref(false);
const loading = ref(false);
const failedImages = ref({});
const imageRefs = ref([]);
const attachmentBlobUrls = ref({});
const isZoomed = ref(false);

// 附件相关
const attachmentsStore = useAttachmentsStore();

// 格式化日期
const formatDate = (dateString, shortFormat = false) => {
  console.log('NoteDetail 格式化日期，原始值:', dateString);
  
  if (!dateString) {
    console.warn('日期字符串为空或无效');
    return '无日期';
  }
  
  // 尝试修复日期格式
  let fixedDateString = dateString;
  if (typeof dateString === 'string' && dateString.includes('T')) {
    // 确保日期字符串有Z或时区
    if (!dateString.includes('Z') && !dateString.includes('+')) {
      fixedDateString = dateString + 'Z';
    }
  }
  
  const date = new Date(fixedDateString);
  console.log('NoteDetail 转换后的日期对象:', date);
  
  if (isNaN(date.getTime())) {
    console.warn('日期解析失败:', dateString);
    return '日期无效';
  }
  
  try {
    // 在移动端使用更简短的日期格式
    if (shortFormat && window.innerWidth < 640) {
      return date.toLocaleDateString('zh-CN', { 
        month: 'numeric', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
    
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (e) {
    console.error('日期格式化错误:', e);
    
    // 降级方案：手动格式化
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    
    if (shortFormat && window.innerWidth < 640) {
      return `${month}月${day}日 ${hours}:${minutes}`;
    }
    return `${year}年${month}月${day}日 ${hours}:${minutes}`;
  }
};

// 附件相关函数
// 判断是否为图片
const isImage = (attachment) => {
  if (!attachment || !attachment.filetype) return false;
  return attachment.filetype.startsWith('image/');
};

// 格式化文件大小
const formatFileSize = (size) => {
  if (!size) return '未知大小';
  if (size < 1024) {
    return size + ' B';
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(1) + ' KB';
  } else if (size < 1024 * 1024 * 1024) {
    return (size / (1024 * 1024)).toFixed(1) + ' MB';
  } else {
    return (size / (1024 * 1024 * 1024)).toFixed(1) + ' GB';
  }
};

// 获取附件URL
const getAttachmentUrl = (id) => {
  return attachmentsStore.getAttachmentUrl(id);
};

// 下载附件
const downloadAttachment = async (attachment) => {
  try {
    await attachmentsStore.downloadAttachment(attachment.id, attachment.filename);
    console.log('附件下载成功:', attachment.filename);
  } catch (error) {
    console.error('下载附件时出错:', error);
    alert('下载附件失败: ' + (error.message || '请稍后重试'));
  }
};

// 获取带认证的图片URL
const getAuthImageUrl = (id) => {
  if (!id) return '';
  const token = localStorage.getItem('token');
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
  const url = `${baseUrl}/attachments/${id}?token=${token}`;
  console.log('构建图片URL:', url);
  return url;
};

// 处理图片加载错误
const handleImageError = (attachment) => {
  console.error('图片加载失败:', attachment.id);
  failedImages.value[attachment.id] = true;
};

// 预加载图片，创建Blob URL
const preloadImage = async (attachment) => {
  if (!attachment || !attachment.id || !isImage(attachment)) return;
  
  const token = localStorage.getItem('token');
  if (!token) {
    console.error('未提供认证令牌，无法加载图片');
    failedImages.value[attachment.id] = true;
    return;
  }
  
  console.log('预加载图片:', attachment.id);
  
  try {
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
    const response = await fetch(`${baseUrl}/attachments/${attachment.id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const blob = await response.blob();
    const blobUrl = URL.createObjectURL(blob);
    attachmentBlobUrls.value[attachment.id] = blobUrl;
    failedImages.value[attachment.id] = false;
    
    console.log('图片预加载成功:', attachment.id);
    return blobUrl;
  } catch (error) {
    console.error('预加载图片失败:', error);
    failedImages.value[attachment.id] = true;
    return null;
  }
};

// 重新加载图片
const reloadImage = async (attachment) => {
  console.log('重新加载图片:', attachment.id);
  failedImages.value[attachment.id] = false;
  
  // 如果之前有创建的Blob URL，需要释放
  if (attachmentBlobUrls.value[attachment.id]) {
    URL.revokeObjectURL(attachmentBlobUrls.value[attachment.id]);
    attachmentBlobUrls.value[attachment.id] = null;
  }
  
  await preloadImage(attachment);
};

// 打开图片查看器
const openImageViewer = (attachment) => {
  if (attachment.filetype && attachment.filetype.startsWith('image/')) {
    currentImageAttachment.value = attachment;
    showImageViewer.value = true;
    isZoomed.value = false; // 重置缩放状态
    
    // 防止滚动
    document.body.style.overflow = 'hidden';
    
    // 确保图片已加载
    if (!attachmentBlobUrls.value[attachment.id]) {
      imageLoading.value = true;
      imageError.value = false;
      preloadImage(attachment).then(url => {
        if (url) {
          imageLoading.value = false;
        } else {
          imageError.value = true;
          imageLoading.value = false;
        }
      }).catch(() => {
        imageError.value = true;
        imageLoading.value = false;
      });
    } else {
      imageLoading.value = false;
      imageError.value = false;
    }
  }
};

// 关闭图片查看器
const closeImageViewer = () => {
  showImageViewer.value = false;
  isZoomed.value = false;
  
  // 恢复滚动
  document.body.style.overflow = '';
};

// 监听滚动事件，判断是否显示返回顶部按钮
const handleScroll = () => {
  showScrollToTop.value = window.scrollY > 300;
};

// 返回顶部
const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  });
};

// 组件加载时渲染内容和预加载图片
onMounted(() => {
  window.addEventListener('scroll', handleScroll);
  
  // 监听复制成功事件
  window.addEventListener('markdown-copy-success', handleCopySuccess);
  
  // 预加载笔记中的图片附件
  if (props.note && props.note.attachments) {
    props.note.attachments.forEach(attachment => {
      if (isImage(attachment)) {
        preloadImage(attachment);
      }
    });
  }
});

// 组件卸载时移除事件监听和清理资源
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
  
  // 移除复制成功事件监听器
  window.removeEventListener('markdown-copy-success', handleCopySuccess);
  
  // 释放所有Blob URL
  Object.values(attachmentBlobUrls.value).forEach(url => {
    if (url && url.startsWith('blob:')) {
      URL.revokeObjectURL(url);
    }
  });
  
  // 重置状态并恢复滚动
  if (showImageViewer.value) {
    showImageViewer.value = false;
    document.body.style.overflow = '';
  }
});

// 处理复制成功事件
const handleCopySuccess = () => {
  showCopySuccess.value = true;
  setTimeout(() => {
    showCopySuccess.value = false;
  }, 2000);
};

// 处理图片加载完成
const handleImageLoaded = () => {
  console.log('图片加载完成');
  imageLoading.value = false;
};

// 切换图片缩放状态
const toggleZoom = () => {
  isZoomed.value = !isZoomed.value;
  console.log('图片缩放状态:', isZoomed.value ? '已放大' : '已恢复');
};
</script>

<style scoped>
.prose {
  @apply text-text-primary transition-colors;
}

.dark .prose {
  @apply text-gray-100;
}

/* 增强Markdown内容在暗色模式下的可读性 */
:deep(.dark .prose) {
  --tw-prose-headings: theme('colors.gray.100');
  --tw-prose-body: theme('colors.gray.200');
  --tw-prose-bold: theme('colors.gray.50');
  --tw-prose-links: theme('colors.primary.300');
  --tw-prose-code: theme('colors.gray.100');
  --tw-prose-quotes: theme('colors.gray.300');
}

:deep(.prose) {
  --tw-prose-links: theme('colors.primary.600');
}

/* 增强暗色模式下的样式 */
:deep(.dark .prose blockquote) {
  @apply border-l-4 border-primary-700 pl-4 bg-gray-800/50 py-1 rounded-r-md;
}

:deep(.dark .prose h1) {
  @apply text-primary-300 border-b pb-2 border-gray-700;
}

:deep(.dark .prose h2) {
  @apply text-primary-300 opacity-90 border-b pb-1 border-gray-700;
}

:deep(.dark .prose h3) {
  @apply text-primary-300 opacity-80;
}

:deep(.dark .prose strong) {
  @apply text-primary-200;
}

:deep(.dark .prose code) {
  @apply bg-gray-800 text-primary-200 border border-gray-700;
}

:deep(.dark .prose pre) {
  @apply bg-gray-900 border border-gray-700;
}

:deep(.dark .prose pre code) {
  @apply bg-transparent border-0 text-gray-200;
}

:deep(.dark .prose a) {
  @apply text-primary-300 hover:text-primary-200 underline;
}

:deep(.dark .prose table) {
  @apply border-collapse border border-gray-700;
}

:deep(.dark .prose th) {
  @apply bg-gray-800 text-gray-100 border border-gray-700 px-4 py-2;
}

:deep(.dark .prose td) {
  @apply border border-gray-700 px-4 py-2;
}

/* 代码块复制按钮样式 - 这些样式现在由MarkdownRenderer组件提供，可以删除或保留 */
:deep(.code-block-wrapper) {
  position: relative;
}

:deep(.code-copy-btn) {
  transition: opacity 0.2s ease;
}

:deep(.code-copy-btn:hover) {
  @apply bg-primary-600;
}

/* 复制成功提示动画 */
.animate-fade-in-out {
  animation: fadeInOut 2s ease;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translateY(10px); }
  15% { opacity: 1; transform: translateY(0); }
  85% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(-10px); }
}

/* 移动端优化 Markdown 内容 */
@media (max-width: 640px) {
  :deep(.prose img) {
    @apply w-full;
  }
  
  :deep(.prose pre) {
    @apply text-sm overflow-x-auto;
  }
  
  :deep(.prose table) {
    @apply text-sm;
  }
  
  /* 在移动端显示复制按钮 */
  :deep(.code-copy-btn) {
    opacity: 0.7 !important;
  }
}

.animate-fadeIn {
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* 增强暗色模式下标签的可见性 */
.dark :deep(.flex-wrap .bg-primary-600) {
  background-color: rgba(79, 70, 229, 0.9) !important; /* 更亮更饱和的背景 */
  color: #ffffff !important; /* 纯白色文字 */
  font-weight: 600 !important; /* 加粗 */
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.5) !important; /* 添加文字阴影增强可读性 */
  border: 1px solid rgba(255, 255, 255, 0.3) !important; /* 添加亮色边框 */
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3) !important; /* 添加阴影增强层次感 */
  font-size: 0.85rem !important; /* 稍微增大字体 */
  padding: 0.3rem 0.6rem !important; /* 增加内边距 */
}
</style> 