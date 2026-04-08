<template>
  <BaseDashboardPage :is-loading="isLoading">
    <template #header>
      <PageHeader title="Catálogo y Almacén">
        <template #icon><span class="material-symbols-outlined">inventory_2</span></template>
        <template #actions>
          <button class="btn btn-outline btn-sm" @click="loadProductData" :disabled="isLoading">
            <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
            <span>Actualizar</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <div class="module-dashboard-content">
      <!-- 1. KPIs de Resumen -->
      <section class="stats-grid">
        <div class="stat-card clickable" @click="router.push('/products')">
          <div class="stat-icon blue"><span class="material-symbols-outlined">inventory_2</span></div>
          <div class="stat-info">
            <span class="stat-label">Productos Totales</span>
            <span class="stat-value">{{ counts.totalProducts }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/master-data/brands')">
          <div class="stat-icon yellow"><span class="material-symbols-outlined">branding_watermark</span></div>
          <div class="stat-info">
            <span class="stat-label">Marcas</span>
            <span class="stat-value">{{ counts.brands }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/master-data/product-groups')">
          <div class="stat-icon green"><span class="material-symbols-outlined">category</span></div>
          <div class="stat-info">
            <span class="stat-label">Familias</span>
            <span class="stat-value">{{ counts.groups }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/master-data/attributes')">
          <div class="stat-icon purple"><span class="material-symbols-outlined">tune</span></div>
          <div class="stat-info">
            <span class="stat-label">Atributos</span>
            <span class="stat-value">{{ counts.attributes }}</span>
          </div>
        </div>
      </section>

      <!-- 2. Accesos a Listados -->
      <section class="listings-grid">
        <RouterLink to="/products" class="listing-link">
          <span class="material-symbols-outlined">list_alt</span>
          <span>Listado de Productos</span>
        </RouterLink>
        <RouterLink to="/master-data/product-groups" class="listing-link">
          <span class="material-symbols-outlined">account_tree</span>
          <span>Familias y Categorías</span>
        </RouterLink>
        <RouterLink to="/master-data/brands" class="listing-link">
          <span class="material-symbols-outlined">branding_watermark</span>
          <span>Marcas Registradas</span>
        </RouterLink>
        <RouterLink to="/master-data/attributes" class="listing-link">
          <span class="material-symbols-outlined">tune</span>
          <span>Atributos Técnicos</span>
        </RouterLink>
      </section>

      <!-- 3. Actividad Reciente -->
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
          <span class="material-symbols-outlined">bolt</span>
          <h2>Operaciones</h2>
        </div>
        <div class="quick-actions-list">
          <RouterLink to="/products/new" class="admin-card clickable">
            <span class="material-symbols-outlined text-primary">add_box</span>
            <div class="admin-card-info">
              <strong>Nuevo Producto</strong>
              <p>Alta de referencia base</p>
            </div>
          </RouterLink>
          <RouterLink to="/master-data/attributes" class="admin-card clickable mt-2">
            <span class="material-symbols-outlined text-secondary">tune</span>
            <div class="admin-card-info">
              <strong>Gestionar Atributos</strong>
              <p>Tallas, colores, materiales</p>
            </div>
          </RouterLink>
        </div>
      </section>
    </template>
  </BaseDashboardPage>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, RouterLink } from 'vue-router';
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

function formatPrice(v) { return v ? new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(v) : '—'; }

onMounted(loadProductData);
</script>

<style scoped>
@import "@/design-system/_sections.css";

.module-dashboard-content { display: flex; flex-direction: column; gap: 1.5rem; }

.stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0.75rem; }
.stat-card { background: white; padding: 0.75rem 1rem; border-radius: 10px; border: 1px solid var(--color-border); display: flex; align-items: center; gap: 0.75rem; position: relative; transition: 0.2s; cursor: pointer; }
.stat-card:hover { transform: translateY(-2px); box-shadow: var(--box-shadow-md); border-color: var(--color-primary); }
.stat-icon { width: 40px; height: 40px; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
.stat-icon .material-symbols-outlined { font-size: 22px; }
.stat-icon.blue { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.stat-icon.yellow { background: rgba(230, 184, 0, 0.1); color: #E6B800; }
.stat-icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.stat-icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
.stat-info { display: flex; flex-direction: column; gap: 0.25rem; }
.stat-label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.stat-value { font-size: 1.25rem; font-weight: 700; }

.listings-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.75rem; }
.listing-link { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; text-decoration: none; color: var(--color-text-primary); font-size: 0.85rem; font-weight: 600; transition: 0.2s; }
.listing-link:hover { background: var(--color-background); border-color: var(--color-secondary); color: var(--color-secondary); transform: translateY(-1px); }
.listing-link .material-symbols-outlined { color: var(--color-secondary); font-size: 1.25rem; }

.dashboard-section { background: white; padding: 0.75rem 1rem; border-radius: 10px; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.section-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.75rem; padding-bottom: 0.5rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.85rem; font-weight: 700; text-transform: uppercase; margin: 0; }

.quick-actions-list { display: flex; flex-direction: column; gap: 0.75rem; }
.admin-card { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem; background: var(--color-background); border-radius: 8px; border: 1px solid transparent; text-decoration: none; color: var(--color-text-primary); transition: 0.2s; }
.admin-card:hover { background: white; border-color: var(--color-primary); transform: translateX(4px); box-shadow: var(--box-shadow-sm); }
.admin-card-info strong { font-size: 0.8rem; display: block; }
.admin-card-info p { font-size: 0.65rem; color: var(--color-text-secondary); margin: 0; }

.code-badge { background: var(--color-background); padding: 0.15rem 0.35rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.75rem; font-weight: 700; color: var(--color-secondary); }

@media (max-width: 1180px) {
  .stats-grid, .listings-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 720px) {
  .stats-grid, .listings-grid { grid-template-columns: 1fr; }
}
</style>
