import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../../stores/auth'

// Mock auth service
vi.mock('../../services/auth', () => ({
  authService: {
    login: vi.fn(async (email, password) => {
      if (email === 'valid@test.com' && password === 'ValidPass123!') {
        return {
          accessToken: 'mock-access-token',
          refreshToken: 'mock-refresh-token',
          usuario: { id: '1', email, created_at: new Date(), updated_at: new Date() },
          expiresIn: 3600
        }
      }
      throw new Error('Invalid credentials')
    }),
    logout: vi.fn(async () => {}),
    refreshToken: vi.fn(async () => ({
      accessToken: 'new-mock-token',
      refreshToken: 'new-refresh-token',
      expiresIn: 3600
    })),
    getCurrentUser: vi.fn(async () => ({ id: '1', email: 'valid@test.com', created_at: new Date(), updated_at: new Date() }))
  }
}))

describe('useAuthStore (Pinia)', () => {
  // Helper to generate a valid JWT token for testing
  const generateTestToken = (expiresIn: number = 3600): string => {
    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({ 
      sub: '1', 
      email: 'test@test.com',
      iat: Math.floor(Date.now() / 1000),
      exp: Math.floor(Date.now() / 1000) + expiresIn
    }))
    const signature = 'test-signature'
    return `${header}.${payload}.${signature}`
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('initial state should be unauthenticated', () => {
    const store = useAuthStore()
    
    expect(store.isAuthenticated).toBe(false)
    expect(store.user).toBeNull()
    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
  })

  it('login() action updates state with tokens and user', async () => {
    const store = useAuthStore()
    
    await store.login('valid@test.com', 'ValidPass123!')
    
    expect(store.isAuthenticated).toBe(true)
    expect(store.user).toBeDefined()
    expect(store.user?.email).toBe('valid@test.com')
    expect(store.accessToken).toBe('mock-access-token')
    expect(store.refreshToken).toBe('mock-refresh-token')
  })

  it('login() with invalid credentials should throw error', async () => {
    const store = useAuthStore()
    
    try {
      await store.login('invalid@test.com', 'wrongpassword')
      expect.fail('Should have thrown error')
    } catch (error) {
      expect(error).toBeDefined()
      expect(store.isAuthenticated).toBe(false)
    }
  })

  it('logout() action clears state completely', async () => {
    const store = useAuthStore()
    
    // Setup: logged in
    store.accessToken = 'some-token'
    store.refreshToken = 'some-refresh'
    store.user = { id: '1', email: 'test@test.com', created_at: new Date(), updated_at: new Date() }
    
    expect(store.isAuthenticated).toBe(true)
    
    // Logout
    await store.logout()
    
    // Verify cleared
    expect(store.isAuthenticated).toBe(false)
    expect(store.user).toBeNull()
    expect(store.accessToken).toBeNull()
    expect(store.refreshToken).toBeNull()
  })

  it('refreshAccessToken() updates access token', async () => {
    const store = useAuthStore()
    
    // Setup
    store.accessToken = 'old-token'
    store.refreshToken = 'refresh-token'
    store.user = { id: '1', email: 'test@test.com', created_at: new Date(), updated_at: new Date() }
    
    // Refresh
    await store.refreshAccessToken()
    
    // Verify updated
    expect(store.accessToken).toBe('new-mock-token')
    expect(store.refreshToken).toBe('new-refresh-token')
    expect(store.isAuthenticated).toBe(true)
  })

  it('checkAuthStatus() restores session from localStorage', () => {
    // Setup: localStorage with valid session
    const userData = { id: '1', email: 'test@test.com', created_at: new Date(), updated_at: new Date() }
    const token = generateTestToken()
    localStorage.setItem('tramatex_auth_token', token)
    localStorage.setItem('tramatex_refresh_token', 'stored-refresh')
    localStorage.setItem('tramatex_user', JSON.stringify(userData))
    
    const store = useAuthStore()
    store.checkAuthStatus()
    
    // Verify restored
    expect(store.isAuthenticated).toBe(true)
    expect(store.accessToken).toBe(token)
    expect(store.refreshToken).toBe('stored-refresh')
    expect(store.user?.email).toBe('test@test.com')
  })

  it('checkAuthStatus() clears if tokens expired', () => {
    // Setup: localStorage with expired token marker
    const expiredToken = generateTestToken(-1) // -1 second = expired
    localStorage.setItem('tramatex_auth_token', expiredToken)
    localStorage.setItem('tokenExpired', 'true')
    
    const store = useAuthStore()
    store.checkAuthStatus()
    
    // Verify cleared (token was expired)
    expect(store.isAuthenticated).toBe(false)
    expect(store.accessToken).toBeNull()
    expect(store.user).toBeNull()
  })

  it('isAuthenticated getter reflects current state', () => {
    const store = useAuthStore()
    
    expect(store.isAuthenticated).toBe(false)
    
    // Both token and user must be set for authenticated state
    store.accessToken = 'token'
    store.user = { id: '1', email: 'test@example.com', created_at: new Date(), updated_at: new Date() }
    expect(store.isAuthenticated).toBe(true)
    
    // Clear either one
    store.accessToken = null
    expect(store.isAuthenticated).toBe(false)
    
    // Reset and test the other way
    store.accessToken = 'token'
    store.user = null
    expect(store.isAuthenticated).toBe(false)
  })

  it('currentUser getter returns user or null', () => {
    const store = useAuthStore()
    
    expect(store.currentUser).toBeNull()
    
    const testUser = { id: '1', email: 'test@test.com', created_at: new Date(), updated_at: new Date() }
    store.user = testUser
    
    expect(store.currentUser).toEqual(testUser)
    expect(store.currentUser?.email).toBe('test@test.com')
  })

  it('state persists to localStorage on login', async () => {
    const store = useAuthStore()
    
    await store.login('valid@test.com', 'ValidPass123!')
    
    // Verify localStorage
    expect(localStorage.getItem('tramatex_auth_token')).toBe('mock-access-token')
    expect(localStorage.getItem('tramatex_refresh_token')).toBe('mock-refresh-token')
    expect(JSON.parse(localStorage.getItem('tramatex_user') || '{}')).toHaveProperty('email', 'valid@test.com')
  })

  it('state clears from localStorage on logout', async () => {
    const store = useAuthStore()
    
    // Setup: logged in
    await store.login('valid@test.com', 'ValidPass123!')
    expect(localStorage.getItem('tramatex_auth_token')).toBeTruthy()
    
    // Logout
    await store.logout()
    
    // Verify localStorage cleared
    expect(localStorage.getItem('tramatex_auth_token')).toBeNull()
    expect(localStorage.getItem('tramatex_refresh_token')).toBeNull()
    expect(localStorage.getItem('tramatex_user')).toBeNull()
  })
})
