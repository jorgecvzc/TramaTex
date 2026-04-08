<script setup lang="ts">
const props = defineProps<{
  title: string
  icon: string
  description?: string
}>()
</script>

<template>
  <section class="card form-section">
    <div class="section-header">
      <div class="header-icon">
        <span class="material-symbols-outlined">{{ props.icon }}</span>
      </div>
      <div class="header-text">
        <h2>{{ props.title }}</h2>
        <p v-if="props.description" class="header-desc">{{ props.description }}</p>
      </div>
    </div>
    <div class="section-body">
      <slot></slot>
    </div>
  </section>
</template>

<style scoped>
.form-section {
  padding: 0;
  overflow: visible;
  margin-bottom: 1.5rem;
}

.section-header {
  padding: 1.25rem 1.5rem;
  background: var(--color-background);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
  background: white;
  padding: 0.5rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  min-width: 42px;
  min-height: 42px;
}

.header-icon .material-symbols-outlined {
  font-size: 24px;
}

.header-text h2 {
  font-size: 1rem;
  font-weight: 700;
  margin: 0;
  color: var(--color-text-primary);
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.header-desc {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  margin: 0.15rem 0 0;
}

.section-body {
  padding: 1.5rem;
}
</style>
