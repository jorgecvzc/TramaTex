<template>
  
  
  <BaseEntityPage v-if="isLoading">
    <template #header>
      <PageHeader title="Cargando..." :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos' }]" />
    </template>
    <div class="loading-state card">
      <div class="spinner"></div>
      <p>Cargando información del presupuesto...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="error">
    <template #header>
      <PageHeader title="Error" :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-icon-wrapper error">
        <span class="material-symbols-outlined">error</span>
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/sales/quotes')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="quote || mode === 'create'">
    <!-- 1. IDENTITY HEADER -->
    <template #header>
      <PageHeader 
        :title="mode === 'create' ? 'Nuevo Presupuesto' : (mode === 'edit' ? `Editando Presupuesto ${quote?.quoteNumber}` : `Presupuesto ${quote?.quoteNumber}`)" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos', to: '/sales/quotes' }, { label: mode === 'create' ? 'Crear' : quote?.quoteNumber }]"
      >
        <template #icon>
          <span class="material-symbols-outlined">description</span>
        </template>
        <template #actions>
          <template v-if="mode === 'detail'">
            <button v-if="quote?.status !== 'BORRADOR'" class="btn btn-outline" @click="printQuote">
              <span class="material-symbols-outlined">print</span> <span>Imprimir</span>
            </button>
            <button v-if="canEdit" class="btn btn-primary" @click="enterEditMode">
              <span class="material-symbols-outlined">edit</span> <span>Editar Presupuesto</span>
            </button>
          </template>
          <template v-else>
            <button class="btn btn-outline" @click="exitEditMode" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-secondary" @click="saveQuote" :disabled="isSaving">
              <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
              <span>{{ isSaving ? 'Guardando...' : 'Guardar Presupuesto' }}</span>
            </button>
          </template>
        </template>
      </PageHeader>
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
          <button v-if="['BORRADOR', 'DRAFT'].includes(quote.status)" class="btn btn-success btn-sm" @click="confirmIssueQuote">
            <span class="material-symbols-outlined">send</span> <span>Emitir a Cliente</span>
          </button>
          <button v-if="['EMITIDA', 'ISSUED'].includes(quote.status) && !isExpired" class="btn btn-success btn-sm" @click="showConvertModal = true">
            <span class="material-symbols-outlined">check_circle</span> <span>Aceptar y Crear Pedido</span>
          </button>
          <button v-if="['EMITIDA', 'ISSUED'].includes(quote.status)" class="btn btn-danger btn-sm" @click="rejectQuote">
            <span class="material-symbols-outlined">cancel</span> <span>Rechazar</span>
          </button>
          <button v-if="['RECHAZADA', 'REJECTED'].includes(quote.status)" class="btn btn-primary btn-sm" @click="reactivateQuote">
            <span class="material-symbols-outlined">refresh</span> <span>Reactivar</span>
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
          <div class="tag-content"><label>Fecha Emisión</label><strong>{{ formatDate(mode === 'create' ? new Date() : quote?.quoteDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><span class="material-symbols-outlined">event_busy</span></div>
          <div class="tag-content"><label>Válido Hasta</label><strong :class="{'text-danger': isExpired}">{{ formatDate(mode === 'create' ? formData.expirationDate : quote?.expirationDate) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><span class="material-symbols-outlined">payments</span></div>
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
          <div class="tag-icon"><span class="material-symbols-outlined">shopping_cart</span></div>
          <div class="tag-content">
            <label>Pedido Generado</label>
            <strong>{{ quote.generatedOrderNumber || 'Ver Pedido' }}</strong>
          </div>
          <span class="material-symbols-outlined jump-icon">open_in_new</span>
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
            <span class="material-symbols-outlined">add</span> <span>Añadir Trabajo MES</span>
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
                  <select v-model="ref.workSetupId" class="form-input-sm w-full">
                    <option :value="null">-- Personalizado --</option>
                    <option v-for="setup in availableMesSetups" :key="setup.id" :value="setup.id">{{ setup.name }}</option>
                  </select>
                </td>
                <td class="w-full">
                  <input v-model="ref.description" type="text" class="form-input-sm w-full" placeholder="Especificaciones técnicas..." required />
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

    <FormSection title="Líneas del Presupuesto" icon="list_alt">
      <div v-if="mode !== 'detail'" class="mb-4">
        <button type="button" class="btn btn-primary btn-sm" @click="openVariantSelector">
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
            <tr v-for="(item, idx) in (mode === 'detail' ? quote.lineItems : formData.lineItems)" :key="idx">
              <td><code class="code-badge">{{ item.variantSku || formatVariantId(item.productVariantId || item.productVariantID) }}</code></td>
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
                <input v-else v-model.number="item.discountPercent" type="number" step="0.01" class="form-input-sm w-16 text-center" />
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
        <DataRow label="Subtotal" :value="salesApi.formatMoney(quote.subtotal)" />
        <DataRow label="Impuestos" :value="salesApi.formatMoney(quote.taxAmount)" />
        <DataRow label="TOTAL PRESUPUESTO" :value="salesApi.formatMoney(quote.total)" highlight />
      </div>
      <div v-else class="totals-checkout-layout">
        <section class="totals-checkout-card">
          <div class="total-row"><label>Subtotal:</label><span>{{ salesApi.formatMoney(liveTotals.subtotal) }}</span></div>
          <div class="total-row"><label>IVA (21%):</label><span>{{ salesApi.formatMoney(liveTotals.taxAmount) }}</span></div>
          <div class="total-row final"><label>TOTAL ESTIMADO:</label><span class="total-value">{{ salesApi.formatMoney(liveTotals.total) }}</span></div>
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

  <BaseEntityPage v-else>
    <template #header>
      <PageHeader title="Estado Indeterminado" :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-content">
        <h3>No hay datos para mostrar</h3>
        <p>El presupuesto solicitado no ha cargado correctamente o la sesión ha expirado.</p>
        <button class="btn btn-primary mt-4" @click="initComponent">Reintentar Carga</button>
      </div>
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
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue';
import { useRoute, useRouter, RouterLink } from 'vue-router';
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { mesApi } from '@/services/mesApi';

const route = useRoute();
const router = useRouter();

const mode = ref('detail');
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

const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewTimer = null;
const showVariantSelector = ref(false);
const showConvertModal = ref(false);
const isConverting = ref(false);

const lineItemsFingerprint = computed(() =>
  mode.value !== 'detail'
    ? formData.lineItems.map(i => `${i.productVariantId}|${i.quantity}|${i.unitPrice}|${i.discountPercent}`).join('§')
    : ''
);
watch(lineItemsFingerprint, (val) => {
  if (val && !isSaving.value) calculateTotals();
});

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
});

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
    lineItems: (quote.value.lineItems || []).map(i => ({ productVariantId: i.productVariantId || i.productVariantID, variantSku: i.variantSku, displayName: buildDisplayName(i), quantity: i.quantity, unitPrice: i.listUnitPrice?.amount ?? i.unitPrice?.amount ?? 0, discountPercent: i.discountPercent || 0 }))
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

function openVariantSelector() { 
  showVariantSelector.value = true; 
}

function handleVariantSelected(payload) {
  const variant = payload.variant || payload;
  formData.lineItems.push({
    productVariantId: variant.id,
    variantSku: variant.sku,
    displayName: (variant.product_name || 'Producto') + (variant.option_configuration ? ' - ' + Object.values(variant.option_configuration).join(', ') : ''),
    quantity: 1,
    unitPrice: variant.product_base_price || 0,
    discountPercent: partyDefaultDiscount.value || 0
  });
  showVariantSelector.value = false;
}

function removeLineItem(idx) { formData.lineItems.splice(idx, 1); }
function addMesWorkRef() { formData.mesWorkRefs.push({ workSetupId: null, description: '' }); }
function removeMesWorkRef(idx) { formData.mesWorkRefs.splice(idx, 1); }

function calculateTotals() {
  clearTimeout(previewTimer);
  previewTimer = setTimeout(fetchPreviewCalculation, 400);
}

async function fetchPreviewCalculation() {
  const partyId = mode.value === 'create' ? formData.partyId : quote.value?.partyId;
  const items = formData.lineItems.map(i => ({ productVariantId: i.productVariantId, quantity: i.quantity, unitPrice: { amount: i.unitPrice, currency: 'EUR' }, discountPercent: i.discountPercent }));
  if (!partyId || !items.length) { previewResult.value = null; return; }
  isPreviewLoading.value = true;
  try { previewResult.value = await salesApi.previewQuoteCalculation(partyId, items); } catch (err) {}
  finally { isPreviewLoading.value = false; }
}

const liveTotals = computed(() => {
  if (mode.value === 'detail' && quote.value) return { subtotal: quote.value.subtotal, taxAmount: quote.value.taxAmount, total: quote.value.total };
  if (previewResult.value) return { subtotal: previewResult.value.subtotal, taxAmount: previewResult.value.taxAmount, total: previewResult.value.total };
  return { subtotal: { amount: 0, currency: 'EUR' }, taxAmount: { amount: 0, currency: 'EUR' }, total: { amount: 0, currency: 'EUR' } };
});

function calculateLineSubtotal(idx) {
  const item = formData.lineItems[idx];
  if (!item) return 0;
  if (previewResult.value?.lineItems) {
    const calculated = previewResult.value.lineItems[idx];
    if (calculated?.subtotal?.amount !== undefined) return calculated.subtotal.amount;
  }
  return item.quantity * item.unitPrice * (1 - (item.discountPercent || 0) / 100);
}

async function saveQuote() {
  if (!formData.partyId) { alert('Seleccione un cliente'); return; }
  if (formData.lineItems.length === 0) { alert('Añada al menos un producto'); return; }
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
      items: formData.lineItems.map(i => ({ productVariantId: i.productVariantId, quantity: Number(i.quantity), unitPrice: { amount: Number(i.unitPrice), currency: 'EUR' }, discountPercent: Number(i.discountPercent || 0) })),
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
  } catch (err) { alert('Error al guardar: ' + err.message); }
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

async function confirmIssueQuote() { try { await salesApi.changeQuoteStatus(quote.value.id, 'EMITIDA'); await fetchQuote(); } catch (err) { alert(err.message); } }
async function rejectQuote() { if (confirm('¿Rechazar?')) { try { await salesApi.changeQuoteStatus(quote.value.id, 'RECHAZADA'); await fetchQuote(); } catch (err) { alert(err.message); } } }
async function reactivateQuote() { try { await salesApi.changeQuoteStatus(quote.value.id, 'BORRADOR'); await fetchQuote(); } catch (err) { alert(err.message); } }

async function convertToOrder() {
  isConverting.value = true;
  try {
    const order = await salesApi.createOrderFromQuote(quote.value.id);
    router.push(`/sales/orders/${order.id}`);
  } catch (err) { alert(err.message); }
  finally { isConverting.value = false; }
}

function printQuote() { window.print(); }
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
.jump-icon { font-size: 18px; color: var(--color-text-secondary); opacity: 0.5; margin-left: auto; }

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

.w-64 { width: 16rem; } .w-16 { width: 4rem; } .w-24 { width: 6rem; }
.w-full { width: 100%; } .fixed-layout { table-layout: fixed; }
.text-danger { color: var(--color-danger); }
.audit-info { color: var(--color-text-secondary); font-size: 0.8rem; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; }
</style>
