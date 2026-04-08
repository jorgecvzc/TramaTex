<template>
  <BaseEntityPage v-if="isLoading">
    <template #header>
      <BasePageHeader title="Cargando..." :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Cargando' }]" />
    </template>
    <div class="loading-state card">
      <div class="spinner"></div>
      <p>Consultando ficha de entidad...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="error">
    <template #header>
      <BasePageHeader title="Error" :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Error' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-icon-wrapper error">
        <span class="material-symbols-outlined">error</span>
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/parties')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="party || mode === 'create'">
    <!-- 1. IDENTITY HEADER -->
    <template #header>
      <BasePageHeader 
        :title="mode === 'create' ? 'Nueva Entidad' : (mode === 'edit' ? `Editando ${party?.name}` : party?.name)" 
        :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: mode === 'create' ? 'Crear' : party?.name }]"
        show-back
      >
        <template #icon>
          <span class="material-symbols-outlined">{{ (party?.has_person || formData.hasPerson) ? 'person' : 'domain' }}</span>
        </template>
        <template #actions>
          <template v-if="mode === 'detail'">
            <button class="btn btn-primary btn-sm" @click="enterEditMode">
              <span class="material-symbols-outlined">edit</span> <span>Editar</span>
            </button>
          </template>
          <template v-else>
            <button class="btn btn-outline btn-sm" @click="exitEditMode" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-primary btn-sm" @click="saveParty" :disabled="isSaving">
              <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
              <span>{{ isSaving ? 'Guardando...' : 'Guardar' }}</span>
            </button>
          </template>
        </template>
      </BasePageHeader>
    </template>

    <!-- 2. TOOLBAR -->
    <template #toolbar v-if="mode === 'detail' && party">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', party.status === 'ACTIVE' ? 'status-success' : 'status-secondary']">
            {{ party.status === 'ACTIVE' ? 'Activo' : 'Inactivo' }}
          </span>
          <span class="type-badge-inline">{{ formatRole(party.role) }}</span>
        </div>
        <div class="toolbar-buttons">
          <button class="btn btn-outline btn-sm" @click="toggleStatus">
            <span class="material-symbols-outlined">{{ party.status === 'ACTIVE' ? 'block' : 'check_circle' }}</span>
            <span>{{ party.status === 'ACTIVE' ? 'Desactivar' : 'Activar' }}</span>
          </button>
          <button v-if="party.can_delete" class="btn btn-outline btn-sm btn-danger" @click="deletePartyConfirm">
            <span class="material-symbols-outlined">delete</span>
            <span>Eliminar</span>
          </button>
        </div>
      </div>
    </template>

    <!-- 3. SUMMARY -->
    <template #summary v-if="mode === 'detail' && party">
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><span class="material-symbols-outlined">fingerprint</span></div>
          <div class="tag-content"><label>Identificación</label><strong>{{ party.tax_id || 'Sin NIF/CIF' }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><span class="material-symbols-outlined">percent</span></div>
          <div class="tag-content"><label>Dto. Comercial</label><strong>{{ party.default_discount_percentage || 0 }}%</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><span class="material-symbols-outlined">calendar_today</span></div>
          <div class="tag-content"><label>Fecha Alta</label><strong>{{ formatDate(party.created_at) }}</strong></div>
        </div>
      </div>
    </template>

    <!-- 4. RELATED (Solo para contactos) -->
    <template #related v-if="mode === 'detail' && party?.role === 'CONTACT' && relatedEntities.length > 0">
      <div class="related-history-grid">
        <router-link v-for="entity in relatedEntities" :key="entity.id" :to="`/parties/${entity.id}`" class="related-tag-card highlight-info">
          <div class="tag-icon"><span class="material-symbols-outlined">{{ entity.role === 'SUPPLIER' ? 'factory' : 'person' }}</span></div>
          <div class="tag-content">
            <label>Empresa Vinculada</label>
            <strong>{{ entity.name }}</strong>
          </div>
          <span class="material-symbols-outlined jump-icon">open_in_new</span>
        </router-link>
      </div>
    </template>

    <!-- 5. MAIN CONTENT -->
    <FormSection title="Información Básica" icon="info">
      <div v-if="mode === 'detail'">
        <DataRow label="Nombre / Razón Social" :value="party?.name" icon="person" />
        <DataRow label="Tipo de Entidad" :value="party?.has_person ? 'Persona Física' : 'Organización / Empresa'" icon="category" />
        <DataRow label="Identificación Fiscal" icon="fingerprint">
          <code class="code-badge">{{ party?.tax_id || 'No proporcionado' }}</code>
          <span class="text-xs text-muted ml-2">({{ party?.tax_id_type || 'NIF' }})</span>
        </DataRow>
      </div>
      <div v-else>
        <div class="form-row">
          <div class="form-group">
            <label>Nombre o Razón Social *</label>
            <input v-model="formData.name" type="text" class="form-input" required placeholder="Ej: Textil del Norte S.L." />
          </div>
          <div class="form-group">
            <label>Tipo de Entidad *</label>
            <select v-model="formData.hasPerson" class="form-input" :disabled="mode === 'edit'">
              <option :value="false">Organización / Empresa</option>
              <option :value="true">Persona Física</option>
            </select>
            <span v-if="mode === 'edit'" class="text-xs text-muted mt-1">No modificable tras la creación</span>
          </div>
        </div>
        <div class="form-row mt-4">
          <div class="form-group">
            <label>Identificación Fiscal (NIF/CIF)</label>
            <input v-model="formData.taxId" type="text" class="form-input" placeholder="B12345678" />
          </div>
          <div class="form-group">
            <label>Rol de Negocio *</label>
            <select v-model="formData.role" class="form-input">
              <option value="CLIENT">Cliente</option>
              <option value="SUPPLIER">Proveedor</option>
              <option value="BOTH">Ambos (Cliente/Prov)</option>
              <option value="CONTACT">Solo Contacto</option>
            </select>
          </div>
        </div>
      </div>
    </FormSection>

    <FormSection title="Datos de Contacto" icon="contact_mail">
      <div v-if="mode === 'detail'">
        <DataRow label="Teléfono" :value="party?.phone" icon="call" />
        <DataRow label="Email" icon="mail">
          <a v-if="party?.email" :href="`mailto:${party.email}`" class="link-primary">{{ party.email }}</a>
          <span v-else>—</span>
        </DataRow>
        <DataRow label="Sitio Web" icon="language">
          <a v-if="party?.website" :href="party.website" target="_blank" class="link-primary">{{ party.website }}</a>
          <span v-else>—</span>
        </DataRow>
      </div>
      <div v-else>
        <div class="form-row">
          <div class="form-group">
            <label>Teléfono Principal</label>
            <input v-model="formData.phone" type="tel" class="form-input" />
          </div>
          <div class="form-group">
            <label>Correo Electrónico</label>
            <input v-model="formData.email" type="email" class="form-input" />
          </div>
        </div>
        <div class="form-group mt-4">
          <label>Sitio Web</label>
          <input v-model="formData.website" type="url" class="form-input" placeholder="https://..." />
        </div>
      </div>
    </FormSection>

    <FormSection title="Configuración Comercial" icon="settings_suggest">
      <div v-if="mode === 'detail'">
        <DataRow label="Descuento Comercial" icon="percent">
          <strong class="text-primary" style="font-size: 1.25rem">{{ party?.default_discount_percentage || 0 }}%</strong>
          <span class="text-xs text-muted ml-2">Aplicado por defecto en ventas</span>
        </DataRow>
        <DataRow label="Observaciones Internas" icon="notes">
          <p class="notes-text">{{ party?.notes || 'Sin observaciones.' }}</p>
        </DataRow>
      </div>
      <div v-else>
        <div class="form-group">
          <label>Descuento Comercial (%)</label>
          <input v-model.number="formData.defaultDiscountPercentage" type="number" step="0.01" class="form-input w-32" />
        </div>
        <div class="form-group mt-4">
          <label>Observaciones e Instrucciones Internas</label>
          <textarea v-model="formData.notes" class="form-textarea" rows="3"></textarea>
        </div>
      </div>
    </FormSection>

    <!-- GESTORES DINÁMICOS (Solo en Detalle/Edición) -->
    <div v-if="mode !== 'create' && party?.id" class="mt-8">
      <PersonManager v-if="party.role !== 'CONTACT'" :party-id="party.id" />
      <AddressManager :party-id="party.id" class="mt-8" />
    </div>

    <template #footer v-if="mode === 'detail' && party">
      <div class="audit-info">
        <p>Entidad registrada en el sistema el {{ formatDate(party.created_at) }}.</p>
        <p>UUID único: <code>{{ party.id }}</code></p>
      </div>
    </template>
  </BaseEntityPage>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { partyApi } from '@/services/partyApi';
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import BasePageHeader from '@/components/shared/BasePageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import PersonManager from './PersonManager.vue';
import AddressManager from './AddressManager.vue';

const route = useRoute();
const router = useRouter();

const mode = ref('detail');
const isLoading = ref(true);
const isSaving = ref(false);
const error = ref('');

const party = ref(null);
const relatedEntities = ref([]);

const formData = reactive({
  name: '', role: 'CLIENT', hasPerson: false, taxId: '', taxIdType: 'NIF',
  website: '', phone: '', email: '', notes: '', defaultDiscountPercentage: 0,
});

watch(() => route.params.id, async (newId) => {
  if (newId && newId !== 'new') {
    mode.value = 'detail';
    await fetchParty();
  } else {
    mode.value = 'create';
    resetForm();
    isLoading.value = false;
  }
}, { immediate: true });

onMounted(async () => {
  if (route.params.id !== 'new') await fetchParty();
  else isLoading.value = false;
});

async function fetchParty() {
  const id = route.params.id;
  if (!id || id === 'new') return;
  isLoading.value = true; error.value = '';
  try {
    const data = await partyApi.getParty(id);
    party.value = data;
    if (data.role === 'CONTACT') fetchRelatedEntities(id);
  } catch (err) { error.value = err?.message || 'Error al cargar'; }
  finally { isLoading.value = false; }
}

async function fetchRelatedEntities(id) {
  try {
    const relationships = await partyApi.listRelationships(id);
    const entityIds = (relationships || []).map(rel => rel.to_party_id === id ? rel.from_party_id : rel.to_party_id);
    const entities = await Promise.all(entityIds.map(id => partyApi.getParty(id).catch(() => null)));
    relatedEntities.value = entities.filter(e => e !== null);
  } catch (err) {}
}

function resetForm() {
  Object.assign(formData, {
    name: '', role: 'CLIENT', hasPerson: false, taxId: '', taxIdType: 'NIF',
    website: '', phone: '', email: '', notes: '', defaultDiscountPercentage: 0,
  });
  party.value = null;
}

function enterEditMode() {
  Object.assign(formData, {
    name: party.value.name,
    role: party.value.role,
    hasPerson: party.value.has_person,
    taxId: party.value.tax_id,
    taxIdType: party.value.tax_id_type || 'NIF',
    website: party.value.website,
    phone: party.value.phone,
    email: party.value.email,
    notes: party.value.notes,
    defaultDiscountPercentage: party.value.default_discount_percentage || 0
  });
  mode.value = 'edit';
}

function exitEditMode() {
  if (mode.value === 'edit') mode.value = 'detail';
  else router.push('/parties');
}

async function saveParty() {
  if (!formData.name) { alert('El nombre es obligatorio'); return; }
  isSaving.value = true;
  try {
    const payload = { 
      ...formData, 
      default_discount_percentage: Number(formData.defaultDiscountPercentage),
      has_person: formData.hasPerson,
      tax_id: formData.taxId,
      tax_id_type: formData.taxIdType
    };

    if (mode.value === 'create') {
      const newParty = await partyApi.createParty(payload);
      await router.push(`/parties/${newParty.id}`);
    } else {
      await partyApi.updateParty(party.value.id, payload);
      await fetchParty();
      mode.value = 'detail';
    }
  } catch (err) { alert('Error al guardar: ' + err.message); }
  finally { isSaving.value = false; }
}

async function toggleStatus() {
  const newStatus = party.value.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE';
  try {
    const updated = await partyApi.changePartyStatus(party.value.id, newStatus);
    party.value = updated;
  } catch (err) { alert(err.message); }
}

async function deletePartyConfirm() {
  if (!confirm(`¿Eliminar "${party.value.name}"? Esta acción no se puede deshacer.`)) return;
  try {
    await partyApi.deleteParty(party.value.id);
    router.push('/parties');
  } catch (err) { alert('No se pudo eliminar: ' + (err?.message || 'Error desconocido')); }
}

function formatRole(r) { const map = { CLIENT: 'Cliente', SUPPLIER: 'Proveedor', BOTH: 'Cliente/Prov.', CONTACT: 'Contacto' }; return map[r] || r; }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '—'; }
</script>

<style scoped>
@import "@/design-system/_sections.css";

.overview-tags-row, .related-history-grid { display: flex; flex-wrap: wrap; gap: 1rem; }
.related-history-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); }

.summary-tag { flex: 1; min-width: 240px; padding: 0.6rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 0.75rem; box-shadow: var(--box-shadow-sm); }
.related-tag-card { padding: 0.6rem 1rem; background: var(--color-background); border: 1px solid var(--color-border); border-left: 4px solid var(--color-secondary); border-radius: 10px; display: flex; align-items: center; gap: 0.75rem; text-decoration: none; position: relative; transition: all 0.2s ease; }
.related-tag-card.highlight-info { border-left-color: #2563eb; }
.related-tag-card:hover { background: white; transform: translateX(2px) translateY(-1px); box-shadow: var(--box-shadow-md); }
.related-tag-card:hover strong { color: var(--color-primary); text-decoration: underline; }

.tag-icon { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; background: rgba(0,0,0,0.03); color: var(--color-text-secondary); }
.tag-icon .material-symbols-outlined { font-size: 22px; }

.icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }

.tag-content { display: flex; flex-direction: column; gap: 0.15rem; line-height: 1.2; }
.tag-content label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.025em; }
.tag-content strong { font-size: 0.95rem; color: var(--color-text-primary); }

.action-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; box-shadow: var(--box-shadow-sm); }
.status-badge { padding: 0.4rem 1rem; font-size: 0.85rem; font-weight: 800; }
.toolbar-buttons { display: flex; gap: 0.75rem; }

.type-badge-inline { margin-left: 1rem; padding: 0.4rem 0.75rem; font-size: 0.75rem; font-weight: 800; border-radius: 6px; background: var(--color-background); color: var(--color-text-secondary); border: 1px solid var(--color-border); text-transform: uppercase; }

.notes-text { background: var(--color-background); padding: 1rem; border-radius: 8px; font-style: italic; color: var(--color-text-secondary); line-height: 1.6; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.5rem; border-radius: 4px; font-family: var(--font-family-mono); font-weight: 700; }

.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }

.link-primary { color: var(--color-secondary); font-weight: 600; text-decoration: none; }
.link-primary:hover { text-decoration: underline; }
.w-32 { width: 8rem; }
</style>
