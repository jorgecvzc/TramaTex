<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Administración / Impresión</p>
          <h1>Perfil fiscal de impresión</h1>
          <p class="subtitle">Configura los datos del emisor para documentos impresos de Sales.</p>
        </div>
        <button class="btn btn-secondary" @click="reloadProfile">
          Recargar
        </button>
      </header>

      <div v-if="!isAdmin" class="alert-warning">
        Solo el rol admin puede gestionar la configuración de impresión.
      </div>

      <section v-else class="card">
        <div class="card-header">
          <div>
            <h2>Datos del emisor</h2>
            <p>Estos datos se muestran en cabecera/pie de impresión en presupuestos, pedidos, albaranes y facturas.</p>
          </div>
          <span class="badge">Admin</span>
        </div>

        <form class="form-grid" @submit.prevent="saveProfile">
          <div>
            <label>Nombre comercial</label>
            <input v-model.trim="form.displayName" type="text" placeholder="TramaTex" required />
          </div>

          <div>
            <label>Etiqueta fiscal</label>
            <input v-model.trim="form.taxLabel" type="text" placeholder="CIF" />
          </div>

          <div>
            <label>Identificador fiscal</label>
            <input v-model.trim="form.taxId" type="text" placeholder="B12345678" />
          </div>

          <div>
            <label>Dirección</label>
            <input v-model.trim="form.addressLine" type="text" placeholder="C/ Ejemplo 123" />
          </div>

          <div>
            <label>Ciudad / CP</label>
            <input v-model.trim="form.cityLine" type="text" placeholder="28001 Madrid" />
          </div>

          <div>
            <label>Contacto</label>
            <input v-model.trim="form.contactLine" type="text" placeholder="info@tramatex.com · +34 900 000 000" />
          </div>

          <div class="form-actions">
            <button type="submit" class="btn btn-primary">
              Guardar
            </button>
            <button type="button" class="btn btn-secondary" @click="resetProfile">
              Restablecer
            </button>
            <button type="button" class="btn btn-secondary" @click="printPreview">
              Vista previa de impresión
            </button>
          </div>

          <p v-if="message" class="helper-text success">{{ message }}</p>
        </form>
      </section>

      <section v-if="isAdmin" class="card print-preview-card">
        <div class="card-header">
          <div>
            <h2>Vista previa</h2>
            <p>Simulación del encabezado y pie que aparecerá en documentos impresos.</p>
          </div>
        </div>

        <div class="print-preview-sheet">
          <div class="print-doc-header">
            <p class="print-brand">{{ form.displayName || 'TramaTex' }}</p>
            <p v-if="form.taxId" class="print-issuer-line">{{ form.taxLabel || 'CIF' }}: {{ form.taxId }}</p>
            <p v-if="form.addressLine || form.cityLine" class="print-issuer-line">{{ form.addressLine }}<span v-if="form.addressLine && form.cityLine"> · </span>{{ form.cityLine }}</p>
            <p v-if="form.contactLine" class="print-issuer-line">{{ form.contactLine }}</p>
            <h3>Documento de ejemplo SALES-0001</h3>
            <div class="print-doc-meta">
              <span>Cliente: Cliente Demo S.L.</span>
              <span>Fecha: 21 febrero 2026</span>
              <span>Estado: Borrador</span>
            </div>
          </div>

          <div class="print-preview-body">
            Contenido del documento...
          </div>

          <div class="print-doc-footer">
            <span>Documento generado por {{ form.displayName || 'TramaTex' }}</span>
            <span>Referencia: SALES-0001</span>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { useAuthStore } from '@/stores/auth'
import {
  getPrintIssuerProfile,
  savePrintIssuerProfile,
  resetPrintIssuerProfile,
  type PrintIssuerProfile,
} from '@/services/printIssuerProfile'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
const message = ref('')

const form = reactive<PrintIssuerProfile>(getPrintIssuerProfile())

function copyProfileToForm(profile: PrintIssuerProfile) {
  form.displayName = profile.displayName
  form.taxLabel = profile.taxLabel
  form.taxId = profile.taxId
  form.addressLine = profile.addressLine
  form.cityLine = profile.cityLine
  form.contactLine = profile.contactLine
}

function reloadProfile() {
  copyProfileToForm(getPrintIssuerProfile())
  message.value = 'Perfil recargado'
}

function saveProfile() {
  const saved = savePrintIssuerProfile({ ...form })
  copyProfileToForm(saved)
  message.value = 'Perfil guardado correctamente'
}

function resetProfile() {
  const reset = resetPrintIssuerProfile()
  copyProfileToForm(reset)
  message.value = 'Perfil restablecido'
}

function printPreview() {
  window.print()
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

.page-header {
  margin-bottom: 2rem;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.page-header h1 {
  color: #1b3a6b;
  font-size: 2rem;
  margin: 0 0 0.5rem;
}

.breadcrumb,
.subtitle {
  margin: 0;
  color: #64748b;
}

.breadcrumb {
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 1.5rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1rem;
}

.card-header h2 {
  margin: 0 0 0.35rem;
  color: #0f172a;
}

.card-header p {
  margin: 0;
  color: #64748b;
}

.badge {
  align-self: flex-start;
  background: #fef3c7;
  color: #92400e;
  border-radius: 999px;
  padding: 0.25rem 0.65rem;
  font-size: 0.75rem;
  font-weight: 600;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.form-grid > div {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.form-grid label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #334155;
}

.form-grid input {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 0.65rem 0.75rem;
  font-size: 0.9rem;
}

.form-grid input:focus {
  outline: none;
  border-color: #e6b800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.15);
}

.form-actions {
  grid-column: 1 / -1;
  display: flex;
  gap: 0.75rem;
}

.helper-text {
  grid-column: 1 / -1;
  margin: 0;
  font-size: 0.875rem;
}

.helper-text.success {
  color: #166534;
}

.alert-warning {
  background: #fef3c7;
  border: 1px solid #f59e0b;
  color: #78350f;
  border-radius: 8px;
  padding: 0.85rem 1rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.55rem 0.95rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
}

.btn-primary {
  background: #e6b800;
  color: #0f172a;
}

.btn-secondary {
  background: #e2e8f0;
  color: #0f172a;
}

.print-preview-card {
  margin-top: 1.25rem;
}

.print-preview-sheet {
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
  padding: 1rem;
  background: #ffffff;
}

.print-doc-header,
.print-doc-footer {
  border: 1px solid #d1d5db;
  padding: 0.75rem 1rem;
  background: white;
}

.print-doc-footer {
  margin-top: 0.75rem;
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: #4b5563;
}

.print-brand {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #374151;
}

.print-issuer-line {
  margin: 0;
  font-size: 0.75rem;
  color: #4b5563;
}

.print-doc-header h3 {
  margin: 0.35rem 0;
  font-size: 1.1rem;
  font-weight: 700;
  color: #111827;
}

.print-doc-meta {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  font-size: 0.8rem;
  color: #4b5563;
}

.print-preview-body {
  margin-top: 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  min-height: 80px;
  padding: 0.75rem;
  color: #6b7280;
  font-size: 0.875rem;
}

@media print {
  .dashboard :deep(.navbar),
  .page-header,
  .form-grid,
  .card-header,
  .helper-text,
  .alert-warning,
  .print-preview-body {
    display: none !important;
  }

  .dashboard,
  .dashboard-content,
  .card,
  .print-preview-sheet {
    background: white;
    border: none;
    box-shadow: none;
    padding: 0;
    margin: 0;
    max-width: none;
  }

  .print-doc-header,
  .print-doc-footer {
    margin: 0 0 0.75rem;
  }
}

@media (max-width: 768px) {
  .dashboard-content {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
