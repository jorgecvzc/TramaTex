<template>
  <BaseDashboardPage :is-loading="isLoading">
    <template #header>
      <PageHeader title="Ventas y Facturación">
        <template #icon><span class="material-symbols-outlined">payments</span></template>
        <template #actions>
          <button class="btn btn-outline btn-sm" @click="loadSalesData" :disabled="isLoading">
            <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
            <span>Actualizar</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <div class="module-dashboard-content">
      <!-- 1. KPIs de Resumen -->
      <section class="stats-grid">
        <div class="stat-card clickable" @click="router.push('/sales/quotes')">
          <div class="stat-icon yellow"><span class="material-symbols-outlined">description</span></div>
          <div class="stat-info">
            <span class="stat-label">Presupuestos</span>
            <span class="stat-value">{{ counts.pendingQuotes }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/sales/orders')">
          <div class="stat-icon blue"><span class="material-symbols-outlined">shopping_cart</span></div>
          <div class="stat-info">
            <span class="stat-label">Pedidos Activos</span>
            <span class="stat-value">{{ counts.activeOrders }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/sales/delivery-notes')">
          <div class="stat-icon green"><span class="material-symbols-outlined">local_shipping</span></div>
          <div class="stat-info">
            <span class="stat-label">Pend. Entrega</span>
            <span class="stat-value">{{ counts.pendingDelivery }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/sales/invoices')">
          <div class="stat-icon purple"><span class="material-symbols-outlined">receipt_long</span></div>
          <div class="stat-info">
            <span class="stat-label">Facturas Mes</span>
            <span class="stat-value">{{ counts.monthlyInvoices }}</span>
          </div>
        </div>
      </section>

      <!-- 2. Accesos a Listados -->
      <section class="listings-grid">
        <RouterLink to="/sales/quotes" class="listing-link">
          <span class="material-symbols-outlined">request_quote</span>
          <span>Listado de Presupuestos</span>
        </RouterLink>
        <RouterLink to="/sales/orders" class="listing-link">
          <span class="material-symbols-outlined">inventory</span>
          <span>Listado de Pedidos</span>
        </RouterLink>
        <RouterLink to="/sales/delivery-notes" class="listing-link">
          <span class="material-symbols-outlined">conveyor_belt</span>
          <span>Listado de Albaranes</span>
        </RouterLink>
        <RouterLink to="/sales/invoices" class="listing-link">
          <span class="material-symbols-outlined">receipt_long</span>
          <span>Listado de Facturas</span>
        </RouterLink>
      </section>

      <!-- 3. Actividad Reciente -->
      <section class="dashboard-section">
        <div class="section-header">
          <span class="material-symbols-outlined text-primary">history</span>
          <h2>Últimos Movimientos</h2>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Pedido</th>
                <th>Cliente</th>
                <th>Total</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in recentOrders" :key="order.id" class="row-clickable" @click="router.push(`/sales/orders/${order.id}`)">
                <td><code class="code-badge">{{ order.orderNumber }}</code></td>
                <td>{{ order.partyName || '...' }}</td>
                <td><strong>{{ salesApi.formatMoney(order.total) }}</strong></td>
                <td>
                  <span :class="['status-badge', `status-${salesApi.getStatusClass(order.status)}`]">
                    {{ salesApi.getStatusLabel(order.status) }}
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
          <h2>Accesos Rápidos</h2>
        </div>
        <div class="quick-actions-list">
          <RouterLink to="/sales/tickets/new" class="admin-card clickable highlight">
            <span class="material-symbols-outlined text-primary">point_of_sale</span>
            <div class="admin-card-info">
              <strong>Venta Directa</strong>
              <p>TPV mostrador</p>
            </div>
          </RouterLink>
          <RouterLink to="/sales/quotes/new" class="admin-card clickable">
            <span class="material-symbols-outlined text-secondary">add_notes</span>
            <div class="admin-card-info">
              <strong>Crear Presupuesto</strong>
              <p>Nueva oferta comercial</p>
            </div>
          </RouterLink>
          <RouterLink to="/sales/orders/new" class="admin-card clickable">
            <span class="material-symbols-outlined">add_shopping_cart</span>
            <div class="admin-card-info">
              <strong>Crear Pedido</strong>
              <p>Nuevo pedido de venta</p>
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
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';

const router = useRouter();
const isLoading = ref(true);
const counts = ref({ activeOrders: 0, pendingQuotes: 0, pendingDelivery: 0, monthlyInvoices: 0 });
const recentOrders = ref([]);

async function loadSalesData() {
  isLoading.value = true;
  try {
    const [orders, quotes, dnotes, invoices] = await Promise.all([
      salesApi.listOrders({ limit: 5 }),
      // Do not request quotes with limit=1: when backend returns array (without total),
      // salesApi falls back to rawData.length and the counter becomes incorrectly 1.
      salesApi.listQuotes({}),
      salesApi.listDeliveryNotes({ status: 'PENDIENTE', limit: 1 }),
      salesApi.listInvoices({ limit: 1 })
    ]);
    
    recentOrders.value = orders.data || [];
    counts.value.activeOrders = orders.total || 0;
    counts.value.pendingQuotes = quotes.total || 0;
    counts.value.pendingDelivery = dnotes.total || 0;
    counts.value.monthlyInvoices = invoices.total || 0;

    const partyIds = [...new Set(recentOrders.value.map(o => o.partyId))];
    if (partyIds.length > 0) {
      const parties = await partyApi.getPartiesBatch(partyIds);
      recentOrders.value.forEach(o => o.partyName = parties[o.partyId]?.name);
    }
  } catch (err) {
    console.error('Error dashboard ventas:', err);
  } finally {
    isLoading.value = false;
  }
}

onMounted(loadSalesData);
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
.stat-icon.yellow { background: rgba(230, 184, 0, 0.1); color: #E6B800; }
.stat-icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.stat-icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
.stat-info { display: flex; flex-direction: column; gap: 0.25rem; }
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
.admin-card.highlight { background: rgba(230, 184, 0, 0.05); border: 1px dashed var(--color-primary); }
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
