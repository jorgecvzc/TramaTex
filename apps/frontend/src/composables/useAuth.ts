import { useAuthStore } from '@/stores/auth'
import { computed } from 'vue'

export function useAuth() {
  const authStore = useAuthStore()

  const usuario = computed(() => authStore.currentUser)
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const isLoading = computed(() => authStore.isLoading)

  async function login(email: string, password: string) {
    // Validación básica
    if (!email || !password) {
      throw new Error('Email y contraseña son requeridos')
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(email)) {
      throw new Error('Email inválido')
    }

    await authStore.login(email, password)
    return {
      user: authStore.currentUser,
      accessToken: authStore.accessToken
    }
  }

  async function logout() {
    return await authStore.logout()
  }

  async function checkAuth() {
    return await authStore.checkAuthStatus()
  }

  return {
    usuario,
    isAuthenticated,
    isLoading,
    login,
    logout,
    checkAuth
  }
}
