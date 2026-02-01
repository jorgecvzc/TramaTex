export type UserRole = 'admin' | 'commercial' | 'designer' | 'workshop'

export interface Usuario {
  id: string
  email: string
  role: UserRole
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  user: Usuario
}

export interface TokenClaims {
  sub: string
  email: string
  role: UserRole
  exp: number
  iat: number
}

export type ErrorType = 'validation' | 'auth' | 'network' | 'server' | 'unknown'

export interface AuthError {
  type: ErrorType
  message: string
  timestamp: number
}

export interface AuthState {
  user: Usuario | null
  accessToken: string | null
  refreshToken: string | null
  isLoading: boolean
  error: string | null
  tokenExpiresAt: number | null
}
