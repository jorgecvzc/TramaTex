<template>
  <BaseDashboardPage :is-loading="isLoading">
    <template #header>
      <BasePageHeader 
        title="Panel de Control" 
        :subtitle="`Bienvenido, ${userName}. Gestión de TramaTex para el ${todayDate}.`"
      >
        <template #actions>
          <div class="context-tags">
            <span class="context-tag">Vista Ejecutiva</span>
            <span class="context-tag">Operaciones</span>
            <span class="context-tag" :class="{ loading: isLoading }">Datos en tiempo real</span>
          </div>
          <button @click="loadStats" class="btn-refresh" :disabled="isLoading" title="Actualizar datos">
            <RefreshCw :size="20" :class="{ 'spin': isLoading }" />
          </button>
        </template>
      </BasePageHeader>
    </template>

    <div class="dashboard-content-wrapper">
      <!-- KPIs INTERACTIVOS -->
      <section class="kpi-grid">
        <div class="kpi-card clickable" @click="navigateTo('/sales/dashboard')">
          <div class="kpi-icon blue"><CreditCard :size="24" /></div>
          <div class="kpi-data">
            <label>Ventas Mes</label>
            <strong>{{ salesStats.monthlyTotal }}</strong>
          </div>
        </div>
        <div class="kpi-card clickable" @click="navigateTo('/sales/orders?status=PENDING')">
          <div class="kpi-icon green"><ShoppingCart :size="24" /></div>
          <div class="kpi-data">
            <label>Pedidos Pendientes</label>
            <strong>{{ salesStats.pendingOrders }}</strong>
          </div>
        </div>
        <div class="kpi-card clickable" @click="navigateTo('/mes/work-orders?status=IN_PROGRESS')">
          <div class="kpi-icon yellow"><Cpu :size="24" /></div>
          <div class="kpi-data">
            <label>Órdenes en Taller</label>
            <strong>{{ mesStats.activeWorkOrders }}</strong>
          </div>
        </div>
        <div class="kpi-card clickable" @click="navigateTo('/parties')">
          <div class="kpi-icon purple"><Users :size="24" /></div>
          <div class="kpi-data">
            <label>Entidades</label>
            <strong>{{ partyStats.totalParties }}</strong>
          </div>
        </div>
      </section>

      <!-- ACCESOS DIRECTOS -->
      <section class="ops-section section-block section-spacing">
        <div class="section-title-alt">
          <Rocket :size="18" />
          <h2>Accesos Directos</h2>
          <span class="section-tag">Atajos de operación</span>
        </div>
        <div class="ops-grid">
          <RouterLink to="/sales/orders/new" class="op-item">
            <div class="op-icon"><ShoppingCart :size="28" /></div>
            <span>Nuevo Pedido</span>
          </RouterLink>
          <RouterLink to="/sales/tickets/new" class="op-item highlight">
            <div class="op-icon"><Receipt :size="28" /></div>
            <span>Venta Directa</span>
          </RouterLink>
          <RouterLink to="/products/new" class="op-item">
            <div class="op-icon"><PlusSquare :size="28" /></div>
            <span>Nuevo Producto</span>
          </RouterLink>
          <RouterLink to="/parties/new" class="op-item">
            <div class="op-icon"><UserPlus :size="28" /></div>
            <span>Nueva Entidad</span>
          </RouterLink>
        </div>
      </section>

      <!-- DASHBOARDS DE MÓDULO -->
      <section class="modules-section section-block section-spacing">
        <div class="section-title-alt">
          <LayoutGrid :size="18" />
          <h2>Módulos Principales</h2>
          <span class="section-tag">Gestión por dominio</span>
        </div>
        <div class="modules-grid">
          <RouterLink to="/sales/dashboard" class="module-link-card">
            <div class="m-icon blue"><Wallet :size="22" /></div>
            <div class="m-info">
              <strong>Ventas</strong>
              <p>Gestión de presupuestos, pedidos y facturación.</p>
            </div>
            <ChevronRight class="arrow" :size="18" />
          </RouterLink>
          <RouterLink to="/products/dashboard" class="module-link-card">
            <div class="m-icon yellow"><Package :size="22" /></div>
            <div class="m-info">
              <strong>Catálogo</strong>
              <p>Control de productos, stock y atributos.</p>
            </div>
            <ChevronRight class="arrow" :size="18" />
          </RouterLink>
          <RouterLink to="/parties/dashboard" class="module-link-card">
            <div class="m-icon green"><Users :size="22" /></div>
            <div class="m-info">
              <strong>Entidades</strong>
              <p>Base de datos de clientes y proveedores.</p>
            </div>
            <ChevronRight class="arrow" :size="18" />
          </RouterLink>
          <RouterLink to="/mes/dashboard" class="module-link-card">
            <div class="m-icon purple"><Cpu :size="22" /></div>
            <div class="m-info">
              <strong>Taller (MES)</strong>
              <p>Monitorización de producción y órdenes de trabajo.</p>
            </div>
            <ChevronRight class="arrow" :size="18" />
          </RouterLink>
        </div>
      </section>
    </div>

    <template #sidebar>
      <section v-if="isAdmin" class="admin-side-section">
        <div class="ops-header">
          <ShieldCheck :size="18" />
          <h2>Sistema</h2>
          <span class="section-tag">Admin</span>
        </div>
        <div class="side-links">
          <RouterLink to="/admin/users" class="side-link-item">
            <UserCog :size="20" />
            <span>Gestión de Usuarios</span>
          </RouterLink>
          <RouterLink to="/admin/print-profile" class="side-link-item">
            <Receipt :size="20" />
            <span>Perfil de Impresión</span>
          </RouterLink>
          <RouterLink to="/dev/design-system" class="side-link-item dev-link">
            <Palette :size="20" />
            <span>Design System</span>
          </RouterLink>
        </div>
      </section>
    </template>
  </BaseDashboardPage>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import { 
  RefreshCw, 
  CreditCard, 
  ShoppingCart, 
  Cpu, 
  Users, 
  Rocket, 
  Receipt, 
  PlusSquare, 
  UserPlus, 
  LayoutGrid, 
  Wallet, 
  Package, 
  ChevronRight, 
  ShieldCheck, 
  UserCog, 
  Palette 
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import salesApi from '@/services/salesApi'
import { mesApi } from '@/services/mesApi'
import { partyApi } from '@/services/partyApi'

const router = useRouter()
const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)
const userName = computed(() => authStore.user?.email.split('@')[0] || 'Usuario')
const todayDate = computed(() => new Date().toLocaleDateString('es-ES', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }))
const isLoading = ref(true)
const salesStats = ref({ monthlyTotal: '0,00 €', pendingOrders: 0 })
const mesStats = ref({ activeWorkOrders: 0 })
const partyStats = ref({ totalParties: 0 })

async function loadStats() {
  isLoading.value = true
  console.log('[Dashboard] Loading statistics...');
  
  try {
    const now = new Date()
    const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)

    // Execute in parallel with individual capture
    const [ordersRes, workOrdersRes, partiesRes, invoicesRes] = await Promise.allSettled([
      salesApi.listOrders({ status: 'PENDING', limit: 1 }),
      mesApi.listWorkOrders({ status: 'IN_PROGRESS' }),
      partyApi.listParties({ pageSize: 1, pageNumber: 1 }),
      salesApi.listInvoices({}) // No initial date filter for compatibility
    ])
    
    // 1. Pending Orders
    if (ordersRes.status === 'fulfilled' && ordersRes.value) {
      salesStats.value.pendingOrders = ordersRes.value.total || 0
      console.log('[Dashboard] Orders loaded:', salesStats.value.pendingOrders);
    } else {
      console.warn('[Dashboard] Failed to load orders:', (ordersRes as any).reason);
    }
    
    // 2. Monthly Sales (Invoices for the current month)
    if (invoicesRes.status === 'fulfilled' && invoicesRes.value) {
      const invData = invoicesRes.value.data || []
      const monthInvoices = invData.filter((inv: any) => {
        const d = new Date(inv.issueDate || inv.invoiceDate || inv.invoice_date)
        const isThisMonth = d >= firstDayOfMonth
        const isActive = inv.status !== 'CANCELLED' && inv.status !== 'VOID'
        return isThisMonth && isActive
      })

      const totalAmount = monthInvoices.reduce((acc: number, inv: any) => {
        const amount = typeof inv.total === 'object' ? inv.total.amount : (inv.totalAmount || 0)
        return acc + (Number(amount) || 0)
      }, 0)

      salesStats.value.monthlyTotal = salesApi.formatMoney(totalAmount)
      console.log('[Dashboard] Monthly sales calculated:', salesStats.value.monthlyTotal, `(${monthInvoices.length} invoices)`);
    } else {

      console.warn('[Dashboard] Failed to load invoices:', (invoicesRes as any).reason);
    }

    // 3. Work Orders (MES)
    if (workOrdersRes.status === 'fulfilled' && workOrdersRes.value) {
      const workOrders = workOrdersRes.value
      mesStats.value.activeWorkOrders = Array.isArray(workOrders) ? workOrders.length : 0
      console.log('[Dashboard] MES orders loaded:', mesStats.value.activeWorkOrders);
    } else {
      console.warn('[Dashboard] Failed to load MES orders:', (workOrdersRes as any).reason);
    }
    
    // 4. Entities
    if (partiesRes.status === 'fulfilled' && partiesRes.value) {
      partyStats.value.totalParties = partiesRes.value.total || 0
      console.log('[Dashboard] Entities loaded:', partyStats.value.totalParties);
    } else {
      console.warn('[Dashboard] Failed to load entities:', (partiesRes as any).reason);
    }

  } catch (err) { 
    console.error('[Dashboard] Unexpected critical error:', err) 
  } finally { 
    isLoading.value = false 
  }
}
function navigateTo(path: string) { router.push(path) }
onMounted(loadStats)
</script>

<style scoped>
.dashboard-content-wrapper {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.btn-refresh { 
  background: white; 
  border: 1px solid var(--color-border); 
  padding: 0.5rem; 
  border-radius: 8px; 
  cursor: pointer; 
  color: var(--color-text-secondary); 
  transition: 0.2s; 
  display: flex;
  align-items: center;
  justify-content: center;
}
.btn-refresh:hover { color: var(--color-primary); border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }

.section-spacing { margin-top: var(--spacing-md); }

.context-tags { display: flex; flex-wrap: wrap; gap: var(--spacing-xs); align-items: center; }
.context-tag {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  padding: 0.24rem 0.6rem;
  border-radius: 999px;
  background: #e2e8f0;
  color: #334155;
}
.context-tag.loading { background: #dbeafe; color: #1d4ed8; }

.kpi-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1.1rem; }
.kpi-card { background: white; padding: 1.25rem; border-radius: 12px; display: flex; align-items: center; gap: 1rem; border: 1px solid var(--color-border); }
.kpi-card.clickable { cursor: pointer; transition: 0.2s; }
.kpi-card.clickable:hover { transform: translateY(-3px); border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }
.kpi-icon { width: 48px; height: 48px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.kpi-icon :deep(svg) { width: 24px; height: 24px; }
.kpi-icon.blue { background: rgba(37, 99, 235, 0.08); color: #2563eb; }
.kpi-icon.green { background: rgba(22, 163, 74, 0.08); color: #16a34a; }
.kpi-icon.yellow { background: rgba(217, 119, 6, 0.08); color: #d97706; }
.kpi-icon.purple { background: rgba(147, 51, 234, 0.08); color: #9333ea; }
.kpi-data { display: flex; flex-direction: column; }
.kpi-data label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }
.kpi-data strong { font-size: 1.4rem; color: var(--color-text-primary); font-weight: 800; }

.section-block {
  background: white;
  border: 1px solid var(--color-border);
  border-radius: 14px;
  box-shadow: var(--box-shadow-sm);
}
.section-title-alt { display: flex; align-items: center; gap: var(--spacing-xs); margin-bottom: 0.85rem; }
.section-title-alt h2 { font-size: 0.9rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-secondary); }
.section-title-alt :deep(svg) { color: var(--color-text-secondary); }
.section-tag {
  margin-left: auto;
  font-size: 0.66rem;
  padding: 0.18rem 0.55rem;
  border-radius: 999px;
  background: var(--color-background);
  color: var(--color-text-secondary);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.ops-section,
.modules-section { padding: clamp(1rem, 2vw, 1.35rem); }
.ops-header { display: flex; align-items: center; gap: var(--spacing-xs); margin-bottom: 1.5rem; }
.ops-header h2 { font-size: 0.9rem; font-weight: 800; text-transform: uppercase; margin: 0; color: var(--color-text-primary); }
.ops-header :deep(svg) { color: var(--color-primary); }
.ops-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(140px, 1fr)); gap: var(--spacing-sm); }
.op-item { display: flex; flex-direction: column; align-items: center; padding: 1.25rem; background: var(--color-background); border-radius: 12px; text-decoration: none; color: var(--color-text-primary); font-weight: 700; font-size: 0.8rem; transition: 0.2s; text-align: center; gap: 0.75rem; border: 1px solid var(--color-border); }
.op-item:hover { transform: translateY(-3px); border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); color: var(--color-primary); }
.op-icon { color: var(--color-primary); }
.op-icon :deep(svg) { width: 28px; height: 28px; }
.op-item.highlight { background: rgba(230, 184, 0, 0.05); }

.modules-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: var(--spacing-sm); }
.module-link-card { display: flex; align-items: center; gap: 1rem; padding: 1.1rem; background: white; border-radius: 12px; border: 1px solid var(--color-border); text-decoration: none; transition: 0.2s; }
.module-link-card:hover { transform: translateX(5px); border-color: var(--color-secondary); box-shadow: var(--box-shadow-sm); }
.m-icon { width: 44px; height: 44px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.m-icon :deep(svg) { width: 22px; height: 22px; }
.m-icon.blue { background: rgba(37, 99, 235, 0.05); color: #2563eb; }
.m-icon.yellow { background: rgba(217, 119, 6, 0.05); color: #d97706; }
.m-icon.green { background: rgba(22, 163, 74, 0.05); color: #16a34a; }
.m-icon.purple { background: rgba(147, 51, 234, 0.05); color: #9333ea; }
.m-info { min-width: 0; }
.m-info strong { display: block; font-size: 0.9rem; color: var(--color-text-primary); }
.m-info p { font-size: 0.75rem; color: var(--color-text-secondary); margin: 0.1rem 0 0; }
.module-link-card .arrow { margin-left: auto; color: var(--color-border); font-size: 18px; }

.side-links { display: flex; flex-direction: column; gap: var(--spacing-sm); }
.side-link-item { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: var(--color-background); border-radius: 8px; text-decoration: none; color: var(--color-text-primary); font-size: 0.85rem; font-weight: 600; border: 1px solid transparent; transition: 0.2s; }
.side-link-item:hover { background: white; border-color: var(--color-primary); color: var(--color-primary); transform: translateX(3px); }
.side-link-item :deep(svg) { width: 20px; height: 20px; color: var(--color-text-secondary); }
.dev-link { opacity: 0.7; }

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

@media (max-width: 760px) {
  .ops-grid { grid-template-columns: 1fr; }
  .modules-grid { grid-template-columns: 1fr; }
}
</style>
