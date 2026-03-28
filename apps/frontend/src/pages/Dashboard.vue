<template>
  <Navbar />
  
  <div class="dashboard-page">
    <div class="dashboard-container">
      <header class="dashboard-header mb-10">
        <div class="welcome-section">
          <h1 class="font-brand">¡Hola de nuevo, {{ userName }}!</h1>
          <p class="text-muted">Aquí tienes el resumen operativo de TramaTex para hoy.</p>
        </div>
        <div class="header-date">
          <span class="material-symbols-outlined">calendar_today</span>
          <span>{{ todayDate }}</span>
        </div>
      </header>

      <!-- SECCIÓN: KPIs GLOBALES -->
      <section class="stats-grid mb-10">
        <div class="stat-card blue">
          <div class="stat-icon"><span class="material-symbols-outlined">payments</span></div>
          <div class="stat-data">
            <span class="stat-value">{{ salesStats.monthlyTotal }}</span>
            <span class="stat-label">Ventas del Mes</span>
          </div>
        </div>
        <div class="stat-card green">
          <div class="stat-icon"><span class="material-symbols-outlined">shopping_cart</span></div>
          <div class="stat-data">
            <span class="stat-value">{{ salesStats.pendingOrders }}</span>
            <span class="stat-label">Pedidos Pendientes</span>
          </div>
        </div>
        <div class="stat-card yellow">
          <div class="stat-icon"><span class="material-symbols-outlined">precision_manufacturing</span></div>
          <div class="stat-data">
            <span class="stat-value">{{ mesStats.activeWorkOrders }}</span>
            <span class="stat-label">Órdenes en Taller</span>
          </div>
        </div>
        <div class="stat-card purple">
          <div class="stat-icon"><span class="material-symbols-outlined">groups</span></div>
          <div class="stat-data">
            <span class="stat-value">{{ partyStats.totalParties }}</span>
            <span class="stat-label">Entidades</span>
          </div>
        </div>
      </section>

      <!-- SECCIÓN: ACCESOS DIRECTOS OPERATIVOS -->
      <section class="dashboard-section mb-10">
        <div class="section-header">
          <span class="material-symbols-outlined">bolt</span>
          <h2>Operaciones Rápidas</h2>
        </div>
        <div class="actions-grid">
          <RouterLink to="/sales/orders/new" class="action-card">
            <span class="material-symbols-outlined">add_shopping_cart</span>
            <span>Nuevo Pedido</span>
          </RouterLink>
          <RouterLink to="/sales/quotes/new" class="action-card">
            <span class="material-symbols-outlined">add_notes</span>
            <span>Presupuesto</span>
          </RouterLink>
          <RouterLink to="/sales/tickets/new" class="action-card highlight">
            <span class="material-symbols-outlined">point_of_sale</span>
            <span>Venta Directa</span>
          </RouterLink>
          <RouterLink to="/products/new" class="action-card">
            <span class="material-symbols-outlined">add_box</span>
            <span>Nuevo Producto</span>
          </RouterLink>
        </div>
      </section>

      <!-- SECCIÓN: DASHBOARDS DE MÓDULO -->
      <section class="dashboard-section mb-10">
        <div class="section-header">
          <span class="material-symbols-outlined">monitoring</span>
          <h2>Paneles de Control de Módulo</h2>
        </div>
        <div class="module-dashboards-grid">
          <RouterLink to="/sales/dashboard" class="module-card">
            <div class="module-icon blue"><span class="material-symbols-outlined">payments</span></div>
            <div class="module-info">
              <strong>Dashboard de Ventas</strong>
              <p>KPIs de facturación, pedidos y entregas.</p>
            </div>
            <span class="material-symbols-outlined arrow">chevron_right</span>
          </RouterLink>

          <RouterLink to="/products/dashboard" class="module-card">
            <div class="module-icon yellow"><span class="material-symbols-outlined">inventory_2</span></div>
            <div class="module-info">
              <strong>Gestión de Catálogo</strong>
              <p>Métricas de productos, marcas y atributos.</p>
            </div>
            <span class="material-symbols-outlined arrow">chevron_right</span>
          </RouterLink>

          <RouterLink to="/parties/dashboard" class="module-card">
            <div class="module-icon green"><span class="material-symbols-outlined">groups</span></div>
            <div class="module-info">
              <strong>Control de Entidades</strong>
              <p>Visión general de clientes y proveedores.</p>
            </div>
            <span class="material-symbols-outlined arrow">chevron_right</span>
          </RouterLink>

          <RouterLink to="/mes/dashboard" class="module-card">
            <div class="module-icon purple"><span class="material-symbols-outlined">precision_manufacturing</span></div>
            <div class="module-info">
              <strong>Monitor de Taller (MES)</strong>
              <p>Estado de producción y órdenes activas.</p>
            </div>
            <span class="material-symbols-outlined arrow">chevron_right</span>
          </RouterLink>
        </div>
      </section>

      <div class="dashboard-grid-main">
        <!-- Listados de actividad -->
        <section class="dashboard-section main-area">
          <div class="section-header">
            <span class="material-symbols-outlined">history</span>
            <h2>Actividad Reciente</h2>
          </div>
          <div v-if="isLoading" class="loading-state">
            <div class="spinner"></div>
            <p>Sincronizando datos operativos...</p>
          </div>
          <div v-else class="activity-content">
            <p class="text-muted italic">Últimos movimientos registrados en el sistema.</p>
            <!-- Aquí iría un feed de actividad en el futuro -->
          </div>
        </section>

        <!-- Sidebar Administrativo -->
        <aside class="dashboard-sidebar">
          <section class="dashboard-section">
            <div class="section-header">
              <span class="material-symbols-outlined">admin_panel_settings</span>
              <h2>Administración</h2>
            </div>
            <div class="admin-links">
              <RouterLink to="/admin/users" class="admin-card">
                <span class="material-symbols-outlined">manage_accounts</span>
                <div class="admin-card-info">
                  <strong>Gestión de Usuarios</strong>
                  <p>Roles y accesos</p>
                </div>
                <span class="material-symbols-outlined arrow">chevron_right</span>
              </RouterLink>
              <RouterLink to="/admin/print-profile" class="admin-card mt-3">
                <span class="material-symbols-outlined">receipt_long</span>
                <div class="admin-card-info">
                  <strong>Perfil de Impresión</strong>
                  <p>Datos fiscales para PDF</p>
                </div>
                <span class="material-symbols-outlined arrow">chevron_right</span>
              </RouterLink>
            </div>
          </section>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { useAuthStore } from '@/stores/auth'
import salesApi from '@/services/salesApi'
import { mesApi } from '@/services/mesApi'
import { partyApi } from '@/services/partyApi'

const authStore = useAuthStore()
const userName = computed(() => authStore.user?.email.split('@')[0] || 'Usuario')
const todayDate = computed(() => new Date().toLocaleDateString('es-ES', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }))

const isLoading = ref(true)
const salesStats = ref({ monthlyTotal: '0,00 €', pendingOrders: 0 })
const mesStats = ref({ activeWorkOrders: 0 })
const partyStats = ref({ totalParties: 0 })

async function loadStats() {
  isLoading.value = true
  try {
    const [orders, workOrders, parties] = await Promise.all([
      salesApi.listOrders({ limit: 1 }),
      mesApi.listWorkOrders({ status: 'IN_PROGRESS' }),
      partyApi.listParties({ limit: 1 })
    ])
    
    salesStats.value.pendingOrders = orders.total || 0
    mesStats.value.activeWorkOrders = Array.isArray(workOrders) ? workOrders.length : 0
    partyStats.value.totalParties = parties.total || 0
  } catch (err) {
    console.error('Error loading dashboard stats:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(loadStats)
</script>

<style scoped>
.dashboard-page { background-color: var(--color-background); min-height: 100vh; padding-bottom: 4rem; }
.dashboard-container { max-width: 1400px; margin: 0 auto; padding: 2rem; }

.dashboard-header { display: flex; justify-content: space-between; align-items: flex-start; }
.dashboard-header h1 { font-size: 2rem; color: var(--color-secondary); margin: 0; }
.header-date { display: flex; align-items: center; gap: 0.5rem; color: var(--color-text-secondary); font-size: 0.9rem; font-weight: 600; padding: 0.5rem 1rem; background: white; border-radius: 100px; box-shadow: var(--box-shadow-sm); border: 1px solid var(--color-border); }

/* KPI Cards */
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 1.5rem; }
.stat-card { background: white; padding: 1.5rem; border-radius: 16px; display: flex; align-items: center; gap: 1.25rem; box-shadow: var(--box-shadow-sm); border: 1px solid var(--color-border); }
.stat-icon { width: 56px; height: 56px; border-radius: 12px; display: flex; align-items: center; justify-content: center; }
.stat-icon .material-symbols-outlined { font-size: 32px; }

.stat-card.blue .stat-icon { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.stat-card.green .stat-icon { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.stat-card.yellow .stat-icon { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.stat-card.purple .stat-icon { background: rgba(168, 85, 247, 0.1); color: #9333ea; }

.stat-data { display: flex; flex-direction: column; }
.stat-value { font-size: 1.75rem; font-weight: 800; color: var(--color-text-primary); }
.stat-label { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }

/* Dashboard Sections */
.dashboard-section { background: white; padding: 2rem; border-radius: 16px; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.section-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; padding-bottom: 1rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 1rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-primary); letter-spacing: 0.05em; }
.section-header .material-symbols-outlined { color: var(--color-primary); font-size: 24px; }

/* Actions Grid */
.actions-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
.action-card { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; padding: 2rem; background: var(--color-background); border-radius: 12px; text-decoration: none; color: var(--color-text-primary); font-weight: 700; transition: all 0.2s ease; border: 1px solid transparent; }
.action-card:hover { transform: translateY(-4px); background: white; border-color: var(--color-primary); box-shadow: var(--box-shadow-md); color: var(--color-primary); }
.action-card .material-symbols-outlined { font-size: 36px; color: var(--color-primary); }
.action-card.highlight { background: rgba(230, 184, 0, 0.1); border: 1px dashed var(--color-primary); }

/* Module Dashboards Grid */
.module-dashboards-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; }
.module-card { display: flex; align-items: center; gap: 1.25rem; padding: 1.25rem; background: var(--color-background); border-radius: 12px; text-decoration: none; transition: all 0.2s; border: 1px solid transparent; position: relative; }
.module-card:hover { background: white; border-color: var(--color-secondary); transform: translateX(4px); box-shadow: var(--box-shadow-md); }
.module-icon { width: 48px; height: 48px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.module-icon .material-symbols-outlined { font-size: 24px; }
.module-icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.module-icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.module-icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.module-icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }

.module-info { flex: 1; }
.module-info strong { display: block; font-size: 0.95rem; color: var(--color-text-primary); }
.module-info p { font-size: 0.75rem; color: var(--color-text-secondary); margin: 0.25rem 0 0; }
.module-card .arrow { color: var(--color-border); transition: 0.2s; }
.module-card:hover .arrow { color: var(--color-secondary); transform: translateX(3px); }

/* Main Grid Layout */
.dashboard-grid-main { display: grid; grid-template-columns: 1fr 380px; gap: 1.5rem; align-items: start; }
.admin-card { display: flex; align-items: center; gap: 1rem; padding: 1rem; background: var(--color-background); border-radius: 10px; border: 1px solid transparent; text-decoration: none; color: var(--color-text-primary); transition: 0.2s; }
.admin-card:hover { background: white; border-color: var(--color-primary); transform: translateX(4px); box-shadow: var(--box-shadow-sm); }
.admin-card-info strong { font-size: 0.85rem; display: block; }
.admin-card-info p { font-size: 0.7rem; color: var(--color-text-secondary); margin: 0; }

.help-notice { padding: 1.25rem; background: rgba(59, 130, 246, 0.05); border-radius: 12px; border: 1px dashed rgba(59, 130, 246, 0.3); }
.notice-header { display: flex; align-items: center; gap: 0.5rem; color: #2563eb; font-size: 0.85rem; font-weight: 700; text-transform: uppercase; }
.help-text { font-size: 0.8rem; color: var(--color-text-secondary); margin-top: 0.5rem; line-height: 1.5; }

.spinner { width: 40px; height: 40px; margin: 2rem auto; border: 3px solid #f3f4f6; border-top-color: var(--color-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 1024px) {
  .dashboard-grid-main { grid-template-columns: 1fr; }
  .module-dashboards-grid { grid-template-columns: 1fr; }
}
</style>
