<script setup lang="ts">
/**
 * BaseDialog.vue - Maestro de Diálogos TramaTex
 * 
 * Proporciona un contenedor modal estandarizado con:
 * - Backdrop con efecto blur
 * - Cabecera con icono y botón de cierre
 * - Acciones consistentes (Cancelar/Aceptar)
 */
import { watch, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  show: boolean
  title: string
  icon?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  confirmText?: string
  confirmClass?: string
  isConfirming?: boolean
  disableConfirm?: boolean
  hideActions?: boolean
}>()

const emit = defineEmits(['close', 'confirm'])

function handleBackdropClick() {
  emit('close')
}

function handleEscapeKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.show) {
    emit('close')
  }
}

// Escuchador global de teclado
onMounted(() => {
  window.addEventListener('keydown', handleEscapeKey)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleEscapeKey)
})

// Bloquear scroll del body cuando el modal está abierto
watch(() => props.show, (newVal) => {
  if (newVal) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}, { immediate: true })
</script>

<template>
  <Teleport to="body">
    <Transition name="dialog-fade">
      <div v-if="show" class="dialog-fixed-overlay" @click.self="handleBackdropClick">
        <div :class="['dialog-surface', size || 'md']">
          <!-- CABECERA -->
          <header class="dialog-header">
            <div class="dialog-title-group">
              <span v-if="icon" class="material-symbols-outlined dialog-icon">
                {{ icon }}
              </span>
              <h2>{{ title }}</h2>
            </div>
            <button class="btn-close" @click="emit('close')" aria-label="Cerrar">
              <span class="material-symbols-outlined">close</span>
            </button>
          </header>

          <!-- CUERPO -->
          <main class="dialog-body custom-scrollbar">
            <slot></slot>
          </main>

          <!-- ACCIONES -->
          <footer v-if="!hideActions" class="dialog-footer">
            <slot name="actions">
              <button class="btn btn-outline" @click="emit('close')" :disabled="isConfirming">
                Cancelar
              </button>
              <button 
                :class="['btn', confirmClass || 'btn-primary']" 
                @click="emit('confirm')" 
                :disabled="isConfirming || disableConfirm"
              >
                <span v-if="isConfirming" class="material-symbols-outlined spin mr-2">sync</span>
                <span>{{ confirmText || 'Confirmar' }}</span>
              </button>
            </slot>
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style>
/* ESTILOS GLOBALES PARA EL OVERLAY (Fuera de scoped para garantizar el fixed) */
.dialog-fixed-overlay {
  position: fixed !important;
  inset: 0 !important;
  width: 100vw !important;
  height: 100vh !important;
  background: rgba(15, 23, 42, 0.75) !important;
  backdrop-filter: blur(8px) !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  z-index: 10000 !important;
  padding: 1.5rem;
}

.dialog-fade-enter-active, .dialog-fade-leave-active {
  transition: opacity 0.3s ease;
}
.dialog-fade-enter-from, .dialog-fade-leave-to {
  opacity: 0;
}
</style>

<style scoped>
.dialog-surface {
  background: white !important;
  border-radius: 16px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  max-height: 90vh;
  width: 100%;
  position: relative;
}

/* Tamaños */
.sm { max-width: 400px; }
.md { max-width: 600px; }
.lg { max-width: 800px; }
.xl { max-width: 1140px; }

.dialog-header {
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.dialog-title-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.dialog-title-group h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: #1e293b;
}

.dialog-icon {
  color: var(--color-primary);
  font-size: 24px;
}

.btn-close {
  background: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 8px;
  display: flex;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #f1f5f9;
  color: #ef4444;
}

.dialog-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.dialog-footer {
  padding: 1.25rem 1.5rem;
  background: #f8fafc;
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  border-bottom-left-radius: 16px;
  border-bottom-right-radius: 16px;
}

.spin { animation: rotate 1s linear infinite; }
@keyframes rotate { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.mr-2 { margin-right: 0.5rem; }

.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
</style>
