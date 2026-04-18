<template>
  <div class="page-layout">
    <BaseCatalog
      title="Gestión de Albaranes"
      icon="local_shipping"
      :breadcrumbs="[{ label: 'Ventas', to: '/sales/dashboard' }, { label: 'Albaranes' }]"
      :items="deliveryNotes"
      :is-loading="isLoading"
      :has-filters="hasFilters"
      empty-icon="local_shipping"
      empty-text="No se han encontrado albaranes"
      @clear-filters="clearFilters"
      @refresh="fetchNotes"
      @click-item="(item) => navigateToDetail(item.id)"
    >
      <template #header-actions>
        <button class="btn btn-outline" @click="router.push('/sales/orders')">
          <span class="material-symbols-outlined">list_alt</span>
          <span>Catálogo de Pedidos</span>
        </button>
      </template>

      <template #filters>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input v-model="filters.searchText" type="text" placeholder="Ref. albarán..." />
        </div>

        <div class="filter-group">
          <label>Estado</label>
          <select v-model="filters.status" class="form-input-sm">
            <option value="">Todos</option>
            <option value="PENDING">Pendiente</option>
            <option value="DELIVERED">Entregado</option>
            <option value="CANCELLED">Anulado</option>
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
        <th>Fecha Entrega</th>
        <th>Pedido Relacionado</th>
        <th>Cliente</th>
        <th>Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><span class="order-ref">{{ item.deliveryNoteNumber }}</span></td>
        <td>{{ formatDate(item.deliveryDate) }}</td>
        <td>
          <router-link v-if="item.salesOrderId" :to="`/sales/orders/${item.salesOrderId}`" class="link-sm" @click.stop>
            #{{ formatOrderNumber(item.salesOrderId) }}
          </router-link>
          <span v-else class="text-muted">—</span>
        </td>
        <td>{{ formatPartyName(item.partyId) }}</td>
        <td>
          <span :class="['status-badge', `status-${salesApi.getStatusClass(item.status)}`]">
            {{ salesApi.getStatusLabel(item.status) }}
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button class="btn-icon" @click="navigateToDetail(item.id)" title="Ver detalle">
              <span class="material-symbols-outlined">visibility</span>
            </button>
          </div>
        </td>
      </template>
    </BaseCatalog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';

import BaseCatalog from '@/components/shared/BaseCatalog.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';

const router = useRouter();
const deliveryNotes = ref([]);
const isLoading = ref(false);
const filters = ref({ searchText: '', status: '', fromDate: '', toDate: '' });
const partiesCache = ref({});
const ordersCache = ref({});

const hasFilters = computed(() => filters.value.searchText || filters.value.status || filters.value.fromDate || filters.value.toDate);

let searchDebounceTimer = null;
watch(() => [filters.value.searchText, filters.value.status, filters.value.fromDate, filters.value.toDate], () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
  searchDebounceTimer = setTimeout(() => fetchNotes(), 350);
});

onMounted(() => fetchNotes());

async function fetchNotes() {
  isLoading.value = true;
  try {
    const res = await salesApi.listDeliveryNotes(filters.value);
    deliveryNotes.value = res.data || (Array.isArray(res) ? res : []);
    
    if (deliveryNotes.value.length > 0) {
      await Promise.all([loadPartyNames(), loadOrderNumbers()]);
    }
  } catch (err) {
    console.error('Error fetching delivery notes:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyNames() {
  const ids = [...new Set(deliveryNotes.value.map(n => n.partyId).filter(id => id && !partiesCache.value[id]))];
  if (ids.length === 0) return;
  try {
    const map = await partyApi.getPartiesBatch(ids);
    ids.forEach(id => partiesCache.value[id] = map[id]?.name || 'Desconocido');
  } catch (err) {}
}

async function loadOrderNumbers() {
  const ids = [...new Set(deliveryNotes.value.map(n => n.salesOrderId).filter(id => id && !ordersCache.value[id]))];
  if (ids.length === 0) return;
  try {
    const results = await Promise.allSettled(ids.map(id => salesApi.getOrder(id)));
    results.forEach((r, i) => {
      if (r.status === 'fulfilled') ordersCache.value[ids[i]] = r.value.orderNumber;
    });
  } catch (err) {}
}

function formatPartyName(id) { return partiesCache.value[id] || 'Cargando...'; }
function formatOrderNumber(id) { return ordersCache.value[id] || id?.substring(0, 8) || '...'; }

function clearFilters() { filters.value = { searchText: '', status: '', fromDate: '', toDate: '' }; fetchNotes(); }
function navigateToDetail(id) { router.push(`/sales/delivery-notes/${id}`); }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: '2-digit', day: '2-digit' }) : '—'; }
</script>

<style scoped>
.order-ref { font-family: var(--font-family-mono); font-weight: 700; color: var(--color-secondary); }
.align-right { text-align: right; }
.link-sm { font-size: 0.85rem; color: var(--color-secondary); text-decoration: none; font-weight: 600; background: var(--color-background); padding: 0.2rem 0.5rem; border-radius: 4px; border: 1px solid var(--color-border); }
.link-sm:hover { background: white; border-color: var(--color-secondary); }

.action-buttons { display: flex; justify-content: flex-end; }
.btn-icon { color: var(--color-text-secondary); transition: 0.2s; padding: 0.4rem; border-radius: 6px; border: none; background: transparent; cursor: pointer; }
.btn-icon:hover { color: var(--color-text-primary); background: rgba(0,0,0,0.05); }
</style>