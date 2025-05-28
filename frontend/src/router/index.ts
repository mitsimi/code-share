import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/snippets',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('@/views/auth/SignupView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/snippets',
      name: 'snippets',
      component: () => import('@/views/snippets/SnippetsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/snippets/:snippetId',
      name: 'snippet-details',
      component: () => import('@/views/snippets/SnippetDetailsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue'),
    },
    {
      path: '/_/sandbox',
      name: 'sandbox',
      component: () => import('@/views/Sandbox.vue'),
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
