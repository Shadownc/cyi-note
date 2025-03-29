import { get, post, del } from '@/utils/request'

// 获取所有标签
export function getTags() {
  return get('/tags')
}

// 创建标签
export function createTag(data) {
  return post('/tags', data)
}

// 删除标签
export function deleteTag(id) {
  return del(`/tags/${id}`)
}

// 给笔记添加标签
export function addTagToNote(tagId, noteId) {
  return post(`/tags/${tagId}/notes/${noteId}`)
}

// 从笔记中移除标签
export function removeTagFromNote(tagId, noteId) {
  return del(`/tags/${tagId}/notes/${noteId}`)
} 