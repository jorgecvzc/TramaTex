<template>
  <div class="page-layout">
    
    
    <BaseCatalog
      title="Gestión de Facturas"
      :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Facturas' }]"
      :items="invoices"
      :is-loading="isLoading"
      :error="error"
      :has-filters="hasFilters"
      empty-icon="receipt_long"
      empty-text="No hay facturas registradas"
      @clear-filters="clearFilters"
      @refresh="fetchInvoices"
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
          <input v-model="filters.searchText" type="text" placeholder="Ref. factura..." />
        </div>

        <div class="filter-group">
          <label>Tipo</label>
          <select v-model="filters.type">
            <option value="">Todos los tipos</option>
            <option value="STANDARD">Estándar</option>
            <option value="SIMPLIFIED">Simplificada</option>
          </select>
        </div>

        <div class="filter-group">
          <label>Estado</label>
          <select v-model="filters.status">
            <option value="">Cualquier estado</option>
            <option value="DRAFT">Borrador</option>
            <option value="ISSUED">Emitida</option>
            <option value="PAID">Pagada</option>
            <option value="CANCELLED">Cancelada</option>
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
        <th>Tipo</th>
        <th>Fecha Factura</th>
        <th>Cliente</th>
        <th>Estado</th>
        <th class="align-right">Total</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><span class="order-ref">{{ item.invoiceNumber }}</span></td>
        <td>
          <span :class="['type-tag', item.type === 'SIMPLIFIED' ? 'simplified' : 'standard']">
            {{ item.type === 'SIMPLIFIED' ? 'Simplificada' : 'Estándar' }}
          </span>
        </td>
        <td>{{ formatDate(item.invoiceDate) }}</td>
        <td>{{ formatPartyName(item.partyId) }}</td>
        <td>
          <span :class="['status-badge', `status-${getStatusClass(item.status)}`]">
            {{ getStatusLabel(item.status) }}
          </span>
        </td>
        <td class="align-right"><strong>{{ salesApi.formatMoney(item.total) }}</strong></td>
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
const invoices = ref([]);
const isLoading = ref(false);
const error = ref('');
const filters = ref({ searchText: '', status: '', type: '', fromDate: '', toDate: '' });
const partiesCache = ref({});

const hasFilters = computed(() => filters.value.searchText || filters.value.status || filters.value.type || filters.value.fromDate || filters.value.toDate);

// Lógica de filtrado automática con debounce
let searchDebounceTimer = null;
watch(filters, () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
  searchDebounceTimer = setTimeout(() => fetchInvoices(), 350);
}, { deep: true });

onMounted(() => fetchInvoices());

async function fetchInvoices() {
  isLoading.value = true;
  error.value = '';
  try {
    const res = await salesApi.listInvoices(filters.value);
    invoices.value = res.data || (Array.isArray(res) ? res : []);
    
    if (invoices.value.length > 0) {
      await loadPartyNames();
    }
  } catch (err) {
    console.error('Error fetching invoices:', err);
    error.value = 'No se han podido cargar las facturas.';
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyNames() {
  const ids = [...new Set(invoices.value.map(i => i.partyId).filter(id => id && !partiesCache.value[id]))];
  if (ids.length === 0) return;
  try {
    const map = await partyApi.getPartiesBatch(ids);
    ids.forEach(id => partiesCache.value[id] = map[id]?.name || 'Desconocido');
  } catch (err) {}
}

function formatPartyName(id) { return partiesCache.value[id] || 'Cargando...'; }

function clearFilters() { 
  filters.value = { searchText: '', status: '', type: '', fromDate: '', toDate: '' }; 
  // El watch se encargará de llamar a fetchInvoices al cambiar la referencia
}

function navigateToDetail(id) { router.push(`/sales/invoices/${id}`); }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: '2-digit', day: '2-digit' }) : '—'; }

function getStatusLabel(s) { return salesApi.getStatusLabel(s); }
function getStatusClass(s) { return salesApi.getStatusClass(s); }
</script>

<style scoped>
.order-ref { font-family: var(--font-family-mono); font-weight: 700; color: var(--color-secondary); }
.align-right { text-align: right; }

.type-tag {
  font-size: 0.7rem;
  font-weight: 800;
  text-transform: uppercase;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  border: 1px solid transparent;
}
.type-tag.standard { background: rgba(37, 99, 235, 0.1); color: #2563eb; border-color: rgba(37, 99, 235, 0.2); }
.type-tag.simplified { background: rgba(217, 119, 6, 0.1); color: #d97706; border-color: rgba(217, 119, 6, 0.2); }

.action-buttons { display: flex; justify-content: flex-end; }
.btn-icon { color: var(--color-text-secondary); transition: 0.2s; padding: 0.4rem; border-radius: 6px; border: none; background: transparent; cursor: pointer; }
.btn-icon:hover { color: var(--color-text-primary); background: rgba(0,0,0,0.05); }
</style>
