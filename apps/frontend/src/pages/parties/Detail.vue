<template>
  <div class="party-page-master">
    <PartyDetail />
  </div>
</template>

<script setup>
/**
 * Detail.vue - Wrapper de la Entidad Party
 * Delega toda la estructura al componente PartyDetail que ya implementa BaseEntityPage.
 */
import PartyDetail from '@/components/party/PartyDetail.vue'
</script>

<style scoped>
.party-page-master {
  min-height: 100vh;
  background-color: var(--color-background);
}
</style>
