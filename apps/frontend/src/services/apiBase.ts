const DEFAULT_API_BASE = '/api'

export function getApiBase(): string {
  const raw = import.meta.env.VITE_API_URL || DEFAULT_API_BASE

  if (raw.startsWith('/')) {
    return raw
  }

  const normalized = raw.replace(/\/+$/, '')
  return normalized.endsWith('/api') ? normalized : `${normalized}/api`
}
