<template>
  <div :class="['party-selector', { 'dark-mode': darkMode }]">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    
    <div class="selector-container">
      <!-- Search Mode: Dropdown with search -->
      <div class="search-selector">
        <input
          :id="inputId"
          ref="searchInput"
          v-model="searchTerm"
          type="text"
          :placeholder="placeholder || 'Buscar por nombre o referencia...'"
          class="form-input"
          @input="handleInput"
          @focus="onFocus"
          @blur="handleBlur"
          @keydown.enter.prevent="selectFirst"
          @keydown.down.prevent="navigateDown"
          @keydown.up.prevent="navigateUp"
          :required="required"
          autocomplete="off"
        />
        
        <!-- Dropdown Results -->
        <div v-if="showDropdown && searchTerm" class="dropdown-results">
          <div v-if="isSearching" class="dropdown-item loading">
            <span class="spinner-small"></span>
            Buscando...
          </div>
          <template v-else-if="filteredParties.length > 0">
            <div
              v-for="(party, index) in filteredParties"
              :key="party.id"
              :class="['dropdown-item', { active: index === activeIndex, selected: party.id === modelValue }]"
              @mousedown.prevent="selectParty(party)"
              @mouseenter="activeIndex = index"
            >
              <div class="party-info">
                <span class="party-name">{{ party.name }}</span>
                <div class="party-meta">
                  <span v-if="party.tax_id" class="party-tax">{{ party.tax_id }}</span>
                  <span class="party-role-badge">{{ getRoleLabel(party.role) }}</span>
                </div>
              </div>
              <span v-if="party.id === modelValue" class="selected-indicator">✓</span>
            </div>
          </template>
          <div v-else class="dropdown-item empty">
            No se encontraron resultados para "{{ searchTerm }}"
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
    
    <input type="hidden" :value="modelValue" :name="name" />
    <span v-if="helpText" class="help-text">{{ helpText }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';
import { partyApi } from '@/services/partyApi';

const props = defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: '' },
  placeholder: { type: String, default: 'Buscar por nombre o referencia...' },
  required: { type: Boolean, default: false },
  roleFilter: { type: String, default: null },
  name: { type: String, default: 'partyId' },
  helpText: { type: String, default: '' },
  darkMode: { type: Boolean, default: false },
});

const emit = defineEmits(['update:modelValue', 'select']);

const searchTerm = ref('');
const allParties = ref<any[]>([]);
const externalParty = ref<any | null>(null);
const showDropdown = ref(false);
const isSearching = ref(false);
const activeIndex = ref(0);
const inputId = computed(() => `party-selector-${Math.random().toString(36).substr(2, 9)}`);

let searchTimer: any = null;

const selectedParty = computed(() => {
  if (!props.modelValue) return null;
  return allParties.value.find(p => p.id === props.modelValue) || externalParty.value || null;
});

const filteredParties = computed(() => {
  return allParties.value.slice(0, 50);
});

async function loadParties(name = '') {
  isSearching.value = true;
  try {
    const filters: any = { limit: 100 };
    if (props.roleFilter) filters.role = props.roleFilter;
    if (name) filters.name = name;
    
    const response = await partyApi.listParties(filters);
    allParties.value = response.data || (Array.isArray(response) ? response : []);
    
    if (props.modelValue && !name && !selectedParty.value) {
      // Si tenemos un ID pero no lo encontramos en la lista, traerlo específicamente
      try {
        externalParty.value = await partyApi.getParty(props.modelValue);
        if (externalParty.value) searchTerm.value = externalParty.value.name;
      } catch (e) { console.error("Error trayendo entidad seleccionada", e); }
    } else if (selectedParty.value) {
      searchTerm.value = selectedParty.value.name;
    }
  } catch (error) {
    console.error('Error loading parties:', error);
  } finally {
    isSearching.value = false;
  }
}

const searchInput = ref<HTMLInputElement | null>(null);

function handleInput() {
  showDropdown.value = true;
  activeIndex.value = 0;
  
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    loadParties(searchTerm.value);
  }, 300);
}

function onFocus() {
  showDropdown.value = true;
  if (allParties.value.length === 0) loadParties();
}

function handleBlur() {
  setTimeout(() => {
    showDropdown.value = false;
    if (selectedParty.value) searchTerm.value = selectedParty.value.name;
    else if (!props.modelValue) searchTerm.value = '';
  }, 250);
}

function selectParty(party: any) {
  emit('update:modelValue', party.id);
  emit('select', party);
  searchTerm.value = party.name;
  showDropdown.value = false;
}

function selectFirst() {
  if (filteredParties.value.length > 0) selectParty(filteredParties.value[activeIndex.value]);
}

function navigateDown() { if (activeIndex.value < filteredParties.value.length - 1) activeIndex.value++; }
function navigateUp() { if (activeIndex.value > 0) activeIndex.value--; }

function handleGlobalEsc() {
  showDropdown.value = false;
}

onMounted(() => {
  loadParties();
  window.addEventListener('tramatex-esc', handleGlobalEsc);
  
  // Auto-focus if requested or standard
  nextTick(() => {
    searchInput.value?.focus();
  });
});

onBeforeUnmount(() => {
  window.removeEventListener('tramatex-esc', handleGlobalEsc);
});
</script>

<style scoped>
.party-selector { display: flex; flex-direction: column; gap: 0.5rem; }
.form-label { font-weight: 700; color: var(--color-text-secondary); font-size: 0.75rem; text-transform: uppercase; }
.required { color: var(--color-danger); margin-left: 0.25rem; }
.selector-container { position: relative; }
.search-selector { position: relative; }

.form-input { 
  width: 100%; padding: 0.75rem 1rem; border: 1px solid var(--color-border-strong); border-radius: 8px; font-size: 0.9rem; transition: all 0.2s; 
  background: white; box-shadow: var(--box-shadow-sm);
}
.form-input:focus { outline: none; border-color: var(--color-primary); box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.1); }

/* Dark Mode Styles */
.party-selector.dark-mode .form-input { background: #0f172a; border-color: #334155; color: white; }
.party-selector.dark-mode .dropdown-results { background: #1e293b; border-color: #334155; }
.party-selector.dark-mode .dropdown-item { border-bottom-color: #334155; color: #e2e8f0; }
.party-selector.dark-mode .dropdown-item:hover, .party-selector.dark-mode .dropdown-item.active { background-color: #0f172a; }
.party-selector.dark-mode .party-name { color: white; }
.party-selector.dark-mode .selected-party { background-color: #0f172a; border-color: var(--color-primary); }

.dropdown-results {
  position: absolute; top: 100%; left: 0; right: 0; margin-top: 0.5rem; background: white; border: 1px solid var(--color-border-strong);
  border-radius: 8px; box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1); max-height: 350px; overflow-y: auto; z-index: 2000;
}

.dropdown-item { padding: 0.75rem 1rem; cursor: pointer; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--color-background); }
.dropdown-item:last-child { border-bottom: none; }
.dropdown-item:hover, .dropdown-item.active { background-color: var(--color-background); }
.dropdown-item.selected { background-color: rgba(0, 35, 149, 0.05); }

.party-info { display: flex; flex-direction: column; gap: 0.15rem; flex: 1; }
.party-name { font-weight: 600; color: var(--color-text-primary); }
.party-meta { display: flex; align-items: center; gap: 0.75rem; }
.party-tax { font-size: 0.75rem; color: var(--color-text-secondary); font-family: var(--font-family-mono); }
.party-role-badge { font-size: 0.65rem; padding: 0.1rem 0.4rem; background: rgba(0, 35, 149, 0.1); color: var(--color-primary); border-radius: 4px; font-weight: 700; text-transform: uppercase; }

.selected-indicator { color: var(--color-success); font-weight: bold; }

.selected-party {
  display: flex; align-items: center; justify-content: space-between; padding: 0.75rem 1rem; background-color: white; 
  border: 2px solid var(--color-primary); border-radius: 8px; box-shadow: var(--box-shadow-sm);
}
.selected-party .party-name { color: var(--color-primary); }

.btn-clear { background: none; border: none; color: var(--color-text-secondary); cursor: pointer; padding: 0.25rem; font-size: 1.1rem; }
.btn-clear:hover { color: var(--color-danger); }

.spinner-small { display: inline-block; width: 1rem; height: 1rem; border: 2px solid rgba(0,0,0,0.1); border-top-color: var(--color-primary); border-radius: 50%; animation: spin 0.6s linear infinite; margin-right: 0.5rem; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
