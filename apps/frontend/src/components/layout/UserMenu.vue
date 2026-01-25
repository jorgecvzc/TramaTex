<template>
  <div class="user-menu-container">
    <button class="user-menu-toggle" @click="isOpen = !isOpen">
      <span>{{ usuario?.nombre || 'Usuario' }}</span>
      <span class="arrow">▼</span>
    </button>

    <div v-if="isOpen" class="user-menu-dropdown">
      <div class="user-info">
        <div class="user-email">{{ usuario?.email }}</div>
        <div class="user-role">{{ usuario?.rol }}</div>
      </div>

      <hr class="menu-divider" />

      <button class="menu-item" @click="handleLogout">
        Cerrar Sesión
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables'

const router = useRouter()
const { usuario, logout } = useAuth()
const isOpen = ref(false)

async function handleLogout() {
  await logout()
  isOpen.value = false
  router.push('/login')
}
</script>

<style scoped>
.user-menu-container {
  position: relative;
}

.user-menu-toggle {
  background-color: rgba(255, 215, 0, 0.1);
  border: 1px solid rgba(255, 215, 0, 0.3);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  transition: background-color 0.2s;
}

.user-menu-toggle:hover {
  background-color: rgba(255, 215, 0, 0.2);
}

.arrow {
  font-size: 0.7rem;
  transition: transform 0.2s;
}

.user-menu-toggle[aria-expanded='true'] .arrow {
  transform: rotate(180deg);
}

.user-menu-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  min-width: 200px;
  margin-top: 0.5rem;
  z-index: 100;
  overflow: hidden;
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.user-info {
  padding: 1rem;
  color: #333;
}

.user-email {
  font-size: 0.9rem;
  font-weight: 500;
}

.user-role {
  font-size: 0.8rem;
  color: #666;
  text-transform: capitalize;
}

.menu-divider {
  border: none;
  border-top: 1px solid #eee;
  margin: 0;
}

.menu-item {
  width: 100%;
  padding: 0.75rem 1rem;
  border: none;
  background-color: transparent;
  color: #333;
  text-align: left;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.2s;
}

.menu-item:hover {
  background-color: #f5f5f5;
}

.menu-item:active {
  background-color: #efefef;
}
</style>
