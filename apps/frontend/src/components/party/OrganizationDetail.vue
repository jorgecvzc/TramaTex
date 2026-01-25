<template>
  <div class="organization-detail">
    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Loading organization...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="error-message">
      <span>✗ {{ error }}</span>
      <router-link to="/organizations" class="btn btn-secondary">
        ← Back to Organizations
      </router-link>
    </div>

    <!-- Organization Detail -->
    <div v-if="!isLoading && organization" class="detail-container">
      <!-- Header -->
      <div class="detail-header">
        <div class="header-content">
          <h1>{{ organization.name }}</h1>
          <div class="header-badges">
            <span class="badge" :class="`role-${organization.role.toLowerCase()}`">
              {{ formatRole(organization.role) }}
            </span>
            <span class="badge" :class="`status-${organization.status.toLowerCase()}`">
              {{ organization.status }}
            </span>
          </div>
        </div>
        <router-link to="/organizations" class="btn btn-secondary">
          ← Back
        </router-link>
      </div>

      <!-- Organization Info -->
      <div class="info-section">
        <h2>Organization Information</h2>
        
        <div class="info-grid">
          <div class="info-item">
            <label>Organization Name</label>
            <p>{{ organization.name }}</p>
          </div>
          
          <div class="info-item">
            <label>Role</label>
            <p>{{ formatRole(organization.role) }}</p>
          </div>
          
          <div class="info-item">
            <label>Status</label>
            <p>
              {{ organization.status }}
              <button
                @click="toggleStatus"
                :disabled="isUpdating"
                class="btn btn-small"
              >
                {{ organization.status === 'ACTIVE' ? 'Deactivate' : 'Activate' }}
              </button>
            </p>
          </div>

          <div v-if="organization.tax_id" class="info-item">
            <label>Tax ID</label>
            <p>{{ organization.tax_id }}</p>
          </div>

          <div v-if="organization.website" class="info-item">
            <label>Website</label>
            <p>
              <a :href="organization.website" target="_blank" rel="noopener">
                {{ organization.website }}
              </a>
            </p>
          </div>

          <div class="info-item">
            <label>Created</label>
            <p>{{ formatDate(organization.created_at) }}</p>
          </div>

          <div v-if="organization.modified_at" class="info-item">
            <label>Last Modified</label>
            <p>{{ formatDate(organization.modified_at) }}</p>
          </div>
        </div>

        <!-- Edit Button -->
        <div class="info-actions">
          <button v-if="!isEditing" @click="isEditing = true" class="btn btn-primary">
            ✎ Edit Organization
          </button>

          <div v-if="isEditing" class="edit-form">
            <div class="form-group">
              <label for="editName">Name</label>
              <input
                id="editName"
                v-model="editForm.name"
                type="text"
              />
            </div>

            <div class="form-group">
              <label for="editWebsite">Website</label>
              <input
                id="editWebsite"
                v-model="editForm.website"
                type="url"
              />
            </div>

            <div class="form-group">
              <label for="editNotes">Notes</label>
              <textarea
                id="editNotes"
                v-model="editForm.notes"
                rows="3"
              />
            </div>

            <div class="edit-actions">
              <button
                @click="submitEdit"
                :disabled="isUpdating"
                class="btn btn-primary"
              >
                {{ isUpdating ? 'Saving...' : 'Save Changes' }}
              </button>
              <button
                @click="isEditing = false"
                class="btn btn-secondary"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Contacts Section -->
      <person-manager :organization-id="organizationId" />

      <!-- Addresses Section -->
      <address-manager :organization-id="organizationId" />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { partyApi } from '@/services/partyApi';
import PersonManager from './PersonManager.vue';
import AddressManager from './AddressManager.vue';

const route = useRoute();
const router = useRouter();

const organizationId = route.params.id;

const organization = ref(null);
const isLoading = ref(false);
const isUpdating = ref(false);
const error = ref('');
const isEditing = ref(false);

const editForm = reactive({
  name: '',
  website: '',
  notes: '',
});

onMounted(() => {
  fetchOrganization();
});

async function fetchOrganization() {
  isLoading.value = true;
  error.value = '';

  try {
    const data = await partyApi.getOrganization(organizationId);
    organization.value = data;
    
    // Initialize edit form
    editForm.name = data.name;
    editForm.website = data.website || '';
    editForm.notes = data.notes || '';
  } catch (err) {
    error.value = err.message || 'Organization not found';
  } finally {
    isLoading.value = false;
  }
}

async function submitEdit() {
  if (!editForm.name.trim()) {
    alert('Organization name is required');
    return;
  }

  isUpdating.value = true;

  try {
    const updated = await partyApi.updateOrganization(organizationId, {
      name: editForm.name,
      website: editForm.website,
      notes: editForm.notes,
    });

    organization.value = updated;
    isEditing.value = false;
  } catch (err) {
    error.value = err.message || 'Failed to update organization';
  } finally {
    isUpdating.value = false;
  }
}

async function toggleStatus() {
  if (!organization.value) return;

  isUpdating.value = true;

  try {
    const newStatus = organization.value.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE';
    const updated = await partyApi.changeOrganizationStatus(organizationId, newStatus);
    organization.value = updated;
  } catch (err) {
    error.value = err.message || 'Failed to update status';
  } finally {
    isUpdating.value = false;
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
    month: 'long',
    day: 'numeric',
  });
}
</script>

<style scoped>
.organization-detail {
  padding: 1.5rem;
  background: var(--color-background);
  min-height: 100vh;
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
  background: rgba(244, 67, 54, 0.1);
  border: 1px solid var(--color-error);
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  color: var(--color-error);
}

.detail-container {
  max-width: 1000px;
  margin: 0 auto;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  background: var(--color-surface);
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content h1 {
  color: var(--color-text-primary);
  margin: 0 0 1rem 0;
  font-size: 2rem;
}

.header-badges {
  display: flex;
  gap: 1rem;
}

.badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 12px;
  font-weight: 500;
  font-size: 0.9rem;
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

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  display: inline-block;
}

.btn-primary {
  background-color: var(--primary-color);
  color: var(--color-text-on-primary);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--primary-color-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(230, 184, 0, 0.3);
}

.btn-secondary {
  background-color: var(--color-secondary);
  color: var(--color-text-primary);
}

.btn-secondary:hover {
  background-color: var(--color-border);
}

.btn-small {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  margin-left: 0.5rem;
}

.info-section {
  background: var(--color-surface);
  padding: 2rem;
  border-radius: 8px;
  margin-bottom: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.info-section h2 {
  color: var(--color-text-primary);
  margin-top: 0;
  padding-bottom: 1rem;
  border-bottom: 2px solid var(--color-border);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.info-item label {
  display: block;
  font-weight: 600;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
  text-transform: uppercase;
}

.info-item p {
  color: var(--color-text-primary);
  margin: 0;
  font-size: 1.1rem;
}

.info-item a {
  color: var(--primary-color);
  text-decoration: none;
}

.info-item a:hover {
  text-decoration: underline;
}

.info-actions {
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.edit-form {
  background: rgba(230, 184, 0, 0.05);
  padding: 1.5rem;
  border-radius: 6px;
  margin-top: 1rem;
  border: 1px solid rgba(230, 184, 0, 0.2);
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.edit-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .detail-header {
    flex-direction: column;
    gap: 1rem;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .edit-actions {
    flex-direction: column;
  }
}
</style>
