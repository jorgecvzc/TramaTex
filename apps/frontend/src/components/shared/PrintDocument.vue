<script setup lang="ts">
/**
 * PrintDocument.vue
 * Template profesional para la impresión de presupuestos, pedidos y facturas.
 * Diseñado para cumplir estándares industriales y estéticos.
 */
import { computed } from 'vue'

interface PrintData {
  type: 'PRESUPUESTO' | 'PEDIDO' | 'ALBARÁN' | 'FACTURA'
  number: string
  date: string
  expiryDate?: string
  party: {
    name: string
    taxId?: string
    address?: string
    email?: string
    phone?: string
  }
  items: Array<{
    sku?: string
    name: string
    quantity: number
    unitPrice: number
    discount?: number
    subtotal: number
  }>
  subtotal: number
  taxAmount: number
  total: number
  notes?: string
}

const props = defineProps<{
  data: PrintData
}>()

const formatCurrency = (val: number) => {
  return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(val)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('es-ES', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  })
}
</script>

<template>
  <div class="print-container">
    <!-- Header: Company Info & Document Type -->
    <header class="document-header">
      <div class="company-branding">
        <div class="logo-placeholder">TramaTex</div>
        <div class="company-details">
          <p class="company-name">TRAMATEX S.L.</p>
          <p>Polígono Industrial "El Álamo", Nave 14</p>
          <p>28000 Madrid, España</p>
          <p>CIF: B-12345678</p>
          <p>soporte@tramatex.es | +34 912 345 678</p>
        </div>
      </div>
      
      <div class="document-meta">
        <h1 class="doc-type">{{ data.type }}</h1>
        <div class="meta-grid">
          <div class="meta-item">
            <label>Nº Documento:</label>
            <strong>{{ data.number }}</strong>
          </div>
          <div class="meta-item">
            <label>Fecha:</label>
            <span>{{ formatDate(data.date) }}</span>
          </div>
          <div v-if="data.expiryDate" class="meta-item">
            <label>Válido hasta:</label>
            <span>{{ formatDate(data.expiryDate) }}</span>
          </div>
        </div>
      </div>
    </header>

    <div class="address-section">
      <div class="party-address">
        <label>CLIENTE / DESTINATARIO</label>
        <h3>{{ data.party.name }}</h3>
        <p v-if="data.party.taxId">NIF/CIF: {{ data.party.taxId }}</p>
        <p v-if="data.party.address">{{ data.party.address }}</p>
        <p v-if="data.party.email">{{ data.party.email }}</p>
        <p v-if="data.party.phone">{{ data.party.phone }}</p>
      </div>
    </div>

    <!-- Table: Items -->
    <table class="items-table">
      <thead>
        <tr>
          <th class="text-left">DESCRIPCIÓN / REFERENCIA</th>
          <th class="text-center">CANT.</th>
          <th class="text-right">PRECIO UD.</th>
          <th class="text-center">DTO. %</th>
          <th class="text-right">SUBTOTAL</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, idx) in data.items" :key="idx">
          <td>
            <div class="item-desc">
              <span v-if="item.sku" class="sku">[{{ item.sku }}]</span>
              <span class="name">{{ item.name }}</span>
            </div>
          </td>
          <td class="text-center">{{ item.quantity }}</td>
          <td class="text-right">{{ formatCurrency(item.unitPrice) }}</td>
          <td class="text-center">{{ item.discount ? item.discount + '%' : '—' }}</td>
          <td class="text-right"><strong>{{ formatCurrency(item.subtotal) }}</strong></td>
        </tr>
      </tbody>
    </table>

    <!-- Footer: Notes & Totals -->
    <footer class="document-footer">
      <div class="notes-area">
        <template v-if="data.notes">
          <label>OBSERVACIONES</label>
          <p>{{ data.notes }}</p>
        </template>
      </div>

      <div class="totals-area">
        <div class="total-row">
          <label>Base Imponible:</label>
          <span>{{ formatCurrency(data.subtotal) }}</span>
        </div>
        <div class="total-row">
          <label>IVA (21%):</label>
          <span>{{ formatCurrency(data.taxAmount) }}</span>
        </div>
        <div class="total-row final">
          <label>TOTAL DOCUMENTO:</label>
          <span class="total-value">{{ formatCurrency(data.total) }}</span>
        </div>
      </div>
    </footer>

    <div class="legal-disclaimer">
      <p>De acuerdo con la Ley Orgánica de Protección de Datos (LOPD), sus datos están incorporados en un fichero de TRAMATEX S.L.</p>
      <p>Este documento no es válido como factura a menos que se indique explícitamente.</p>
    </div>
  </div>
</template>

<style scoped>
/* Estilos solo para pantalla */
.print-container {
  display: none; /* Oculto por defecto en la UI */
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  color: #1e293b;
  line-height: 1.5;
  background: white;
  padding: 40px;
}

/* @media print activa la visualización */
@media print {
  .print-container {
    display: block !important;
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    min-height: 100%;
    z-index: 9999;
  }

  /* Ocultar elementos de la UI principal durante la impresión */
  :global(#app > *:not(.print-container)),
  :global(.entity-page-container),
  :global(.navbar),
  :global(.side-navbar) {
    display: none !important;
  }
}

.document-header {
  display: flex;
  justify-content: space-between;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 20px;
  margin-bottom: 30px;
}

.logo-placeholder {
  font-size: 28px;
  font-weight: 900;
  color: #1b3a6b;
  margin-bottom: 10px;
  letter-spacing: -1px;
}

.company-name {
  font-weight: 700;
  margin-bottom: 4px;
}

.company-details p {
  margin: 0;
  font-size: 12px;
  color: #64748b;
}

.doc-type {
  font-size: 24px;
  font-weight: 800;
  text-align: right;
  margin: 0 0 15px 0;
  color: #1b3a6b;
}

.meta-grid {
  display: flex;
  flex-direction: column;
  gap: 5px;
  align-items: flex-end;
}

.meta-item {
  font-size: 13px;
}

.meta-item label {
  font-weight: 600;
  margin-right: 8px;
  color: #64748b;
}

.address-section {
  margin-bottom: 40px;
}

.party-address {
  width: 50%;
  border: 1px solid #e2e8f0;
  padding: 20px;
  border-radius: 8px;
}

.party-address label {
  font-size: 10px;
  font-weight: 800;
  color: #94a3b8;
  display: block;
  margin-bottom: 8px;
}

.party-address h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
}

.party-address p {
  margin: 0;
  font-size: 13px;
  color: #475569;
}

/* Items Table */
.items-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 40px;
}

.items-table th {
  background: #f8fafc;
  padding: 12px;
  font-size: 11px;
  font-weight: 800;
  text-transform: uppercase;
  color: #64748b;
  border-bottom: 2px solid #e2e8f0;
}

.items-table td {
  padding: 12px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 13px;
}

.item-desc .sku {
  font-family: monospace;
  font-weight: 700;
  color: #1b3a6b;
  margin-right: 8px;
}

.text-left { text-align: left; }
.text-center { text-align: center; }
.text-right { text-align: right; }

/* Footer */
.document-footer {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 40px;
  margin-bottom: 60px;
}

.notes-area label {
  font-size: 10px;
  font-weight: 800;
  color: #94a3b8;
  display: block;
  margin-bottom: 8px;
}

.notes-area p {
  font-size: 12px;
  color: #475569;
  font-style: italic;
}

.totals-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.total-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.total-row.final {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 2px solid #1b3a6b;
  font-weight: 800;
  font-size: 18px;
}

.total-value {
  color: #1b3a6b;
}

.legal-disclaimer {
  border-top: 1px solid #e2e8f0;
  padding-top: 20px;
  font-size: 10px;
  color: #94a3b8;
  text-align: center;
}

.legal-disclaimer p {
  margin: 4px 0;
}
</style>
