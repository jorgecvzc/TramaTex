import { ref } from 'vue'

export function useAuthError() {
  const error = ref<string | null>(null)
  const errorType = ref<string | null>(null)
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  function setError(message: string, durationMs: number = 5000) {
    error.value = message
    errorType.value = 'default'

    // Clear any existing timeout
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    // Auto-limpiar después del duration especificado
    timeoutId = setTimeout(() => {
      clearError()
    }, durationMs)
  }

  function clearError() {
    error.value = null
    errorType.value = null
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
  }

  return {
    error,
    errorType,
    setError,
    clearError
  }
}
