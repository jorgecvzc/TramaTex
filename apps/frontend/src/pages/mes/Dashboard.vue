<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Panel</p>
          <h1>Monitoreo de producción</h1>
          <p class="subtitle">Resumen de trabajos, estados y vencimientos.</p>
        </div>
        <div class="header-actions">
          <RouterLink to="/mes/terminal" class="btn btn-primary">Terminal de taller</RouterLink>
          <RouterLink to="/mes/work-orders" class="btn btn-secondary">Ver órdenes</RouterLink>
        </div>
      </header>

      <section class="stats-grid">
        <RouterLink to="/mes/work-orders" class="card stat-card stat-card-link">
          <p class="stat-label">Total trabajos</p>
          <p class="stat-value">{{ stats?.total ?? 0 }}</p>
        </RouterLink>
        <RouterLink to="/mes/work-orders" class="card stat-card stat-card-link">
          <p class="stat-label">Vencidos</p>
          <p class="stat-value danger">{{ stats?.overdue ?? 0 }}</p>
        </RouterLink>
        <RouterLink to="/mes/work-orders" class="card stat-card stat-card-link">
          <p class="stat-label">Vencen hoy</p>
          <p class="stat-value warning">{{ stats?.due_today ?? 0 }}</p>
        </RouterLink>
        <RouterLink to="/mes/work-orders?status=IN_PROGRESS" class="card stat-card stat-card-link">
          <p class="stat-label">En progreso</p>
          <p class="stat-value">{{ stats?.by_status?.IN_PROGRESS ?? 0 }}</p>
        </RouterLink>
        <RouterLink to="/mes/terminal?status=PENDING" class="card stat-card stat-card-link">
          <p class="stat-label">Tareas en cola</p>
          <p class="stat-value" :class="pendingTasksCount > 0 ? 'warning' : ''">{{ pendingTasksCount }}</p>
        </RouterLink>
        <a class="card stat-card stat-card-link" @click.prevent="openCompletedDialog">
          <p class="stat-label">Terminadas</p>
          <p class="stat-value success">{{ stats?.by_status?.COMPLETED ?? 0 }}</p>
        </a>
        <a class="card stat-card stat-card-link" @click.prevent="openCancelledDialog">
          <p class="stat-label">Canceladas</p>
          <p class="stat-value muted">{{ stats?.by_status?.CANCELLED ?? 0 }}</p>
        </a>
      </section>

      <!-- Dedicated PENDING section with conditional actions -->
      <section class="card">
        <div class="section-header">
          <h2>Pendientes <span class="status-badge badge-pending">{{ pendingOrders.length }}</span></h2>
          <RouterLink to="/mes/work-orders?status=PENDING" class="btn-link section-link">Ver todas →</RouterLink>
        </div>
        <div v-if="isLoadingByStatus" class="empty-state">Cargando órdenes...</div>
        <div v-else-if="pendingOrders.length === 0" class="empty-state">No hay órdenes pendientes.</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Número</th>
              <th>Nombre</th>
              <th>Pedido</th>
              <th>Prioridad</th>
              <th>Vencimiento</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="wo in pendingOrders" :key="wo.id">
              <td>
                <RouterLink :to="`/mes/work-orders/${wo.id}`" class="btn-link"><strong>{{ wo.work_number }}</strong></RouterLink>
              </td>
              <td>{{ wo.work_name }}</td>
              <td>
                <RouterLink v-if="wo.sales_order_id" :to="`/sales/orders/${wo.sales_order_id}`" class="btn-link">{{ wo.sales_order_number }}</RouterLink>
                <span v-else class="text-muted">—</span>
              </td>
              <td>{{ mesApi.getPriorityLabel(wo.priority) }}</td>
              <td :class="{ 'date-urgent': isDeliveryUrgent(wo.due_date) }">{{ formatDate(wo.due_date) }}</td>
              <td class="actions-cell">
                <template v-if="!wo.work_setup_id">
                  <button class="btn btn-sm btn-primary" @click="openSetupDialog(wo)">Definir configuración</button>
                </template>
                <template v-else>
                  <button class="btn btn-sm btn-success" @click="startWorkOrder(wo)">Iniciar</button>
                  <button class="btn btn-sm btn-warning" @click="suspendWorkOrder(wo)">Suspender</button>
                </template>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="card">
        <div class="section-header">
          <h2>Órdenes vencidas</h2>
          <RouterLink to="/mes/work-orders" class="btn-link section-link">Ver todas →</RouterLink>
        </div>
        <div v-if="isLoading" class="empty-state">Cargando panel...</div>
        <div v-else-if="error" class="alert">{{ error }}</div>
        <div v-else-if="overdueWorks.length === 0" class="empty-state">No hay órdenes vencidas.</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Número</th>
              <th>Nombre</th>
              <th>Pedido</th>
              <th>Estado</th>
              <th>Vencimiento</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="work in overdueWorks" :key="work.id">
              <td>
                <RouterLink :to="`/mes/work-orders/${work.id}`" class="btn-link"><strong>{{ work.work_number }}</strong></RouterLink>
              </td>
              <td>{{ work.work_name }}</td>
              <td>
                <RouterLink v-if="work.sales_order_id" :to="`/sales/orders/${work.sales_order_id}`" class="btn-link">{{ work.sales_order_number }}</RouterLink>
                <span v-else class="text-muted">—</span>
              </td>
              <td>{{ mesApi.getWorkStatusLabel(work.status) }}</td>
              <td>{{ formatDate(work.due_date) }}</td>
              <td>
                <RouterLink :to="`/mes/work-orders/${work.id}`" class="btn-link">Ver detalle</RouterLink>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <section v-for="section in statusSections" :key="section.status" class="card">
        <div class="section-header">
          <h2>{{ section.title }} <span class="status-badge" :class="section.badgeClass">{{ section.orders.length }}</span></h2>
          <RouterLink :to="`/mes/work-orders?status=${section.status}`" class="btn-link section-link">Ver todas →</RouterLink>
        </div>
        <div v-if="isLoadingByStatus" class="empty-state">Cargando órdenes...</div>
        <div v-else-if="section.orders.length === 0" class="empty-state">No hay órdenes {{ section.emptyLabel }}.</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Número</th>
              <th>Nombre</th>
              <th>Pedido</th>
              <th>Prioridad</th>
              <th>Vencimiento</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="wo in section.orders" :key="wo.id">
              <td>
                <RouterLink :to="`/mes/work-orders/${wo.id}`" class="btn-link"><strong>{{ wo.work_number }}</strong></RouterLink>
              </td>
              <td>{{ wo.work_name }}</td>
              <td>
                <RouterLink v-if="wo.sales_order_id" :to="`/sales/orders/${wo.sales_order_id}`" class="btn-link">{{ wo.sales_order_number }}</RouterLink>
                <span v-else class="text-muted">—</span>
              </td>
              <td>{{ mesApi.getPriorityLabel(wo.priority) }}</td>
              <td :class="{ 'date-urgent': isDeliveryUrgent(wo.due_date) }">{{ formatDate(wo.due_date) }}</td>
              <td>
                <RouterLink :to="`/mes/work-orders/${wo.id}`" class="btn-link">Ver detalle</RouterLink>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <!-- Dedicated SUSPENDED section with cancel action -->
      <section class="card">
        <div class="section-header">
          <h2>Suspendidas <span class="status-badge badge-suspended">{{ suspendedOrders.length }}</span></h2>
          <RouterLink to="/mes/work-orders?status=SUSPENDED" class="btn-link section-link">Ver todas →</RouterLink>
        </div>
        <div v-if="isLoadingByStatus" class="empty-state">Cargando órdenes...</div>
        <div v-else-if="suspendedOrders.length === 0" class="empty-state">No hay órdenes suspendidas.</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Número</th>
              <th>Nombre</th>
              <th>Pedido</th>
              <th>Prioridad</th>
              <th>Vencimiento</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="wo in suspendedOrders" :key="wo.id">
              <td>
                <RouterLink :to="`/mes/work-orders/${wo.id}`" class="btn-link"><strong>{{ wo.work_number }}</strong></RouterLink>
              </td>
              <td>{{ wo.work_name }}</td>
              <td>
                <RouterLink v-if="wo.sales_order_id" :to="`/sales/orders/${wo.sales_order_id}`" class="btn-link">{{ wo.sales_order_number }}</RouterLink>
                <span v-else class="text-muted">—</span>
              </td>
              <td>{{ mesApi.getPriorityLabel(wo.priority) }}</td>
              <td :class="{ 'date-urgent': isDeliveryUrgent(wo.due_date) }">{{ formatDate(wo.due_date) }}</td>
              <td class="actions-cell">
                <button class="btn btn-sm btn-pending" @click="reactivateWorkOrder(wo)">Reactivar</button>
                <button class="btn btn-sm btn-danger" @click="cancelWorkOrder(wo)">Cancelar</button>
              </td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>

    <!-- Completed Orders Dialog -->
    <div v-if="showCompletedDialog" class="modal-overlay" @click.self="showCompletedDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Órdenes terminadas</h3>
          <button @click="showCompletedDialog = false" class="btn-close">✕</button>
        </div>
        <div class="modal-body">
          <div v-if="isLoadingCompleted" class="empty-state">Cargando órdenes terminadas...</div>
          <div v-else-if="completedOrders.length === 0" class="empty-state">No hay órdenes terminadas.</div>
          <table v-else class="data-table">
            <thead>
              <tr>
                <th>Número</th>
                <th>Nombre</th>
                <th>Pedido</th>
                <th>Completada</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in completedOrders" :key="wo.id">
                <td><RouterLink :to="`/mes/work-orders/${wo.id}`" class="btn-link"><strong>{{ wo.work_number }}</strong></RouterLink></td>
                <td>{{ wo.work_name }}</td>
                <td>
                  <RouterLink v-if="wo.sales_order_id" :to="`/sales/orders/${wo.sales_order_id}`" class="btn-link">{{ wo.sales_order_number }}</RouterLink>
                  <span v-else class="text-muted">—</span>
                </td>
                <td>{{ formatDate(wo.completed_date) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Cancelled Orders Dialog -->
    <div v-if="showCancelledDialog" class="modal-overlay" @click.self="showCancelledDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Órdenes canceladas</h3>
          <button @click="showCancelledDialog = false" class="btn-close">✕</button>
        </div>
        <div class="modal-body">
          <div v-if="isLoadingCancelled" class="empty-state">Cargando órdenes canceladas...</div>
          <div v-else-if="cancelledOrders.length === 0" class="empty-state">No hay órdenes canceladas.</div>
          <table v-else class="data-table">
            <thead>
              <tr>
                <th>Número</th>
                <th>Nombre</th>
                <th>Pedido</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in cancelledOrders" :key="wo.id">
                <td><RouterLink :to="`/mes/work-orders/${wo.id}`" class="btn-link"><strong>{{ wo.work_number }}</strong></RouterLink></td>
                <td>{{ wo.work_name }}</td>
                <td>
                  <RouterLink v-if="wo.sales_order_id" :to="`/sales/orders/${wo.sales_order_id}`" class="btn-link">{{ wo.sales_order_number }}</RouterLink>
                  <span v-else class="text-muted">—</span>
                </td>
                <td class="actions-cell">
                  <button class="btn btn-sm btn-pending" @click="reactivateCancelledOrder(wo)">Reactivar</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- WorkSetup Assignment Dialog -->
    <div v-if="showSetupDialog" class="modal-overlay" @click.self="closeSetupDialog">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Configuración de trabajo — {{ setupDialogOrder?.work_number }}</h3>
          <button @click="closeSetupDialog" class="btn-close">✕</button>
        </div>

        <div class="modal-body">
          <div v-if="setupDialogError" class="alert">{{ setupDialogError }}</div>

          <!-- WorkOrder notes as reference -->
          <div v-if="setupDialogOrder?.notes" class="info-banner">
            <strong>Observaciones del pedido:</strong> {{ setupDialogOrder.notes }}
          </div>

          <!-- Mode tabs -->
          <div class="tab-bar">
            <button :class="['tab-btn', { active: setupDialogMode === 'assign' }]" @click="setupDialogMode = 'assign'">Asignar existente</button>
            <button :class="['tab-btn', { active: setupDialogMode === 'create' }]" @click="switchToCreateMode">Crear nueva</button>
          </div>

          <!-- Assign existing mode -->
          <div v-if="setupDialogMode === 'assign'">
            <div v-if="isLoadingSetups" class="empty-state">Cargando configuraciones...</div>
            <div v-else-if="existingWorkSetups.length === 0" class="empty-state">No hay configuraciones activas para este cliente.</div>
            <div v-else class="setup-list">
              <label v-for="ws in existingWorkSetups" :key="ws.id" class="setup-option" :class="{ selected: selectedWorkSetupId === ws.id }">
                <input type="radio" v-model="selectedWorkSetupId" :value="ws.id" name="work_setup" />
                <div class="setup-info">
                  <strong>{{ ws.name }}</strong>
                  <span v-if="ws.description" class="setup-desc">{{ ws.description }}</span>
                  <span class="setup-lines">{{ ws.lines.length }} línea(s)</span>
                </div>
              </label>
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="closeSetupDialog">Cancelar</button>
              <button class="btn btn-primary" :disabled="!selectedWorkSetupId || setupDialogSaving" @click="assignExistingSetup">
                {{ setupDialogSaving ? 'Asignando...' : 'Asignar' }}
              </button>
            </div>
          </div>

          <!-- Create new mode -->
          <div v-if="setupDialogMode === 'create'" class="create-form">
            <label class="label">Nombre *</label>
            <input v-model="newSetupForm.name" type="text" class="input" placeholder="Ej: Uniformes Empresa XYZ" />

            <PartySelector
              v-model="newSetupForm.party_id"
              label="Cliente *"
              placeholder="Buscar cliente por nombre..."
              role-filter="CLIENT"
              :required="true"
            />

            <label class="label">Categoría tangible *</label>
            <select v-model="newSetupForm.tangible_group_id" class="input">
              <option value="">Seleccionar categoría tangible</option>
              <option v-for="group in tangibleGroups" :key="group.id" :value="group.id">
                {{ group.name }}
              </option>
            </select>

            <label class="label">Descripción</label>
            <textarea v-model="newSetupForm.description" class="input" rows="2" placeholder="Descripción opcional" />

            <div class="lines-block">
              <div class="lines-header">
                <h4>Líneas de configuración</h4>
                <button @click="addSetupLine" type="button" class="btn btn-secondary btn-sm">+ Añadir línea</button>
              </div>
              <div v-for="(line, index) in newSetupForm.lines" :key="index" class="line-row">
                <select v-model="line.work_type_id" class="input">
                  <option value="">Tipo de trabajo</option>
                  <option v-for="wt in workTypes" :key="wt.id" :value="wt.id">{{ wt.name }}</option>
                </select>
                <select v-model="line.position_id" class="input">
                  <option value="">Posición</option>
                  <option v-for="pos in positions" :key="pos.id" :value="pos.id">{{ pos.name }} ({{ pos.code }})</option>
                </select>
                <input v-model.number="line.sequence" type="number" min="1" class="input seq" placeholder="Seq" />
                <input v-model="line.design_file_path" type="text" class="input file-path-input" placeholder="Ruta del archivo (opcional)" title="Ruta completa del archivo de diseño" />
                <button type="button" class="btn btn-secondary btn-sm file-pick-btn" title="Seleccionar archivo" @click="pickFile(index)">📂</button>
                <input :ref="el => { fileInputRefs[index] = el as HTMLInputElement }" type="file" class="file-input-hidden" @change="onFileSelected($event, index)" />
                <button @click="removeSetupLine(index)" type="button" class="btn btn-danger btn-sm">Quitar</button>
              </div>
              <p v-if="newSetupForm.lines.length === 0" class="empty-state">Sin líneas. Añade al menos una.</p>
            </div>

            <div class="modal-footer">
              <button class="btn btn-secondary" @click="closeSetupDialog">Cancelar</button>
              <button class="btn btn-primary" :disabled="setupDialogSaving" @click="createAndAssignSetup">
                {{ setupDialogSaving ? 'Creando...' : 'Crear y asignar' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import PartySelector from '@/components/party/PartySelector.vue'
import { mesApi } from '@/services/mesApi'
import { productApi } from '@/services/productApi'
import type { WorkOrder, WorkOrderDashboardStats, WorkSetup, MESWorkType, MESPosition } from '@/types/mes'

const isLoading = ref(false)
const isLoadingByStatus = ref(false)
const error = ref('')
const stats = ref<WorkOrderDashboardStats | null>(null)
const overdueWorks = ref<WorkOrder[]>([])

const suspendedOrders = ref<WorkOrder[]>([])
const pendingOrders = ref<WorkOrder[]>([])
const inProgressOrders = ref<WorkOrder[]>([])
const onHoldOrders = ref<WorkOrder[]>([])

const statusSections = computed(() => [
  { status: 'IN_PROGRESS', title: 'En progreso', emptyLabel: 'en progreso', badgeClass: 'badge-progress', orders: inProgressOrders.value },
  { status: 'ON_HOLD', title: 'En espera', emptyLabel: 'en espera', badgeClass: 'badge-hold', orders: onHoldOrders.value },
])

// Count of PENDING tasks across all IN_PROGRESS work orders (visible in the terminal)
const pendingTasksCount = computed(() => {
  let count = 0
  for (const wo of inProgressOrders.value) {
    for (const line of wo.lines || []) {
      for (const task of line.tasks || []) {
        if (task.status === 'PENDING') count++
      }
    }
  }
  return count
})

// --- Completed / Cancelled dialog state ---
const showCompletedDialog = ref(false)
const showCancelledDialog = ref(false)
const completedOrders = ref<WorkOrder[]>([])
const cancelledOrders = ref<WorkOrder[]>([])
const isLoadingCompleted = ref(false)
const isLoadingCancelled = ref(false)

async function openCompletedDialog() {
  showCompletedDialog.value = true
  isLoadingCompleted.value = true
  try {
    completedOrders.value = await mesApi.listWorkOrders({ status: 'COMPLETED' })
  } catch (err: any) {
    console.error('Error loading completed orders:', err)
    completedOrders.value = []
  } finally {
    isLoadingCompleted.value = false
  }
}

async function openCancelledDialog() {
  showCancelledDialog.value = true
  isLoadingCancelled.value = true
  try {
    cancelledOrders.value = await mesApi.listWorkOrders({ status: 'CANCELLED' })
  } catch (err: any) {
    console.error('Error loading cancelled orders:', err)
    cancelledOrders.value = []
  } finally {
    isLoadingCancelled.value = false
  }
}

async function reactivateWorkOrder(wo: WorkOrder) {
  try {
    await mesApi.updateWorkOrder(wo.id, { status: 'PENDING' })
    await loadOrdersByStatus()
  } catch (err: any) {
    console.error('Error reactivating work order:', err)
  }
}

async function reactivateCancelledOrder(wo: WorkOrder) {
  try {
    await mesApi.updateWorkOrder(wo.id, { status: 'PENDING' })
    cancelledOrders.value = cancelledOrders.value.filter(o => o.id !== wo.id)
    await Promise.all([loadOrdersByStatus(), loadDashboard()])
  } catch (err: any) {
    console.error('Error reactivating cancelled order:', err)
  }
}

// --- WorkSetup assignment dialog state ---
type ProductGroupOption = { id: string; name: string; type: string }

const showSetupDialog = ref(false)
const setupDialogOrder = ref<WorkOrder | null>(null)
const setupDialogMode = ref<'assign' | 'create'>('assign')
const setupDialogError = ref('')
const setupDialogSaving = ref(false)

// Assign mode state
const existingWorkSetups = ref<WorkSetup[]>([])
const selectedWorkSetupId = ref('')
const isLoadingSetups = ref(false)

// Create mode state
const workTypes = ref<MESWorkType[]>([])
const positions = ref<MESPosition[]>([])
const productGroups = ref<ProductGroupOption[]>([])
const tangibleGroups = computed(() => productGroups.value.filter(g => g.type === 'TANGIBLE'))
const newSetupForm = ref({
  name: '',
  party_id: '',
  tangible_group_id: '',
  description: '',
  is_active: true,
  lines: [] as Array<{ work_type_id: string; position_id: string; sequence: number; design_file_path: string }>,
})

function openSetupDialog(wo: WorkOrder) {
  setupDialogOrder.value = wo
  setupDialogMode.value = 'assign'
  setupDialogError.value = ''
  setupDialogSaving.value = false
  selectedWorkSetupId.value = ''
  existingWorkSetups.value = []
  showSetupDialog.value = true
  loadExistingSetups(wo.party_id)
}

function closeSetupDialog() {
  showSetupDialog.value = false
  setupDialogOrder.value = null
}

async function loadExistingSetups(partyId: string) {
  isLoadingSetups.value = true
  try {
    existingWorkSetups.value = await mesApi.listWorkSetups({ party_id: partyId, is_active: true })
  } catch (err: any) {
    console.error('Error loading work setups:', err)
    existingWorkSetups.value = []
  } finally {
    isLoadingSetups.value = false
  }
}

async function loadCreateFormOptions() {
  try {
    const [wt, pos, groups] = await Promise.all([
      mesApi.listWorkTypes({ is_active: true }),
      mesApi.listPositions({ is_active: true }),
      productApi.listProductGroups({ isActive: true }),
    ])
    workTypes.value = wt
    positions.value = pos
    productGroups.value = groups.data
  } catch (err: any) {
    setupDialogError.value = err.message || 'Error cargando opciones del formulario'
  }
}

function switchToCreateMode() {
  setupDialogMode.value = 'create'
  const wo = setupDialogOrder.value
  newSetupForm.value = {
    name: '',
    party_id: wo?.party_id || '',
    tangible_group_id: '',
    description: '',
    is_active: true,
    lines: [],
  }
  if (workTypes.value.length === 0) {
    loadCreateFormOptions()
  }
}

function addSetupLine() {
  newSetupForm.value.lines.push({ work_type_id: '', position_id: '', sequence: newSetupForm.value.lines.length + 1, design_file_path: '' })
}

function removeSetupLine(index: number) {
  newSetupForm.value.lines.splice(index, 1)
}

const fileInputRefs = ref<HTMLInputElement[]>([])

function pickFile(index: number) {
  fileInputRefs.value[index]?.click()
}

function onFileSelected(event: Event, index: number) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    newSetupForm.value.lines[index].design_file_path = file.name
  }
}

async function assignExistingSetup() {
  if (!selectedWorkSetupId.value || !setupDialogOrder.value) return
  setupDialogError.value = ''
  setupDialogSaving.value = true
  try {
    await mesApi.updateWorkOrder(setupDialogOrder.value.id, { work_setup_id: selectedWorkSetupId.value })
    closeSetupDialog()
    await loadOrdersByStatus()
  } catch (err: any) {
    setupDialogError.value = err.message || 'Error al asignar configuración'
  } finally {
    setupDialogSaving.value = false
  }
}

async function createAndAssignSetup() {
  const wo = setupDialogOrder.value
  if (!wo) return
  setupDialogError.value = ''

  const f = newSetupForm.value
  if (!f.name.trim()) { setupDialogError.value = 'El nombre es obligatorio'; return }
  if (!f.party_id.trim()) { setupDialogError.value = 'El cliente es obligatorio'; return }
  if (!f.tangible_group_id.trim()) { setupDialogError.value = 'La categoría tangible es obligatoria'; return }
  const invalidLine = f.lines.find(l => !l.work_type_id || !l.position_id || l.sequence < 1)
  if (invalidLine) { setupDialogError.value = 'Todas las líneas deben tener tipo de trabajo, posición y secuencia válida'; return }
  if (f.lines.length === 0) { setupDialogError.value = 'Debe añadir al menos una línea'; return }

  setupDialogSaving.value = true
  try {
    const created = await mesApi.createWorkSetup({
      name: f.name.trim(),
      party_id: f.party_id.trim(),
      tangible_group_id: f.tangible_group_id.trim(),
      description: f.description.trim() || undefined,
      is_active: f.is_active,
      lines: f.lines.map(l => ({
        work_type_id: l.work_type_id,
        position_id: l.position_id,
        sequence: l.sequence,
        design_file_path: l.design_file_path || undefined,
      })),
    })
    await mesApi.updateWorkOrder(wo.id, { work_setup_id: created.id })
    closeSetupDialog()
    await loadOrdersByStatus()
  } catch (err: any) {
    setupDialogError.value = err.message || 'Error al crear y asignar configuración'
  } finally {
    setupDialogSaving.value = false
  }
}

async function startWorkOrder(wo: WorkOrder) {
  try {
    await mesApi.updateWorkOrder(wo.id, { status: 'IN_PROGRESS' })
    await loadOrdersByStatus()
  } catch (err: any) {
    console.error('Error starting work order:', err)
  }
}

async function suspendWorkOrder(wo: WorkOrder) {
  try {
    await mesApi.updateWorkOrder(wo.id, { status: 'SUSPENDED' })
    await loadOrdersByStatus()
  } catch (err: any) {
    console.error('Error suspending work order:', err)
  }
}

async function cancelWorkOrder(wo: WorkOrder) {
  try {
    await mesApi.updateWorkOrder(wo.id, { status: 'CANCELLED' })
    await loadOrdersByStatus()
  } catch (err: any) {
    console.error('Error cancelling work order:', err)
  }
}

function formatDate(value?: string) {
  if (!value) {
    return '—'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleDateString('es-ES')
}

function isDeliveryUrgent(dateStr?: string): boolean {
  if (!dateStr) return false
  const date = new Date(dateStr)
  const now = new Date()
  const diffDays = (date.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
  return diffDays <= 3
}

async function loadDashboard() {
  isLoading.value = true
  error.value = ''

  try {
    const [statsResult, overdueResult] = await Promise.all([
      mesApi.getWorkOrderDashboardStats(),
      mesApi.listOverdueWorkOrders(20),
    ])

    stats.value = statsResult
    overdueWorks.value = overdueResult
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar el panel'
  } finally {
    isLoading.value = false
  }
}

async function loadOrdersByStatus() {
  isLoadingByStatus.value = true
  try {
    const [suspended, pending, inProgress, onHold] = await Promise.all([
      mesApi.listWorkOrders({ status: 'SUSPENDED' }),
      mesApi.listWorkOrders({ status: 'PENDING' }),
      mesApi.listWorkOrders({ status: 'IN_PROGRESS' }),
      mesApi.listWorkOrders({ status: 'ON_HOLD' }),
    ])
    suspendedOrders.value = suspended
    pendingOrders.value = pending
    inProgressOrders.value = inProgress
    onHoldOrders.value = onHold
  } catch (err: any) {
    console.error('Error loading orders by status:', err)
  } finally {
    isLoadingByStatus.value = false
  }
}

onMounted(() => {
  loadDashboard()
  loadOrdersByStatus()
})
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1200px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.header-actions { display: flex; gap: .5rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: .75rem; }
.stat-card { display: flex; flex-direction: column; gap: .25rem; }
.stat-label { margin: 0; color: #64748b; font-size: .85rem; }
.stat-value { margin: 0; font-size: 1.8rem; font-weight: 700; color: #1e293b; }
.stat-value.danger { color: #b91c1c; }
.stat-value.warning { color: #b45309; }
.stat-value.success { color: #065f46; }
.stat-value.muted { color: #94a3b8; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.btn-sm { padding: .35rem .75rem; font-size: .8rem; }
.btn-disabled { background: #e2e8f0; color: #94a3b8; cursor: not-allowed; pointer-events: none; }
.btn-link { color: #1d4ed8; text-decoration: none; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: .75rem; border-bottom: 1px solid #e2e8f0; text-align: left; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
.section-subtitle { color: #64748b; font-size: .85rem; margin: 0 0 .75rem; }
.pending-badge { display: inline-block; padding: .15rem .5rem; border-radius: 9999px; font-size: .75rem; font-weight: 500; background: #fef3c7; color: #92400e; }
.date-urgent { color: #b91c1c; font-weight: 600; }
.text-muted { color: #94a3b8; }
.observations-cell { max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #6b7280; font-size: .85rem; }
.stat-card-link { text-decoration: none; transition: box-shadow .15s, transform .15s; }
.stat-card-link:hover { box-shadow: 0 2px 8px rgba(0,0,0,.1); transform: translateY(-2px); }
.section-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: .5rem; }
.section-header h2 { margin: 0; display: flex; align-items: center; gap: .5rem; }
.section-link { font-size: .85rem; white-space: nowrap; padding-top: .25rem; }
.status-badge { display: inline-flex; align-items: center; justify-content: center; min-width: 1.6rem; height: 1.6rem; border-radius: 9999px; font-size: .75rem; font-weight: 600; }
.badge-suspended { background: #fef3c7; color: #92400e; }
.badge-pending { background: #e0e7ff; color: #3730a3; }
.badge-progress { background: #d1fae5; color: #065f46; }
.badge-hold { background: #fee2e2; color: #991b1b; }
@media (max-width: 900px) { .stats-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }

/* Action buttons */
.actions-cell { display: flex; gap: .4rem; flex-wrap: wrap; }
.btn-success { background: #d1fae5; color: #065f46; border: 1px solid #a7f3d0; font-weight: 600; }
.btn-warning { background: #fef3c7; color: #92400e; border: 1px solid #fde68a; font-weight: 600; }
.btn-danger { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; }
.btn-pending { background: #e0e7ff; color: #3730a3; border: 1px solid #c7d2fe; font-weight: 600; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.4); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: #fff; border-radius: 12px; width: 90vw; max-width: 1400px; max-height: 88vh; overflow-y: auto; box-shadow: 0 8px 32px rgba(0,0,0,.15); }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 1rem 1.25rem; border-bottom: 1px solid #e2e8f0; }
.modal-header h3 { margin: 0; font-size: 1.05rem; }
.btn-close { background: none; border: none; font-size: 1.2rem; cursor: pointer; color: #64748b; padding: .25rem .5rem; }
.modal-body { padding: 1.25rem; display: flex; flex-direction: column; gap: .75rem; }
.modal-footer { display: flex; justify-content: flex-end; gap: .5rem; padding-top: .75rem; border-top: 1px solid #e2e8f0; margin-top: .5rem; }

/* Info banner */
.info-banner { background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 8px; padding: .75rem; color: #1e40af; font-size: .85rem; }

/* Tab bar */
.tab-bar { display: flex; gap: 0; border-bottom: 2px solid #e2e8f0; }
.tab-btn { background: none; border: none; padding: .5rem 1rem; font: inherit; cursor: pointer; color: #64748b; border-bottom: 2px solid transparent; margin-bottom: -2px; }
.tab-btn.active { color: #1d4ed8; border-bottom-color: #1d4ed8; font-weight: 600; }

/* Setup list (assign mode) */
.setup-list { display: flex; flex-direction: column; gap: .5rem; }
.setup-option { display: flex; align-items: flex-start; gap: .75rem; padding: .75rem; border: 1px solid #e2e8f0; border-radius: 8px; cursor: pointer; transition: border-color .15s; }
.setup-option.selected { border-color: #1d4ed8; background: #eff6ff; }
.setup-option input[type="radio"] { margin-top: .15rem; }
.setup-info { display: flex; flex-direction: column; gap: .15rem; }
.setup-desc { font-size: .8rem; color: #64748b; }
.setup-lines { font-size: .75rem; color: #94a3b8; }

/* Create form */
.create-form { display: flex; flex-direction: column; gap: .75rem; }
.create-form .label { font-size: .85rem; color: #334155; font-weight: 600; margin: 0; }
.create-form .input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .5rem .75rem; font: inherit; }

/* Lines block */
.lines-block { border: 1px solid #e2e8f0; border-radius: 8px; padding: .75rem; display: flex; flex-direction: column; gap: .5rem; }
.lines-header { display: flex; justify-content: space-between; align-items: center; }
.lines-header h4 { margin: 0; font-size: .9rem; color: #1e293b; }
.line-row { display: grid; grid-template-columns: 1fr 1fr 70px 1fr auto auto auto; gap: .4rem; align-items: center; }
.file-pick-btn { padding: .4rem .6rem; font-size: .85rem; }
.file-input-hidden { display: none; }
.seq { text-align: center; }
</style>
