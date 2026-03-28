<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="show" class="dialog-overlay" @mousedown="$emit('close')">
        <div
          ref="dialogRef"
          class="dialog"
          role="dialog"
          aria-modal="true"
          :class="[size, { 'has-icon': !!icon }]"
          @mousedown.stop
        >
          <header v-if="!hideHeader" class="dialog-header">
            <div class="header-content">
              <span v-if="icon" class="material-symbols-outlined">{{ icon }}</span>
              <h2>{{ title }}</h2>
            </div>
            <button @click="$emit('close')" class="btn-icon" aria-label="Cerrar">
              <span class="material-symbols-outlined">close</span>
            </button>
          </header>

          <main class="dialog-body" :class="contentClass">
            <slot />
          </main>

          <footer v-if="!hideActions" class="dialog-footer">
            <button class="btn btn-outline" @click="$emit('close')">Cancelar</button>
            <button
              class="btn"
              :class="confirmClass"
              @click="$emit('confirm')"
              :disabled="isConfirming"
            >
              <span v-if="isConfirming" class="spinner"></span>
              <span v-else>{{ confirmText }}</span>
            </button>
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted, nextTick } from 'vue'

const dialogRef = ref<HTMLElement | null>(null)
const props = withDefaults(defineProps<{
  show: boolean
  title: string
  icon?: string
  size?: 'md' | 'lg' | 'xl'
  contentClass?: string
  hideHeader?: boolean
  hideActions?: boolean
  isConfirming?: boolean
  confirmText?: string
  confirmClass?: 'btn-primary' | 'btn-secondary' | 'btn-danger'
  initialFocus?: HTMLElement | null
}>(), {
  size: 'md',
  confirmText: 'Confirmar',
  confirmClass: 'btn-primary',
})

const emit = defineEmits(['close', 'confirm'])

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
  }
}

watch(() => props.show, (isShown) => {
  if (isShown) {
    document.addEventListener('keydown', handleKeydown)
    nextTick(() => {
      if (props.initialFocus) {
        props.initialFocus.focus()
      } else {
        dialogRef.value?.focus()
      }
    })
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
.fade-enter-active .dialog,
.fade-leave-active .dialog {
  transition: transform 0.2s ease;
}
.fade-enter-from .dialog,
.fade-leave-to .dialog {
  transform: scale(0.95);
}

.dialog {
  background: white;
  border-radius: 14px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  display: flex;
  flex-direction: column;
  max-height: 90vh;
  width: 90%;
}
.dialog:focus {
  outline: none;
}
.dialog.md { max-width: 550px; }
.dialog.lg { max-width: 800px; }
.dialog.xl { max-width: 1100px; }

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.75rem;
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}
.header-content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--color-text-primary);
}
.header-content .material-symbols-outlined { color: var(--color-primary); }
.btn-icon { background: none; border: none; padding: 0.5rem; border-radius: 50%; cursor: pointer; color: var(--color-text-secondary); }
.btn-icon:hover { background: var(--color-background-soft); }

.dialog-body {
  padding: 1.5rem 1.75rem;
  overflow-y: auto;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1rem 1.75rem;
  border-top: 1px solid var(--color-border);
  background: var(--color-background);
  border-bottom-left-radius: 14px;
  border-bottom-right-radius: 14px;
  flex-shrink: 0;
}
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  border: 1px solid transparent;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-outline {
  background: white;
  border-color: var(--color-border-strong);
  color: var(--color-text-secondary);
}
.btn-outline:hover { background: var(--color-background-soft); }

.btn-primary { background: var(--color-primary); color: white; }
.btn-primary:hover { background: var(--color-primary-dark); }
.btn-secondary { background: var(--color-secondary); color: white; }
.btn-secondary:hover { background: var(--color-secondary-dark); }
.btn-danger { background: var(--color-danger); color: white; }
.btn-danger:hover { background: var(--color-danger-dark); }

.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.5);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
