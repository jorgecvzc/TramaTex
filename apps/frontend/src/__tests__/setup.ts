import '@testing-library/jest-dom/vitest'

// Global test setup for Vitest + Vue Test Utils + Testing Library
// This file is automatically loaded before each test file

// Mock global fetch if needed (already done in individual test files)
global.fetch = global.fetch || vi.fn()
