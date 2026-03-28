<template>
  <Navbar />
  
  <div class="main-container">
    <PageHeader 
      title="Nuevo Pedido de Venta" 
      :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Crear Nuevo' }]"
    >
      <template #actions>
        <button @click="router.push('/sales/orders')" class="btn btn-outline">
          <span class="material-symbols-outlined">list_alt</span>
          <span>Ir al catálogo</span>
        </button>
      </template>
    </PageHeader>

    <div class="order-form-layout">
      <form @submit.prevent="handleSubmit">
        <!-- Customer Section -->
        <section class="card mb-6">
          <div class="section-header">
            <span class="material-symbols-outlined">person</span>
            <h2>Selección de Cliente</h2>
          </div>
          <div class="section-body">
            <PartySelector
              v-model="formData.partyId"
              label="Cliente *"
              placeholder="Buscar por nombre, NIF o referencia..."
              role-filter="CLIENT"
              :required="true"
              @select="onPartySelected"
            />
          </div>
        </section>

        <!-- Order Details Section -->
        <section class="card mb-6">
          <div class="section-header">
            <span class="material-symbols-outlined">event_note</span>
            <h2>Detalles del Pedido</h2>
          </div>
          <div class="section-body">
            <div class="form-row">
              <div class="form-group">
                <label for="orderDate">Fecha de Pedido *</label>
                <input id="orderDate" v-model="formData.orderDate" type="date" class="form-input" required />
              </div>
              <div class="form-group">
                <label for="deliveryDate">Fecha de Entrega Estimada *</label>
                <input id="deliveryDate" v-model="formData.deliveryDate" type="date" class="form-input" :min="minDeliveryDate" required />
              </div>
            </div>
            <div class="form-group mt-4">
              <label for="notes">Observaciones del Pedido</label>
              <textarea id="notes" v-model="formData.notes" class="form-textarea" rows="3" placeholder="Instrucciones especiales para logística o taller..."></textarea>
            </div>
          </div>
        </section>

        <!-- MES Work Configurations -->
        <section v-if="formData.partyId" class="card mb-6 highlight-blue">
          <div class="section-header">
            <span class="material-symbols-outlined">precision_manufacturing</span>
            <h2>Configuraciones MES (Producción)</h2>
          </div>
          <div class="section-body">
            <div v-if="isLoadingMesWorks" class="loading-inline">
              <div class="spinner-tiny"></div> Cargando configuraciones...
            </div>
            <div v-else class="mes-ref-list">
              <div v-for="(config, idx) in formData.mesWorkRefs" :key="idx" class="mes-config-entry card mb-3">
                <div class="form-row">
                  <div class="form-group">
                    <label>Elegir configuración base</label>
                    <select class="form-input" :value="config.workSetupId || ''" @change="onSetupSelect(idx, $event.target.value)">
                      <option value="">— Configuración Personalizada —</option>
                      <option v-for="ws in mesWorkSetups" :key="ws.id" :value="ws.id">{{ ws.name }}</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label>Descripción del encargo</label>
                    <div class="input-with-action">
                      <input v-model="config.description" type="text" class="form-input" placeholder="Ej: Bordado logo espalda" />
                      <button type="button" class="btn btn-outline btn-sm text-danger" @click="removeConfig(idx)"><span class="material-symbols-outlined">delete</span></button>
                    </div>
                  </div>
                </div>
              </div>
              <button type="button" class="btn btn-secondary btn-sm" @click="addConfig">
                <span class="material-symbols-outlined">add</span> Añadir Configuración MES
              </button>
            </div>
          </div>
        </section>

        <!-- Line Items Section -->
        <section class="card table-card mb-6">
          <div class="table-header-row">
            <div class="section-header">
              <span class="material-symbols-outlined">shopping_basket</span>
              <h2>Líneas del Pedido</h2>
            </div>
            <button type="button" class="btn btn-secondary btn-sm" @click="addLineItem">
              <span class="material-symbols-outlined">add</span> Añadir Producto
            </button>
          </div>

          <div class="table-wrapper">
            <table class="data-table">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Referencia / Producto</th>
                  <th class="text-center">Cant.</th>
                  <th class="align-right">P. Tarifa</th>
                  <th class="align-right">P. Venta</th>
                  <th class="text-center">Dto %</th>
                  <th class="align-right">Subtotal</th>
                  <th class="text-center">Acciones</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in formData.lineItems" :key="index">
                  <td class="text-muted">{{ index + 1 }}</td>
                  <td>
                    <div class="variant-search-group">
                      <div v-if="!item.productVariantId" class="input-with-icon">
                        <span class="material-symbols-outlined icon-start">barcode_scanner</span>
                        <input v-model="item.quickSearchQuery" type="text" class="form-input-sm" placeholder="SKU o código..." @keyup.enter="inlineSmartSearch(index)" />
                        <button type="button" class="btn-icon-inside" @click="openVariantSelector(index)"><span class="material-symbols-outlined">search</span></button>
                      </div>
                      <div v-else class="variant-selected-tag" @click="clearLineVariant(index)">
                        <span class="variant-sku">{{ item.selectedVariantName }}</span>
                        <span class="variant-name">{{ buildDisplayName(item) }}</span>
                        <span class="material-symbols-outlined">close</span>
                      </div>
                      <small v-if="item.inlineSearchError" class="text-danger">{{ item.inlineSearchError }}</small>
                    </div>
                  </td>
                  <td class="text-center"><input v-model.number="item.quantity" type="number" min="1" class="form-input-sm w-16" @input="calculateTotals" /></td>
                  <td class="align-right text-muted">{{ item.listPrice != null ? item.listPrice.toFixed(2) : '—' }}</td>
                  <td class="align-right"><input v-model.number="item.unitPrice" type="number" step="0.01" class="form-input-sm w-20 text-right" @input="calculateTotals" /></td>
                  <td class="text-center"><input v-model.number="item.discountPercent" type="number" step="0.01" class="form-input-sm w-16 text-center" @input="calculateTotals" /></td>
                  <td class="align-right">
                    <span v-if="isPreviewLoading" class="spinner-tiny"></span>
                    <strong v-else>{{ formatMoney(calculateLineSubtotal(index)) }}</strong>
                  </td>
                  <td class="text-center">
                    <button type="button" class="btn-icon text-danger" @click="removeLineItem(index)"><span class="material-symbols-outlined">delete</span></button>
                  </td>
                </tr>
                <tr v-if="formData.lineItems.length === 0">
                  <td colspan="8" class="empty-row">No hay productos en el pedido. Pulse "Añadir Producto" para comenzar.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <!-- Totals & Actions -->
        <div class="form-footer-layout">
          <section class="card totals-summary-card">
            <div class="total-row">
              <label>Subtotal:</label>
              <span>{{ formatMoney(calculatedTotals.subtotal) }}</span>
            </div>
            <div class="total-row">
              <label>IVA (21%):</label>
              <span>{{ formatMoney(calculatedTotals.tax) }}</span>
            </div>
            <div class="total-row total-final">
              <label>TOTAL PEDIDO:</label>
              <span class="total-value">{{ formatMoney(calculatedTotals.total) }}</span>
            </div>
          </section>

          <div class="form-actions mt-6">
            <button type="button" class="btn btn-outline btn-lg" @click="goBack">Cancelar</button>
            <button type="submit" class="btn btn-primary btn-lg" :disabled="!isFormValid || isSubmitting">
              <span class="material-symbols-outlined">{{ isSubmitting ? 'sync' : 'shopping_cart_checkout' }}</span>
              <span>{{ isSubmitting ? "Procesando..." : "Confirmar y Crear Pedido" }}</span>
            </button>
          </div>
        </div>
      </form>

      <div v-if="submitError" class="alert alert-error mt-4">
        <span class="material-symbols-outlined">error</span>
        <p>{{ submitError }}</p>
      </div>
    </div>

    <!-- Variant Selector Modal -->
    <Transition name="fade">
      <div v-if="showVariantSelector" class="modal-backdrop">
        <div class="modal card w-modal-xl">
          <div class="modal-header">
            <span class="material-symbols-outlined">inventory_2</span>
            <h2>Seleccionar Producto</h2>
            <button class="btn-icon-inside ml-auto" @click="showVariantSelector = false"><span class="material-symbols-outlined">close</span></button>
          </div>
          <div class="modal-body overflow-y">
            <VariantSelector
              :key="variantSelectorQuery + '-' + editingLineIndex"
              :product-id="null"
              :initial-query="variantSelectorQuery"
              initial-mode="quick"
              title=""
              @variant-selected="handleVariantSelected"
            />
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import salesApi from '@/services/salesApi';
import { productApi } from '@/services/productApi';
import { calculateBaseSalesPrice } from '@/services/pricingApi';
import { mesApi } from '@/services/mesApi';

const router = useRouter();

const formData = ref({
  partyId: '',
  orderDate: new Date().toISOString().split('T')[0],
  deliveryDate: new Date().toISOString().split('T')[0],
  notes: '',
  mesWorkRefs: [],
  lineItems: [],
});

const isSubmitting = ref(false);
const submitError = ref('');
const mesWorkSetups = ref([]);
const isLoadingMesWorks = ref(false);
const partyDefaultDiscount = ref(null);

const calculatedTotals = ref({ subtotal: 0, tax: 0, total: 0 });
const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewDebounceTimer = null;

const minDeliveryDate = computed(() => {
  const tomorrow = new Date(); tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().split('T')[0];
});

const isFormValid = computed(() => (
  formData.value.partyId && formData.value.orderDate && 
  formData.value.lineItems.length > 0 && 
  formData.value.lineItems.every(i => i.productVariantId && i.quantity > 0)
));

watch(() => formData.value.partyId, (partyId) => { loadMesWorksForParty(partyId); });

function onPartySelected(party) { partyDefaultDiscount.value = party?.default_discount_percentage || null; }

const showVariantSelector = ref(false);
const editingLineIndex = ref(null);
const variantSelectorQuery = ref('');

function addLineItem() {
  formData.value.lineItems.push({
    productVariantId: '', productId: '', selectedVariantName: '', productName: '',
    optionConfiguration: {}, quickSearchQuery: '', inlineSearchError: '',
    quantity: 1, listPrice: null, unitPrice: null, discountPercent: partyDefaultDiscount.value || null,
  });
}

async function loadMesWorksForParty(partyId) {
  if (!partyId) { mesWorkSetups.value = []; formData.value.mesWorkRefs = []; return; }
  isLoadingMesWorks.value = true;
  try { mesWorkSetups.value = await mesApi.listWorkSetups({ party_id: partyId }); }
  catch (err) { mesWorkSetups.value = []; }
  finally { isLoadingMesWorks.value = false; }
}

function addConfig() { formData.value.mesWorkRefs.push({ workSetupId: null, description: '' }); }
function removeConfig(idx) { formData.value.mesWorkRefs.splice(idx, 1); }
function onSetupSelect(idx, setupId) { formData.value.mesWorkRefs[idx].workSetupId = setupId || null; }

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
    item.productId = variant.product_id || '';
    item.selectedVariantName = variant.sku;
    item.productName = variant.product_name || '';
    item.optionConfiguration = variant.option_configuration || {};
    item.listPrice = variant.product_base_price ?? variant.base_cost ?? null;
    item.unitPrice = item.listPrice;
    fetchPriceForLineItem(item);
  }
  showVariantSelector.value = false;
  calculateTotals();
}

async function fetchPriceForLineItem(item) {
  if (!item.productVariantId || !item.productId) return;
  try {
    const result = await calculateBaseSalesPrice(item.productId, item.productVariantId);
    item.listPrice = result.baseSalesPrice?.amount ?? item.listPrice;
    item.unitPrice = item.listPrice;
    calculateTotals();
  } catch (err) {}
}

async function inlineSmartSearch(index) {
  const item = formData.value.lineItems[index];
  const query = item.quickSearchQuery?.trim();
  if (!query) return;
  try {
    const result = await productApi.smartSearch(query);
    if (result.type === 'exact_variant' && result.variant) {
      handleVariantSelected(result.variant);
    } else {
      openVariantSelector(index);
    }
  } catch (err) { item.inlineSearchError = 'No encontrado'; }
}

function clearLineVariant(index) {
  formData.value.lineItems[index] = {
    ...formData.value.lineItems[index], productVariantId: '', productId: '', selectedVariantName: '', productName: '',
    listPrice: null, unitPrice: null, discountPercent: null
  };
  calculateTotals();
}

function removeLineItem(index) { formData.value.lineItems.splice(index, 1); calculateTotals(); }

function calculateLineSubtotal(i) {
  if (!previewResult.value) return 0;
  let pIdx = 0; for(let j=0; j<i; j++) if(formData.value.lineItems[j].productVariantId) pIdx++;
  return previewResult.value.lineItems[pIdx]?.subtotal?.amount ?? 0;
}

function calculateTotals() {
  clearTimeout(previewDebounceTimer);
  previewDebounceTimer = setTimeout(fetchPreviewCalculation, 400);
}

async function fetchPreviewCalculation() {
  const partyId = formData.value.partyId;
  const items = formData.value.lineItems.filter(i => i.productVariantId).map(i => ({
    productVariantId: i.productVariantId, quantity: i.quantity,
    unitPrice: i.unitPrice ? { amount: i.unitPrice, currency: 'EUR' } : undefined,
    discountPercent: i.discountPercent || undefined
  }));
  if (!partyId || items.length === 0) { calculatedTotals.value = { subtotal: 0, tax: 0, total: 0 }; return; }
  isPreviewLoading.value = true;
  try {
    previewResult.value = await salesApi.previewOrderCalculation(partyId, items);
    calculatedTotals.value = { subtotal: previewResult.value.subtotal.amount, tax: previewResult.value.taxAmount.amount, total: previewResult.value.total.amount };
  } catch (err) {}
  finally { isPreviewLoading.value = false; }
}

async function handleSubmit() {
  if (!isFormValid.value || isSubmitting.value) return;
  isSubmitting.value = true;
  try {
    const orderData = {
      partyId: formData.value.partyId,
      deliveryDate: salesApi.formatDateForAPI(new Date(formData.value.deliveryDate)),
      items: formData.value.lineItems.map(i => ({
        productVariantId: i.productVariantId, quantity: i.quantity,
        unitPrice: i.unitPrice ? { amount: i.unitPrice, currency: 'EUR' } : undefined,
        discountPercent: i.discountPercent || undefined
      })),
      notes: formData.value.notes,
      mesWorkRefs: formData.value.mesWorkRefs
    };
    const newOrder = await salesApi.createOrder(orderData);
    router.push(`/sales/orders/${newOrder.id}`);
  } catch (err) { submitError.value = err.message; }
  finally { isSubmitting.value = false; }
}

function goBack() { router.back(); }
function buildDisplayName(i) { return i.productName + (i.optionConfiguration ? ' - ' + Object.values(i.optionConfiguration).join(', ') : ''); }
function formatMoney(a) { return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(a); }
</script>

<style scoped>
.main-container { padding-bottom: 4rem; }

.section-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.25rem; }
.section-header h2 { font-size: 1rem; margin: 0; color: var(--color-text-primary); text-transform: uppercase; letter-spacing: 0.025em; }
.section-header .material-symbols-outlined { color: var(--color-text-secondary); font-size: 22px; }

.section-body { padding: 1.5rem; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 600; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }

input, select, textarea { width: 100%; padding: 0.75rem 1rem; border-radius: var(--border-radius-md); border: 1px solid var(--color-border); font-size: var(--font-size-sm); }

/* MES */
.highlight-blue { border-top: 4px solid var(--color-secondary); }
.mes-config-entry { padding: 1rem; background: var(--color-background); border: 1px solid var(--color-border); }

/* Table */
.table-card { padding: 0; overflow: hidden; }
.table-header-row { display: flex; justify-content: space-between; align-items: center; padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--color-background); }
.data-table th { background: var(--color-background); padding: 1rem; font-size: var(--font-size-xs); }
.data-table td { padding: 1rem; border-bottom: 1px solid var(--color-background); vertical-align: middle; }

/* Variant Search */
.variant-search-group { min-width: 300px; }
.variant-selected-tag { 
  display: flex; align-items: center; gap: 0.75rem; padding: 0.5rem 0.75rem; 
  background: rgba(34, 197, 94, 0.1); border: 1px solid #86efac; border-radius: 6px; cursor: pointer;
}
.variant-sku { font-family: var(--font-family-mono); font-weight: 700; color: #166534; font-size: 0.85rem; }
.variant-name { font-size: 0.85rem; flex: 1; }

.input-with-icon { position: relative; display: flex; align-items: center; }
.icon-start { position: absolute; left: 0.75rem; font-size: 20px; color: var(--color-text-secondary); }
.input-with-icon input { padding-left: 2.5rem; }
.btn-icon-inside { background: transparent; border: none; padding: 0.25rem; color: var(--color-text-secondary); cursor: pointer; }

/* Footer & Totals */
.form-footer-layout { display: flex; flex-direction: column; align-items: flex-end; }
.totals-summary-card { width: 400px; padding: 1.5rem; background: var(--color-background); }
.total-row { display: flex; justify-content: space-between; margin-bottom: 0.75rem; font-size: 0.9rem; }
.total-final { margin-top: 1rem; padding-top: 1rem; border-top: 2px solid var(--color-border); font-weight: 700; font-size: 1.25rem; }
.total-value { color: var(--color-primary); }

.w-16 { width: 4rem; } .w-20 { width: 5rem; } .w-24 { width: 6rem; }
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.w-modal-xl { width: 90%; max-width: 1000px; }
</style>