<template>
  <BaseTerminalPage 
    title="Terminal de Venta Directa (Ticket)" 
    station-id="CAJA PRINCIPAL"
    icon="point_of_sale"
    :is-loading="isSubmitting"
    @close="router.push('/sales/dashboard')"
    @refresh="clearTicket"
  >
    <!-- BOTON DE ATAJOS EN CABECERA -->
    <template #extra-actions>
      <div class="shortcut-header-group">
        <button 
          class="btn-terminal btn-sync btn-shortcuts"
          @click.stop="showShortcuts = !showShortcuts"
        >
          <span class="material-symbols-outlined">keyboard</span>
          <span>Atajos</span>
        </button>
        
        <!-- DESPLEGABLE DE ATAJOS -->
        <Teleport to="body">
          <div v-if="showShortcuts" class="hotkeys-modal-overlay" @click="showShortcuts = false">
            <div class="hotkeys-modal-box" @click.stop>
              <header class="modal-box-header">
                <div class="flex items-center gap-3">
                  <span class="material-symbols-outlined">keyboard</span>
                  <h3>ATAJOS DE TECLADO</h3>
                </div>
                <button class="btn-close-modal" @click="showShortcuts = false">×</button>
              </header>
              <div class="modal-box-body">
                <div class="hotkey-entry"><span>Buscar Producto</span> <kbd>F3</kbd></div>
                <div class="hotkey-entry"><span>Seleccionar Cliente</span> <kbd>F4</kbd></div>
                <div class="hotkey-entry"><span>Cobrar e Imprimir</span> <kbd>F12</kbd></div>
                <div class="hotkey-entry"><span>Limpiar / Cerrar</span> <kbd>Esc</kbd></div>
                <div class="hotkey-entry"><span>Más Cantidad</span> <kbd>Num +</kbd></div>
                <div class="hotkey-entry"><span>Menos Cantidad</span> <kbd>Num -</kbd></div>
                <div class="hotkey-entry"><span>Eliminar ítem</span> <kbd>Supr</kbd></div>
              </div>
            </div>
          </div>
        </Teleport>
      </div>
    </template>

    <!-- LAYOUT PRINCIPAL -->
    <div class="tpv-grid">
      <!-- IZQUIERDA: CARRITO -->
      <main class="tpv-main-content">
        <section class="tpv-card-search">
          <div class="search-flex">
            <span class="material-symbols-outlined">barcode_scanner</span>
            <input 
              ref="productSearchInput"
              v-model="productSearch" 
              type="text" 
              placeholder="Escanee código o busque artículo..." 
              @keyup.enter="handleSearch"
            />
            <button class="btn-add-item" @click="handleSearch">
              <span class="material-symbols-outlined">add</span>
              <span>AÑADIR</span>
            </button>
          </div>
        </section>

        <section class="tpv-card-cart">
          <header class="cart-title">
            <span class="material-symbols-outlined">shopping_cart</span>
            <h2>Artículos ({{ lineItems.length }})</h2>
          </header>
          
          <div class="cart-scroll-area">
            <div v-for="item in lineItems" :key="item.productVariantId" class="cart-item-row">
              <div class="item-info">
                <div class="meta">
                  <code class="sku">{{ item.variantSku }}</code>
                  <span class="name">{{ item.productName }}</span>
                  <span class="options">{{ item.optionDescription }}</span>
                </div>
              </div>
              
              <div class="item-controls">
                <div class="qty-control">
                  <button @click="updateQtyByItem(item, -1)">-</button>
                  <input v-model.number="item.quantity" type="number" @change="refreshLinePrice(item)" />
                  <button @click="updateQtyByItem(item, 1)">+</button>
                </div>

                <div class="price-unit">{{ salesApi.formatMoney({ amount: item.unitPrice, currency: 'EUR' }) }}</div>

                <!-- DESCUENTO EN LINEA -->
                <div class="discount-control">
                  <div class="input-inline">
                    <input v-model.number="item.discountPercent" type="number" min="0" max="100" />
                    <span class="symbol">% DTO.</span>
                  </div>
                </div>

                <div class="price-total">{{ salesApi.formatMoney({ amount: (item.unitPrice * item.quantity) * (1 - item.discountPercent / 100), currency: 'EUR' }) }}</div>

                <button class="btn-remove" @click="removeLine(item)">
                  <span class="material-symbols-outlined">close</span>
                </button>
              </div>
            </div>

            <div v-if="lineItems.length === 0" class="empty-cart">
              <span class="material-symbols-outlined">point_of_sale</span>
              <p>Esperando primer producto...</p>
            </div>
          </div>
        </section>
      </main>

      <!-- DERECHA: PANEL AZUL -->
      <aside class="tpv-side-panel">
        <div class="side-content">
          <section class="tpv-card-side">
            <header class="side-header">
              <span class="material-symbols-outlined">person</span>
              <h2>Cliente</h2>
            </header>
            <div class="p-4">
              <PartySelector
                v-model="partyId"
                label=""
                placeholder="Buscar cliente..."
                role-filter="CLIENT"
                :dark-mode="true"
              />
              <div class="client-detail mt-4" v-if="partyId !== CONSUMIDOR_FINAL_ID">
                <div class="row">
                  <label>Descuento cliente:</label>
                  <strong>{{ customerDiscount }}%</strong>
                </div>
              </div>
            </div>
          </section>
        </div>

        <!-- TOTALES FIJOS -->
        <section class="tpv-checkout-footer">
          <div class="checkout-summary">
            <div class="summary-line">
              <label>SUBTOTAL</label>
              <span>{{ salesApi.formatMoney({ amount: subtotal, currency: 'EUR' }) }}</span>
            </div>
            <div class="summary-line">
              <label>IVA (21%)</label>
              <span>{{ salesApi.formatMoney({ amount: taxAmount, currency: 'EUR' }) }}</span>
            </div>
            <div class="grand-total-box">
              <label>TOTAL TICKET</label>
              <div class="total-value">{{ salesApi.formatMoney({ amount: total, currency: 'EUR' }) }}</div>
            </div>
          </div>

          <div class="checkout-actions">
            <button 
              class="btn-checkout-primary" 
              :disabled="lineItems.length === 0 || isSubmitting" 
              @click="processTicket"
            >
              <span class="material-symbols-outlined">print</span>
              <span>COBRAR (F12)</span>
            </button>
            <button class="btn-checkout-secondary" @click="clearTicketPrompt">
              <span class="material-symbols-outlined">delete_sweep</span>
              <span>ANULAR TICKET</span>
            </button>
          </div>
        </section>
      </aside>
    </div>

    <!-- DIALOGO BUSQUEDA -->
    <BaseDialog :show="showVariantSelector" title="Buscador" icon="search" size="xl" hide-actions @close="closeVariantSelector">
      <div class="p-4">
        <VariantSelector :initial-query="productSearch" @variant-selected="handleVariantSelected" />
      </div>
    </BaseDialog>

    <!-- IMPRESION -->
    <Teleport to="body">
      <div id="tpv-print-area" v-if="lastProcessedTicket">
        <PrintTicket
          :number="lastProcessedTicket.number"
          :date="lastProcessedTicket.date"
          :items="lastProcessedTicket.items"
          :totals="lastProcessedTicket.totals"
          :customer-name="lastProcessedTicket.customerName"
        />
      </div>
    </Teleport>
  </BaseTerminalPage>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import BaseTerminalPage from '@/components/shared/BaseTerminalPage.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import PrintTicket from '@/components/sales/PrintTicket.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { productApi } from '@/services/productApi';
import { pricingApi } from '@/services/pricingApi';

const router = useRouter();
const productSearchInput = ref(null);
const productSearch = ref("");
const CONSUMIDOR_FINAL_ID = "00000000-0000-0000-0000-000000000001";
const partyId = ref(CONSUMIDOR_FINAL_ID);
const customerDiscount = ref(0);
const lineItems = ref([]);
const isSubmitting = ref(false);
const showShortcuts = ref(false);
const showVariantSelector = ref(false);
const lastProcessedTicket = ref(null);

const subtotal = computed(() => lineItems.value.reduce((acc, item) => {
  const lineVal = (item.unitPrice * item.quantity);
  const discount = lineVal * (item.discountPercent / 100);
  return acc + (lineVal - discount);
}, 0));
const taxAmount = computed(() => subtotal.value * 0.21);
const total = computed(() => subtotal.value + taxAmount.value);

watch(partyId, async (newId) => {
  if (!newId) return;
  try {
    const party = await partyApi.getParty(newId);
    customerDiscount.value = party?.default_discount_percentage || 0;
    lineItems.value.forEach(item => { item.discountPercent = customerDiscount.value; });
  } catch (err) { console.error("Error loading party:", err); }
  focusSearch();
});

function focusSearch() { nextTick(() => productSearchInput.value?.focus()); }

function handleGlobalKeydown(e) {
  const isInput = ["INPUT", "TEXTAREA"].includes(document.activeElement.tagName);
  if (e.key === "F3") { e.preventDefault(); focusSearch(); }
  if (e.key === "F4") { e.preventDefault(); document.querySelector(".party-selector input")?.focus(); }
  if (e.key === "F12") { e.preventDefault(); if (lineItems.value.length > 0 && !isSubmitting.value) processTicket(); }
  if (e.key === "Escape") {
    if (showShortcuts.value) showShortcuts.value = false;
    else if (showVariantSelector.value) closeVariantSelector();
    else { productSearch.value = ""; focusSearch(); }
  }
}

onMounted(async () => {
  window.addEventListener("keydown", handleGlobalKeydown);
  await loadDefaultCustomer();
  focusSearch();
});
onBeforeUnmount(() => window.removeEventListener("keydown", handleGlobalKeydown));

async function loadDefaultCustomer() {
  try {
    const p = await partyApi.getParty(CONSUMIDOR_FINAL_ID);
    if (p) { partyId.value = p.id; customerDiscount.value = p.default_discount_percentage || 0; }
  } catch (e) { console.error(e); }
}

async function handleSearch() {
  const q = productSearch.value.trim();
  if (!q) { showVariantSelector.value = true; return; }
  try {
    const res = await productApi.smartSearch(q);
    if (res.type === "exact_variant") handleVariantSelected(res.variant);
    else showVariantSelector.value = true;
  } catch (e) { showVariantSelector.value = true; }
}

async function handleVariantSelected(v) {
  const variant = v.variant || v;
  const existing = lineItems.value.find(i => i.productVariantId === variant.id);
  if (existing) {
    existing.quantity++;
    await refreshLinePrice(existing);
  } else {
    const item = {
      productVariantId: variant.id, 
      variantSku: variant.sku,
      productName: variant.product_name || variant.productName || "Producto",
      optionDescription: "",
      quantity: 1,
      unitPrice: 0, 
      discountPercent: customerDiscount.value
    };
    lineItems.value.push(item);
    await refreshLinePrice(item);
  }
  productSearch.value = ""; showVariantSelector.value = false; focusSearch();
}

function closeVariantSelector() { showVariantSelector.value = false; focusSearch(); }

async function refreshLinePrice(item) {
  try {
    const res = await pricingApi.calculateFinalSalePrice([{ productVariantId: item.productVariantId, quantity: item.quantity }], partyId.value, new Date());
    if (res.calculatedItems?.length) {
      item.unitPrice = (res.calculatedItems[0].baseSalesPrice?.amount ?? res.calculatedItems[0].baseSalesPrice) || 0;
    }
  } catch (e) { console.error(e); }
}

async function updateQtyByItem(item, delta) { item.quantity = Math.max(1, item.quantity + delta); await refreshLinePrice(item); }
function removeLine(item) { lineItems.value = lineItems.value.filter(i => i.productVariantId !== item.productVariantId); }
function clearTicket() { lineItems.value = []; productSearch.value = ""; partyId.value = CONSUMIDOR_FINAL_ID; focusSearch(); }
function clearTicketPrompt() { if (confirm("Anular ticket actual?")) clearTicket(); }

async function processTicket() {
  if (!lineItems.value.length) return;
  isSubmitting.value = true;
  try {
    const req = {
      partyId: partyId.value,
      invoiceDate: new Date().toISOString(),
      type: "SIMPLIFIED",
      items: lineItems.value.map(i => ({ productVariantId: i.productVariantId, quantity: i.quantity, discountPercent: Number(i.discountPercent || 0) }))
    };
    const inv = await salesApi.createSimplifiedInvoice(req);
    const items = inv.lineItems || inv.line_items || [];
    lastProcessedTicket.value = {
      number: inv.invoice_number || inv.invoiceNumber,
      date: inv.issue_date || inv.invoiceDate,
      items: items.map(li => ({
        productName: li.productName || li.product_name,
        variantSku: li.variantSku || li.variant_sku, 
        quantity: li.quantity,
        unitPrice: (li.unitPrice?.amount ?? li.unit_price) || 0,
        discountPercent: li.discountPercent ?? li.discount_percentage,
        subtotal: (li.subtotal?.amount ?? li.subtotal) || 0
      })),
      totals: { 
        subtotal: (inv.subtotal?.amount ?? inv.subtotal) || 0, 
        taxAmount: (inv.taxAmount?.amount ?? inv.tax_total) || 0, 
        total: (inv.total?.amount ?? inv.total) || 0 
      },
      customerName: inv.partyName || inv.party_name || "CLIENTE"
    };
    await new Promise(r => setTimeout(r, 600));
    window.print();
    clearTicket();
  } catch (err) { alert("Error: " + err.message); }
  finally { isSubmitting.value = false; lastProcessedTicket.value = null; }
}
</script>

<style scoped>
/* RESET & TPV LAYOUT */
.tpv-grid { display: grid; grid-template-columns: 1fr 360px; gap: 1rem; height: 100%; overflow: hidden; }
.tpv-main-content, .tpv-side-panel { display: flex; flex-direction: column; height: 100%; overflow: hidden; }

/* PANEL AZUL DERECHO */
.tpv-side-panel { background: #1e293b; border-left: 2px solid #334155; }
.side-content { flex: 1; display: flex; flex-direction: column; padding: 1rem; }
.side-header { padding: 0.5rem 1rem; display: flex; align-items: center; gap: 0.75rem; border-bottom: 1px solid #334155; margin-bottom: 1rem; }
.side-header h2 { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: #94a3b8; margin: 0; }
.side-header .material-symbols-outlined { color: var(--color-primary); font-size: 1.1rem; }

/* BUSQUEDA */
.tpv-card-search { background: #1e293b; border-radius: 12px; border: 1px solid #334155; margin-bottom: 1rem; padding: 0.4rem; }
.search-flex { display: flex; align-items: center; gap: 0.75rem; padding: 0.2rem 0.75rem; }
.search-flex .material-symbols-outlined { font-size: 1.75rem; color: #475569; }
.search-flex input { flex: 1; background: #0f172a; border: 1px solid #334155; border-radius: 8px; padding: 0.6rem 1rem; font-size: 1.1rem; color: white; font-weight: 600; outline: none; }
.search-flex input:focus { border-color: var(--color-primary); }
.btn-add-item { height: 40px; padding: 0 1.25rem; border-radius: 8px; border: none; background: var(--color-primary); color: black; font-weight: 900; font-size: 0.85rem; cursor: pointer; display: flex; align-items: center; gap: 0.4rem; }

/* CARRITO */
.tpv-card-cart { flex: 1; background: #1e293b; border-radius: 12px; border: 1px solid #334155; display: flex; flex-direction: column; overflow: hidden; }
.cart-title { padding: 0.6rem 1rem; display: flex; align-items: center; gap: 0.75rem; border-bottom: 1px solid #334155; background: #0f172a; }
.cart-title h2 { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: #94a3b8; margin: 0; }
.cart-scroll-area { flex: 1; overflow-y: auto; padding: 0.5rem; display: flex; flex-direction: column; gap: 0.4rem; }

/* FILA COMPACTA */
.cart-item-row { background: #0f172a; padding: 0.55rem 0.85rem; border-radius: 8px; border: 1px solid #334155; display: flex; align-items: center; gap: 1rem; border-left: 4px solid var(--color-primary); }
.item-info { flex: 1; min-width: 260px; display: flex; align-items: center; }
.meta { display: flex; align-items: center; gap: 0.65rem; min-width: 0; flex-wrap: nowrap; }
.sku { font-size: 0.72rem; background: var(--color-primary); color: #000; padding: 0.15rem 0.45rem; border-radius: 4px; font-weight: 900; line-height: 1; flex-shrink: 0; }
.name { font-size: 1rem; font-weight: 800; color: white; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; line-height: 1.2; }
.options { font-size: 0.82rem; color: #94a3b8; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.item-controls { display: flex; align-items: center; gap: 1rem; }
.qty-control { display: flex; align-items: center; background: #1e293b; border-radius: 8px; border: 1px solid #334155; height: 34px; overflow: hidden; }
.qty-control button { width: 28px; height: 100%; border: none; background: #162033; color: white; font-size: 1rem; font-weight: 900; cursor: pointer; }
.qty-control input { width: 42px; text-align: center; background: transparent; border: none; color: #f8fafc; font-size: 0.98rem; font-weight: 900; }
input::-webkit-outer-spin-button, input::-webkit-inner-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }

.price-unit { font-size: 0.92rem; font-weight: 700; color: #cbd5e1; text-align: right; width: 82px; }

.discount-control { display: flex; align-items: center; }
.input-inline { display: flex; align-items: center; background: #1e293b; border: 1px solid #334155; border-radius: 8px; height: 34px; padding: 0 0.5rem; width: 92px; position: relative; }
.input-inline input { width: 100%; background: transparent; border: none; color: var(--color-primary); font-size: 0.92rem; font-weight: 900; text-align: right; padding-right: 1.9rem; outline: none; }
.input-inline .symbol { position: absolute; right: 0.45rem; font-weight: 900; color: var(--color-primary); font-size: 0.66rem; pointer-events: none; }

.price-total { font-size: 1rem; font-weight: 900; color: #16a34a; text-align: right; width: 96px; }
.btn-remove { background: transparent; border: none; color: #475569; cursor: pointer; }
.btn-remove:hover { color: #ef4444; }

/* CHECKOUT */
.tpv-checkout-footer { background: #0f172a; border-top: 2px solid #334155; flex-shrink: 0; }
.checkout-summary { padding: 0.75rem 1.25rem; display: flex; flex-direction: column; gap: 0.25rem; }
.summary-line { display: flex; justify-content: space-between; font-size: 0.8rem; color: #94a3b8; font-weight: 600; }
.grand-total-box { margin-top: 0.4rem; padding-top: 0.4rem; border-top: 1px solid #334155; display: flex; flex-direction: column; align-items: center; }
.grand-total-box label { font-size: 0.55rem; font-weight: 900; color: #64748b; letter-spacing: 0.1em; }
.total-value { font-size: 1.75rem; font-weight: 950; color: var(--color-primary); line-height: 1; }

.checkout-actions { padding: 0.5rem 1.25rem 1rem; }
.btn-checkout-primary { width: 100%; height: 50px; border-radius: 10px; border: none; background: #16a34a; color: white; font-size: 1rem; font-weight: 950; display: flex; align-items: center; justify-content: center; gap: 0.6rem; cursor: pointer; }
.btn-checkout-secondary { width: 100%; height: 36px; border-radius: 8px; border: 1px solid #334155; background: transparent; color: #ef4444; font-size: 0.75rem; font-weight: 800; margin-top: 0.6rem; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 0.4rem; }

/* MODAL ATAJOS */
.hotkeys-modal-overlay { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.7); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: 99999; }
.hotkeys-modal-box { width: 400px; background: #1e293b; border: 2px solid var(--color-primary); border-radius: 16px; box-shadow: 0 30px 100px rgba(0,0,0,0.8); overflow: hidden; }
.modal-box-header { padding: 1rem 1.25rem; background: #0f172a; border-bottom: 1px solid #334155; display: flex; justify-content: space-between; align-items: center; }
.modal-box-header h3 { margin: 0; font-size: 0.9rem; font-weight: 900; color: white; }
.btn-close-modal { background: transparent; border: none; color: #94a3b8; font-size: 1.75rem; line-height: 1; cursor: pointer; }
.modal-box-body { padding: 1rem; display: flex; flex-direction: column; gap: 0.75rem; }
.hotkey-entry { display: flex; justify-content: space-between; align-items: center; padding-bottom: 0.5rem; border-bottom: 1px solid rgba(51, 65, 85, 0.4); }
.hotkey-entry span { font-size: 0.85rem; font-weight: 700; color: #cbd5e1; }
.hotkey-entry kbd { background: #0f172a; color: var(--color-primary); padding: 0.3rem 0.6rem; border-radius: 6px; font-family: monospace; font-weight: 900; font-size: 0.85rem; border: 1px solid #334155; min-width: 50px; text-align: center; }

/* UTILS */
.empty-cart { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; opacity: 0.2; gap: 0.5rem; }
.mt-4 { margin-top: 1rem; }
.p-4 { padding: 1rem; }
.bg-dark-alt { background: #0f172a; }
.rounded-lg { border-radius: 8px; }

#tpv-print-area { display: none; }

:deep(.btn-terminal.btn-sync) { background: #334155 !important; color: #cbd5e1 !important; height: 40px !important; border: 1px solid #475569 !important; }
:deep(.btn-terminal.btn-sync:hover) { background: #475569 !important; color: white !important; border-color: var(--color-primary) !important; }
:deep(.btn-terminal.btn-sync.btn-shortcuts) { padding: 0.75rem 1.5rem !important; border-radius: 12px !important; }

@media print {
  body * {
    visibility: hidden !important;
  }

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
