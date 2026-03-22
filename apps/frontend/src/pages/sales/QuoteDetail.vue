<template>
  <Navbar />
  <div class="quote-detail-container">
    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando presupuesto...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchQuote">Reintentar</button>
      <button class="btn btn-secondary" @click="goBack">Volver</button>
    </div>

    <!-- Quote Detail -->
    <div v-else-if="quote" ref="quoteContentRef" class="quote-content">
      <!-- Header -->
      <div class="page-header">
        <div>
          <button class="btn-back" @click="goBack">← Volver</button>
          <h1>Presupuesto {{ quote.quoteNumber }}</h1>
          <span :class="['status-badge', `status-${getStatusClass(quote.status)}`]">
            {{ getStatusLabel(quote.status) }}
          </span>
        </div>
        <div class="header-actions">
          <button
            v-if="quote.status !== 'BORRADOR'"
            class="btn btn-secondary"
            @click="printQuote"
            title="Imprimir presupuesto"
          >
            🖨️ Imprimir
          </button>
          <!-- DRAFT/ISSUED: Editar, Emitir, Eliminar -->
          <button 
            v-if="canEdit && !isEditing" 
            class="btn btn-primary" 
            @click="enterEditMode"
            title="Editar presupuesto"
          >
            ✏️ Editar
          </button>
          <button 
            v-if="isEditing" 
            class="btn btn-success" 
            @click="saveQuote"
            :disabled="isSaving"
            title="Guardar cambios"
          >
            {{ isSaving ? 'Guardando...' : '💾 Guardar' }}
          </button>
          <button 
            v-if="isEditing" 
            class="btn btn-secondary" 
            @click="cancelEdit"
            title="Cancelar edición"
          >
            ✕ Cancelar
          </button>
          <button 
            v-if="quote.status === 'BORRADOR' && !isEditing" 
            class="btn btn-primary" 
            @click="confirmIssueQuote"
            title="Emitir presupuesto al cliente"
          >
            📧 Emitir
          </button>

          <button 
            v-if="quote.status === 'BORRADOR' && !isEditing" 
            class="btn btn-danger" 
            @click="deleteQuote"
            title="Eliminar presupuesto"
          >
            🗑️ Eliminar
          </button>
          
          <!-- ISSUED: Editar, Aceptar (abre modal de conversión), Rechazar -->
          <button 
            v-if="quote.status === 'EMITIDA' && !isExpired && !isEditing" 
            class="btn btn-success" 
            @click="showConvertModal = true"
            title="Aceptar presupuesto y generar pedido"
          >
            ✓ Aceptar y Generar Pedido
          </button>
          <button 
            v-if="quote.status === 'EMITIDA' && !isEditing" 
            class="btn btn-danger" 
            @click="rejectQuote"
            title="Marcar como rechazado"
          >
            ✕ Rechazar
          </button>

          <!-- REJECTED: Reactivar -->
          <button 
            v-if="quote.status === 'RECHAZADA'" 
            class="btn btn-primary" 
            @click="reactivateQuote"
            title="Volver a borrador para editar y reemitir"
          >
            🔄 Reactivar
          </button>
        </div>
      </div>

      <div class="print-doc-header">
        <p class="print-brand">{{ issuerProfile.displayName }}</p>
        <p v-if="issuerProfile.taxId" class="print-issuer-line">{{ issuerProfile.taxLabel }}: {{ issuerProfile.taxId }}</p>
        <p v-if="issuerProfile.addressLine || issuerProfile.cityLine" class="print-issuer-line">{{ issuerProfile.addressLine }}<span v-if="issuerProfile.addressLine && issuerProfile.cityLine"> · </span>{{ issuerProfile.cityLine }}</p>
        <p v-if="issuerProfile.contactLine" class="print-issuer-line">{{ issuerProfile.contactLine }}</p>
        <h2>Presupuesto {{ quote.quoteNumber }}</h2>
        <div class="print-doc-meta">
          <span>Cliente: {{ partyName }}</span>
          <span>Fecha: {{ formatDate(quote.quoteDate) }}</span>
          <span>Estado: {{ getStatusLabel(quote.status) }}</span>
        </div>
      </div>

      <!-- Expiration Warning (ISSUED and approaching expiration) -->
      <div v-if="quote.status === 'EMITIDA' && daysUntilExpiration <= 7 && daysUntilExpiration > 0" class="warning-card">
        ⚠️ Este presupuesto vence en <strong>{{ daysUntilExpiration }}</strong> día{{ daysUntilExpiration !== 1 ? 's' : '' }}.
      </div>
      <div v-if="isExpired && quote.status === 'EMITIDA'" class="error-card">
        ❌ Este presupuesto ha <strong>EXPIRADO</strong>. No se puede convertir a pedido.
      </div>

      <!-- Generated Order Link -->
      <div v-if="quote.generatedOrderId" class="success-card">
        📦 Este presupuesto ha generado el pedido
        <router-link :to="`/sales/orders/${quote.generatedOrderId}`" class="order-link">
          {{ quote.generatedOrderNumber }}
        </router-link>
      </div>

      <!-- Quote Info Cards -->
      <div class="info-grid">
        <div class="info-card">
          <h3>Información del Cliente</h3>
          <div class="info-row">
            <span class="label">Cliente:</span>
            <span class="value">{{ partyName }}</span>
          </div>
          <div class="info-row">
            <span class="label">Party ID:</span>
            <span class="value party-id">{{ formatPartyId(quote.partyId) }}</span>
          </div>
        </div>

        <div class="info-card">
          <h3>Fechas</h3>
          <div class="info-row">
            <span class="label">Fecha de Presupuesto:</span>
            <span class="value">{{ formatDate(quote.quoteDate) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Válido Hasta:</span>
            <span v-if="!isEditing" class="value" :class="{'text-danger': isExpired, 'text-warning': daysUntilExpiration <= 7 && daysUntilExpiration > 0}">
              {{ formatDate(quote.expirationDate) }}
            </span>
            <input v-else v-model="editForm.validUntil" type="date" class="form-input form-input-inline" />
          </div>
        </div>

        <div class="info-card">
          <h3>Totales</h3>
          <div class="info-row">
            <span class="label">Subtotal:</span>
            <span class="value">{{ salesApi.formatMoney(quote.subtotal) }}</span>
          </div>
          <div class="info-row">
            <span class="label">IVA:</span>
            <span class="value">{{ salesApi.formatMoney(quote.taxAmount) }}</span>
          </div>
          <div class="info-row total">
            <span class="label">Total:</span>
            <span class="value">{{ salesApi.formatMoney(quote.total) }}</span>
          </div>
        </div>
      </div>

      <!-- Observaciones -->
      <div v-if="quote.notes || isEditing" class="notes-card">
        <h3>Observaciones</h3>
        <p v-if="!isEditing">{{ quote.notes }}</p>
        <textarea v-else v-model="editForm.notes" class="form-textarea" rows="3" placeholder="Observaciones sobre el presupuesto..."></textarea>
      </div>

      <!-- MES Work References (Document-level) -->
      <div v-if="(quote.mesWorkRefs && quote.mesWorkRefs.length > 0) || isEditing" class="notes-card">
        <h3>Configuraciones MES</h3>
        <template v-if="!isEditing">
          <div v-for="mesRef in quote.mesWorkRefs" :key="mesRef.id" class="mes-ref-view">
            <div class="mes-ref-header">
              <span v-if="mesRef.workOrderId" class="mes-status-badge status-en_proceso">Con orden</span>
              <span v-else-if="mesRef.workSetupId" class="mes-status-badge status-pendiente">Configurado</span>
              <span v-else class="mes-status-badge status-borrador">Sin configurar</span>
            </div>
            <p v-if="mesRef.description" class="mes-description">{{ mesRef.description }}</p>
          </div>
        </template>
        <template v-else>
          <div v-if="isLoadingMesWorks" class="mes-loading">Cargando configuraciones MES...</div>
          <div v-else class="mes-ref-list">
            <div v-for="(config, idx) in editForm.mesWorkRefs" :key="idx" class="mes-config-entry">
              <div class="form-row mes-config-row">
                <div class="form-group">
                  <label>Configuración base</label>
                  <select class="form-input" :value="config.workSetupId || ''"
                    @change="onSetupSelect(idx, $event.target.value)">
                    <option value="">— Personalizada —</option>
                    <option v-for="ws in mesWorkSetups" :key="ws.id" :value="ws.id">{{ ws.name }}</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Descripción</label>
                  <div class="input-with-action">
                    <textarea class="form-textarea" rows="2" v-model="config.description"
                      placeholder="Descripción de la configuración..."></textarea>
                    <button type="button" class="btn-remove" @click="removeEditConfig(idx)" title="Eliminar configuración">✕</button>
                  </div>
                </div>
              </div>
            </div>
            <button type="button" class="btn btn-secondary" @click="addEditConfig">+ Agregar configuración</button>
          </div>
        </template>
      </div>

      <!-- Line Items -->
      <div class="line-items-section">
        <div class="section-header">
          <h2>Líneas del Presupuesto</h2>
          <button v-if="isEditing" type="button" class="btn btn-primary" @click="addEditLineItem">
            + Agregar Línea
          </button>
        </div>

        <!-- View mode -->
        <template v-if="!isEditing">
        <div v-if="!quote.lineItems || quote.lineItems.length === 0" class="empty-state">
          <p>No hay líneas en este presupuesto</p>
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
              <tr v-for="item in quote.lineItems" :key="item.id">
                <td class="variant-id">
                  <span v-if="item.variantSku">{{ item.variantSku }}</span>
                  <span v-else>{{ formatVariantId(item.productVariantID) }}</span>
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
            <tfoot>
              <tr class="totals-row">
                <td colspan="6" class="totals-label">Subtotal:</td>
                <td class="amount">{{ salesApi.formatMoney(quote.subtotal) }}</td>
              </tr>
              <tr class="totals-row">
                <td colspan="6" class="totals-label">IVA:</td>
                <td class="amount">{{ salesApi.formatMoney(quote.taxAmount) }}</td>
              </tr>
              <tr class="totals-row total">
                <td colspan="6" class="totals-label">Total:</td>
                <td class="amount">{{ salesApi.formatMoney(quote.total) }}</td>
              </tr>
            </tfoot>
          </table>
        </div>
        </template>

        <!-- Edit mode -->
        <template v-else>
        <div v-if="editForm.lineItems.length === 0" class="empty-state">
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
              <tr v-for="(item, idx) in editForm.lineItems" :key="idx">
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
                  <span v-else class="subtotal-value">{{ formatMoneyAmount(getEditLineSubtotal(idx)) }}</span>
                </td>
                <td class="actions-cell">
                  <button type="button" class="btn-icon danger" @click="removeEditLineItem(idx)" title="Eliminar">🗑️</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Totals Summary -->
        <div v-if="editForm.lineItems.length > 0" class="totals-section">
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
        <span>Válido hasta: {{ formatDate(quote.expirationDate) }}</span>
      </div>
    </div>

    <!-- Post-Issue Actions Modal -->
    <div v-if="showPostIssueModal" class="modal-overlay" @click="showPostIssueModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Presupuesto Emitido</h3>
          <button class="btn-close" @click="showPostIssueModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p class="modal-description">
            El presupuesto <strong>{{ quote?.quoteNumber }}</strong> se ha emitido correctamente.
          </p>
          <p class="post-issue-prompt">¿Qué desea hacer a continuación?</p>
          <div class="post-issue-actions">
            <button class="btn btn-primary post-issue-btn" @click="postIssuePrint">
              🖨️ Imprimir
            </button>
            <button class="btn btn-primary post-issue-btn" @click="postIssueEmail">
              📧 Descargar PDF y enviar por correo
            </button>
            <button class="btn btn-primary post-issue-btn" @click="postIssueBoth">
              🖨️📧 Descargar PDF, imprimir y enviar
            </button>
          </div>
          <p class="post-issue-hint">El PDF se descargará automáticamente para adjuntarlo al correo.</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showPostIssueModal = false">Cerrar sin acción</button>
        </div>
      </div>
    </div>

    <!-- Convert to Order Modal -->
    <div v-if="showConvertModal" class="modal-overlay" @click="showConvertModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Aceptar Presupuesto y Generar Pedido</h3>
          <button class="btn-close" @click="showConvertModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p class="modal-description">
            Presupuesto: <strong>{{ quote?.quoteNumber }}</strong><br>
            Total: <strong>{{ quote ? salesApi.formatMoney(quote.total) : '' }}</strong>
          </p>
          <div class="form-group">
            <label for="deliveryDate">Fecha de Entrega *</label>
            <input
              id="deliveryDate"
              v-model="convertForm.deliveryDate"
              type="date"
              class="form-input"
              :min="minDeliveryDate"
              required
            />
            <small class="help-text">La fecha de entrega debe ser al menos mañana</small>
          </div>
          <div class="form-group">
            <label for="orderNotes">Notas del Pedido</label>
            <textarea
              id="orderNotes"
              v-model="convertForm.notes"
              class="form-textarea"
              rows="3"
              placeholder="Notas adicionales para el pedido..."
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showConvertModal = false">Cancelar</button>
          <button 
            class="btn btn-primary" 
            @click="convertToOrder"
            :disabled="!convertForm.deliveryDate || isConverting"
          >
            {{ isConverting ? 'Convirtiendo...' : 'Crear Pedido' }}
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

const quoteContentRef = ref(null);
const quote = ref(null);
const isLoading = ref(false);
const error = ref('');
const partyName = ref('Cargando...');
const partyDefaultDiscount = ref(null);
const mesWorkSetups = ref([]);
const isLoadingMesWorks = ref(false);
const issuerProfile = getPrintIssuerProfile();

const showPostIssueModal = ref(false);
const showConvertModal = ref(false);
const isConverting = ref(false);
const convertForm = ref({
  deliveryDate: '',
  notes: '',
});

// Edit mode state
const isEditing = ref(false);
const isSaving = ref(false);
const showVariantSelector = ref(false);
const editForm = ref({
  validUntil: '',
  notes: '',
  mesWorkRefs: [],
  lineItems: [],
});

const minDeliveryDate = computed(() => {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().split('T')[0];
});

const canEdit = computed(() => {
  const s = quote.value?.status;
  return s === 'BORRADOR' || s === 'EMITIDA';
});

const isExpired = computed(() => {
  if (!quote.value?.expirationDate) return false;
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const validUntil = new Date(quote.value.expirationDate);
  validUntil.setHours(0, 0, 0, 0);
  return validUntil < today;
});

const daysUntilExpiration = computed(() => {
  if (!quote.value?.expirationDate) return null;
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const validUntil = new Date(quote.value.expirationDate);
  validUntil.setHours(0, 0, 0, 0);
  const diffTime = validUntil - today;
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return diffDays;
});

onMounted(() => {
  fetchQuote();
});

async function fetchQuote() {
  const quoteId = route.params.id;
  if (!quoteId) {
    error.value = 'ID de presupuesto no válido';
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    quote.value = await salesApi.getQuote(quoteId);
    await loadPartyName();
    // MES work refs are loaded inline from the order response
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar el presupuesto';
    console.error('Error loading quote:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyName() {
  if (!quote.value?.partyId) {
    partyName.value = 'Desconocido';
    return;
  }
  
  try {
    const party = await partyApi.getParty(quote.value.partyId);
    partyName.value = party.name || 'Sin nombre';
    partyDefaultDiscount.value = party.default_discount_percentage || null;
  } catch (err) {
    console.error('Error loading party name:', err);
    partyName.value = 'Error al cargar';
  }
}

async function confirmIssueQuote() {
  if (!confirm('¿Emitir este presupuesto al cliente? El estado cambiará a "Emitido".')) return;

  try {
    await salesApi.changeQuoteStatus(quote.value.id, 'EMITIDA');
    await fetchQuote();
    showPostIssueModal.value = true;
  } catch (err) {
    alert(err?.message || 'No se pudo emitir el presupuesto');
  }
}

function openMailClient() {
  let email = '';
  try {
    // Best-effort: intentar obtener email del contacto principal
    partyApi.listContacts(quote.value.partyId).then(({ data: contacts }) => {
      const primary = contacts.find(c => c.email);
      if (primary) email = primary.email;
    }).catch(() => {});
  } catch { /* no bloquear */ }

  const q = quote.value;
  const subject = encodeURIComponent(`Presupuesto ${q.quoteNumber || q.id}`);
  const total = q.total?.amount != null ? Number(q.total.amount).toFixed(2) + ' €' : '(pendiente)';
  const body = encodeURIComponent(
    `Estimado/a ${partyName.value},\n\n` +
    `Adjunto le enviamos el presupuesto ${q.quoteNumber || ''} ` +
    `por un total de ${total}.\n\n` +
    `Quedamos a su disposición para cualquier consulta.\n\n` +
    `Un saludo.`
  );
  window.open(`mailto:${email}?subject=${subject}&body=${body}`, '_self');
}

function postIssuePrint() {
  showPostIssueModal.value = false;
  window.print();
}

async function postIssueEmail() {
  showPostIssueModal.value = false;
  await generateQuotePdf();
  openMailClient();
}

async function postIssueBoth() {
  showPostIssueModal.value = false;
  await generateQuotePdf();
  window.print();
  setTimeout(() => openMailClient(), 500);
}

async function acceptQuote() {
  // Now handled by the accept-and-convert modal flow
  showConvertModal.value = true;
}

async function rejectQuote() {
  if (!confirm('¿Marcar este presupuesto como rechazado? Podrá reactivarse más tarde.')) return;

  try {
    await salesApi.changeQuoteStatus(quote.value.id, 'RECHAZADA');
    await fetchQuote();
  } catch (err) {
    alert(err?.message || 'No se pudo rechazar el presupuesto');
  }
}

async function reactivateQuote() {
  if (!confirm('¿Reactivar este presupuesto? Volverá al estado Borrador para poder editarlo y reemitirlo.')) return;

  try {
    await salesApi.changeQuoteStatus(quote.value.id, 'BORRADOR');
    await fetchQuote();
  } catch (err) {
    alert(err?.message || 'No se pudo reactivar el presupuesto');
  }
}

async function deleteQuote() {
  if (!confirm('¿Eliminar este presupuesto? Esta acción no se puede deshacer.')) return;

  try {
    await salesApi.deleteQuote(quote.value.id);
    router.push('/sales/quotes');
  } catch (err) {
    alert(err?.message || 'No se pudo eliminar el presupuesto');
  }
}

async function convertToOrder() {
  if (!convertForm.value.deliveryDate) {
    alert('Por favor, especifica una fecha de entrega');
    return;
  }

  if (isExpired.value) {
    alert('No se puede convertir un presupuesto expirado a pedido');
    return;
  }

  isConverting.value = true;

  try {
    let order;
    if (quote.value.status === 'EMITIDA') {
      // Accept + convert in one atomic operation
      order = await salesApi.acceptAndConvertQuote(quote.value.id, salesApi.formatDateForAPI(new Date(convertForm.value.deliveryDate)));
    } else {
      // Already APROBADA, just convert
      order = await salesApi.convertQuoteToOrder(quote.value.id, salesApi.formatDateForAPI(new Date(convertForm.value.deliveryDate)));
    }
    
    alert(`✓ Presupuesto convertido exitosamente.\nPedido creado: ${order.orderNumber || order.id}`);
    
    // Navigate to the new order
    if (order.id) {
      router.push(`/sales/orders/${order.id}`);
    } else {
      router.push('/sales/orders');
    }
  } catch (err) {
    alert(err?.message || 'No se pudo convertir el presupuesto a pedido');
    isConverting.value = false;
  }
}

function enterEditMode() {
  editForm.value.validUntil = quote.value.expirationDate ? new Date(quote.value.expirationDate).toISOString().split('T')[0] : '';
  editForm.value.notes = quote.value.notes || '';
  editForm.value.mesWorkRefs = (quote.value.mesWorkRefs || []).map(r => ({
    workSetupId: r.workSetupId || null,
    description: r.description || '',
  }));
  editForm.value.lineItems = (quote.value.lineItems || []).map(item => ({
    productVariantId: item.productVariantId,
    variantSku: item.variantSku || '',
    displayName: buildDisplayName(item),
    quantity: item.quantity,
    listPrice: item.listUnitPrice?.amount ?? null,
    unitPrice: item.unitPrice?.amount ?? null,
    discountPercent: item.discountPercent ?? null,
    taxRate: item.taxRate ?? 21,
  }));
  isEditing.value = true;
  loadMesWorksForParty();
  fetchPreviewCalculation();
}

function cancelEdit() {
  isEditing.value = false;
  showVariantSelector.value = false;
}

async function loadMesWorksForParty() {
  const partyId = quote.value?.partyId;
  if (!partyId) return;
  isLoadingMesWorks.value = true;
  try {
    mesWorkSetups.value = await mesApi.listWorkSetups({ party_id: partyId });
  } catch (err) {
    console.error('Error loading MES work setups for party:', err);
    mesWorkSetups.value = [];
  } finally {
    isLoadingMesWorks.value = false;
  }
}

function onSetupSelect(idx, setupId) {
  const config = editForm.value.mesWorkRefs[idx];
  if (!config) return;
  if (setupId) {
    config.workSetupId = setupId;
  } else {
    config.workSetupId = null;
  }
}

function addEditConfig() {
  editForm.value.mesWorkRefs.push({ workSetupId: null, description: '' });
}

function removeEditConfig(idx) {
  editForm.value.mesWorkRefs.splice(idx, 1);
}

function addEditLineItem() {
  showVariantSelector.value = true;
}

function handleVariantSelected(payload) {
  const variant = payload?.variant;
  if (variant) {
    const name = variant.product_name || '';
    const config = variant.option_configuration;
    let displayName = name || '—';
    if (config && Object.keys(config).length > 0) {
      displayName = name + ' - ' + Object.values(config).join(', ');
    }
    const newItem = {
      productVariantId: variant.id,
      variantSku: variant.sku || '',
      displayName,
      quantity: 1,
      listPrice: null,
      unitPrice: null,
      discountPercent: partyDefaultDiscount.value || null,
      taxRate: 21,
    };
    editForm.value.lineItems.push(newItem);
    showVariantSelector.value = false;
    // Use the reactive proxy from the array, not the plain object, so Vue detects async mutations
    const reactiveItem = editForm.value.lineItems[editForm.value.lineItems.length - 1];
    fetchEditLinePrice(reactiveItem, variant.product_id || '', variant.product_base_price);
  }
}

async function fetchEditLinePrice(item, productId, basePrice) {
  if (!item.productVariantId || !productId) {
    if (basePrice != null) {
      const price = Math.round(basePrice * 1000) / 1000;
      item.listPrice = price;
      item.unitPrice = price;
    }
    return;
  }
  try {
    const result = await calculateBaseSalesPrice(productId, item.productVariantId);
    const rawBaseCost = result.baseCost?.amount ?? basePrice ?? null;
    const rawSalesPrice = result.baseSalesPrice?.amount ?? basePrice ?? null;
    item.listPrice = rawBaseCost != null ? Math.round(rawBaseCost * 1000) / 1000 : null;
    item.unitPrice = rawSalesPrice != null ? Math.round(rawSalesPrice * 1000) / 1000 : null;
    if (result.taxRate != null) {
      item.taxRate = result.taxRate;
    }
  } catch (err) {
    console.warn('[QuoteDetail] Error fetching sale price:', err.message);
    if (basePrice != null) {
      const price = Math.round(basePrice * 1000) / 1000;
      item.listPrice = price;
      item.unitPrice = price;
    }
  }
}

function removeEditLineItem(index) {
  editForm.value.lineItems.splice(index, 1);
}

// Backend-driven preview: all monetary calculations come from the API.
const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewDebounceTimer = null;

function getEditLineSubtotal(index) {
  if (!previewResult.value) return 0;
  // Map form index to preview index (preview only contains items with productVariantId)
  let previewIdx = 0;
  for (let i = 0; i < index; i++) {
    if (editForm.value.lineItems[i].productVariantId) previewIdx++;
  }
  const previewItem = previewResult.value.lineItems[previewIdx];
  return previewItem?.subtotal?.amount ?? 0;
}

const editCalculatedTotals = computed(() => {
  if (!previewResult.value) return { subtotal: 0, tax: 0, total: 0 };
  return {
    subtotal: previewResult.value.subtotal?.amount ?? 0,
    tax: previewResult.value.taxAmount?.amount ?? 0,
    total: previewResult.value.total?.amount ?? 0,
  };
});

function buildPreviewItems() {
  return editForm.value.lineItems
    .filter(item => item.productVariantId)
    .map(item => {
      const entry = {
        productVariantId: item.productVariantId,
        quantity: item.quantity || 1,
      };
      if (item.unitPrice != null && item.unitPrice > 0) {
        entry.unitPrice = { amount: item.unitPrice, currency: 'EUR' };
      }
      if (item.discountPercent != null) {
        entry.discountPercent = item.discountPercent;
      }
      return entry;
    });
}

async function fetchPreviewCalculation() {
  const partyId = quote.value?.partyId;
  const items = buildPreviewItems();
  if (!partyId || items.length === 0) {
    previewResult.value = null;
    return;
  }
  isPreviewLoading.value = true;
  try {
    previewResult.value = await salesApi.previewQuoteCalculation(partyId, items);
  } catch (err) {
    console.warn('[QuoteDetail] Preview calculation error:', err.message);
  } finally {
    isPreviewLoading.value = false;
  }
}

function debouncedPreview() {
  clearTimeout(previewDebounceTimer);
  previewDebounceTimer = setTimeout(fetchPreviewCalculation, 400);
}

watch(
  () => editForm.value.lineItems.map(i => `${i.productVariantId}|${i.quantity}|${i.unitPrice}|${i.discountPercent}`).join(','),
  () => {
    if (isEditing.value) debouncedPreview();
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

async function saveQuote() {
  if (editForm.value.lineItems.length === 0) {
    alert('Debe haber al menos una línea en el presupuesto');
    return;
  }

  isSaving.value = true;

  try {
    const items = editForm.value.lineItems.map(item => {
      const lineItem = {
        productVariantId: item.productVariantId,
        quantity: item.quantity,
      };

      if (item.unitPrice != null && item.unitPrice > 0) {
        lineItem.unitPrice = {
          amount: item.unitPrice,
          currency: 'EUR',
        };
      }

      if (item.discountPercent != null) {
        lineItem.discountPercent = item.discountPercent;
      }

      return lineItem;
    });

    const updateData = {
      items,
      notes: editForm.value.notes || undefined,
      mesWorkRefs: editForm.value.mesWorkRefs,
    };

    if (editForm.value.validUntil) {
      updateData.expirationDate = salesApi.formatDateForAPI(new Date(editForm.value.validUntil));
    }

    await salesApi.updateQuote(quote.value.id, updateData);
    isEditing.value = false;
    await fetchQuote();
  } catch (err) {
    alert(err?.message || 'No se pudo actualizar el presupuesto');
  } finally {
    isSaving.value = false;
  }
}

function goBack() {
  router.push('/sales/quotes');
}

function printQuote() {
  window.print();
}

async function generateQuotePdf() {
  const el = quoteContentRef.value;
  if (!el) return;
  const { default: html2pdf } = await import('html2pdf.js');
  const filename = `Presupuesto_${quote.value.quoteNumber || quote.value.id}.pdf`;

  // Toggle print-like styling so html2pdf captures the professional layout
  el.classList.add('pdf-rendering');
  // Allow the browser to repaint before capture
  await new Promise(r => setTimeout(r, 100));

  try {
    await html2pdf()
      .set({
        margin: [10, 10, 10, 10],
        filename,
        image: { type: 'jpeg', quality: 0.95 },
        html2canvas: { scale: 2, useCORS: true },
        jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' },
        pagebreak: { mode: ['avoid-all', 'css', 'legacy'] },
      })
      .from(el)
      .save();
  } finally {
    el.classList.remove('pdf-rendering');
  }
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

function getStatusLabel(status) {
  const labels = {
    BORRADOR: 'Borrador',
    EMITIDA: 'Emitida',
    APROBADA: 'Aprobada',
    RECHAZADA: 'Rechazada',
    EXPIRADA: 'Expirada',
    CONVERTIDA_A_PEDIDO: 'Convertida a Pedido',
  };
  return labels[status] || status;
}

function getStatusClass(status) {
  const classes = {
    BORRADOR: 'warning',
    EMITIDA: 'info',
    APROBADA: 'success',
    RECHAZADA: 'danger',
    EXPIRADA: 'secondary',
    CONVERTIDA_A_PEDIDO: 'primary',
  };
  return classes[status] || 'secondary';
}
</script>

<style scoped>
.quote-detail-container {
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
  flex-wrap: wrap;
}

.warning-card,
.error-card,
.success-card {
  padding: 1rem 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
}

.warning-card {
  background: #fef3c7;
  border-left: 4px solid #f59e0b;
  color: #92400e;
}

.error-card {
  background: #fee2e2;
  border-left: 4px solid #ef4444;
  color: #991b1b;
}

.success-card {
  background: #d1fae5;
  border-left: 4px solid #10b981;
  color: #065f46;
}

.success-card .order-link {
  font-weight: 600;
  color: #047857;
  text-decoration: underline;
  margin-left: 0.25rem;
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

.party-id {
  font-family: 'Courier New', monospace;
  font-size: 0.8rem;
}

.text-danger {
  color: #dc2626 !important;
  font-weight: 600;
}

.text-warning {
  color: #f59e0b !important;
  font-weight: 600;
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

.data-table tfoot {
  background: #f9fafb;
  font-weight: 600;
}

.totals-row td {
  padding: 0.75rem 1rem;
  border-top: 1px solid #e5e7eb;
}

.totals-row.total td {
  border-top: 2px solid #d1d5db;
  font-size: 1rem;
  padding-top: 1rem;
}

.totals-label {
  text-align: right;
  color: #6b7280;
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

.mes-name {
  font-weight: 600;
  font-size: 0.875rem;
}

.mes-status-badge {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.mes-status-badge.status-borrador { background: #f3f4f6; color: #6b7280; }
.mes-status-badge.status-pendiente { background: #fef3c7; color: #92400e; }
.mes-status-badge.status-en_proceso { background: #dbeafe; color: #1e40af; }
.mes-status-badge.status-completado { background: #d1fae5; color: #065f46; }
.mes-status-badge.status-cancelado { background: #fee2e2; color: #991b1b; }

.mes-ref-view {
  padding: 0.5rem 0;
  border-bottom: 1px solid #f3f4f6;
}

.mes-ref-view:last-child {
  border-bottom: none;
}

.mes-ref-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.mes-description {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0.25rem 0 0;
  line-height: 1.5;
}

.mes-ref-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.mes-config-entry {
  padding-bottom: 1rem;
  border-bottom: 1px solid #f3f4f6;
}

.mes-config-entry:last-of-type {
  border-bottom: none;
  padding-bottom: 0;
}

.mes-config-row {
  margin-bottom: 0;
}

.input-with-action {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.input-with-action .form-input {
  flex: 1;
}

.mes-loading,
.mes-empty {
  color: #6b7280;
  font-size: 0.875rem;
  font-style: italic;
  padding: 0.5rem 0;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #9ca3af;
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

.modal-description {
  font-size: 0.875rem;
  color: #4b5563;
  line-height: 1.6;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: #f9fafb;
  border-radius: 4px;
}

.post-issue-prompt {
  font-size: 0.95rem;
  color: #374151;
  margin-bottom: 1rem;
}

.post-issue-actions {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.post-issue-btn {
  width: 100%;
  justify-content: center;
  font-size: 0.95rem;
  padding: 0.75rem 1rem;
}

.post-issue-hint {
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: var(--tt-text-secondary, #64748b);
  text-align: center;
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

.form-input,
.form-textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: inherit;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.form-textarea {
  resize: vertical;
}

.help-text {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: #6b7280;
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

.form-input-small {
  width: 100px;
  padding: 0.25rem 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
}

.form-input-small:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.actions-cell {
  white-space: nowrap;
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
  cursor: pointer;
  padding: 0.25rem;
  font-size: 1rem;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn-icon:hover {
  background: #f3f4f6;
}

.btn-icon.danger:hover {
  background: #fee2e2;
}

.modal-wide {
  max-width: 700px;
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
