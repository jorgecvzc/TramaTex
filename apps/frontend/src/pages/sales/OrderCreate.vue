<template>
  <Navbar />
  <div class="order-create-container">
    <div class="page-header">
      <div>
        <button class="btn-back" @click="goBack">← Volver</button>
        <h1>Nuevo Pedido</h1>
      </div>
    </div>

    <div class="form-card">
      <form @submit.prevent="handleSubmit">
        <!-- Customer Selection -->
        <div class="form-section">
          <h2>Cliente</h2>
          <PartySelector
            v-model="formData.partyId"
            label="Cliente"
            placeholder="Buscar cliente por nombre..."
            role-filter="CLIENT"
            :required="true"
            help-text="Seleccione el cliente para este pedido"
          />
        </div>

        <!-- Order Details -->
        <div class="form-section">
          <h2>Detalles del Pedido</h2>
          <div class="form-row">
            <div class="form-group">
              <label for="orderDate">Fecha de Pedido *</label>
              <input
                id="orderDate"
                v-model="formData.orderDate"
                type="date"
                class="form-input"
                required
              />
            </div>
            <div class="form-group">
              <label for="deliveryDate">Fecha de Entrega *</label>
              <input
                id="deliveryDate"
                v-model="formData.deliveryDate"
                type="date"
                class="form-input"
                :min="minDeliveryDate"
                required
              />
            </div>
          </div>
          <div class="form-group">
            <label for="notes">Notas</label>
            <textarea
              id="notes"
              v-model="formData.notes"
              class="form-textarea"
              rows="3"
              placeholder="Notas adicionales sobre el pedido..."
            ></textarea>
          </div>
        </div>

        <!-- Line Items -->
        <div class="form-section">
          <div class="section-header">
            <h2>Líneas del Pedido</h2>
            <button type="button" class="btn btn-secondary" @click="addLineItem">
              + Agregar Línea
            </button>
          </div>

          <div v-if="formData.lineItems.length === 0" class="empty-state">
            <p>No hay líneas agregadas. Agregue al menos una línea para crear el pedido.</p>
          </div>

          <div v-else class="line-items-list">
            <div
              v-for="(item, index) in formData.lineItems"
              :key="index"
              class="line-item-card"
            >
              <div class="line-item-header">
                <span class="line-item-number">Línea {{ index + 1 }}</span>
                <button
                  type="button"
                  class="btn-remove"
                  @click="removeLineItem(index)"
                  title="Eliminar línea"
                >
                  ✕
                </button>
              </div>
              <div class="line-item-fields">
                <div class="form-group">
                  <label>Variante de Producto *</label>
                  <button
                    type="button"
                    class="btn-select-variant"
                    @click="openVariantSelector(index)"
                  >
                    {{ item.selectedVariantName || 'Seleccionar variante...' }}
                  </button>
                  <input
                    v-model="item.productVariantId"
                    type="hidden"
                    required
                  />
                </div>
                <div class="form-group">
                  <label>Cantidad *</label>
                  <input
                    v-model.number="item.quantity"
                    type="number"
                    min="1"
                    class="form-input"
                    required
                  />
                </div>
                <div class="form-group">
                  <label>Precio Manual (EUR)</label>
                  <input
                    v-model.number="item.manualUnitPrice"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="Opcional"
                    class="form-input"
                  />
                </div>
                <div class="form-group">
                  <label>Descuento Manual (EUR)</label>
                  <input
                    v-model.number="item.manualDiscountPerUnit"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="Opcional"
                    class="form-input"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="goBack">
            Cancelar
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="!isFormValid || isSubmitting"
          >
            {{ isSubmitting ? "Creando..." : "Crear Pedido" }}
          </button>
        </div>
      </form>

      <!-- Error Display -->
      <div v-if="submitError" class="error-box">
        {{ submitError }}
      </div>
    </div>

    <!-- Variant Selector Modal -->
    <div v-if="showVariantSelector" class="modal-overlay" @click.self="showVariantSelector = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Seleccionar Variante de Producto</h3>
          <button class="btn-close" @click="showVariantSelector = false">✕</button>
        </div>
        <div class="modal-body">
          <VariantSelector
            :product-id="null"
            title=""
            description="Seleccione una variante de producto"
            @variant-selected="handleVariantSelected"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import salesApi from '@/services/salesApi.js';

const router = useRouter();

const formData = ref({
  partyId: '',
  orderDate: '',
  deliveryDate: '',
  notes: '',
  lineItems: [],
});

const isSubmitting = ref(false);
const submitError = ref('');

const minDeliveryDate = computed(() => {
  return new Date().toISOString().split('T')[0];
});

const isFormValid = computed(() => {
  return (
    formData.value.partyId &&
    formData.value.orderDate &&
    formData.value.deliveryDate &&
    formData.value.lineItems.length > 0 &&
    formData.value.lineItems.every(
      (item) => item.productVariantId && item.quantity > 0
    )
  );
});

onMounted(() => {
  // Set default order date to today
  const today = new Date().toISOString().split('T')[0];
  formData.value.orderDate = today;
  formData.value.deliveryDate = today;
});

const showVariantSelector = ref(false);
const editingLineIndex = ref(null);

function addLineItem() {
  formData.value.lineItems.push({
    productVariantId: '',
    selectedVariantName: '',
    quantity: 1,
    manualUnitPrice: null,
    manualDiscountPerUnit: null,
  });
}

function openVariantSelector(index) {
  editingLineIndex.value = index;
  showVariantSelector.value = true;
}

function handleVariantSelected(variant) {
  if (editingLineIndex.value !== null && variant) {
    const item = formData.value.lineItems[editingLineIndex.value];
    item.productVariantId = variant.id;
    item.selectedVariantName = `${variant.product_name} - ${variant.sku}`;
  }
  showVariantSelector.value = false;
  editingLineIndex.value = null;
}

function removeLineItem(index) {
  formData.value.lineItems.splice(index, 1);
}

async function handleSubmit() {
  if (!isFormValid.value || isSubmitting.value) return;

  isSubmitting.value = true;
  submitError.value = '';

  try {
    // Prepare line items
    const lineItems = formData.value.lineItems.map((item) => {
      const lineItem = {
        productVariantId: item.productVariantId,
        quantity: item.quantity,
      };

      if (item.manualUnitPrice) {
        lineItem.manualUnitPrice = {
          amount: item.manualUnitPrice,
          currency: 'EUR',
        };
      }

      if (item.manualDiscountPerUnit) {
        lineItem.manualDiscountPerUnit = {
          amount: item.manualDiscountPerUnit,
          currency: 'EUR',
        };
      }

      return lineItem;
    });

    // Prepare order data
    const orderData = {
      partyId: formData.value.partyId,
      orderDate: salesApi.formatDateForAPI(new Date(formData.value.orderDate)),
      deliveryDate: salesApi.formatDateForAPI(new Date(formData.value.deliveryDate)),
      lineItems,
    };

    if (formData.value.notes) {
      orderData.notes = formData.value.notes;
    }

    // Create order
    const newOrder = await salesApi.createOrder(orderData);

    // Navigate to order detail
    router.push(`/sales/orders/${newOrder.id}`);
  } catch (err) {
    submitError.value = err?.message || 'No se pudo crear el pedido';
    console.error('Error creating order:', err);
  } finally {
    isSubmitting.value = false;
  }
}

function goBack() {
  router.push('/sales/orders');
}
</script>

<style scoped>
.order-create-container {
  padding: 2rem;
  max-width: 900px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0.5rem 0;
}

.btn-back {
  background: transparent;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  transition: color 0.2s;
}

.btn-back:hover {
  color: #1f2937;
}

.form-card {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid #f3f4f6;
}

.form-section:last-of-type {
  border-bottom: none;
}

.form-section h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1.5rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h2 {
  margin: 0;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.5rem;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: inherit;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.help-text {
  display: block;
  font-size: 0.75rem;
  color: #9ca3af;
  margin-top: 0.25rem;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #9ca3af;
  background: #f9fafb;
  border-radius: 4px;
}

.line-items-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.line-item-card {
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 1rem;
  background: #f9fafb;
}

.line-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.line-item-number {
  font-size: 0.875rem;
  font-weight: 600;
  color: #4b5563;
}

.btn-remove {
  background: transparent;
  border: none;
  color: #dc2626;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn-remove:hover {
  background: rgba(220, 38, 38, 0.1);
}

.line-item-fields {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.btn {
  padding: 0.625rem 1.25rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #E6B800;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #d4a700;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f3f4f6;
  color: #4a5568;
}

.btn-secondary:hover {
  background: #e5e7eb;
}
.btn-select-variant {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background: white;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.875rem;
}

.btn-select-variant:hover {
  border-color: #3b82f6;
  background: #f9fafb;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  color: #1b3a6b;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #6b7280;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: color 0.2s;
}

.btn-close:hover {
  color: #dc2626;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
}

.error-box {
  margin-top: 1rem;
  padding: 1rem;
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 4px;
  color: #991b1b;
  font-size: 0.875rem;
}
</style>
