<template>
  <div class="form-step">
    <h3 class="step-title">Gestión de variantes</h3>
    <p class="step-description">
      Información sobre cómo se gestionarán las variantes de este producto
    </p>

    <form @submit.prevent="handleNext" class="step-form">
      <!-- Info Box -->
      <div class="info-box">
        <span class="material-symbols-outlined info-icon" style="font-size: 20px">info</span>
        <div class="info-content">
          <p>
            <strong>¿Qué son las variantes?</strong> Una variante es una combinación específica de valores de atributos.
            Por ejemplo: Camiseta (Talla: L, Color: Rojo) es una variante del producto base "Camiseta".
          </p>
        </div>
      </div>

      <!-- Scenario: Product with NO attributes -->
      <div v-if="!hasAttributes" class="strategy-info-card simple-product">
        <div class="info-card-header">
          <span class="material-symbols-outlined info-card-icon" style="font-size: 24px">inventory_2</span>
          <div>
            <strong class="info-card-title">Producto Simple (Sin Variantes)</strong>
            <p class="info-card-description">
              Este producto no tiene atributos configurables, por lo tanto no tendrá variantes.
            </p>
          </div>
        </div>
        <div class="info-card-details">
          <p>
            <strong>¿Qué significa esto?</strong>
          </p>
          <ul>
            <li>El producto no requerirá selección de opciones al vender</li>
            <li>Se puede añadir directamente a pedidos usando su SKU base</li>
            <li>Ideal para productos únicos, servicios o artículos sin configuración</li>
          </ul>
        </div>
      </div>

      <!-- Scenario: Product WITH attributes -->
      <div v-else class="strategy-info-card jit-product">
        <div class="info-card-header">
          <span class="material-symbols-outlined info-card-icon" style="font-size: 24px">bolt</span>
          <div>
            <strong class="info-card-title">Producto Configurable (Variantes JIT + Manual)</strong>
            <p class="info-card-description">
              Este producto tiene atributos configurables. Las variantes se crearán bajo demanda o manualmente.
            </p>
          </div>
        </div>
        <div class="info-card-details">
          <p>
            <strong>¿Cómo funciona?</strong>
          </p>
          <ul>
            <li>
              <strong>Just-in-Time (JIT):</strong> Al añadir a un pedido, si la variante no existe, se creará automáticamente con estado PROVISIONAL
            </li>
            <li>
              <strong>Creación Manual:</strong> Desde la vista de detalle del producto, podrás crear variantes específicas manualmente
            </li>
            <li>
              <strong>Combinaciones posibles:</strong> Se pueden generar hasta {{ estimatedVariantCount }} variantes según los atributos seleccionados
            </li>
          </ul>
        </div>
        <div class="info-card-note">
          <span class="note-icon material-symbols-outlined">lightbulb</span>
          <span><strong>Nota:</strong> Solo se crearán las variantes que realmente se utilicen, evitando generar combinaciones innecesarias.</span>
        </div>
      </div>

      <!-- Architecture Note -->
      <div class="architecture-note">
        <p>
          <span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 4px; font-size: 18px">change_history</span>
          <strong>Según la arquitectura:</strong> La estrategia de generación de variantes no es configurable por el usuario.
          Se determina automáticamente según el tipo de producto y sus atributos.
        </p>
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
import { computed } from 'vue'

const props = defineProps({
  attributeCount: {
    type: Number,
    default: 0,
  },
  attributeValueCounts: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['next', 'prev'])

// Determine if product has attributes
const hasAttributes = computed(() => {
  return props.attributeCount > 0
})

// Calculate estimated variant count (cartesian product)
const estimatedVariantCount = computed(() => {
  if (!props.attributeValueCounts || props.attributeValueCounts.length === 0) {
    return 0
  }

  // Cartesian product: multiply all value counts
  return props.attributeValueCounts.reduce((acc, count) => acc * count, 1)
})

// Handle next step
function handleNext() {
  emit('next')
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

/* Strategy Info Cards */
.strategy-info-card {
  padding: 1.5rem;
  border-radius: 8px;
  border: 2px solid;
}

.strategy-info-card.simple-product {
  background-color: #f8fafc;
  border-color: #cbd5e1;
}

.strategy-info-card.jit-product {
  background-color: #fffbeb;
  border-color: #fde047;
}

.info-card-header {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.info-card-icon {
  font-size: 2rem;
  flex-shrink: 0;
}

.info-card-title {
  display: block;
  font-size: 1.1rem;
  color: #1e293b;
  margin-bottom: 0.25rem;
}

.info-card-description {
  font-size: 0.9rem;
  color: #64748b;
  margin: 0;
  line-height: 1.4;
}

.info-card-details {
  font-size: 0.9rem;
  color: #475569;
  line-height: 1.6;
}

.info-card-details p {
  margin: 0 0 0.5rem;
}

.info-card-details ul {
  margin: 0.5rem 0 0;
  padding-left: 1.5rem;
}

.info-card-details li {
  margin-bottom: 0.5rem;
}

.info-card-details li:last-child {
  margin-bottom: 0;
}

.info-card-note {
  display: flex;
  gap: 0.75rem;
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: rgba(255, 255, 255, 0.5);
  border-radius: 6px;
  font-size: 0.85rem;
  color: #475569;
  line-height: 1.5;
}

.note-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

/* Architecture Note */
.architecture-note {
  padding: 1rem;
  background-color: #f1f5f9;
  border-left: 4px solid #64748b;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #475569;
  line-height: 1.5;
}

.architecture-note p {
  margin: 0;
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
