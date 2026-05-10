<template>
  <BaseFormLayout
    title="Nuevo Presupuesto"
    :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Crear Presupuesto' }]"
    :is-submitting="isSubmitting"
    submit-text="Crear Presupuesto"
    submit-icon="add_card"
    catalog-route="/sales/quotes"
    @submit="handleSubmit"
  >
    <!-- Section: Customer -->
    <FormSection title="Selección de Cliente" icon="person" description="Identifica al cliente para el que se emite este presupuesto.">
      <PartySelector
        v-model="formData.partyId"
        label="Cliente *"
        placeholder="Buscar por nombre, NIF o referencia..."
        role-filter="CLIENT"
        :required="true"
        @select="onPartySelected"
      />
    </FormSection>

    <!-- Section: Order Info -->
    <FormSection title="Detalles del Presupuesto" icon="event_note">
      <div class="form-row">
        <div class="form-group">
          <label>Fecha de Emisión *</label>
          <input v-model="formData.quoteDate" type="date" class="form-input" required />
        </div>
        <div class="form-group">
          <label>Válido hasta (Expiración) *</label>
          <input v-model="formData.expirationDate" type="date" class="form-input" :min="formData.quoteDate" required />
        </div>
      </div>
      <div class="form-group mt-4">
        <label>Notas para el Cliente</label>
        <textarea v-model="formData.notes" class="form-textarea" rows="3" placeholder="Información sobre plazos, condiciones de pago, etc."></textarea>
      </div>
    </FormSection>

    <!-- Section: Line Items -->
    <FormSection title="Detalle de Líneas" icon="list_alt" description="Añade los productos o servicios que componen el presupuesto.">
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>Producto / Variante</th>
              <th class="text-center">Cant.</th>
              <th class="align-right">P. Unitario</th>
              <th class="text-center">Dto %</th>
              <th class="align-right">Subtotal</th>
              <th class="text-center">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in formData.lineItems" :key="index">
              <td>
                <div v-if="!item.productVariantId" class="input-with-action">
                  <input v-model="item.quickSearch" type="text" class="form-input-sm" placeholder="SKU o nombre..." @keyup.enter="openVariantSelector(index)" />
                  <button type="button" class="btn btn-outline btn-sm" @click="openVariantSelector(index)"><Search :size="16" /></button>
                </div>
                <div v-else class="selected-item-tag" @click="clearLine(index)">
                  <strong>{{ item.variantSku }}</strong> <span>{{ item.displayName }}</span>
                  <X :size="16" />
                </div>
              </td>
              <td class="text-center">
                <input 
                  v-model.number="item.quantity" 
                  type="number" 
                  min="1" 
                  class="form-input-sm w-16" 
                  :data-row="index"
                  data-col="qty"
                  @input="calculateTotals" 
                  @keydown="handleLineKeyDown($event, index, 'qty', item)"
                />
              </td>
              <td class="align-right">
                <input 
                  v-model.number="item.unitPrice" 
                  type="number" 
                  step="0.01" 
                  class="form-input-sm w-24 text-right" 
                  :data-row="index"
                  data-col="price"
                  @input="calculateTotals" 
                  @keydown="handleLineKeyDown($event, index, 'price', item)"
                />
              </td>
              <td class="text-center">
                <input 
                  v-model.number="item.discountPercent" 
                  type="number" 
                  step="0.01" 
                  class="form-input-sm w-16 text-center" 
                  :data-row="index"
                  data-col="disc"
                  @input="calculateTotals" 
                  @keydown="handleLineKeyDown($event, index, 'disc', item)"
                />
              </td>
              <td class="align-right"><strong>{{ formatMoney(calculateLineSubtotal(index)) }}</strong></td>
              <td class="text-center">
                <button type="button" class="btn-icon text-danger" @click="removeLine(index)"><Trash2 :size="18" /></button>
              </td>
            </tr>
            <tr v-if="formData.lineItems.length === 0">
              <td colspan="6" class="empty-row">Pulse el botón para añadir productos al presupuesto.</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="mt-4">
        <button type="button" class="btn btn-secondary btn-sm" @click="addLineItem">
          <Plus :size="16" /> Añadir Producto
        </button>
      </div>
    </FormSection>

    <!-- Totals Summary Card -->
    <div class="totals-checkout-layout">
      <div class="card totals-checkout-card">
        <div class="total-row">
          <label>Subtotal:</label>
          <span>{{ formatMoney(totals.subtotal) }}</span>
        </div>
        <div class="total-row">
          <label>IVA (21%):</label>
          <span>{{ formatMoney(totals.tax) }}</span>
        </div>
        <div class="total-row final">
          <label>TOTAL PRESUPUESTO:</label>
          <span class="total-value">{{ formatMoney(totals.total) }}</span>
        </div>
      </div>
    </div>
  </BaseFormLayout>

  <!-- Variant Selector Modal -->
  <Transition name="fade">
    <div v-if="showVariantSelector" class="modal-backdrop">
      <div class="modal card w-modal-xl">
        <div class="modal-header">
          <Package :size="24" />
          <h2>Seleccionar Producto</h2>
          <button class="btn-icon ml-auto" @click="showVariantSelector = false"><X :size="20" /></button>
        </div>
        <div class="modal-body overflow-y">
          <VariantSelector :initial-query="variantQuery" @variant-selected="handleVariantSelected" />
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRouter } from 'vue-router';
import { Search, X, Trash2, Plus, Package } from 'lucide-vue-next';
import BaseFormLayout from '@/components/shared/BaseFormLayout.vue';
import FormSection from '@/components/shared/FormSection.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import { useLineNavigation } from '@/composables/useLineNavigation';
import salesApi from '@/services/salesApi';

import { useToastStore } from '@/stores/toast';

const router = useRouter();
const toastStore = useToastStore();
const isSubmitting = ref(false);
const formData = reactive({
  partyId: '',
  quoteDate: new Date().toISOString().split('T')[0],
  expirationDate: new Date(Date.now() + 30*24*60*60*1000).toISOString().split('T')[0],
  notes: '',
  lineItems: []
});

const { handleLineKeyDown, focusLineInput } = useLineNavigation({
  rowCount: () => formData.lineItems.length,
  columns: ['qty', 'price', 'disc'],
  onUpdate: (index, col, val) => {
    const item = formData.lineItems[index];
    if (col === 'qty') item.quantity = val;
    else if (col === 'price') item.unitPrice = val;
    else if (col === 'disc') item.discountPercent = val;
    calculateTotals();
  },
  onLastFieldEnter: () => addLineItem(),
  onAddField: () => addLineItem()
});

onMounted(() => {
  window.addEventListener('tramatex-save', handleGlobalSave);
  window.addEventListener('tramatex-esc', () => router.push('/sales/quotes'));
});

onBeforeUnmount(() => {
  window.removeEventListener('tramatex-save', handleGlobalSave);
  window.removeEventListener('tramatex-esc', () => router.push('/sales/quotes'));
});

function handleGlobalSave() {
  if (!isSubmitting.value && formData.lineItems.length > 0) handleSubmit();
}
const partyDefaultDiscount = ref(null);
const showVariantSelector = ref(false);
const editingIdx = ref(null);
const variantQuery = ref('');

function onPartySelected(party) { partyDefaultDiscount.value = party?.default_discount_percentage || null; }

function addLineItem() {
  formData.lineItems.push({ productVariantId: '', variantSku: '', displayName: '', quantity: 1, unitPrice: 0, discountPercent: partyDefaultDiscount.value || 0, quickSearch: '' });
}

function removeLine(idx) { formData.lineItems.splice(idx, 1); calculateTotals(); }
function clearLine(idx) { formData.lineItems[idx] = { ...formData.lineItems[idx], productVariantId: '', variantSku: '', displayName: '' }; calculateTotals(); }

function openVariantSelector(idx) {
  editingIdx.value = idx; variantQuery.value = formData.lineItems[idx].quickSearch; showVariantSelector.value = true;
}

function handleVariantSelected(v) {
  const variant = v.variant || v;
  const item = formData.lineItems[editingIdx.value];
  item.productVariantId = variant.id;
  item.variantSku = variant.sku;
  item.displayName = variant.product_name + (variant.option_configuration ? ' - ' + Object.values(variant.option_configuration).join(', ') : '');
  item.unitPrice = null;
  item.listPrice = null;
  item._autoPrice = true;
  showVariantSelector.value = false;
  
  // Trigger immediate calculation to get the correct sale price with margins
  calculateTotals(true);
}

let previewTimer = null;
const isPreviewLoading = ref(false);

function calculateTotals(immediate = false) {
  clearTimeout(previewTimer);
  const delay = immediate ? 0 : 400;
  
  previewTimer = setTimeout(async () => {
    const items = formData.lineItems.filter(i => i.productVariantId).map(i => ({
      productVariantId: i.productVariantId, 
      quantity: Number(i.quantity || 0),
      ...(i._autoPrice === false ? { unitPrice: { amount: Number(i.unitPrice || 0), currency: 'EUR' } } : {}), 
      discountPercent: Number(i.discountPercent || 0)
    }));
    
    if (!formData.partyId || items.length === 0) { 
      Object.assign(totals, { subtotal: 0, tax: 0, total: 0 }); 
      return; 
    }
    
    isPreviewLoading.value = true;
    try {
      const res = await salesApi.previewQuoteCalculation(formData.partyId, items);
      if (res) {
        totals.subtotal = res.subtotal.amount; 
        totals.tax = res.taxAmount.amount; 
        totals.total = res.total.amount;
        
        // Update line prices for items that haven't been manually overridden
        formData.lineItems.forEach((item, idx) => {
          if (item.productVariantId && item._autoPrice !== false && res.lineItems?.[idx]) {
            item.unitPrice = res.lineItems[idx].unitPrice.amount;
          }
        });
      }
    } catch (err) {
      console.error("Error calculating preview:", err);
    } finally {
      isPreviewLoading.value = false;
    }
  }, delay);
}

function calculateLineSubtotal(idx) {
  const i = formData.lineItems[idx];
  return (i.unitPrice * i.quantity) * (1 - (i.discountPercent / 100));
}

async function handleSubmit() {
  if (formData.lineItems.length === 0) return;
  isSubmitting.value = true;
  try {
    const payload = {
      partyId: formData.partyId,
      notes: formData.notes,
      expirationDate: salesApi.formatDateForAPI(new Date(formData.expirationDate)),
      items: formData.lineItems.map(i => ({
        productVariantId: i.productVariantId, quantity: i.quantity,
        unitPrice: { amount: i.unitPrice, currency: 'EUR' }, discountPercent: i.discountPercent
      }))
    };
    const res = await salesApi.createQuote(payload);
    toastStore.success('Presupuesto creado correctamente');
    router.push(`/sales/quotes/${res.id}`);
  } catch (err: any) { 
    toastStore.error(err.message || 'Error al crear el presupuesto'); 
  }
  finally { isSubmitting.value = false; }
}

function formatMoney(a) { return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(a); }
</script>

<style scoped>
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); }

.input-with-action { display: flex; gap: 0.5rem; }
.selected-item-tag { display: flex; align-items: center; gap: 0.75rem; padding: 0.5rem 1rem; background: rgba(34, 197, 94, 0.1); border: 1px solid #86efac; border-radius: 8px; cursor: pointer; }

.totals-checkout-layout { display: flex; justify-content: flex-end; margin-top: 1rem; }
.totals-checkout-card { width: 450px; padding: 1.5rem 2rem; background: var(--color-background); border-top: 4px solid var(--color-primary); }
.total-row { display: flex; justify-content: space-between; margin-bottom: 0.75rem; font-size: 0.95rem; }
.total-row.final { margin-top: 1rem; padding-top: 1rem; border-top: 2px solid var(--color-border); font-weight: 800; font-size: 1.25rem; }
.total-value { color: var(--color-primary); }

.w-16 { width: 4rem; } .w-24 { width: 6rem; }
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.w-modal-xl { width: 90%; max-width: 1100px; }
.btn-icon { color: var(--color-text-secondary); cursor: pointer; }
</style>