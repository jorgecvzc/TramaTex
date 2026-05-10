import type { Router, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTokenManager } from '@/composables'

export async function requireAuth(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized
) {
  const authStore = useAuthStore()
  const { isTokenExpired } = useTokenManager()

  if (!authStore.isAuthenticated) {
    return { name: 'Login', query: { redirect: to.fullPath } }
  }

  // Check if token expired
  if (authStore.accessToken && isTokenExpired(authStore.accessToken)) {
    try {
      await authStore.refreshAccessToken()
    } catch (err) {
      return { name: 'Login' }
    }
  }
}

export async function requireGuest(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized
) {
  const authStore = useAuthStore()

  if (authStore.isAuthenticated) {
    return { name: 'Dashboard' }
  }
}

export function setupAuthGuards(router: Router) {
  router.beforeEach(async (to, from) => {
    const authStore = useAuthStore()

    // Restore session if available
    if (!authStore.isAuthenticated && !from.name) {
      await authStore.checkAuthStatus()
    }

    const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
    const requiresGuest = to.matched.some((record) => record.meta.requiresGuest)
    const requiresAdmin = to.matched.some((record) => record.meta.requiresAdmin)

    if (requiresAuth && !authStore.isAuthenticated) {
      return { name: 'Login', query: { redirect: to.fullPath } }
    }

    if (requiresAdmin && !authStore.isAdmin) {
      return { name: 'Dashboard' }
    }

    if (requiresGuest && authStore.isAuthenticated) {
      return { name: 'Dashboard' }
    }
  })
}
