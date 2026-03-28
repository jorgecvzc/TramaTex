<template>
  <Navbar />
  
  <BaseDashboardPage :is-loading="isLoading">
    <template #header>
      <PageHeader title="Panel de Gestión de Catálogo" :breadcrumbs="[{ label: 'Catálogo', to: '/products' }, { label: 'Dashboard' }]">
        <template #icon><span class="material-symbols-outlined">inventory_2</span></template>
        <template #actions>
          <button class="btn btn-primary" @click="router.push('/products/new')">
            <span class="material-symbols-outlined">add_box</span>
            <span>Nuevo Producto</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <div class="product-dashboard-content">
      <!-- KPIs Superiores -->
      <section class="stats-grid mb-12">
        <div class="stat-card clickable" @click="router.push('/products')">
          <div class="stat-icon blue"><span class="material-symbols-outlined">inventory_2</span></div>
          <div class="stat-info">
            <span class="stat-label">Productos Totales</span>
            <span class="stat-value">{{ counts.totalProducts }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
        <div class="stat-card clickable" @click="router.push('/master-data/brands')">
          <div class="stat-icon yellow"><span class="material-symbols-outlined">branding_watermark</span></div>
          <div class="stat-info">
            <span class="stat-label">Marcas</span>
            <span class="stat-value">{{ counts.brands }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
        <div class="stat-card clickable" @click="router.push('/master-data/product-groups')">
          <div class="stat-icon green"><span class="material-symbols-outlined">category</span></div>
          <div class="stat-info">
            <span class="stat-label">Familias</span>
            <span class="stat-value">{{ counts.groups }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
        <div class="stat-card clickable" @click="router.push('/master-data/attributes')">
          <div class="stat-icon purple"><span class="material-symbols-outlined">tune</span></div>
          <div class="stat-info">
            <span class="stat-label">Atributos</span>
            <span class="stat-value">{{ counts.attributes }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
      </section>

      <!-- Productos Recientes -->
      <section class="dashboard-section">
        <div class="section-header">
          <span class="material-symbols-outlined text-primary">new_releases</span>
          <h2>Últimas Incorporaciones</h2>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>SKU</th>
                <th>Producto</th>
                <th>Precio Base</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="prod in recentProducts" :key="prod.id" class="row-clickable" @click="router.push(`/products/${prod.id}`)">
                <td><code class="code-badge">{{ prod.sku }}</code></td>
                <td><strong>{{ prod.name }}</strong></td>
                <td>{{ formatPrice(prod.base_price) }}</td>
                <td>
                  <span :class="['status-badge', prod.is_active ? 'status-success' : 'status-secondary']">
                    {{ prod.is_active ? 'Activo' : 'Inactivo' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <template #sidebar>
      <section class="sidebar-section">
        <div class="section-header">
          <span class="material-symbols-outlined">settings</span>
          <h2>Configuración Técnica</h2>
        </div>
        <div class="quick-links-vertical">
          <RouterLink to="/master-data/attributes" class="admin-card clickable">
            <span class="material-symbols-outlined text-secondary">tune</span>
            <div class="admin-card-info">
              <strong>Gestión de Atributos</strong>
              <p>Tallas, colores y materiales</p>
            </div>
          </RouterLink>
          <RouterLink to="/master-data/product-groups" class="admin-card clickable mt-3">
            <span class="material-symbols-outlined text-secondary">account_tree</span>
            <div class="admin-card-info">
              <strong>Árbol de Familias</strong>
              <p>Categorización jerárquica</p>
            </div>
          </RouterLink>
          <RouterLink to="/master-data/brands" class="admin-card clickable mt-3">
            <span class="material-symbols-outlined text-secondary">branding_watermark</span>
            <div class="admin-card-info">
              <strong>Marcas y Logos</strong>
              <p>Fabricantes y personalizaciones</p>
            </div>
          </RouterLink>
        </div>
      </section>

      <section class="help-notice mt-10">
        <div class="notice-header">
          <span class="material-symbols-outlined">layers</span>
          <h3>Variantes JIT</h3>
        </div>
        <p class="help-text">
          TramaTex genera variantes automáticamente al recibir pedidos para combinaciones no registradas, optimizando el tamaño del catálogo.
        </p>
      </section>
    </template>
  </BaseDashboardPage>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, RouterLink } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import { productApi } from '@/services/productApi';

const router = useRouter();
const isLoading = ref(true);
const counts = ref({ totalProducts: 0, brands: 0, groups: 0, attributes: 0 });
const recentProducts = ref([]);

async function loadProductData() {
  isLoading.value = true;
  try {
    const [prods, brands, groups, attrs] = await Promise.all([
      productApi.listProducts({ limit: 5 }),
      productApi.listBrands({ limit: 1 }),
      productApi.listProductGroups({ limit: 1 }),
      productApi.listAttributes({ limit: 1 })
    ]);
    
    recentProducts.value = prods.data || [];
    counts.value.totalProducts = prods.total || 0;
    counts.value.brands = brands.total || 0;
    counts.value.groups = groups.total || 0;
    counts.value.attributes = attrs.total || 0;
  } catch (err) {
    console.error('Error dashboard productos:', err);
  } finally {
    isLoading.value = false;
  }
}

function formatPrice(v) { return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(v || 0); }

onMounted(loadProductData);
</script>

<style scoped>
@import "@/design-system/_sections.css";
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
.stat-card { background: white; padding: 1.25rem; border-radius: 12px; border: 1px solid var(--color-border); display: flex; align-items: center; gap: 1rem; position: relative; transition: 0.2s; cursor: pointer; }
.stat-card:hover { transform: translateY(-3px); box-shadow: var(--box-shadow-md); border-color: var(--color-primary); }
.stat-icon { width: 44px; height: 44px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.stat-icon .material-symbols-outlined { font-size: 24px; }
.stat-icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.stat-icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.stat-icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.stat-icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
.stat-info { display: flex; flex-direction: column; }
.stat-label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.stat-value { font-size: 1.5rem; font-weight: 700; }
.stat-link-arrow { position: absolute; right: 1rem; color: var(--color-border); font-size: 18px; }
.dashboard-section { background: white; padding: 1.5rem; border-radius: 12px; border: 1px solid var(--color-border); }
.section-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; padding-bottom: 0.75rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.9rem; font-weight: 700; text-transform: uppercase; margin: 0; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.admin-card { display: flex; align-items: center; gap: 1rem; padding: 1rem; background: var(--color-background); border-radius: 10px; border: 1px solid transparent; text-decoration: none; color: var(--color-text-primary); transition: 0.2s; }
.admin-card:hover { background: white; border-color: var(--color-primary); transform: translateX(4px); box-shadow: var(--box-shadow-sm); }
.help-notice { padding: 1.25rem; background: var(--color-background); border-radius: 12px; border: 1px dashed var(--color-border-strong); }
.notice-header { display: flex; align-items: center; gap: 0.5rem; color: var(--color-text-secondary); font-size: 0.8rem; font-weight: 700; text-transform: uppercase; }
.help-text { font-size: 0.8rem; color: var(--color-text-secondary); margin-top: 0.5rem; line-height: 1.5; }
</style>
