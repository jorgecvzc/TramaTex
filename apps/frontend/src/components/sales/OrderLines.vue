<template>
  <div class="order-lines-component">
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>Producto / Variante</th>
            <th class="text-center" style="width: 100px">{{ quantityLabel }}</th>
            <th v-if="showPrices" class="text-right" style="width: 140px">P. Tarifa</th>
            <th v-if="showPrices" class="text-right" style="width: 140px">P. Venta</th>
            <th v-if="showPrices" class="text-center" style="width: 100px">Dto %</th>
            <th v-if="showTotal" class="text-right" style="width: 140px">Subtotal</th>
            <th v-if="isEditing" style="width: 50px"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(line, index) in lines" :key="line.productVariantId || line.id || index">
            <!-- Producto / Variante (SKU al lado del nombre) -->
            <td>
              <div class="product-info-cell">
                <div class="name-with-sku">
                  <span class="sku-label" v-if="line.variantSku">{{ line.variantSku }}</span>
                  <span class="product-name">{{ line.productName || line.displayName || line.description || 'Producto' }}</span>
                </div>
              </div>
            </td>

            <!-- Cantidad -->
            <td class="text-center">
              <input 
                v-if="isEditing" 
                :value="line.quantity" 
                type="number" 
                class="form-input text-center" 
                min="1" 
                :data-row="index"
                data-col="quantity"
                @input="updateLineField(index, 'quantity', $event.target.value)" 
                @keydown="handleLineKeyDown($event, index, 'quantity', line)"
              />
              <span v-else>{{ line.quantity }}</span>
            </td>

            <!-- Precio Tarifa -->
            <td v-if="showPrices" class="text-right">
              <span class="text-muted">{{ formatMoney(line.listUnitPrice || line.listPrice) }}</span>
            </td>

            <!-- Precio Venta -->
            <td v-if="showPrices" class="text-right">
              <input 
                v-if="isEditing" 
                :value="line.unitPrice" 
                type="number" 
                step="0.01" 
                class="form-input text-right" 
                :data-row="index"
                data-col="unitPrice"
                @input="onManualPriceChange(index, $event.target.value)" 
                @keydown="handleLineKeyDown($event, index, 'unitPrice', line)"
              />
              <span v-else>{{ formatMoney(line.unitPrice) }}</span>
            </td>

            <!-- Descuento -->
            <td v-if="showPrices" class="text-center">
              <input 
                v-if="isEditing" 
                :value="line.discountPercent" 
                type="number" 
                step="0.01" 
                class="form-input text-center" 
                min="0" 
                max="100" 
                :data-row="index"
                data-col="discountPercent"
                @input="updateLineField(index, 'discountPercent', $event.target.value)" 
                @keydown="handleLineKeyDown($event, index, 'discountPercent', line)"
              />
              <span v-else>{{ line.discountPercent || 0 }}%</span>
            </td>

            <!-- Subtotal -->
            <td v-if="showTotal" class="text-right">
              <strong class="text-primary">{{ formatMoney(calculateSubtotal(line)) }}</strong>
            </td>

            <!-- Borrar -->
            <td v-if="isEditing" class="text-center">
              <button class="btn-delete" type="button" @click="removeLine(index)" title="Quitar línea">
                <Trash2 :size="18" />
              </button>
            </td>
          </tr>
          <tr v-if="!lines || lines.length === 0">
            <td :colspan="isEditing ? 7 : 6" class="empty-msg">
              <p>No hay líneas cargadas en el documento.</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { computed, toRefs } from 'vue'
import { Trash2 } from 'lucide-vue-next'
import { useLineNavigation } from '@/composables/useLineNavigation'
import salesApi from '@/services/salesApi'

const props = defineProps({
  lines: { type: Array, default: () => [] },
  isEditing: { type: Boolean, default: false },
  showPrices: { type: Boolean, default: true },
  showTotal: { type: Boolean, default: true },
  quantityLabel: { type: String, default: 'Cant.' }
})

const { lines, isEditing, showPrices, showTotal, quantityLabel } = toRefs(props)

const emit = defineEmits(['update:lines', 'add-line', 'last-field-tab'])

// DEBUG LOGGING
import { watch } from 'vue'
watch(lines, (newVal) => {
  console.log('[OrderLines] lines updated:', newVal?.length, newVal)
}, { immediate: true, deep: true })

const columns = computed(() => {
  const cols = ['quantity']
  if (showPrices.value) {
    cols.push('unitPrice')
    cols.push('discountPercent')
  }
  return cols
})

const { handleLineKeyDown } = useLineNavigation({
  rowCount: () => Array.isArray(lines.value) ? lines.value.length : 0,
  columns: columns.value,
  onUpdate: (index, col, val) => updateLineField(index, col, val),
  onRemoveField: (index) => removeLine(index),
  onLastFieldTab: () => emit('last-field-tab'),
  onLastFieldEnter: () => emit('add-line'),
  onAddField: () => emit('add-line')
})

function removeLine(index) {
  if (!Array.isArray(lines.value)) return
  const newLines = [...lines.value]
  newLines.splice(index, 1)
  emit('update:lines', newLines)
}

function updateLineField(index, field, value) {
  if (!Array.isArray(lines.value)) return
  const newLines = lines.value.map((line, i) => {
    if (i === index) {
      return { ...line, [field]: value === '' ? 0 : Number(value) }
    }
    return { ...line }
  })
  emit('update:lines', newLines)
}

function onManualPriceChange(index, value) {
  if (!Array.isArray(lines.value)) return
  const newLines = lines.value.map((line, i) => {
    if (i === index) {
      return { 
        ...line, 
        unitPrice: value === '' ? 0 : Number(value),
        _autoPrice: false // Mark as manual override
      }
    }
    return { ...line }
  })
  emit('update:lines', newLines)
}

function calculateSubtotal(line) {
  if (!isEditing.value && line.subtotal !== undefined) {
    return line.subtotal?.amount ?? line.subtotal;
  }
  
  const price = Number(line.unitPrice ?? 0)
  const qty = Number(line.quantity ?? 0)
  const disc = Number(line.discountPercent ?? 0)
  return (price * qty) * (1 - disc / 100)
}

function formatMoney(v) { return salesApi.formatMoney(v) }
</script>

<style scoped>
.order-lines-component { width: 100%; }
.table-container { border: 1px solid #cbd5e1; border-radius: 8px; background: white; overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { background: #f8fafc; padding: 0.75rem 1rem; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: #64748b; text-align: left; border-bottom: 2px solid #cbd5e1; }
.data-table td { padding: 0.75rem 1rem; border-bottom: 1px solid #f1f5f9; vertical-align: middle; }

/* Product Info */
.product-info-cell { display: flex; align-items: center; }
.name-with-sku { display: flex; align-items: center; gap: 0.75rem; }
.sku-label { font-family: 'Fira Code', monospace; font-size: 0.7rem; font-weight: 700; color: #1b3a6b; background: #eff6ff; padding: 0.15rem 0.4rem; border-radius: 4px; white-space: nowrap; }
.product-name { font-weight: 600; color: #1e293b; }

/* Inline Inputs */
.form-input { 
  width: 100%; padding: 0.5rem; border: 1px solid #cbd5e1; border-radius: 6px; 
  font-size: 0.9rem; font-family: inherit; transition: border-color 0.2s;
}
.form-input:focus { outline: none; border-color: #1b3a6b; box-shadow: 0 0 0 2px rgba(27, 58, 107, 0.1); }

/* Buttons */
.btn-delete { background: none; border: none; cursor: pointer; color: #94a3b8; padding: 0.4rem; border-radius: 6px; transition: 0.2s; display: flex; align-items: center; justify-content: center; }
.btn-delete:hover { color: #ef4444; background: #fef2f2; }

/* Utils */
.text-right { text-align: right; }
.text-center { text-align: center; }
.text-primary { color: #1b3a6b; }
.empty-msg { text-align: center; padding: 3rem !important; color: #94a3b8; font-style: italic; }
</style>
