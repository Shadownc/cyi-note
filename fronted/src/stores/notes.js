import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api'

export const useNotesStore = defineStore('notes', () => {
  // 状态
  const notes = ref([])
  const publicNotes = ref([]) // 公开笔记
  const currentNote = ref(null)
  const isLoading = ref(false)
  const totalNotes = ref(0)
  const totalPublicNotes = ref(0) // 公开笔记总数
  const currentPage = ref(1)
  const publicNotesPage = ref(1) // 公开笔记分页
  const pageSize = ref(10)
  const activeTag = ref(null)
  const searchKeyword = ref('')

  // 计算属性
  const hasNotes = computed(() => notes.value.length > 0)
  const hasPublicNotes = computed(() => publicNotes.value.length > 0)
  const totalPages = computed(() => Math.ceil(totalNotes.value / pageSize.value))
  const totalPublicPages = computed(() => Math.ceil(totalPublicNotes.value / pageSize.value))

  // 获取用户自己的笔记列表
  async function fetchNotes(page = 1, size = 10, tagId = null, sortOption = null) {
    isLoading.value = true
    try {
      currentPage.value = page
      pageSize.value = size
      
      const params = {
        page,
        pageSize: size
      }
      
      // 添加排序参数
      if (sortOption) {
        params.sortBy = sortOption
      }
      
      if (tagId) {
        params.tagId = tagId
        activeTag.value = tagId
      } else if (searchKeyword.value) {
        const searchResult = await api.notes.searchNotes(
          searchKeyword.value, 
          page, 
          size,
          sortOption // 传递排序参数
        )
        notes.value = searchResult.notes
        totalNotes.value = searchResult.total
        return
      }
      
      const response = await api.notes.getNotes(params)
      notes.value = response.notes
      totalNotes.value = response.total
    } catch (error) {
      console.error('获取笔记失败:', error)
    } finally {
      isLoading.value = false
    }
  }

  // 获取公开笔记列表
  async function fetchPublicNotes(page = 1, size = 10) {
    isLoading.value = true
    try {
      publicNotesPage.value = page
      
      const params = {
        page,
        pageSize: size
      }
      
      const response = await api.notes.getPublicNotes(params)
      publicNotes.value = response.notes
      totalPublicNotes.value = response.total
    } catch (error) {
      console.error('获取公开笔记失败:', error)
    } finally {
      isLoading.value = false
    }
  }

  // 获取单个笔记
  async function fetchNote(id) {
    isLoading.value = true
    try {
      const response = await api.notes.getNote(id)
      console.log('获取笔记详情响应:', response)
      
      // 如果响应存在
      if (response) {
        // 响应已经过拦截器处理，直接使用
        // 注意：拦截器会从 {success:true, data:{笔记对象}} 中提取data字段返回给我们
        console.log('设置当前笔记:', response)
        currentNote.value = response
        return response
      } else {
        console.warn('获取笔记详情返回空数据')
        currentNote.value = null
      }
    } catch (error) {
      console.error(`获取笔记 ${id} 失败:`, error)
      currentNote.value = null
    } finally {
      isLoading.value = false
    }
  }

  // 创建笔记
  async function createNote(noteData) {
    isLoading.value = true
    try {
      // 调用API创建笔记
      const response = await api.notes.createNote(noteData)
      console.log('API创建笔记响应:', response)
      
      // 规范化响应数据：处理可能的不同响应格式
      // 1. 统一格式: response为笔记对象本身 (后端处理了data部分)
      // 2. 旧格式: response.note 包含笔记
      // 3. 新的统一格式: { success: true, data: 笔记对象 } - 已在拦截器处理
      const newNote = response || {}
      
      // 确保我们有有效的笔记对象
      if (newNote && newNote.id) {
        // 添加到本地数据集合
        notes.value.unshift(newNote)
        totalNotes.value++
      } else {
        console.warn('API返回的笔记数据格式不符合预期:', response)
      }
      
      // 返回原始响应，让调用者决定如何处理
      return response
    } catch (error) {
      console.error('创建笔记失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 更新笔记
  async function updateNote(id, noteData) {
    isLoading.value = true
    try {
      const updatedNote = await api.notes.updateNote(id, noteData)
      
      // 更新本地状态
      if (currentNote.value && currentNote.value.id === id) {
        currentNote.value = updatedNote
      }
      
      const index = notes.value.findIndex(note => note.id === id)
      if (index !== -1) {
        notes.value[index] = updatedNote
      }
      
      return updatedNote
    } catch (error) {
      console.error(`更新笔记 ${id} 失败:`, error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 删除笔记
  async function deleteNote(id) {
    isLoading.value = true
    try {
      await api.notes.deleteNote(id)
      
      // 更新本地状态
      notes.value = notes.value.filter(note => note.id !== id)
      if (currentNote.value && currentNote.value.id === id) {
        currentNote.value = null
      }
      totalNotes.value--
      
      return true
    } catch (error) {
      console.error(`删除笔记 ${id} 失败:`, error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 搜索笔记
  async function searchNotes(keyword) {
    searchKeyword.value = keyword
    return fetchNotes(1, pageSize.value)
  }

  // 清除搜索
  function clearSearch() {
    searchKeyword.value = ''
    activeTag.value = null
    return fetchNotes(1, pageSize.value)
  }

  // 根据标签筛选
  function filterByTag(tagId) {
    searchKeyword.value = ''
    return fetchNotes(1, pageSize.value, tagId)
  }

  return {
    notes,
    publicNotes,
    currentNote,
    isLoading,
    totalNotes,
    totalPublicNotes,
    currentPage,
    publicNotesPage,
    pageSize,
    activeTag,
    searchKeyword,
    hasNotes,
    hasPublicNotes,
    totalPages,
    totalPublicPages,
    fetchNotes,
    fetchPublicNotes,
    fetchNote,
    createNote,
    updateNote,
    deleteNote,
    searchNotes,
    clearSearch,
    filterByTag
  }
}) 