<template>
  <Navbar class="no-print" />
  
  <BaseEntityPage v-if="isLoading" class="no-print">
    <template #header>
      <PageHeader title="Cargando..." :breadcrumbs="[{ label: 'Catálogo', to: '/products' }, { label: 'Productos' }]" />
    </template>
    <div class="loading-state card">
      <div class="spinner"></div>
      <p>Consultando ficha técnica...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="error" class="no-print">
    <template #header>
      <PageHeader title="Error" :breadcrumbs="[{ label: 'Catálogo', to: '/products' }, { label: 'Productos' }]" />
    </template>
    <div class="alert-card card">
      <div class="alert-icon-wrapper error">
        <span class="material-symbols-outlined">error</span>
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/products')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="product || mode === 'create'" class="no-print">
    <!-- CAPA 1: IDENTIDAD + PESTAÑAS STICKY -->
    <template #header>
      <div class="sticky-identity-wrapper">
        <PageHeader 
          :title="mode === 'create' ? 'Nuevo Producto' : (mode === 'edit' ? `Editando ${product?.name}` : product?.name)" 
          :breadcrumbs="[{ label: 'Operaciones', to: '/products' }, { label: 'Productos', to: '/products' }, { label: mode === 'create' ? 'Alta' : product?.sku }]"
        >
          <template #icon>
            <span class="material-symbols-outlined">{{ product?.product_type === 'SERVICE' ? 'precision_manufacturing' : 'inventory_2' }}</span>
          </template>
          <template #actions>
            <template v-if="mode === 'detail'">
              <button class="btn btn-primary" @click="enterEditMode">
                <span class="material-symbols-outlined">edit</span> <span>Editar Producto</span>
              </button>
            </template>
            <template v-else>
              <button class="btn btn-outline" @click="exitEditMode" :disabled="isSaving">Cancelar</button>
              <button class="btn btn-secondary" @click="saveProduct" :disabled="isSaving">
                <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
                <span>{{ isSaving ? 'Guardando...' : 'Guardar Producto' }}</span>
              </button>
            </template>
          </template>
        </PageHeader>

        <!-- PESTAÑAS (Solo visibles si el producto ya existe) -->
        <nav v-if="mode !== 'create'" class="tabs-navigation-bar">
          <button
            v-for="tab in availableTabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="['tab-link', { active: activeTab === tab.id }]"
          >
            <span class="material-symbols-outlined">{{ tab.icon }}</span>
            <span class="tab-label">{{ tab.label }}</span>
            <span v-if="tab.count !== undefined" class="tab-count">{{ tab.count }}</span>
          </button>
        </nav>
      </div>
    </template>

    <!-- CAPA 2: CONTEXTO (KPIs de Producto) -->
    <template #summary v-if="mode !== 'create' && product">
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><span class="material-symbols-outlined">payments</span></div>
          <div class="tag-content"><label>Precio Base</label><strong>{{ formatPrice(product.base_price) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><span class="material-symbols-outlined">layers</span></div>
          <div class="tag-content"><label>Variantes</label><strong>{{ variants.length }} activas</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><span class="material-symbols-outlined">category</span></div>
          <div class="tag-content"><label>Marca</label><strong>{{ brand?.name || 'Genérica' }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><span class="material-symbols-outlined">verified</span></div>
          <div class="tag-content">
            <label>Estado</label>
            <strong :class="product.is_active ? 'text-success' : 'text-secondary'">{{ product.is_active ? 'Activo' : 'Inactivo' }}</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- CAPA 3: TRABAJO (Contenido según Pestaña) -->
    <div class="product-work-area">
      <!-- PESTAÑA: FICHA TÉCNICA -->
      <div v-if="activeTab === 'info'">
        <div v-if="mode === 'detail'">
          <ProductDetailInfo 
            :product="product" 
            :brand="brand" 
            :groups="groups" 
            readonly
          />
        </div>
        <div v-else>
          <!-- Formulario de Edición/Creación inyectado directamente para control maestro -->
          <FormSection title="Información Básica" icon="description">
            <div class="form-row">
              <div class="form-group">
                <label>Nombre del Producto *</label>
                <input v-model="formData.name" type="text" class="form-input" required placeholder="Nombre comercial" />
              </div>
              <div class="form-group">
                <label>SKU / Referencia *</label>
                <input v-model="formData.sku" type="text" class="form-input" required placeholder="Código único" />
              </div>
            </div>
            <div class="form-group mt-4">
              <label>Descripción Larga (Catálogo)</label>
              <textarea v-model="formData.longName" class="form-textarea" rows="2"></textarea>
            </div>
          </FormSection>

          <FormSection title="Clasificación y Tipo" icon="category" class="mt-8">
            <div class="form-row">
              <div class="form-group">
                <label>Tipo de Producto</label>
                <select v-model="formData.productType" class="form-input" :disabled="mode === 'edit'">
                  <option value="TANGIBLE">Tangible (Ropa, EPIs, etc.)</option>
                  <option value="SERVICE">Servicio (Bordado, Arreglo, etc.)</option>
                </select>
              </div>
              <div class="form-group">
                <label>Marca</label>
                <select v-model="formData.brandId" class="form-input">
                  <option :value="null">-- Ninguna --</option>
                  <option v-for="b in brands" :key="b.id" :value="b.id">{{ b.name }}</option>
                </select>
              </div>
            </div>
          </FormSection>

          <FormSection title="Configuración Comercial" icon="payments" class="mt-8">
            <div class="form-group">
              <label>Precio Base de Venta (€)</label>
              <input v-model.number="formData.basePrice" type="number" step="0.01" class="form-input w-48" />
            </div>
          </FormSection>
        </div>
      </div>

      <!-- PESTAÑA: VARIANTES (JIT) -->
      <div v-if="activeTab === 'variants' && product">
        <VariantTable
          :product-id="product.id"
          :product-sku="product.sku"
          :variants="variants"
          :is-loading="isLoadingVariants"
          @refresh="fetchVariants"
        />
      </div>

      <!-- PESTAÑA: ATRIBUTOS (MATRIZ) -->
      <div v-if="activeTab === 'attributes' && product">
        <AttributesPanel
          :product="product"
          :calculated-attributes="calculatedAttributes"
          :is-loading="isLoadingAttributes"
          @refresh="fetchCalculatedAttributes"
        />
      </div>

      <!-- PESTAÑA: PRECIOS -->
      <div v-if="activeTab === 'pricing' && product">
        <PricingPanel
          :product="product"
          @refresh="fetchProduct"
        />
      </div>
    </div>

    <!-- FOOTER AUDITORÍA -->
    <template #footer v-if="mode === 'detail' && product">
      <div class="audit-info">
        <p>Producto registrado el {{ formatDate(product.created_at) }}.</p>
        <p>UUID: <code>{{ product.id }}</code></p>
      </div>
    </template>
  </BaseEntityPage>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import ProductDetailInfo from '@/components/product/ProductDetailInfo.vue';
import VariantTable from '@/components/product/VariantTable.vue';
import AttributesPanel from '@/components/product/AttributesPanel.vue';
import PricingPanel from '@/components/product/PricingPanel.vue';
import { productApi } from '@/services/productApi';

const route = useRoute();
const router = useRouter();

const mode = ref('detail');
const activeTab = ref('info');
const isLoading = ref(true);
const isSaving = ref(false);
const error = ref('');

const product = ref(null);
const brand = ref(null);
const groups = ref([]);
const variants = ref([]);
const calculatedAttributes = ref([]);
const brands = ref([]);

const isLoadingVariants = ref(false);
const isLoadingAttributes = ref(false);

const formData = reactive({
  name: '', sku: '', longName: '', productType: 'TANGIBLE', brandId: null, basePrice: 0,
});

const availableTabs = computed(() => {
  if (mode.value === 'create') return [{ id: 'info', label: 'General', icon: 'info' }];
  return [
    { id: 'info', label: 'Ficha Técnica', icon: 'description' },
    { id: 'variants', label: 'Variantes (JIT)', icon: 'layers', count: variants.value.length },
    { id: 'attributes', label: 'Atributos', icon: 'settings_input_component' },
    { id: 'pricing', label: 'Precios y Costes', icon: 'payments' }
  ];
});

watch(() => route.params.id, async (newId) => {
  if (newId && newId !== 'new') {
    mode.value = 'detail';
    await fetchProduct();
    await Promise.all([fetchVariants(), fetchCalculatedAttributes()]);
  } else {
    mode.value = 'create';
    resetForm();
    await loadBrands();
    isLoading.value = false;
  }
}, { immediate: true });

async function fetchProduct() {
  const id = route.params.id;
  if (!id || id === 'new') return;
  isLoading.value = true; error.value = '';
  try {
    const data = await productApi.getProduct(id);
    product.value = data;
    if (data.brand_id) brand.value = await productApi.getBrand(data.brand_id);
    if (data.group_ids?.length > 0) {
      const fetched = await Promise.all(data.group_ids.map(gid => productApi.getProductGroup(gid).catch(() => null)));
      groups.value = fetched.filter(g => g !== null);
    }
  } catch (err) { error.value = 'No se pudo cargar el producto.'; }
  finally { isLoading.value = false; }
}

async function fetchVariants() {
  if (!product.value) return;
  isLoadingVariants.value = true;
  try { const data = await productApi.listProductVariants(product.value.id); variants.value = data.variants || data || []; }
  catch (err) { variants.value = []; } finally { isLoadingVariants.value = false; }
}

async function fetchCalculatedAttributes() {
  if (!product.value) return;
  isLoadingAttributes.value = true;
  try { const data = await productApi.getCalculatedAttributes(product.value.id); calculatedAttributes.value = data.attributes || data || []; }
  catch (err) { calculatedAttributes.value = []; } finally { isLoadingAttributes.value = false; }
}

async function loadBrands() {
  try { const res = await productApi.listBrands({ isActive: true }); brands.value = res.data || []; } catch (err) {}
}

function resetForm() {
  Object.assign(formData, { name: '', sku: '', longName: '', productType: 'TANGIBLE', brandId: null, basePrice: 0 });
  product.value = null;
}

function enterEditMode() {
  Object.assign(formData, {
    name: product.value.name,
    sku: product.value.sku,
    longName: product.value.long_name,
    productType: product.value.product_type,
    brandId: product.value.brand_id,
    basePrice: product.value.base_price
  });
  loadBrands();
  mode.value = 'edit';
}

function exitEditMode() {
  if (mode.value === 'edit') mode.value = 'detail';
  else router.push('/products');
}

async function saveProduct() {
  if (!formData.name || !formData.sku) { alert('Nombre y SKU son obligatorios'); return; }
  isSaving.value = true;
  try {
    const payload = {
      name: formData.name, sku: formData.sku, long_name: formData.longName,
      product_type: formData.productType, brand_id: formData.brandId || undefined,
      base_price: Number(formData.basePrice)
    };

    if (mode.value === 'create') {
      const newProd = await productApi.createProduct(payload);
      await router.push(`/products/${newProd.id}`);
    } else {
      await productApi.updateProduct(product.value.id, payload);
      await fetchProduct();
      mode.value = 'detail';
    }
  } catch (err) { alert('Error al guardar: ' + err.message); }
  finally { isSaving.value = false; }
}

function formatPrice(v) { return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(v || 0); }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'long', day: 'numeric' }) : '—'; }
</script>

<style scoped>
@import "@/design-system/_sections.css";

.sticky-identity-wrapper {
  background: white;
  margin-top: -1.5rem;
  padding-top: 1.5rem;
  border-bottom: 1px solid var(--color-border);
}

.tabs-navigation-bar {
  display: flex;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem;
  gap: 0.5rem;
}

.tab-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  color: var(--color-text-secondary);
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.9rem;
  margin-bottom: -1px;
}

.tab-link:hover { color: var(--color-text-primary); background: rgba(0,0,0,0.02); }
.tab-link.active { border-bottom-color: var(--color-primary); color: var(--color-primary); background: rgba(0, 35, 149, 0.03); }
.tab-link .material-symbols-outlined { font-size: 20px; }

.tab-count {
  background: var(--color-background);
  color: var(--color-text-secondary);
  padding: 0.1rem 0.5rem;
  border-radius: 20px;
  font-size: 0.7rem;
}
.tab-link.active .tab-count { background: var(--color-primary); color: white; }

.overview-tags-row { display: flex; flex-wrap: wrap; gap: 1rem; }
.summary-tag { flex: 1; min-width: 220px; padding: 0.75rem 1.25rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 1rem; box-shadow: var(--box-shadow-sm); }
.summary-tag .icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.summary-tag .icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.summary-tag .icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.summary-tag .icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.summary-tag .icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }

.tag-content { display: flex; flex-direction: column; line-height: 1.2; }
.tag-content label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.tag-content strong { font-size: 1rem; color: var(--color-text-primary); }

.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }

.w-48 { width: 12rem; }
.link-primary { color: var(--color-secondary); font-weight: 600; text-decoration: none; }
.notes-text { font-style: italic; color: var(--color-text-secondary); }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.85rem; font-weight: 700; }
</style>
