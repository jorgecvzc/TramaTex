<template>
  <div class="product-list">
    <!-- Filters and Search -->
    <div class="filters">
      <div>
        <label>Buscar por nombre o SKU</label>
        <input
          v-model="filters.search"
          type="text"
          placeholder="Buscar producto..."
          @input="debouncedSearch"
        />
      </div>

      <div>
        <label>Filtrar por marca</label>
        <select v-model="filters.brandId" @change="applyFilters">
          <option value="">Todas las marcas</option>
          <option v-for="brand in brands" :key="brand.id" :value="brand.id">
            {{ brand.name }}
          </option>
        </select>
      </div>

      <div>
        <label>Filtrar por categoría</label>
        <select v-model="filters.groupId" @change="applyFilters">
          <option value="">Todas las categorías</option>
          <option v-for="group in productGroups" :key="group.id" :value="group.id">
            {{ group.name }}
          </option>
        </select>
      </div>

      <div>
        <label>Filtrar por estado</label>
        <select v-model="filters.isActive" @change="applyFilters">
          <option value="">Todos</option>
          <option value="true">Activo</option>
          <option value="false">Inactivo</option>
        </select>
      </div>

      <div>
        <label>Filtrar por tipo</label>
        <select v-model="filters.productType" @change="applyFilters">
          <option value="">Todos</option>
          <option value="TANGIBLE">Tangible</option>
          <option value="SERVICE">Servicio</option>
        </select>
      </div>

      <div class="filter-actions">
        <button @click="clearFilters" class="btn btn-secondary">
          Limpiar filtros
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Cargando productos...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="alert-error">
      {{ error }}
      <button @click="fetchProducts" class="btn btn-outline retry-btn">
        Reintentar
      </button>
    </div>

    <!-- Products Table -->
    <div v-if="!isLoading && products.length > 0" class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>SKU</th>
            <th>Nombre</th>
            <th>Marca</th>
            <th>Categoría</th>
            <th>Tipo</th>
            <th class="align-right">Precio Base</th>
            <th>Variantes</th>
            <th>Estado</th>
            <th class="align-right">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="product in products"
            :key="product.id"
            :class="{ inactive: !product.is_active }"
          >
            <td>
              <code class="sku-code">{{ product.sku || '—' }}</code>
            </td>
            <td>
              <router-link :to="`/products/${product.id}`" class="product-link">
                {{ product.name }}
              </router-link>
              <div v-if="product.long_name" class="product-long-name">
                {{ product.long_name }}
              </div>
            </td>
            <td>
              <span class="brand-name">{{ getBrandName(product.brand_id) }}</span>
            </td>
            <td>
              <span class="group-name">{{ getGroupName(product.group_ids?.[0]) }}</span>
            </td>
            <td>
              <span class="type-pill" :class="`type-${product.product_type?.toLowerCase()}`">
                {{ formatProductType(product.product_type) }}
              </span>
            </td>
            <td class="align-right">
              <span class="price-value">
                {{ formatPrice(product.base_price) }}
              </span>
            </td>
            <td>
              <span class="variants-badge">
                {{ product.variants_count || 0 }} var{{ product.variants_count !== 1 ? 's' : '' }}
              </span>
            </td>
            <td>
              <span class="status-pill" :class="`status-${product.is_active ? 'active' : 'inactive'}`">
                {{ product.is_active ? 'Activo' : 'Inactivo' }}
              </span>
            </td>
            <td class="align-right">
              <div class="action-buttons">
                <router-link :to="`/products/${product.id}`" class="btn btn-outline">
                  Ver detalles
                </router-link>
                <button class="btn btn-secondary" @click="toggleStatus(product)">
                  {{ product.is_active ? 'Desactivar' : 'Activar' }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="!isLoading && products.length === 0" class="empty-state-block">
      <Package :size="64" class="empty-icon" />
      <p v-if="hasFilters">
        No se encontraron productos con los filtros aplicados.
      </p>
      <p v-else>
        Aún no hay productos en el catálogo.
      </p>
      <p v-if="hasFilters" class="empty-hint">
        Prueba ajustando los filtros para ver más resultados.
      </p>
      <p v-else class="empty-hint">
        Crea tu primer producto para empezar a gestionar tu catálogo.
      </p>
      <button v-if="hasFilters" @click="clearFilters" class="btn btn-primary">
        Limpiar filtros
      </button>
      <router-link v-else to="/products/new" class="btn btn-primary">
        + Crear primer producto
      </router-link>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        :disabled="currentPage === 1"
        @click="previousPage"
        class="btn btn-secondary"
      >
        ← Anterior
      </button>
      <span class="page-info">Página {{ currentPage }} de {{ totalPages }}</span>
      <button
        :disabled="currentPage === totalPages"
        @click="nextPage"
        class="btn btn-secondary"
      >
        Siguiente →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import { productApi } from '@/services/productApi';
import { Package } from 'lucide-vue-next';

const products = ref([]);
const brands = ref([]);
const productGroups = ref([]);
const isLoading = ref(false);
const error = ref('');
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
let searchTimeout = null;

const filters = reactive({
  search: '',
  brandId: '',
  groupId: '',
  isActive: '',
  productType: '',
});

const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

const hasFilters = computed(
  () =>
    filters.search.trim() !== '' ||
    filters.brandId !== '' ||
    filters.groupId !== '' ||
    filters.isActive !== '' ||
    filters.productType !== ''
);

onMounted(() => {
  fetchBrands();
  fetchProductGroups();
  fetchProducts();
});

onUnmounted(() => {
  if (searchTimeout) clearTimeout(searchTimeout);
});

async function fetchProducts() {
  isLoading.value = true;
  error.value = '';

  try {
    const response = await productApi.listProducts({
      search: filters.search,
      brandId: filters.brandId,
      groupId: filters.groupId,
      isActive: filters.isActive,
      productType: filters.productType,
      pageNumber: currentPage.value,
      pageSize: pageSize.value,
    });

    products.value = response.data || [];
    total.value = response.total || 0;
  } catch (err) {
    error.value = err?.message || 'No se pudieron cargar los productos';
    products.value = [];
  } finally {
    isLoading.value = false;
  }
}

async function fetchBrands() {
  try {
    const response = await productApi.listBrands({ isActive: true });
    brands.value = response.data || [];
  } catch (err) {
    console.error('Error loading brands:', err);
  }
}

async function fetchProductGroups() {
  try {
    const response = await productApi.listProductGroups({ isActive: true });
    productGroups.value = response.data || [];
  } catch (err) {
    console.error('Error loading product groups:', err);
  }
}

function applyFilters() {
  currentPage.value = 1;
  fetchProducts();
}

function debouncedSearch() {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    applyFilters();
  }, 350);
}

function clearFilters() {
  filters.search = '';
  filters.brandId = '';
  filters.groupId = '';
  filters.isActive = '';
  filters.productType = '';
  currentPage.value = 1;
  fetchProducts();
}

function previousPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchProducts();
    scrollToTop();
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchProducts();
    scrollToTop();
  }
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' });
}

async function toggleStatus(product) {
  const newStatus = !product.is_active;
  const action = newStatus ? 'activar' : 'desactivar';

  if (!confirm(`¿Estás seguro de ${action} el producto "${product.name}"?`)) {
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    await productApi.changeProductStatus(product.id, newStatus);
    product.is_active = newStatus;
  } catch (err) {
    error.value = err?.message || `No se pudo ${action} el producto`;
  } finally {
    isLoading.value = false;
  }
}

function getBrandName(brandId) {
  if (!brandId) return '—';
  const brand = brands.value.find((b) => b.id === brandId);
  return brand?.name || '—';
}

function getGroupName(groupId) {
  if (!groupId) return '—';
  const group = productGroups.value.find((g) => g.id === groupId);
  return group?.name || '—';
}

function formatProductType(type) {
  const map = {
    TANGIBLE: 'Tangible',
    SERVICE: 'Servicio',
  };
  return map[type] || type || '—';
}

function formatPrice(value) {
  if (value === null || value === undefined) return '—';
  return new Intl.NumberFormat('es-ES', {
    style: 'currency',
    currency: 'EUR',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value);
}
</script>

<style scoped>
.product-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.filters > div {
  flex: 1;
  min-width: 200px;
}

.filter-actions {
  display: flex;
  align-items: flex-end;
}

label {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin-bottom: 0.4rem;
}

input,
select {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.6rem 0.8rem;
  font-size: 0.9rem;
  color: #1e293b;
  background: #ffffff;
}

input:focus,
select:focus {
  outline: none;
  border-color: #1b3a6b;
  box-shadow: 0 0 0 3px rgba(27, 58, 107, 0.12);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  gap: 1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(27, 58, 107, 0.12);
  border-top-color: #1b3a6b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.alert-error {
  background: #fee2e2;
  border: 1px solid #ef4444;
  color: #991b1b;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.retry-btn {
  white-space: nowrap;
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

thead {
  background: #f8fafc;
  color: #64748b;
}

th {
  padding: 0.85rem 1rem;
  text-align: left;
  font-weight: 600;
  border-bottom: 2px solid #e2e8f0;
  white-space: nowrap;
}

tbody tr {
  border-bottom: 1px solid #e2e8f0;
  transition: background-color 0.15s ease;
}

tbody tr:hover {
  background-color: #f8fafc;
}

tbody tr.inactive {
  opacity: 0.6;
}

td {
  padding: 0.85rem 1rem;
  color: #1e293b;
}

.sku-code {
  background: #f1f5f9;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  font-weight: 600;
  color: #475569;
}

.product-link {
  color: #1b3a6b;
  text-decoration: none;
  font-weight: 600;
  display: block;
}

.product-link:hover {
  text-decoration: underline;
}

.product-long-name {
  font-size: 0.8rem;
  color: #64748b;
  margin-top: 0.2rem;
}

.brand-name,
.group-name {
  color: #64748b;
  font-size: 0.9rem;
}

.type-pill,
.status-pill,
.variants-badge {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.type-pill.type-tangible {
  background: rgba(33, 150, 243, 0.1);
  color: #2196f3;
}

.type-pill.type-service {
  background: rgba(255, 152, 0, 0.1);
  color: #ff9800;
}

.variants-badge {
  background: #f1f5f9;
  color: #64748b;
}

.status-pill.status-active {
  background: rgba(76, 175, 80, 0.1);
  color: #4caf50;
}

.status-pill.status-inactive {
  background: rgba(158, 158, 158, 0.1);
  color: #9e9e9e;
}

.price-value {
  font-weight: 600;
  color: #16a34a;
  font-size: 0.95rem;
}

.align-right {
  text-align: right;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.empty-state-block {
  text-align: center;
  color: #64748b;
  padding: 3rem 1.5rem;
  border: 2px dashed #e2e8f0;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.empty-icon {
  font-size: 3rem;
  opacity: 0.5;
}

.empty-hint {
  font-size: 0.9rem;
  color: #94a3b8;
  margin-top: -0.5rem;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}

.page-info {
  color: #64748b;
  font-weight: 500;
  min-width: 140px;
  text-align: center;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.05s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.btn-primary {
  background: #f4d03f;
  color: #1e293b;
  font-weight: 700;
  box-shadow: 0 2px 4px rgba(244, 208, 63, 0.3);
}

.btn-primary:hover {
  background: #e6c230;
  box-shadow: 0 4px 8px rgba(244, 208, 63, 0.4);
}

.btn-primary:active {
  transform: translateY(1px);
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-outline {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #1e293b;
  font-size: 0.8rem;
  padding: 0.5rem 0.8rem;
}

.btn-outline:hover {
  background: #f8fafc;
  border-color: #1b3a6b;
  color: #1b3a6b;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
  }

  .filters > div {
    min-width: 100%;
  }

  table {
    font-size: 0.85rem;
  }

  th,
  td {
    padding: 0.65rem 0.5rem;
  }

  .action-buttons {
    flex-direction: column;
    gap: 0.3rem;
  }

  .btn-outline {
    font-size: 0.75rem;
  }
}
</style>
