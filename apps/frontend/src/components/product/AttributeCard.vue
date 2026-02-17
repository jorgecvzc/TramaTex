<template>
  <div class="attribute-card" :class="`source-${source}`">
    <div class="card-header">
      <div class="attribute-name">
        <span class="attribute-icon">{{ getSourceIcon(source) }}</span>
        <h5>{{ attribute.name }}</h5>
      </div>
      <code class="attribute-code">{{ attribute.code }}</code>
    </div>

    <div class="card-body">
      <div class="attribute-info">
        <span class="info-label">Orden:</span>
        <span class="info-value">{{ attribute.order || '—' }}</span>
      </div>

      <div class="attribute-info">
        <span class="info-label">Origen:</span>
        <span class="info-value source-label" :class="`source-${source}`">
          {{ getSourceLabel(source) }}
        </span>
      </div>
    </div>

    <div class="card-values">
      <div class="values-header">
        <span class="values-label">Valores disponibles</span>
        <span class="values-count">{{ attribute.values?.length || 0 }}</span>
      </div>

      <div v-if="attribute.values && attribute.values.length > 0" class="values-list">
        <span
          v-for="value in attribute.values"
          :key="value.id || value.code"
          class="value-tag"
        >
          {{ value.value }}
          <code class="value-code">{{ value.code }}</code>
        </span>
      </div>

      <p v-else class="no-values">Sin valores configurados</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  attribute: {
    type: Object,
    required: true,
  },
  source: {
    type: String,
    required: true,
    validator: (value) => ['direct', 'brand-group', 'group', 'brand', 'generic'].includes(value),
  },
})

function getSourceIcon(source) {
  const icons = {
    'direct': '📌',
    'brand-group': '🏢',
    'group': '📁',
    'brand': '🏷️',
    'generic': '🌐',
  }
  return icons[source] || '❓'
}

function getSourceLabel(source) {
  const labels = {
    'direct': 'Directo',
    'brand-group': 'Marca + Categoría',
    'group': 'Categoría',
    'brand': 'Marca',
    'generic': 'Genérico',
  }
  return labels[source] || source
}
</script>

<style scoped>
.attribute-card {
  background: #ffffff;
  border-radius: 8px;
  border: 2px solid #e2e8f0;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  transition: all 0.2s ease;
}

.attribute-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: #cbd5e1;
}

.attribute-card.source-direct {
  border-left: 4px solid #3b82f6;
}

.attribute-card.source-brand-group {
  border-left: 4px solid #8b5cf6;
}

.attribute-card.source-group {
  border-left: 4px solid #10b981;
}

.attribute-card.source-brand {
  border-left: 4px solid #f59e0b;
}

.attribute-card.source-generic {
  border-left: 4px solid #64748b;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 0.75rem;
}

.attribute-name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
}

.attribute-icon {
  font-size: 1.2rem;
}

.attribute-name h5 {
  margin: 0;
  color: #1b3a6b;
  font-size: 0.95rem;
  font-weight: 700;
}

.attribute-code {
  background: #f1f5f9;
  color: #475569;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.75rem 0;
  border-top: 1px solid #f1f5f9;
  border-bottom: 1px solid #f1f5f9;
}

.attribute-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.info-value {
  font-size: 0.85rem;
  color: #1e293b;
  font-weight: 500;
}

.source-label {
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
}

.source-label.source-direct {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.source-label.source-brand-group {
  background: rgba(139, 92, 246, 0.1);
  color: #8b5cf6;
}

.source-label.source-group {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.source-label.source-brand {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.source-label.source-generic {
  background: rgba(100, 116, 139, 0.1);
  color: #64748b;
}

.card-values {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.values-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.values-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.values-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.5rem;
  height: 1.5rem;
  padding: 0 0.4rem;
  background: #1b3a6b;
  color: #ffffff;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 700;
}

.values-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.value-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: #f1f5f9;
  padding: 0.35rem 0.65rem;
  border-radius: 6px;
  font-size: 0.8rem;
  color: #1e293b;
  font-weight: 500;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
}

.value-tag:hover {
  background: #e2e8f0;
  border-color: #cbd5e1;
}

.value-code {
  background: #cbd5e1;
  color: #475569;
  padding: 0.15rem 0.35rem;
  border-radius: 3px;
  font-size: 0.7rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.no-values {
  color: #94a3b8;
  font-style: italic;
  font-size: 0.85rem;
  margin: 0;
  text-align: center;
  padding: 0.5rem;
}
</style>
