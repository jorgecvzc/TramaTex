<template>
  <div class="attributes-panel">
    <div class="panel-header">
      <div>
        <h3>Atributos Configurables</h3>
        <p class="panel-description">
          Atributos heredados de marca y categorías, más atributos específicos asignados directamente.
        </p>
      </div>
      <button @click="$emit('refresh')" class="btn btn-secondary" :disabled="isLoading">
        <span class="material-symbols-outlined" style="font-size: 16px">refresh</span>
        Actualizar
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Cargando atributos...</p>
    </div>

    <!-- Empty State -->
    <div v-if="!isLoading && calculatedAttributes.length === 0" class="empty-state">
      <span class="material-symbols-outlined empty-icon" style="font-size: 64px">sell</span>
      <p>No hay atributos configurados para este producto.</p>
      <p class="empty-hint">
        Los atributos definen las características configurables del producto (talla, color, etc.)
        y se heredan de la marca y las categorías, o se asignan directamente.
      </p>
    </div>

    <!-- Attributes Display -->
    <div v-if="!isLoading && calculatedAttributes.length > 0" class="attributes-container">
      <!-- Direct Attributes Section -->
      <div v-if="directAttributes.length > 0" class="attribute-section">
        <div class="section-header">
          <h4>
            <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 6px; font-size: 18px">push_pin</span>
            Atributos Directos
          </h4>
          <span class="section-badge">{{ directAttributes.length }}</span>
        </div>
        <p class="section-description">
          Atributos asignados específicamente a este producto.
          Tienen la mayor prioridad.
        </p>
        <div class="attributes-grid">
          <AttributeCard
            v-for="attr in directAttributes"
            :key="attr.id"
            :attribute="attr"
            source="direct"
          />
        </div>
      </div>

      <!-- Brand + Group Attributes -->
      <div v-if="brandGroupAttributes.length > 0" class="attribute-section">
        <div class="section-header">
          <h4>
            <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 6px; font-size: 18px">corporate_fare</span>
            Marca + Categoría
          </h4>
          <span class="section-badge">{{ brandGroupAttributes.length }}</span>
        </div>
        <p class="section-description">
          Atributos heredados de la combinación de marca y categoría.
        </p>
        <div class="attributes-grid">
          <AttributeCard
            v-for="attr in brandGroupAttributes"
            :key="attr.id"
            :attribute="attr"
            source="brand-group"
          />
        </div>
      </div>

      <!-- Group Attributes -->
      <div v-if="groupAttributes.length > 0" class="attribute-section">
        <div class="section-header">
          <h4>
            <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 6px; font-size: 18px">folder</span>
            Categoría
          </h4>
          <span class="section-badge">{{ groupAttributes.length }}</span>
        </div>
        <p class="section-description">
          Atributos heredados de las categorías del producto.
        </p>
        <div class="attributes-grid">
          <AttributeCard
            v-for="attr in groupAttributes"
            :key="attr.id"
            :attribute="attr"
            source="group"
          />
        </div>
      </div>

      <!-- Brand Attributes -->
      <div v-if="brandAttributes.length > 0" class="attribute-section">
        <div class="section-header">
          <h4>
            <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 6px; font-size: 18px">sell</span>
            Marca
          </h4>
          <span class="section-badge">{{ brandAttributes.length }}</span>
        </div>
        <p class="section-description">
          Atributos heredados de la marca del producto.
        </p>
        <div class="attributes-grid">
          <AttributeCard
            v-for="attr in brandAttributes"
            :key="attr.id"
            :attribute="attr"
            source="brand"
          />
        </div>
      </div>

      <!-- Generic Attributes -->
      <div v-if="genericAttributes.length > 0" class="attribute-section">
        <div class="section-header">
          <h4>
            <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 6px; font-size: 18px">public</span>
            Genéricos
          </h4>
          <span class="section-badge">{{ genericAttributes.length }}</span>
        </div>
        <p class="section-description">
          Atributos genéricos aplicables a todos los productos.
        </p>
        <div class="attributes-grid">
          <AttributeCard
            v-for="attr in genericAttributes"
            :key="attr.id"
            :attribute="attr"
            source="generic"
          />
        </div>
      </div>

      <!-- Info Box -->
      <div class="info-box">
        <div class="info-icon">
          <span class="material-symbols-outlined" style="font-size: 20px">info</span>
        </div>
        <div>
          <strong>Jerarquía de Atributos</strong>
          <p>
            Los atributos se aplican según su nivel de especificidad:
            <strong>Directo</strong> &gt; <strong>Marca+Categoría</strong> &gt;
            <strong>Categoría</strong> &gt; <strong>Marca</strong> &gt; <strong>Genérico</strong>.
            Si un atributo con el mismo código existe en múltiples niveles, se usa el más específico.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import AttributeCard from './AttributeCard.vue'

const props = defineProps({
  product: {
    type: Object,
    required: true,
  },
  calculatedAttributes: {
    type: Array,
    default: () => [],
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['refresh'])

// Computed: Group attributes by source
const directAttributes = computed(() => {
  return props.calculatedAttributes.filter(
    attr => attr.source === 'direct' || attr.scope_type === 'direct'
  )
})

const brandGroupAttributes = computed(() => {
  return props.calculatedAttributes.filter(
    attr => attr.source === 'brand-group' ||
           (attr.scope_brand_id && attr.scope_group_id)
  )
})

const groupAttributes = computed(() => {
  return props.calculatedAttributes.filter(
    attr => attr.source === 'group' ||
           (attr.scope_group_id && !attr.scope_brand_id)
  )
})

const brandAttributes = computed(() => {
  return props.calculatedAttributes.filter(
    attr => attr.source === 'brand' ||
           (attr.scope_brand_id && !attr.scope_group_id)
  )
})

const genericAttributes = computed(() => {
  return props.calculatedAttributes.filter(
    attr => attr.source === 'generic' ||
           (!attr.scope_brand_id && !attr.scope_group_id && attr.source !== 'direct')
  )
})
</script>

<style scoped>
.attributes-panel {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.panel-header h3 {
  color: #1b3a6b;
  margin: 0 0 0.5rem 0;
  font-size: 1.2rem;
}

.panel-description {
  color: #64748b;
  margin: 0;
  font-size: 0.9rem;
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
  color: #64748b;
  gap: 1rem;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.empty-state p {
  margin: 0.25rem 0;
  font-size: 0.95rem;
}

.empty-hint {
  font-size: 0.85rem;
  color: #94a3b8;
  max-width: 600px;
  line-height: 1.5;
}

.attributes-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.attribute-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.5rem;
  background: #f8fafc;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.section-header h4 {
  color: #1b3a6b;
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
}

.section-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.75rem;
  height: 1.75rem;
  padding: 0 0.5rem;
  background: #1b3a6b;
  color: #ffffff;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
}

.section-description {
  color: #64748b;
  margin: 0;
  font-size: 0.85rem;
  line-height: 1.5;
}

.attributes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  gap: 0.5rem;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.info-box {
  display: flex;
  gap: 1rem;
  padding: 1.5rem;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 10px;
  color: #1e40af;
}

.info-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.info-box strong {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.95rem;
}

.info-box p {
  margin: 0;
  font-size: 0.85rem;
  line-height: 1.6;
  color: #1e40af;
}

@media (max-width: 768px) {
  .panel-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .attributes-grid {
    grid-template-columns: 1fr;
  }
}
</style>
