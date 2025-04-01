<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useAttachmentsStore } from '@/stores/attachments';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import LibraryStats from '@/components/LibraryStats.vue';

// 路由器
const router = useRouter();

// 附件store
const attachmentsStore = useAttachmentsStore();
const { libraryAttachments, dateGroups, totalAttachments, isLoading } = storeToRefs(attachmentsStore);

// 环境变量
const isDev = import.meta.env.DEV || false;

// 分页和过滤状态
const currentPage = ref(1);
const pageSize = ref(20);
const selectedFileType = ref('');
const searchKeyword = ref('');

// 文件类型选项
const fileTypeOptions = [
  { value: '', label: '全部文件' },
  { value: 'image/', label: '图片' },
  { value: 'application/pdf', label: 'PDF文档' },
  { value: 'application/msword', label: 'Word文档' },
  { value: 'application/vnd.ms-excel', label: 'Excel表格' },
  { value: 'text/', label: '文本文件' },
  { value: 'application/zip', label: '压缩文件' },
];

// 图片查看器状态
const showImageViewer = ref(false);
const currentImage = ref(null);

// 视图类型（网格/列表）
const viewType = ref('grid'); // 'grid' 或 'list'

// 获取资源库文件
const fetchLibrary = async () => {
  try {
    const response = await attachmentsStore.fetchLibraryAttachments(
      currentPage.value,
      pageSize.value,
      selectedFileType.value
    );
    
    // 如果当前页超过了总页数，回到第一页重新获取
    if (response && response.total > 0 && currentPage.value > Math.ceil(response.total / pageSize.value)) {
      currentPage.value = 1;
      await fetchLibrary();
    }
  } catch (error) {
    console.error('获取资源库文件失败:', error);
  }
};

// 当过滤条件变更时重新加载文件
watch([selectedFileType], () => {
  currentPage.value = 1; // 重置到第一页
  fetchLibrary();
});

// 获取文件类型图标
const getFileIcon = (filetype) => {
  if (!filetype) return 'fas fa-file';
  
  if (filetype.startsWith('image/')) return 'fas fa-file-image';
  if (filetype.startsWith('audio/')) return 'fas fa-file-audio';
  if (filetype.startsWith('video/')) return 'fas fa-file-video';
  if (filetype.startsWith('text/')) return 'fas fa-file-alt';
  if (filetype.includes('pdf')) return 'fas fa-file-pdf';
  if (filetype.includes('word') || filetype.includes('document')) return 'fas fa-file-word';
  if (filetype.includes('excel') || filetype.includes('sheet')) return 'fas fa-file-excel';
  if (filetype.includes('powerpoint') || filetype.includes('presentation')) return 'fas fa-file-powerpoint';
  if (filetype.includes('zip') || filetype.includes('compressed')) return 'fas fa-file-archive';
  
  return 'fas fa-file';
};

// 获取文件类型标签
const getFileTypeLabel = (filetype) => {
  if (!filetype) return '未知类型';
  
  if (filetype.startsWith('image/')) return '图片';
  if (filetype.startsWith('audio/')) return '音频';
  if (filetype.startsWith('video/')) return '视频';
  if (filetype.startsWith('text/')) return '文本';
  if (filetype.includes('pdf')) return 'PDF';
  if (filetype.includes('word') || filetype.includes('document')) return 'Word';
  if (filetype.includes('excel') || filetype.includes('sheet')) return 'Excel';
  if (filetype.includes('powerpoint') || filetype.includes('presentation')) return 'PPT';
  if (filetype.includes('zip') || filetype.includes('compressed')) return '压缩包';
  
  // 从MIME类型中提取简化格式
  const parts = filetype.split('/');
  if (parts.length === 2) {
    return parts[1].toUpperCase();
  }
  
  return '文件';
};

// 获取文件颜色
const getFileColor = (filetype) => {
  if (!filetype) return 'text-gray-500';
  
  if (filetype.startsWith('image/')) {
    return 'text-green-500';
  } else if (filetype.includes('pdf')) {
    return 'text-red-500';
  } else if (filetype.includes('word') || filetype.includes('document')) {
    return 'text-blue-500';
  } else if (filetype.includes('excel') || filetype.includes('spreadsheet')) {
    return 'text-green-600';
  } else if (filetype.includes('zip') || filetype.includes('compressed')) {
    return 'text-yellow-500';
  } else if (filetype.includes('text')) {
    return 'text-gray-600';
  }
  
  return 'text-gray-500';
};

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// 处理文件点击
const handleFileClick = (attachment) => {
  if (attachment.filetype && attachment.filetype.startsWith('image/')) {
    // 查看图片
    currentImage.value = attachment;
    showImageViewer.value = true;
  } else {
    // 下载其他类型文件
    downloadFile(attachment);
  }
};

// 下载文件
const downloadFile = async (attachment) => {
  try {
    await attachmentsStore.downloadAttachment(attachment.id, attachment.filename);
  } catch (error) {
    console.error('下载文件失败:', error);
  }
};

// 关闭图片查看器
const closeImageViewer = () => {
  showImageViewer.value = false;
  currentImage.value = null;
};

// 是否为图片
const isImage = (filetype) => {
  return filetype && filetype.startsWith('image/');
};

// 按日期过滤的附件
const attachmentsByDate = computed(() => {
  const result = {};
  
  if (!dateGroups.value || !libraryAttachments.value) return result;
  
  // 创建日期索引
  dateGroups.value.forEach(group => {
    // 确保日期格式一致，处理可能的完整时间格式
    let date = group.date;
    if (date && date.includes('T')) {
      date = date.split('T')[0];
    }
    
    result[date] = {
      displayDate: group.displayDate,
      attachments: []
    };
  });
  
  // 将附件分配到日期组
  libraryAttachments.value.forEach(attachment => {
    // 从完整时间中提取日期部分
    const createdAt = attachment.created_at;
    let date = createdAt;
    
    if (createdAt && createdAt.includes('T')) {
      date = createdAt.split('T')[0];
    }
    
    if (result[date]) {
      result[date].attachments.push(attachment);
    } else {
      // 如果找不到匹配的日期组，尝试查找包含日期的组
      const matchingDate = Object.keys(result).find(groupDate => 
        date && groupDate && date.includes(groupDate)
      );
      
      if (matchingDate) {
        result[matchingDate].attachments.push(attachment);
      }
    }
  });
  
  console.log('处理后的日期组:', result);
  return result;
});

// 搜索过滤后的日期组
const filteredDateGroups = computed(() => {
  if (!searchKeyword.value) return dateGroups.value;
  
  const keyword = searchKeyword.value.toLowerCase();
  const filteredDates = [];
  
  for (const date in attachmentsByDate.value) {
    const group = attachmentsByDate.value[date];
    const matchingAttachments = group.attachments.filter(att => 
      att.filename.toLowerCase().includes(keyword)
    );
    
    if (matchingAttachments.length > 0) {
      filteredDates.push({
        date,
        displayDate: group.displayDate,
        count: matchingAttachments.length
      });
    }
  }
  
  return filteredDates;
});

// 搜索过滤后的附件
const filteredAttachments = computed(() => {
  if (!searchKeyword.value) return libraryAttachments.value;
  
  const keyword = searchKeyword.value.toLowerCase();
  return libraryAttachments.value.filter(att => 
    att.filename.toLowerCase().includes(keyword)
  );
});

// 当搜索词改变时
watch(searchKeyword, () => {
  // 可以在这里添加搜索逻辑
});

// 将ISO日期格式化为年月日
const formatISODate = (isoDate) => {
  if (!isoDate) return '';
  
  console.log('格式化ISO日期:', isoDate);
  
  try {
    // 处理不同的日期格式
    let dateStr = isoDate;
    if (dateStr.includes('T')) {
      dateStr = dateStr.split('T')[0];
    }
    
    // 创建日期对象
    const date = new Date(dateStr);
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      console.warn('无效的日期格式:', isoDate);
      
      // 尝试手动解析
      const [year, month, day] = dateStr.split('-').map(Number);
      if (!isNaN(year) && !isNaN(month) && !isNaN(day)) {
        return `${year}年${month}月${day}日`;
      }
      
      return isoDate;
    }
    
    // 使用toLocaleDateString进行格式化
    try {
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
    } catch (e) {
      // 降级格式化方案
      const year = date.getFullYear();
      const month = date.getMonth() + 1;
      const day = date.getDate();
      return `${year}年${month}月${day}日`;
    }
  } catch (error) {
    console.error('日期格式化错误:', error, isoDate);
    return isoDate;
  }
};

// 页面加载时获取数据
onMounted(() => {
  console.log('资源库页面加载，开始获取数据');
  fetchLibrary();
});

// 调试按钮 - 直接显示文件，用于排查问题
const forceShowFiles = ref(false);
const toggleForceShow = () => {
  forceShowFiles.value = !forceShowFiles.value;
  console.log('强制显示文件模式:', forceShowFiles.value);
};

// 根据文件类型获取不同的渐变背景
const getGradientByFileType = (filetype) => {
  if (!filetype) return 'from-gray-300 to-gray-400 dark:from-gray-700 dark:to-gray-600';
  
  if (filetype.startsWith('image/')) {
    return 'from-green-400 to-green-500 dark:from-green-500 dark:to-green-600';
  } else if (filetype.includes('pdf')) {
    return 'from-red-400 to-red-500 dark:from-red-500 dark:to-red-600';
  } else if (filetype.includes('word') || filetype.includes('document')) {
    return 'from-blue-400 to-blue-500 dark:from-blue-500 dark:to-blue-600';
  } else if (filetype.includes('excel') || filetype.includes('spreadsheet')) {
    return 'from-green-500 to-green-600 dark:from-green-600 dark:to-green-700';
  } else if (filetype.includes('zip') || filetype.includes('compressed')) {
    return 'from-yellow-400 to-yellow-500 dark:from-yellow-500 dark:to-yellow-600';
  } else if (filetype.includes('text')) {
    return 'from-gray-400 to-gray-500 dark:from-gray-500 dark:to-gray-600';
  }
  
  return 'from-purple-400 to-purple-500 dark:from-purple-500 dark:to-purple-600';
};

// 根据文件类型获取图标背景色
const getIconBgByFileType = (filetype) => {
  if (!filetype) return 'bg-gray-500 dark:bg-gray-600';
  
  if (filetype.startsWith('image/')) {
    return 'bg-green-500 dark:bg-green-600';
  } else if (filetype.includes('pdf')) {
    return 'bg-red-500 dark:bg-red-600';
  } else if (filetype.includes('word') || filetype.includes('document')) {
    return 'bg-blue-500 dark:bg-blue-600';
  } else if (filetype.includes('excel') || filetype.includes('spreadsheet')) {
    return 'bg-green-600 dark:bg-green-700';
  } else if (filetype.includes('zip') || filetype.includes('compressed')) {
    return 'bg-yellow-500 dark:bg-yellow-600';
  } else if (filetype.includes('text')) {
    return 'bg-gray-500 dark:bg-gray-600';
  }
  
  return 'bg-purple-500 dark:bg-purple-600';
};

// 根据文件类型获取图标路径
const getIconPathByFileType = (filetype) => {
  if (!filetype) return 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z';
  
  if (filetype.startsWith('image/')) {
    return 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z';
  } else if (filetype.includes('pdf')) {
    return 'M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z';
  } else if (filetype.includes('word') || filetype.includes('document')) {
    return 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z';
  } else if (filetype.includes('excel') || filetype.includes('spreadsheet')) {
    return 'M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z';
  } else if (filetype.includes('zip') || filetype.includes('compressed')) {
    return 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4';
  } else if (filetype.includes('text')) {
    return 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z';
  }
  
  return 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z';
};
</script>

<template>
  <div class="px-4 py-6 bg-gray-50 dark:bg-gray-900 min-h-screen">
    <!-- 顶部导航和搜索区域 -->
    <div class="max-w-7xl mx-auto mb-6">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-md overflow-hidden transition-all duration-300 hover:shadow-lg border border-gray-100 dark:border-gray-700">
        <div class="p-5">
          <div class="flex flex-col md:flex-row md:items-center gap-4">
            <!-- 左侧标题 -->
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-gradient-to-br from-blue-500 to-indigo-600 dark:from-blue-600 dark:to-indigo-700 rounded-lg p-3 mr-4 shadow-md">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
                </svg>
              </div>
              <div>
                <h1 class="text-2xl font-bold text-gray-800 dark:text-white bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-500 bg-clip-text text-transparent">资源库</h1>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">共 {{ totalAttachments }} 个文件</p>
              </div>
            </div>
            
            <!-- 右侧搜索和过滤 -->
            <div class="flex-grow mt-4 md:mt-0">
              <div class="flex flex-col sm:flex-row gap-3">
                <!-- 搜索框 -->
                <div class="relative flex-grow">
                  <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                  </div>
                  <input 
                    v-model="searchKeyword" 
                    type="text" 
                    placeholder="搜索文件名..." 
                    class="block w-full pl-10 pr-3 py-2.5 border border-gray-200 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200 transition-all duration-200"
                  />
                </div>
                
                <div class="flex gap-2">
                  <!-- 过滤选择器 -->
                  <div class="relative">
                    <select 
                      v-model="selectedFileType" 
                      class="block appearance-none w-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-200 py-2.5 px-4 pr-8 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200"
                    >
                      <option v-for="option in fileTypeOptions" :key="option.value" :value="option.value">
                        {{ option.label }}
                      </option>
                    </select>
                    <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-400">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </div>
                  </div>
                  
                  <!-- 视图切换 -->
                  <div class="bg-gray-100 dark:bg-gray-700 rounded-lg flex overflow-hidden shadow-sm">
                    <button
                      @click="viewType = 'grid'"
                      :class="[
                        'px-3 py-2.5 transition-all duration-300', 
                        viewType === 'grid' 
                          ? 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white shadow-md' 
                          : 'text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                      ]"
                      title="网格视图"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                      </svg>
                    </button>
                    <button
                      @click="viewType = 'list'"
                      :class="[
                        'px-3 py-2.5 transition-all duration-300', 
                        viewType === 'list' 
                          ? 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white shadow-md' 
                          : 'text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                      ]"
                      title="列表视图"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="max-w-7xl mx-auto">
      <!-- 资源统计 -->
      <LibraryStats v-if="!isLoading && libraryAttachments.length > 0" :attachments="libraryAttachments" />
      
      <!-- 加载状态 -->
      <div v-if="isLoading" class="flex justify-center items-center py-20">
        <div class="loader-container">
          <div class="ripple-loader">
            <div></div>
            <div></div>
          </div>
          <span class="ml-4 text-gray-600 dark:text-gray-300 font-medium">加载资源中...</span>
        </div>
      </div>
      
      <!-- 没有文件 -->
      <div v-else-if="!libraryAttachments || libraryAttachments.length === 0" class="text-center py-20 bg-white dark:bg-gray-800 rounded-2xl shadow-md transition-all duration-300">
        <div class="text-primary-200 dark:text-primary-900 text-8xl mb-6 opacity-50">
          <i class="fas fa-cloud-upload-alt"></i>
        </div>
        <h3 class="text-xl font-medium text-gray-600 dark:text-gray-300 mb-2">暂无文件</h3>
        <p class="text-gray-500 dark:text-gray-400">在笔记中添加附件后，文件将显示在这里</p>
      </div>
      
      <!-- 文件列表（按日期分组） -->
      <div v-else>
        <!-- 如果有附件但没有日期组或强制显示所有文件 -->
        <div v-if="forceShowFiles || (dateGroups.length === 0 && libraryAttachments.length > 0)" class="mb-10">
          <!-- 日期标题与计数 -->
          <div class="mb-5 flex items-center">
            <div class="flex items-center mr-3">
              <div class="w-1.5 h-10 bg-blue-500 dark:bg-blue-600 rounded-full mr-2"></div>
              <h2 class="text-xl font-semibold text-gray-700 dark:text-gray-200">
                <span class="mr-2">全部文件</span>
                <span class="text-xs font-normal px-2 py-0.5 bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 rounded-full">
                  {{ libraryAttachments.length }} 文件
                </span>
              </h2>
            </div>
          </div>
          
          <!-- 网格视图 -->
          <div v-if="viewType === 'grid'" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5">
            <div 
              v-for="attachment in libraryAttachments" 
              :key="attachment.id"
              v-show="!searchKeyword || attachment.filename.toLowerCase().includes(searchKeyword.toLowerCase())"
              class="cursor-pointer"
              @click="handleFileClick(attachment)"
            >
              <!-- 图片预览 -->
              <div v-if="isImage(attachment.filetype)" class="group relative h-0 pb-[100%] rounded-xl overflow-hidden bg-gray-100 dark:bg-gray-900 shadow-md hover:shadow-lg transition-all duration-500">
                <!-- 图片 -->
                <img 
                  :src="attachmentsStore.getAttachmentUrl(attachment.id)" 
                  :alt="attachment.filename"
                  class="absolute inset-0 w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
                  loading="lazy"
                />
                
                <!-- 悬停叠加层 -->
                <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/30 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex flex-col justify-end">
                  <!-- 文件信息 -->
                  <div class="p-3.5">
                    <div class="truncate text-white text-sm font-medium mb-1">{{ attachment.filename }}</div>
                    <div class="flex justify-between items-center text-xs text-white/80">
                      <span>{{ formatFileSize(attachment.filesize) }}</span>
                      <span class="px-1.5 py-0.5 bg-white/20 backdrop-blur-sm rounded-full text-xs">
                        {{ getFileTypeLabel(attachment.filetype) }}
                      </span>
                    </div>
                  </div>
                </div>
                
                <!-- 操作按钮 -->
                <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-all duration-300 transform group-hover:translate-y-0 -translate-y-2">
                  <button 
                    @click.stop="downloadFile(attachment)"
                    class="p-2 bg-white/90 hover:bg-white text-gray-900 rounded-full shadow-lg hover:shadow-xl transition-all duration-200"
                    title="下载文件"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                  </button>
                </div>
              </div>
              
              <!-- 文件卡片 -->
              <div v-else class="group bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-md hover:shadow-lg border border-gray-100 dark:border-gray-700 transition-all duration-300 h-full transform hover:-translate-y-1">
                <!-- 顶部色带 -->
                <div class="h-2 bg-gradient-to-r" :class="getGradientByFileType(attachment.filetype)"></div>
                
                <div class="p-4">
                  <!-- 图标和文件名 -->
                  <div class="flex items-start mb-3">
                    <!-- 文件类型图标 -->
                    <div class="flex-shrink-0 p-2.5 mr-3 rounded-full shadow-sm" :class="getIconBgByFileType(attachment.filetype)">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white drop-shadow-sm" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getIconPathByFileType(attachment.filetype)" />
                      </svg>
                    </div>
                    
                    <div class="flex-1 min-w-0">
                      <div class="truncate font-medium text-gray-800 dark:text-gray-200 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">{{ attachment.filename }}</div>
                      <div class="flex items-center mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                        <span class="mr-2">{{ formatFileSize(attachment.filesize) }}</span>
                        <span class="px-1.5 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
                          {{ getFileTypeLabel(attachment.filetype) }}
                        </span>
                      </div>
                    </div>
                  </div>
                  
                  <!-- 下载按钮 -->
                  <div class="flex justify-end mt-3">
                    <button 
                      @click.stop="downloadFile(attachment)"
                      class="p-1.5 text-gray-400 hover:text-blue-500 dark:hover:text-blue-400 bg-gray-100 dark:bg-gray-700 hover:bg-blue-50 dark:hover:bg-blue-900/30 rounded-full transition-all group-hover:shadow-sm"
                      title="下载文件"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 列表视图 -->
          <div v-else>
            <div class="bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-md border border-gray-100 dark:border-gray-700 divide-y divide-gray-100 dark:divide-gray-700">
              <div 
                v-for="attachment in libraryAttachments" 
                :key="attachment.id"
                v-show="!searchKeyword || attachment.filename.toLowerCase().includes(searchKeyword.toLowerCase())"
                class="group border-b border-gray-100 dark:border-gray-700 last:border-0 hover:bg-gray-50 dark:hover:bg-gray-750 transition-all duration-300 cursor-pointer"
                @click="handleFileClick(attachment)"
              >
                <div class="p-3.5 flex items-center">
                  <!-- 左侧图标/缩略图 -->
                  <div class="mr-4 flex-shrink-0">
                    <!-- 图片缩略图 -->
                    <div v-if="isImage(attachment.filetype)" class="w-14 h-14 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 shadow-sm group-hover:shadow transition-all duration-300">
                      <img 
                        :src="attachmentsStore.getAttachmentUrl(attachment.id)" 
                        :alt="attachment.filename"
                        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                        loading="lazy"
                      />
                    </div>
                    
                    <!-- 文件图标 -->
                    <div v-else class="w-12 h-12 rounded-lg flex items-center justify-center shadow-sm group-hover:shadow transition-all duration-300" :class="getIconBgByFileType(attachment.filetype)">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white drop-shadow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getIconPathByFileType(attachment.filetype)" />
                      </svg>
                    </div>
                  </div>
                  
                  <!-- 文件信息 -->
                  <div class="flex-grow min-w-0">
                    <div class="flex items-center mb-1">
                      <div class="truncate font-medium text-gray-800 dark:text-gray-200 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">{{ attachment.filename }}</div>
                    </div>
                    <div class="flex items-center text-xs text-gray-500 dark:text-gray-400">
                      <span class="mr-3">{{ formatFileSize(attachment.filesize) }}</span>
                      <span class="px-2 py-0.5 rounded-full text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 group-hover:bg-blue-50 dark:group-hover:bg-blue-900/20 group-hover:text-blue-700 dark:group-hover:text-blue-300 transition-colors">
                        {{ getFileTypeLabel(attachment.filetype) }}
                      </span>
                      <span class="ml-3 text-gray-400 dark:text-gray-500">{{ formatISODate(attachment.created_at) }}</span>
                    </div>
                  </div>
                  
                  <!-- 下载按钮 -->
                  <button 
                    @click.stop="downloadFile(attachment)"
                    class="ml-3 p-2 text-gray-400 hover:text-blue-500 dark:hover:text-blue-400 bg-gray-100 dark:bg-gray-700 hover:bg-blue-50 dark:hover:bg-blue-900/30 rounded-full transition-all opacity-0 group-hover:opacity-100 shadow-sm"
                    title="下载文件"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 日期分组显示（仅当有日期组且不强制显示所有文件时） -->
        <div v-else class="mb-10 space-y-8">
          <div v-for="group in filteredDateGroups" :key="group.date">
            <!-- 日期标题与计数 -->
            <div class="mb-5 flex items-center">
              <div class="flex items-center mr-3">
                <div class="w-1.5 h-10 bg-blue-500 dark:bg-blue-600 rounded-full mr-2"></div>
                <h2 class="text-xl font-semibold text-gray-700 dark:text-gray-200">
                  {{ group.displayDate || formatISODate(group.date) }}
                </h2>
              </div>
              <span class="px-2.5 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 text-xs font-medium rounded-full">
                {{ group.count }} 个文件
              </span>
            </div>
            
            <!-- 网格视图 -->
            <div v-if="viewType === 'grid'" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5">
              <div 
                v-for="attachment in (attachmentsByDate[group.date?.split('T')[0] || group.date]?.attachments || [])" 
                :key="attachment.id"
                v-show="!searchKeyword || attachment.filename.toLowerCase().includes(searchKeyword.toLowerCase())"
                class="cursor-pointer"
                @click="handleFileClick(attachment)"
              >
                <!-- 图片预览 -->
                <div v-if="isImage(attachment.filetype)" class="group relative h-0 pb-[100%] rounded-xl overflow-hidden bg-gray-100 dark:bg-gray-900 shadow-md hover:shadow-lg transition-all duration-500">
                  <!-- 图片 -->
                  <img 
                    :src="attachmentsStore.getAttachmentUrl(attachment.id)" 
                    :alt="attachment.filename"
                    class="absolute inset-0 w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
                    loading="lazy"
                  />
                  
                  <!-- 悬停叠加层 -->
                  <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/30 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex flex-col justify-end">
                    <!-- 文件信息 -->
                    <div class="p-3.5">
                      <div class="truncate text-white text-sm font-medium mb-1">{{ attachment.filename }}</div>
                      <div class="flex justify-between items-center text-xs text-white/80">
                        <span>{{ formatFileSize(attachment.filesize) }}</span>
                        <span class="px-1.5 py-0.5 bg-white/20 backdrop-blur-sm rounded-full text-xs">
                          {{ getFileTypeLabel(attachment.filetype) }}
                        </span>
                      </div>
                    </div>
                  </div>
                  
                  <!-- 操作按钮 -->
                  <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-all duration-300 transform group-hover:translate-y-0 -translate-y-2">
                    <button 
                      @click.stop="downloadFile(attachment)"
                      class="p-2 bg-white/90 hover:bg-white text-gray-900 rounded-full shadow-lg hover:shadow-xl transition-all duration-200"
                      title="下载文件"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                <!-- 文件卡片 -->
                <div v-else class="group bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-md hover:shadow-lg border border-gray-100 dark:border-gray-700 transition-all duration-300 h-full transform hover:-translate-y-1">
                  <!-- 顶部色带 -->
                  <div class="h-2 bg-gradient-to-r" :class="getGradientByFileType(attachment.filetype)"></div>
                  
                  <div class="p-4">
                    <!-- 图标和文件名 -->
                    <div class="flex items-start mb-3">
                      <!-- 文件类型图标 -->
                      <div class="flex-shrink-0 p-2.5 mr-3 rounded-full shadow-sm" :class="getIconBgByFileType(attachment.filetype)">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white drop-shadow-sm" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getIconPathByFileType(attachment.filetype)" />
                        </svg>
                      </div>
                      
                      <div class="flex-1 min-w-0">
                        <div class="truncate font-medium text-gray-800 dark:text-gray-200 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">{{ attachment.filename }}</div>
                        <div class="flex items-center mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                          <span class="mr-2">{{ formatFileSize(attachment.filesize) }}</span>
                          <span class="px-1.5 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
                            {{ getFileTypeLabel(attachment.filetype) }}
                          </span>
                        </div>
                      </div>
                    </div>
                    
                    <!-- 下载按钮 -->
                    <div class="flex justify-end mt-3">
                      <button 
                        @click.stop="downloadFile(attachment)"
                        class="p-1.5 text-gray-400 hover:text-blue-500 dark:hover:text-blue-400 bg-gray-100 dark:bg-gray-700 hover:bg-blue-50 dark:hover:bg-blue-900/30 rounded-full transition-all group-hover:shadow-sm"
                        title="下载文件"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 列表视图 -->
            <div v-else class="bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-md border border-gray-100 dark:border-gray-700 divide-y divide-gray-100 dark:divide-gray-700">
              <div 
                v-for="attachment in (attachmentsByDate[group.date?.split('T')[0] || group.date]?.attachments || [])" 
                :key="attachment.id"
                v-show="!searchKeyword || attachment.filename.toLowerCase().includes(searchKeyword.toLowerCase())"
                class="group border-b border-gray-100 dark:border-gray-700 last:border-0 hover:bg-gray-50 dark:hover:bg-gray-750 transition-all duration-300 cursor-pointer"
                @click="handleFileClick(attachment)"
              >
                <div class="p-3.5 flex items-center">
                  <!-- 左侧图标/缩略图 -->
                  <div class="mr-4 flex-shrink-0">
                    <!-- 图片缩略图 -->
                    <div v-if="isImage(attachment.filetype)" class="w-14 h-14 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 shadow-sm group-hover:shadow transition-all duration-300">
                      <img 
                        :src="attachmentsStore.getAttachmentUrl(attachment.id)" 
                        :alt="attachment.filename"
                        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                        loading="lazy"
                      />
                    </div>
                    
                    <!-- 文件图标 -->
                    <div v-else class="w-12 h-12 rounded-lg flex items-center justify-center shadow-sm group-hover:shadow transition-all duration-300" :class="getIconBgByFileType(attachment.filetype)">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white drop-shadow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getIconPathByFileType(attachment.filetype)" />
                      </svg>
                    </div>
                  </div>
                  
                  <!-- 文件信息 -->
                  <div class="flex-grow min-w-0">
                    <div class="flex items-center mb-1">
                      <div class="truncate font-medium text-gray-800 dark:text-gray-200 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">{{ attachment.filename }}</div>
                    </div>
                    <div class="flex items-center text-xs text-gray-500 dark:text-gray-400">
                      <span class="mr-3">{{ formatFileSize(attachment.filesize) }}</span>
                      <span class="px-2 py-0.5 rounded-full text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 group-hover:bg-blue-50 dark:group-hover:bg-blue-900/20 group-hover:text-blue-700 dark:group-hover:text-blue-300 transition-colors">
                        {{ getFileTypeLabel(attachment.filetype) }}
                      </span>
                      <span class="ml-3 text-gray-400 dark:text-gray-500">{{ formatISODate(attachment.created_at) }}</span>
                    </div>
                  </div>
                  
                  <!-- 下载按钮 -->
                  <button 
                    @click.stop="downloadFile(attachment)"
                    class="ml-3 p-2 text-gray-400 hover:text-blue-500 dark:hover:text-blue-400 bg-gray-100 dark:bg-gray-700 hover:bg-blue-50 dark:hover:bg-blue-900/30 rounded-full transition-all opacity-0 group-hover:opacity-100 shadow-sm"
                    title="下载文件"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 分页控制 -->
        <div v-if="libraryAttachments.length > pageSize" class="mt-8 flex justify-center">
          <div class="inline-flex rounded-lg overflow-hidden bg-white dark:bg-gray-800 shadow-md divide-x divide-gray-200 dark:divide-gray-700 border border-gray-100 dark:border-gray-700">
            <!-- 上一页 -->
            <button 
              @click="currentPage > 1 ? (currentPage--, fetchLibrary()) : null" 
              :disabled="currentPage <= 1"
              class="px-4 py-2.5 flex items-center transition-all duration-300"
              :class="currentPage <= 1 ? 'text-gray-400 cursor-not-allowed' : 'hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400'"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              <span>上一页</span>
            </button>
            
            <!-- 页码显示 -->
            <div class="px-4 py-2.5 bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-medium flex items-center shadow-inner">
              {{ currentPage }} / {{ Math.ceil(libraryAttachments.length / pageSize) }}
            </div>
            
            <!-- 下一页 -->
            <button 
              @click="currentPage < Math.ceil(libraryAttachments.length / pageSize) ? (currentPage++, fetchLibrary()) : null" 
              :disabled="currentPage >= Math.ceil(libraryAttachments.length / pageSize)"
              class="px-4 py-2.5 flex items-center transition-all duration-300"
              :class="currentPage >= Math.ceil(libraryAttachments.length / pageSize) ? 'text-gray-400 cursor-not-allowed' : 'hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400'"
            >
              <span>下一页</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </div>
        
        <!-- 只在开发模式下显示调试按钮 -->
        <div v-if="isDev" class="fixed bottom-4 right-4 z-10">
          <button 
            @click="toggleForceShow" 
            class="bg-yellow-500 hover:bg-yellow-600 text-white text-xs px-3 py-1.5 rounded-full shadow-lg transition-colors"
          >
            {{ forceShowFiles ? '使用日期分组' : '显示全部文件' }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- 图片查看器 -->
    <div v-if="showImageViewer" class="fixed inset-0 z-50" @click.self="closeImageViewer">
      <!-- 玻璃态背景 -->
      <div class="absolute inset-0 bg-black/90 backdrop-blur-md"></div>
      
      <!-- 主内容区 -->
      <div class="relative z-10 w-full h-full flex flex-col">
        
        <!-- 顶部工具栏 -->
        <div class="flex items-center justify-between p-4 bg-gradient-to-b from-black/80 to-transparent">
          <h3 class="text-white text-lg font-medium truncate max-w-[60%] opacity-90" v-if="currentImage">
            {{ currentImage.filename }}
          </h3>
          
          <div class="flex space-x-3">
            <button 
              v-if="currentImage"
              @click.stop="downloadFile(currentImage)"
              class="bg-blue-600/90 hover:bg-blue-700 text-white p-2.5 rounded-full transition-all shadow-lg hover:shadow-xl hover:scale-105"
              title="下载图片"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
            </button>
            <button 
              @click.stop="closeImageViewer"
              class="bg-red-600/90 hover:bg-red-700 text-white p-2.5 rounded-full transition-all shadow-lg hover:shadow-xl hover:scale-105"
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
          <!-- 图片显示 -->
          <div class="relative group">
            <img
              v-if="currentImage"
              :src="attachmentsStore.getAttachmentUrl(currentImage.id)" 
              :alt="currentImage.filename"
              class="max-h-[85vh] max-w-[90vw] object-contain rounded-lg shadow-2xl transition-all duration-500"
            />
            <div class="absolute inset-0 bg-gradient-to-b from-transparent to-black/20 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          </div>
        </div>
      </div>
      
      <!-- 移动端专用底部工具栏 -->
      <div class="fixed bottom-8 left-0 right-0 flex justify-center z-30 sm:hidden">
        <div class="flex items-center px-5 py-2.5 rounded-full bg-gray-900/70 backdrop-blur-sm shadow-xl">
          <div class="text-white/90 text-sm truncate max-w-[60vw]" v-if="currentImage">
            {{ currentImage.filename }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 确保顺滑的过渡效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

.animate-fadeIn {
  animation: fadeIn 0.3s ease-out forwards;
}

/* 使图片保持正方形比例 */
.aspect-square {
  aspect-ratio: 1 / 1;
}

/* 波纹加载动画 */
.loader-container {
  display: flex;
  align-items: center;
  justify-content: center;
}

.ripple-loader {
  display: inline-block;
  position: relative;
  width: 80px;
  height: 80px;
}

.ripple-loader div {
  position: absolute;
  border: 4px solid;
  @apply border-blue-500 dark:border-blue-400;
  opacity: 1;
  border-radius: 50%;
  animation: ripple-loader 1.5s cubic-bezier(0, 0.2, 0.8, 1) infinite;
}

.ripple-loader div:nth-child(2) {
  animation-delay: -0.5s;
}

@keyframes ripple-loader {
  0% {
    top: 36px;
    left: 36px;
    width: 0;
    height: 0;
    opacity: 0;
  }
  4.9% {
    top: 36px;
    left: 36px;
    width: 0;
    height: 0;
    opacity: 0;
  }
  5% {
    top: 36px;
    left: 36px;
    width: 0;
    height: 0;
    opacity: 1;
  }
  100% {
    top: 0px;
    left: 0px;
    width: 72px;
    height: 72px;
    opacity: 0;
  }
}

/* 自定义滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.03);
  border-radius: 8px;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}

/* 暗黑模式滚动条 */
.dark ::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.03);
}

.dark ::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style> 