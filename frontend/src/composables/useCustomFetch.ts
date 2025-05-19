import { createFetch } from '@vueuse/core'

export const useCustomFetch = createFetch({
  fetchOptions: {
    mode: 'cors',
  },
})
