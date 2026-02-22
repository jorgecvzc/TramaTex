<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Operaciones / Productos</p>
          <h1>Detalle de producto</h1>
          <p class="subtitle">Información completa, variantes, atributos y precios.</p>
        </div>
        <RouterLink to="/products" class="btn btn-secondary">
          Volver al catálogo
        </RouterLink>
      </header>

      <!-- Loading State -->
      <div v-if="isLoading" class="loading">
        <div class="spinner"></div>
        <p>Cargando producto...</p>
      </div>

      <!-- Error State -->
      <div v-if="error" class="alert-error">
        {{ error }}
        <RouterLink to="/products" class="btn btn-outline">
          ← Volver al catálogo
        </RouterLink>
      </div>

      <!-- Product Detail -->
      <div v-if="!isLoading && product" class="detail-container">
        <!-- Header Card -->
        <div class="detail-header card">
          <div class="header-content">
            <div>
              <div class="header-title">
                <h2>{{ product.name }}</h2>
                <code v-if="product.sku" class="sku-badge">{{ product.sku }}</code>
              </div>
              <p v-if="product.long_name" class="long-name">{{ product.long_name }}</p>
            </div>
            <div class="header-badges">
              <span class="pill" :class="`type-${product.product_type.toLowerCase()}`">
                {{ product.product_type === 'TANGIBLE' ? 'Tangible' : 'Servicio' }}
              </span>
              <span class="pill" :class="product.is_active ? 'active' : 'inactive'">
                {{ product.is_active ? 'Activo' : 'Inactivo' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Tabs Navigation -->
        <div class="tabs-container card">
          <div class="tabs-nav">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              :class="['tab-button', { active: activeTab === tab.id }]"
            >
              <component :is="tab.icon" :size="20" class="tab-icon" />
              <span class="tab-label">{{ tab.label }}</span>
              <span v-if="tab.count !== undefined" class="tab-count">{{ tab.count }}</span>
            </button>
          </div>

          <!-- Tab Content -->
          <div class="tab-content">
            <!-- Tab 1: General Information -->
            <ProductDetailInfo
              v-if="activeTab === 'info'"
              :product="product"
              :brand="brand"
              :groups="groups"
              @update="handleProductUpdate"
              @toggle-status="handleToggleStatus"
            />

            <!-- Tab 2: Variants -->
            <VariantTable
              v-if="activeTab === 'variants'"
              :product-id="productId"
              :product-sku="product?.sku || ''"
              :variants="variants"
              :is-loading="isLoadingVariants"
              @refresh="fetchVariants"
            />

            <!-- Tab 3: Attributes -->
            <AttributesPanel
              v-if="activeTab === 'attributes'"
              :product="product"
              :calculated-attributes="calculatedAttributes"
              :is-loading="isLoadingAttributes"
              @refresh="fetchCalculatedAttributes"
            />

            <!-- Tab 4: Pricing -->
            <PricingPanel
              v-if="activeTab === 'pricing'"
              :product-id="productId"
              :product-name="product?.name || ''"
              :variants="variants"
              :is-loading-variants="isLoadingVariants"
            />

            <!-- Tab 5: History -->
            <div v-if="activeTab === 'history'" class="history-tab">
              <h3>Historial de Cambios</h3>
              <div class="empty-state">
                <ClipboardList :size="48" class="empty-icon" />
                <p>El historial de auditoría estará disponible próximamente.</p>
                <p class="empty-hint">
                  Aquí podrás ver todos los cambios realizados en este producto,
                  incluyendo quién y cuándo se modificó.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import ProductDetailInfo from '@/components/product/ProductDetailInfo.vue'
import VariantTable from '@/components/product/VariantTable.vue'
import AttributesPanel from '@/components/product/AttributesPanel.vue'
import PricingPanel from '@/components/product/PricingPanel.vue'
import { productApi } from '@/services/productApi'
import { FileText, Hash, Tag, DollarSign, ClipboardList } from 'lucide-vue-next'

const route = useRoute()
const productId = route.params.id

// State
const product = ref(null)
const brand = ref(null)
const groups = ref([])
const variants = ref([])
const calculatedAttributes = ref([])
const isLoading = ref(false)
const isLoadingVariants = ref(false)
const isLoadingAttributes = ref(false)
const error = ref('')
const activeTab = ref('info')

// Tabs configuration
const tabs = computed(() => [
  {
    id: 'info',
    label: 'Información',
    icon: FileText,
  },
  {
    id: 'variants',
    label: 'Variantes',
    icon: Hash,
    count: variants.value.length,
  },
  {
    id: 'attributes',
    label: 'Atributos',
    icon: Tag,
    count: calculatedAttributes.value.length,
  },
  {
    id: 'pricing',
    label: 'Precios',
    icon: DollarSign,
  },
  {
    id: 'history',
    label: 'Historial',
    icon: ClipboardList,
  },
])

// Lifecycle
onMounted(async () => {
  await fetchProduct()
  await Promise.all([
    fetchVariants(),
    fetchCalculatedAttributes(),
  ])
})

// Methods
async function fetchProduct() {
  isLoading.value = true
  error.value = ''

  try {
    const data = await productApi.getProduct(productId)
    product.value = data

    // Fetch related data
    if (data.brand_id) {
      brand.value = await productApi.getBrand(data.brand_id)
    }

    if (data.group_ids && data.group_ids.length > 0) {
      const groupPromises = data.group_ids.map(id =>
        productApi.getProductGroup(id).catch(() => null)
      )
      const fetchedGroups = await Promise.all(groupPromises)
      groups.value = fetchedGroups.filter(g => g !== null)
    }
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar el producto'
  } finally {
    isLoading.value = false
  }
}

async function fetchVariants() {
  isLoadingVariants.value = true

  try {
    const data = await productApi.listProductVariants(productId)
    variants.value = data.variants || data || []
  } catch (err) {
    console.error('Error fetching variants:', err)
    variants.value = []
  } finally {
    isLoadingVariants.value = false
  }
}

async function fetchCalculatedAttributes() {
  isLoadingAttributes.value = true

  try {
    const data = await productApi.getCalculatedAttributes(productId)
    calculatedAttributes.value = data.attributes || data || []
  } catch (err) {
    console.error('Error fetching attributes:', err)
    calculatedAttributes.value = []
  } finally {
    isLoadingAttributes.value = false
  }
}

async function handleProductUpdate(updatedData) {
  try {
    const updated = await productApi.updateProduct(productId, updatedData)
    product.value = { ...product.value, ...updated }
  } catch (err) {
    alert(err?.message || 'No se pudo actualizar el producto')
  }
}

async function handleToggleStatus() {
  if (!product.value) return

  try {
    const newStatus = !product.value.is_active
    await productApi.changeProductStatus(productId, newStatus)
    product.value.is_active = newStatus
  } catch (err) {
    alert(err?.message || 'No se pudo cambiar el estado del producto')
  }
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.page-header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.page-header h1 {
  color: #1b3a6b;
  margin: 0.25rem 0 0;
}

.breadcrumb {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin: 0;
}

.subtitle {
  color: #64748b;
  margin: 0.5rem 0 0;
  font-size: 0.95rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover {
  background: #f8fafc;
}

.btn-outline {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  gap: 1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(27, 58, 107, 0.12);
  border-top-color: #1b3a6b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.alert-error {
  background: #fee2e2;
  border: 1px solid #ef4444;
  border-radius: 8px;
  padding: 1rem;
  text-align: center;
  color: #991b1b;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  align-items: center;
}

.detail-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  width: 100%;
  gap: 2rem;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.header-title h2 {
  color: #1b3a6b;
  margin: 0;
  font-size: 1.6rem;
}

.sku-badge {
  background: #f1f5f9;
  color: #475569;
  padding: 0.25rem 0.75rem;
  border-radius: 6px;
  font-size: 0.85rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 600;
}

.long-name {
  color: #64748b;
  margin: 0.5rem 0 0;
  font-size: 0.95rem;
}

.header-badges {
  display: flex;
  gap: 0.75rem;
  flex-shrink: 0;
}

.pill {
  display: inline-block;
  padding: 0.35rem 0.85rem;
  border-radius: 999px;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.pill.type-tangible {
  background-color: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.pill.type-service {
  background-color: rgba(139, 92, 246, 0.1);
  color: #8b5cf6;
}

.pill.active {
  background-color: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.pill.inactive {
  background-color: rgba(148, 163, 184, 0.1);
  color: #94a3b8;
}

.tabs-container {
  padding: 0;
  overflow: hidden;
}

.tabs-nav {
  display: flex;
  gap: 0;
  border-bottom: 2px solid #e2e8f0;
  padding: 0 1.5rem;
  background: #f8fafc;
}

.tab-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  cursor: pointer;
  color: #64748b;
  font-size: 0.9rem;
  font-weight: 600;
  transition: all 0.2s ease;
  margin-bottom: -2px;
}

.tab-button:hover {
  color: #1b3a6b;
  background: rgba(27, 58, 107, 0.05);
}

.tab-button.active {
  color: #1b3a6b;
  border-bottom-color: #f4d03f;
  background: #ffffff;
}

.tab-icon {
  font-size: 1.1rem;
}

.tab-label {
  white-space: nowrap;
}

.tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.5rem;
  height: 1.5rem;
  padding: 0 0.4rem;
  background: #e2e8f0;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 700;
  color: #475569;
}

.tab-button.active .tab-count {
  background: #f4d03f;
  color: #1e293b;
}

.tab-content {
  padding: 1.5rem;
  min-height: 400px;
}

.history-tab {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.history-tab h3 {
  color: #1b3a6b;
  margin: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
  color: #64748b;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state p {
  margin: 0.5rem 0;
  font-size: 0.95rem;
}

.empty-hint {
  font-size: 0.85rem;
  color: #94a3b8;
  max-width: 500px;
}

@media (max-width: 768px) {
  .dashboard-content {
    padding: 1.5rem;
  }

  .header-content {
    flex-direction: column;
    gap: 1rem;
  }

  .header-badges {
    width: 100%;
  }

  .tabs-nav {
    overflow-x: auto;
    padding: 0 1rem;
  }

  .tab-button {
    padding: 1rem;
    font-size: 0.85rem;
  }

  .tab-label {
    display: none;
  }

  .tab-icon {
    font-size: 1.3rem;
  }
}
</style>
