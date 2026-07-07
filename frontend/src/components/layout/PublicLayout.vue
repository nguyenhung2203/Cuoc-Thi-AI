<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
  <div class="public-layout" style="min-height: 100vh; display: flex; flex-direction: column; background-color: var(--background);">
    <!-- Simple Header -->
    <header style="padding: 16px 32px; border-bottom: 1px solid var(--border); background-color: var(--surface); display: flex; justify-content: space-between; align-items: center;">
      <div style="cursor: pointer; display: flex; align-items: center;" @click="router.push('/login')">
        <img src="/images/logo.png" alt="Logo" style="height: 56px; object-fit: contain;" />
      </div>
      <div>
        <button class="btn btn-ghost" @click="router.push('/login')">Đăng nhập</button>
      </div>
    </header>

    <!-- Main Content -->
    <main style="flex: 1; padding: 40px; max-width: 1000px; width: 100%; margin: 0 auto;">
      <router-view />
    </main>

    <!-- Footer -->
    <footer style="padding: 24px; text-align: center; color: var(--text-muted); font-size: 14px; border-top: 1px solid var(--border);">
      © 2026 Interview AI. All rights reserved.
    </footer>
  </div>
</template>
