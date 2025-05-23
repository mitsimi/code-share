import { createRouter, createWebHistory } from 'vue-router'
import SnippetsView from '../views/SnippetsView.vue'
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/snippets',
    },
    {
      path: '/snippets',
      name: 'snippets',
      component: SnippetsView,
    },
    {
      path: '/snippets/:snippetId',
      name: 'snippet-details',
      component: () => import('@/views/SnippetDetailsView.vue'),
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue'),
    },
  ],
})

export default router
