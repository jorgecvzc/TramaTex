<template>
  <div class="dashboard-page">
    <!-- CABECERA: Bienvenida y Fecha -->
    <header class="dashboard-header">
      <div class="welcome-text">
        <h1 class="font-brand">Panel de Control</h1>
        <p class="text-muted">Bienvenido, <strong>{{ userName }}</strong>. Gestión de TramaTex para el {{ todayDate }}.</p>
        <div class="context-tags">
          <span class="context-tag">Vista Ejecutiva</span>
          <span class="context-tag">Operaciones</span>
          <span class="context-tag" :class="{ loading: isLoading }">Datos en tiempo real</span>
        </div>
      </div>
      <div class="header-actions">
        <button @click="loadStats" class="btn-refresh" :disabled="isLoading" title="Actualizar datos">
          <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">sync</span>
        </button>
      </div>
    </header>

    <!-- KPIs INTERACTIVOS -->
    <section class="kpi-grid section-spacing">
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
    
    <div class="dashboard-main-layout section-spacing">
      <!-- COLUMNA PRINCIPAL -->
      <main class="main-column">
        <!-- ACCESOS DIRECTOS -->
        <section class="ops-section section-block">
          <div class="section-title-alt">
            <span class="material-symbols-outlined">rocket_launch</span>
            <h2>Accesos Directos</h2>
            <span class="section-tag">Atajos de operación</span>
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

        <!-- DASHBOARDS DE MÓDULO -->
        <section class="modules-section section-block">
          <div class="section-title-alt">
            <span class="material-symbols-outlined">grid_view</span>
            <h2>Módulos Principales</h2>
            <span class="section-tag">Gestión por dominio</span>
          </div>
          <div class="modules-grid">
            <RouterLink to="/sales/dashboard" class="module-link-card">
              <div class="m-icon blue"><span class="material-symbols-outlined">account_balance_wallet</span></div>
              <div class="m-info">
                <strong>Ventas</strong>
                <p>Gestión de presupuestos, pedidos y facturación.</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </RouterLink>
            <RouterLink to="/products/dashboard" class="module-link-card">
              <div class="m-icon yellow"><span class="material-symbols-outlined">inventory_2</span></div>
              <div class="m-info">
                <strong>Catálogo</strong>
                <p>Control de productos, stock y atributos.</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </RouterLink>
            <RouterLink to="/parties/dashboard" class="module-link-card">
              <div class="m-icon green"><span class="material-symbols-outlined">groups</span></div>
              <div class="m-info">
                <strong>Entidades</strong>
                <p>Base de datos de clientes y proveedores.</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </RouterLink>
            <RouterLink to="/mes/dashboard" class="module-link-card">
              <div class="m-icon purple"><span class="material-symbols-outlined">precision_manufacturing</span></div>
              <div class="m-info">
                <strong>Taller (MES)</strong>
                <p>Monitorización de producción y órdenes de trabajo.</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </RouterLink>
          </div>
        </section>
      </main>
      
      <!-- COLUMNA LATERAL -->
      <aside class="side-column">
        <section v-if="isAdmin" class="card admin-side-card section-block sticky-admin">
          <div class="ops-header">
            <span class="material-symbols-outlined">admin_panel_settings</span>
            <h2>Sistema</h2>
            <span class="section-tag">Administración</span>
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
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
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
  } catch (err) { console.error('Error dashboard:', err) }
  finally { isLoading.value = false }
}
function navigateTo(path: string) { router.push(path) }
onMounted(loadStats)
</script>

<style scoped>
.dashboard-page { 
  max-width: 1300px; 
  margin: 0 auto; 
  padding: 1rem;
  
  /* Variables de espaciado locales equilibradas */
  --dashboard-section-gap: 1.5rem;
  --dashboard-column-gap: 1.5rem;
  --dashboard-tag-gap: 0.5rem;
  --dashboard-card-gap: 0.75rem;
}

.dashboard-header { display: flex; justify-content: space-between; align-items: flex-end; gap: 1.5rem; }
.dashboard-header h1 { font-size: clamp(1.8rem, 2.5vw, 2.25rem); color: var(--color-text-primary); margin: 0; font-weight: 800; }
.btn-refresh { background: white; border: 1px solid var(--color-border); padding: 0.5rem; border-radius: 8px; cursor: pointer; color: var(--color-text-secondary); transition: 0.2s; }
.btn-refresh:hover { color: var(--color-primary); border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }

.section-spacing { margin-top: var(--dashboard-section-gap); }

.context-tags { display: flex; flex-wrap: wrap; gap: var(--dashboard-tag-gap); margin-top: 0.7rem; }
.context-tag {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  padding: 0.24rem 0.6rem;
  border-radius: 999px;
  background: #e2e8f0;
  color: #334155;
}
.context-tag.loading { background: #dbeafe; color: #1d4ed8; }

.kpi-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1.1rem; }
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

.dashboard-main-layout { display: grid; grid-template-columns: minmax(0, 1fr) 320px; gap: var(--dashboard-column-gap); align-items: start; }
.main-column, .side-column { display: flex; flex-direction: column; gap: var(--dashboard-section-gap); }
.section-block {
  background: white;
  border: 1px solid var(--color-border);
  border-radius: 14px;
  box-shadow: var(--box-shadow-sm);
}
.section-title-alt { display: flex; align-items: center; gap: var(--dashboard-tag-gap); margin-bottom: 0.85rem; }
.section-title-alt h2 { font-size: 0.9rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-secondary); }
.section-title-alt .material-symbols-outlined { color: var(--color-text-secondary); }
.section-tag {
  margin-left: auto;
  font-size: 0.66rem;
  padding: 0.18rem 0.55rem;
  border-radius: 999px;
  background: var(--color-background);
  color: var(--color-text-secondary);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.ops-section,
.modules-section { padding: clamp(1rem, 2vw, 1.35rem); }
.ops-header { display: flex; align-items: center; gap: var(--dashboard-tag-gap); margin-bottom: 0.85rem; }
.ops-header h2 { font-size: 0.9rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-primary); }
.ops-header .material-symbols-outlined { color: var(--color-primary); }
.ops-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: var(--dashboard-card-gap); }
.op-item { display: flex; flex-direction: column; align-items: center; padding: 1.25rem; background: var(--color-background); border-radius: 12px; text-decoration: none; color: var(--color-text-primary); font-weight: 700; font-size: 0.8rem; transition: 0.2s; text-align: center; gap: 0.75rem; border: 1px solid var(--color-border); }
.op-item:hover { transform: translateY(-3px); border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); color: var(--color-primary); }
.op-icon { color: var(--color-primary); }
.op-icon .material-symbols-outlined { font-size: 28px; }
.op-item.highlight { background: rgba(230, 184, 0, 0.05); }

.modules-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: var(--dashboard-card-gap); }
.module-link-card { display: flex; align-items: center; gap: 1rem; padding: 1.1rem; background: white; border-radius: 12px; border: 1px solid var(--color-border); text-decoration: none; transition: 0.2s; }
.module-link-card:hover { transform: translateX(5px); border-color: var(--color-secondary); box-shadow: var(--box-shadow-sm); }
.m-icon { width: 44px; height: 44px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.m-icon .material-symbols-outlined { font-size: 22px; }
.m-icon.blue { background: rgba(37, 99, 235, 0.05); color: #2563eb; }
.m-icon.yellow { background: rgba(217, 119, 6, 0.05); color: #d97706; }
.m-icon.green { background: rgba(22, 163, 74, 0.05); color: #16a34a; }
.m-icon.purple { background: rgba(147, 51, 234, 0.05); color: #9333ea; }
.m-info { min-width: 0; }
.m-info strong { display: block; font-size: 0.9rem; color: var(--color-text-primary); }
.m-info p { font-size: 0.75rem; color: var(--color-text-secondary); margin: 0.1rem 0 0; }
.module-link-card .arrow { margin-left: auto; color: var(--color-border); font-size: 18px; }

.admin-side-card { padding: 1.2rem; }
.sticky-admin { position: sticky; top: 88px; }
.side-links { display: flex; flex-direction: column; gap: var(--dashboard-tag-gap); }
.side-link-item { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: var(--color-background); border-radius: 8px; text-decoration: none; color: var(--color-text-primary); font-size: 0.85rem; font-weight: 600; border: 1px solid transparent; transition: 0.2s; }
.side-link-item:hover { background: white; border-color: var(--color-primary); color: var(--color-primary); transform: translateX(3px); }
.side-link-item .material-symbols-outlined { font-size: 20px; color: var(--color-text-secondary); }
.side-links .mt-2,
.dev-link { margin-top: 0; }
.dev-link { opacity: 0.7; }

@media (max-width: 1180px) {
  .dashboard-main-layout { grid-template-columns: 1fr; }
  .sticky-admin { position: static; top: auto; }
  .ops-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .modules-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (max-width: 760px) {
  :root {
    --dashboard-section-gap: 2rem;
  }
  .dashboard-header { flex-direction: column; align-items: flex-start; }
  .ops-grid { grid-template-columns: 1fr; }
  .modules-grid { grid-template-columns: 1fr; }
  .section-title-alt { flex-wrap: wrap; }
  .section-tag { margin-left: 0; }
}

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>
