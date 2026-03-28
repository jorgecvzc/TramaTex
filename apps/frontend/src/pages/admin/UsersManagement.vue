<template>
  <Navbar />
  
  <div class="main-container">
    <PageHeader 
      title="Gestión de Usuarios" 
      :breadcrumbs="[{ label: 'Administración', to: '/admin/users' }, { label: 'Usuarios' }]"
    >
      <template #actions>
        <button class="btn btn-secondary" @click="loadUsers" :disabled="isLoading">
          <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
          Actualizar Lista
        </button>
      </template>
    </PageHeader>

    <div v-if="!isAdmin" class="alert alert-warning card">
      <span class="material-symbols-outlined">warning</span>
      <p>Acceso restringido. Solo el personal administrador puede gestionar cuentas de usuario.</p>
    </div>

    <div v-else class="admin-layout">
      <!-- Top Section: Creation and Info -->
      <div class="cards-grid">
        <section class="card form-card">
          <div class="card-header">
            <div class="header-icon">
              <span class="material-symbols-outlined">person_add</span>
            </div>
            <div>
              <h2>Alta de Nuevo Usuario</h2>
              <p>Crea una cuenta con permisos específicos.</p>
            </div>
          </div>

          <form class="user-form" @submit.prevent="createUser">
            <div class="form-group">
              <label>Correo Electrónico</label>
              <div class="input-with-icon">
                <span class="material-symbols-outlined icon-start">mail</span>
                <input
                  v-model.trim="newUser.email"
                  type="email"
                  placeholder="usuario@tramatex.local"
                  required
                />
              </div>
            </div>
            
            <div class="form-group">
              <label>Contraseña de Acceso</label>
              <div class="input-with-icon">
                <span class="material-symbols-outlined icon-start">lock</span>
                <input
                  v-model="newUser.password"
                  type="password"
                  placeholder="Mínimo 8 caracteres"
                  required
                  minlength="8"
                />
              </div>
            </div>

            <div class="form-group">
              <label>Rol Asignado</label>
              <select v-model="newUser.role">
                <option value="admin">Administrador (Total)</option>
                <option value="commercial">Comercial (Ventas)</option>
                <option value="designer">Diseñador (Productos)</option>
                <option value="workshop">Taller (MES)</option>
              </select>
            </div>

            <div class="form-actions">
              <button type="submit" class="btn btn-primary" :disabled="isCreating">
                <span class="material-symbols-outlined">{{ isCreating ? 'sync' : 'save' }}</span>
                <span>{{ isCreating ? 'Creando...' : 'Registrar Usuario' }}</span>
              </button>
              <button type="button" class="btn btn-outline" @click="resetNewUser">
                Limpiar
              </button>
            </div>
          </form>
        </section>

        <section class="card help-card">
          <div class="card-header">
            <div class="header-icon secondary">
              <span class="material-symbols-outlined">verified_user</span>
            </div>
            <div>
              <h2>Permisos por Rol</h2>
              <p>Resumen de capacidades del sistema.</p>
            </div>
          </div>
          
          <ul class="roles-guide">
            <li>
              <strong>admin</strong>
              <span>Gestión total, configuración fiscal y borrado de registros.</span>
            </li>
            <li>
              <strong>commercial</strong>
              <span>Gestión de clientes, presupuestos y pedidos de venta.</span>
            </li>
            <li>
              <strong>designer</strong>
              <span>Edición de catálogo de productos, variantes y atributos.</span>
            </li>
            <li>
              <strong>workshop</strong>
              <span>Acceso al terminal de taller y reporte de operaciones MES.</span>
            </li>
          </ul>
        </section>
      </div>

      <!-- Bottom Section: Users Table -->
      <section class="card table-card mt-6">
        <div class="table-header-filters">
          <div class="filter-group">
            <label>Buscar usuario</label>
            <div class="input-with-icon">
              <span class="material-symbols-outlined icon-start">search</span>
              <input v-model.trim="search" type="text" placeholder="Email o nombre..." />
            </div>
          </div>
          <div class="filter-group">
            <label>Filtrar por rol</label>
            <select v-model="roleFilter">
              <option value="">Todos los roles</option>
              <option value="admin">Administrador</option>
              <option value="commercial">Comercial</option>
              <option value="designer">Diseñador</option>
              <option value="workshop">Taller</option>
            </select>
          </div>
        </div>

        <div v-if="error" class="alert alert-error">
          <span class="material-symbols-outlined">error</span>
          <p>{{ error }}</p>
        </div>

        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Usuario (Email)</th>
                <th>Rol Actual</th>
                <th class="align-right">Acciones de Gestión</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in filteredUsers" :key="user.id" class="row-hover">
                <td>
                  <div class="user-cell">
                    <span class="material-symbols-outlined user-avatar">account_circle</span>
                    <strong>{{ user.email }}</strong>
                  </div>
                </td>
                <td>
                  <span class="role-badge">{{ user.role }}</span>
                </td>
                <td class="align-right">
                  <div class="action-buttons">
                    <button class="btn btn-outline btn-sm" @click="openRoleModal(user)">
                      <span class="material-symbols-outlined">key</span>
                      Cambiar Rol
                    </button>
                    <button class="btn btn-danger btn-sm" @click="confirmDelete(user)">
                      <span class="material-symbols-outlined">delete</span>
                      Eliminar
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!isLoading && filteredUsers.length === 0">
                <td colspan="3" class="empty-row">No se han encontrado usuarios que coincidan con la búsqueda.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <!-- Role Assignment Modal -->
    <Transition name="fade">
      <div v-if="showModal" class="modal-backdrop">
        <div class="modal card">
          <div class="modal-header">
            <span class="material-symbols-outlined">manage_accounts</span>
            <h2>Modificar Rol de Usuario</h2>
          </div>
          
          <div class="modal-body">
            <p class="modal-info">Estás actualizando los permisos para: <br><strong>{{ selectedUser?.email }}</strong></p>
            
            <div class="form-group">
              <label>Nuevo Rol</label>
              <select v-model="selectedRole">
                <option value="admin">Administrador</option>
                <option value="commercial">Comercial</option>
                <option value="designer">Diseñador</option>
                <option value="workshop">Taller</option>
              </select>
            </div>
          </div>

          <div class="modal-actions">
            <button class="btn btn-outline" @click="closeModal">Cancelar</button>
            <button class="btn btn-primary" @click="saveRole" :disabled="isSaving">
              <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'check' }}</span>
              <span>{{ isSaving ? 'Guardando...' : 'Confirmar Cambio' }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import PageHeader from '@/components/layout/PageHeader.vue'
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
    error.value = err?.response?.data?.error || err?.message || 'Error al cargar usuarios.'
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
    users.value = users.value.map((u) => u.id === result.userId ? { ...u, role: result.role } : u)
    closeModal()
  } catch (err: any) {
    error.value = err?.response?.data?.error || err?.message || 'Error al asignar rol.'
  } finally {
    isSaving.value = false
  }
}

const resetNewUser = () => { newUser.value = { email: '', password: '', role: 'commercial' } }

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
    error.value = err?.response?.data?.error || err?.message || 'Error al crear usuario.'
  } finally {
    isCreating.value = false
  }
}

const confirmDelete = async (user: Usuario) => {
  if (!isAdmin.value || isDeleting.value) return
  if (!window.confirm(`¿Estás seguro de eliminar permanentemente a ${user.email}?`)) return
  isDeleting.value = true
  error.value = null
  try {
    await iamService.deleteUser(user.id)
    users.value = users.value.filter((item) => item.id !== user.id)
  } catch (err: any) {
    error.value = err?.message || 'Error al eliminar usuario.'
  } finally {
    isDeleting.value = false
  }
}

onMounted(() => loadUsers())
</script>

<style scoped>
.main-container { padding-bottom: 4rem; }

.admin-layout { display: flex; flex-direction: column; gap: 1.5rem; }

.cards-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

/* Card Styling */
.card-header {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--color-background);
}

.header-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  background: rgba(230, 184, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-primary);
}

.header-icon.secondary { background: rgba(0, 35, 149, 0.1); color: var(--color-secondary); }
.header-icon .material-symbols-outlined { font-size: 28px; }

.card h2 { font-size: 1.1rem; margin: 0; color: var(--color-text-primary); }
.card p { margin: 0.25rem 0 0; font-size: 0.85rem; color: var(--color-text-secondary); }

/* Form */
.user-form { display: grid; gap: 1.25rem; }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: var(--font-size-xs); font-weight: 600; text-transform: uppercase; color: var(--color-text-secondary); }

input, select {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: var(--border-radius-md);
  border: 1px solid var(--color-border);
  font-size: var(--font-size-sm);
}

.input-with-icon { position: relative; display: flex; align-items: center; }
.icon-start { position: absolute; left: 0.85rem; font-size: 20px; color: var(--color-text-secondary); }
.input-with-icon input { padding-left: 2.75rem; }

.form-actions { display: flex; gap: 1rem; margin-top: 0.5rem; }
.form-actions .btn-primary { flex: 2; }
.form-actions .btn-outline { flex: 1; }

/* Roles Guide */
.roles-guide { list-style: none; padding: 0; margin: 0; display: grid; gap: 0.75rem; }
.roles-guide li {
  padding: 0.75rem 1rem;
  background: var(--color-background);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
.roles-guide strong { color: var(--color-secondary); font-size: 0.9rem; }
.roles-guide span { font-size: 0.8rem; color: var(--color-text-secondary); }

/* Table Section */
.table-header-filters {
  display: flex;
  gap: 1.5rem;
  padding: 1.25rem 1.5rem;
  background: var(--color-background);
  border-bottom: 1px solid var(--color-border);
}

.table-header-filters .filter-group { flex: 1; margin: 0; }

.table-wrapper { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; text-align: left; }
.data-table th {
  padding: 1rem;
  font-size: var(--font-size-xs);
  font-weight: 600;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  border-bottom: 1px solid var(--color-border);
}
.data-table td { padding: 1rem; border-bottom: 1px solid var(--color-border); font-size: var(--font-size-sm); vertical-align: middle; }

.user-cell { display: flex; align-items: center; gap: 0.75rem; }
.user-avatar { color: var(--color-border); font-size: 32px; }
.role-badge { 
  display: inline-block; padding: 0.2rem 0.6rem; background: var(--color-background); 
  border-radius: 99px; font-size: 0.75rem; font-weight: 600; color: var(--color-text-secondary);
}

.align-right { text-align: right; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.5rem; }

/* Modal */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal { width: 100%; max-width: 450px; padding: 0; }
.modal-header { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--color-border); display: flex; align-items: center; gap: 0.75rem; }
.modal-header .material-symbols-outlined { color: var(--color-primary); }
.modal-body { padding: 1.5rem; }
.modal-info { background: var(--color-background); padding: 1rem; border-radius: 8px; margin-bottom: 1.5rem; font-size: 0.9rem; }
.modal-actions { padding: 1.25rem 1.5rem; border-top: 1px solid var(--color-border); display: flex; justify-content: flex-end; gap: 1rem; }

/* Transitions */
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.mt-6 { margin-top: 1.5rem; }

@media (max-width: 1024px) { .cards-grid { grid-template-columns: 1fr; } }
@media (max-width: 768px) {
  .table-header-filters { flex-direction: column; }
  .action-buttons { flex-direction: column; align-items: stretch; }
}
</style>