import { ref } from 'vue'

export type ToastType = 'info' | 'success' | 'error'

interface Toast {
  id: number
  message: string
  type: ToastType
  timeout?: number
  timer?: number
}

const toasts = ref<Toast[]>([])
let nextId = 1

export function useToast() {
  const showToast = (message: string, type: ToastType = 'info', timeout = 3000) => {
    const id = nextId++
    const toast: Toast = {
      id,
      message,
      type,
      timeout,
    }

    toasts.value.push(toast)

    if (timeout > 0) {
      toast.timer = window.setTimeout(() => {
        removeToast(id)
      }, timeout)
    }

    return id
  }

  const removeToast = (id: number) => {
    const index = toasts.value.findIndex((t) => t.id === id)
    if (index !== -1) {
      const toast = toasts.value[index]
      if (toast.timer) {
        clearTimeout(toast.timer)
      }
      toasts.value.splice(index, 1)
    }
  }

  const pauseToast = (id: number) => {
    const toast = toasts.value.find((t) => t.id === id)
    if (toast && toast.timer) {
      clearTimeout(toast.timer)
      toast.timer = undefined
    }
  }

  const resumeToast = (id: number) => {
    const toast = toasts.value.find((t) => t.id === id)
    if (toast && toast.timeout && !toast.timer) {
      toast.timer = window.setTimeout(() => {
        removeToast(id)
      }, toast.timeout)
    }
  }

  return {
    toasts,
    showToast,
    removeToast,
    pauseToast,
    resumeToast,
  }
}
