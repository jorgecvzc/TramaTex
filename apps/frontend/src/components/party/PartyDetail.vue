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
          
          <div v-if="party.phone" class="info-item">
            <label>Teléfono</label>
            <p>{{ party.phone }}</p>
          </div>
          
          <div v-if="party.email" class="info-item">
            <label>Email</label>
            <p>
              <a :href="`mailto:${party.email}`">
                {{ party.email }}
              </a>
            </p>
          </div>

          <div v-if="party.role === 'CLIENT' || party.role === 'BOTH'" class="info-item">
            <label>Bonificación por defecto</label>
            <p>{{ party.default_discount_percentage != null ? party.default_discount_percentage + '%' : '0%' }}</p>
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
              <label for="editRole">Rol</label>
              <select
                id="editRole"
                v-model="editForm.role"
              >
                <option value="CLIENT">Cliente</option>
                <option value="SUPPLIER">Proveedor</option>
                <option value="BOTH">Cliente y proveedor</option>
                <option value="CONTACT">Contacto</option>
              </select>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="editTaxId">NIF/CIF</label>
                <input
                  id="editTaxId"
                  v-model="editForm.taxId"
                  type="text"
                  @blur="validateEditField('taxId')"
                />
                <span v-if="editErrors.taxId" class="error">{{ editErrors.taxId }}</span>
              </div>

              <div class="form-group">
                <label for="editTaxIdType">Tipo de NIF/CIF</label>
                <select
                  id="editTaxIdType"
                  v-model="editForm.taxIdType"
                >
                  <option value="NIF">NIF</option>
                  <option value="CIF">CIF</option>
                  <option value="VAT">VAT</option>
                </select>
              </div>
            </div>

            <div class="form-group">
              <label for="editWebsite">Sitio web</label>
              <input
                id="editWebsite"
                v-model="editForm.website"
                type="text"
                placeholder="example.com"
                @blur="validateEditField('website')"
              />
              <span v-if="editErrors.website" class="error">{{ editErrors.website }}</span>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label for="editPhone">Teléfono</label>
                <input
                  id="editPhone"
                  v-model="editForm.phone"
                  type="tel"
                  @blur="validateEditField('phone')"
                />
                <span v-if="editErrors.phone" class="error">{{ editErrors.phone }}</span>
              </div>
              
              <div class="form-group">
                <label for="editEmail">Email</label>
                <input
                  id="editEmail"
                  v-model="editForm.email"
                  type="email"
                  @blur="validateEditField('email')"
                />
                <span v-if="editErrors.email" class="error">{{ editErrors.email }}</span>
              </div>
            </div>

            <div v-if="editForm.role === 'CLIENT' || editForm.role === 'BOTH'" class="form-group">
              <label for="editDiscount">Bonificación por defecto (%)</label>
              <input
                id="editDiscount"
                v-model.number="editForm.defaultDiscountPercentage"
                type="number"
                step="0.01"
                min="0"
                max="100"
                placeholder="0.00"
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

      <!-- Related Entities Section (for CONTACT role only) -->
      <div v-if="party.role === 'CONTACT'" class="card related-entities-section">
        <h3>Entidades Vinculadas</h3>
        
        <div v-if="isLoadingRelations" class="loading-relations">
          <div class="spinner-small"></div>
          <p>Cargando entidades...</p>
        </div>

        <div v-else-if="relatedEntities.length > 0" class="related-entities-list">
          <div 
            v-for="entity in relatedEntities" 
            :key="entity.id" 
            class="entity-card"
          >
            <div class="entity-header">
              <div class="entity-info">
                <div class="info-header">
                  <h4>{{ entity.name }}</h4>
                  <span class="badge" :class="`role-${entity.role.toLowerCase()}`">
                    {{ formatRole(entity.role) }}
                  </span>
                </div>
                <p v-if="entity.email" class="email">📧 {{ entity.email }}</p>
                <p v-if="entity.phone" class="phone">📞 {{ entity.phone }}</p>
                <p v-if="entity.tax_id" class="tax-id">🆔 {{ entity.tax_id }}</p>
              </div>
              <div class="entity-badges">
                <button 
                  type="button"
                  class="btn btn-primary" 
                  @click="navigateToEntity(entity.id)"
                >
                  Ver detalles
                </button>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="empty-state">
          <p>Este contacto no está vinculado a ninguna entidad aún.</p>
        </div>

        <div v-if="relationError" class="error-message-inline">
          <span>⚠️ {{ relationError }}</span>
        </div>
      </div>

      <!-- Contacts Section -->
      <person-manager
        v-if="party.role !== 'CONTACT'"
        :party-id="partyId"
      />

      <!-- Addresses Section -->
      <address-manager :party-id="partyId" />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { partyApi } from '@/services/partyApi';
import PersonManager from './PersonManager.vue';
import AddressManager from './AddressManager.vue';

const route = useRoute();
const router = useRouter();

const partyId = ref(route.params.id);

const party = ref(null);
const isLoading = ref(false);
const isUpdating = ref(false);
const error = ref('');
const isEditing = ref(false);

// Related entities (for CONTACT role)
const relatedEntities = ref([]);
const isLoadingRelations = ref(false);
const relationError = ref('');

const editForm = reactive({
  name: '',
  role: 'CLIENT',
  taxId: '',
  taxIdType: 'NIF',
  website: '',
  phone: '',
  email: '',
  notes: '',
  defaultDiscountPercentage: 0,
});

const editErrors = reactive({
  taxId: '',
  phone: '',
  email: '',
  website: '',
});

onMounted(() => {
  fetchParty();
});

async function fetchParty() {
  isLoading.value = true;
  error.value = '';

  try {
    const data = await partyApi.getParty(partyId.value);
    party.value = data;
    
    // Initialize edit form
    editForm.name = data.name;
    editForm.role = data.role || 'CLIENT';
    editForm.taxId = data.tax_id || '';
    editForm.taxIdType = data.tax_id_type || 'NIF';
    editForm.website = data.website || '';
    editForm.phone = data.phone || '';
    editForm.email = data.email || '';
    editForm.notes = data.notes || '';
    editForm.defaultDiscountPercentage = data.default_discount_percentage ?? 0;

    // Load related entities if this is a CONTACT
    if (data.role === 'CONTACT') {
      await fetchRelatedEntities();
    }
  } catch (err) {
    error.value = err?.message || 'Entidad no encontrada';
  } finally {
    isLoading.value = false;
  }
}

async function fetchRelatedEntities() {
  isLoadingRelations.value = true;
  relationError.value = '';
  
  try {
    // Get relationships for the contact
    const relationships = await partyApi.listRelationships(partyId.value);
    
    if (!relationships || relationships.length === 0) {
      relatedEntities.value = [];
      return;
    }

    // Extract entity IDs (from_party_id where contact is to_party_id, or vice versa)
    const entityIds = relationships
      .map(rel => {
        // If this contact is to_party_id, get from_party_id
        if (rel.to_party_id === partyId.value) {
          return rel.from_party_id;
        }
        // If this contact is from_party_id, get to_party_id
        if (rel.from_party_id === partyId.value) {
          return rel.to_party_id;
        }
        return null;
      })
      .filter(id => id !== null);

    // Fetch details for each related entity
    const entityPromises = entityIds.map(entityId => 
      partyApi.getParty(entityId).catch(() => null)
    );
    
    const entities = await Promise.all(entityPromises);
    relatedEntities.value = entities.filter(e => e !== null);
  } catch (err) {
    relationError.value = err?.message || 'No se pudieron cargar las entidades vinculadas';
  } finally {
    isLoadingRelations.value = false;
  }
}

function navigateToEntity(entityId) {
  partyId.value = entityId;
  router.push({ name: 'PartyDetail', params: { id: entityId } });
}

// Watch for route changes to reload data
watch(() => route.params.id, (newId) => {
  if (newId) {
    partyId.value = newId;
    fetchParty();
  }
});

async function submitEdit() {
  if (!editForm.name.trim()) {
    alert('El nombre es obligatorio');
    return;
  }

  // Validar todos los campos
  validateEditField('taxId');
  validateEditField('email');
  validateEditField('phone');
  validateEditField('website');

  // Verificar si hay errores
  if (editErrors.taxId || editErrors.email || editErrors.phone || editErrors.website) {
    alert('Por favor, corrija los errores en el formulario antes de guardar');
    return;
  }

  // Normalizar URL si se proporciona
  if (editForm.website && editForm.website.trim()) {
    editForm.website = normalizeUrl(editForm.website);
  }

  isUpdating.value = true;

  try {
    const updatePayload = {
      name: editForm.name,
      role: editForm.role,
      taxId: editForm.taxId,
      taxIdType: editForm.taxIdType,
      website: editForm.website || '',
      phone: editForm.phone,
      email: editForm.email,
      notes: editForm.notes,
    };
    if (editForm.role === 'CLIENT' || editForm.role === 'BOTH') {
      updatePayload.default_discount_percentage = editForm.defaultDiscountPercentage || 0;
    }
    const updated = await partyApi.updateParty(partyId.value, updatePayload);

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
    const updated = await partyApi.changePartyStatus(partyId.value, newStatus);
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
    CONTACT: 'Contacto',
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

function normalizeUrl(url) {
  if (!url || !url.trim()) return url;
  
  const trimmedUrl = url.trim();
  // Si ya tiene protocolo, devolver tal cual
  if (/^https?:\/\//i.test(trimmedUrl)) {
    return trimmedUrl;
  }
  
  // Si no tiene protocolo, añadir https://
  return `https://${trimmedUrl}`;
}

function isValidUrl(string) {
  if (!string || !string.trim()) return true; // Empty is valid
  
  try {
    const url = new URL(string);
    
    // Validar que el hostname tenga al menos un punto (dominio válido)
    // Esto rechaza casos como "https://asdf" y acepta "https://example.com"
    if (!url.hostname.includes('.')) {
      return false;
    }
    
    return true;
  } catch (_) {
    return false;
  }
}

function isValidEmail(string) {
  const emailRegex = /^[a-zA-Z0-9.+_%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(string.trim());
}

function isValidPhone(string) {
  const phoneRegex = /^[\+]?[\d\s\-()]{8,}$/;
  return phoneRegex.test(string.trim());
}

function isValidTaxId(taxId, taxIdType) {
  if (!taxId || !taxId.trim()) return true; // Empty is valid
  
  const trimmed = taxId.trim().toUpperCase();
  
  if (taxIdType === 'NIF') {
    // NIF español: 8 dígitos + letra
    const nifRegex = /^[0-9]{8}[A-Z]$/;
    return nifRegex.test(trimmed);
  } else if (taxIdType === 'CIF') {
    // CIF español: letra + 7 dígitos + dígito o letra
    const cifRegex = /^[A-Z][0-9]{7}[0-9A-Z]$/;
    return cifRegex.test(trimmed);
  } else if (taxIdType === 'VAT') {
    // VAT genérico: al menos 2 caracteres
    return trimmed.length >= 2;
  }
  
  return true;
}

const editValidationRules = {
  taxId: (value) => {
    if (value && !isValidTaxId(value, editForm.taxIdType)) {
      if (editForm.taxIdType === 'NIF') {
        return 'Formato de NIF inválido (debe ser 8 dígitos seguidos de una letra)';
      } else if (editForm.taxIdType === 'CIF') {
        return 'Formato de CIF inválido (debe ser letra + 7 dígitos + dígito o letra)';
      }
      return 'Formato inválido';
    }
    return '';
  },
  website: (value) => {
    if (value && value.trim()) {
      const normalized = normalizeUrl(value);
      if (!isValidUrl(normalized)) {
        return 'Formato de URL inválido';
      }
    }
    return '';
  },
  phone: (value) => {
    if (value && value.trim() && !isValidPhone(value)) {
      return 'Formato inválido. Debe tener al menos 8 dígitos y puede incluir +, espacios, guiones y paréntesis';
    }
    return '';
  },
  email: (value) => {
    if (value && value.trim() && !isValidEmail(value)) {
      return 'Formato de email inválido';
    }
    return '';
  },
};

function validateEditField(fieldName) {
  const validator = editValidationRules[fieldName];
  if (validator) {
    editErrors[fieldName] = validator(editForm[fieldName]);
  }
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

.info-note h3 {
  margin: 0 0 0.5rem;
  color: #1b3a6b;
}

.info-note p {
  margin: 0;
  color: #64748b;
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
  background: #eac54f;
  color: #1e293b;
  border: none;
  font-weight: 600;
}

.btn-primary:hover:not(:disabled) {
  background: #d4a41d;
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
.form-group select,
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
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group .error {
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.edit-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

/* Related Entities Section */
.related-entities-section {
  margin-top: 1.5rem;
}

.related-entities-section h3 {
  color: #1b3a6b;
  margin: 0 0 1rem 0;
}

.loading-relations {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  justify-content: center;
  color: #64748b;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(230, 184, 0, 0.2);
  border-top-color: #e6b800;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.related-entities-list {
  display: grid;
  gap: 1rem;
}

.entity-card {
  padding: 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.entity-card:hover {
  border-color: #002395;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.08);
}

.entity-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.entity-info h4 {
  color: #1e293b;
  margin: 0 0 0.5rem 0;
}

.entity-info {
  flex: 1;
}

.info-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.info-header h4 {
  margin: 0;
}

.entity-info p {
  color: #64748b;
  margin: 0.25rem 0;
  font-size: 0.9rem;
}

.entity-badges {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  flex-wrap: wrap;
  flex-shrink: 0;
}

.error-message-inline {
  color: #991b1b;
  background-color: #fee2e2;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-top: 1rem;
  border: 1px solid #ef4444;
}

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: #64748b;
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

  .form-row {
    grid-template-columns: 1fr;
  }

  .entity-info {
    flex-direction: column;
    align-items: flex-start;
  }

  .entity-meta {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style>
