<template>
  <div class="flex flex-col gap-4">
    <!-- 编辑器部分 -->
    <div class="w-full mb-2">
      <div class="mb-3 flex flex-wrap justify-between items-center gap-2 bg-gray-50 dark:bg-gray-750 p-2 rounded-t-md border border-gray-200 dark:border-gray-600">
        <div class="flex flex-wrap gap-1">
          <button class="tool-btn" @click="insertMd('# ')" title="一级标题">H1</button>
          <button class="tool-btn" @click="insertMd('## ')" title="二级标题">H2</button>
          <button class="tool-btn" @click="insertMd('### ')" title="三级标题">H3</button>
          <span class="tool-divider"></span>
          <button class="tool-btn" @click="insertMd('**粗体**')" title="粗体">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 10h8" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 14h8" />
              <path d="M12 4.5V19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
          <button class="tool-btn" @click="insertMd('*斜体*')" title="斜体">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
            </svg>
          </button>
          <span class="tool-divider"></span>
          <button class="tool-btn" @click="insertMd('[链接](https://example.com)')" title="链接">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.172 13.828a4 4 0 015.656 0l4 4a4 4 0 01-5.656 5.656l-1.102-1.101" />
            </svg>
          </button>
          <button class="tool-btn" @click="insertMd('![图片描述](图片链接)')" title="图片">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </button>
          <button class="tool-btn" @click="insertMd('\n```js\n// 代码块\n```\n')" title="代码块">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
            </svg>
          </button>
          <span class="tool-divider"></span>
          <button class="tool-btn" @click="insertMd('\n- 列表项1\n- 列表项2\n')" title="无序列表">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <button class="tool-btn" @click="insertMd('\n1. 列表项1\n2. 列表项2\n')" title="有序列表">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </button>
        </div>
        <div class="flex space-x-2">
          <button class="tool-btn-action" @click="clearContent" title="清空内容">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            清空
          </button>
        </div>
      </div>
      <textarea 
        ref="textareaRef"
        v-model="content" 
        class="w-full border border-gray-300 dark:border-gray-600 rounded-b-md p-3 min-h-[300px] sm:min-h-[400px] h-full resize-none focus:outline-none focus:ring-2 focus:ring-primary-500 font-mono dark:bg-gray-800 dark:text-gray-100 text-sm leading-relaxed"
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
import Prism from 'prismjs';

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

// 一个简单的JavaScript代码高亮函数
const highlightJS = (code) => {
  // 仅针对JS代码进行简单的语法高亮
  return code
    // 字符串 - 支持多行字符串
    .replace(/(["'`])([\s\S]*?)\1/g, '<span class="code-string">$1$2$1</span>')
    // 关键字
    .replace(/\b(const|let|var|function|return|if|else|for|while|class|new|import|export|from|this|async|await|try|catch|throw|true|false|null|undefined)\b/g, 
      '<span class="code-keyword">$1</span>')
    // 数字
    .replace(/\b(\d+(\.\d+)?)\b/g, '<span class="code-number">$1</span>')
    // 注释（处理单行注释）
    .replace(/(\/\/[^\n]*)/g, '<span class="code-comment">$1</span>')
    // 注释（处理多行注释）
    .replace(/(\/\*[\s\S]*?\*\/)/g, '<span class="code-comment">$1</span>');
};

// 监听窗口大小变化
onMounted(() => {
  window.addEventListener('resize', () => {
    isLargeScreen.value = window.innerWidth >= 1024;
    if (isLargeScreen.value) {
      mobileView.value = 'edit'; // 恢复默认视图
    }
  });
});

// 应用Prism进行代码高亮
const applyPrismHighlight = () => {
  if (previewRef.value) {
    // 修复Prism.js选择器问题
    const customPrismCSS = `
      .token.operator { background: transparent !important; }
      code .token { background: transparent !important; }
    `;
    
    // 添加自定义样式
    if (!document.getElementById('custom-prism-fix')) {
      const styleEl = document.createElement('style');
      styleEl.id = 'custom-prism-fix';
      styleEl.textContent = customPrismCSS;
      document.head.appendChild(styleEl);
    }
    
    const codeBlocks = previewRef.value.querySelectorAll('pre code');
    codeBlocks.forEach((block) => {
      Prism.highlightElement(block);
    });
  }
};

// 自定义Markdown渲染器，添加复制按钮到代码块
const setupMarked = () => {
  const renderer = new marked.Renderer();
  
  // 保存原始代码块渲染方法
  const originalCodeRenderer = renderer.code;
  
  // 自定义代码块渲染
  renderer.code = function(code, language, isEscaped) {
    // 处理新版本marked传递的对象格式
    if (typeof code === 'object' && code !== null) {
      // 如果code是对象，直接提取text属性
      if (code.text !== undefined) {
        language = code.lang || language;
        code = code.text;
      } else {
        // 尝试将整个对象转为字符串
        try {
          code = JSON.stringify(code);
        } catch (e) {
          console.error('无法将代码对象转换为字符串:', e);
          code = '';
        }
      }
    }
    
    // 确保code是字符串类型
    if (code === null || code === undefined) {
      code = '';
    } else if (typeof code !== 'string') {
      // 如果code不是字符串，尝试将其转换为字符串
      try {
        code = String(code);
      } catch (e) {
        console.error('无法将代码转换为字符串:', e);
        code = '';
      }
    }
    
    // 处理语言标识
    let displayLang = language || 'text';
    // 处理一些语言别名，以匹配Prism的语言名称
    if (displayLang === 'js') displayLang = 'javascript';
    if (displayLang === 'ts') displayLang = 'typescript';
    if (displayLang === 'html') displayLang = 'markup';
    if (displayLang === 'shell') displayLang = 'bash';
    
    // 通过转义字符来确保代码正确显示
    let escapedCode = code
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
    
    // 添加语言标识和固定显示的复制按钮
    return `
      <div class="code-block-wrapper">
        <div class="code-header">
          <span class="code-language">${displayLang}</span>
          <button class="code-copy-btn" data-code="${encodeURIComponent(code)}">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
            </svg>
            <span>复制</span>
          </button>
        </div>
        <pre class="code-block"><code class="language-${displayLang}">${escapedCode}</code></pre>
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
  
  // 确保markdown是字符串
  if (markdown === null || markdown === undefined) {
    markdown = '';
  } else if (typeof markdown !== 'string') {
    try {
      markdown = String(markdown);
      console.log('将非字符串markdown转换为字符串:', markdown);
    } catch (e) {
      console.error('无法将markdown转换为字符串:', e);
      markdown = '';
    }
  }
  
  return marked(markdown);
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
  
  // 更新渲染内容后添加复制功能和应用语法高亮
  nextTick(() => {
    addCopyFunctionality();
    applyPrismHighlight(); // 应用Prism高亮
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

// 清空内容
const clearContent = () => {
  if (confirm('确定要清空所有内容吗？此操作无法撤销。')) {
    content.value = '';
    updateContent();
  }
};
</script>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.animate-fade-in-out {
  animation: fadeInOut 2s forwards;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translateY(10px); }
  10% { opacity: 1; transform: translateY(0); }
  90% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(10px); }
}

/* 新增的工具按钮样式 */
.tool-btn {
  @apply px-2 py-1 rounded text-xs text-gray-700 dark:text-gray-200 bg-gray-100 dark:bg-gray-700 
  hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors;
  min-width: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.tool-btn-action {
  @apply px-2 py-1 rounded text-xs text-gray-700 dark:text-gray-200 bg-gray-100 dark:bg-gray-700 
  hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors flex items-center gap-1;
}

.tool-divider {
  @apply w-px h-5 bg-gray-300 dark:bg-gray-600 mx-1;
}

/* 响应式设计优化 */
@media (max-width: 640px) {
  .tool-btn {
    @apply p-1;
    min-width: 24px;
  }
  
  .tool-text {
    display: none;
  }
}

/* 编辑器和预览区域 */
.prose {
  max-width: 100%;
}

.prose img {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 1rem auto;
}

.prose pre {
  background-color: #f8f8f8;
  padding: 1rem;
  border-radius: 4px;
  overflow-x: auto;
}

/* 预览区代码块样式优化 */
.prose code {
  font-size: 0.9em;
  background: #f5f5f5;
  padding: 0.1em 0.3em;
  border-radius: 3px;
}

.prose pre code {
  background: transparent;
  padding: 0;
}

/* 代码块样式 - 与MarkdownRenderer组件保持一致 */
.prose .code-block-wrapper {
  margin: 1.5em 0;
  border-radius: 0.5rem;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  background: #f8f9fa;
}

.dark .prose .code-block-wrapper {
  background: #1f2937;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

/* 代码块头部 */
.prose .code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  background: #2d3748;
  color: #e2e8f0;
}

.dark .prose .code-header {
  background: #2d3748; /* 改为更深色以增加对比度 */
  border-bottom: 1px solid #4b5563; /* 增强边框对比度 */
  color: #f3f4f6; /* 更亮的文字 */
}

.prose .code-language {
  font-family: ui-monospace, monospace;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
}

.dark .prose .code-language {
  background: rgba(255, 255, 255, 0.25); /* 提高亮度以增加对比度 */
  color: #ffffff; /* 确保文字清晰可见 */
  font-weight: 600; /* 加粗以增强可读性 */
  letter-spacing: 0.03em; /* 调整字间距 */
}

/* 复制按钮 */
.prose .code-copy-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: #e2e8f0;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 0.75rem;
  transition: all 0.2s ease;
}

.dark .prose .code-copy-btn {
  background: rgba(255, 255, 255, 0.15); /* 提高亮度 */
  color: #ffffff; /* 确保文字清晰可见 */
}

.prose .code-copy-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.dark .prose .code-copy-btn:hover {
  background: rgba(255, 255, 255, 0.25); /* 提高悬停亮度 */
}

.prose .code-copy-btn svg {
  width: 0.875rem;
  height: 0.875rem;
}

/* 代码块本身 */
.prose .code-block {
  margin: 0;
  padding: 1rem;
  overflow-x: auto;
  background: #f8f9fa;
  font-size: 0.9rem;
  line-height: 1.5;
  white-space: pre;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  color: #24292e;
}

.dark .prose .code-block {
  background: #1f2937;
  color: #e5e7eb;
}

/* 代码语法高亮 */
.prose .code-string {
  color: #e83e8c;
}

.dark .prose .code-string {
  color: #ff79c6;
}

.prose .code-keyword {
  color: #0076ff;
  font-weight: bold;
}

.dark .prose .code-keyword {
  color: #ff79c6;
}

.prose .code-comment {
  color: #6a737d;
  font-style: italic;
}

.dark .prose .code-comment {
  color: #6272a4;
}

.prose .code-number {
  color: #005cc5;
}

.dark .prose .code-number {
  color: #bd93f9;
}

/* Prism.js样式 - 与MarkdownRenderer组件保持一致 */
/* 增强暗色模式下的代码高亮可视性 */
.dark .prose .token.comment,
.dark .prose .token.prolog,
.dark .prose .token.doctype,
.dark .prose .token.cdata {
  color: #a8bbd0; /* 提高注释可读性 */
}

.dark .prose .token.punctuation {
  color: #ffffff; /* 让标点更明显 */
}

.dark .prose .token.operator {
  color: #ffffff; /* 让操作符(包括等号)更明显 */
  background: transparent !important; /* 确保没有背景色 */
}

.dark .prose .token.string {
  color: #7efa9b; /* 提高字符串亮度 */
}

.dark .prose .token.keyword {
  color: #ff8dd7; /* 稍微提高关键字亮度 */
  font-weight: bold; /* 加粗关键字 */
}

.dark .prose .token.function {
  color: #67e2ff; /* 提高函数名亮度 */
}

.dark .prose .token.number {
  color: #d3a4ff; /* 提高数字亮度 */
}

/* 修复等号被选中的问题 */
.dark .prose code .token {
  background: transparent !important; /* 移除所有token元素的背景 */
}

.dark .prose .token.operator {
  background: transparent !important; /* 特别指定操作符没有背景 */
}

/* 提高整体代码区域亮度 */
.dark .prose .code-block {
  background: #1a2234; /* 略微调整代码块背景以增加对比度 */
}

.dark .prose .code-block code {
  color: #f1f5f9; /* 提高默认文本亮度 */
}
</style> 