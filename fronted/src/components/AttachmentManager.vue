<template>
  <div class="attachment-manager" @click.stop @submit.prevent>
    <!-- 文件上传区域 -->
    <div class="file-upload-area mb-4" v-if="canUpload" @click.stop>
      <div
        class="border-2 border-dashed dark:border-gray-600 rounded-lg p-4 text-center hover:border-primary-400 dark:hover:border-primary-500 transition-colors"
        :class="{ 'bg-primary-50 dark:bg-primary-900/20': isDragging }"
        @dragenter.prevent.stop="isDragging = true"
        @dragover.prevent.stop="isDragging = true"
        @dragleave.prevent.stop="isDragging = false"
        @drop.prevent.stop="handleFileDrop"
        @click.stop
      >
        <input
          type="file"
          ref="fileInput"
          @change.stop="handleFileChange"
          class="hidden"
          multiple
        />
        <div class="flex flex-col items-center justify-center py-3" @click.stop>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-400 dark:text-gray-500 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          <p class="text-gray-700 dark:text-gray-300 mb-1 font-medium">拖放文件到此处上传</p>
          <p class="text-gray-500 dark:text-gray-400 text-sm mb-2">或</p>
          <button
            type="button"
            @click.stop.prevent="triggerFileInput"
            class="btn btn-primary btn-sm text-sm"
            :disabled="isUploading"
          >
            {{ isUploading ? "上传中..." : "选择文件" }}
          </button>
        </div>

        <!-- 上传进度条 -->
        <div v-if="isUploading" class="mt-4">
          <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full bg-primary-500 transition-all duration-300"
              :style="{ width: `${uploadProgress}%` }"
            ></div>
          </div>
          <div class="text-xs text-gray-500 dark:text-gray-400 mt-1 text-right">
            {{ uploadProgress }}%
          </div>
        </div>
      </div>
    </div>

    <!-- 附件列表 -->
    <div v-if="displayAttachments.length > 0">
      <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">附件列表</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
        <div
          v-for="attachment in displayAttachments"
          :key="attachment.id"
          class="attachment-item border dark:border-gray-700 rounded-lg overflow-hidden bg-white dark:bg-gray-800 shadow-sm hover:shadow-md transition-shadow"
        >
          <!-- 图片附件 -->
          <div v-if="isImage(attachment)" class="relative">
            <div class="w-full h-32 bg-gray-100 dark:bg-gray-700 flex items-center justify-center overflow-hidden">
              <!-- 加载中状态 -->
              <img
                v-if="imageLoading[attachment.id]"
                class="w-6 h-6 animate-spin text-primary-500"
                src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJjdXJyZW50Q29sb3IiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIj48Y2lyY2xlIGN4PSIxMiIgY3k9IjEyIiByPSIxMCIgb3BhY2l0eT0iMC4yNSIvPjxwYXRoIGQ9Ik0xNi4yNCAxNy43OEE5IDkgMCAwIDEgNS43NiA2LjIyIiBvcGFjaXR5PSIwLjc1Ii8+PC9zdmc+"
                alt="加载中"
              />
              <!-- 图片内容 - 正常模式 -->
              <img
                v-else-if="!imageErrors[attachment.id] && attachmentBlobUrls[attachment.id] && !tempMode"
                :src="attachmentBlobUrls[attachment.id]"
                :alt="attachment.filename"
                class="w-full h-full object-cover cursor-pointer"
                @click.stop.prevent="openImageViewer(attachment)"
                loading="lazy"
                @error="handleAttachmentImageError(attachment)"
              />
              <!-- 图片内容 - 临时模式 -->
              <img
                v-else-if="tempMode && attachment.previewUrl"
                :src="attachment.previewUrl"
                :alt="attachment.filename"
                class="w-full h-full object-cover cursor-pointer"
                @click.stop.prevent="openImageViewer(attachment)"
                loading="lazy"
              />
              <!-- 加载失败状态 -->
              <div v-else-if="imageErrors[attachment.id]" class="absolute inset-0 flex flex-col items-center justify-center bg-gray-200 dark:bg-gray-700">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-400 dark:text-gray-500 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <p class="text-sm text-gray-500 dark:text-gray-400">图片加载失败</p>
                <button
                  @click.stop.prevent="reloadAttachmentImage(attachment)"
                  class="mt-2 text-xs text-white bg-primary-500 hover:bg-primary-600 px-2 py-1 rounded"
                >
                  重试
                </button>
              </div>
            </div>
            <div class="absolute top-1 right-1 flex space-x-1">
              <button
                type="button"
                @click.stop.prevent="downloadAttachment(attachment)"
                class="bg-gray-800/70 hover:bg-gray-900/70 text-white rounded p-1 transition-colors"
                title="下载"
                v-if="!tempMode"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </button>
              <button
                type="button"
                v-if="canDelete"
                @click.stop.prevent="confirmDeleteAttachment(attachment)"
                class="bg-red-500/70 hover:bg-red-600/70 text-white rounded p-1 transition-colors"
                title="删除"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
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
            <div class="flex space-x-1">
              <button
                type="button"
                v-if="!tempMode"
                @click.stop.prevent="downloadAttachment(attachment)"
                class="bg-primary-500 hover:bg-primary-600 text-white rounded p-1 transition-colors"
                title="下载"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </button>
              <button
                type="button"
                v-if="canDelete"
                @click.stop.prevent="confirmDeleteAttachment(attachment)"
                class="bg-red-500 hover:bg-red-600 text-white rounded p-1 transition-colors"
                title="删除"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 没有附件时的提示 -->
    <div v-else-if="!isLoading && !isUploading" class="text-center py-3">
      <p class="text-gray-500 dark:text-gray-400 text-sm">
        {{ canUpload ? "暂无附件，可以上传文件" : "暂无附件" }}
      </p>
    </div>

    <!-- 删除确认对话框 -->
    <div v-if="showDeleteConfirm" class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50 p-4" @click.self.stop="showDeleteConfirm = false">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl p-4 sm:p-6 max-w-md w-full border dark:border-gray-700" @click.stop>
        <h3 class="text-lg font-bold mb-3 text-gray-800 dark:text-gray-100">确认删除</h3>
        <p class="mb-4 text-gray-600 dark:text-gray-300">
          您确定要删除附件 "{{ attachmentToDelete?.filename }}" 吗？此操作无法撤销。
        </p>
        <div class="flex justify-end space-x-2">
          <button 
            type="button"
            @click.stop.prevent="showDeleteConfirm = false" 
            class="btn btn-secondary dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600"
          >
            取消
          </button>
          <button 
            type="button"
            @click.stop.prevent="deleteAttachmentConfirmed" 
            class="btn bg-red-600 dark:bg-red-700 text-white hover:bg-red-700 dark:hover:bg-red-800 transition-colors"
          >
            确认删除
          </button>
        </div>
      </div>
    </div>

    <!-- 图片查看器 -->
    <div v-if="showImageViewer" class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-90 p-4" @click.self="closeImageViewer">
      <div class="relative w-full h-full flex flex-col items-center justify-center">
        <div class="absolute top-4 right-4 flex space-x-2 z-10">
          <button
            type="button"
            @click.stop="downloadCurrentImage"
            class="bg-gray-800/50 hover:bg-gray-700/50 text-white rounded-full p-3 transition-colors"
            title="下载图片"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
          </button>
          <button
            type="button"
            @click.stop="closeImageViewer"
            class="bg-gray-800/50 hover:bg-gray-700/50 text-white rounded-full p-3 transition-colors"
            title="关闭"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- 加载中提示 -->
        <div v-if="imageLoading" class="flex flex-col items-center">
          <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-white"></div>
          <p class="text-white mt-4">图片加载中...</p>
        </div>
        
        <!-- 图片加载错误提示 -->
        <div v-if="imageError" class="text-center text-white">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto mb-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <p>图片加载失败</p>
          <button 
            @click.stop="reloadImage" 
            class="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded text-white transition-colors"
          >
            重新加载
          </button>
        </div>
        
        <!-- 图片 -->
        <img
          v-if="!imageLoading && !imageError"
          :src="currentImageUrl"
          :alt="currentImageAttachment?.filename"
          class="max-w-full max-h-full object-contain touch-manipulation"
          @load="handleImageLoaded"
          @error="handleImageError"
        />
        
        <!-- 图片名称 -->
        <div v-if="!imageLoading && !imageError && currentImageAttachment" class="absolute bottom-4 text-center text-white">
          <p class="px-10">{{ currentImageAttachment.filename }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, onMounted, computed, onUnmounted } from 'vue';
import { useAttachmentsStore } from '@/stores/attachments';
import { storeToRefs } from 'pinia';

const props = defineProps({
  noteId: {
    type: [Number, String],
    required: true
  },
  canUpload: {
    type: Boolean,
    default: true
  },
  canDelete: {
    type: Boolean,
    default: true
  },
  tempMode: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['upload-success', 'upload-error', 'delete-success']);

const attachmentsStore = useAttachmentsStore();
const { attachments, isLoading, uploadProgress } = storeToRefs(attachmentsStore);

// 本地状态
const fileInput = ref(null);
const isDragging = ref(false);
const isUploading = ref(false);
const showDeleteConfirm = ref(false);
const attachmentToDelete = ref(null);
const showImageViewer = ref(false);
const currentImageAttachment = ref(null);
const imageLoading = ref({});
const imageErrors = ref({});
const imageError = ref(false);
const currentImageUrl = ref('');
const attachmentBlobUrls = ref({});
const tempAttachments = ref([]); // 暂存附件数组

// 触发文件选择
const triggerFileInput = (e) => {
  // 防止事件冒泡和表单提交
  e?.preventDefault();
  e?.stopPropagation();
  
  console.log('触发文件选择');
  // 直接访问DOM元素来触发点击
  if (fileInput.value) {
    console.log('点击文件输入框');
    // 使用延迟执行，确保DOM更新完成
    setTimeout(() => {
      fileInput.value.click();
    }, 10);
  } else {
    console.error('文件输入框不存在');
  }
};

// 处理文件拖放
const handleFileDrop = (event) => {
  event.preventDefault();
  event.stopPropagation();
  
  isDragging.value = false;
  const files = event.dataTransfer.files;
  console.log('文件拖放，文件数量:', files.length);
  if (files.length > 0) {
    uploadFiles(files);
  }
};

// 处理文件选择
const handleFileChange = (event) => {
  event.preventDefault();
  event.stopPropagation();
  
  const files = event.target.files;
  console.log('文件选择，文件数量:', files.length);
  if (files.length > 0) {
    uploadFiles(files);
  }
};

// 上传文件
const uploadFiles = async (files) => {
  if ((!props.noteId && !props.tempMode) || isUploading.value) {
    console.log('无法上传: 笔记ID不存在或已有上传进行中');
    return;
  }

  console.log('开始上传文件，数量:', files.length);
  isUploading.value = true;
  
  try {
    // 一次上传一个文件
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      console.log(`上传文件 ${i+1}/${files.length}:`, file.name);
      
      if (props.tempMode) {
        // 临时模式，只在本地保存文件
        console.log('临时模式，暂存文件:', file.name);
        
        // 生成预览URL
        const previewUrl = URL.createObjectURL(file);
        
        // 创建临时附件对象
        const tempAttachment = {
          id: 'temp_' + Date.now() + '_' + i, // 生成临时ID
          file: file, // 存储原始文件对象
          filename: file.name,
          filetype: file.type,
          filesize: file.size,
          previewUrl: previewUrl,
          createdAt: new Date().toISOString()
        };
        
        // 添加到临时数组
        tempAttachments.value.push(tempAttachment);
        
        // 触发上传成功事件，将文件数据传给父组件
        emit('upload-success', null, tempAttachment);
      } else {
        // 正常模式，上传到服务器
        await attachmentsStore.uploadAttachment(file, props.noteId);
        
        // 触发上传成功事件
        emit('upload-success');
      }
    }
    
    console.log('所有文件处理完成');
    
    // 清空文件输入控件，以便下次可以选择相同的文件
    if (fileInput.value) {
      fileInput.value.value = ''; 
    }
  } catch (error) {
    console.error('上传文件失败:', error);
    emit('upload-error', error);
    alert('上传文件失败: ' + (error.message || '请稍后重试'));
  } finally {
    isUploading.value = false;
  }
};

// 确认删除附件
const confirmDeleteAttachment = (attachment, event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  console.log('确认删除附件:', attachment.id, attachment.filename);
  attachmentToDelete.value = attachment;
  showDeleteConfirm.value = true;
};

// 确认删除的回调函数
const deleteAttachmentConfirmed = async (event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  if (!attachmentToDelete.value) {
    console.log('没有选择要删除的附件');
    return;
  }
  
  try {
    if (props.tempMode) {
      // 临时模式下删除附件
      console.log('临时模式，删除暂存附件:', attachmentToDelete.value.id);
      
      // 从临时数组中移除
      tempAttachments.value = tempAttachments.value.filter(
        attachment => attachment.id !== attachmentToDelete.value.id
      );
      
      // 如果有blob URL，释放它
      if (attachmentToDelete.value.previewUrl) {
        URL.revokeObjectURL(attachmentToDelete.value.previewUrl);
      }
      
      console.log('临时附件删除成功');
      showDeleteConfirm.value = false;
      
      // 保存删除的附件引用
      const deletedAttachment = { ...attachmentToDelete.value };
      attachmentToDelete.value = null;
      
      // 通知父组件删除成功，并传递删除的附件数据
      emit('delete-success', null, deletedAttachment);
    } else {
      // 常规模式下删除附件
      console.log('开始删除附件:', attachmentToDelete.value.id);
      await attachmentsStore.deleteAttachment(attachmentToDelete.value.id);
      console.log('附件删除成功');
      showDeleteConfirm.value = false;
      attachmentToDelete.value = null;
      // 仅通知父组件删除成功，不进行页面跳转
      emit('delete-success');
    }
  } catch (error) {
    console.error('删除附件失败:', error);
    alert('删除附件失败: ' + (error.message || '请稍后重试'));
  }
};

// 下载附件
const downloadAttachment = async (attachment, event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  if (!attachment) {
    console.error('无法下载：附件对象为空');
    return;
  }
  
  console.log('下载附件:', attachment.id, attachment.filename);
  
  try {
    // 使用附件存储中的带认证下载方法
    await attachmentsStore.downloadAttachment(attachment.id, attachment.filename);
    console.log('附件下载成功');
  } catch (error) {
    console.error('下载附件时出错:', error);
    alert('下载附件失败: ' + (error.message || '请稍后重试'));
  }
};

// 打开图片查看器
const openImageViewer = (attachment, event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  console.log('打开图片查看器:', attachment.id, attachment.filename);
  
  // 重置状态
  imageLoading.value = true;
  imageError.value = false;
  currentImageAttachment.value = attachment;
  
  if (props.tempMode && attachment.previewUrl) {
    // 临时模式，直接使用预览URL
    console.log('使用临时预览URL:', attachment.previewUrl);
    currentImageUrl.value = attachment.previewUrl;
    imageLoading.value = false;
    
    // 打开查看器
    showImageViewer.value = true;
    
    // 禁止body滚动
    document.body.style.overflow = 'hidden';
    return;
  }
  
  // 获取带认证的图片URL
  const token = localStorage.getItem('token');
  if (!token) {
    console.error('未提供认证令牌，无法预览图片');
    imageError.value = true;
    imageLoading.value = false;
    return;
  }
  
  // 创建一个带认证信息的Blob URL
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
  const xhr = new XMLHttpRequest();
  xhr.open('GET', `${baseUrl}/attachments/${attachment.id}`, true);
  xhr.setRequestHeader('Authorization', `Bearer ${token}`);
  xhr.responseType = 'blob';
  
  xhr.onload = function() {
    if (this.status === 200) {
      const blob = new Blob([this.response], { type: attachment.filetype });
      currentImageUrl.value = URL.createObjectURL(blob);
      imageLoading.value = false;
    } else {
      console.error('获取图片失败:', this.status);
      imageError.value = true;
      imageLoading.value = false;
    }
  };
  
  xhr.onerror = function() {
    console.error('图片加载请求失败');
    imageError.value = true;
    imageLoading.value = false;
  };
  
  xhr.send();
  
  // 打开查看器
  showImageViewer.value = true;
  
  // 禁止body滚动
  document.body.style.overflow = 'hidden';
};

// 关闭图片查看器
const closeImageViewer = (event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  console.log('关闭图片查看器');
  showImageViewer.value = false;
  currentImageAttachment.value = null;
  currentImageUrl.value = '';
  imageLoading.value = false;
  imageError.value = false;
  
  // 恢复body滚动
  document.body.style.overflow = '';
};

// 下载当前查看的图片
const downloadCurrentImage = (event) => {
  // 阻止事件冒泡和表单提交
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }
  
  if (currentImageAttachment.value) {
    console.log('下载当前查看的图片:', currentImageAttachment.value.id);
    downloadAttachment(currentImageAttachment.value);
  }
};

// 获取附件URL
const getAttachmentUrl = (id) => {
  if (!id) {
    console.warn('无法获取附件URL：ID不存在');
    return '';
  }
  
  const url = attachmentsStore.getAttachmentUrl(id);
  console.log('获取附件URL:', id, url);
  return url;
};

// 判断是否为图片
const isImage = (attachment) => {
  if (!attachment || !attachment.filetype) return false;
  return attachment.filetype.startsWith('image/');
};

// 格式化文件大小
const formatFileSize = (size) => {
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

// 计算属性：显示的附件列表（根据模式选择服务器附件或临时附件）
const displayAttachments = computed(() => {
  if (props.tempMode) {
    return tempAttachments.value;
  } else {
    return attachments.value;
  }
});

// 加载笔记的附件
onMounted(async () => {
  console.log('附件管理器挂载，笔记ID:', props.noteId, '临时模式:', props.tempMode);
  if (props.noteId && !props.tempMode) {
    try {
      await attachmentsStore.fetchAttachments(props.noteId);
      console.log('附件加载完成，数量:', attachments.value.length);
      
      // 预加载所有图片类型的附件
      for (const attachment of attachments.value) {
        if (isImage(attachment)) {
          preloadAttachmentImage(attachment);
        }
      }
    } catch (error) {
      console.error('加载附件失败:', error);
    }
  }
});

// 预加载图片附件，创建Blob URL
const preloadAttachmentImage = (attachment) => {
  if (!attachment || !attachment.id) return;
  
  imageLoading.value[attachment.id] = true;
  imageErrors.value[attachment.id] = false;
  
  const token = localStorage.getItem('token');
  if (!token) {
    console.error('未提供认证令牌，无法加载图片');
    imageErrors.value[attachment.id] = true;
    imageLoading.value[attachment.id] = false;
    return;
  }
  
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
  const xhr = new XMLHttpRequest();
  xhr.open('GET', `${baseUrl}/attachments/${attachment.id}`, true);
  xhr.setRequestHeader('Authorization', `Bearer ${token}`);
  xhr.responseType = 'blob';
  
  xhr.onload = function() {
    if (this.status === 200) {
      const blob = new Blob([this.response], { type: attachment.filetype });
      const blobUrl = URL.createObjectURL(blob);
      attachmentBlobUrls.value[attachment.id] = blobUrl;
      imageLoading.value[attachment.id] = false;
      console.log('图片预加载成功:', attachment.id);
    } else {
      console.error('预加载图片失败:', this.status);
      imageErrors.value[attachment.id] = true;
      imageLoading.value[attachment.id] = false;
    }
  };
  
  xhr.onerror = function() {
    console.error('预加载图片请求失败:', attachment.id);
    imageErrors.value[attachment.id] = true;
    imageLoading.value[attachment.id] = false;
  };
  
  xhr.send();
};

// 处理图片加载
const handleImageLoaded = () => {
  console.log('图片加载成功');
  imageLoading.value = false;
};

// 处理图片加载错误
const handleImageError = () => {
  console.error('图片加载失败');
  imageError.value = true;
};

// 重新加载图片
const reloadImage = () => {
  console.log('重新加载图片');
  imageLoading.value = true;
  imageError.value = false;
  currentImageUrl.value = '';
  if (currentImageAttachment.value) {
    console.log('重新加载当前查看的图片:', currentImageAttachment.value.id);
    downloadAttachment(currentImageAttachment.value);
  }
};

// 获取带认证的图片URL
const getAuthImageUrl = (id) => {
  if (!id) return '';
  const token = localStorage.getItem('token');
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
  const url = `${baseUrl}/attachments/${id}?token=${token}`;
  console.log('构建附件图片URL:', url);
  return url;
};

// 处理图片加载错误
const handleAttachmentImageError = (attachment) => {
  console.error('附件图片加载失败:', attachment.id);
  imageErrors.value[attachment.id] = true;
};

// 重新加载图片
const reloadAttachmentImage = (attachment) => {
  console.log('重新加载附件图片:', attachment.id);
  
  // 如果之前有创建的Blob URL，需要释放
  if (attachmentBlobUrls.value[attachment.id]) {
    URL.revokeObjectURL(attachmentBlobUrls.value[attachment.id]);
    attachmentBlobUrls.value[attachment.id] = null;
  }
  
  // 重新加载图片
  preloadAttachmentImage(attachment);
};

// 在组件卸载时释放所有Blob URL
onUnmounted(() => {
  console.log('附件管理器卸载，释放Blob URL资源');
  // 释放所有Blob URL
  Object.values(attachmentBlobUrls.value).forEach(url => {
    if (url && url.startsWith('blob:')) {
      URL.revokeObjectURL(url);
    }
  });
});
</script>

<style scoped>
.file-upload-area {
  transition: all 0.3s ease;
}

.attachment-item {
  transition: all 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style> 