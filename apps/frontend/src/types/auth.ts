export interface Usuario {
  id: string
  email: string
  nombre: string
  rol: 'admin' | 'usuario'
  creadoEn: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  usuario: Usuario
}

export interface TokenClaims {
  sub: string
  email: string
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
