<script setup lang="ts">
import { useRouter } from 'vue-router'

interface Breadcrumb {
  label: string
  to?: string
}

const props = defineProps<{
  title: string
  subtitle?: string
  breadcrumbs?: Breadcrumb[]
  showBack?: boolean
}>()

const router = useRouter()

function goBack() {
  router.back()
}
</script>

<template>
  <header class="page-header">
    <div class="page-header-info">
      <nav v-if="breadcrumbs" class="breadcrumb" aria-label="Breadcrumb">
        <template v-for="(item, index) in breadcrumbs" :key="index">
          <RouterLink v-if="item.to" :to="item.to" class="breadcrumb-item">{{ item.label }}</RouterLink>
          <span v-else class="breadcrumb-item">{{ item.label }}</span>
        </template>
      </nav>

      <div class="title-with-back">
        <button v-if="showBack" @click="goBack" class="btn btn-ghost btn-sm back-button" title="Volver">
          <span class="material-symbols-outlined">arrow_back</span>
        </button>
        <div class="title-content">
          <slot name="icon"></slot>
          <h1>{{ title }}</h1>
        </div>
      </div>
      
      <p v-if="subtitle" class="subtitle">{{ subtitle }}</p>
    </div>
    
    <div class="page-header-actions">
      <slot name="actions"></slot>
    </div>
  </header>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  min-height: inherit; /* Hereda la altura centrada del padre */
}

.page-header-info {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px; /* Pequeño ajuste para cuando hay subtítulo */
}

.title-with-back {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.title-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  color: var(--color-secondary);
}

/* Soporte para iconos tanto en slot como directos */
.title-content :deep(.material-symbols-outlined) {
  color: var(--color-primary);
  font-size: 28px;
}

.back-button {
  margin-left: calc(var(--spacing-md) * -1);
  color: var(--color-text-secondary);
}

h1 {
  margin: 0;
}

.page-header-actions {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}
</style>
