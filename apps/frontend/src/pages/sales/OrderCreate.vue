<template>
  <div class="page-layout">
    <!-- CAPA 1: IDENTIDAD -->
    <BasePageHeader 
      title="Crear Nuevo Pedido" 
      :breadcrumbs="[{ label: 'Ventas', to: '/sales/dashboard' }, { label: 'Pedidos', to: '/sales/orders' }, { label: 'Nuevo' }]"
    >
      <template #icon><ShoppingCart :size="28" /></template>
      <template #actions>
        <button @click="goBack" class="btn btn-outline btn-sm">
          <X :size="16" />
          <span>Cancelar</span>
        </button>
        <button 
          @click="handleSubmit" 
          class="btn btn-primary btn-sm ml-2" 
          :disabled="!isFormValid || isSubmitting"
        >
          <component :is="isSubmitting ? RefreshCw : Save" :size="16" :class="{ 'spin': isSubmitting }" />
          <span>{{ isSubmitting ? 'Procesando...' : 'Crear Pedido' }}</span>
        </button>
      </template>
    </BasePageHeader>

    <main class="page-content">
      <div class="transactional-layout">
        <!-- ÁREA PRINCIPAL: DATOS Y LÍNEAS -->
        <div class="main-form-area">
          
          <!-- SECCIÓN: CLIENTE (CONTEXTO) -->
          <section class="card mb-4">
            <div class="section-header">
              <User :size="20" />
              <h2>Cliente y Datos Fiscales</h2>
            </div>
            <div class="p-4">
              <PartySelector
                v-model="formData.partyId"
                label="Seleccionar Cliente *"
                placeholder="Busca por nombre, NIF o ID..."
                role-filter="CLIENT"
                :required="true"
                @select="onPartySelected"
              />
            </div>
          </section>

          <!-- SECCIÓN: CONDICIONES -->
          <section class="card mb-4">
            <div class="section-header">
              <Calendar :size="20" />
              <h2>Condiciones de Venta</h2>
            </div>
            <div class="p-4">
              <div class="form-grid">
                <div class="form-group">
                  <label>Fecha del Pedido</label>
                  <input v-model="formData.orderDate" type="date" required />
                </div>
                <div class="form-group">
                  <label>Entrega Estimada</label>
                  <input v-model="formData.deliveryDate" type="date" :min="minDeliveryDate" required />
                </div>
              </div>
              <div class="form-group mt-2">
                <label>Observaciones del Cliente / Logística</label>
                <textarea v-model="formData.notes" rows="2" placeholder="Instrucciones especiales..."></textarea>
              </div>
            </div>
          </section>

          <!-- SECCIÓN: LÍNEAS DE PRODUCTO (TRABAJO) -->
          <section class="card table-card mb-4">
            <div class="table-header-row">
              <div class="section-header">
                <List :size="20" />
                <h2>Líneas de Producto</h2>
              </div>
            </div>

            <div class="p-4">
              <OrderLines
                :lines="formData.lineItems"
                :is-editing="true"
                @update:lines="(newLines) => formData.lineItems = newLines"
                @add-line="handleAddLineRequest"
                @last-field-tab="focusAddButton"
              />
              <div class="mt-4">
                <button 
                  ref="addProductBtnRef"
                  type="button" 
                  class="btn btn-secondary btn-sm" 
                  @click="handleAddLineRequest"
                >
                  <Plus :size="16" /> Añadir Producto (Ins)
                </button>
              </div>
            </div>
          </section>
        </div>

        <!-- SIDEBAR: TOTALES Y MES -->
        <aside class="side-summary-area">
          <!-- RESUMEN ECONÓMICO -->
          <section class="card totals-card mb-4">
            <div class="section-header">
              <CreditCard :size="20" />
              <h2>Resumen de Venta</h2>
            </div>
            <div class="totals-body">
              <div class="total-line">
                <label>Base Imponible</label>
                <span>{{ formatMoney(calculatedTotals.subtotal) }}</span>
              </div>
              <div class="total-line">
                <label>Impuestos (21%)</label>
                <span>{{ formatMoney(calculatedTotals.tax) }}</span>
              </div>
              <div class="total-line final">
                <label>Total Pedido</label>
                <strong>{{ formatMoney(calculatedTotals.total) }}</strong>
              </div>
            </div>
          </section>

          <!-- CONFIGURACIONES MES -->
          <section v-if="formData.partyId" class="card mes-summary-card">
            <div class="section-header">
              <Factory :size="20" />
              <h2>Configuración Taller</h2>
            </div>
            <div class="mes-body p-3">
              <div v-for="(config, idx) in formData.mesWorkRefs" :key="idx" class="mes-item-row">
                <div class="mes-item-info">
                  <select 
                    v-model="config.workSetupId" 
                    class="mes-select"
                    :data-mes-row="idx"
                    data-mes-col="setup"
                    @keydown="handleMesKeyDown($event, idx, 'setup', config)"
                  >
                    <option :value="null">Sin configuración base</option>
                    <option v-for="ws in mesWorkSetups" :key="ws.id" :value="ws.id">{{ ws.name }}</option>
                  </select>
                  <input 
                    v-model="config.description" 
                    type="text" 
                    placeholder="Notas taller..." 
                    class="mes-input" 
                    :data-mes-row="idx"
                    data-mes-col="desc"
                    @keydown="handleMesKeyDown($event, idx, 'desc', config)"
                  />
                </div>
                <button @click="removeConfig(idx)" class="btn-icon text-danger"><X :size="16" /></button>
              </div>
              <button class="btn btn-secondary btn-sm w-full mt-2" @click="addConfig">
                <Plus :size="16" /> Añadir Trabajo Taller
              </button>
            </div>
          </section>
        </aside>
      </div>
    </main>

    <!-- MODAL: SELECTOR DE PRODUCTOS -->
    <BaseDialog :show="showVariantSelector" title="Buscar Producto" icon="inventory_2" size="lg" hide-actions @close="showVariantSelector = false">
      <div class="p-4 overflow-y-auto" style="max-height: 70vh">
        <VariantSelector
          :key="variantSelectorQuery + '-' + editingLineIndex"
          :product-id="null"
          :initial-query="variantSelectorQuery"
          initial-mode="quick"
          @variant-selected="handleVariantSelected"
        />
      </div>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { 
  ShoppingCart, 
  X, 
  RefreshCw, 
  Save, 
  User, 
  Calendar, 
  List, 
  Plus, 
  ScanBarcode, 
  Search, 
  Trash2, 
  CreditCard, 
  Factory 
} from 'lucide-vue-next';
import BasePageHeader from '@/components/shared/BasePageHeader.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import OrderLines from '@/components/sales/OrderLines.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import { useLineNavigation } from '@/composables/useLineNavigation';
import salesApi from '@/services/salesApi';
// ... (rest of imports)

const { handleLineKeyDown, focusLineInput } = useLineNavigation({
  rowCount: () => formData.value.lineItems.length,
  columns: ['qty', 'price', 'disc'],
  onUpdate: (index, col, val) => {
    const item = formData.value.lineItems[index];
    if (col === 'qty') item.quantity = val;
    else if (col === 'price') item.unitPrice = val;
    else if (col === 'disc') item.discountPercent = val;
    calculateTotals();
  },
  onRemoveField: (index) => removeLineItem(index),
  onAddField: () => addLineItem()
});

const { handleLineKeyDown: handleMesKeyDown, focusLineInput: focusMesInput } = useLineNavigation({
  rowCount: () => formData.value.mesWorkRefs.length,
  columns: ['setup', 'desc'],
  prefix: 'mes',
  onRemoveField: (index) => removeConfig(index),
  onLastFieldEnter: () => addConfig(),
  onAddField: () => addConfig()
});

onMounted(() => {
  window.addEventListener('tramatex-save', handleGlobalSave);
  window.addEventListener('tramatex-esc', goBack);
});

onBeforeUnmount(() => {
  window.removeEventListener('tramatex-save', handleGlobalSave);
  window.removeEventListener('tramatex-esc', goBack);
});

function handleGlobalSave() {
  if (isFormValid.value && !isSubmitting.value) handleSubmit();
}

import { pricingApi } from '@/services/pricingApi';
import { mesApi } from '@/services/mesApi';
import { useToastStore } from '@/stores/toast';

const router = useRouter();
const toastStore = useToastStore();

// --- STATE ---
const formData = ref({
  partyId: '',
  orderDate: new Date().toISOString().split('T')[0],
  deliveryDate: new Date().toISOString().split('T')[0],
  notes: '',
  mesWorkRefs: [],
  lineItems: [],
});

const isSubmitting = ref(false);
const isLoadingMesWorks = ref(false);
const mesWorkSetups = ref([]);
const partyDefaultDiscount = ref(null);
const calculatedTotals = ref({ subtotal: 0, tax: 0, total: 0 });
const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewDebounceTimer = null;

const showVariantSelector = ref(false);
const editingLineIndex = ref(null);
const variantSelectorQuery = ref('');

// --- COMPUTED ---
const minDeliveryDate = computed(() => {
  const tomorrow = new Date(); tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().split('T')[0];
});

const isFormValid = computed(() => (
  formData.value.partyId && formData.value.orderDate && 
  formData.value.lineItems.length > 0 && 
  formData.value.lineItems.every(i => i.productVariantId && i.quantity > 0)
));

// --- WATCHERS ---
watch(() => formData.value.partyId, (id) => { if(id) loadMesWorksForParty(id); });

// --- METHODS ---
function onPartySelected(party) {
  partyDefaultDiscount.value = party?.default_discount_percentage || 0;
  // Aplicar descuento por defecto a todas las líneas actuales
  formData.value.lineItems.forEach(item => {
    item.discountPercent = partyDefaultDiscount.value;
  });
  calculateTotals();
}

async function loadMesWorksForParty(partyId) {
  try { mesWorkSetups.value = await mesApi.listWorkSetups({ party_id: partyId }); }
  catch (err) { mesWorkSetups.value = []; }
}

function addLineItem() {
  formData.value.lineItems.push({
    productVariantId: '', productId: '', selectedVariantName: '', productName: '',
    optionConfiguration: {}, quickSearchQuery: '', quantity: 1, unitPrice: 0, discountPercent: partyDefaultDiscount.value || 0,
  });
}

function openVariantSelector(index) {
  editingLineIndex.value = index;
  variantSelectorQuery.value = formData.value.lineItems[index]?.quickSearchQuery?.trim() || '';
  showVariantSelector.value = true;
}

function handleVariantSelected(payload) {
  const variant = payload?.variant || payload;
  if (editingLineIndex.value !== null && variant) {
    const item = formData.value.lineItems[editingLineIndex.value];
    item.productVariantId = variant.id;
    item.productId = variant.product_id;
    item.selectedVariantName = variant.sku;
    item.productName = variant.product_name;
    item.optionConfiguration = variant.option_configuration || {};
    item.unitPrice = null;
    item._autoPrice = true;
  }
  showVariantSelector.value = false;
  calculateTotals(true);
}

function calculateTotals(immediate = false) {
  clearTimeout(previewDebounceTimer);
  const delay = immediate ? 0 : 400;
  previewDebounceTimer = setTimeout(fetchPreviewCalculation, delay);
}

async function fetchPreviewCalculation() {
  const items = formData.value.lineItems.filter(i => i.productVariantId).map(i => ({
    productVariantId: i.productVariantId, 
    quantity: Number(i.quantity || 0),
    ...(i._autoPrice === false ? { unitPrice: { amount: i.unitPrice.toString(), currency: 'EUR' } } : {}),
    discountPercent: Number(i.discountPercent || 0)
  }));
  
  if (!formData.value.partyId || items.length === 0) { 
    calculatedTotals.value = { subtotal: 0, tax: 0, total: 0 }; 
    return; 
  }
  
  isPreviewLoading.value = true;
  try {
    const res = await salesApi.previewOrderCalculation(formData.value.partyId, items);
    if (res) {
      previewResult.value = res;
      calculatedTotals.value = { 
        subtotal: res.subtotal.amount, 
        tax: res.taxAmount.amount, 
        total: res.total.amount 
      };
      
      // Update line unit prices from pricing engine for auto-priced items
      formData.value.lineItems.forEach((item, idx) => {
        if (item.productVariantId && item._autoPrice !== false && res.lineItems?.[idx]) {
          item.unitPrice = res.lineItems[idx].unitPrice.amount;
        }
      });
    }
  } catch (err) {
    console.error("Error calculating order preview:", err);
  } finally {
    isPreviewLoading.value = false;
  }
}

function calculateLineSubtotal(i) {
  if (!previewResult.value) return 0;
  let pIdx = 0; for(let j=0; j<i; j++) if(formData.value.lineItems[j].productVariantId) pIdx++;
  return previewResult.value.lineItems[pIdx]?.subtotal?.amount ?? 0;
}

async function handleSubmit() {
  if (!isFormValid.value || isSubmitting.value) return;
  isSubmitting.value = true;
  try {
    const data = {
      partyId: formData.value.partyId,
      deliveryDate: salesApi.formatDateForAPI(new Date(formData.value.deliveryDate)),
      items: formData.value.lineItems.map(i => ({
        productVariantId: i.productVariantId, quantity: i.quantity,
        unitPrice: { amount: i.unitPrice.toString(), currency: 'EUR' },
        discountPercent: i.discountPercent || 0
      })),
      notes: formData.value.notes,
      mesWorkRefs: formData.value.mesWorkRefs
    };
    const order = await salesApi.createOrder(data);
    toastStore.success('Pedido creado correctamente');
    router.push(`/sales/orders/${order.id}`);
  } catch (err: any) { 
    toastStore.error(err.message || 'Error al crear el pedido'); 
  }
  finally { isSubmitting.value = false; }
}

const goBack = () => router.back();
const buildDisplayName = (i) => i.productName + (i.optionConfiguration ? ' (' + Object.values(i.optionConfiguration).join(', ') + ')' : '');
const formatMoney = (a) => salesApi.formatMoney({ amount: a?.toString() || '0', currency: 'EUR' });
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.page-content { max-width: 1300px; margin: 0 auto; padding: 1.5rem; }

.transactional-layout { display: grid; grid-template-columns: 1fr 380px; gap: 1.5rem; align-items: start; }

.section-header { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 1rem; padding-bottom: 0.6rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.85rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-secondary); }
.section-header :deep(svg) { color: var(--color-secondary); }

/* Table specific */
.table-header-row { display: flex; justify-content: space-between; align-items: center; padding: 1rem 1.5rem; }
.data-table th { background: var(--color-background-soft); font-size: 0.7rem; font-weight: 800; color: var(--color-text-secondary); text-transform: uppercase; }
.variant-search-cell { min-width: 280px; }
.variant-active-tag { display: flex; align-items: center; gap: 0.6rem; padding: 0.4rem 0.75rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; cursor: pointer; transition: 0.2s; }
.variant-active-tag:hover { background: white; border-color: var(--color-error); }
.variant-active-tag .sku { font-weight: 800; color: #166534; font-family: var(--font-family-mono); font-size: 0.8rem; }
.variant-active-tag .name { font-size: 1rem; color: #1e293b; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; flex: 1; }
.variant-active-tag .remove-icon { color: #94a3b8; }

.qty-input { width: 70px !important; text-align: center; font-weight: 700; }
.price-input { width: 100px !important; text-align: right; font-weight: 700; color: var(--color-secondary); }
.subtotal-text { color: var(--color-text-primary); font-size: 0.95rem; }

/* Totals */
.totals-body { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.75rem; }
.total-line { display: flex; justify-content: space-between; font-size: 0.9rem; color: var(--color-text-secondary); }
.total-line.final { margin-top: 1rem; padding-top: 1rem; border-top: 2px solid var(--color-background); color: var(--color-text-primary); font-size: 1.1rem; }
.total-line.final strong { color: var(--color-primary); font-size: 1.5rem; font-weight: 900; }

/* MES Sidebar */
.mes-item-row { display: flex; gap: 0.5rem; padding: 0.75rem; background: var(--color-background-soft); border-radius: 8px; margin-bottom: 0.5rem; border: 1px solid var(--color-border); }
.mes-item-info { flex: 1; display: flex; flex-direction: column; gap: 0.4rem; }
.mes-select, .mes-input { padding: 0.4rem 0.6rem !important; font-size: 0.8rem !important; border-radius: 6px !important; }

.btn-inside { position: absolute; right: 0.5rem; top: 50%; transform: translateY(-50%); background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); }
.btn-icon { background: transparent; border: none; cursor: pointer; padding: 0.4rem; border-radius: 6px; }
.btn-icon:hover { background: rgba(0,0,0,0.05); }
.w-full { width: 100%; }

@media (max-width: 1180px) {
  .transactional-layout { grid-template-columns: 1fr; }
  .side-summary-area { order: -1; }
}

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>