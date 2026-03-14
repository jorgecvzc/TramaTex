<template>
  <Navbar />
  <div class="order-detail-container">
    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando pedido...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchOrder">Reintentar</button>
      <button class="btn btn-secondary" @click="goBack">Volver</button>
    </div>

    <!-- Order Detail -->
    <div v-else-if="order" class="order-content">
      <!-- Header -->
      <div class="page-header">
        <div>
          <button class="btn-back" @click="goBack">← Volver</button>
          <h1>Pedido {{ order.orderNumber }}</h1>
          <span :class="['status-badge', `status-${salesApi.getStatusClass(order.status)}`]">
            {{ salesApi.getStatusLabel(order.status) }}
          </span>
        </div>
        <div class="header-actions">
          <button
            class="btn btn-secondary"
            @click="printOrder"
            title="Imprimir pedido"
          >
            🖨️ Imprimir
          </button>
          <button 
            v-if="canEdit && !isEditingHeader" 
            class="btn btn-primary" 
            @click="enterHeaderEditMode"
            title="Editar datos del pedido"
          >
            ✏️ Editar
          </button>
          <button 
            v-if="isEditingHeader" 
            class="btn btn-success" 
            @click="saveOrderHeader"
            :disabled="isSavingHeader"
            title="Guardar cambios"
          >
            {{ isSavingHeader ? 'Guardando...' : '💾 Guardar' }}
          </button>
          <button 
            v-if="isEditingHeader" 
            class="btn btn-secondary" 
            @click="cancelHeaderEdit"
            title="Cancelar edición"
          >
            ✕ Cancelar
          </button>
          <button 
            v-if="order.status === 'PENDING'" 
            class="btn btn-success" 
            @click="confirmOrder"
          >
            ✓ Confirmar
          </button>
          <button 
            v-if="canCreateDeliveryNote" 
            class="btn btn-primary" 
            @click="showDeliveryNoteModal = true"
            title="Crear albarán de entrega"
          >
            📦 Crear Albarán
          </button>
          <button 
            v-if="canCreateInvoice" 
            class="btn btn-success" 
            @click="createInvoiceFromOrder"
            :disabled="isCreatingInvoice"
            title="Crear factura para este pedido"
          >
            🧾 {{ isCreatingInvoice ? 'Creando...' : 'Crear Factura' }}
          </button>
          <button 
            v-if="canCancel" 
            class="btn btn-danger" 
            @click="cancelOrder"
          >
            Cancelar Pedido
          </button>
          <button
            v-if="order.status === 'CANCELLED'"
            class="btn btn-success"
            @click="reactivateOrder"
          >
            ♻️ Reactivar Pedido
          </button>
        </div>
      </div>

      <div class="print-doc-header">
        <p class="print-brand">{{ issuerProfile.displayName }}</p>
        <p v-if="issuerProfile.taxId" class="print-issuer-line">{{ issuerProfile.taxLabel }}: {{ issuerProfile.taxId }}</p>
        <p v-if="issuerProfile.addressLine || issuerProfile.cityLine" class="print-issuer-line">{{ issuerProfile.addressLine }}<span v-if="issuerProfile.addressLine && issuerProfile.cityLine"> · </span>{{ issuerProfile.cityLine }}</p>
        <p v-if="issuerProfile.contactLine" class="print-issuer-line">{{ issuerProfile.contactLine }}</p>
        <h2>Pedido {{ order.orderNumber }}</h2>
        <div class="print-doc-meta">
          <span>Cliente: {{ formatPartyId(order.partyId) }}</span>
          <span>Fecha pedido: {{ formatDate(order.orderDate) }}</span>
          <span>Estado: {{ salesApi.getStatusLabel(order.status) }}</span>
        </div>
      </div>

      <!-- Order Info Cards -->
      <div class="info-grid">
        <div class="info-card">
          <h3>Información General</h3>
          <div class="info-row">
            <span class="label">Cliente:</span>
            <span class="value">{{ formatPartyId(order.partyId) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Fecha de Pedido:</span>
            <span class="value">{{ formatDate(order.orderDate) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Fecha de Entrega:</span>
            <span v-if="!isEditingHeader" class="value">{{ formatDate(order.deliveryDate) }}</span>
            <input v-else v-model="headerEditForm.deliveryDate" type="date" class="form-input form-input-inline" />
          </div>
        </div>

        <div class="info-card">
          <h3>Totales</h3>
          <div class="info-row">
            <span class="label">Subtotal:</span>
            <span class="value">{{ salesApi.formatMoney(order.subtotal) }}</span>
          </div>
          <div class="info-row">
            <span class="label">IVA:</span>
            <span class="value">{{ salesApi.formatMoney(order.taxAmount) }}</span>
          </div>
          <div class="info-row total">
            <span class="label">Total:</span>
            <span class="value">{{ salesApi.formatMoney(order.total) }}</span>
          </div>
        </div>
      </div>

      <!-- Notes -->
      <div v-if="order.notes || isEditingHeader" class="notes-card">
        <h3>Notas</h3>
        <p v-if="!isEditingHeader">{{ order.notes }}</p>
        <textarea v-else v-model="headerEditForm.notes" class="form-textarea" rows="3" placeholder="Notas del pedido..."></textarea>
      </div>

      <!-- MES Work References (Document-level) -->
      <div v-if="(order.mesWorkRefs && order.mesWorkRefs.length > 0) || isEditingHeader" class="notes-card">
        <h3>Trabajos MES</h3>
        <template v-if="!isEditingHeader">
          <div v-for="mesRef in order.mesWorkRefs" :key="mesRef.mesWorkId" class="mes-ref-view">
            <button class="mes-link" @click="goToMesWork(mesRef.mesWorkId)" type="button"
              :title="getMesWorkTooltip(mesRef.mesWorkId)">
              {{ formatMesWorkId(mesRef.mesWorkId) }}
            </button>
            <span v-if="mesRef.observations" class="mes-observations-text">— {{ mesRef.observations }}</span>
          </div>
        </template>
        <template v-else>
          <div v-if="isLoadingMesWorks" class="mes-loading">Cargando trabajos MES...</div>
          <div v-else-if="mesWorks.length === 0" class="mes-empty">No hay trabajos MES disponibles para este cliente.</div>
          <div v-else class="mes-ref-list">
            <div v-for="work in mesWorks" :key="work.id" class="mes-ref-item">
              <label class="mes-checkbox-label">
                <input type="checkbox" :checked="isEditMesWorkSelected(work.id)" @change="toggleEditMesWork(work.id)" />
                <span>{{ work.work_number }} - {{ work.work_name }}</span>
              </label>
              <textarea v-if="isEditMesWorkSelected(work.id)" class="mes-observations" rows="2"
                placeholder="Observaciones para este trabajo MES..."
                :value="getEditMesWorkObservations(work.id)"
                @input="setEditMesWorkObservations(work.id, $event.target.value)"
              ></textarea>
            </div>
          </div>
        </template>
      </div>

      <!-- Delivery Notes Section -->
      <div v-if="deliveryNotes.length > 0" class="delivery-notes-section">
        <h3>Albaranes de Entrega</h3>
        <div class="delivery-notes-list">
          <div 
            v-for="note in deliveryNotes" 
            :key="note.id" 
            class="delivery-note-item"
            @click="navigateToDeliveryNote(note.id)"
          >
            <div class="note-icon">📦</div>
            <div class="note-info">
              <span class="note-number">{{ note.deliveryNoteNumber }}</span>
              <span class="note-date">{{ formatDate(note.deliveryDate) }}</span>
            </div>
            <div class="note-action">→</div>
          </div>
        </div>
      </div>

      <!-- Line Items -->
      <div class="line-items-section">
        <div class="section-header">
          <h2>Líneas del Pedido</h2>
          <button 
            v-if="canEdit && isEditingHeader" 
            class="btn btn-primary" 
            @click="addEditLineItem"
          >
            + Agregar Línea
          </button>
        </div>

        <!-- View mode -->
        <template v-if="!isEditingHeader">
        <div v-if="!order.lineItems || order.lineItems.length === 0" class="empty-state">
          <p>No hay líneas en este pedido</p>
        </div>

        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Referencia</th>
                <th>Nombre</th>
                <th>Cantidad</th>
                <th>P. Tarifa</th>
                <th>P. Venta</th>
                <th>Dto. %</th>
                <th>Subtotal</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in order.lineItems" :key="item.id">
                <td class="variant-id">
                  <span v-if="item.variantSku">{{ item.variantSku }}</span>
                  <span v-else>{{ formatVariantId(item.productVariantId) }}</span>
                </td>
                <td>{{ buildDisplayName(item) }}</td>
                <td>{{ item.quantity }}</td>
                <td>
                  {{ salesApi.formatUnitPrice(item.listUnitPrice) }}
                </td>
                <td>
                  {{ salesApi.formatUnitPrice(item.unitPrice) }}
                </td>
                <td>{{ item.discountPercent ? item.discountPercent.toFixed(2) + '%' : '—' }}</td>
                <td class="amount">{{ salesApi.formatMoney(item.subtotal) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        </template>

        <!-- Edit mode (inline) -->
        <template v-else>
        <div v-if="editLineItems.length === 0" class="empty-state">
          <p>No hay líneas. Agregue al menos una línea.</p>
        </div>

        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Referencia</th>
                <th>Nombre</th>
                <th>Cantidad</th>
                <th>P. Tarifa</th>
                <th>P. Venta</th>
                <th>Dto. %</th>
                <th>Subtotal</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in editLineItems" :key="idx">
                <td class="variant-id">
                  <span v-if="item.variantSku">{{ item.variantSku }}</span>
                  <span v-else>{{ formatVariantId(item.productVariantId) }}</span>
                </td>
                <td>{{ item.displayName || '—' }}</td>
                <td>
                  <input v-model.number="item.quantity" type="number" min="1" class="form-input form-input-small" />
                </td>
                <td class="col-readonly">
                  {{ item.listPrice != null ? item.listPrice.toFixed(3) : '—' }}
                </td>
                <td>
                  <input v-model.number="item.unitPrice" type="number" step="0.001" min="0" placeholder="Auto" class="form-input form-input-small" />
                </td>
                <td>
                  <input v-model.number="item.discountPercent" type="number" step="0.01" min="0" max="100" placeholder="0" class="form-input form-input-small" />
                </td>
                <td class="col-subtotal">
                  <span v-if="isPreviewLoading" class="subtotal-loading">…</span>
                  <span v-else class="subtotal-value">{{ formatMoneyAmount(calculateEditLineSubtotal(item)) }}</span>
                </td>
                <td class="actions-cell">
                  <button type="button" class="btn-icon danger" @click="removeEditLineItem(idx)" title="Eliminar">🗑️</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Totals Summary -->
        <div v-if="editLineItems.length > 0" class="totals-section">
          <h3>Resumen de Totales</h3>
          <div class="totals-grid">
            <div class="total-row">
              <span class="total-label">Subtotal:</span>
              <span class="total-value">{{ formatMoneyAmount(editCalculatedTotals.subtotal) }}</span>
            </div>
            <div class="total-row">
              <span class="total-label">IVA:</span>
              <span class="total-value">{{ formatMoneyAmount(editCalculatedTotals.tax) }}</span>
            </div>
            <div class="total-row total-final">
              <span class="total-label">Total:</span>
              <span class="total-value">{{ formatMoneyAmount(editCalculatedTotals.total) }}</span>
            </div>
          </div>
        </div>
        </template>
      </div>

      <div class="print-doc-footer">
        <span>Documento generado por {{ issuerProfile.displayName }}</span>
        <span>Entrega prevista: {{ formatDate(order.deliveryDate) }}</span>
      </div>
    </div>

    <!-- Delivery Note Creation Modal -->
    <div v-if="showDeliveryNoteModal" class="modal-overlay" @click="showDeliveryNoteModal = false">
      <div class="modal-content modal-wide" @click.stop>
        <div class="modal-header">
          <h3>Crear Albarán de Entrega</h3>
          <button class="btn-close" @click="showDeliveryNoteModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="deliveryDate">Fecha de Entrega *</label>
            <input
              id="deliveryDate"
              v-model="deliveryNoteForm.deliveryDate"
              type="date"
              class="form-input"
              :min="minDeliveryDate"
              required
            />
          </div>
          
          <div class="form-group">
            <label>Tipo de Entrega *</label>
            <div class="radio-group">
              <label class="radio-label">
                <input
                  v-model="deliveryNoteForm.type"
                  type="radio"
                  value="TOTAL"
                />
                <span>Entrega Total</span>
              </label>
              <label class="radio-label">
                <input
                  v-model="deliveryNoteForm.type"
                  type="radio"
                  value="PARTIAL"
                />
                <span>Entrega Parcial</span>
              </label>
            </div>
          </div>

          <!-- Partial Delivery: Item Selection -->
          <div v-if="deliveryNoteForm.type === 'PARTIAL'" class="form-section">
            <h4>Seleccionar Items a Entregar</h4>
            <div class="items-selection">
              <div 
                v-for="(item, index) in deliveryNoteForm.items" 
                :key="item.lineItemId" 
                class="item-selection-row"
              >
                <div class="item-info">
                  <input
                    v-model="item.selected"
                    type="checkbox"
                    :id="`item-${index}`"
                  />
                  <label :for="`item-${index}`" class="item-label">
                    <span class="item-sku">{{ item.variantSku || formatVariantId(item.productVariantId) }}</span>
                    <span class="item-name">{{ item.productName }}</span>
                    <span class="item-available">Disponible: {{ item.availableQuantity }}</span>
                  </label>
                </div>
                <div v-if="item.selected" class="item-quantity">
                  <label>Cantidad:</label>
                  <input
                    v-model.number="item.quantityToDeliver"
                    type="number"
                    min="1"
                    :max="item.availableQuantity"
                    class="form-input-small"
                  />
                </div>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label for="deliveryAddress">Dirección de Entrega</label>
            <textarea
              id="deliveryAddress"
              v-model="deliveryNoteForm.deliveryAddress"
              class="form-textarea"
              rows="3"
              placeholder="Calle, ciudad, código postal..."
            ></textarea>
          </div>

          <div class="form-group">
            <label for="notes">Notas / Observaciones</label>
            <textarea
              id="notes"
              v-model="deliveryNoteForm.notes"
              class="form-textarea"
              rows="2"
              placeholder="Observaciones sobre la entrega..."
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showDeliveryNoteModal = false">Cancelar</button>
          <button 
            class="btn btn-primary" 
            @click="createDeliveryNote"
            :disabled="!isDeliveryNoteFormValid || isCreatingDeliveryNote"
          >
            {{ isCreatingDeliveryNote ? 'Creando...' : 'Crear Albarán' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Variant Selector Modal (Edit mode) -->
    <div v-if="showVariantSelector" class="modal-overlay" @click.self="showVariantSelector = false">
      <div class="modal-content modal-wide" @click.stop>
        <div class="modal-header">
          <h3>Seleccionar Variante de Producto</h3>
          <button class="btn-close" @click="showVariantSelector = false">✕</button>
        </div>
        <div class="modal-body">
          <VariantSelector
            :product-id="null"
            title=""
            description="Seleccione una variante de producto"
            @variant-selected="handleVariantSelected"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { mesApi } from '@/services/mesApi';
import { getPrintIssuerProfile } from '@/services/printIssuerProfile';
import { calculateBaseSalesPrice } from '@/services/pricingApi';
import '@/assets/sales-print.css';

const route = useRoute();
const router = useRouter();

const order = ref(null);
const isLoading = ref(false);
const error = ref('');
const deliveryNotes = ref([]);
const mesWorksCache = ref({});
const mesWorks = ref([]);
const isLoadingMesWorks = ref(false);
const issuerProfile = getPrintIssuerProfile();

const showVariantSelector = ref(false);
const partyDefaultDiscount = ref(null);

const showDeliveryNoteModal = ref(false);
const isCreatingDeliveryNote = ref(false);
const deliveryNoteForm = ref({
  deliveryDate: '',
  type: 'TOTAL',
  items: [],
  deliveryAddress: '',
  notes: '',
});

// Inline edit state for line items (populated on enterHeaderEditMode)
const editLineItems = ref([]);
const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewDebounceTimer = null;

// Header edit state
const isEditingHeader = ref(false);
const isSavingHeader = ref(false);
const headerEditForm = ref({
  deliveryDate: '',
  notes: '',
  mesWorkRefs: [],
});

const canEdit = computed(() => {
  return order.value && ['PENDING', 'CONFIRMED'].includes(order.value.status);
});

const canCancel = computed(() => {
  return order.value && ['PENDING', 'CONFIRMED'].includes(order.value.status);
});

const canCreateDeliveryNote = computed(() => {
  return order.value && ['CONFIRMED', 'PARTIALLY_DELIVERED'].includes(order.value.status);
});

const canCreateInvoice = computed(() => {
  return order.value && ['DELIVERED', 'PARTIALLY_INVOICED'].includes(order.value.status);
});

const isCreatingInvoice = ref(false);

const minDeliveryDate = computed(() => {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().split('T')[0];
});

const isDeliveryNoteFormValid = computed(() => {
  if (!deliveryNoteForm.value.deliveryDate) return false;
  
  if (deliveryNoteForm.value.type === 'PARTIAL') {
    const selectedItems = deliveryNoteForm.value.items.filter(item => item.selected);
    if (selectedItems.length === 0) return false;
    
    return selectedItems.every(item => 
      item.quantityToDeliver > 0 && 
      item.quantityToDeliver <= item.availableQuantity
    );
  }
  
  return true;
});

onMounted(() => {
  fetchOrder();
});

async function fetchOrder() {
  const orderId = route.params.id;
  if (!orderId) {
    error.value = 'ID de pedido no válido';
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    order.value = await salesApi.getOrder(orderId);
    await loadDeliveryNotes();
    await loadMesWorksForOrder();
    await loadPartyDiscount();
    initializeDeliveryNoteForm();
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar el pedido';
    console.error('Error loading order:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyDiscount() {
  if (!order.value?.partyId) return;
  try {
    const party = await partyApi.getParty(order.value.partyId);
    partyDefaultDiscount.value = party.default_discount_percentage || null;
  } catch (err) {
    console.warn('[OrderDetail] Could not load party discount:', err.message);
  }
}

async function loadDeliveryNotes() {
  if (!order.value?.id) return;
  
  try {
    const response = await salesApi.listDeliveryNotes({ orderId: order.value.id });
    deliveryNotes.value = Array.isArray(response) ? response : (response.data || []);
  } catch (err) {
    console.error('Error loading delivery notes:', err);
    // Non-critical, don't show error to user
  }
}

async function loadMesWorksForOrder() {
  const refs = order.value?.mesWorkRefs;
  if (!Array.isArray(refs) || refs.length === 0) return;

  const mesWorkIds = refs.map(r => r.mesWorkId).filter(Boolean);
  const uncachedIds = mesWorkIds.filter(id => !mesWorksCache.value[id]);
  if (uncachedIds.length === 0) return;

  const results = await Promise.allSettled(uncachedIds.map(id => mesApi.getWorkDefinition(id)));
  results.forEach((result, index) => {
    const mesWorkId = uncachedIds[index];
    mesWorksCache.value[mesWorkId] = result.status === 'fulfilled' ? result.value : null;
  });
}

function initializeDeliveryNoteForm() {
  if (!order.value?.lineItems) return;
  
  // Calculate already-delivered quantities from existing delivery notes
  const deliveredByLineItem = {};
  for (const note of deliveryNotes.value) {
    if (note.status === 'CANCELLED') continue;
    for (const noteItem of (note.lineItems || [])) {
      const key = noteItem.salesOrderLineItemId;
      deliveredByLineItem[key] = (deliveredByLineItem[key] || 0) + noteItem.deliveredQuantity;
    }
  }

  const today = new Date().toISOString().split('T')[0];
  deliveryNoteForm.value.deliveryDate = today;
  deliveryNoteForm.value.type = 'TOTAL';
  deliveryNoteForm.value.items = order.value.lineItems.map(item => {
    const alreadyDelivered = deliveredByLineItem[item.id] || 0;
    const available = Math.max(0, item.quantity - alreadyDelivered);
    return {
      lineItemId: item.id,
      productVariantId: item.productVariantId,
      productName: item.productName || '',
      variantSku: item.variantSku || '',
      availableQuantity: available,
      quantityToDeliver: available,
      selected: false,
    };
  });
}

async function createDeliveryNote() {
  if (!isDeliveryNoteFormValid.value || isCreatingDeliveryNote.value) return;
  
  isCreatingDeliveryNote.value = true;
  
  try {
    let items;
    if (deliveryNoteForm.value.type === 'PARTIAL') {
      items = deliveryNoteForm.value.items
        .filter(item => item.selected)
        .map(item => ({
          salesOrderLineItemId: item.lineItemId,
          deliveredQuantity: item.quantityToDeliver,
        }));
    } else {
      // TOTAL: include all line items with remaining quantity (skip fully delivered)
      items = deliveryNoteForm.value.items
        .filter(item => item.availableQuantity > 0)
        .map(item => ({
          salesOrderLineItemId: item.lineItemId,
          deliveredQuantity: item.availableQuantity,
        }));
    }

    const deliveryNoteData = {
      salesOrderId: order.value.id,
      deliveryDate: salesApi.formatDateForAPI(new Date(deliveryNoteForm.value.deliveryDate)),
      items,
    };
    
    if (deliveryNoteForm.value.notes) {
      deliveryNoteData.notes = deliveryNoteForm.value.notes;
    }
    
    await salesApi.createDeliveryNote(deliveryNoteData);
    
    // Refresh order and delivery notes
    await fetchOrder();
    
    // Close modal and reset form
    showDeliveryNoteModal.value = false;
    initializeDeliveryNoteForm();
    
    alert('✓ Albarán creado exitosamente');
  } catch (err) {
    alert(err?.message || 'No se pudo crear el albarán');
    console.error('Error creating delivery note:', err);
  } finally {
    isCreatingDeliveryNote.value = false;
  }
}

function navigateToDeliveryNote(noteId) {
  router.push(`/sales/delivery-notes/${noteId}`);
}

async function createInvoiceFromOrder() {
  if (!confirm('¿Crear factura para este pedido?')) return;
  
  isCreatingInvoice.value = true;
  try {
    const now = new Date();
    const dueDate = new Date(now);
    dueDate.setDate(dueDate.getDate() + 30);

    const data = {
      partyId: order.value.partyId,
      salesOrderIds: [order.value.id],
      invoiceDate: now.toISOString(),
      dueDate: dueDate.toISOString(),
    };
    
    const newInvoice = await salesApi.createInvoice(data);
    router.push(`/sales/invoices/${newInvoice.id}`);
  } catch (err) {
    alert(err?.message || 'No se pudo crear la factura');
    console.error('Error creating invoice:', err);
  } finally {
    isCreatingInvoice.value = false;
  }
}

async function confirmOrder() {
  if (!confirm('¿Confirmar este pedido?')) return;

  try {
    await salesApi.changeOrderStatus(order.value.id, 'CONFIRMED');
    await fetchOrder();
  } catch (err) {
    alert(err?.message || 'No se pudo confirmar el pedido');
  }
}

async function cancelOrder() {
  if (!confirm('¿Cancelar este pedido? Esta acción no se puede deshacer.')) return;

  try {
    await salesApi.changeOrderStatus(order.value.id, 'CANCELLED');
    await fetchOrder();
  } catch (err) {
    alert(err?.message || 'No se pudo cancelar el pedido');
  }
}

async function reactivateOrder() {
  if (!confirm('¿Reactivar este pedido? Volverá a estado Pendiente.')) return;

  try {
    await salesApi.changeOrderStatus(order.value.id, 'PENDING');
    await fetchOrder();
  } catch (err) {
    alert(err?.message || 'No se pudo reactivar el pedido');
  }
}

function addEditLineItem() {
  showVariantSelector.value = true;
}

function handleVariantSelected(payload) {
  const variant = payload?.variant || payload;
  if (variant) {
    const name = variant.product_name || '';
    const config = variant.option_configuration;
    let displayName = name || '—';
    if (config && Object.keys(config).length > 0) {
      displayName = name + ' - ' + Object.values(config).join(', ');
    }
    const newItem = {
      id: null,
      productVariantId: variant.id,
      variantSku: variant.sku || '',
      displayName,
      quantity: 1,
      unitPrice: null,
      discountPercent: partyDefaultDiscount.value || null,
      effectiveUnitPrice: 0,
    };
    editLineItems.value.push(newItem);
    showVariantSelector.value = false;
    fetchEditLinePrice(newItem, variant.product_id || '', variant.product_base_price ?? variant.base_cost);
  }
}

async function fetchEditLinePrice(item, productId, basePrice) {
  if (!item.productVariantId || !productId) {
    if (basePrice != null) {
      item.unitPrice = Math.round(basePrice * 1000) / 1000;
    }
    return;
  }
  try {
    const result = await calculateBaseSalesPrice(productId, item.productVariantId);
    const rawPrice = result.baseSalesPrice?.amount ?? basePrice ?? null;
    item.unitPrice = rawPrice != null ? Math.round(rawPrice * 1000) / 1000 : 0;
  } catch (err) {
    console.warn('[OrderDetail] Error fetching sale price:', err.message);
    if (basePrice != null) {
      item.unitPrice = Math.round(basePrice * 1000) / 1000;
    }
  }
}

function removeEditLineItem(index) {
  editLineItems.value.splice(index, 1);
}

function calculateEditLineSubtotal(item) {
  if (!previewResult.value) return 0;
  const match = previewResult.value.lineItems.find(
    li => li.productVariantId === item.productVariantId
  );
  return match?.subtotal?.amount ?? 0;
}

const editCalculatedTotals = computed(() => {
  if (!previewResult.value) return { subtotal: 0, tax: 0, total: 0 };
  return {
    subtotal: previewResult.value.subtotal?.amount ?? 0,
    tax: previewResult.value.taxAmount?.amount ?? 0,
    total: previewResult.value.total?.amount ?? 0,
  };
});

function debouncedPreview() {
  clearTimeout(previewDebounceTimer);
  previewDebounceTimer = setTimeout(fetchPreviewCalculation, 400);
}

function buildPreviewItems() {
  return editLineItems.value
    .filter(item => item.productVariantId)
    .map(item => ({
      productVariantId: item.productVariantId,
      quantity: item.quantity || 1,
      unitPrice: item.unitPrice ? { amount: item.unitPrice, currency: 'EUR' } : undefined,
      discountPercent: item.discountPercent != null ? item.discountPercent : undefined,
    }));
}

async function fetchPreviewCalculation() {
  const partyId = order.value?.partyId;
  const items = buildPreviewItems();
  if (!partyId || items.length === 0) {
    previewResult.value = null;
    return;
  }
  isPreviewLoading.value = true;
  try {
    previewResult.value = await salesApi.previewOrderCalculation(partyId, items);
  } catch (err) {
    console.warn('[OrderDetail] Preview calculation error:', err.message);
  } finally {
    isPreviewLoading.value = false;
  }
}

watch(
  () => editLineItems.value.map(i => `${i.productVariantId}|${i.quantity}|${i.unitPrice}|${i.discountPercent}`).join(','),
  () => {
    if (isEditingHeader.value) debouncedPreview();
  }
);

function formatMoneyAmount(amount) {
  return new Intl.NumberFormat('es-ES', {
    style: 'currency',
    currency: 'EUR',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

function enterHeaderEditMode() {
  headerEditForm.value.deliveryDate = order.value.deliveryDate 
    ? new Date(order.value.deliveryDate).toISOString().split('T')[0] 
    : '';
  headerEditForm.value.notes = order.value.notes || '';
  headerEditForm.value.mesWorkRefs = (order.value.mesWorkRefs || []).map(r => ({
    mesWorkId: r.mesWorkId,
    observations: r.observations || '',
  }));
  // Populate inline edit line items from current order
  editLineItems.value = (order.value.lineItems || []).map(item => ({
    id: item.id,
    productVariantId: item.productVariantId,
    variantSku: item.variantSku || '',
    displayName: buildDisplayName(item),
    quantity: item.quantity,
    listPrice: item.listUnitPrice?.amount || null,
    unitPrice: item.unitPrice?.amount || null,
    discountPercent: item.discountPercent ?? null,
    effectiveUnitPrice: item.unitPrice?.amount || 0,
  }));
  isEditingHeader.value = true;
  loadMesWorksForParty();
}

function cancelHeaderEdit() {
  isEditingHeader.value = false;
  showVariantSelector.value = false;
  editLineItems.value = [];
  previewResult.value = null;
}

async function loadMesWorksForParty() {
  const partyId = order.value?.partyId;
  if (!partyId) return;
  isLoadingMesWorks.value = true;
  try {
    mesWorks.value = await mesApi.listWorkDefinitions({ party_id: partyId });
  } catch (err) {
    console.error('Error loading MES works for party:', err);
    mesWorks.value = [];
  } finally {
    isLoadingMesWorks.value = false;
  }
}

function isEditMesWorkSelected(workId) {
  return headerEditForm.value.mesWorkRefs.some(r => r.mesWorkId === workId);
}

function toggleEditMesWork(workId) {
  const idx = headerEditForm.value.mesWorkRefs.findIndex(r => r.mesWorkId === workId);
  if (idx >= 0) {
    headerEditForm.value.mesWorkRefs.splice(idx, 1);
  } else {
    headerEditForm.value.mesWorkRefs.push({ mesWorkId: workId, observations: '' });
  }
}

function getEditMesWorkObservations(workId) {
  return headerEditForm.value.mesWorkRefs.find(r => r.mesWorkId === workId)?.observations || '';
}

function setEditMesWorkObservations(workId, value) {
  const ref = headerEditForm.value.mesWorkRefs.find(r => r.mesWorkId === workId);
  if (ref) ref.observations = value;
}

async function saveOrderHeader() {
  if (editLineItems.value.length === 0) {
    alert('Debe haber al menos una línea en el pedido');
    return;
  }

  isSavingHeader.value = true;

  try {
    // 1. Update header (deliveryDate, notes, mesWorkRefs)
    const updateData = {
      notes: headerEditForm.value.notes || undefined,
      mesWorkRefs: headerEditForm.value.mesWorkRefs,
    };

    if (headerEditForm.value.deliveryDate) {
      updateData.deliveryDate = salesApi.formatDateForAPI(new Date(headerEditForm.value.deliveryDate));
    }

    await salesApi.updateOrder(order.value.id, updateData);

    // 2. Sync line items: remove deleted, update existing, add new
    const originalIds = new Set((order.value.lineItems || []).map(i => i.id));
    const editIds = new Set(editLineItems.value.filter(i => i.id).map(i => i.id));

    // Remove lines that were deleted in edit mode
    for (const origId of originalIds) {
      if (!editIds.has(origId)) {
        await salesApi.removeOrderLineItem(order.value.id, origId);
      }
    }

    // Update existing lines and add new ones
    for (const item of editLineItems.value) {
      const lineData = {
        productVariantId: item.productVariantId,
        quantity: item.quantity,
      };

      if (item.unitPrice != null && item.unitPrice > 0) {
        lineData.unitPrice = {
          amount: item.unitPrice,
          currency: 'EUR',
        };
      }

      if (item.discountPercent != null) {
        lineData.discountPercent = item.discountPercent;
      }

      if (item.id) {
        // Update existing line
        await salesApi.updateOrderLineItem(order.value.id, item.id, lineData);
      } else {
        // Add new line
        await salesApi.addOrderLineItem(order.value.id, lineData);
      }
    }

    isEditingHeader.value = false;
    editLineItems.value = [];
    await fetchOrder();
  } catch (err) {
    alert(err?.message || 'No se pudo actualizar el pedido');
  } finally {
    isSavingHeader.value = false;
  }
}

function goBack() {
  router.push('/sales/orders');
}

function printOrder() {
  window.print();
}

function formatDate(dateString) {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return date.toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

function formatPartyId(partyId) {
  if (!partyId) return '—';
  return partyId.substring(0, 8) + '...';
}

function formatVariantId(variantId) {
  if (!variantId) return '—';
  return variantId.substring(0, 8) + '...';
}

function buildDisplayName(item) {
  const name = item.productName || '';
  const config = item.optionConfiguration;
  if (!name) return '—';
  if (!config || Object.keys(config).length === 0) return name;
  return name + ' - ' + Object.values(config).join(', ');
}

function formatMesWorkId(mesWorkId) {
  if (!mesWorkId) return '—';
  const mesWork = mesWorksCache.value[mesWorkId];
  if (mesWork?.work_number) return mesWork.work_number;
  return mesWorkId.substring(0, 8) + '...';
}

function getMesWorkTooltip(mesWorkId) {
  if (!mesWorkId) return 'Definición de trabajo MES';
  const mesWork = mesWorksCache.value[mesWorkId];
  if (mesWork?.work_number && mesWork?.work_name) {
    return `${mesWork.work_number} · ${mesWork.work_name}`;
  }
  if (mesWork?.work_number) return mesWork.work_number;
  return `Definición MES ${mesWorkId}`;
}

function goToMesWork(mesWorkId) {
  if (!mesWorkId) return;
  router.push(`/mes/work-definitions/${mesWorkId}`);
}
</script>

<style scoped>
.order-detail-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 3rem 1rem;
  background: white;
  border-radius: 8px;
}

.spinner {
  width: 40px;
  height: 40px;
  margin: 0 auto 1rem;
  border: 3px solid #f3f4f6;
  border-top-color: #E6B800;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-message {
  color: #dc2626;
  margin-bottom: 1rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0.5rem 0;
}

.btn-back {
  background: transparent;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  transition: color 0.2s;
}

.btn-back:hover {
  color: #1f2937;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.info-card,
.notes-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.info-card h3,
.notes-card h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  font-size: 0.875rem;
}

.info-row.total {
  margin-top: 0.5rem;
  padding-top: 0.75rem;
  border-top: 2px solid #f3f4f6;
  font-weight: 600;
  font-size: 1rem;
}

.info-row .label {
  color: #6b7280;
}

.info-row .value {
  color: #1f2937;
  font-weight: 500;
}

.notes-card {
  margin-bottom: 2rem;
}

.notes-card p {
  font-size: 0.875rem;
  color: #4b5563;
  line-height: 1.6;
}

.line-items-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #E6B800;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #d4a700;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f3f4f6;
  color: #4a5568;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.btn-success {
  background: #10b981;
  color: white;
}

.btn-success:hover {
  background: #059669;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.status-warning {
  background: #fef3c7;
  color: #92400e;
}

.status-info {
  background: #dbeafe;
  color: #1e40af;
}

.status-primary {
  background: #e0e7ff;
  color: #3730a3;
}

.status-success {
  background: #d1fae5;
  color: #065f46;
}

.status-danger {
  background: #fee2e2;
  color: #991b1b;
}

.status-secondary {
  background: #f3f4f6;
  color: #4b5563;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table thead {
  background: #f9fafb;
}

.data-table th {
  text-align: left;
  padding: 0.75rem 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #6b7280;
  letter-spacing: 0.05em;
}

.data-table td {
  padding: 1rem;
  border-top: 1px solid #f3f4f6;
  font-size: 0.875rem;
  color: #1f2937;
}

.variant-id {
  color: #374151;
}

.variant-sku {
  color: #6b7280;
  font-family: 'Courier New', monospace;
  font-size: 0.85em;
}

.amount {
  font-weight: 600;
  text-align: right;
}

.mes-link {
  background: transparent;
  border: none;
  color: #2563eb;
  cursor: pointer;
  padding: 0;
  font-size: 0.875rem;
  text-decoration: underline;
}

.mes-link:hover {
  color: #1d4ed8;
}

.mes-ref-view {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  padding: 0.25rem 0;
}

.mes-observations-text {
  color: #6b7280;
  font-size: 0.875rem;
}

.mes-ref-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.mes-ref-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.mes-checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background 0.15s, border-color 0.15s;
}

.mes-checkbox-label:hover {
  background: #f0f7ff;
  border-color: #93c5fd;
}

.mes-checkbox-label input[type="checkbox"] {
  accent-color: #2563eb;
}

.mes-observations {
  margin-left: 1.75rem;
  padding: 0.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 0.875rem;
  resize: vertical;
}

.mes-loading,
.mes-empty {
  color: #6b7280;
  font-size: 0.875rem;
  font-style: italic;
  padding: 0.5rem 0;
}

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

.col-subtotal {
  text-align: right;
}

.subtotal-value {
  font-size: 0.875rem;
  color: #1f2937;
  font-weight: 600;
  white-space: nowrap;
}

.totals-section {
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  margin-top: 1rem;
}

.totals-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1rem;
}

.totals-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.total-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.total-row.total-final {
  margin-top: 0.5rem;
  padding-top: 0.75rem;
  border-top: 2px solid #d1d5db;
  font-size: 1.125rem;
  font-weight: 700;
}

.total-label {
  color: #6b7280;
}

.total-row.total-final .total-label {
  color: #1f2937;
}

.total-value {
  color: #1f2937;
  font-weight: 600;
}

.total-row.total-final .total-value {
  color: #E6B800;
}

.btn-icon {
  background: transparent;
  border: none;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  font-size: 1.25rem;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.btn-icon:hover {
  opacity: 1;
}

.btn-icon.danger:hover {
  color: #dc2626;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #9ca3af;
}

/* Delivery Notes Section */
.delivery-notes-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.delivery-notes-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
}

.delivery-notes-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.delivery-note-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.delivery-note-item:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
  transform: translateX(4px);
}

.note-icon {
  font-size: 1.5rem;
}

.note-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.note-number {
  font-weight: 600;
  color: #1f2937;
  font-size: 0.875rem;
}

.note-date {
  font-size: 0.75rem;
  color: #6b7280;
}

.note-action {
  color: #3b82f6;
  font-size: 1.25rem;
}

/* Modal Styles */
.modal-wide {
  max-width: 700px;
}

.form-section {
  margin: 1.5rem 0;
  padding: 1rem;
  background: #f9fafb;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.form-section h4 {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1rem;
}

.radio-group {
  display: flex;
  gap: 1.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  color: #4a5568;
}

.radio-label input[type="radio"] {
  cursor: pointer;
}

.items-selection {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.item-selection-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
}

.item-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
}

.item-info input[type="checkbox"] {
  cursor: pointer;
}

.item-label {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  cursor: pointer;
}

.item-sku {
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
  color: #1f2937;
  font-weight: 500;
}

.item-name {
  font-size: 0.8rem;
  color: #374151;
}

.item-available {
  font-size: 0.75rem;
  color: #6b7280;
}

.item-quantity {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.item-quantity label {
  font-size: 0.75rem;
  color: #6b7280;
  white-space: nowrap;
}

.form-input-small {
  width: 80px;
  padding: 0.375rem 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
}

.form-input-small:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.form-textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: inherit;
  resize: vertical;
}

.form-textarea:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-wide {
  max-width: 900px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #f3f4f6;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
}

.btn-close {
  background: transparent;
  border: none;
  font-size: 1.5rem;
  color: #9ca3af;
  cursor: pointer;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.modal-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.25rem;
}

.form-input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
}

.form-input:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.form-input:disabled {
  background: #f3f4f6;
  color: #9ca3af;
  cursor: not-allowed;
}

.modal-footer {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
  padding: 1.5rem;
  border-top: 1px solid #f3f4f6;
}

/* Print styles: see @/assets/sales-print.css */

.form-input-inline {
  max-width: 180px;
  padding: 0.25rem 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
}

.form-input-inline:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.tax-notice {
  margin-top: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: #fffbea;
  border-left: 3px solid #E6B800;
  border-radius: 4px;
  color: #92710c;
}
</style>
