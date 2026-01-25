import { describe, it, expect } from 'vitest'
import { useTokenManager } from '../useTokenManager'

describe('useTokenManager composable', () => {
  
  it('decodeJWT() extracts claims from valid token', () => {
    const { decodeJWT } = useTokenManager()
    
    // Mock JWT token (header.payload.signature)
    // Payload: { sub: '1', email: 'test@test.com', iat: 1234567890, exp: 9999999999 }
    const mockToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiaWF0IjoxMjM0NTY3ODkwLCJleHAiOjk5OTk5OTk5OTl9.fake-signature'
    
    const claims = decodeJWT(mockToken)
    
    expect(claims).toBeDefined()
    expect(claims.sub).toBe('1')
    expect(claims.email).toBe('test@test.com')
  })

  it('isTokenExpired() returns true for expired token', () => {
    const { isTokenExpired, decodeJWT } = useTokenManager()
    
    // Create token with expiry in the past (1 hour ago)
    const pastTimestamp = Math.floor(Date.now() / 1000) - 3600
    const expiredToken = `header.${btoa(JSON.stringify({ exp: pastTimestamp }))}.signature`
    
    const result = isTokenExpired(expiredToken)
    
    expect(result).toBe(true)
  })

  it('isTokenExpired() returns false for valid token', () => {
    const { isTokenExpired } = useTokenManager()
    
    // Create token with expiry in the future (1 hour from now)
    const futureTimestamp = Math.floor(Date.now() / 1000) + 3600
    const validToken = `header.${btoa(JSON.stringify({ exp: futureTimestamp }))}.signature`
    
    const result = isTokenExpired(validToken)
    
    expect(result).toBe(false)
  })

  it('isRefreshNeeded() returns true for tokens expiring soon', () => {
    const { isRefreshNeeded } = useTokenManager()
    
    // Token expiring in 5 minutes (threshold: refresh if <10 min remaining)
    const soonTimestamp = Math.floor(Date.now() / 1000) + 300 // 5 minutes
    const soonToken = `header.${btoa(JSON.stringify({ exp: soonTimestamp }))}.signature`
    
    const result = isRefreshNeeded(soonToken)
    
    expect(result).toBe(true)
  })

  it('isRefreshNeeded() returns false for fresh tokens', () => {
    const { isRefreshNeeded } = useTokenManager()
    
    // Token expiring in 30 minutes (well above 10 min threshold)
    const freshTimestamp = Math.floor(Date.now() / 1000) + 1800 // 30 minutes
    const freshToken = `header.${btoa(JSON.stringify({ exp: freshTimestamp }))}.signature`
    
    const result = isRefreshNeeded(freshToken)
    
    expect(result).toBe(false)
  })
})
