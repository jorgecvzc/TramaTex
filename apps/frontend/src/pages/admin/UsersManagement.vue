<template>
  <div class="min-h-screen bg-slate-950/5 p-6">
    <div class="mx-auto max-w-6xl space-y-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 class="text-3xl font-semibold text-slate-900">Gestión de usuarios</h1>
          <p class="text-sm text-slate-600">
            Administración de usuarios, roles y accesos.
          </p>
        </div>
        <div class="flex items-center gap-2">
          <button
            class="rounded-full border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm hover:border-slate-300 hover:bg-slate-50"
            @click="loadUsers"
            :disabled="isLoading"
          >
            {{ isLoading ? 'Cargando...' : 'Refrescar' }}
          </button>
        </div>
      </div>

      <div v-if="!isAdmin" class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-amber-800">
        Solo el rol admin puede gestionar usuarios.
      </div>

      <div v-else class="space-y-6">
        <div class="grid gap-4 lg:grid-cols-[1.1fr_1fr]">
          <div class="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
            <div class="mb-4 flex items-center justify-between">
              <div>
                <h2 class="text-lg font-semibold text-slate-900">Alta de usuario</h2>
                <p class="text-xs text-slate-500">Crea cuentas con roles predefinidos.</p>
              </div>
              <span class="rounded-full bg-slate-900 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-white">
                Admin
              </span>
            </div>

            <form class="grid gap-3" @submit.prevent="createUser">
              <div>
                <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-slate-500">
                  Email
                </label>
                <input
                  v-model.trim="newUser.email"
                  type="email"
                  placeholder="usuario@tramatex.local"
                  class="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 focus:border-slate-400 focus:outline-none"
                  required
                />
              </div>
              <div>
                <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-slate-500">
                  Contraseña
                </label>
                <input
                  v-model="newUser.password"
                  type="password"
                  placeholder="Mínimo 8 caracteres"
                  class="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 focus:border-slate-400 focus:outline-none"
                  required
                  minlength="8"
                />
              </div>
              <div>
                <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-slate-500">
                  Rol
                </label>
                <select
                  v-model="newUser.role"
                  class="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700"
                >
                  <option value="admin">admin</option>
                  <option value="commercial">commercial</option>
                  <option value="designer">designer</option>
                  <option value="workshop">workshop</option>
                </select>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <button
                  type="submit"
                  class="rounded-lg bg-[#002395] px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-[#001c77]"
                  :disabled="isCreating"
                >
                  {{ isCreating ? 'Creando...' : 'Crear usuario' }}
                </button>
                <button
                  type="button"
                  class="rounded-lg border border-slate-200 px-4 py-2 text-sm text-slate-600 hover:bg-slate-50"
                  @click="resetNewUser"
                >
                  Limpiar
                </button>
              </div>
              <p class="text-xs text-slate-500">
                Las contraseñas se almacenan con hash bcrypt. El acceso se controla por rol.
              </p>
            </form>
          </div>

          <div class="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
            <h2 class="text-lg font-semibold text-slate-900">Roles disponibles</h2>
            <p class="mb-4 text-xs text-slate-500">Referencia rápida de permisos (MVP).</p>
            <ul class="grid gap-2 text-sm text-slate-700 sm:grid-cols-2">
              <li class="rounded-lg bg-slate-50 p-3">
                <strong class="text-[#002395]">admin</strong>
                <p class="text-xs text-slate-500">Acceso total y administración.</p>
              </li>
              <li class="rounded-lg bg-slate-50 p-3">
                <strong class="text-[#002395]">commercial</strong>
                <p class="text-xs text-slate-500">Gestión comercial y clientes.</p>
              </li>
              <li class="rounded-lg bg-slate-50 p-3">
                <strong class="text-[#002395]">designer</strong>
                <p class="text-xs text-slate-500">Diseño de pedidos.</p>
              </li>
              <li class="rounded-lg bg-slate-50 p-3">
                <strong class="text-[#002395]">workshop</strong>
                <p class="text-xs text-slate-500">Operaciones de taller.</p>
              </li>
            </ul>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <input
            v-model.trim="search"
            type="text"
            placeholder="Buscar por email o rol"
            class="w-full rounded-full border border-slate-200 bg-white px-4 py-2 text-sm text-slate-700 focus:border-slate-400 focus:outline-none md:max-w-md"
          />
        </div>

        <div v-if="error" class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">
          {{ error }}
        </div>

        <div class="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
          <table class="min-w-full divide-y divide-slate-200 text-sm">
            <thead class="bg-slate-50 text-slate-600">
              <tr>
                <th class="px-4 py-3 text-left font-semibold">Email</th>
                <th class="px-4 py-3 text-left font-semibold">Rol</th>
                <th class="px-4 py-3 text-right font-semibold">Acciones</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="user in filteredUsers" :key="user.id">
                <td class="px-4 py-3 text-slate-800">{{ user.email }}</td>
                <td class="px-4 py-3">
                  <span class="rounded-full bg-slate-100 px-2 py-1 text-xs font-medium text-slate-700">
                    {{ user.role }}
                  </span>
                </td>
                <td class="px-4 py-3 text-right">
                  <div class="flex justify-end gap-2">
                    <button
                      class="rounded-md border border-slate-200 px-3 py-1 text-xs font-medium text-slate-700 hover:bg-slate-50"
                      @click="openRoleModal(user)"
                    >
                      Asignar rol
                    </button>
                    <button
                      class="rounded-md border border-red-200 px-3 py-1 text-xs font-medium text-red-600 hover:bg-red-50"
                      @click="confirmDelete(user)"
                    >
                      Eliminar
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!isLoading && filteredUsers.length === 0">
                <td class="px-4 py-6 text-center text-slate-500" colspan="3">
                  No hay usuarios para mostrar.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-lg">
        <h2 class="mb-2 text-lg font-semibold text-slate-900">Asignar rol</h2>
        <p class="mb-4 text-sm text-slate-600">
          Usuario: <span class="font-medium text-slate-800">{{ selectedUser?.email }}</span>
        </p>

        <label class="mb-2 block text-xs font-semibold uppercase tracking-wide text-slate-500">Rol</label>
        <select
          v-model="selectedRole"
          class="mb-4 w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700"
        >
          <option value="admin">admin</option>
          <option value="commercial">commercial</option>
          <option value="designer">designer</option>
          <option value="workshop">workshop</option>
        </select>

        <div class="flex justify-end gap-2">
          <button
            class="rounded-lg border border-slate-200 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50"
            @click="closeModal"
          >
            Cancelar
          </button>
          <button
            class="rounded-lg bg-slate-900 px-4 py-2 text-sm text-white hover:bg-slate-800"
            @click="saveRole"
            :disabled="isSaving"
          >
            {{ isSaving ? 'Guardando...' : 'Guardar' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { Usuario, UserRole } from '@/types/auth'
import { iamService } from '@/services/iam'

const authStore = useAuthStore()

const users = ref<Usuario[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const search = ref('')

const showModal = ref(false)
const selectedUser = ref<Usuario | null>(null)
const selectedRole = ref<UserRole>('commercial')
const isSaving = ref(false)
const isCreating = ref(false)
const isDeleting = ref(false)

const newUser = ref({
  email: '',
  password: '',
  role: 'commercial' as UserRole
})

const isAdmin = computed(() => authStore.isAdmin)

const filteredUsers = computed(() => {
  const term = search.value.toLowerCase()
  if (!term) return users.value
  return users.value.filter((user) =>
    user.email.toLowerCase().includes(term) || user.role.toLowerCase().includes(term)
  )
})

const loadUsers = async () => {
  if (!isAdmin.value) return

  isLoading.value = true
  error.value = null

  try {
    users.value = await iamService.listUsers()
  } catch (err: any) {
    error.value = err?.response?.data?.error || err?.message || 'No se pudo cargar el listado.'
  } finally {
    isLoading.value = false
  }
}

const openRoleModal = (user: Usuario) => {
  selectedUser.value = user
  selectedRole.value = user.role
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedUser.value = null
}

const saveRole = async () => {
  if (!selectedUser.value) return

  isSaving.value = true
  try {
    const result = await iamService.assignRole(selectedUser.value.id, selectedRole.value)
    users.value = users.value.map((user) =>
      user.id === result.userId ? { ...user, role: result.role } : user
    )
    closeModal()
  } catch (err: any) {
    error.value = err?.response?.data?.error || err?.message || 'No se pudo asignar el rol.'
  } finally {
    isSaving.value = false
  }
}

const resetNewUser = () => {
  newUser.value = {
    email: '',
    password: '',
    role: 'commercial'
  }
}

const createUser = async () => {
  if (!isAdmin.value) return

  isCreating.value = true
  error.value = null

  try {
    const created = await iamService.createUser({
      email: newUser.value.email,
      password: newUser.value.password,
      role: newUser.value.role
    })
    users.value = [created, ...users.value]
    resetNewUser()
  } catch (err: any) {
    error.value = err?.response?.data?.error || err?.message || 'No se pudo crear el usuario.'
  } finally {
    isCreating.value = false
  }
}

const confirmDelete = async (user: Usuario) => {
  if (!isAdmin.value || isDeleting.value) return

  const confirmed = window.confirm(`¿Eliminar al usuario ${user.email}?`)
  if (!confirmed) return

  isDeleting.value = true
  error.value = null

  try {
    await iamService.deleteUser(user.id)
    users.value = users.value.filter((item) => item.id !== user.id)
  } catch (err: any) {
    error.value = err?.response?.data?.error || err?.message || 'No se pudo eliminar el usuario.'
  } finally {
    isDeleting.value = false
  }
}

onMounted(() => {
  loadUsers()
})
</script>
