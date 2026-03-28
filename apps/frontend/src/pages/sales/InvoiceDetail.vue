<template>
  <Navbar class="no-print" />
  
  <BaseEntityPage v-if="isLoading" class="no-print">
    <template #header>
      <PageHeader title="Cargando..." :breadcrumbs="[{ label: 'Ventas', to: '/sales/invoices' }, { label: 'Facturas' }]" />
    </template>
    <div class="loading-state card">
      <div class="spinner"></div>
      <p>Cargando información de la factura...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="error" class="no-print">
    <template #header>
      <PageHeader title="Error" :breadcrumbs="[{ label: 'Ventas', to: '/sales/invoices' }, { label: 'Facturas' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-icon-wrapper error">
        <span class="material-symbols-outlined">error</span>
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/sales/invoices')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="invoice" class="no-print">
    <!-- 1. IDENTITY HEADER -->
    <template #header>
      <PageHeader 
        :title="mode === 'edit' ? `Editando Factura ${invoice.invoiceNumber}` : `Factura ${invoice.invoiceNumber}`" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/invoices' }, { label: 'Facturas', to: '/sales/invoices' }, { label: invoice.invoiceNumber }]"
      >
        <template #icon>
          <span class="material-symbols-outlined">receipt</span>
        </template>
        <template #actions>
          <template v-if="mode === 'detail'">
            <button class="btn btn-outline" @click="printInvoice">
              <span class="material-symbols-outlined">print</span> <span>Imprimir</span>
            </button>
            <button v-if="invoice.status === 'BORRADOR'" class="btn btn-primary" @click="enterEditMode">
              <span class="material-symbols-outlined">edit</span> <span>Editar</span>
            </button>
          </template>
          <template v-else>
            <button class="btn btn-outline" @click="mode = 'detail'" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-secondary" @click="saveInvoice" :disabled="isSaving">
              <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
              <span>{{ isSaving ? 'Guardando...' : 'Guardar Cambios' }}</span>
            </button>
          </template>
        </template>
      </PageHeader>
    </template>

    <!-- 2. TOOLBAR -->
    <template #toolbar v-if="mode === 'detail'">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', `status-${salesApi.getStatusClass(invoice.status)}`]">
            {{ salesApi.getStatusLabel(invoice.status) }}
          </span>
          <span :class="['type-badge-inline', `type-${invoice.type?.toLowerCase()}`]">
            {{ getTypeLabel(invoice.type) }}
          </span>
        </div>
        <div class="toolbar-buttons">
          <button
            v-if="invoice.status === 'BORRADOR'"
            class="btn btn-success btn-sm"
            @click="emitInvoice"
            :disabled="isChangingStatus"
          >
            <span class="material-symbols-outlined">send</span> <span>Emitir Factura</span>
          </button>
          <button
            v-if="invoice.status === 'EMITIDA' || invoice.status === 'VENCIDA'"
            class="btn btn-success btn-sm"
            @click="markAsPaid"
            :disabled="isChangingStatus"
          >
            <span class="material-symbols-outlined">payments</span> <span>Registrar Cobro</span>
          </button>
        </div>
      </div>
    </template>

    <!-- 3. SUMMARY -->
    <template #summary>
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><span class="material-symbols-outlined">person</span></div>
          <div class="tag-content"><label>Cliente</label><strong>{{ partyName }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><span class="material-symbols-outlined">calendar_today</span></div>
          <div class="tag-content"><label>Fecha Emisión</label><strong>{{ formatDate(invoice.issueDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><span class="material-symbols-outlined">event_busy</span></div>
          <div class="tag-content"><label>Vencimiento</label><strong>{{ formatDate(invoice.dueDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><span class="material-symbols-outlined">payments</span></div>
          <div class="tag-content">
            <label>Total Factura</label>
            <strong class="amount">{{ salesApi.formatMoney(invoice.total) }}</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- 4. RELATED -->
    <template #related v-if="mode === 'detail' && (relatedOrders.length > 0 || relatedDeliveryNotes.length > 0)">
      <div class="related-history-grid">
        <router-link v-for="order in relatedOrders" :key="order.id" :to="`/sales/orders/${order.id}`" class="related-tag-card highlight-info">
          <div class="tag-icon"><span class="material-symbols-outlined">shopping_cart</span></div>
          <div class="tag-content">
            <label>Pedido de Origen</label>
            <strong>{{ order.orderNumber || formatId(order.id) }}</strong>
          </div>
          <span class="material-symbols-outlined jump-icon">open_in_new</span>
        </router-link>

        <router-link v-for="dn in relatedDeliveryNotes" :key="dn.id" :to="`/sales/delivery-notes/${dn.id}`" class="related-tag-card">
          <div class="tag-icon"><span class="material-symbols-outlined">local_shipping</span></div>
          <div class="tag-content">
            <label>Albarán de Origen</label>
            <strong>{{ dn.deliveryNoteNumber || formatId(dn.id) }}</strong>
          </div>
          <span class="material-symbols-outlined jump-icon">open_in_new</span>
        </router-link>
      </div>
    </template>

    <!-- 5. MAIN CONTENT -->
    <div ref="invoiceContentRef">
      <FormSection title="Identificación del Cliente" icon="person">
        <DataRow label="Nombre del Cliente" :value="partyName" icon="person" />
        <DataRow v-if="invoice.taxId" label="NIF/CIF" :value="invoice.taxId" is-mono />
        <DataRow label="ID de Cliente" :value="invoice.partyId" is-mono />
      </FormSection>

      <FormSection title="Condiciones y Notas" icon="description">
        <div v-if="mode === 'detail'">
          <DataRow v-if="invoice.paymentTerms" label="Condiciones de Pago" icon="payments">
            <p class="notes-text">{{ invoice.paymentTerms }}</p>
          </DataRow>
          <DataRow v-if="invoice.notes" label="Observaciones" icon="notes">
            <p class="notes-text">{{ invoice.notes }}</p>
          </DataRow>
        </div>
        <div v-else>
          <div class="form-group">
            <label>Fecha de Vencimiento *</label>
            <input v-model="formData.dueDate" type="date" class="form-input" required />
          </div>
          <div class="form-group mt-4">
            <label>Condiciones de Pago</label>
            <textarea v-model="formData.paymentTerms" class="form-textarea" rows="2" placeholder="Ej: Transferencia 30 días..."></textarea>
          </div>
          <div class="form-group mt-4">
            <label>Observaciones de Factura</label>
            <textarea v-model="formData.notes" class="form-textarea" rows="3" placeholder="Información adicional..."></textarea>
          </div>
        </div>
      </FormSection>

      <FormSection title="Líneas de la Factura" icon="list_alt">
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Producto / Referencia</th>
                <th class="text-center">Cant.</th>
                <th class="align-right">P. Unitario</th>
                <th class="text-center">Dto. %</th>
                <th class="text-center">IVA %</th>
                <th class="align-right">Subtotal</th>
                <th class="align-right">Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in invoice.lineItems" :key="item.id || item.productVariantID">
                <td>
                  <div class="product-info-cell">
                    <span class="material-symbols-outlined icon-secondary">inventory_2</span>
                    <div class="content">
                      <strong>{{ buildDisplayName(item) }}</strong>
                      <code class="code-badge ml-2">{{ item.variantSku || formatId(item.productVariantId || item.productVariantID) }}</code>
                    </div>
                  </div>
                </td>
                <td class="text-center">{{ item.quantity }}</td>
                <td class="align-right">{{ salesApi.formatUnitPrice(item.unitPrice) }}</td>
                <td class="text-center">{{ item.discountAmount && item.unitPrice?.amount ? ((item.discountAmount.amount / item.unitPrice.amount) * 100).toFixed(2) + '%' : '—' }}</td>
                <td class="text-center">{{ typeof item.taxRate === 'number' ? `${item.taxRate}%` : '21%' }}</td>
                <td class="align-right">{{ salesApi.formatMoney(item.subtotal) }}</td>
                <td class="align-right"><strong>{{ salesApi.formatMoney(item.total) }}</strong></td>
              </tr>
            </tbody>
          </table>
        </div>
      </FormSection>

      <FormSection title="Resumen Económico" icon="payments">
        <DataRow label="Base Imponible" :value="salesApi.formatMoney(invoice.subtotal)" />
        <DataRow label="Impuestos (IVA)" :value="salesApi.formatMoney(invoice.taxAmount)" />
        <DataRow label="TOTAL FACTURA" :value="salesApi.formatMoney(invoice.total)" highlight />
      </FormSection>
    </div>

    <!-- 6. FOOTER -->
    <template #footer v-if="mode === 'detail' && invoice">
      <div class="audit-info">
        <p>Factura legal generada por TramaTex ERP.</p>
        <p>Código de verificación del documento: <code>{{ invoice.id }}</code></p>
      </div>
    </template>
  </BaseEntityPage>

  <!-- CAPA DE IMPRESIÓN (SOLO VISIBLE AL IMPRIMIR) -->
  <div class="print-container">
    <PrintDocument
      v-if="invoice"
      type="INVOICE"
      :number="invoice.invoiceNumber"
      :date="invoice.invoiceDate"
      :customer-name="partyName"
      :customer-tax-id="invoice.taxId"
      :items="invoice.lineItems"
      :totals="{ subtotal: invoice.subtotal, taxAmount: invoice.taxAmount, total: invoice.total }"
      :notes="invoice.notes"
    />
  </div>

  <!-- Post-Issue Actions Modal -->
  <Transition name="fade">
    <div v-if="showPostIssueModal" class="modal-backdrop no-print">
      <div class="modal card w-modal-md">
        <div class="modal-header">
          <span class="material-symbols-outlined text-success">check_circle</span>
          <h2>Factura Emitida con Éxito</h2>
          <button class="btn-icon ml-auto" @click="showPostIssueModal = false"><span class="material-symbols-outlined">close</span></button>
        </div>
        <div class="modal-body">
          <p class="mb-4">La factura <strong>{{ invoice?.invoiceNumber }}</strong> ya es oficial y puede ser enviada al cliente.</p>
          <div class="post-issue-actions">
            <button class="btn btn-primary w-full justify-center mb-3" @click="postIssuePrint">
              <span class="material-symbols-outlined">print</span> <span>Imprimir Factura</span>
            </button>
            <button class="btn btn-outline w-full justify-center mb-2" @click="postIssueEmail">
              <span class="material-symbols-outlined">mail</span> <span>Enviar por Email</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter, RouterLink } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import PrintDocument from '@/components/sales/PrintDocument.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';

const route = useRoute();
const router = useRouter();

const mode = ref('detail');
const invoice = ref(null);
const isLoading = ref(false);
const isSaving = ref(false);
const isChangingStatus = ref(false);
const error = ref('');
const partyName = ref('');
const relatedOrders = ref([]);
const relatedDeliveryNotes = ref([]);
const showPostIssueModal = ref(false);

const formData = reactive({
  dueDate: '',
  paymentTerms: '',
  notes: ''
});

onMounted(() => {
  fetchInvoice();
});

async function fetchInvoice() {
  const invoiceId = route.params.id;
  isLoading.value = true;
  try {
    invoice.value = await salesApi.getInvoice(invoiceId);
    await Promise.all([loadPartyName(), loadRelatedOrders(), loadRelatedDeliveryNotes()]);
  } catch (err) {
    error.value = err?.message || 'Error al cargar factura';
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyName() {
  if (!invoice.value?.partyId) return;
  try {
    const party = await partyApi.getParty(invoice.value.partyId);
    partyName.value = party.name || 'Sin nombre';
  } catch { partyName.value = 'Desconocido'; }
}

async function loadRelatedOrders() {
  const ids = invoice.value?.salesOrderIds;
  if (!ids?.length) return;
  relatedOrders.value = await Promise.all(ids.map(id => salesApi.getOrder(id).catch(() => ({ id, orderNumber: null }))));
}

async function loadRelatedDeliveryNotes() {
  const ids = invoice.value?.deliveryNoteIds;
  if (!ids?.length) return;
  relatedDeliveryNotes.value = await Promise.all(ids.map(id => salesApi.getDeliveryNote(id).catch(() => ({ id, deliveryNoteNumber: null }))));
}

function enterEditMode() {
  if (!invoice.value) return;
  formData.dueDate = invoice.value.dueDate ? new Date(invoice.value.dueDate).toISOString().split('T')[0] : '';
  formData.paymentTerms = invoice.value.paymentTerms || '';
  formData.notes = invoice.value.notes || '';
  mode.value = 'edit';
}

async function saveInvoice() {
  isSaving.value = true;
  try {
    const payload = {
      dueDate: new Date(formData.dueDate).toISOString(),
      paymentTerms: formData.paymentTerms,
      notes: formData.notes
    };
    alert('Función de actualización de factura en desarrollo (Backend MVP limitado)');
    mode.value = 'detail';
    await fetchInvoice();
  } catch (err) {
    alert('Error al guardar: ' + err.message);
  } finally {
    isSaving.value = false;
  }
}

async function emitInvoice() {
  if (!confirm('¿Emitir esta factura? Esta acción generará el número de factura definitivo.')) return;
  isChangingStatus.value = true;
  try {
    invoice.value = await salesApi.changeInvoiceStatus(invoice.value.id, 'EMITIDA');
    showPostIssueModal.value = true;
  } catch (err) { alert(err?.message); }
  finally { isChangingStatus.value = false; }
}

async function markAsPaid() {
  if (!confirm('¿Registrar el cobro de esta factura?')) return;
  isChangingStatus.value = true;
  try {
    invoice.value = await salesApi.changeInvoiceStatus(invoice.value.id, 'PAGADA');
  } catch (err) { alert(err?.message); }
  finally { isChangingStatus.value = false; }
}

function printInvoice() { window.print(); }

function postIssuePrint() { showPostIssueModal.value = false; window.print(); }
async function postIssueEmail() { showPostIssueModal.value = false; await generateInvoicePdf(); openMailClient(); }

const invoiceContentRef = ref(null);

async function generateInvoicePdf() {
  const { default: html2pdf } = await import('html2pdf.js');
  const el = document.querySelector('.print-container');
  if (!el) return;
  const filename = `Factura_${invoice.value.invoiceNumber || invoice.value.id}.pdf`;
  await html2pdf().set({ margin: 10, filename, image: { type: 'jpeg', quality: 0.98 }, html2canvas: { scale: 2 }, jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' } }).from(el).save();
}

function openMailClient() {
  const inv = invoice.value;
  const subject = encodeURIComponent(`Factura ${inv.invoiceNumber || inv.id} de TramaTex`);
  const body = encodeURIComponent(`Estimado/a cliente,\n\nAdjunto enviamos su factura correspondiente a su pedido.\n\nGracias por su confianza.\n\nUn saludo.`);
  window.open(`mailto:?subject=${subject}&body=${body}`, '_self');
}

function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '—'; }
function formatId(id) { return id ? id.substring(0, 8) : '—'; }
function getTypeLabel(t) { const map = { STANDARD: 'Venta Estándar', SIMPLIFIED: 'Venta Simplificada' }; return map[t] || t; }
function buildDisplayName(item) { return (item.productName || item.displayName || 'Producto') + (item.optionConfiguration ? ' - ' + Object.values(item.optionConfiguration).join(', ') : ''); }
</script>

<style scoped>
@import "@/design-system/_sections.css";

.overview-tags-row, .related-history-grid { display: flex; flex-wrap: wrap; gap: 1rem; }
.related-history-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); }

.summary-tag { flex: 1; min-width: 240px; padding: 0.6rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 0.75rem; box-shadow: var(--box-shadow-sm); }
.related-tag-card { padding: 0.6rem 1rem; background: var(--color-background); border: 1px solid var(--color-border); border-left: 4px solid var(--color-secondary); border-radius: 10px; display: flex; align-items: center; gap: 0.75rem; text-decoration: none; position: relative; transition: all 0.2s ease; }
.related-tag-card.highlight-info { border-left-color: #2563eb; }
.related-tag-card:hover { background: white; transform: translateX(2px) translateY(-1px); box-shadow: var(--box-shadow-md); }
.related-tag-card:hover strong { color: var(--color-primary); text-decoration: underline; }

.tag-icon { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; background: rgba(0,0,0,0.03); color: var(--color-text-secondary); }
.tag-icon .material-symbols-outlined { font-size: 22px; }

.icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }

.tag-content { display: flex; flex-direction: column; gap: 0.15rem; line-height: 1.2; }
.tag-content label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.025em; }
.tag-content strong { font-size: 0.95rem; color: var(--color-text-primary); }
.amount { color: #16a34a !important; font-size: 1.15rem !important; }

.jump-icon { font-size: 18px; color: var(--color-text-secondary); opacity: 0.5; margin-left: auto; transition: all 0.2s; }
.related-tag-card:hover .jump-icon { opacity: 1; color: var(--color-primary); transform: scale(1.1); }

.action-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; box-shadow: var(--box-shadow-sm); margin-bottom: 0; }
.status-badge { padding: 0.4rem 1rem; font-size: 0.85rem; font-weight: 800; }
.toolbar-buttons { display: flex; gap: 0.75rem; }

.type-badge-inline { margin-left: 1rem; padding: 0.4rem 0.75rem; font-size: 0.75rem; font-weight: 800; border-radius: 6px; background: var(--color-background); color: var(--color-text-secondary); border: 1px solid var(--color-border); text-transform: uppercase; }
.type-standard { border-left: 3px solid #2563eb; }
.type-simplified { border-left: 3px solid #d97706; }

.product-info-cell { display: flex; align-items: center; gap: 0.75rem; }
.product-info-cell .content { display: flex; flex-direction: column; }

.notes-text { font-style: italic; color: var(--color-text-secondary); margin: 0; }
.audit-info { color: var(--color-text-secondary); font-size: 0.8rem; font-style: italic; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; }

.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.w-modal-md { width: 90%; max-width: 500px; }
.post-issue-actions { padding-top: 1rem; }

.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }

/* ESTILOS DE IMPRESIÓN PROFESIONAL */
.print-container { display: none; }

@media print {
  .no-print { display: none !important; }
  .print-container { display: block !important; position: absolute; left: 0; top: 0; width: 100%; }
}
</style>
