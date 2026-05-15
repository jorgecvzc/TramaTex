<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ChevronRight, ArrowLeft } from 'lucide-vue-next'

interface Breadcrumb {
  label: string
  to?: string
}

interface Shortcut {
  label: string
  key: string
}

const props = defineProps<{
  title: string
  subtitle?: string
  breadcrumbs?: Breadcrumb[]
  showBack?: boolean
  shortcuts?: Shortcut[]
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
          <div v-if="shortcuts && shortcuts.length > 0" class="header-shortcuts">
            <div v-for="shortcut in shortcuts" :key="shortcut.key" class="shortcut-item">
              <span class="shortcut-label">{{ shortcut.label }}</span>
              <kbd>{{ shortcut.key }}</kbd>
            </div>
          </div>
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
  border: 1px solid var(--color-border-strong, #cbd5e1);
  border-radius: var(--border-radius-md, 8px);
  color: var(--color-text-secondary);
  padding: var(--spacing-xs, 0.5rem);
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.back-button:hover {
  background-color: var(--color-background-soft, #f8fafc);
  color: var(--color-primary);
  border-color: var(--color-primary);
  transform: translateX(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.back-button:active {
  transform: translateX(0);
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.06);
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
  gap: var(--spacing-lg);
}

.header-shortcuts {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding-right: var(--spacing-md);
  border-right: 1px solid var(--color-border);
}

.shortcut-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.shortcut-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  font-weight: var(--font-weight-medium);
}

/* Industrial Keyboard Shortcut Style */
kbd {
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  padding: 0.1rem 0.4rem;
  font-family: var(--font-family-mono, monospace);
  font-size: 0.75rem;
  box-shadow: 0 2px 0 #cbd5e1;
  color: var(--color-text-secondary);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.2rem;
}

/* Icons in slots */
.title-with-icon :deep(svg) {
  color: var(--color-primary);
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
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }

  .header-shortcuts {
    padding-right: 0;
    border-right: none;
    padding-bottom: var(--spacing-xs);
    border-bottom: 1px solid var(--color-border);
    width: 100%;
  }
}
</style>
