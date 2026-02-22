<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <h1>Bienvenido, {{ usuario?.email || 'Usuario' }}</h1>
          <p class="subtitle">Panel de control del sistema TramaTex ERP/MES</p>
        </div>
      </header>

      <div class="areas-grid">
        <!-- Área de Entidades -->
        <div class="area-card">
          <div class="area-header">
            <div class="area-icon">
              <Users :size="40" />
            </div>
            <h2>Entidades</h2>
            <span class="area-badge">Disponible</span>
          </div>
          <p class="area-description">Gestiona clientes y proveedores</p>
          <div class="area-links">
            <RouterLink to="/parties" class="link-primary">
              <span class="link-icon">
                <Users :size="24" />
              </span>
              <span>Entidades</span>
            </RouterLink>
          </div>
        </div>

        <!-- Área de Productos -->
        <div class="area-card">
          <div class="area-header">
            <div class="area-icon">
              <Package :size="40" />
            </div>
            <h2>Productos</h2>
            <span class="area-badge">Disponible</span>
          </div>
          <p class="area-description">Gestiona el catálogo de productos, atributos, marcas y categorías</p>
          <div class="area-links">
            <RouterLink to="/products" class="link-primary">
              <span class="link-icon">
                <Package :size="24" />
              </span>
              <span>Catálogo de Productos</span>
            </RouterLink>
            <RouterLink to="/master-data/attributes" class="link-secondary">
              <span class="link-icon">
                <Zap :size="24" />
              </span>
              <span>Atributos</span>
            </RouterLink>
            <RouterLink to="/master-data/brands" class="link-secondary">
              <span class="link-icon">
                <Tag :size="24" />
              </span>
              <span>Marcas</span>
            </RouterLink>
            <RouterLink to="/master-data/product-groups" class="link-secondary">
              <span class="link-icon">
                <Folder :size="24" />
              </span>
              <span>Categorías</span>
            </RouterLink>
          </div>
        </div>

        <!-- Área de Ventas -->
        <div class="area-card">
          <div class="area-header">
            <div class="area-icon">
              <DollarSign :size="40" />
            </div>
            <h2>Ventas</h2>
            <span class="area-badge">Disponible</span>
          </div>
          <p class="area-description">Gestión de cotizaciones, pedidos y facturación</p>
          <div class="area-links">
            <RouterLink to="/sales/quotes" class="link-primary">
              <span class="link-icon">
                <Clipboard :size="24" />
              </span>
              <span>Presupuestos</span>
            </RouterLink>
            <RouterLink to="/sales/orders" class="link-secondary">
              <span class="link-icon">
                <ShoppingCart :size="24" />
              </span>
              <span>Pedidos</span>
            </RouterLink>
            <RouterLink to="/sales/delivery-notes" class="link-secondary">
              <span class="link-icon">
                <ScrollText :size="24" />
              </span>
              <span>Albaranes</span>
            </RouterLink>
            <RouterLink to="/sales/invoices" class="link-secondary">
              <span class="link-icon">
                <Receipt :size="24" />
              </span>
              <span>Facturas</span>
            </RouterLink>
            <RouterLink to="/sales/tickets/new" class="link-secondary">
              <span class="link-icon">
                <Ticket :size="24" />
              </span>
              <span>Nuevo Ticket</span>
            </RouterLink>
          </div>
        </div>

        <!-- Área de Producción/MES -->
        <div class="area-card">
          <div class="area-header">
            <div class="area-icon">
              <Factory :size="40" />
            </div>
            <h2>Producción</h2>
            <span class="area-badge">Disponible</span>
          </div>
          <p class="area-description">Manufacturing Execution System - Control de producción</p>
          <div class="area-links">
            <RouterLink to="/mes/dashboard" class="link-primary">
              <span class="link-icon">
                <Settings :size="24" />
              </span>
              <span>Dashboard MES</span>
            </RouterLink>
            <RouterLink to="/mes/works" class="link-secondary">
              <span class="link-icon">
                <Wrench :size="24" />
              </span>
              <span>Trabajos MES</span>
            </RouterLink>
            <RouterLink to="/mes/terminal" class="link-secondary">
              <span class="link-icon">
                <CheckCircle :size="24" />
              </span>
              <span>Terminal Tablet</span>
            </RouterLink>
          </div>
        </div>

        <!-- Área de Informes -->
        <div class="area-card area-disabled">
          <div class="area-header">
            <div class="area-icon">
              <BarChart :size="40" />
            </div>
            <h2>Informes</h2>
            <span class="area-badge disabled">Post-MVP</span>
          </div>
          <p class="area-description">Reportes, análisis y estadísticas del sistema</p>
          <div class="area-links">
            <div class="link-primary link-disabled">
              <span class="link-icon">
                <LineChart :size="24" />
              </span>
              <span>Dashboard Analítico</span>
            </div>
            <div class="link-secondary link-disabled">
              <span class="link-icon">
                <FileText :size="24" />
              </span>
              <span>Reportes</span>
            </div>
          </div>
        </div>

        <!-- Área de Administración -->
        <div v-if="isAdmin" class="area-card">
          <div class="area-header">
            <div class="area-icon">
              <UserCog :size="40" />
            </div>
            <h2>Administración</h2>
            <span class="area-badge admin">Admin</span>
          </div>
          <p class="area-description">Gestión de usuarios, roles y configuración del sistema</p>
          <div class="area-links">
            <RouterLink to="/admin/users" class="link-primary">
              <span class="link-icon">
                <User :size="24" />
              </span>
              <span>Gestión de Usuarios</span>
            </RouterLink>
            <RouterLink to="/admin/print-profile" class="link-secondary">
              <span class="link-icon">
                <FileText :size="24" />
              </span>
              <span>Perfil Fiscal de Impresión</span>
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { useAuth } from '@/composables'
import { useAuthStore } from '@/stores/auth'
import { Users, Package, DollarSign, Factory, BarChart, UserCog, Zap, Tag, Folder, Clipboard, ShoppingCart, ScrollText, Receipt, Ticket, Settings, Wrench, CheckCircle, LineChart, FileText, User } from 'lucide-vue-next';

const { usuario } = useAuth()
const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
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
}

.page-header {
  margin-bottom: 2.5rem;
}

.page-header h1 {
  color: #1b3a6b;
  font-size: 2rem;
  margin: 0 0 0.5rem;
}

.subtitle {
  color: #64748b;
  font-size: 1rem;
  margin: 0;
}

/* Areas Grid */
.areas-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;
}

/* Area Card */
.area-card {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 2px solid #e2e8f0;
  transition: all 0.3s ease;
}

.area-card:not(.area-disabled):hover {
  border-color: #f4c430;
  box-shadow: 0 8px 20px rgba(244, 196, 48, 0.15);
  transform: translateY(-2px);
}

.area-disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Area Header */
.area-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.75rem;
}

.area-icon :deep(svg) {
  width: 2.5rem; /* Equivalent to 40px */
  height: 2.5rem; /* Equivalent to 40px */
}

.area-header h2 {
  color: #1b3a6b;
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
  flex: 1;
}

.area-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 0.875rem;
  font-size: 0.7rem;
  font-weight: 700;
  border-radius: 16px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.area-badge:not(.disabled):not(.admin) {
  background: linear-gradient(135deg, #dcfce7 0%, #bbf7d0 100%);
  color: #166534;
}

.area-badge.disabled {
  background: #f1f5f9;
  color: #94a3b8;
}

.area-badge.admin {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  color: #78350f;
}

/* Area Description */
.area-description {
  color: #64748b;
  font-size: 0.95rem;
  line-height: 1.6;
  margin: 0 0 1.5rem 0;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
}

/* Area Links */
.area-links {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

/* Primary Link (más grande y destacado) */
.link-primary {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 1rem 1.25rem;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  font-size: 1.05rem;
  color: #1b3a6b;
  background: linear-gradient(135deg, #f4c430 0%, #f4d03f 100%);
  transition: all 0.2s ease;
  border: 2px solid transparent;
}

.link-primary:not(.link-disabled):hover {
  background: linear-gradient(135deg, #f4d03f 0%, #fde68a 100%);
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(244, 196, 48, 0.3);
}

.link-primary .link-icon :deep(svg) {
  width: 1.5rem; /* Equivalent to 24px */
  height: 1.5rem; /* Equivalent to 24px */
}

/* Secondary Links (más pequeños) */
.link-secondary {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.95rem;
  color: #475569;
  background: #f8fafc;
  transition: all 0.2s ease;
  border: 1px solid #e2e8f0;
}

.link-secondary:not(.link-disabled):hover {
  background: #f1f5f9;
  color: #1b3a6b;
  border-color: #cbd5e1;
  transform: translateX(4px);
}

.link-secondary .link-icon :deep(svg) {
  width: 1.25rem; /* Equivalent to 20px */
  height: 1.25rem; /* Equivalent to 20px */
}

/* Disabled Links */
.link-disabled {
  cursor: not-allowed;
  opacity: 0.4;
  background: #f1f5f9 !important;
}

.link-disabled:hover {
  transform: none !important;
  box-shadow: none !important;
}

/* Responsive */
@media (max-width: 768px) {
  .areas-grid {
    grid-template-columns: 1fr;
  }
  
  .area-card {
    padding: 1.5rem;
  }
  
  .area-icon :deep(svg) {
    width: 2rem;
    height: 2rem;
  }
  
  .area-header h2 {
    font-size: 1.25rem;
  }
}
</style>
