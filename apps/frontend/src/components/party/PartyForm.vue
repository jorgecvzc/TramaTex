<template>
  <div class="party-form">
    <div class="form-container">
      <h2>{{ isEditing ? 'Editar entidad' : 'Crear entidad' }}</h2>
      
      <form @submit.prevent="submitForm">
        <!-- Basic Information -->
        <fieldset>
          <legend>Información básica</legend>
          
          <div class="form-group">
            <label for="name">
              Nombre de la entidad *
              <span class="required">obligatorio</span>
            </label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              placeholder="Ingresa el nombre"
              required
              @blur="validateField('name')"
            />
            <span v-if="errors.name" class="error">{{ errors.name }}</span>
          </div>

          <div class="form-group">
            <label for="role">
              Rol de la entidad *
              <span class="required">obligatorio</span>
            </label>
            <select
              id="role"
              v-model="form.role"
              required
              @change="validateField('role')"
            >
              <option value="">-- Selecciona rol --</option>
              <option value="CLIENT">Cliente</option>
              <option value="SUPPLIER">Proveedor</option>
              <option value="BOTH">Cliente y proveedor</option>
              <option value="CONTACT">Contacto</option>
            </select>
            <span v-if="errors.role" class="error">{{ errors.role }}</span>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="taxId">NIF/CIF</label>
              <input
                id="taxId"
                v-model="form.taxId"
                type="text"
                placeholder="p. ej., 12345678A"
                @blur="validateField('taxId')"
              />
              <span v-if="errors.taxId" class="error">{{ errors.taxId }}</span>
            </div>

            <div class="form-group">
              <label for="taxIdType">Tipo de NIF/CIF</label>
              <select id="taxIdType" v-model="form.taxIdType">
                <option value="NIF">NIF</option>
                <option value="CIF">CIF</option>
                <option value="VAT">VAT</option>
              </select>
            </div>
          </div>

          <div class="form-group">
              <label for="website">Sitio web</label>
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
          <legend>Información adicional</legend>
          
          <div class="form-group">
            <label for="notes">Notas</label>
            <textarea
              id="notes"
              v-model="form.notes"
              placeholder="Agrega notas adicionales..."
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
            <span v-if="isSubmitting">{{ isEditing ? 'Actualizando...' : 'Creando...' }}</span>
            <span v-else>{{ isEditing ? 'Actualizar entidad' : 'Crear entidad' }}</span>
          </button>
          <button
            type="button"
            @click="resetForm"
            class="btn btn-secondary"
          >
            Reiniciar
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
  partyId: {
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

const isEditing = computed(() => !!props.partyId);

// Initialize form with data
if (props.initialData) {
  const initial = props.initialData;
  Object.assign(form, {
    ...initial,
    role: initial.role || form.role,
    taxId: initial.taxId ?? initial.tax_id ?? form.taxId,
    taxIdType: initial.taxIdType ?? initial.tax_id_type ?? form.taxIdType,
  });

  if (isEditing.value && !form.role) {
    form.role = 'CLIENT';
  }
}

// Validation rules
const validationRules = {
  name: (value) => {
    if (!value || value.trim().length === 0) {
      return 'El nombre es obligatorio';
    }
    if (value.length < 3) {
      return 'El nombre debe tener al menos 3 caracteres';
    }
    if (value.length > 100) {
      return 'El nombre no debe superar 100 caracteres';
    }
    return '';
  },
  role: (value) => {
    if (!value) {
      return 'El rol es obligatorio';
    }
    if (!['CLIENT', 'SUPPLIER', 'BOTH', 'CONTACT'].includes(value)) {
      return 'Rol inválido';
    }
    return '';
  },
  taxId: (value) => {
    if (value && (value.length < 5 || value.length > 20)) {
      return 'El NIF/CIF debe tener entre 5 y 20 caracteres';
    }
    return '';
  },
  website: (value) => {
    if (value && !isValidUrl(value)) {
      return 'Formato de URL inválido';
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
    errorMessage.value = 'Corrige los errores indicados';
    return;
  }

  isSubmitting.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    let result;
    
    if (isEditing.value) {
      result = await partyApi.updateParty(props.partyId, {
        name: form.name,
        role: form.role,
        taxId: form.taxId,
        taxIdType: form.taxIdType,
        website: form.website,
        notes: form.notes,
      });
      successMessage.value = 'Entidad actualizada correctamente';
      emit('update', result);
    } else {
      result = await partyApi.createParty({
        id: form.id || generateId(),
        name: form.name,
        role: form.role,
        taxId: form.taxId,
        taxIdType: form.taxIdType,
        website: form.website,
      });
      successMessage.value = 'Entidad creada correctamente';
      resetForm();
      emit('submit', result);
    }
  } catch (error) {
    errorMessage.value = error?.data?.message || error?.message || 'No se pudo guardar la entidad';
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
  return `party-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}
</script>

<style scoped>
.party-form {
  width: 100%;
}

.form-container {
  max-width: 100%;
  margin: 0;
  background: transparent;
  padding: 0;
  border-radius: 0;
  box-shadow: none;
  border: none;
}

.form-container h2 {
  color: #1b3a6b;
  margin: 0 0 1.5rem;
  font-size: 1.4rem;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 0.75rem;
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
  font-size: 0.95rem;
  font-weight: 600;
  color: #1e293b;
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
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin-bottom: 0.4rem;
}

.required {
  font-size: 0.75rem;
  color: #ef4444;
  font-weight: 600;
}

input,
select,
textarea {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.6rem 0.8rem;
  font-size: 0.9rem;
  color: #1e293b;
  font-family: inherit;
}

input:focus,
select:focus,
textarea:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
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
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
  flex: 1;
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

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

.message {
  margin-top: 1.5rem;
  padding: 1rem;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message.success {
  background-color: rgba(76, 175, 80, 0.1);
  color: #166534;
  border-left: 4px solid #22c55e;
}

.message.error {
  background-color: rgba(244, 67, 54, 0.1);
  color: #991b1b;
  border-left: 4px solid #ef4444;
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
