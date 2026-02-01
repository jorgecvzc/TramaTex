<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Administración / Usuarios</p>
          <h1>Gestión de usuarios</h1>
          <p class="subtitle">Administración de usuarios, roles y accesos.</p>
        </div>
        <button class="btn btn-secondary" @click="loadUsers" :disabled="isLoading">
          {{ isLoading ? 'Cargando...' : 'Refrescar' }}
        </button>
      </header>

      <div v-if="!isAdmin" class="alert-warning">
        Solo el rol admin puede gestionar usuarios.
      </div>

      <div v-else class="cards-grid">
        <section class="card">
          <div class="card-header">
            <div>
              <h2>Alta de usuario</h2>
              <p>Crea cuentas con roles predefinidos.</p>
            </div>
            <span class="badge">Admin</span>
          </div>

          <form class="form-grid" @submit.prevent="createUser">
            <div>
              <label>Email</label>
              <input
                v-model.trim="newUser.email"
                type="email"
                placeholder="usuario@tramatex.local"
                required
              />
            </div>
            <div>
              <label>Contraseña</label>
              <input
                v-model="newUser.password"
                type="password"
                placeholder="Mínimo 8 caracteres"
                required
                minlength="8"
              />
            </div>
            <div>
              <label>Rol</label>
              <select v-model="newUser.role">
                <option value="admin">admin</option>
                <option value="commercial">commercial</option>
                <option value="designer">designer</option>
                <option value="workshop">workshop</option>
              </select>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" :disabled="isCreating">
                {{ isCreating ? 'Creando...' : 'Crear usuario' }}
              </button>
              <button type="button" class="btn btn-secondary" @click="resetNewUser">
                Limpiar
              </button>
            </div>
            <p class="helper-text">
              Las contraseñas se almacenan con hash bcrypt. El acceso se controla por rol.
            </p>
          </form>
        </section>

        <section class="card">
          <h2>Roles disponibles</h2>
          <p class="card-subtitle">Referencia rápida de permisos (MVP).</p>
          <ul class="roles-list">
            <li>
              <strong>admin</strong>
              <span>Acceso total y administración.</span>
            </li>
            <li>
              <strong>commercial</strong>
              <span>Gestión comercial y clientes.</span>
            </li>
            <li>
              <strong>designer</strong>
              <span>Diseño de pedidos.</span>
            </li>
            <li>
              <strong>workshop</strong>
              <span>Operaciones de taller.</span>
            </li>
          </ul>
        </section>
      </div>

      <section v-if="isAdmin" class="card card-full">
        <div class="filters">
          <div>
            <label>Buscar por nombre o email</label>
            <input
              v-model.trim="search"
              type="text"
              placeholder="Buscar por email"
            />
          </div>
          <div>
            <label>Filtrar por rol</label>
            <select v-model="roleFilter">
              <option value="">Todos</option>
              <option value="admin">admin</option>
              <option value="commercial">commercial</option>
              <option value="designer">designer</option>
              <option value="workshop">workshop</option>
            </select>
          </div>
        </div>

        <div v-if="error" class="alert-error">
          {{ error }}
        </div>

        <div class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Email</th>
                <th>Rol</th>
                <th class="align-right">Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in filteredUsers" :key="user.id">
                <td>{{ user.email }}</td>
                <td>
                  <span class="role-pill">{{ user.role }}</span>
                </td>
                <td class="align-right">
                  <div class="action-buttons">
                    <button class="btn btn-outline" @click="openRoleModal(user)">
                      Asignar rol
                    </button>
                    <button class="btn btn-danger" @click="confirmDelete(user)">
                      Eliminar
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="!isLoading && filteredUsers.length === 0">
                <td colspan="3" class="empty-state">No hay usuarios para mostrar.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <div v-if="showModal" class="modal-backdrop">
      <div class="modal">
        <h2>Asignar rol</h2>
        <p>
          Usuario: <strong>{{ selectedUser?.email }}</strong>
        </p>

        <label>Rol</label>
        <select v-model="selectedRole">
          <option value="admin">admin</option>
          <option value="commercial">commercial</option>
          <option value="designer">designer</option>
          <option value="workshop">workshop</option>
        </select>

        <div class="modal-actions">
          <button class="btn btn-secondary" @click="closeModal">Cancelar</button>
          <button class="btn btn-primary" @click="saveRole" :disabled="isSaving">
            {{ isSaving ? 'Guardando...' : 'Guardar' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
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

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.page-header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.page-header h1 {
  color: #1b3a6b;
  margin: 0.25rem 0 0;
}

.breadcrumb {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin: 0;
}

.subtitle {
  color: #64748b;
  margin: 0.5rem 0 0;
  font-size: 0.95rem;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
}

.card-full {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.card h2 {
  color: #1b3a6b;
  margin-bottom: 0.5rem;
}

.card p {
  color: #64748b;
  font-size: 0.9rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: center;
  margin-bottom: 1rem;
}

.badge {
  background-color: #002395;
  color: #ffffff;
  padding: 0.35rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.form-grid {
  display: grid;
  gap: 1rem;
}

label {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin-bottom: 0.4rem;
}

input,
select {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.6rem 0.8rem;
  font-size: 0.9rem;
  color: #1e293b;
}

input:focus,
select:focus {
  outline: none;
  border-color: #002395;
  box-shadow: 0 0 0 3px rgba(0, 35, 149, 0.12);
}

.form-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.helper-text {
  font-size: 0.8rem;
  color: #64748b;
  margin: 0;
}

.roles-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 0.75rem;
}

.roles-list li {
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 0.75rem 1rem;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.roles-list strong {
  color: #002395;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.filters > div {
  flex: 1;
  min-width: 220px;
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

thead {
  background: #f8fafc;
  color: #64748b;
}

th,
td {
  padding: 0.85rem 1rem;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}

.align-right {
  text-align: right;
}

.role-pill {
  display: inline-block;
  background-color: #e2e8f0;
  color: #1e293b;
  font-size: 0.75rem;
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-weight: 600;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #e6b800;
  color: #1e293b;
  font-weight: 700;
}

.btn-primary:hover {
  background: #d6aa00;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-outline {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-danger {
  background: #fee2e2;
  color: #991b1b;
  border: 1px solid #fecaca;
}

.alert-warning {
  background: #fef3c7;
  border: 1px solid #f59e0b;
  color: #92400e;
  padding: 1rem;
  border-radius: 8px;
}

.alert-error {
  background: #fee2e2;
  border: 1px solid #ef4444;
  color: #991b1b;
  padding: 0.8rem 1rem;
  border-radius: 8px;
}

.empty-state {
  text-align: center;
  color: #64748b;
  padding: 1.5rem;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  z-index: 50;
}

.modal {
  background: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.2);
  display: grid;
  gap: 1rem;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

@media (max-width: 768px) {
  .dashboard-content {
    padding: 1.5rem;
  }

  .action-buttons {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
