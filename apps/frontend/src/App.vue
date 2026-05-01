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
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, RouterView } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import SideNavbar from '@/components/layout/SideNavbar.vue'
import ToastContainer from '@/components/shared/ToastContainer.vue'

const route = useRoute()
const showAppChrome = computed(() => route.path !== '/login')
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
