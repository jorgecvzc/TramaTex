<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserMenu from './UserMenu.vue'
import GlobalSearch from '../shared/GlobalSearch.vue'

const authStore = useAuthStore()
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

function handleDocumentClick() {
  closeAllDropdowns()
}

function handleShortcuts(e: KeyboardEvent) {
  if (e.ctrlKey && e.key === 'k') {
    e.preventDefault()
    openSearch()
  }
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
          <RouterLink to="/dashboard" class="nav-link" active-class="active" title="Dashboard">
            <span class="material-symbols-outlined">dashboard</span>
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/parties" class="nav-link" active-class="active" title="Entidades">
            <span class="material-symbols-outlined">groups</span>
          </RouterLink>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" title="Producto" @click.stop="toggleProduct">
            <span class="material-symbols-outlined">inventory_2</span>
            <span class="material-symbols-outlined chevron">expand_more</span>
          </button>
          <ul v-if="showProduct" class="dropdown-menu">
            <li>
              <RouterLink to="/products" class="dropdown-item" title="Catálogo de Productos" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">list_alt</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/products/new" class="dropdown-item" title="Nuevo Producto" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">add_box</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" title="Ventas" @click.stop="toggleSales">
            <span class="material-symbols-outlined">payments</span>
            <span class="material-symbols-outlined chevron">expand_more</span>
          </button>
          <ul v-if="showSales" class="dropdown-menu">
            <li>
              <RouterLink to="/sales/quotes" class="dropdown-item" title="Presupuestos" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">request_quote</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/orders" class="dropdown-item" title="Pedidos" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">shopping_cart</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/delivery-notes" class="dropdown-item" title="Albaranes" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">local_shipping</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/invoices" class="dropdown-item" title="Facturas" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">receipt_long</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" title="MES" @click="toggleMES">
            <span class="material-symbols-outlined">precision_manufacturing</span>
            <span class="material-symbols-outlined chevron">expand_more</span>
          </button>
          <ul v-if="showMES" class="dropdown-menu">
            <li>
              <RouterLink to="/mes/dashboard" class="dropdown-item" title="Panel de control" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">monitoring</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/terminal" class="dropdown-item" title="Terminal Taller" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">tablet_mac</span>
              </RouterLink>
            </li>
          </ul>
        </li>
      </ul>

      <div class="navbar-actions">
        <button @click="openSearch" class="search-btn" title="Búsqueda Global (Ctrl+K)">
          <span class="material-symbols-outlined">search</span>
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
  background-color: #1b3a6b;
  color: white;
  padding: 0.5rem 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}
.navbar-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: #E6B800;
  text-decoration: none;
  font-family: var(--font-family-brand);
  font-style: italic;
  letter-spacing: -0.025em;
  margin-right: 2rem;
}
.nav-menu {
  display: flex;
  list-style: none;
  gap: 0.5rem;
  margin: 0;
  padding: 0;
}
.nav-link {
  color: white;
  text-decoration: none;
  transition: all 0.2s;
  padding: 0.75rem;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 48px;
  gap: 0.25rem;
}
.nav-link:hover, .nav-link.active {
  color: #E6B800;
  background-color: rgba(230, 184, 0, 0.1);
}
.chevron { opacity: 0.5; font-size: 18px !important; }
.dropdown { position: relative; }
.dropdown-toggle { cursor: pointer; user-select: none; border: none; background: transparent; font: inherit; color: inherit; }
.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background-color: #1b3a6b;
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
  min-width: 150px;
}
.dropdown-item:hover {
  background-color: rgba(230, 184, 0, 0.1);
  color: #E6B800;
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
  gap: 0.75rem;
  background: rgba(255,255,255,0.1);
  border: 1px solid rgba(255,255,255,0.2);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  color: rgba(255,255,255,0.7);
  cursor: pointer;
  width: 100%;
  max-width: 400px;
  transition: all 0.2s;
}
.search-btn:hover {
  background: rgba(255,255,255,0.15);
  border-color: rgba(255,255,255,0.3);
  color: white;
}
.search-placeholder { flex: 1; text-align: left; font-size: 0.9rem; }
.kbd-shortcut {
  font-size: 0.75rem;
  font-family: monospace;
  background: rgba(0,0,0,0.2);
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
}
</style>
