<template>
  <BaseEntityPage class="no-print" v-if="isLoading">
    <template #header>
      <PageHeader title="Cargando..." :breadcrumbs="[{ label: 'Ventas', to: '/sales/delivery-notes' }, { label: 'Albaranes' }]" />
    </template>
    <div class="loading-state card">
      <div class="spinner"></div>
      <p>Cargando información del albarán...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage class="no-print" v-else-if="error">
    <template #header>
      <PageHeader title="Error" :breadcrumbs="[{ label: 'Ventas', to: '/sales/delivery-notes' }, { label: 'Albaranes' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-icon-wrapper error">
        <AlertCircle :size="32" />
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/sales/delivery-notes')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage class="no-print" v-else-if="deliveryNote">
    <!-- 1. IDENTITY HEADER -->
    <template #header>
      <PageHeader 
        :title="mode === 'edit' ? `Editando Albarán ${deliveryNote.deliveryNoteNumber}` : `Albarán ${deliveryNote.deliveryNoteNumber}`" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/dashboard' }, { label: 'Albaranes', to: '/sales/delivery-notes' }, { label: deliveryNote.deliveryNoteNumber }]"
      >
        <template #icon>
          <Truck :size="28" />
        </template>
        <template #actions>
          <template v-if="mode === 'detail'">
            <button class="btn btn-outline btn-sm" @click="printDeliveryNote">
              <Printer :size="18" /> <span>Imprimir</span>
            </button>
            <button v-if="['PENDING', 'PENDIENTE'].includes(deliveryNote.status)" class="btn btn-primary btn-sm" @click="enterEditMode">
              <Pencil :size="18" /> <span>Editar</span>
            </button>
          </template>
          <template v-else>
            <button class="btn btn-outline btn-sm" @click="mode = 'detail'" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-secondary btn-sm" @click="saveDeliveryNote" :disabled="isSaving">
              <component :is="isSaving ? RefreshCw : Save" :size="18" :class="{ 'spin': isSaving }" />
              <span>{{ isSaving ? 'Guardando...' : 'Guardar' }}</span>
            </button>
          </template>
        </template>
      </PageHeader>
    </template>

    <!-- 2. TOOLBAR -->
    <template #toolbar v-if="mode === 'detail'">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', `status-${salesApi.getStatusClass(deliveryNote.status)}`]">
            {{ salesApi.getStatusLabel(deliveryNote.status) }}
          </span>
        </div>
        <div class="toolbar-buttons">
          <button
            v-if="['PENDING', 'PENDIENTE'].includes(deliveryNote.status)"
            class="btn btn-success btn-sm"
            @click="markAsDelivered"
            :disabled="isChangingStatus"
          >
            <CheckCircle :size="18" /> <span>Confirmar Entrega</span>
          </button>

          <button
            v-if="!relatedInvoice"
            class="btn btn-danger btn-sm"
            @click="confirmDeleteDeliveryNote"
            :disabled="isChangingStatus"
          >
            <Trash2 :size="18" /> <span>Eliminar Albarán</span>
          </button>
          
          <button
            v-if="['DELIVERED', 'ENTREGADO'].includes(deliveryNote.status) && !relatedInvoice"
            class="btn btn-success btn-sm"
            @click="createInvoiceFromDeliveryNote"
            :disabled="isCreatingInvoice"
          >
            <Receipt :size="18" /> <span>Facturar Albarán</span>
          </button>
        </div>
      </div>
    </template>

    <!-- 3. SUMMARY -->
    <template #summary>
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><User :size="20" /></div>
          <div class="tag-content"><label>Cliente</label><strong>{{ partyName }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><Calendar :size="20" /></div>
          <div class="tag-content"><label>Fecha Entrega</label><strong>{{ formatDate(deliveryNote.deliveryDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><ShoppingCart :size="20" /></div>
          <div class="tag-content">
            <label>Pedido Origen</label>
            <strong>{{ orderNumber || formatOrderId(deliveryNote.salesOrderId) }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><Package :size="20" /></div>
          <div class="tag-content">
            <label>Bultos / Líneas</label>
            <strong>{{ deliveryNote.lineItems?.length || 0 }} ítems</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- 4. RELATED -->
    <template #related v-if="mode === 'detail'">
      <div class="related-history-grid">
        <router-link :to="`/sales/orders/${deliveryNote.salesOrderId}`" class="related-tag-card highlight-info">
          <div class="tag-icon"><ShoppingCart :size="18" /></div>
          <div class="tag-content">
            <label>Pedido de Venta Origen</label>
            <strong>{{ orderNumber || 'Ver Pedido' }}</strong>
          </div>
          <ExternalLink :size="18" class="jump-icon" />
        </router-link>

        <router-link v-if="relatedInvoice" :to="`/sales/invoices/${relatedInvoice.id}`" class="related-tag-card">
          <div class="tag-icon success"><Receipt :size="18" /></div>
          <div class="tag-content">
            <label>Factura Generada</label>
            <strong>{{ relatedInvoice.invoiceNumber }}</strong>
          </div>
          <ExternalLink :size="18" class="jump-icon" />
        </router-link>
      </div>
    </template>

    <!-- 5. MAIN CONTENT -->
    <FormSection title="Identificación del Cliente" icon="person">
      <DataRow label="Nombre del Cliente" :value="partyName" icon="person" />
      <DataRow v-if="deliveryNote.taxId" label="NIF/CIF" :value="deliveryNote.taxId" is-mono />
    </FormSection>

    <FormSection title="Detalles de Entrega" icon="local_shipping">
      <div v-if="mode === 'detail'">
        <DataRow label="Fecha de Entrega" :value="formatDate(deliveryNote.deliveryDate)" icon="calendar_today" />
        <DataRow v-if="deliveryNote.deliveryAddress" label="Dirección de Entrega" icon="location_on">
          <div class="address-content">
            <p>{{ deliveryNote.deliveryAddress.street }}</p>
            <p>{{ deliveryNote.deliveryAddress.postalCode }} {{ deliveryNote.deliveryAddress.city }}</p>
            <p v-if="deliveryNote.deliveryAddress.province">{{ deliveryNote.deliveryAddress.province }}</p>
            <p v-if="deliveryNote.deliveryAddress.country">{{ deliveryNote.deliveryAddress.country }}</p>
          </div>
        </DataRow>
        <DataRow label="Observaciones del Albarán" icon="notes">
          <p class="notes-text">{{ deliveryNote.notes || 'Sin observaciones.' }}</p>
        </DataRow>
      </div>
      <div v-else>
        <div class="form-row">
          <div class="form-group">
            <label>Fecha de Entrega *</label>
            <input v-model="formData.deliveryDate" type="date" class="form-input" required />
          </div>
        </div>
        <div class="address-edit-group mt-4">
          <label class="form-label">Dirección de Entrega</label>
          <input v-model="formData.address.street" type="text" class="form-input mb-2" placeholder="Calle y número" />
          <div class="form-row">
            <input v-model="formData.address.postalCode" type="text" class="form-input" placeholder="CP" />
            <input v-model="formData.address.city" type="text" class="form-input" placeholder="Ciudad" />
          </div>
        </div>
        <div class="form-group mt-4">
          <label>Observaciones</label>
          <textarea v-model="formData.notes" class="form-textarea" rows="3" placeholder="Instrucciones para el transportista..."></textarea>
        </div>
      </div>
    </FormSection>

    <FormSection title="Líneas del Albarán" icon="list_alt">
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>Producto / Referencia</th>
              <th class="text-center">Cant. Entregada</th>
              <th>Estado / Notas</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in deliveryNote.lineItems" :key="item.id">
              <td>
                <div class="product-info-cell">
                  <Package :size="18" class="icon-secondary" />
                  <div class="content">
                    <strong>{{ item.productName || formatVariantId(item.productVariantId) }}</strong>
                    <code v-if="item.variantSku" class="code-badge ml-2">{{ item.variantSku }}</code>
                  </div>
                </div>
              </td>
              <td class="text-center">
                <template v-if="mode === 'detail'">
                  <strong class="text-success" style="font-size: 1.1rem">{{ item.deliveredQuantity }}</strong>
                </template>
                <input v-else v-model.number="item.deliveredQuantity" type="number" class="form-input-sm w-24 text-center" />
              </td>
              <td class="text-muted italic">Correcto</td>
            </tr>
          </tbody>
        </table>
      </div>
    </FormSection>

    <FormSection v-if="mode === 'detail'" title="Conformidad y Firmas" icon="history_edu">
      <div class="signatures-grid">
        <div class="signature-box">
          <label class="form-label">Recibido por (Cliente)</label>
          <div class="signature-area">
            <template v-if="deliveryNote.signatures?.customer">
              <CheckCircle :size="18" class="text-success" />
              <div class="sig-info">
                <strong>{{ deliveryNote.signatures.customer.name }}</strong>
                <small>{{ formatDateTime(deliveryNote.signatures.customer.timestamp) }}</small>
              </div>
            </template>
            <span v-else class="text-muted">Pendiente de firma del receptor</span>
          </div>
        </div>
        <div class="signature-box">
          <label class="form-label">Entregado por (Logística)</label>
          <div class="signature-area">
            <template v-if="deliveryNote.signatures?.driver">
              <CheckCircle :size="18" class="text-success" />
              <div class="sig-info">
                <strong>{{ deliveryNote.signatures.driver.name }}</strong>
                <small>{{ formatDateTime(deliveryNote.signatures.driver.timestamp) }}</small>
              </div>
            </template>
            <span v-else class="text-muted">Pendiente de firma del transportista</span>
          </div>
        </div>
      </div>
    </FormSection>

    <!-- 6. FOOTER -->
    <template #footer v-if="mode === 'detail' && deliveryNote">
      <div class="audit-info">
        <p>Documento oficial de entrega generado por TramaTex.</p>
        <p>ID único del documento: <code>{{ deliveryNote.id }}</code></p>
      </div>
    </template>
    </BaseEntityPage>

    <!-- PORTAL DE IMPRESIÓN (Solo visible en @media print) -->
    <div class="print-container">
    <PrintDocument
      v-if="deliveryNote"
      type="DELIVERY_NOTE"
      :number="deliveryNote.deliveryNoteNumber || deliveryNote.delivery_note_number"
      :date="deliveryNote.deliveryDate || deliveryNote.delivery_date"
      :customer-name="partyName"
      :customer-tax-id="deliveryNote.taxId"
      :items="deliveryNote.lineItems"
      :totals="{ subtotal: 0, taxAmount: 0, total: 0 }"
      :notes="deliveryNote.notes"
    />
    </div>

    <!-- DIÁLOGOS DE CONFIRMACIÓN -->
    <BaseDialog
      :show="showStatusConfirm"
      :title="statusConfirmTitle"
      :icon="statusConfirmIcon"
      :confirm-text="statusConfirmText"
      :confirm-class="statusConfirmClass"
      :is-confirming="isChangingStatus"
      @close="showStatusConfirm = false"
      @confirm="executeStatusChange"
    >
      <p v-html="statusConfirmMessage"></p>
    </BaseDialog>

    <!-- DIÁLOGO DE CONFIRMACIÓN DE FACTURACIÓN -->
    <BaseDialog
      :show="showInvoiceConfirm"
      title="Confirmar Facturación"
      icon="receipt_long"
      confirm-text="Generar Factura"
      confirm-class="btn-success"
      :is-confirming="isCreatingInvoice"
      @close="showInvoiceConfirm = false"
      @confirm="confirmCreateInvoice"
    >
      <div class="confirm-dialog-body">
        <p>Está a punto de <strong>generar una factura oficial</strong> para este albarán.</p>
        <div class="info-notice mt-4">
          <Info :size="18" />
          <p>Esta operación creará un nuevo documento contable y vinculará permanentemente este albarán.</p>
        </div>
        <p class="mt-4 text-secondary italic">¿Desea continuar?</p>
      </div>
    </BaseDialog>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute, useRouter, RouterLink } from 'vue-router';
import { 
  AlertCircle, 
  Truck, 
  Printer, 
  Pencil, 
  RefreshCw, 
  Save, 
  CheckCircle, 
  Trash2, 
  Receipt, 
  User, 
  Calendar, 
  ShoppingCart, 
  Package, 
  ExternalLink, 
  Info 
} from 'lucide-vue-next';

import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import PrintDocument from '@/components/sales/PrintDocument.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { useToastStore } from '@/stores/toast';
import '@/assets/sales-print.css';

const route = useRoute();
const router = useRouter();
const toastStore = useToastStore();

const mode = ref('detail');
const deliveryNote = ref(null);
const isLoading = ref(false);
const isSaving = ref(false);
const isChangingStatus = ref(false);
const isCreatingInvoice = ref(false);
const showInvoiceConfirm = ref(false);
const error = ref('');
const partyName = ref('Cargando...');
const orderNumber = ref(null);
const relatedInvoice = ref(null);

// Status change confirm dialog state
const showStatusConfirm = ref(false);
const statusConfirmTitle = ref('');
const statusConfirmMessage = ref('');
const statusConfirmIcon = ref('');
const statusConfirmText = ref('');
const statusConfirmClass = ref('');
const pendingStatus = ref('');

const formData = reactive({
  deliveryDate: '',
  notes: '',
  address: { street: '', city: '', postalCode: '', province: '', country: '' }
});

onMounted(() => {
  fetchDeliveryNote();
});

async function fetchDeliveryNote() {
  const noteId = route.params.id;
  if (!noteId) return;
  isLoading.value = true;
  try {
    deliveryNote.value = await salesApi.getDeliveryNote(noteId);
    await Promise.all([loadPartyName(), loadOrderNumber(), loadRelatedInvoice()]);
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar el albarán';
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyName() {
  if (!deliveryNote.value?.partyId) return;
  try {
    const party = await partyApi.getParty(deliveryNote.value.partyId);
    partyName.value = party.name || 'Sin nombre';
  } catch (err) {}
}

async function loadOrderNumber() {
  if (!deliveryNote.value?.salesOrderId) return;
  try {
    const order = await salesApi.getOrder(deliveryNote.value.salesOrderId);
    orderNumber.value = order.orderNumber;
  } catch (err) {}
}

async function loadRelatedInvoice() {
  if (!deliveryNote.value?.id) return;
  try {
    const res = await salesApi.listInvoices({ deliveryNoteId: deliveryNote.value.id });
    const invoices = Array.isArray(res) ? res : (res.data || []);
    if (invoices.length > 0) relatedInvoice.value = invoices[0];
  } catch (err) {
    console.error('Error loading related invoice:', err);
  }
}

function enterEditMode() {
  if (!deliveryNote.value) return;
  formData.deliveryDate = deliveryNote.value.deliveryDate ? new Date(deliveryNote.value.deliveryDate).toISOString().split('T')[0] : '';
  formData.notes = deliveryNote.value.notes || '';
  if (deliveryNote.value.deliveryAddress) {
    Object.assign(formData.address, deliveryNote.value.deliveryAddress);
  }
  mode.value = 'edit';
}

async function saveDeliveryNote() {
  isSaving.value = true;
  try {
    toastStore.info('Función de actualización de albarán en desarrollo (Backend MVP limitado)');
    mode.value = 'detail';
    await fetchDeliveryNote();
  } catch (err) {
    toastStore.error('Error al guardar: ' + err.message);
  } finally {
    isSaving.value = false;
  }
}

function markAsDelivered() {
  statusConfirmTitle.value = 'Confirmar Entrega';
  statusConfirmMessage.value = '¿Desea marcar este albarán como <strong>entregado</strong>? Esta acción confirmará que el cliente ha recibido la mercancía.';
  statusConfirmIcon.value = 'check_circle';
  statusConfirmText.value = 'Confirmar Entrega';
  statusConfirmClass.value = 'btn-success';
  pendingStatus.value = 'ENTREGADO';
  showStatusConfirm.value = true;
}

function confirmDeleteDeliveryNote() {
  statusConfirmTitle.value = 'Eliminar Albarán';
  statusConfirmMessage.value = '¿Realmente desea <strong>ELIMINAR</strong> este albarán? Las cantidades volverán a estar disponibles en el pedido original para ser albaraneadas de nuevo.';
  statusConfirmIcon.value = 'delete';
  statusConfirmText.value = 'Eliminar Permanentemente';
  statusConfirmClass.value = 'btn-danger';
  pendingStatus.value = 'DELETE_ACTION'; // Special marker for delete
  showStatusConfirm.value = true;
}

async function executeStatusChange() {
  isChangingStatus.value = true;
  try {
    if (pendingStatus.value === 'DELETE_ACTION') {
      const orderId = deliveryNote.value.salesOrderId;
      await salesApi.deleteDeliveryNote(deliveryNote.value.id);
      showStatusConfirm.value = false;
      router.push(orderId ? `/sales/orders/${orderId}` : '/sales/delivery-notes');
    } else {
      deliveryNote.value = await salesApi.changeDeliveryNoteStatus(deliveryNote.value.id, pendingStatus.value);
      showStatusConfirm.value = false;
    }
  } catch (err) {
    toastStore.error(err?.message || 'Error al procesar la solicitud');
  } finally {
    isChangingStatus.value = false;
  }
}

function createInvoiceFromDeliveryNote() {
  showInvoiceConfirm.value = true;
}

async function confirmCreateInvoice() {
  isCreatingInvoice.value = true;
  try {
    const now = new Date();
    const dueDate = new Date(now);
    dueDate.setDate(dueDate.getDate() + 30);
    const newInvoice = await salesApi.createInvoice({
      partyId: deliveryNote.value.partyId,
      deliveryNoteIds: [deliveryNote.value.id],
      invoiceDate: now.toISOString(),
      dueDate: dueDate.toISOString(),
    });
    showInvoiceConfirm.value = false;
    router.push(`/sales/invoices/${newInvoice.id}`);
  } catch (err) {
    toastStore.error(err?.message || 'Error al crear factura');
  } finally {
    isCreatingInvoice.value = false;
  }
}

function printDeliveryNote() { window.print(); }

function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '—'; }
function formatDateTime(d) { return d ? new Date(d).toLocaleString('es-ES', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }) : '—'; }
function formatOrderId(id) { return id ? id.substring(0, 8) : '—'; }
function formatVariantId(id) { return id ? id.substring(0, 8) : '—'; }
</script>

<style scoped>
@import "@/design-system/_sections.css";

.action-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; box-shadow: var(--box-shadow-sm); margin: 0; }
.status-badge { padding: 0.4rem 1rem; font-size: 0.85rem; font-weight: 800; letter-spacing: 0.05em; }
.toolbar-buttons { display: flex; gap: 0.75rem; }

.address-content p { margin: 0.2rem 0; color: var(--color-text-primary); }
.notes-text { font-style: italic; color: var(--color-text-secondary); margin: 0; }

.product-info-cell { display: flex; align-items: center; gap: 0.75rem; }
.product-info-cell .content { display: flex; flex-direction: column; }

.signatures-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 1.5rem; }
.signature-box { padding: 1.25rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; }
.signature-area { min-height: 80px; background: var(--color-background); border: 2px dashed var(--color-border); border-radius: 8px; display: flex; align-items: center; justify-content: center; gap: 1rem; padding: 1rem; }
.sig-info { display: flex; flex-direction: column; line-height: 1.2; }
.sig-info strong { font-size: 0.9rem; color: var(--color-text-primary); }
.sig-info small { font-size: 0.75rem; color: var(--color-text-secondary); }
.form-label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.75rem; }

.audit-info { color: var(--color-text-secondary); font-size: 0.8rem; font-style: italic; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; }

.info-notice { display: flex; gap: 0.75rem; padding: 1rem; background: rgba(59, 130, 246, 0.1); border-radius: 8px; color: #1e40af; font-size: 0.9rem; }
.info-notice svg { color: #2563eb; }

/* ESTILOS DE IMPRESIÓN PROFESIONAL */
.print-container { display: none; }

@media print {
  .no-print { display: none !important; }
  .print-container { display: block !important; position: absolute; left: 0; top: 0; width: 100%; }
}
</style>
