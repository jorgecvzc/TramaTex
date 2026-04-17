import { api } from './api'
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

class MESApiService {
  private readonly moduleBase = '/mes'

  private async handleError(error: any, defaultMessage: string): Promise<never> {
    const errorData = error.response?.data
    let message = errorData?.error || errorData?.message || error.message || defaultMessage

    if (error.response?.status === 401) {
      message = 'Sesión expirada o no autenticada. Inicia sesión nuevamente.'
    } else if (error.response?.status === 403) {
      message = 'No tienes permisos para ejecutar esta acción en MES.'
    } else if (error.response?.status === 404) {
      message = 'Recurso no encontrado en el módulo MES.'
    }

    throw new Error(message)
  }

  async listTasks(filters: ListMESFilters = {}): Promise<MESTask[]> {
    const params: any = {}
    if (filters.search) params.search = filters.search
    if (filters.is_active !== undefined) params.is_active = filters.is_active

    try {
      const response = await api.get(`${this.moduleBase}/tasks`, { params })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las tareas MES')
    }
  }

  async createTask(data: CreateTaskRequest): Promise<MESTask> {
    try {
      const response = await api.post(`${this.moduleBase}/tasks`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la tarea')
    }
  }

  async getTask(id: string): Promise<MESTask> {
    try {
      const response = await api.get(`${this.moduleBase}/tasks/${encodeURIComponent(id)}`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo cargar la tarea')
    }
  }

  async updateTask(id: string, data: UpdateTaskRequest): Promise<MESTask> {
    try {
      const response = await api.put(`${this.moduleBase}/tasks/${encodeURIComponent(id)}`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la tarea')
    }
  }

  async deleteTask(id: string): Promise<void> {
    try {
      await api.delete(`${this.moduleBase}/tasks/${encodeURIComponent(id)}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar la tarea')
    }
  }

  async listPositions(filters: ListMESFilters = {}): Promise<MESPosition[]> {
    const params: any = {}
    if (filters.search) params.search = filters.search
    if (filters.is_active !== undefined) params.is_active = filters.is_active

    try {
      const response = await api.get(`${this.moduleBase}/positions`, { params })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las posiciones MES')
    }
  }

  async createPosition(data: CreatePositionRequest): Promise<MESPosition> {
    try {
      const response = await api.post(`${this.moduleBase}/positions`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la posición')
    }
  }

  async getPosition(id: string): Promise<MESPosition> {
    try {
      const response = await api.get(`${this.moduleBase}/positions/${encodeURIComponent(id)}`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo cargar la posición')
    }
  }

  async updatePosition(id: string, data: UpdatePositionRequest): Promise<MESPosition> {
    try {
      const response = await api.put(`${this.moduleBase}/positions/${encodeURIComponent(id)}`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la posición')
    }
  }

  async deletePosition(id: string): Promise<void> {
    try {
      await api.delete(`${this.moduleBase}/positions/${encodeURIComponent(id)}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar la posición')
    }
  }

  async listWorkTypes(filters: ListMESFilters = {}): Promise<MESWorkType[]> {
    const params: any = {}
    if (filters.search) params.search = filters.search
    if (filters.is_active !== undefined) params.is_active = filters.is_active

    try {
      const response = await api.get(`${this.moduleBase}/work-types`, { params })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los tipos de trabajo MES')
    }
  }

  async createWorkType(data: CreateWorkTypeRequest): Promise<MESWorkType> {
    try {
      const response = await api.post(`${this.moduleBase}/work-types`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear el tipo de trabajo MES')
    }
  }

  async getWorkType(id: string): Promise<MESWorkType> {
    try {
      const response = await api.get(`${this.moduleBase}/work-types/${encodeURIComponent(id)}`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo cargar el tipo de trabajo')
    }
  }

  async updateWorkType(id: string, data: UpdateWorkTypeRequest): Promise<MESWorkType> {
    try {
      const response = await api.put(`${this.moduleBase}/work-types/${encodeURIComponent(id)}`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar el tipo de trabajo')
    }
  }

  async deleteWorkType(id: string): Promise<void> {
    try {
      await api.delete(`${this.moduleBase}/work-types/${encodeURIComponent(id)}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar el tipo de trabajo')
    }
  }

  async listWorkOrders(filters: ListWorkOrderFilters = {}): Promise<WorkOrder[]> {
    const params: any = {}
    if (filters.search) params.search = filters.search
    if (filters.status) params.status = filters.status
    if (filters.party_id) params.party_id = filters.party_id

    try {
      const response = await api.get(`${this.moduleBase}/work-orders`, { params })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las órdenes de trabajo MES')
    }
  }

  async getWorkOrder(id: string): Promise<WorkOrder> {
    try {
      const response = await api.get(`${this.moduleBase}/work-orders/${id}`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo cargar la orden de trabajo MES')
    }
  }

  async createWorkOrder(data: CreateWorkOrderRequest): Promise<WorkOrder> {
    try {
      const response = await api.post(`${this.moduleBase}/work-orders`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la orden de trabajo MES')
    }
  }

  async updateWorkOrder(id: string, data: UpdateWorkOrderRequest): Promise<WorkOrder> {
    try {
      const response = await api.put(`${this.moduleBase}/work-orders/${id}`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la orden de trabajo MES')
    }
  }

  async getWorkOrderDashboardStats(): Promise<WorkOrderDashboardStats> {
    try {
      const response = await api.get(`${this.moduleBase}/work-orders/dashboard/stats`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las estadísticas MES')
    }
  }

  async listOverdueWorkOrders(limit?: number): Promise<WorkOrder[]> {
    const params: any = {}
    if (limit && limit > 0) params.limit = limit

    try {
      const response = await api.get(`${this.moduleBase}/work-orders/overdue`, { params })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar los trabajos vencidos')
    }
  }

  async updateWorkOrderTaskStatus(workId: string, taskId: string, data: UpdateWorkOrderTaskStatusRequest): Promise<WorkOrder> {
    try {
      const response = await api.patch(`${this.moduleBase}/work-orders/${workId}/tasks/${taskId}/status`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar el estado de la tarea MES')
    }
  }

  // --- WorkSetup methods ---

  async listWorkSetups(filters: ListWorkSetupFilters = {}): Promise<WorkSetup[]> {
    const params: any = {}
    if (filters.search) params.search = filters.search
    if (filters.is_active !== undefined) params.is_active = filters.is_active
    if (filters.party_id) params.party_id = filters.party_id

    try {
      const response = await api.get(`${this.moduleBase}/work-setups`, { params })
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron cargar las configuraciones de cliente')
    }
  }

  async getWorkSetup(id: string): Promise<WorkSetup> {
    try {
      const response = await api.get(`${this.moduleBase}/work-setups/${encodeURIComponent(id)}`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo cargar la configuración de cliente')
    }
  }

  async createWorkSetup(data: CreateWorkSetupRequest): Promise<WorkSetup> {
    try {
      const response = await api.post(`${this.moduleBase}/work-setups`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo crear la configuración de cliente')
    }
  }

  async updateWorkSetup(id: string, data: UpdateWorkSetupRequest): Promise<WorkSetup> {
    try {
      const response = await api.put(`${this.moduleBase}/work-setups/${encodeURIComponent(id)}`, data)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudo actualizar la configuración de cliente')
    }
  }

  async deleteWorkSetup(id: string): Promise<void> {
    try {
      await api.delete(`${this.moduleBase}/work-setups/${encodeURIComponent(id)}`)
    } catch (e) {
      await this.handleError(e, 'No se pudo eliminar la configuración de cliente')
    }
  }

  // --- Pending WorkSetups (from confirmed Sales orders) ---

  async listPendingWorkSetups(): Promise<PendingWorkSetup[]> {
    try {
      const response = await api.get(`${this.moduleBase}/pending-work-setups`)
      return response.data
    } catch (e) {
      await this.handleError(e, 'No se pudieron obtener las configuraciones pendientes')
    }
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

