<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Operaciones / Entidades</p>
          <h1>Crear entidad</h1>
          <p class="subtitle">Registra una nueva entidad con sus datos principales.</p>
        </div>
        <RouterLink to="/parties" class="btn btn-secondary">
          Volver al listado
        </RouterLink>
      </header>

      <section class="card">
        <PartyForm @submit="handleSubmit" />
      </section>
    </div>
  </div>
</template>

<script setup>
import { useRouter, RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import PartyForm from '@/components/party/PartyForm.vue'

const router = useRouter()

async function handleSubmit(party) {
  await router.push(`/parties/${party.id}`)
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.page-header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.page-header h1 {
  color: #1b3a6b;
  margin: 0.25rem 0 0;
}

.breadcrumb {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin: 0;
}

.subtitle {
  color: #64748b;
  margin: 0.5rem 0 0;
  font-size: 0.95rem;
}

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
  max-width: 720px;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover {
  background: #f8fafc;
}

@media (max-width: 768px) {
  .dashboard-content {
    padding: 1.5rem;
  }

  .card {
    max-width: 100%;
  }
}
</style>
