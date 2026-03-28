<template>
  <Navbar />
  <div class="main-container">
    <PageHeader 
      title="Nueva Entidad" 
      :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Crear nueva' }]"
    >
      <template #actions>
        <button @click="router.push('/parties')" class="btn btn-outline">
          <span class="material-symbols-outlined">list_alt</span>
          <span>Ir al catálogo</span>
        </button>
      </template>
    </PageHeader>

    <div class="form-wrapper">
      <PartyForm @submit="handleSubmit" />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import PageHeader from '@/components/layout/PageHeader.vue'
import PartyForm from '@/components/party/PartyForm.vue'

const router = useRouter()

async function handleSubmit(party) {
  await router.push(`/parties/${party.id}`)
}
</script>

<style scoped>
.main-container {
  padding-bottom: 3rem;
}

.form-wrapper {
  width: 100%;
}
</style>