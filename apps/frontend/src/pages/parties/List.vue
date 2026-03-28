<template>
  <div class="parties-page">
    <Navbar />
    <PartyList />
  </div>
</template>

<script setup>
import Navbar from '@/components/layout/Navbar.vue'
import PartyList from '@/components/party/PartyList.vue'
</script>

<style scoped>
.parties-page {
  min-height: 100vh;
  background-color: var(--color-background);
}
</style>