<template>
  <div class="min-h-screen bg-[#f1f5f9]">
    <div class="flex min-h-screen">
      <aside class="hidden w-64 flex-col border-r border-[#e2e8f0] bg-white lg:flex">
        <div class="bg-[#002395] px-6 py-5">
          <span class="text-xl font-bold italic text-white">TramaTex</span>
        </div>
        <nav class="flex-1 px-4 py-6 text-sm text-[#1e293b]">
          <p class="mb-3 text-xs font-semibold uppercase tracking-wide text-[#64748b]">Administración</p>
          <ul class="space-y-2">
            <li class="rounded-lg bg-[#002395] px-3 py-2 text-white">Usuarios</li>
            <li class="rounded-lg px-3 py-2 text-[#1e293b] opacity-70">Roles y permisos</li>
            <li class="rounded-lg px-3 py-2 text-[#1e293b] opacity-70">Auditoría</li>
          </ul>
        </nav>
      </aside>

      <main class="flex-1">
        <header class="border-b border-[#e2e8f0] bg-white">
          <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4 px-6 py-4">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wide text-[#64748b]">Administración</p>
              <h1 class="text-2xl font-bold text-[#002395]">Gestión de usuarios</h1>
            </div>
            <div class="flex items-center gap-2">
              <button
                class="rounded-full border border-[#e2e8f0] bg-white px-4 py-2 text-sm font-medium text-[#1e293b] shadow-sm hover:border-[#cbd5f5]"
                @click="loadUsers"
                :disabled="isLoading"
              >
                {{ isLoading ? 'Cargando...' : 'Refrescar' }}
              </button>
            </div>
          </div>
        </header>

        <section class="mx-auto max-w-6xl space-y-6 px-6 py-6">

          <div v-if="!isAdmin" class="rounded-lg border border-[#f59e0b] bg-[#fef3c7] p-4 text-sm text-[#92400e]">
            Solo el rol admin puede gestionar usuarios.
          </div>

          <div v-else class="space-y-6">
            <div class="grid gap-4 xl:grid-cols-[1.15fr_0.85fr]">
              <div class="rounded-2xl border border-[#e2e8f0] bg-white p-6 shadow-sm">
                <div class="mb-4 flex items-center justify-between">
                  <div>
                    <h2 class="text-lg font-semibold text-[#1e293b]">Alta de usuario</h2>
                    <p class="text-xs text-[#64748b]">Crea cuentas con roles predefinidos.</p>
                  </div>
                  <span class="rounded-full bg-[#002395] px-3 py-1 text-xs font-semibold uppercase tracking-wide text-white">
                    Admin
                  </span>
                </div>

                <form class="grid gap-4" @submit.prevent="createUser">
                  <div>
                    <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-[#64748b]">
                      Email
                    </label>
                    <input
                      v-model.trim="newUser.email"
                      type="email"
                      placeholder="usuario@tramatex.local"
                      class="w-full rounded-lg border border-[#e2e8f0] bg-white px-3 py-2 text-sm text-[#1e293b] focus:border-[#002395] focus:outline-none"
                      required
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-[#64748b]">
                      Contraseña
                    </label>
                    <input
                      v-model="newUser.password"
                      type="password"
                      placeholder="Mínimo 8 caracteres"
                      class="w-full rounded-lg border border-[#e2e8f0] bg-white px-3 py-2 text-sm text-[#1e293b] focus:border-[#002395] focus:outline-none"
                      required
                      minlength="8"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-[#64748b]">
                      Rol
                    </label>
                    <select
                      v-model="newUser.role"
                      class="w-full rounded-lg border border-[#e2e8f0] bg-white px-3 py-2 text-sm text-[#1e293b]"
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
                      class="rounded-lg bg-[#E6B800] px-4 py-2 text-sm font-semibold text-[#1e293b] shadow-sm hover:bg-[#d6aa00]"
                      :disabled="isCreating"
                    >
                      {{ isCreating ? 'Creando...' : 'Crear usuario' }}
                    </button>
                    <button
                      type="button"
                      class="rounded-lg border border-[#e2e8f0] px-4 py-2 text-sm text-[#64748b] hover:bg-[#f1f5f9]"
                      @click="resetNewUser"
                    >
                      Limpiar
                    </button>
                  </div>
                  <p class="text-xs text-[#64748b]">
                    Las contraseñas se almacenan con hash bcrypt. El acceso se controla por rol.
                  </p>
                </form>
              </div>

              <div class="rounded-2xl border border-[#e2e8f0] bg-white p-6 shadow-sm">
                <h2 class="text-lg font-semibold text-[#1e293b]">Roles disponibles</h2>
                <p class="mb-4 text-xs text-[#64748b]">Referencia rápida de permisos (MVP).</p>
                <ul class="grid gap-3 text-sm text-[#1e293b]">
                  <li class="rounded-lg border border-[#e2e8f0] bg-[#f8fafc] p-3">
                    <strong class="text-[#002395]">admin</strong>
                    <p class="text-xs text-[#64748b]">Acceso total y administración.</p>
                  </li>
                  <li class="rounded-lg border border-[#e2e8f0] bg-[#f8fafc] p-3">
                    <strong class="text-[#002395]">commercial</strong>
                    <p class="text-xs text-[#64748b]">Gestión comercial y clientes.</p>
                  </li>
                  <li class="rounded-lg border border-[#e2e8f0] bg-[#f8fafc] p-3">
                    <strong class="text-[#002395]">designer</strong>
                    <p class="text-xs text-[#64748b]">Diseño de pedidos.</p>
                  </li>
                  <li class="rounded-lg border border-[#e2e8f0] bg-[#f8fafc] p-3">
                    <strong class="text-[#002395]">workshop</strong>
                    <p class="text-xs text-[#64748b]">Operaciones de taller.</p>
                  </li>
                </ul>
              </div>
            </div>

            <div class="rounded-2xl border border-[#e2e8f0] bg-white p-5 shadow-sm">
              <div class="flex flex-wrap items-center gap-3">
                <div class="flex-1">
                  <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-[#64748b]">
                    Buscar por nombre/email
                  </label>
                  <input
                    v-model.trim="search"
                    type="text"
                    placeholder="Buscar por email"
                    class="w-full rounded-full border border-[#e2e8f0] bg-white px-4 py-2 text-sm text-[#1e293b] focus:border-[#002395] focus:outline-none"
                  />
                </div>
                <div class="min-w-[180px]">
                  <label class="mb-1 block text-xs font-semibold uppercase tracking-wide text-[#64748b]">
                    Filtrar por rol
                  </label>
                  <select
                    v-model="roleFilter"
                    class="w-full rounded-full border border-[#e2e8f0] bg-white px-4 py-2 text-sm text-[#1e293b]"
                  >
                    <option value="">Todos</option>
                    <option value="admin">admin</option>
                    <option value="commercial">commercial</option>
                    <option value="designer">designer</option>
                    <option value="workshop">workshop</option>
                  </select>
                </div>
              </div>
            </div>

            <div v-if="error" class="rounded-lg border border-[#ef4444] bg-[#fee2e2] p-3 text-sm text-[#991b1b]">
              {{ error }}
            </div>

            <div class="overflow-hidden rounded-2xl border border-[#e2e8f0] bg-white shadow-sm">
              <table class="min-w-full divide-y divide-[#e2e8f0] text-sm">
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
        </section>
      </main>
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
const roleFilter = ref<UserRole | ''>('')

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
  return users.value.filter((user) => {
    const matchesTerm = !term || user.email.toLowerCase().includes(term)
    const matchesRole = !roleFilter.value || user.role === roleFilter.value
    return matchesTerm && matchesRole
  })
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
