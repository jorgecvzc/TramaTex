import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from "vue-router";
import Login from "@/pages/auth/Login.vue";
import Dashboard from "@/pages/Dashboard.vue";
import NotFound from "@/pages/NotFound.vue";
import StyleGuide from "@/components/StyleGuide.vue";
import PartiesList from "@/pages/parties/List.vue";
import PartiesCreate from "@/pages/parties/Create.vue";
import PartiesDetail from "@/pages/parties/Detail.vue";
import ProductsList from "@/pages/products/List.vue";
import UsersManagement from "@/pages/admin/UsersManagement.vue";
import PrintIssuerProfile from "@/pages/admin/PrintIssuerProfile.vue";
import { setupAuthGuards } from "./guards";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: {
      requiresGuest: true,
      title: "Login - TramaTex",
    },
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: Dashboard,
    meta: {
      requiresAuth: true,
      title: "Dashboard - TramaTex",
    },
  },
  {
    path: "/orders",
    name: "Orders",
    component: Dashboard,
    meta: {
      requiresAuth: true,
      title: "Pedidos - TramaTex",
    },
  },
  {
    path: "/inventory",
    name: "Inventory",
    component: Dashboard,
    meta: {
      requiresAuth: true,
      title: "Inventario - TramaTex",
    },
  },
  {
    path: "/customers",
    name: "Customers",
    component: Dashboard,
    meta: {
      requiresAuth: true,
      title: "Clientes - TramaTex",
    },
  },
  {
    path: "/parties",
    name: "Parties",
    component: PartiesList,
    meta: {
      requiresAuth: true,
      title: "Entidades - TramaTex",
    },
  },
  {
    path: "/parties/new",
    name: "CreateParty",
    component: PartiesCreate,
    meta: {
      requiresAuth: true,
      title: "Crear Entidad - TramaTex",
    },
  },
  {
    path: "/parties/:id",
    name: "PartyDetail",
    component: PartiesDetail,
    meta: {
      requiresAuth: true,
      title: "Detalle de Entidad - TramaTex",
    },
  },
  {
    path: "/products",
    name: "Products",
    component: ProductsList,
    meta: {
      requiresAuth: true,
      title: "Catálogo de Productos - TramaTex",
    },
  },
  {
    path: "/products/new",
    name: "CreateProduct",
    component: () => import("@/pages/products/Create.vue"),
    meta: {
      requiresAuth: true,
      title: "Crear Producto - TramaTex",
    },
  },
  {
    path: "/products/:id",
    name: "ProductDetail",
    component: () => import("@/pages/products/Detail.vue"),
    meta: {
      requiresAuth: true,
      title: "Detalle de Producto - TramaTex",
    },
  },
  {
    path: "/sales/quotes",
    name: "QuoteList",
    component: () => import("@/pages/sales/QuoteList.vue"),
    meta: {
      requiresAuth: true,
      title: "Presupuestos - TramaTex",
    },
  },
  {
    path: "/sales/quotes/new",
    name: "CreateQuote",
    component: () => import("@/pages/sales/QuoteCreate.vue"),
    meta: {
      requiresAuth: true,
      title: "Nuevo Presupuesto - TramaTex",
    },
  },
  {
    path: "/sales/quotes/:id",
    name: "QuoteDetail",
    component: () => import("@/pages/sales/QuoteDetail.vue"),
    meta: {
      requiresAuth: true,
      title: "Detalle Presupuesto - TramaTex",
    },
  },
  {
    path: "/sales/orders",
    name: "OrderList",
    component: () => import("@/pages/sales/OrderList.vue"),
    meta: {
      requiresAuth: true,
      title: "Pedidos - TramaTex",
    },
  },
  {
    path: "/sales/orders/new",
    name: "CreateOrder",
    component: () => import("@/pages/sales/OrderCreate.vue"),
    meta: {
      requiresAuth: true,
      title: "Nuevo Pedido - TramaTex",
    },
  },
  {
    path: "/sales/orders/:id",
    name: "OrderDetail",
    component: () => import("@/pages/sales/OrderDetail.vue"),
    meta: {
      requiresAuth: true,
      title: "Detalle Pedido - TramaTex",
    },
  },
  {
    path: "/sales/invoices",
    name: "InvoiceList",
    component: () => import("@/pages/sales/InvoiceList.vue"),
    meta: {
      requiresAuth: true,
      title: "Facturas - TramaTex",
    },
  },
  {
    path: "/sales/invoices/:id",
    name: "InvoiceDetail",
    component: () => import("@/pages/sales/InvoiceDetail.vue"),
    meta: {
      requiresAuth: true,
      title: "Detalle Factura - TramaTex",
    },
  },
  {
    path: "/sales/delivery-notes",
    name: "DeliveryNoteList",
    component: () => import("@/pages/sales/DeliveryNoteList.vue"),
    meta: {
      requiresAuth: true,
      title: "Albaranes - TramaTex",
    },
  },
  {
    path: "/sales/delivery-notes/:id",
    name: "DeliveryNoteDetail",
    component: () => import("@/pages/sales/DeliveryNoteDetail.vue"),
    meta: {
      requiresAuth: true,
      title: "Detalle Albarán - TramaTex",
    },
  },
  {
    path: "/sales/tickets/new",
    name: "CreateTicket",
    component: () => import("@/pages/sales/TicketCreate.vue"),
    meta: {
      requiresAuth: true,
      title: "Nuevo Ticket - TramaTex",
    },
  },
  {
    path: "/master-data/brands",
    name: "BrandsList",
    component: () => import("@/pages/master-data/brands/List.vue"),
    meta: {
      requiresAuth: true,
      title: "Marcas - TramaTex",
    },
  },
  {
    path: "/master-data/product-groups",
    name: "ProductGroupsList",
    component: () => import("@/pages/master-data/product-groups/List.vue"),
    meta: {
      requiresAuth: true,
      title: "Categorías - TramaTex",
    },
  },
  {
    path: "/master-data/attributes",
    name: "AttributesList",
    component: () => import("@/pages/master-data/attributes/List.vue"),
    meta: {
      requiresAuth: true,
      title: "Atributos - TramaTex",
    },
  },
  {
    path: "/mes/dashboard",
    name: "MESDashboard",
    component: () => import("@/pages/mes/Dashboard.vue"),
    meta: {
      requiresAuth: true,
      title: "Dashboard MES - TramaTex",
    },
  },
  {
    path: "/mes/tasks",
    name: "MESTasksList",
    component: () => import("@/pages/mes/tasks/List.vue"),
    meta: {
      requiresAuth: true,
      title: "Tareas MES - TramaTex",
    },
  },
  {
    path: "/mes/tasks/new",
    name: "MESCreateTask",
    component: () => import("@/pages/mes/tasks/Create.vue"),
    meta: {
      requiresAuth: true,
      title: "Nueva Tarea MES - TramaTex",
    },
  },
  {
    path: "/mes/positions",
    name: "MESPositionsList",
    component: () => import("@/pages/mes/positions/List.vue"),
    meta: {
      requiresAuth: true,
      title: "Puestos MES - TramaTex",
    },
  },
  {
    path: "/mes/positions/new",
    name: "MESCreatePosition",
    component: () => import("@/pages/mes/positions/Create.vue"),
    meta: {
      requiresAuth: true,
      title: "Nuevo Puesto MES - TramaTex",
    },
  },
  {
    path: "/mes/service-groups",
    name: "MESServiceGroupsList",
    component: () => import("@/pages/mes/service-groups/List.vue"),
    meta: {
      requiresAuth: true,
      title: "Grupos de Servicio MES - TramaTex",
    },
  },
  {
    path: "/mes/service-groups/new",
    name: "MESCreateServiceGroup",
    component: () => import("@/pages/mes/service-groups/Create.vue"),
    meta: {
      requiresAuth: true,
      title: "Nuevo Grupo de Servicio MES - TramaTex",
    },
  },
  {
    path: "/mes/works",
    name: "MESWorksList",
    component: () => import("@/pages/mes/works/List.vue"),
    meta: {
      requiresAuth: true,
      title: "Trabajos MES - TramaTex",
    },
  },
  {
    path: "/mes/works/new",
    name: "MESCreateWork",
    component: () => import("@/pages/mes/works/Create.vue"),
    meta: {
      requiresAuth: true,
      title: "Nuevo Trabajo MES - TramaTex",
    },
  },
  {
    path: "/mes/works/:id",
    name: "MESWorkDetail",
    component: () => import("@/pages/mes/works/Detail.vue"),
    meta: {
      requiresAuth: true,
      title: "Detalle Trabajo MES - TramaTex",
    },
  },
  {
    path: "/mes/terminal",
    name: "MESTabletTerminal",
    component: () => import("@/pages/mes/terminal/Tablet.vue"),
    meta: {
      requiresAuth: true,
      title: "Terminal MES Tablet - TramaTex",
    },
  },
  {
    path: "/admin/users",
    name: "UsersManagement",
    component: UsersManagement,
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: "Usuarios - TramaTex",
    },
  },
  {
    path: "/admin/print-profile",
    name: "PrintIssuerProfile",
    component: PrintIssuerProfile,
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: "Perfil Fiscal Impresión - TramaTex",
    },
  },
  {
    path: "/catalogos",
    name: "Catalogos",
    component: () => import("@/pages/CatalogosPage.vue"),
    meta: {
      requiresAuth: true,
      title: "Catálogos - TramaTex",
    },
  },
  {
    path: "/style-guide",
    name: "StyleGuide",
    component: StyleGuide,
    meta: {
      title: "Guía de Estilos - TramaTex",
    },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: NotFound,
    meta: {
      title: "Página no encontrada - TramaTex",
    },
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

// Configurar guards de autenticación
setupAuthGuards(router);

// Actualizar título de la página
router.afterEach((to) => {
  document.title = (to.meta.title as string) || "TramaTex";
});

export default router;
