<template>
  
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
        <AlertCircle :size="32" />
      </div>
      <div class="alert-content">
        <h3>Error al cargar</h3>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm mt-4" @click="router.push('/products')">Volver al catálogo</button>
      </div>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="product || mode === 'create'" class="no-print">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <div class="sticky-header-container">
        <PageHeader
          :title="mode === 'create' ? 'Nuevo Producto' : (mode === 'edit' ? `Editando ${product?.name}` : product?.name)"
          :breadcrumbs="[{ label: 'Catálogo', to: '/products/dashboard' }, { label: 'Productos', to: '/products' }, { label: mode === 'create' ? 'Alta' : product?.sku }]"
          show-back
        >

          <template #icon>
            <component :is="(product?.product_type === 'SERVICE' || formData.productType === 'SERVICE') ? Cpu : Package" :size="28" />
          </template>
          <template #actions>
            <template v-if="mode === 'detail'">
              <button class="btn btn-primary" @click="enterEditMode">
                <Pencil :size="18" /> <span>Editar Producto</span>
              </button>
            </template>
            <template v-else>
              <button class="btn btn-outline" @click="exitEditMode" :disabled="isSaving">
                <X :size="18" />
                <span>Cancelar</span>
              </button>
              <button class="btn btn-secondary" @click="saveProduct" :disabled="isSaving">
                <component :is="isSaving ? RefreshCw : Save" :size="18" :class="{ 'spin': isSaving }" />
                <span>Guardar Producto</span>
              </button>
            </template>
          </template>
        </PageHeader>

        <!-- NAVEGACIÓN POR PESTAÑAS -->
        <nav v-if="mode !== 'create'" class="entity-tabs">
          <button 
            v-for="tab in availableTabs" 
            :key="tab.id" 
            @click="activeTab = tab.id"
            :class="['tab-btn', { active: activeTab === tab.id }]"
          >
            <component :is="tab.icon" :size="18" />
            <span>{{ tab.label }}</span>
            <span v-if="tab.count !== undefined" class="tab-badge">{{ tab.count }}</span>
          </button>
        </nav>
      </div>
    </template>

    <!-- CAPA 2: CONTEXTO (Resumen) -->
    <template #summary v-if="mode !== 'create' && product">
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><Banknote :size="20" /></div>
          <div class="tag-content"><label>Precio Base</label><strong>{{ formatPrice(product.base_price) }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon green"><Layers :size="20" /></div>
          <div class="tag-content"><label>Variantes</label><strong>{{ variants.length }} activas</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><Tag :size="20" /></div>
          <div class="tag-content"><label>Marca</label><strong>{{ brand?.name || 'Genérica' }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><BadgeCheck :size="20" /></div>
          <div class="tag-content">
            <label>Estado</label>
            <strong :class="product.is_active ? 'text-success' : 'text-secondary'">{{ product.is_active ? 'Activo' : 'Inactivo' }}</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- CAPA 3: TRABAJO (Contenido) -->
    <div class="product-master-content">
      <div v-if="activeTab === 'info'" class="tab-fade-in">
        <!-- SECCIÓN: BÁSICO -->
        <FormSection title="Información Básica" icon="description">
          <div v-if="mode === 'detail'">
            <DataRow label="Nombre Comercial" :value="product?.name" icon="label" />
            <DataRow label="Referencia (SKU)" icon="fingerprint">
              <code class="code-badge">{{ product?.sku }}</code>
            </DataRow>
            <DataRow label="Descripción Larga" icon="notes">
              <p class="notes-text">{{ product?.long_name || 'Sin descripción adicional.' }}</p>
            </DataRow>
          </div>
          <div v-else>
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
              <textarea v-model="formData.longName" class="form-textarea" rows="2" placeholder="Detalles técnicos..."></textarea>
            </div>
          </div>
        </FormSection>

        <!-- SECCIÓN: CLASIFICACIÓN (Multi-grupo) -->
        <FormSection title="Clasificación y Familias" icon="category" class="mt-8">
          <div v-if="mode === 'detail'">
            <DataRow label="Tipo de Producto" :value="product?.product_type === 'TANGIBLE' ? 'Tangible / Stock' : 'Servicio / Taller'" icon="inventory" />
            <DataRow label="Marca / Fabricante" :value="brand?.name || 'Genérica'" icon="branding_watermark" />
            <DataRow label="Familias / Categorías" icon="account_tree">
              <div class="tags-cloud">
                <span v-for="g in groups" :key="g.id" class="status-badge status-secondary">{{ g.name }}</span>
                <span v-if="groups.length === 0" class="text-muted italic">Sin categorías asignadas</span>
              </div>
            </DataRow>
          </div>
          <div v-else>
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
                  <option :value="null">-- Ninguna / Genérica --</option>
                  <option v-for="b in brands" :key="b.id" :value="b.id">{{ b.name }}</option>
                </select>
              </div>
            </div>
            
            <div class="form-group mt-6">
              <label class="mb-3 block">Familias / Grupos de Producto (Seleccione uno o más)</label>
              <div class="selection-grid-container">
                <label v-for="group in productGroups" :key="group.id" class="selection-card">
                  <input type="checkbox" :value="group.id" v-model="formData.groupIds" />
                  <div class="card-content">
                    <strong>{{ group.name }}</strong>
                    <p class="text-xs text-muted">{{ group.description || 'Sin descripción' }}</p>
                  </div>
                </label>
              </div>
            </div>
          </div>
        </FormSection>

        <!-- SECCIÓN: ATRIBUTOS (Solo en creación/edición) -->
        <FormSection title="Atributos de Configuración" icon="tune" class="mt-8">
          <div v-if="mode === 'detail'">
            <DataRow label="Atributos Vinculados" icon="settings_input_component">
              <div class="tags-cloud">
                <span v-for="attrId in product?.attribute_ids" :key="attrId" class="status-badge status-info">
                  {{ allAttributes.find(a => a.id === attrId)?.name || attrId }}
                </span>
                <span v-if="!product?.attribute_ids?.length" class="text-muted italic">No define atributos (Producto único)</span>
              </div>
            </DataRow>
          </div>
          <div v-else>
            <p class="help-text mb-4">Seleccione qué atributos definen las variantes de este producto (ej: Talla, Color).</p>
            <div class="selection-grid-container">
              <label v-for="attr in allAttributes" :key="attr.id" class="selection-card">
                <input type="checkbox" :value="attr.id" v-model="formData.attributeIds" />
                <div class="card-content">
                  <strong>{{ attr.name }}</strong>
                  <code class="text-xs">{{ attr.code }}</code>
                </div>
              </label>
            </div>
          </div>
        </FormSection>

        <!-- SECCIÓN: PRECIOS -->
        <FormSection title="Estrategia de Precios" icon="payments" class="mt-8">
          <div v-if="mode === 'detail'">
            <DataRow label="Precio Base de Venta" icon="sell">
              <strong class="text-primary" style="font-size: 1.25rem">{{ formatPrice(product?.base_price) }}</strong>
              <span class="text-xs text-muted ml-2">PVP base para variantes</span>
            </DataRow>
          </div>
          <div v-else>
            <div class="form-group">
              <label>Precio Base de Venta (€) *</label>
              <input v-model.number="formData.basePrice" type="number" step="0.01" class="form-input w-48" required />
            </div>
          </div>
        </FormSection>
      </div>

      <!-- RESTO DE PESTAÑAS (Solo en modo detalle) -->
      <div v-if="activeTab === 'variants' && product" class="tab-fade-in">
        <VariantTable :product-id="product.id" :product-sku="product.sku" :variants="variants" :is-loading="isLoadingVariants" @refresh="fetchVariants" />
      </div>
      <div v-if="activeTab === 'attributes' && product" class="tab-fade-in">
        <AttributesPanel :product="product" :calculated-attributes="calculatedAttributes" :is-loading="isLoadingAttributes" @refresh="fetchCalculatedAttributes" />
      </div>
      <div v-if="activeTab === 'pricing' && product" class="tab-fade-in">
        <PricingPanel
          :product-id="product.id"
          :product-name="product.name"
          :variants="variants"
          :is-loading-variants="isLoadingVariants"
          @refresh="fetchProduct"
        />
      </div>
    </div>

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
import { 
  AlertCircle, 
  Cpu, 
  Package, 
  Pencil, 
  X, 
  RefreshCw, 
  Save, 
  FileText, 
  Layers, 
  Settings2, 
  Banknote, 
  Tag, 
  BadgeCheck, 
  Fingerprint, 
  PlusSquare 
} from 'lucide-vue-next';
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import FormSection from '@/components/shared/FormSection.vue';
import DataRow from '@/components/shared/DataRow.vue';
import VariantTable from '@/components/product/VariantTable.vue';
import AttributesPanel from '@/components/product/AttributesPanel.vue';
import PricingPanel from '@/components/product/PricingPanel.vue';
import { productApi } from '@/services/productApi';
import { useToastStore } from '@/stores/toast';

const router = useRouter();
const route = useRoute();
const toastStore = useToastStore();

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
const productGroups = ref([]);
const allAttributes = ref([]);

const isLoadingVariants = ref(false);
const isLoadingAttributes = ref(false);

const formData = reactive({
  name: '', sku: '', longName: '', productType: 'TANGIBLE', brandId: null, basePrice: 0,
  groupIds: [], attributeIds: []
});

const availableTabs = computed(() => {
  if (mode.value === 'create') return [{ id: 'info', label: 'Alta de Producto', icon: PlusSquare }];
  return [
    { id: 'info', label: 'Ficha Técnica', icon: FileText },
    { id: 'variants', label: 'Variantes (JIT)', icon: Layers, count: variants.value.length },
    { id: 'attributes', label: 'Matriz Atributos', icon: Settings2 },
    { id: 'pricing', label: 'Precios y Costes', icon: Banknote }
  ];
});

watch(() => route.params.id, async (newId) => {
  if (newId && newId !== 'new') {
    mode.value = 'detail';
    await fetchProduct();
    await Promise.all([fetchVariants(), fetchCalculatedAttributes(), loadMasters()]);
  } else {
    mode.value = 'create';
    resetForm();
    await loadMasters();
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
  } catch (err) { error.value = 'No se ha podido cargar el producto.'; }
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

async function loadMasters() {
  try {
    const [bRes, gRes, aRes] = await Promise.all([
      productApi.listBrands({ isActive: true }),
      productApi.listProductGroups({ isActive: true }),
      productApi.listAttributes({})
    ]);
    brands.value = bRes.data || [];
    productGroups.value = gRes.data || [];
    allAttributes.value = aRes.data || aRes || [];
  } catch (err) { console.error('Error maestros', err); }
}

function resetForm() {
  Object.assign(formData, { name: '', sku: '', longName: '', productType: 'TANGIBLE', brandId: null, basePrice: 0, groupIds: [], attributeIds: [] });
  product.value = null;
}

function enterEditMode() {
  Object.assign(formData, {
    name: product.value.name,
    sku: product.value.sku,
    longName: product.value.long_name,
    productType: product.value.product_type,
    brandId: product.value.brand_id,
    basePrice: product.value.base_price,
    groupIds: product.value.group_ids || [],
    attributeIds: product.value.attribute_ids || []
  });
  mode.value = 'edit';
}

function exitEditMode() {
  if (mode.value === 'edit') mode.value = 'detail';
  else router.push('/products');
}

async function saveProduct() {
  if (!formData.name || !formData.sku) {
    toastStore.addToast('El nombre y el SKU son obligatorios', 'warning')
    return
  }
  
  isSaving.value = true
  try {
    const payload = {
      name: formData.name, sku: formData.sku, long_name: formData.longName,
      product_type: formData.productType, 
      brand_id: formData.brandId || null,
      base_price: Number(formData.basePrice),
      group_ids: formData.groupIds,
      attribute_ids: formData.attributeIds
    };

    if (mode.value === 'create') {
      const newProd = await productApi.createProduct(payload);
      toastStore.addToast('Producto creado exitosamente', 'success')
      await router.push(`/products/${newProd.id}`);
    } else {
      await productApi.updateProduct(product.value.id, payload);
      toastStore.addToast('Ficha técnica actualizada correctamente', 'success')
      await fetchProduct();
      mode.value = 'detail';
    }
  } catch (err) {
    toastStore.addToast('Error al guardar: ' + err.message, 'error')
  } finally {
    isSaving.value = false
  }
}

function formatPrice(v) { return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(v || 0); }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'long', day: 'numeric' }) : '—'; }
</script>

<style scoped>
@import "@/design-system/_sections.css";

.sticky-header-container { background: white; margin-top: -1.5rem; padding-top: 1.5rem; border-bottom: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.entity-tabs { display: flex; max-width: 1300px; margin: 0 auto; padding: 0 2rem; gap: 0.25rem; }
.tab-btn { display: flex; align-items: center; gap: 0.6rem; padding: 1rem 1.25rem; background: transparent; border: none; border-bottom: 3px solid transparent; color: var(--color-text-secondary); font-weight: 700; cursor: pointer; transition: all 0.2s; font-size: 0.85rem; text-transform: uppercase; letter-spacing: 0.025em; margin-bottom: -1px; }
.tab-btn:hover { color: var(--color-primary); background: rgba(0,0,0,0.02); }
.tab-btn.active { border-bottom-color: var(--color-secondary); color: var(--color-secondary); background: rgba(0, 35, 149, 0.03); }
.tab-btn .material-symbols-outlined { font-size: 18px; }
.tab-badge { background: var(--color-background); color: var(--color-text-secondary); padding: 0.1rem 0.5rem; border-radius: 20px; font-size: 0.7rem; }
.tab-btn.active .tab-badge { background: var(--color-secondary); color: white; }

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

/* Selection Grid */
.selection-grid-container { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 1rem; }
.selection-card { display: flex; align-items: flex-start; gap: 1rem; padding: 1rem; background: var(--color-background); border: 1px solid var(--color-border); border-radius: 10px; cursor: pointer; transition: all 0.2s; }
.selection-card:hover { background: white; border-color: var(--color-primary); }
.selection-card input[type="checkbox"] { width: 18px; height: 18px; margin-top: 0.2rem; cursor: pointer; }
.selection-card .card-content { display: flex; flex-direction: column; gap: 0.2rem; }
.selection-card strong { font-size: 0.85rem; color: var(--color-text-primary); }

.tags-cloud { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.w-48 { width: 12rem; }
.tab-fade-in { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
