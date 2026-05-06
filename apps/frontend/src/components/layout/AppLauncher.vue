<template>
  <div v-if="isOpen" class="app-launcher-overlay" @click.self="$emit('close')">
    <div class="app-launcher-card">
      <header class="launcher-header">
        <div class="logo-brand">
          <span class="logo-icon">TX</span>
          <span class="logo-text">TramaTex</span>
        </div>
        <button class="close-button" @click="$emit('close')">
          <X :size="24" />
        </button>
      </header>

      <div class="launcher-grid">
        <!-- SECTION: SALES -->
        <section class="launcher-section">
          <h3 class="section-title">
            <CreditCard :size="20" class="section-icon" />
            Ventas
          </h3>
          <div class="section-links">
            <RouterLink to="/sales/quotes" class="launcher-link" @click="$emit('close')">Presupuestos</RouterLink>
            <RouterLink to="/sales/orders" class="launcher-link" @click="$emit('close')">Pedidos</RouterLink>
            <RouterLink to="/sales/delivery-notes" class="launcher-link" @click="$emit('close')">Albaranes</RouterLink>
            <RouterLink to="/sales/invoices" class="launcher-link" @click="$emit('close')">Facturas</RouterLink>
            <RouterLink to="/sales/tickets/new" class="launcher-link highlight-link" @click="$emit('close')">Venta Directa (TPV)</RouterLink>
          </div>
        </section>

        <!-- SECTION: PRODUCTION -->
        <section class="launcher-section">
          <h3 class="section-title">
            <Factory :size="20" class="section-icon" />
            Producción
          </h3>
          <div class="section-links">
            <RouterLink to="/mes/work-orders" class="launcher-link" @click="$emit('close')">Órdenes de Trabajo</RouterLink>
            <RouterLink to="/mes/tasks" class="launcher-link" @click="$emit('close')">Tareas Pendientes</RouterLink>
            <RouterLink to="/mes/terminal" class="launcher-link highlight-link" @click="$emit('close')">Terminal de Operario</RouterLink>
            <RouterLink to="/mes/positions" class="launcher-link" @click="$emit('close')">Gestión de Posiciones</RouterLink>
          </div>
        </section>

        <!-- SECTION: ENTITIES -->
        <section class="launcher-section">
          <h3 class="section-title">
            <Users :size="20" class="section-icon" />
            Entidades
          </h3>
          <div class="section-links">
            <RouterLink to="/parties" class="launcher-link" @click="$emit('close')">Clientes y Proveedores</RouterLink>
            <RouterLink to="/parties#groups" class="launcher-link" @click="$emit('close')">Grupos de Entidad</RouterLink>
          </div>

        </section>

        <!-- SECTION: CATALOG (Moved here) -->
        <section class="launcher-section">
          <h3 class="section-title">
            <Package :size="20" class="section-icon" />
            Catálogo
          </h3>
          <div class="section-links">
            <RouterLink to="/products" class="launcher-link" @click="$emit('close')">Listado de Productos</RouterLink>
            <RouterLink to="/master-data/product-groups" class="launcher-link" @click="$emit('close')">Categorías y Familias</RouterLink>
            <RouterLink to="/master-data/brands" class="launcher-link" @click="$emit('close')">Marcas</RouterLink>
            <RouterLink to="/master-data/attributes" class="launcher-link" @click="$emit('close')">Atributos Técnicos</RouterLink>
          </div>
        </section>

        <!-- SECTION: ADMINISTRATION -->
        <section v-if="isAdmin" class="launcher-section admin-section">
          <h3 class="section-title">
            <ShieldCheck :size="20" class="section-icon" />
            Sistema
          </h3>
          <div class="section-links">
            <RouterLink to="/admin/users" class="launcher-link" @click="$emit('close')">Usuarios y Permisos</RouterLink>
            <RouterLink to="/admin/print-profile" class="launcher-link" @click="$emit('close')">Perfil Fiscal</RouterLink>
            <RouterLink to="/dev/design-system" class="launcher-link dev-link" @click="$emit('close')">Design System</RouterLink>
          </div>
        </section>
      </div>

      <footer class="launcher-footer">
        <p>TramaTex ERP v1.0 &copy; 2026</p>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { 
  X, 
  CreditCard, 
  Factory, 
  Users, 
  Package, 
  ShieldCheck 
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)

const handleEscape = (e) => {
  if (e.key === 'Escape' && props.isOpen) {
    emit('close')
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleEscape)
})

// Lock body scroll when open
watch(() => props.isOpen, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>

<style scoped>
.app-launcher-overlay {
  position: fixed;
  top: 76px;
  left: 64px;
  width: calc(100vw - 64px);
  height: calc(100vh - 76px);
  background-color: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  z-index: 3000;
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  padding: 1rem 1.5rem 1.5rem 1.5rem; /* Added top padding */
}

.app-launcher-card {
  background: white;
  width: 100%;
  max-width: 900px;
  border-radius: 16px; /* Rounded corners on all corners */
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  max-height: 82vh;
  overflow: hidden;
  animation: slide-in 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slide-in {
  from { opacity: 0; transform: translateX(-20px); }
  to { opacity: 1; transform: translateX(0); }
}

.launcher-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #f1f5f9;
}

.logo-brand {
  display: flex;
  align-items: center;
}

.logo-icon {
  width: 2.5rem;
  height: 2.5rem;
  background-color: var(--color-secondary);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  font-weight: 900;
  font-size: 1.1rem;
}

.logo-text {
  margin-left: 1rem;
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--color-text-primary);
}

.close-button {
  background: var(--color-background);
  border: none;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all 0.2s;
}

.close-button:hover {
  background: var(--color-border);
  color: var(--color-text-primary);
  transform: rotate(90deg);
}

.launcher-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 2rem;
  padding: 2rem;
  overflow-y: auto;
}

.launcher-section {
  display: flex;
  flex-direction: column;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--color-background);
}

.section-icon {
  color: var(--color-info);
}

.section-links {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.launcher-link {
  text-decoration: none;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
  padding: 0.4rem 0.75rem;
  border-radius: 6px;
  transition: all 0.2s;
}

.launcher-link:hover {
  background-color: var(--color-background);
  color: var(--color-info);
  padding-left: 1rem;
}

.admin-section .section-icon {
  color: var(--color-primary);
}

.launcher-footer {
  padding: 1rem 2rem;
  background: var(--color-background);
  border-top: 1px solid var(--color-border);
  text-align: center;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}
</style>
