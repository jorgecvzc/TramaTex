<script setup lang="ts">
/**
 * PrintTicket.vue - Plantilla de Ticket Térmico (80mm)
 * 
 * Específico para el terminal de venta directa (Tickets).
 * Formato compacto, fuentes optimizadas para impresoras térmicas.
 */
import salesApi from '@/services/salesApi'
import { computed } from 'vue'
import { getPrintIssuerProfile } from '@/services/printIssuerProfile'

const props = defineProps<{
  number: string
  date: string
  items: any[]
  totals: any
  customerName?: string
}>()

const issuer = computed(() => getPrintIssuerProfile())

const formatDate = (d: string) => d ? new Date(d).toLocaleString('es-ES') : '—'
</script>

<template>
  <div class="ticket-container">
    <header class="ticket-header">
      <h1>{{ issuer.displayName }}</h1>
      <p>{{ issuer.taxLabel }}: {{ issuer.taxId }}</p>
      <p>{{ issuer.addressLine }}</p>
      <p>{{ issuer.cityLine }}</p>
      <p>{{ issuer.contactLine }}</p>
    </header>

    <div class="ticket-info">
      <p><strong>TICKET:</strong> {{ number }}</p>
      <p><strong>FECHA:</strong> {{ formatDate(date) }}</p>
      <p v-if="customerName"><strong>CTE:</strong> {{ customerName }}</p>
    </div>

    <div class="divider">********************************</div>

    <table class="ticket-table">
      <thead>
        <tr>
          <th>DESCRIPCIÓN</th>
          <th style="text-align: center">CANT</th>
          <th style="text-align: right">IMPORTE</th>
        </tr>
      </thead>
      <tbody>
        <template v-for="(item, idx) in items" :key="idx">
          <tr>
            <td colspan="3" style="padding-bottom: 0"><strong>{{ item.productName }}</strong></td>
          </tr>
          <tr>
            <td style="font-size: 9px; color: #444; padding-top: 0">
              {{ salesApi.formatMoney({ amount: item.unitPrice, currency: 'EUR' }) }}
              <span v-if="item.discountPercent > 0"> (-{{ item.discountPercent }}%)</span>
            </td>
            <td style="text-align: center; padding-top: 0">{{ item.quantity }}</td>
            <td style="text-align: right; padding-top: 0">{{ salesApi.formatMoney({ amount: item.subtotal, currency: 'EUR' }) }}</td>
          </tr>
        </template>
      </tbody>
    </table>

    <div class="divider">********************************</div>

    <div class="ticket-totals">
      <div class="total-row">
        <span>SUBTOTAL:</span>
        <span>{{ salesApi.formatMoney({ amount: totals.subtotal, currency: 'EUR' }) }}</span>
      </div>
      <div class="total-row">
        <span>IVA (21%):</span>
        <span>{{ salesApi.formatMoney({ amount: totals.taxAmount, currency: 'EUR' }) }}</span>
      </div>
      <div class="total-row grand-total">
        <span>TOTAL:</span>
        <span>{{ salesApi.formatMoney({ amount: totals.total, currency: 'EUR' }) }}</span>
      </div>
    </div>

    <footer class="ticket-footer">
      <p>¡GRACIAS POR SU VISITA!</p>
      <p>No se admiten devoluciones sin ticket.</p>
      <p>www.tramatex.local</p>
    </footer>
  </div>
</template>

<style scoped>
.ticket-container {
  width: 80mm; /* Ancho estándar de impresora térmica */
  margin: 0 auto;
  background: white;
  color: black;
  font-family: 'Courier New', Courier, monospace; /* Fuente mono para tickets */
  font-size: 12px;
  padding: 5mm;
}

.ticket-header { text-align: center; margin-bottom: 5mm; }
.ticket-header h1 { font-size: 18px; margin: 0 0 2mm 0; }
.ticket-header p { margin: 1mm 0; }

.ticket-info { margin-bottom: 3mm; font-size: 11px; }
.ticket-info p { margin: 1mm 0; }

.divider { text-align: center; margin: 2mm 0; }

.ticket-table { width: 100%; border-collapse: collapse; font-size: 11px; }
.ticket-table th { border-bottom: 1px dashed #000; padding: 2mm 0; text-align: left; }
.ticket-table td { padding: 2mm 0; vertical-align: top; }

.ticket-totals { margin-top: 3mm; }
.total-row { display: flex; justify-content: space-between; margin-bottom: 1mm; }
.grand-total { font-weight: bold; font-size: 14px; margin-top: 2mm; border-top: 1px solid #000; padding-top: 2mm; }

.ticket-footer { text-align: center; margin-top: 10mm; font-size: 10px; }
.ticket-footer p { margin: 1mm 0; }

@media print {
  /* Ocultar absolutamente todo excepto el ticket */
  :global(body) { background: white !important; padding: 0 !important; margin: 0 !important; }
  :global(#app), :global(nav), :global(header), :global(footer) { display: none !important; }
  
  .ticket-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 80mm;
    margin: 0;
    padding: 2mm;
    z-index: 999999;
    visibility: visible !important;
    display: block !important;
    background: white !important;
    color: black !important;
  }
}
</style>
