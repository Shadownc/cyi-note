import { get, post, put, del } from '@/utils/request'
import { normalizeNoteDates } from '@/utils/dateUtils'

// 标准化API响应中的笔记列表
function normalizeNotes(response) {
  if (!response || !response.notes) return response;
  
  response.notes = response.notes.map(note => normalizeNoteDates(note));
  return response;
}

// 获取笔记列表，支持分页和标签过滤
export function getNotes(params) {
  return get('/notes', params).then(normalizeNotes);
}

// 获取公开笔记列表
export function getPublicNotes(params) {
  return get('/public/notes', params).then(normalizeNotes);
}

// 创建笔记
export function createNote(data) {
  return post('/notes', data).then(normalizeNoteDates);
}

// 获取指定笔记
export function getNote(id) {
  return get(`/notes/${id}`).then(normalizeNoteDates);
}

// 更新笔记
export function updateNote(id, data) {
  return put(`/notes/${id}`, data).then(normalizeNoteDates);
}

// 删除笔记
export function deleteNote(id) {
  return del(`/notes/${id}`);
}

// 搜索笔记
export function searchNotes(keyword, page = 1, pageSize = 10, sortBy = null) {
  const params = { keyword, page, pageSize };
  if (sortBy) {
    params.sortBy = sortBy;
  }
  return get('/notes/search', params).then(normalizeNotes);
} 