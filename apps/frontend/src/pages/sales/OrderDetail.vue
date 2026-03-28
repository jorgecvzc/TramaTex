<template>
  <BaseEntityPage :is-loading="isLoading" :error="error">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <PageHeader 
        :title="pageTitle" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/dashboard' }, { label: 'Pedidos', to: '/sales/orders' }, { label: headerLabel }]"
      >
        <template #icon>
          <span class="material-symbols-outlined">shopping_cart</span>
        </template>
        <template #actions>
          <div v-if="order" class="header-actions-group">
            <button class="btn btn-outline" @click="printOrder">
              <span class="material-symbols-outlined">print</span>
              <span>Imprimir</span>
            </button>
            <button v-if="mode === 'detail'" class="btn btn-primary" @click="enterEditMode">
              <span class="material-symbols-outlined">edit</span>
              <span>Editar</span>
            </button>
            <button v-else class="btn btn-primary" @click="saveOrder" :disabled="isSaving">
              <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
              <span>{{ isSaving ? 'Guardando...' : 'Guardar' }}</span>
            </button>
          </div>
        </template>
      </PageHeader>
    </template>
    
    <!-- CAPA 2: CONTEXTO -->
    <template #summary v-if="order && mode === 'detail'">
      <div class="summary-grid">
        <div class="summary-item">
          <label>Total Pedido</label>
          <strong class="text-lg text-primary">{{ formatMoney(order.totalAmount) }}</strong>
        </div>
        <div class="summary-item">
          <label>Estado</label>
          <span :class="['status-badge', statusClass]">{{ statusLabel }}</span>
        </div>
        <div class="summary-item">
          <label>Cliente</label>
          <RouterLink :to="`/parties/${order.party.id}`" class="link-primary">{{ order.party.name }}</RouterLink>
        </div>
        <div class="summary-item">
          <label>Fecha Pedido</label>
          <span>{{ formatDate(order.orderDate) }}</span>
        </div>
      </div>
    </template>
    
    <!-- CAPA 3: TRABAJO -->
    <div v-if="order" class="order-detail-content">
      <!-- Sección Cliente y Condiciones -->
      <FormSection title="Cliente y Condiciones" icon="groups">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <DataRow label="Cliente" icon="person">
            <template v-if="mode !== 'detail'">
              <select v-model="editableOrder.partyId" class="form-input">
                <option v-for="party in parties" :key="party.id" :value="party.id">{{ party.name }}</option>
              </select>
            </template>
            <template v-else>{{ order.party.name }} <code class="ml-2">#{{ order.party.id.substring(0,8) }}</code></template>
          </DataRow>
          <DataRow label="Fecha del Pedido" icon="event">
            <template v-if="mode !== 'detail'">
              <input type="date" v-model="editableOrder.orderDate" class="form-input" />
            </template>
            <template v-else>{{ formatDate(order.orderDate) }}</template>
          </DataRow>
          <DataRow label="Referencia Cliente" icon="badge">
            <template v-if="mode !== 'detail'">
              <input v-model="editableOrder.partyReference" class="form-input" placeholder="Ej: PO-12345" />
            </template>
            <template v-else>{{ order.partyReference || '—' }}</template>
          </DataRow>
        </div>
      </FormSection>
      
      <!-- Sección Líneas de Pedido -->
      <FormSection title="Líneas de Pedido" icon="list">
        <OrderLines
          :lines="editableOrder.lines"
          :is-editing="mode !== 'detail'"
          @update:lines="updateLines"
        />
      </FormSection>
      
      <!-- Sección de Totales -->
      <div class="totals-section">
        <div class="total-row">
          <span class="label">Subtotal</span>
          <span class="value">{{ formatMoney(subtotal) }}</span>
        </div>
        <div class="total-row">
          <span class="label">IVA ({{ order.taxRate }}%)</span>
          <span class="value">{{ formatMoney(taxAmount) }}</span>
        </div>
        <div class="total-row grand-total">
          <span class="label">Total Pedido</span>
          <span class="value">{{ formatMoney(totalAmount) }}</span>
        </div>
      </div>
    </div>
  </BaseEntityPage>
  
  <PrintDocument ref="printComponent" :document="order" type="Pedido" />
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue'
import PageHeader from '@/components/layout/PageHeader.vue'
import FormSection from '@/components/shared/FormSection.vue'
import DataRow from '@/components/shared/DataRow.vue'
import OrderLines from '@/components/sales/OrderLines.vue'
import PrintDocument from '@/components/sales/PrintDocument.vue'
import salesApi from '@/services/salesApi'
import { partyApi } from '@/services/partyApi'

const route = useRoute()
const router = useRouter()
const order = ref(null)
const editableOrder = ref({ lines: [] })
const parties = ref([])
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref('')
const mode = ref('detail') // 'detail' | 'edit' | 'create'
const printComponent = ref(null)

const orderId = computed(() => route.params.id)
const isCreateMode = computed(() => !orderId.value)
const pageTitle = computed(() => {
  if (isCreateMode.value) return 'Nuevo Pedido'
  return mode.value === 'edit' ? `Editando Pedido #${order.value?.orderNumber}` : `Pedido #${order.value?.orderNumber}`
})
const headerLabel = computed(() => {
  if (isCreateMode.value) return 'Nuevo'
  return `#${order.value?.orderNumber}`
})

const statusLabel = computed(() => salesApi.getStatusLabel(order.value?.status))
const statusClass = computed(() => salesApi.getStatusClass(order.value?.status))

const subtotal = computed(() => editableOrder.value.lines.reduce((acc, line) => acc + (line.quantity * line.unitPrice), 0))
const taxAmount = computed(() => subtotal.value * (order.value?.taxRate / 100 || 0.21))
const totalAmount = computed(() => subtotal.value + taxAmount.value)

async function loadOrder() {
  if (isCreateMode.value) {
    mode.value = 'create'
    order.value = salesApi.getNewOrder()
    editableOrder.value = JSON.parse(JSON.stringify(order.value))
    await loadParties()
    return
  }

  isLoading.value = true
  try {
    const data = await salesApi.getOrder(orderId.value)
    order.value = data
    editableOrder.value = JSON.parse(JSON.stringify(data))
    mode.value = 'detail'
  } catch (e) {
    error.value = 'No se pudo cargar el pedido.'
  } finally {
    isLoading.value = false
  }
}

async function loadParties() {
  parties.value = await partyApi.listParties({ isActive: true })
}

function enterEditMode() {
  loadParties()
  mode.value = 'edit'
}

async function saveOrder() {
  isSaving.value = true
  try {
    const payload = { ...editableOrder.value }
    if (isCreateMode.value) {
      const newOrder = await salesApi.createOrder(payload)
      router.push(`/sales/orders/${newOrder.id}`)
    } else {
      const updatedOrder = await salesApi.updateOrder(orderId.value, payload)
      order.value = updatedOrder
      editableOrder.value = JSON.parse(JSON.stringify(updatedOrder))
      mode.value = 'detail'
    }
  } catch(e) {
    alert('Error al guardar el pedido.')
  } finally {
    isSaving.value = false
  }
}

function updateLines(newLines) {
  editableOrder.value.lines = newLines
}

function printOrder() {
  printComponent.value.print()
}

function formatDate(dateString) { return dateString ? new Date(dateString).toLocaleDateString('es-ES') : '' }
function formatMoney(amount) { return salesApi.formatMoney(amount) }

watch(orderId, loadOrder, { immediate: true })
</script>

<style scoped>
.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}
.summary-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
.summary-item label {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  font-weight: 600;
}
.totals-section {
  max-width: 350px;
  margin-left: auto;
  margin-top: 2rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.total-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
}
.total-row.grand-total {
  font-size: 1.25rem;
  font-weight: bold;
  border-top: 2px solid var(--color-border);
  margin-top: 0.5rem;
  padding-top: 1rem;
}
</style>
