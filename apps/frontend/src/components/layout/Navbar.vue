<template>
  <nav class="navbar">
    <div class="navbar-container">
      <div class="navbar-brand">
        <RouterLink to="/dashboard" class="logo">TramaTex</RouterLink>
      </div>

      <ul class="nav-menu">
        <li>
          <RouterLink to="/orders" class="nav-link" active-class="active">
            Pedidos
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/inventory" class="nav-link" active-class="active">
            Inventario
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/customers" class="nav-link" active-class="active">
            Clientes
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/parties" class="nav-link" active-class="active">
            Entidades
          </RouterLink>
        </li>
        <li>
          <RouterLink to="/catalogos" class="nav-link" active-class="active">
            Catálogos
          </RouterLink>
        </li>
        <li v-if="isAdmin">
          <RouterLink to="/admin/users" class="nav-link" active-class="active">
            Usuarios
          </RouterLink>
        </li>
      </ul>

      <UserMenu />
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserMenu from './UserMenu.vue'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
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
  transition: color 0.2s;
  padding: 0.5rem 1rem;
  border-radius: 4px;
}

.nav-link:hover {
  color: #ffd700;
}

.nav-link.active {
  background-color: rgba(255, 215, 0, 0.2);
  color: #ffd700;
}

@media (max-width: 768px) {
  .nav-menu {
    gap: 1rem;
  }

  .nav-link {
    padding: 0.25rem 0.5rem;
    font-size: 0.9rem;
  }
}
</style>
