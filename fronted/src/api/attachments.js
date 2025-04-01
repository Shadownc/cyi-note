import { post, get, del } from '@/utils/request'
import axios from 'axios'

// API 基础URL
const getBaseUrl = () => {
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
};

// 上传附件
export function uploadAttachment(file, noteId) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('note_id', noteId)
  // 添加原始文件名
  formData.append('filename', file.name)
  
  console.log('正在上传附件:', file.name, '到笔记ID:', noteId);
  
  return post('/attachments', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    timeout: 60000 // 上传文件可能需要更长的超时时间
  })
}

// 获取附件列表
export function getAttachments(noteId) {
  return get(`/notes/${noteId}/attachments`)
}

// 获取附件
export function getAttachment(id) {
  return get(`/attachments/${id}`)
}

// 删除附件
export function deleteAttachment(id) {
  return del(`/attachments/${id}`)
}

// 获取附件下载链接
export function getAttachmentUrl(id) {
  if (!id) return '';
  // 从环境变量获取API基础URL, 或使用默认值
  const baseUrl = getBaseUrl();
  // 使用API路径直接访问附件
  const url = `${baseUrl}/attachments/${id}`;
  console.log('生成附件URL:', url);
  return url;
}

// 使用认证下载附件
export async function downloadWithAuth(id, filename) {
  if (!id) return false;
  
  const token = localStorage.getItem('token');
  if (!token) {
    console.error('未提供认证令牌，无法下载');
    throw new Error('未提供认证令牌');
  }

  try {
    // 使用axios直接获取文件，并带上认证头
    const response = await axios.get(`${getBaseUrl()}/attachments/${id}`, {
      responseType: 'blob',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename || `attachment-${id}`);
    document.body.appendChild(link);
    link.click();
    
    // 清理
    setTimeout(() => {
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    }, 100);
    
    return true;
  } catch (error) {
    console.error('下载附件失败:', error);
    throw error;
  }
}

// 获取资源库文件（按日期分组）
export function getAttachmentsLibrary(params) {
  return get('/attachments/library', params)
    .then(response => {
      // 检查响应是否是标准格式（带有success字段的包装响应）
      if (response && typeof response.success !== 'undefined') {
        console.log('获取到标准响应结构:', response);
        return response.data; // 返回data字段中的实际数据
      }
      
      // 非标准格式，直接返回
      console.log('获取到非标准响应结构:', response);
      return response;
    })
    .catch(error => {
      console.error('资源库API调用出错:', error);
      throw error;
    });
} 