import axios, { AxiosInstance } from 'axios'
import type { Usuario, UserRole } from '@/types/auth'
import { getApiBase } from './apiBase'

const API_URL = getApiBase()

class IamService {
  private apiClient: AxiosInstance

  constructor() {
    this.apiClient = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    this.apiClient.interceptors.request.use((config) => {
      const token = localStorage.getItem('tramatex_auth_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    this.apiClient.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          window.location.href = '/login'
        }
        return Promise.reject(error)
      }
    )
  }

  async listUsers(): Promise<Usuario[]> {
    const response = await this.apiClient.get('/auth/users')
    return response.data.users
  }

  async assignRole(userId: string, role: UserRole): Promise<{ userId: string; role: UserRole }> {
    const response = await this.apiClient.post('/auth/assign-role', {
      user_id: userId,
      role
    })
    return {
      userId: response.data.user_id,
      role: response.data.role
    }
  }

  async createUser(payload: { email: string; password: string; role: UserRole }): Promise<Usuario> {
    const response = await this.apiClient.post('/auth/users', payload)
    return response.data
  }

  async deleteUser(userId: string): Promise<void> {
    await this.apiClient.delete(`/auth/users/${userId}`)
  }
}

export const iamService = new IamService()
