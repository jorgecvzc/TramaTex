<script setup lang="ts">
import { getIcon } from '@/utils/icons'

const props = defineProps<{
  label: string
  value?: string | number | null
  icon?: any
  isMono?: boolean
  highlight?: boolean
}>()
</script>

<template>
  <div class="data-row" :class="{ 'highlight-row': props.highlight }">
    <div class="data-label">
      <div v-if="props.icon" class="row-icon">
        <component :is="getIcon(props.icon)" :size="18" />
      </div>
      <label>{{ props.label }}</label>
    </div>
    <div class="data-value" :class="{ 'text-mono': props.isMono }">
      <slot>
        {{ props.value || '—' }}
      </slot>
    </div>
  </div>
</template>

<style scoped>
.data-row {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--color-background);
  align-items: baseline;
  gap: 2rem;
}

.data-row:last-child {
  border-bottom: none;
}

.data-label {
  flex: 0 0 250px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.row-icon {
  display: flex;
  align-items: center;
  color: var(--color-text-secondary);
  opacity: 0.7;
}

.row-icon .material-symbols-outlined {
  font-size: 18px;
}

.data-label label {
  font-size: var(--font-size-xs);
  font-weight: 700;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  letter-spacing: 0.05em;
}

.data-value {
  flex: 1;
  font-size: 1rem;
  color: var(--color-text-primary);
  font-weight: 500;
}

.text-mono {
  font-family: var(--font-family-mono);
  font-size: 0.9rem;
}

.highlight-row .data-value {
  color: var(--color-primary);
  font-weight: 700;
  font-size: 1.1rem;
}

@media (max-width: 768px) {
  .data-row {
    flex-direction: column;
    gap: 0.25rem;
  }
  .data-label {
    flex: none;
  }
}
</style>
