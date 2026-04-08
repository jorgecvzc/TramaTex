<template>
  <div class="catalog-page">
    <ProductList />
  </div>
</template>

<script setup>
import ProductList from '@/components/product/ProductList.vue'
</script>

<style scoped>
.catalog-page {
  min-height: 100vh;
  background-color: var(--color-background);
}
</style>
