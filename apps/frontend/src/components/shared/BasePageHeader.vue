<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ChevronRight, ArrowLeft } from 'lucide-vue-next'

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
  <header class="base-page-header">
    <div class="header-main-content">
      <!-- Breadcrumbs Layer -->
      <nav v-if="breadcrumbs && breadcrumbs.length > 0" class="breadcrumbs" aria-label="Breadcrumb">
        <ol class="breadcrumb-list">
          <li v-for="(item, index) in breadcrumbs" :key="index" class="breadcrumb-item">
            <template v-if="index > 0">
              <ChevronRight class="breadcrumb-separator" :size="14" />
            </template>
            
            <RouterLink v-if="item.to" :to="item.to" class="breadcrumb-link">
              {{ item.label }}
            </RouterLink>
            <span v-else class="breadcrumb-current">{{ item.label }}</span>
          </li>
        </ol>
      </nav>

      <!-- Title & Back Button Layer -->
      <div class="title-row">
        <div class="title-group">
          <button 
            v-if="showBack" 
            @click="goBack" 
            class="back-button" 
            title="Volver"
          >
            <ArrowLeft :size="20" />
          </button>
          
          <div class="title-container">
            <div class="title-with-icon">
              <slot name="icon"></slot>
              <h1 class="page-title">{{ title }}</h1>
            </div>
            <p v-if="subtitle" class="page-subtitle">{{ subtitle }}</p>
          </div>
        </div>

        <!-- Actions Slot -->
        <div class="header-actions">
          <slot name="actions"></slot>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.base-page-header {
  background-color: transparent;
  padding: var(--spacing-sm) 0 var(--spacing-md) 0;
  width: 100%;
}

.header-main-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

/* Breadcrumbs Styles */
.breadcrumbs {
  margin-bottom: var(--spacing-xs);
}

.breadcrumb-list {
  display: flex;
  align-items: center;
  list-style: none;
  padding: 0;
  margin: 0;
  flex-wrap: wrap;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.breadcrumb-separator {
  margin: 0 var(--spacing-xs);
  color: var(--color-border-strong);
  font-size: 14px;
  flex-shrink: 0;
}

.breadcrumb-link {
  color: var(--color-info);
  text-decoration: none;
  transition: color 0.2s;
}

.breadcrumb-link:hover {
  color: var(--color-secondary);
  text-decoration: underline;
}

.breadcrumb-current {
  color: var(--color-text-secondary);
  font-weight: var(--font-weight-medium);
}

/* Title Row Styles */
.title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.title-group {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.back-button {
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-md);
  color: var(--color-text-secondary);
  padding: var(--spacing-xs);
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.back-button:hover {
  background-color: var(--color-background);
  color: var(--color-text-primary);
  border-color: var(--color-border-strong);
}

.title-container {
  display: flex;
  flex-direction: column;
}

.title-with-icon {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.page-title {
  margin: 0;
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  line-height: 1.2;
}

.page-subtitle {
  margin: var(--spacing-xs) 0 0 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

/* Icons in slots */
.title-with-icon :deep(.material-symbols-outlined) {
  color: var(--color-primary);
  font-size: 28px;
  flex-shrink: 0;
}

@media (max-width: 640px) {
  .title-row {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-start;
    margin-top: var(--spacing-sm);
  }
}
</style>
