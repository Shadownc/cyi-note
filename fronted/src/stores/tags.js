import { ref } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api'

export const useTagsStore = defineStore('tags', () => {
  // 状态
  const tags = ref([])
  const isLoading = ref(false)

  // 获取所有标签
  async function fetchTags() {
    isLoading.value = true
    try {
      tags.value = await api.tags.getTags()
    } catch (error) {
      console.error('获取标签失败:', error)
    } finally {
      isLoading.value = false
    }
  }

  // 创建标签
  async function createTag(tagData) {
    isLoading.value = true
    try {
      const newTag = await api.tags.createTag(tagData)
      tags.value.push(newTag)
      return newTag
    } catch (error) {
      console.error('创建标签失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 删除标签
  async function deleteTag(id) {
    isLoading.value = true
    try {
      await api.tags.deleteTag(id)
      tags.value = tags.value.filter(tag => tag.id !== id)
      return true
    } catch (error) {
      console.error(`删除标签 ${id} 失败:`, error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 向笔记添加标签
  async function addTagToNote(tagId, noteId) {
    isLoading.value = true
    try {
      await api.tags.addTagToNote(tagId, noteId)
      return true
    } catch (error) {
      console.error(`向笔记 ${noteId} 添加标签 ${tagId} 失败:`, error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 从笔记移除标签
  async function removeTagFromNote(tagId, noteId) {
    isLoading.value = true
    try {
      await api.tags.removeTagFromNote(tagId, noteId)
      return true
    } catch (error) {
      console.error(`从笔记 ${noteId} 移除标签 ${tagId} 失败:`, error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 通过AI生成标签建议
  async function generateTagSuggestions(content) {
    isLoading.value = true
    try {
      return await api.ai.generateTags(content)
    } catch (error) {
      console.error('生成标签建议失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  return {
    tags,
    isLoading,
    fetchTags,
    createTag,
    deleteTag,
    addTagToNote,
    removeTagFromNote,
    generateTagSuggestions
  }
}) 