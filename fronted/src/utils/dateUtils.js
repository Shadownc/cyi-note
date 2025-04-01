/**
 * 格式化日期为本地格式
 * @param {string|Date} dateInput - 日期字符串或日期对象
 * @param {boolean} showTime - 是否显示时间部分
 * @returns {string} 格式化后的日期字符串
 */
export function formatDate(dateInput, showTime = false) {
  if (!dateInput) {
    console.warn('日期输入为空');
    return '无日期';
  }
  
  // 尝试修复日期格式
  let fixedDate = dateInput;
  if (typeof dateInput === 'string') {
    if (dateInput.includes('T') && !dateInput.includes('Z') && !dateInput.includes('+')) {
      fixedDate = dateInput + 'Z';
    }
  }
  
  const date = new Date(fixedDate);
  
  if (isNaN(date.getTime())) {
    console.warn('日期格式无效:', dateInput);
    return '日期无效';
  }
  
  try {
    // 使用标准的本地化格式
    if (showTime) {
      return date.toLocaleDateString('zh-CN', { 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    });
  } catch (e) {
    console.error('日期格式化错误，使用降级方案:', e);
    
    // 降级方案：手动格式化
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    
    if (showTime) {
      return `${month}月${day}日 ${hours}:${minutes}`;
    }
    return `${year}年${month}月${day}日`;
  }
}

/**
 * 标准化日期格式，确保前后端日期处理一致
 * @param {Object} note - 笔记对象
 * @returns {Object} 标准化后的笔记对象
 */
export function normalizeNoteDates(note) {
  if (!note) return note;
  
  // 转换日期字段格式
  if (note.created_at && !note.createdAt) {
    note.createdAt = note.created_at;
  } else if (note.createdAt && !note.created_at) {
    note.created_at = note.createdAt;
  }
  
  if (note.updated_at && !note.updatedAt) {
    note.updatedAt = note.updated_at;
  } else if (note.updatedAt && !note.updated_at) {
    note.updated_at = note.updatedAt;
  }
  
  return note;
} 