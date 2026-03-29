<template>
  <div class="page-layout">
    <BaseCatalog
      title="Gestión de Usuarios"
      icon="manage_accounts"
      :breadcrumbs="[{ label: 'Administración', to: '/admin/users' }, { label: 'Usuarios' }]"
      :items="filteredUsers"
      :is-loading="isLoading"
      :error="error || undefined"
      create-text="Nuevo Usuario"
      @refresh="loadUsers"
    >
      <!-- CAPA 2: CONTEXTO (Filtros y Alta Rápida) -->
      <template #filters>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input v-model.trim="search" type="text" placeholder="Email..." />
        </div>
        <div class="filter-group">
          <label>Rol</label>
          <select v-model="roleFilter">
            <option value="">Todos los roles</option>
            <option value="admin">Administrador</option>
            <option value="commercial">Comercial</option>
            <option value="designer">Diseñador</option>
            <option value="workshop">Taller</option>
          </select>
        </div>
        <div class="filter-actions-inline ml-auto">
          <button class="btn btn-primary btn-sm" @click="showCreateDialog = true">
            <span class="material-symbols-outlined">person_add</span>
            Nuevo Usuario
          </button>
        </div>
      </template>

      <!-- CAPA 3: TRABAJO (Tabla) -->
      <template #table-header>
        <th>Usuario</th>
        <th>Rol Asignado</th>
        <th>Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item: user }">
        <td>
          <div class="user-cell">
            <span class="material-symbols-outlined avatar">account_circle</span>
            <strong>{{ user.email }}</strong>
          </div>
        </td>
        <td>
          <span :class="['role-pill', `role-${user.role}`]">{{ user.role }}</span>
        </td>
        <td><span class="status-badge status-success">Activo</span></td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button class="btn-icon" @click="openRoleModal(user)" title="Cambiar permisos">
              <span class="material-symbols-outlined">key</span>
            </button>
            <button class="btn-icon text-danger" @click="confirmDelete(user)" title="Eliminar cuenta">
              <span class="material-symbols-outlined">delete</span>
            </button>
          </div>
        </td>
      </template>

      <!-- Sidebar Informativo (Inyectado vía slot si el componente lo permite o mantenido en el layout) -->
    </BaseCatalog>

    <!-- DIÁLOGOS DE GESTIÓN -->
    <BaseDialog :show="showCreateDialog" title="Registrar Nuevo Usuario" icon="person_add" confirm-text="Crear Usuario" :is-confirming="isCreating" @close="showCreateDialog = false" @confirm="createUser">
      <div class="form-grid">
        <div class="form-group">
          <label>Email Corporativo</label>
          <input v-model.trim="newUser.email" type="email" placeholder="ejemplo@tramatex.com" />
        </div>
        <div class="form-group">
          <label>Contraseña</label>
          <input v-model="newUser.password" type="password" placeholder="Mínimo 8 caracteres" />
        </div>
        <div class="form-group">
          <label>Rol Inicial</label>
          <select v-model="newUser.role">
            <option value="admin">Administrador</option>
            <option value="commercial">Comercial</option>
            <option value="designer">Diseñador</option>
            <option value="workshop">Taller</option>
          </select>
        </div>
      </div>
    </BaseDialog>

    <BaseDialog :show="showModal" title="Modificar Permisos" icon="key" confirm-text="Actualizar" :is-confirming="isSaving" @close="closeModal" @confirm="saveRole">
      <p class="mb-4">Cambiando permisos para: <strong>{{ selectedUser?.email }}</strong></p>
      <div class="form-group">
        <label>Nuevo Rol</label>
        <select v-model="selectedRole">
          <option value="admin">Administrador</option>
          <option value="commercial">Comercial</option>
          <option value="designer">Diseñador</option>
          <option value="workshop">Taller</option>
        </select>
      </div>
    </BaseDialog>

    <BaseDialog :show="showDeleteConfirm" title="Eliminar Usuario" icon="warning" confirm-text="Eliminar" confirm-class="btn-danger" :is-confirming="isDeleting" @close="showDeleteConfirm = false" @confirm="executeDelete">
      <p>¿Seguro que deseas eliminar permanentemente a <strong>{{ userToDelete?.email }}</strong>?</p>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive, computed } from 'vue'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
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
const showCreateDialog = ref(false)
const showDeleteConfirm = ref(false)
const selectedUser = ref<Usuario | null>(null)
const userToDelete = ref<Usuario | null>(null)
const selectedRole = ref<UserRole>('commercial')
const isSaving = ref(false)
const isCreating = ref(false)
const isDeleting = ref(false)
const newUser = reactive({ email: '', password: '', role: 'commercial' as UserRole })

const filteredUsers = computed(() => {
  const term = search.value.toLowerCase()
  return users.value.filter((u) => (!term || u.email.toLowerCase().includes(term)) && (!roleFilter.value || u.role === roleFilter.value))
})

async function loadUsers() {
  isLoading.value = true; error.value = null
  try { users.value = await iamService.listUsers() }
  catch (err: any) { error.value = err.message }
  finally { isLoading.value = false }
}

function openRoleModal(user: Usuario) { selectedUser.value = user; selectedRole.value = user.role; showModal.value = true; }
function closeModal() { showModal.value = false; selectedUser.value = null; }

async function saveRole() {
  if (!selectedUser.value) return
  isSaving.value = true
  try {
    const result = await iamService.assignRole(selectedUser.value.id, selectedRole.value)
    users.value = users.value.map((u) => u.id === result.userId ? { ...u, role: result.role } : u)
    closeModal()
  } catch (err: any) { alert(err.message) }
  finally { isSaving.value = false }
}

async function createUser() {
  isCreating.value = true
  try {
    const created = await iamService.createUser(newUser)
    users.value = [created, ...users.value]
    showCreateDialog.value = false
    Object.assign(newUser, { email: '', password: '', role: 'commercial' })
  } catch (err: any) { alert(err.message) }
  finally { isCreating.value = false }
}

function confirmDelete(user: Usuario) { userToDelete.value = user; showDeleteConfirm.value = true; }
async function executeDelete() {
  if (!userToDelete.value) return
  isDeleting.value = true
  try {
    await iamService.deleteUser(userToDelete.value.id)
    users.value = users.value.filter((item) => item.id !== userToDelete.value!.id)
    showDeleteConfirm.value = false
  } catch (err: any) { alert(err.message) }
  finally { isDeleting.value = false }
}

onMounted(loadUsers)
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.user-cell { display: flex; align-items: center; gap: 0.75rem; }
.avatar { color: var(--color-border); font-size: 24px; }
.role-pill { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; background: var(--color-background-soft); border: 1px solid var(--color-border); }
.role-admin { border-left: 3px solid var(--color-error); color: var(--color-error); }
.role-commercial { border-left: 3px solid #2563eb; color: #2563eb; }
.role-designer { border-left: 3px solid #9333ea; color: #9333ea; }
.role-workshop { border-left: 3px solid #d97706; color: #d97706; }
.align-right { text-align: right; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; transition: 0.2s; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }
</style>
