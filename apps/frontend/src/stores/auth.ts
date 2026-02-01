import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Usuario, AuthState } from '@/types/auth'
import { authService } from '@/services/auth'

/**
 * TramaTex Authentication Store (Pinia)
 * 
 * Manages:
 * - User authentication state
 * - JWT tokens (access + refresh)
 * - Token expiration and renewal
 * - Local storage persistence
 * - Authentication errors
 */
export const useAuthStore = defineStore('auth', () => {
  // =========================================================================
  // STATE
  // =========================================================================

  const user = ref<Usuario | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const isRefreshing = ref(false)
  const error = ref<string | null>(null)
  const tokenExpiresAt = ref<number | null>(null)
  let tokenRefreshTimer: NodeJS.Timeout | null = null

  // =========================================================================
  // COMPUTED PROPERTIES (GETTERS)
  // =========================================================================

  /**
   * Checks if user is authenticated (has both user and valid token)
   */
  const isAuthenticated = computed(() => {
    return user.value !== null && accessToken.value !== null
  })

  /**
   * Returns current authenticated user or null
   */
  const currentUser = computed(() => user.value)

  /**
   * Returns token expiration date or null
   */
  const tokenExpiry = computed(() => {
    if (!tokenExpiresAt.value) return null
    return new Date(tokenExpiresAt.value)
  })

  /**
   * Checks if token is about to expire (within 5 minutes)
   */
  const isTokenExpiringSoon = computed(() => {
    if (!tokenExpiresAt.value) return false
    const now = Date.now()
    const fiveMinutesMs = 5 * 60 * 1000
    return tokenExpiresAt.value - now < fiveMinutesMs
  })

  /**
   * Returns user role or null
   */
  const userRole = computed(() => user.value?.role || null)

  /**
   * Check if user has admin role
   */
  const isAdmin = computed(() => user.value?.role === 'admin')

  // =========================================================================
  // ACTIONS
  // =========================================================================

  /**
   * Login user with email and password
   * 
   * Flow:
   * 1. Call auth service with credentials
   * 2. Store tokens and user info
   * 3. Persist to localStorage
   * 4. Setup token refresh timer
   * 
   * @param email - User email
   * @param password - User password
   * @throws Error if login fails
   */
  async function login(email: string, password: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await authService.login(email, password)

      // Store response data
      accessToken.value = response.accessToken
      refreshToken.value = response.refreshToken
      user.value = response.user
      tokenExpiresAt.value = Date.now() + response.expiresIn * 1000

      // Persist to localStorage for session recovery
      persistAuthState()

      // Setup automatic token refresh
      setupTokenRefreshTimer()

      return response
    } catch (err: any) {
      // Extract error message from backend response
      error.value = err.response?.data?.message || err.message || 'Error al iniciar sesión'
      
      // Clear auth state on failure
      clearAuthState()

      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Logout user and clean up state
   * 
   * - Calls logout endpoint (if exists)
   * - Clears all state
   * - Removes localStorage data
   * - Cancels token refresh timer
   */
  async function logout() {
    try {
      // Call backend logout endpoint (for audit/cleanup)
      await authService.logout()
    } catch (err) {
      // Logout on client side even if backend fails
      console.error('[AuthStore] Logout error:', err)
    } finally {
      // Clear all auth state
      clearAuthState()
    }
  }

  /**
   * Refresh access token using refresh token
   * 
   * Called when:
   * - Token is about to expire
   * - 401 response received (token expired)
   * 
   * @throws Error if refresh fails
   */
  async function refreshAccessToken() {
    if (!refreshToken.value) {
      await logout()
      return
    }

    // Prevent multiple simultaneous refresh attempts
    if (isRefreshing.value) {
      return
    }

    isRefreshing.value = true

    try {
      const response = await authService.refreshToken(refreshToken.value)

      // Update tokens
      accessToken.value = response.accessToken
      // Keep existing refresh token (backend does not return a new one)
      refreshToken.value = refreshToken.value
      tokenExpiresAt.value = Date.now() + response.expiresIn * 1000

      // Persist updated tokens
      persistAuthState()

      // Restart refresh timer
      setupTokenRefreshTimer()

      return response
    } catch (err) {
      // If refresh fails, logout user
      console.error('[AuthStore] Token refresh failed:', err)
      await logout()
      throw err
    } finally {
      isRefreshing.value = false
    }
  }

  /**
   * Check authentication status on app startup
   * 
   * - Restore session from localStorage
   * - Verify token validity
   * - Refresh if needed
   */
  async function checkAuthStatus() {
    const storedToken = localStorage.getItem('tramatex_auth_token')
    const storedUser = localStorage.getItem('tramatex_user')
    const storedRefreshToken = localStorage.getItem('tramatex_refresh_token')

    if (!storedToken || !storedUser) {
      clearAuthState()
      return
    }

    try {
      // Restore from localStorage
      accessToken.value = storedToken
      user.value = JSON.parse(storedUser)
      refreshToken.value = storedRefreshToken

      // Decode and verify token expiration
      const isExpired = isTokenExpired(storedToken)

      if (isExpired) {
        // Token expired, try to refresh
        await refreshAccessToken()
      } else {
        // Token still valid, restore expiry time and setup timer
        const expiresAt = decodeTokenExpiry(storedToken)
        if (expiresAt) {
          tokenExpiresAt.value = expiresAt
          setupTokenRefreshTimer()
        }
      }
    } catch (err) {
      // Any error during session restore → logout
      console.error('[AuthStore] Session restore failed:', err)
      clearAuthState()
    }
  }

  /**
   * Clear all authentication state
   */
  function clearAuthState(): void {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    tokenExpiresAt.value = null
    error.value = null
    clearTokenRefreshTimer()

    // Remove from localStorage
    localStorage.removeItem('tramatex_auth_token')
    localStorage.removeItem('tramatex_refresh_token')
    localStorage.removeItem('tramatex_user')
  }

  // =========================================================================
  // PRIVATE HELPER FUNCTIONS
  // =========================================================================

  /**
   * Persist current auth state to localStorage
   */
  function persistAuthState(): void {
    if (accessToken.value) {
      localStorage.setItem('tramatex_auth_token', accessToken.value)
    }
    if (refreshToken.value) {
      localStorage.setItem('tramatex_refresh_token', refreshToken.value)
    }
    if (user.value) {
      localStorage.setItem('tramatex_user', JSON.stringify(user.value))
    }
  }

  /**
   * Check if JWT token is expired
   * 
   * Decodes token payload and checks exp claim
   */
  function isTokenExpired(token: string): boolean {
    try {
      const parts = token.split('.')
      if (parts.length !== 3) return true

      const payload = JSON.parse(atob(parts[1]))
      const expiresAt = payload.exp * 1000 // Convert to milliseconds

      return Date.now() > expiresAt
    } catch (err) {
      console.error('[AuthStore] Token decode error:', err)
      return true
    }
  }

  /**
   * Decode token expiration time (ms)
   */
  function decodeTokenExpiry(token: string): number | null {
    try {
      const parts = token.split('.')
      if (parts.length !== 3) return null

      const payload = JSON.parse(atob(parts[1]))
      return payload.exp * 1000 // Convert to milliseconds
    } catch (err) {
      console.error('[AuthStore] Token expiry decode error:', err)
      return null
    }
  }

  /**
   * Setup automatic token refresh
   * 
   * Refreshes token 2 minutes before expiration
   */
  function setupTokenRefreshTimer(): void {
    clearTokenRefreshTimer()

    if (!tokenExpiresAt.value) return

    const now = Date.now()
    const expiresAt = tokenExpiresAt.value
    const twoMinutesMs = 2 * 60 * 1000
    const timeUntilRefresh = expiresAt - now - twoMinutesMs

    if (timeUntilRefresh > 0) {
      tokenRefreshTimer = setTimeout(() => {
        refreshAccessToken().catch(err => {
          console.error('[AuthStore] Automatic token refresh failed:', err)
        })
      }, timeUntilRefresh)
    }
  }

  /**
   * Clear token refresh timer
   */
  function clearTokenRefreshTimer(): void {
    if (tokenRefreshTimer) {
      clearTimeout(tokenRefreshTimer)
      tokenRefreshTimer = null
    }
  }

  // =========================================================================
  // PUBLIC API
  // =========================================================================

  return {
    // State
    user,
    accessToken,
    refreshToken,
    isLoading,
    isRefreshing,
    error,
    tokenExpiresAt,

    // Computed
    isAuthenticated,
    currentUser,
    tokenExpiry,
    isTokenExpiringSoon,
    userRole,
    isAdmin,

    // Actions
    login,
    logout,
    refreshAccessToken,
    checkAuthStatus,
    clearAuthState
  }
})
