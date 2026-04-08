<template>
  <BaseEntityPage>
    <template #header>
      <BasePageHeader 
        title="Nueva Entidad" 
        :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Crear nueva' }]"
        show-back
      >
        <template #icon>
          <span class="material-symbols-outlined">person_add</span>
        </template>
        <template #actions>
          <button @click="triggerSubmit" class="btn btn-primary btn-sm">
            <span class="material-symbols-outlined">save</span>
            <span>Crear Entidad</span>
          </button>
          <button @click="router.push('/parties')" class="btn btn-outline btn-sm">
            <span class="material-symbols-outlined">close</span>
            <span>Cancelar</span>
          </button>
        </template>
      </BasePageHeader>
    </template>

    <div class="form-wrapper">
      <PartyForm ref="partyFormRef" hide-actions hide-header @submit="handleSubmit" />
    </div>
  </BaseEntityPage>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import PartyForm from '@/components/party/PartyForm.vue'

const router = useRouter()
const partyFormRef = ref(null)

function triggerSubmit() {
  if (partyFormRef.value) {
    partyFormRef.value.submitForm()
  }
}

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