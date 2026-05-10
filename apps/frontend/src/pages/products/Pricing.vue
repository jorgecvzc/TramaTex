<template>
  <BaseDashboardPage :is-loading="isLoadingProducts" class="pricing-dashboard">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <BasePageHeader 
        title="Consulta de Precios" 
        :breadcrumbs="[{ label: 'Catálogo', to: '/products' }, { label: 'Precios' }]"
      >
        <template #icon><Banknote :size="24" /></template>
        <template #actions>
          <button 
            v-if="selectedProductId && variants.length > 0"
            class="btn btn-primary btn-sm" 
            :disabled="isCalculatingBase || isCalculatingFinal"
            @click="calculateAllPrices"
          >
            <Calculator :size="18" />
            Calcular Todo
          </button>
        </template>
      </BasePageHeader>
    </template>

    <!-- CAPA 3: TRABAJO (Resultados) -->
    <div class="pricing-results-area">
      <div v-if="calculationError" class="alert-error mb-4">
        <AlertCircle :size="18" />
        {{ calculationError }}
      </div>

      <div class="dashboard-section">
        <div class="section-header">
          <BarChart4 :size="18" />
          <h2>Precios por Variante</h2>
          <span v-if="variants.length" class="header-tag">{{ variants.length }} variantes</span>
        </div>

        <!-- Estados vacíos y carga -->
        <div v-if="!selectedProductId" class="empty-state">
          <Package :size="48" />
          <p>Selecciona un producto en el panel lateral para iniciar el cálculo.</p>
        </div>

        <div v-else-if="isLoadingVariants" class="loading-state">
          <div class="spinner"></div>
          <p>Cargando variantes del producto...</p>
        </div>

        <div v-else-if="variants.length === 0" class="empty-state">
          <SearchX :size="48" />
          <p>Este producto no tiene variantes configuradas.</p>
        </div>

        <!-- Tabla de Resultados -->
        <div v-else class="table-wrapper">
          <table class="data-table compact">
            <thead>
              <tr>
                <th>SKU / Atributos</th>
                <th class="align-right">Precio Ud.</th>
                <th class="align-right">Cant.</th>
                <th v-if="clientDiscount > 0" class="align-right">Dto.</th>
                <th class="align-right">Base Imp.</th>
                <th class="align-right">Total c/IVA</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="variant in variants" :key="variant.id" class="row-hover">
                <td>
                  <div class="variant-cell">
                    <code class="code-badge">{{ variant.sku }}</code>
                    <span class="attr-text">{{ getAttributeSummary(variant) }}</span>
                  </div>
                </td>
                <td class="align-right">
                  <span v-if="variantPrices[variant.id]?.baseSalesPrice" class="price-val">
                    {{ formatMoney(variantPrices[variant.id].baseSalesPrice) }}
                  </span>
                  <span v-else class="text-muted">—</span>
                </td>
                <td class="align-right text-muted">{{ quantity }}</td>
                <td v-if="clientDiscount > 0" class="align-right text-danger">
                  -{{ calcDiscountAmount(variant.id) }}
                </td>
                <td class="align-right"><strong>{{ calcTaxBase(variant.id) }}</strong></td>
                <td class="align-right"><span class="total-val">{{ calcLineTotal(variant.id) }}</span></td>
                <td>
                  <span :class="['pill', variant.is_active ? 'active' : 'inactive']">
                    {{ variant.is_active ? 'Activo' : 'No' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- CAPA 2: CONTEXTO (Filtros en Sidebar) -->
    <template #sidebar>
      <div class="sidebar-filters">
        <section class="sidebar-section">
          <div class="section-header">
            <Filter :size="18" />
            <h2>Condiciones</h2>
          </div>
          
          <div class="form-group">
            <label>Producto Base</label>
            <select v-model="selectedProductId" @change="onProductChange">
              <option value="">Seleccionar...</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <PartySelector
              v-model="selectedClientId"
              label="Cliente"
              placeholder="Buscar cliente..."
              role-filter="CLIENT"
              :required="false"
            />
            <div v-if="selectedClientId && clientDiscount > 0" class="client-info-mini">
              <Sparkles :size="16" />
              Bonificación activa: <strong>{{ clientDiscount }}%</strong>
            </div>
          </div>

          <div class="form-grid">
            <div class="form-group">
              <label>Cantidad</label>
              <input v-model.number="quantity" type="number" min="1" />
            </div>
            <div class="form-group">
              <label>Fecha Venta</label>
              <input v-model="saleDate" type="date" />
            </div>
          </div>
        </section>

        <section class="help-notice">
          <div class="notice-header">
            <div class="flex items-center gap-2 mb-2">
              <Info :size="16" class="text-secondary" />
              <h3>Cálculo Dinámico</h3>
            </div>
          </div>
          <p class="help-text">
            Los precios se calculan en base a las reglas de margen, tarifas por cliente y volumen de compra definidos en el sistema.
          </p>
        </section>
      </div>
    </template>
  </BaseDashboardPage>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { 
  Banknote, 
  Calculator, 
  AlertCircle, 
  BarChart4, 
  Package, 
  SearchX, 
  Filter, 
  Sparkles, 
  Info 
} from 'lucide-vue-next'
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import PartySelector from '@/components/party/PartySelector.vue'
import { productApi } from '@/services/productApi'
import { pricingApi } from '@/services/pricingApi'
import { partyApi } from '@/services/partyApi'
import salesApi from '@/services/salesApi'

// State
const products = ref([])
const isLoadingProducts = ref(false)
const selectedProductId = ref('')
const variants = ref([])
const isLoadingVariants = ref(false)
const variantPrices = ref({})
const isCalculatingBase = ref(false)
const isCalculatingFinal = ref(false)
const calculationError = ref('')
const clientDiscount = ref(0)
const quantity = ref(100)
const saleDate = ref(new Date().toISOString().split('T')[0])
const selectedClientId = ref('')

// Watchers
watch(selectedClientId, async (newId) => {
  if (newId) {
    try {
      const client = await partyApi.getParty(newId)
      clientDiscount.value = client?.default_discount_percentage ?? 0
      if (variants.value.length > 0) calculateFinalPrices()
    } catch { clientDiscount.value = 0 }
  } else { clientDiscount.value = 0 }
})

// Lifecycle
onMounted(async () => {
  isLoadingProducts.value = true
  try {
    const result = await productApi.listProducts()
    products.value = result.data || []
  } finally { isLoadingProducts.value = false }
})

// Methods
async function onProductChange() {
  if (!selectedProductId.value) { variants.value = []; return }
  isLoadingVariants.value = true
  try {
    const data = await productApi.listProductVariants(selectedProductId.value)
    variants.value = data.variants || data || []
    if (variants.value.length > 0) {
      await calculateBasePrices()
      if (selectedClientId.value) await calculateFinalPrices()
    }
  } finally { isLoadingVariants.value = false }
}

async function calculateAllPrices() {
  await calculateBasePrices()
  if (selectedClientId.value) await calculateFinalPrices()
}

async function calculateBasePrices() {
  isCalculatingBase.value = true
  for (const variant of variants.value) {
    try {
      const res = await pricingApi.calculateBaseSalesPrice(selectedProductId.value, variant.id)
      variantPrices.value[variant.id] = { ...variantPrices.value[variant.id], baseSalesPrice: res?.baseSalesPrice }
    } catch (err) { variantPrices.value[variant.id] = { baseError: 'Error' } }
  }
  isCalculatingBase.value = false
}

async function calculateFinalPrices() {
  isCalculatingFinal.value = true
  try {
    const saleItems = variants.value.map(v => ({ productVariantId: v.id, quantity: quantity.value || 1 }))
    const res = await pricingApi.calculateFinalSalePrice(saleItems, selectedClientId.value, new Date(saleDate.value))
    if (res?.calculatedItems) {
      res.calculatedItems.forEach(item => {
        variantPrices.value[item.productVariantId] = { ...variantPrices.value[item.productVariantId], ...item }
      })
    }
  } finally { isCalculatingFinal.value = false }
}

const getAttributeSummary = (v) => {
  const cfg = v.option_configuration || v.attribute_values || {}
  return Object.entries(cfg).map(([k, val]) => `${val}`).join(', ') || 'Sin atributos'
}

const formatMoney = (m) => salesApi.formatMoney(m)
const calcSubtotal = (id) => formatMoney({ amount: (parseFloat(variantPrices.value[id]?.baseSalesPrice?.amount || 0) * (quantity.value || 1)).toString() })
const calcTaxBase = (id) => formatMoney(variantPrices.value[id]?.finalPrice || { amount: (parseFloat(variantPrices.value[id]?.baseSalesPrice?.amount || 0) * (quantity.value || 1)).toString() })
const calcLineTotal = (id) => formatMoney(variantPrices.value[id]?.finalPriceWithTax || { amount: (parseFloat(variantPrices.value[id]?.finalPrice?.amount || 0) * 1.21).toString() })
const calcDiscountAmount = (id) => {
  const base = parseFloat(variantPrices.value[id]?.baseSalesPrice?.amount || 0) * (quantity.value || 1)
  const final = parseFloat(variantPrices.value[id]?.finalPrice?.amount || 0) * (quantity.value || 1)
  return formatMoney({ amount: (base - final).toString() })
}
</script>

<style scoped>
.pricing-results-area { display: flex; flex-direction: column; gap: 1rem; }
.dashboard-section { background: white; padding: 1rem; border-radius: 12px; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.section-header { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 1rem; padding-bottom: 0.6rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.9rem; font-weight: 700; text-transform: uppercase; margin: 0; flex: 1; }
.header-tag { font-size: 0.65rem; font-weight: 800; padding: 0.15rem 0.6rem; background: var(--color-background); color: var(--color-secondary); border-radius: 20px; }

.variant-cell { display: flex; flex-direction: column; gap: 0.15rem; }
.attr-text { font-size: 0.75rem; color: var(--color-text-secondary); font-weight: 500; }
.price-val { color: var(--color-secondary); font-weight: 600; }
.total-val { color: var(--color-success); font-weight: 800; font-size: 1rem; }

.sidebar-filters { display: flex; flex-direction: column; gap: 1.5rem; }
.client-info-mini { margin-top: 0.5rem; font-size: 0.75rem; color: #166534; background: #f0fdf4; padding: 0.4rem 0.6rem; border-radius: 6px; display: flex; align-items: center; gap: 0.3rem; }
.client-info-mini svg { color: var(--color-success); }

.alert-error { background: #fef2f2; border: 1px solid #fecaca; color: #991b1b; padding: 0.75rem 1rem; border-radius: 8px; display: flex; align-items: center; gap: 0.5rem; font-size: 0.9rem; }

.empty-state, .loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 4rem 2rem; text-align: center; color: var(--color-text-secondary); }
.empty-state svg { opacity: 0.3; margin-bottom: 1rem; }

.pill { padding: 0.2rem 0.5rem; border-radius: 4px; font-size: 0.65rem; font-weight: 800; text-transform: uppercase; }
.pill.active { background: #dcfce7; color: #166534; }
.pill.inactive { background: #fee2e2; color: #991b1b; }

.align-right { text-align: right; }
.help-notice { padding: 1rem; background: var(--color-background); border-radius: 10px; border: 1px dashed var(--color-border); }
.help-text { font-size: 0.75rem; color: var(--color-text-secondary); line-height: 1.4; margin-top: 0.5rem; }
</style>
