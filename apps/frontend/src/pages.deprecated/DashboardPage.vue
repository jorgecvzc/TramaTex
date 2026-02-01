<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <h1>Bienvenido, {{ usuario?.email || 'Usuario' }}</h1>
      <div class="cards-grid">
        <div class="card">
          <h2>Pedidos</h2>
          <p>Gestiona tus pedidos</p>
        </div>
        <div class="card">
          <h2>Inventario</h2>
          <p>Controla tu inventario</p>
        </div>
        <div class="card">
          <h2>Clientes</h2>
          <p>Gestiona tus clientes</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Navbar from '@/components/layout/Navbar.vue'
import { useAuth } from '@/composables'

const { usuario } = useAuth()
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.dashboard-content h1 {
  color: #1b3a6b;
  margin-bottom: 2rem;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.card {
  background-color: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.card h2 {
  color: #1b3a6b;
  margin-bottom: 0.5rem;
  font-size: 1.2rem;
}

.card p {
  color: #666;
  font-size: 0.9rem;
}
</style>
