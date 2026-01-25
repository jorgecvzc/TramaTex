<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { onMounted } from 'vue'

const authStore = useAuthStore()

onMounted(async () => {
  await authStore.checkAuthStatus()
})
</script>

<template>
  <div id="app">
    <RouterView />
  </div>
</template>

<style>
@import '@/design-system/theme.css';
</style>
