<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- 页面标题 -->
    <div class="text-center mb-8">
      <h1 class="text-3xl sm:text-4xl font-bold text-gray-800 dark:text-white mb-3">
        热门<span class="text-primary-600 dark:text-primary-400">推荐</span>
      </h1>
      <p class="text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
        发现精彩知识，探索创意灵感，分享你的见解
      </p>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 dark:border-primary-300"></div>
    </div>
    
    <!-- 无笔记状态 -->
    <div v-else-if="publicNotes.length === 0" class="bg-white dark:bg-gray-800 rounded-xl p-8 text-center shadow-md border border-gray-100 dark:border-gray-700">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto mb-4 text-gray-400 dark:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
      </svg>
      <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">暂无公开笔记</h3>
      <p class="text-gray-500 dark:text-gray-400 mb-4 max-w-md mx-auto">
        成为第一个分享知识的用户吧！创建笔记并标记为公开，让更多人受益。
      </p>
      <router-link 
        :to="isAuthenticated ? '/notes/new' : '/auth/login'" 
        class="px-6 py-2.5 bg-primary-600 text-white rounded-lg font-medium hover:bg-primary-700 transition-colors shadow-sm"
      >
        {{ isAuthenticated ? '创建笔记' : '立即加入' }}
      </router-link>
    </div>

    <!-- 笔记内容列表 -->
    <div v-else>
      <!-- 笔记卡片 -->
      <div v-for="(note, index) in publicNotes" :key="note.id" 
        class="mb-6 bg-white dark:bg-gray-800 rounded-xl shadow-sm overflow-hidden border border-gray-200 dark:border-gray-700 transition-all duration-300 hover:shadow-md">
        
        <!-- 笔记顶部信息栏：标题和作者（微博样式） -->
        <div class="flex items-start p-4">
          <!-- 用户头像（微博样式圆形头像） -->
          <div class="w-10 h-10 rounded-full bg-primary-500 dark:bg-primary-600 flex items-center justify-center text-white font-medium text-sm mr-3 flex-shrink-0 border-2 border-white dark:border-gray-700 shadow-sm">
            {{ getUserInitial(note) }}
          </div>
          
          <!-- 用户信息区域 -->
          <div class="flex-1">
            <div class="flex flex-col">
              <!-- 用户名与认证标识 -->
              <div class="flex items-center">
                <span class="font-medium text-gray-800 dark:text-gray-200">{{ getUserDisplayName(note) }}</span>
                <span class="ml-1 text-primary-500 dark:text-primary-400">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                </span>
              </div>
              
              <!-- 发布时间 -->
              <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ getFormattedDate(note.created_at) }}</div>
            </div>
            
            <!-- 笔记标题 (加粗突出显示) -->
            <h2 class="text-lg font-bold text-gray-800 dark:text-white mt-2 mb-2">{{ note.title }}</h2>
            
            <!-- 笔记内容预览 -->
            <div class="prose dark:prose-invert max-w-none prose-sm line-clamp-3 mt-2">
              <MarkdownRenderer 
                :content="getPreviewContent(note.content)" 
                :simple-mode="true" 
              />
            </div>
          </div>
        </div>
        
        <!-- 标签区域（微博话题样式） -->
        <div v-if="note.tags && note.tags.length > 0" class="flex flex-wrap gap-1.5 px-4 pb-3">
          <span 
            v-for="tag in note.tags" 
            :key="tag.id" 
            class="px-3 py-1 bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-300 text-xs rounded-full inline-flex items-center hover:bg-blue-100 dark:hover:bg-blue-800/50 cursor-pointer transition-colors"
          >
            <span class="mr-0.5">#</span>
            {{ getTagName(tag) }}
          </span>
        </div>
        
        <!-- 附件区域（微博图片/附件展示样式） -->
        <div v-if="note.attachments && note.attachments.length > 0" class="px-4 pb-3 mt-1">
          <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-2">
            <div class="flex flex-wrap gap-2">
              <a 
                v-for="attachment in note.attachments" 
                :key="attachment.id"
                :href="attachment.file_url"
                target="_blank"
                class="inline-flex items-center px-3 py-1.5 bg-white dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg text-xs text-gray-700 dark:text-gray-300 transition-colors border border-gray-200 dark:border-gray-600"
              >
                <!-- 文件图标 - 根据文件类型显示不同图标 -->
                <svg v-if="attachment.filetype.includes('image')" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1.5 text-blue-500 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <svg v-else-if="attachment.filetype.includes('pdf')" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1.5 text-red-500 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <svg v-else-if="attachment.filetype.includes('zip') || attachment.filetype.includes('rar')" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1.5 text-orange-500 dark:text-orange-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1.5 text-gray-500 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                </svg>
                
                <!-- 文件名称 -->
                <span class="truncate max-w-[150px]" :title="attachment.filename">
                  {{ attachment.filename }}
                </span>
                
                <!-- 文件大小 -->
                <span class="ml-1 text-gray-500 dark:text-gray-400">
                  ({{ formatFileSize(attachment.filesize) }})
                </span>
              </a>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 分页与加载更多控制 -->
      <div v-if="publicNotes.length > 0 && !isLoading" class="mt-8 mb-10">
        <!-- 加载更多按钮 -->
        <div class="text-center" v-if="hasMore">
          <button 
            @click="loadMore" 
            class="px-6 py-2.5 bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-full shadow-sm text-gray-700 dark:text-gray-200 font-medium hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors flex items-center mx-auto"
            :disabled="isLoadingMore"
          >
            <span v-if="!isLoadingMore" class="flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
              查看更多
            </span>
            <div v-else class="flex items-center">
              <div class="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-primary-500 mr-2"></div>
              加载中...
            </div>
          </button>
        </div>
        
        <!-- 全部加载完毕提示 -->
        <div v-else class="text-center text-sm text-gray-500 dark:text-gray-400 border-t border-gray-200 dark:border-gray-700 pt-6">
          <div class="flex items-center justify-center">
            <span class="h-px bg-gray-300 dark:bg-gray-600 w-16"></span>
            <span class="mx-4">已显示全部内容</span>
            <span class="h-px bg-gray-300 dark:bg-gray-600 w-16"></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useUserStore } from '@/stores/user';
import { useNotesStore } from '@/stores/notes';
import { storeToRefs } from 'pinia';
import MarkdownRenderer from '@/components/MarkdownRenderer.vue';

// 获取用户状态
const userStore = useUserStore();
const notesStore = useNotesStore();
const isAuthenticated = userStore.isAuthenticated;

// 从store获取公开笔记状态
const { totalPublicNotes, isLoading: storeLoading } = storeToRefs(notesStore);

// 本地状态
const publicNotes = ref([]);
const isLoading = ref(true);
const isLoadingMore = ref(false);
const currentPage = ref(1);
const pageSize = ref(8); // 每页显示笔记数量
const hasMore = ref(true);
const expandedNotes = reactive({}); // 跟踪展开的笔记

// 获取用户信息函数
const getUserInitial = (note) => {
  // 检查note.user是否存在以及如何访问
  console.log('笔记的用户对象:', note.user, typeof note.user);
  
  // 如果用户对象不存在
  if (!note.user) return '?';
  
  // 后端预加载了username字段
  const username = note.user.username;
  return username ? username.charAt(0).toUpperCase() : '?';
};

const getUserDisplayName = (note) => {
  // 确保获取用户信息
  if (!note.user) return '匿名用户';
  
  // 根据用户类型显示不同样式（用于显示认证标志等）
  console.log('用户信息:', note.user);
  return note.user.username || '匿名用户';
};

// 获取标签名称
const getTagName = (tag) => {
  if (typeof tag === 'string') return tag;
  return tag.name || '';
};

// 截取内容预览
const getPreviewContent = (content) => {
  if (!content) return '';
  
  // 限制预览内容长度，保留简短摘要
  const maxLength = 300;
  if (content.length <= maxLength) return content;
  return content.substring(0, maxLength) + '...';
};

// 正确的日期格式化函数
const getFormattedDate = (dateString, shortFormat = false) => {
  if (!dateString) return '未知时间';
  
  try {
    const date = new Date(dateString);
    
    // 确保日期有效
    if (isNaN(date.getTime())) {
      console.error('无效的日期:', dateString);
      return '日期无效';
    }
    
    // 获取当前时间
    const now = new Date();
    
    // 计算时间差（单位：毫秒）
    const diffMs = now.getTime() - date.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHour = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHour / 24);
    
    // 微博风格的时间显示
    if (diffSec < 60) {
      return '刚刚';
    } else if (diffMin < 60) {
      return `${diffMin}分钟前`;
    } else if (diffHour < 24) {
      return `${diffHour}小时前`;
    } else if (diffDay < 7) {
      return `${diffDay}天前`;
    } else {
      // 超过一周，显示日期
      const year = date.getFullYear();
      const month = (date.getMonth() + 1).toString().padStart(2, '0');
      const day = date.getDate().toString().padStart(2, '0');
      
      if (shortFormat) {
        return `${month}-${day}`;
      } else {
        // 如果是不同年份的日期，显示年份
        if (date.getFullYear() !== now.getFullYear()) {
          return `${year}-${month}-${day}`;
        } else {
          return `${month}-${day}`;
        }
      }
    }
  } catch (e) {
    console.error('日期格式化错误:', e);
    return '日期无效';
  }
};

// 切换笔记展开/折叠状态
const expandNote = (noteId) => {
  expandedNotes[noteId] = !expandedNotes[noteId];
};

// 加载公开笔记
async function loadPublicNotes() {
  isLoading.value = true;
  try {
    await notesStore.fetchPublicNotes(currentPage.value, pageSize.value);
    publicNotes.value = notesStore.publicNotes;
    
    // 调试输出
    console.log('加载的笔记数据:', publicNotes.value);
    
    hasMore.value = notesStore.totalPublicNotes > (currentPage.value * pageSize.value);
  } catch (error) {
    console.error('加载公开笔记失败', error);
  } finally {
    isLoading.value = false;
  }
}

// 加载更多笔记
async function loadMore() {
  if (isLoadingMore.value || !hasMore.value) return;
  
  isLoadingMore.value = true;
  currentPage.value++;
  
  try {
    await notesStore.fetchPublicNotes(currentPage.value, pageSize.value);
    // 添加新笔记到已有列表
    publicNotes.value = [...publicNotes.value, ...notesStore.publicNotes];
    // 检查是否还有更多
    hasMore.value = notesStore.totalPublicNotes > (currentPage.value * pageSize.value);
  } catch (error) {
    console.error('加载更多笔记失败', error);
    currentPage.value--; // 恢复页码
  } finally {
    isLoadingMore.value = false;
  }
}

// 文件大小格式化
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B';
  
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return parseFloat((bytes / Math.pow(1024, i)).toFixed(1)) + ' ' + sizes[i];
};

// 初始化加载
onMounted(async () => {
  await loadPublicNotes();
});
</script>

<style scoped>
/* 笔记卡片悬浮效果 */
.hover\:shadow-md:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* 暗色模式下卡片悬浮效果 */
.dark .hover\:shadow-md:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2), 0 2px 4px -1px rgba(0, 0, 0, 0.1);
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
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

/* 动画效果 */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}

.animate-fadeIn {
  animation: fadeIn 0.3s ease-in-out;
}

/* 简化的Markdown样式，针对预览模式 */
:deep(.prose-sm) {
  font-size: 0.875rem;
  line-height: 1.5;
}

:deep(.prose-sm p) {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
}

:deep(.prose-sm h1),
:deep(.prose-sm h2),
:deep(.prose-sm h3) {
  margin-top: 1em;
  margin-bottom: 0.5em;
}

:deep(.prose-sm img) {
  max-height: 150px;
  width: auto;
}
</style>
