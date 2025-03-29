import { post } from '@/utils/request'

// 生成标签推荐
export function generateTags(content) {
  return post('/ai/tags', { content })
}

// 生成内容摘要
export function generateSummary(content, noteId = null) {
  const params = { content }
  if (noteId) {
    params.note_id = noteId
  }
  return post('/ai/summary', params)
} 