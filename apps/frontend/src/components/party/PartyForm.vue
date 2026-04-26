<template>
  <div class="party-form">
    <header v-if="!hideHeader" class="form-header">
      <span class="material-symbols-outlined header-icon">{{ isEditing ? 'edit_note' : 'person_add' }}</span>
      <div>
        <h1>{{ isEditing ? 'Editar entidad' : 'Crear entidad' }}</h1>
        <p class="subtitle">{{ isEditing ? `ID: ${props.partyId}` : 'Completa los datos para dar de alta una nueva entidad' }}</p>
      </div>
    </header>

    <form @submit.prevent="submitForm" class="form-layout">
      <!-- Card: Basic Identification -->
      <section class="card form-section">
        <div class="section-header">
          <span class="material-symbols-outlined">badge</span>
          <h2>Identificación y Rol</h2>
        </div>
        
        <div class="section-body">
          <div class="form-row">
            <div class="form-group">
              <label for="role">
                Rol de la entidad *
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
              <span v-if="errors.role" class="error-msg">{{ errors.role }}</span>
            </div>

            <div class="form-group">
              <label for="entityType">
                Tipo de entidad *
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
              <span v-if="errors.entityType" class="error-msg">{{ errors.entityType }}</span>
            </div>
          </div>

          <div v-if="form.role === 'CONTACT' || (form.role === 'CONTACT' && form.entityType !== 'PERSON')" class="alert-info">
            <span class="material-symbols-outlined">info</span>
            <p>Los contactos deben ser únicamente <strong>personas físicas</strong>.</p>
          </div>

          <!-- Conditional fields based on entityType -->
          <div v-if="form.entityType === 'PERSON'" class="form-row mt-4">
            <div class="form-group">
              <label for="firstName">Nombre *</label>
              <input
                id="firstName"
                v-model="form.firstName"
                type="text"
                placeholder="Ej: Juan"
                required
                @blur="validateField('firstName')"
              />
              <span v-if="errors.firstName" class="error-msg">{{ errors.firstName }}</span>
            </div>

            <div class="form-group">
              <label for="lastName">Apellido(s) *</label>
              <input
                id="lastName"
                v-model="form.lastName"
                type="text"
                placeholder="Ej: García López"
                required
                @blur="validateField('lastName')"
              />
              <span v-if="errors.lastName" class="error-msg">{{ errors.lastName }}</span>
            </div>
          </div>

          <div v-else-if="form.entityType === 'ORGANIZATION'" class="form-group mt-4">
            <label for="name">Nombre de la organización *</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              placeholder="Ej: Acme Corporation S.L."
              required
              @blur="validateField('name')"
            />
            <span v-if="errors.name" class="error-msg">{{ errors.name }}</span>
          </div>
        </div>
      </section>

      <!-- Card: Legal and Contact -->
      <section class="card form-section">
        <div class="section-header">
          <span class="material-symbols-outlined">contact_mail</span>
          <h2>Datos Legales y Contacto</h2>
        </div>
        
        <div class="section-body">
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
              <label for="taxId">Número de identificación</label>
              <input
                id="taxId"
                v-model="form.taxId"
                type="text"
                placeholder="p. ej., 12345678A"
                class="text-mono"
                @blur="validateField('taxId')"
              />
              <span v-if="errors.taxId" class="error-msg">{{ errors.taxId }}</span>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="phone">Teléfono de contacto</label>
              <div class="input-with-icon">
                <span class="material-symbols-outlined icon-start">call</span>
                <input
                  id="phone"
                  v-model="form.phone"
                  type="tel"
                  placeholder="+34 000 000 000"
                  @blur="validateField('phone')"
                />
              </div>
              <span v-if="errors.phone" class="error-msg">{{ errors.phone }}</span>
            </div>
            
            <div class="form-group">
              <label for="email">Correo electrónico</label>
              <div class="input-with-icon">
                <span class="material-symbols-outlined icon-start">mail</span>
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  placeholder="ejemplo@correo.com"
                  @blur="validateField('email')"
                />
              </div>
              <span v-if="errors.email" class="error-msg">{{ errors.email }}</span>
            </div>
          </div>

          <div class="form-group">
            <label for="website">Sitio web / URL</label>
            <div class="input-with-icon">
              <span class="material-symbols-outlined icon-start">language</span>
              <input
                id="website"
                v-model="form.website"
                type="text"
                placeholder="www.empresa.com"
                @blur="validateField('website')"
              />
            </div>
            <span v-if="errors.website" class="error-msg">{{ errors.website }}</span>
          </div>
        </div>
      </section>

      <!-- Card: Configuration and Notes -->
      <section class="card form-section">
        <div class="section-header">
          <span class="material-symbols-outlined">settings_suggest</span>
          <h2>Configuración y Notas</h2>
        </div>
        
        <div class="section-body">
          <div v-if="form.role === 'CLIENT' || form.role === 'BOTH'" class="form-group">
            <label for="defaultDiscount">Bonificación comercial por defecto (%)</label>
            <div class="input-with-icon">
              <span class="material-symbols-outlined icon-start">percent</span>
              <input
                id="defaultDiscount"
                v-model.number="form.defaultDiscountPercentage"
                type="number"
                step="0.01"
                min="0"
                max="100"
                placeholder="0.00"
              />
            </div>
            <p class="help-text">Descuento automático que se aplicará en presupuestos y pedidos.</p>
          </div>

          <div class="form-group">
            <label for="notes">Notas internas</label>
            <textarea
              id="notes"
              v-model="form.notes"
              placeholder="Notas privadas sobre la entidad..."
              rows="4"
            />
          </div>
        </div>
      </section>

      <!-- Form Actions -->
      <footer v-if="!hideActions" class="form-footer">
        <button
          type="button"
          @click="resetForm"
          class="btn btn-outline"
          :disabled="isSubmitting"
        >
          <span class="material-symbols-outlined">restart_alt</span>
          Reiniciar
        </button>
        
        <button
          type="submit"
          :disabled="isSubmitting"
          class="btn btn-primary btn-grow"
        >
          <span class="material-symbols-outlined">{{ isSubmitting ? 'sync' : 'save' }}</span>
          <span>{{ isSubmitting ? (isEditing ? 'Actualizando...' : 'Creando...') : (isEditing ? 'Actualizar Entidad' : 'Crear Entidad') }}</span>
        </button>
      </footer>
    </form>

    <!-- Feedback Messages -->
    <Transition name="fade">
      <div v-if="successMessage" class="feedback-toast success">
        <span class="material-symbols-outlined">check_circle</span>
        <p>{{ successMessage }}</p>
        <button @click="successMessage = ''" class="toast-close">&times;</button>
      </div>
    </Transition>
    
    <Transition name="fade">
      <div v-if="errorMessage" class="feedback-toast error">
        <span class="material-symbols-outlined">error</span>
        <p>{{ errorMessage }}</p>
        <button @click="errorMessage = ''" class="toast-close">&times;</button>
      </div>
    </Transition>
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
  hideActions: {
    type: Boolean,
    default: false,
  },
  hideHeader: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['submit', 'update']);

// Expose methods for parent components (like header buttons)
defineExpose({
  submitForm,
  resetForm
});

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
    if (!value) return 'El tipo de entidad es obligatorio';
    if (!['PERSON', 'ORGANIZATION'].includes(value)) return 'Tipo de entidad inválido';
    if (form.role === 'CONTACT' && value !== 'PERSON') return 'Los contactos solo pueden ser personas físicas';
    return '';
  },
  firstName: (value) => {
    if (form.entityType === 'PERSON') {
      if (!value || value.trim().length === 0) return 'El nombre es obligatorio';
      if (value.length < 2) return 'El nombre debe tener al menos 2 caracteres';
    }
    return '';
  },
  lastName: (value) => {
    if (form.entityType === 'PERSON') {
      if (!value || value.trim().length === 0) return 'El apellido es obligatorio';
    }
    return '';
  },
  name: (value) => {
    if (form.entityType === 'ORGANIZATION') {
      if (!value || value.trim().length === 0) return 'El nombre es obligatorio';
      if (value.trim().length < 3) return 'El nombre debe tener al menos 3 caracteres';
    }
    return '';
  },
  role: (value) => {
    if (!value) return 'El rol es obligatorio';
    return '';
  },
  taxId: (value) => {
    if (value && !isValidTaxId(value, form.taxIdType)) {
			return 'Formato inválido de NIF/CIF/VAT';
    }
    return '';
  },
  website: (value) => {
    if (value && !isValidUrl(value)) return 'URL inválido';
    return '';
  },
  phone: (value) => {
    if (value && value.trim()) {
      const phoneRegex = /^[\+]?[\d\s\-()]{8,}$/;
      if (!phoneRegex.test(value.trim())) return 'Formato de teléfono inválido';
    }
    return '';
  },
  email: (value) => {
    if (value && value.trim() && !isValidEmail(value)) return 'Email no válido';
    return '';
  },
};

// Helper functions
function normalizeUrl(url) {
  if (!url || !url.trim()) return url;
  const trimmedUrl = url.trim();
  if (/^https?:\/\//i.test(trimmedUrl)) return trimmedUrl;
  return `https://${trimmedUrl}`;
}

function isValidUrl(string) {
  if (!string || !string.trim()) return true;
  try {
    const normalized = normalizeUrl(string);
    const url = new URL(normalized);
    return url.hostname.includes('.');
  } catch (_) {
    return false;
  }
}

function isValidEmail(string) {
  const emailRegex = /^[a-zA-Z0-9.+_%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(string.trim());
}

function isValidTaxId(taxId, taxIdType) {
  if (!taxId || !taxId.trim()) return true;
  const trimmed = taxId.trim().toUpperCase();
  if (taxIdType === 'NIF') return /^[0-9]{8}[A-Z]$/.test(trimmed);
  if (taxIdType === 'CIF') return /^[A-Z][0-9]{7}[0-9A-Z]$/.test(trimmed);
  return trimmed.length >= 2;
}

function validateField(fieldName) {
  const validator = validationRules[fieldName];
  if (validator) errors[fieldName] = validator(form[fieldName]);
}

function validateForm() {
  Object.keys(validationRules).forEach((field) => validateField(field));
  return Object.values(errors).every((err) => !err);
}

async function submitForm() {
  if (!validateForm()) {
    errorMessage.value = 'Corrige los errores antes de continuar';
    return;
  }

  isSubmitting.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    let result;
    if (isEditing.value) {
      const updatePayload = {
        name: form.entityType === 'PERSON' ? `${form.firstName} ${form.lastName}` : form.name,
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
      successMessage.value = 'Cambios guardados con éxito';
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
        notes: form.notes,
      };
      if (form.role === 'CLIENT' || form.role === 'BOTH') {
        requestData.default_discount_percentage = form.defaultDiscountPercentage ?? 0;
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
    errorMessage.value = error?.data?.message || error?.message || 'Error al procesar la solicitud';
  } finally {
    isSubmitting.value = false;
  }
}

function resetForm() {
  Object.assign(form, {
    id: '', entityType: '', name: '', firstName: '', lastName: '', role: '',
    taxId: '', taxIdType: 'NIF', website: '', phone: '', email: '', notes: '',
    defaultDiscountPercentage: 0
  });
  Object.keys(errors).forEach(k => errors[k] = '');
}

function onEntityTypeChange() {
  form.name = '';
  form.firstName = '';
  // Auto-select the most common tax ID type for the chosen entity type
  if (form.entityType === 'ORGANIZATION') {
    form.taxIdType = 'CIF';
  } else if (form.entityType === 'PERSON') {
    form.taxIdType = 'NIF';
  }
  validateField('entityType');
  validateField('taxId');
}

function onRoleChange() {
  if (form.role === 'CONTACT' && form.entityType !== 'PERSON') {
    form.entityType = 'PERSON';
    onEntityTypeChange();
  }
  validateField('role');
}

function generateId() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16);
  });
}
</script>

<style scoped>
.party-form {
  max-width: 1300px;
  margin: 0 auto;
  width: 100%;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  margin-bottom: 2rem;
}

.header-icon {
  font-size: 3rem;
  color: var(--color-primary);
  background: rgba(230, 184, 0, 0.1);
  padding: 0.75rem;
  border-radius: 12px;
}

.form-header h1 {
  font-size: 1.75rem;
  color: var(--color-text-primary);
  margin: 0;
  font-family: var(--font-family-brand);
}

.subtitle {
  color: var(--color-text-secondary);
  margin: 0.25rem 0 0;
}

.form-layout {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.card {
  background: var(--color-surface);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--color-border);
  box-shadow: var(--box-shadow-sm);
  overflow: hidden;
}

.section-header {
  padding: 1.25rem 1.5rem;
  background: var(--color-background);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.section-header h2 {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  color: var(--color-text-primary);
}

.section-header .material-symbols-outlined {
  font-size: 20px;
  color: var(--color-text-secondary);
}

.section-body {
  padding: 1.5rem;
}

/* Groups and Rows */
.form-group {
  margin-bottom: 1.25rem;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
}

.mt-4 { margin-top: 1rem; }

label {
  display: block;
  font-size: var(--font-size-xs);
  font-weight: 600;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  margin-bottom: 0.5rem;
  letter-spacing: 0.05em;
}

/* Inputs */
input, select, textarea {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: var(--border-radius-md);
  border: 1px solid var(--color-border);
  background: white;
  color: var(--color-text-primary);
  font-size: var(--font-size-sm);
  transition: all 0.2s;
  font-family: inherit;
}

input:focus, select:focus, textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

input:disabled, select:disabled {
  background-color: var(--color-background);
  cursor: not-allowed;
  opacity: 0.7;
}

textarea { resize: vertical; }

.text-mono { font-family: var(--font-family-mono); }

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
}

.icon-start {
  position: absolute;
  left: 0.85rem;
  font-size: 20px;
  color: var(--color-text-secondary);
}

.input-with-icon input { padding-left: 2.75rem; }

/* Feedback */
.error-msg {
  color: var(--color-error);
  font-size: 0.75rem;
  margin-top: 0.4rem;
  display: block;
  font-weight: 500;
}

.help-text {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-top: 0.5rem;
}

.alert-info {
  background: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.2);
  padding: 0.75rem 1rem;
  border-radius: var(--border-radius-md);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 1rem;
}

.alert-info p { margin: 0; font-size: 0.85rem; color: #1e40af; }
.alert-info .material-symbols-outlined { color: #3b82f6; font-size: 20px; }

/* Footer Actions */
.form-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1rem 0;
  margin-top: 1rem;
}

.btn-grow {
  flex: 2;
}

@media (max-width: 768px) {
  .form-row { grid-template-columns: 1fr; }
  .form-footer { flex-direction: column-reverse; }
  .btn { width: 100%; }
}

/* Toasts */
.feedback-toast {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  padding: 1rem 1.5rem;
  border-radius: var(--border-radius-lg);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  box-shadow: var(--box-shadow-lg);
  z-index: 1000;
  min-width: 300px;
}

.feedback-toast.success { background: #16a34a; color: white; }
.feedback-toast.error { background: #dc2626; color: white; }
.feedback-toast p { margin: 0; font-weight: 500; flex: 1; }
.toast-close { background: none; border: none; color: white; font-size: 1.5rem; cursor: pointer; }

/* Transitions */
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s, transform 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(10px); }

@media (max-width: 768px) {
  .form-row { grid-template-columns: 1fr; }
  .form-footer { flex-direction: column-reverse; }
  .btn { width: 100%; }
}
</style>
