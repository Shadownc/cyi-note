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
    </div>

    <!-- 笔记标签 -->
    <div class="flex flex-wrap gap-2 mb-5 sm:mb-6">
      <span 
        v-for="tag in note.tags" 
        :key="tag.id" 
        class="px-2 py-1 bg-primary-100 dark:bg-primary-700/50 text-primary-800 dark:text-primary-100 text-xs rounded-full inline-flex items-center"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-primary-600 dark:bg-primary-300 mr-1 inline-block"></span>
        {{ tag.name }}
      </span>
    </div>

    <!-- 笔记内容 - Markdown渲染 -->
    <div class="prose dark:prose-invert sm:prose-sm md:prose-base max-w-none overflow-hidden dark:bg-gray-800/30 dark:p-4 dark:rounded-lg animate-fadeIn" v-html="renderedContent" ref="contentRef"></div>
    
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
import { ref, defineProps, defineEmits, onMounted, watch, onUnmounted, nextTick } from 'vue';
import { marked } from 'marked';

const props = defineProps({
  note: {
    type: Object,
    required: true
  }
});

defineEmits(['edit', 'delete']);

const renderedContent = ref('');
const showScrollToTop = ref(false);
const showCopySuccess = ref(false);
const contentRef = ref(null);

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

// 自定义Markdown渲染器，添加复制按钮到代码块
const setupMarked = () => {
  const renderer = new marked.Renderer();
  
  // 保存原始代码块渲染方法
  const originalCodeRenderer = renderer.code;
  
  // 自定义代码块渲染
  renderer.code = function(code, language, isEscaped) {
    // 调用原始渲染方法获取HTML
    const html = originalCodeRenderer.call(this, code, language, isEscaped);
    
    // 添加复制按钮包装器
    return `
      <div class="code-block-wrapper relative group">
        <button class="code-copy-btn absolute top-2 right-2 bg-gray-700 dark:bg-gray-600 text-white rounded p-1 opacity-0 group-hover:opacity-100 transition-opacity" data-code="${encodeURIComponent(code)}">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
          </svg>
        </button>
        ${html}
      </div>
    `;
  };
  
  marked.setOptions({
    renderer: renderer,
    highlight: function(code, lang) {
      return code;
    },
    gfm: true,
    breaks: true
  });
};

// 为代码块添加复制功能
const addCopyFunctionality = () => {
  const copyButtons = document.querySelectorAll('.code-copy-btn');
  
  copyButtons.forEach(btn => {
    btn.addEventListener('click', (e) => {
      const code = decodeURIComponent(btn.getAttribute('data-code'));
      copyToClipboard(code);
      e.stopPropagation();
    });
  });
};

// 复制文本到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text);
    showCopySuccess.value = true;
    setTimeout(() => {
      showCopySuccess.value = false;
    }, 2000);
  } catch (err) {
    console.error('无法复制到剪贴板:', err);
    // 降级方案
    const textArea = document.createElement('textarea');
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);
    showCopySuccess.value = true;
    setTimeout(() => {
      showCopySuccess.value = false;
    }, 2000);
  }
};

// 渲染Markdown内容
const renderContent = () => {
  if (props.note && props.note.content) {
    setupMarked();
    renderedContent.value = marked(props.note.content);
    
    nextTick(() => {
      addCopyFunctionality();
    });
  }
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

// 监听props变化重新渲染
watch(() => props.note, () => {
  renderContent();
}, { deep: true });

// 组件加载时渲染内容
onMounted(() => {
  renderContent();
  window.addEventListener('scroll', handleScroll);
});

// 组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
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

/* 代码块复制按钮样式 */
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
</style> 