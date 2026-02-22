import { getApiBase } from './apiBase'
import type {
  CreateMESWorkRequest,
  CreatePositionRequest,
  CreateServiceGroupRequest,
  CreateTaskRequest,
  ListMESWorkFilters,
  ListMESFilters,
  MESWorkDashboardStats,
  MESWork,
  MESPosition,
  MESServiceGroup,
  MESTask,
  UpdateMESWorkTaskStatusRequest,
} from '@/types/mes'

const API_BASE = getApiBase()

interface MESApiError extends Error {
  status?: number
  data?: unknown
  cause?: Error
}

class MESApiService {
  private tasksUrl = `${API_BASE}/mes/tasks`
  private positionsUrl = `${API_BASE}/mes/positions`
  private serviceGroupsUrl = `${API_BASE}/mes/service-groups`
  private worksUrl = `${API_BASE}/mes/works`

  private getHeaders(additionalHeaders: Record<string, string> = {}): Record<string, string> {
    const token = localStorage.getItem('tramatex_auth_token')
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-User-ID': this.getCurrentUserId(),
      ...additionalHeaders,
    }

    if (token) {
      headers.Authorization = `Bearer ${token}`
    }

    return headers
  }

  private getCurrentUserId(): string {
    try {
      const userStr = localStorage.getItem('tramatex_user')
      if (userStr) {
        const user = JSON.parse(userStr)
        return user.id || 'anonymous'
      }
    } catch (error) {
      console.error('[MESApi] Error parsing user:', error)
    }

    return 'anonymous'
  }

  private async handleError(response: Response, message: string): Promise<never> {
    let errorData: { error?: string; message?: string } | undefined
    let rawBody = ''

    try {
      errorData = await response.clone().json()
    } catch {
      try {
        rawBody = (await response.text()).trim()
      } catch {
        rawBody = ''
      }
    }

    const responseMessage = errorData?.error || errorData?.message || rawBody

    let fallbackMessage = `${message} (HTTP ${response.status})`
    if (response.status === 401) {
      fallbackMessage = 'Sesión expirada o no autenticada. Inicia sesión nuevamente.'
    } else if (response.status === 403) {
      fallbackMessage = 'No tienes permisos para ejecutar esta acción en MES.'
    } else if (response.status === 404 && response.url.includes('/mes/')) {
      fallbackMessage = 'El backend activo no expone endpoints MES. Verifica que la API correcta esté levantada.'
    }

    const errorMessage = responseMessage || fallbackMessage
    const error = new Error(errorMessage) as MESApiError
    error.status = response.status
    error.data = errorData || rawBody
    throw error
  }

  private async safeFetch(url: string, options: RequestInit, fallbackMessage?: string): Promise<Response> {
    try {
      return await fetch(url, options)
    } catch (error) {
      const message =
        fallbackMessage ||
        `No se pudo conectar con el servidor. Verifica tu conexión o que la API esté activa. (URL: ${url})`
      const err = new Error(message) as MESApiError
      err.cause = error as Error
      throw err
    }
  }

  private buildQuery(filters: ListMESFilters = {}): string {
    const params = new URLSearchParams()

    if (filters.search) {
      params.append('search', filters.search)
    }
    if (filters.is_active !== undefined) {
      params.append('is_active', String(filters.is_active))
    }

    const query = params.toString()
    return query ? `?${query}` : ''
  }

  private buildWorkQuery(filters: ListMESWorkFilters = {}): string {
    const params = new URLSearchParams()

    if (filters.search) {
      params.append('search', filters.search)
    }
    if (filters.status) {
      params.append('status', filters.status)
    }
    if (filters.party_id) {
      params.append('party_id', filters.party_id)
    }

    const query = params.toString()
    return query ? `?${query}` : ''
  }

  async listTasks(filters: ListMESFilters = {}): Promise<MESTask[]> {
    const response = await this.safeFetch(`${this.tasksUrl}${this.buildQuery(filters)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las tareas MES')
    }

    return response.json()
  }

  async createTask(data: CreateTaskRequest): Promise<MESTask> {
    const response = await this.safeFetch(this.tasksUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la tarea MES')
    }

    return response.json()
  }

  async listPositions(filters: ListMESFilters = {}): Promise<MESPosition[]> {
    const response = await this.safeFetch(`${this.positionsUrl}${this.buildQuery(filters)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las posiciones MES')
    }

    return response.json()
  }

  async createPosition(data: CreatePositionRequest): Promise<MESPosition> {
    const response = await this.safeFetch(this.positionsUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la posición MES')
    }

    return response.json()
  }

  async listServiceGroups(filters: ListMESFilters = {}): Promise<MESServiceGroup[]> {
    const response = await this.safeFetch(`${this.serviceGroupsUrl}${this.buildQuery(filters)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los grupos de servicio MES')
    }

    return response.json()
  }

  async createServiceGroup(data: CreateServiceGroupRequest): Promise<MESServiceGroup> {
    const response = await this.safeFetch(this.serviceGroupsUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear el grupo de servicio MES')
    }

    return response.json()
  }

  async listWorks(filters: ListMESWorkFilters = {}): Promise<MESWork[]> {
    const response = await this.safeFetch(`${this.worksUrl}${this.buildWorkQuery(filters)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los trabajos MES')
    }

    return response.json()
  }

  async getWork(id: string): Promise<MESWork> {
    const response = await this.safeFetch(`${this.worksUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cargar el trabajo MES')
    }

    return response.json()
  }

  async createWork(data: CreateMESWorkRequest): Promise<MESWork> {
    const response = await this.safeFetch(this.worksUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear el trabajo MES')
    }

    return response.json()
  }

  async getWorkDashboardStats(): Promise<MESWorkDashboardStats> {
    const response = await this.safeFetch(`${this.worksUrl}/dashboard/stats`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las estadísticas MES')
    }

    return response.json()
  }

  async listOverdueWorks(limit?: number): Promise<MESWork[]> {
    const params = new URLSearchParams()
    if (limit && limit > 0) {
      params.append('limit', String(limit))
    }
    const query = params.toString()

    const response = await this.safeFetch(`${this.worksUrl}/overdue${query ? `?${query}` : ''}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los trabajos vencidos')
    }

    return response.json()
  }

  async updateWorkTaskStatus(workId: string, taskId: string, data: UpdateMESWorkTaskStatusRequest): Promise<MESWork> {
    const response = await this.safeFetch(`${this.worksUrl}/${workId}/tasks/${taskId}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el estado de la tarea MES')
    }

    return response.json()
  }
}

export const mesApi = new MESApiService()
