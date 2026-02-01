<template>
  <div class="min-h-screen bg-gray-50 p-6">
    <div class="mx-auto max-w-5xl">
      <div class="mb-6 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">Gestión de usuarios</h1>
          <p class="text-sm text-slate-600">Administración de usuarios y asignación de roles.</p>
        </div>
        <button
          class="rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
          @click="loadUsers"
          :disabled="isLoading"
        >
          {{ isLoading ? 'Cargando...' : 'Refrescar' }}
        </button>
      </div>

      <div v-if="!isAdmin" class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-amber-800">
        Solo el rol admin puede gestionar usuarios.
      </div>

      <div v-else>
        <div class="mb-4 flex items-center gap-3">
          <input
            v-model.trim="search"
            type="text"
            placeholder="Buscar por email o rol"
            class="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 focus:border-slate-400 focus:outline-none"
          />
        </div>

        <div v-if="error" class="mb-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">
          {{ error }}
        </div>

        <div class="overflow-hidden rounded-lg border border-slate-200 bg-white">
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
                  <button
                    class="rounded-md border border-slate-200 px-3 py-1 text-xs font-medium text-slate-700 hover:bg-slate-50"
                    @click="openRoleModal(user)"
                  >
                    Asignar rol
                  </button>
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

        <div class="mt-6 rounded-lg border border-slate-200 bg-white p-4">
          <h2 class="mb-2 text-sm font-semibold text-slate-800">Permisos (MVP)</h2>
          <p class="mb-3 text-sm text-slate-600">
            La gestión fina de permisos está planificada; por ahora se controla por rol.
          </p>
          <ul class="grid gap-2 text-sm text-slate-700 sm:grid-cols-2">
            <li class="rounded-md bg-slate-50 p-2"><strong>admin</strong>: acceso total</li>
            <li class="rounded-md bg-slate-50 p-2"><strong>commercial</strong>: gestión comercial</li>
            <li class="rounded-md bg-slate-50 p-2"><strong>designer</strong>: diseño de pedidos</li>
            <li class="rounded-md bg-slate-50 p-2"><strong>workshop</strong>: operaciones de taller</li>
          </ul>
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

onMounted(() => {
  loadUsers()
})
</script>
