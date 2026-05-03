import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from "vue-router";
import Login from "@/pages/auth/Login.vue";
import Dashboard from "@/pages/Dashboard.vue";
import NotFound from "@/pages/NotFound.vue";
import PartiesList from "@/pages/parties/List.vue";
import PartiesCreate from "@/pages/parties/Create.vue";
import PartiesDetail from "@/pages/parties/Detail.vue";
import ProductsList from "@/pages/products/List.vue";
import UsersManagement from "@/pages/admin/UsersManagement.vue";
import PrintIssuerProfile from "@/pages/admin/PrintIssuerProfile.vue";
import { setupAuthGuards } from "./guards";

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    redirect: "/dashboard",
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: { requiresGuest: true, title: "Login - TramaTex" },
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: Dashboard,
    meta: { requiresAuth: true, title: "Dashboard - TramaTex" },
  },
  
  // --- MÓDULO: ENTIDADES (PARTY) ---
  {
    path: "/parties",
    name: "Parties",
    component: PartiesList,
    meta: { requiresAuth: true, title: "Listado de Entidades - TramaTex" },
  },
  {
    path: "/parties/dashboard",
    name: "PartyDashboard",
    component: () => import("@/pages/parties/PartyDashboard.vue"),
    meta: { requiresAuth: true, title: "Gestión de Entidades - TramaTex" },
  },
  {
    path: "/parties/new",
    name: "CreateParty",
    component: PartiesDetail,
    meta: { requiresAuth: true, title: "Crear Entidad - TramaTex" },
  },
  {
    path: "/parties/:id",
    name: "PartyDetail",
    component: PartiesDetail,
    meta: { requiresAuth: true, title: "Detalle de Entidad - TramaTex" },
  },

  // --- MÓDULO: PRODUCTOS (PRODUCT) ---
  {
    path: "/products",
    name: "Products",
    component: ProductsList,
    meta: { requiresAuth: true, title: "Catálogo de Productos - TramaTex" },
  },
  {
    path: "/products/dashboard",
    name: "ProductDashboard",
    component: () => import("@/pages/products/ProductDashboard.vue"),
    meta: { requiresAuth: true, title: "Gestión de Catálogo - TramaTex" },
  },
  {
    path: "/products/new",
    name: "CreateProduct",
    component: () => import("@/pages/products/Detail.vue"),
    meta: { requiresAuth: true, title: "Crear Producto - TramaTex" },
  },
  {
    path: "/products/:id",
    name: "ProductDetail",
    component: () => import("@/pages/products/Detail.vue"),
    meta: { requiresAuth: true, title: "Detalle de Producto - TramaTex" },
  },
  {
    path: "/products/pricing",
    name: "ProductPricing",
    component: () => import("@/pages/products/Pricing.vue"),
    meta: { requiresAuth: true, title: "Gestión de Precios - TramaTex" },
  },

  // --- MÓDULO: VENTAS (SALES) ---
  {
    path: "/sales",
    redirect: "/sales/dashboard",
  },
  {
    path: "/sales/dashboard",
    name: "SalesDashboard",
    component: () => import("@/pages/sales/SalesDashboard.vue"),
    meta: { requiresAuth: true, title: "Panel de Ventas - TramaTex" },
  },
  {
    path: "/sales/quotes",
    name: "QuoteList",
    component: () => import("@/pages/sales/QuoteList.vue"),
    meta: { requiresAuth: true, title: "Presupuestos - TramaTex" },
  },
  {
    path: "/sales/quotes/new",
    name: "CreateQuote",
    component: () => import("@/pages/sales/QuoteDetail.vue"),
    meta: { requiresAuth: true, title: "Nuevo Presupuesto - TramaTex" },
  },
  {
    path: "/sales/quotes/:id",
    name: "QuoteDetail",
    component: () => import("@/pages/sales/QuoteDetail.vue"),
    meta: { requiresAuth: true, title: "Detalle Presupuesto - TramaTex" },
  },
  {
    path: "/sales/orders",
    name: "OrderList",
    component: () => import("@/pages/sales/OrderList.vue"),
    meta: { requiresAuth: true, title: "Pedidos de Venta - TramaTex" },
  },
  {
    path: "/sales/orders/new",
    name: "CreateOrder",
    component: () => import("@/pages/sales/OrderDetail.vue"),
    meta: { requiresAuth: true, title: "Nuevo Pedido - TramaTex" },
  },
  {
    path: "/sales/orders/:id",
    name: "OrderDetail",
    component: () => import("@/pages/sales/OrderDetail.vue"),
    meta: { requiresAuth: true, title: "Detalle Pedido - TramaTex" },
  },
  {
    path: "/sales/invoices",
    name: "InvoiceList",
    component: () => import("@/pages/sales/InvoiceList.vue"),
    meta: { requiresAuth: true, title: "Facturas - TramaTex" },
  },
  {
    path: "/sales/invoices/:id",
    name: "InvoiceDetail",
    component: () => import("@/pages/sales/InvoiceDetail.vue"),
    meta: { requiresAuth: true, title: "Detalle Factura - TramaTex" },
  },
  {
    path: "/sales/delivery-notes",
    name: "DeliveryNoteList",
    component: () => import("@/pages/sales/DeliveryNoteList.vue"),
    meta: { requiresAuth: true, title: "Albaranes - TramaTex" },
  },
  {
    path: "/sales/delivery-notes/:id",
    name: "DeliveryNoteDetail",
    component: () => import("@/pages/sales/DeliveryNoteDetail.vue"),
    meta: { requiresAuth: true, title: "Detalle Albarán - TramaTex" },
  },
  {
    path: "/sales/tickets/new",
    name: "CreateTicket",
    component: () => import("@/pages/sales/TicketCreate.vue"),
    meta: { requiresAuth: true, title: "Nuevo Ticket - TramaTex" },
  },

  // --- DATOS MAESTROS (SECUNDARIOS) ---
  {
    path: "/master-data/brands",
    name: "BrandsList",
    component: () => import("@/pages/master-data/brands/List.vue"),
    meta: { requiresAuth: true, title: "Marcas - TramaTex" },
  },
  {
    path: "/master-data/product-groups",
    name: "ProductGroupsList",
    component: () => import("@/pages/master-data/product-groups/List.vue"),
    meta: { requiresAuth: true, title: "Categorías - TramaTex" },
  },
  {
    path: "/master-data/attributes",
    name: "AttributesList",
    component: () => import("@/pages/master-data/attributes/List.vue"),
    meta: { requiresAuth: true, title: "Atributos - TramaTex" },
  },

  // --- MÓDULO: PRODUCCIÓN (MES) ---
  {
    path: "/mes",
    redirect: "/mes/dashboard",
  },
  {
    path: "/mes/dashboard",
    name: "MESDashboard",
    component: () => import("@/pages/mes/Dashboard.vue"),
    meta: { requiresAuth: true, title: "Monitoreo de Producción - TramaTex" },
  },
  {
    path: "/mes/tasks",
    name: "MESTasksList",
    component: () => import("@/pages/mes/tasks/List.vue"),
    meta: { requiresAuth: true, title: "Tareas - TramaTex" },
  },
  {
    path: "/mes/positions",
    name: "MESPositionsList",
    component: () => import("@/pages/mes/positions/List.vue"),
    meta: { requiresAuth: true, title: "Posiciones - TramaTex" },
  },
  {
    path: "/mes/work-orders",
    name: "MESWorkOrdersList",
    component: () => import("@/pages/mes/works/List.vue"),
    meta: { requiresAuth: true, title: "Órdenes de Trabajo - TramaTex" },
  },
  {
    path: "/mes/work-orders/:id",
    name: "MESWorkOrderDetail",
    component: () => import("@/pages/mes/works/Detail.vue"),
    meta: { requiresAuth: true, title: "Detalle Orden de Trabajo - TramaTex" },
  },
  {
    path: "/mes/work-setups",
    name: "MESWorkSetupsList",
    component: () => import("@/pages/mes/work-setups/List.vue"),
    meta: { requiresAuth: true, title: "Configuraciones de Cliente - TramaTex" },
  },
  {
    path: "/mes/work-setups/new",
    name: "MESCreateWorkSetup",
    component: () => import("@/pages/mes/work-setups/Create.vue"),
    meta: { requiresAuth: true, title: "Nueva Configuración - TramaTex" },
  },
  {
    path: "/mes/work-setups/:id/edit",
    name: "MESEditWorkSetup",
    component: () => import("@/pages/mes/work-setups/Edit.vue"),
    meta: { requiresAuth: true, title: "Editar Configuración - TramaTex" },
  },
  {
    path: "/mes/terminal",
    name: "MESTabletTerminal",
    component: () => import("@/pages/mes/terminal/Tablet.vue"),
    meta: { requiresAuth: true, title: "Terminal de Taller - TramaTex" },
  },

  // --- ADMINISTRACIÓN Y SISTEMA ---
  {
    path: "/admin/users",
    name: "UsersManagement",
    component: UsersManagement,
    meta: { requiresAuth: true, requiresAdmin: true, title: "Usuarios - TramaTex" },
  },
  {
    path: "/admin/print-profile",
    name: "PrintIssuerProfile",
    component: PrintIssuerProfile,
    meta: { requiresAuth: true, requiresAdmin: true, title: "Perfil Fiscal Impresión - TramaTex" },
  },
  {
    path: "/dev/design-system",
    name: "DesignSystem",
    component: () => import("@/pages/dev/DesignSystem.vue"),
    meta: { title: "Sistema de Diseño - TramaTex" },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: NotFound,
    meta: { title: "Página no encontrada - TramaTex" },
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

setupAuthGuards(router);

router.afterEach((to) => {
  document.title = (to.meta.title as string) || "TramaTex";
});

export default router;
