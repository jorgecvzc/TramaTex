<template>
  <div class="party-selector">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    
    <div class="selector-container">
      <!-- Search Mode: Dropdown with search -->
      <div class="search-selector">
        <input
          :id="inputId"
          v-model="searchTerm"
          type="text"
          :placeholder="placeholder || 'Buscar por nombre...'"
          class="form-input"
          @input="handleSearch"
          @focus="showDropdown = true"
          @blur="handleBlur"
          @keydown.enter.prevent="selectFirst"
          @keydown.down.prevent="navigateDown"
          @keydown.up.prevent="navigateUp"
          :required="required"
        />
        
        <!-- Dropdown Results -->
        <div v-if="showDropdown && (filteredParties.length > 0 || isSearching)" class="dropdown-results">
          <div v-if="isSearching" class="dropdown-item loading">
            <span class="spinner-small"></span>
            Buscando...
          </div>
          <div
            v-for="(party, index) in filteredParties"
            :key="party.id"
            :class="['dropdown-item', { active: index === activeIndex, selected: party.id === modelValue }]"
            @mousedown.prevent="selectParty(party)"
            @mouseenter="activeIndex = index"
          >
            <div class="party-info">
              <span class="party-name">{{ party.name }}</span>
              <span v-if="party.tax_id" class="party-tax">{{ party.tax_id }}</span>
              <span class="party-role">{{ getRoleLabel(party.role) }}</span>
            </div>
            <span v-if="party.id === modelValue" class="selected-indicator">✓</span>
          </div>
          <div v-if="!isSearching && filteredParties.length === 0" class="dropdown-item empty">
            No se encontraron resultados
          </div>
        </div>
      </div>
      
      <!-- Selected Party Display -->
      <div v-if="selectedParty && !showDropdown" class="selected-party">
        <div class="selected-party-info">
          <span class="party-name">{{ selectedParty.name }}</span>
          <span v-if="selectedParty.tax_id" class="party-detail">{{ selectedParty.tax_id }}</span>
        </div>
        <button
          type="button"
          class="btn-clear"
          @click="clearSelection"
          title="Limpiar selección"
        >
          ✕
        </button>
      </div>
    </div>
    
    <!-- Hidden input for form compatibility -->
    <input type="hidden" :value="modelValue" :name="name" />
    
    <!-- Help Text -->
    <span v-if="helpText" class="help-text">{{ helpText }}</span>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import partyApi from '@/services/partyApi';

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: 'Buscar cliente...',
  },
  required: {
    type: Boolean,
    default: false,
  },
  roleFilter: {
    type: String,
    default: 'CLIENT', // CLIENT, SUPPLIER, BOTH, or null for all
  },
  name: {
    type: String,
    default: 'partyId',
  },
  helpText: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:modelValue', 'select']);

// Component state
const searchTerm = ref('');
const allParties = ref([]);
const showDropdown = ref(false);
const isSearching = ref(false);
const activeIndex = ref(0);
const inputId = computed(() => `party-selector-${Math.random().toString(36).substr(2, 9)}`);

const selectedParty = computed(() => {
  if (!props.modelValue) return null;
  return allParties.value.find(p => p.id === props.modelValue) || null;
});

const filteredParties = computed(() => {
  if (!searchTerm.value.trim()) {
    return allParties.value.slice(0, 50); // Show first 50 if no search
  }
  
  const term = searchTerm.value.toLowerCase();
  return allParties.value.filter(p => 
    p.name.toLowerCase().includes(term) ||
    (p.tax_id && p.tax_id.toLowerCase().includes(term)) ||
    p.id.toLowerCase().includes(term)
  ).slice(0, 50);
});

// Methods
async function loadParties() {
  isSearching.value = true;
  try {
    const filters = {};
    if (props.roleFilter) {
      filters.role = props.roleFilter;
    }
    filters.pageSize = 500; // Load more parties for search
    
    const response = await partyApi.listParties(filters);
    allParties.value = response.data || [];
    
    // If modelValue is set, load that party's name into search
    if (props.modelValue && selectedParty.value) {
      searchTerm.value = selectedParty.value.name;
    }
  } catch (error) {
    console.error('Error loading parties:', error);
  } finally {
    isSearching.value = false;
  }
}

function handleSearch() {
  showDropdown.value = true;
  activeIndex.value = 0;
}

function handleBlur() {
  // Delay to allow click event on dropdown
  setTimeout(() => {
    showDropdown.value = false;
    // Reset search term to selected party name or clear
    if (selectedParty.value) {
      searchTerm.value = selectedParty.value.name;
    } else if (!props.modelValue) {
      searchTerm.value = '';
    }
  }, 200);
}

function selectParty(party) {
  emit('update:modelValue', party.id);
  emit('select', party);
  searchTerm.value = party.name;
  showDropdown.value = false;
  activeIndex.value = 0;
}

function selectFirst() {
  if (filteredParties.value.length > 0) {
    selectParty(filteredParties.value[0]);
  }
}

function navigateDown() {
  if (activeIndex.value < filteredParties.value.length - 1) {
    activeIndex.value++;
  }
}

function navigateUp() {
  if (activeIndex.value > 0) {
    activeIndex.value--;
  }
}

function clearSelection() {
  emit('update:modelValue', '');
  emit('select', null);
  searchTerm.value = '';
  showDropdown.value = false;
}

function getRoleLabel(role) {
  const labels = {
    CLIENT: 'Cliente',
    SUPPLIER: 'Proveedor',
    BOTH: 'Cliente/Proveedor',
  };
  return labels[role] || role;
}

// Watch for external changes to modelValue
watch(() => props.modelValue, (newVal) => {
  if (newVal && selectedParty.value) {
    searchTerm.value = selectedParty.value.name;
  } else if (!newVal) {
    searchTerm.value = '';
  }
});

onMounted(() => {
  loadParties();
});
</script>

<style scoped>
.party-selector {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-weight: 500;
  color: #374151;
  font-size: 0.875rem;
}

.required {
  color: #dc2626;
  margin-left: 0.25rem;
}

.selector-container {
  position: relative;
}

.search-selector {
  position: relative;
}

.form-input {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dropdown-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 0.25rem;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  max-height: 300px;
  overflow-y: auto;
  z-index: 1000;
}

.dropdown-item {
  padding: 0.75rem 1rem;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: background-color 0.15s;
}

.dropdown-item:hover,
.dropdown-item.active {
  background-color: #f3f4f6;
}

.dropdown-item.selected {
  background-color: #eff6ff;
}

.dropdown-item.loading,
.dropdown-item.empty {
  cursor: default;
  color: #6b7280;
  justify-content: center;
}

.party-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

.party-name {
  font-weight: 500;
  color: #111827;
}

.party-tax,
.party-detail {
  font-size: 0.75rem;
  color: #6b7280;
}

.party-role {
  font-size: 0.75rem;
  color: #3b82f6;
  font-weight: 500;
}

.selected-indicator {
  color: #10b981;
  font-weight: bold;
}

.selected-party {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.625rem 0.875rem;
  background-color: #eff6ff;
  border: 1px solid #3b82f6;
  border-radius: 0.375rem;
}

.selected-party-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.btn-clear {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem;
  font-size: 1.25rem;
  line-height: 1;
  transition: color 0.2s;
}

.btn-clear:hover {
  color: #dc2626;
}

.help-text {
  font-size: 0.75rem;
  color: #6b7280;
}

.spinner-small {
  display: inline-block;
  width: 1rem;
  height: 1rem;
  border: 2px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  margin-right: 0.5rem;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
