<template>
  

  <BaseCatalog
    title="Pedidos de Venta"
    icon="shopping_cart"
    :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Pedidos' }]"
    :items="orders"
    :is-loading="isLoading"
    :error="error"
    :has-filters="hasFilters"
    create-route="/sales/orders/new"
    create-text="Nuevo Pedido"
    empty-icon="shopping_basket"
    empty-text="No se han encontrado pedidos"
    @clear-filters="clearFilters"
    @refresh="fetchOrders"
    @click-item="(item) => navigateToDetail(item.id)"
  >
    <template #filters>
      <div class="filter-group">
        <PartySelector
          v-model="filters.partyId"
          label="Cliente"
          placeholder="Filtrar por cliente..."
          role-filter="CLIENT"
          :required="false"
        />
      </div>

      <div class="filter-group">
        <label>Búsqueda</label>
        <input v-model="filters.searchText" type="text" placeholder="Ref. o nombre..." />
      </div>

      <div class="filter-group">
        <label>Estado</label>
        <select v-model="filters.status">
          <option value="">Cualquier estado</option>
          <option value="PENDIENTE">Pendiente</option>
          <option value="EN_PREPARACION">En Preparación</option>
          <option value="ENTREGADO">Entregado</option>
          <option value="CANCELADO">Cancelado</option>
          <option value="FACTURADO">Facturado</option>
        </select>
      </div>

      <div class="filter-group">
        <label>Desde</label>
        <input v-model="filters.fromDate" type="date" />
      </div>

      <div class="filter-group">
        <label>Hasta</label>
        <input v-model="filters.toDate" type="date" />
      </div>
    </template>

    <template #table-header>
      <th>Número</th>
      <th>Cliente</th>
      <th>Fecha</th>
      <th>Entrega</th>
      <th>Estado</th>
      <th class="align-right">Total</th>
      <th>MES</th>
      <th class="align-right">Acciones</th>
    </template>

    <template #item="{ item }">
      <td><span class="order-ref">{{ item.orderNumber }}</span></td>
      <td>{{ formatPartyId(item.partyId) }}</td>
      <td>{{ formatDate(item.orderDate) }}</td>
      <td>{{ formatDate(item.deliveryDate) }}</td>
      <td>
        <span :class="['status-badge', `status-${salesApi.getStatusClass(item.status)}`]">
          {{ salesApi.getStatusLabel(item.status) }}
        </span>
      </td>
      <td class="align-right"><strong>{{ salesApi.formatMoney(item.total) }}</strong></td>
      <td @click.stop>
        <span v-if="getMesWorkIds(item.lineItems).length > 0" class="mes-badge">
          <Factory :size="16" />
          {{ getMesSummary(item.lineItems) }}
        </span>
        <span v-else class="text-muted">—</span>
      </td>
      <td class="align-right" @click.stop>
        <div class="action-buttons">
          <router-link :to="`/sales/orders/${item.id}`" class="btn-icon" title="Ver detalle">
            <Eye :size="18" />
          </router-link>
        </div>
      </td>
    </template>
  </BaseCatalog>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRouter } from 'vue-router';
import { Factory, Eye } from 'lucide-vue-next';
import BaseCatalog from '@/components/shared/BaseCatalog.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { mesApi } from '@/services/mesApi';

const router = useRouter();
const orders = ref([]);
const isLoading = ref(false);
const error = ref('');
const partiesCache = ref({});
const mesWorksCache = ref({});

const filters = ref({ partyId: '', searchText: '', status: '', fromDate: '', toDate: '', limit: 50 });

const hasFilters = computed(() => {
  return (
    filters.value.partyId !== '' || 
    filters.value.searchText !== '' || 
    filters.value.status !== '' || 
    filters.value.fromDate !== ''
  );
});

let searchDebounceTimer = null;
watch(() => [filters.value.searchText, filters.value.partyId, filters.value.status, filters.value.fromDate, filters.value.toDate], () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
  searchDebounceTimer = setTimeout(() => fetchOrders(), 350);
});

onMounted(() => {
  const today = new Date();
  const ago = new Date(today); ago.setDate(today.getDate() - 90);
  filters.value.fromDate = ago.toISOString().split('T')[0];
  filters.value.toDate = today.toISOString().split('T')[0];
  fetchOrders();
});

onBeforeUnmount(() => { if (searchDebounceTimer) clearTimeout(searchDebounceTimer); });

async function fetchOrders() {
  isLoading.value = true;
  try {
    const res = await salesApi.listOrders(filters.value);
    const rawData = Array.isArray(res) ? res : (res.data || []);
    // Ordenar por fecha descendente (más recientes primero)
    orders.value = rawData.sort((a, b) => new Date(b.orderDate) - new Date(a.orderDate));
    await loadPartyNames();
    await loadMesWorksForOrders();
  } catch (err) { error.value = err.message; }
  finally { isLoading.value = false; }
}

async function loadPartyNames() {
  const ids = [...new Set(orders.value.map(o => o.partyId).filter(id => id && !partiesCache.value[id]))];
  if (ids.length === 0) return;
  try {
    const map = await partyApi.getPartiesBatch(ids);
    ids.forEach(id => partiesCache.value[id] = map[id]?.name || 'Desconocido');
  } catch (err) {}
}

async function loadMesWorksForOrders() {
  const ids = [...new Set(orders.value.flatMap(o => (o.lineItems || []).map(i => i.mesWorkId)).filter(id => id && !mesWorksCache.value[id]))];
  if (ids.length === 0) return;
  const results = await Promise.allSettled(ids.map(id => mesApi.getWorkOrder(id)));
  results.forEach((r, i) => mesWorksCache.value[ids[i]] = r.status === 'fulfilled' ? r.value : null);
}

function clearFilters() { filters.value = { partyId: '', searchText: '', status: '', fromDate: '', toDate: '', limit: 50 }; fetchOrders(); }
function navigateToDetail(id) { router.push(`/sales/orders/${id}`); }

function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: '2-digit', day: '2-digit' }) : '—'; }
function formatPartyId(id) { return partiesCache.value[id] || '...'; }

function getMesSummary(items) {
  const ids = getMesWorkIds(items); if (ids.length === 0) return '—';
  const first = mesWorksCache.value[ids[0]]?.work_number || 'MES';
  return ids.length === 1 ? first : `${first} +${ids.length-1}`;
}
function getMesWorkIds(items) { return [...new Set((items || []).map(i => i.mesWorkId).filter(id => !!id))]; }
</script>

<style scoped>
.order-ref { font-family: var(--font-family-mono); font-weight: 700; color: var(--color-secondary); }
.align-right { text-align: right; }

.mes-badge { display: inline-flex; align-items: center; gap: 0.35rem; padding: 0.25rem 0.5rem; background: rgba(0, 35, 149, 0.05); border-radius: 4px; color: var(--color-secondary); font-size: 0.75rem; font-weight: 600; }

.action-buttons { display: flex; justify-content: flex-end; }
.btn-icon { color: var(--color-text-secondary); transition: 0.2s; padding: 0.4rem; border-radius: 6px; border: none; background: transparent; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; }
.btn-icon:hover { color: var(--color-text-primary); background: rgba(0,0,0,0.05); }
</style>