<template>
  <aside class="side-navbar" :class="{ 'is-expanded': isExpanded }">
    <div class="navbar-header">
      <button class="logo-trigger" @click="isLauncherOpen = true" title="Abrir Menú Principal">
        <span class="logo-icon">TX</span>
        <span class="logo-text">TramaTex</span>
      </button>
    </div>

    <div class="menu-section">
      <div class="menu">
        <RouterLink to="/dashboard" class="menu-item" active-class="is-active" title="Dashboard General">
          <LayoutDashboard :size="24" class="icon" />
          <span class="text">Dashboard</span>
        </RouterLink>

        <RouterLink to="/sales/dashboard" class="menu-item" active-class="is-active" title="Ventas">
          <CreditCard :size="24" class="icon" />
          <span class="text">Ventas</span>
        </RouterLink>

        <RouterLink to="/products/dashboard" class="menu-item" active-class="is-active" title="Catálogo">
          <Package :size="24" class="icon" />
          <span class="text">Catálogo</span>
        </RouterLink>

        <RouterLink to="/parties/dashboard" class="menu-item" active-class="is-active" title="Entidades">
          <Users :size="24" class="icon" />
          <span class="text">Entidades</span>
        </RouterLink>

        <RouterLink to="/mes/dashboard" class="menu-item" active-class="is-active" title="Taller (MES)">
          <Factory :size="24" class="icon" />
          <span class="text">Taller</span>
        </RouterLink>
        
        <RouterLink v-if="isAdmin" to="/admin/users" class="menu-item admin-item" active-class="is-active" title="Administración">
          <ShieldCheck :size="24" class="icon" />
          <span class="text">Administración</span>
        </RouterLink>
      </div>
    </div>

    <div class="flex"></div>

    <div class="menu-footer">
      <!-- Unificado: Menú de Ayuda -->
      <HelpMenu 
        :sidebar-expanded="isExpanded" 
        @open-shortcuts="handleOpenShortcuts"
        @open-contextual-help="handleOpenHelp"
      />

      <div class="menu-toggle-wrap">
        <button class="menu-toggle" @click="toggleMenu" title="Colapsar/Expandir menú">
          <ChevronsRight :size="20" class="toggle-icon" />
        </button>
      </div>
    </div>
  </aside>

  <!-- App Launcher (Menú Clásico) -->
  <AppLauncher :is-open="isLauncherOpen" @close="isLauncherOpen = false" />
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { 
  LayoutDashboard, 
  CreditCard, 
  Package, 
  Users, 
  Factory, 
  ShieldCheck, 
  ChevronsRight 
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import AppLauncher from './AppLauncher.vue'
import HelpMenu from './HelpMenu.vue'

const router = useRouter()
const isExpanded = ref(true)
const isLauncherOpen = ref(false)
const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)

function handleOpenShortcuts() {
  window.dispatchEvent(new KeyboardEvent('keydown', { key: '?' }))
}

function handleOpenHelp() {
  window.dispatchEvent(new KeyboardEvent('keydown', { key: 'F1' }))
}

function handleGlobalKeydown(e) {
  // 1. Toggle Sidebar: Ctrl + B
  if (e.ctrlKey && e.key === 'b') {
    e.preventDefault()
    toggleMenu()
  }

  // 2. Module Navigation: Alt + 1, 2, 3, 4, 5
  if (e.altKey) {
    const map = {
      '1': '/dashboard',
      '2': '/sales/dashboard',
      '3': '/products/dashboard',
      '4': '/parties/dashboard',
      '5': '/mes/dashboard'
    }
    if (map[e.key]) {
      e.preventDefault()
      router.push(map[e.key])
    }
  }
}

onMounted(() => {
  const saved = localStorage.getItem('sidebar-expanded')
  if (saved !== null) {
    isExpanded.value = saved === 'true'
  }
  window.addEventListener('keydown', handleGlobalKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})

const toggleMenu = () => {
  isExpanded.value = !isExpanded.value
  localStorage.setItem('sidebar-expanded', isExpanded.value)
}
</script>

<style scoped>
.side-navbar {
  --sidebar-bg-color: var(--color-secondary-light);
  --sidebar-item-hover: rgba(255, 255, 255, 0.1);
  --sidebar-item-active: var(--color-primary);
  --sidebar-text-color: #ffffff;
  --header-height: 76px;
  
  display: flex;
  flex-direction: column;
  background-color: var(--sidebar-bg-color);
  color: var(--sidebar-text-color);
  width: 64px;
  height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  overflow: hidden;
  padding: 0.75rem 0.5rem;
  padding-top: calc(var(--header-height) + 0.75rem);
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 2000;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
}

.side-navbar.is-expanded {
  width: 240px;
}

.navbar-header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-bottom: 2rem;
  width: 100%;
  padding-left: 0.35rem; /* Ajuste para alinear con items */
}

.logo-trigger {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  width: 100%;
  color: inherit;
  transition: all 0.3s;
}

.logo-icon {
  min-width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.2);
  color: var(--color-primary);
  font-size: 1rem;
  font-weight: 900;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.logo-text {
  color: var(--sidebar-text-color);
  font-size: 1.25rem;
  font-weight: 800;
  margin-left: 1rem;
  white-space: nowrap;
  display: none;
}

.is-expanded .logo-text {
  display: block;
  animation: fade-in 0.4s ease forwards;
}

@keyframes fade-in {
  from { opacity: 0; transform: translateX(-10px); }
  to { opacity: 1; transform: translateX(0); }
}

.menu-section {
  width: 100%;
}

.menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  width: 100%;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  text-decoration: none;
  padding: 0.75rem 0.35rem; /* Alineado a la izquierda */
  color: var(--sidebar-text-color);
  border-radius: 10px;
  transition: all 0.2s ease;
  width: 100%;
}

.menu-item .icon {
  color: var(--color-primary);
  flex-shrink: 0;
  width: 2.5rem; /* Mismo ancho que el logo para alinear centros */
  display: flex;
  justify-content: center;
}

.menu-item .text {
  font-size: 0.95rem;
  font-weight: 600;
  margin-left: 0.5rem;
  color: #ffffff;
  display: none;
}

.is-expanded .menu-item .text {
  display: block;
  animation: fade-in 0.4s ease forwards;
}

.menu-item:hover {
  background-color: var(--sidebar-item-hover);
}

.menu-item.is-active {
  background-color: rgba(255, 255, 255, 0.2);
  color: var(--sidebar-item-active);
}

.admin-item {
  margin-top: 2rem;
}

.menu-footer {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.menu-toggle-wrap {
  display: flex;
  justify-content: flex-start; /* Toggle también a la izquierda */
  padding-left: 0.5rem;
  margin-bottom: 20px;
}

.menu-toggle {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: #ffffff;
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: 0.3s ease;
}

.menu-toggle:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.1);
}

.menu-toggle .toggle-icon {
  transition: 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.is-expanded .menu-toggle .toggle-icon {
  transform: rotate(-180deg);
}

.flex {
  flex: 1;
}
</style>
