<template>
  <div class="markdown-renderer" v-html="renderedContent" ref="contentRef"></div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick, defineProps } from 'vue';
import { marked } from 'marked';
import Prism from 'prismjs';
import { loadAuthenticatedImage } from '@/utils/imageHelper';
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
  
  // 自定义图片渲染
  renderer.image = function(href, title, text) {
    let imgClass = 'markdown-image';
    let loadingClass = '';
    
    console.log('图片渲染参数:', { href: href, title: title, text: text, hrefType: typeof href });
    
    // 特殊处理object类型的href
    if (href && typeof href === 'object') {
      if (href.href) {
        console.log('从对象中提取href属性:', href.href);
        href = href.href;
      } else if (href.src) {
        console.log('从对象中提取src属性:', href.src);
        href = href.src;
      } else if (href.url) {
        console.log('从对象中提取url属性:', href.url);
        href = href.url;
      } else {
        // 尝试转换对象为字符串
        try {
          const objString = String(href);
          console.log('对象转换为字符串:', objString);
          // 如果不是[object Object]，则使用转换后的字符串
          if (objString !== '[object Object]') {
            href = objString;
          } else {
            // 如果是[object Object]，尝试序列化
            const jsonString = JSON.stringify(href);
            console.log('对象序列化为JSON:', jsonString);
            href = '';
          }
        } catch (e) {
          console.error('无法将href对象转换为字符串:', e);
          href = '';
        }
      }
    }
    
    // 确保href是字符串类型
    if (href === null || href === undefined) {
      href = '';
    } else if (typeof href !== 'string') {
      try {
        href = String(href);
        console.log('强制将非字符串href转换为字符串:', href);
      } catch (e) {
        console.error('无法将href转换为字符串:', e);
        href = '';
      }
    }
    
    // 检测临时占位符
    if (href && href.includes('img_placeholder_')) {
      loadingClass = 'image-loading';
      imgClass += ' placeholder-image';
    }
    
    // 特殊处理blob URL
    if (href && href.startsWith('blob:')) {
      console.log('处理blob URL:', href);
    }
    
    // 添加错误处理和加载状态
    return `
      <div class="image-container ${loadingClass}">
        <img 
          src="${href}" 
          alt="${text || ''}" 
          title="${title || text || ''}" 
          class="${imgClass}" 
          onerror="this.onerror=null; this.classList.add('image-error'); this.src='data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9ImN1cnJlbnRDb2xvciIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiIGNsYXNzPSJsdWNpZGUgbHVjaWRlLWFsZXJ0LXRyaWFuZ2xlIj48cGF0aCBkPSJtMjEuNzMgMTgtOC0xNGEyIDIgMCAwIDAtMy40NiAwbC04IDE0QTIgMiAwIDAgMCA0IDE5LjVoMTZhMiAyIDAgMCAwIDEuNzMtMS41WiIvPjxsaW5lIHgxPSIxMiIgeTE9IjgiIHgyPSIxMiIgeTI9IjEyIi8+PGxpbmUgeDE9IjEyIiB5MT0iMTYiIHgyPSIxMi4wMSIgeTI9IjE2Ii8+PC9zdmc+'; this.parentNode.classList.add('image-error-container'); console.error('图片加载失败:', this.src);"
          onload="this.parentNode.classList.remove('image-loading'); console.log('图片加载成功:', this.src);"
        />
        ${loadingClass ? '<div class="image-loading-indicator"><div class="spinner"></div></div>' : ''}
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
  
  // 处理图片URL，确保路径正确 
  const processedMarkdown = markdown.replace(
    /!\[(.*?)\]\((.*?)\)/g,
    (match, alt, url) => {
      // 确保url是字符串
      if (url === null || url === undefined) {
        url = '';
      } else if (typeof url !== 'string') {
        try {
          url = String(url);
          console.log('URL处理: 非字符串URL转换为字符串:', url);
        } catch (e) {
          console.error('URL处理: 无法将URL转换为字符串:', e);
          url = '';
        }
      }
      
      // 记录处理前的URL
      console.log('URL处理: 原始URL:', url);
      
      // 检查是否是blob URL，如果是则尝试从文件名提取ID
      if (url && url.startsWith('blob:')) {
        console.log('URL处理: 检测到blob URL:', url);
        // 尝试在内容中查找这个 blob URL 对应的附件 ID
        // 这可能需要根据您的数据结构调整
        const blobFileName = alt || '';
        if (blobFileName) {
          // 记录这个blob URL以便后续处理
          console.log('URL处理: blob URL对应的文件名:', blobFileName);
          // 返回原始URL，但在renderMarkdown后会尝试特殊处理这个URL
          return `![${alt}](${url})`;
        }
      }
      
      // 特殊处理img_placeholder标识符 - 保持不变
      if (url && url.includes('img_placeholder_')) {
        console.log('URL处理: 保留占位符URL:', url);
        return `![${alt}](${url})`;
      }
      
      // 如果URL是附件API格式，确保正确处理
      if (url && url.includes('/api/attachments/')) {
        // 对于已经是完整API格式的URL，确保添加认证信息
        if (!url.startsWith('blob:')) {
          // 尝试加载已认证的图片
          loadAuthenticatedImageAsync(url).then(blobUrl => {
            if (blobUrl) {
              console.log(`URL处理: 已加载认证图片并创建blob URL: ${url} -> ${blobUrl}`);
              
              // 找到渲染后的图片并替换src
              setTimeout(() => {
                const renderedImages = document.querySelectorAll(`img[src="${url}"]`);
                renderedImages.forEach(img => {
                  img.src = blobUrl;
                  console.log(`已替换渲染图片的src: ${url} -> ${blobUrl}`);
                });
              }, 100);
            }
          }).catch(err => {
            console.error('URL处理: 加载认证图片失败:', err);
          });
        }
        return `![${alt}](${url})`;
      }
      
      // 如果URL是相对路径且不以/开头，添加/api前缀
      if (url && !url.startsWith('http') && !url.startsWith('/') && !url.startsWith('blob:')) {
        const newUrl = `/api/${url}`;
        console.log(`URL处理: 相对路径添加/api前缀: ${url} -> ${newUrl}`);
        return `![${alt}](${newUrl})`;
      }
      
      // 其他情况保持URL不变
      return `![${alt}](${url})`;
    }
  );
  
  return marked(processedMarkdown);
};

// 加载认证图片的异步函数
async function loadAuthenticatedImageAsync(url) {
  const token = localStorage.getItem('token');
  if (!token) {
    console.error('未提供认证令牌，无法加载图片');
    return null;
  }
  
  try {
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
    let fullUrl = url;
    
    // 处理相对URL
    if (url.startsWith('/api/')) {
      fullUrl = url.replace('/api/', `${baseUrl}/`);
    }
    
    // 创建XHR请求
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();
      xhr.open('GET', fullUrl, true);
      xhr.setRequestHeader('Authorization', `Bearer ${token}`);
      xhr.responseType = 'blob';
      
      xhr.onload = function() {
        if (this.status === 200) {
          const blob = new Blob([this.response]);
          const blobUrl = URL.createObjectURL(blob);
          console.log('已创建认证图片的blob URL:', blobUrl);
          resolve(blobUrl);
        } else {
          console.error('获取图片失败:', this.status);
          reject(new Error(`HTTP错误 ${this.status}`));
        }
      };
      
      xhr.onerror = function() {
        console.error('图片加载请求失败');
        reject(new Error('网络请求失败'));
      };
      
      xhr.send();
    });
  } catch (error) {
    console.error('加载认证图片异步处理失败:', error);
    return null;
  }
}

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

// 处理图片URL
const processImageUrl = (url) => {
  if (!url) return '';
  
  // 记录原始URL
  console.log('处理图片URL:', url);
  
  // 如果是Blob URL或data URL，保持不变
  if (url.startsWith('blob:') || url.startsWith('data:')) {
    console.log('保留原始Blob/Data URL:', url);
    return url;
  }
  
  // 以#或?开头的锚点或查询参数，不处理
  if (url.startsWith('#') || url.startsWith('?')) {
    return url;
  }
  
  // 确保API附件路径正确，但不添加token，因为我们会在图片加载时处理认证
  if (url.startsWith('/api/attachments/')) {
    console.log('保留完整API路径:', url);
    return url;
  }
  
  // 如果是相对路径但不是以/开头，添加/api前缀
  if (!url.startsWith('http') && !url.startsWith('/')) {
    const newUrl = `/api/${url}`;
    console.log(`添加API前缀: ${url} -> ${newUrl}`);
    return newUrl;
  }
  
  // 如果是以/开头但不是以/api开头，添加/api前缀
  if (url.startsWith('/') && !url.startsWith('/api/')) {
    const newUrl = `/api${url}`;
    console.log(`添加API前缀: ${url} -> ${newUrl}`);
    return newUrl;
  }
  
  // 其他情况，保留原URL
  return url;
};

// 处理需要认证的图片
const handleAuthenticatedImage = (img) => {
  const src = img.getAttribute('src');
  
  // 只处理需要认证的API路径
  if (src && src.startsWith('/api/attachments/')) {
    console.log('处理需要认证的图片URL:', src);
    
    // 添加加载中的样式
    const parent = img.parentNode;
    if (parent) {
      parent.classList.add('image-loading');
    }
    
    // 使用工具函数加载认证图片
    loadAuthenticatedImage(src)
      .then(blobUrl => {
        if (blobUrl !== src) {
          // 替换为blob URL
          img.src = blobUrl;
          console.log('已替换为认证图片的blob URL:', blobUrl);
          
          // 移除加载状态
          if (parent) {
            parent.classList.remove('image-loading');
            parent.classList.add('image-loaded');
          }
        }
      })
      .catch(error => {
        console.error('加载认证图片失败:', error);
        img.onerror();
      });
  }
};

// 处理复制按钮点击事件
const handleCopyClick = (e) => {
  const btn = e.currentTarget;
  const code = decodeURIComponent(btn.getAttribute('data-code'));
  copyToClipboard(code);
  e.stopPropagation();
};

// 尝试从blob URL中提取特征信息
const extractInfoFromBlobUrl = (url) => {
  if (!url || typeof url !== 'string') return null;
  
  // 一般blob URL格式: blob:http://domain/uuid
  try {
    const parts = url.split('/');
    const uuid = parts[parts.length - 1];
    return uuid;
  } catch (error) {
    console.error('无法从blob URL提取信息:', error);
    return null;
  }
};

// 保存和加载blob URL映射
const saveBlobMapping = (blobId, apiUrl) => {
  if (!blobId || !apiUrl) return;
  
  try {
    // 从localStorage获取现有映射
    const blobMap = JSON.parse(localStorage.getItem('blobUrlMap') || '{}');
    
    // 更新映射
    blobMap[blobId] = apiUrl;
    
    // 保存回localStorage
    localStorage.setItem('blobUrlMap', JSON.stringify(blobMap));
    console.log(`已保存blob映射: ${blobId} -> ${apiUrl}`);
  } catch (error) {
    console.error('保存blob URL映射失败:', error);
  }
};

// 获取blob URL的映射
const getBlobMapping = (blobId) => {
  if (!blobId) return null;
  
  try {
    // 从localStorage获取映射
    const blobMap = JSON.parse(localStorage.getItem('blobUrlMap') || '{}');
    
    // 返回对应的API URL
    return blobMap[blobId] || null;
  } catch (error) {
    console.error('获取blob URL映射失败:', error);
    return null;
  }
};

// 清理过期的blob映射
const cleanupBlobMappings = () => {
  try {
    const blobMap = JSON.parse(localStorage.getItem('blobUrlMap') || '{}');
    const now = Date.now();
    let count = 0;
    
    // 检查localStorage中是否已有过期时间记录
    const lastCleanup = parseInt(localStorage.getItem('lastBlobMapCleanup') || '0');
    
    // 如果距离上次清理不足1天，则跳过
    if (now - lastCleanup < 24 * 60 * 60 * 1000) {
      return;
    }
    
    // 遍历所有映射并检查其有效性
    Object.keys(blobMap).forEach(blobId => {
      // 这里我们可以添加额外的有效性检查
      // 或者设置一个过期时间（例如7天）
      // 不过由于映射是轻量级的，我们可以选择保留它们
    });
    
    // 记录清理时间
    localStorage.setItem('lastBlobMapCleanup', now.toString());
  } catch (error) {
    console.error('清理blob映射失败:', error);
  }
};

// 更新渲染内容
const updateContent = () => {
  try {
    console.log('准备渲染内容, 类型:', typeof props.content, '值:', props.content?.substring?.(0, 100));
    let contentToRender = props.content;
    if (contentToRender === null || contentToRender === undefined) {
      contentToRender = '';
    } else if (typeof contentToRender !== 'string') {
      contentToRender = String(contentToRender);
    }
    
    // 检查内容中的图片URL
    const imageMatches = contentToRender.match(/!\[.*?\]\((.*?)\)/g);
    if (imageMatches && imageMatches.length > 0) {
      console.log(`内容中包含 ${imageMatches.length} 个图片链接:`);
      imageMatches.forEach((match, index) => {
        const urlMatch = match.match(/!\[.*?\]\((.*?)\)/);
        const url = urlMatch?.[1] || '';
        console.log(`图片 #${index + 1}: ${url}`);
        
        // 尝试预处理图片URL
        if (url) {
          const processedUrl = processImageUrl(url);
          if (url !== processedUrl) {
            // 替换内容中的图片URL
            contentToRender = contentToRender.replace(
              `![](${url})`, 
              `![](${processedUrl})`
            ).replace(
              `![${url}](${url})`, 
              `![${url}](${processedUrl})`
            );
            console.log(`URL已预处理: ${url} -> ${processedUrl}`);
          }
        }
      });
    } else {
      console.log('内容中未找到图片链接');
    }
    
    // 检查URL类型
    if (contentToRender.includes('blob:')) {
      console.log('内容包含blob URL');
    }
    if (contentToRender.includes('/api/attachments/')) {
      console.log('内容包含API附件URL');
    }
    
    renderedContent.value = renderMarkdown(contentToRender);
    
    // 在下一个tick处理页面上的图片
    nextTick(() => {
      try {
        applyPrismHighlight(); // 应用Prism高亮
        
        // 处理页面上的所有图片
        handleRenderedImages();
        
        // 添加代码块复制功能
        const copyButtons = contentRef.value.querySelectorAll('.code-copy-btn');
        copyButtons.forEach(btn => {
          btn.addEventListener('click', handleCopyClick);
        });
      } catch (error) {
        console.error('处理渲染后内容时出错:', error);
      }
    });
  } catch (error) {
    console.error('渲染内容时出错:', error);
    renderedContent.value = `<div class="text-red-500">渲染内容失败: ${error.message}</div>`;
  }
};

// 在页面渲染后处理所有图片
const handleRenderedImages = () => {
  // 查找所有图片元素
  const images = contentRef.value.querySelectorAll('img');
  console.log(`处理已渲染的图片，共 ${images.length} 个`);
  
  // 为每个图片添加错误处理
  images.forEach((img, index) => {
    const src = img.getAttribute('src');
    console.log(`检查图片 #${index+1}，src:`, src);
    
    // 添加加载状态的CSS类
    img.classList.add('loading-image');
    
    // 为图片添加加载成功的监听器
    img.onload = function() {
      this.classList.remove('loading-image');
      this.classList.add('loaded-image');
      console.log(`图片 #${index+1} 加载成功:`, src);
      
      // 如果是blob URL，记录下来以便调试
      if (src.startsWith('blob:')) {
        this.setAttribute('data-loaded-blob', src);
        console.log(`Blob URL 图片加载成功:`, src);
      }
    };
    
    // 如果是blob URL，添加高级错误处理
    if (src && src.startsWith('blob:')) {
      console.log(`为blob URL图片 #${index+1} 添加错误处理`);
      
      // 保存原始blob URL以便调试
      img.setAttribute('data-original-blob', src);
      
      // 提取blob URL的特征信息
      const blobId = extractInfoFromBlobUrl(src);
      if (blobId) {
        img.setAttribute('data-blob-id', blobId);
      }
      
      // 添加错误处理函数
      img.onerror = async function() {
        console.error(`Blob URL图片加载失败: ${src}`);
        this.classList.remove('loading-image');
        this.classList.add('error-image');
        
        // 设置一个基本的错误指示器图片
        this.src = 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9ImN1cnJlbnRDb2xvciIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiIGNsYXNzPSJsdWNpZGUgbHVjaWRlLWFsZXJ0LXRyaWFuZ2xlIj48cGF0aCBkPSJtMjEuNzMgMTgtOC0xNGEyIDIgMCAwIDAtMy40NiAwbC04IDE0QTIgMiAwIDAgMCA0IDE5LjVoMTZhMiAyIDAgMCAwIDEuNzMtMS41WiIvPjxsaW5lIHgxPSIxMiIgeTE9IjgiIHgyPSIxMiIgeTI9IjEyIi8+PGxpbmUgeDE9IjEyIiB5MT0iMTYiIHgyPSIxMi4wMSIgeTI9IjE2Ii8+PC9zdmc+';
        
        // 尝试恢复策略 1: 通过alt属性查找附件ID
        const alt = this.getAttribute('alt');
        let attachmentId = null;
        let apiUrl = null;
        
        // 如果alt包含数字，可能是附件ID
        if (alt) {
          const numberMatch = alt.match(/(\d+)/);
          if (numberMatch && numberMatch[1]) {
            attachmentId = numberMatch[1];
            apiUrl = `/api/attachments/${attachmentId}`;
            console.log(`从alt属性中发现可能的附件ID: ${attachmentId}`);
          } else {
            // 尝试根据alt作为文件名查找附件
            try {
              attachmentId = await fetchAttachmentsByFilename(alt);
              if (attachmentId) {
                apiUrl = `/api/attachments/${attachmentId}`;
                console.log(`通过文件名 "${alt}" 找到附件ID: ${attachmentId}`);
              }
            } catch (error) {
              console.error('通过文件名查找附件失败:', error);
            }
          }
        }
        
        // 尝试恢复策略 2: 从blob URL特征中查找附件ID
        if (!apiUrl && blobId) {
          // 尝试从blob URL特征中查找附件ID
          const cachedUrl = getBlobMapping(blobId);
          if (cachedUrl) {
            apiUrl = cachedUrl;
            console.log(`从缓存的blob映射中找到匹配的API URL: ${apiUrl}`);
          }
        }
        
        // 如果找到了API URL，尝试重新加载
        if (apiUrl) {
          console.log(`尝试从API URL重新加载图片: ${apiUrl}`);
          
          // 显示加载指示器
          this.classList.remove('error-image');
          this.classList.add('loading-image');
          
          // 尝试使用认证方式加载
          loadAuthenticatedImage(apiUrl)
            .then(blobUrl => {
              if (blobUrl) {
                console.log(`成功加载替代图片: ${blobUrl}`);
                this.src = blobUrl;
                this.classList.remove('loading-image');
                this.classList.remove('error-image');
                this.classList.add('loaded-image');
                
                // 记录成功的blob URL映射，以便下次使用
                if (blobId && apiUrl) {
                  saveBlobMapping(blobId, apiUrl);
                }
              }
            })
            .catch(error => {
              console.error('重新加载图片失败:', error);
              this.classList.remove('loading-image');
              this.classList.add('error-image');
              showImageError(this);
            });
        } else {
          // 如果找不到附件ID，显示错误指示器
          showImageError(this);
        }
      };
      
      // 检查图片是否已经加载失败
      if (img.complete && (img.naturalWidth === 0 || img.naturalHeight === 0)) {
        console.log(`图片 #${index+1} 已失败，触发错误处理`);
        img.onerror();
      }
    } 
    // 如果是api附件URL，确保正确加载
    else if (src && src.includes('/api/attachments/')) {
      console.log(`为API附件图片 #${index+1} 添加认证加载处理`);
      
      // 提取附件ID
      const idMatch = src.match(/\/api\/attachments\/(\d+)/);
      if (idMatch && idMatch[1]) {
        const attachmentId = idMatch[1];
        console.log(`从URL提取的附件ID: ${attachmentId}`);
        
        // 加载认证图片
        loadAuthenticatedImage(src)
          .then(blobUrl => {
            if (blobUrl) {
              console.log(`成功加载认证图片: ${blobUrl}`);
              img.src = blobUrl;
              img.classList.remove('loading-image');
              img.classList.add('loaded-image');
              
              // 存储blob URL和原始URL的映射关系
              const blobId = extractInfoFromBlobUrl(blobUrl);
              if (blobId && src) {
                saveBlobMapping(blobId, src);
              }
            }
          })
          .catch(error => {
            console.error('加载认证图片失败:', error);
            img.classList.remove('loading-image');
            img.classList.add('error-image');
            showImageError(img);
          });
      }
    }
  });
};

// 显示图片错误指示器
const showImageError = (img) => {
  // 创建容器
  const parent = img.parentNode;
  if (!parent) return;
  
  parent.classList.add('image-error-container');
  
  // 检查是否已存在错误指示器
  if (parent.querySelector('.image-error-indicator')) return;
  
  // 添加错误标记
  const errorIndicator = document.createElement('div');
  errorIndicator.className = 'image-error-indicator';
  errorIndicator.innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
    </svg>
    <span class="text-sm text-red-500">图片加载失败</span>
    <span class="text-xs text-gray-500 mt-1">刷新后的blob链接无效</span>
  `;
  parent.appendChild(errorIndicator);
};

// 通过文件名查找附件ID
const fetchAttachmentsByFilename = async (filename) => {
  try {
    // 尝试从全局状态或者API获取附件列表
    // 这里的实现需要根据您的应用逻辑调整
    
    // 方法1: 尝试从window上的附件store获取
    if (window.$attachmentsStore) {
      const attachments = window.$attachmentsStore.attachments;
      if (attachments && attachments.length) {
        const attachment = attachments.find(a => a.filename === filename);
        if (attachment) {
          return attachment.id;
        }
      }
    }
    
    // 方法2: 尝试使用API查询
    // 这里需要根据您的API调整
    const token = localStorage.getItem('token');
    if (token) {
      const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
      const response = await fetch(`${baseUrl}/attachments/search?filename=${encodeURIComponent(filename)}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        if (data && data.data && data.data.length > 0) {
          return data.data[0].id;
        }
      }
    }
    
    // 如果找不到，尝试从文件名中提取数字作为ID
    const matches = filename.match(/(\d+)/);
    if (matches && matches[1]) {
      return matches[1];
    }
    
    return null;
  } catch (error) {
    console.error('查找附件失败:', error);
    return null;
  }
};

// 监听内容变化
watch(() => props.content, (newVal) => {
  updateContent();
});

// 组件加载时渲染初始内容
onMounted(() => {
  console.log('MarkdownRenderer组件挂载，marked版本:', marked.version || '未知');
  console.log('初始内容，类型:', typeof props.content, '长度:', props.content?.length || 0);
  
  // 检查marked配置
  const markedDefaults = marked.getDefaults ? marked.getDefaults() : 'getDefaults方法不可用';
  console.log('Marked默认配置:', markedDefaults);
  
  // 尝试清理过期的blob映射
  cleanupBlobMappings();
  
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
  transition: all 0.3s ease;
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

/* 图片容器样式 */
.image-container {
  position: relative;
  display: inline-block;
  max-width: 100%;
  margin: 0.5rem 0;
  border-radius: 0.375rem;
  overflow: hidden;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease-in-out;
}

.image-container:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* 图片样式 */
.markdown-image {
  display: block;
  max-width: 100%;
  height: auto;
  transition: filter 0.3s ease-in-out;
}

/* 加载状态 */
.image-loading {
  position: relative;
  min-height: 100px;
  background-color: rgba(229, 231, 235, 0.5);
}

.dark .image-loading {
  background-color: rgba(55, 65, 81, 0.3);
}

.image-loading .markdown-image {
  opacity: 0.5;
  filter: blur(3px);
}

.image-loading-indicator {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 30px;
  height: 30px;
  border: 3px solid rgba(99, 102, 241, 0.2);
  border-radius: 50%;
  border-top-color: rgba(99, 102, 241, 0.8);
  animation: spin 1s linear infinite;
}

/* 错误状态 */
.image-error-container {
  position: relative;
  min-height: 100px;
  background-color: rgba(254, 226, 226, 0.5);
  border: 1px dashed rgba(220, 38, 38, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.dark .image-error-container {
  background-color: rgba(127, 29, 29, 0.2);
}

.image-error {
  opacity: 0.1;
  filter: grayscale(100%);
}

.image-error-indicator {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  text-align: center;
  gap: 0.5rem;
}

.dark .image-error-indicator {
  background-color: #374151;
  color: #e5e7eb;
}

/* 代码块样式 */
.code-block-wrapper {
  position: relative;
  margin: 1rem 0;
  border-radius: 0.375rem;
  overflow: hidden;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  background-color: rgba(243, 244, 246, 0.8);
  border-bottom: 1px solid rgba(229, 231, 235, 0.8);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.75rem;
}

.dark .code-header {
  background-color: rgba(31, 41, 55, 0.8);
  border-bottom: 1px solid rgba(55, 65, 81, 0.8);
}

.code-language {
  font-weight: 500;
  color: #6366F1;
  text-transform: uppercase;
}

.dark .code-language {
  color: #A5B4FC;
}

.code-copy-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  background-color: rgba(255, 255, 255, 0.5);
  border-radius: 0.25rem;
  font-size: 0.75rem;
  color: #4B5563;
  transition: all 0.2s ease-in-out;
}

.dark .code-copy-btn {
  background-color: rgba(55, 65, 81, 0.5);
  color: #E5E7EB;
}

.code-copy-btn:hover {
  background-color: rgba(255, 255, 255, 0.8);
  color: #111827;
}

.dark .code-copy-btn:hover {
  background-color: rgba(75, 85, 99, 0.8);
  color: #F9FAFB;
}

.code-block {
  margin: 0;
  padding: 1rem;
  overflow-x: auto;
  background-color: rgba(249, 250, 251, 0.8);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.875rem;
  line-height: 1.5;
}

.dark .code-block {
  background-color: rgba(17, 24, 39, 0.8);
}

/* 其他样式增强 */
.markdown-renderer h1,
.markdown-renderer h2,
.markdown-renderer h3,
.markdown-renderer h4,
.markdown-renderer h5,
.markdown-renderer h6 {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
  line-height: 1.25;
  color: #111827;
}

.dark .markdown-renderer h1,
.dark .markdown-renderer h2,
.dark .markdown-renderer h3,
.dark .markdown-renderer h4,
.dark .markdown-renderer h5,
.dark .markdown-renderer h6 {
  color: #F9FAFB;
}

.markdown-renderer a {
  color: #4F46E5;
  text-decoration: none;
  transition: color 0.2s ease-in-out;
}

.dark .markdown-renderer a {
  color: #818CF8;
}

.markdown-renderer a:hover {
  text-decoration: underline;
  color: #6366F1;
}

.dark .markdown-renderer a:hover {
  color: #A5B4FC;
}

.markdown-renderer ul,
.markdown-renderer ol {
  padding-left: 1.5rem;
  margin: 1rem 0;
}

.markdown-renderer blockquote {
  border-left: 4px solid #E5E7EB;
  padding-left: 1rem;
  margin: 1rem 0;
  font-style: italic;
  color: #4B5563;
}

.dark .markdown-renderer blockquote {
  border-left-color: #4B5563;
  color: #9CA3AF;
}

.markdown-renderer hr {
  border: 0;
  border-top: 1px solid #E5E7EB;
  margin: 1.5rem 0;
}

.dark .markdown-renderer hr {
  border-top-color: #374151;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 复制成功提示 */
.copy-success-tooltip {
  position: fixed;
  bottom: 1rem;
  right: 1rem;
  padding: 0.5rem 1rem;
  background-color: #10B981;
  color: white;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  z-index: 50;
  animation: fadeInOut 2s ease-in-out;
}

@keyframes fadeInOut {
  0%, 100% { opacity: 0; transform: translateY(10px); }
  10%, 90% { opacity: 1; transform: translateY(0); }
}
</style> 