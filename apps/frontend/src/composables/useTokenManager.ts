import { ref } from 'vue'

export function useTokenManager() {
  const refreshTimer = ref<NodeJS.Timeout | null>(null)

  function decodeJWT(token: string) {
    try {
      const payload = token.split('.')[1]
      const decoded = JSON.parse(atob(payload))
      return decoded
    } catch (err) {
      return null
    }
  }

  function isTokenExpired(token: string): boolean {
    const decoded = decodeJWT(token)
    if (!decoded || !decoded.exp) return true
    return Date.now() >= decoded.exp * 1000
  }

  function getRemainingTime(token: string): number {
    const decoded = decodeJWT(token)
    if (!decoded || !decoded.exp) return 0
    return Math.max(0, decoded.exp * 1000 - Date.now())
  }

  function isRefreshNeeded(token: string): boolean {
    const remaining = getRemainingTime(token)
    const tenMinutesMs = 10 * 60 * 1000
    return remaining > 0 && remaining < tenMinutesMs
  }

  function startRefreshTimer(callback: () => void) {
    if (refreshTimer.value) clearInterval(refreshTimer.value)

    refreshTimer.value = setInterval(() => {
      callback()
    }, 60000) // Cada 60 segundos
  }

  function clearTimer() {
    if (refreshTimer.value) {
      clearInterval(refreshTimer.value)
      refreshTimer.value = null
    }
  }

  return {
    isTokenExpired,
    getRemainingTime,
    isRefreshNeeded,
    startRefreshTimer,
    clearTimer,
    decodeJWT
  }
}
