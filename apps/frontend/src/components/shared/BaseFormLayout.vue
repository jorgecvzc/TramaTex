<script setup lang="ts">
import { useRouter } from 'vue-router'
import { List, Save, RefreshCw } from 'lucide-vue-next'
import { getIcon } from '@/utils/icons'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'

interface Breadcrumb {
  label: string
  to?: string
}

const props = defineProps<{
  title: string
  breadcrumbs: Breadcrumb[]
  isSubmitting?: boolean
  submitText?: string
  submitIcon?: string
  cancelRoute?: string
  catalogRoute?: string
}>()

const emit = defineEmits(['submit', 'cancel'])
const router = useRouter()

function handleCancel() {
  if (props.cancelRoute) router.push(props.cancelRoute)
  else router.back()
  emit('cancel')
}
</script>

<template>
  <div class="main-container">
    <BasePageHeader :title="props.title" :breadcrumbs="props.breadcrumbs">
      <template #actions>
        <button v-if="props.catalogRoute" @click="router.push(props.catalogRoute)" class="btn btn-outline">
          <List :size="18" />
          <span>Ir al catálogo</span>
        </button>
      </template>
    </BasePageHeader>

    <form @submit.prevent="emit('submit')" class="form-standard-layout">
      <div class="form-content">
        <slot></slot>
      </div>

      <footer class="form-footer-actions card">
        <slot name="actions">
          <button type="button" @click="handleCancel" class="btn btn-outline btn-lg" :disabled="props.isSubmitting">
            Cancelar
          </button>
          <button type="submit" class="btn btn-primary btn-lg" :disabled="props.isSubmitting">
            <component 
              :is="props.isSubmitting ? RefreshCw : (getIcon(props.submitIcon) || Save)" 
              :size="20" 
              :class="{ 'spin': props.isSubmitting }" 
            />
            <span>{{ props.isSubmitting ? 'Procesando...' : (props.submitText || 'Guardar Cambios') }}</span>
          </button>
        </slot>
      </footer>
    </form>
  </div>
</template>

<style scoped>
/* ... (rest of CSS) */

.lucide.spin {
  animation: spin 1s linear infinite;
}

.form-footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1.5rem;
  padding: 1.5rem 2rem;
  margin-top: 1rem;
  position: sticky;
  bottom: 1.5rem;
  z-index: 10;
  border-top: 4px solid var(--color-primary);
}

.form-footer-actions .btn {
  min-width: 180px;
}

.lucide.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .form-footer-actions {
    flex-direction: column-reverse;
    bottom: 0;
    margin: 0 -1rem;
    border-radius: 0;
  }
  .form-footer-actions .btn {
    width: 100%;
  }
}
</style>
