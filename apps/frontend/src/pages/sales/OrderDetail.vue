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
            v-if="canCancel" 
            class="btn btn-danger" 
            @click="cancelOrder"
          >
            Cancelar Pedido
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
            <span class="value">{{ formatDate(order.deliveryDate) }}</span>
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
      <div v-if="order.notes" class="notes-card">
        <h3>Notas</h3>
        <p>{{ order.notes }}</p>
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
            v-if="canEdit" 
            class="btn btn-primary" 
            @click="showAddItemModal = true"
          >
            + Agregar Línea
          </button>
        </div>

        <div v-if="!order.lineItems || order.lineItems.length === 0" class="empty-state">
          <p>No hay líneas en este pedido</p>
        </div>

        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Variante</th>
                <th>Trabajo MES</th>
                <th>Cantidad</th>
                <th>Precio Unitario</th>
                <th>Descuento</th>
                <th>Subtotal</th>
                <th v-if="canEdit">Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in order.lineItems" :key="item.id">
                <td class="variant-id">{{ formatVariantId(item.productVariantID) }}</td>
                 <td>
                   <button
                     v-if="item.mesWorkId"
                     class="mes-link"
                     @click="goToMesWork(item.mesWorkId)"
                     type="button"
                     :title="getMesWorkTooltip(item.mesWorkId)"
                   >
                     {{ formatMesWorkId(item.mesWorkId) }}
                   </button>
                   <span v-else>—</span>
                 </td>
                <td>{{ item.quantity }}</td>
                <td>
                  {{ salesApi.formatMoney(item.finalUnitPrice) }}
                  <span v-if="item.manualUnitPrice" class="manual-badge">Manual</span>
                </td>
                <td>{{ item.finalDiscountPerUnit ? salesApi.formatMoney(item.finalDiscountPerUnit) : '—' }}</td>
                <td class="amount">{{ salesApi.formatMoney(item.subtotal) }}</td>
                <td v-if="canEdit" class="actions-cell">
                  <button 
                    class="btn-icon" 
                    @click="editLineItem(item)"
                    title="Editar"
                  >
                    ✏️
                  </button>
                  <button 
                    class="btn-icon danger" 
                    @click="removeLineItem(item.id)"
                    title="Eliminar"
                  >
                    🗑️
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
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
                    <span class="item-variant">{{ formatVariantId(item.productVariantId) }}</span>
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

    <!-- Add/Edit Line Item Modal -->
    <div v-if="showAddItemModal || showEditItemModal" class="modal-overlay" @click="closeModals">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ showEditItemModal ? 'Editar Línea' : 'Agregar Línea' }}</h3>
          <button class="btn-close" @click="closeModals">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>Variante de Producto *</label>
            <input
              v-model="lineItemForm.productVariantId"
              type="text"
              placeholder="UUID de la variante"
              class="form-input"
              :disabled="showEditItemModal"
            />
          </div>
          <div class="form-group">
            <label>Cantidad *</label>
            <input
              v-model.number="lineItemForm.quantity"
              type="number"
              min="1"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>Precio Unitario Manual (opcional)</label>
            <input
              v-model.number="lineItemForm.manualUnitPrice"
              type="number"
              step="0.01"
              min="0"
              placeholder="Dejar vacío para precio calculado"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>Descuento Por Unidad Manual (opcional)</label>
            <input
              v-model.number="lineItemForm.manualDiscountPerUnit"
              type="number"
              step="0.01"
              min="0"
              placeholder="Dejar vacío para descuento calculado"
              class="form-input"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeModals">Cancelar</button>
          <button 
            class="btn btn-primary" 
            @click="showEditItemModal ? updateLineItem() : addLineItem()"
            :disabled="!isLineItemFormValid"
          >
            {{ showEditItemModal ? 'Actualizar' : 'Agregar' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import salesApi from '@/services/salesApi';
import { mesApi } from '@/services/mesApi';
import { getPrintIssuerProfile } from '@/services/printIssuerProfile';

const route = useRoute();
const router = useRouter();

const order = ref(null);
const isLoading = ref(false);
const error = ref('');
const deliveryNotes = ref([]);
const mesWorksCache = ref({});
const issuerProfile = getPrintIssuerProfile();

const showAddItemModal = ref(false);
const showEditItemModal = ref(false);
const editingItemId = ref(null);

const showDeliveryNoteModal = ref(false);
const isCreatingDeliveryNote = ref(false);
const deliveryNoteForm = ref({
  deliveryDate: '',
  type: 'TOTAL',
  items: [],
  deliveryAddress: '',
  notes: '',
});

const lineItemForm = ref({
  productVariantId: '',
  quantity: 1,
  manualUnitPrice: null,
  manualDiscountPerUnit: null,
});

const canEdit = computed(() => {
  return order.value && ['PENDING', 'CONFIRMED'].includes(order.value.status);
});

const canCancel = computed(() => {
  return order.value && ['PENDING', 'CONFIRMED'].includes(order.value.status);
});

const canCreateDeliveryNote = computed(() => {
  return order.value && ['IN_PROGRESS', 'COMPLETED', 'CONFIRMED'].includes(order.value.status);
});

const minDeliveryDate = computed(() => {
  const today = new Date();
  return today.toISOString().split('T')[0];
});

const isLineItemFormValid = computed(() => {
  return lineItemForm.value.productVariantId && lineItemForm.value.quantity > 0;
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
    initializeDeliveryNoteForm();
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar el pedido';
    console.error('Error loading order:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadDeliveryNotes() {
  if (!order.value?.id) return;
  
  try {
    const response = await salesApi.listDeliveryNotes({ salesOrderId: order.value.id });
    deliveryNotes.value = Array.isArray(response) ? response : (response.data || []);
  } catch (err) {
    console.error('Error loading delivery notes:', err);
    // Non-critical, don't show error to user
  }
}

async function loadMesWorksForOrder() {
  const lineItems = order.value?.lineItems;
  if (!Array.isArray(lineItems) || lineItems.length === 0) return;

  const mesWorkIds = [...new Set(
    lineItems
      .map((item) => item?.mesWorkId)
      .filter((mesWorkId) => typeof mesWorkId === 'string' && mesWorkId.length > 0),
  )];

  const uncachedIds = mesWorkIds.filter((id) => !mesWorksCache.value[id]);
  if (uncachedIds.length === 0) return;

  const results = await Promise.allSettled(uncachedIds.map((id) => mesApi.getWork(id)));
  results.forEach((result, index) => {
    const mesWorkId = uncachedIds[index];
    mesWorksCache.value[mesWorkId] = result.status === 'fulfilled' ? result.value : null;
  });
}

function initializeDeliveryNoteForm() {
  if (!order.value?.lineItems) return;
  
  const today = new Date().toISOString().split('T')[0];
  deliveryNoteForm.value.deliveryDate = today;
  deliveryNoteForm.value.type = 'TOTAL';
  deliveryNoteForm.value.items = order.value.lineItems.map(item => ({
    lineItemId: item.id,
    productVariantId: item.productVariantID,
    availableQuantity: item.quantity,
    quantityToDeliver: item.quantity,
    selected: false,
  }));
}

async function createDeliveryNote() {
  if (!isDeliveryNoteFormValid.value || isCreatingDeliveryNote.value) return;
  
  isCreatingDeliveryNote.value = true;
  
  try {
    const deliveryNoteData = {
      salesOrderId: order.value.id,
      deliveryDate: salesApi.formatDateForAPI(new Date(deliveryNoteForm.value.deliveryDate)),
    };
    
    if (deliveryNoteForm.value.type === 'PARTIAL') {
      const selectedItems = deliveryNoteForm.value.items
        .filter(item => item.selected)
        .map(item => ({
          lineItemId: item.lineItemId,
          quantityDelivered: item.quantityToDeliver,
        }));
      
      deliveryNoteData.lineItems = selectedItems;
    }
    
    if (deliveryNoteForm.value.deliveryAddress) {
      deliveryNoteData.deliveryAddress = deliveryNoteForm.value.deliveryAddress;
    }
    
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

async function addLineItem() {
  try {
    const itemData = {
      productVariantId: lineItemForm.value.productVariantId,
      quantity: lineItemForm.value.quantity,
    };

    if (lineItemForm.value.manualUnitPrice) {
      itemData.manualUnitPrice = {
        amount: lineItemForm.value.manualUnitPrice,
        currency: 'EUR',
      };
    }

    if (lineItemForm.value.manualDiscountPerUnit) {
      itemData.manualDiscountPerUnit = {
        amount: lineItemForm.value.manualDiscountPerUnit,
        currency: 'EUR',
      };
    }

    await salesApi.addOrderLineItem(order.value.id, itemData);
    await fetchOrder();
    closeModals();
  } catch (err) {
    alert(err?.message || 'No se pudo agregar la línea');
  }
}

function editLineItem(item) {
  editingItemId.value = item.id;
  lineItemForm.value = {
    productVariantId: item.productVariantID,
    quantity: item.quantity,
    manualUnitPrice: item.manualUnitPrice?.amount || null,
    manualDiscountPerUnit: item.manualDiscountPerUnit?.amount || null,
  };
  showEditItemModal.value = true;
}

async function updateLineItem() {
  try {
    const updates = {
      quantity: lineItemForm.value.quantity,
    };

    if (lineItemForm.value.manualUnitPrice) {
      updates.manualUnitPrice = {
        amount: lineItemForm.value.manualUnitPrice,
        currency: 'EUR',
      };
    }

    if (lineItemForm.value.manualDiscountPerUnit) {
      updates.manualDiscountPerUnit = {
        amount: lineItemForm.value.manualDiscountPerUnit,
        currency: 'EUR',
      };
    }

    await salesApi.updateOrderLineItem(order.value.id, editingItemId.value, updates);
    await fetchOrder();
    closeModals();
  } catch (err) {
    alert(err?.message || 'No se pudo actualizar la línea');
  }
}

async function removeLineItem(itemId) {
  if (!confirm('¿Eliminar esta línea del pedido?')) return;

  try {
    await salesApi.removeOrderLineItem(order.value.id, itemId);
    await fetchOrder();
  } catch (err) {
    alert(err?.message || 'No se pudo eliminar la línea');
  }
}

function closeModals() {
  showAddItemModal.value = false;
  showEditItemModal.value = false;
  editingItemId.value = null;
  lineItemForm.value = {
    productVariantId: '',
    quantity: 1,
    manualUnitPrice: null,
    manualDiscountPerUnit: null,
  };
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

function formatMesWorkId(mesWorkId) {
  if (!mesWorkId) return '—';
  const mesWork = mesWorksCache.value[mesWorkId];
  if (mesWork?.work_number) return mesWork.work_number;
  return mesWorkId.substring(0, 8) + '...';
}

function getMesWorkTooltip(mesWorkId) {
  if (!mesWorkId) return 'Trabajo MES';
  const mesWork = mesWorksCache.value[mesWorkId];
  if (mesWork?.work_number && mesWork?.work_name) {
    return `${mesWork.work_number} · ${mesWork.work_name}`;
  }
  if (mesWork?.work_number) return mesWork.work_number;
  return `Trabajo MES ${mesWorkId}`;
}

function goToMesWork(mesWorkId) {
  if (!mesWorkId) return;
  router.push(`/mes/works/${mesWorkId}`);
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
  font-family: 'Courier New', monospace;
  color: #6b7280;
}

.amount {
  font-weight: 600;
  text-align: right;
}

.manual-badge {
  display: inline-block;
  margin-left: 0.5rem;
  padding: 0.125rem 0.375rem;
  background: #e0e7ff;
  color: #3730a3;
  border-radius: 4px;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
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

.actions-cell {
  display: flex;
  gap: 0.5rem;
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

.item-variant {
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
  color: #1f2937;
  font-weight: 500;
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
  max-width: 700px;
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

.print-doc-header,
.print-doc-footer {
  display: none;
}

.print-brand {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #374151;
}

.print-doc-header h2 {
  margin: 0.35rem 0;
  font-size: 1.3rem;
  font-weight: 700;
  color: #111827;
}

.print-issuer-line {
  margin: 0;
  font-size: 0.75rem;
  color: #4b5563;
}

.print-doc-meta {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  font-size: 0.8rem;
  color: #4b5563;
}

@media print {
  :deep(.navbar),
  :deep(nav),
  .page-header,
  .btn-back,
  .header-actions,
  .modal-overlay,
  .btn,
  .btn-icon,
  .actions-cell {
    display: none !important;
  }

  .order-detail-container {
    padding: 0;
    max-width: none;
  }

  .print-doc-header,
  .print-doc-footer {
    display: block;
    border: 1px solid #d1d5db;
    padding: 0.75rem 1rem;
    margin-bottom: 0.75rem;
    background: white;
  }

  .print-doc-footer {
    margin-top: 0.75rem;
    margin-bottom: 0;
    display: flex;
    justify-content: space-between;
    font-size: 0.75rem;
    color: #4b5563;
  }

  .info-card,
  .notes-card,
  .line-items-section,
  .table-container,
  .delivery-notes-section {
    box-shadow: none !important;
    border: 1px solid #d1d5db;
  }

  .info-grid {
    grid-template-columns: 1fr 1fr;
    gap: 0.75rem;
  }

  .status-badge {
    border: 1px solid #9ca3af;
    color: #111827 !important;
    background: transparent !important;
  }

  .data-table {
    font-size: 0.75rem;
  }
}
</style>
