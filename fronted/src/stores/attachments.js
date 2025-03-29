import { ref } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api'

export const useAttachmentsStore = defineStore('attachments', () => {
  // 状态
  const attachments = ref([])
  const isLoading = ref(false)
  const uploadProgress = ref(0)

  // 获取笔记的附件
  async function fetchAttachments(noteId) {
    if (!noteId) return []
    
    isLoading.value = true
    try {
      const response = await api.attachments.getAttachments(noteId)
      attachments.value = response
      return response
    } catch (error) {
      console.error(`获取笔记 ${noteId} 的附件失败:`, error)
      return []
    } finally {
      isLoading.value = false
    }
  }

  // 上传附件
  async function uploadAttachment(file, noteId) {
    if (!file || !noteId) return null
    
    isLoading.value = true
    uploadProgress.value = 0
    
    try {
      const newAttachment = await api.attachments.uploadAttachment(file, noteId)
      attachments.value.push(newAttachment)
      return newAttachment
    } catch (error) {
      console.error('上传附件失败:', error)
      throw error
    } finally {
      isLoading.value = false
      uploadProgress.value = 0
    }
  }

  // 删除附件
  async function deleteAttachment(id) {
    if (!id) return false
    
    isLoading.value = true
    try {
      await api.attachments.deleteAttachment(id)
      attachments.value = attachments.value.filter(attachment => attachment.id !== id)
      return true
    } catch (error) {
      console.error(`删除附件 ${id} 失败:`, error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 获取附件下载链接
  function getAttachmentUrl(id) {
    if (!id) return ''
    return api.attachments.getAttachmentUrl(id)
  }

  // 更新上传进度
  function updateUploadProgress(progress) {
    uploadProgress.value = progress
  }

  return {
    attachments,
    isLoading,
    uploadProgress,
    fetchAttachments,
    uploadAttachment,
    deleteAttachment,
    getAttachmentUrl,
    updateUploadProgress
  }
}) 