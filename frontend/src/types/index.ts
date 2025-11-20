export interface Snippet {
  id: string
  title: string
  content: string
  language: string
  author: User
  views: number
  likes: number
  isLiked: boolean
  isSaved: boolean
  createdAt: string
  updatedAt: string
}

export interface User {
  id: string
  username: string
  avatar: string
  email: string
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
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

// Backend API Response structure matching the Go struct
export interface APIResponse<T = unknown> {
  statusCode: number
  message: string
  data?: T
  error?: string
  timestamp?: string
}
