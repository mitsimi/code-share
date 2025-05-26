export interface User {
  id: string
  username: string
  email: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  token: string
  user: User
  expires_at: number
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
