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

    const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
    const requiresGuest = to.matched.some((record) => record.meta.requiresGuest)

    if (requiresAuth && !authStore.isAuthenticated) {
      return next({ name: 'Login', query: { redirect: to.fullPath } })
    }

    if (requiresGuest && authStore.isAuthenticated) {
      return next({ name: 'Dashboard' })
    }

    next()
  })
}
