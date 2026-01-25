<template>
  <div class="person-manager">
    <div class="manager-header">
      <h3>Contacts & People</h3>
      <button @click="showForm = !showForm" class="btn btn-primary">
        {{ showForm ? '✕ Close' : '+ Add Contact' }}
      </button>
    </div>

    <!-- Add/Edit Form -->
    <div v-if="showForm" class="form-section">
      <form @submit.prevent="submitForm">
        <div class="form-row">
          <div class="form-group">
            <label for="firstName">First Name *</label>
            <input
              id="firstName"
              v-model="form.firstName"
              type="text"
              placeholder="First name"
              required
            />
          </div>
          <div class="form-group">
            <label for="lastName">Last Name *</label>
            <input
              id="lastName"
              v-model="form.lastName"
              type="text"
              placeholder="Last name"
              required
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="email">Email *</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="email@example.com"
              required
            />
          </div>
          <div class="form-group">
            <label for="phone">Phone</label>
            <input
              id="phone"
              v-model="form.phone"
              type="tel"
              placeholder="+34 123 456 789"
            />
          </div>
        </div>

        <div class="form-group">
          <label for="jobTitle">Job Title</label>
          <input
            id="jobTitle"
            v-model="form.jobTitle"
            type="text"
            placeholder="e.g., Manager"
          />
        </div>

        <div class="form-group checkbox">
          <input
            id="isPrimary"
            v-model="form.isPrimary"
            type="checkbox"
          />
          <label for="isPrimary">Mark as primary contact</label>
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="isSubmitting" class="btn btn-primary">
            {{ isSubmitting ? 'Adding...' : 'Add Contact' }}
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
            <span v-if="person.is_primary" class="badge primary">Primary</span>
            <span class="badge date">{{ formatDate(person.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="persons.length === 0 && !showForm" class="empty-state">
      <p>No contacts yet. Add your first contact to get started.</p>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Loading contacts...</p>
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

const persons = ref([]);
const isLoading = ref(false);
const isSubmitting = ref(false);
const showForm = ref(false);
const formError = ref('');

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

watch(() => props.organizationId, () => {
  if (props.organizationId) {
    fetchPersons();
  }
});

async function fetchPersons() {
  if (!props.organizationId) return;

  isLoading.value = true;
  try {
    const response = await partyApi.listPersons(props.organizationId);
    persons.value = response.data || [];
  } catch (error) {
    formError.value = error.message || 'Failed to load contacts';
  } finally {
    isLoading.value = false;
  }
}

async function submitForm() {
  if (!form.firstName || !form.lastName || !form.email) {
    formError.value = 'First name, last name, and email are required';
    return;
  }

  isSubmitting.value = true;
  formError.value = '';

  try {
    await partyApi.addPerson(props.organizationId, {
      id: `person-${Date.now()}`,
      firstName: form.firstName,
      lastName: form.lastName,
      email: form.email,
      phone: form.phone,
      jobTitle: form.jobTitle,
      isPrimary: form.isPrimary,
    });

    resetForm();
    await fetchPersons();
  } catch (error) {
    formError.value = error.message || 'Failed to add contact';
  } finally {
    isSubmitting.value = false;
  }
}

function resetForm() {
  form.firstName = '';
  form.lastName = '';
  form.email = '';
  form.phone = '';
  form.jobTitle = '';
  form.isPrimary = false;
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
}

.form-group label {
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
}

.form-group input:not([type="checkbox"]) {
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

.persons-list {
  display: grid;
  gap: 1rem;
}

.person-card {
  padding: 1rem;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: var(--color-background);
  transition: all 0.2s ease;
}

.person-card:hover {
  border-color: var(--primary-color);
  box-shadow: 0 2px 8px rgba(230, 184, 0, 0.1);
}

.person-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.person-info h4 {
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
}

.person-info p {
  color: var(--color-text-secondary);
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

  .person-header {
    flex-direction: column;
    gap: 1rem;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
