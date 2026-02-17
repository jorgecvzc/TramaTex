<template>
  <nav class="navbar">
    <div class="navbar-container">
      <div class="navbar-brand">
        <RouterLink to="/dashboard" class="logo">TramaTex</RouterLink>
      </div>

      <ul class="nav-menu">
        <li>
          <RouterLink to="/dashboard" class="nav-link" active-class="active" title="Dashboard">
            <span class="icon">🏠</span>
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/products" class="nav-link" active-class="active" title="Productos">
            <span class="icon">📦</span>
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/parties" class="nav-link" active-class="active" title="Entidades">
            <span class="icon">👥</span>
          </RouterLink>
        </li>
        <li class="dropdown" @mouseenter="showSales = true" @mouseleave="showSales = false">
          <span class="nav-link dropdown-toggle" title="Ventas">
            <span class="icon">💰</span>
          </span>
          <ul v-if="showSales" class="dropdown-menu">
            <li>
              <RouterLink to="/sales/quotes" class="dropdown-item" title="Presupuestos">
                <span class="icon">📝</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/orders" class="dropdown-item" title="Pedidos">
                <span class="icon">🛒</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/delivery-notes" class="dropdown-item" title="Albaranes">
                <span class="icon">📋</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/sales/invoices" class="dropdown-item" title="Facturas">
                <span class="icon">🧾</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li class="dropdown" @mouseenter="showMasterData = true" @mouseleave="showMasterData = false">
          <span class="nav-link dropdown-toggle" title="Datos Maestros">
            <span class="icon">🗂️</span>
          </span>
          <ul v-if="showMasterData" class="dropdown-menu">
            <li>
              <RouterLink to="/master-data/brands" class="dropdown-item" title="Marcas">
                <span class="icon">🏷️</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/master-data/product-groups" class="dropdown-item" title="Categorías">
                <span class="icon">📁</span>
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/master-data/attributes" class="dropdown-item" title="Atributos">
                <span class="icon">⚡</span>
              </RouterLink>
            </li>
          </ul>
        </li>
        <li v-if="isAdmin">
          <RouterLink to="/admin/users" class="nav-link" active-class="active" title="Usuarios">
            <span class="icon">👤</span>
          </RouterLink>
        </li>
      </ul>

      <UserMenu />
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserMenu from './UserMenu.vue'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
const showMasterData = ref(false)
const showSales = ref(false)
</script>

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

.nav-link .icon {
  font-size: 1.5rem;
  line-height: 1;
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

.dropdown-item .icon {
  font-size: 1.25rem;
  line-height: 1;
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
  
  .nav-link .icon {
    font-size: 1.25rem;
  }
  
  .dropdown-item .icon {
    font-size: 1rem;
  }
}
</style>
