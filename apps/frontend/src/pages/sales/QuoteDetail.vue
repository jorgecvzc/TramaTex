<template>
  <BaseEntityPage class="no-print" v-if="isLoading">
    <template #header>
      <BasePageHeader title="Cargando..." :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos' }]" show-back />
    </template>
    <div class="loading-state card">
      <div class="spinner"></div>
      <p>Cargando información del presupuesto...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage class="no-print" v-else-if="error">
    <template #header>
      <BasePageHeader title="Error" :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-icon-wrapper error">
        <AlertCircle :size="24" />
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/sales/quotes')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage class="no-print" v-else-if="quote || mode === 'create'">
    <!-- 1. IDENTITY HEADER -->
    <template #header>
      <BasePageHeader 
        :title="mode === 'create' ? 'Nuevo Presupuesto' : (mode === 'edit' ? `Editando Presupuesto ${quote?.quoteNumber}` : `Presupuesto ${quote?.quoteNumber}`)" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos', to: '/sales/quotes' }, { label: mode === 'create' ? 'Crear' : quote?.quoteNumber }]"
        show-back
      >
        <template #icon>
          <FileText :size="28" />
        </template>
        <template #actions>
          <div v-if="quote || mode === 'create'" class="header-actions-group">
            <template v-if="mode === 'detail'">
              <button v-if="quote?.status !== 'BORRADOR'" class="btn btn-outline" @click="printQuote">
                <Printer :size="16" /> <span>Imprimir</span>
              </button>
              <button v-if="canEdit" class="btn btn-primary" @click="enterEditMode">
                <Pencil :size="16" /> <span>Editar Presupuesto</span>
              </button>
            </template>
            <template v-else>
              <button class="btn btn-outline" @click="exitEditMode" :disabled="isSaving">Cancelar</button>
              <button class="btn btn-secondary" @click="saveQuote" :disabled="isSaving">
                <component :is="isSaving ? RefreshCw : Save" :size="16" :class="{ 'spin': isSaving }" />
                <span>{{ isSaving ? 'Guardando...' : 'Guardar Presupuesto' }}</span>
              </button>
            </template>
          </div>
        </template>
      </BasePageHeader>
    </template>

    <!-- 2. TOOLBAR -->
    <template #toolbar v-if="mode === 'detail' && quote">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', `status-${getStatusClass(quote.status)}`]">
            {{ getStatusLabel(quote.status) }}
          </span>
        </div>
        <div class="toolbar-buttons">
          <button v-if="['BORRADOR', 'DRAFT'].includes(quote.status?.toUpperCase())" class="btn btn-success btn-sm" @click="promptIssueQuote">
            <Send :size="16" /> <span>Emitir a Cliente</span>
          </button>
          <!-- Robust check for conversion button -->
          <button 
            v-if="['ISSUED', 'EMITIDA', 'EMITIDO', 'APPROVED', 'APROBADA', 'APROBADO', 'ACCEPTED'].includes(quote.status?.toUpperCase()) && !quote.generatedOrderId" 
            class="btn btn-success btn-sm" 
            @click="showConvertModal = true"
          >
            <ShoppingCart :size="16" /> <span>Convertir a Pedido</span>
          </button>
          <button v-if="['EMITIDA', 'EMITIDO', 'ISSUED'].includes(quote.status?.toUpperCase())" class="btn btn-danger btn-sm" @click="promptRejectQuote">
            <XCircle :size="16" /> <span>Rechazar</span>
          </button>
          <button v-if="['RECHAZADA', 'RECHAZADO', 'REJECTED'].includes(quote.status?.toUpperCase())" class="btn btn-primary btn-sm" @click="promptReactivateQuote">
            <RefreshCw :size="16" /> <span>Reactivar</span>
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
          <div class="tag-content"><label>Fecha Emisión</label><strong>{{ formatDate(mode === 'create' ? new Date() : quote?.quoteDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><CalendarOff :size="20" /></div>
          <div class="tag-content"><label>Válido Hasta</label><strong :class="{'text-danger': isExpired}">{{ formatDate(mode === 'create' ? formData.expirationDate : quote?.expirationDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><CreditCard :size="20" /></div>
          <div class="tag-content">
            <label>Total Presupuesto</label>
            <strong class="amount">{{ salesApi.formatMoney(liveTotals.total) }}</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- 4. RELATED -->
    <template #related v-if="mode === 'detail' && quote?.generatedOrderId">
      <div class="related-history-grid">
        <router-link :to="`/sales/orders/${quote.generatedOrderId}`" class="related-tag-card highlight-info">
          <div class="tag-icon"><ShoppingCart :size="20" /></div>
          <div class="tag-content">
            <label>Pedido Generado</label>
            <strong>{{ quote.generatedOrderNumber || 'Ver Pedido' }}</strong>
          </div>
          <ExternalLink :size="14" class="jump-icon" />
        </router-link>
      </div>
    </template>

    <!-- 5. MAIN CONTENT -->
    <FormSection title="Identificación del Cliente" icon="person">
      <div v-if="(mode === 'create' || mode === 'edit') && !isSaving">
        <PartySelector
          v-model="formData.partyId"
          label="Seleccionar Cliente *"
          placeholder="Buscar por nombre o referencia..."
          role-filter="CLIENT"
          :required="true"
          @select="onPartySelected"
        />
      </div>
      <div v-else>
        <DataRow label="Nombre del Cliente" :value="partyName" icon="person" />
        <DataRow v-if="quote?.taxId" label="NIF/CIF" :value="quote.taxId" is-mono />
      </div>
    </FormSection>

    <FormSection title="Información del Presupuesto" icon="info">
      <div v-if="mode === 'detail'">
        <DataRow label="Fecha de Emisión" :value="formatDate(quote?.quoteDate)" icon="calendar_today" />
        <DataRow label="Válido Hasta" :value="formatDate(quote?.expirationDate)" icon="event_busy" :class="{'text-danger fw-bold': isExpired}" />
        <DataRow label="Observaciones Comerciales" icon="notes">
          <p class="notes-text">{{ quote?.notes || 'Sin observaciones.' }}</p>
        </DataRow>
      </div>
      <div v-else>
        <div class="form-row">
          <div class="form-group">
            <label>Válido Hasta *</label>
            <input v-model="formData.expirationDate" type="date" class="form-input" required />
          </div>
        </div>
        <div class="form-group mt-4">
          <label>Observaciones del Presupuesto</label>
          <textarea v-model="formData.notes" class="form-textarea" rows="3" placeholder="Plazos, condiciones, instrucciones..."></textarea>
        </div>
      </div>
    </FormSection>

    <FormSection title="Configuración Técnica (MES)" icon="precision_manufacturing">
      <div v-if="mode === 'detail'" class="table-wrapper">
        <table v-if="quote?.mesWorkRefs?.length > 0" class="data-table">
          <thead>
            <tr>
              <th>Proceso / Configuración</th>
              <th>Configuración de Trabajo</th>
              <th>Descripción</th>
              <th class="text-right">Estado MES</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="mesRef in quote.mesWorkRefs" :key="mesRef.id">
              <td class="w-64">
                <div class="mes-config-info">
                  <Settings2 :size="18" class="icon-secondary" />
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
                <span class="status-badge status-info">
                  {{ mesRef.workSetupId ? 'Configurado' : 'Pendiente' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        <p v-else class="text-muted p-4">No se han definido requerimientos técnicos.</p>
      </div>

      <!-- MODO EDICIÓN / CREACIÓN -->
      <div v-else>
        <div class="mb-4">
          <button type="button" class="btn btn-primary btn-sm" @click="addMesWorkRef">
            <Plus :size="16" /> <span>Añadir Trabajo MES</span>
          </button>
        </div>
        <div class="table-wrapper">
          <table class="data-table fixed-layout">
            <thead>
              <tr>
                <th style="width: 50px">#</th>
                <th style="width: 250px">Configuración</th>
                <th>Descripción / Notas *</th>
                <th class="text-center" style="width: 80px">Borrar</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(ref, idx) in formData.mesWorkRefs" :key="idx">
                <td class="text-muted">{{ idx + 1 }}</td>
                <td>
                  <select 
                    v-model="ref.workSetupId" 
                    class="form-input-sm w-full"
                    :data-mes-row="idx"
                    data-mes-col="setup"
                    @keydown="handleMesKeyDown($event, idx, 'setup', ref)"
                  >
                    <option :value="null">-- Personalizado --</option>
                    <option v-for="setup in availableMesSetups" :key="setup.id" :value="setup.id">{{ setup.name }}</option>
                  </select>
                </td>
                <td class="w-full">
                  <input 
                    v-model="ref.description" 
                    type="text" 
                    class="form-input-sm w-full" 
                    placeholder="Especificaciones técnicas..." 
                    required 
                    :data-mes-row="idx"
                    data-mes-col="desc"
                    @keydown="handleMesKeyDown($event, idx, 'desc', ref)"
                  />
                </td>
                <td class="text-center">
                  <button type="button" class="btn-icon text-danger" @click="removeMesWorkRef(idx)">
                    <Trash2 :size="18" />
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

    <FormSection title="Líneas del Presupuesto" icon="list_alt">
      <OrderLines
        :lines="mode === 'detail' ? quote.lineItems : formData.lineItems"
        :is-editing="mode !== 'detail'"
        @update:lines="(newLines) => formData.lineItems = newLines"
        @add-line="openVariantSelector"
        @last-field-tab="addProductBtn?.focus()"
      />
      <div v-if="mode !== 'detail'" class="mt-4">
        <button 
          ref="addProductBtn"
          type="button" 
          class="btn btn-primary btn-sm" 
          @click="openVariantSelector"
        >
          <Plus :size="16" /> <span>Añadir Producto (Ins)</span>
        </button>
      </div>
    </FormSection>

    <FormSection title="Resumen Económico" icon="payments">
      <div class="totals-summary-container">
        <section class="totals-summary-card" :class="{ 'is-loading-overlay': isPreviewLoading }">
          <div class="summary-row">
            <label>Subtotal:</label>
            <span>{{ salesApi.formatMoney(liveTotals.subtotal) }}</span>
          </div>
          <div class="summary-row">
            <label>Impuestos (IVA 21%):</label>
            <span>{{ salesApi.formatMoney(liveTotals.taxAmount) }}</span>
          </div>
          <div class="summary-row grand-total">
            <label>{{ mode === 'detail' ? 'TOTAL PRESUPUESTO:' : 'TOTAL ESTIMADO:' }}</label>
            <span>{{ salesApi.formatMoney(liveTotals.total) }}</span>
          </div>
          
          <!-- Overlay de carga solo en edición -->
          <div v-if="isPreviewLoading && mode !== 'detail'" class="mini-spinner-overlay">
            <div class="mini-spinner"></div>
          </div>
        </section>
      </div>
    </FormSection>

    <template #footer v-if="mode === 'detail' && quote">
      <div class="audit-info">
        <p>Documento generado por el sistema TramaTex.</p>
        <p v-if="quote.id">ID único: <code>{{ quote.id }}</code></p>
      </div>
    </template>
  </BaseEntityPage>

  <BaseEntityPage class="no-print" v-else>
    <template #header>
      <BasePageHeader title="Estado Indeterminado" :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-content">
        <h3>No hay datos para mostrar</h3>
        <p>El presupuesto solicitado no ha cargado correctamente o la sesión ha expirado.</p>
        <button class="btn btn-primary mt-4" @click="initComponent">Reintentar Carga</button>
      </div>
    </div>
  </BaseEntityPage>

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

  <!-- MODAL: CONVERTIR A PEDIDO -->
  <BaseDialog
    :show="showConvertModal"
    title="Convertir a Pedido"
    icon="shopping_cart"
    confirm-text="Confirmar y Crear Pedido"
    confirm-class="btn-success"
    :is-confirming="isConverting"
    @close="showConvertModal = false"
    @confirm="convertToOrder"
  >
    <p>¿Está seguro de que desea convertir este presupuesto en un pedido en firme?</p>
    <p class="mt-2 text-muted italic">Esta acción no se puede deshacer y el presupuesto quedará marcado como convertido.</p>
  </BaseDialog>

  <!-- PORTAL DE IMPRESIÓN (Solo visible en @media print) -->
  <div class="print-container">
    <PrintDocument
      v-if="quote"
      type="QUOTE"
      :number="quote.quoteNumber"
      :date="quote.quoteDate"
      :expiry-date="quote.expirationDate"
      :customer-name="partyName"
      :customer-tax-id="quote.taxId"
      :items="quote.lineItems"
      :totals="{ subtotal: quote.subtotal, taxAmount: quote.taxAmount, total: quote.total }"
      :notes="quote.notes"
    />
  </div>

  <!-- Post-Issue Actions Modal -->
  <Transition name="fade">
    <div v-if="showPostIssueModal" class="modal-backdrop no-print">
      <div class="modal card w-modal-md animate-fade-in">
        <div class="modal-header border-none pb-0">
          <div class="icon-circle success">
            <CheckCircle :size="32" />
          </div>
          <button class="btn-icon ml-auto" @click="showPostIssueModal = false"><X :size="20" /></button>
        </div>
        <div class="modal-body text-center p-6 pt-2">
          <h2 class="mb-2">¡Presupuesto Emitido!</h2>
          <p class="text-secondary mb-6">El presupuesto <strong>{{ quote?.quoteNumber }}</strong> se ha emitido correctamente. ¿Qué desea hacer ahora?</p>
          
          <div class="post-issue-actions">
            <button class="btn btn-primary w-full justify-center mb-3 py-3" @click="postIssuePrint">
              <Printer :size="16" /> <span>Imprimir Presupuesto</span>
            </button>
            <button class="btn btn-outline w-full justify-center" @click="showPostIssueModal = false">
              <span>Continuar trabajando</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue';
import { useRoute, useRouter, RouterLink } from 'vue-router';
import { 
  AlertCircle, 
  FileText, 
  Printer, 
  Pencil, 
  RefreshCw, 
  Save, 
  Send, 
  CheckCircle, 
  ShoppingCart, 
  XCircle, 
  User, 
  Calendar, 
  CalendarOff, 
  CreditCard, 
  ExternalLink, 
  Settings2, 
  Plus, 
  Trash2, 
  X,
  AlertTriangle
} from 'lucide-vue-next';
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import BasePageHeader from '@/components/shared/BasePageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import OrderLines from '@/components/sales/OrderLines.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import PrintDocument from '@/components/sales/PrintDocument.vue';
import { useLineNavigation } from '@/composables/useLineNavigation';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { mesApi } from '@/services/mesApi';
import { useToastStore } from '@/stores/toast';
import '@/assets/sales-print.css';

const route = useRoute();
const router = useRouter();
const toastStore = useToastStore();

const isCreateMode = computed(() => !route.params.id || route.params.id === 'new')
const mode = ref(isCreateMode.value ? 'create' : 'detail');
const isLoading = ref(true); 
const isSaving = ref(false);
const error = ref('');

const quote = ref(null);
const partyName = ref('');
const partyDefaultDiscount = ref(0);
const mesWorksCache = ref({});
const availableMesSetups = ref([]);

const workTypesCache = ref({});
const positionsCache = ref({});

const formData = reactive({
  partyId: '', expirationDate: '', notes: '', mesWorkRefs: [], lineItems: []
});

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

function promptIssueQuote() {
  confirmDialog.title = 'Emitir Presupuesto'
  confirmDialog.message = '¿Desea marcar este presupuesto como EMITIDO? Se generará el documento oficial para el cliente.'
  confirmDialog.icon = Send
  confirmDialog.confirmText = 'Sí, Emitir'
  confirmDialog.confirmClass = 'btn-success'
  confirmDialog.action = confirmIssueQuote
  confirmDialog.show = true
}

function promptRejectQuote() {
  confirmDialog.title = 'Rechazar Presupuesto'
  confirmDialog.message = '¿Desea marcar este presupuesto como RECHAZADO? El proceso comercial se detendrá.'
  confirmDialog.icon = XCircle
  confirmDialog.confirmText = 'Rechazar'
  confirmDialog.confirmClass = 'btn-danger'
  confirmDialog.action = rejectQuote
  confirmDialog.show = true
}

function promptReactivateQuote() {
  confirmDialog.title = 'Reactivar Presupuesto'
  confirmDialog.message = '¿Desea volver a poner este presupuesto en estado BORRADOR para poder editarlo?'
  confirmDialog.icon = RefreshCw
  confirmDialog.confirmText = 'Reactivar'
  confirmDialog.confirmClass = 'btn-primary'
  confirmDialog.action = reactivateQuote
  confirmDialog.show = true
}

async function handleConfirmDialog() {
  if (confirmDialog.action) {
    await confirmDialog.action()
  }
  confirmDialog.show = false
}

const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewTimer = null;
const showVariantSelector = ref(false);
const showConvertModal = ref(false);
const isConverting = ref(false);
const showPostIssueModal = ref(false);
const addProductBtn = ref(null);

// Watcher para cambios en el formulario que requieran recalcular totales
watch(() => [formData.partyId, formData.lineItems], () => {
  if (mode.value !== 'detail' && !isSaving.value) {
    calculateTotals();
  }
}, { deep: true });

watch(() => route.params.id, async (newId) => {
  if (newId && newId !== 'new') {
    mode.value = 'detail';
    await fetchQuote();
  } else {
    mode.value = 'create';
    resetForm();
    isLoading.value = false;
  }
});

onMounted(async () => {
  await initComponent();
  await loadMesMasters();
  window.addEventListener('tramatex-save', handleGlobalSave);
  window.addEventListener('keydown', handleQuoteKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener('tramatex-save', handleGlobalSave);
  window.removeEventListener('keydown', handleQuoteKeydown);
});

function handleQuoteKeydown(e) {
  if (e.ctrlKey && e.key === 'e') {
    e.preventDefault();
    if (mode.value === 'detail' && canEdit.value) {
      enterEditMode();
    }
  }
}

function handleGlobalSave() {
  if (mode.value !== 'detail' && !isSaving.value) {
    saveQuote();
  }
}

async function initComponent() {
  const id = route.params.id;
  const isCreate = id === 'new' || route.name === 'CreateQuote' || !id;

  if (isCreate) {
    mode.value = 'create';
    resetForm();
    isLoading.value = false;
  } else {
    mode.value = 'detail';
    await fetchQuote();
  }
}

function resetForm() {
  partyName.value = 'Cliente no seleccionado';
  Object.assign(formData, {
    partyId: '', 
    expirationDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    notes: '', 
    mesWorkRefs: [], 
    lineItems: []
  });
  quote.value = null;
}

async function fetchQuote() {
  const id = route.params.id;
  if (!id || id === 'new') {
    isLoading.value = false;
    return;
  }

  isLoading.value = true;
  error.value = '';
  try {
    const res = await salesApi.getQuote(id);
    const data = res?.data || res;
    if (!data || (!data.id && !data.ID && !data.quoteNumber)) throw new Error('Presupuesto no encontrado');
    if (data.ID && !data.id) data.id = data.ID;
    quote.value = data;
    if (quote.value.partyId) {
      try {
        const party = await partyApi.getParty(quote.value.partyId);
        partyName.value = party?.name || party?.displayName || 'Cliente sin nombre';
        partyDefaultDiscount.value = party?.default_discount_percentage || 0;
      } catch (pErr) { partyName.value = 'Cliente (ID: ' + quote.value.partyId.substring(0,8) + ')'; }
    }
    loadMesWorksNames().catch(() => {});
  } catch (err) { 
    console.error('Error fetchQuote:', err);
    error.value = 'Error al cargar el presupuesto: ' + err.message; 
  }
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
  const refs = quote.value?.mesWorkRefs || [];
  if (!refs.length) return;
  const ids = refs.map(r => r.workSetupId).filter(Boolean);
  const results = await Promise.allSettled(ids.map(id => mesApi.getWorkSetup(id)));
  results.forEach((r, i) => { if (r.status === 'fulfilled') mesWorksCache.value[ids[i]] = r.value; });
}

async function enterEditMode() {
  if (!quote.value) return;
  const data = {
    partyId: quote.value.partyId,
    expirationDate: quote.value.expirationDate ? new Date(quote.value.expirationDate).toISOString().split('T')[0] : '',
    notes: quote.value.notes || '',
    mesWorkRefs: (quote.value.mesWorkRefs || []).map(r => ({ id: r.id, workSetupId: r.workSetupId || null, description: r.description || '' })),
    lineItems: (quote.value.lineItems || []).map(i => ({
      productVariantId: i.productVariantId || i.productVariantID,
      variantSku: i.variantSku,
      productName: buildDisplayName(i),
      quantity: i.quantity,
      listPrice: i.listUnitPrice?.amount ?? i.unitPrice?.amount ?? 0,
      unitPrice: i.unitPrice?.amount ?? i.listUnitPrice?.amount ?? 0,
      _autoPrice: false,
      discountPercent: i.discountPercent || 0,
    }))
  };
  Object.assign(formData, data);
  mode.value = 'edit';
  loadAvailableSetups(quote.value.partyId);
  calculateTotals();
}

function exitEditMode() { if (mode.value === 'edit') mode.value = 'detail'; else router.push('/sales/quotes'); }

async function loadAvailableSetups(partyId) {
  if (!partyId) return;
  try { availableMesSetups.value = await mesApi.listWorkSetups({ party_id: partyId }); }
  catch (err) { availableMesSetups.value = []; }
}

function onPartySelected(party) {
  partyName.value = party?.name || 'Cliente no seleccionado';
  partyDefaultDiscount.value = party?.default_discount_percentage || 0;
  formData.partyId = party?.id || '';
  loadAvailableSetups(party?.id);
  calculateTotals();
}

const { handleLineKeyDown: handleMesKeyDown, focusLineInput: focusMesInput } = useLineNavigation({
  rowCount: () => formData.mesWorkRefs.length,
  columns: ['setup', 'desc'],
  prefix: 'mes',
  onRemoveField: (index) => removeMesWorkRef(index),
  onLastFieldEnter: () => addMesWorkRef(),
  onAddField: () => addMesWorkRef()
});

function openVariantSelector() { 
  showVariantSelector.value = true; 
}

function handleVariantSelected(payload) {
  const variant = payload.variant || payload;
  const newItem = {
    productVariantId: variant.id,
    variantSku: variant.sku,
    productName: (variant.product_name || 'Producto') + (variant.option_configuration ? ' - ' + Object.values(variant.option_configuration).join(', ') : ''),
    quantity: 1,
    listPrice: null,
    unitPrice: null,
    _autoPrice: true,
    discountPercent: partyDefaultDiscount.value || 0
  };

  formData.lineItems.push(newItem);
  showVariantSelector.value = false;

  // Position focus on the quantity of the new line
  nextTick(() => {
    fetchPreviewCalculation();
    const lastIdx = formData.lineItems.length - 1
    const el = document.querySelector(`input[data-row="${lastIdx}"][data-col="quantity"]`)
    if (el) {
      el.focus()
      el.select()
    }
  });
}

function addMesWorkRef() { formData.mesWorkRefs.push({ workSetupId: null, description: '' }); }
function removeMesWorkRef(idx) { formData.mesWorkRefs.splice(idx, 1); }

function calculateTotals() {
  clearTimeout(previewTimer);
  // Reset preview so UI falls back to instant local totals while new preview is fetched.
  previewResult.value = null;
  previewTimer = setTimeout(fetchPreviewCalculation, 400);
}

async function fetchPreviewCalculation() {
  const partyId = formData.partyId || quote.value?.partyId;
  const items = formData.lineItems.map(i => ({ 
    productVariantId: i.productVariantId, 
    quantity: Number(i.quantity || 0), 
    ...(!i._autoPrice ? { unitPrice: { amount: Number(i.unitPrice || 0), currency: 'EUR' } } : {}),
    discountPercent: Number(i.discountPercent || 0) 
  }));

  if (!partyId || !items.length) { 
    previewResult.value = null; 
    return; 
  }

  isPreviewLoading.value = true;
  try { 
    const res = await salesApi.previewQuoteCalculation(partyId, items);
    if (res) {
      previewResult.value = res;
      // Populate unit prices from pricing engine for auto-priced items
      formData.lineItems.forEach((item, idx) => {
        if (item._autoPrice && res.lineItems?.[idx]) {
          item.unitPrice = res.lineItems[idx].unitPrice.amount;
          item.listPrice = res.lineItems[idx].listUnitPrice?.amount ?? item.listPrice;
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

function calculateLineSubtotalLocal(item) {
  if (!item) return 0;
  const quantity = Number(item.quantity || 0);
  const unitPrice = Number(item.unitPrice || 0);
  const discountPercent = Number(item.discountPercent || 0);
  return quantity * unitPrice * (1 - discountPercent / 100);
}

const localDraftTotals = computed(() => {
  const subtotalAmount = formData.lineItems.reduce((acc, item) => acc + calculateLineSubtotalLocal(item), 0);
  const taxAmount = formData.lineItems.reduce((acc, item) => {
    const lineSubtotal = calculateLineSubtotalLocal(item);
    const taxRate = Number(item.taxRate ?? 21);
    return acc + (lineSubtotal * taxRate / 100);
  }, 0);

  return {
    subtotal: { amount: subtotalAmount, currency: 'EUR' },
    taxAmount: { amount: taxAmount, currency: 'EUR' },
    total: { amount: subtotalAmount + taxAmount, currency: 'EUR' },
  };
});

const liveTotals = computed(() => {
  if (mode.value === 'detail' && quote.value) return { subtotal: quote.value.subtotal, taxAmount: quote.value.taxAmount, total: quote.value.total };
  if (previewResult.value && !isPreviewLoading.value) return { subtotal: previewResult.value.subtotal, taxAmount: previewResult.value.taxAmount, total: previewResult.value.total };
  return localDraftTotals.value;
});

async function saveQuote() {
  if (!formData.partyId) { toastStore.error('Seleccione un cliente'); return; }
  if (formData.lineItems.length === 0) { toastStore.error('Añada al menos un producto'); return; }
  isSaving.value = true;
  try {
    let isoExpiration = undefined;
    if (formData.expirationDate) {
      const d = new Date(formData.expirationDate);
      d.setHours(23, 59, 59);
      isoExpiration = d.toISOString();
    }
    const payload = {
      partyId: formData.partyId,
      expirationDate: isoExpiration,
      notes: formData.notes || '',
      items: formData.lineItems.map(i => ({ productVariantId: i.productVariantId, quantity: Number(i.quantity), ...(!i._autoPrice ? { unitPrice: { amount: Number(i.unitPrice || 0), currency: 'EUR' } } : {}), discountPercent: Number(i.discountPercent || 0) })),
      mesWorkRefs: formData.mesWorkRefs.filter(r => r.workSetupId || r.description).map(r => ({ workSetupId: r.workSetupId || undefined, description: r.description || '' }))
    };
    if (mode.value === 'create') {
      const newQuote = await salesApi.createQuote(payload);
      await router.push(`/sales/quotes/${newQuote.id}`);
    } else {
      await salesApi.updateQuote(quote.value.id, payload);
      await fetchQuote();
      mode.value = 'detail';
    }
  } catch (err) { toastStore.error('Error al guardar: ' + err.message); }
  finally { isSaving.value = false; }
}

const isExpired = computed(() => quote.value?.expirationDate && new Date(quote.value.expirationDate) < new Date());
const canEdit = computed(() => {
  if (!quote.value) return false;
  const status = quote.value.status;
  return ['DRAFT', 'BORRADOR', 'ISSUED', 'EMITIDA', 'PENDING', 'PENDIENTE'].includes(status);
});

function getStatusLabel(s) { return salesApi.getStatusLabel(s); }
function getStatusClass(s) { return salesApi.getStatusClass(s); }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '—'; }
function buildDisplayName(i) { return (i.productName || i.displayName || 'Producto') + (i.optionConfiguration ? ' - ' + Object.values(i.optionConfiguration).join(', ') : ''); }
function formatVariantId(id) { return id ? id.substring(0, 8) : '—'; }
function formatMesWorkId(id) { return id ? (mesWorksCache.value[id]?.name || id.substring(0, 8)) : 'Sin config.'; }

async function confirmIssueQuote() { 
  try { 
    await salesApi.changeQuoteStatus(quote.value.id, 'ISSUED'); 
    await fetchQuote();
    showPostIssueModal.value = true;
    toastStore.success('Presupuesto emitido');
  } catch (err) { 
    toastStore.error(err.message); 
  } 
}

function postIssuePrint() { 
  showPostIssueModal.value = false; 
  window.print(); 
}
async function rejectQuote() { try { await salesApi.changeQuoteStatus(quote.value.id, 'REJECTED'); await fetchQuote(); toastStore.info('Presupuesto rechazado'); } catch (err) { toastStore.error(err.message); } }
async function reactivateQuote() { try { await salesApi.changeQuoteStatus(quote.value.id, 'DRAFT'); await fetchQuote(); toastStore.success('Presupuesto reactivado'); } catch (err) { toastStore.error(err.message); } }

async function convertToOrder() {
  isConverting.value = true;
  try {
    // Set default delivery date: 15 days from today
    const deliveryDateObj = new Date();
    deliveryDateObj.setDate(deliveryDateObj.getDate() + 15);
    const deliveryDate = deliveryDateObj.toISOString().split('T')[0];
    
    // The backend auto-approves the quote (DRAFT→ISSUED→APPROVED→CONVERTED)
    const order = await salesApi.createOrderFromQuote(quote.value.id, deliveryDate);
    toastStore.success('Pedido generado con éxito');
    router.push(`/sales/orders/${order.id}`);
  } catch (err) { 
    toastStore.error('Error al convertir presupuesto: ' + err.message); 
  }
  finally { isConverting.value = false; }
}

function printQuote() { window.print(); }
</script>

<style scoped>
@import "@/design-system/_sections.css";

.action-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; box-shadow: var(--box-shadow-sm); }
.status-badge { padding: 0.4rem 1rem; font-size: 0.85rem; font-weight: 800; }
.toolbar-buttons { display: flex; gap: 0.75rem; }

.mes-config-info { display: flex; align-items: center; gap: 0.75rem; }
.mes-specs-inline { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.spec-pill { display: inline-flex; align-items: center; gap: 0.35rem; padding: 0.2rem 0.5rem; background: var(--color-background); border-radius: 4px; border: 1px solid var(--color-border); font-size: 0.8rem; }
.spec-pill .label { color: var(--color-text-secondary); font-weight: 600; }
.spec-pill .value { color: var(--color-text-primary); font-weight: 700; }
.notes-text-sm { font-size: 0.85rem; font-style: italic; color: var(--color-text-secondary); margin: 0; }

.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-input-sm { padding: 0.5rem; border: 1px solid var(--color-border); border-radius: 4px; font-size: 0.85rem; }

.totals-checkout-layout { display: flex; justify-content: flex-end; margin-top: 1rem; }
.totals-checkout-card { width: 400px; padding: 1.5rem; background: white; border: 1px solid var(--color-border-strong); border-radius: 12px; box-shadow: var(--box-shadow-md); }
.total-row { display: flex; justify-content: space-between; margin-bottom: 0.75rem; }
.total-row.final { margin-top: 1rem; padding-top: 1rem; border-top: 2px solid var(--color-border); font-weight: 800; font-size: 1.25rem; }
.total-value { color: var(--color-secondary); }

/* Totals Loading State */
.totals-checkout-card { position: relative; transition: opacity 0.3s ease; }
.is-loading-overlay { opacity: 0.7; pointer-events: none; }
.mini-spinner-overlay { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; background: rgba(255, 255, 255, 0.4); border-radius: 12px; z-index: 5; }
.mini-spinner { width: 24px; height: 24px; border: 3px solid var(--color-border); border-top-color: var(--color-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.w-64 { width: 16rem; } .w-16 { width: 4rem; } .w-24 { width: 6rem; }
.w-full { width: 100%; } .fixed-layout { table-layout: fixed; }
.text-danger { color: var(--color-danger); }
.audit-info { color: var(--color-text-secondary); font-size: 0.8rem; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; }

/* ESTILOS DE IMPRESIÓN PROFESIONAL */
.print-container { display: none; }

@media print {
  .no-print { display: none !important; }
  .print-container { display: block !important; position: absolute; left: 0; top: 0; width: 100%; }
}
</style>>>>