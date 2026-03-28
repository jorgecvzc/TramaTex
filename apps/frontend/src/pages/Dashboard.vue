<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { useAuth } from '@/composables'
import { useAuthStore } from '@/stores/auth'

// API Services
import { partyApi } from '@/services/partyApi'
import { productApi } from '@/services/productApi'
import salesApi from '@/services/salesApi'
import { mesApi } from '@/services/mesApi'

const router = useRouter()
const authStore = useAuthStore()
const { user } = useAuth()

const isAdmin = computed(() => authStore.isAdmin)

// Reactive Stats
const counts = ref({
  parties: 0,
  products: 0,
  ordersToday: 0,
  mesTasks: 0
})

const isLoading = ref(true)

const stats = computed(() => [
  { label: 'Entidades', value: counts.value.parties.toString(), icon: 'groups_2', color: 'blue', link: '/parties' },
  { label: 'Productos', value: counts.value.products.toString(), icon: 'inventory_2', color: 'green', link: '/products' },
  { label: 'Pedidos Hoy', value: counts.value.ordersToday.toString(), icon: 'shopping_cart', color: 'yellow', link: '/sales/orders' },
  { label: 'Tareas MES', value: counts.value.mesTasks.toString(), icon: 'precision_manufacturing', color: 'purple', link: '/mes/dashboard' }
])

async function fetchDashboardData() {
  isLoading.value = true
  try {
    const [partiesRes, productsRes, salesRes, mesRes] = await Promise.all([
      partyApi.listParties({ pageSize: 1000 }),
      productApi.listProducts({ pageSize: 1000 }),
      salesApi.listOrders({ 
        pageSize: 1000,
        startDate: new Date().toISOString().split('T')[0]
      }),
      mesApi.getWorkOrderDashboardStats().catch(() => ({ total: 0, by_status: { IN_PROGRESS: 0 } }))
    ])

    counts.value.parties = partiesRes.data?.length || 0
    counts.value.products = productsRes.data?.length || 0
    counts.value.ordersToday = salesRes.data?.length || 0
    
    const mesStats = mesRes as any
    counts.value.mesTasks = (mesStats.by_status?.IN_PROGRESS || 0) + (mesStats.by_status?.PENDING || 0)
  } catch (error) {
    console.error('Error cargando datos del dashboard:', error)
  } finally {
    isLoading.value = false
  }
}

function navigateTo(link: string) {
  router.push(link)
}

onMounted(() => {
  fetchDashboardData()
})
</script>

<template>
  <div class="dashboard">
    <Navbar />
    
    <main class="dashboard-content">
      <header class="dashboard-header">
        <div class="header-title">
          <h1>Bienvenido, {{ user?.name || user?.email || 'Usuario' }}</h1>
          <p class="subtitle">Panel de control general de TramaTex ERP/MES</p>
        </div>
        <button @click="fetchDashboardData" class="btn btn-outline btn-sm" :disabled="isLoading">
          <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
          Actualizar datos
        </button>
      </header>

      <!-- Stats Grid -->
      <div class="stats-grid">
        <div 
          v-for="stat in stats" 
          :key="stat.label" 
          class="stat-card clickable" 
          :class="{ 'loading-pulse': isLoading }"
          @click="navigateTo(stat.link)"
        >
          <div class="stat-icon" :class="stat.color">
            <span class="material-symbols-outlined">{{ stat.icon }}</span>
          </div>
          <div class="stat-info">
            <span class="stat-label">{{ stat.label }}</span>
            <span class="stat-value">{{ isLoading ? '...' : stat.value }}</span>
          </div>
          <div class="stat-link-arrow">
            <span class="material-symbols-outlined">arrow_forward</span>
          </div>
        </div>
      </div>

      <div class="dashboard-grid">
        <!-- Main Actions (Shortcuts) -->
        <section class="dashboard-section main-actions">
          <div class="section-header">
            <span class="material-symbols-outlined">bolt</span>
            <h2>Accesos Directos</h2>
            <span class="header-tag">Fijos</span>
          </div>
          <div class="actions-grid">
            <RouterLink to="/sales/orders/new" class="action-card">
              <span class="material-symbols-outlined icon-secondary">add_shopping_cart</span>
              <span>Nuevo Pedido</span>
            </RouterLink>
            <RouterLink to="/products/new" class="action-card">
              <span class="material-symbols-outlined icon-secondary">add_box</span>
              <span>Nuevo Producto</span>
            </RouterLink>
            <RouterLink to="/parties/new" class="action-card">
              <span class="material-symbols-outlined icon-secondary">person_add</span>
              <span>Nueva Entidad</span>
            </RouterLink>
            <RouterLink to="/mes/terminal" class="action-card highlight">
              <span class="material-symbols-outlined icon-secondary">tablet_mac</span>
              <span>Terminal de Taller</span>
            </RouterLink>
          </div>
          <div class="shortcuts-footer">
            <span class="material-symbols-outlined">info</span>
            <p>La personalización de accesos directos por usuario estará disponible Post-MVP.</p>
          </div>
        </section>

        <!-- Avisos (Post-MVP) -->
        <section class="dashboard-section notices">
          <div class="section-header">
            <span class="material-symbols-outlined">notifications</span>
            <h2>Avisos y Notificaciones</h2>
          </div>
          <div class="notices-placeholder">
            <div class="placeholder-content">
              <span class="material-symbols-outlined icon-muted">construction</span>
              <h3>Próximamente</h3>
              <p>Módulo centralizado de alertas automáticas (stock, urgencias, MES).</p>
            </div>
          </div>
        </section>

        <!-- Admin Quick Access -->
        <section v-if="isAdmin" class="dashboard-section admin">
          <div class="section-header">
            <span class="material-symbols-outlined">admin_panel_settings</span>
            <h2>Administración</h2>
          </div>
          <div class="admin-links">
            <div class="admin-card clickable" @click="navigateTo('/admin/users')">
              <span class="material-symbols-outlined icon-muted">manage_accounts</span>
              <div class="admin-card-info">
                <strong>Gestión de Usuarios</strong>
                <p>Accesos, roles y seguridad</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </div>
            <div class="admin-card clickable" @click="navigateTo('/admin/print-profile')">
              <span class="material-symbols-outlined icon-muted">business</span>
              <div class="admin-card-info">
                <strong>Datos de la Empresa</strong>
                <p>Perfil fiscal y configuración</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </div>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: var(--color-background);
}

.dashboard-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.dashboard-header h1 {
  font-size: 2rem;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem;
  font-family: var(--font-family-brand);
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 1.1rem;
  margin: 0;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2.5rem;
}

.stat-card {
  background: var(--color-surface);
  padding: 1.5rem;
  border-radius: var(--border-radius-lg);
  display: flex;
  align-items: center;
  gap: 1.25rem;
  position: relative;
  box-shadow: var(--box-shadow-sm);
  border: 1px solid var(--color-border);
  transition: all 0.2s ease;
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable:hover {
  transform: translateY(-3px);
  box-shadow: var(--box-shadow-md);
  border-color: var(--color-primary);
}

.stat-card.clickable:hover .stat-link-arrow {
  color: var(--color-primary);
  transform: translateX(3px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon .material-symbols-outlined { font-size: 32px; }
.stat-icon.blue { background-color: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.stat-icon.green { background-color: rgba(34, 197, 94, 0.1); color: #22c55e; }
.stat-icon.yellow { background-color: rgba(230, 184, 0, 0.1); color: #E6B800; }
.stat-icon.purple { background-color: rgba(168, 85, 247, 0.1); color: #a855f7; }

.stat-info { display: flex; flex-direction: column; }
.stat-label { font-size: 0.75rem; color: var(--color-text-secondary); font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; }
.stat-value { font-size: 1.75rem; font-weight: 700; color: var(--color-text-primary); }

.stat-link-arrow { 
  position: absolute; 
  right: 1.25rem; 
  color: var(--color-border); 
  transition: all 0.2s; 
}

/* Sections Layout */
.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
  align-items: start;
}

.dashboard-section {
  background: var(--color-surface);
  padding: 1.5rem;
  border-radius: var(--border-radius-lg);
  box-shadow: var(--box-shadow-sm);
  border: 1px solid var(--color-border);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--color-background);
}

.section-header h2 {
  font-size: 1rem;
  font-weight: 700;
  margin: 0;
  color: var(--color-text-primary);
  text-transform: uppercase;
  letter-spacing: 0.025em;
  flex: 1;
}

.section-header .material-symbols-outlined { color: var(--color-text-secondary); font-size: 22px; }

.header-tag {
  font-size: 0.6rem;
  font-weight: 700;
  padding: 0.2rem 0.5rem;
  background: var(--color-background);
  color: var(--color-text-secondary);
  border-radius: 4px;
}

/* Actions Grid */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 1rem;
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1.5rem;
  background-color: var(--color-background);
  border-radius: 12px;
  text-decoration: none;
  color: var(--color-text-primary);
  font-weight: 600;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.icon-secondary { color: var(--color-secondary); font-size: 32px; }

.action-card:hover {
  transform: translateY(-2px);
  background-color: white;
  box-shadow: var(--box-shadow-md);
  border-color: var(--color-primary);
}

.action-card.highlight { background-color: rgba(230, 184, 0, 0.05); border: 1px dashed var(--color-primary); }

.shortcuts-footer {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-background);
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--color-text-secondary);
}

.shortcuts-footer .material-symbols-outlined { font-size: 16px; }
.shortcuts-footer p { margin: 0; font-size: 0.75rem; font-style: italic; }

/* Notices Placeholder */
.notices-placeholder {
  padding: 2rem 1rem;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 12px;
  border: 1px dashed var(--color-border);
}

.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 0.75rem;
}

.icon-muted { color: var(--color-border); font-size: 48px; }
.placeholder-content h3 { margin: 0; font-size: 0.9rem; color: var(--color-text-secondary); }
.placeholder-content p { margin: 0; font-size: 0.8rem; color: var(--color-text-secondary); max-width: 200px; }

/* Admin Section */
.admin-links { display: flex; flex-direction: column; gap: 0.75rem; }
.admin-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background-color: var(--color-background);
  border-radius: 10px;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.admin-card.clickable { cursor: pointer; }

.admin-card.clickable:hover {
  background-color: white;
  border-color: var(--color-primary);
  box-shadow: var(--box-shadow-sm);
  transform: translateX(4px);
}

.admin-card .material-symbols-outlined:first-child { font-size: 24px; color: var(--color-text-secondary); }
.admin-card-info { flex: 1; }
.admin-card-info strong { display: block; font-size: 0.9rem; }
.admin-card-info p { font-size: 0.75rem; color: var(--color-text-secondary); margin: 0; }
.admin-card .arrow { color: var(--color-border); transition: color 0.2s; font-size: 18px; }
.admin-card.clickable:hover .arrow { color: var(--color-primary); }

/* Helpers */
.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
.loading-pulse { opacity: 0.7; }

@media (max-width: 1024px) { .dashboard-grid { grid-template-columns: 1fr; } }
</style>
