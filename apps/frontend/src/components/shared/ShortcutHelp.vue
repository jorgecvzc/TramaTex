<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="show" class="shortcuts-overlay" @mousedown="$emit('close')">
        <div class="shortcuts-modal" @mousedown.stop>
          <header class="modal-header">
            <Keyboard :size="24" />
            <h2>Atajos de Teclado</h2>
            <button class="btn-close" @click="$emit('close')"><X :size="20" /></button>
          </header>
          
          <div class="modal-body">
            <section class="shortcut-group">
              <h3>Globales</h3>
              <div class="shortcut-item"><span>Guardar cambios</span> <kbd>Ctrl</kbd> + <kbd>S</kbd></div>
              <div class="shortcut-item"><span>Cerrar / Atrás</span> <kbd>Esc</kbd></div>
              <div class="shortcut-item"><span>Buscar en listado</span> <kbd>/</kbd> o <kbd>Ctrl</kbd> + <kbd>K</kbd></div>
              <div class="shortcut-item"><span>Colapsar barra lateral</span> <kbd>Ctrl</kbd> + <kbd>B</kbd></div>
              <div class="shortcut-item"><span>Ayuda atajos</span> <kbd>?</kbd></div>
            </section>

            <section class="shortcut-group">
              <h3>Navegación de Módulos</h3>
              <div class="shortcut-item"><span>Dashboard / Ventas / Productos</span> <kbd>Alt</kbd> + <kbd>1</kbd>...<kbd>3</kbd></div>
              <div class="shortcut-item"><span>Entidades / Taller</span> <kbd>Alt</kbd> + <kbd>4</kbd>...<kbd>5</kbd></div>
            </section>

            <section class="shortcut-group">
              <h3>Documentos con Líneas</h3>
              <div class="shortcut-item"><span>Subir valor</span> <kbd>↑</kbd></div>
              <div class="shortcut-item"><span>Bajar valor</span> <kbd>↓</kbd></div>
              <div class="shortcut-item"><span>Siguiente campo</span> <kbd>Enter</kbd></div>
              <div class="shortcut-item"><span>Añadir línea</span> <kbd>Insert</kbd></div>
              <div class="shortcut-item"><span>Eliminar línea</span> <kbd>Ctrl</kbd> + <kbd>Supr</kbd></div>
            </section>

            <section class="shortcut-group">
              <h3>Listados y Selectores</h3>
              <div class="shortcut-item"><span>Navegar resultados</span> <kbd>↑</kbd> <kbd>↓</kbd></div>
              <div class="shortcut-item"><span>Seleccionar</span> <kbd>Enter</kbd></div>
              <div class="shortcut-item"><span>Cerrar desplegable</span> <kbd>Esc</kbd></div>
            </section>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { Keyboard, X } from 'lucide-vue-next'

defineProps<{ show: boolean }>()
defineEmits(['close'])
</script>

<style scoped>
.shortcuts-overlay {
  position: fixed; inset: 0; background: rgba(15, 23, 42, 0.7);
  backdrop-filter: blur(4px); display: flex; align-items: center;
  justify-content: center; z-index: 10000;
}

.shortcuts-modal {
  width: 100%; max-width: 500px; background: white;
  border-radius: 16px; box-shadow: var(--box-shadow-lg);
  overflow: hidden; border: 1px solid var(--color-border);
}

.modal-header {
  padding: 1.25rem 1.5rem; background: var(--color-background-soft);
  border-bottom: 1px solid var(--color-border); display: flex;
  align-items: center; gap: 0.75rem;
}
.modal-header h2 { margin: 0; font-size: 1.1rem; flex: 1; }
.btn-close { background: none; border: none; color: var(--color-text-secondary); cursor: pointer; }

.modal-body { padding: 1.5rem; display: flex; flex-direction: column; gap: 1.5rem; }

.shortcut-group h3 {
  font-size: 0.75rem; font-weight: 800; text-transform: uppercase;
  color: var(--color-text-secondary); letter-spacing: 0.05em;
  margin-bottom: 0.75rem; border-bottom: 1px solid var(--color-border-soft);
  padding-bottom: 0.25rem;
}

.shortcut-item {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 0.5rem; font-size: 0.9rem;
}

.shortcut-item span { color: var(--color-text-primary); }
kbd { font-size: 0.75rem; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
