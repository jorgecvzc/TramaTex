<template>
  <main class="page organizations-create-page">
    <div class="page-wrapper">
      <router-link to="/organizations" class="breadcrumb">← Organizations</router-link>
      <organization-form @submit="handleSubmit" />
    </div>
  </main>
</template>

<script setup>
import { useRouter } from 'vue-router';
import OrganizationForm from '@/components/party/OrganizationForm.vue';

const router = useRouter();

async function handleSubmit(organization) {
  // Navigate to the new organization's detail page
  await router.push(`/organizations/${organization.id}`);
}
</script>

<style scoped>
.page {
  width: 100%;
}

.page-wrapper {
  max-width: 600px;
  margin: 0 auto;
  padding: 1.5rem;
}

.breadcrumb {
  display: inline-block;
  color: var(--primary-color);
  text-decoration: none;
  margin-bottom: 1.5rem;
  font-weight: 500;
}

.breadcrumb:hover {
  text-decoration: underline;
}
</style>
