import axios, { AxiosInstance } from 'axios'
import type { Usuario, UserRole } from '@/types/auth'
import { getAuthBase } from './apiBase'

const API_URL = getAuthBase()

function toServiceError(error: any): Error {
  const backendError = error?.response?.data?.error || error?.response?.data?.message
  const fallback = error?.message || 'Error de comunicación con IAM'
  return new Error(backendError || fallback)
}

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
    try {
      const response = await this.apiClient.get('/auth/users')
      return response.data.users
    } catch (error: any) {
      throw toServiceError(error)
    }
  }

  async assignRole(userId: string, role: UserRole): Promise<{ userId: string; role: UserRole }> {
    try {
      const response = await this.apiClient.post('/auth/assign-role', {
        user_id: userId,
        role
      })
      return {
        userId: response.data.user_id,
        role: response.data.role
      }
    } catch (error: any) {
      throw toServiceError(error)
    }
  }

  async createUser(payload: { email: string; password: string; role: UserRole }): Promise<Usuario> {
    try {
      const response = await this.apiClient.post('/auth/users', payload)
      return response.data
    } catch (error: any) {
      throw toServiceError(error)
    }
  }

  async updateUser(userId: string, payload: { email?: string; password?: string }): Promise<Usuario> {
    try {
      const response = await this.apiClient.put(`/auth/users/${userId}`, payload)
      return response.data
    } catch (error: any) {
      throw toServiceError(error)
    }
  }

  async deleteUser(userId: string): Promise<void> {
    try {
      await this.apiClient.delete(`/auth/users/${userId}`)
    } catch (error: any) {
      throw toServiceError(error)
    }
  }
}

export const iamService = new IamService()
