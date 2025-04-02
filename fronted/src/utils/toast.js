/**
 * 简单的Toast通知工具
 */

// 创建并显示toast元素
function createToast(message, type = 'info', duration = 3000) {
  // 检查toast容器是否存在，不存在则创建
  let toastContainer = document.getElementById('toast-container');
  if (!toastContainer) {
    toastContainer = document.createElement('div');
    toastContainer.id = 'toast-container';
    toastContainer.style.position = 'fixed';
    toastContainer.style.bottom = '20px';
    toastContainer.style.right = '20px';
    toastContainer.style.zIndex = '10000';
    document.body.appendChild(toastContainer);
  }
  
  // 创建toast元素
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.textContent = message;
  
  // 设置样式
  toast.style.minWidth = '250px';
  toast.style.margin = '5px 0';
  toast.style.padding = '10px 15px';
  toast.style.borderRadius = '4px';
  toast.style.boxShadow = '0 2px 10px rgba(0, 0, 0, 0.2)';
  toast.style.fontFamily = 'Arial, sans-serif';
  toast.style.fontSize = '14px';
  toast.style.opacity = '0';
  toast.style.transition = 'opacity 0.3s, transform 0.3s';
  toast.style.transform = 'translateY(20px)';
  
  // 根据类型设置不同的颜色
  switch (type) {
    case 'success':
      toast.style.backgroundColor = '#4caf50';
      toast.style.color = 'white';
      break;
    case 'error':
      toast.style.backgroundColor = '#f44336';
      toast.style.color = 'white';
      break;
    case 'warning':
      toast.style.backgroundColor = '#ff9800';
      toast.style.color = 'white';
      break;
    default:
      toast.style.backgroundColor = '#2196f3';
      toast.style.color = 'white';
  }
  
  // 添加到容器
  toastContainer.appendChild(toast);
  
  // 显示toast
  setTimeout(() => {
    toast.style.opacity = '1';
    toast.style.transform = 'translateY(0)';
  }, 10);
  
  // 设置自动关闭
  setTimeout(() => {
    // 淡出动画
    toast.style.opacity = '0';
    toast.style.transform = 'translateY(20px)';
    
    // 移除元素
    setTimeout(() => {
      toast.remove();
      
      // 如果没有更多toast，移除容器
      if (toastContainer.children.length === 0) {
        toastContainer.remove();
      }
    }, 300);
  }, duration);
  
  return toast;
}

export default {
  success(message, duration) {
    return createToast(message, 'success', duration);
  },
  
  error(message, duration) {
    return createToast(message, 'error', duration);
  },
  
  warning(message, duration) {
    return createToast(message, 'warning', duration);
  },
  
  info(message, duration) {
    return createToast(message, 'info', duration);
  }
}; 