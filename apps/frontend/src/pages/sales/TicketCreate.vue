<template>
  <Navbar class="no-print" />
  
  <BaseDashboardPage :is-loading="isLoading" class="no-print">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <PageHeader 
        title="Terminal de Venta Directa (Ticket)" 
        :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Nuevo Ticket' }]"
      >
        <template #icon>
          <span class="material-symbols-outlined">point_of_sale</span>
        </template>
        <template #actions>
          <button class="btn btn-outline" @click="router.push('/sales/orders')">
            <span class="material-symbols-outlined">list_alt</span>
            <span>Ver Pedidos</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <!-- CAPA 3: TRABAJO (Área Principal) -->
    <div class="ticket-main-area">
      <section class="card mb-6">
        <div class="section-header">
          <span class="material-symbols-outlined">search</span>
          <h2>Búsqueda de Productos</h2>
        </div>
        <div class="section-body">
          <div class="input-with-icon">
            <span class="material-symbols-outlined icon-start">barcode_scanner</span>
            <input 
              v-model="productSearch" 
              type="text" 
              class="form-input-lg" 
              placeholder="Escanee código de barras o escriba SKU/Nombre..." 
              @keyup.enter="handleSearch"
            />
            <button class="btn btn-primary" @click="handleSearch">
              <span class="material-symbols-outlined">add</span>
              <span>Añadir</span>
            </button>
          </div>
        </div>
      </section>

      <section class="card table-card">
        <div class="table-header-row">
          <div class="section-header">
            <span class="material-symbols-outlined">shopping_cart</span>
            <h2>Productos en el Ticket</h2>
          </div>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Producto</th>
                <th class="text-center">Cant.</th>
                <th class="align-right">P. Unitario</th>
                <th class="align-right">Subtotal</th>
                <th class="text-center">Eliminar</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in lineItems" :key="idx">
                <td>
                  <strong>{{ item.productName }}</strong>
                  <div class="text-muted text-xs">{{ item.variantSku }}</div>
                </td>
                <td class="text-center">
                  <div class="qty-control">
                    <button class="btn-qty" @click="updateQty(idx, -1)">-</button>
                    <input v-model.number="item.quantity" type="number" class="qty-input" />
                    <button class="btn-qty" @click="updateQty(idx, 1)">+</button>
                  </div>
                </td>
                <td class="align-right">{{ salesApi.formatMoney({ amount: item.unitPrice, currency: 'EUR' }) }}</td>
                <td class="align-right"><strong>{{ salesApi.formatMoney({ amount: item.unitPrice * item.quantity, currency: 'EUR' }) }}</strong></td>
                <td class="text-center">
                  <button class="btn-icon text-danger" @click="removeLine(idx)"><span class="material-symbols-outlined">delete</span></button>
                </td>
              </tr>
              <tr v-if="lineItems.length === 0">
                <td colspan="5" class="empty-row-msg">El ticket está vacío. Busque un producto para comenzar.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <!-- CAPA 2: CONTEXTO (Panel Lateral) -->
    <template #sidebar>
      <div class="ticket-sidebar-content">
        <section class="mb-8">
          <div class="section-header">
            <span class="material-symbols-outlined">person</span>
            <h2>Cliente</h2>
          </div>
          <div class="mt-4">
            <PartySelector
              v-model="partyId"
              label=""
              placeholder="Buscar cliente..."
              role-filter="CLIENT"
              :required="false"
            />
          </div>
        </section>

        <section class="totals-checkout-panel">
          <div class="checkout-header">
            <span class="material-symbols-outlined">payments</span>
            <h3>Resumen de Venta</h3>
          </div>
          <div class="checkout-body mt-6">
            <div class="total-line">
              <label>Subtotal:</label>
              <span>{{ salesApi.formatMoney({ amount: subtotal, currency: 'EUR' }) }}</span>
            </div>
            <div class="total-line">
              <label>IVA (21%):</label>
              <span>{{ salesApi.formatMoney({ amount: taxAmount, currency: 'EUR' }) }}</span>
            </div>
            <div class="total-line final">
              <label>TOTAL A COBRAR:</label>
              <span class="total-value">{{ salesApi.formatMoney({ amount: total, currency: 'EUR' }) }}</span>
            </div>
          </div>
          <div class="checkout-footer mt-10">
            <button class="btn btn-secondary btn-lg w-full" :disabled="lineItems.length === 0 || isSubmitting" @click="processTicket">
              <span class="material-symbols-outlined">{{ isSubmitting ? 'sync' : 'print' }}</span>
              <span>{{ isSubmitting ? 'PROCESANDO...' : 'COBRAR E IMPRIMIR' }}</span>
            </button>
            <button class="btn btn-outline btn-sm w-full mt-4" @click="clearTicket">Limpiar Ticket</button>
          </div>
        </section>
      </div>
    </template>
  </BaseDashboardPage>

  <!-- MODAL DE PRODUCTOS -->
  <BaseDialog
    :show="showVariantSelector"
    title="Catálogo de Productos"
    icon="inventory_2"
    size="xl"
    hide-actions
    @close="showVariantSelector = false"
  >
    <VariantSelector :initial-query="productSearch" @variant-selected="handleVariantSelected" />
  </BaseDialog>

  <!-- CAPA DE IMPRESIÓN -->
  <div class="print-ticket-container" v-if="lastProcessedTicket">
    <PrintTicket
      :number="lastProcessedTicket.number"
      :date="lastProcessedTicket.date"
      :items="lastProcessedTicket.items"
      :totals="lastProcessedTicket.totals"
      :customer-name="lastProcessedTicket.customerName"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import PrintTicket from '@/components/sales/PrintTicket.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { productApi } from '@/services/productApi';

const router = useRouter();
const productSearch = ref('');
const partyId = ref('');
const lineItems = ref([]);
const isSubmitting = ref(false);
const isLoading = ref(true);
const showVariantSelector = ref(false);
const lastProcessedTicket = ref(null);

const subtotal = computed(() => lineItems.value.reduce((acc, item) => acc + (item.unitPrice * item.quantity), 0));
const taxAmount = computed(() => subtotal.value * 0.21);
const total = computed(() => subtotal.value + taxAmount.value);

onMounted(async () => {
  await loadDefaultCustomer();
  isLoading.value = false;
});

async function loadDefaultCustomer() {
  try {
    const res = await partyApi.listParties({ searchText: 'CONSUMIDOR FINAL', limit: 1 });
    const parties = res.data || (Array.isArray(res) ? res : []);
    if (parties.length > 0) partyId.value = parties[0].id;
  } catch (err) {}
}

async function handleSearch() {
  const query = productSearch.value.trim();
  if (!query) { showVariantSelector.value = true; return; }
  try {
    const result = await productApi.smartSearch(query);
    if (result.type === 'exact_variant') handleVariantSelected(result.variant);
    else showVariantSelector.value = true;
  } catch (err) { showVariantSelector.value = true; }
}

function handleVariantSelected(v) {
  const variant = v.variant || v;
  const existing = lineItems.value.find(i => i.productVariantId === variant.id);
  if (existing) existing.quantity++;
  else {
    lineItems.value.push({
      productVariantId: variant.id, variantSku: variant.sku,
      productName: variant.product_name || variant.productName || 'Producto', 
      quantity: 1,
      unitPrice: variant.product_base_price || variant.base_cost || variant.price || 0
    });
  }
  productSearch.value = '';
  showVariantSelector.value = false;
}

function updateQty(idx, delta) {
  const item = lineItems.value[idx];
  item.quantity = Math.max(1, item.quantity + delta);
}

function removeLine(idx) { lineItems.value.splice(idx, 1); }
function clearTicket() { if (confirm('¿Vaciar el ticket actual?')) lineItems.value = []; }

async function processTicket() {
  if (lineItems.value.length === 0) return;
  isSubmitting.value = true;
  try {
    lastProcessedTicket.value = {
      number: 'T-' + Math.floor(100000 + Math.random() * 900000),
      date: new Date().toISOString(),
      items: [...lineItems.value],
      totals: { subtotal: subtotal.value, taxAmount: taxAmount.value, total: total.value },
      customerName: 'CONSUMIDOR FINAL'
    };
    await new Promise(resolve => setTimeout(resolve, 500));
    window.print();
    alert('✓ Venta procesada correctamente');
    lineItems.value = [];
    productSearch.value = '';
  } catch (err) { alert('Error al procesar'); }
  finally { isSubmitting.value = false; }
}
</script>

<style scoped>
.section-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.25rem; }
.section-header h2 { font-size: 1rem; margin: 0; color: var(--color-text-primary); text-transform: uppercase; letter-spacing: 0.025em; }
.section-header .material-symbols-outlined { color: var(--color-primary); font-size: 22px; }
.section-body { padding: 0.5rem 0; }

/* Product Search */
.form-input-lg { flex: 1; padding: 1rem 1.5rem; border-radius: 8px; border: 1px solid var(--color-border); font-size: 1.1rem; }
.input-with-icon { display: flex; gap: 1rem; align-items: center; }

/* Table */
.table-card { padding: 0; overflow: hidden; }
.table-header-row { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--color-background); }
.empty-row-msg { padding: 4rem; text-align: center; color: var(--color-text-secondary); font-style: italic; }

.qty-control { display: flex; align-items: center; justify-content: center; gap: 0.5rem; }
.btn-qty { width: 32px; height: 32px; border-radius: 4px; border: 1px solid var(--color-border); background: white; cursor: pointer; font-weight: bold; }
.qty-input { width: 50px; text-align: center; border: none; font-weight: bold; background: transparent; }

/* Checkout Panel */
.checkout-header { display: flex; align-items: center; gap: 0.75rem; }
.checkout-header h3 { margin: 0; font-size: 1.25rem; }
.checkout-header .material-symbols-outlined { font-size: 28px; color: var(--color-secondary); }

.total-line { display: flex; justify-content: space-between; margin-bottom: 1.25rem; color: var(--color-text-secondary); }
.total-line.final { margin-top: 2rem; padding-top: 2rem; border-top: 2px solid var(--color-border); color: var(--color-text-primary); font-weight: 800; font-size: 1.25rem; }
.total-value { color: var(--color-secondary); font-size: 2rem; }

.w-full { width: 100%; }

/* ESTILOS DE IMPRESIÓN DE TICKET */
.print-ticket-container { display: none; }

@media print {
  .no-print { display: none !important; }
  .print-ticket-container { display: block !important; position: absolute; left: 0; top: 0; width: 80mm; }
}
</style>
