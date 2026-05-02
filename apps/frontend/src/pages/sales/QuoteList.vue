<template>
  <div class="page-layout">
    <BaseCatalog
      title="Presupuestos de Venta"
      icon="description"
      :breadcrumbs="[{ label: 'Ventas', to: '/sales/quotes' }, { label: 'Presupuestos' }]"
      :items="quotes"
      :is-loading="isLoading"
      :error="error"
      :has-filters="hasFilters"
      create-route="/sales/quotes/new"
      create-text="Nuevo Presupuesto"
      empty-icon="description"
      empty-text="No hay presupuestos registrados"
      @clear-filters="clearFilters"
      @refresh="fetchQuotes"
      @click-item="(item) => navigateToDetail(item.id)"
    >
      <!-- CAPA 2: CONTEXTO (Filtros) -->
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
            <option value="BORRADOR">Borrador</option>
            <option value="EMITIDA">Emitida</option>
            <option value="APROBADA">Aprobada</option>
            <option value="RECHAZADA">Rechazada</option>
            <option value="EXPIRADA">Expirada</option>
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

      <!-- CAPA 3: TRABAJO (Tabla) -->
      <template #table-header>
        <th>Número</th>
        <th>Cliente</th>
        <th>Fecha</th>
        <th>Vencimiento</th>
        <th>Estado</th>
        <th class="align-right">Importe Total</th>
        <th>MES</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><span class="order-ref">{{ item.quoteNumber }}</span></td>
        <td>{{ formatPartyId(item.partyId) }}</td>
        <td>{{ formatDate(item.quoteDate) }}</td>
        <td :class="{'text-danger fw-bold': isQuoteExpired(item.expirationDate)}">{{ formatDate(item.expirationDate) }}</td>
        <td>
          <span :class="['status-badge', `status-${getStatusClass(item.status)}`]">
            {{ getStatusLabel(item.status) }}
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
            <button class="btn-icon" @click="navigateToDetail(item.id)" title="Ver detalle">
              <Eye :size="18" />
            </button>
          </div>
        </td>
      </template>
    </BaseCatalog>
  </div>
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
const quotes = ref([]);
const isLoading = ref(false);
const error = ref('');
const partiesCache = ref({});
const mesWorksCache = ref({});

const filters = ref({ partyId: '', searchText: '', status: '', fromDate: '', toDate: '', limit: 50 });

const hasFilters = computed(() => (filters.value.partyId || filters.value.searchText || filters.value.status || filters.value.fromDate));
const isDateRangeValid = computed(() => !filters.value.fromDate || !filters.value.toDate || filters.value.fromDate <= filters.value.toDate);

// Auto-fetch logic con limpieza robusta
let searchDebounceTimer = null;
watch(() => ({ ...filters.value }), (newFilters, oldFilters) => {
  if (JSON.stringify(newFilters) === JSON.stringify(oldFilters)) return;
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
  searchDebounceTimer = setTimeout(() => fetchQuotes(), 350);
}, { deep: true });

onMounted(async () => {
  fetchQuotes();
  loadPartyNames().catch(console.error);
  loadMesWorksForQuotes().catch(console.error);
});
onBeforeUnmount(() => { if (searchDebounceTimer) clearTimeout(searchDebounceTimer); });

async function fetchQuotes() {
  isLoading.value = true;
  error.value = '';
  try {
    const res = await salesApi.listQuotes(filters.value);
    if (Array.isArray(res)) {
      quotes.value = res;
    } else if (res && res.data && Array.isArray(res.data)) {
      quotes.value = res.data;
    } else if (res && typeof res === 'object') {
      const possibleArray = Object.values(res).find(v => Array.isArray(v));
      quotes.value = possibleArray || [];
    } else {
      quotes.value = [];
    }
    
    if (quotes.value.length > 0) {
      await loadPartyNames();
      await loadMesWorksForQuotes();
    }
  } catch (err) { 
    console.error('Error fetchQuotes:', err);
    error.value = err.message; 
  }
  finally { isLoading.value = false; }
}

async function loadPartyNames() {
  const ids = [...new Set(quotes.value.map(q => q.partyId).filter(id => id && !partiesCache.value[id]))];
  if (ids.length === 0) return;
  try {
    const map = await partyApi.getPartiesBatch(ids);
    ids.forEach(id => partiesCache.value[id] = map[id]?.name || 'Desconocido');
  } catch (err) {}
}

async function loadMesWorksForQuotes() {
  const ids = [...new Set(quotes.value.flatMap(q => (q.lineItems || []).map(i => i.mesWorkId)).filter(id => id && !mesWorksCache.value[id]))];
  if (ids.length === 0) return;
  const results = await Promise.allSettled(ids.map(id => mesApi.getWorkOrder(id)));
  results.forEach((r, i) => mesWorksCache.value[ids[i]] = r.status === 'fulfilled' ? r.value : null);
}

function clearFilters() { filters.value = { partyId: '', searchText: '', status: '', fromDate: '', toDate: '', limit: 50 }; fetchQuotes(); }
function navigateToDetail(id) { router.push(`/sales/quotes/${id}`); }

function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: '2-digit', day: '2-digit' }) : '—'; }
function formatPartyId(id) { return partiesCache.value[id] || '...'; }
function isQuoteExpired(d) { return d && new Date(d) < new Date(); }

function getMesSummary(items) {
  const ids = getMesWorkIds(items); if (ids.length === 0) return '—';
  const first = mesWorksCache.value[ids[0]]?.work_number || 'MES';
  return ids.length === 1 ? first : `${first} +${ids.length-1}`;
}
function getMesWorkIds(items) { return [...new Set((items || []).map(i => i.mesWorkId).filter(id => !!id))]; }

function getStatusLabel(s) { return salesApi.getStatusLabel(s); }
function getStatusClass(s) { return salesApi.getStatusClass(s); }
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.order-ref { font-family: var(--font-family-mono); font-weight: 700; color: var(--color-secondary); }
.align-right { text-align: right; }

.mes-badge { display: inline-flex; align-items: center; gap: 0.35rem; padding: 0.25rem 0.5rem; background: rgba(0, 35, 149, 0.05); border-radius: 4px; color: var(--color-secondary); font-size: 0.75rem; font-weight: 600; }

.action-buttons { display: flex; justify-content: flex-end; }
.btn-icon { color: var(--color-text-secondary); transition: 0.2s; padding: 0.4rem; border-radius: 6px; border: none; background: transparent; cursor: pointer; }
.btn-icon:hover { color: var(--color-text-primary); background: rgba(0,0,0,0.05); }
</style>