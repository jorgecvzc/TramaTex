<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  title: string
  stationId?: string
  icon?: string
  isLoading?: boolean
}>()

const emit = defineEmits(['refresh', 'close'])

// Al montar, forzamos al body a ignorar el scroll y avisamos del modo terminal
onMounted(() => {
  document.body.classList.add('terminal-open')
})

onUnmounted(() => {
  document.body.classList.remove('terminal-open')
})
</script>

<template>
  <div class="base-terminal-overlay">
    <!-- CABECERA INDUSTRIAL -->
    <header class="terminal-header">
      <div class="header-brand">
        <span class="material-symbols-outlined logo-icon">{{ icon || 'terminal' }}</span>
        <div class="title-stack">
          <h1>{{ title }}</h1>
          <span v-if="stationId" class="station-label">ESTACIÓN: {{ stationId }}</span>
        </div>
      </div>

      <div class="header-actions">
        <button class="btn-terminal btn-sync" @click="$emit('refresh')" :disabled="isLoading">
          <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
          <span>Sincronizar</span>
        </button>
        <button class="btn-terminal btn-exit" @click="$emit('close')">
          <span class="material-symbols-outlined">logout</span>
          <span>Salir</span>
        </button>
      </div>
    </header>

    <!-- CONTENIDO TÁCTIL -->
    <main class="terminal-body">
      <slot></slot>
    </main>

    <!-- BARRA DE ESTADO INFERIOR (Opcional) -->
    <footer v-if="$slots.footer" class="terminal-footer">
      <slot name="footer"></slot>
    </footer>
  </div>
</template>

<style scoped>
.base-terminal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: #0f172a;
  color: #f8fafc;
  display: flex;
  flex-direction: column;
  z-index: 5000; /* Por encima de la navegación estándar */
}

.terminal-header {
  background-color: #1e293b;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 2px solid #334155;
  box-shadow: 0 4px 15px rgba(0,0,0,0.3);
}

.header-brand { display: flex; align-items: center; gap: 1.5rem; }
.logo-icon { font-size: 3rem; color: var(--color-primary); }
.title-stack h1 { font-size: 1.5rem; margin: 0; font-weight: 800; color: white; text-transform: uppercase; }
.station-label { font-size: 0.75rem; color: var(--color-primary); font-weight: 700; letter-spacing: 0.1em; }

.header-actions { display: flex; gap: 1rem; }
.btn-terminal { 
  display: flex; align-items: center; gap: 0.75rem;
  padding: 0.75rem 1.5rem; border-radius: 12px; font-weight: 700; cursor: pointer; border: none;
  transition: all 0.2s;
}
.btn-sync { background: #334155; color: white; }
.btn-sync:active { background: #475569; }
.btn-exit { background: transparent; border: 2px solid #334155; color: #94a3b8; }

.terminal-body { 
  flex: 1; 
  display: flex; 
  flex-direction: column; 
  overflow: hidden; 
  padding: 1.5rem; 
}

.terminal-footer {
  background: #1e293b;
  padding: 0.75rem 2rem;
  border-top: 2px solid #334155;
}

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>
