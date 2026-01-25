import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { useAuthError } from '../useAuthError'

describe('useAuthError composable', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('setError() stores error message', () => {
    const { setError, error } = useAuthError()
    
    const testMessage = 'Invalid credentials'
    setError(testMessage)
    
    expect(error.value).toBe(testMessage)
  })

  it('clearError() removes error immediately', () => {
    const { setError, clearError, error } = useAuthError()
    
    setError('Test error')
    expect(error.value).toBe('Test error')
    
    clearError()
    
    expect(error.value).toBeNull()
  })

  it('setError() with duration auto-clears after timeout', () => {
    const { setError, error } = useAuthError()
    
    setError('Temporary error', 1000) // 1 second
    
    expect(error.value).toBe('Temporary error')
    
    // Fast-forward time
    vi.advanceTimersByTime(1000)
    
    expect(error.value).toBeNull()
  })

  it('setError() default duration is 5 seconds', () => {
    const { setError, error } = useAuthError()
    
    setError('Error with default duration')
    
    expect(error.value).toBe('Error with default duration')
    
    // Advance 4 seconds - should still be there
    vi.advanceTimersByTime(4000)
    expect(error.value).toBe('Error with default duration')
    
    // Advance 1 more second - should clear
    vi.advanceTimersByTime(1000)
    expect(error.value).toBeNull()
  })

  it('clearing error before timeout prevents auto-clear', () => {
    const { setError, clearError, error } = useAuthError()
    
    setError('Will be cleared manually', 5000)
    
    vi.advanceTimersByTime(2000) // After 2 seconds
    clearError()
    
    expect(error.value).toBeNull()
    
    // Even if we advance more, should stay null
    vi.advanceTimersByTime(5000)
    expect(error.value).toBeNull()
  })
})
