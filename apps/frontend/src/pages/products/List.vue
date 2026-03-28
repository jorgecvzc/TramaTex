<template>
  <div class="catalog-page">
    <Navbar style="z-index: 2000;" />
    <ProductList />
  </div>
</template>

<script setup>
import Navbar from '@/components/layout/Navbar.vue'
import ProductList from '@/components/product/ProductList.vue'
</script>

<style scoped>
.catalog-page {
  min-height: 100vh;
  background-color: var(--color-background);
}
</style>
