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
          <label for="existingContact">Seleccionar contacto existente *</label>
          <select id="existingContact" v-model="selectedContactId" required>
            <option value="">-- Selecciona un contacto --</option>
            <option v-for="contact in availableContacts" :key="contact.id" :value="contact.id">
              {{ contact.first_name }} {{ contact.last_name }} {{ contact.email ? `(${contact.email})` : '' }}
            </option>
          </select>
          <small v-if="availableContacts.length === 0" class="hint">
            No hay contactos disponibles para vincular.
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
              type="email"
              placeholder="correo@ejemplo.com"
              required
            />
          </div>
          <div class="form-group">
            <label for="phone">Teléfono</label>
            <input
              id="phone"
              v-model="form.phone"
              type="tel"
              placeholder="+34 123 456 789"
            />
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
            <p class="email">📧 {{ person.email }}</p>
            <p v-if="person.phone" class="phone">📞 {{ person.phone }}</p>
            <p v-if="person.job_title" class="job">💼 {{ person.job_title }}</p>
          </div>
          <div class="person-badges">
            <span v-if="person.is_primary" class="badge primary">Principal</span>
            <span class="badge date">{{ formatDate(person.created_at) }}</span>
            <button
              type="button"
              class="btn btn-danger"
              :disabled="isRemovingId === person.id"
              @click="handleRemoveContact(person)"
            >
              {{ isRemovingId === person.id ? 'Eliminando...' : 'Eliminar' }}
            </button>
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
import { ref, reactive, onMounted, watch } from 'vue';
import { partyApi } from '@/services/partyApi';

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

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  jobTitle: '',
  isPrimary: false,
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
      if (!form.firstName || !form.lastName || !form.email) {
        formError.value = 'Nombre, apellido y correo son obligatorios';
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
    console.log('Contactos disponibles:', availableContacts.value.length);
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

.form-group input:not([type="checkbox"]) {
  padding: 0.6rem 0.8rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.9rem;
  color: #1e293b;
}

.form-group select {
  padding: 0.6rem 0.8rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.9rem;
  color: #1e293b;
  background: #ffffff;
}

.form-group input:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
}

.form-group select:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
}

.form-group.checkbox {
  flex-direction: row;
  align-items: center;
  margin: 1rem 0;
}

.form-group.checkbox input {
  width: 20px;
  height: 20px;
  margin-right: 0.5rem;
}

.form-group.checkbox label {
  margin: 0;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.form-mode {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.form-mode .btn.active {
  border-color: #002395;
  box-shadow: 0 0 0 2px rgba(0, 35, 149, 0.12);
}

.hint {
  margin-top: 0.5rem;
  color: #64748b;
  font-size: 0.85rem;
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

.btn-danger {
  background: #ffffff;
  border: 1px solid #fecaca;
  color: #b91c1c;
}

.btn-danger:hover:not(:disabled) {
  background: #fef2f2;
}

.error-message {
  color: #991b1b;
  background-color: #fee2e2;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-top: 1rem;
  border: 1px solid #ef4444;
}

.persons-list {
  display: grid;
  gap: 1rem;
}

.person-card {
  padding: 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.person-card:hover {
  border-color: #002395;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.08);
}

.person-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.person-info h4 {
  color: #1e293b;
  margin: 0 0 0.5rem 0;
}

.person-info p {
  color: #64748b;
  margin: 0.25rem 0;
  font-size: 0.9rem;
}

.person-badges {
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

  .person-header {
    flex-direction: column;
    gap: 1rem;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
