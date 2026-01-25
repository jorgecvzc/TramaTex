import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuth } from '../useAuth'
import { useAuthStore } from '../../stores/auth'
import { setActivePinia, createPinia } from 'pinia'

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

describe('useAuth composable', () => {
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
    vi.clearAllMocks()
  })

  it('login() with valid credentials should succeed', async () => {
    const { login } = useAuth()
    const result = await login('valid@test.com', 'ValidPass123!')
    
    expect(result).toBeDefined()
    expect(result.user).toBeDefined()
    expect(result.user.email).toBe('valid@test.com')
  })

  it('login() with invalid email should throw error', async () => {
    const { login } = useAuth()
    
    try {
      await login('invalid@test.com', 'ValidPass123!')
      expect.fail('Should have thrown error')
    } catch (error) {
      expect(error).toBeDefined()
      expect(error.message).toContain('Invalid')
    }
  })

  it('login() with invalid password should throw error', async () => {
    const { login } = useAuth()
    
    try {
      await login('valid@test.com', 'wrongpass')
      expect.fail('Should have thrown error')
    } catch (error) {
      expect(error).toBeDefined()
    }
  })

  it('logout() should clear authentication', async () => {
    const store = useAuthStore()
    const { logout } = useAuth()
    
    // Setup: manually set authenticated state
    store.accessToken = 'some-token'
    store.user = { id: '1', email: 'test@test.com', created_at: new Date(), updated_at: new Date() }
    
    expect(store.isAuthenticated).toBe(true)
    
    // Call logout
    await logout()
    
    // Verify cleared
    expect(store.isAuthenticated).toBe(false)
    expect(store.user).toBeNull()
    expect(store.accessToken).toBeNull()
  })

  it('checkAuth() restores session from valid tokens', async () => {
    const store = useAuthStore()
    const { checkAuth } = useAuth()
    
    // Setup: localStorage with valid tokens
    const token = generateTestToken()
    localStorage.setItem('tramatex_auth_token', token)
    localStorage.setItem('tramatex_refresh_token', 'valid-refresh-token')
    localStorage.setItem('tramatex_user', JSON.stringify({
      id: '1',
      email: 'test@test.com',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }))
    
    // Call checkAuth
    checkAuth()
    
    // Verify restored
    expect(store.isAuthenticated).toBe(true)
    expect(store.user).toBeDefined()
  })

  it('checkAuth() clears expired tokens', () => {
    const store = useAuthStore()
    const { checkAuth } = useAuth()
    
    // Setup: localStorage with expired token
    const expiredToken = generateTestToken(-1) // -1 second = expired
    localStorage.setItem('tramatex_auth_token', expiredToken)
    
    // Call checkAuth
    checkAuth()
    
    // Verify cleared
    expect(store.accessToken).toBeNull()
  })
})
