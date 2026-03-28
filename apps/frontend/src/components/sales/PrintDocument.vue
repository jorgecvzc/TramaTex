<script setup lang="ts">
/**
 * PrintDocument.vue - Plantilla Profesional A4 TramaTex
 * 
 * Genera un documento formal (Factura, Albarán, Presupuesto)
 * siguiendo el estándar industrial A4.
 */
import { computed } from 'vue'
import salesApi from '@/services/salesApi'

const props = defineProps<{
  type: 'QUOTE' | 'ORDER' | 'DELIVERY_NOTE' | 'INVOICE'
  number: string
  date: string
  customerName: string
  customerTaxId?: string
  address?: any
  items: any[]
  totals?: any
  notes?: string
}>()

const docTitle = computed(() => {
  const map = {
    QUOTE: 'PRESUPUESTO',
    ORDER: 'PEDIDO DE VENTA',
    DELIVERY_NOTE: 'ALBARÁN DE ENTREGA',
    INVOICE: 'FACTURA'
  }
  return map[props.type] || 'DOCUMENTO'
})

const formatDate = (d: string) => d ? new Date(d).toLocaleDateString('es-ES') : '—'
</script>

<template>
  <div class="print-document-container">
    <!-- CAPA 1: CABECERA CORPORATIVA -->
    <header class="print-header">
      <div class="issuer-info">
        <h1>TRAMATEX</h1>
        <p>Vestuario Laboral y EPIs</p>
        <p>C/ Industria, 42 - Pol. Ind. El Trama</p>
        <p>28001 Madrid - España</p>
        <p>NIF: B-12345678</p>
      </div>
      <div class="document-info">
        <h2>{{ docTitle }}</h2>
        <div class="doc-meta">
          <p><strong>Nº:</strong> {{ number }}</p>
          <p><strong>Fecha:</strong> {{ formatDate(date) }}</p>
        </div>
      </div>
    </header>

    <!-- CAPA 2: DIRECCIONES -->
    <div class="print-addresses">
      <div class="address-block">
        <label>Cliente</label>
        <strong>{{ customerName }}</strong>
        <p v-if="customerTaxId">NIF/CIF: {{ customerTaxId }}</p>
        <div v-if="address" class="mt-2">
          <p>{{ address.street }}</p>
          <p>{{ address.postalCode }} {{ address.city }}</p>
          <p>{{ address.province || '' }} {{ address.country || '' }}</p>
        </div>
      </div>
    </div>

    <!-- CAPA 3: CUERPO DE DATOS -->
    <table class="print-table">
      <thead>
        <tr>
          <th style="width: 15%">Ref.</th>
          <th>Descripción</th>
          <th style="width: 10%; text-align: center">Cant.</th>
          <th style="width: 15%; text-align: right">Precio</th>
          <th style="width: 10%; text-align: center">Dto.</th>
          <th style="width: 15%; text-align: right">Subtotal</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, idx) in items" :key="idx">
          <td><code class="mono">{{ item.variantSku || '—' }}</code></td>
          <td>{{ item.productName || item.displayName }}</td>
          <td style="text-align: center">{{ item.quantity || item.deliveredQuantity }}</td>
          <td style="text-align: right">{{ salesApi.formatMoney(item.unitPrice || item.listUnitPrice) }}</td>
          <td style="text-align: center">{{ item.discountPercent ? item.discountPercent + '%' : '—' }}</td>
          <td style="text-align: right"><strong>{{ salesApi.formatMoney(item.subtotal || (item.unitPrice?.amount * item.quantity)) }}</strong></td>
        </tr>
      </tbody>
    </table>

    <!-- CAPA 4: TOTALES Y NOTAS -->
    <div class="print-footer-content">
      <div class="print-notes" v-if="notes">
        <label>Observaciones:</label>
        <p>{{ notes }}</p>
      </div>
      <div class="print-totals" v-if="totals">
        <div class="print-total-row">
          <span>Base Imponible:</span>
          <span>{{ salesApi.formatMoney(totals.subtotal) }}</span>
        </div>
        <div class="print-total-row">
          <span>IVA (21%):</span>
          <span>{{ salesApi.formatMoney(totals.taxAmount) }}</span>
        </div>
        <div class="print-total-row final">
          <span>TOTAL DOCUMENTO:</span>
          <span>{{ salesApi.formatMoney(totals.total) }}</span>
        </div>
      </div>
    </div>

    <!-- PIE DE PÁGINA REGLAMENTARIO -->
    <footer class="print-legal-footer">
      <p>Inscrita en el Registro Mercantil de Madrid, Tomo 12345, Folio 67, Hoja M-123456.</p>
      <p>Documento generado electrónicamente por TramaTex ERP.</p>
    </footer>
  </div>
</template>

<style scoped>
/* Estos estilos solo se activan al imprimir o en modo preview */
.print-document-container {
  font-family: 'Inter', sans-serif;
  color: black;
  line-height: 1.5;
  background: white;
  padding: 0;
}

.print-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 40px;
  border-bottom: 3px solid #1b3a6b;
  padding-bottom: 20px;
}

.issuer-info h1 { margin: 0; font-size: 24px; color: #1b3a6b; }
.issuer-info p { margin: 2px 0; font-size: 12px; color: #444; }

.document-info { text-align: right; }
.document-info h2 { margin: 0; font-size: 20px; color: #333; }
.doc-meta { margin-top: 10px; font-size: 14px; }

.print-addresses { margin-bottom: 40px; }
.address-block label { 
  display: block; 
  font-size: 10px; 
  text-transform: uppercase; 
  font-weight: bold; 
  color: #666;
  border-bottom: 1px solid #eee;
  margin-bottom: 5px;
}
.address-block strong { font-size: 16px; display: block; margin-bottom: 5px; }
.address-block p { margin: 2px 0; font-size: 13px; }

.print-table { width: 100%; border-collapse: collapse; margin-bottom: 30px; }
.print-table th { 
  background: #f8fafc; 
  padding: 10px; 
  font-size: 11px; 
  text-transform: uppercase; 
  text-align: left;
  border-bottom: 2px solid #333;
}
.print-table td { padding: 10px; font-size: 12px; border-bottom: 1px solid #eee; }
.mono { font-family: 'Fira Code', monospace; background: #f1f5f9; padding: 2px 4px; border-radius: 3px; font-size: 10px; }

.print-footer-content { display: flex; justify-content: space-between; gap: 40px; page-break-inside: avoid; }
.print-notes { flex: 1; }
.print-notes label { font-size: 10px; font-weight: bold; text-transform: uppercase; color: #666; }
.print-notes p { font-size: 12px; font-style: italic; color: #444; margin-top: 5px; white-space: pre-wrap; }

.print-totals { width: 300px; }
.print-total-row { display: flex; justify-content: space-between; padding: 5px 0; font-size: 13px; }
.print-total-row.final { 
  margin-top: 10px; 
  padding-top: 10px; 
  border-top: 2px solid #1b3a6b; 
  font-weight: bold; 
  font-size: 16px; 
  color: #1b3a6b; 
}

.print-legal-footer {
  margin-top: 50px;
  border-top: 1px solid #eee;
  padding-top: 10px;
  text-align: center;
  font-size: 9px;
  color: #999;
}

@media screen {
  .print-document-container {
    max-width: 800px;
    margin: 2rem auto;
    padding: 40px;
    box-shadow: 0 0 20px rgba(0,0,0,0.1);
    border-radius: 8px;
  }
}
</style>
