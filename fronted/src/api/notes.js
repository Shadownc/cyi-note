import { get, post, put, del } from '@/utils/request'

// 获取笔记列表，支持分页和标签过滤
export function getNotes(params) {
  return get('/notes', params)
}

// 创建笔记
export function createNote(data) {
  return post('/notes', data)
}

// 获取指定笔记
export function getNote(id) {
  return get(`/notes/${id}`)
}

// 更新笔记
export function updateNote(id, data) {
  return put(`/notes/${id}`, data)
}

// 删除笔记
export function deleteNote(id) {
  return del(`/notes/${id}`)
}

// 搜索笔记
export function searchNotes(keyword, page = 1, pageSize = 10, sortBy = null) {
  const params = { keyword, page, pageSize };
  if (sortBy) {
    params.sortBy = sortBy;
  }
  return get('/notes/search', params);
} 