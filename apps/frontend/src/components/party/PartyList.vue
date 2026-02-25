<template>
  <div class="party-list">
    <!-- Filters and Search -->
    <div class="filters">
      <div>
        <label>Buscar por nombre</label>
        <input
          v-model="filters.name"
          type="text"
          placeholder="Buscar por nombre"
          @input="applyFilters"
        />
      </div>

      <div>
        <label>Filtrar por rol</label>
        <select v-model="filters.role" @change="applyFilters">
          <option value="">Todos</option>
          <option value="CLIENT">Clientes</option>
          <option value="SUPPLIER">Proveedores</option>
          <option value="BOTH">Ambos</option>
          <option value="CONTACT">Contactos</option>
        </select>
      </div>

      <div>
        <label>Filtrar por estado</label>
        <select v-model="filters.status" @change="applyFilters">
          <option value="">Todos</option>
          <option value="ACTIVE">Activo</option>
          <option value="INACTIVE">Inactivo</option>
        </select>
      </div>

      <div class="filter-actions">
        <button @click="clearFilters" class="btn btn-secondary">
          Limpiar filtros
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Cargando entidades...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="alert-error">
      {{ error }}
    </div>

    <!-- Parties Table -->
    <div v-if="!isLoading && parties.length > 0" class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Nombre</th>
            <th>Rol</th>
            <th>Estado</th>
            <th>NIF/CIF</th>
            <th>Creado</th>
            <th class="align-right">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="party in parties" :key="party.id" :class="{ inactive: party.status === 'INACTIVE' }">
            <td>
              <router-link :to="`/parties/${party.id}`" class="party-link">
                {{ party.name }}
              </router-link>
            </td>
            <td>
              <span class="role-pill">{{ formatRole(party.role) }}</span>
            </td>
            <td>
              <span class="status-pill" :class="`status-${party.status.toLowerCase()}`">
                {{ party.status === 'ACTIVE' ? 'Activo' : 'Inactivo' }}
              </span>
            </td>
            <td>{{ party.tax_id || '—' }}</td>
            <td>{{ formatDate(party.created_at) }}</td>
            <td class="align-right">
              <div class="action-buttons">
                <router-link :to="`/parties/${party.id}`" class="btn btn-outline">
                  Ver detalles
                </router-link>
                <button class="btn btn-secondary" @click="toggleStatus(party)">
                  {{ party.status === 'ACTIVE' ? 'Desactivar' : 'Activar' }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="!isLoading && parties.length === 0">
            <td colspan="6" class="empty-state">No hay entidades para mostrar.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="!isLoading && parties.length === 0" class="empty-state-block">
      <p>
        {{
          filters.name || filters.role || filters.status
            ? 'Prueba ajustando los filtros.'
            : 'Crea tu primera entidad para empezar.'
        }}
      </p>
      <router-link v-if="!hasFilters" to="/parties/new" class="btn btn-primary">
        Crear entidad
      </router-link>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        :disabled="currentPage === 1"
        @click="previousPage"
        class="btn btn-secondary"
      >
        ← Anterior
      </button>
      <span class="page-info">Página {{ currentPage }} de {{ totalPages }}</span>
      <button
        :disabled="currentPage === totalPages"
        @click="nextPage"
        class="btn btn-secondary"
      >
        Siguiente →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { partyApi } from '@/services/partyApi';

const parties = ref([]);
const isLoading = ref(false);
const error = ref('');
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

const filters = reactive({
  name: '',
  role: '',
  status: '',
});

const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

const hasFilters = computed(
  () => filters.name.trim() !== '' || filters.role !== '' || filters.status !== ''
);

onMounted(() => {
  fetchParties();
});

async function fetchParties() {
  isLoading.value = true;
  error.value = '';

  try {
    const response = await partyApi.listParties({
      name: filters.name,
      role: filters.role,
      status: filters.status,
      pageNumber: currentPage.value,
      pageSize: pageSize.value,
    });

    parties.value = response.data || [];
    total.value = response.total || 0;
  } catch (err) {
    error.value = err?.message || 'No se pudieron cargar las entidades';
  } finally {
    isLoading.value = false;
  }
}

function applyFilters() {
  currentPage.value = 1;
  fetchParties();
}

function clearFilters() {
  filters.name = '';
  filters.role = '';
  filters.status = '';
  currentPage.value = 1;
  fetchParties();
}

function previousPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchParties();
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchParties();
  }
}

async function toggleStatus(party) {
  const newStatus = party.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE';
  isLoading.value = true;
  error.value = '';

  try {
    await partyApi.changePartyStatus(party.id, newStatus);
    party.status = newStatus;
  } catch (err) {
    error.value = err?.message || `No se pudo cambiar el estado`;
  } finally {
    isLoading.value = false;
  }
}

function formatRole(role) {
  const map = {
    CLIENT: 'Cliente',
    SUPPLIER: 'Proveedor',
    BOTH: 'Ambos',
    CONTACT: 'Contacto',
  };
  return map[role] || role;
}

function formatDate(dateString) {
  if (!dateString) return '—';
  return new Date(dateString).toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}
</script>

<style scoped>
.party-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.filters > div {
  flex: 1;
  min-width: 220px;
}

.filter-actions {
  display: flex;
  align-items: flex-end;
}

label {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin-bottom: 0.4rem;
}

input,
select {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.6rem 0.8rem;
  font-size: 0.9rem;
  color: #1e293b;
}

input:focus,
select:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  gap: 1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(0, 35, 149, 0.12);
  border-top-color: #002395;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.alert-error {
  background: #fee2e2;
  border: 1px solid #ef4444;
  color: #991b1b;
  padding: 0.8rem 1rem;
  border-radius: 8px;
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

thead {
  background: #f8fafc;
  color: #64748b;
}

th {
  padding: 0.85rem 1rem;
  text-align: left;
  font-weight: 600;
  border-bottom: 1px solid #e2e8f0;
}

tbody tr {
  border-bottom: 1px solid #e2e8f0;
}

tbody tr:hover {
  background-color: #f8fafc;
}

tbody tr.inactive {
  opacity: 0.6;
}

td {
  padding: 0.85rem 1rem;
  color: #1e293b;
}

.party-link {
  color: #002395;
  text-decoration: none;
  font-weight: 600;
}

.party-link:hover {
  text-decoration: underline;
}

.role-pill,
.status-pill {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.role-pill {
  background-color: #e2e8f0;
  color: #1e293b;
}

.status-pill.status-active {
  background-color: rgba(76, 175, 80, 0.1);
  color: #4caf50;
}

.status-pill.status-inactive {
  background-color: rgba(158, 158, 158, 0.1);
  color: #9e9e9e;
}

.align-right {
  text-align: right;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.empty-state {
  text-align: center;
  color: #64748b;
  padding: 1.5rem;
}

.empty-state-block {
  text-align: center;
  color: #64748b;
  padding: 1.5rem;
  border: 1px dashed #e2e8f0;
  border-radius: 10px;
  display: grid;
  gap: 0.75rem;
  justify-items: center;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}

.page-info {
  color: #64748b;
  font-weight: 500;
  min-width: 120px;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.btn-primary {
  background: #e6b800;
  color: #1e293b;
  font-weight: 700;
}

.btn-primary:hover {
  background: #d6aa00;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-outline {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
  }

  table {
    font-size: 0.875rem;
  }

  th,
  td {
    padding: 0.75rem 0.5rem;
  }

  .action-buttons {
    flex-direction: column;
  }
}
</style>
