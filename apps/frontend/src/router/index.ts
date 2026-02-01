import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import LoginPage from '@/pages.deprecated/auth/LoginPage.vue'
import DashboardPage from '@/pages.deprecated/DashboardPage.vue'
import NotFoundPage from '@/pages.deprecated/NotFoundPage.vue'
import StyleGuide from '@/components/StyleGuide.vue'
import OrganizationsList from '@/pages/organizations/List.vue'
import OrganizationsCreate from '@/pages/organizations/Create.vue'
import OrganizationsDetail from '@/pages/organizations/Detail.vue'
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
    path: '/organizations',
    name: 'Organizations',
    component: OrganizationsList,
    meta: {
      requiresAuth: true,
      title: 'Organizations - TramaTex'
    }
  },
  {
    path: '/organizations/new',
    name: 'CreateOrganization',
    component: OrganizationsCreate,
    meta: {
      requiresAuth: true,
      title: 'Create Organization - TramaTex'
    }
  },
  {
    path: '/organizations/:id',
    name: 'OrganizationDetail',
    component: OrganizationsDetail,
    meta: {
      requiresAuth: true,
      title: 'Organization Details - TramaTex'
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
