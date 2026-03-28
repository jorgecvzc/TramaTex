<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserMenu from './UserMenu.vue'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
const showProduct = ref(false)
const showSales = ref(false)
const showMES = ref(false)

function closeAllDropdowns() {
  showProduct.value = false
  showSales.value = false
  showMES.value = false
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

onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleDocumentClick)
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
            <span class="material-symbols-outlined">dashboard_customize</span>
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/parties" class="nav-link" active-class="active" title="Entidades">
            <span class="material-symbols-outlined">groups_2</span>
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
            <li>
              <RouterLink to="/master-data/attributes" class="dropdown-item" title="Atributos" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">rule</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/master-data/brands" class="dropdown-item" title="Marcas" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">branding_watermark</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/master-data/product-groups" class="dropdown-item" title="Categorías" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">category</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/products/pricing" class="dropdown-item" title="Consulta de Precios" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">price_check</span>
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
            <li>
              <RouterLink to="/sales/tickets/new" class="dropdown-item" title="Nuevo Ticket" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">point_of_sale</span>
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
              <RouterLink to="/mes/tasks" class="dropdown-item" title="Tareas" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">assignment</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/positions" class="dropdown-item" title="Posiciones" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">factory</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/work-types" class="dropdown-item" title="Tipos de trabajo" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">account_tree</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/work-setups" class="dropdown-item" title="Configuraciones" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">settings_input_component</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/work-orders" class="dropdown-item" title="Órdenes de trabajo" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">pending_actions</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/terminal" class="dropdown-item" title="Terminal Taller" @click="closeAllDropdowns">
                <span class="material-symbols-outlined">tablet_mac</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li v-if="isAdmin">
          <RouterLink to="/admin/users" class="nav-link" active-class="active" title="Usuarios">
            <span class="material-symbols-outlined">manage_accounts</span>
          </RouterLink>
        </li>
      </ul>

      <UserMenu />
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  background-color: #1b3a6b;
  color: white;
  padding: 0.5rem 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000; /* Incrementado para ser la capa superior absoluta */
}

.navbar-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-brand {
  display: flex;
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

.nav-link .material-symbols-outlined {
  font-size: 24px;
}

.nav-link:hover {
  color: #E6B800;
  background-color: rgba(230, 184, 0, 0.1);
}

.nav-link.active {
  color: #E6B800;
  background-color: rgba(230, 184, 0, 0.15);
}

.chevron {
  opacity: 0.5;
  font-size: 18px !important;
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
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background-color: #1b3a6b;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  list-style: none;
  margin: 0.5rem 0 0 0;
  padding: 0.5rem;
  display: flex;
  gap: 0.25rem;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.25);
  z-index: 1000;
}

.dropdown-item {
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  text-decoration: none;
  padding: 0.75rem;
  border-radius: 4px;
  transition: all 0.2s;
  min-width: 48px;
}

.dropdown-item .material-symbols-outlined {
  font-size: 20px;
}

.dropdown-item:hover {
  background-color: rgba(230, 184, 0, 0.1);
  color: #E6B800;
}

@media (max-width: 768px) {
  .nav-menu {
    gap: 0.25rem;
  }

  .nav-link {
    padding: 0.5rem;
    min-width: 40px;
  }
}
</style>
