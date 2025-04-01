import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import NotesView from '../views/NotesView.vue'
import NoteDetailView from '../views/NoteDetailView.vue'
import NoteEditView from '../views/NoteEditView.vue'
import TagsView from '../views/TagsView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ProfileView from '../views/ProfileView.vue'
import LibraryView from '../views/LibraryView.vue'
import { useUserStore } from '../stores/user'

// 白名单路由 - 这些路由不需要登录就能访问
const whiteList = ['home', 'login', 'register']

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/notes',
      name: 'notes',
      component: NotesView,
      meta: { requiresAuth: true }
    },
    {
      path: '/notes/new',
      name: 'note-new',
      component: NoteEditView,
      meta: { requiresAuth: true }
    },
    {
      path: '/notes/:id',
      name: 'note-detail',
      component: NoteDetailView,
      meta: { requiresAuth: true }
    },
    {
      path: '/notes/:id/edit',
      name: 'note-edit',
      component: NoteEditView,
      meta: { requiresAuth: true }
    },
    {
      path: '/tags',
      name: 'tags',
      component: TagsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/library',
      name: 'library',
      component: LibraryView,
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    },
  ],
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  // 使用Pinia存储的认证状态
  const userStore = useUserStore()
  const hasToken = !!userStore.token
  const isPublicPage = whiteList.includes(to.name)
  
  console.log('路由导航守卫 - 目标:', to.path, '是否公开页面:', isPublicPage, '是否有token:', hasToken)
  
  // 如果有token
  if (hasToken) {
    // 已登录用户访问登录/注册页面，重定向到笔记页
    if (to.name === 'login' || to.name === 'register') {
      console.log('已登录用户访问登录/注册页面，重定向到笔记页')
      next({ name: 'notes' })
    } else {
      // 其他页面，先获取用户信息（如果还没有）
      // 只有在user对象为空时才请求用户信息，避免频繁请求
      if (Object.keys(userStore.user).length === 0) {
        try {
          await userStore.fetchUserInfo()
        } catch (error) {
          // 如果获取用户信息失败（如token过期），清除token并跳转到登录页
          console.error('获取用户信息失败:', error)
          userStore.logout()
          next({ name: 'login', query: { redirect: to.fullPath } })
          return
        }
      }
      // 正常访问
      next()
    }
  } else {
    // 未登录状态
    if (isPublicPage) {
      // 公开页面允许访问
      console.log('未登录用户访问公开页面，允许访问')
      next()
    } else {
      // 非公开页面，重定向到登录页
      console.log('未登录用户访问受保护页面，重定向到登录页')
      next({ name: 'login', query: { redirect: to.fullPath } })
    }
  }
})

export default router
