<template>
  <div class="page-layout">
    <!-- CAPA 1: IDENTIDAD -->
    <BasePageHeader 
      title="Perfil Fiscal de Impresión" 
      :breadcrumbs="[{ label: 'Administración', to: '/admin/users' }, { label: 'Configuración Fiscal' }]"
    >
      <template #icon><Receipt :size="28" /></template>
      <template #actions>
        <button class="btn btn-outline btn-sm" @click="reloadProfile">
          <RefreshCw :size="16" /> Recargar
        </button>
        <button class="btn btn-primary btn-sm ml-2" @click="saveProfile">
          <Save :size="16" /> Guardar Cambios
        </button>
      </template>
    </BasePageHeader>

    <main class="page-content">
      <div v-if="!isAdmin" class="alert-error">
        <AlertTriangle :size="20" />
        Solo los administradores pueden gestionar el perfil fiscal del emisor.
      </div>

      <div v-else class="fiscal-layout">
        <!-- ÁREA DE EDICIÓN -->
        <div class="fiscal-form-area">
          <FormSection title="Datos del Emisor" icon="domain">
            <div class="form-grid">
              <div class="form-group">
                <label>Nombre Comercial</label>
                <input v-model.trim="form.displayName" type="text" placeholder="TramaTex" />
              </div>
              <div class="form-group">
                <label>Identificador Fiscal (NIF/CIF)</label>
                <div class="input-row">
                  <input v-model.trim="form.taxLabel" type="text" placeholder="CIF" class="label-input" />
                  <input v-model.trim="form.taxId" type="text" placeholder="B12345678" />
                </div>
              </div>
              <div class="form-group full-width">
                <label>Dirección Postal</label>
                <input v-model.trim="form.addressLine" type="text" placeholder="C/ Ejemplo 123" />
              </div>
              <div class="form-group">
                <label>Ciudad y Código Postal</label>
                <input v-model.trim="form.cityLine" type="text" placeholder="28001 Madrid" />
              </div>
              <div class="form-group">
                <label>Línea de Contacto</label>
                <input v-model.trim="form.contactLine" type="text" placeholder="email / teléfono" />
              </div>
              <div class="form-group full-width">
                <label>Datos Registro Mercantil</label>
                <input v-model.trim="form.mercantileRegistry" type="text" placeholder="R.M. Madrid, Tomo..." />
              </div>
            </div>
          </FormSection>

          <div v-if="message" class="status-message">
            <CheckCircle :size="18" />
            {{ message }}
          </div>
        </div>

        <!-- ÁREA DE PREVISUALIZACIÓN -->
        <div class="fiscal-preview-area">
          <div class="section-title-mini">
            <Printer :size="16" />
            Previsualización de Documento
          </div>
          <div class="print-mockup">
            <div class="mock-header">
              <div class="mock-brand">{{ form.displayName || 'TramaTex' }}</div>
              <div class="mock-fiscal">
                {{ form.taxLabel || 'CIF' }}: {{ form.taxId }}<br>
                {{ form.addressLine }} {{ form.cityLine }}<br>
                {{ form.contactLine }}
              </div>
              <div class="mock-doc-title">PEDIDO DE VENTA #0001</div>
            </div>
            <div class="mock-body">Contenido del documento...</div>
            <div class="mock-footer">
              <p>{{ form.mercantileRegistry }}</p>
              <p>Generado por {{ form.displayName || 'TramaTex' }} · TramaTex ERP</p>
            </div>
          </div>
          <button class="btn btn-outline btn-sm w-full mt-4" @click="printPreview">
            <ExternalLink :size="16" /> Probar Impresión Real
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { Receipt, RefreshCw, Save, AlertTriangle, CheckCircle, Printer, ExternalLink } from 'lucide-vue-next'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import FormSection from '@/components/shared/FormSection.vue'
import { useAuthStore } from '@/stores/auth'
import { getPrintIssuerProfile, savePrintIssuerProfile, resetPrintIssuerProfile, type PrintIssuerProfile } from '@/services/printIssuerProfile'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
const message = ref('')
const form = reactive<PrintIssuerProfile>(getPrintIssuerProfile())

function copyProfileToForm(p: PrintIssuerProfile) {
  form.displayName = p.displayName; 
  form.taxLabel = p.taxLabel; 
  form.taxId = p.taxId;
  form.addressLine = p.addressLine; 
  form.cityLine = p.cityLine; 
  form.contactLine = p.contactLine;
  form.mercantileRegistry = p.mercantileRegistry;
}
function reloadProfile() { copyProfileToForm(getPrintIssuerProfile()); message.value = 'Datos recargados' }
function saveProfile() { savePrintIssuerProfile({ ...form }); message.value = 'Perfil guardado con éxito' }
function printPreview() { window.print() }
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }

/* Corrección de altura de cabecera: Resetear altura fija del header heredada de layouts de dashboard si existiera */
:deep(.page-header) {
  min-height: auto;
  padding: 1.5rem 1rem;
}

.page-content { max-width: 1300px; margin: 0 auto; padding: 1rem; }

.fiscal-layout { display: grid; grid-template-columns: 1fr 400px; gap: 2rem; align-items: flex-start; }

.fiscal-form-area { background: white; border-radius: 12px; padding: 1.5rem; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }

.input-row { display: grid; grid-template-columns: 80px 1fr; gap: 0.5rem; }
.full-width { grid-column: 1 / -1; }

.status-message { margin-top: 1.5rem; display: flex; align-items: center; gap: 0.5rem; color: #166534; background: #f0fdf4; padding: 0.75rem; border-radius: 8px; font-weight: 600; font-size: 0.9rem; }

.fiscal-preview-area { position: sticky; top: 90px; }
.section-title-mini { display: flex; align-items: center; gap: 0.5rem; font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 1rem; }

.print-mockup { background: white; border: 1px solid #d1d5db; border-radius: 4px; padding: 1.5rem; box-shadow: 0 10px 25px rgba(0,0,0,0.05); font-family: sans-serif; }
.mock-brand { font-weight: 900; font-size: 1rem; border-bottom: 2px solid #000; margin-bottom: 0.5rem; padding-bottom: 0.25rem; }
.mock-fiscal { font-size: 0.7rem; color: #4b5563; line-height: 1.4; margin-bottom: 1.5rem; }
.mock-doc-title { font-weight: 800; font-size: 0.85rem; color: #111827; }
.mock-body { margin: 2rem 0; padding: 1rem 0; border-top: 1px solid #e5e7eb; border-bottom: 1px solid #e5e7eb; color: #9ca3af; font-size: 0.75rem; text-align: center; }
.mock-footer { font-size: 0.6rem; color: #9ca3af; text-align: center; margin-top: 1rem; }

.w-full { width: 100%; }

@media (max-width: 1100px) {
  .fiscal-layout { grid-template-columns: 1fr; }
  .fiscal-preview-area { position: static; margin-top: 2rem; }
}

@media print {
  .page-layout > *:not(.page-content), .fiscal-form-area, .section-title-mini, .fiscal-preview-area > button { display: none !important; }
  .page-content { padding: 0; margin: 0; max-width: none; }
  .fiscal-layout { display: block; }
  .print-mockup { border: none; box-shadow: none; padding: 0; }
}
</style>