<template>
  <div class="organization-form">
    <div class="form-container">
      <h2>{{ isEditing ? 'Edit Organization' : 'Create Organization' }}</h2>
      
      <form @submit.prevent="submitForm">
        <!-- Basic Information -->
        <fieldset>
          <legend>Basic Information</legend>
          
          <div class="form-group">
            <label for="name">
              Organization Name *
              <span class="required">required</span>
            </label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              placeholder="Enter organization name"
              required
              @blur="validateField('name')"
            />
            <span v-if="errors.name" class="error">{{ errors.name }}</span>
          </div>

          <div class="form-group">
            <label for="role">
              Organization Role *
              <span class="required">required</span>
            </label>
            <select
              id="role"
              v-model="form.role"
              required
              @change="validateField('role')"
            >
              <option value="">-- Select Role --</option>
              <option value="CLIENT">Client</option>
              <option value="SUPPLIER">Supplier</option>
              <option value="BOTH">Both Client & Supplier</option>
            </select>
            <span v-if="errors.role" class="error">{{ errors.role }}</span>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="taxId">Tax ID</label>
              <input
                id="taxId"
                v-model="form.taxId"
                type="text"
                placeholder="e.g., 12345678A"
                @blur="validateField('taxId')"
              />
              <span v-if="errors.taxId" class="error">{{ errors.taxId }}</span>
            </div>

            <div class="form-group">
              <label for="taxIdType">Tax ID Type</label>
              <select id="taxIdType" v-model="form.taxIdType">
                <option value="NIF">NIF</option>
                <option value="CIF">CIF</option>
                <option value="VAT">VAT</option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label for="website">Website</label>
            <input
              id="website"
              v-model="form.website"
              type="url"
              placeholder="https://example.com"
              @blur="validateField('website')"
            />
            <span v-if="errors.website" class="error">{{ errors.website }}</span>
          </div>
        </fieldset>

        <!-- Additional Information -->
        <fieldset v-if="isEditing">
          <legend>Additional Information</legend>
          
          <div class="form-group">
            <label for="notes">Notes</label>
            <textarea
              id="notes"
              v-model="form.notes"
              placeholder="Add any additional notes..."
              rows="4"
            />
          </div>
        </fieldset>

        <!-- Form Actions -->
        <div class="form-actions">
          <button
            type="submit"
            :disabled="isSubmitting"
            class="btn btn-primary"
          >
            <span v-if="isSubmitting">{{ isEditing ? 'Updating...' : 'Creating...' }}</span>
            <span v-else>{{ isEditing ? 'Update Organization' : 'Create Organization' }}</span>
          </button>
          <button
            type="button"
            @click="resetForm"
            class="btn btn-secondary"
          >
            Reset
          </button>
        </div>
      </form>

      <!-- Success/Error Messages -->
      <div v-if="successMessage" class="message success">
        <span>✓ {{ successMessage }}</span>
        <button @click="successMessage = ''" class="close">&times;</button>
      </div>
      <div v-if="errorMessage" class="message error">
        <span>✗ {{ errorMessage }}</span>
        <button @click="errorMessage = ''" class="close">&times;</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue';
import { partyApi } from '@/services/partyApi';

const props = defineProps({
  organizationId: {
    type: String,
    default: null,
  },
  initialData: {
    type: Object,
    default: null,
  },
});

const emit = defineEmits(['submit', 'update']);

// Form state
const form = reactive({
  id: '',
  name: '',
  role: '',
  taxId: '',
  taxIdType: 'NIF',
  website: '',
  notes: '',
});

const errors = reactive({});
const isSubmitting = ref(false);
const successMessage = ref('');
const errorMessage = ref('');

const isEditing = computed(() => !!props.organizationId);

// Initialize form with data
if (props.initialData) {
  Object.assign(form, props.initialData);
}

// Validation rules
const validationRules = {
  name: (value) => {
    if (!value || value.trim().length === 0) {
      return 'Organization name is required';
    }
    if (value.length < 3) {
      return 'Organization name must be at least 3 characters';
    }
    if (value.length > 100) {
      return 'Organization name must not exceed 100 characters';
    }
    return '';
  },
  role: (value) => {
    if (!value) {
      return 'Role is required';
    }
    if (!['CLIENT', 'SUPPLIER', 'BOTH'].includes(value)) {
      return 'Invalid role selected';
    }
    return '';
  },
  taxId: (value) => {
    if (value && (value.length < 5 || value.length > 20)) {
      return 'Tax ID must be between 5 and 20 characters';
    }
    return '';
  },
  website: (value) => {
    if (value && !isValidUrl(value)) {
      return 'Invalid URL format';
    }
    return '';
  },
};

// Helper functions
function isValidUrl(string) {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
}

function validateField(fieldName) {
  const validator = validationRules[fieldName];
  if (validator) {
    errors[fieldName] = validator(form[fieldName]);
  }
}

function validateForm() {
  Object.keys(validationRules).forEach((field) => {
    validateField(field);
  });
  return Object.values(errors).every((err) => !err);
}

async function submitForm() {
  if (!validateForm()) {
    errorMessage.value = 'Please fix the errors above';
    return;
  }

  isSubmitting.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    let result;
    
    if (isEditing.value) {
      result = await partyApi.updateOrganization(props.organizationId, {
        name: form.name,
        website: form.website,
        notes: form.notes,
      });
      successMessage.value = 'Organization updated successfully';
      emit('update', result);
    } else {
      result = await partyApi.createOrganization({
        id: form.id || generateId(),
        name: form.name,
        role: form.role,
        taxId: form.taxId,
        taxIdType: form.taxIdType,
        website: form.website,
      });
      successMessage.value = 'Organization created successfully';
      resetForm();
      emit('submit', result);
    }
  } catch (error) {
    errorMessage.value = error.data?.message || error.message || 'Failed to save organization';
  } finally {
    isSubmitting.value = false;
  }
}

function resetForm() {
  form.id = '';
  form.name = '';
  form.role = '';
  form.taxId = '';
  form.taxIdType = 'NIF';
  form.website = '';
  form.notes = '';
  Object.keys(errors).forEach((key) => {
    errors[key] = '';
  });
}

function generateId() {
  return `org-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}
</script>

<style scoped>
.organization-form {
  padding: 1.5rem;
  background: var(--color-background);
}

.form-container {
  max-width: 600px;
  margin: 0 auto;
  background: var(--color-surface);
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.form-container h2 {
  color: var(--color-text-primary);
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
  border-bottom: 2px solid var(--primary-color);
  padding-bottom: 0.5rem;
}

fieldset {
  border: none;
  padding: 1rem 0;
  margin: 1rem 0;
  border-top: 1px solid var(--color-border);
}

fieldset:first-of-type {
  border-top: none;
  margin-top: 0;
  padding-top: 0;
}

legend {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
  display: flex;
  flex-direction: column;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

label {
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.required {
  font-size: 0.75rem;
  color: var(--color-error);
  font-weight: 600;
}

input,
select,
textarea {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
  transition: all 0.2s ease;
}

input:focus,
select:focus,
textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

input[type="url"]::placeholder,
input[type="text"]::placeholder {
  color: var(--color-text-secondary);
}

textarea {
  resize: vertical;
  font-family: inherit;
}

.error {
  color: var(--color-error);
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  flex: 1;
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

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: var(--color-secondary);
  color: var(--color-text-primary);
}

.btn-secondary:hover:not(:disabled) {
  background-color: var(--color-border);
}

.message {
  margin-top: 1.5rem;
  padding: 1rem;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message.success {
  background-color: rgba(76, 175, 80, 0.1);
  color: var(--color-success);
  border-left: 4px solid var(--color-success);
}

.message.error {
  background-color: rgba(244, 67, 54, 0.1);
  color: var(--color-error);
  border-left: 4px solid var(--color-error);
}

.message .close {
  background: none;
  border: none;
  color: inherit;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
}

@media (max-width: 768px) {
  .form-container {
    padding: 1.5rem;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
