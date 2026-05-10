<template>
  <template v-if="showAppChrome">
    <div class="app-shell">
      <Navbar />
      <div class="app-layout">
        <SideNavbar class="app-sidebar" />
        <main class="app-main">
          <RouterView />
        </main>
      </div>
    </div>
  </template>
  <RouterView v-else />
  <ToastContainer />
  <ShortcutHelp :show="showShortcutHelp" @close="showShortcutHelp = false" />
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, RouterView } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import SideNavbar from '@/components/layout/SideNavbar.vue'
import ToastContainer from '@/components/shared/ToastContainer.vue'
import ShortcutHelp from '@/components/shared/ShortcutHelp.vue'

const route = useRoute()
const showAppChrome = computed(() => route.path !== '/login')
const showShortcutHelp = ref(false)

function handleGlobalKeydown(e) {
  // 1. Guardado Rápido: Ctrl + S
  if (e.ctrlKey && e.key === 's') {
    e.preventDefault()
    window.dispatchEvent(new CustomEvent('tramatex-save'))
  }

  // 2. Búsqueda Instantánea: Ctrl + K o / (si no estamos en un input)
  const isInput = ['INPUT', 'TEXTAREA', 'SELECT'].includes(document.activeElement?.tagName)
  
  if ((e.ctrlKey && e.key === 'k') || (e.key === '/' && !isInput)) {
    e.preventDefault()
    window.dispatchEvent(new CustomEvent('tramatex-search'))
  }

  // 3. Ayuda de Atajos: ? (si no estamos en un input)
  if (e.key === '?' && !isInput) {
    showShortcutHelp.value = true
  }

  // 4. Atrás / Cerrar: Esc
  if (e.key === 'Escape') {
    if (showShortcutHelp.value) {
      showShortcutHelp.value = false
    } else {
      window.dispatchEvent(new CustomEvent('tramatex-esc'))
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<style>
.app-shell {
  min-height: 100vh;
  background: var(--color-background);
}

.app-layout {
  display: flex;
  align-items: stretch;
}

.app-sidebar {
  flex-shrink: 0;
  /* Al ser fixed, no ocupa espacio en el flujo, el main necesita margen */
}

.app-main {
  flex: 1;
  min-width: 0;
  margin-left: 64px; /* Ancho de la barra colapsada */
  transition: margin-left 0.3s ease;
}

.app-layout:has(.is-expanded) .app-main {
  margin-left: 240px; /* Ancho de la barra expandida */
}

@media (max-width: 1200px) {
  .app-sidebar {
    display: none;
  }
}
</style>
