<script setup lang="ts">
import { useRouter } from 'vue-router'
import PageHeader from '@/components/layout/PageHeader.vue'

interface Breadcrumb {
  label: string
  to?: string
}

const props = defineProps<{
  title: string
  icon?: string
  breadcrumbs: Breadcrumb[]
  catalogRoute?: string
  catalogText?: string
}>()

const router = useRouter()
</script>

<template>
  <div class="main-container">
    <PageHeader :title="props.title" :breadcrumbs="props.breadcrumbs">
      <template #icon v-if="props.icon">
        <span class="material-symbols-outlined">{{ props.icon }}</span>
      </template>
      
      <template #actions>
        <div class="header-actions-group">
          <slot name="actions"></slot>
          <button v-if="props.catalogRoute" @click="router.push(props.catalogRoute)" class="btn btn-outline">
            <span class="material-symbols-outlined">list_alt</span>
            <span>{{ props.catalogText || 'Ir al catálogo' }}</span>
          </button>
        </div>
      </template>
    </PageHeader>

    <div class="detail-vertical-layout">
      <div v-if="$slots.toolbar" class="detail-toolbar-area">
        <slot name="toolbar"></slot>
      </div>

      <div v-if="$slots.top" class="detail-overview-area">
        <slot name="top"></slot>
      </div>

      <div class="detail-main-content">
        <slot></slot>
      </div>

      <div v-if="$slots.footer" class="detail-footer-area">
        <slot name="footer"></slot>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail-vertical-layout {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
}

.detail-toolbar-area {
  width: 100%;
}

.detail-overview-area {
  width: 100%;
}

.detail-main-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
}

.detail-footer-area {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.header-actions-group {
  display: flex;
  gap: 0.75rem;
}
</style>
