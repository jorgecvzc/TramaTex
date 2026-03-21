import { getApiBase } from './apiBase'
import type {
  CreateWorkOrderRequest,
  CreatePositionRequest,
  CreateWorkTypeRequest,
  CreateTaskRequest,
  CreateWorkSetupRequest,
  ListWorkOrderFilters,
  ListMESFilters,
  ListWorkSetupFilters,
  WorkOrderDashboardStats,
  WorkOrder,
  MESPosition,
  MESWorkType,
  MESTask,
  UpdateWorkOrderRequest,
  UpdateWorkOrderTaskStatusRequest,
  UpdatePositionRequest,
  UpdateWorkTypeRequest,
  UpdateTaskRequest,
  UpdateWorkSetupRequest,
  WorkSetup,
  PendingWorkSetup,
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
  private workTypesUrl = `${API_BASE}/mes/work-types`
  private workOrdersUrl = `${API_BASE}/mes/work-orders`
  private workSetupsUrl = `${API_BASE}/mes/work-setups`

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

  private buildWorkQuery(filters: ListWorkOrderFilters = {}): string {
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
      await this.handleError(response, 'No se pudo crear la tarea')
    }

    return response.json()
  }

  async getTask(id: string): Promise<MESTask> {
    const response = await this.safeFetch(`${this.tasksUrl}/${encodeURIComponent(id)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cargar la tarea')
    }

    return response.json()
  }

  async updateTask(id: string, data: UpdateTaskRequest): Promise<MESTask> {
    const response = await this.safeFetch(`${this.tasksUrl}/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la tarea')
    }

    return response.json()
  }

  async deleteTask(id: string): Promise<void> {
    const response = await this.safeFetch(`${this.tasksUrl}/${encodeURIComponent(id)}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar la tarea')
    }
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
      await this.handleError(response, 'No se pudo crear la posición')
    }

    return response.json()
  }

  async getPosition(id: string): Promise<MESPosition> {
    const response = await this.safeFetch(`${this.positionsUrl}/${encodeURIComponent(id)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cargar la posición')
    }

    return response.json()
  }

  async updatePosition(id: string, data: UpdatePositionRequest): Promise<MESPosition> {
    const response = await this.safeFetch(`${this.positionsUrl}/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la posición')
    }

    return response.json()
  }

  async deletePosition(id: string): Promise<void> {
    const response = await this.safeFetch(`${this.positionsUrl}/${encodeURIComponent(id)}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar la posición')
    }
  }

  async listWorkTypes(filters: ListMESFilters = {}): Promise<MESWorkType[]> {
    const response = await this.safeFetch(`${this.workTypesUrl}${this.buildQuery(filters)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los tipos de trabajo MES')
    }

    return response.json()
  }

  async createWorkType(data: CreateWorkTypeRequest): Promise<MESWorkType> {
    const response = await this.safeFetch(this.workTypesUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear el tipo de trabajo MES')
    }

    return response.json()
  }

  async getWorkType(id: string): Promise<MESWorkType> {
    const response = await this.safeFetch(`${this.workTypesUrl}/${encodeURIComponent(id)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cargar el tipo de trabajo')
    }

    return response.json()
  }

  async updateWorkType(id: string, data: UpdateWorkTypeRequest): Promise<MESWorkType> {
    const response = await this.safeFetch(`${this.workTypesUrl}/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el tipo de trabajo')
    }

    return response.json()
  }

  async deleteWorkType(id: string): Promise<void> {
    const response = await this.safeFetch(`${this.workTypesUrl}/${encodeURIComponent(id)}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar el tipo de trabajo')
    }
  }

  async listWorkOrders(filters: ListWorkOrderFilters = {}): Promise<WorkOrder[]> {
    const response = await this.safeFetch(`${this.workOrdersUrl}${this.buildWorkQuery(filters)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las órdenes de trabajo MES')
    }

    return response.json()
  }

  async getWorkOrder(id: string): Promise<WorkOrder> {
    const response = await this.safeFetch(`${this.workOrdersUrl}/${id}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo cargar la orden de trabajo MES')
    }

    return response.json()
  }

  async createWorkOrder(data: CreateWorkOrderRequest): Promise<WorkOrder> {
    const response = await this.safeFetch(this.workOrdersUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la orden de trabajo MES')
    }

    return response.json()
  }

  async updateWorkOrder(id: string, data: UpdateWorkOrderRequest): Promise<WorkOrder> {
    const response = await this.safeFetch(`${this.workOrdersUrl}/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la orden de trabajo MES')
    }

    return response.json()
  }

  async getWorkOrderDashboardStats(): Promise<WorkOrderDashboardStats> {
    const response = await this.safeFetch(`${this.workOrdersUrl}/dashboard/stats`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las estadísticas MES')
    }

    return response.json()
  }

  async listOverdueWorkOrders(limit?: number): Promise<WorkOrder[]> {
    const params = new URLSearchParams()
    if (limit && limit > 0) {
      params.append('limit', String(limit))
    }
    const query = params.toString()

    const response = await this.safeFetch(`${this.workOrdersUrl}/overdue${query ? `?${query}` : ''}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar los trabajos vencidos')
    }

    return response.json()
  }

  async updateWorkOrderTaskStatus(workId: string, taskId: string, data: UpdateWorkOrderTaskStatusRequest): Promise<WorkOrder> {
    const response = await this.safeFetch(`${this.workOrdersUrl}/${workId}/tasks/${taskId}/status`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar el estado de la tarea MES')
    }

    return response.json()
  }

  // --- WorkSetup methods ---

  private buildWorkSetupQuery(filters: ListWorkSetupFilters = {}): string {
    const params = new URLSearchParams()
    if (filters.search) params.append('search', filters.search)
    if (filters.is_active !== undefined) params.append('is_active', String(filters.is_active))
    if (filters.party_id) params.append('party_id', filters.party_id)
    const query = params.toString()
    return query ? `?${query}` : ''
  }

  async listWorkSetups(filters: ListWorkSetupFilters = {}): Promise<WorkSetup[]> {
    const response = await this.safeFetch(`${this.workSetupsUrl}${this.buildWorkSetupQuery(filters)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })
    if (!response.ok) {
      await this.handleError(response, 'No se pudieron cargar las configuraciones de cliente')
    }
    return response.json()
  }

  async getWorkSetup(id: string): Promise<WorkSetup> {
    const response = await this.safeFetch(`${this.workSetupsUrl}/${encodeURIComponent(id)}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })
    if (!response.ok) {
      await this.handleError(response, 'No se pudo cargar la configuración de cliente')
    }
    return response.json()
  }

  async createWorkSetup(data: CreateWorkSetupRequest): Promise<WorkSetup> {
    const response = await this.safeFetch(this.workSetupsUrl, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })
    if (!response.ok) {
      await this.handleError(response, 'No se pudo crear la configuración de cliente')
    }
    return response.json()
  }

  async updateWorkSetup(id: string, data: UpdateWorkSetupRequest): Promise<WorkSetup> {
    const response = await this.safeFetch(`${this.workSetupsUrl}/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    })
    if (!response.ok) {
      await this.handleError(response, 'No se pudo actualizar la configuración de cliente')
    }
    return response.json()
  }

  async deleteWorkSetup(id: string): Promise<void> {
    const response = await this.safeFetch(`${this.workSetupsUrl}/${encodeURIComponent(id)}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })
    if (!response.ok) {
      await this.handleError(response, 'No se pudo eliminar la configuración de cliente')
    }
  }

  // --- Pending WorkSetups (from confirmed Sales orders) ---

  async listPendingWorkSetups(): Promise<PendingWorkSetup[]> {
    const response = await this.safeFetch(`${API_BASE}/mes/pending-work-setups`, {
      method: 'GET',
      headers: this.getHeaders(),
    })
    if (!response.ok) {
      await this.handleError(response, 'No se pudieron obtener las configuraciones pendientes')
    }
    return response.json()
  }

  // --- Status & priority label helpers ---

  getWorkStatusLabel(status: string): string {
    const labels: Record<string, string> = {
      'PENDING': 'Pendiente',
      'IN_PROGRESS': 'En progreso',
      'ON_HOLD': 'En espera',
      'SUSPENDED': 'Suspendida',
      'COMPLETED': 'Completado',
      'CANCELLED': 'Cancelado',
    }
    return labels[status] || status
  }

  getPriorityLabel(priority: string): string {
    const labels: Record<string, string> = {
      'LOW': 'Baja',
      'NORMAL': 'Normal',
      'HIGH': 'Alta',
      'URGENT': 'Urgente',
    }
    return labels[priority] || priority
  }

  getTaskStatusLabel(status: string): string {
    const labels: Record<string, string> = {
      'PENDING': 'Pendiente',
      'IN_PROGRESS': 'En progreso',
      'PAUSED': 'Pausada',
      'COMPLETED': 'Completada',
      'BLOCKED': 'Bloqueada',
    }
    return labels[status] || status
  }
}

export const mesApi = new MESApiService()
