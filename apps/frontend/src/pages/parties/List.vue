<template>
  <div class="parties-page">
    <PartyList />
  </div>
</template>

<script setup>
import PartyList from '@/components/party/PartyList.vue'
</script>

<style scoped>
.parties-page {
  min-height: 100vh;
  background-color: var(--color-background);
}
</style>