<template>
  <div class="form-step">
    <h3 class="step-title">Preview y confirmación</h3>
    <p class="step-description">
      Revisa todos los datos antes de crear el producto
    </p>

    <div class="preview-container">
      <!-- Step 1: Basic Info -->
      <div class="preview-section">
        <div class="section-header">
          <h4 class="section-title">
            <FileText :size="18" class="section-icon" />
            Información básica
          </h4>
          <button
            type="button"
            @click="$emit('edit', 0)"
            class="btn-edit"
          >
            Editar
          </button>
        </div>
        <div class="section-content">
          <div class="preview-row">
            <span class="label">Tipo:</span>
            <span class="value">
              <span class="pill" :class="`type-${formData.productType?.toLowerCase()}`">
                {{ formData.productType === 'TANGIBLE' ? 'Tangible' : 'Servicio' }}
              </span>
            </span>
          </div>
          <div class="preview-row">
            <span class="label">SKU:</span>
            <code class="value code">{{ formData.sku }}</code>
          </div>
          <div class="preview-row">
            <span class="label">Nombre:</span>
            <span class="value">{{ formData.name }}</span>
          </div>
          <div v-if="formData.longName" class="preview-row">
            <span class="label">Nombre completo:</span>
            <span class="value">{{ formData.longName }}</span>
          </div>
          <div v-if="formData.description" class="preview-row">
            <span class="label">Descripción:</span>
            <span class="value description">{{ formData.description }}</span>
          </div>
        </div>
      </div>

      <!-- Step 2: Classification -->
      <div class="preview-section">
        <div class="section-header">
          <h4 class="section-title">
            <Tag :size="18" class="section-icon" />
            Clasificación
          </h4>
          <button
            type="button"
            @click="$emit('edit', 1)"
            class="btn-edit"
          >
            Editar
          </button>
        </div>
        <div class="section-content">
          <div class="preview-row">
            <span class="label">Marca:</span>
            <span class="value">{{ brandName || '(No especificada)' }}</span>
          </div>
          <div class="preview-row">
            <span class="label">Categorías:</span>
            <span v-if="groupNames.length > 0" class="value">
              <span v-for="(name, index) in groupNames" :key="index" class="tag">
                {{ name }}
              </span>
            </span>
            <span v-else class="value empty">(Ninguna)</span>
          </div>
        </div>
      </div>

      <!-- Step 3: Attributes -->
      <div class="preview-section">
        <div class="section-header">
          <h4 class="section-title">
            <Settings :size="18" class="section-icon" />
            Atributos configurables
          </h4>
          <button
            type="button"
            @click="$emit('edit', 2)"
            class="btn-edit"
          >
            Editar
          </button>
        </div>
        <div class="section-content">
          <div class="preview-row">
            <span class="label">Atributos directos:</span>
            <span v-if="attributeNames.length > 0" class="value">
              <span v-for="(name, index) in attributeNames" :key="index" class="tag">
                {{ name }}
              </span>
            </span>
            <span v-else class="value empty">(Ninguno seleccionado)</span>
          </div>
          <div v-if="attributeNames.length === 0" class="info-note">
            <Lightbulb :size="18" class="note-icon" />
            <span>
              Este producto heredará atributos de su marca y categorías.
            </span>
          </div>
        </div>
      </div>

      <!-- Step 4: Variant Management -->
      <div class="preview-section">
        <div class="section-header">
          <h4 class="section-title">
            <Hash :size="18" class="section-icon" />
            Gestión de variantes
          </h4>
          <span class="auto-badge">Automático</span>
        </div>
        <div class="section-content">
          <div class="preview-row">
            <span class="label">Tipo de producto:</span>
            <span class="value">
              <span class="badge" :class="`strategy-${formData.strategy}`">
                {{ getStrategyLabel(formData.strategy) }}
              </span>
            </span>
          </div>
          <div v-if="formData.strategy === 'jit'" class="info-note">
            <Zap :size="18" class="note-icon" />
            <div>
              <strong>Producto configurable con variantes JIT + Manual:</strong>
              <ul>
                <li>Las variantes se crearán automáticamente al añadirlas a pedidos (estado PROVISIONAL)</li>
                <li>También podrás crear variantes manualmente desde la vista de detalle del producto</li>
                <li>Solo se generarán las combinaciones que realmente se utilicen</li>
              </ul>
            </div>
          </div>
          <div v-else-if="formData.strategy === 'none'" class="info-note">
            <Package :size="18" class="note-icon" />
            <div>
              <strong>Producto simple sin variantes:</strong>
              <ul>
                <li>Este producto no tiene atributos configurables</li>
                <li>Se puede añadir directamente a pedidos usando su SKU base</li>
                <li>Ideal para productos únicos o servicios sin configuración</li>
              </ul>
            </div>
          </div>
          <div class="architecture-info">
            <Triangle :size="18" class="info-icon" />
            <span>
              La estrategia se determina automáticamente según la arquitectura del sistema y los atributos del producto.
            </span>
          </div>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button
          type="button"
          @click="$emit('prev')"
          class="btn btn-secondary"
          :disabled="isSubmitting"
        >
          ← Anterior
        </button>
        <button
          type="button"
          @click="handleSubmit"
          class="btn btn-primary"
          :disabled="isSubmitting"
        >
          <span v-if="isSubmitting">Creando producto...</span>
          <span v-else>
            <Check :size="16" />
            Crear producto
          </span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { 
  FileText, Tag, Settings, Lightbulb, Hash, Zap, Package, Triangle, Check 
} from 'lucide-vue-next'

const props = defineProps({
  formData: {
    type: Object,
    required: true,
  },
  brandName: {
    type: String,
    default: '',
  },
  groupNames: {
    type: Array,
    default: () => [],
  },
  attributeNames: {
    type: Array,
    default: () => [],
  },
  isSubmitting: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['edit', 'prev', 'submit'])

// Get strategy label
function getStrategyLabel(strategy) {
  const labels = {
    jit: 'Configurable (JIT + Manual)',
    none: 'Simple (Sin variantes)',
  }
  return labels[strategy] || strategy
}

// Handle submit
function handleSubmit() {
  emit('submit')
}
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

.preview-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.preview-section {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
  color: #1b3a6b;
  margin: 0;
}

.section-icon {
  font-size: 1.25rem;
}

.auto-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.3rem 0.65rem;
  background-color: #e0e7ff;
  border: 1px solid #c7d2fe;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  color: #4338ca;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.btn-edit {
  padding: 0.4rem 0.75rem;
  font-size: 0.85rem;
  font-weight: 500;
  color: #1b3a6b;
  background-color: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-edit:hover {
  background-color: #f1f5f9;
  border-color: #94a3b8;
}

.section-content {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.preview-row {
  display: grid;
  grid-template-columns: 150px 1fr;
  gap: 1rem;
  align-items: start;
}

.label {
  font-weight: 600;
  color: #475569;
  font-size: 0.9rem;
}

.value {
  color: #1e293b;
  font-size: 0.9rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  align-items: center;
}

.value.code {
  font-family: 'Courier New', monospace;
  background-color: #f1f5f9;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  border: 1px solid #cbd5e1;
}

.value.description {
  line-height: 1.5;
  color: #64748b;
}

.value.empty {
  color: #94a3b8;
  font-style: italic;
}

.pill {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.75rem;
  font-size: 0.8rem;
  font-weight: 600;
  border-radius: 12px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.pill.type-tangible {
  background-color: #dbeafe;
  color: #1e40af;
}

.pill.type-service {
  background-color: #e0e7ff;
  color: #4338ca;
}

.tag {
  display: inline-flex;
  align-items: center;
  padding: 0.35rem 0.65rem;
  background-color: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.85rem;
  color: #1e293b;
}

.badge {
  display: inline-flex;
  align-items: center;
  padding: 0.35rem 0.75rem;
  font-size: 0.85rem;
  font-weight: 600;
  border-radius: 6px;
}

.badge.strategy-automatic {
  background-color: #dcfce7;
  color: #166534;
}

.badge.strategy-jit {
  background-color: #fef3c7;
  color: #78350f;
}

.badge.strategy-manual {
  background-color: #e0e7ff;
  color: #4338ca;
}

.badge.strategy-none {
  background-color: #f1f5f9;
  color: #475569;
}

.info-note {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.75rem;
  background-color: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 6px;
  font-size: 0.85rem;
  color: #1e40af;
  line-height: 1.5;
}

.info-note div {
  flex: 1;
}

.info-note strong {
  display: block;
  margin-bottom: 0.5rem;
  color: #1e40af;
}

.info-note ul {
  margin: 0;
  padding-left: 1.25rem;
  list-style-type: disc;
}

.info-note li {
  margin-bottom: 0.35rem;
}

.info-note li:last-child {
  margin-bottom: 0;
}

.info-note.success {
  background-color: #f0fdf4;
  border-color: #bbf7d0;
  color: #166534;
}

.architecture-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background-color: #f1f5f9;
  border-left: 4px solid #64748b;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #475569;
  line-height: 1.5;
  margin-top: 0.5rem;
}

.info-icon {
  font-size: 1rem;
  flex-shrink: 0;
}

.note-icon {
  font-size: 1.2rem;
  flex-shrink: 0;
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
  min-width: 180px;
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

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

@media (max-width: 768px) {
  .form-step {
    max-width: 100%;
  }

  .preview-row {
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .btn {
    width: 100%;
  }
}
</style>
