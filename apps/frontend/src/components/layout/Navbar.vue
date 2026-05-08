<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { 
  LayoutDashboard, 
  Users, 
  Package, 
  List, 
  PlusSquare, 
  CreditCard, 
  FileText, 
  ShoppingCart, 
  Truck, 
  ReceiptText, 
  Factory, 
  LineChart, 
  Tablet, 
  UserCog, 
  Search 
} from 'lucide-vue-next'
import UserMenu from './UserMenu.vue'
import GlobalSearch from '../shared/GlobalSearch.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()
const isAdmin = computed(() => authStore.isAdmin)
const showProduct = ref(false)
const showSales = ref(false)
const showMES = ref(false)
const showSearch = ref(false)

function closeAllDropdowns() {
  showProduct.value = false
  showSales.value = false
  showMES.value = false
}

function openSearch() {
  closeAllDropdowns()
  showSearch.value = true
}

function toggleProduct() {
  const next = !showProduct.value
  closeAllDropdowns()
  showProduct.value = next
}

function toggleSales() {
  const next = !showSales.value
  closeAllDropdowns()
  showSales.value = next
}

function toggleMES() {
  const next = !showMES.value
  closeAllDropdowns()
  showMES.value = next
}

function handleDocumentClick(event: MouseEvent) {
  // Solo cerrar si el clic es fuera de un toggle
  const target = event.target as HTMLElement
  if (!target.closest('.dropdown-toggle')) {
    closeAllDropdowns()
  }
}

function handleShortcuts(e: KeyboardEvent) {
  // 1. BUSQUEDA GLOBAL: Ctrl + K
  if (e.ctrlKey && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    openSearch()
    return
  }

  // 2. NUEVO ELEMENTO: Alt + N
  if (e.altKey && e.key.toLowerCase() === 'n') {
    e.preventDefault()
    handleNewAction()
    return
  }

  // 3. GUARDAR / ENVIAR: Ctrl + Enter
  if (e.ctrlKey && e.key === 'Enter') {
    e.preventDefault()
    triggerSubmit()
    return
  }

  // 4. REFRESCAR DATOS: Alt + R
  if (e.altKey && e.key.toLowerCase() === 'r') {
    e.preventDefault()
    triggerRefresh()
    return
  }
}

function handleNewAction() {
  const path = route.path
  if (path.startsWith('/parties')) {
    router.push('/parties/new')
  } else if (path.startsWith('/products')) {
    router.push('/products/new')
  } else if (path.startsWith('/sales/orders')) {
    router.push('/sales/orders/new')
  } else if (path.startsWith('/sales/quotes')) {
    router.push('/sales/quotes/new')
  } else if (path.startsWith('/sales/tickets')) {
    router.push('/sales/tickets/new')
  } else {
    // Si no estamos en un módulo específico, podemos abrir el buscador o un menú de creación rápida
    toastStore.info('Usa Alt+N dentro de un módulo para crear un nuevo elemento')
  }
}

function triggerSubmit() {
  // Buscamos el botón de submit principal (convención: btn-primary o btn-secondary en footer/actions)
  // O el botón de tipo submit dentro de un formulario
  const submitBtn = document.querySelector('button[type="submit"]') as HTMLButtonElement || 
                    document.querySelector('.form-actions .btn-secondary') as HTMLButtonElement ||
                    document.querySelector('.header-actions .btn-primary') as HTMLButtonElement

  if (submitBtn && !submitBtn.disabled) {
    submitBtn.click()
  } else {
    console.log('No submit button found for Ctrl+Enter')
  }
}

function triggerRefresh() {
  // Emitimos un evento global que los componentes pueden escuchar
  window.dispatchEvent(new CustomEvent('tramatex-refresh'))
  
  // Feedback visual
  toastStore.info('Refrescando datos...')
}

onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
  document.addEventListener('keydown', handleShortcuts)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleDocumentClick)
  document.removeEventListener('keydown', handleShortcuts)
})
</script>

<template>
  <nav class="navbar">
    <div class="navbar-container">
      <div class="navbar-brand">
        <RouterLink to="/dashboard" class="logo">TramaTex</RouterLink>
      </div>

      <ul class="nav-menu">
        <li>
          <RouterLink to="/dashboard" class="nav-link" active-class="active">
            <LayoutDashboard :size="20" />
            <span class="nav-label">Dashboard</span>
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/parties" class="nav-link" active-class="active">
            <Users :size="20" />
            <span class="nav-label">Entidades</span>
          </RouterLink>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" @click.stop="toggleProduct">
            <Package :size="20" />
            <span class="nav-label">Catálogo</span>
          </button>
          <ul v-if="showProduct" class="dropdown-menu">
            <li>
              <RouterLink to="/products" class="dropdown-item" @click="closeAllDropdowns">
                <List :size="20" />
                <span>Listado de Productos</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/products/new" class="dropdown-item" @click="closeAllDropdowns">
                <PlusSquare :size="20" />
                <span>Nuevo Producto</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" @click.stop="toggleSales">
            <CreditCard :size="20" />
            <span class="nav-label">Ventas</span>
          </button>
          <ul v-if="showSales" class="dropdown-menu">
            <li>
              <RouterLink to="/sales/quotes" class="dropdown-item" @click="closeAllDropdowns">
                <FileText :size="20" />
                <span>Presupuestos</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/orders" class="dropdown-item" @click="closeAllDropdowns">
                <ShoppingCart :size="20" />
                <span>Pedidos</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/delivery-notes" class="dropdown-item" @click="closeAllDropdowns">
                <Truck :size="20" />
                <span>Albaranes</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/invoices" class="dropdown-item" @click="closeAllDropdowns">
                <ReceiptText :size="20" />
                <span>Facturas</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" @click.stop="toggleMES">
            <Factory :size="20" />
            <span class="nav-label">Taller</span>
          </button>
          <ul v-if="showMES" class="dropdown-menu">
            <li>
              <RouterLink to="/mes/dashboard" class="dropdown-item" @click="closeAllDropdowns">
                <LineChart :size="20" />
                <span>Panel de control</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/terminal" class="dropdown-item" @click="closeAllDropdowns">
                <Tablet :size="20" />
                <span>Terminal Taller</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li v-if="isAdmin">
          <RouterLink to="/admin/users" class="nav-link" active-class="active">
            <UserCog :size="20" />
            <span class="nav-label">Admin</span>
          </RouterLink>
        </li>
      </ul>

      <div class="navbar-actions">
        <button @click="openSearch" class="search-btn" title="Búsqueda Global (Ctrl+K)">
          <Search :size="18" />
          <span class="search-placeholder">Buscar...</span>
          <span class="kbd-shortcut">Ctrl+K</span>
        </button>
      </div>

      <UserMenu />
    </div>
  </nav>

  <GlobalSearch :show="showSearch" @close="showSearch = false" />
</template>

<style scoped>
.navbar {
  background-color: var(--color-secondary);
  color: white;
  height: 76px; /* Altura fija para alineación exacta */
  display: flex;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 2100;
}

.navbar-container {
  width: 100%;
  max-width: 1300px;
  margin: 0 auto;
  padding: 0 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--color-primary);
  text-decoration: none;
  font-family: var(--font-family-brand);
  font-style: italic;
  letter-spacing: -0.025em;
  margin-right: 2rem;
}

.nav-menu {
  display: flex;
  list-style: none;
  gap: 0.25rem;
  margin: 0;
  padding: 0;
}

.nav-link {
  color: white;
  text-decoration: none;
  transition: all 0.2s;
  padding: 0.5rem;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.2rem;
  width: 80px;
  height: 60px;
}

.nav-link:hover, 
.nav-link.active {
  color: var(--color-primary);
  background-color: rgba(255, 255, 255, 0.1);
}

.nav-label {
  font-size: 0.6rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.dropdown {
  position: relative;
}

.dropdown-toggle {
  cursor: pointer;
  user-select: none;
  border: none;
  background: transparent;
  font: inherit;
  color: inherit;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background-color: var(--color-secondary);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  list-style: none;
  margin: 0.5rem 0 0 0;
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.25);
  z-index: 1000;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: white;
  text-decoration: none;
  padding: 0.75rem 1rem;
  border-radius: 4px;
  transition: all 0.2s;
  min-width: 180px;
}

.dropdown-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: var(--color-primary);
}

.navbar-actions {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 0 2rem;
}

.search-btn {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 0.35rem 0.85rem;
  min-height: 36px;
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  width: 100%;
  max-width: 400px;
  transition: all 0.2s;
}

.search-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
  color: white;
}

.search-placeholder {
  flex: 1;
  text-align: left;
  font-size: 0.84rem;
}

.kbd-shortcut {
  font-size: 0.68rem;
  font-family: monospace;
  background: rgba(0, 0, 0, 0.2);
  padding: 0.08rem 0.35rem;
  border-radius: 4px;
}

@media (min-width: 1201px) {
  .nav-menu {
    display: none;
  }

  .navbar-actions {
    justify-content: center;
    padding-left: 0;
  }

  .logo {
    margin-right: 0.75rem;
  }
}
</style>
