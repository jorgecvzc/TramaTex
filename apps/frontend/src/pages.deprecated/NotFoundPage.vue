<template>
  <div class="not-found">
    <div class="container">
      <h1>404</h1>
      <p>Página no encontrada</p>
      <RouterLink to="/dashboard" class="btn">Ir al Dashboard</RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1b3a6b 0%, #2c5aa0 100%);
}

.container {
  text-align: center;
  color: white;
}

.container h1 {
  font-size: 6rem;
  margin: 0;
  font-weight: bold;
}

.container p {
  font-size: 1.5rem;
  margin: 1rem 0 2rem;
}

.btn {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  background-color: #ffd700;
  color: #1b3a6b;
  text-decoration: none;
  border-radius: 4px;
  font-weight: 500;
  transition: background-color 0.2s;
}

.btn:hover {
  background-color: #ffed4e;
}
</style>
