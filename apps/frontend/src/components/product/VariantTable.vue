<template>
  <div class="variant-table-container">
    <div class="table-wrapper">
      <table class="data-table">
        <thead>
          <tr>
            <th>SKU / Referencia</th>
            <th>Configuración</th>
            <th class="text-right">Precio Base</th>
            <th class="text-center">Estado</th>
            <th class="text-right">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="variant in variants" :key="variant.id" class="row-hover">
            <td>
              <div class="sku-cell">
                <code class="code-badge">{{ variant.sku }}</code>
                <span v-if="variant.is_default" class="default-pill">Principal</span>
              </div>
            </td>
            <td>
              <div class="options-stack">
                <span v-for="(val, attr) in variant.option_configuration" :key="attr" class="option-pill">
                  <strong>{{ attr }}:</strong> {{ val }}
                </span>
              </div>
            </td>
            <td class="text-right">
              <strong class="price-text">{{ formatMoney(variant.price || basePrice) }}</strong>
            </td>
            <td class="text-center">
              <span :class="['status-badge', variant.is_active ? 'status-success' : 'status-secondary']">
                {{ variant.is_active ? 'Activa' : 'Inactiva' }}
              </span>
            </td>
            <td class="text-right" @click.stop>
              <div class="action-buttons">
                <button class="btn-icon" @click="emit('edit', variant)" title="Editar variante"><Pencil :size="18" /></button>
                <button 
                  class="btn-icon" 
                  :class="variant.is_active ? 'text-warning' : 'text-success'"
                  @click="promptToggleActive(variant)" 
                  :title="variant.is_active ? 'Desactivar' : 'Activar'"
                >
                  <component :is="variant.is_active ? Ban : CheckCircle" :size="18" />
                </button>
                <button class="btn-icon text-danger" @click="promptDelete(variant)" title="Eliminar"><Trash2 :size="18" /></button>
              </div>
            </td>
          </tr>
          <tr v-if="variants.length === 0">
            <td colspan="5" class="empty-msg">No hay variantes generadas para este producto.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- MODAL DE CONFIRMACIÓN -->
    <BaseDialog
      :show="confirmDialog.show"
      :title="confirmDialog.title"
      :icon="confirmDialog.icon"
      :confirm-text="confirmDialog.confirmText"
      :confirm-class="confirmDialog.confirmClass"
      :is-confirming="isProcessing"
      @close="confirmDialog.show = false"
      @confirm="handleConfirmAction"
    >
      <p>{{ confirmDialog.message }}</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { Pencil, Ban, CheckCircle, Trash2, AlertTriangle } from 'lucide-vue-next'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import salesApi from '@/services/salesApi'
import { productApi } from '@/services/productApi'
import { useToastStore } from '@/stores/toast'

const props = defineProps({
  variants: { type: Array, default: () => [] },
  basePrice: { type: Object, default: null },
  productId: { type: String, required: true }
})

const emit = defineEmits(['edit', 'refresh'])
const toastStore = useToastStore()
const isProcessing = ref(false)

// --- Confirm Dialog Logic ---
const confirmDialog = reactive({
  show: false,
  title: '',
  message: '',
  icon: AlertTriangle,
  confirmText: 'Confirmar',
  confirmClass: 'btn-primary',
  variant: null,
  actionType: null
})

function promptToggleActive(variant) {
  const action = variant.is_active ? 'DESACTIVAR' : 'ACTIVAR'
  confirmDialog.title = `${action.charAt(0) + action.slice(1).toLowerCase()} Variante`
  confirmDialog.message = `¿Estás seguro de que deseas ${action.toLowerCase()} la variante ${variant.sku}?`
  confirmDialog.icon = variant.is_active ? Ban : CheckCircle
  confirmDialog.confirmText = action.charAt(0) + action.slice(1).toLowerCase()
  confirmDialog.confirmClass = variant.is_active ? 'btn-warning' : 'btn-success'
  confirmDialog.variant = variant
  confirmDialog.actionType = 'toggle'
  confirmDialog.show = true
}

function promptDelete(variant) {
  confirmDialog.title = 'Eliminar Variante'
  confirmDialog.message = `¿Realmente deseas eliminar permanentemente la variante ${variant.sku}? Esta acción no se puede deshacer.`
  confirmDialog.icon = AlertTriangle
  confirmDialog.confirmText = 'Eliminar permanentemente'
  confirmDialog.confirmClass = 'btn-danger'
  confirmDialog.variant = variant
  confirmDialog.actionType = 'delete'
  confirmDialog.show = true
}

async function handleConfirmAction() {
  const { variant, actionType } = confirmDialog
  if (!variant) return

  isProcessing.value = true
  try {
    if (actionType === 'toggle') {
      await productApi.updateVariant(props.productId, variant.id, { is_active: !variant.is_active })
      toastStore.success(`Variante ${variant.is_active ? 'desactivada' : 'activada'}`)
    } else if (actionType === 'delete') {
      await productApi.deleteVariant(props.productId, variant.id)
      toastStore.success('Variante eliminada correctamente')
    }
    emit('refresh')
    confirmDialog.show = false
  } catch (err) {
    toastStore.error(err.message || 'Error al procesar la acción')
  } finally {
    isProcessing.value = false
  }
}

const formatMoney = (v) => salesApi.formatMoney(v)
</script>

<style scoped>
.variant-table-container { width: 100%; }
.sku-cell { display: flex; align-items: center; gap: 0.75rem; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.default-pill { background: #eff6ff; color: #1e40af; font-size: 0.6rem; font-weight: 800; text-transform: uppercase; padding: 0.15rem 0.4rem; border-radius: 4px; border: 1px solid #bfdbfe; }

.options-stack { display: flex; flex-wrap: wrap; gap: 0.4rem; }
.option-pill { font-size: 0.75rem; background: var(--color-background-soft); padding: 0.15rem 0.5rem; border-radius: 4px; border: 1px solid var(--color-border); color: var(--color-text-secondary); }
.option-pill strong { color: var(--color-text-primary); }

.price-text { color: var(--color-primary); }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; transition: 0.2s; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.empty-msg { text-align: center; padding: 3rem !important; color: var(--color-text-secondary); font-style: italic; }
</style>
