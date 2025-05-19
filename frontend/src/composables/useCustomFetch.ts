import { createFetch } from '@vueuse/core'

export const useCustomFetch = createFetch({
  baseUrl: 'http://localhost:8080/api',
  fetchOptions: {
    mode: 'cors',
  },
})
