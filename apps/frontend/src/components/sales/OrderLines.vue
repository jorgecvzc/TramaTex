<template>
  <div class="order-lines-component">
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>Producto / Variante</th>
            <th class="text-center" style="width: 100px">Cant.</th>
            <th class="text-right" style="width: 140px">Precio Unit.</th>
            <th class="text-center" style="width: 100px">Dto %</th>
            <th class="text-right" style="width: 140px">Subtotal</th>
            <th v-if="isEditing" style="width: 50px"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(line, index) in lines" :key="index">
            <!-- Producto / Variante (SKU al lado del nombre) -->
            <td>
              <div class="product-info-cell">
                <div class="name-with-sku">
                  <span class="sku-label" v-if="line.variant_sku || line.variantSku">{{ line.variant_sku || line.variantSku }}</span>
                  <span class="product-name">{{ line.product_name || line.productName || line.description }}</span>
                </div>
              </div>
            </td>

            <!-- Cantidad -->
            <td class="text-center">
              <input v-if="isEditing" v-model.number="line.quantity" type="number" class="form-input text-center" min="1" @input="emitLines" />
              <span v-else>{{ line.quantity }}</span>
            </td>

            <!-- Precio Unitario -->
            <td class="text-right">
              <input v-if="isEditing" v-model.number="line.unit_price" type="number" step="0.01" class="form-input text-right" @input="emitLines" />
              <span v-else>{{ formatMoney(line.unit_price || line.unitPrice) }}</span>
            </td>

            <!-- Descuento -->
            <td class="text-center">
              <input v-if="isEditing" v-model.number="line.discount_percent" type="number" step="0.01" class="form-input text-center" min="0" max="100" @input="emitLines" />
              <span v-else>{{ line.discount_percent || line.discountPercent || 0 }}%</span>
            </td>

            <!-- Subtotal -->
            <td class="text-right">
              <strong class="text-primary">{{ formatMoney(calculateSubtotal(line)) }}</strong>
            </td>

            <!-- Borrar -->
            <td v-if="isEditing" class="text-center">
              <button class="btn-delete" type="button" @click="removeLine(index)" title="Quitar línea">
                <span class="material-symbols-outlined">delete</span>
              </button>
            </td>
          </tr>
          <tr v-if="lines.length === 0">
            <td :colspan="isEditing ? 6 : 5" class="empty-msg">
              <p>No hay líneas cargadas en el pedido.</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import salesApi from '@/services/salesApi'

const props = defineProps({
  lines: { type: Array, default: () => [] },
  isEditing: { type: Boolean, default: false }
})

const emit = defineEmits(['update:lines'])

function removeLine(index) {
  const newLines = [...props.lines]
  newLines.splice(index, 1)
  emit('update:lines', newLines)
}

function calculateSubtotal(line) {
  const price = Number(line.unit_price || line.unitPrice || 0)
  const qty = Number(line.quantity || 0)
  const disc = Number(line.discount_percent || line.discountPercent || 0)
  return (price * qty) * (1 - disc / 100)
}

function emitLines() {
  emit('update:lines', props.lines.map(line => ({ ...line })))
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
