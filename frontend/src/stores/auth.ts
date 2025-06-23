import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'
import { authService } from '@/services/auth'
import { usersService } from '@/services/users'
import { toast } from 'vue-sonner'
import { useRouter } from 'vue-router'

interface AuthenticatedUser {
  user: User
  token: string
  refreshToken: string
  expiresAt: number
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const expiresAt = ref<number | null>(null)

  const router = useRouter()

  // Add a ref to store the refresh timer
  const refreshTimer = ref<number | null>(null)

  // Function to schedule token refresh
  const scheduleTokenRefresh = (expiresAt: number) => {
    if (refreshTimer.value) {
      window.clearTimeout(refreshTimer.value)
    }

    // Calculate time until 5 minutes before expiration
    const refreshTime = expiresAt * 1000 - 5 * 60 * 1000 - Date.now()

    // Only schedule if the refresh time is in the future
    if (refreshTime > 0) {
      refreshTimer.value = window.setTimeout(async () => {
        try {
          await refreshAccessToken()
        } catch (error) {
          clearAuth()
          router.push({
            name: 'login',
            query: { redirect: router.currentRoute.value.fullPath },
          })
        }
      }, refreshTime)
    }
  }

  // Load auth state from localStorage.
  // Returns true if it could load from local storage.
  const loadAuthFromStorage = (): boolean => {
    const storedToken = localStorage.getItem('token')
    const storedRefreshToken = localStorage.getItem('refreshToken')
    const storedExpiresAt = localStorage.getItem('expiresAt')
    const storedUser = localStorage.getItem('user')

    if (storedToken && storedRefreshToken && storedExpiresAt && storedUser) {
      setAuth({
        token: storedToken,
        refreshToken: storedRefreshToken,
        expiresAt: Number(storedExpiresAt),
        user: JSON.parse(storedUser),
      })
      return true
    }
    return false
  }

  // Save auth state to localStorage
  const saveAuthToStorage = (auth: AuthenticatedUser) => {
    if (!auth || !auth.token || !auth.refreshToken || !auth.expiresAt || !auth.user) {
      throw new Error('Invalid auth data provided to store in local storage')
    }
    localStorage.setItem('token', auth.token)
    localStorage.setItem('refreshToken', auth.refreshToken)
    localStorage.setItem('expiresAt', String(auth.expiresAt))
    localStorage.setItem('user', JSON.stringify(auth.user))
  }

  // Clear auth state from localStorage
  const clearAuthFromStorage = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('expiresAt')
    localStorage.removeItem('user')
  }

  const setAuth = (auth: AuthenticatedUser) => {
    user.value = auth.user
    token.value = auth.token
    refreshToken.value = auth.refreshToken
    expiresAt.value = auth.expiresAt
    saveAuthToStorage(auth)
    scheduleTokenRefresh(auth.expiresAt)
  }

  const setUser = (data: User) => {
    user.value = data
    localStorage.setItem('user', JSON.stringify(data))
  }

  const clearAuth = () => {
    if (refreshTimer.value) {
      window.clearTimeout(refreshTimer.value)
      refreshTimer.value = null
    }

    user.value = null
    token.value = null
    refreshToken.value = null
    expiresAt.value = null
    clearAuthFromStorage()
  }

  const isAuthenticated = () => {
    // Checks if the user is authenticated:
    // - token.value: access token is present
    // - user.value: user object is present
    // - expiresAt.value: token expiration timestamp is present
    // - expiresAt.value > Date.now() / 1000: token has not expired
    return !!token.value && !!user.value && !!expiresAt.value && expiresAt.value > Date.now() / 1000
  }

  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await authService.refreshToken(refreshToken.value)
      setAuth({
        user: response.user,
        token: response.token,
        refreshToken: response.refreshToken,
        expiresAt: response.expiresAt,
      })
      return response
    } catch (error) {
      clearAuth()
      throw error
    }
  }

  const initializeAuth = async () => {
    const hasStoredAuth = loadAuthFromStorage()
    if (hasStoredAuth && isAuthenticated()) {
      try {
        // Verify the session is still valid by calling the backend
        const currentUser = await usersService.getMe()

        // Update user data in case it changed on the server
        setUser(currentUser)

        // Schedule token refresh since session is valid
        if (expiresAt.value) {
          scheduleTokenRefresh(expiresAt.value)
        }
      } catch (error) {
        // Session is invalid, try to refresh the token
        try {
          await refreshAccessToken()
        } catch (refreshError) {
          // If refresh fails, clear auth
          clearAuth()
          throw refreshError
        }
      }
    } else if (hasStoredAuth) {
      // If stored auth exists but is expired, try to refresh
      try {
        await refreshAccessToken()
      } catch (error) {
        clearAuth()
      }
    }
  }

  const logout = async () => {
    try {
      await authService.logout()
      clearAuth()
      toast.success('You have been logged out successfully')
      router.push('/login')
    } catch (error) {
      toast.error('Failed to logout')
      // Still clear local auth state even if the server request fails
      clearAuth()
      router.push('/login')
    }
  }

  return {
    user,
    token,
    refreshToken,
    expiresAt,
    setAuth,
    setUser,
    isAuthenticated,
    refreshAccessToken,
    initializeAuth,
    logout,
    clearAuth,
  }
})
