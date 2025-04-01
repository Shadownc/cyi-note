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

// 布局组件
import DefaultLayout from '../layouts/DefaultLayout.vue'
import BlankLayout from '../layouts/BlankLayout.vue'
import AdminLayout from '../layouts/AdminLayout.vue'

// 管理后台子视图
import AdminUsersView from '../views/admin/UsersView.vue'
import AdminNotesView from '../views/admin/NotesView.vue'
import AdminTagsView from '../views/admin/TagsView.vue'
import AdminSettingsView from '../views/admin/SettingsView.vue'

import { useUserStore } from '../stores/user'

// 白名单路由 - 这些路由不需要登录就能访问
const whiteList = ['home', 'login', 'register']

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: DefaultLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView,
        },
        {
          path: 'notes',
          name: 'notes',
          component: NotesView,
          meta: { requiresAuth: true }
        },
        {
          path: 'notes/new',
          name: 'note-new',
          component: NoteEditView,
          meta: { requiresAuth: true }
        },
        {
          path: 'notes/:id',
          name: 'note-detail',
          component: NoteDetailView,
          meta: { requiresAuth: true }
        },
        {
          path: 'notes/:id/edit',
          name: 'note-edit',
          component: NoteEditView,
          meta: { requiresAuth: true }
        },
        {
          path: 'tags',
          name: 'tags',
          component: TagsView,
          meta: { requiresAuth: true }
        },
        {
          path: 'library',
          name: 'library',
          component: LibraryView,
          meta: { requiresAuth: true }
        },
        {
          path: 'profile',
          name: 'profile',
          component: ProfileView,
          meta: { requiresAuth: true }
        },
      ]
    },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          redirect: '/admin/users'
        },
        {
          path: 'users',
          name: 'admin-users',
          component: AdminUsersView,
          meta: { requiresAuth: true, requiresAdmin: true }
        },
        {
          path: 'notes',
          name: 'admin-notes',
          component: AdminNotesView,
          meta: { requiresAuth: true, requiresAdmin: true }
        },
        {
          path: 'tags',
          name: 'admin-tags',
          component: AdminTagsView,
          meta: { requiresAuth: true, requiresAdmin: true }
        },
        {
          path: 'settings',
          name: 'admin-settings',
          component: AdminSettingsView,
          meta: { requiresAuth: true, requiresAdmin: true }
        }
      ]
    },
    {
      path: '/auth',
      component: BlankLayout,
      children: [
        {
          path: 'login',
          name: 'login',
          component: LoginView
        },
        {
          path: 'register',
          name: 'register',
          component: RegisterView
        },
      ]
    },
    // 重定向旧的登录/注册路径到新路径
    {
      path: '/login',
      redirect: '/auth/login'
    },
    {
      path: '/register',
      redirect: '/auth/register'
    }
  ],
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  // 获取用户存储
  const userStore = useUserStore()
  
  // 如果已登录但没有用户信息，尝试获取
  if (userStore.isAuthenticated && Object.keys(userStore.user).length === 0) {
    try {
      await userStore.fetchUserInfo()
    } catch (error) {
      console.error('获取用户信息失败:', error)
      userStore.logout()
      next({ path: '/auth/login', query: { redirect: to.fullPath } })
      return
    }
  }
  
  // 如果需要认证的路由
  if (to.meta.requiresAuth) {
    // 检查用户是否已登录
    if (!userStore.isAuthenticated) {
      next({
        path: '/auth/login',
        query: { redirect: to.fullPath }
      })
      return
    }
    
    // 如果需要管理员权限
    if (to.meta.requiresAdmin && userStore.user.role !== 'admin') {
      // 不是管理员，跳转到主页
      next({ path: '/' })
      return
    }
  }
  
  // 如果已登录且访问登录/注册页面，重定向到首页
  if (userStore.isAuthenticated && (to.path === '/auth/login' || to.path === '/auth/register')) {
    next({ path: '/' })
    return
  }
  
  // 其他情况允许访问
  next()
})

export default router
