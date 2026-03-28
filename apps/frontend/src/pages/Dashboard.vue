<template>
  <Navbar />
  
  <div class="dashboard-page">
    <div class="dashboard-container">
      <!-- CABECERA: Bienvenida y Fecha -->
      <header class="dashboard-header mb-8">
        <div class="welcome-text">
          <h1 class="font-brand">Panel de Control</h1>
          <p class="text-muted">Bienvenido, <strong>{{ userName }}</strong>. Gestión de TramaTex para el {{ todayDate }}.</p>
        </div>
        <div class="header-actions">
          <button @click="loadStats" class="btn-refresh" :disabled="isLoading" title="Actualizar datos">
            <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">sync</span>
          </button>
        </div>
      </header>

      <!-- CAPA 1: KPIs INTERACTIVOS -->
      <section class="kpi-grid mb-12">
        <div class="kpi-card clickable" @click="navigateTo('/sales/dashboard')">
          <div class="kpi-icon blue"><span class="material-symbols-outlined">payments</span></div>
          <div class="kpi-data">
            <label>Ventas Mes</label>
            <strong>{{ salesStats.monthlyTotal }}</strong>
          </div>
        </div>
        <div class="kpi-card clickable" @click="navigateTo('/sales/orders?status=PENDIENTE')">
          <div class="kpi-icon green"><span class="material-symbols-outlined">shopping_cart</span></div>
          <div class="kpi-data">
            <label>Pedidos Pendientes</label>
            <strong>{{ salesStats.pendingOrders }}</strong>
          </div>
        </div>
        <div class="kpi-card clickable" @click="navigateTo('/mes/work-orders?status=IN_PROGRESS')">
          <div class="kpi-icon yellow"><span class="material-symbols-outlined">precision_manufacturing</span></div>
          <div class="kpi-data">
            <label>Órdenes en Taller</label>
            <strong>{{ mesStats.activeWorkOrders }}</strong>
          </div>
        </div>
        <div class="kpi-card clickable" @click="navigateTo('/parties')">
          <div class="kpi-icon purple"><span class="material-symbols-outlined">groups</span></div>
          <div class="kpi-data">
            <label>Entidades</label>
            <strong>{{ partyStats.totalParties }}</strong>
          </div>
        </div>
      </section>

      <div class="dashboard-main-layout">
        <!-- COLUMNA PRINCIPAL: OPERACIONES Y MÓDULOS -->
        <main class="main-column">
          <!-- Gestión por Módulos -->
          <section class="modules-section">
            <div class="section-title-alt">
              <span class="material-symbols-outlined">grid_view</span>
              <h2>Gestión de Módulos</h2>
            </div>
            <div class="modules-grid">
              <RouterLink to="/sales/dashboard" class="module-link-card">
                <div class="m-icon blue"><span class="material-symbols-outlined">account_balance_wallet</span></div>
                <div class="m-info">
                  <strong>Ventas y Facturación</strong>
                  <p>Presupuestos, pedidos y albaranes.</p>
                </div>
                <span class="material-symbols-outlined arrow">chevron_right</span>
              </RouterLink>
              <RouterLink to="/products/dashboard" class="module-link-card">
                <div class="m-icon yellow"><span class="material-symbols-outlined">inventory_2</span></div>
                <div class="m-info">
                  <strong>Catálogo y Almacén</strong>
                  <p>Variantes JIT, marcas y atributos.</p>
                </div>
                <span class="material-symbols-outlined arrow">chevron_right</span>
              </RouterLink>
              <RouterLink to="/parties/dashboard" class="module-link-card">
                <div class="m-icon green"><span class="material-symbols-outlined">groups</span></div>
                <div class="m-info">
                  <strong>Terceros y CRM</strong>
                  <p>Base de datos de clientes y proveedores.</p>
                </div>
                <span class="material-symbols-outlined arrow">chevron_right</span>
              </RouterLink>
              <RouterLink to="/mes/dashboard" class="module-link-card">
                <div class="m-icon purple"><span class="material-symbols-outlined">precision_manufacturing</span></div>
                <div class="m-info">
                  <strong>Monitor de Producción</strong>
                  <p>Estado de taller y órdenes MES.</p>
                </div>
                <span class="material-symbols-outlined arrow">chevron_right</span>
              </RouterLink>
            </div>
          </section>
        </main>

        <!-- COLUMNA LATERAL: ACCIONES Y ADMIN -->
        <aside class="side-column">
          <!-- Centro de Operaciones -->
          <section class="ops-section card mb-8">
            <div class="ops-header">
              <span class="material-symbols-outlined">rocket_launch</span>
              <h2>Accesos Directos</h2>
            </div>
            <div class="ops-grid">
              <RouterLink to="/sales/orders/new" class="op-item">
                <div class="op-icon"><span class="material-symbols-outlined">add_shopping_cart</span></div>
                <span>Nuevo Pedido</span>
              </RouterLink>
              <RouterLink to="/sales/tickets/new" class="op-item highlight">
                <div class="op-icon"><span class="material-symbols-outlined">point_of_sale</span></div>
                <span>Venta Directa</span>
              </RouterLink>
              <RouterLink to="/products/new" class="op-item">
                <div class="op-icon"><span class="material-symbols-outlined">add_box</span></div>
                <span>Nuevo Producto</span>
              </RouterLink>
              <RouterLink to="/parties/new" class="op-item">
                <div class="op-icon"><span class="material-symbols-outlined">person_add</span></div>
                <span>Nueva Entidad</span>
              </RouterLink>
            </div>
          </section>

          <section v-if="isAdmin" class="card admin-side-card">
            <div class="ops-header">
              <span class="material-symbols-outlined">admin_panel_settings</span>
              <h2>Sistema</h2>
            </div>
            <div class="side-links">
              <RouterLink to="/admin/users" class="side-link-item">
                <span class="material-symbols-outlined">manage_accounts</span>
                <span>Gestión de Usuarios</span>
              </RouterLink>
              <RouterLink to="/admin/print-profile" class="side-link-item mt-2">
                <span class="material-symbols-outlined">receipt_long</span>
                <span>Perfil de Impresión</span>
              </RouterLink>
              <RouterLink to="/dev/design-system" class="side-link-item mt-2 dev-link">
                <span class="material-symbols-outlined">architecture</span>
                <span>Design System</span>
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

const router = useRouter()
const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
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
      salesApi.listOrders({ status: 'PENDIENTE' }),
      mesApi.listWorkOrders({ status: 'IN_PROGRESS' }),
      partyApi.listParties({ limit: 1 })
    ])
    salesStats.value.pendingOrders = orders.total || 0
    mesStats.value.activeWorkOrders = Array.isArray(workOrders) ? workOrders.length : 0
    partyStats.value.totalParties = parties.total || 0
  } catch (err) {
    console.error('Error dashboard:', err)
  } finally {
    isLoading.value = false
  }
}

function navigateTo(path: string) {
  router.push(path)
}

onMounted(loadStats)
</script>

<style scoped>
.dashboard-page { background-color: var(--color-background); min-height: 100vh; }
.dashboard-container { max-width: 1400px; margin: 0 auto; padding: 2rem; }

.dashboard-header { display: flex; justify-content: space-between; align-items: flex-end; }
.dashboard-header h1 { font-size: 2.25rem; color: var(--color-text-primary); margin: 0; font-weight: 800; }
.btn-refresh { background: white; border: 1px solid var(--color-border); padding: 0.5rem; border-radius: 8px; cursor: pointer; color: var(--color-text-secondary); transition: 0.2s; }
.btn-refresh:hover { color: var(--color-primary); border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }

.kpi-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1.25rem; }
.kpi-card { background: white; padding: 1.25rem; border-radius: 12px; display: flex; align-items: center; gap: 1rem; border: 1px solid var(--color-border); }
.kpi-card.clickable { cursor: pointer; transition: 0.2s; }
.kpi-card.clickable:hover { transform: translateY(-3px); border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }
.kpi-icon { width: 48px; height: 48px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.kpi-icon .material-symbols-outlined { font-size: 24px; }
.kpi-icon.blue { background: rgba(37, 99, 235, 0.08); color: #2563eb; }
.kpi-icon.green { background: rgba(22, 163, 74, 0.08); color: #16a34a; }
.kpi-icon.yellow { background: rgba(217, 119, 6, 0.08); color: #d97706; }
.kpi-icon.purple { background: rgba(147, 51, 234, 0.08); color: #9333ea; }
.kpi-data { display: flex; flex-direction: column; }
.kpi-data label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }
.kpi-data strong { font-size: 1.4rem; color: var(--color-text-primary); font-weight: 800; }

.dashboard-main-layout { display: grid; grid-template-columns: 1fr 380px; gap: 2rem; align-items: start; }

.ops-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; }
.ops-header h2 { font-size: 0.9rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-primary); }
.ops-header .material-symbols-outlined { color: var(--color-primary); }
.ops-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 1rem; }
.op-item { display: flex; flex-direction: column; align-items: center; padding: 1.5rem; background: var(--color-background); border-radius: 12px; text-decoration: none; color: var(--color-text-primary); font-weight: 700; font-size: 0.85rem; transition: 0.2s; text-align: center; gap: 0.75rem; border: 1px solid transparent; }
.op-item:hover { transform: translateY(-3px); background: white; border-color: var(--color-primary); box-shadow: var(--box-shadow-md); color: var(--color-primary); }
.op-icon { color: var(--color-primary); }
.op-icon .material-symbols-outlined { font-size: 32px; }
.op-item.highlight { background: rgba(230, 184, 0, 0.05); border: 1px dashed var(--color-primary); }

.section-title-alt { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.25rem; }
.section-title-alt h2 { font-size: 0.9rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-secondary); }
.modules-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; }
.module-link-card { display: flex; align-items: center; gap: 1.25rem; padding: 1.25rem; background: white; border-radius: 12px; border: 1px solid var(--color-border); text-decoration: none; transition: 0.2s; }
.module-link-card:hover { transform: translateX(5px); border-color: var(--color-secondary); box-shadow: var(--box-shadow-sm); }
.m-icon { width: 44px; height: 44px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.m-icon .material-symbols-outlined { font-size: 22px; }
.m-icon.blue { background: rgba(37, 99, 235, 0.05); color: #2563eb; }
.m-icon.yellow { background: rgba(217, 119, 6, 0.05); color: #d97706; }
.m-icon.green { background: rgba(22, 163, 74, 0.05); color: #16a34a; }
.m-icon.purple { background: rgba(147, 51, 234, 0.05); color: #9333ea; }
.m-info strong { display: block; font-size: 0.9rem; color: var(--color-text-primary); }
.m-info p { font-size: 0.75rem; color: var(--color-text-secondary); margin: 0.1rem 0 0; }
.module-link-card .arrow { margin-left: auto; color: var(--color-border); font-size: 18px; }

.side-link-item { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: var(--color-background); border-radius: 8px; text-decoration: none; color: var(--color-text-primary); font-size: 0.85rem; font-weight: 600; border: 1px solid transparent; transition: 0.2s; }
.side-link-item:hover { background: white; border-color: var(--color-primary); color: var(--color-primary); transform: translateX(3px); }
.side-link-item .material-symbols-outlined { font-size: 20px; color: var(--color-text-secondary); }
.dev-link { margin-top: 1.5rem; opacity: 0.7; }

.notice-card { background: white; border: 1px solid var(--color-border); border-left: 4px solid var(--color-primary); }
.notice-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.75rem; color: var(--color-primary); font-weight: 800; font-size: 0.8rem; text-transform: uppercase; }
.notice-card p { font-size: 0.8rem; line-height: 1.5; margin: 0; }

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.mb-12 { margin-bottom: 3.5rem; }

@media (max-width: 1024px) {
  .dashboard-main-layout { grid-template-columns: 1fr; }
  .modules-grid { grid-template-columns: 1fr; }
}
</style>
