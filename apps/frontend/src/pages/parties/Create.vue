<template>
  <BaseEntityPage>
    <template #header>
      <BasePageHeader 
        title="Nueva Entidad" 
        :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Crear nueva' }]"
        show-back
      >
        <template #icon>
          <UserPlus :size="28" />
        </template>
        <template #actions>
          <button @click="triggerSubmit" class="btn btn-primary btn-sm">
            <Save :size="16" />
            <span>Crear Entidad</span>
          </button>
          <button @click="router.push('/parties')" class="btn btn-outline btn-sm">
            <X :size="16" />
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
import { UserPlus, Save, X } from 'lucide-vue-next'
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