<template>
  <aside class="side-navbar" :class="{ 'is-expanded': isExpanded }">
    <div class="navbar-header">
      <RouterLink to="/dashboard" class="logo">
        <img src="/logo.svg" alt="TramaTex Logo" class="logo-icon" />
        <span class="logo-text">TramaTex</span>
      </RouterLink>
    </div>

    <div class="menu-toggle-wrap">
      <button class="menu-toggle" @click="toggleMenu">
        <span class="material-symbols-outlined">double_arrow</span>
      </button>
    </div>

    <h3 class="menu-heading">Menú</h3>
    <div class="menu">
      <RouterLink to="/dashboard" class="menu-item" active-class="is-active">
        <span class="material-symbols-outlined">dashboard</span>
        <span class="text">Dashboard</span>
      </RouterLink>
      <RouterLink to="/sales/dashboard" class="menu-item" active-class="is-active">
        <span class="material-symbols-outlined">payments</span>
        <span class="text">Ventas</span>
      </RouterLink>
      <RouterLink to="/products" class="menu-item" active-class="is-active">
        <span class="material-symbols-outlined">inventory_2</span>
        <span class="text">Catálogo</span>
      </RouterLink>
      <RouterLink to="/parties" class="menu-item" active-class="is-active">
        <span class="material-symbols-outlined">groups</span>
        <span class="text">Entidades</span>
      </RouterLink>
      <RouterLink to="/mes/dashboard" class="menu-item" active-class="is-active">
        <span class="material-symbols-outlined">precision_manufacturing</span>
        <span class="text">Taller (MES)</span>
      </RouterLink>
    </div>

    <div class="flex"></div>

    <div class="menu">
       <RouterLink v-if="isAdmin" to="/admin/users" class="menu-item" active-class="is-active">
        <span class="material-symbols-outlined">admin_panel_settings</span>
        <span class="text">Admin</span>
      </RouterLink>
      <div class="menu-item">
         <UserMenu />
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserMenu from './UserMenu.vue'

const isExpanded = ref(localStorage.getItem('sidebar-expanded') === 'true')
const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)


const toggleMenu = () => {
  isExpanded.value = !isExpanded.value
  localStorage.setItem('sidebar-expanded', isExpanded.value)
}
</script>

<style scoped>
:root {
	--sidebar-bg-color: #1e293b;
	--sidebar-item-hover: #334155;
	--sidebar-item-active: #f4c430;
  --sidebar-text-color: #f1f5f9;
}

.side-navbar {
  display: flex;
  flex-direction: column;
  background-color: var(--sidebar-bg-color);
  color: var(--sidebar-text-color);
  width: calc(2rem + 32px);
  min-height: 100vh;
  overflow: hidden;
  padding: 1rem;
  transition: 0.2s ease-in-out;
}
.side-navbar.is-expanded {
  width: 280px;
}

.navbar-header {
  margin-bottom: 1rem;
}
.logo {
  display: flex;
  align-items: center;
  text-decoration: none;
}
.logo-icon {
  width: 2rem;
  height: 2rem;
}
.logo-text {
  color: var(--sidebar-text-color);
  font-size: 1.25rem;
  font-weight: bold;
  margin-left: 0.5rem;
  opacity: 0;
  transition: opacity 0.2s;
}
.is-expanded .logo-text {
  opacity: 1;
}

.menu-toggle-wrap {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
  position: relative;
  top: 0;
  transition: 0.2s ease-in-out;
}
.menu-toggle {
  transition: 0.2s ease-in-out;
  background: none;
  border: none;
  color: var(--sidebar-text-color);
  cursor: pointer;
}
.menu-toggle .material-symbols-outlined {
  font-size: 2rem;
  transition: 0.2s ease-out;
}
.menu-toggle:hover .material-symbols-outlined {
  color: var(--sidebar-item-active);
  transform: translateX(0.5rem);
}
.is-expanded .menu-toggle {
  transform: rotate(-180deg);
}

.menu-heading {
  color: #64748b;
  font-size: 0.75rem;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
  opacity: 0;
  transition: opacity 0.3s;
}
.is-expanded .menu-heading {
  opacity: 1;
}

.menu .menu-item {
  position: relative;
  display: flex;
  align-items: center;
  text-decoration: none;
  padding: 0.75rem 1rem;
  color: var(--sidebar-text-color);
  transition: 0.2s ease-in-out;
}
.menu-item .material-symbols-outlined {
  font-size: 2rem;
  margin-right: 1rem;
  transition: 0.2s ease-in-out;
}
.menu-item .text {
  opacity: 0;
  transition: opacity 0.3s;
}
.is-expanded .menu-item .text {
  opacity: 1;
}
.menu-item:hover, .menu-item.is-active {
  background-color: var(--sidebar-item-hover);
  border-right: 5px solid var(--sidebar-item-active);
}
.menu-item.is-active .material-symbols-outlined,
.menu-item.is-active .text {
  color: var(--sidebar-item-active);
}

.flex {
  flex: 1 1 0%;
}
</style>
