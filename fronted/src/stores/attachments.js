import { ref } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api'

export const useAttachmentsStore = defineStore('attachments', () => {
  // 状态
  const attachments = ref([])
  const libraryAttachments = ref([]) // 资源库附件
  const dateGroups = ref([]) // 日期分组
  const totalAttachments = ref(0) // 总附件数
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
  
  // 下载附件（带认证）
  async function downloadAttachment(id, filename) {
    if (!id) return false
    try {
      return await api.attachments.downloadWithAuth(id, filename)
    } catch (error) {
      console.error(`下载附件 ${id} 失败:`, error)
      throw error
    }
  }

  // 更新上传进度
  function updateUploadProgress(progress) {
    uploadProgress.value = progress
  }

  // 上传临时附件（不关联到笔记）
  async function uploadTempAttachment(file) {
    if (!file) return null;
    
    isLoading.value = true;
    uploadProgress.value = 0;
    
    try {
      const formData = new FormData();
      formData.append('file', file);
      
      // 获取认证令牌
      const token = localStorage.getItem('token');
      const headers = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      
      // 使用相对路径和认证头
      const response = await fetch('/api/attachments/temp', {
        method: 'POST',
        headers: headers,
        body: formData
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `上传失败: ${response.status} ${response.statusText}`);
      }
      
      const responseData = await response.json();
      console.log('临时附件上传响应:', responseData);
      
      // 获取响应中的数据
      // backend返回的数据结构为：{ data: { id, url, filename, filetype, filesize, success }, message, status }
      // 或者是：{ id, url, filename, filetype, filesize, success } (取决于后端实现)
      const data = responseData.data || responseData;
      
      // 确认返回数据包含成功标志
      if (!data.success) {
        throw new Error(responseData.error || data.error || '服务器未返回成功状态');
      }
      
      // 返回处理后的结果
      return {
        id: data.id,
        url: data.url || '', // 确保url存在，即使为空字符串
        filename: data.filename,
        success: true
      };
    } catch (error) {
      console.error('上传临时附件失败:', error);
      throw error;
    } finally {
      isLoading.value = false;
      uploadProgress.value = 0;
    }
  }
  
  // 关联临时附件到笔记
  async function associateTempAttachment(tempId, noteId) {
    if (!tempId || !noteId) return false
    
    isLoading.value = true
    try {
      // 获取认证令牌
      const token = localStorage.getItem('token')
      const headers = {
        'Content-Type': 'application/json'
      }
      if (token) {
        headers['Authorization'] = `Bearer ${token}`
      }
      
      // 使用相对路径访问API
      const response = await fetch(`/api/attachments/temp/${tempId}/associate`, {
        method: 'POST',
        headers: headers,
        body: JSON.stringify({ noteId })
      })
      
      if (!response.ok) {
        throw new Error(`关联失败: ${response.status} ${response.statusText}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error(`关联临时附件 ${tempId} 到笔记 ${noteId} 失败:`, error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 获取资源库文件（按日期分组）
  async function fetchLibraryAttachments(page = 1, pageSize = 20, fileType = '') {
    isLoading.value = true
    try {
      console.log('请求资源库数据，参数:', { page, pageSize, filetype: fileType });
      
      const response = await api.attachments.getAttachmentsLibrary({
        page,
        pageSize,
        filetype: fileType
      });
      
      console.log('资源库数据返回:', response);
      
      // 确保 dateGroups 和 attachments 数据存在并且是数组
      libraryAttachments.value = Array.isArray(response.attachments) ? response.attachments : [];
      dateGroups.value = Array.isArray(response.dateGroups) ? response.dateGroups : [];
      totalAttachments.value = response.total || 0;
      
      // 检查日期组是否为空
      if (dateGroups.value.length === 0) {
        console.warn('未获取到日期分组数据');
      } else {
        console.log('日期分组数据:', dateGroups.value);
      }
      
      // 检查附件数据是否为空
      if (libraryAttachments.value.length === 0) {
        console.warn('未获取到附件数据');
      } else {
        console.log('附件数据:', libraryAttachments.value);
      }
      
      return response;
    } catch (error) {
      console.error('获取资源库文件失败:', error);
      return { attachments: [], dateGroups: [], total: 0 };
    } finally {
      isLoading.value = false;
    }
  }

  return {
    attachments,
    libraryAttachments,
    dateGroups,
    totalAttachments,
    isLoading,
    uploadProgress,
    fetchAttachments,
    uploadAttachment,
    deleteAttachment,
    getAttachmentUrl,
    downloadAttachment,
    updateUploadProgress,
    fetchLibraryAttachments,
    uploadTempAttachment,
    associateTempAttachment
  }
}) 