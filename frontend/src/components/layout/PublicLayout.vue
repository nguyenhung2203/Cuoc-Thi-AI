<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppFooter from './AppFooter.vue'
import AppLogo from '../common/AppLogo.vue'
import { langStore } from '../../stores/lang.store'
import { authStore } from '../../stores/auth.store'
import { ChevronDown, User as UserIcon, LogOut, Globe } from 'lucide-vue-next'

const router = useRouter()
const showProfileMenu = ref(false)

onMounted(async () => {
  if (authStore.isAuthenticated && !authStore.user) {
    await authStore.init()
  }
})
</script>

<template>
  <div class="public-layout" style="min-height: 100vh; display: flex; flex-direction: column; background-color: var(--background);">
    <!-- Simple Header -->
    <header style="padding: 16px 32px; border-bottom: 1px solid var(--border); background-color: var(--surface); display: flex; justify-content: space-between; align-items: center;">
      <div style="cursor: pointer; display: flex; align-items: center;" @click="router.push(authStore.isAuthenticated ? '/home' : '/login')">
        <AppLogo size="md" />
      </div>
      <div class="flex items-center gap-4">
        <!-- Language Switcher in Public Layout -->
        <button
          @click="langStore.setLang(langStore.lang === 'vi' ? 'en' : 'vi')"
          class="px-3 py-1.5 rounded-full border border-gray-300 dark:border-gray-700 text-sm font-semibold hover:border-primary transition-colors flex items-center gap-2">
          <Globe size="16" /> {{ langStore.lang === 'vi' ? 'VI' : 'EN' }}
        </button>
        
        <template v-if="authStore.isAuthenticated">
          <div style="position: relative; cursor: pointer" @click="showProfileMenu = !showProfileMenu">
            <div style="display: flex; align-items: center; gap: 8px; padding: 4px 10px 4px 4px; border-radius: 24px; border: 1px solid var(--border); transition: border-color 0.2s" class="hover-border">
              <div style="width: 32px; height: 32px; flex-shrink: 0; border-radius: 50%; background-color: var(--primary); color: white; display: flex; align-items: center; justify-content: center; font-weight: bold">
                {{ authStore.user?.full_name ? authStore.user.full_name[0].toUpperCase() : 'C' }}
              </div>
              <div style="display: flex; flex-direction: column; max-width: 130px;">
                <span style="font-size: 12px; font-weight: 600; color: var(--text-main); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ authStore.user?.full_name || 'Ứng viên' }}</span>
              </div>
              <ChevronDown size="14" color="var(--text-muted)" />
            </div>

            <div v-if="showProfileMenu" class="dropdown-menu" style="position: absolute; top: 100%; right: 0; margin-top: 12px; width: 240px; background-color: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: 0 4px 20px rgba(0,0,0,0.08); z-index: 50;">
              <div style="padding: 16px; border-bottom: 1px solid var(--border)">
                <div style="font-weight: 600; color: var(--text-main)">{{ authStore.user?.full_name || 'Ứng viên' }}</div>
                <div style="font-size: 13px; color: var(--text-muted); margin-top: 4px">{{ authStore.user?.email || 'Chưa cập nhật' }}</div>
              </div>
              <div style="padding: 8px">
                <button class="dropdown-item" @click="router.push('/home')" style="width: 100%; display: flex; align-items: center; gap: 10px; padding: 10px; border-radius: 6px; text-align: left; background: transparent; border: none; cursor: pointer; color: var(--text-main); font-weight: 500;">
                  <UserIcon size="16" /> Về trang cá nhân
                </button>
                <button class="dropdown-item text-danger" @click="authStore.logout()" style="width: 100%; display: flex; align-items: center; gap: 10px; padding: 10px; border-radius: 6px; text-align: left; background: transparent; border: none; cursor: pointer; color: var(--danger); margin-top: 4px; font-weight: 500;">
                  <LogOut size="16" /> {{ langStore.t('nav', 'logout') }}
                </button>
              </div>
            </div>
          </div>
        </template>
        <template v-else>
          <button class="btn btn-primary" @click="router.push('/login')">{{ langStore.lang === 'vi' ? 'Đăng nhập / Đăng ký' : 'Login / Register' }}</button>
        </template>
      </div>
    </header>

    <!-- Main Content -->
    <main style="flex: 1; padding: 40px; max-width: 1152px; width: 100%; margin: 0 auto;">
      <router-view />
    </main>

    <!-- Footer -->
    <AppFooter />
  </div>
</template>

<style scoped>
.hover-border:hover {
  border-color: var(--primary) !important;
}
.dropdown-item {
  transition: all 0.2s ease;
}
.dropdown-item:hover {
  background-color: var(--surface-soft) !important;
}
.dropdown-item.text-danger:hover {
  background-color: rgba(220, 38, 38, 0.05) !important;
}
</style>
