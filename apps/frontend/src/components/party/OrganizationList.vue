<template>
  <div class="organization-list">
    <!-- Header with Actions -->
    <div class="list-header">
      <h2>Organizations</h2>
      <router-link to="/organizations/new" class="btn btn-primary">
        + Create Organization
      </router-link>
    </div>

    <!-- Filters and Search -->
    <div class="filters">
      <div class="filter-group">
        <input
          v-model="filters.name"
          type="text"
          placeholder="Search by name..."
          class="search-input"
          @input="applyFilters"
        />
      </div>

      <div class="filter-group">
        <select v-model="filters.role" @change="applyFilters">
          <option value="">All Roles</option>
          <option value="CLIENT">Clients</option>
          <option value="SUPPLIER">Suppliers</option>
          <option value="BOTH">Both</option>
        </select>
      </div>

      <div class="filter-group">
        <select v-model="filters.status" @change="applyFilters">
          <option value="">All Status</option>
          <option value="ACTIVE">Active</option>
          <option value="INACTIVE">Inactive</option>
        </select>
      </div>

      <button @click="clearFilters" class="btn btn-secondary">
        Clear Filters
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Loading organizations...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="error-message">
      <span>✗ {{ error }}</span>
      <button @click="error = ''" class="close">&times;</button>
    </div>

    <!-- Organizations Table -->
    <div v-if="!isLoading && organizations.length > 0" class="organizations-table">
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Role</th>
            <th>Status</th>
            <th>Tax ID</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="org in organizations" :key="org.id" :class="{ inactive: org.status === 'INACTIVE' }">
            <td>
              <router-link :to="`/organizations/${org.id}`" class="org-link">
                {{ org.name }}
              </router-link>
            </td>
            <td>
              <span class="badge" :class="`role-${org.role.toLowerCase()}`">
                {{ formatRole(org.role) }}
              </span>
            </td>
            <td>
              <span class="badge" :class="`status-${org.status.toLowerCase()}`">
                {{ org.status }}
              </span>
            </td>
            <td>{{ org.tax_id || '—' }}</td>
            <td>{{ formatDate(org.created_at) }}</td>
            <td class="actions">
              <router-link
                :to="`/organizations/${org.id}`"
                class="action-btn view"
                title="View Details"
              >
                👁️
              </router-link>
              <button
                @click="toggleStatus(org)"
                :title="org.status === 'ACTIVE' ? 'Deactivate' : 'Activate'"
                class="action-btn toggle"
              >
                {{ org.status === 'ACTIVE' ? '✓' : '○' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="!isLoading && organizations.length === 0" class="empty-state">
      <div class="empty-icon">📦</div>
      <h3>No organizations found</h3>
      <p>{{ filters.name || filters.role || filters.status ? 'Try adjusting your filters' : 'Create your first organization to get started' }}</p>
      <router-link v-if="!hasFilters" to="/organizations/new" class="btn btn-primary">
        Create Organization
      </router-link>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        :disabled="currentPage === 1"
        @click="previousPage"
        class="btn btn-secondary"
      >
        ← Previous
      </button>
      <span class="page-info">Page {{ currentPage }} of {{ totalPages }}</span>
      <button
        :disabled="currentPage === totalPages"
        @click="nextPage"
        class="btn btn-secondary"
      >
        Next →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { partyApi } from '@/services/partyApi';

const router = useRouter();

const organizations = ref([]);
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
  fetchOrganizations();
});

async function fetchOrganizations() {
  isLoading.value = true;
  error.value = '';

  try {
    const response = await partyApi.listOrganizations({
      name: filters.name,
      role: filters.role,
      status: filters.status,
      pageNumber: currentPage.value,
      pageSize: pageSize.value,
    });

    organizations.value = response.data || [];
    total.value = response.total || 0;
  } catch (err) {
    error.value = err.message || 'Failed to load organizations';
  } finally {
    isLoading.value = false;
  }
}

function applyFilters() {
  currentPage.value = 1;
  fetchOrganizations();
}

function clearFilters() {
  filters.name = '';
  filters.role = '';
  filters.status = '';
  currentPage.value = 1;
  fetchOrganizations();
}

function previousPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchOrganizations();
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchOrganizations();
  }
}

async function toggleStatus(org) {
  const newStatus = org.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE';
  isLoading.value = true;
  error.value = '';

  try {
    await partyApi.changeOrganizationStatus(org.id, newStatus);
    org.status = newStatus;
  } catch (err) {
    error.value = err.message || `Failed to change status`;
  } finally {
    isLoading.value = false;
  }
}

function formatRole(role) {
  const map = {
    CLIENT: 'Client',
    SUPPLIER: 'Supplier',
    BOTH: 'Both',
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
.organization-list {
  padding: 1.5rem;
  background: var(--color-background);
  min-height: 100vh;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.list-header h2 {
  color: var(--color-text-primary);
  font-size: 2rem;
  margin: 0;
}

.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 1rem;
  background: var(--color-surface);
  border-radius: 8px;
}

.filter-group {
  display: flex;
  flex-direction: column;
}

.filter-group input,
.filter-group select {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 0.95rem;
}

.filter-group input:focus,
.filter-group select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.search-input {
  width: 100%;
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
  border: 4px solid rgba(230, 184, 0, 0.2);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-message {
  padding: 1rem;
  background-color: rgba(244, 67, 54, 0.1);
  color: var(--color-error);
  border-left: 4px solid var(--color-error);
  border-radius: 4px;
  margin-bottom: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.error-message .close {
  background: none;
  border: none;
  color: inherit;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
}

.organizations-table {
  background: var(--color-surface);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: var(--color-background);
  border-bottom: 2px solid var(--color-border);
}

th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: var(--color-text-primary);
}

tbody tr {
  border-bottom: 1px solid var(--color-border);
  transition: background-color 0.2s ease;
}

tbody tr:hover {
  background-color: rgba(230, 184, 0, 0.05);
}

tbody tr.inactive {
  opacity: 0.6;
}

td {
  padding: 1rem;
  color: var(--color-text-primary);
}

.org-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.org-link:hover {
  text-decoration: underline;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 500;
}

.badge.role-client {
  background-color: rgba(33, 150, 243, 0.1);
  color: #2196f3;
}

.badge.role-supplier {
  background-color: rgba(76, 175, 80, 0.1);
  color: #4caf50;
}

.badge.role-both {
  background-color: rgba(230, 184, 0, 0.1);
  color: var(--primary-color);
}

.badge.status-active {
  background-color: rgba(76, 175, 80, 0.1);
  color: #4caf50;
}

.badge.status-inactive {
  background-color: rgba(158, 158, 158, 0.1);
  color: #9e9e9e;
}

.actions {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 4px;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background-color: rgba(230, 184, 0, 0.1);
}

.action-btn.view {
  text-decoration: none;
}

.empty-state {
  text-align: center;
  padding: 3rem 1.5rem;
  background: var(--color-surface);
  border-radius: 8px;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.empty-state h3 {
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
}

.empty-state p {
  color: var(--color-text-secondary);
  margin-bottom: 1.5rem;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
}

.page-info {
  color: var(--color-text-secondary);
  font-weight: 500;
  min-width: 120px;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary {
  background-color: var(--primary-color);
  color: var(--color-text-on-primary);
}

.btn-primary:hover {
  background-color: var(--primary-color-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(230, 184, 0, 0.3);
}

.btn-secondary {
  background-color: var(--color-secondary);
  color: var(--color-text-primary);
}

.btn-secondary:hover:not(:disabled) {
  background-color: var(--color-border);
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .list-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .filters {
    grid-template-columns: 1fr;
  }

  table {
    font-size: 0.875rem;
  }

  th,
  td {
    padding: 0.75rem 0.5rem;
  }

  .actions {
    flex-direction: column;
  }

  .action-btn {
    width: 28px;
    height: 28px;
    font-size: 0.875rem;
  }
}
</style>
