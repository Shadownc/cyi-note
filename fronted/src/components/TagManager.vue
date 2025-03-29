<template>
  <div class="w-full">
    <div class="flex flex-wrap items-center gap-2 mb-4">
      <div class="relative">
        <input 
          v-model="newTagName" 
          type="text" 
          placeholder="新建标签..." 
          class="border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm dark:bg-gray-700 dark:text-white"
          @keyup.enter="addTag"
        >
        <button 
          @click="addTag"
          class="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-500 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>
    </div>

    <div v-if="tags.length === 0" class="text-gray-500 dark:text-gray-400 text-center my-8">
      暂无标签，请添加新标签
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
import { ref, defineProps, defineEmits } from 'vue';

const props = defineProps({
  tags: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['add-tag', 'remove-tag']);

const newTagName = ref('');

const addTag = () => {
  if (newTagName.value.trim()) {
    emit('add-tag', {
      name: newTagName.value.trim()
    });
    newTagName.value = '';
  }
};

const removeTag = (tagId) => {
  emit('remove-tag', tagId);
};
</script> 