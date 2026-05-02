<template>
  <div class="address-manager">
    <div class="manager-header">
      <h3>Direcciones</h3>
      <button @click="showForm = !showForm" class="btn btn-primary">
        {{ showForm ? '✕ Cerrar' : '+ Agregar dirección' }}
      </button>
    </div>

    <!-- Add/Edit Form -->
    <div v-if="showForm" class="form-section">
      <h4 class="form-title">{{ isEditing ? 'Editar dirección' : 'Agregar nueva dirección' }}</h4>
      <form @submit.prevent="submitForm">
        <div class="form-group">
          <label for="street">Calle y número *</label>
          <input
            id="street"
            v-model="form.street"
            type="text"
            placeholder="Calle y número"
            required
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="city">Ciudad *</label>
            <input
              id="city"
              v-model="form.city"
              type="text"
              placeholder="Ciudad"
              required
            />
          </div>
          <div class="form-group">
            <label for="province">Provincia/Estado *</label>
            <input
              id="province"
              v-model="form.province"
              type="text"
              placeholder="Provincia"
              required
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="postalCode">Código postal *</label>
            <input
              id="postalCode"
              v-model="form.postalCode"
              type="text"
              placeholder="28001"
              required
            />
          </div>
          <div class="form-group">
            <label for="country">País</label>
            <input
              id="country"
              v-model="form.country"
              type="text"
              placeholder="España"
            />
          </div>
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="isSubmitting" class="btn btn-primary">
            {{ isSubmitting 
              ? (isEditing ? 'Actualizando...' : 'Agregando...') 
              : (isEditing ? 'Actualizar' : 'Agregar') 
            }}
          </button>
          <button type="button" @click="resetForm" class="btn btn-secondary">
            Cancelar
          </button>
        </div>
      </form>

      <div v-if="formError" class="error-message">
        <span>✗ {{ formError }}</span>
      </div>
    </div>

    <!-- Addresses List -->
    <div v-if="addresses.length > 0" class="addresses-list">
      <div v-for="address in addresses" :key="address.id" class="address-card">
        <div class="address-header">
          <div class="address-info">
            <h4>{{ address.street }}</h4>
            <p class="location">
              <MapPin :size="18" />
              {{ address.city }}, {{ address.province }} {{ address.postal_code }}
            </p>
            <p v-if="address.country" class="country">
              <Globe :size="18" />
              {{ address.country }}
            </p>
          </div>
          <div class="address-badges-group">
            <div class="address-badges">
              <span v-if="address.is_primary" class="badge status-success">Principal</span>
              <span class="badge status-secondary">{{ formatDate(address.created_at) }}</span>
            </div>
            <div class="action-buttons">
              <button 
                @click="editAddress(address)" 
                class="btn-icon"
                title="Editar dirección"
              >
                <Pencil :size="18" />
              </button>
              <button 
                @click="deleteAddress(address.id)" 
                class="btn-icon text-danger"
                title="Eliminar dirección"
              >
                <Trash2 :size="18" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="addresses.length === 0 && !showForm" class="empty-state">
      <p>No hay direcciones aún. Agrega la primera para empezar.</p>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Cargando direcciones...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import { MapPin, Globe, Pencil, Trash2 } from 'lucide-vue-next';
import { partyApi } from '@/services/partyApi';

const props = defineProps({
  partyId: {
    type: String,
    required: true,
  },
});

const addresses = ref([]);
const isLoading = ref(false);
const isSubmitting = ref(false);
const showForm = ref(false);
const formError = ref('');
const isEditing = ref(false);
const editingAddressId = ref(null);

const form = reactive({
  street: '',
  city: '',
  province: '',
  postalCode: '',
  country: 'Spain',
});

onMounted(() => {
  fetchAddresses();
});

watch(() => props.partyId, () => {
  if (props.partyId) {
    fetchAddresses();
  }
});

async function fetchAddresses() {
  if (!props.partyId) return;

  isLoading.value = true;
  try {
    const response = await partyApi.listPartyAddresses(props.partyId);
    addresses.value = response.data || [];
  } catch (error) {
    formError.value = error?.message || 'No se pudieron cargar las direcciones';
  } finally {
    isLoading.value = false;
  }
}

async function submitForm() {
  if (!form.street || !form.city || !form.province || !form.postalCode) {
    formError.value = 'Calle, ciudad, provincia y código postal son obligatorios';
    return;
  }

  isSubmitting.value = true;
  formError.value = '';

  try {
    const addressData = {
      id: isEditing.value ? editingAddressId.value : `addr-${Date.now()}`,
      street: form.street,
      city: form.city,
      province: form.province,
      postalCode: form.postalCode,
      country: form.country || 'Spain',
    };

    if (isEditing.value) {
      await partyApi.updatePartyAddress(props.partyId, editingAddressId.value, addressData);
    } else {
      await partyApi.addPartyAddress(props.partyId, addressData);
    }

    resetForm();
    await fetchAddresses();
  } catch (error) {
    formError.value = error?.message || `No se pudo ${isEditing.value ? 'actualizar' : 'agregar'} la dirección`;
  } finally {
    isSubmitting.value = false;
  }
}

function editAddress(address) {
  form.street = address.street;
  form.city = address.city;
  form.province = address.province;
  form.postalCode = address.postal_code;
  form.country = address.country || 'Spain';
  
  isEditing.value = true;
  editingAddressId.value = address.id;
  showForm.value = true;
  formError.value = '';
}

async function deleteAddress(addressId) {
  if (!confirm('¿Estás seguro de que quieres eliminar esta dirección?')) {
    return;
  }

  try {
    await partyApi.deletePartyAddress(props.partyId, addressId);
    await fetchAddresses();
  } catch (error) {
    formError.value = error?.message || 'No se pudo eliminar la dirección';
  }
}

function resetForm() {
  form.street = '';
  form.city = '';
  form.province = '';
  form.postalCode = '';
  form.country = 'Spain';
  formError.value = '';
  showForm.value = false;
  isEditing.value = false;
  editingAddressId.value = null;
}

function formatDate(dateString) {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}
</script>

<style scoped>
.address-manager {
  padding: var(--spacing-lg);
  background: var(--color-surface);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--color-border);
  box-shadow: var(--box-shadow-sm);
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.manager-header h3 {
  color: var(--color-secondary);
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: 700;
}

.form-section {
  background: var(--color-background);
  padding: var(--spacing-lg);
  border-radius: var(--border-radius-md);
  margin-bottom: var(--spacing-lg);
  border: 1px solid var(--color-border);
}

.form-title {
  color: var(--color-text-primary);
  margin: 0 0 var(--spacing-md) 0;
  font-size: 1rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.form-group label {
  display: block;
  font-size: var(--font-size-xs);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-secondary);
}

.form-group input {
  padding: 0.75rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-md);
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  background: white;
  transition: all 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.form-actions {
  display: flex;
  gap: var(--spacing-md);
  margin-top: var(--spacing-md);
}

.error-message {
  color: var(--color-error);
  background-color: var(--color-background);
  padding: 0.75rem 1rem;
  border-radius: var(--border-radius-md);
  margin-top: var(--spacing-md);
  border: 1px solid var(--color-error);
  font-size: var(--font-size-sm);
}

.addresses-list {
  display: grid;
  gap: var(--spacing-md);
}

.address-card {
  padding: var(--spacing-md);
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-md);
  background: white;
  transition: all 0.2s ease;
}

.address-card:hover {
  border-color: var(--color-border-strong);
  box-shadow: var(--box-shadow-md);
}

.address-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-md);
}

.address-info h4 {
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  font-weight: 700;
}

.address-info p {
  color: var(--color-text-secondary);
  margin: 0.25rem 0;
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.icon-sm { font-size: 18px; color: var(--color-text-secondary); }

.address-badges-group {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--spacing-sm);
}

.address-badges {
  display: flex;
  gap: 0.5rem;
}

.action-buttons {
  display: flex;
  gap: 0.25rem;
}

.btn-icon {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--color-text-secondary);
  padding: 0.4rem;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: var(--color-background);
  color: var(--color-text-primary);
}

.text-danger { color: var(--color-error); }
.text-danger:hover { background: var(--color-primary-light); }

.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
}

.status-success { background: rgba(22, 163, 74, 0.1); color: var(--color-success); }
.status-secondary { background: var(--color-background); color: var(--color-text-secondary); }

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: var(--color-text-secondary);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  gap: 1rem;
}

.spinner {
  width: 24px; height: 24px; border: 2px solid var(--color-border);
  border-top-color: var(--color-primary); border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 768px) {
  .form-row { grid-template-columns: 1fr; }
  .address-header { flex-direction: column; }
  .address-badges-group { align-items: flex-start; }
  .form-actions { flex-direction: column; }
}
</style>
