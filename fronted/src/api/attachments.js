import { post, get, del } from '@/utils/request'
import axios from 'axios'

// 上传附件
export function uploadAttachment(file, noteId) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('note_id', noteId)
  
  return post('/attachments', formData)
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
  return `${axios.defaults.baseURL}/attachments/${id}`
} 