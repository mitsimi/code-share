import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'
import { authService } from '@/services/auth'
import { toast } from 'sonner'
import { useRouter } from 'vue-router'

interface AuthState {
  user: User | null
  token: string | null
  expiresAt: number | null
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const expiresAt = ref<number | null>(null)
  const router = useRouter()

  const setAuth = (auth: AuthState) => {
    user.value = auth.user
    token.value = auth.token
    expiresAt.value = auth.expiresAt
  }

  const clearAuth = () => {
    user.value = null
    token.value = null
    expiresAt.value = null
  }

  const isAuthenticated = () => {
    return !!token.value && !!expiresAt.value && expiresAt.value > Date.now() / 1000
  }

  const refreshToken = async () => {
    try {
      const response = await authService.refreshToken()
      setAuth({
        user: response.user,
        token: response.token,
        expiresAt: response.expires_at,
      })
      return response
    } catch (error) {
      clearAuth()
      throw error
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
    expiresAt,
    setAuth,
    clearAuth,
    isAuthenticated,
    refreshToken,
    logout,
  }
})
