<template>
  <div class="order-lines">
    <div class="table-wrapper">
      <table class="data-table">
        <thead>
          <tr>
            <th>Producto</th>
            <th class="text-center">Cantidad</th>
            <th class="align-right">Precio Unit.</th>
            <th class="align-right">Subtotal</th>
            <th v-if="isEditing" class="text-center">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(line, index) in normalizedLines" :key="line.id || index">
            <td>
              <template v-if="isEditing">
                <input
                  v-model="line.description"
                  class="form-input"
                  type="text"
                  placeholder="Descripción de la línea"
                  @input="emitLines"
                />
              </template>
              <template v-else>
                <strong>{{ line.description || 'Sin descripción' }}</strong>
              </template>
            </td>
            <td class="text-center">
              <template v-if="isEditing">
                <input
                  v-model.number="line.quantity"
                  class="form-input form-input-small"
                  type="number"
                  min="1"
                  @input="emitLines"
                />
              </template>
              <template v-else>{{ line.quantity }}</template>
            </td>
            <td class="align-right">
              <template v-if="isEditing">
                <input
                  v-model.number="line.unitPrice"
                  class="form-input form-input-small align-right"
                  type="number"
                  step="0.01"
                  min="0"
                  @input="emitLines"
                />
              </template>
              <template v-else>{{ formatMoney(line.unitPrice) }}</template>
            </td>
            <td class="align-right">
              <strong>{{ formatMoney(line.quantity * line.unitPrice) }}</strong>
            </td>
            <td v-if="isEditing" class="text-center">
              <button class="btn-icon text-danger" type="button" @click="removeLine(index)">
                <span class="material-symbols-outlined">delete</span>
              </button>
            </td>
          </tr>
          <tr v-if="normalizedLines.length === 0">
            <td :colspan="isEditing ? 5 : 4" class="empty-row">No hay líneas de pedido.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="isEditing" class="actions-row">
      <button class="btn btn-outline btn-sm" type="button" @click="addLine">
        <span class="material-symbols-outlined">add</span>
        <span>Añadir línea</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  lines: {
    type: Array,
    default: () => [],
  },
  isEditing: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:lines'])

const normalizedLines = ref([])

function normalizeLine(line) {
  return {
    ...line,
    description: line?.description || line?.productName || line?.variantSku || '',
    quantity: Number(line?.quantity || 1),
    unitPrice: Number(typeof line?.unitPrice === 'object' ? line?.unitPrice?.amount : line?.unitPrice || 0),
  }
}

watch(
  () => props.lines,
  (value) => {
    normalizedLines.value = (value || []).map(normalizeLine)
  },
  { immediate: true, deep: true }
)

function emitLines() {
  emit('update:lines', normalizedLines.value.map((line) => ({ ...line })))
}

function addLine() {
  normalizedLines.value.push({ description: '', quantity: 1, unitPrice: 0 })
  emitLines()
}

function removeLine(index) {
  normalizedLines.value.splice(index, 1)
  emitLines()
}

function formatMoney(value) {
  return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(Number(value || 0))
}
</script>

<style scoped>
.table-wrapper {
  overflow-x: auto;
}

.actions-row {
  margin-top: 1rem;
}

.form-input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 6px;
}

.form-input-small {
  max-width: 120px;
}

.align-right {
  text-align: right;
}

.text-center {
  text-align: center;
}

.btn-icon {
  background: transparent;
  border: none;
  cursor: pointer;
}

.empty-row {
  text-align: center;
  color: var(--color-text-secondary);
}
</style>
