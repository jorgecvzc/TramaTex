import axios, { AxiosInstance } from 'axios'
import type { LoginResponse, Usuario } from '@/types/auth'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

class AuthService {
  private apiClient: AxiosInstance

  constructor() {
    this.apiClient = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Interceptor para agregar token
    this.apiClient.interceptors.request.use((config) => {
      const token = localStorage.getItem('tramatex_auth_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    // Interceptor para manejo de errores
    this.apiClient.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('tramatex_auth_token')
          window.location.href = '/login'
        }
        return Promise.reject(error)
      }
    )
  }

  async login(email: string, password: string): Promise<LoginResponse> {
    const response = await this.apiClient.post('/auth/login', {
      email,
      password
    })
    const data = response.data
    return {
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
      expiresIn: data.expires_in,
      user: data.user
    }
  }

  async refreshToken(refreshToken: string): Promise<{ accessToken: string; expiresIn: number }> {
    const response = await this.apiClient.post('/auth/refresh', {
      refreshToken
    })
    const data = response.data
    return {
      accessToken: data.access_token,
      expiresIn: data.expires_in
    }
  }

  async getCurrentUser(): Promise<Usuario> {
    const response = await this.apiClient.get<Usuario>('/auth/me')
    return response.data
  }

  async logout(): Promise<void> {
    try {
      await this.apiClient.post('/auth/logout')
    } finally {
      localStorage.removeItem('tramatex_auth_token')
      localStorage.removeItem('tramatex_refresh_token')
    }
  }
}

export const authService = new AuthService()
