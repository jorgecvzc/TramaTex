<template>
  <div class="auth-layout">
    <div class="auth-container">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.auth-layout {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1b3a6b 0%, #2c5aa0 100%);
}

.auth-container {
  width: 100%;
  max-width: 500px;
  padding: 2rem;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

@media (max-width: 640px) {
  .auth-layout {
    padding: 1rem;
  }

  .auth-container {
    max-width: 100%;
  }
}
</style>
