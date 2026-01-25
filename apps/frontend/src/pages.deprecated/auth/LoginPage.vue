<template>
  <div style="
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1.5rem;
    background-color: #FDFDFD;
  ">
    <!-- Main Card -->
    <div style="width: 100%; max-width: 400px;">
      <div style="
        background: white;
        border-radius: 0.5rem;
        box-shadow: 0 8px 30px rgba(0, 0, 0, 0.04);
        border: 1px solid #f3f4f6;
        padding: 2rem;
        display: flex;
        flex-direction: column;
        align-items: center;
      ">
        <!-- Logo -->
        <h1 style="
          font-family: 'Calibri', 'Candara', 'Segoe', 'Segoe UI', 'Optima', 'Arial', sans-serif;
          font-size: 2.25rem;
          font-weight: bold;
          color: #002395;
          font-style: italic;
          letter-spacing: -0.02em;
          margin-bottom: 2.5rem;
        ">
          TramaTex
        </h1>

        <!-- Form Content -->
        <div style="width: 100%; display: flex; flex-direction: column; gap: 1.5rem;">
          <LoginForm />
        </div>

        <!-- Footer -->
        <div style="margin-top: 2.5rem; text-align: center;">
          <p style="font-size: 0.75rem; color: #9ca3af; font-weight: 500;">
            Enterprise Resource Planning for Textile Micro-Enterprises
          </p>
          <div style="margin-top: 1rem; display: flex; justify-content: center; gap: 1rem; font-size: 0.625rem; text-transform: uppercase; letter-spacing: -0.05em; color: rgba(0, 35, 149, 0.6); font-weight: bold;">
            <a href="#" style="color: rgba(0, 35, 149, 0.6); text-decoration: none; transition: color 0.3s ease;" onmouseover="this.style.color='#002395'" onmouseout="this.style.color='rgba(0, 35, 149, 0.6)'">Terms of Service</a>
            <span style="color: #d1d5db;">|</span>
            <a href="#" style="color: rgba(0, 35, 149, 0.6); text-decoration: none; transition: color 0.3s ease;" onmouseover="this.style.color='#002395'" onmouseout="this.style.color='rgba(0, 35, 149, 0.6)'">Support</a>
          </div>
        </div>
      </div>

      <!-- Bottom Accent -->
      <div style="margin-top: 2rem; display: flex; justify-content: center;">
        <div style="width: 128px; height: 6px; background-color: #d1d5db; border-radius: 9999px;"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import LoginForm from '@/components/auth/LoginForm.vue'
</script>

<style scoped>
.calibri {
  font-family: 'Calibri', 'Candara', 'Segoe', 'Segoe UI', 'Optima', 'Arial', sans-serif;
}

.text-royal-blue {
  color: #002395;
}

.text-royal-blue-60 {
  color: rgba(0, 35, 149, 0.6);
}

.hover\:text-royal-blue:hover {
  color: #002395;
}
</style>
