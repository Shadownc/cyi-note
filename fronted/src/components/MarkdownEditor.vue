<template>
  <div class="flex flex-col gap-4">
    <!-- 编辑器部分 -->
    <div class="w-full mb-2">
      <div class="mb-2 flex flex-wrap justify-between items-center gap-2">
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200">编辑</h3>
        <div class="flex flex-wrap gap-1">
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('# ')">H1</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('## ')">H2</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('### ')">H3</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('**粗体**')">粗体</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('*斜体*')">斜体</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('[链接](https://example.com)')">链接</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('\n```js\n// 代码块\n```\n')">代码块</button>
        </div>
      </div>
      <textarea 
        ref="textareaRef"
        v-model="content" 
        class="w-full border border-gray-300 dark:border-gray-600 rounded-md p-3 min-h-[200px] sm:min-h-[300px] h-full resize-none focus:outline-none focus:ring-2 focus:ring-primary-500 font-mono dark:bg-gray-800 dark:text-gray-100"
        placeholder="在这里输入Markdown内容..."
        @input="updateContent"
      ></textarea>
    </div>
    
    <!-- 编辑/预览切换按钮（移动端） -->
    <div class="lg:hidden flex justify-center mb-2">
      <div class="inline-flex rounded-md shadow-sm" role="group">
        <button 
          type="button" 
          @click="mobileView = 'edit'"
          :class="[
            'px-4 py-2 text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-l-lg',
            mobileView === 'edit' ? 'bg-primary-500 text-white dark:bg-primary-600' : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200'
          ]"
        >
          编辑
        </button>
        <button 
          type="button" 
          @click="mobileView = 'preview'"
          :class="[
            'px-4 py-2 text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-r-lg',
            mobileView === 'preview' ? 'bg-primary-500 text-white dark:bg-primary-600' : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200'
          ]"
        >
          预览
        </button>
      </div>
    </div>
    
    <!-- 预览部分 (桌面端始终显示，移动端仅在预览模式下显示) -->
    <div class="w-full" :class="{'hidden': mobileView === 'edit' && !isLargeScreen}">
      <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">预览</h3>
      <div class="border border-gray-300 dark:border-gray-600 rounded-md p-3 min-h-[200px] sm:min-h-[300px] prose dark:prose-invert max-w-none overflow-auto dark:bg-gray-800 dark:bg-opacity-60" v-html="renderedContent" ref="previewRef"></div>
    </div>
    
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
import { ref, defineProps, defineEmits, onMounted, watch, nextTick } from 'vue';
import { marked } from 'marked';

const props = defineProps({
  initialValue: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:content']);

const content = ref(props.initialValue);
const renderedContent = ref('');
const textareaRef = ref(null);
const previewRef = ref(null);
const mobileView = ref('edit'); // 'edit' 或 'preview'
const isLargeScreen = ref(window.innerWidth >= 1024); // 初始化时检查屏幕大小
const showCopySuccess = ref(false);

// 监听窗口大小变化
onMounted(() => {
  window.addEventListener('resize', () => {
    isLargeScreen.value = window.innerWidth >= 1024;
    if (isLargeScreen.value) {
      mobileView.value = 'edit'; // 恢复默认视图
    }
  });
});

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

// 将Markdown渲染为HTML
const renderMarkdown = (markdown) => {
  setupMarked();
  return marked(markdown || '');
};

// 为代码块添加复制功能
const addCopyFunctionality = () => {
  if (!previewRef.value) return;
  
  const copyButtons = previewRef.value.querySelectorAll('.code-copy-btn');
  
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

// 更新内容并通知父组件
const updateContent = () => {
  renderedContent.value = renderMarkdown(content.value);
  emit('update:content', content.value);
  
  // 更新渲染内容后添加复制功能
  nextTick(() => {
    addCopyFunctionality();
  });
};

// 在特定位置插入Markdown语法
const insertMd = (syntax) => {
  if (!textareaRef.value) return;
  
  const textarea = textareaRef.value;
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  
  content.value = content.value.substring(0, start) + syntax + content.value.substring(end);
  
  // 更新后重新设置光标位置
  setTimeout(() => {
    textarea.focus();
    textarea.selectionStart = start + syntax.length;
    textarea.selectionEnd = start + syntax.length;
    updateContent();
  }, 0);
};

// 监听初始值变化
watch(() => props.initialValue, (newVal) => {
  content.value = newVal;
  updateContent();
});

// 组件加载时渲染初始内容
onMounted(() => {
  updateContent();
});
</script>

<style>
/* 确保Markdown内容正确渲染 */
.prose {
  width: 100%;
}

.prose h1 {
  @apply text-2xl font-bold my-4 text-gray-900 dark:text-gray-100;
}

.prose h2 {
  @apply text-xl font-bold my-3 text-gray-800 dark:text-gray-200;
}

.prose h3 {
  @apply text-lg font-bold my-2 text-gray-800 dark:text-gray-200;
}

.prose p {
  @apply my-2 text-gray-700 dark:text-gray-300;
}

.prose ul {
  @apply list-disc pl-6 my-2 text-gray-700 dark:text-gray-300;
}

.prose ol {
  @apply list-decimal pl-6 my-2 text-gray-700 dark:text-gray-300;
}

.prose blockquote {
  @apply border-l-4 border-gray-300 dark:border-primary-700 pl-4 italic my-4 text-gray-600 dark:text-gray-400 dark:bg-gray-800 dark:bg-opacity-50 dark:py-1 dark:rounded-r-md;
}

.prose code {
  @apply bg-gray-100 dark:bg-gray-800 px-1 rounded text-sm text-gray-800 dark:text-primary-200 dark:border dark:border-gray-700;
}

.prose pre {
  @apply bg-gray-100 dark:bg-gray-900 p-3 rounded my-3 overflow-auto dark:border dark:border-gray-700;
}

.dark .prose pre code {
  @apply text-gray-200 bg-transparent border-0;
}

.dark .prose a {
  @apply text-primary-300 hover:text-primary-200 underline;
}

.dark .prose h1 {
  @apply text-primary-300 border-b pb-2 border-gray-700;
}

.dark .prose h2 {
  @apply text-primary-300 opacity-90 border-b pb-1 border-gray-700;
}

.dark .prose h3 {
  @apply text-primary-300 opacity-80;
}

.dark .prose strong {
  @apply text-primary-200;
}

.dark .prose table {
  @apply border-collapse border border-gray-700;
}

.dark .prose th {
  @apply bg-gray-800 text-gray-100 border border-gray-700 px-4 py-2;
}

.dark .prose td {
  @apply border border-gray-700 px-4 py-2;
}

/* 代码块复制按钮样式 */
.code-block-wrapper {
  position: relative;
}

.code-copy-btn {
  transition: opacity 0.2s ease;
}

.code-copy-btn:hover {
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

/* 移动端优化 */
@media (max-width: 768px) {
  .prose {
    font-size: 0.95rem;
  }
  
  .prose pre {
    max-width: 100%;
    overflow-x: auto;
  }
  
  .prose img {
    max-width: 100%;
    height: auto;
  }
  
  .btn {
    padding: 0.25rem 0.5rem;
  }
  
  /* 在移动端显示复制按钮 */
  .code-copy-btn {
    opacity: 0.7 !important;
  }
}
</style> 