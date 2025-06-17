export interface Snippet {
  id: string
  title: string
  content: string
  language: string
  author: string
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
