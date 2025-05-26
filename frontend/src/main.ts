import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'

const app = createApp(App)
const pinia = createPinia()

app.use(router).use(pinia).use(VueQueryPlugin)

app.mount('#app')
