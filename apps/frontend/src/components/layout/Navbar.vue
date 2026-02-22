<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserMenu from './UserMenu.vue'
import { Home, Package, Users, DollarSign, Clipboard, ShoppingCart, ScrollText, Receipt, FolderOpen, Tag, Folder, Zap, User, Wrench } from 'lucide-vue-next'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
const showMasterData = ref(false)
const showSales = ref(false)
const showMES = ref(false)

function closeAllDropdowns() {
  showSales.value = false
  showMasterData.value = false
  showMES.value = false
}

function toggleSales() {
  const next = !showSales.value
  closeAllDropdowns()
  showSales.value = next
}

function toggleMasterData() {
  const next = !showMasterData.value
  closeAllDropdowns()
  showMasterData.value = next
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
            <Home :size="24" />
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/products" class="nav-link" active-class="active" title="Productos">
            <Package :size="24" />
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/parties" class="nav-link" active-class="active" title="Entidades">
            <Users :size="24" />
          </RouterLink>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" title="Ventas" @click.stop="toggleSales">
            <DollarSign :size="24" />
          </button>
          <ul v-if="showSales" class="dropdown-menu">
            <li>
              <RouterLink to="/sales/quotes" class="dropdown-item" title="Presupuestos" @click="closeAllDropdowns">
                <Clipboard :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/orders" class="dropdown-item" title="Pedidos" @click="closeAllDropdowns">
                <ShoppingCart :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/delivery-notes" class="dropdown-item" title="Albaranes" @click="closeAllDropdowns">
                <ScrollText :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/invoices" class="dropdown-item" title="Facturas" @click="closeAllDropdowns">
                <Receipt :size="20" />
              </RouterLink>
            </li>
          </ul>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" title="Datos Maestros" @click.stop="toggleMasterData">
            <FolderOpen :size="24" />
          </button>
          <ul v-if="showMasterData" class="dropdown-menu">
            <li>
              <RouterLink to="/master-data/brands" class="dropdown-item" title="Marcas" @click="closeAllDropdowns">
                <Tag :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/master-data/product-groups" class="dropdown-item" title="Categorías" @click="closeAllDropdowns">
                <Folder :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/master-data/attributes" class="dropdown-item" title="Atributos" @click="closeAllDropdowns">
                <Zap :size="20" />
              </RouterLink>
            </li>
          </ul>
        </li>
        <li class="dropdown" @click.stop>
          <button type="button" class="nav-link dropdown-toggle" title="MES" @click.stop="toggleMES">
            <Zap :size="24" />
          </button>
          <ul v-if="showMES" class="dropdown-menu">
            <li>
              <RouterLink to="/mes/dashboard" class="dropdown-item" title="Dashboard MES" @click="closeAllDropdowns">
                <Home :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/tasks" class="dropdown-item" title="Tareas MES" @click="closeAllDropdowns">
                <Clipboard :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/positions" class="dropdown-item" title="Puestos MES" @click="closeAllDropdowns">
                <Users :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/service-groups" class="dropdown-item" title="Grupos de Servicio MES" @click="closeAllDropdowns">
                <Folder :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/works" class="dropdown-item" title="Trabajos MES" @click="closeAllDropdowns">
                <ShoppingCart :size="20" />
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/mes/terminal" class="dropdown-item" title="Terminal Taller" @click="closeAllDropdowns">
                <Wrench :size="20" />
              </RouterLink>
            </li>
          </ul>
        </li>
        <li v-if="isAdmin">
          <RouterLink to="/admin/users" class="nav-link" active-class="active" title="Usuarios">
            <User :size="24" />
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
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.navbar-container {
  max-width: 1200px;
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
  color: #ffd700;
  text-decoration: none;
  transition: color 0.2s;
}

.logo:hover {
  color: #ffed4e;
}

.nav-menu {
  display: flex;
  list-style: none;
  gap: 2rem;
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
}

.nav-link :deep(svg) { /* Target Lucide SVG directly */
  width: 1.5rem; /* Equivalent to 24px */
  height: 1.5rem; /* Equivalent to 24px */
}

.nav-link:hover {
  color: #ffd700;
  background-color: rgba(255, 215, 0, 0.1);
}

.nav-link.active {
  background-color: rgba(255, 215, 0, 0.2);
  color: #ffd700;
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
  gap: 0.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  z-index: 1000;
}

.dropdown-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: white;
  text-decoration: none;
  padding: 0.75rem;
  transition: all 0.2s;
  min-width: 48px;
}

.dropdown-item :deep(svg) { /* Target Lucide SVG directly */
  width: 1.25rem; /* Equivalent to 20px */
  height: 1.25rem; /* Equivalent to 20px */
}

.dropdown-item:hover {
  background-color: rgba(255, 215, 0, 0.1);
  color: #ffd700;
}

@media (max-width: 768px) {
  .nav-menu {
    gap: 0.5rem;
  }

  .nav-link {
    padding: 0.5rem;
    min-width: 40px;
  }
  
  .nav-link :deep(svg) {
    width: 1.25rem;
    height: 1.25rem;
  }
  
  .dropdown-item :deep(svg) {
    width: 1rem;
    height: 1rem;
  }
}
</style>
