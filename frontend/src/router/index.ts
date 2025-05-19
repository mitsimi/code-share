import { createRouter, createWebHistory } from 'vue-router'
import SnippetsView from '../views/SnippetsView.vue'
import SnippetDetailsView from '../views/SnippetDetailsView.vue'
import AboutView from '../views/AboutView.vue'
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
      component: SnippetDetailsView,
    },
    {
      path: '/about',
      name: 'about',
      component: AboutView,
    },
  ],
})

export default router
