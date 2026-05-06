<template>
  <BaseTerminalPage 
    title="Terminal de Venta Directa (Ticket)" 
    station-id="CAJA PRINCIPAL"
    :icon="Store"
    :is-loading="isSubmitting"
    @close="router.push('/sales/dashboard')"
    @refresh="initTicket"
  >
    <div class="ticket-create-layout">
      <!-- PANEL IZQUIERDO: BÚSQUEDA Y PRODUCTOS -->
      <section class="catalog-panel">
        <div class="search-bar">
          <div class="input-with-icon">
            <Search :size="24" />
            <input 
              ref="barcodeInput"
              v-model="searchQuery" 
              type="text" 
              placeholder="Escanear código de barras o buscar..." 
              class="tpv-input"
              @keyup.enter="handleBarcodeSearch"
            />
          </div>
          <button class="btn-clear-search" @click="searchQuery = ''" v-if="searchQuery"><X :size="20" /></button>
        </div>

        <div class="product-grid-container">
          <div v-if="isSearching" class="tpv-loader">
            <RefreshCw class="spin" :size="48" />
            <p>Buscando productos...</p>
          </div>

          <div v-else-if="catalog.length === 0" class="empty-catalog">
            <PackageSearch :size="64" />
            <p>No se encontraron productos.</p>
          </div>

          <div v-else class="tpv-product-grid">
            <div 
              v-for="product in catalog" 
              :key="product.id" 
              class="tpv-product-card"
              @click="handleProductClick(product)"
            >
              <div class="product-img">
                <Package :size="32" />
              </div>
              <div class="product-info">
                <span class="sku">{{ product.sku }}</span>
                <p class="name">{{ product.name }}</p>
                <span class="price">{{ formatMoney(product.price) }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- PANEL DERECHO: TICKET ACTUAL -->
      <aside class="ticket-panel">
        <div class="ticket-header">
          <div class="ticket-meta">
            <h2 class="flex items-center gap-2"><Receipt :size="20" /> Ticket Actual</h2>
            <span class="item-count">{{ lineCount }} artículos</span>
          </div>
          <button class="btn-tpv btn-tpv-outline text-danger" @click="promptClearTicket" :disabled="lines.length === 0">
            <Trash2 :size="18" /> Vaciar
          </button>
        </div>

        <div class="ticket-lines no-scrollbar">
          <div v-for="(line, index) in lines" :key="index" class="ticket-line">
            <div class="line-main">
              <span class="line-name">{{ line.name }}</span>
              <span class="line-sku">{{ line.sku }}</span>
            </div>
            <div class="line-controls">
              <div class="qty-stepper">
                <button @click="updateQty(index, -1)">-</button>
                <input type="number" v-model.number="line.quantity" readonly />
                <button @click="updateQty(index, 1)">+</button>
              </div>
              <span class="line-total">{{ formatMoney(line.total) }}</span>
              <button class="btn-remove-line" @click="removeLine(index)"><X :size="16" /></button>
            </div>
          </div>
          
          <div v-if="lines.length === 0" class="empty-ticket">
            <ShoppingCart :size="48" />
            <p>Ticket vacío. Empieza a añadir productos.</p>
          </div>
        </div>

        <div class="ticket-footer">
          <div class="summary-area">
            <div class="summary-row">
              <label>Subtotal</label>
              <span>{{ formatMoney(subtotal) }}</span>
            </div>
            <div class="summary-row">
              <label>IVA (21%)</label>
              <span>{{ formatMoney(tax) }}</span>
            </div>
            <div class="summary-total">
              <label>TOTAL A PAGAR</label>
              <span class="amount">{{ formatMoney(total) }}</span>
            </div>
          </div>

          <div class="payment-actions">
            <button 
              class="btn-tpv btn-tpv-primary btn-checkout" 
              :disabled="lines.length === 0 || isSubmitting"
              @click="processTicket"
            >
              <CreditCard :size="24" />
              <span>COBRAR TICKET (F10)</span>
            </button>
          </div>
        </div>
      </aside>
    </div>

    <!-- MODAL: SELECCIÓN DE VARIANTE (Si el producto tiene varias) -->
    <BaseDialog 
      :show="showVariantModal" 
      title="Seleccionar Variante" 
      size="md" 
      hide-actions
      @close="showVariantModal = false"
    >
      <div class="variant-selection-list">
        <button 
          v-for="v in activeVariants" 
          :key="v.id" 
          class="variant-btn"
          @click="addVariantToTicket(v)"
        >
          <div class="v-info">
            <strong>{{ v.sku }}</strong>
            <span>{{ formatVariantOptions(v.option_configuration) }}</span>
          </div>
          <span class="v-price">{{ formatMoney(v.price || selectedProduct?.price) }}</span>
        </button>
      </div>
    </BaseDialog>

    <!-- MODAL DE CONFIRMACIÓN -->
    <BaseDialog
      :show="confirmDialog.show"
      :title="confirmDialog.title"
      :icon="confirmDialog.icon"
      confirm-text="Sí, Vaciar Ticket"
      confirm-class="btn-danger"
      @close="confirmDialog.show = false"
      @confirm="confirmDialog.action(); confirmDialog.show = false"
    >
      <p>{{ confirmDialog.message }}</p>
    </BaseDialog>

    <!-- ÁREA DE IMPRESIÓN OCULTA (80mm) -->
    <div id="tpv-print-area" v-if="lastProcessedTicket">
      <div class="print-ticket">
        <div class="print-header">
          <h1>TramaTex ERP</h1>
          <p>CIF: B12345678</p>
          <p>Calle Falsa 123, 28001 Madrid</p>
          <p>Tlf: 910 000 000</p>
        </div>
        <div class="print-divider">--------------------------------</div>
        <div class="print-meta">
          <p>TICKET: {{ lastProcessedTicket.invoiceNumber }}</p>
          <p>FECHA: {{ new Date().toLocaleString() }}</p>
        </div>
        <div class="print-divider">--------------------------------</div>
        <div class="print-lines">
          <div v-for="l in lastProcessedTicket.lineItems" :key="l.id" class="print-line">
            <div class="l-name">{{ l.productName }} ({{ l.variantSku }})</div>
            <div class="l-math">
              <span>{{ l.quantity }} x {{ formatMoney(l.unitPrice) }}</span>
              <span class="l-sub">{{ formatMoney(l.total) }}</span>
            </div>
          </div>
        </div>
        <div class="print-divider">--------------------------------</div>
        <div class="print-totals">
          <div class="t-row"><span>SUBTOTAL:</span> <span>{{ formatMoney(lastProcessedTicket.subtotal) }}</span></div>
          <div class="t-row"><span>IVA (21%):</span> <span>{{ formatMoney(lastProcessedTicket.taxAmount) }}</span></div>
          <div class="t-row bold"><span>TOTAL:</span> <span>{{ formatMoney(lastProcessedTicket.total) }}</span></div>
        </div>
        <div class="print-footer">
          <p>¡Gracias por su compra!</p>
          <p>TramaTex - www.tramatex.com</p>
        </div>
      </div>
    </div>
  </BaseTerminalPage>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { 
  Store, 
  Search, 
  X, 
  RefreshCw, 
  PackageSearch, 
  Package, 
  Receipt, 
  Trash2, 
  ShoppingCart, 
  CreditCard,
  AlertTriangle 
} from 'lucide-vue-next';
import BaseTerminalPage from '@/components/shared/BaseTerminalPage.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import salesApi from '@/services/salesApi';
import { productApi } from '@/services/productApi';
import { useToastStore } from '@/stores/toast';

const router = useRouter();
const toastStore = useToastStore();

// --- STATE ---
const searchQuery = ref('');
const catalog = ref([]);
const lines = ref([]);
const isSearching = ref(false);
const isSubmitting = ref(false);
const barcodeInput = ref(null);

const showVariantModal = ref(false);
const selectedProduct = ref(null);
const activeVariants = ref([]);

const lastProcessedTicket = ref(null);

// --- Confirm Dialog Logic ---
const confirmDialog = reactive({
  show: false,
  title: '',
  message: '',
  icon: AlertTriangle,
  action: null
})

function promptClearTicket() {
  confirmDialog.title = 'Vaciar Ticket';
  confirmDialog.message = '¿Estás seguro de que deseas eliminar todos los artículos del ticket actual?';
  confirmDialog.action = clearTicket;
  confirmDialog.show = true;
}

// --- COMPUTED ---
const subtotal = computed(() => lines.value.reduce((acc, l) => acc + l.total, 0));
const tax = computed(() => subtotal.value * 0.21);
const total = computed(() => subtotal.value + tax.value);
const lineCount = computed(() => lines.value.reduce((acc, l) => acc + l.quantity, 0));

// --- METHODS ---
async function initTicket() {
  lines.value = [];
  searchQuery.value = '';
  lastProcessedTicket.value = null;
  await loadRecommendedProducts();
  focusInput();
}

async function loadRecommendedProducts() {
  try {
    const res = await productApi.listProducts({ limit: 12 });
    catalog.value = res.data || [];
  } catch (err) {}
}

async function handleBarcodeSearch() {
  const query = searchQuery.value.trim();
  if (!query) return;

  isSearching.value = true;
  try {
    const res = await productApi.searchVariants(query);
    if (res && res.length === 1) {
      addVariantToTicket(res[0]);
      searchQuery.value = '';
    } else if (res && res.length > 1) {
      activeVariants.value = res;
      showVariantModal.value = true;
    } else {
      toastStore.warning(`No se encontró el producto: ${query}`);
    }
  } catch (err) {
    toastStore.error("Error en la búsqueda.");
  } finally {
    isSearching.value = false;
    focusInput();
  }
}

async function handleProductClick(product) {
  selectedProduct.value = product;
  isSearching.value = true;
  try {
    const variants = await productApi.listVariants(product.id);
    if (variants.length === 1) {
      addVariantToTicket(variants[0]);
    } else {
      activeVariants.value = variants;
      showVariantModal.value = true;
    }
  } catch (err) {
    toastStore.error("Error al cargar variantes.");
  } finally {
    isSearching.value = false;
  }
}

function addVariantToTicket(variant) {
  const existingIdx = lines.value.findIndex(l => l.variantId === variant.id);
  if (existingIdx !== -1) {
    updateQty(existingIdx, 1);
  } else {
    lines.value.push({
      variantId: variant.id,
      sku: variant.sku,
      name: variant.product_name || selectedProduct.value?.name || 'Producto',
      price: variant.price || selectedProduct.value?.price || 0,
      quantity: 1,
      total: variant.price || selectedProduct.value?.price || 0
    });
  }
  showVariantModal.value = false;
  toastStore.success(`${variant.sku} añadido`, 1000);
}

function updateQty(idx, delta) {
  const line = lines.value[idx];
  line.quantity += delta;
  if (line.quantity <= 0) {
    removeLine(idx);
  } else {
    line.total = line.quantity * line.price;
  }
}

function removeLine(idx) {
  lines.value.splice(idx, 1);
}

function clearTicket() {
  lines.value = [];
}

async function processTicket() {
  if (lines.value.length === 0) return;
  
  isSubmitting.value = true;
  try {
    const payload = {
      type: 'SIMPLIFIED',
      items: lines.value.map(l => ({
        productVariantId: l.variantId,
        quantity: l.quantity,
        unitPrice: { amount: l.price, currency: 'EUR' }
      }))
    };
    
    const invoice = await salesApi.createInvoice(payload);
    lastProcessedTicket.value = invoice;
    
    toastStore.success("Ticket cobrado con éxito");
    
    // Auto-print
    await nextTick();
    window.print();
    
    // Reset for the next customer
    initTicket();
    
  } catch (err) {
    toastStore.error(err.message || "Error al procesar el cobro.");
  } finally {
    isSubmitting.value = false;
  }
}

// --- UTILS ---
const formatMoney = (v) => salesApi.formatMoney(v);
const formatVariantOptions = (opts) => opts ? Object.values(opts).join(' / ') : 'Estándar';
const focusInput = () => barcodeInput.value?.focus();

// --- KEYBOARD SHORTCUTS ---
const handleGlobalKeys = (e) => {
  if (e.key === 'F10') { 
    e.preventDefault(); 
    processTicket(); 
  }
  if (e.key === 'Escape') { 
    searchQuery.value = ''; 
    focusInput(); 
  }
  
  // Atajos rápidos para cantidad en el ticket (+/-)
  // Solo si no estamos escribiendo en el buscador
  if (document.activeElement !== barcodeInput.value && lines.value.length > 0) {
    if (e.key === '+' || e.key === '=') {
      e.preventDefault();
      updateQty(lines.value.length - 1, 1);
    }
    if (e.key === '-') {
      e.preventDefault();
      updateQty(lines.value.length - 1, -1);
    }
  }
};

onMounted(() => {
  initTicket();
  window.addEventListener('keydown', handleGlobalKeys);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeys);
});
</script>

<style scoped>
/* Estilos industriales táctiles */
.ticket-create-layout { display: grid; grid-template-columns: 1fr 450px; height: 100%; gap: 0; }

/* Catalog Panel */
.catalog-panel { display: flex; flex-direction: column; background: #0f172a; border-right: 2px solid #1e293b; }
.search-bar { padding: 1.5rem; background: #1e293b; display: flex; gap: 1rem; }
.input-with-icon { flex: 1; position: relative; }
.input-with-icon :deep(svg) { position: absolute; left: 1rem; top: 50%; transform: translateY(-50%); color: #64748b; }
.tpv-input { width: 100%; background: #0f172a; border: 2px solid #334155; color: white; padding: 1rem 1rem 1rem 3.5rem; border-radius: 12px; font-size: 1.25rem; font-weight: 700; }
.tpv-input:focus { outline: none; border-color: var(--color-primary); box-shadow: 0 0 0 4px rgba(230, 184, 0, 0.2); }

.product-grid-container { flex: 1; padding: 1.5rem; overflow-y: auto; }
.tpv-product-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 1rem; }
.tpv-product-card { background: #1e293b; border-radius: 12px; padding: 1rem; border: 2px solid transparent; cursor: pointer; transition: 0.2s; display: flex; flex-direction: column; align-items: center; text-align: center; }
.tpv-product-card:hover { border-color: var(--color-primary); background: #2d3a4f; transform: translateY(-3px); }
.product-img { width: 80px; height: 80px; background: #0f172a; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin-bottom: 1rem; color: var(--color-primary); }
.tpv-product-card .sku { font-size: 0.7rem; font-weight: 800; color: #64748b; text-transform: uppercase; }
.tpv-product-card .name { font-size: 0.95rem; font-weight: 700; color: white; margin: 0.25rem 0; height: 2.8rem; overflow: hidden; }
.tpv-product-card .price { font-size: 1.2rem; font-weight: 900; color: var(--color-primary); }

/* Ticket Panel */
.ticket-panel { display: flex; flex-direction: column; background: #1e293b; }
.ticket-header { padding: 1.5rem; border-bottom: 2px solid #334155; display: flex; justify-content: space-between; align-items: center; }
.ticket-meta h2 { color: white; font-size: 1.1rem; text-transform: uppercase; font-weight: 800; margin: 0; }
.item-count { font-size: 0.8rem; color: var(--color-primary); font-weight: 700; }

.ticket-lines { flex: 1; overflow-y: auto; padding: 1rem; display: flex; flex-direction: column; gap: 0.5rem; }
.ticket-line { background: #0f172a; border-radius: 8px; padding: 1rem; display: flex; flex-direction: column; gap: 0.75rem; border: 1px solid #334155; }
.line-main { display: flex; flex-direction: column; }
.line-name { color: white; font-weight: 700; }
.line-sku { font-size: 0.75rem; color: #64748b; font-family: monospace; }
.line-controls { display: flex; align-items: center; justify-content: space-between; }

.qty-stepper { display: flex; align-items: center; background: #1e293b; border-radius: 8px; overflow: hidden; }
.qty-stepper button { width: 40px; height: 40px; border: none; background: transparent; color: white; font-size: 1.5rem; cursor: pointer; }
.qty-stepper button:active { background: #334155; }
.qty-stepper input { width: 50px; background: transparent; border: none; color: white; text-align: center; font-weight: 800; font-size: 1.1rem; }
.line-total { font-weight: 900; color: white; font-size: 1.1rem; }

.ticket-footer { padding: 1.5rem; background: #0f172a; border-top: 4px solid var(--color-primary); }
.summary-row { display: flex; justify-content: space-between; color: #94a3b8; font-weight: 600; margin-bottom: 0.5rem; }
.summary-total { display: flex; justify-content: space-between; align-items: flex-end; margin-top: 1rem; padding-top: 1rem; border-top: 1px solid #334155; }
.summary-total label { font-weight: 900; color: white; font-size: 0.9rem; }
.summary-total .amount { font-size: 2.5rem; font-weight: 900; color: var(--color-primary); line-height: 1; }

.payment-actions { margin-top: 1.5rem; }
.btn-checkout { width: 100%; height: 80px; font-size: 1.5rem; gap: 1rem; box-shadow: 0 10px 25px rgba(230, 184, 0, 0.2); }

.empty-ticket { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; color: #334155; text-align: center; padding: 2rem; }
.empty-ticket p { margin-top: 1rem; font-weight: 700; font-size: 1.1rem; }

/* Variants */
.variant-selection-list { display: flex; flex-direction: column; gap: 0.75rem; padding: 1rem; }
.variant-btn { display: flex; justify-content: space-between; align-items: center; padding: 1.25rem; background: #f8fafc; border: 2px solid #e2e8f0; border-radius: 12px; cursor: pointer; transition: 0.2s; }
.variant-btn:hover { border-color: var(--color-primary); background: white; }
.variant-btn .v-info { display: flex; flex-direction: column; text-align: left; }
.variant-btn .v-price { font-weight: 900; color: var(--color-secondary); font-size: 1.2rem; }

/* Utils */
.btn-tpv { display: flex; align-items: center; justify-content: center; border: none; border-radius: 12px; font-weight: 800; cursor: pointer; transition: 0.2s; }
.btn-tpv-primary { background: var(--color-primary); color: #0f172a; }
.btn-tpv-primary:hover { background: #facc15; transform: translateY(-2px); }
.btn-tpv-primary:active { transform: translateY(0); }
.btn-tpv-primary:disabled { background: #334155; color: #475569; cursor: not-allowed; transform: none; }
.btn-tpv-outline { background: transparent; border: 2px solid #334155; color: white; padding: 0.5rem 1rem; }
.btn-tpv-outline:hover { background: #334155; }
.btn-remove-line { background: transparent; border: none; color: #64748b; cursor: pointer; padding: 0.5rem; }
.btn-remove-line:hover { color: #ef4444; }

/* PRINT STYLES */
@media print {
  body * { visibility: hidden; }
  #tpv-print-area,
  #tpv-print-area * {
    visibility: visible !important;
  }

  #tpv-print-area {
    display: block !important;
    position: absolute;
    inset: 0;
    width: 80mm;
    margin: 0 auto;
    background: white;
  }
}
</style>