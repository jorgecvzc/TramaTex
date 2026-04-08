<template>
  <BaseTerminalPage 
    title="Terminal de Venta Directa (Ticket)" 
    station-id="CAJA PRINCIPAL"
    icon="point_of_sale"
    :is-loading="isSubmitting"
    @close="router.push('/sales/dashboard')"
    @refresh="clearTicket"
  >
    <!-- LAYOUT PRINCIPAL DEL TERMINAL (DOS COLUMNAS) -->
    <div class="terminal-grid">
      
      <!-- COLUMNA IZQUIERDA: BÚSQUEDA Y CARRITO -->
      <div class="main-column">
        <!-- ÁREA DE BÚSQUEDA INDUSTRIAL -->
        <section class="terminal-card mb-4 search-section">
          <div class="terminal-search-row">
            <span class="material-symbols-outlined search-icon">barcode_scanner</span>
            <input 
              ref="productSearchInput"
              v-model="productSearch" 
              type="text" 
              class="terminal-input-giant" 
              placeholder="Escanee código o busque producto..." 
              @keyup.enter="handleSearch"
            />
            <button class="btn-terminal-action primary" @click="handleSearch">
              <span class="material-symbols-outlined">add</span>
              <span>Añadir</span>
            </button>
          </div>
        </section>

        <!-- LISTADO DE TICKET (CARRITO) -->
        <section class="terminal-card ticket-list-section">
          <header class="card-header">
            <span class="material-symbols-outlined">shopping_cart</span>
            <h2>Artículos en Ticket ({{ lineItems.length }})</h2>
          </header>
          
          <div class="ticket-items-container">
            <div v-for="item in lineItems" :key="item.productVariantId" class="ticket-item-row">
              <!-- IDENTIDAD (Nombre y SKU) -->
              <div class="item-identity">
                <span class="item-name">{{ item.productName }}</span>
                <div class="item-meta">
                  <code class="item-sku">{{ item.variantSku }}</code>
                  <span v-if="item.optionDescription" class="item-attributes">{{ item.optionDescription }}</span>
                </div>
              </div>
              
              <!-- GRID DE DATOS EN UNA LÍNEA -->
              <div class="item-data-line">
                <!-- CANTIDAD -->
                <div class="data-col qty-col">
                  <div class="qty-stepper-mini">
                    <button @click="updateQtyByItem(item, -1)">-</button>
                    <input v-model.number="item.quantity" type="number" @change="refreshLinePrice(item)" />
                    <button @click="updateQtyByItem(item, 1)">+</button>
                  </div>
                </div>

                <!-- PRECIO VENTA (PVP calculado con márgenes) -->
                <div class="data-col price-col">
                  <span class="value-text">{{ salesApi.formatMoney({ amount: item.unitPrice, currency: 'EUR' }) }}</span>
                </div>

                <!-- DESCUENTO -->
                <div class="data-col discount-col">
                  <div class="discount-input-compact">
                    <input v-model.number="item.discountPercent" type="number" min="0" max="100" />
                    <span class="pct-symbol">%</span>
                  </div>
                </div>

                <!-- TOTAL LÍNEA -->
                <div class="data-col total-col">
                  <span class="value-total">{{ salesApi.formatMoney({ amount: (item.unitPrice * item.quantity) * (1 - item.discountPercent / 100), currency: 'EUR' }) }}</span>
                </div>

                <!-- ACCIONES -->
                <div class="data-col action-col">
                  <button class="btn-remove-item" @click="removeLine(item)">
                    <span class="material-symbols-outlined">close</span>
                  </button>
                </div>
              </div>
            </div>

            <div v-if="lineItems.length === 0" class="terminal-empty-state">
              <span class="material-symbols-outlined">point_of_sale</span>
              <p>Escanee un producto para comenzar la venta</p>
            </div>
          </div>
        </section>
      </div>

      <!-- COLUMNA DERECHA: CLIENTE Y COBRO -->
      <aside class="side-column">
        <!-- SELECCIÃ“N DE CLIENTE -->
        <section class="terminal-card mb-4">
          <header class="card-header">
            <span class="material-symbols-outlined">person</span>
            <h2>Cliente</h2>
          </header>
          <div class="p-4 bg-dark-alt rounded-lg">
            <PartySelector
              v-model="partyId"
              label=""
              placeholder="Buscar cliente..."
              role-filter="CLIENT"
              :dark-mode="true"
              :required="false"
            />
          </div>
        </section>

        <!-- PANEL DE TOTALES Y COBRO -->
        <section class="terminal-card checkout-panel mb-4">
          <div class="totals-area">
            <div class="total-line">
              <label>Subtotal</label>
              <span>{{ salesApi.formatMoney({ amount: subtotal, currency: 'EUR' }) }}</span>
            </div>
            <div class="total-line">
              <label>IVA (21%)</label>
              <span>{{ salesApi.formatMoney({ amount: taxAmount, currency: 'EUR' }) }}</span>
            </div>
            <div class="total-main">
              <label>TOTAL A COBRAR</label>
              <div class="grand-total">{{ salesApi.formatMoney({ amount: total, currency: 'EUR' }) }}</div>
            </div>
          </div>

          <div class="checkout-actions">
            <button 
              class="btn-giant-checkout success" 
              :disabled="lineItems.length === 0 || isSubmitting" 
              @click="processTicket"
            >
              <span class="material-symbols-outlined">print</span>
              <span>COBRAR E IMPRIMIR (F12)</span>
            </button>
            <button class="btn-giant-checkout danger mt-4" @click="clearTicketPrompt">
              <span class="material-symbols-outlined">delete_sweep</span>
              <span>ANULAR TICKET</span>
            </button>
          </div>
        </section>

        <!-- LEYENDA DE ATAJOS -->
        <section class="terminal-card shortcut-legend">
          <header class="card-header">
            <span class="material-symbols-outlined">keyboard</span>
            <h2>Atajos de Teclado</h2>
          </header>
          <div class="shortcuts-grid">
            <div class="shortcut-item"><kbd>F3</kbd> <span>Buscar Producto</span></div>
            <div class="shortcut-item"><kbd>F4</kbd> <span>Seleccionar Cliente</span></div>
            <div class="shortcut-item"><kbd>F12</kbd> <span>Cobrar e Imprimir</span></div>
            <div class="shortcut-item"><kbd>Esc</kbd> <span>Cerrar / Limpiar</span></div>
            <div class="shortcut-item"><kbd>Num +</kbd> <span>Más Cantidad</span></div>
            <div class="shortcut-item"><kbd>Num -</kbd> <span>Menos Cantidad</span></div>
            <div class="shortcut-item"><kbd>Supr</kbd> <span>Eliminar ítem</span></div>
          </div>
        </section>
      </aside>
    </div>

    <!-- MODAL DE PRODUCTOS (ADAPTADO) -->
    <BaseDialog
      :show="showVariantSelector"
      title="Buscador de Productos"
      icon="search"
      size="xl"
      hide-actions
      @close="closeVariantSelector"
    >
      <div class="terminal-dialog-fix">
        <VariantSelector :initial-query="productSearch" @variant-selected="handleVariantSelected" />
      </div>
    </BaseDialog>

    <!-- CAPA DE IMPRESIÃ“N (Teleportada a la raÃ­z para evitar interferencias de la UI) -->
    <Teleport to="body">
      <div id="tpv-print-area" class="print-ticket-container" v-if="lastProcessedTicket">
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
const productSearch = ref('');
const CONSUMIDOR_FINAL_ID = '00000000-0000-0000-0000-000000000001';
const partyId = ref(CONSUMIDOR_FINAL_ID);
const customerDiscount = ref(0);
const lineItems = ref([]);
const isSubmitting = ref(false);
const showVariantSelector = ref(false);
const lastProcessedTicket = ref(null);

const subtotal = computed(() => lineItems.value.reduce((acc, item) => {
  const lineTotal = (item.unitPrice * item.quantity);
  const discount = lineTotal * (item.discountPercent / 100);
  return acc + (lineTotal - discount);
}, 0));
const taxAmount = computed(() => subtotal.value * 0.21);
const total = computed(() => subtotal.value + taxAmount.value);

watch(partyId, async (newId) => {
  if (newId) {
    try {
      const party = await partyApi.getParty(newId);
      customerDiscount.value = party?.default_discount_percentage || 0;
      // Recalcular precios para todas las lÃ­neas con el nuevo cliente
      if (lineItems.value.length > 0) {
        for (const item of lineItems.value) {
          item.discountPercent = customerDiscount.value;
          await refreshLinePrice(item);
        }
      }
    } catch (err) {
      console.error("Error al cargar descuento del cliente:", err);
    }
  } else {
    customerDiscount.value = 0;
  }
  focusSearch();
});

function focusSearch() {
  nextTick(() => {
    productSearchInput.value?.focus();
  });
}

function handleGlobalKeydown(e) {
  // Solo procesar si no estamos en un input que no sea el de búsqueda principal
  // a menos que sean teclas de función específicas.
  const isInput = ['INPUT', 'TEXTAREA', 'SELECT'].includes(document.activeElement.tagName);
  const isSearchInput = document.activeElement === productSearchInput.value;

  // F3: Buscar Producto
  if (e.key === 'F3') {
    e.preventDefault();
    focusSearch();
  }
  
  // F4: Seleccionar Cliente
  if (e.key === 'F4') {
    e.preventDefault();
    const partyInput = document.querySelector('.party-selector input');
    partyInput?.focus();
  }

  // F12: Cobrar
  if (e.key === 'F12') {
    e.preventDefault();
    if (lineItems.value.length > 0 && !isSubmitting.value) {
      processTicket();
    }
  }

  // Esc: Cerrar diálogos o limpiar búsqueda
  if (e.key === 'Escape') {
    if (showVariantSelector.value) {
      closeVariantSelector();
    } else {
      productSearch.value = '';
      focusSearch();
    }
  }

  // Atajos para el carrito (solo si no estamos escribiendo en un input de datos)
  if (!isInput || isSearchInput) {
    // Num + : Incrementar última línea
    if (e.key === '+' || e.key === 'Add') {
      if (lineItems.value.length > 0) {
        e.preventDefault();
        updateQtyByItem(lineItems.value[lineItems.value.length - 1], 1);
      }
    }

    // Num - : Decrementar última línea
    if (e.key === '-' || e.key === 'Subtract') {
      if (lineItems.value.length > 0) {
        e.preventDefault();
        updateQtyByItem(lineItems.value[lineItems.value.length - 1], -1);
      }
    }

    // Supr : Eliminar última línea
    if (e.key === 'Delete') {
      if (lineItems.value.length > 0) {
        e.preventDefault();
        removeLine(lineItems.value[lineItems.value.length - 1]);
      }
    }
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleGlobalKeydown);
  // Aseguramos que cargue consumidor final al inicio
  if (!partyId.value) partyId.value = CONSUMIDOR_FINAL_ID;
  await loadDefaultCustomer();
  focusSearch();
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleGlobalKeydown);
});

async function loadDefaultCustomer() {
  // Ya inicializamos con el ID constante de CONSUMIDOR FINAL
  try {
    const party = await partyApi.getParty(CONSUMIDOR_FINAL_ID);
    if (party) {
      partyId.value = party.id;
      customerDiscount.value = party.default_discount_percentage || 0;
    }
  } catch (err) {
    console.error("Error cargando consumidor final:", err);
  }
}

async function handleSearch() {
  const query = productSearch.value.trim();
  if (!query) { 
    showVariantSelector.value = true; 
    return; 
  }
  try {
    const result = await productApi.smartSearch(query);
    if (result.type === 'exact_variant') {
      handleVariantSelected(result.variant);
    } else {
      showVariantSelector.value = true;
    }
  } catch (err) { 
    console.error("Error en búsqueda rápida:", err);
    showVariantSelector.value = true; 
  }
}

async function handleVariantSelected(v) {
  const variant = v.variant || v;
  const existing = lineItems.value.find(i => i.productVariantId === variant.id);
  
  // Formateamos los atributos (Color, Talla, etc.)
  let optionDescription = '';
  const options = variant.option_configuration || variant.optionConfiguration;
  if (options) {
    optionDescription = Object.entries(options)
      .map(([key, val]) => `${key}: ${val}`)
      .join(' | ');
  }

  if (existing) {
    existing.quantity++;
    await refreshLinePrice(existing);
  } else {
    // Creamos el item pero CON PRECIO CERO inicialmente para forzar la carga
    const newItem = {
      productVariantId: variant.id, 
      variantSku: variant.sku,
      productName: variant.product_name || variant.productName || 'Producto', 
      optionDescription: optionDescription,
      quantity: 1,
      unitPrice: 0, 
      discountPercent: customerDiscount.value 
    };
    
    // CARGA CRÃTICA: Esperamos a que el motor de pricing nos dÃ© el "Precio Venta" (PVP con margen)
    await refreshLinePrice(newItem);
    
    // Solo lo aÃ±adimos cuando tenemos el precio real del motor
    lineItems.value.push(newItem);
  }
  productSearch.value = '';
  showVariantSelector.value = false;
  focusSearch();
}

function closeVariantSelector() {
  showVariantSelector.value = false;
  focusSearch();
}

async function refreshLinePrice(item) {
  if (!partyId.value) return;
  try {
    const result = await pricingApi.calculateFinalSalePrice(
      [{ productVariantId: item.productVariantId, quantity: item.quantity }],
      partyId.value,
      new Date()
    );
    
    if (result.calculatedItems && result.calculatedItems.length > 0) {
      const calc = result.calculatedItems[0];
      
      // PRECIO VENTA = BaseSalesPrice (Coste + Suplementos + Margen)
      // Usamos parentheses para evitar el error de mezcla de operadores ?? y ||
      item.unitPrice = (calc.baseSalesPrice?.amount ?? calc.baseSalesPrice) || 0;
      
      // DESCUENTO: Prioridad al del cliente, salvo que el motor traiga una regla mÃ¡s fuerte
      const engineDiscount = calc.discountPercent || 0;
      if (Math.abs(engineDiscount - customerDiscount.value) < 0.5) {
        item.discountPercent = customerDiscount.value;
      } else {
        item.discountPercent = Math.round(engineDiscount * 100) / 100;
      }
    }
  } catch (err) {
    console.error("Error en motor de pricing:", err);
    item.discountPercent = customerDiscount.value;
  }
}

async function updateQtyByItem(item, delta) {
  item.quantity = Math.max(1, item.quantity + delta);
  await refreshLinePrice(item);
}

function removeLine(item) { 
  lineItems.value = lineItems.value.filter(i => i.productVariantId !== item.productVariantId); 
}

function clearTicket() { 
  lineItems.value = []; 
  productSearch.value = ''; 
  partyId.value = CONSUMIDOR_FINAL_ID;
  focusSearch();
}
function clearTicketPrompt() { if (confirm('¿Vaciar el ticket actual?')) clearTicket(); }

async function processTicket() {
  if (lineItems.value.length === 0) return;
  if (!partyId.value) {
    alert('Debes seleccionar un cliente (o cargar el por defecto).');
    return;
  }

  isSubmitting.value = true;
  try {
    console.log('Registrando venta en el backend...');
    
    const request = {
      partyId: partyId.value,
      invoiceDate: new Date().toISOString(), // ISO String para el backend en Go
      items: lineItems.value.map(item => ({
        productVariantId: item.productVariantId,
        quantity: item.quantity,
        discountPercent: Number(item.discountPercent || 0)
      }))
    };

    const invoice = await salesApi.createSimplifiedInvoice(request);
    console.log('Venta registrada con éxito:', invoice.invoice_number);

    // Mapeamos la respuesta manejando camelCase (API) vs snake_case (Frontend)
    const apiLineItems = invoice.lineItems || invoice.line_items || [];
    
    lastProcessedTicket.value = {
      number: invoice.invoice_number || invoice.invoiceNumber,
      date: invoice.issue_date || invoice.invoiceDate || new Date().toISOString(),
      items: apiLineItems.map(li => ({
        productName: li.productName || li.product_name || 'Producto',
        variantSku: li.variantSku || li.variant_sku || '', 
        quantity: li.quantity,
        unitPrice: (li.unitPrice?.amount !== undefined ? li.unitPrice.amount : li.unit_price) || 0,
        discountPercent: li.discountPercent || li.discount_percentage || 0,
        subtotal: (li.subtotal?.amount !== undefined ? li.subtotal.amount : li.subtotal) || 0
      })),
      totals: { 
        subtotal: (invoice.subtotal?.amount !== undefined ? invoice.subtotal.amount : invoice.subtotal) || 0, 
        taxAmount: (invoice.taxAmount?.amount !== undefined ? invoice.taxAmount.amount : invoice.tax_total) || 0, 
        total: (invoice.total?.amount !== undefined ? invoice.total.amount : invoice.total) || 0 
      },
      customerName: invoice.partyName || invoice.party_name || 'CONSUMIDOR FINAL'
    };

    // Pequeño retardo para asegurar que el DOM de PrintTicket se actualiza
    await new Promise(resolve => setTimeout(resolve, 500));
    window.print();

    // Limpieza tras imprimir
    lineItems.value = [];
    productSearch.value = '';
    lastProcessedTicket.value = null;
    focusSearch();
    
  } catch (err) {
    console.error('Error al procesar la venta:', err);
    alert('Error al registrar la venta: ' + (err.message || 'Error desconocido'));
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<style>
/* REGLAS DE IMPRESIÃ“N GLOBALES (Fuera de scoped para poder ocultar #app) */
@media print {
  /* Ocultamos absolutamente todos los hijos directos del body */
  body > * {
    display: none !important;
  }

  /* Centramos el contenido del body para que el ticket salga en medio */
  body {
    display: flex !important;
    justify-content: center !important;
    align-items: flex-start !important;
    background: white !important;
    width: 100% !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  /* Exceptuamos el Ã¡rea de impresiÃ³n del TPV que estÃ¡ teleportada al body */
  body > #tpv-print-area {
    display: block !important;
    width: 80mm !important;
    margin: 0 auto !important;
    padding: 0 !important;
    background: white !important;
    box-shadow: none !important;
  }

  /* Reset de mÃ¡rgenes de pÃ¡gina para impresoras tÃ©rmicas */
  @page {
    margin: 0;
    size: 80mm auto;
  }
}
</style>

<style scoped>
.terminal-grid {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 1.5rem;
  height: 100%;
  overflow: hidden;
}

.main-column, .side-column {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.side-column {
  overflow-y: auto;
}

/* Industrial Card */
.terminal-card {
  background: #1e293b;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  border: 1px solid #334155;
  /* Eliminamos overflow hidden para permitir que los dropdowns salgan */
}

.card-header {
  padding: 0.75rem 1.25rem; /* Reducido */
  background: #0f172a;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  border-bottom: 1px solid #334155;
  border-radius: 16px 16px 0 0;
}

.card-header h2 { font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: #94a3b8; margin: 0; letter-spacing: 0.05em; }
.card-header .material-symbols-outlined { color: var(--color-primary); font-size: 1.1rem; }

/* Giant Search Row */
.terminal-search-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem; /* Reducido */
}

.search-icon { font-size: 2rem; color: #475569; }
.terminal-input-giant {
  flex: 1;
  background: #0f172a;
  border: 2px solid #334155;
  border-radius: 12px;
  padding: 1rem 1.25rem; /* Reducido */
  font-size: 1.25rem; /* Reducido */
  color: white;
  font-weight: 700;
}
.terminal-input-giant:focus { border-color: var(--color-primary); outline: none; }

.btn-terminal-action {
  padding: 1rem 1.5rem; /* Reducido */
  border-radius: 12px;
  border: none;
  font-weight: 900;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.btn-terminal-action.primary { background: var(--color-primary); color: black; }

/* Items List */
.ticket-list-section { flex: 1; overflow: hidden; }
.ticket-items-container { flex: 1; overflow-y: auto; padding: 0.75rem; }

.ticket-item-row {
  background: #0f172a;
  margin-bottom: 0.5rem;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  border-left: 4px solid var(--color-primary);
  border: 1px solid #334155;
}

.item-identity { flex: 1; min-width: 200px; display: flex; flex-direction: column; gap: 0.15rem; }
.item-name { font-size: 0.95rem; font-weight: 700; color: white; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.item-meta { display: flex; align-items: center; gap: 0.5rem; }
.item-attributes { font-size: 0.7rem; color: var(--color-primary); font-weight: 600; opacity: 0.8; }
.item-sku { 
  font-size: 0.65rem; 
  color: #1e293b; 
  background: #94a3b8; 
  padding: 0.05rem 0.3rem; 
  border-radius: 4px; 
  width: fit-content;
  font-weight: 700;
  font-family: monospace; 
}

.item-data-line {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  flex-shrink: 0;
}

.data-col { display: flex; align-items: center; }

/* Global Hide Spinners */
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
input[type=number] {
  -moz-appearance: textfield;
}

/* Stepper Mini */
.qty-stepper-mini { display: flex; align-items: center; background: #1e293b; border-radius: 6px; overflow: hidden; border: 1px solid #334155; height: 36px; }
.qty-stepper-mini button { width: 32px; height: 100%; border: none; background: transparent; color: white; font-size: 1.2rem; cursor: pointer; }
.qty-stepper-mini button:active { background: #334155; }
.qty-stepper-mini input { width: 40px; text-align: center; background: transparent; border: none; color: white; font-size: 1rem; font-weight: 800; }

/* Precio y Total */
.value-text { font-size: 1rem; font-weight: 600; color: #94a3b8; width: 80px; text-align: right; }
.value-total { font-size: 1.1rem; font-weight: 900; color: #16a34a; width: 100px; text-align: right; }

/* Descuento Compacto */
.discount-input-compact {
  position: relative;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  height: 36px;
  display: flex;
  align-items: center;
  padding: 0 0.5rem;
  width: 70px;
}
.discount-input-compact input {
  width: 100%;
  background: transparent;
  border: none;
  color: var(--color-primary);
  font-size: 1rem;
  font-weight: 800;
  text-align: right;
  padding-right: 0.75rem;
}
.discount-input-compact input:focus { outline: none; }
.pct-symbol { position: absolute; right: 0.4rem; font-size: 0.8rem; font-weight: 800; color: var(--color-primary); }

.btn-remove-item { background: transparent; border: none; color: #64748b; cursor: pointer; display: flex; align-items: center; }
.btn-remove-item:hover { color: #ef4444; }
.btn-remove-item .material-symbols-outlined { font-size: 1.25rem; }

/* Totals and Checkout */
.checkout-panel { flex: 1; min-height: 0; display: flex; flex-direction: column; }
.totals-area { padding: 1.25rem; display: flex; flex-direction: column; gap: 0.5rem; }
.total-line { display: flex; justify-content: space-between; font-size: 0.9rem; color: #94a3b8; }
.total-main { margin-top: 0.75rem; padding-top: 0.75rem; border-top: 2px solid #334155; display: flex; flex-direction: column; align-items: center; gap: 0.5rem; }
.total-main label { font-size: 0.7rem; font-weight: 800; color: #64748b; letter-spacing: 0.1em; }
.grand-total { font-size: 2.25rem; font-weight: 900; color: var(--color-primary); line-height: 1; }

.checkout-actions { padding: 1rem; margin-top: auto; }
.btn-giant-checkout {
  width: 100%;
  height: 64px; /* Reducido de 80px */
  border-radius: 16px;
  border: none;
  font-size: 1.1rem;
  font-weight: 900;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  cursor: pointer;
}
.btn-giant-checkout.success { background: #16a34a; color: white; }
.btn-giant-checkout.danger { background: transparent; border: 2px solid #334155; color: #ef4444; height: 60px; font-size: 1rem; }

.terminal-empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; opacity: 0.3; }
.terminal-empty-state .material-symbols-outlined { font-size: 5rem; margin-bottom: 1rem; }

/* Shortcut Legend */
.shortcut-legend {
  margin-top: auto;
}

.shortcuts-grid {
  padding: 1rem;
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.5rem;
}

.shortcut-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.8rem;
  color: #94a3b8;
}

.shortcut-item kbd {
  background: #334155;
  color: var(--color-primary);
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-family: monospace;
  font-weight: 800;
  min-width: 40px;
  text-align: center;
  border-bottom: 2px solid #0f172a;
}

/* Utils */
.bg-dark-alt { background: #0f172a; }
.p-4 { padding: 1rem; }
.rounded-lg { border-radius: 8px; }
.mt-4 { margin-top: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mb-6 { margin-bottom: 1.5rem; }

.print-ticket-container { display: none; }

@media print {
  /* Ocultamos absolutamente todo desde la raÃ­z */
  #app, 
  .app-shell, 
  .app-layout, 
  .app-main, 
  .base-terminal-overlay,
  header, 
  nav, 
  aside, 
  main, 
  footer,
  .no-print { 
    display: none !important; 
    height: 0 !important;
    overflow: hidden !important;
  }

  /* El body debe estar limpio para el ticket */
  body {
    background: white !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  /* Mostramos solo el contenedor del ticket y forzamos su visibilidad */
  .print-ticket-container { 
    display: block !important; 
    position: fixed !important; 
    left: 0 !important; 
    top: 0 !important; 
    width: 80mm !important; 
    z-index: 9999999 !important;
    background: white !important;
    color: black !important;
  }
}

/* Fix para diálogos heredados */
.terminal-dialog-fix { color: #1e293b; }
:deep(.form-label) { color: #64748b; font-weight: 700; font-size: 0.75rem; text-transform: uppercase; margin-bottom: 0.5rem; display: block; }
</style>
