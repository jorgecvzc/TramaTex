<template>
  <form @submit.prevent="handleSubmit">
    <!-- Email Input -->
    <div style="margin-bottom: 1rem;">
      <input
        id="email"
        v-model.trim="form.email"
        type="email"
        placeholder="Correo electrónico"
        style="
          display: block;
          width: 100%;
          padding: 0.75rem;
          border: 2px solid rgba(0, 35, 149, 0.3);
          border-radius: 0.375rem;
          background: white;
          color: #333;
          font-size: 0.875rem;
          box-sizing: border-box;
        "
        :disabled="isLoading"
        @blur="validateField('email')"
        required
      />
    </div>

    <!-- Password Input -->
    <div style="margin-bottom: 1rem;">
      <input
        id="password"
        v-model="form.password"
        type="password"
        placeholder="Contraseña"
        style="
          display: block;
          width: 100%;
          padding: 0.75rem;
          border: 2px solid rgba(0, 35, 149, 0.3);
          border-radius: 0.375rem;
          background: white;
          color: #333;
          font-size: 0.875rem;
          box-sizing: border-box;
        "
        :disabled="isLoading"
        @blur="validateField('password')"
        required
      />
    </div>

    <!-- Remember Me & Forgot Password -->
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; font-size: 0.75rem;">
      <label style="display: flex; align-items: center; gap: 0.5rem; cursor: pointer;">
        <input
          v-model="form.rememberMe"
          type="checkbox"
          style="cursor: pointer;"
          :disabled="isLoading"
        />
        Recuérdame
      </label>
      <a href="#" style="color: #002395; text-decoration: none; font-weight: 600;">¿Olvidaste tu contraseña?</a>
    </div>

    <!-- Error Alert -->
    <div v-if="generalError" style="margin-bottom: 1rem; padding: 0.75rem; background-color: #ffebee; border: 1px solid #ef5350; border-radius: 0.375rem; color: #c92a2a; font-size: 0.875rem;">
      {{ generalError }}
    </div>

    <!-- Login Button -->
    <button
      type="submit"
      style="
        width: 100%;
        background-color: #E6B800;
        color: white;
        font-weight: bold;
        padding: 0.875rem;
        border: none;
        border-radius: 0.375rem;
        cursor: pointer;
        font-size: 0.875rem;
        text-transform: uppercase;
      "
      :disabled="isLoading || !isFormValid"
    >
      {{ isLoading ? 'Cargando...' : 'Ingresar' }} →
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import type { LoginFormData } from '@/types/auth'

const router = useRouter()
const { login, isLoading } = useAuth()

const form = reactive<LoginFormData>({
  email: '',
  password: '',
  rememberMe: false
})

const errors = reactive<Record<string, string | null>>({
  email: null,
  password: null
})

const generalError = ref<string | null>(null)

const isFormValid = computed(() => {
  return (
    form.email.trim() !== '' &&
    form.password !== '' &&
    !errors.email &&
    !errors.password
  )
})

const focusStyle = (event: Event) => {
  const input = event.target as HTMLInputElement
  input.style.borderColor = '#002395'
  input.style.boxShadow = '0 0 0 4px rgba(0, 35, 149, 0.1)'
}

const validateField = (field: keyof typeof form) => {
  if (field === 'email') {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    errors.email = form.email && !emailRegex.test(form.email) ? 'Correo inválido' : null
  } else if (field === 'password') {
    errors.password = form.password && form.password.length < 6 ? 'Mínimo 6 caracteres' : null
  }
}

const handleSubmit = async () => {
  validateField('email')
  validateField('password')

  if (!isFormValid.value) return

  try {
    generalError.value = null
    await login(form.email, form.password)
    
    if (form.rememberMe) {
      localStorage.setItem('tramatex_remember_email', form.email)
    }

    router.push('/dashboard')
  } catch (error) {
    generalError.value = 'Credenciales inválidas. Intenta de nuevo.'
  }
}

// Load remembered email
const rememberedEmail = localStorage.getItem('tramatex_remember_email')
if (rememberedEmail) {
  form.email = rememberedEmail
  form.rememberMe = true
}
</script>

<style scoped>
.spinner {
  width: 16px;
  height: 16px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
