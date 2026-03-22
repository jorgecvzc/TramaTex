<template>
  <Navbar />
  <div class="delivery-notes-list-container">
    <!-- Header -->
    <div class="page-header">
      <h1>Albaranes</h1>
    </div>

    <!-- Filters -->
    <div class="filters-card">
      <div class="filters-grid">
        <div class="filter-group">
          <label>Búsqueda</label>
          <input
            v-model="filters.searchText"
            type="text"
            placeholder="Buscar por referencia o nombre"
            class="filter-input"
          />
        </div>

        <div class="filter-group">
          <label>Estado</label>
          <select v-model="filters.status" class="filter-input">
            <option value="">Todos</option>
            <option value="PENDING">Pendiente</option>
            <option value="DELIVERED">Entregado</option>
            <option value="CANCELLED">Cancelado</option>
          </select>
        </div>

        <div class="filter-group">
          <label>Desde</label>
          <input
            v-model="filters.fromDate"
            type="date"
            class="filter-input"
          />
        </div>

        <div class="filter-group">
          <label>Hasta</label>
          <input
            v-model="filters.toDate"
            type="date"
            class="filter-input"
          />
        </div>
      </div>

      <div class="filters-actions">
        <div class="limit-group">
          <label>Mostrar</label>
          <select v-model.number="filters.limit" class="filter-select limit-select" @change="applyFilters">
            <option :value="25">25</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
            <option :value="0">Todos</option>
          </select>
          <span class="limit-label">registros</span>
        </div>
        <button class="btn btn-secondary" @click="clearFilters" v-if="hasFilters">
          Limpiar Filtros
        </button>
        <button class="btn btn-primary" @click="applyFilters" :disabled="!isDateRangeValid" :title="!isDateRangeValid ? 'Completa ambas fechas o vacía ambas para buscar' : ''">
          Buscar
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando albaranes...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchDeliveryNotes">Reintentar</button>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredDeliveryNotes.length === 0" class="empty-state">
      <p>No se encontraron albaranes</p>
      <p class="hint">Los albaranes se crean desde los pedidos</p>
    </div>

    <!-- Delivery Notes Table -->
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>Número</th>
            <th>Cliente</th>
            <th>Pedido</th>
            <th>Fecha Entrega</th>
            <th>Estado</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="note in filteredDeliveryNotes" :key="note.id" @click="navigateToDetail(note.id)" class="clickable-row">
            <td class="note-number">{{ note.deliveryNoteNumber }}</td>
            <td>{{ partiesCache[note.partyId] || '—' }}</td>
            <td>{{ ordersCache[note.salesOrderId] || formatOrderId(note.salesOrderId) }}</td>
            <td>{{ formatDate(note.deliveryDate) }}</td>
            <td><span :class="['status-badge', 'status-' + salesApi.getStatusClass(note.status)]">{{ salesApi.getStatusLabel(note.status) }}</span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';

const router = useRouter();

const deliveryNotes = ref([]);
const isLoading = ref(false);
const error = ref('');
const partiesCache = ref({});
const ordersCache = ref({});

const filters = ref({
  searchText: '',
  status: '',
  fromDate: '',
  toDate: '',
  limit: 50,
});

const hasFilters = computed(() => {
  return filters.value.searchText.trim() !== '' || 
         filters.value.status !== '' ||
         filters.value.fromDate !== '' || 
         filters.value.toDate !== '';
});

const isDateRangeValid = computed(() => {
  const hasFromDate = Boolean(filters.value.fromDate);
  const hasToDate = Boolean(filters.value.toDate);
  return hasFromDate === hasToDate;
});

const filteredDeliveryNotes = computed(() => {
  return deliveryNotes.value;
});

let searchDebounceTimer = null;
let autoFetchEnabled = false;

function scheduleDeliveryNotesFetch() {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }

  searchDebounceTimer = setTimeout(() => {
    fetchDeliveryNotes();
  }, 300);
}

watch(
  () => filters.value.searchText,
  (newSearch, oldSearch) => {
    const normalizedNew = (newSearch || '').trim();
    const normalizedOld = (oldSearch || '').trim();
    if (normalizedNew === normalizedOld) return;

    if (!autoFetchEnabled) return;

    scheduleDeliveryNotesFetch();
  },
);

watch(
  () => filters.value.status,
  (newStatus, oldStatus) => {
    if (!autoFetchEnabled) return;
    if (newStatus === oldStatus) return;
    scheduleDeliveryNotesFetch();
  },
);

watch(
  () => [filters.value.fromDate, filters.value.toDate],
  (newValues, oldValues) => {
    if (!autoFetchEnabled || !oldValues) return;

    const [newFromDate, newToDate] = newValues;
    const [oldFromDate, oldToDate] = oldValues;
    if (newFromDate === oldFromDate && newToDate === oldToDate) return;

    const hasFromDate = Boolean(newFromDate);
    const hasToDate = Boolean(newToDate);
    if (hasFromDate !== hasToDate) return;

    scheduleDeliveryNotesFetch();
  },
);

onBeforeUnmount(() => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }
});

onMounted(() => {
  autoFetchEnabled = true;
  
  fetchDeliveryNotes();
});

async function fetchDeliveryNotes() {
  isLoading.value = true;
  error.value = '';

  try {
    const apiFilters = {};
    
    if (filters.value.searchText) apiFilters.searchText = filters.value.searchText;
    if (filters.value.status) apiFilters.status = filters.value.status;
    if (filters.value.fromDate) apiFilters.fromDate = filters.value.fromDate;
    if (filters.value.toDate) apiFilters.toDate = filters.value.toDate;
    if (filters.value.limit) apiFilters.limit = filters.value.limit;

    const response = await salesApi.listDeliveryNotes(apiFilters);
    deliveryNotes.value = Array.isArray(response) ? response : (response.data || []);
    await Promise.all([loadPartyNames(), loadOrderNumbers()]);
  } catch (err) {
    error.value = err?.message || 'No se pudieron cargar los albaranes';
    console.error('Error loading delivery notes:', err);
  } finally {
    isLoading.value = false;
  }
}

function applyFilters() {
  if (!isDateRangeValid.value) return;
  fetchDeliveryNotes();
}

function clearFilters() {
  filters.value.searchText = '';
  filters.value.status = '';
  filters.value.fromDate = '';
  filters.value.toDate = '';
  filters.value.limit = 50;
  fetchDeliveryNotes();
}

async function loadPartyNames() {
  const partyIds = [...new Set(deliveryNotes.value.map((note) => note.partyId).filter(Boolean))];
  const uncachedIds = partyIds.filter((id) => !partiesCache.value[id]);

  if (uncachedIds.length === 0) {
    return;
  }

  try {
    const partiesMap = await partyApi.getPartiesBatch(uncachedIds);
    for (const partyId of uncachedIds) {
      partiesCache.value[partyId] = partiesMap[partyId]?.name || 'Sin nombre';
    }
  } catch (loadError) {
    console.error('Error loading delivery note party names:', loadError);
  }
}

async function loadOrderNumbers() {
  const orderIds = [...new Set(deliveryNotes.value.map((note) => note.salesOrderId).filter(Boolean))];
  const uncachedIds = orderIds.filter((id) => !ordersCache.value[id]);

  if (uncachedIds.length === 0) {
    return;
  }

  const results = await Promise.allSettled(uncachedIds.map((id) => salesApi.getOrder(id)));
  results.forEach((result, index) => {
    const orderId = uncachedIds[index];
    if (result.status === 'fulfilled') {
      ordersCache.value[orderId] = result.value.orderNumber || formatOrderId(orderId);
    }
  });
}

function navigateToDetail(noteId) {
  router.push(`/sales/delivery-notes/${noteId}`);
}

function formatDate(dateString) {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return date.toLocaleDateString('es-ES', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

function formatOrderId(orderId) {
  if (!orderId) return '—';
  return orderId.substring(0, 8) + '...';
}
</script>

<style scoped>
.delivery-notes-list-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.filters-card {
  background: white;
  border-radius: 8px;
  padding: 1rem 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  display: flex;
  flex-wrap: nowrap;
  gap: 0.5rem;
  align-items: flex-end;
}

.filters-grid {
  display: contents;
}

.filter-group {
  display: flex;
  flex-direction: column;
  flex: 1 1 0;
  min-width: 0;
}

.filter-group label {
  font-size: 0.8rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.25rem;
}

.filter-input,
.filter-select {
  padding: 0.4rem 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.85rem;
}

.filters-actions {
  display: contents;
}

.limit-group {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.limit-group label {
  font-size: 0.8rem;
  color: #6b7280;
  white-space: nowrap;
}

.limit-select {
  width: 70px;
}

.limit-label {
  font-size: 0.8rem;
  color: #6b7280;
  white-space: nowrap;
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 3rem;
  background: white;
  border-radius: 8px;
}

.hint {
  font-size: 0.875rem;
  color: #6b7280;
  margin-top: 0.5rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top-color: #E6B800;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.table-container {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table thead {
  background: #f9fafb;
}

.data-table th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  font-size: 0.875rem;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

.data-table td {
  padding: 1rem;
  border-bottom: 1px solid #f3f4f6;
  font-size: 0.875rem;
}

.clickable-row {
  cursor: pointer;
  transition: background 0.2s;
}

.clickable-row:hover {
  background: #f9fafb;
}

.note-number {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #1e40af;
}

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn-icon:hover {
  background: #f3f4f6;
}

.btn {
  padding: 0.625rem 1.25rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #E6B800;
  color: white;
}

.btn-primary:hover {
  background: #d4a700;
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-warning {
  background: #fef3c7;
  color: #92400e;
}

.status-success {
  background: #d1fae5;
  color: #065f46;
}

.status-danger {
  background: #fee2e2;
  color: #991b1b;
}

.status-primary {
  background: #dbeafe;
  color: #1e40af;
}
</style>
