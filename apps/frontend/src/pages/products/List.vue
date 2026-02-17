<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Catálogo / Productos</p>
          <h1>Catálogo de Productos</h1>
          <p class="subtitle">Gestión de productos, variantes y precios.</p>
        </div>
        <RouterLink to="/products/new" class="btn btn-primary">
          + Nuevo Producto
        </RouterLink>
      </header>

      <section class="card card-full">
        <ProductList />
      </section>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import ProductList from '@/components/product/ProductList.vue'
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 1400px;
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
  font-size: 2rem;
  font-weight: 700;
}

.breadcrumb {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin: 0;
  font-weight: 600;
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
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
}

.card-full {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.65rem 1.2rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.05s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
}

.btn-primary {
  background: #f4d03f;
  color: #1e293b;
  box-shadow: 0 2px 4px rgba(244, 208, 63, 0.3);
}

.btn-primary:hover {
  background: #e6c230;
  box-shadow: 0 4px 8px rgba(244, 208, 63, 0.4);
}

.btn-primary:active {
  transform: translateY(1px);
}

@media (max-width: 768px) {
  .dashboard-content {
    padding: 1.5rem;
  }

  .page-header h1 {
    font-size: 1.5rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .btn {
    width: 100%;
  }
}
</style>
