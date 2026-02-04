import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import LoginPage from '@/pages.deprecated/auth/LoginPage.vue'
import DashboardPage from '@/pages.deprecated/DashboardPage.vue'
import NotFoundPage from '@/pages.deprecated/NotFoundPage.vue'
import StyleGuide from '@/components/StyleGuide.vue'
import PartiesList from '@/pages/parties/List.vue'
import PartiesCreate from '@/pages/parties/Create.vue'
import PartiesDetail from '@/pages/parties/Detail.vue'
import UsersManagement from '@/pages/admin/UsersManagement.vue'
import { setupAuthGuards } from './guards'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
    meta: {
      requiresGuest: true,
      title: 'Login - TramaTex'
    }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardPage,
    meta: {
      requiresAuth: true,
      title: 'Dashboard - TramaTex'
    }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: DashboardPage,
    meta: {
      requiresAuth: true,
      title: 'Pedidos - TramaTex'
    }
  },
  {
    path: '/inventory',
    name: 'Inventory',
    component: DashboardPage,
    meta: {
      requiresAuth: true,
      title: 'Inventario - TramaTex'
    }
  },
  {
    path: '/customers',
    name: 'Customers',
    component: DashboardPage,
    meta: {
      requiresAuth: true,
      title: 'Clientes - TramaTex'
    }
  },
  {
    path: '/parties',
    name: 'Parties',
    component: PartiesList,
    meta: {
      requiresAuth: true,
      title: 'Entidades - TramaTex'
    }
  },
  {
    path: '/parties/new',
    name: 'CreateParty',
    component: PartiesCreate,
    meta: {
      requiresAuth: true,
      title: 'Crear Entidad - TramaTex'
    }
  },
  {
    path: '/parties/:id',
    name: 'PartyDetail',
    component: PartiesDetail,
    meta: {
      requiresAuth: true,
      title: 'Detalle de Entidad - TramaTex'
    }
  },
  {
    path: '/admin/users',
    name: 'UsersManagement',
    component: UsersManagement,
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Usuarios - TramaTex'
    }
  },
  {
    path: '/catalogos',
    name: 'Catalogos',
    component: () => import('@/pages/CatalogosPage.vue'),
    meta: {
      requiresAuth: true,
      title: 'Catálogos - TramaTex'
    }
  },
  {
    path: '/style-guide',
    name: 'StyleGuide',
    component: StyleGuide,
    meta: {
      title: 'Guía de Estilos - TramaTex'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFoundPage,
    meta: {
      title: 'Página no encontrada - TramaTex'
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Configurar guards de autenticación
setupAuthGuards(router)

// Actualizar título de la página
router.afterEach((to) => {
  document.title = (to.meta.title as string) || 'TramaTex'
})

export default router
