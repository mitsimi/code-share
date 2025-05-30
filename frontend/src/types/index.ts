export interface Snippet {
  id: string
  title: string
  content: string
  author: string
  likes: number
  isLiked: boolean
  createdAt: string
  updatedAt: string
  language?: string
  isSaved?: boolean
}

export interface User {
  id: string
  username: string
  email: string
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface SignupRequest {
  username: string
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  refreshToken: string
  user: User
  expiresAt: number
}

export interface ApiResponse<T> {
  data: T
  message?: string
  error?: string
}
