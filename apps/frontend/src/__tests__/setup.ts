import '@testing-library/jest-dom/vitest'
import { setActivePinia, createPinia } from 'pinia'
import { beforeEach, vi } from 'vitest'

// Global test setup for Vitest + Vue Test Utils + Testing Library
// This file is automatically loaded before each test file

beforeEach(() => {
  setActivePinia(createPinia())
})

// Mock global fetch
global.fetch = vi.fn()

// Mock axios globally
vi.mock('axios', () => {
  return {
    default: {
      create: vi.fn().mockReturnValue({
        interceptors: {
          request: { use: vi.fn(), eject: vi.fn() },
          response: { use: vi.fn(), eject: vi.fn() },
        },
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        patch: vi.fn(),
        delete: vi.fn(),
      }),
    },
  }
})

// Mock the base api service to use the same mock pattern as fetch
// This allows legacy tests using global.fetch to still work if we bridge them,
// but better to mock the api service directly.
vi.mock('@/services/api', () => {
  const mockAxiosInstance = {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn(), eject: vi.fn() },
      response: { use: vi.fn(), eject: vi.fn() },
    },
  }
  return {
    api: mockAxiosInstance
  }
})
