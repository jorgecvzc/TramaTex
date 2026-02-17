<template>
  <div class="form-step">
    <h3 class="step-title">Atributos configurables</h3>
    <p class="step-description">
      Asigna atributos específicos a este producto (Talla, Color, Material, etc.)
    </p>

    <form @submit.prevent="handleNext" class="step-form">
      <!-- Info Box -->
      <div class="info-box">
        <span class="info-icon">ℹ️</span>
        <div class="info-content">
          <p>
            <strong>¿Qué son los atributos?</strong> Los atributos definen las características configurables
            del producto (ej. Talla: S, M, L). Cada combinación de valores de atributo genera una variante.
          </p>
          <p>
            Los atributos se heredan desde la <strong>Marca</strong> y las <strong>Categorías</strong>,
            pero puedes asignar atributos adicionales directos a este producto.
          </p>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoadingAttributes" class="loading-state">
        <div class="spinner"></div>
        <p>Cargando atributos disponibles...</p>
      </div>

      <!-- Error State -->
      <div v-if="loadError" class="alert-error">
        {{ loadError }}
      </div>

      <!-- Attributes Selection -->
      <div v-if="!isLoadingAttributes && !loadError" class="attributes-section">
        <!-- Generic Attributes -->
        <div v-if="genericAttributes.length > 0" class="attribute-category">
          <h4 class="category-title">
            Atributos genéricos
            <span class="category-badge">{{ genericAttributes.length }}</span>
          </h4>
          <p class="category-description">
            Atributos disponibles para todos los productos.
          </p>
          <div class="checkbox-group">
            <label
              v-for="attr in genericAttributes"
              :key="attr.id"
              class="checkbox-label"
            >
              <input
                type="checkbox"
                :value="attr.id"
                v-model="localData.directAttributeIds"
              />
              <span class="checkbox-text">
                <strong>{{ attr.name }}</strong>
                <code class="attr-code">{{ attr.code }}</code>
                <span v-if="attr.values && attr.values.length > 0" class="attr-values">
                  ({{ attr.values.length }} valores)
                </span>
              </span>
            </label>
          </div>
        </div>

        <!-- Brand-Specific Attributes -->
        <div v-if="brandAttributes.length > 0" class="attribute-category">
          <h4 class="category-title">
            Atributos de la marca
            <span class="category-badge brand">{{ brandAttributes.length }}</span>
          </h4>
          <p class="category-description">
            Atributos específicos de la marca seleccionada (se heredarán automáticamente).
          </p>
          <div class="inherited-list">
            <div
              v-for="attr in brandAttributes"
              :key="attr.id"
              class="inherited-item"
            >
              <span class="inherited-icon">🔗</span>
              <span class="inherited-text">
                <strong>{{ attr.name }}</strong>
                <code class="attr-code">{{ attr.code }}</code>
                <span v-if="attr.values && attr.values.length > 0" class="attr-values">
                  ({{ attr.values.length }} valores)
                </span>
              </span>
            </div>
          </div>
        </div>

        <!-- Group-Specific Attributes -->
        <div v-if="groupAttributes.length > 0" class="attribute-category">
          <h4 class="category-title">
            Atributos de las categorías
            <span class="category-badge group">{{ groupAttributes.length }}</span>
          </h4>
          <p class="category-description">
            Atributos específicos de las categorías seleccionadas (se heredarán automáticamente).
          </p>
          <div class="inherited-list">
            <div
              v-for="attr in groupAttributes"
              :key="attr.id"
              class="inherited-item"
            >
              <span class="inherited-icon">🔗</span>
              <span class="inherited-text">
                <strong>{{ attr.name }}</strong>
                <code class="attr-code">{{ attr.code }}</code>
                <span v-if="attr.values && attr.values.length > 0" class="attr-values">
                  ({{ attr.values.length }} valores)
                </span>
              </span>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="allAttributes.length === 0" class="empty-state">
          <span class="empty-icon">📦</span>
          <p>No hay atributos disponibles.</p>
          <p class="empty-hint">
            Puedes crear atributos en la sección de "Gestión de Atributos" y luego asignarlos a este producto.
          </p>
        </div>
      </div>

      <!-- Selected Summary -->
      <div v-if="localData.directAttributeIds.length > 0" class="summary-box">
        <h4>Atributos directos seleccionados</h4>
        <div class="summary-tags">
          <span
            v-for="id in localData.directAttributeIds"
            :key="id"
            class="tag"
          >
            {{ getAttributeName(id) }}
          </span>
        </div>
      </div>

      <!-- Variant Strategy Info -->
      <div class="variant-info-box">
        <div class="variant-info-header">
          <span class="variant-info-icon">🔄</span>
          <strong>Gestión de Variantes</strong>
        </div>
        <div v-if="hasAttributes" class="variant-info-content jit">
          <p>
            <strong>⚡ Producto Configurable:</strong> Este producto tendrá variantes que se gestionarán automáticamente.
          </p>
          <ul>
            <li><strong>Just-in-Time (JIT):</strong> Al añadir a un pedido, las variantes se crearán automáticamente con estado PROVISIONAL</li>
            <li><strong>Creación Manual:</strong> Desde la vista de detalle del producto podrás crear variantes específicas manualmente</li>
            <li>Solo se generarán las combinaciones que realmente se utilicen</li>
          </ul>
        </div>
        <div v-else class="variant-info-content simple">
          <p>
            <strong>📦 Producto Simple:</strong> Sin atributos configurables, este producto no tendrá variantes.
          </p>
          <ul>
            <li>Se podrá añadir directamente a pedidos usando su SKU base</li>
            <li>Ideal para productos únicos o servicios sin configuración</li>
          </ul>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button
          type="button"
          @click="$emit('prev')"
          class="btn btn-secondary"
        >
          ← Anterior
        </button>
        <button
          type="submit"
          class="btn btn-primary"
        >
          Siguiente: Preview →
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted, watch, ref } from 'vue'
import { productApi } from '@/services/productApi'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
  brandId: {
    type: String,
    default: '',
  },
  groupIds: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'next', 'prev'])

// Local copy of data
const localData = reactive({
  directAttributeIds: props.modelValue.directAttributeIds || [],
})

// State
const allAttributes = reactive([])
const isLoadingAttributes = ref(false)
const loadError = ref('')

// Watch local data changes and emit to parent
watch(localData, (newValue) => {
  emit('update:modelValue', { ...newValue })
}, { deep: true })

// Computed: Categorize attributes by scope
const genericAttributes = computed(() => {
  return allAttributes.filter(attr =>
    !attr.scope_brand_id && !attr.scope_group_id
  )
})

const brandAttributes = computed(() => {
  return allAttributes.filter(attr =>
    attr.scope_brand_id && !attr.scope_group_id
  )
})

const groupAttributes = computed(() => {
  return allAttributes.filter(attr =>
    !attr.scope_brand_id && attr.scope_group_id
  )
})

// Computed: Check if product has any attributes (inherited or direct)
const hasAttributes = computed(() => {
  return localData.directAttributeIds.length > 0 ||
         brandAttributes.value.length > 0 ||
         groupAttributes.value.length > 0
})

// Get attribute name by ID
function getAttributeName(id) {
  const attr = allAttributes.find(a => a.id === id)
  return attr ? attr.name : id
}

// Handle next step
function handleNext() {
  emit('next')
}

// Load attributes from API
async function loadAttributes() {
  isLoadingAttributes.value = true
  loadError.value = ''
  allAttributes.length = 0

  try {
    // Load all attributes (generic + brand + groups)
    const promises = [
      productApi.listAttributes({}), // Generic
    ]

    // Add brand-specific if brand is selected
    if (props.brandId) {
      promises.push(productApi.listAttributes({ scopeBrandId: props.brandId }))
    }

    // Add group-specific for each selected group
    if (props.groupIds && props.groupIds.length > 0) {
      props.groupIds.forEach(groupId => {
        promises.push(productApi.listAttributes({ scopeGroupId: groupId }))
      })
    }

    const results = await Promise.all(promises)

    // Merge all results and deduplicate by ID
    const attributeMap = new Map()
    results.forEach(result => {
      result.data.forEach(attr => {
        if (!attributeMap.has(attr.id)) {
          attributeMap.set(attr.id, attr)
        }
      })
    })

    allAttributes.push(...attributeMap.values())
  } catch (error) {
    loadError.value = 'No se pudieron cargar los atributos. Intenta recargar la página.'
    console.error('Error loading attributes:', error)
  } finally {
    isLoadingAttributes.value = false
  }
}

// Load data on mount
onMounted(() => {
  loadAttributes()
})
</script>

<style scoped>
.form-step {
  max-width: 800px;
  margin: 0 auto;
}

.step-title {
  color: #1b3a6b;
  font-size: 1.5rem;
  margin: 0 0 0.5rem;
}

.step-description {
  color: #64748b;
  margin: 0 0 2rem;
  font-size: 0.95rem;
}

.step-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-box {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background-color: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
}

.info-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.info-content {
  font-size: 0.9rem;
  color: #1e40af;
  line-height: 1.5;
}

.info-content p {
  margin: 0 0 0.5rem;
}

.info-content p:last-child {
  margin-bottom: 0;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  text-align: center;
  color: #64748b;
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e2e8f0;
  border-top-color: #1b3a6b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.alert-error {
  padding: 1rem;
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #991b1b;
  font-size: 0.9rem;
}

.attributes-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.attribute-category {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1rem;
}

.category-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
  color: #1b3a6b;
  margin: 0 0 0.5rem;
}

.category-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 0.5rem;
  font-size: 0.75rem;
  font-weight: 600;
  background-color: #e2e8f0;
  color: #475569;
  border-radius: 12px;
}

.category-badge.brand {
  background-color: #dbeafe;
  color: #1e40af;
}

.category-badge.group {
  background-color: #dcfce7;
  color: #166534;
}

.category-description {
  font-size: 0.85rem;
  color: #64748b;
  margin: 0 0 1rem;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.checkbox-label:hover {
  background-color: #f8fafc;
  border-color: #cbd5e1;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
  flex-shrink: 0;
}

.checkbox-text {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  color: #1e293b;
  flex-wrap: wrap;
}

.attr-code {
  font-size: 0.75rem;
  padding: 0.15rem 0.4rem;
  background-color: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  color: #475569;
}

.attr-values {
  font-size: 0.8rem;
  color: #64748b;
}

.inherited-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.inherited-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
}

.inherited-icon {
  font-size: 1rem;
  flex-shrink: 0;
}

.inherited-text {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  color: #1e293b;
  flex-wrap: wrap;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 2rem;
  text-align: center;
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
}

.empty-icon {
  font-size: 3rem;
  opacity: 0.5;
}

.empty-state p {
  margin: 0;
  color: #64748b;
  font-size: 0.95rem;
}

.empty-hint {
  font-size: 0.85rem !important;
  color: #94a3b8 !important;
}

.summary-box {
  background-color: #fefce8;
  border: 1px solid #fde047;
  border-radius: 8px;
  padding: 1rem;
}

.summary-box h4 {
  margin: 0 0 0.75rem;
  font-size: 0.9rem;
  color: #854d0e;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.summary-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 0.75rem;
  background-color: #fef3c7;
  border: 1px solid #fde047;
  border-radius: 6px;
  font-size: 0.85rem;
  color: #78350f;
  font-weight: 500;
}

/* Variant Info Box */
.variant-info-box {
  border-radius: 8px;
  padding: 1.25rem;
  border: 2px solid;
}

.variant-info-box:has(.jit) {
  background-color: #fffbeb;
  border-color: #fde047;
}

.variant-info-box:has(.simple) {
  background-color: #f8fafc;
  border-color: #cbd5e1;
}

.variant-info-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  font-size: 1rem;
  color: #1e293b;
}

.variant-info-icon {
  font-size: 1.25rem;
}

.variant-info-content p {
  margin: 0 0 0.75rem;
  font-size: 0.9rem;
  color: #475569;
  line-height: 1.5;
}

.variant-info-content strong {
  color: #1e293b;
}

.variant-info-content ul {
  margin: 0;
  padding-left: 1.5rem;
  list-style-type: disc;
}

.variant-info-content li {
  margin-bottom: 0.5rem;
  font-size: 0.85rem;
  color: #64748b;
  line-height: 1.5;
}

.variant-info-content li:last-child {
  margin-bottom: 0;
}

.variant-info-content.jit {
  color: #1e40af;
}

.variant-info-content.simple {
  color: #475569;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: space-between;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
  margin-top: 1rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.65rem 1.25rem;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #f4d03f;
  color: #1e293b;
  box-shadow: 0 2px 4px rgba(244, 208, 63, 0.3);
}

.btn-primary:not(:disabled):hover {
  background: #f0c929;
  box-shadow: 0 4px 8px rgba(244, 208, 63, 0.4);
  transform: translateY(-1px);
}

.btn-primary:not(:disabled):active {
  transform: translateY(0);
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover {
  background: #f8fafc;
}

@media (max-width: 768px) {
  .form-step {
    max-width: 100%;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .btn {
    width: 100%;
  }
}
</style>
