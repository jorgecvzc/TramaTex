export interface MESTask {
  id: string
  name: string
  description?: string
  is_active: boolean
}

export interface MESPosition {
  id: string
  name: string
  code: string
  description?: string
  is_active: boolean
}

export interface MESServiceGroupTask {
  task_id: string
  sequence: number
}

export type MESServiceTemplateTask = MESServiceGroupTask

export interface MESServiceGroup {
  id: string
  name: string
  description?: string
  product_group_id?: string
  is_active: boolean
  tasks: MESServiceGroupTask[]
}

export type MESServiceTemplate = MESServiceGroup

export interface ListMESFilters {
  is_active?: boolean
  search?: string
}

export interface CreateTaskRequest {
  name: string
  description?: string
  is_active?: boolean
}

export interface CreatePositionRequest {
  name: string
  code: string
  description?: string
  is_active?: boolean
}

export interface ServiceGroupTaskInput {
  task_id: string
  sequence: number
}

export type ServiceTemplateTaskInput = ServiceGroupTaskInput

export interface CreateServiceGroupRequest {
  name: string
  description?: string
  product_group_id?: string
  is_active?: boolean
  task_assignments: ServiceGroupTaskInput[]
}

export type CreateServiceTemplateRequest = CreateServiceGroupRequest

export interface MESWorkTask {
  id: string
  task_id: string
  sequence: number
  status: string
  assigned_to?: string
  started_at?: string
  completed_at?: string
  notes?: string
}

export interface MESWorkServiceGroup {
  id: string
  service_group_id: string
  position_id: string
  design_file_path?: string
  notes?: string
  sequence: number
  tasks: MESWorkTask[]
}

export type MESWorkServiceTemplate = MESWorkServiceGroup

export interface MESWork {
  id: string
  work_number: string
  work_name: string
  party_id: string
  tangible_group_id: string
  garment_notes?: string
  status: string
  priority: string
  start_date?: string
  due_date?: string
  completed_date?: string
  service_groups: MESWorkServiceGroup[]
}

export type MESWorkDefinition = MESWork
export type MESWorkExecution = MESWork

export interface CreateMESWorkServiceGroupInput {
  service_group_id: string
  position_id: string
  design_file_path?: string
  notes?: string
  sequence: number
}

export type CreateMESWorkServiceTemplateInput = CreateMESWorkServiceGroupInput

export interface CreateMESWorkRequest {
  work_name: string
  party_id: string
  tangible_group_id: string
  garment_notes?: string
  status?: string
  priority?: string
  service_group_assignments: CreateMESWorkServiceGroupInput[]
}

export type CreateMESWorkDefinitionRequest = CreateMESWorkRequest

export interface UpdateMESWorkRequest {
  work_name?: string
  party_id?: string
  tangible_group_id?: string
  garment_notes?: string
  status?: string
  priority?: string
  due_date?: string
}

export type UpdateMESWorkDefinitionRequest = UpdateMESWorkRequest

export interface ListMESWorkFilters {
  status?: string
  search?: string
  party_id?: string
}

export interface MESWorkDashboardStats {
  total: number
  by_status: Record<string, number>
  overdue: number
  due_today: number
}

export type MESWorkTaskAction = 'START' | 'PAUSE' | 'COMPLETE' | 'BLOCK'

export interface UpdateMESWorkTaskStatusRequest {
  action: MESWorkTaskAction
  notes?: string
}
