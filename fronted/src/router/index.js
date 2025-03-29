import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import NotesView from '../views/NotesView.vue'
import NoteDetailView from '../views/NoteDetailView.vue'
import NoteEditView from '../views/NoteEditView.vue'
import TagsView from '../views/TagsView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ProfileView from '../views/ProfileView.vue'

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
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const isAuthenticated = localStorage.getItem('user') !== null
  
  if (requiresAuth && !isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
