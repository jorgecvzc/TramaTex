<template>
  <div class="party-detail">
    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Cargando entidad...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="error-message">
      <span>✗ {{ error }}</span>
      <router-link to="/parties" class="btn btn-secondary">
        ← Volver a Entidades
      </router-link>
    </div>

    <!-- Party Detail -->
    <div v-if="!isLoading && party" class="detail-container">
      <!-- Header -->
      <div class="detail-header card">
        <div class="header-content">
          <h2>{{ party.name }}</h2>
          <div class="header-badges">
            <span class="badge" :class="`role-${party.role.toLowerCase()}`">
              {{ formatRole(party.role) }}
            </span>
            <span class="badge" :class="`status-${party.status.toLowerCase()}`">
              {{ party.status === 'ACTIVE' ? 'Activo' : 'Inactivo' }}
            </span>
          </div>
        </div>
      </div>

      <!-- Información de la parte -->
      <div class="info-section card">
        <h3>Información de la entidad</h3>
        
        <div class="info-grid">
          <div class="info-item">
            <label>Nombre</label>
            <p>{{ party.name }}</p>
          </div>
          
          <div class="info-item">
            <label>Rol</label>
            <p>{{ formatRole(party.role) }}</p>
          </div>
          
          <div class="info-item">
            <label>Estado</label>
            <p>
              {{ party.status }}
              <button
                @click="toggleStatus"
                :disabled="isUpdating"
                class="btn btn-small"
              >
                {{ party.status === 'ACTIVE' ? 'Desactivar' : 'Activar' }}
              </button>
            </p>
          </div>

          <div v-if="party.tax_id" class="info-item">
            <label>NIF/CIF</label>
            <p>{{ party.tax_id }}</p>
          </div>

          <div v-if="party.website" class="info-item">
            <label>Sitio web</label>
            <p>
              <a :href="party.website" target="_blank" rel="noopener">
                {{ party.website }}
              </a>
            </p>
          </div>

          <div class="info-item">
            <label>Creado</label>
            <p>{{ formatDate(party.created_at) }}</p>
          </div>

          <div v-if="party.modified_at" class="info-item">
            <label>Última modificación</label>
            <p>{{ formatDate(party.modified_at) }}</p>
          </div>
        </div>

        <!-- Edit Button -->
        <div class="info-actions">
          <button v-if="!isEditing" @click="isEditing = true" class="btn btn-primary">
            ✎ Editar entidad
          </button>

          <div v-if="isEditing" class="edit-form">
            <div class="form-group">
              <label for="editName">Nombre</label>
              <input
                id="editName"
                v-model="editForm.name"
                type="text"
              />
            </div>

            <div class="form-group">
              <label for="editWebsite">Sitio web</label>
              <input
                id="editWebsite"
                v-model="editForm.website"
                type="url"
              />
            </div>

            <div class="form-group">
              <label for="editNotes">Notas</label>
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
                {{ isUpdating ? 'Guardando...' : 'Guardar cambios' }}
              </button>
              <button
                @click="isEditing = false"
                class="btn btn-secondary"
              >
                Cancelar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Contacts Section -->
      <person-manager :party-id="partyId" />

      <!-- Addresses Section -->
      <address-manager :party-id="partyId" />
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

const partyId = route.params.id;

const party = ref(null);
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
  fetchParty();
});

async function fetchParty() {
  isLoading.value = true;
  error.value = '';

  try {
    const data = await partyApi.getParty(partyId);
    party.value = data;
    
    // Initialize edit form
    editForm.name = data.name;
    editForm.website = data.website || '';
    editForm.notes = data.notes || '';
  } catch (err) {
    error.value = err?.message || 'Entidad no encontrada';
  } finally {
    isLoading.value = false;
  }
}

async function submitEdit() {
  if (!editForm.name.trim()) {
    alert('El nombre es obligatorio');
    return;
  }

  isUpdating.value = true;

  try {
    const updated = await partyApi.updateParty(partyId, {
      name: editForm.name,
      website: editForm.website,
      notes: editForm.notes,
    });

    party.value = updated;
    isEditing.value = false;
  } catch (err) {
    error.value = err?.message || 'No se pudo actualizar la entidad';
  } finally {
    isUpdating.value = false;
  }
}

async function toggleStatus() {
  if (!party.value) return;

  isUpdating.value = true;

  try {
    const newStatus = party.value.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE';
    const updated = await partyApi.changePartyStatus(partyId, newStatus);
    party.value = updated;
  } catch (err) {
    error.value = err?.message || 'No se pudo actualizar el estado';
  } finally {
    isUpdating.value = false;
  }
}

function formatRole(role) {
  const map = {
    CLIENT: 'Cliente',
    SUPPLIER: 'Proveedor',
    BOTH: 'Ambos',
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
.party-detail {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
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

.error-message {
  background: #fee2e2;
  border: 1px solid #ef4444;
  border-radius: 8px;
  padding: 1rem;
  text-align: center;
  color: #991b1b;
}

.detail-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
}

.header-content h2 {
  color: #1b3a6b;
  margin: 0 0 1rem 0;
  font-size: 1.6rem;
}

.header-badges {
  display: flex;
  gap: 1rem;
}

.badge {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-weight: 600;
  font-size: 0.75rem;
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

.btn-primary:hover:not(:disabled) {
  background: #d6aa00;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover {
  background: #f8fafc;
}

.btn-small {
  padding: 0.4rem 0.75rem;
  font-size: 0.75rem;
  margin-left: 0.5rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.info-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-section h3 {
  color: #1b3a6b;
  margin: 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.info-item label {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin-bottom: 0.4rem;
}

.info-item p {
  color: #1e293b;
  margin: 0;
  font-size: 1rem;
}

.info-item a {
  color: #002395;
  text-decoration: none;
}

.info-item a:hover {
  text-decoration: underline;
}

.info-actions {
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

.edit-form {
  background: #f8fafc;
  padding: 1.5rem;
  border-radius: 10px;
  margin-top: 1rem;
  border: 1px solid #e2e8f0;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin-bottom: 0.4rem;
}

.form-group input,
.form-group textarea {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.6rem 0.8rem;
  font-size: 0.9rem;
  color: #1e293b;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
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
