<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <h1 class="not-found-title">404</h1>
      <h2 class="not-found-subtitle">Página no encontrada</h2>
      <p class="not-found-description">
        Lo sentimos, la página que buscas no existe o ha sido movida.
      </p>
      <RouterLink to="/dashboard" class="btn-home">
        <ArrowLeft :size="24" />
        <span>Ir al Dashboard</span>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { ArrowLeft } from 'lucide-vue-next';
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1b3a6b 0%, #2c5aa0 100%);
  padding: 2rem;
}

.not-found-content {
  text-align: center;
  max-width: 600px;
}

.not-found-title {
  font-size: 8rem;
  font-weight: 900;
  color: #f4c430;
  margin: 0 0 1rem;
  line-height: 1;
  text-shadow: 0 4px 12px rgba(244, 196, 48, 0.3);
}

.not-found-subtitle {
  font-size: 2rem;
  font-weight: 700;
  color: white;
  margin: 0 0 1rem;
}

.not-found-description {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.85);
  margin: 0 0 2.5rem;
  line-height: 1.6;
}

.btn-home {
  display: inline-flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 2rem;
  background: #f4c430;
  color: #1b3a6b;
  font-size: 1.125rem;
  font-weight: 700;
  text-decoration: none;
  border-radius: 12px;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(244, 196, 48, 0.3);
}

.btn-home:hover {
  background: #fde68a;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(244, 196, 48, 0.4);
}

.btn-icon {
  font-size: 1.5rem;
}

@media (max-width: 768px) {
  .not-found-title {
    font-size: 5rem;
  }
  
  .not-found-subtitle {
    font-size: 1.5rem;
  }
}
</style>
