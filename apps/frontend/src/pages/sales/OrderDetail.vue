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
          <div v-if="order || mode === 'create'" class="header-actions-group">
            <button v-if="mode === 'detail'" class="btn btn-outline btn-sm" @click="printOrder">
              <span class="material-symbols-outlined">print</span>
              <span>Imprimir</span>
            </button>
            <button v-if="mode === 'detail'" class="btn btn-primary btn-sm" @click="enterEditMode">
              <span class="material-symbols-outlined">edit</span>
              <span>Editar</span>
            </button>
            <template v-else>
              <button class="btn btn-outline btn-sm" @click="exitEditMode" :disabled="isSaving">Cancelar</button>
              <button class="btn btn-primary btn-sm ml-2" @click="saveOrder" :disabled="isSaving">
                <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
                <span>{{ isSaving ? 'Guardando...' : 'Guardar Pedido' }}</span>
              </button>
            </template>
          </div>
        </template>
      </PageHeader>
    </template>
    
    <!-- CAPA 2: CONTEXTO -->
    <template #summary v-if="mode === 'detail' && order">
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><span class="material-symbols-outlined">person</span></div>
          <div class="tag-content">
            <label>Cliente</label>
            <strong>{{ order.party_name || order.partyName }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><span class="material-symbols-outlined">calendar_today</span></div>
          <div class="tag-content">
            <label>Fecha Pedido</label>
            <strong>{{ formatDate(order.order_date || order.orderDate) }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><span class="material-symbols-outlined">local_shipping</span></div>
          <div class="tag-content">
            <label>Fecha Entrega</label>
            <strong>{{ formatDate(order.delivery_date || order.deliveryDate) }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><span class="material-symbols-outlined">payments</span></div>
          <div class="tag-content">
            <label>Total Pedido</label>
            <strong class="amount">{{ formatMoney(totalAmount) }}</strong>
          </div>
        </div>
      </div>
    </template>
    
    <!-- CAPA 3: TRABAJO -->
    <div v-if="order || mode === 'create'" class="order-detail-content">
      <!-- Sección Cliente y Condiciones -->
      <FormSection title="Identificación del Cliente y Plazos" icon="person">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
          <div class="form-item-wrapper">
            <template v-if="mode !== 'detail'">
              <PartySelector
                v-model="editableOrder.party_id"
                label="Cliente *"
                placeholder="Buscar cliente por nombre o NIF..."
                required
                @select="handlePartyChange"
              />
            </template>
            <DataRow v-else label="Cliente" icon="person">
              <strong>{{ order.party_name || order.partyName }}</strong> 
              <code class="ml-2 text-xs">#{{ (order.party_id || order.partyId || '').substring(0,8) }}</code>
            </DataRow>
          </div>

          <div class="form-group">
            <label v-if="mode !== 'detail'">Referencia Cliente</label>
            <template v-if="mode !== 'detail'">
              <input v-model="editableOrder.party_reference" class="form-input" placeholder="Ej: PO-12345" />
            </template>
            <DataRow v-else label="Referencia Cliente" icon="badge">
              {{ order.party_reference || order.partyReference || '—' }}
            </DataRow>
          </div>

          <div class="form-group">
            <label v-if="mode !== 'detail'">Fecha del Pedido *</label>
            <template v-if="mode !== 'detail'">
              <input type="date" v-model="editableOrder.order_date" class="form-input" required />
            </template>
            <DataRow v-else label="Fecha del Pedido" icon="event">
              {{ formatDate(order.order_date || order.orderDate) }}
            </DataRow>
          </div>

          <div class="form-group">
            <label v-if="mode !== 'detail'">Fecha Estimada de Entrega</label>
            <template v-if="mode !== 'detail'">
              <input type="date" v-model="editableOrder.delivery_date" class="form-input" />
            </template>
            <DataRow v-else label="Fecha de Entrega" icon="local_shipping">
              {{ formatDate(order.delivery_date || order.deliveryDate) }}
            </DataRow>
          </div>
        </div>
      </FormSection>

      <!-- SECCIÓN MES -->
      <FormSection title="Configuración Técnica (MES)" icon="precision_manufacturing">
        <div v-if="mode === 'detail'" class="table-wrapper">
          <table v-if="(order.mes_work_refs || order.mesWorkRefs || []).length > 0" class="data-table">
            <thead>
              <tr>
                <th>Proceso / Configuración</th>
                <th>Descripción</th>
                <th class="text-right">Estado</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mesRef in (order.mes_work_refs || order.mesWorkRefs)" :key="mesRef.id">
                <td>
                  <div class="flex items-center gap-2">
                    <span class="material-symbols-outlined text-secondary">settings_suggest</span>
                    <strong>{{ formatMesWorkId(mesRef.work_setup_id || mesRef.workSetupId) }}</strong>
                  </div>
                </td>
                <td><p class="text-sm italic m-0">{{ mesRef.description || '—' }}</p></td>
                <td class="text-right">
                  <span class="status-badge status-info">Vinculado</span>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-else class="text-muted p-4 text-center italic">No hay requerimientos técnicos definidos para este pedido.</p>
        </div>

        <div v-else>
          <div class="mb-4">
            <button type="button" class="btn btn-outline-secondary btn-sm" @click="addMesWorkRef">
              <span class="material-symbols-outlined">add</span> <span>Añadir Requerimiento Técnico</span>
            </button>
          </div>
          <div class="table-wrapper border rounded-lg overflow-hidden">
            <table class="data-table fixed-layout">
              <thead>
                <tr>
                  <th style="width: 50px">#</th>
                  <th style="width: 250px">Configuración Cliente</th>
                  <th>Notas Técnicas / Instrucciones *</th>
                  <th class="text-center" style="width: 80px">Borrar</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(ref, idx) in editableOrder.mes_work_refs" :key="idx">
                  <td class="text-muted font-bold">{{ idx + 1 }}</td>
                  <td>
                    <select v-model="ref.work_setup_id" class="form-input w-full">
                      <option :value="null">-- Personalizado --</option>
                      <option v-for="setup in availableMesSetups" :key="setup.id" :value="setup.id">{{ setup.name }}</option>
                    </select>
                  </td>
                  <td>
                    <input v-model="ref.description" type="text" class="form-input w-full" placeholder="Ej: Color especial..." required />
                  </td>
                  <td class="text-center">
                    <button type="button" class="btn-icon text-danger" @click="removeMesWorkRef(idx)">
                      <span class="material-symbols-outlined">delete</span>
                    </button>
                  </td>
                </tr>
                <tr v-if="editableOrder.mes_work_refs.length === 0">
                  <td colspan="4" class="text-muted text-center p-6 italic">No se han añadido trabajos MES. Haz clic en "Añadir".</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </FormSection>
      
      <!-- Sección Líneas de Pedido -->
      <FormSection title="Líneas de Pedido" icon="list_alt">
        <div v-if="mode !== 'detail'" class="mb-4">
          <button type="button" class="btn btn-primary btn-sm" @click="showVariantSelector = true">
            <span class="material-symbols-outlined">add</span> <span>Añadir Producto</span>
          </button>
        </div>
        <OrderLines
          :lines="mode === 'detail' ? (order.line_items || order.lineItems) : editableOrder.line_items"
          :is-editing="mode !== 'detail'"
          @update:lines="updateLines"
        />
      </FormSection>
      
      <!-- Sección de Totales -->
      <FormSection title="Resumen Económico" icon="payments">
        <div class="totals-checkout-layout">
          <section class="totals-checkout-card">
            <div class="total-row">
              <label>Subtotal:</label>
              <span>{{ formatMoney(subtotal) }}</span>
            </div>
            <div class="total-row">
              <label>IVA (21%):</label>
              <span>{{ formatMoney(taxAmount) }}</span>
            </div>
            <div class="total-row grand-total">
              <label>TOTAL PEDIDO:</label>
              <span class="total-value">{{ formatMoney(totalAmount) }}</span>
            </div>
          </section>
        </div>
      </FormSection>

      <FormSection v-if="mode !== 'detail'" title="Observaciones internas" icon="notes" class="mt-6">
        <div class="form-group">
          <textarea v-model="editableOrder.notes" class="form-textarea" rows="3" placeholder="Estas notas no se imprimen en el documento oficial..."></textarea>
        </div>
      </FormSection>
    </div>
  </BaseEntityPage>
  
  <!-- MODAL: SELECCIÓN DE PRODUCTO -->
  <BaseDialog
    :show="showVariantSelector"
    title="Seleccionar Producto"
    icon="inventory_2"
    size="xl"
    hide-actions
    @close="showVariantSelector = false"
  >
    <VariantSelector @variant-selected="handleVariantSelected" />
  </BaseDialog>

  <!-- PORTAL DE IMPRESIÓN (Solo visible en @media print) -->
  <PrintDocument v-if="order" :data="printData" />
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue'
import PageHeader from '@/components/layout/PageHeader.vue'
import FormSection from '@/components/shared/FormSection.vue'
import DataRow from '@/components/shared/DataRow.vue'
import PartySelector from '@/components/party/PartySelector.vue'
import OrderLines from '@/components/sales/OrderLines.vue'
import VariantSelector from '@/components/product/VariantSelector.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import PrintDocument from '@/components/shared/PrintDocument.vue'
import salesApi from '@/services/salesApi'
import { partyApi } from '@/services/partyApi'
import { mesApi } from '@/services/mesApi'

const route = useRoute()
const router = useRouter()
const order = ref(null)
const editableOrder = ref({ line_items: [], mes_work_refs: [] })
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref('')
const mode = ref('detail')
const showVariantSelector = ref(false)
const printComponent = ref(null)

const availableMesSetups = ref([])
const mesWorksCache = ref({})

const orderId = computed(() => route.params.id)
const isCreateMode = computed(() => !orderId.value || orderId.value === 'new')

const pageTitle = computed(() => {
  if (isCreateMode.value) return 'Nuevo Pedido'
  const num = order.value?.order_number || order.value?.orderNumber || '...'
  return mode.value === 'edit' ? `Editando Pedido #${num}` : `Pedido #${num}`
})

const headerLabel = computed(() => {
  if (isCreateMode.value) return 'Nuevo'
  return `#${order.value?.order_number || order.value?.orderNumber || '...'}`
})

const statusLabel = computed(() => salesApi.getStatusLabel(order.value?.status))
const statusClass = computed(() => salesApi.getStatusClass(order.value?.status))

const subtotal = computed(() => {
  if (mode.value === 'detail' && order.value) {
    return order.value.subtotal?.amount ?? order.value.subtotal ?? 0
  }
  return (editableOrder.value.line_items || []).reduce((acc, line) => {
    return acc + (Number(line.quantity || 0) * Number(line.unit_price || 0)) * (1 - (Number(line.discount_percent || 0) / 100))
  }, 0)
})

const taxAmount = computed(() => {
  if (mode.value === 'detail' && order.value) {
    return order.value.tax_total?.amount ?? order.value.tax_amount?.amount ?? order.value.tax_amount ?? 0
  }
  return subtotal.value * 0.21
})

const totalAmount = computed(() => {
  if (mode.value === 'detail' && order.value) {
    return order.value.total?.amount ?? order.value.total ?? 0
  }
  return subtotal.value + taxAmount.value
})

async function loadOrder() {
  if (isCreateMode.value) {
    mode.value = 'create'
    resetForm()
    return
  }

  isLoading.value = true
  error.value = ''
  try {
    const data = await salesApi.getOrder(orderId.value)
    order.value = data
    syncEditableOrder(data)
    if (data.party_id || data.partyId) {
      loadAvailableSetups(data.party_id || data.partyId)
    }
    mode.value = 'detail'
  } catch (e) {
    console.error('Error loading order:', e)
    error.value = 'No se pudo cargar el pedido.'
  } finally {
    isLoading.value = false
  }
}

function resetForm() {
  editableOrder.value = {
    party_id: '',
    order_date: new Date().toISOString().split('T')[0],
    delivery_date: '',
    party_reference: '',
    mes_work_refs: [],
    line_items: [],
    notes: ''
  }
  order.value = null
}

function syncEditableOrder(data) {
  editableOrder.value = {
    party_id: data.party_id || data.partyId,
    order_date: data.order_date ? new Date(data.order_date).toISOString().split('T')[0] : (data.orderDate ? new Date(data.orderDate).toISOString().split('T')[0] : ''),
    delivery_date: data.delivery_date ? new Date(data.delivery_date).toISOString().split('T')[0] : (data.deliveryDate ? new Date(data.deliveryDate).toISOString().split('T')[0] : ''),
    party_reference: data.party_reference || data.partyReference || '',
    mes_work_refs: (data.mes_work_refs || data.mesWorkRefs || []).map(r => ({
      id: r.id,
      work_setup_id: r.work_setup_id || r.workSetupId || null,
      description: r.description || ''
    })),
    line_items: (data.line_items || data.lineItems || []).map(li => ({
      ...li,
      unit_price: li.unit_price?.amount ?? li.unit_price ?? li.unitPrice?.amount ?? li.unitPrice ?? 0,
      quantity: li.quantity || 0,
      discount_percent: li.discount_percent || li.discountPercent || 0,
      product_variant_id: li.product_variant_id || li.productVariantId
    })),
    notes: data.notes || ''
  }
}

function enterEditMode() {
  if (editableOrder.value.party_id) {
    loadAvailableSetups(editableOrder.value.party_id)
  }
  mode.value = 'edit'
}

function exitEditMode() {
  if (isCreateMode.value) router.push('/sales/orders')
  else mode.value = 'detail'
}

function handlePartyChange(party) {
  if (party) {
    editableOrder.value.party_id = party.id
    loadAvailableSetups(party.id)
  }
}

async function loadAvailableSetups(partyId) {
  if (!partyId) return
  try {
    const res = await mesApi.listWorkSetups({ party_id: partyId })
    availableMesSetups.value = res.data || res || []
    availableMesSetups.value.forEach(s => { mesWorksCache.value[s.id] = s })
  } catch (err) {
    console.error('Error cargando setups MES:', err)
  }
}

function addMesWorkRef() {
  editableOrder.value.mes_work_refs.push({ work_setup_id: null, description: '' })
}

function removeMesWorkRef(idx) {
  editableOrder.value.mes_work_refs.splice(idx, 1)
}

function formatMesWorkId(id) {
  if (!id) return 'Sin config.'
  return mesWorksCache.value[id]?.name || id.substring(0, 8)
}

function handleVariantSelected(payload) {
  const variant = payload.variant || payload
  editableOrder.value.line_items.push({
    product_variant_id: variant.id,
    variant_sku: variant.sku,
    product_name: variant.product_name || variant.name,
    quantity: 1,
    unit_price: variant.product_base_price || variant.basePrice || variant.price || 0,
    discount_percent: 0
  })
  showVariantSelector.value = false
}

async function saveOrder() {
  if (!editableOrder.value.party_id) { alert('Debe seleccionar un cliente'); return; }
  if (!editableOrder.value.line_items.length) { alert('El pedido debe tener al menos una línea'); return; }
  
  isSaving.value = true
  try {
    const formatToRFC3339 = (dateStr) => {
      if (!dateStr) return null
      return new Date(dateStr).toISOString()
    }

    const payload = { 
      partyId: editableOrder.value.party_id,
      orderDate: formatToRFC3339(editableOrder.value.order_date),
      deliveryDate: formatToRFC3339(editableOrder.value.delivery_date) || formatToRFC3339(editableOrder.value.order_date),
      partyReference: editableOrder.value.party_reference || "",
      notes: editableOrder.value.notes || "",
      mesWorkRefs: (editableOrder.value.mes_work_refs || []).map(r => ({
        workSetupId: r.work_setup_id || undefined,
        description: r.description || ''
      })),
      items: editableOrder.value.line_items.map(li => ({
        productVariantId: li.product_variant_id || li.productVariantId,
        quantity: Number(li.quantity),
        unitPrice: { amount: Number(li.unit_price), currency: 'EUR' },
        discountPercent: Number(li.discount_percent || 0)
      }))
    }

    if (isCreateMode.value) {
      const newOrder = await salesApi.createOrder(payload)
      router.push(`/sales/orders/${newOrder.id}`)
    } else {
      // 1. Actualizar Cabecera
      await salesApi.updateOrder(orderId.value, payload)
      
      // 2. Gestionar Líneas (Sincronización granular)
      const originalLines = order.value.line_items || order.value.lineItems || []
      const currentLines = editableOrder.value.line_items
      
      // ELIMINAR líneas que ya no están
      const toDelete = originalLines.filter(ol => !currentLines.some(cl => cl.id === ol.id))
      for (const line of toDelete) {
        await salesApi.removeOrderLineItem(orderId.value, line.id)
      }
      
      // AÑADIR o ACTUALIZAR
      for (const line of currentLines) {
        const itemPayload = {
          productVariantId: line.product_variant_id || line.productVariantId,
          quantity: Number(line.quantity),
          unitPrice: { amount: Number(line.unit_price), currency: 'EUR' },
          discountPercent: Number(line.discount_percent || 0)
        }
        
        if (!line.id || line.id.length < 10) { // Es nueva (ID temporal o inexistente)
          await salesApi.addOrderLineItem(orderId.value, itemPayload)
        } else {
          const original = originalLines.find(ol => ol.id === line.id)
          // Solo actualizar si algo ha cambiado para ahorrar peticiones
          if (original && (original.quantity !== line.quantity || original.unit_price !== line.unit_price || original.discount_percent !== line.discount_percent)) {
            await salesApi.updateOrderLineItem(orderId.value, line.id, itemPayload)
          }
        }
      }

      await loadOrder()
      mode.value = 'detail'
    }
  } catch(e) {
    console.error('Error saving order:', e)
    alert('Error al guardar el pedido: ' + (e.message || 'Error desconocido'))
  } finally {
    isSaving.value = false
  }
}

function updateLines(newLines) {
  editableOrder.value.line_items = newLines
}

function printOrder() {
  window.print()
}

const printData = computed(() => {
  if (!order.value) return null;
  return {
    type: 'PEDIDO',
    number: order.value.order_number || order.value.orderNumber || '—',
    date: order.value.order_date || order.value.orderDate,
    expiryDate: order.value.delivery_date || order.value.deliveryDate,
    party: {
      name: order.value.party_name || order.value.partyName,
      taxId: order.value.tax_id || order.value.taxId,
      address: order.value.address,
    },
    items: (order.value.line_items || order.value.lineItems || []).map(li => ({
      sku: li.variant_sku || li.variantSku,
      name: li.product_name || li.productName,
      quantity: li.quantity,
      unitPrice: li.unit_price?.amount || li.unitPrice?.amount || li.unit_price || 0,
      discount: li.discount_percent || li.discountPercent,
      subtotal: li.subtotal?.amount || li.subtotal || 0
    })),
    subtotal: subtotal.value,
    taxAmount: taxAmount.value,
    total: totalAmount.value,
    notes: order.value.notes
  }
});

function formatDate(dateString) { return dateString ? new Date(dateString).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '' }
function formatMoney(amount) { return salesApi.formatMoney(amount) }

watch(() => route.params.id, loadOrder, { immediate: true })
</script>

<style scoped>
@import "@/design-system/_sections.css";

.order-detail-content { display: flex; flex-direction: column; gap: 1.5rem; }

.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.025em; margin-bottom: 0.5rem; }

.form-input, .form-textarea {
  padding: 0.75rem 1rem; border: 1px solid var(--color-border-strong); border-radius: 8px;
  font-size: 0.9rem; font-family: inherit; transition: all 0.2s; color: var(--color-text-primary);
}
.form-input:focus, .form-textarea:focus { outline: none; border-color: var(--color-secondary); box-shadow: 0 0 0 3px rgba(27, 58, 107, 0.1); }

.overview-tags-row { display: flex; flex-wrap: wrap; gap: 1rem; margin-bottom: 0.5rem; }
.summary-tag { flex: 1; min-width: 200px; padding: 0.75rem 1.25rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 1rem; box-shadow: var(--box-shadow-sm); }

.icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.icon .material-symbols-outlined { font-size: 24px; }
.icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }

.tag-content { display: flex; flex-direction: column; gap: 0.2rem; }
.tag-content label { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }
.tag-content strong { font-size: 1rem; color: var(--color-text-primary); font-weight: 700; }
.amount { color: var(--color-success) !important; font-size: 1.25rem !important; font-family: var(--font-family-mono); }

.totals-checkout-layout { display: flex; justify-content: flex-end; }
.totals-checkout-card { 
  background: white; border: 1px solid var(--color-border-strong); border-radius: 12px; padding: 1.5rem; 
  width: 100%; max-width: 400px; box-shadow: var(--box-shadow-md); 
}
.total-row { display: flex; justify-content: space-between; padding: 0.5rem 0; font-size: 0.95rem; }
.total-row label { color: var(--color-text-secondary); font-weight: 600; }
.total-row.grand-total { margin-top: 1rem; padding-top: 1rem; border-top: 2px solid var(--color-border); font-weight: 800; font-size: 1.25rem; }
.total-row.grand-total span { color: var(--color-secondary); }

.status-badge { padding: 0.2rem 0.6rem; border-radius: 4px; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; }

.print-only { display: none; }
@media print {
  .print-only { display: block !important; }
  :deep(.page-layout), :deep(.navbar), :deep(.side-navbar), :deep(.app-header), :deep(header) { display: none !important; }
}
</style>
