<script setup lang="ts">
/**
 * PrintDocument.vue - Plantilla Profesional A4 TramaTex
 * 
 * Genera un documento formal (Factura, Albarán, Presupuesto)
 * siguiendo el estándar industrial A4 para ser enviado a clientes.
 */
import salesApi from '@/services/salesApi'
import { computed } from 'vue'
import { getPrintIssuerProfile } from '@/services/printIssuerProfile'

const props = defineProps<{
  type: 'QUOTE' | 'ORDER' | 'DELIVERY_NOTE' | 'INVOICE' | 'WORK_ORDER'
  number: string
  date: string
  expiryDate?: string
  customerName: string
  customerTaxId?: string
  address?: any
  items: any[]
  totals?: any
  notes?: string
}>()

const issuer = computed(() => getPrintIssuerProfile())

const docTitle = computed(() => {
  const map = {
    QUOTE: 'PRESUPUESTO',
    ORDER: 'CONFIRMACIÓN DE PEDIDO',
    DELIVERY_NOTE: 'ALBARÁN DE ENTREGA',
    INVOICE: 'FACTURA ORDINARIA',
    WORK_ORDER: 'ORDEN DE TRABAJO (TALLER)'
  }
  return map[props.type] || 'DOCUMENTO'
})

const formatDate = (d: string) => d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'long', day: 'numeric' }) : '—'

const calculateLineSubtotal = (item: any) => {
  const price = item.unitPrice?.amount ?? item.unit_price?.amount ?? item.unitPrice ?? item.unit_price ?? item.listUnitPrice?.amount ?? item.listUnitPrice ?? 0
  const qty = item.quantity ?? item.deliveredQuantity ?? item.delivered_quantity ?? 0
  const discount = item.discountPercent ?? item.discount_percent ?? 0
  return (price * qty) * (1 - (discount / 100))
}

const getGeneralConditions = computed(() => {
  const conditions: Record<string, string> = {
    'QUOTE': 'Este presupuesto tiene una validez de 15 días naturales. Los plazos de entrega indicados son estimativos desde la fecha de aceptación formal y están sujetos a disponibilidad de materias primas y carga de trabajo en el taller de TramaTex.',
    'ORDER': 'Documento de confirmación de pedido. Los precios y condiciones quedan fijados según lo acordado. La fecha de entrega prevista está sujeta a la aprobación definitiva del diseño y el cumplimiento de los hitos de producción. TramaTex se reserva el derecho de ajustar la planificación ante incidencias técnicas mayores.',
    'DELIVERY_NOTE': 'Documento de entrega de mercancía. El receptor declara haber recibido los bultos indicados en perfecto estado, salvo anotación expresa en este documento. La firma del albarán supone la aceptación de la mercancía. No se admitirán reclamaciones transcurridos 48 horas desde la recepción.',
    'INVOICE': 'Factura oficial. El pago deberá realizarse según las condiciones acordadas y en la fecha de vencimiento indicada. En caso de demora, TramaTex se reserva el derecho de aplicar los intereses legales correspondientes según la Ley 3/2004 de lucha contra la morosidad.'
  }
  return conditions[props.type] || conditions['QUOTE']
})
</script>

<template>
  <div class="print-document-container">
    <!-- CABECERA: Membrete Corporativo -->
    <header class="print-header">
      <div class="issuer-logo-area">
        <div class="logo-placeholder">{{ issuer.displayName }}</div>
        <p class="slogan">Vestuario Laboral y Equipamiento de Seguridad</p>
      </div>
      <div class="issuer-details">
        <strong>{{ issuer.displayName }}</strong>
        <p>{{ issuer.addressLine }}</p>
        <p>{{ issuer.cityLine }}</p>
        <p>{{ issuer.taxLabel }}: {{ issuer.taxId }}</p>
        <p>{{ issuer.contactLine }}</p>
      </div>
    </header>

    <div class="document-separator"></div>

    <!-- TÍTULO Y META: Info del documento -->
    <section class="document-main-meta">
      <div class="doc-title-box">
        <h1>{{ docTitle }}</h1>
        <div class="doc-number">Referencia: <strong>{{ number }}</strong></div>
      </div>
      <div class="doc-dates-box">
        <p><span>Fecha Emisión:</span> <strong>{{ formatDate(date) }}</strong></p>
        <p v-if="expiryDate" class="expiry-highlight"><span>Válido hasta:</span> <strong>{{ formatDate(expiryDate) }}</strong></p>
      </div>
    </section>

    <!-- DIRECCIONES: Emisor y Receptor -->
    <section class="address-section">
      <div class="address-block customer">
        <label>DESTINATARIO</label>
        <div class="address-content">
          <strong>{{ customerName }}</strong>
          <p v-if="customerTaxId">NIF/CIF: {{ customerTaxId }}</p>
          <div v-if="address" class="full-address">
            <p>{{ address.street }}</p>
            <p>{{ address.postalCode }} {{ address.city }}</p>
            <p>{{ address.province || '' }} {{ address.country || '' }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- TABLA DE ARTÍCULOS -->
    <table class="print-table">
      <thead>
        <tr>
          <th style="width: 12%">REFERENCIA</th>
          <th>DESCRIPCIÓN</th>
          <th style="width: 8%; text-align: center">CANT.</th>
          <th style="width: 12%; text-align: right">PRECIO</th>
          <th style="width: 8%; text-align: center">DTO.</th>
          <th style="width: 15%; text-align: right">TOTAL LÍNEA</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, idx) in items" :key="idx">
          <td class="mono">{{ item.variantSku || item.variant_sku || '—' }}</td>
          <td>
            <div class="product-name">{{ item.productName || item.product_name || item.displayName || 'Producto' }}</div>
            <div v-if="item.description" class="product-desc">{{ item.description }}</div>
          </td>
          <td style="text-align: center">{{ item.quantity || item.deliveredQuantity || item.delivered_quantity }}</td>
          <td style="text-align: right">{{ salesApi.formatMoney(item.unitPrice || item.unit_price || item.listUnitPrice || item.list_unit_price) }}</td>
          <td style="text-align: center">{{ (item.discountPercent || item.discount_percent) ? (item.discountPercent || item.discount_percent) + '%' : '—' }}</td>
          <td style="text-align: right"><strong>{{ salesApi.formatMoney(item.total || item.subtotal || calculateLineSubtotal(item)) }}</strong></td>
        </tr>
      </tbody>
    </table>

    <!-- TOTALES Y CONDICIONES -->
    <div class="footer-summary">
      <div class="notes-conditions">
        <div v-if="notes" class="notes-block">
          <label>OBSERVACIONES Y CONDICIONES:</label>
          <p>{{ notes }}</p>
        </div>
        <div class="legal-notice">
          <p><strong>Condiciones Generales:</strong> {{ getGeneralConditions }}</p>
          <p>De conformidad con la LOPD, le informamos que sus datos están incorporados en un fichero responsabilidad de TRAMATEX S.L. para la gestión comercial.</p>
        </div>
      </div>
      
      <div v-if="totals && type !== 'DELIVERY_NOTE'" class="totals-block">
        <div class="total-row">
          <span>Suma Importes:</span>
          <span>{{ salesApi.formatMoney(totals.subtotal) }}</span>
        </div>
        <div class="total-row">
          <span>I.V.A. (21%):</span>
          <span>{{ salesApi.formatMoney(totals.taxAmount || totals.tax_amount) }}</span>
        </div>
        <div class="total-row final-amount">
          <span>TOTAL DOCUMENTO:</span>
          <span class="price-big">{{ salesApi.formatMoney(totals.total) }}</span>
        </div>
      </div>
    </div>

    <!-- ESPACIADOR PARA EL FOOTER FIJO (Solo en print) -->
    <div class="print-footer-spacer"></div>

    <!-- PIE DE PÁGINA (Fijo al final del papel en impresión) -->
    <footer class="print-footer">
      <div class="separator-line"></div>
      <div class="footer-info-grid">
        <div class="footer-col">
          <strong>{{ issuer.displayName }}</strong>
          <p>{{ issuer.taxLabel }}: {{ issuer.taxId }}</p>
        </div>
        <div class="footer-col">
          <strong>Registro Mercantil</strong>
          <p>{{ issuer.mercantileRegistry }}</p>
        </div>
        <div class="footer-col">
          <strong>Contacto</strong>
          <p>{{ issuer.contactLine }}</p>
        </div>
      </div>
      <p class="page-number">Página 1 de 1</p>
    </footer>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;800&family=Fira+Code&display=swap');

.print-document-container {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  color: #1a1a1a;
  line-height: 1.4;
  background: white;
  width: 100%;
  display: flex;
  flex-direction: column;
}

/* Cabecera Membrete */
.print-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.logo-placeholder {
  font-size: 32px;
  font-weight: 800;
  color: #1b3a6b;
  letter-spacing: -1.5px;
  line-height: 1;
}

.slogan {
  font-size: 10px;
  color: #666;
  text-transform: uppercase;
  margin-top: 5px;
  letter-spacing: 1px;
}

.issuer-details {
  text-align: right;
  font-size: 11px;
  color: #444;
}

.issuer-details strong {
  display: block;
  font-size: 13px;
  color: #1b3a6b;
  margin-bottom: 4px;
}

.issuer-details p { margin: 1px 0; }

.document-separator {
  height: 2px;
  background: #1b3a6b;
  margin-bottom: 30px;
}

/* Meta del documento */
.document-main-meta {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 40px;
}

.doc-title-box h1 {
  font-size: 28px;
  font-weight: 300;
  color: #1b3a6b;
  margin: 0;
  text-transform: uppercase;
}

.doc-number {
  font-size: 14px;
  color: #666;
  margin-top: 5px;
}

.doc-dates-box {
  text-align: right;
  font-size: 13px;
}

.doc-dates-box p { margin: 3px 0; }
.doc-dates-box span { color: #666; width: 120px; display: inline-block; }
.expiry-highlight { color: #dc2626; }

/* Direcciones */
.address-section {
  margin-bottom: 50px;
}

.address-block {
  max-width: 350px;
}

.address-block label {
  display: block;
  font-size: 10px;
  font-weight: 800;
  color: #1b3a6b;
  border-bottom: 1px solid #eee;
  padding-bottom: 4px;
  margin-bottom: 10px;
  letter-spacing: 1px;
}

.address-content strong {
  font-size: 16px;
  display: block;
  margin-bottom: 5px;
}

.address-content p {
  margin: 2px 0;
  font-size: 13px;
  color: #333;
}

/* Tabla */
.print-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 40px;
}

.print-table th {
  background: #f8fafc;
  color: #1b3a6b;
  font-size: 10px;
  font-weight: 800;
  padding: 12px 10px;
  text-align: left;
  border-top: 1px solid #1b3a6b;
  border-bottom: 2px solid #1b3a6b;
}

.print-table td {
  padding: 12px 10px;
  font-size: 12px;
  border-bottom: 1px solid #eee;
  vertical-align: top;
}

.mono { font-family: 'Fira Code', monospace; font-size: 11px; }
.product-name { font-weight: 600; color: #000; }
.product-desc { font-size: 11px; color: #666; margin-top: 2px; }

/* Resumen y Totales */
.footer-summary {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 50px;
  page-break-inside: avoid;
}

.notes-block label {
  font-size: 10px;
  font-weight: 800;
  color: #666;
  display: block;
  margin-bottom: 8px;
}

.notes-block p {
  font-size: 12px;
  font-style: italic;
  white-space: pre-wrap;
  color: #444;
  background: #fdfdfd;
  padding: 10px;
  border: 1px solid #eee;
}

.legal-notice {
  margin-top: 20px;
  font-size: 9px;
  color: #999;
}

.total-row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 13px;
}

.final-amount {
  margin-top: 15px;
  padding: 15px 0;
  border-top: 2px solid #1b3a6b;
  color: #1b3a6b;
  font-weight: 800;
}

.price-big { font-size: 22px; }

/* Pie Legal */
.print-footer {
  margin-top: 60px;
  text-align: center;
  font-size: 9px;
  color: #999;
}

.footer-info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 10px;
  text-align: center; /* Cambiado de left a center para equilibrar visualmente */
}

.footer-col strong {
  display: block;
  color: #333; /* Un poco más oscuro para mejor legibilidad */
  margin-bottom: 2px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.separator-line {
  height: 1px;
  background: #1b3a6b; /* Usar el azul corporativo para la línea también */
  margin-bottom: 15px;
  opacity: 0.3;
}

.page-number { margin-top: 5px; }

.print-footer-spacer {
  display: none;
  height: 100px;
}

/* Preview en pantalla */
@media screen {
  .print-document-container {
    max-width: 210mm; /* Ancho A4 */
    margin: 3rem auto;
    padding: 20mm; /* Margen A4 */
    box-shadow: 0 20px 50px rgba(0,0,0,0.15);
    border-radius: 4px;
    border: 1px solid #ddd;
    position: relative;
    min-height: 297mm;
  }
  .print-document-container::before {
    content: "VISTA PREVIA DE IMPRESIÓN";
    position: absolute;
    top: 10px;
    left: 50%;
    transform: translateX(-50%);
    font-size: 10px;
    color: #999;
    letter-spacing: 2px;
  }
  .footer-summary {
    margin-bottom: auto; /* Empuja el footer al final en pantalla */
  }
}

@media print {
  .print-document-container {
    padding: 0;
    margin: 0;
    width: 100%;
    min-height: auto;
  }
  .print-footer {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    background: white;
    padding: 15mm 20mm; /* Aumentado para centrar más el contenido visualmente */
    box-sizing: border-box;
    border-top: 1px solid #eee;
  }
  .print-footer-spacer {
    display: block; /* Ocupa el espacio que tapa el footer fixed */
  }
  .footer-summary {
    margin-bottom: 0;
  }
}
</style>
