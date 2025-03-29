<template>
  <div class="w-full">
    <div class="flex flex-wrap items-center gap-2 mb-4">
      <div class="relative w-full sm:w-64">
        <input 
          v-model="tagSearchText" 
          @input="debouncedSearch"
          type="text" 
          placeholder="搜索或添加标签..." 
          class="w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm dark:bg-gray-700 dark:text-white"
          @keydown.enter="addTagFromInput"
          @focus="showSuggestions = true"
          @blur="setTimeout(() => showSuggestions = false, 200)"
        >
        <button 
          @click="addTagFromInput"
          class="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-500 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
        
        <!-- 标签建议下拉框 -->
        <div v-if="showSuggestions && filteredAvailableTags.length > 0" 
             class="absolute z-10 mt-1 w-full bg-white dark:bg-gray-800 shadow-lg rounded-md border border-gray-200 dark:border-gray-700 max-h-48 overflow-y-auto">
          <div v-for="tag in filteredAvailableTags" 
               :key="tag.id"
               class="px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer text-sm"
               @mousedown="selectTag(tag)">
            {{ tag.name }}
          </div>
          <div v-if="tagSearchText && !tagExists && filteredAvailableTags.length > 0" 
               class="px-3 py-2 border-t border-gray-200 dark:border-gray-700 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer text-sm text-primary-600 dark:text-primary-400"
               @mousedown="createNewTag">
            + 创建新标签 "{{ tagSearchText }}"
          </div>
        </div>
        
        <div v-else-if="showSuggestions && tagSearchText && !tagExists"
             class="absolute z-10 mt-1 w-full bg-white dark:bg-gray-800 shadow-lg rounded-md border border-gray-200 dark:border-gray-700">
          <div class="px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer text-sm text-primary-600 dark:text-primary-400"
               @mousedown="createNewTag">
            + 创建新标签 "{{ tagSearchText }}"
          </div>
        </div>
      </div>
      
      <button 
        @click="$emit('generate-tags')"
        v-if="$attrs['ongenerate-tags']"
        class="text-sm py-1 px-3 bg-primary-50 dark:bg-primary-900 text-primary-600 dark:text-primary-300 border border-primary-200 dark:border-primary-700 rounded-md hover:bg-primary-100 dark:hover:bg-primary-800 transition-colors"
      >
        AI生成标签
      </button>
    </div>

    <div v-if="tags.length === 0" class="text-gray-500 dark:text-gray-400 text-center my-4">
      暂无标签，请添加或选择标签
    </div>

    <div class="flex flex-wrap gap-2">
      <div 
        v-for="tag in tags" 
        :key="tag.id"
        class="inline-flex items-center bg-primary-100 dark:bg-primary-700 dark:bg-opacity-50 text-primary-800 dark:text-primary-100 rounded-full px-3 py-1 text-sm group"
      >
        <span class="mr-1">{{ tag.name }}</span>
        <button 
          @click="removeTag(tag.id)" 
          class="text-primary-500 dark:text-primary-300 hover:text-red-500 dark:hover:text-red-400 opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, computed, inject, watch } from 'vue';
import { useTagsStore } from '@/stores/tags';

const props = defineProps({
  tags: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['add-tag', 'remove-tag', 'generate-tags']);

// 获取所有可用标签
const tagsStore = useTagsStore();
const availableTags = computed(() => tagsStore.tags || []);

const tagSearchText = ref('');
const showSuggestions = ref(false);

// 过滤已添加的标签，避免重复添加
const filteredAvailableTags = computed(() => {
  const currentTagIds = props.tags.map(tag => tag.id);
  let filtered = availableTags.value.filter(tag => !currentTagIds.includes(tag.id));
  
  // 如果有搜索文本，按名称过滤
  if (tagSearchText.value) {
    const search = tagSearchText.value.toLowerCase();
    filtered = filtered.filter(tag => tag.name.toLowerCase().includes(search));
  }
  
  return filtered;
});

// 检查当前输入的标签名是否已存在
const tagExists = computed(() => {
  if (!tagSearchText.value) return false;
  const search = tagSearchText.value.toLowerCase();
  return availableTags.value.some(tag => tag.name.toLowerCase() === search);
});

// 防抖函数
let searchTimeout;
const debouncedSearch = () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    showSuggestions.value = !!tagSearchText.value;
  }, 300);
};

// 从已有标签中选择一个
const selectTag = (tag) => {
  emit('add-tag', tag);
  tagSearchText.value = '';
  showSuggestions.value = false;
};

// 创建新标签
const createNewTag = async () => {
  if (!tagSearchText.value.trim()) return;
  
  try {
    const newTag = await tagsStore.createTag({ name: tagSearchText.value.trim() });
    emit('add-tag', newTag);
    tagSearchText.value = '';
    showSuggestions.value = false;
  } catch (error) {
    console.error('创建标签失败:', error);
  }
};

// 从输入框添加标签
const addTagFromInput = async () => {
  if (!tagSearchText.value.trim()) return;
  
  // 查找是否已有同名标签
  const existingTag = availableTags.value.find(tag => 
    tag.name.toLowerCase() === tagSearchText.value.trim().toLowerCase()
  );
  
  if (existingTag) {
    // 如果已有同名标签，直接选择
    selectTag(existingTag);
  } else {
    // 否则创建新标签
    await createNewTag();
  }
};

const removeTag = (tagId) => {
  emit('remove-tag', tagId);
};

// 在组件挂载时加载标签
if (availableTags.value.length === 0) {
  tagsStore.fetchTags();
}
</script> 