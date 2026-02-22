<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Operaciones / Productos</p>
          <h1>Crear producto</h1>
          <p class="subtitle">Define un nuevo producto en el catálogo con sus variantes.</p>
        </div>
        <RouterLink to="/products" class="btn btn-secondary">
          Volver al catálogo
        </RouterLink>
      </header>

      <!-- Stepper -->
      <div class="stepper-container card">
        <div class="stepper">
          <div
            v-for="(step, index) in steps"
            :key="index"
            class="step"
            :class="{
              active: currentStep === index,
              completed: currentStep > index,
            }"
          >
            <div class="step-indicator">
              <Check v-if="currentStep > index" :size="16" class="check-icon" />
              <span v-else class="step-number">{{ index + 1 }}</span>
            </div>
            <div class="step-label">
              <span class="step-title">{{ step.title }}</span>
              <span class="step-description">{{ step.description }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Form Content -->
      <div class="form-container card">
        <!-- Step 0: Basic Info -->
        <ProductFormBasic
          v-if="currentStep === 0"
          v-model="formData.basic"
          @next="handleNext"
          @cancel="handleCancel"
        />

        <!-- Step 1: Classification -->
        <ProductFormClassification
          v-if="currentStep === 1"
          v-model="formData.classification"
          @next="handleNext"
          @prev="handlePrev"
        />

        <!-- Step 2: Attributes -->
        <ProductFormAttributes
          v-if="currentStep === 2"
          v-model="formData.attributes"
          :brand-id="formData.classification.brandId"
          :group-ids="formData.classification.groupIds"
          @next="handleNext"
          @prev="handlePrev"
        />

        <!-- Step 3: Preview -->
        <ProductFormPreview
          v-if="currentStep === 3"
          :form-data="allFormData"
          :brand-name="brandName"
          :group-names="groupNames"
          :attribute-names="attributeNames"
          :is-submitting="isSubmitting"
          @edit="goToStep"
          @prev="handlePrev"
          @submit="handleSubmit"
        />
      </div>

      <!-- Error Message -->
      <div v-if="errorMessage" class="alert alert-error">
        <X :size="20" class="alert-icon" />
        <div class="alert-content">
          <strong>Error al crear el producto</strong>
          <p>{{ errorMessage }}</p>
        </div>
        <button @click="errorMessage = ''" class="alert-close">&times;</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import ProductFormBasic from '@/components/product/ProductFormBasic.vue'
import ProductFormClassification from '@/components/product/ProductFormClassification.vue'
import ProductFormAttributes from '@/components/product/ProductFormAttributes.vue'
import ProductFormPreview from '@/components/product/ProductFormPreview.vue'
import { productApi } from '@/services/productApi'
import { Check, X } from 'lucide-vue-next'

const router = useRouter()

// Stepper configuration
const steps = [
  {
    title: 'Información básica',
    description: 'Datos principales',
  },
  {
    title: 'Clasificación',
    description: 'Marca y categorías',
  },
  {
    title: 'Atributos',
    description: 'Características configurables',
  },
  {
    title: 'Preview',
    description: 'Revisar y confirmar',
  },
]

// Current step
const currentStep = ref(0)

// Form data
const formData = reactive({
  basic: {
    productType: '',
    sku: '',
    name: '',
    longName: '',
    description: '',
    basePrice: 0,
  },
  classification: {
    brandId: '',
    groupIds: [],
  },
  attributes: {
    directAttributeIds: [],
  },
})

// State
const isSubmitting = ref(false)
const errorMessage = ref('')

// Cached data for preview
const brands = ref([])
const groups = ref([])
const attributes = ref([])

// Computed: Total attribute count (inherited + direct)
// For MVP, we only count direct attributes since inheritance is pending
const totalAttributeCount = computed(() => {
  return formData.attributes.directAttributeIds.length
})

// Computed: Determine variant strategy automatically based on architecture
const variantStrategy = computed(() => {
  // Rule: If product has no attributes → No variants (simple product)
  //       If product has attributes → JIT + Manual
  return totalAttributeCount.value > 0 ? 'jit' : 'none'
})

// Computed: All form data combined
const allFormData = computed(() => ({
  productType: formData.basic.productType,
  sku: formData.basic.sku,
  name: formData.basic.name,
  longName: formData.basic.longName,
  description: formData.basic.description,
  basePrice: formData.basic.basePrice,
  brandId: formData.classification.brandId,
  groupIds: formData.classification.groupIds,
  directAttributeIds: formData.attributes.directAttributeIds,
  strategy: variantStrategy.value, // Automatically determined
}))

// Get brand name by ID
const brandName = computed(() => {
  if (!formData.classification.brandId) return ''
  const brand = brands.value.find(b => b.id === formData.classification.brandId)
  return brand ? brand.name : ''
})

// Get group names by IDs
const groupNames = computed(() => {
  if (!formData.classification.groupIds || formData.classification.groupIds.length === 0) {
    return []
  }
  return formData.classification.groupIds
    .map(id => {
      const group = groups.value.find(g => g.id === id)
      return group ? group.name : null
    })
    .filter(Boolean)
})

// Get attribute names by IDs
const attributeNames = computed(() => {
  if (!formData.attributes.directAttributeIds || formData.attributes.directAttributeIds.length === 0) {
    return []
  }
  return formData.attributes.directAttributeIds
    .map(id => {
      const attr = attributes.value.find(a => a.id === id)
      return attr ? attr.name : null
    })
    .filter(Boolean)
})

// Navigation handlers
function handleNext() {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function handlePrev() {
  if (currentStep.value > 0) {
    currentStep.value--
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function goToStep(stepIndex) {
  if (stepIndex >= 0 && stepIndex < steps.length) {
    currentStep.value = stepIndex
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function handleCancel() {
  router.push('/products')
}

// Load cached data for preview
async function loadCachedData() {
  try {
    const [brandsRes, groupsRes, attrsRes] = await Promise.all([
      productApi.listBrands({ isActive: true }),
      productApi.listProductGroups({ isActive: true }),
      productApi.listAttributes({}),
    ])
    brands.value = brandsRes.data
    groups.value = groupsRes.data
    attributes.value = attrsRes.data
  } catch (error) {
    console.error('Error loading cached data:', error)
  }
}

// Submit handler
async function handleSubmit() {
  isSubmitting.value = true
  errorMessage.value = ''

  try {
    // Validate required fields
    const emptyUuid = '00000000-0000-0000-0000-000000000000'
    console.log('[Create] Validating brandId:', formData.classification.brandId)
    if (!formData.classification.brandId || formData.classification.brandId === emptyUuid) {
      errorMessage.value = 'La marca es obligatoria. Por favor, vuelve al paso de Clasificación y selecciona una marca.'
      window.scrollTo({ top: 0, behavior: 'smooth' })
      isSubmitting.value = false
      return
    }

    // Create the product - Only send brandId if it exists
    const productPayload = {
      sku: formData.basic.sku,
      name: formData.basic.name,
      longName: formData.basic.longName || undefined,
      description: formData.basic.description || undefined,
      productType: formData.basic.productType,
      brandId: formData.classification.brandId, // Already validated above
      groupIds: formData.classification.groupIds || [],
      directAttributeIds: formData.attributes.directAttributeIds || [],
      base_price: formData.basic.basePrice || 0,
      isActive: true,
    }

    console.log('[Create] Sending payload:', productPayload)

    const createdProduct = await productApi.createProduct(productPayload)

    // Note: Variants will be created Just-in-Time when added to orders
    // or manually from the product detail view. No automatic generation here.

    // Success: Show message and navigate to product list
    console.log('[Create] Product created successfully:', createdProduct)
    alert(`✅ Producto "${createdProduct.name}" creado exitosamente con SKU: ${createdProduct.sku}`)
    await router.push('/products')
  } catch (error) {
    console.error('Error creating product:', error)
    errorMessage.value = error.message || 'Ocurrió un error inesperado al crear el producto.'
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } finally {
    isSubmitting.value = false
  }
}

// Load cached data on mount
loadCachedData()
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
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

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
}

.stepper-container {
  padding: 1.25rem;
}

.stepper {
  display: flex;
  justify-content: space-between;
  position: relative;
}

.stepper::before {
  content: '';
  position: absolute;
  top: 20px;
  left: 5%;
  right: 5%;
  height: 2px;
  background-color: #e2e8f0;
  z-index: 0;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
  position: relative;
  z-index: 1;
}

.step-indicator {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #ffffff;
  border: 2px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
  color: #94a3b8;
  transition: all 0.3s ease;
}

.step.active .step-indicator {
  background-color: #f4d03f;
  border-color: #f4d03f;
  color: #1e293b;
  box-shadow: 0 0 0 4px rgba(244, 208, 63, 0.2);
}

.step.completed .step-indicator {
  background-color: #10b981;
  border-color: #10b981;
  color: #ffffff;
}

.check-icon {
  font-size: 1.2rem;
}

.step-label {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  max-width: 120px;
}

.step-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: #1e293b;
}

.step-description {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 0.15rem;
}

.step.active .step-title {
  color: #1b3a6b;
}

.form-container {
  padding: 2rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover {
  background: #f8fafc;
}

.alert {
  display: flex;
  gap: 1rem;
  padding: 1rem 1.25rem;
  border-radius: 8px;
  align-items: flex-start;
}

.alert-error {
  background-color: #fef2f2;
  border: 1px solid #fecaca;
}

.alert-icon {
  font-size: 1.25rem;
  color: #ef4444;
  flex-shrink: 0;
}

.alert-content {
  flex: 1;
}

.alert-content strong {
  display: block;
  color: #991b1b;
  margin-bottom: 0.25rem;
  font-size: 0.95rem;
}

.alert-content p {
  margin: 0;
  color: #dc2626;
  font-size: 0.9rem;
  line-height: 1.5;
}

.alert-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #dc2626;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  flex-shrink: 0;
}

.alert-close:hover {
  color: #991b1b;
}

@media (max-width: 768px) {
  .dashboard-content {
    padding: 1.5rem;
  }

  .stepper {
    flex-direction: column;
    gap: 1rem;
  }

  .stepper::before {
    display: none;
  }

  .step {
    flex-direction: row;
    justify-content: flex-start;
    align-items: center;
  }

  .step-label {
    align-items: flex-start;
    text-align: left;
    max-width: none;
  }

  .form-container {
    padding: 1.5rem;
  }
}
</style>
