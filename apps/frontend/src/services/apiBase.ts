const DEFAULT_API_BASE = '/api'
const DEFAULT_AUTH_BASE = ''

export function getApiBase(): string {
  const raw = import.meta.env.VITE_API_URL || DEFAULT_API_BASE

  if (raw.startsWith('/')) {
    return raw
  }

  const normalized = raw.replace(/\/+$/, '')
  return normalized.endsWith('/api') ? normalized : `${normalized}/api`
}

export function getAuthBase(): string {
  const raw = import.meta.env.VITE_AUTH_URL || import.meta.env.VITE_API_URL || DEFAULT_AUTH_BASE

  if (!raw) {
    return DEFAULT_AUTH_BASE
  }

  if (raw.startsWith('/')) {
    const normalized = raw.replace(/\/+$/, '')
    return normalized.endsWith('/api') ? normalized.slice(0, -4) || '' : normalized
  }

  const normalized = raw.replace(/\/+$/, '')
  return normalized.endsWith('/api') ? normalized.slice(0, -4) : normalized
}
