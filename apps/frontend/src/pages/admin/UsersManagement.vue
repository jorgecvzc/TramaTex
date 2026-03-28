<template>
  <Navbar class="no-print" />
  
  <BaseDashboardPage :is-loading="isLoading" class="no-print">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <PageHeader 
        title="Gestión de Cuentas de Usuario" 
        :breadcrumbs="[{ label: 'Administración', to: '/admin/users' }, { label: 'Usuarios' }]"
      >
        <template #icon>
          <span class="material-symbols-outlined">admin_panel_settings</span>
        </template>
        <template #actions>
          <button class="btn btn-outline" @click="loadUsers" :disabled="isLoading">
            <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
            <span>Sincronizar Lista</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <!-- CAPA 3: TRABAJO (Área Principal) -->
    <div class="admin-main-content">
      <div v-if="!isAdmin" class="alert-card card error mb-6">
        <span class="material-symbols-outlined">warning</span>
        <p>Acceso restringido. Solo administradores pueden gestionar cuentas.</p>
      </div>

      <!-- Alta de Usuario -->
      <FormSection title="Alta de Nuevo Usuario" icon="person_add">
        <form @submit.prevent="createUser" class="user-creation-form">
          <div class="form-row">
            <div class="form-group">
              <label>Correo Electrónico *</label>
              <div class="input-with-icon">
                <span class="material-symbols-outlined icon-start">mail</span>
                <input v-model.trim="newUser.email" type="email" class="form-input" placeholder="usuario@tramatex.local" required />
              </div>
            </div>
            <div class="form-group">
              <label>Contraseña de Acceso *</label>
              <div class="input-with-icon">
                <span class="material-symbols-outlined icon-start">lock</span>
                <input v-model="newUser.password" type="password" class="form-input" placeholder="Mínimo 8 caracteres" required minlength="8" />
              </div>
            </div>
          </div>
          <div class="form-row mt-4">
            <div class="form-group">
              <label>Rol Asignado</label>
              <select v-model="newUser.role" class="form-input">
                <option value="admin">Administrador (Total)</option>
                <option value="commercial">Comercial (Ventas)</option>
                <option value="designer">Diseñador (Productos)</option>
                <option value="workshop">Taller (MES)</option>
              </select>
            </div>
            <div class="form-actions-inline align-end">
              <button type="submit" class="btn btn-primary" :disabled="isCreating">
                <span class="material-symbols-outlined">{{ isCreating ? 'sync' : 'save' }}</span>
                <span>{{ isCreating ? 'Procesando...' : 'Registrar Usuario' }}</span>
              </button>
              <button type="button" class="btn btn-outline" @click="resetNewUser">Limpiar</button>
            </div>
          </div>
        </form>
      </FormSection>

      <!-- Listado de Usuarios -->
      <section class="dashboard-section mt-8">
        <div class="section-header">
          <span class="material-symbols-outlined">groups</span>
          <h2>Cuentas Registradas</h2>
          <div class="header-filters ml-auto">
            <input v-model.trim="search" type="text" class="form-input-sm" placeholder="Buscar por email..." />
            <select v-model="roleFilter" class="form-input-sm ml-2">
              <option value="">Todos los roles</option>
              <option value="admin">Administrador</option>
              <option value="commercial">Comercial</option>
              <option value="designer">Diseñador</option>
              <option value="workshop">Taller</option>
            </select>
          </div>
        </div>

        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Usuario</th>
                <th>Rol</th>
                <th>Estado</th>
                <th class="align-right">Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in filteredUsers" :key="user.id" class="row-hover">
                <td>
                  <div class="user-info-cell">
                    <span class="material-symbols-outlined user-avatar">account_circle</span>
                    <strong>{{ user.email }}</strong>
                  </div>
                </td>
                <td>
                  <span :class="['role-tag', `role-${user.role}`]">{{ user.role }}</span>
                </td>
                <td>
                  <span class="status-badge status-success">Activo</span>
                </td>
                <td class="align-right">
                  <div class="action-buttons">
                    <button class="btn-icon" @click="openRoleModal(user)" title="Cambiar permisos">
                      <span class="material-symbols-outlined">key</span>
                    </button>
                    <button class="btn-icon text-danger" @click="confirmDelete(user)" title="Eliminar cuenta">
                      <span class="material-symbols-outlined">delete</span>
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredUsers.length === 0">
                <td colspan="4" class="empty-row-msg">No se han encontrado usuarios con estos criterios.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <!-- CAPA 2: CONTEXTO (Panel Lateral) -->
    <template #sidebar>
      <div class="admin-sidebar-content">
        <section class="sidebar-section">
          <div class="section-header">
            <span class="material-symbols-outlined">verified_user</span>
            <h2>Guía de Permisos</h2>
          </div>
          <div class="roles-guide-list mt-4">
            <div class="role-guide-item">
              <span class="dot admin"></span>
              <div class="content">
                <strong>Administrador</strong>
                <p>Gestión total, configuración fiscal y borrado.</p>
              </div>
            </div>
            <div class="role-guide-item">
              <span class="dot commercial"></span>
              <div class="content">
                <strong>Comercial</strong>
                <p>Gestión de clientes y ciclo de ventas.</p>
              </div>
            </div>
            <div class="role-guide-item">
              <span class="dot designer"></span>
              <div class="content">
                <strong>Diseñador</strong>
                <p>Edición de productos y variantes JIT.</p>
              </div>
            </div>
            <div class="role-guide-item">
              <span class="dot workshop"></span>
              <div class="content">
                <strong>Operario Taller</strong>
                <p>Acceso a terminal y reporte MES.</p>
              </div>
            </div>
          </div>
        </section>

        <section class="help-notice mt-10">
          <div class="notice-header">
            <span class="material-symbols-outlined">info</span>
            <h3>Seguridad</h3>
          </div>
          <p class="help-text">
            Los cambios de rol afectan a los permisos de acceso de forma inmediata. Asegúrese de que cada usuario tiene el nivel mínimo de privilegio necesario.
          </p>
        </section>
      </div>
    </template>
  </BaseDashboardPage>

  <!-- MODAL: CAMBIO DE ROL -->
  <BaseDialog
    :show="showModal"
    title="Modificar Permisos de Usuario"
    icon="manage_accounts"
    confirm-text="Actualizar Rol"
    :is-confirming="isSaving"
    @close="closeModal"
    @confirm="saveRole"
  >
    <div class="modal-info-box mb-4" v-if="selectedUser">
      <p>Actualizando cuenta:</p>
      <strong>{{ selectedUser.email }}</strong>
    </div>
    <div class="form-group">
      <label>Nuevo Rol de Usuario</label>
      <select v-model="selectedRole" class="form-input">
        <option value="admin">Administrador</option>
        <option value="commercial">Comercial</option>
        <option value="designer">Diseñador</option>
        <option value="workshop">Taller</option>
      </select>
    </div>
  </BaseDialog>

  <!-- MODAL: ELIMINAR USUARIO -->
  <BaseDialog
    :show="showDeleteConfirm"
    title="Eliminar Cuenta de Usuario"
    icon="warning"
    confirm-text="Eliminar Definitivamente"
    confirm-class="btn-danger"
    :is-confirming="isDeleting"
    @close="showDeleteConfirm = false"
    @confirm="executeDelete"
  >
    <p>¿Está seguro de que desea eliminar permanentemente la cuenta de <strong>{{ userToDelete?.email }}</strong>?</p>
    <p class="mt-2 text-danger italic">Esta acción no se puede deshacer.</p>
  </BaseDialog>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive, computed } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import PageHeader from '@/components/layout/PageHeader.vue'
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue'
import FormSection from '@/components/shared/FormSection.vue'
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
const showDeleteConfirm = ref(false)
const selectedUser = ref<Usuario | null>(null)
const userToDelete = ref<Usuario | null>(null)
const selectedRole = ref<UserRole>('commercial')
const isSaving = ref(false)
const isCreating = ref(false)
const isDeleting = ref(false)

const newUser = reactive({ email: '', password: '', role: 'commercial' as UserRole })

const isAdmin = computed(() => authStore.isAdmin)

const filteredUsers = computed(() => {
  const term = search.value.toLowerCase()
  return users.value.filter((user) => {
    const matchesTerm = !term || user.email.toLowerCase().includes(term)
    const matchesRole = !roleFilter.value || user.role === roleFilter.value
    return matchesTerm && matchesRole
  })
})

async function loadUsers() {
  if (!isAdmin.value) return
  isLoading.value = true; error.value = null
  try { users.value = await iamService.listUsers() }
  catch (err: any) { error.value = err?.message || 'Error al cargar usuarios.' }
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

function resetNewUser() { Object.assign(newUser, { email: '', password: '', role: 'commercial' }) }

async function createUser() {
  if (!isAdmin.value) return
  isCreating.value = true; error.value = null
  try {
    const created = await iamService.createUser({ email: newUser.email, password: newUser.password, role: newUser.role })
    users.value = [created, ...users.value]
    resetNewUser()
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
@import "@/design-system/_sections.css";

.admin-main-content { display: flex; flex-direction: column; gap: 0; }

.dashboard-section { background: white; padding: 1.5rem; border-radius: 12px; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.section-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; padding-bottom: 0.75rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.9rem; font-weight: 700; margin: 0; text-transform: uppercase; letter-spacing: 0.05em; }
.header-tag { font-size: 0.65rem; font-weight: 800; padding: 0.2rem 0.6rem; background: var(--color-background); color: var(--color-secondary); border-radius: 20px; }

.user-info-cell { display: flex; align-items: center; gap: 0.75rem; }
.user-avatar { color: var(--color-border); font-size: 32px; }

.role-tag { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; background: var(--color-background); border: 1px solid var(--color-border); }
.role-admin { border-left: 3px solid #dc2626; color: #dc2626; }
.role-commercial { border-left: 3px solid #2563eb; color: #2563eb; }
.role-designer { border-left: 3px solid #9333ea; color: #9333ea; }
.role-workshop { border-left: 3px solid #d97706; color: #d97706; }

.role-guide-item { display: flex; align-items: flex-start; gap: 1rem; padding: 0.75rem 0; }
.dot { width: 10px; height: 10px; border-radius: 50%; margin-top: 5px; flex-shrink: 0; }
.dot.admin { background: #dc2626; }
.dot.commercial { background: #2563eb; }
.dot.designer { background: #9333ea; }
.dot.workshop { background: #d97706; }

.role-guide-item strong { font-size: 0.85rem; display: block; }
.role-guide-item p { font-size: 0.75rem; color: var(--color-text-secondary); margin: 0; }

.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-actions-inline { display: flex; gap: 1rem; }
.align-end { align-items: flex-end; }

.input-with-icon { position: relative; display: flex; align-items: center; width: 100%; }
.icon-start { position: absolute; left: 0.85rem; font-size: 20px; color: var(--color-text-secondary); }
.input-with-icon .form-input { padding-left: 2.75rem; }

.form-input { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-input-sm { padding: 0.4rem 0.75rem; border-radius: 6px; border: 1px solid var(--color-border); font-size: 0.85rem; }

.modal-info-box { padding: 1rem; background: var(--color-background); border-radius: 8px; font-size: 0.9rem; }
.help-notice { padding: 1.25rem; background: rgba(59, 130, 246, 0.05); border-radius: 12px; border: 1px dashed rgba(59, 130, 246, 0.3); }
.notice-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.75rem; color: #2563eb; }
.notice-header h3 { margin: 0; font-size: 0.85rem; text-transform: uppercase; }
.help-text { font-size: 0.8rem; color: var(--color-text-secondary); line-height: 1.5; margin: 0; }

.align-right { text-align: right; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.5rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>
