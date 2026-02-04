<template>
  <div class="address-manager">
    <div class="manager-header">
      <h3>Direcciones</h3>
      <button @click="showForm = !showForm" class="btn btn-primary">
        {{ showForm ? '✕ Cerrar' : '+ Agregar dirección' }}
      </button>
    </div>

    <!-- Add Form -->
    <div v-if="showForm" class="form-section">
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
            {{ isSubmitting ? 'Agregando...' : 'Agregar dirección' }}
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
            <p class="location">📍 {{ address.city }}, {{ address.province }} {{ address.postal_code }}</p>
            <p v-if="address.country" class="country">🌍 {{ address.country }}</p>
          </div>
          <div class="address-badges">
            <span v-if="address.is_primary" class="badge primary">Principal</span>
            <span class="badge date">{{ formatDate(address.created_at) }}</span>
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
    await partyApi.addPartyAddress(props.partyId, {
      id: `addr-${Date.now()}`,
      street: form.street,
      city: form.city,
      province: form.province,
      postalCode: form.postalCode,
      country: form.country || 'Spain',
    });

    resetForm();
    await fetchAddresses();
  } catch (error) {
    formError.value = error?.message || 'No se pudo agregar la dirección';
  } finally {
    isSubmitting.value = false;
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
  padding: 1.5rem;
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.manager-header h3 {
  color: #1b3a6b;
  margin: 0;
}

.form-section {
  background: #f8fafc;
  padding: 1.5rem;
  border-radius: 10px;
  margin-bottom: 1.5rem;
  border: 1px solid #e2e8f0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
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

.form-group input {
  padding: 0.6rem 0.8rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.9rem;
  color: #1e293b;
}

.form-group input:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
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

.error-message {
  color: #991b1b;
  background-color: #fee2e2;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-top: 1rem;
  border: 1px solid #ef4444;
}

.addresses-list {
  display: grid;
  gap: 1rem;
}

.address-card {
  padding: 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.address-card:hover {
  border-color: #002395;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.08);
}

.address-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.address-info h4 {
  color: #1e293b;
  margin: 0 0 0.5rem 0;
}

.address-info p {
  color: #64748b;
  margin: 0.25rem 0;
  font-size: 0.9rem;
}

.address-badges {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
}

.badge.primary {
  background-color: rgba(230, 184, 0, 0.2);
  color: #1e293b;
}

.badge.date {
  background-color: rgba(0, 0, 0, 0.05);
  color: #64748b;
}

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: #64748b;
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
  width: 30px;
  height: 30px;
  border: 3px solid rgba(230, 184, 0, 0.2);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .address-header {
    flex-direction: column;
    gap: 1rem;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
