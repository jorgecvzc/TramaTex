<template>
  <div class="dashboard">
    
    <BaseFormLayout
      title="Nuevo Producto"
      :breadcrumbs="[{ label: 'Operaciones', to: '/products' }, { label: 'Productos', to: '/products' }, { label: 'Crear Nuevo' }]"
      :is-submitting="isSubmitting"
      catalog-route="/products"
      @submit="nextStep"
    >
      <!-- Stepper Progress Header -->
      <section class="card stepper-header mb-6">
        <div class="stepper-steps">
          <div v-for="step in 3" :key="step" :class="['step-item', { active: currentStep === step, completed: currentStep > step }]">
            <div class="step-number">
              <span v-if="currentStep > step" class="material-symbols-outlined">check</span>
              <span v-else>{{ step }}</span>
            </div>
            <span class="step-label">{{ getStepLabel(step) }}</span>
          </div>
        </div>
      </section>

      <!-- Step 1: Basic Information -->
      <div v-if="currentStep === 1" class="step-content">
        <FormSection title="Información Básica" icon="inventory_2" description="Datos principales que identifican el producto en el catálogo.">
          <div class="form-row">
            <div class="form-group">
              <label>Nombre del Producto *</label>
              <input v-model="formData.name" type="text" placeholder="Ej: Camiseta Técnica Algodón" required />
            </div>
            <div class="form-group">
              <label>Referencia / SKU *</label>
              <input v-model="formData.sku" type="text" placeholder="Ej: CAM-001" required />
            </div>
          </div>
          <div class="form-group mt-4">
            <label>Descripción Larga (Comercial)</label>
            <textarea v-model="formData.long_name" rows="2" placeholder="Nombre completo para documentos oficiales..."></textarea>
          </div>
        </FormSection>

        <FormSection title="Clasificación y Tipo" icon="category">
          <div class="form-row">
            <div class="form-group">
              <label>Tipo de Producto</label>
              <select v-model="formData.product_type">
                <option value="TANGIBLE">Producto Físico (Tangible)</option>
                <option value="SERVICE">Servicio / Mano de Obra</option>
              </select>
            </div>
            <div class="form-group">
              <label>Marca</label>
              <select v-model="formData.brand_id">
                <option :value="null">Sin marca específica</option>
                <option v-for="brand in brands" :key="brand.id" :value="brand.id">{{ brand.name }}</option>
              </select>
            </div>
          </div>
          <div class="form-group mt-4">
            <label>Categorías</label>
            <div class="category-grid">
              <label v-for="group in productGroups" :key="group.id" class="checkbox-pill">
                <input type="checkbox" :value="group.id" v-model="formData.group_ids" />
                <span>{{ group.name }}</span>
              </label>
            </div>
          </div>
        </FormSection>
      </div>

      <!-- Step 2: Variants & Attributes -->
      <div v-if="currentStep === 2" class="step-content">
        <FormSection title="Configuración de Variantes" icon="layers" description="Define las combinaciones de atributos (tallas, colores, etc.) para este producto.">
          <div class="alert-info mb-4">
            <span class="material-symbols-outlined">info</span>
            <p>Se generará automáticamente una variante base. Podrás añadir más atributos para crear combinaciones múltiples.</p>
          </div>
          <!-- Simplified variant logic for the refactor, keeping original intent -->
          <div class="table-wrapper">
            <table class="data-table">
              <thead>
                <tr>
                  <th>SKU Variante</th>
                  <th>Atributos</th>
                  <th class="align-right">Precio Base</th>
                  <th class="text-center">Estado</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><code class="code-badge">{{ formData.sku }}-BASE</code></td>
                  <td><span class="text-muted">Variante inicial predeterminada</span></td>
                  <td class="align-right">0.00 €</td>
                  <td class="text-center"><span class="status-pill active">Activa</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </FormSection>
      </div>

      <!-- Step 3: Review -->
      <div v-if="currentStep === 3" class="step-content">
        <FormSection title="Resumen del Nuevo Producto" icon="fact_check">
          <div class="review-grid">
            <div class="review-item">
              <label>Nombre</label>
              <p>{{ formData.name }}</p>
            </div>
            <div class="review-item">
              <label>SKU Principal</label>
              <p><code>{{ formData.sku }}</code></p>
            </div>
            <div class="review-item">
              <label>Tipo</label>
              <p>{{ formData.product_type }}</p>
            </div>
            <div class="review-item">
              <label>Categorías</label>
              <p>{{ formData.group_ids.length }} seleccionadas</p>
            </div>
          </div>
        </FormSection>
      </div>

      <!-- Custom Footer Actions for Stepper -->
      <template #actions>
        <button type="button" class="btn btn-outline btn-lg" @click="prevStep" v-if="currentStep > 1">
          <span class="material-symbols-outlined">arrow_back</span>
          Anterior
        </button>
        <button type="button" class="btn btn-outline btn-lg" @click="router.push('/products')" v-if="currentStep === 1">
          Cancelar
        </button>
        <button type="submit" class="btn btn-primary btn-lg">
          <span>{{ currentStep === 3 ? 'Finalizar y Crear' : 'Siguiente Paso' }}</span>
          <span class="material-symbols-outlined">{{ currentStep === 3 ? 'check_circle' : 'arrow_forward' }}</span>
        </button>
      </template>
    </BaseFormLayout>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseFormLayout from '@/components/shared/BaseFormLayout.vue'
import FormSection from '@/components/shared/FormSection.vue'
import { productApi } from '@/services/productApi'

const router = useRouter()
const currentStep = ref(1)
const isSubmitting = ref(false)

const brands = ref([])
const productGroups = ref([])

const formData = reactive({
  name: '',
  sku: '',
  long_name: '',
  product_type: 'TANGIBLE',
  brand_id: null,
  group_ids: [],
  is_active: true
})

onMounted(async () => {
  try {
    const [bRes, gRes] = await Promise.all([
      productApi.listBrands({ isActive: true }),
      productApi.listProductGroups({ isActive: true })
    ])
    brands.value = bRes.data || []
    productGroups.value = gRes.data || []
  } catch (err) {}
})

function getStepLabel(step) {
  const labels = ['', 'Datos Generales', 'Variantes y Precios', 'Revisión Final']
  return labels[step]
}

function nextStep() {
  if (currentStep.value < 3) {
    currentStep.value++
    window.scrollTo(0, 0)
  } else {
    createProduct()
  }
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
    window.scrollTo(0, 0)
  }
}

async function createProduct() {
  isSubmitting.value = true
  try {
    const res = await productApi.createProduct(formData)
    router.push(`/products/${res.id}`)
  } catch (err) {
    alert(err.message || 'Error al crear el producto')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
/* Stepper UI */
.stepper-header { padding: 1.5rem 2rem; background: var(--color-background); }
.stepper-steps { display: flex; justify-content: space-between; position: relative; max-width: 800px; margin: 0 auto; }
.stepper-steps::before { content: ''; position: absolute; top: 20px; left: 0; right: 0; height: 2px; background: var(--color-border); z-index: 1; }

.step-item { position: relative; z-index: 2; display: flex; flex-direction: column; align-items: center; gap: 0.75rem; flex: 1; }
.step-number { 
  width: 40px; height: 40px; border-radius: 50%; background: white; border: 2px solid var(--color-border); 
  display: flex; align-items: center; justify-content: center; font-weight: 700; color: var(--color-text-secondary); transition: 0.3s;
}
.step-label { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; color: var(--color-text-secondary); }

.step-item.active .step-number { border-color: var(--color-primary); background: var(--color-primary); color: var(--color-text-on-primary); transform: scale(1.1); }
.step-item.active .step-label { color: var(--color-text-primary); }
.step-item.completed .step-number { border-color: var(--color-success); background: var(--color-success); color: white; }

/* Form Helper Classes */
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
input, select, textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }

.category-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 0.75rem; }
.checkbox-pill { 
  display: flex; align-items: center; gap: 0.5rem; padding: 0.5rem 1rem; border: 1px solid var(--color-border); 
  border-radius: 99px; cursor: pointer; transition: 0.2s; background: white; font-size: 0.85rem;
}
.checkbox-pill:has(input:checked) { border-color: var(--color-primary); background: rgba(230, 184, 0, 0.05); }

.review-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 2rem; }
.review-item label { display: block; font-size: 0.7rem; font-weight: 700; color: var(--color-text-secondary); text-transform: uppercase; margin-bottom: 0.5rem; }
.review-item p { font-size: 1.1rem; font-weight: 600; margin: 0; }

.alert-info { display: flex; align-items: center; gap: 0.75rem; padding: 1rem; background: rgba(59, 130, 246, 0.05); border: 1px solid rgba(59, 130, 246, 0.2); border-radius: 8px; }
.alert-info p { margin: 0; font-size: 0.9rem; color: #1e40af; }

.code-badge { background: var(--color-background); padding: 0.2rem 0.5rem; border-radius: 4px; font-family: var(--font-family-mono); }
</style>