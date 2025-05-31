import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

type Theme = 'light' | 'dark' | 'system'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref<Theme>('system')
  const systemTheme = ref<'light' | 'dark'>('light')

  // Computed property that returns the actual theme being applied
  const resolvedTheme = computed<'light' | 'dark'>(() => {
    if (theme.value === 'system') {
      return systemTheme.value
    }
    return theme.value
  })

  // Initialize theme based on user preference or saved setting
  const initTheme = () => {
    // Check localStorage first
    const savedTheme = localStorage.getItem('theme') as Theme | null
    if (savedTheme) {
      theme.value = savedTheme
    }

    // Initialize system theme
    updateSystemTheme()

    // Listen for system theme changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateSystemTheme)

    // Apply theme
    applyTheme()
  }

  // Update system theme based on user preference
  const updateSystemTheme = () => {
    systemTheme.value = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    if (theme.value === 'system') {
      applyTheme()
    }
  }

  // Change theme to a specific value
  const changeTheme = (newTheme: Theme) => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    applyTheme()
  }

  // Apply theme to body element
  const applyTheme = () => {
    const currentTheme = resolvedTheme.value
    document.body.classList.remove('light', 'dark')
    document.body.classList.add(currentTheme)
  }

  return {
    theme,
    resolvedTheme,
    initTheme,
    changeTheme,
  }
})
