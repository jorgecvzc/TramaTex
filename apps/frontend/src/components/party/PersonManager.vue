<template>
  <div class="person-manager">
    <div class="manager-header">
      <h3>Contactos</h3>
      <button @click="showForm = !showForm" class="btn btn-primary">
        {{ showForm ? '✕ Cerrar' : '+ Agregar contacto' }}
      </button>
    </div>

    <!-- Add/Edit Form -->
    <div v-if="showForm" class="form-section">
      <form @submit.prevent="submitForm">
        <div class="form-mode">
          <button
            type="button"
            class="btn btn-secondary"
            :class="{ active: formMode === 'new' }"
            @click="setFormMode('new')"
          >
            Nuevo contacto
          </button>
          <button
            type="button"
            class="btn btn-secondary"
            :class="{ active: formMode === 'existing' }"
            @click="setFormMode('existing')"
          >
            Contacto existente
          </button>
        </div>

        <div v-if="formMode === 'existing'" class="form-group">
          <label for="existingContact">Seleccionar persona física existente *</label>
          <select id="existingContact" v-model="selectedContactId" required>
            <option value="">-- Selecciona una persona física --</option>
            <option v-for="contact in availableContacts" :key="contact.id" :value="contact.id">
              {{ contact.first_name }} {{ contact.last_name }} {{ contact.email ? `(${contact.email})` : '' }}
            </option>
          </select>
          <small v-if="availableContacts.length === 0" class="hint">
            No hay personas físicas disponibles para vincular.
          </small>
        </div>

        <template v-if="formMode === 'new'">
        <div class="form-row">
          <div class="form-group">
            <label for="firstName">Nombre *</label>
            <input
              id="firstName"
              v-model="form.firstName"
              type="text"
              placeholder="Nombre"
              required
            />
          </div>
          <div class="form-group">
            <label for="lastName">Apellido *</label>
            <input
              id="lastName"
              v-model="form.lastName"
              type="text"
              placeholder="Apellido"
              required
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="email">Correo *</label>
            <input
              id="email"
              v-model="form.email"
              type="text"
              placeholder="correo@ejemplo.com"
              @blur="validateField('email')"
              required
            />
            <span v-if="formErrors.email" class="error">{{ formErrors.email }}</span>
          </div>
          <div class="form-group">
            <label for="phone">Teléfono</label>
            <input
              id="phone"
              v-model="form.phone"
              type="text"
              placeholder="+34 123 456 789"
              @blur="validateField('phone')"
            />
            <span v-if="formErrors.phone" class="error">{{ formErrors.phone }}</span>
          </div>
        </div>
        </template>

        <div v-if="formMode === 'new'" class="form-group">
          <label for="jobTitle">Cargo</label>
          <input
            id="jobTitle"
            v-model="form.jobTitle"
            type="text"
            placeholder="p. ej., Gerente"
          />
        </div>

        <div v-else class="form-group">
          <label for="jobTitleExisting">Cargo</label>
          <input
            id="jobTitleExisting"
            v-model="form.jobTitle"
            type="text"
            placeholder="p. ej., Gerente"
          />
        </div>

        <div class="form-group checkbox">
          <input
            id="isPrimary"
            v-model="form.isPrimary"
            type="checkbox"
          />
          <label for="isPrimary">Marcar como contacto principal</label>
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="isSubmitting" class="btn btn-primary">
            {{
              isSubmitting
                ? (formMode === 'existing' ? 'Vinculando...' : 'Agregando...')
                : (formMode === 'existing' ? 'Vincular contacto' : 'Agregar contacto')
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

    <!-- Contacts List -->
    <div v-if="persons.length > 0" class="persons-list">
      <div v-for="person in persons" :key="person.id" class="person-card">
        <div class="person-header">
          <div class="person-info">
            <h4>{{ person.first_name }} {{ person.last_name }}</h4>
            <p class="email">
              <span class="material-symbols-outlined icon-sm">mail</span>
              {{ person.email }}
            </p>
            <p v-if="person.phone" class="phone">
              <span class="material-symbols-outlined icon-sm">call</span>
              {{ person.phone }}
            </p>
            <button v-if="editingJobTitleId !== person.id" type="button" class="editable-chip" @click="startEditJobTitle(person)" :title="person.job_title ? 'Clic para editar cargo' : 'Clic para añadir cargo'">
              <span class="material-symbols-outlined chip-icon">work</span>
              <span :class="person.job_title ? '' : 'chip-placeholder'">{{ person.job_title || 'Añadir cargo...' }}</span>
              <span class="material-symbols-outlined chip-edit">edit</span>
            </button>
            <div v-else class="inline-edit">
              <span class="material-symbols-outlined chip-icon">work</span>
              <input
                ref="jobTitleInput"
                v-model="editingJobTitleValue"
                type="text"
                placeholder="p. ej., Gerente"
                class="inline-input"
                @keydown.enter="saveJobTitle(person)"
                @keydown.escape="cancelEditJobTitle"
                @blur="saveJobTitle(person)"
              />
            </div>
          </div>
          <div class="person-badges">
            <span v-if="person.is_primary" class="badge status-success">Principal</span>
            <span class="badge status-secondary">{{ formatDate(person.created_at) }}</span>
            <div class="action-buttons">
              <button
                type="button"
                class="btn-icon"
                @click="navigateToContact(person.id)"
                title="Ver detalles"
              >
                <span class="material-symbols-outlined">visibility</span>
              </button>
              <button
                type="button"
                class="btn-icon text-danger"
                :disabled="isRemovingId === person.id"
                @click="handleRemoveContact(person)"
                title="Eliminar"
              >
                <span class="material-symbols-outlined">delete</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="persons.length === 0 && !showForm" class="empty-state">
      <p>No hay contactos aún. Agrega el primero para empezar.</p>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Cargando contactos...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { partyApi } from '@/services/partyApi';

const router = useRouter();

const props = defineProps({
  partyId: {
    type: String,
    required: true,
  },
});

const persons = ref([]);
const isLoading = ref(false);
const isSubmitting = ref(false);
const isRemovingId = ref('');
const showForm = ref(false);
const formError = ref('');
const formMode = ref('new');
const availableContacts = ref([]);
const selectedContactId = ref('');
const editingJobTitleId = ref('');
const editingJobTitleValue = ref('');
const jobTitleInput = ref(null);
const isSavingJobTitle = ref(false);

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  jobTitle: '',
  isPrimary: false,
});

const formErrors = reactive({
  email: '',
  phone: '',
});

onMounted(() => {
  fetchPersons();
});

watch(() => props.partyId, () => {
  if (props.partyId) {
    fetchPersons();
  }
});

watch(showForm, (isVisible) => {
  if (isVisible) {
    loadAvailableContacts();
  }
});

async function fetchPersons() {
  if (!props.partyId) return;

  isLoading.value = true;
  try {
    const response = await partyApi.listContacts(props.partyId);
    persons.value = response.data || [];
  } catch (error) {
    formError.value = error?.message || 'No se pudieron cargar los contactos';
  } finally {
    isLoading.value = false;
  }
}

function validateEmail(email) {
  if (!email || !email.trim()) return true; // Email is optional
  const emailRegex = /^[a-zA-Z0-9.+_%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(email.trim());
}

function navigateToContact(partyId) {
  router.push({ name: 'PartyDetail', params: { id: partyId } });
}

function validatePhone(phone) {
  if (!phone || !phone.trim()) return true; // Phone is optional
  const phoneRegex = /^[\+]?[\d\s\-()]{8,}$/;
  return phoneRegex.test(phone.trim());
}

const validationRules = {
  email: (value) => {
    if (!value || !value.trim()) return ''; // Email is required but error is handled elsewhere
    if (!validateEmail(value)) {
      return 'Formato de email inválido';
    }
    return '';
  },
  phone: (value) => {
    if (value && value.trim() && !validatePhone(value)) {
      return 'Formato inválido. Debe tener al menos 8 dígitos y puede incluir +, espacios, guiones y paréntesis';
    }
    return '';
  },
};

function validateField(fieldName) {
  const validator = validationRules[fieldName];
  if (validator) {
    formErrors[fieldName] = validator(form[fieldName]);
  }
}

async function submitForm() {
  isSubmitting.value = true;
  formError.value = '';

  try {
    if (formMode.value === 'existing') {
      if (!selectedContactId.value) {
        formError.value = 'Selecciona un contacto existente';
        return;
      }

      const selected = availableContacts.value.find((contact) => contact.id === selectedContactId.value);

      await partyApi.linkExistingContact(props.partyId, selectedContactId.value, {
        jobTitle: form.jobTitle,
        email: selected?.email || '',
        phone: selected?.phone || '',
        isPrimary: form.isPrimary,
      });
    } else {
      // Validaciones obligatorias
      if (!form.firstName || !form.lastName || !form.email) {
        formError.value = 'Nombre, apellido y correo son obligatorios';
        return;
      }

      // Validar todos los campos
      validateField('email');
      validateField('phone');

      // Verificar si hay errores
      if (formErrors.email || formErrors.phone) {
        formError.value = 'Por favor, corrija los errores en el formulario antes de guardar';
        return;
      }

      await partyApi.addContact(props.partyId, {
        id: `person-${Date.now()}`,
        firstName: form.firstName,
        lastName: form.lastName,
        email: form.email,
        phone: form.phone,
        jobTitle: form.jobTitle,
        isPrimary: form.isPrimary,
      });
    }

    resetForm();
    await fetchPersons();
    await loadAvailableContacts();
  } catch (error) {
    formError.value = error?.message || 'No se pudo agregar el contacto';
  } finally {
    isSubmitting.value = false;
  }
}

function setFormMode(mode) {
  formMode.value = mode;
  formError.value = '';
}

async function loadAvailableContacts() {
  if (!props.partyId) {
    return;
  }

  try {
    availableContacts.value = await partyApi.listAvailableContactsForParty(props.partyId);
  } catch (error) {
    console.error('Error al cargar contactos disponibles:', error);
    formError.value = 'No se pudieron cargar los contactos disponibles';
    availableContacts.value = [];
  }
}

async function handleRemoveContact(person) {
  if (!person?.id) {
    return;
  }

  // Check if this contact has other references
  let hasOtherReferences = false;
  try {
    const contactParty = await partyApi.getParty(person.id);
    const relationships = await partyApi.listRelationships(person.id);
    // Check if contact has other employment relationships
    hasOtherReferences = relationships.some(rel => 
      rel.type === 'IS_EMPLOYEE_OF' && rel.to_party_id !== props.partyId
    );
  } catch (error) {
    // Assume it has references if we can't check
    hasOtherReferences = true;
  }

  let message = `¿Eliminar el contacto ${person.first_name} ${person.last_name}?`;
  if (!hasOtherReferences) {
    message = `¿Eliminar el contacto ${person.first_name} ${person.last_name}?\n\nEste contacto no está asociado a otras entidades y será eliminado completamente del sistema.`;
  } else {
    message = `¿Eliminar el contacto ${person.first_name} ${person.last_name}?\n\nEste contacto será desvinculado de esta entidad, pero permanecerá en el sistema porque está asociado a otras entidades.`;
  }

  const confirmed = window.confirm(message);
  if (!confirmed) {
    return;
  }

  isRemovingId.value = person.id;
  formError.value = '';

  try {
    // Pass true to delete the party if it has no other references
    await partyApi.removeContact(props.partyId, person.id, !hasOtherReferences);
    await fetchPersons();
  } catch (error) {
    formError.value = error?.message || 'No se pudo eliminar el contacto';
  } finally {
    isRemovingId.value = '';
  }
}

function startEditJobTitle(person) {
  editingJobTitleId.value = person.id;
  editingJobTitleValue.value = person.job_title || '';
  nextTick(() => {
    if (jobTitleInput.value) {
      const el = Array.isArray(jobTitleInput.value) ? jobTitleInput.value[0] : jobTitleInput.value;
      el?.focus();
    }
  });
}

function cancelEditJobTitle() {
  editingJobTitleId.value = '';
  editingJobTitleValue.value = '';
}

async function saveJobTitle(person) {
  if (isSavingJobTitle.value) return;
  if (editingJobTitleId.value !== person.id) return;

  const newValue = editingJobTitleValue.value.trim();
  if (newValue === (person.job_title || '')) {
    cancelEditJobTitle();
    return;
  }

  let contactDetailsId = person.contact_details_id;

  // If contact_details_id is not cached, look it up from the API
  if (!contactDetailsId) {
    try {
      const contacts = await partyApi.listContactDetails(props.partyId);
      const match = contacts.find((c) => c.related_party_id === person.id);
      if (match?.id) {
        contactDetailsId = match.id;
        person.contact_details_id = match.id;
      }
    } catch {
      // ignore lookup failure
    }
  }

  if (!contactDetailsId) {
    // No contact-details record exists yet — create one
    isSavingJobTitle.value = true;
    try {
      const created = await partyApi.createContactDetails(props.partyId, {
        type_description: newValue || 'Contacto',
        related_party_id: person.id,
      });
      if (created?.id) {
        person.contact_details_id = created.id;
      }
      person.job_title = newValue || 'Contacto';
      cancelEditJobTitle();
    } catch (error) {
      formError.value = error?.message || 'No se pudo crear el cargo';
      cancelEditJobTitle();
    } finally {
      isSavingJobTitle.value = false;
    }
    return;
  }

  isSavingJobTitle.value = true;
  try {
    await partyApi.updateContactJobTitle(props.partyId, contactDetailsId, newValue || 'Contacto');
    person.job_title = newValue || 'Contacto';
    cancelEditJobTitle();
  } catch (error) {
    formError.value = error?.message || 'No se pudo actualizar el cargo';
    cancelEditJobTitle();
  } finally {
    isSavingJobTitle.value = false;
  }
}

function resetForm() {
  form.firstName = '';
  form.lastName = '';
  form.email = '';
  form.phone = '';
  form.jobTitle = '';
  form.isPrimary = false;
  selectedContactId.value = '';
  formMode.value = 'new';
  formError.value = '';
  formErrors.email = '';
  formErrors.phone = '';
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
.person-manager {
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

.form-group input:not([type="checkbox"]), .form-group select {
  padding: 0.75rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-md);
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  background: white;
  transition: all 0.2s;
}

.form-group input:focus, .form-group select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.form-group.checkbox {
  flex-direction: row;
  align-items: center;
  gap: var(--spacing-sm);
  margin: var(--spacing-md) 0;
}

.form-group.checkbox input {
  width: 18px;
  height: 18px;
}

.form-group .error {
  color: var(--color-error);
  font-size: var(--font-size-xs);
}

.form-actions {
  display: flex;
  gap: var(--spacing-md);
  margin-top: var(--spacing-md);
}

.form-mode {
  display: flex;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-md);
}

.persons-list {
  display: grid;
  gap: var(--spacing-md);
}

.person-card {
  padding: var(--spacing-md);
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-md);
  background: white;
  transition: all 0.2s ease;
}

.person-card:hover {
  border-color: var(--color-border-strong);
  box-shadow: var(--box-shadow-md);
}

.person-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-md);
}

.person-info h4 {
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  font-weight: 700;
}

.person-info p {
  color: var(--color-text-secondary);
  margin: 0.25rem 0;
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.icon-sm { font-size: 18px; color: var(--color-text-secondary); }

.person-badges {
  display: flex;
  gap: 0.5rem;
  align-items: center;
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

.spinner {
  width: 24px; height: 24px; border: 2px solid var(--color-border);
  border-top-color: var(--color-primary); border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 768px) {
  .form-row { grid-template-columns: 1fr; }
  .person-header { flex-direction: column; }
  .form-actions { flex-direction: column; }
}

.editable-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.3rem 0.6rem;
  margin: 0.25rem 0;
  border: 1px dashed var(--color-border);
  border-radius: var(--border-radius-md);
  background: var(--color-background);
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: all 0.15s ease;
}

.editable-chip:hover {
  border-color: var(--color-primary);
  background: rgba(230, 184, 0, 0.06);
  color: var(--color-text-primary);
}

.chip-icon {
  font-size: 16px;
}

.chip-edit {
  font-size: 14px;
  opacity: 0.5;
}

.editable-chip:hover .chip-edit {
  opacity: 1;
  color: var(--color-primary);
}

.chip-placeholder {
  font-style: italic;
}

.inline-edit {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0.25rem 0;
}

.inline-input {
  padding: 0.3rem 0.5rem;
  border: 1px solid var(--color-primary);
  border-radius: var(--border-radius-md);
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  background: white;
  outline: none;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
  width: 200px;
}
</style>
