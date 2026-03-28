<template>
  <Navbar />
  
  <BaseDashboardPage :is-loading="isLoading">
    <template #header>
      <PageHeader title="Gestión de Entidades" :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Dashboard' }]">
        <template #icon><span class="material-symbols-outlined">groups</span></template>
        <template #actions>
          <button class="btn btn-primary" @click="router.push('/parties/new')">
            <span class="material-symbols-outlined">person_add</span>
            <span>Nueva Entidad</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <div class="party-dashboard-content">
      <!-- KPIs Superiores -->
      <section class="stats-grid mb-8">
        <div class="stat-card clickable" @click="router.push('/parties/list')">
          <div class="stat-icon blue"><span class="material-symbols-outlined">groups</span></div>
          <div class="stat-info">
            <span class="stat-label">Total Entidades</span>
            <span class="stat-value">{{ counts.total }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
        <div class="stat-card clickable" @click="router.push('/parties/list?role=CLIENT')">
          <div class="stat-icon green"><span class="material-symbols-outlined">person</span></div>
          <div class="stat-info">
            <span class="stat-label">Clientes</span>
            <span class="stat-value">{{ counts.clients }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
        <div class="stat-card clickable" @click="router.push('/parties/list?role=SUPPLIER')">
          <div class="stat-icon purple"><span class="material-symbols-outlined">factory</span></div>
          <div class="stat-info">
            <span class="stat-label">Proveedores</span>
            <span class="stat-value">{{ counts.suppliers }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
        <div class="stat-card clickable" @click="router.push('/parties/list?status=INACTIVE')">
          <div class="stat-icon red"><span class="material-symbols-outlined">block</span></div>
          <div class="stat-info">
            <span class="stat-label">Inactivos</span>
            <span class="stat-value">{{ counts.inactive }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
      </section>

      <div class="dashboard-grid-main">
        <!-- Entidades Recientes -->
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
    </div>

    <template #sidebar>
      <section class="sidebar-section">
        <div class="section-header">
          <span class="material-symbols-outlined">contact_support</span>
          <h2>Ayuda y Soporte</h2>
        </div>
        <div class="help-notice">
          <p class="help-text">
            Las entidades representan a cualquier persona u organización con la que TramaTex interactúa. Un contacto puede ser cliente y proveedor simultáneamente.
          </p>
        </div>
      </section>
    </template>
  </BaseDashboardPage>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
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
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
.stat-card { background: white; padding: 1.25rem; border-radius: 12px; border: 1px solid var(--color-border); display: flex; align-items: center; gap: 1rem; position: relative; transition: 0.2s; cursor: pointer; }
.stat-card:hover { transform: translateY(-3px); box-shadow: var(--box-shadow-md); border-color: var(--color-primary); }
.stat-icon { width: 44px; height: 44px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.stat-icon .material-symbols-outlined { font-size: 24px; }
.stat-icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.stat-icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.stat-icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.stat-icon.red { background: rgba(239, 68, 68, 0.1); color: #dc2626; }
.stat-info { display: flex; flex-direction: column; }
.stat-label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.stat-value { font-size: 1.5rem; font-weight: 700; }
.stat-link-arrow { position: absolute; right: 1rem; color: var(--color-border); font-size: 18px; }
.dashboard-section { background: white; padding: 1.5rem; border-radius: 12px; border: 1px solid var(--color-border); }
.section-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; padding-bottom: 0.75rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.9rem; font-weight: 700; text-transform: uppercase; margin: 0; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; }
.help-notice { padding: 1.25rem; background: rgba(59, 130, 246, 0.05); border-radius: 12px; border: 1px dashed rgba(59, 130, 246, 0.3); }
.help-text { font-size: 0.8rem; color: var(--color-text-secondary); line-height: 1.5; }
</style>
