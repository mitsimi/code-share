export interface Snippet {
  id: string
  title: string
  content: string
  author: string
  likes: number
  is_liked: boolean
  created_at: string
  updated_at: string
}

export interface User {
  id: string
  username: string
  email: string
  created_at: string
  updated_at: string
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
  refresh_token: string
  user: User
  expires_at: number
}

export interface ApiResponse<T> {
  data: T
  message?: string
  error?: string
}
