<template>
  <div class="flex flex-col lg:flex-row gap-4">
    <!-- 编辑器部分 -->
    <div class="lg:w-1/2 w-full">
      <div class="mb-2 flex justify-between items-center">
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200">编辑</h3>
        <div class="flex flex-wrap gap-1">
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('# ')">H1</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('## ')">H2</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('### ')">H3</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('**粗体**')">粗体</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('*斜体*')">斜体</button>
          <button class="btn btn-secondary text-xs md:text-sm dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600" @click="insertMd('[链接](https://example.com)')">链接</button>
        </div>
      </div>
      <textarea 
        v-model="content" 
        class="w-full border border-gray-300 dark:border-gray-600 rounded-md p-3 min-h-[300px] h-full resize-none focus:outline-none focus:ring-2 focus:ring-primary-500 font-mono dark:bg-gray-800 dark:text-gray-100"
        placeholder="在这里输入Markdown内容..."
        @input="updateContent"
      ></textarea>
    </div>
    
    <!-- 预览部分 -->
    <div class="lg:w-1/2 w-full">
      <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">预览</h3>
      <div class="border border-gray-300 dark:border-gray-600 rounded-md p-3 min-h-[300px] prose dark:prose-invert max-w-none overflow-auto dark:bg-gray-800 dark:bg-opacity-60" v-html="renderedContent"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, onMounted, watch } from 'vue';
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

// 将Markdown渲染为HTML
const renderMarkdown = (markdown) => {
  return marked(markdown);
};

// 更新内容并通知父组件
const updateContent = () => {
  renderedContent.value = renderMarkdown(content.value);
  emit('update:content', content.value);
};

// 在特定位置插入Markdown语法
const insertMd = (syntax) => {
  const textarea = document.querySelector('textarea');
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
</style> 