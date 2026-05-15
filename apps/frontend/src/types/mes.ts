export interface MESTask {
  id: string
  name: string
  reference?: string
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

export interface MESWorkTypeTask {
  task_id: string
  sequence: number
}

export interface MESWorkType {
  id: string
  name: string
  reference?: string
  description?: string
  is_active: boolean
  tasks: MESWorkTypeTask[]
}

export interface ListMESFilters {
  is_active?: boolean
  search?: string
}

export interface CreateTaskRequest {
  name: string
  reference?: string
  description?: string
  is_active?: boolean
}

export interface UpdateTaskRequest {
  name?: string
  reference?: string
  description?: string
  is_active?: boolean
}

export interface CreatePositionRequest {
  name: string
  code: string
  description?: string
  is_active?: boolean
}

export interface UpdatePositionRequest {
  name?: string
  code?: string
  description?: string
  is_active?: boolean
}

export interface WorkTypeTaskInput {
  task_id: string
  sequence: number
}

export interface CreateWorkTypeRequest {
  name: string
  reference?: string
  description?: string
  is_active?: boolean
  task_assignments: WorkTypeTaskInput[]
}

export interface UpdateWorkTypeRequest {
  name?: string
  reference?: string
  description?: string
  is_active?: boolean
  task_assignments?: WorkTypeTaskInput[]
}

export interface WorkOrderTask {
  id: string
  task_id: string
  sequence: number
  status: string
  assigned_to?: string // Post-MVP: task assignment to operators
  started_at?: string
  completed_at?: string
  notes?: string
}

export interface WorkOrderLine {
  id: string
  work_type_id: string
  position_id: string
  design_file_path?: string
  notes?: string
  sequence: number
  tasks: WorkOrderTask[]
}

export interface WorkOrder {
  id: string
  work_number: string
  work_name: string
  party_id: string
  work_setup_id: string
  notes?: string
  status: string
  priority: string
  start_date?: string
  due_date?: string
  completed_date?: string
  lines: WorkOrderLine[]
  sales_order_id?: string
  sales_order_number?: string
}

export interface CreateWorkOrderRequest {
  work_name: string
  party_id: string
  work_setup_id: string
  notes?: string
  priority?: string
  due_date?: string
  order_work_setup_id?: string
}

export interface UpdateWorkOrderRequest {
  work_name?: string
  notes?: string
  status?: string
  priority?: string
  due_date?: string
  work_setup_id?: string
}

export interface ListWorkOrderFilters {
  status?: string
  search?: string
  party_id?: string
}

export interface WorkOrderDashboardStats {
  total: number
  by_status: Record<string, number>
  overdue: number
  due_today: number
}

export type WorkOrderTaskAction = 'START' | 'PAUSE' | 'COMPLETE' | 'BLOCK'

export interface UpdateWorkOrderTaskStatusRequest {
  action: WorkOrderTaskAction
  notes?: string
}

// --- WorkSetup types ---

export interface WorkSetupLine {
  id: string
  work_type_id: string
  position_id: string
  design_file_path?: string
  notes?: string
  sequence: number
}

export interface WorkSetup {
  id: string
  name: string
  party_id: string
  tangible_group_id: string
  description?: string
  is_active: boolean
  lines: WorkSetupLine[]
}

export interface WorkSetupLineInput {
  work_type_id: string
  position_id: string
  design_file_path?: string
  notes?: string
  sequence: number
}

export interface CreateWorkSetupRequest {
  name: string
  party_id: string
  tangible_group_id: string
  description?: string
  is_active?: boolean
  lines: WorkSetupLineInput[]
}

export interface UpdateWorkSetupRequest {
  name?: string
  party_id?: string
  tangible_group_id?: string
  description?: string
  is_active?: boolean
  lines?: WorkSetupLineInput[]
}

export interface ListWorkSetupFilters {
  is_active?: boolean
  search?: string
  party_id?: string
}

// --- Pending WorkSetup (from confirmed Sales orders) ---

export interface PendingWorkSetup {
  id: string
  work_setup_id?: string | null
  description: string
  order_id: string
  order_number: string
  delivery_date: string
  party_id: string
}
