import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/home/HomeView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/snippets',
      name: 'snippets',
      component: () => import('@/views/snippets/SnippetsView.vue'),
      //meta: { requiresAuth: true },
    },
    {
      path: '/snippets/:snippetId',
      name: 'snippet-details',
      component: () => import('@/views/snippets/SnippetDetailsView.vue'),
      //meta: { requiresAuth: true },
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/profile/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Check if the route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated()) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  // Check if the route requires guest (not authenticated)
  if (to.meta.requiresGuest && authStore.isAuthenticated()) {
    next({ name: 'snippets' })
    return
  }

  next()
})

export default router
