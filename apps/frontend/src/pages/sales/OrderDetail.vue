<template>
  <BaseEntityPage class="no-print" :is-loading="isLoading" :error="error">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <BasePageHeader 
        :title="pageTitle" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/dashboard' }, { label: 'Pedidos', to: '/sales/orders' }, { label: headerLabel }]"
        show-back
      >
        <template #icon>
          <ShoppingCart :size="28" />
        </template>
        <template #actions>
          <div v-if="order || mode === 'create'" class="header-actions-group">
            <template v-if="mode === 'detail'">
              <button class="btn btn-outline btn-sm" @click="printOrder">
                <Printer :size="16" />
                <span>Imprimir</span>
              </button>
              <button v-if="!hasActiveDeliveryNotes" class="btn btn-primary btn-sm" @click="enterEditMode">
                <Pencil :size="16" />
                <span>Editar Pedido</span>
              </button>
            </template>
            <template v-else>
              <button class="btn btn-outline btn-sm" @click="exitEditMode" :disabled="isSaving">Cancelar</button>
              <button class="btn btn-secondary btn-sm" @click="saveOrder" :disabled="isSaving">
                <component :is="isSaving ? RefreshCw : Save" :size="16" :class="{ 'spin': isSaving }" />
                <span>{{ isSaving ? 'Guardando...' : 'Guardar Pedido' }}</span>
              </button>
            </template>
          </div>
        </template>
      </BasePageHeader>
    </template>

    <!-- 2. TOOLBAR: ACCIONES DE FLUJO -->
    <template #toolbar v-if="mode === 'detail' && order">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', `status-${statusClass}`]">
            {{ statusLabel }}
          </span>
        </div>
        <div class="toolbar-buttons">
          <!-- Lanzar a Producción (Solo si está pendiente Y tiene trabajos MES definidos) -->
          <button 
            v-if="order.status === 'PENDING' && (order.mes_work_refs || order.mesWorkRefs || []).length > 0" 
            class="btn btn-primary btn-sm" 
            @click="promptLaunchProduction"
          >
            <Rocket :size="16" />
            <span>Lanzar a Producción</span>
          </button>

          <!-- Albaranar (de Confirmado a Albaranado) -->
          <button 
            v-if="['PENDING', 'IN_PREPARATION', 'READY_FOR_PRODUCTION', 'PARTIALLY_DELIVERED'].includes(order.status)" 
            class="btn btn-success btn-sm" 
            @click="createDeliveryNote"
          >
            <Truck :size="16" />
            <span>Albaranar</span>
          </button>
          
          <!-- Reactivar (Si está anulado) -->
          <button 
            v-if="order.status === 'CANCELLED'" 
            class="btn btn-primary btn-sm" 
            @click="promptReactivate"
          >
            <RefreshCw :size="16" />
            <span>Reactivar Pedido</span>
          </button>

          <!-- Anular (Solo si no está entregado ni facturado ni ya anulado) -->
          <button 
            v-if="!hasActiveDeliveryNotes && !['DELIVERED', 'PARTIALLY_DELIVERED', 'INVOICED', 'PARTIALLY_INVOICED', 'CANCELLED'].includes(order.status)" 
            class="btn btn-danger btn-sm" 
            @click="promptCancelOrder"
          >
            <Ban :size="16" />
            <span>Anular Pedido</span>
          </button>
        </div>
      </div>
    </template>
    
    <!-- CAPA 2: CONTEXTO -->
    <template #summary v-if="mode === 'detail' && order">
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><User :size="20" /></div>
          <div class="tag-content">
            <label>Cliente</label>
            <strong>{{ order.party_name || order.partyName }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><Calendar :size="20" /></div>
          <div class="tag-content">
            <label>Fecha Pedido</label>
            <strong>{{ formatDate(order.order_date || order.orderDate) }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><Truck :size="20" /></div>
          <div class="tag-content">
            <label>Fecha Entrega</label>
            <strong>{{ formatDate(order.delivery_date || order.deliveryDate) }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><CreditCard :size="20" /></div>
          <div class="tag-content">
            <label>Total Pedido</label>
            <strong class="amount">{{ formatMoney(totalAmount) }}</strong>
          </div>
        </div>
      </div>
    </template>

    <template #related v-if="mode === 'detail' && order">
      <div class="related-history-grid">
        <!-- 1. Presupuesto Origen -->
        <router-link v-if="relatedQuote" :to="`/sales/quotes/${relatedQuote.id}`" class="related-tag-card highlight-info">
          <div class="tag-icon"><FileQuestion :size="20" /></div>
          <div class="tag-content">
            <label>Presupuesto Origen</label>
            <strong>{{ relatedQuote.quoteNumber || relatedQuote.quote_number }}</strong>
          </div>
          <ExternalLink :size="14" class="jump-icon" />
        </router-link>

        <!-- 2. Albaranes -->
        <router-link v-for="dn in relatedDeliveryNotes" :key="dn.id" :to="`/sales/delivery-notes/${dn.id}`" class="related-tag-card">
          <div class="tag-icon"><Truck :size="20" /></div>
          <div class="tag-content">
            <label>Albarán Generado</label>
            <strong>{{ dn.deliveryNoteNumber || dn.delivery_note_number }}</strong>
          </div>
          <ExternalLink :size="14" class="jump-icon" />
        </router-link>

        <!-- 3. Factura -->
        <router-link v-if="relatedInvoice" :to="`/sales/invoices/${relatedInvoice.id}`" class="related-tag-card">
          <div class="tag-icon success"><Receipt :size="20" /></div>
          <div class="tag-content">
            <label>Factura Vinculada</label>
            <strong>{{ relatedInvoice.invoiceNumber || relatedInvoice.invoice_number }}</strong>
          </div>
          <ExternalLink :size="14" class="jump-icon" />
        </router-link>
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
        <div v-if="mode === 'detail' && ['PENDING', 'PENDIENTE', 'CONFIRMED', 'CONFIRMADO', 'EN_PREPARACION'].includes(order.status)" class="info-notice mb-4">
          <Info :size="20" />
          <div>
            <strong>A la espera de lanzamiento operativo.</strong>
            <p class="m-0 text-xs">El taller no visualizará este pedido hasta que se pulse el botón "Lanzar a Producción".</p>
          </div>
        </div>

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
                    <Settings2 :size="18" class="text-secondary" />
                    <strong>{{ formatMesWorkId(mesRef.work_setup_id || mesRef.workSetupId) }}</strong>
                  </div>
                </td>
                <td><p class="text-sm italic m-0">{{ mesRef.description || '—' }}</p></td>
                <td class="text-right">
                  <template v-if="mesRef.work_order_id || mesRef.workOrderId">
                    <button 
                      class="status-badge clickable-status" 
                      :class="`status-${getMesStatusClass(mesOrdersStatus[mesRef.work_order_id || mesRef.workOrderId])}`"
                      @click="router.push(`/mes/work-orders/${mesRef.work_order_id || mesRef.workOrderId}`)"
                    >
                      {{ mesApi.getWorkStatusLabel(mesOrdersStatus[mesRef.work_order_id || mesRef.workOrderId]) || 'Cargando...' }}
                      <ExternalLink :size="14" />
                    </button>
                  </template>
                  <span v-else class="status-badge status-secondary">Sin lanzar</span>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-else class="text-muted p-4 text-center italic">No hay requerimientos técnicos definidos para este pedido.</p>
        </div>

        <div v-else>
          <div class="mb-4">
            <button type="button" class="btn btn-primary btn-sm" @click="addMesWorkRef">
              <Plus :size="16" /> <span>Añadir Trabajo Taller</span>
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
                      <Trash2 :size="18" />
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
        <OrderLines
          :lines="mode === 'detail' ? (order.line_items || order.lineItems) : editableOrder.line_items"
          :is-editing="mode !== 'detail'"
          @update:lines="updateLines"
          @add-line-request="handleAddLineRequest"
          @last-field-tab="focusAddButton"
        />
        <div v-if="mode !== 'detail'" class="mt-4">
          <button 
            ref="addProductBtnRef"
            type="button" 
            class="btn btn-primary btn-sm" 
            @click="showVariantSelector = true"
          >
            <Plus :size="16" /> <span>Añadir Producto (Ins)</span>
          </button>
        </div>
      </FormSection>
      
      <!-- Sección de Totales -->
      <FormSection title="Resumen Económico" icon="payments">
        <div class="totals-summary-container">
          <section class="totals-summary-card" :class="{ 'is-loading-overlay': isPreviewLoading }">
            <div class="summary-row">
              <label>Subtotal:</label>
              <span>{{ formatMoney(subtotal) }}</span>
            </div>
            <div class="summary-row">
              <label>Impuestos (IVA 21%):</label>
              <span>{{ formatMoney(taxAmount) }}</span>
            </div>
            <div class="summary-row grand-total">
              <label>{{ mode === 'detail' ? 'TOTAL PEDIDO:' : 'TOTAL ESTIMADO:' }}</label>
              <span>{{ formatMoney(totalAmount) }}</span>
            </div>

            <!-- Overlay de carga solo en edición -->
            <div v-if="isPreviewLoading && mode !== 'detail'" class="mini-spinner-overlay">
              <div class="mini-spinner"></div>
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

    <!-- MODALES DE CONFIRMACIÓN (REEMPLAZO DE confirm()) -->
    <BaseDialog
      :show="confirmDialog.show"
      :title="confirmDialog.title"
      :icon="confirmDialog.icon"
      :confirm-text="confirmDialog.confirmText"
      :confirm-class="confirmDialog.confirmClass"
      :is-confirming="isSaving"
      @close="confirmDialog.show = false"
      @confirm="handleConfirmDialog"
    >
      <p>{{ confirmDialog.message }}</p>
    </BaseDialog>

    <!-- MODAL: ALBARANADO PARCIAL -->
    <BaseDialog
      :show="showDnDialog"
      title="Generar Albarán de Salida"
      icon="local_shipping"
      size="xl"
      :is-loading="isSaving"
      confirm-text="Generar Albarán"
      @close="showDnDialog = false"
      @confirm="submitDeliveryNote"
    >
      <div class="dn-dialog-content">
        <div class="mb-6 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="form-group">
            <label>Fecha de Entrega / Salida</label>
            <input type="date" v-model="dnForm.deliveryDate" class="form-input" />
          </div>
          <div class="flex items-end gap-2">
            <button class="btn btn-outline-secondary btn-sm" @click="deliverAll">Albaranar Todo</button>
            <button class="btn btn-outline-secondary btn-sm" @click="deliverNone">Limpiar Cantidades</button>
          </div>
        </div>

        <div class="table-wrapper border rounded-lg overflow-hidden">
          <table class="data-table">
            <thead>
              <tr>
                <th>Producto / Variante</th>
                <th class="text-center">Total Pedido</th>
                <th class="text-center">Ya Albaranado</th>
                <th class="text-center">Pendiente</th>
                <th class="text-center" style="width: 150px">A Entregar Ahora</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in dnForm.items" :key="item.orderLineItemId">
                <td>
                  <div class="flex flex-col">
                    <strong>{{ item.productName }}</strong>
                    <small class="text-muted">{{ item.variantSku }}</small>
                  </div>
                </td>
                <td class="text-center">{{ item.totalQuantity }}</td>
                <td class="text-center">
                  <span :class="{'text-success': item.alreadyDelivered > 0}">{{ item.alreadyDelivered }}</span>
                </td>
                <td class="text-center font-bold text-primary">{{ item.pendingQuantity }}</td>
                <td class="text-center">
                  <div class="flex items-center justify-center gap-2">
                    <input 
                      type="number" 
                      v-model.number="item.quantityToDeliver" 
                      min="0" 
                      :max="item.pendingQuantity"
                      class="form-input text-center" 
                      style="width: 80px"
                    />
                    <span class="text-xs text-muted">/ {{ item.pendingQuantity }}</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <p class="mt-4 text-sm text-muted italic">
          * Solo se muestran los ítems que tienen cantidades pendientes de entrega.
        </p>
      </div>
    </BaseDialog>

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
    <div class="print-container">
      <PrintDocument
        v-if="order"
        type="ORDER"
        :number="order.order_number || order.orderNumber"
        :date="order.order_date || order.orderDate"
        :customer-name="order.party_name || order.partyName"
        :customer-tax-id="order.tax_id || order.taxId"
        :items="order.line_items || order.lineItems"
        :totals="{ subtotal: subtotal, taxAmount: taxAmount, total: totalAmount }"
        :notes="order.notes"
      />
    </div>
  </BaseEntityPage>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick, reactive } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { 
  ShoppingCart, 
  Printer, 
  Pencil, 
  RefreshCw, 
  Save, 
  Rocket, 
  Truck, 
  Ban, 
  User, 
  Calendar, 
  CreditCard, 
  FileQuestion, 
  Receipt, 
  ExternalLink, 
  Info, 
  Settings2, 
  Plus, 
  Trash2,
  AlertTriangle,
  Play
} from 'lucide-vue-next'
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import FormSection from '@/components/shared/FormSection.vue'
import DataRow from '@/components/shared/DataRow.vue'
import PartySelector from '@/components/party/PartySelector.vue'
import OrderLines from '@/components/sales/OrderLines.vue'
import VariantSelector from '@/components/product/VariantSelector.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import PrintDocument from '@/components/sales/PrintDocument.vue'
import salesApi from '@/services/salesApi'
import { partyApi } from '@/services/partyApi'
import { mesApi } from '@/services/mesApi'
import { useToastStore } from '@/stores/toast'
import '@/assets/sales-print.css'

const route = useRoute()
const router = useRouter()
const toastStore = useToastStore()
const order = ref(null)
const editableOrder = ref({ line_items: [], mes_work_refs: [] })
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref('')
const mode = ref('detail')
const showVariantSelector = ref(false)
const printComponent = ref(null)

const showDnDialog = ref(false)
const dnForm = ref({
  deliveryDate: new Date().toISOString().split('T')[0],
  items: []
})

const availableMesSetups = ref([])
const mesWorksCache = ref({})
const mesOrdersStatus = ref({}) // Almacena workOrderId -> status
const partyDefaultDiscount = ref(0)
const relatedQuote = ref(null)
const relatedDeliveryNotes = ref([])
const relatedInvoice = ref(null)

// --- Confirm Dialog Logic ---
const confirmDialog = reactive({
  show: false,
  title: '',
  message: '',
  icon: 'help-circle',
  confirmText: 'Confirmar',
  confirmClass: 'btn-primary',
  action: null
})

function promptCancelOrder() {
  confirmDialog.title = 'Anular Pedido'
  confirmDialog.message = '¿Realmente deseas ANULAR este pedido? Esta acción es irreversible y detendrá cualquier proceso vinculado.'
  confirmDialog.icon = AlertTriangle
  confirmDialog.confirmText = 'Sí, Anular Pedido'
  confirmDialog.confirmClass = 'btn-danger'
  confirmDialog.action = cancelOrder
  confirmDialog.show = true
}

function promptReactivate() {
  confirmDialog.title = 'Reactivar Pedido'
  confirmDialog.message = '¿Deseas reactivar este pedido y devolverlo al estado borrador para realizar cambios?'
  confirmDialog.icon = RefreshCw
  confirmDialog.confirmText = 'Reactivar Pedido'
  confirmDialog.confirmClass = 'btn-primary'
  confirmDialog.action = reactivateOrder
  confirmDialog.show = true
}

function promptLaunchProduction() {
  confirmDialog.title = 'Lanzar a Producción'
  confirmDialog.message = '¿Lanzar este pedido a producción? Esta acción notificará al taller y permitirá iniciar los trabajos MES.'
  confirmDialog.icon = Rocket
  confirmDialog.confirmText = 'Lanzar Ahora'
  confirmDialog.confirmClass = 'btn-secondary'
  confirmDialog.action = launchOrderToProduction
  confirmDialog.show = true
}

async function handleConfirmDialog() {
  if (confirmDialog.action) {
    await confirmDialog.action()
  }
  confirmDialog.show = false
}

const hasActiveDeliveryNotes = computed(() =>
  relatedDeliveryNotes.value.some(dn => dn.status !== 'CANCELLED')
)

function getMesStatusClass(status) {
  if (!status) return 'secondary'
  const map = {
    'PENDING': 'warning',
    'IN_PROGRESS': 'info',
    'COMPLETED': 'success',
    'ON_HOLD': 'secondary',
    'SUSPENDED': 'secondary',
    'CANCELLED': 'danger'
  }
  return map[status] || 'secondary'
}

async function fetchMesWorkStatuses(mesRefs) {
  if (!mesRefs || mesRefs.length === 0) return
  
  // Limpiamos estados previos
  mesOrdersStatus.value = {}
  
  // Consultamos el estado de cada orden vinculada
  const promises = mesRefs
    .filter(r => r.work_order_id || r.workOrderId)
    .map(async (r) => {
      const id = r.work_order_id || r.workOrderId
      try {
        const wo = await mesApi.getWorkOrder(id)
        if (wo) {
          mesOrdersStatus.value[id] = wo.status
        }
      } catch (err) {
        console.error(`Error cargando estado MES ${id}:`, err)
      }
    })
    
  await Promise.all(promises)
}

const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewTimer = null;

const orderId = computed(() => route.params.id)
const isCreateMode = computed(() => !orderId.value || orderId.value === 'new')

// Watcher para cambios en el formulario que requieran recalcular totales
watch(() => [editableOrder.value.party_id, editableOrder.value.line_items], () => {
  if (mode.value !== 'detail' && !isSaving.value) {
    calculateTotals();
  }
}, { deep: true });

function calculateTotals() {
  clearTimeout(previewTimer);
  // Reset the previous result so that computed values use the local fallback 
  // instantaneously while waiting for the server response (avoids frozen values).
  previewResult.value = null;
  previewTimer = setTimeout(fetchPreviewCalculation, 400);
}

async function fetchPreviewCalculation() {
  const partyId = editableOrder.value.party_id;
  const items = (editableOrder.value.line_items || []).map(i => ({ 
    productVariantId: i.product_variant_id || i.productVariantId, 
    quantity: Number(i.quantity ?? 0), 
    ...(i._autoPrice === false ? { manualUnitPrice: { amount: Number(i.unit_price ?? 0), currency: 'EUR' } } : {}),
    manualDiscountPercent: Number(i.discount_percent ?? 0) 
  }));
  
  if (!partyId || !items.length) { 
    previewResult.value = null; 
    return; 
  }
  
  isPreviewLoading.value = true;
  try { 
    const res = await salesApi.previewOrderCalculation(partyId, items);
    if (res) {
      previewResult.value = res;
      // Populate unit prices from pricing engine for auto-priced items
      const serverItems = res.lineItems || res.line_items || [];
      editableOrder.value.line_items.forEach((item, idx) => {
        if (item._autoPrice !== false && serverItems[idx]) {
          const sItem = serverItems[idx];
          item.unit_price = sItem.unitPrice?.amount ?? sItem.unit_price?.amount ?? sItem.unit_price ?? 0;
        }
      });
    }
  } catch (err) {
    console.error('Error en vista previa de cálculos:', err);
    previewResult.value = null;
  } finally { 
    isPreviewLoading.value = false; 
  }
}

async function loadRelatedDocs(data) {
  if (!data.id) return;
  try {
    // 1. Presupuesto origen
    const qId = data.quoteId || data.quote_id;
    if (qId) {
      relatedQuote.value = await salesApi.getQuote(qId);
    } else {
      relatedQuote.value = null;
    }

    // 2. Albaranes relacionados
    const dns = await salesApi.listDeliveryNotes({ orderId: data.id });
    relatedDeliveryNotes.value = dns.data || dns || [];
    
    // 3. Factura relacionada (si hay albaranes, buscamos la del primero)
    const firstDnWithInvoice = relatedDeliveryNotes.value.find(dn => dn.invoiceId || dn.invoice_id);
    if (firstDnWithInvoice) {
      relatedInvoice.value = await salesApi.getInvoice(firstDnWithInvoice.invoiceId || firstDnWithInvoice.invoice_id);
    } else {
      relatedInvoice.value = null;
    }
  } catch (err) {
    console.error('Error cargando documentos relacionados:', err);
  }
}

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
    const val = order.value.subtotal || order.value.subTotal;
    return val?.amount ?? val ?? 0;
  }
  if (previewResult.value && !isPreviewLoading.value) {
    return previewResult.value.subtotal?.amount ?? 0;
  }
  return (editableOrder.value.line_items || []).reduce((acc, line) => {
    // Usamos ?? para que el 0 no sea ignorado en el cálculo local
    const price = Number(line.unit_price ?? line.unitPrice ?? 0);
    const qty = Number(line.quantity ?? 0);
    const disc = Number(line.discount_percent ?? line.discountPercent ?? 0);
    return acc + (qty * price) * (1 - (disc / 100))
  }, 0)
})

const taxAmount = computed(() => {
  if (mode.value === 'detail' && order.value) {
    const val = order.value.taxAmount || order.value.tax_amount || order.value.tax_total;
    return val?.amount ?? val ?? 0;
  }
  if (previewResult.value && !isPreviewLoading.value) {
    return previewResult.value.taxAmount?.amount ?? 0;
  }
  return subtotal.value * 0.21
})

const totalAmount = computed(() => {
  if (mode.value === 'detail' && order.value) {
    const val = order.value.total;
    return val?.amount ?? val ?? 0;
  }
  if (previewResult.value && !isPreviewLoading.value) {
    return previewResult.value.total?.amount ?? 0;
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
    const pId = data.party_id || data.partyId
    if (pId) {
      loadAvailableSetups(pId)
      loadRelatedDocs(data)
      partyApi.getParty(pId).then(p => {
        partyDefaultDiscount.value = p?.default_discount_percentage || 0
      })
    }
    
    // CARGAR ESTADOS MES (NUEVO)
    const mesRefs = data.mes_work_refs || data.mesWorkRefs || []
    if (mesRefs.length > 0) {
      fetchMesWorkStatuses(mesRefs)
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
      product_variant_id: li.product_variant_id || li.productVariantId,
      _autoPrice: false
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
    partyDefaultDiscount.value = party.default_discount_percentage || 0
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

const addProductBtnRef = ref(null)

function focusAddButton() {
  if (addProductBtnRef.value) {
    addProductBtnRef.value.focus()
  }
}

function handleAddLineRequest() {
  showVariantSelector.value = true
}

function handleVariantSelected(payload) {
  const variant = payload.variant || payload
  const newItem = {
    product_variant_id: variant.id,
    variant_sku: variant.sku,
    product_name: variant.product_name || variant.name,
    quantity: 1,
    unit_price: null,
    discount_percent: partyDefaultDiscount.value || 0,
    _autoPrice: true
  }
  editableOrder.value.line_items.push(newItem)
  showVariantSelector.value = false
  
  // Posicionar foco en la cantidad de la nueva línea tras el renderizado
  nextTick(() => {
    fetchPreviewCalculation();
    const lastIdx = editableOrder.value.line_items.length - 1
    const el = document.querySelector(`input[data-row="${lastIdx}"][data-col="qty"]`)
    if (el) {
      el.focus()
      el.select()
    }
  })
}

async function saveOrder() {
  if (!editableOrder.value.party_id) { toastStore.error('Debe seleccionar un cliente'); return; }
  if (!editableOrder.value.line_items.length) { toastStore.error('El pedido debe tener al menos una línea'); return; }
  
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
    toastStore.error('Error al guardar el pedido: ' + (e.message || 'Error desconocido'))
  } finally {
    isSaving.value = false
  }
}

async function cancelOrder() {
  isSaving.value = true
  try {
    await salesApi.cancelOrder(order.value.id)
    await loadOrder()
    toastStore.success('Pedido anulado correctamente')
  } catch (e) {
    toastStore.error('Error al anular: ' + (e.message || 'Error desconocido'))
  } finally {
    isSaving.value = false
  }
}

async function reactivateOrder() {
  isSaving.value = true
  try {
    await salesApi.reactivateOrder(order.value.id)
    await loadOrder()
    toastStore.success('Pedido reactivado')
  } catch (e) {
    toastStore.error('Error al reactivar: ' + (e.message || 'Error desconocido'))
  } finally {
    isSaving.value = false
  }
}

async function launchOrderToProduction() {
  isSaving.value = true
  try {
    await salesApi.changeOrderStatus(order.value.id, 'READY_FOR_PRODUCTION')
    await loadOrder()
    toastStore.success('Pedido lanzado a taller')
  } catch (e) {
    toastStore.error('Error al lanzar a producción: ' + (e.message || 'Error desconocido'))
  } finally {
    isSaving.value = false
  }
}

async function createDeliveryNote() {
  // Preparamos los ítems para el diálogo basándonos en lo pendiente
  dnForm.value.items = (order.value.line_items || order.value.lineItems).map(li => {
    const delivered = li.delivered_quantity || li.deliveredQuantity || 0;
    const pending = li.quantity - delivered;
    return {
      orderLineItemId: li.id,
      productName: li.product_name || li.productName,
      variantSku: li.variant_sku || li.variantSku,
      totalQuantity: li.quantity,
      alreadyDelivered: delivered,
      pendingQuantity: pending,
      quantityToDeliver: pending // Por defecto proponemos entregar todo lo pendiente
    };
  }).filter(item => item.pendingQuantity > 0);

  if (dnForm.value.items.length === 0) {
    toastStore.warning('No hay ítems pendientes de albaranar en este pedido.');
    return;
  }

  showDnDialog.value = true;
}

async function submitDeliveryNote() {
  const itemsToDeliver = dnForm.value.items.filter(i => i.quantityToDeliver > 0);
  if (itemsToDeliver.length === 0) {
    toastStore.warning('Debes indicar al menos una cantidad a entregar.');
    return;
  }

  // Validación de cantidades
  for (const item of itemsToDeliver) {
    if (item.quantityToDeliver > item.pendingQuantity) {
      toastStore.error(`La cantidad a entregar de ${item.productName} no puede superar la pendiente (${item.pendingQuantity}).`);
      return;
    }
  }

  isSaving.value = true
  try {
    const payload = {
      salesOrderId: order.value.id,
      deliveryDate: new Date(dnForm.value.deliveryDate).toISOString(),
      items: itemsToDeliver.map(li => ({
        salesOrderLineItemId: li.orderLineItemId,
        deliveredQuantity: Number(li.quantityToDeliver)
      }))
    };
    const newDn = await salesApi.createDeliveryNote(payload);
    showDnDialog.value = false;
    router.push(`/sales/delivery-notes/${newDn.id}`);
  } catch (e) {
    console.error('Error al generar albarán:', e);
    toastStore.error('Error al generar albarán: ' + (e.message || 'Error desconocido'));
  } finally {
    isSaving.value = false
  }
}

function deliverAll() {
  dnForm.value.items.forEach(item => {
    item.quantityToDeliver = item.pendingQuantity;
  });
}

function deliverNone() {
  dnForm.value.items.forEach(item => {
    item.quantityToDeliver = 0;
  });
}

function updateLines(newLines) {
  editableOrder.value.line_items = newLines
}

function printOrder() {
  window.print()
}

function formatDate(dateString) { return dateString ? new Date(dateString).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '' }
function formatMoney(amount) { return salesApi.formatMoney(amount) }

onMounted(() => {
  window.addEventListener('tramatex-save', handleGlobalSave);
  window.addEventListener('tramatex-esc', handleGlobalEsc);
});

onBeforeUnmount(() => {
  window.removeEventListener('tramatex-save', handleGlobalSave);
  window.removeEventListener('tramatex-esc', handleGlobalEsc);
});

function handleGlobalSave() {
  if (mode.value === 'edit' && !isSaving.value) {
    saveOrder();
  }
}

function handleGlobalEsc() {
  if (mode.value === 'edit') {
    exitEditMode();
  } else {
    router.push('/sales/orders');
  }
}

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
.icon :deep(svg) { width: 24px; height: 24px; }
.icon.blue { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }

.tag-content { display: flex; flex-direction: column; gap: 0.2rem; }
.tag-content label { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }
.tag-content strong { font-size: 1rem; color: var(--color-text-primary); font-weight: 700; }
.amount { color: var(--color-success) !important; font-size: 1.25rem !important; font-family: var(--font-family-mono); }

.status-badge { padding: 0.2rem 0.6rem; border-radius: 4px; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; border: none; }

.clickable-status {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}
.clickable-status:hover {
  filter: brightness(0.9);
  transform: translateY(-1px);
  box-shadow: var(--box-shadow-sm);
}
.clickable-status :deep(svg) {
  width: 0.9rem;
  height: 0.9rem;
  opacity: 0.8;
}

/* ESTILOS DE IMPRESIÓN PROFESIONAL */
.print-container { display: none; }

@media print {
  .no-print { display: none !important; }
  .print-container { display: block !important; position: absolute; left: 0; top: 0; width: 100%; }
}

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>