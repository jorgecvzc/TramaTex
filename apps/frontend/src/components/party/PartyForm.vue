<template>
  <div class="party-form">
    <div class="form-container">
      <h2>{{ isEditing ? 'Editar entidad' : 'Crear entidad' }}</h2>
      
      <form @submit.prevent="submitForm">
        <!-- Basic Information -->
        <fieldset>
          <legend>Información básica</legend>
          
          <div class="form-group">
            <label for="role">
              Rol de la entidad *
              <span class="required">obligatorio</span>
            </label>
            <select
              id="role"
              v-model="form.role"
              required
              @change="onRoleChange"
            >
              <option value="">-- Selecciona rol --</option>
              <option value="CLIENT">Cliente</option>
              <option value="SUPPLIER">Proveedor</option>
              <option value="BOTH">Cliente y proveedor</option>
              <option value="CONTACT">Contacto</option>
            </select>
            <span v-if="errors.role" class="error">{{ errors.role }}</span>
            <small v-if="form.role === 'CONTACT'" class="help-text warning">
              ⚠️ Los contactos deben ser <strong>personas físicas</strong> únicamente.
            </small>
          </div>

          <div class="form-group">
            <label for="entityType">
              Tipo de entidad *
              <span class="required">obligatorio</span>
            </label>
            <select
              id="entityType"
              v-model="form.entityType"
              required
              @change="onEntityTypeChange"
              :disabled="form.role === 'CONTACT'"
            >
              <option value="">-- Selecciona tipo --</option>
              <option value="PERSON">Persona Física</option>
              <option value="ORGANIZATION" :disabled="form.role === 'CONTACT'">
                Persona Jurídica (Organización)
              </option>
            </select>
            <span v-if="errors.entityType" class="error">{{ errors.entityType }}</span>
            <small v-if="form.role !== 'CONTACT'" class="help-text">
              • <strong>Persona Física:</strong> Individuo (cliente, contacto, empleado)<br>
              • <strong>Persona Jurídica:</strong> Empresa, organización, entidad legal
            </small>
            <small v-else class="help-text warning">
              ⚠️ Los contactos solo pueden ser <strong>personas físicas</strong>
            </small>
          </div>

          <!-- Fields for PERSON entity type -->
          <template v-if="form.entityType === 'PERSON'">
            <div class="form-row">
              <div class="form-group">
                <label for="firstName">
                  Nombre *
                  <span class="required">obligatorio</span>
                </label>
                <input
                  id="firstName"
                  v-model="form.firstName"
                  type="text"
                  placeholder="Ej: Juan"
                  required
                  @blur="validateField('firstName')"
                />
                <span v-if="errors.firstName" class="error">{{ errors.firstName }}</span>
              </div>

              <div class="form-group">
                <label for="lastName">
                  Apellido(s) *
                  <span class="required">obligatorio</span>
                </label>
                <input
                  id="lastName"
                  v-model="form.lastName"
                  type="text"
                  placeholder="Ej: García López"
                  required
                  @blur="validateField('lastName')"
                />
                <span v-if="errors.lastName" class="error">{{ errors.lastName }}</span>
              </div>
            </div>
          </template>

          <!-- Fields for ORGANIZATION entity type -->
          <div v-if="form.entityType === 'ORGANIZATION'" class="form-group">
            <label for="name">
              Nombre de la organización *
              <span class="required">obligatorio</span>
            </label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              placeholder="Ej: Acme Corporation S.L."
              required
              @blur="validateField('name')"
            />
            <span v-if="errors.name" class="error">{{ errors.name }}</span>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="taxIdType">Tipo de NIF/CIF</label>
              <select id="taxIdType" v-model="form.taxIdType" @change="validateField('taxId')">
                <option value="NIF">NIF</option>
                <option value="CIF">CIF</option>
                <option value="VAT">VAT</option>
              </select>
            </div>

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
          </div>

          <div class="form-group">
            <label for="website">Sitio web</label>
            <input
              id="website"
              v-model="form.website"
              type="text"
              placeholder="example.com"
              @blur="validateField('website')"
            />
            <span v-if="errors.website" class="error">{{ errors.website }}</span>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label for="phone">Teléfono</label>
              <input
                id="phone"
                v-model="form.phone"
                type="tel"
                placeholder="p. ej., +34 123 456 789"
                @blur="validateField('phone')"
              />
              <span v-if="errors.phone" class="error">{{ errors.phone }}</span>
            </div>
            
            <div class="form-group">
              <label for="email">Email</label>
              <input
                id="email"
                v-model="form.email"
                type="email"
                placeholder="p. ej., contacto@empresa.com"
                @blur="validateField('email')"
              />
              <span v-if="errors.email" class="error">{{ errors.email }}</span>
            </div>
          </div>

          <div v-if="form.role === 'CLIENT' || form.role === 'BOTH'" class="form-group">
            <label for="defaultDiscount">Bonificación por defecto (%)</label>
            <input
              id="defaultDiscount"
              v-model.number="form.defaultDiscountPercentage"
              type="number"
              step="0.01"
              min="0"
              max="100"
              placeholder="0.00"
            />
            <small class="help-text">Porcentaje de descuento que se aplicará por defecto en las ventas a este cliente (0-100)</small>
          </div>
        </fieldset>
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
  entityType: '',
  name: '',
  firstName: '',
  lastName: '',
  role: '',
  taxId: '',
  taxIdType: 'NIF',
  website: '',
  phone: '',
  email: '',
  notes: '',
  defaultDiscountPercentage: 0,
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
    defaultDiscountPercentage: initial.defaultDiscountPercentage ?? initial.default_discount_percentage ?? 0,
  });

  if (isEditing.value && !form.role) {
    form.role = 'CLIENT';
  }
}

// Validation rules
const validationRules = {
  entityType: (value) => {
    if (!value) {
      return 'El tipo de entidad es obligatorio';
    }
    if (!['PERSON', 'ORGANIZATION'].includes(value)) {
      return 'Tipo de entidad inválido';
    }
    // Validate contact must be person
    if (form.role === 'CONTACT' && value !== 'PERSON') {
      return 'Los contactos solo pueden ser personas físicas';
    }
    return '';
  },
  firstName: (value) => {
    if (form.entityType === 'PERSON') {
      if (!value || value.trim().length === 0) {
        return 'El nombre es obligatorio';
      }
      if (value.length < 2) {
        return 'El nombre debe tener al menos 2 caracteres';
      }
      if (value.length > 50) {
        return 'El nombre no debe superar 50 caracteres';
      }
    }
    return '';
  },
  lastName: (value) => {
    if (form.entityType === 'PERSON') {
      if (!value || value.trim().length === 0) {
        return 'El apellido es obligatorio';
      }
      if (value.length < 2) {
        return 'El apellido debe tener al menos 2 caracteres';
      }
      if (value.length > 50) {
        return 'El apellido no debe superar 50 caracteres';
      }
    }
    return '';
  },
  name: (value) => {
    if (form.entityType === 'ORGANIZATION') {
      if (!value || value.trim().length === 0) {
        return 'El nombre es obligatorio';
      }
      if (value.length < 3) {
        return 'El nombre debe tener al menos 3 caracteres';
      }
      if (value.length > 100) {
        return 'El nombre no debe superar 100 caracteres';
      }
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
    // Validate contact must be person
    if (value === 'CONTACT' && form.entityType === 'ORGANIZATION') {
      return 'Los contactos solo pueden ser personas físicas';
    }
    return '';
  },
  taxId: (value) => {
    if (value && !isValidTaxId(value, form.taxIdType)) {
      if (form.taxIdType === 'NIF') {
        return 'Formato de NIF inválido (debe ser 8 dígitos seguidos de una letra)';
      } else if (form.taxIdType === 'CIF') {
        return 'Formato de CIF inválido (debe ser letra + 7 dígitos + dígito o letra)';
      } else if (form.taxIdType === 'VAT') {
        return 'Formato de VAT inválido (al menos 2 caracteres)';
      }
      return 'Formato inválido';
    }
    return '';
  },
  website: (value) => {
    if (value && !isValidUrl(value)) {
      return 'Formato de URL inválido';
    }
    return '';
  },
  phone: (value) => {
    if (value && value.trim()) {
      // Backend phone regex: ^[\+]?[\d\s\-()]{8,}$
      const phoneRegex = /^[\+]?[\d\s\-()]{8,}$/;
      if (!phoneRegex.test(value.trim())) {
        return 'Formato inválido. Debe tener al menos 8 dígitos y puede incluir +, espacios, guiones y paréntesis';
      }
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

// Helper functions
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
    const normalized = normalizeUrl(string);
    const url = new URL(normalized);
    
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
  // Backend email regex: ^[a-zA-Z0-9.+_%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$
  const emailRegex = /^[a-zA-Z0-9.+_%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(string.trim());
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

  // Additional validation: block if CONTACT role with ORGANIZATION type
  if (form.role === 'CONTACT' && form.entityType !== 'PERSON') {
    errorMessage.value = 'Error: Los contactos solo pueden ser personas físicas. Por favor, selecciona "Persona Física" como tipo de entidad.';
    return;
  }

  isSubmitting.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    let result;
    
    if (isEditing.value) {
      const updatePayload = {
        name: form.entityType === 'PERSON' 
          ? `${form.firstName} ${form.lastName}` 
          : form.name,
        role: form.role,
        taxId: form.taxId,
        taxIdType: form.taxIdType,
        website: normalizeUrl(form.website),
        phone: form.phone,
        email: form.email,
        notes: form.notes,
      };
      if (form.role === 'CLIENT' || form.role === 'BOTH') {
        updatePayload.default_discount_percentage = form.defaultDiscountPercentage || 0;
      }
      result = await partyApi.updateParty(props.partyId, updatePayload);
      successMessage.value = 'Entidad actualizada correctamente';
      emit('update', result);
    } else {
      const requestData = {
        id: form.id || generateId(),
        role: form.role,
        taxId: form.taxId,
        taxIdType: form.taxIdType,
        website: normalizeUrl(form.website),
        phone: form.phone,
        email: form.email,
        entityType: form.entityType,
      };

      if ((form.role === 'CLIENT' || form.role === 'BOTH') && form.defaultDiscountPercentage) {
        requestData.default_discount_percentage = form.defaultDiscountPercentage;
      }

      if (form.entityType === 'PERSON') {
        requestData.firstName = form.firstName;
        requestData.lastName = form.lastName;
        requestData.name = `${form.firstName} ${form.lastName}`;
      } else {
        requestData.name = form.name;
      }

      result = await partyApi.createParty(requestData);
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
  form.entityType = '';
  form.name = '';
  form.firstName = '';
  form.lastName = '';
  form.role = '';
  form.taxId = '';
  form.taxIdType = 'NIF';
  form.website = '';
  form.phone = '';
  form.email = '';
  form.notes = '';
  form.defaultDiscountPercentage = 0;
  Object.keys(errors).forEach((key) => {
    errors[key] = '';
  });
}

function onEntityTypeChange() {
  // Clear fields when entity type changes
  form.name = '';
  form.firstName = '';
  validateField('entityType');
}

function onRoleChange() {
  // If role is CONTACT, force entity type to PERSON
  if (form.role === 'CONTACT') {
    if (form.entityType !== 'PERSON') {
      form.entityType = 'PERSON';
      onEntityTypeChange();
    }
  }
  validateField('role');
  form.lastName = '';
  errors.name = '';
  errors.firstName = '';
  errors.lastName = '';
}

function generateId() {
  return crypto.randomUUID();
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

.help-text {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: #64748b;
  line-height: 1.4;
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
