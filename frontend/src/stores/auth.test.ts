import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'
import { authService } from '@/services/auth'
import { usersService } from '@/services/users'
import { useRouter, type Router } from 'vue-router'

// Mock dependencies
vi.mock('@/services/auth', () => ({
  authService: {
    refreshToken: vi.fn(),
    logout: vi.fn(),
  },
}))

vi.mock('@/services/users', () => ({
  usersService: {
    getMe: vi.fn(),
  },
}))

vi.mock('vue-router', () => ({
  useRouter: vi.fn(),
}))

vi.mock('vue-sonner', () => ({
  toast: {
    error: vi.fn(),
  },
}))

// Mock WebSocket store import since it's dynamically imported in logout
vi.mock('./websocket', () => ({
  useWebSocketStore: () => ({
    cleanupAuthenticatedSubscriptions: vi.fn(),
  }),
}))

describe('Auth Store', () => {
  const mockRouter = {
    push: vi.fn(),
    currentRoute: { value: { fullPath: '/current' } },
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
    vi.mocked(useRouter).mockReturnValue(mockRouter as unknown as Router)
  })

  const mockUser = {
    id: '1',
    username: 'testuser',
    avatar: '',
    email: 'test@example.com',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  }

  const mockAuthData = {
    user: mockUser,
    token: 'access-token',
    refreshToken: 'refresh-token',
    expiresAt: Math.floor(Date.now() / 1000) + 3600, // 1 hour from now
  }

  it('initializes from localStorage if valid', async () => {
    // Setup localStorage
    localStorage.setItem('token', mockAuthData.token)
    localStorage.setItem('refreshToken', mockAuthData.refreshToken)
    localStorage.setItem('expiresAt', String(mockAuthData.expiresAt))
    localStorage.setItem('user', JSON.stringify(mockAuthData.user))

    // Mock getMe success
    vi.mocked(usersService.getMe).mockResolvedValue(mockUser)

    const store = useAuthStore()
    await store.initializeAuth()

    expect(store.user).toEqual(mockUser)
    expect(store.isAuthenticated()).toBe(true)
    expect(usersService.getMe).toHaveBeenCalled()
  })

  it('sets auth state correctly', () => {
    const store = useAuthStore()
    store.setAuth(mockAuthData)

    expect(store.user).toEqual(mockUser)
    expect(store.token).toBe(mockAuthData.token)
    expect(store.isAuthenticated()).toBe(true)

    // Verify persistence
    expect(localStorage.getItem('token')).toBe(mockAuthData.token)
  })

  it('logs out correctly', async () => {
    const store = useAuthStore()
    store.setAuth(mockAuthData)

    await store.logout()

    expect(store.user).toBeNull()
    expect(store.token).toBeNull()
    expect(localStorage.getItem('token')).toBeNull()
    expect(authService.logout).toHaveBeenCalled()
  })

  it('refreshes token successfully', async () => {
    const store = useAuthStore()
    store.setAuth(mockAuthData)

    const newAuthData = {
      ...mockAuthData,
      token: 'new-token',
    }

    vi.mocked(authService.refreshToken).mockResolvedValue(newAuthData)

    await store.refreshAccessToken()

    expect(store.token).toBe('new-token')
  })

  it('handles failed refresh by clearing auth', async () => {
    const store = useAuthStore()
    store.setAuth(mockAuthData)

    vi.mocked(authService.refreshToken).mockRejectedValue(new Error('Failed'))

    await expect(store.refreshAccessToken()).rejects.toThrow('Failed')

    expect(store.user).toBeNull()
    expect(store.token).toBeNull()
  })
})
