<template>
  <div class="markdown-renderer" v-html="renderedContent" ref="contentRef"></div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue';
import { marked } from 'marked';
import Prism from 'prismjs';
// 导入基本主题样式，您可以更换为其他主题
import 'prismjs/themes/prism.css';
// 导入常用语言
import 'prismjs/components/prism-javascript';
import 'prismjs/components/prism-typescript';
import 'prismjs/components/prism-css';
import 'prismjs/components/prism-markup';
import 'prismjs/components/prism-python';
import 'prismjs/components/prism-java';
import 'prismjs/components/prism-c';
import 'prismjs/components/prism-cpp';
import 'prismjs/components/prism-csharp';
import 'prismjs/components/prism-jsx';
import 'prismjs/components/prism-json';
import 'prismjs/components/prism-bash';
import 'prismjs/components/prism-yaml';
import 'prismjs/components/prism-markdown';
import 'prismjs/components/prism-sql';

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
});

const renderedContent = ref('');
const contentRef = ref(null);
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

// 自定义Markdown渲染器
const setupMarked = () => {
  const renderer = new marked.Renderer();
  
  // 自定义代码块渲染
  renderer.code = function(code, language, isEscaped) {
    // 记录代码块信息以便调试
    console.log('代码块渲染信息:', { 
      codeType: typeof code, 
      codeValue: code,
      language: language,
      isEscaped: isEscaped 
    });
    
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
  if (!contentRef.value) return;
  
  const copyButtons = contentRef.value.querySelectorAll('.code-copy-btn');
  
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
    // 确保复制的是原始代码（无HTML）
    await navigator.clipboard.writeText(text);
    
    // 显示复制成功提示
    const tooltip = document.createElement('div');
    tooltip.className = 'copy-success-tooltip';
    tooltip.textContent = '已复制到剪贴板';
    document.body.appendChild(tooltip);
    
    // 2秒后移除提示
    setTimeout(() => {
      document.body.removeChild(tooltip);
    }, 2000);
  } catch (err) {
    console.error('无法复制到剪贴板:', err);
    // 降级方案
    const textArea = document.createElement('textarea');
    textArea.value = text;
    textArea.style.position = 'fixed';
    textArea.style.top = '0';
    textArea.style.left = '0';
    textArea.style.width = '2em';
    textArea.style.height = '2em';
    textArea.style.opacity = '0';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);
    
    // 显示复制成功提示
    const tooltip = document.createElement('div');
    tooltip.className = 'copy-success-tooltip';
    tooltip.textContent = '已复制到剪贴板';
    document.body.appendChild(tooltip);
    
    // 2秒后移除提示
    setTimeout(() => {
      document.body.removeChild(tooltip);
    }, 2000);
  }
};

// 应用Prism进行代码高亮
const applyPrismHighlight = () => {
  if (contentRef.value) {
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
    
    const codeBlocks = contentRef.value.querySelectorAll('pre code');
    codeBlocks.forEach((block) => {
      Prism.highlightElement(block);
    });
  }
};

// 更新渲染内容
const updateContent = () => {
  try {
    console.log('准备渲染内容, 类型:', typeof props.content, '值:', props.content);
    renderedContent.value = renderMarkdown(props.content);
    
    // 更新渲染内容后添加复制功能和语法高亮
    nextTick(() => {
      addCopyFunctionality();
      applyPrismHighlight(); // 应用Prism高亮
    });
  } catch (error) {
    console.error('渲染Markdown内容时出错:', error);
    // 在发生错误时显示一个友好的错误消息
    renderedContent.value = `<div class="markdown-error">无法渲染内容: ${error.message}</div>`;
  }
};

// 监听内容变化
watch(() => props.content, (newVal) => {
  updateContent();
});

// 组件加载时渲染初始内容
onMounted(() => {
  updateContent();
});
</script>

<style>
/* Markdown内容全局样式 */
.markdown-renderer {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  line-height: 1.6;
  color: #374151;
}

.dark .markdown-renderer {
  color: #e5e7eb;
}

/* 标题样式 */
.markdown-renderer h1 {
  font-size: 2em;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 700;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 0.3em;
}

.markdown-renderer h2 {
  font-size: 1.5em;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 0.3em;
}

.markdown-renderer h3 {
  font-size: 1.25em;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
}

.markdown-renderer h4 {
  font-size: 1em;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
}

/* 段落样式 */
.markdown-renderer p {
  margin-top: 0.75em;
  margin-bottom: 0.75em;
}

/* 列表样式 */
.markdown-renderer ul, .markdown-renderer ol {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
  padding-left: 2em;
}

.markdown-renderer ul {
  list-style-type: disc;
}

.markdown-renderer ol {
  list-style-type: decimal;
}

.markdown-renderer li {
  margin-top: 0.25em;
  margin-bottom: 0.25em;
}

/* 链接样式 */
.markdown-renderer a {
  color: #3b82f6;
  text-decoration: none;
}

.markdown-renderer a:hover {
  text-decoration: underline;
}

/* 图片样式 */
.markdown-renderer img {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 1em auto;
  border-radius: 0.375rem;
}

/* 引用样式 */
.markdown-renderer blockquote {
  border-left: 4px solid #e5e7eb;
  padding-left: 1em;
  margin-left: 0;
  margin-right: 0;
  font-style: italic;
  color: #6b7280;
}

.dark .markdown-renderer blockquote {
  border-left-color: #4b5563;
  color: #9ca3af;
}

/* 代码样式 - 改进版 */
.markdown-renderer code {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.875em;
  background-color: #f3f4f6;
  padding: 0.2em 0.4em;
  border-radius: 0.25em;
}

.dark .markdown-renderer code {
  background-color: #374151;
}

/* 内联代码 */
.markdown-renderer p code {
  color: #d63384; /* 内联代码使用粉色调突出显示 */
  background-color: rgba(214, 51, 132, 0.1);
  padding: 0.2em 0.4em;
  border-radius: 0.25em;
}

.dark .markdown-renderer p code {
  color: #e879f9;
  background-color: rgba(232, 121, 249, 0.1);
}

/* 代码块容器 */
.markdown-renderer .code-block-wrapper {
  margin: 1.5em 0;
  border-radius: 0.5rem;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  background: #f8f9fa;
}

.dark .markdown-renderer .code-block-wrapper {
  background: #1f2937;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

/* 代码块头部 */
.markdown-renderer .code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  background: #2d3748;
  color: #e2e8f0;
}

.dark .markdown-renderer .code-header {
  background: #2d3748; /* 改为更深色以增加对比度 */
  border-bottom: 1px solid #4b5563; /* 增强边框对比度 */
  color: #f3f4f6; /* 更亮的文字 */
}

.markdown-renderer .code-language {
  font-family: ui-monospace, monospace;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
}

.dark .markdown-renderer .code-language {
  background: rgba(255, 255, 255, 0.25); /* 提高亮度以增加对比度 */
  color: #ffffff; /* 确保文字清晰可见 */
  font-weight: 600; /* 加粗以增强可读性 */
  letter-spacing: 0.03em; /* 调整字间距 */
}

/* 复制按钮 */
.markdown-renderer .code-copy-btn {
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

.dark .markdown-renderer .code-copy-btn {
  background: rgba(255, 255, 255, 0.15); /* 提高亮度 */
  color: #ffffff; /* 确保文字清晰可见 */
}

.markdown-renderer .code-copy-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.dark .markdown-renderer .code-copy-btn:hover {
  background: rgba(255, 255, 255, 0.25); /* 提高悬停亮度 */
}

.markdown-renderer .code-copy-btn svg {
  width: 0.875rem;
  height: 0.875rem;
}

/* 代码块本身 */
.markdown-renderer .code-block {
  margin: 0;
  padding: 1rem;
  overflow-x: auto;
  background: #f8f9fa;
  font-size: 0.9rem;
  line-height: 1.5;
}

.dark .markdown-renderer .code-block {
  background: #1f2937;
}

.markdown-renderer .code-block code {
  background: transparent;
  padding: 0;
  color: #24292e;
  white-space: pre;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

.dark .markdown-renderer .code-block code {
  color: #e5e7eb;
}

/* Prism.js样式覆盖，使其适应我们的主题 */
/* 确保在暗黑模式下代码高亮正常显示 */
.dark .token.comment,
.dark .token.prolog,
.dark .token.doctype,
.dark .token.cdata {
  color: #8292a2;
}

.dark .token.punctuation {
  color: #f8f8f2;
}

.dark .token.namespace {
  opacity: .7;
}

.dark .token.property,
.dark .token.tag,
.dark .token.constant,
.dark .token.symbol,
.dark .token.deleted {
  color: #ff79c6;
}

.dark .token.boolean,
.dark .token.number {
  color: #bd93f9;
}

.dark .token.selector,
.dark .token.attr-name,
.dark .token.string,
.dark .token.char,
.dark .token.builtin,
.dark .token.inserted {
  color: #50fa7b;
}

.dark .token.operator,
.dark .token.entity,
.dark .token.url,
.dark .language-css .token.string,
.dark .style .token.string {
  color: #f8f8f2;
}

.dark .token.atrule,
.dark .token.attr-value,
.dark .token.keyword {
  color: #ff79c6;
}

.dark .token.function,
.dark .token.class-name {
  color: #8be9fd;
}

.dark .token.regex,
.dark .token.important,
.dark .token.variable {
  color: #f1fa8c;
}

.dark .token.important,
.dark .token.bold {
  font-weight: bold;
}

.dark .token.italic {
  font-style: italic;
}

.dark .token.entity {
  cursor: help;
}

/* 增强暗色模式下的代码高亮可视性 */
/* 增强亮度和对比度 */
.dark .token.comment,
.dark .token.prolog,
.dark .token.doctype,
.dark .token.cdata {
  color: #a8bbd0; /* 提高注释可读性 */
}

.dark .token.punctuation {
  color: #ffffff; /* 让标点更明显 */
}

.dark .token.operator {
  color: #ffffff; /* 让操作符(包括等号)更明显 */
  background: transparent !important; /* 确保没有背景色 */
}

.dark .token.string {
  color: #7efa9b; /* 提高字符串亮度 */
}

.dark .token.keyword {
  color: #ff8dd7; /* 稍微提高关键字亮度 */
  font-weight: bold; /* 加粗关键字 */
}

.dark .token.function {
  color: #67e2ff; /* 提高函数名亮度 */
}

.dark .token.number {
  color: #d3a4ff; /* 提高数字亮度 */
}

/* 修复等号被选中的问题 */
.dark code .token {
  background: transparent !important; /* 移除所有token元素的背景 */
}

.dark .token.operator {
  background: transparent !important; /* 特别指定操作符没有背景 */
}

/* 提高整体代码区域亮度 */
.dark .markdown-renderer .code-block {
  background: #1a2234; /* 略微调整代码块背景以增加对比度 */
}

.dark .markdown-renderer .code-block code {
  color: #f1f5f9; /* 提高默认文本亮度 */
}

/* 水平线样式 */
.markdown-renderer hr {
  border: 0;
  border-top: 1px solid #e5e7eb;
  margin: 2em 0;
}

.dark .markdown-renderer hr {
  border-top-color: #4b5563;
}

/* 表格样式 */
.markdown-renderer table {
  border-collapse: collapse;
  width: 100%;
  margin-top: 1em;
  margin-bottom: 1em;
  font-size: 0.9em;
}

.markdown-renderer th, .markdown-renderer td {
  border: 1px solid #e5e7eb;
  padding: 0.75em 1em;
  text-align: left;
}

.markdown-renderer th {
  background-color: #f9fafb;
  font-weight: 600;
}

.dark .markdown-renderer th, .dark .markdown-renderer td {
  border-color: #4b5563;
}

.dark .markdown-renderer th {
  background-color: #374151;
}

/* 复制成功提示样式 */
.copy-success-tooltip {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  background-color: #10b981;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  z-index: 50;
  font-size: 0.875rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  animation: fadeInOut 2s forwards;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translateY(10px); }
  10% { opacity: 1; transform: translateY(0); }
  90% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(10px); }
}

/* 代码块语法高亮样式 */
.markdown-renderer .code-string {
  color: #e83e8c;
}

.dark .markdown-renderer .code-string {
  color: #ff79c6;
}

.markdown-renderer .code-keyword {
  color: #0076ff;
  font-weight: bold;
}

.dark .markdown-renderer .code-keyword {
  color: #ff79c6;
}

.markdown-renderer .code-comment {
  color: #6a737d;
  font-style: italic;
}

.dark .markdown-renderer .code-comment {
  color: #6272a4;
}

.markdown-renderer .code-number {
  color: #005cc5;
}

.dark .markdown-renderer .code-number {
  color: #bd93f9;
}

/* 错误消息样式 */
.markdown-renderer .markdown-error {
  padding: 1rem;
  margin: 1rem 0;
  border-radius: 0.5rem;
  background-color: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.875rem;
}

.dark .markdown-renderer .markdown-error {
  background-color: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.4);
  color: #f87171;
}
</style> 