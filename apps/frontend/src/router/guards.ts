import type { Router, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTokenManager } from '@/composables'

export async function requireAuth(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: any
) {
  const authStore = useAuthStore()
  const { isTokenExpired } = useTokenManager()

  if (!authStore.isAuthenticated) {
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }

  // Verificar si el token expiró
  if (authStore.accessToken && isTokenExpired(authStore.accessToken)) {
    try {
      await authStore.refreshAccessToken()
    } catch (err) {
      return next({ name: 'Login' })
    }
  }

  next()
}

export async function requireGuest(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: any
) {
  const authStore = useAuthStore()

  if (authStore.isAuthenticated) {
    return next({ name: 'Dashboard' })
  }

  next()
}

export function setupAuthGuards(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // Restaurar sesión si está disponible
    if (!authStore.isAuthenticated && !from.name) {
      await authStore.checkAuthStatus()
    }

    if (import.meta.env.DEV && to.path.startsWith('/parties')) {
      return next()
    }

    const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
    const requiresGuest = to.matched.some((record) => record.meta.requiresGuest)
    const requiresAdmin = to.matched.some((record) => record.meta.requiresAdmin)

    if (requiresAuth && !authStore.isAuthenticated) {
      return next({ name: 'Login', query: { redirect: to.fullPath } })
    }

    if (requiresAdmin && !authStore.isAdmin) {
      return next({ name: 'Dashboard' })
    }

    if (requiresGuest && authStore.isAuthenticated) {
      return next({ name: 'Dashboard' })
    }

    next()
  })
}
