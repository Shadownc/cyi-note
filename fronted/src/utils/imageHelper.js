import { useAttachmentsStore } from '@/stores/attachments';

/**
 * 获取API基础URL
 * @returns {string} API基础URL
 */
export function getApiBaseUrl() {
  // 如果使用开发服务器代理，则直接返回/api前缀
  return '/api';
  // 下面这行在使用代理时不需要，使用环境变量会导致URL变成完整路径，绕过代理
  // return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
}

/**
 * 构建附件URL
 * @param {string} attachmentId 附件ID
 * @returns {string} 完整的附件URL
 */
export function buildAttachmentUrl(attachmentId) {
  if (!attachmentId) return '';
  
  // 检查是否已经是完整URL
  if (attachmentId.startsWith('http') || attachmentId.startsWith('blob:')) {
    console.log('buildAttachmentUrl: 完整URL不做修改:', attachmentId);
    return attachmentId;
  }
  
  // 如果已经是以/api开头的完整路径
  if (attachmentId.startsWith('/api/attachments/')) {
    console.log('buildAttachmentUrl: 已有API路径不做修改:', attachmentId);
    return attachmentId;
  }
  
  // 如果是相对路径，添加API前缀
  if (!attachmentId.startsWith('/')) {
    // 如果已经包含"attachments/"前缀
    if (attachmentId.startsWith('attachments/')) {
      const url = `${getApiBaseUrl()}/${attachmentId}`;
      console.log(`buildAttachmentUrl: 添加API前缀 ${attachmentId} -> ${url}`);
      return url;
    }
    // 否则添加前缀
    const url = `${getApiBaseUrl()}/attachments/${attachmentId}`;
    console.log(`buildAttachmentUrl: 添加完整附件路径 ${attachmentId} -> ${url}`);
    return url;
  }
  
  // 如果以/开头但不是完整路径
  if (!attachmentId.startsWith('/api')) {
    const url = `${getApiBaseUrl()}${attachmentId}`;
    console.log(`buildAttachmentUrl: 添加API基础URL ${attachmentId} -> ${url}`);
    return url;
  }
  
  // 如果是完整API路径（/api/...）
  console.log('buildAttachmentUrl: API路径不做修改:', attachmentId);
  return attachmentId;
}

/**
 * 上传图片并获取URL
 * @param {File} file 图片文件
 * @param {string} noteId 笔记ID
 * @returns {Promise<{success: boolean, url: string, message: string}>} 上传结果
 */
export async function uploadImage(file, noteId) {
  if (!file || !noteId) {
    return { success: false, url: '', message: '缺少文件或笔记ID' };
  }
  
  // 验证文件类型
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
  if (!allowedTypes.includes(file.type)) {
    return { success: false, url: '', message: '只支持 JPG, PNG, GIF 或 WebP 格式的图片' };
  }
  
  // 获取附件存储
  const attachmentsStore = useAttachmentsStore();
  
  try {
    // 上传附件
    console.log('开始上传图片:', file.name);
    const response = await attachmentsStore.uploadAttachment(file, noteId);
    
    if (!response || !response.id) {
      throw new Error('上传响应中未包含有效的附件ID');
    }
    
    // 构建URL
    const imageUrl = buildAttachmentUrl(response.id);
    console.log('图片上传成功，URL:', imageUrl);
    
    return {
      success: true,
      url: imageUrl,
      id: response.id,
      message: '上传成功'
    };
  } catch (error) {
    console.error('图片上传失败:', error);
    return {
      success: false,
      url: '',
      message: error.message || '未知错误'
    };
  }
}

/**
 * 创建临时图片URL
 * @param {File} file 图片文件
 * @returns {string} 临时URL
 */
export function createTempImageUrl(file) {
  if (!file) return '';
  return URL.createObjectURL(file);
}

/**
 * 释放临时URL
 * @param {string} url 临时URL
 */
export function revokeTempUrl(url) {
  if (url && url.startsWith('blob:')) {
    URL.revokeObjectURL(url);
  }
}

/**
 * 加载需要认证的图片
 * @param {string} url 需要加载的URL
 * @returns {Promise<string>} 返回blob URL或原始URL
 */
export async function loadAuthenticatedImage(url) {
  if (!url || !url.startsWith('/api/attachments/')) {
    return url;
  }
  
  console.log('开始加载认证图片:', url);
  
  try {
    // 获取认证令牌
    const token = localStorage.getItem('token');
    if (!token) {
      console.warn('未找到认证令牌，可能无法加载图片:', url);
      return url;
    }
    
    // 创建新的fetch请求获取图片数据
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`图片加载失败: ${response.status} ${response.statusText}`);
    }
    
    const blob = await response.blob();
    // 创建blob URL
    const blobUrl = URL.createObjectURL(blob);
    console.log('已创建认证图片的blob URL:', blobUrl, '原URL:', url);
    return blobUrl;
  } catch (error) {
    console.error('加载认证图片失败:', error, url);
    return url;
  }
} 