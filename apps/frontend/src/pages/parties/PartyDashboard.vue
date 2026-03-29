<template>
  <BaseDashboardPage :is-loading="isLoading">
    <template #header>
      <PageHeader title="Entidades y CRM" :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Panel' }]">
        <template #icon><span class="material-symbols-outlined">groups</span></template>
        <template #actions>
          <button class="btn btn-primary btn-sm" @click="router.push('/parties/new')">
            <span class="material-symbols-outlined">person_add</span>
            <span>Nueva Entidad</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <div class="module-dashboard-content">
      <!-- 1. KPIs de Resumen -->
      <section class="stats-grid">
        <div class="stat-card clickable" @click="router.push('/parties')">
          <div class="stat-icon blue"><span class="material-symbols-outlined">groups</span></div>
          <div class="stat-info">
            <span class="stat-label">Total Entidades</span>
            <span class="stat-value">{{ counts.total }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/parties?role=CLIENT')">
          <div class="stat-icon green"><span class="material-symbols-outlined">person</span></div>
          <div class="stat-info">
            <span class="stat-label">Clientes</span>
            <span class="stat-value">{{ counts.clients }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/parties?role=SUPPLIER')">
          <div class="stat-icon purple"><span class="material-symbols-outlined">factory</span></div>
          <div class="stat-info">
            <span class="stat-label">Proveedores</span>
            <span class="stat-value">{{ counts.suppliers }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/parties?status=INACTIVE')">
          <div class="stat-icon red"><span class="material-symbols-outlined">block</span></div>
          <div class="stat-info">
            <span class="stat-label">Inactivos</span>
            <span class="stat-value">{{ counts.inactive }}</span>
          </div>
        </div>
      </section>

      <!-- 2. Accesos a Listados -->
      <section class="listings-grid">
        <RouterLink to="/parties" class="listing-link">
          <span class="material-symbols-outlined">contact_page</span>
          <span>Listado General</span>
        </RouterLink>
        <RouterLink to="/parties?role=CLIENT" class="listing-link">
          <span class="material-symbols-outlined">person_search</span>
          <span>Filtro Clientes</span>
        </RouterLink>
        <RouterLink to="/parties?role=SUPPLIER" class="listing-link">
          <span class="material-symbols-outlined">precision_manufacturing</span>
          <span>Filtro Proveedores</span>
        </RouterLink>
        <RouterLink to="/parties/dashboard" class="listing-link highlight-subtle">
          <span class="material-symbols-outlined">category</span>
          <span>Grupos de Entidad</span>
        </RouterLink>
      </section>

      <!-- 3. Actividad Reciente -->
      <section class="dashboard-section">
        <div class="section-header">
          <span class="material-symbols-outlined text-primary">recent_actors</span>
          <h2>Últimas Altas</h2>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Nombre</th>
                <th>Rol</th>
                <th>Identificación</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in recentParties" :key="p.id" class="row-clickable" @click="router.push(`/parties/${p.id}`)">
                <td><strong>{{ p.name }}</strong></td>
                <td>{{ formatRole(p.role) }}</td>
                <td><code class="code-badge">{{ p.tax_id || '—' }}</code></td>
                <td>
                  <span :class="['status-badge', p.status === 'ACTIVE' ? 'status-success' : 'status-secondary']">
                    {{ p.status === 'ACTIVE' ? 'Activo' : 'Inactivo' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <template #sidebar>
      <section class="sidebar-section">
        <div class="section-header">
          <span class="material-symbols-outlined">bolt</span>
          <h2>Operaciones</h2>
        </div>
        <div class="quick-actions-list">
          <RouterLink to="/parties/new" class="admin-card clickable">
            <span class="material-symbols-outlined text-primary">person_add</span>
            <div class="admin-card-info">
              <strong>Nueva Entidad</strong>
              <p>Registrar cliente o proveedor</p>
            </div>
          </RouterLink>
        </div>
      </section>
    </template>
  </BaseDashboardPage>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, RouterLink } from 'vue-router';
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import { partyApi } from '@/services/partyApi';

const router = useRouter();
const isLoading = ref(true);
const counts = ref({ total: 0, clients: 0, suppliers: 0, inactive: 0 });
const recentParties = ref([]);

async function loadPartyData() {
  isLoading.value = true;
  try {
    const res = await partyApi.listParties({ limit: 1000 });
    const all = res.data || [];
    recentParties.value = all.slice(0, 5);
    counts.value.total = all.length;
    counts.value.clients = all.filter(p => p.role === 'CLIENT' || p.role === 'BOTH').length;
    counts.value.suppliers = all.filter(p => p.role === 'SUPPLIER' || p.role === 'BOTH').length;
    counts.value.inactive = all.filter(p => p.status === 'INACTIVE').length;
  } catch (err) {
    console.error('Error dashboard entidades:', err);
  } finally {
    isLoading.value = false;
  }
}

function formatRole(r) { const map = { CLIENT: 'Cliente', SUPPLIER: 'Proveedor', BOTH: 'Ambos', CONTACT: 'Contacto' }; return map[r] || r; }

onMounted(loadPartyData);
</script>

<style scoped>
@import "@/design-system/_sections.css";

.module-dashboard-content { display: flex; flex-direction: column; gap: 1.5rem; }

.stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0.75rem; }
.stat-card { background: white; padding: 0.75rem 1rem; border-radius: 10px; border: 1px solid var(--color-border); display: flex; align-items: center; gap: 0.75rem; position: relative; transition: 0.2s; cursor: pointer; }
.stat-card:hover { transform: translateY(-2px); box-shadow: var(--box-shadow-md); border-color: var(--color-primary); }
.stat-icon { width: 40px; height: 40px; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
.stat-icon .material-symbols-outlined { font-size: 22px; }
.stat-icon.blue { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.stat-icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.stat-icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
.stat-icon.red { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.stat-label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.stat-value { font-size: 1.25rem; font-weight: 700; }

.listings-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.75rem; }
.listing-link { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; text-decoration: none; color: var(--color-text-primary); font-size: 0.85rem; font-weight: 600; transition: 0.2s; }
.listing-link:hover { background: var(--color-background); border-color: var(--color-secondary); color: var(--color-secondary); transform: translateY(-1px); }
.listing-link .material-symbols-outlined { color: var(--color-secondary); font-size: 1.25rem; }

.dashboard-section { background: white; padding: 0.75rem 1rem; border-radius: 10px; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.section-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.75rem; padding-bottom: 0.5rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.85rem; font-weight: 700; text-transform: uppercase; margin: 0; }

.quick-actions-list { display: flex; flex-direction: column; gap: 0.75rem; }
.admin-card { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem; background: var(--color-background); border-radius: 8px; border: 1px solid transparent; text-decoration: none; color: var(--color-text-primary); transition: 0.2s; }
.admin-card:hover { background: white; border-color: var(--color-primary); transform: translateX(4px); box-shadow: var(--box-shadow-sm); }
.admin-card-info strong { font-size: 0.8rem; display: block; }
.admin-card-info p { font-size: 0.65rem; color: var(--color-text-secondary); margin: 0; }

.code-badge { background: var(--color-background); padding: 0.15rem 0.35rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.75rem; font-weight: 700; color: var(--color-secondary); }

@media (max-width: 1180px) {
  .stats-grid, .listings-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 720px) {
  .stats-grid, .listings-grid { grid-template-columns: 1fr; }
}
</style>
