<template>
  <div class="product-info">
    <div class="info-header">
      <h3>Información General</h3>
      <button
        v-if="!isEditing"
        @click="startEditing"
        class="btn btn-primary"
      >
        ✎ Editar producto
      </button>
    </div>

    <!-- View Mode -->
    <div v-if="!isEditing" class="info-grid">
      <div class="info-item">
        <label>Nombre</label>
        <p>{{ product.name }}</p>
      </div>

      <div class="info-item">
        <label>Nombre largo</label>
        <p>{{ product.long_name || '—' }}</p>
      </div>

      <div class="info-item">
        <label>SKU Base</label>
        <p>
          <code v-if="product.sku" class="sku-code">{{ product.sku }}</code>
          <span v-else>—</span>
        </p>
      </div>

      <div class="info-item">
        <label>Código de barras</label>
        <p>
          <code v-if="product.barcode" class="barcode">{{ product.barcode }}</code>
          <span v-else>—</span>
        </p>
      </div>

      <div class="info-item">
        <label>Precio base (coste)</label>
        <p>
          <strong class="price-value">{{ formatPrice(product.base_price) }} €</strong>
          <br>
          <small class="hint">Coste base para cálculo de precios de venta</small>
        </p>
      </div>

      <div class="info-item">
        <label>Tipo de producto</label>
        <p>
          <span class="pill" :class="`type-${product.product_type.toLowerCase()}`">
            {{ product.product_type === 'TANGIBLE' ? 'Tangible' : 'Servicio' }}
          </span>
        </p>
      </div>

      <div class="info-item">
        <label>Estado</label>
        <p>
          <span class="pill" :class="product.is_active ? 'active' : 'inactive'">
            {{ product.is_active ? 'Activo' : 'Inactivo' }}
          </span>
          <button
            @click="$emit('toggle-status')"
            class="btn-link"
          >
            {{ product.is_active ? 'Desactivar' : 'Activar' }}
          </button>
        </p>
      </div>

      <div class="info-item">
        <label>Marca</label>
        <p>
          <span v-if="brand" class="brand-name">{{ brand.name }}</span>
          <span v-else class="text-muted">Sin marca</span>
        </p>
      </div>

      <div class="info-item">
        <label>Categorías</label>
        <p v-if="groups && groups.length > 0">
          <span
            v-for="group in groups"
            :key="group.id"
            class="category-tag"
          >
            {{ group.name }}
          </span>
        </p>
        <p v-else class="text-muted">Sin categorías</p>
      </div>

      <div class="info-item full-width">
        <label>Descripción</label>
        <p class="description">{{ product.description || '—' }}</p>
      </div>

      <div class="info-item">
        <label>Fecha de creación</label>
        <p>{{ formatDate(product.created_at) }}</p>
      </div>

      <div class="info-item">
        <label>Última modificación</label>
        <p>{{ formatDate(product.modified_at) }}</p>
      </div>
    </div>

    <!-- Edit Mode -->
    <div v-if="isEditing" class="edit-form">
      <div class="form-row">
        <div class="form-group">
          <label for="edit-name">Nombre *</label>
          <input
            id="edit-name"
            v-model="editForm.name"
            type="text"
            placeholder="Nombre corto del producto"
            required
          />
        </div>

        <div class="form-group">
          <label for="edit-sku">SKU Base</label>
          <input
            id="edit-sku"
            v-model="editForm.sku"
            type="text"
            placeholder="Ej: FYR2040"
          />
        </div>
      </div>

      <div class="form-group">
        <label for="edit-long-name">Nombre largo</label>
        <input
          id="edit-long-name"
          v-model="editForm.long_name"
          type="text"
          placeholder="Nombre completo para facturas y presupuestos"
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="edit-barcode">Código de barras (EAN/UPC)</label>
          <input
            id="edit-barcode"
            v-model="editForm.barcode"
            type="text"
            placeholder="8412345678901"
          />
        </div>

        <div class="form-group">
          <label for="edit-base-price">Precio base (coste) *</label>
          <input
            id="edit-base-price"
            v-model.number="editForm.base_price"
            type="number"
            step="0.01"
            min="0"
            placeholder="0.00"
            required
          />
          <small class="hint">Coste base del producto (EUR)</small>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="edit-tax-rate">Tasa de IVA (%)</label>
          <select id="edit-tax-rate" v-model.number="editForm.tax_rate">
            <option :value="21">21% — General</option>
            <option :value="10">10% — Reducido</option>
            <option :value="4">4% — Superreducido</option>
            <option :value="0">0% — Exento</option>
          </select>
        </div>

        <div class="form-group">
          <label for="edit-type">Tipo de producto</label>
          <select id="edit-type" v-model="editForm.product_type">
            <option value="TANGIBLE">Tangible</option>
            <option value="SERVICE">Servicio</option>
          </select>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="edit-brand">Marca *</label>
          <select id="edit-brand" v-model="editForm.brand_id" :disabled="isLoadingOptions">
            <option value="">Seleccionar marca...</option>
            <option v-for="b in availableBrands" :key="b.id" :value="b.id">
              {{ b.name }}
            </option>
          </select>
          <small v-if="isLoadingOptions" class="hint">Cargando marcas...</small>
        </div>
      </div>

      <div class="form-group">
        <label>Categorías</label>
        <div v-if="isLoadingOptions" class="hint">Cargando categorías...</div>
        <div v-else-if="availableGroups.length === 0" class="hint">No hay categorías disponibles</div>
        <div v-else class="checkbox-group">
          <label
            v-for="group in availableGroups"
            :key="group.id"
            class="checkbox-option"
          >
            <input
              type="checkbox"
              :value="group.id"
              v-model="editForm.group_ids"
            />
            <span>{{ group.name }}</span>
            <span v-if="group.group_type" class="pill-small">{{ group.group_type }}</span>
          </label>
        </div>
      </div>

      <div class="form-group">
        <label for="edit-description">Descripción</label>
        <textarea
          id="edit-description"
          v-model="editForm.description"
          rows="4"
          placeholder="Descripción detallada del producto..."
        ></textarea>
      </div>

      <div class="form-actions">
        <button
          @click="submitEdit"
          :disabled="isUpdating || !isFormValid"
          class="btn btn-primary"
        >
          {{ isUpdating ? 'Guardando...' : 'Guardar cambios' }}
        </button>
        <button
          @click="cancelEdit"
          :disabled="isUpdating"
          class="btn btn-secondary"
        >
          Cancelar
        </button>
      </div>

      <p v-if="editError" class="error-message">{{ editError }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { productApi } from '@/services/productApi'

const props = defineProps({
  product: {
    type: Object,
    required: true,
  },
  brand: {
    type: Object,
    default: null,
  },
  groups: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update', 'toggle-status'])

// State
const isEditing = ref(false)
const isUpdating = ref(false)
const editError = ref('')
const availableBrands = ref([])
const availableGroups = ref([])
const isLoadingOptions = ref(false)

const editForm = reactive({
  name: '',
  long_name: '',
  sku: '',
  barcode: '',
  base_price: 0,
  tax_rate: 21,
  product_type: 'TANGIBLE',
  description: '',
  brand_id: '',
  group_ids: [],
})

// Computed
const isFormValid = computed(() => {
  return editForm.name.trim().length > 0 && editForm.base_price >= 0
})

// Methods
async function startEditing() {
  editForm.name = props.product.name || ''
  editForm.long_name = props.product.long_name || ''
  editForm.sku = props.product.sku || ''
  editForm.barcode = props.product.barcode || ''
  editForm.base_price = props.product.base_price ?? 0
  editForm.tax_rate = props.product.tax_rate ?? 21
  editForm.product_type = props.product.product_type || 'TANGIBLE'
  editForm.description = props.product.description || ''
  editForm.brand_id = props.product.brand_id || ''
  editForm.group_ids = [...(props.product.group_ids || [])]
  isEditing.value = true
  editError.value = ''

  // Load available brands and groups
  isLoadingOptions.value = true
  try {
    const [brandsResult, groupsResult] = await Promise.all([
      productApi.listBrands(),
      productApi.listProductGroups(),
    ])
    availableBrands.value = brandsResult.data || []
    availableGroups.value = groupsResult.data || []
  } catch (err) {
    console.error('Error loading brands/groups:', err)
  } finally {
    isLoadingOptions.value = false
  }
}

function cancelEdit() {
  isEditing.value = false
  editError.value = ''
}

async function submitEdit() {
  if (!isFormValid.value) {
    editError.value = 'El nombre es obligatorio'
    return
  }

  isUpdating.value = true
  editError.value = ''

  try {
    await emit('update', {
      name: editForm.name.trim(),
      longName: editForm.long_name.trim() || null,
      sku: editForm.sku.trim() || null,
      barcode: editForm.barcode.trim() || null,
      basePrice: editForm.base_price,
      taxRate: editForm.tax_rate,
      productType: editForm.product_type,
      description: editForm.description.trim() || null,
      brandId: editForm.brand_id || undefined,
      groupIds: editForm.group_ids.length > 0 ? editForm.group_ids : [],
    })
    isEditing.value = false
  } catch (err) {
    editError.value = err?.message || 'No se pudo actualizar el producto'
  } finally {
    isUpdating.value = false
  }
}

function formatDate(dateString) {
  if (!dateString) return '—'
  return new Date(dateString).toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatPrice(price) {
  if (price === undefined || price === null) return '0,00'
  return new Intl.NumberFormat('es-ES', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(price)
}
</script>

<style scoped>
.product-info {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.info-header h3 {
  color: #1b3a6b;
  margin: 0;
  font-size: 1.2rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-item label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
}

.info-item p {
  color: #1e293b;
  margin: 0;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.sku-code,
.barcode {
  background: #f1f5f9;
  color: #475569;
  padding: 0.25rem 0.6rem;
  border-radius: 4px;
  font-size: 0.85rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 600;
}

.brand-name {
  font-weight: 600;
  color: #1b3a6b;
}

.category-tag {
  display: inline-block;
  background: #f1f5f9;
  color: #475569;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
  margin-right: 0.5rem;
  margin-bottom: 0.25rem;
}

.text-muted {
  color: #94a3b8;
  font-style: italic;
}

.description {
  line-height: 1.6;
  color: #475569;
}

.pill {
  display: inline-block;
  padding: 0.25rem 0.75rem;
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

.btn-primary {
  background: #f4d03f;
  color: #1e293b;
  font-weight: 700;
}

.btn-primary:hover:not(:disabled) {
  background: #e6c530;
  box-shadow: 0 2px 8px rgba(244, 208, 63, 0.3);
}

.btn-primary:disabled {
  background: #e2e8f0;
  color: #94a3b8;
  cursor: not-allowed;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

.btn-link {
  background: none;
  border: none;
  color: #3b82f6;
  font-size: 0.85rem;
  cursor: pointer;
  text-decoration: underline;
  padding: 0;
  font-weight: 500;
}

.btn-link:hover {
  color: #1d4ed8;
}

.edit-form {
  background: #f8fafc;
  padding: 1.5rem;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin-bottom: 0.5rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.6rem 0.8rem;
  font-size: 0.9rem;
  color: #1e293b;
  font-family: inherit;
  background: #ffffff;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1b3a6b;
  box-shadow: 0 0 0 3px rgba(27, 58, 107, 0.1);
}

.form-group textarea {
  resize: vertical;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.error-message {
  color: #dc2626;
  font-size: 0.85rem;
  margin-top: 0.5rem;
  padding: 0.5rem;
  background: #fee2e2;
  border-radius: 6px;
  border: 1px solid #fca5a5;
}

.price-value {
  color: #059669;
  font-size: 1.1rem;
  font-weight: 600;
}

.hint {
  color: #64748b;
  font-size: 0.8rem;
  font-style: italic;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 200px;
  overflow-y: auto;
  padding: 0.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #f8fafc;
}

.checkbox-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  cursor: pointer;
  padding: 0.3rem 0.4rem;
  border-radius: 4px;
  transition: background 0.15s;
}

.checkbox-option:hover {
  background: #e2e8f0;
}

.checkbox-option input[type="checkbox"] {
  width: auto;
  margin: 0;
}

.pill-small {
  font-size: 0.7rem;
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  background: #e2e8f0;
  color: #64748b;
  text-transform: uppercase;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
