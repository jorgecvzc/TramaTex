<template>
  <Navbar class="no-print" />
  
  <BaseEntityPage v-if="isLoading" class="no-print">
    <template #header>
      <PageHeader title="Cargando..." :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Pedidos' }]" />
    </template>
    <div class="loading-state card">
      <div class="spinner"></div>
      <p>Cargando información del pedido...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="error" class="no-print">
    <template #header>
      <PageHeader title="Error" :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Pedidos' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-icon-wrapper error">
        <span class="material-symbols-outlined">error</span>
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/sales/orders')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="order || mode === 'create'" class="no-print">
    <!-- 1. IDENTITY HEADER -->
    <template #header>
      <PageHeader 
        :title="mode === 'create' ? 'Nuevo Pedido' : (mode === 'edit' ? `Editando Pedido ${order?.orderNumber}` : `Pedido ${order?.orderNumber}`)" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Pedidos', to: '/sales/orders' }, { label: mode === 'create' ? 'Crear' : order?.orderNumber }]"
      >
        <template #icon>
          <span class="material-symbols-outlined">shopping_cart</span>
        </template>
        <template #actions>
          <template v-if="mode === 'detail'">
            <button class="btn btn-outline" @click="printOrder">
              <span class="material-symbols-outlined">print</span> <span>Imprimir</span>
            </button>
            <button v-if="canEdit" class="btn btn-primary" @click="enterEditMode">
              <span class="material-symbols-outlined">edit</span> <span>Editar Pedido</span>
            </button>
          </template>
          <template v-else>
            <button class="btn btn-outline" @click="exitEditMode" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-secondary" @click="saveOrder" :disabled="isSaving">
              <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
              <span>{{ isSaving ? 'Guardando...' : 'Guardar Pedido' }}</span>
            </button>
          </template>
        </template>
      </PageHeader>
    </template>

    <!-- 2. TOOLBAR -->
    <template #toolbar v-if="mode === 'detail' && order">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', `status-${salesApi.getStatusClass(order.status)}`]">
            {{ salesApi.getStatusLabel(order.status) }}
          </span>
        </div>
        <div class="toolbar-buttons">
          <button v-if="['PENDING', 'PENDIENTE'].includes(order.status)" class="btn btn-success btn-sm" @click="confirmOrder">
            <span class="material-symbols-outlined">check_circle</span> <span>Confirmar Pedido</span>
          </button>
          <button v-if="canCreateDeliveryNote" class="btn btn-success btn-sm" @click="showDeliveryNoteModal = true">
            <span class="material-symbols-outlined">local_shipping</span> <span>Generar Albarán</span>
          </button>
          <button v-if="canCreateInvoice" class="btn btn-success btn-sm" @click="createInvoiceFromOrder">
            <span class="material-symbols-outlined">receipt_long</span> <span>Facturar</span>
          </button>
          <button v-if="canReopen" class="btn btn-primary btn-sm" @click="triggerReopen">
            <span class="material-symbols-outlined">settings_backup_restore</span> 
            <span>{{ ['CONFIRMED', 'CONFIRMADO', 'EN_PREPARACION'].includes(order.status) ? 'Modificar Pedido' : 'Reabrir' }}</span>
          </button>
          <button v-if="canCancel" class="btn btn-danger btn-sm" @click="triggerCancel">
            <span class="material-symbols-outlined">cancel</span> <span>Anular</span>
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
          <div class="tag-content"><label>Fecha Pedido</label><strong>{{ formatDate(mode === 'detail' ? order?.orderDate : formData.orderDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><span class="material-symbols-outlined">local_shipping</span></div>
          <div class="tag-content"><label>Entrega Est.</label><strong>{{ formatDate(mode === 'detail' ? order?.deliveryDate : formData.deliveryDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><span class="material-symbols-outlined">payments</span></div>
          <div class="tag-content">
            <label>Total Pedido</label>
            <strong class="amount">{{ salesApi.formatMoney(liveTotals.total) }}</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- 4. RELATED -->
    <template #related v-if="mode === 'detail' && (order?.quoteId || deliveryNotes.length > 0 || invoices.length > 0)">
      <div class="related-history-grid">
        <router-link v-if="order?.quoteId" :to="`/sales/quotes/${order.quoteId}`" class="related-tag-card highlight-info">
          <div class="tag-icon"><span class="material-symbols-outlined">description</span></div>
          <div class="tag-content">
            <label>Presupuesto Origen</label>
            <strong>{{ order.sourceQuoteNumber || 'Ver Presupuesto' }}</strong>
          </div>
          <span class="material-symbols-outlined jump-icon">open_in_new</span>
        </router-link>

        <router-link v-for="note in deliveryNotes" :key="note.id" :to="`/sales/delivery-notes/${note.id}`" class="related-tag-card">
          <div class="tag-icon"><span class="material-symbols-outlined">local_shipping</span></div>
          <div class="tag-content">
            <label>Albarán Emitido</label>
            <strong>{{ note.deliveryNoteNumber }}</strong>
          </div>
          <span class="material-symbols-outlined jump-icon">open_in_new</span>
        </router-link>
        
        <router-link v-for="invoice in invoices" :key="invoice.id" :to="`/sales/invoices/${invoice.id}`" class="related-tag-card">
          <div class="tag-icon success"><span class="material-symbols-outlined">receipt</span></div>
          <div class="tag-content">
            <label>Factura de Venta</label>
            <strong>{{ invoice.invoiceNumber }}</strong>
          </div>
          <span class="material-symbols-outlined jump-icon">open_in_new</span>
        </router-link>
      </div>
    </template>

    <!-- 5. MAIN CONTENT -->
    <FormSection title="Identificación del Cliente" icon="person">
      <div v-if="mode === 'detail'">
        <DataRow label="Nombre del Cliente" :value="partyName" icon="person" />
        <DataRow v-if="order?.taxId" label="NIF/CIF" :value="order.taxId" is-mono />
      </div>
      <div v-else>
        <PartySelector
          v-model="formData.partyId"
          label="Seleccionar Cliente *"
          placeholder="Buscar por nombre o referencia..."
          role-filter="CLIENT"
          :required="true"
          @select="onPartySelected"
        />
      </div>
    </FormSection>

    <FormSection title="Logística y Fechas" icon="calendar_today">
      <div v-if="mode === 'detail'">
        <DataRow label="Fecha de Pedido" :value="formatDate(order?.orderDate)" icon="calendar_today" />
        <DataRow label="Entrega Estimada" :value="formatDate(order?.deliveryDate)" icon="local_shipping" />
        <DataRow label="Notas e Instrucciones" icon="notes">
          <p class="notes-text">{{ order?.notes || 'Sin observaciones.' }}</p>
        </DataRow>
      </div>
      <div v-else>
        <div class="form-row">
          <div class="form-group">
            <label>Fecha de Pedido *</label>
            <input v-model="formData.orderDate" type="date" class="form-input" required />
          </div>
          <div class="form-group">
            <label>Entrega Estimada *</label>
            <input v-model="formData.deliveryDate" type="date" class="form-input" required />
          </div>
        </div>
        <div class="form-group mt-4">
          <label>Observaciones del Pedido</label>
          <textarea v-model="formData.notes" class="form-textarea" rows="3" placeholder="Instrucciones especiales..."></textarea>
        </div>
      </div>
    </FormSection>

    <FormSection title="Especificaciones Técnicas (MES)" icon="precision_manufacturing">
      <div v-if="mode === 'detail'" class="table-wrapper">
        <table v-if="order?.mesWorkRefs?.length > 0" class="data-table">
          <thead>
            <tr>
              <th>Proceso / Configuración</th>
              <th>Configuración de Trabajo</th>
              <th>Notas del Pedido</th>
              <th class="text-right">Navegación MES</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="mesRef in order.mesWorkRefs" :key="mesRef.workSetupId || mesRef.id">
              <td class="w-64">
                <div class="mes-config-info">
                  <span class="material-symbols-outlined icon-secondary">settings_suggest</span>
                  <strong>{{ formatMesWorkId(mesRef.workSetupId) }}</strong>
                </div>
              </td>
              <td>
                <div class="mes-specs-inline">
                  <div v-for="line in mesWorksCache[mesRef.workSetupId]?.lines" :key="line.id" class="spec-pill">
                    <span class="label">{{ getWorkTypeName(line.work_type_id) }}:</span>
                    <span class="value">{{ getPositionName(line.position_id) }}</span>
                  </div>
                  <span v-if="!mesWorksCache[mesRef.workSetupId]?.lines?.length" class="text-muted">Sin especificaciones</span>
                </div>
              </td>
              <td><p class="notes-text-sm">{{ mesRef.description || '—' }}</p></td>
              <td class="text-right">
                <router-link v-if="mesRef.workOrderId" :to="`/mes/work-orders/${mesRef.workOrderId}`" class="btn btn-nav btn-sm">
                  <span class="material-symbols-outlined">analytics</span>
                  <span>Ver Orden MES</span>
                </router-link>
                <span v-else class="status-badge status-pending">
                  <span class="material-symbols-outlined" style="font-size:14px">schedule</span>
                  Pendiente
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        <p v-else class="text-muted p-4">No se han definido requerimientos técnicos.</p>
      </div>

      <!-- MODO EDICIÓN -->
      <div v-else>
        <div class="mb-4">
          <button type="button" class="btn btn-primary btn-sm" @click="addMesWorkRef">
            <span class="material-symbols-outlined">add</span> <span>Añadir Trabajo MES</span>
          </button>
        </div>
        <div class="table-wrapper">
          <table class="data-table fixed-layout">
            <thead>
              <tr>
                <th style="width: 50px">#</th>
                <th style="width: 250px">Configuración</th>
                <th>Instrucciones Técnicas / Notas del Pedido *</th>
                <th class="text-center" style="width: 80px">Borrar</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(ref, idx) in formData.mesWorkRefs" :key="idx">
                <td class="text-muted">{{ idx + 1 }}</td>
                <td>
                  <select v-model="ref.workSetupId" class="form-input-sm w-full">
                    <option :value="null">-- Personalizado --</option>
                    <option v-for="setup in availableMesSetups" :key="setup.id" :value="setup.id">{{ setup.name }}</option>
                  </select>
                </td>
                <td class="w-full">
                  <input v-model="ref.description" type="text" class="form-input-sm w-full" placeholder="Escriba aquí las especificaciones detalladas..." required />
                </td>
                <td class="text-center">
                  <button type="button" class="btn-icon text-danger" @click="removeMesWorkRef(idx)">
                    <span class="material-symbols-outlined">delete</span>
                  </button>
                </td>
              </tr>
              <tr v-if="formData.mesWorkRefs.length === 0">
                <td colspan="4" class="text-muted text-center p-4">No hay trabajos MES vinculados.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </FormSection>

    <FormSection title="Líneas del Pedido" icon="list_alt">
      <div v-if="mode !== 'detail'" class="mb-4">
        <button type="button" class="btn btn-primary btn-sm" @click="openVariantSelector" ref="addBtn">
          <span class="material-symbols-outlined">add</span> <span>Añadir Producto</span>
        </button>
      </div>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>Ref.</th>
              <th>Producto</th>
              <th class="text-center">Cant.</th>
              <th class="align-right">P. Tarifa</th>
              <th class="align-right">P. Venta</th>
              <th class="text-center">Dto %</th>
              <th class="align-right">Subtotal</th>
              <th v-if="mode !== 'detail'" class="text-center">Borrar</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, idx) in (mode === 'detail' ? order.lineItems : formData.lineItems)" :key="idx">
              <td><code class="code-badge">{{ item.variantSku || formatVariantId(item.productVariantId) }}</code></td>
              <td>{{ buildDisplayName(item) }}</td>
              <td class="text-center">
                <template v-if="mode === 'detail'">{{ item.quantity }}</template>
                <input v-else v-model.number="item.quantity" type="number" min="1" class="form-input-sm w-16" />
              </td>
              <td class="align-right">{{ salesApi.formatMoney(mode === 'detail' ? item.listUnitPrice : { amount: item.listPrice || item.unitPrice, currency: 'EUR' }) }}</td>
              <td class="align-right">
                <template v-if="mode === 'detail'">{{ salesApi.formatMoney(item.unitPrice) }}</template>
                <input v-else v-model.number="item.unitPrice" type="number" step="0.01" class="form-input-sm w-24 text-right" />
              </td>
              <td class="text-center">
                <template v-if="mode === 'detail'">{{ item.discountPercent ? item.discountPercent.toFixed(2) + '%' : '—' }}</template>
                <input v-else v-model.number="item.discountPercent" type="number" step="0.01" class="form-input-sm w-16 text-center" @keyup.enter="focusAddBtn" />
              </td>
              <td class="align-right">
                <strong>{{ salesApi.formatMoney(mode === 'detail' ? item.subtotal : calculateLineSubtotal(idx)) }}</strong>
              </td>
              <td v-if="mode !== 'detail'" class="text-center">
                <button type="button" class="btn-icon text-danger" @click="removeLineItem(idx)"><span class="material-symbols-outlined">delete</span></button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </FormSection>

    <FormSection title="Resumen Económico" icon="payments">
      <div v-if="mode === 'detail'">
        <DataRow label="Subtotal" :value="salesApi.formatMoney(order.subtotal)" />
        <DataRow label="Impuestos" :value="salesApi.formatMoney(order.taxAmount)" />
        <DataRow label="TOTAL PEDIDO" :value="salesApi.formatMoney(order.total)" highlight />
      </div>
      <div v-else class="totals-checkout-layout">
        <section class="totals-checkout-card">
          <div class="total-row"><label>Subtotal:</label><span>{{ salesApi.formatMoney(liveTotals.subtotal) }}</span></div>
          <div class="total-row"><label>IVA (21%):</label><span>{{ salesApi.formatMoney(liveTotals.taxAmount) }}</span></div>
          <div class="total-row final"><label>TOTAL ESTIMADO:</label><span class="total-value">{{ salesApi.formatMoney(liveTotals.total) }}</span></div>
        </section>
      </div>
    </FormSection>

    <template #footer v-if="mode === 'detail' && order">
      <div class="audit-info">
        <p>Documento generado por el sistema TramaTex.</p>
        <p v-if="order.id">ID único: <code>{{ order.id }}</code></p>
      </div>
    </template>
  </BaseEntityPage>

  <!-- CAPA DE IMPRESIÓN -->
  <div class="print-container">
    <PrintDocument
      v-if="order"
      type="ORDER"
      :number="order.orderNumber"
      :date="order.orderDate"
      :customer-name="partyName"
      :customer-tax-id="order.taxId"
      :items="order.lineItems"
      :totals="{ subtotal: order.subtotal, taxAmount: order.taxAmount, total: order.total }"
      :notes="order.notes"
    />
  </div>

  <!-- MODALES MAESTROS -->
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
      <p>Está a punto de <strong>generar una factura oficial</strong> para este pedido.</p>
      <div class="info-notice mt-4">
        <span class="material-symbols-outlined">info</span>
        <p>Esta operación consolidará los importes y creará un nuevo documento contable vinculado permanentemente a este pedido.</p>
      </div>
      <p class="mt-4 text-secondary italic">¿Desea continuar?</p>
    </div>
  </BaseDialog>

  <!-- DIÁLOGO DE REAPERTURA / MODIFICACIÓN -->
  <BaseDialog
    :show="showReopenConfirm"
    title="Modificar Pedido Confirmado"
    icon="settings_backup_restore"
    confirm-text="Devolver a Pendiente"
    confirm-class="btn-primary"
    :is-confirming="isChangingStatus"
    @close="showReopenConfirm = false"
    @confirm="confirmReopen"
  >
    <div class="confirm-dialog-body">
      <p>Para realizar cambios en un pedido confirmado, primero debe <strong>devolverlo a estado PENDIENTE</strong>.</p>
      <div class="info-notice mt-4">
        <span class="material-symbols-outlined">warning</span>
        <p>Esta acción notificará al sistema que el pedido vuelve a fase administrativa. Deberá volver a confirmarlo tras realizar los cambios.</p>
      </div>
      <p class="mt-4 text-secondary italic">¿Desea proceder con la reapertura?</p>
    </div>
  </BaseDialog>

  <!-- DIÁLOGO DE ANULACIÓN -->
  <BaseDialog
    :show="showCancelConfirm"
    title="Anular Pedido"
    icon="cancel"
    confirm-text="Anular Definitivamente"
    confirm-class="btn-danger"
    :is-confirming="isChangingStatus"
    @close="showCancelConfirm = false"
    @confirm="confirmCancel"
  >
    <div class="confirm-dialog-body">
      <p>¿Está seguro de que desea <strong>anular este pedido</strong>?</p>
      <p class="mt-4 text-muted">Esta acción es definitiva y el pedido dejará de ser válido para facturación o entrega.</p>
    </div>
  </BaseDialog>

  <!-- MODAL: CREAR ALBARÁN -->
  <div v-if="showDeliveryNoteModal" class="modal-backdrop no-print">
    <div class="modal card">
      <div class="modal-header">
        <span class="material-symbols-outlined icon-secondary">local_shipping</span>
        <h2>Generar Albarán</h2>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Fecha de Entrega Real</label>
          <input v-model="deliveryNoteForm.deliveryDate" type="date" class="form-input" />
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn-outline" @click="showDeliveryNoteModal = false">Cancelar</button>
        <button class="btn btn-success" @click="createDeliveryNote" :disabled="isCreatingDeliveryNote">Confirmar Albarán</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue';
import { useRoute, useRouter, RouterLink } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import PrintDocument from '@/components/sales/PrintDocument.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { mesApi } from '@/services/mesApi';

const route = useRoute();
const router = useRouter();

const mode = ref('detail');
const isLoading = ref(false);
const isSaving = ref(false);
const error = ref('');

const order = ref(null);
const partyName = ref('');
const partyDefaultDiscount = ref(0);
const deliveryNotes = ref([]);
const invoices = ref([]);
const mesWorksCache = ref({});
const availableMesSetups = ref([]);
const isLoadingMesSetups = ref(false);

const workTypesCache = ref({});
const positionsCache = ref({});

const formData = reactive({
  partyId: '', orderDate: '', deliveryDate: '', notes: '', mesWorkRefs: [], lineItems: []
});

const addBtn = ref(null);
const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewTimer = null;

const showVariantSelector = ref(false);
const showDeliveryNoteModal = ref(false);
const showInvoiceConfirm = ref(false);
const showReopenConfirm = ref(false);
const showCancelConfirm = ref(false);

const isCreatingDeliveryNote = ref(false);
const isCreatingInvoice = ref(false);
const isChangingStatus = ref(false);

const deliveryNoteForm = ref({ deliveryDate: new Date().toISOString().split('T')[0] });

// Reactividad de líneas
const lineItemsFingerprint = computed(() =>
  mode.value !== 'detail'
    ? formData.lineItems.map(i => `${i.productVariantId}|${i.quantity}|${i.unitPrice}|${i.discountPercent}`).join('§')
    : ''
);
watch(lineItemsFingerprint, (val) => { if (val && !isSaving.value) calculateTotals(); });

onMounted(async () => {
  if (route.name === 'CreateOrder' || !route.params.id) { 
    mode.value = 'create';
    partyName.value = 'Seleccione un cliente';
    const today = new Date();
    formData.orderDate = today.toISOString().split('T')[0];
    const delivery = new Date();
    delivery.setDate(today.getDate() + 15);
    formData.deliveryDate = delivery.toISOString().split('T')[0];
  }
  else { 
    await fetchOrder();
    await loadMesMasters();
  }
});

async function fetchOrder() {
  const id = route.params.id;
  if (!id) return;
  isLoading.value = true;
  try {
    order.value = await salesApi.getOrder(id);
    const party = await partyApi.getParty(order.value.partyId);
    partyName.value = party.name;
    partyDefaultDiscount.value = party.default_discount_percentage || 0;
    await Promise.all([loadDeliveryNotes(), loadInvoices(), loadMesWorksNames()]);
  } catch (err) { error.value = err.message; }
  finally { isLoading.value = false; }
}

async function loadMesMasters() {
  try {
    const [types, pos] = await Promise.all([mesApi.listWorkTypes(), mesApi.listPositions()]);
    types.forEach(t => workTypesCache.value[t.id] = t.name);
    pos.forEach(p => positionsCache.value[p.id] = p.name);
  } catch (err) { console.error('Error maestros MES', err); }
}

function getWorkTypeName(id) { return workTypesCache.value[id] || 'Cargando...'; }
function getPositionName(id) { return positionsCache.value[id] || 'Cargando...'; }

async function loadMesWorksNames() {
  const refs = order.value?.mesWorkRefs || [];
  if (!refs.length) return;
  const ids = refs.map(r => r.workSetupId).filter(Boolean);
  const results = await Promise.allSettled(ids.map(id => mesApi.getWorkSetup(id)));
  results.forEach((r, i) => { if (r.status === 'fulfilled') mesWorksCache.value[ids[i]] = r.value; });
}

async function loadDeliveryNotes() { try { const res = await salesApi.listDeliveryNotes({ orderId: order.value.id }); deliveryNotes.value = res.data || []; } catch (err) {} }
async function loadInvoices() { try { const res = await salesApi.listInvoices({ orderId: order.value.id }); invoices.value = res.data || []; } catch (err) {} }

function enterEditMode() {
  if (!order.value) return;
  const data = {
    partyId: order.value.partyId,
    orderDate: order.value.orderDate?.split('T')[0] || '',
    deliveryDate: order.value.deliveryDate?.split('T')[0] || '',
    notes: order.value.notes || '',
    mesWorkRefs: (order.value.mesWorkRefs || []).map(r => ({ id: r.id, workSetupId: r.workSetupId || null, workOrderId: r.workOrderId || null, description: r.description || '' })),
    lineItems: (order.value.lineItems || []).map(i => ({ id: i.id, productVariantId: i.productVariantId, variantSku: i.variantSku, displayName: buildDisplayName(i), quantity: i.quantity, unitPrice: i.listUnitPrice?.amount ?? i.unitPrice?.amount ?? 0, discountPercent: i.discountPercent || 0 }))
  };
  Object.assign(formData, data);
  mode.value = 'edit';
  loadAvailableSetups(formData.partyId);
  calculateTotals();
}

function exitEditMode() { if (mode.value === 'edit') mode.value = 'detail'; else router.push('/sales/orders'); }

async function loadAvailableSetups(partyId) {
  if (!partyId) return;
  isLoadingMesSetups.value = true;
  try { availableMesSetups.value = await mesApi.listWorkSetups({ party_id: partyId }); }
  catch (err) { availableMesSetups.value = []; } finally { isLoadingMesSetups.value = false; }
}

function onPartySelected(party) {
  partyName.value = party?.name || 'Cliente no seleccionado';
  partyDefaultDiscount.value = party?.default_discount_percentage || 0;
  formData.partyId = party?.id || '';
  loadAvailableSetups(party?.id);
  calculateTotals();
}

function openVariantSelector() { showVariantSelector.value = true; }

function handleVariantSelected(payload) {
  const variantId = payload.variantId || payload.id;
  const variant = payload.variant || payload;
  formData.lineItems.push({ id: null, productVariantId: variantId, variantSku: variant.sku || '', displayName: buildDisplayName(variant), quantity: 1, unitPrice: variant.product_base_price ?? variant.price ?? 0, discountPercent: partyDefaultDiscount.value || 0 });
  showVariantSelector.value = false;
  nextTick(() => { if (addBtn.value) addBtn.value.focus(); });
}

function focusAddBtn() { if (addBtn.value) addBtn.value.focus(); }
function removeLineItem(idx) { formData.lineItems.splice(idx, 1); }
function addMesWorkRef() { formData.mesWorkRefs.push({ workSetupId: null, description: '' }); }
function removeMesWorkRef(idx) { formData.mesWorkRefs.splice(idx, 1); }

function calculateTotals() { clearTimeout(previewTimer); previewTimer = setTimeout(fetchPreviewCalculation, 400); }

async function fetchPreviewCalculation() {
  const partyId = formData.partyId;
  const items = formData.lineItems.filter(i => i.productVariantId).map(i => ({ productVariantId: i.productVariantId, quantity: i.quantity, unitPrice: { amount: i.unitPrice, currency: 'EUR' }, discountPercent: i.discountPercent }));
  if (!partyId || items.length === 0) { previewResult.value = null; return; }
  isPreviewLoading.value = true;
  try { previewResult.value = await salesApi.previewOrderCalculation(partyId, items); } catch (err) {}
  finally { isPreviewLoading.value = false; }
}

const liveTotals = computed(() => {
  if (mode.value === 'detail' && order.value) return { subtotal: order.value.subtotal, taxAmount: order.value.taxAmount, total: order.value.total };
  if (previewResult.value) return { subtotal: previewResult.value.subtotal, taxAmount: previewResult.value.taxAmount, total: previewResult.value.total };
  let subtotal = 0;
  formData.lineItems.forEach((_, idx) => { subtotal += calculateLineSubtotal(idx); });
  const tax = subtotal * 0.21;
  return { subtotal: { amount: subtotal, currency: 'EUR' }, taxAmount: { amount: tax, currency: 'EUR' }, total: { amount: subtotal + tax, currency: 'EUR' } };
});

function calculateLineSubtotal(idx) {
  const item = formData.lineItems[idx];
  if (!item) return 0;
  if (previewResult.value?.lineItems && item.productVariantId) {
    const calculated = previewResult.value.lineItems.find(li => li.productVariantId === item.productVariantId);
    if (calculated?.subtotal?.amount !== undefined) return calculated.subtotal.amount;
  }
  return (Number(item.quantity) || 0) * (Number(item.unitPrice) || 0) * (1 - (Number(item.discountPercent) || 0) / 100);
}

async function saveOrder() {
  if (formData.lineItems.length === 0) { alert('Añada productos'); return; }
  isSaving.value = true;
  try {
    const deliveryDateISO = new Date(formData.deliveryDate).toISOString();
    const mesRefs = formData.mesWorkRefs.map(r => ({ workSetupId: r.workSetupId || undefined, description: r.description || '' }));
    const lineItems = formData.lineItems.map(item => ({ productVariantId: item.productVariantId, quantity: item.quantity, unitPrice: { amount: item.unitPrice, currency: 'EUR' }, discountPercent: item.discountPercent }));

    if (mode.value === 'create') {
      const newOrder = await salesApi.createOrder({ partyId: formData.partyId, deliveryDate: deliveryDateISO, notes: formData.notes || '', mesWorkRefs: mesRefs, items: lineItems });
      await router.replace(`/sales/orders/${newOrder.id}`);
      await fetchOrder();
      mode.value = 'detail';
    } else {
      await salesApi.updateOrder(order.value.id, { deliveryDate: deliveryDateISO, notes: formData.notes || '', mesWorkRefs: mesRefs });
      await fetchOrder();
      mode.value = 'detail';
    }
  } catch (err) { alert('Error al guardar: ' + err.message); }
  finally { isSaving.value = false; }
}

function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '—'; }
function buildDisplayName(i) { return (i.productName || i.displayName || 'Producto') + (i.optionConfiguration ? ' - ' + Object.values(i.optionConfiguration).join(', ') : ''); }
function formatVariantId(id) { return id ? id.substring(0, 8) : '—'; }
function formatMesWorkId(id) { return id ? (mesWorksCache.value[id]?.name || id.substring(0, 8)) : 'Sin config.'; }

async function confirmOrder() { try { await salesApi.changeOrderStatus(order.value.id, 'CONFIRMED'); await fetchOrder(); } catch (err) { alert('Error al confirmar: ' + err.message); } }

function triggerCancel() { showCancelConfirm.value = true; }
async function confirmCancel() {
  isChangingStatus.value = true;
  try { await salesApi.changeOrderStatus(order.value.id, 'CANCELLED'); await fetchOrder(); showCancelConfirm.value = false; }
  catch (err) { alert('Error al anular: ' + err.message); }
  finally { isChangingStatus.value = false; }
}

function triggerReopen() { showReopenConfirm.value = true; }
async function confirmReopen() {
  isChangingStatus.value = true;
  try {
    // Si el backend no acepta PENDING, el servicio gestionará el error
    await salesApi.changeOrderStatus(order.value.id, 'PENDING');
    await fetchOrder();
    showReopenConfirm.value = false;
  } catch (err) {
    alert('No se ha podido reabrir el pedido. Verifique que no tenga albaranes o facturas vinculadas.');
  } finally {
    isChangingStatus.value = false;
  }
}

function printOrder() { window.print(); }

async function createDeliveryNote() {
  isCreatingDeliveryNote.value = true;
  try {
    const items = order.value.lineItems.map(i => ({ salesOrderLineItemId: i.id, deliveredQuantity: i.quantity }));
    await salesApi.createDeliveryNote({ salesOrderId: order.value.id, deliveryDate: new Date(deliveryNoteForm.value.deliveryDate).toISOString(), items });
    await fetchOrder(); showDeliveryNoteModal.value = false;
  } catch (err) { alert(err.message); }
  finally { isCreatingDeliveryNote.value = false; }
}

function createInvoiceFromOrder() { showInvoiceConfirm.value = true; }
async function confirmCreateInvoice() {
  isCreatingInvoice.value = true;
  try {
    const now = new Date();
    const inv = await salesApi.createInvoice({ partyId: order.value.partyId, salesOrderIds: [order.value.id], invoiceDate: now.toISOString(), dueDate: new Date(now.getTime() + 30*24*60*60*1000).toISOString() });
    showInvoiceConfirm.value = false; router.push(`/sales/invoices/${inv.id}`);
  } catch (err) { alert(err?.message); } finally { isCreatingInvoice.value = false; }
}

const canEdit = computed(() => order.value && ['PENDING', 'PENDIENTE'].includes(order.value.status));
const canCreateDeliveryNote = computed(() => order.value && ['CONFIRMED', 'CONFIRMADO', 'EN_PREPARACION', 'PARTIALLY_DELIVERED', 'ENTREGADO_PARCIALMENTE'].includes(order.value.status));
const canCreateInvoice = computed(() => order.value && ['DELIVERED', 'ENTREGADO', 'PARTIALLY_INVOICED', 'FACTURADO_PARCIALMENTE'].includes(order.value.status));
const canCancel = computed(() => order.value && ['PENDING', 'PENDIENTE', 'CONFIRMED', 'CONFIRMADO', 'EN_PREPARACION'].includes(order.value.status));
const canReopen = computed(() => order.value && ['CANCELLED', 'CANCELADO', 'CONFIRMED', 'CONFIRMADO', 'EN_PREPARACION'].includes(order.value.status));
</script>

<style scoped>
@import "@/design-system/_sections.css";
.overview-tags-row, .related-history-grid { display: flex; flex-wrap: wrap; gap: 1rem; }
.related-history-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); }
.summary-tag { flex: 1; min-width: 240px; padding: 0.6rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 0.75rem; box-shadow: var(--box-shadow-sm); }
.related-tag-card { padding: 0.6rem 1rem; background: var(--color-background); border: 1px solid var(--color-border); border-left: 4px solid var(--color-secondary); border-radius: 10px; display: flex; align-items: center; gap: 0.75rem; text-decoration: none; position: relative; transition: all 0.2s ease; }
.related-tag-card.highlight-info { border-left-color: #2563eb; }
.related-tag-card.success { border-left-color: #16a34a; }
.related-tag-card:hover { background: white; transform: translateX(2px) translateY(-1px); box-shadow: var(--box-shadow-md); }
.related-tag-card:hover strong { color: var(--color-primary); text-decoration: underline; }
.tag-icon { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; background: white; }
.summary-tag .tag-icon { width: 40px; height: 40px; border-radius: 10px; }
.related-tag-card .tag-icon { background: rgba(0,0,0,0.03); color: var(--color-text-secondary); }
.icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.icon.green, .tag-icon.success { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.tag-content { display: flex; flex-direction: column; gap: 0.15rem; line-height: 1.2; }
.tag-content label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.025em; }
.tag-content strong { font-size: 0.95rem; color: var(--color-text-primary); }
.amount { color: #16a34a !important; font-size: 1.15rem !important; }
.action-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; box-shadow: var(--box-shadow-sm); margin: 0; }
.status-badge { padding: 0.4rem 1rem; font-size: 0.85rem; font-weight: 800; letter-spacing: 0.05em; }
.toolbar-buttons { display: flex; gap: 0.75rem; }
.mes-config-info { display: flex; align-items: center; gap: 0.75rem; }
.mes-specs-inline { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.spec-pill { display: inline-flex; align-items: center; gap: 0.35rem; padding: 0.2rem 0.5rem; background: var(--color-background); border-radius: 4px; border: 1px solid var(--color-border); font-size: 0.8rem; }
.spec-pill .label { color: var(--color-text-secondary); font-weight: 600; }
.spec-pill .value { color: var(--color-text-primary); font-weight: 700; }
.notes-text-sm { font-size: 0.85rem; font-style: italic; color: var(--color-text-secondary); margin: 0; line-height: 1.4; }
.w-64 { width: 16rem; } .w-full { width: 100%; } .fixed-layout { table-layout: fixed; width: 100%; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-input-sm { padding: 0.5rem; border: 1px solid var(--color-border); border-radius: 4px; font-size: 0.85rem; }
.totals-checkout-layout { display: flex; justify-content: flex-end; margin-top: 1rem; }
.totals-checkout-card { width: 400px; padding: 1.5rem; background: white; border: 1px solid var(--color-border-strong); border-radius: 12px; box-shadow: var(--box-shadow-md); }
.total-row { display: flex; justify-content: space-between; margin-bottom: 0.75rem; font-size: 0.95rem; }
.total-row.final { margin-top: 1rem; padding-top: 1rem; border-top: 2px solid var(--color-border); font-weight: 800; font-size: 1.25rem; }
.total-value { color: var(--color-secondary); }
.audit-info { color: var(--color-text-secondary); font-size: 0.8rem; font-style: italic; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; }
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.w-modal-xl { width: 90%; max-width: 1100px; }
.info-notice { display: flex; gap: 0.75rem; padding: 1rem; background: rgba(59, 130, 246, 0.1); border-radius: 8px; color: #1e40af; font-size: 0.9rem; }
.print-container { display: none; }
@media print {
  .no-print, . identity-header-sticky, .navbar { display: none !important; }
  .print-container { display: block !important; position: absolute; left: 0; top: 0; width: 100%; }
}
</style>
