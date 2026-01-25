<template>
  <div class="address-manager">
    <div class="manager-header">
      <h3>Addresses</h3>
      <button @click="showForm = !showForm" class="btn btn-primary">
        {{ showForm ? '✕ Close' : '+ Add Address' }}
      </button>
    </div>

    <!-- Add Form -->
    <div v-if="showForm" class="form-section">
      <form @submit.prevent="submitForm">
        <div class="form-group">
          <label for="street">Street Address *</label>
          <input
            id="street"
            v-model="form.street"
            type="text"
            placeholder="Street and house number"
            required
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="city">City *</label>
            <input
              id="city"
              v-model="form.city"
              type="text"
              placeholder="City"
              required
            />
          </div>
          <div class="form-group">
            <label for="province">Province/State *</label>
            <input
              id="province"
              v-model="form.province"
              type="text"
              placeholder="Province"
              required
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="postalCode">Postal Code *</label>
            <input
              id="postalCode"
              v-model="form.postalCode"
              type="text"
              placeholder="28001"
              required
            />
          </div>
          <div class="form-group">
            <label for="country">Country</label>
            <input
              id="country"
              v-model="form.country"
              type="text"
              placeholder="Spain"
            />
          </div>
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="isSubmitting" class="btn btn-primary">
            {{ isSubmitting ? 'Adding...' : 'Add Address' }}
          </button>
          <button type="button" @click="resetForm" class="btn btn-secondary">
            Cancel
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
            <span v-if="address.is_primary" class="badge primary">Primary</span>
            <span class="badge date">{{ formatDate(address.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="addresses.length === 0 && !showForm" class="empty-state">
      <p>No addresses yet. Add your first address to get started.</p>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Loading addresses...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import { partyApi } from '@/services/partyApi';

const props = defineProps({
  organizationId: {
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

watch(() => props.organizationId, () => {
  if (props.organizationId) {
    fetchAddresses();
  }
});

async function fetchAddresses() {
  if (!props.organizationId) return;

  isLoading.value = true;
  try {
    const response = await partyApi.listAddresses(props.organizationId);
    addresses.value = response.data || [];
  } catch (error) {
    formError.value = error.message || 'Failed to load addresses';
  } finally {
    isLoading.value = false;
  }
}

async function submitForm() {
  if (!form.street || !form.city || !form.province || !form.postalCode) {
    formError.value = 'Street, city, province, and postal code are required';
    return;
  }

  isSubmitting.value = true;
  formError.value = '';

  try {
    await partyApi.addAddress(props.organizationId, {
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
    formError.value = error.message || 'Failed to add address';
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
  background: var(--color-surface);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.manager-header h3 {
  color: var(--color-text-primary);
  margin: 0;
}

.form-section {
  background: rgba(230, 184, 0, 0.05);
  padding: 1.5rem;
  border-radius: 6px;
  margin-bottom: 1.5rem;
  border: 1px solid rgba(230, 184, 0, 0.2);
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
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
}

.form-group input {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 0.95rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary {
  background-color: var(--primary-color);
  color: var(--color-text-on-primary);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--primary-color-hover);
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(230, 184, 0, 0.3);
}

.btn-secondary {
  background-color: var(--color-secondary);
  color: var(--color-text-primary);
}

.error-message {
  color: var(--color-error);
  background-color: rgba(244, 67, 54, 0.1);
  padding: 0.75rem;
  border-radius: 4px;
  margin-top: 1rem;
  border-left: 3px solid var(--color-error);
}

.addresses-list {
  display: grid;
  gap: 1rem;
}

.address-card {
  padding: 1rem;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: var(--color-background);
  transition: all 0.2s ease;
}

.address-card:hover {
  border-color: var(--primary-color);
  box-shadow: 0 2px 8px rgba(230, 184, 0, 0.1);
}

.address-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.address-info h4 {
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
}

.address-info p {
  color: var(--color-text-secondary);
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
  color: var(--primary-color);
}

.badge.date {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--color-text-secondary);
}

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
