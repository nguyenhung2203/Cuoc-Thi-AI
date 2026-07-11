<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { LayoutDashboard, Users, Building, Settings, LogOut, ChevronLeft, ChevronRight, Menu, Activity, ShieldCheck, Bell, Search, Sparkles } from 'lucide-vue-next'
import { usePlatformStore } from '../../stores/platform.store'

const router = useRouter()
const route = useRoute()
const platformStore = usePlatformStore()

const isCollapsed = ref(false)
const showNotifications = ref(false)

const handleLogout = async () => {
  await authStore.logout()
  router.push('/admin/login')
}

const navItems = [
  { name: 'Tổng quan', path: '/admin/dashboard', icon: LayoutDashboard },
  { name: 'Quản lý Người dùng', path: '/admin/users', icon: Users },
  { name: 'Quản lý Công ty', path: '/admin/companies', icon: Building },
  { name: 'Nhật ký Hệ thống', path: '/admin/logs', icon: Activity },
  { name: 'Cài đặt Hệ thống', path: '/admin/settings', icon: Settings },
]

const isActive = (path) => {
  if (path === '/admin') return route.path === '/admin' || route.path === '/admin/dashboard'
  if (path === '/admin/dashboard') return route.path === '/admin' || route.path === '/admin/dashboard'
  return route.path.startsWith(path)
}

onMounted(() => {
  // Ensure consistent theme
  document.documentElement.classList.remove('dark')
  localStorage.setItem('theme', 'light')
})
</script>

<template>
  <div class="admin-shell">
    <!-- Sidebar -->
    <aside class="admin-sidebar" :class="{ 'is-collapsed': isCollapsed }">
      <!-- Logo & Toggle -->
      <div class="admin-brand">
        <div v-if="!isCollapsed" class="brand-id" @click="router.push('/admin/dashboard')">
          <div class="brand-mark">{{ platformStore.brandName ? platformStore.brandName[0] : 'V' }}</div>
          <div class="brand-text">
            <span class="brand-name">{{ platformStore.brandName }} <span class="brand-accent">Admin</span></span>
            <span class="brand-sub">{{ platformStore.brandBadge }} System Portal</span>
          </div>
        </div>
        <button class="collapse-btn" :class="{ 'mx-auto': isCollapsed }" @click="isCollapsed = !isCollapsed" :title="isCollapsed ? 'Mở rộng menu' : 'Thu gọn menu'">
          <ChevronLeft v-if="!isCollapsed" :size="18" />
          <ChevronRight v-else :size="18" />
        </button>
      </div>

      <!-- Nav Items -->
      <nav class="admin-nav custom-scrollbar">
        <div v-if="!isCollapsed" class="nav-group-label">Quản trị viên</div>
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="admin-nav-item"
          :class="{ 'is-active': isActive(item.path), 'is-collapsed': isCollapsed }"
          :title="isCollapsed ? item.name : ''"
        >
          <component :is="item.icon" :size="19" class="shrink-0" />
          <span v-if="!isCollapsed" class="truncate">{{ item.name }}</span>
        </router-link>
      </nav>

      <!-- Admin Profile & Logout -->
      <div class="admin-foot" :class="{ 'is-collapsed': isCollapsed }">
        <div class="admin-profile" :class="{ 'is-collapsed': isCollapsed }" :title="isCollapsed ? (authStore.user?.full_name || 'Super Admin') : ''">
          <div class="admin-avatar"><ShieldCheck :size="20" /></div>
          <div v-if="!isCollapsed" class="overflow-hidden flex-1">
            <div class="admin-name truncate">{{ authStore.user?.full_name || 'Super Admin' }}</div>
            <div class="admin-role"><span class="dot-live"></span> Toàn quyền hệ thống</div>
          </div>
        </div>
        <button @click="handleLogout" :title="isCollapsed ? 'Đăng xuất' : ''" class="logout-btn" :class="{ 'is-collapsed': isCollapsed }">
          <LogOut :size="16" class="shrink-0" />
          <span v-if="!isCollapsed">Đăng xuất</span>
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="admin-main custom-scrollbar">
      <header class="admin-header">
        <div class="flex items-center gap-4">
          <div class="admin-search">
            <Search :size="15" />
            <span>Tìm kiếm nhà tuyển dụng, user, log... (Ctrl + K)</span>
          </div>
        </div>

        <div class="flex items-center gap-4">
          <div class="system-badge">
            <span class="dot-ping"></span>
            <span>System Normal</span>
          </div>
          <div class="relative cursor-pointer" @click="showNotifications = !showNotifications">
            <div class="icon-btn"><Bell :size="20" /></div>
          </div>
        </div>
      </header>

      <div class="admin-content">
        <div class="max-w-7xl mx-auto">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.admin-shell { display: flex; height: 100vh; overflow: hidden; background: var(--background); color: var(--text-main); font-family: var(--sans); }

.admin-sidebar { width: 256px; background: var(--surface); border-right: 1px solid var(--border); display: flex; flex-direction: column; transition: width 0.25s ease; box-shadow: var(--shadow-sm); z-index: 20; flex-shrink: 0; }
.admin-sidebar.is-collapsed { width: 68px; }

.admin-brand { display: flex; align-items: center; justify-content: space-between; height: 64px; padding: 0 16px; border-bottom: 1px solid var(--border); }
.brand-id { display: flex; align-items: center; gap: 12px; cursor: pointer; overflow: hidden; }
.brand-mark { width: 36px; height: 36px; border-radius: var(--radius); background: var(--primary); color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 18px; flex-shrink: 0; }
.brand-name { font-weight: 700; font-size: 15px; color: var(--text-main); white-space: nowrap; }
.brand-accent { color: var(--primary); font-weight: 800; }
.brand-sub { display: block; font-size: 10px; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.06em; font-weight: 600; }
.collapse-btn { display: flex; align-items: center; justify-content: center; padding: 7px; border: none; background: transparent; border-radius: var(--radius); color: var(--text-muted); cursor: pointer; transition: all 0.2s ease; flex-shrink: 0; }
.collapse-btn:hover { background: var(--surface-soft); color: var(--text-main); }

.admin-nav { flex: 1; overflow-y: auto; padding: 16px 12px; display: flex; flex-direction: column; gap: 4px; }
.nav-group-label { padding: 4px 12px; font-size: 11px; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.06em; }
.admin-nav-item { display: flex; align-items: center; gap: 12px; padding: 10px 14px; border-radius: var(--radius); font-size: 14px; font-weight: 500; color: var(--text-secondary); transition: all 0.18s ease; }
.admin-nav-item.is-collapsed { justify-content: center; padding: 10px 0; }
.admin-nav-item :deep(svg) { color: var(--text-muted); transition: color 0.18s ease; flex-shrink: 0; }
.admin-nav-item:hover { background: var(--surface-soft); color: var(--primary); }
.admin-nav-item:hover :deep(svg) { color: var(--primary); }
.admin-nav-item.is-active { background: var(--primary); color: #fff; font-weight: 600; box-shadow: var(--shadow-sm); }
.admin-nav-item.is-active :deep(svg) { color: #fff; }

.admin-foot { border-top: 1px solid var(--border); background: var(--surface-soft); padding: 16px; }
.admin-foot.is-collapsed { padding: 8px; }
.admin-profile { display: flex; align-items: center; gap: 12px; padding: 8px; border-radius: var(--radius); margin-bottom: 8px; cursor: pointer; border: 1px solid transparent; transition: all 0.2s ease; }
.admin-profile.is-collapsed { justify-content: center; padding: 4px; }
.admin-profile:hover { background: var(--surface); border-color: var(--border); }
.admin-avatar { width: 40px; height: 40px; border-radius: var(--radius); background: var(--primary-light); color: var(--primary); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.admin-name { font-size: 14px; font-weight: 600; color: var(--text-main); }
.admin-role { display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--text-secondary); font-weight: 500; }
.dot-live { width: 6px; height: 6px; border-radius: 50%; background: var(--success); }
.logout-btn { width: 100%; display: flex; align-items: center; gap: 8px; padding: 9px 12px; border: none; background: transparent; border-radius: var(--radius); font-size: 14px; font-weight: 500; color: var(--danger); cursor: pointer; transition: background 0.2s ease; }
.logout-btn.is-collapsed { justify-content: center; padding: 9px 0; }
.logout-btn:hover { background: rgba(220,38,38,0.08); }

.admin-main { flex: 1; display: flex; flex-direction: column; height: 100vh; overflow-y: auto; background: var(--background); }
.admin-header { position: sticky; top: 0; height: 64px; display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 0 32px; border-bottom: 1px solid var(--border); background: rgba(255,255,255,0.85); backdrop-filter: blur(8px); z-index: 40; flex-shrink: 0; }
.icon-btn { display: flex; align-items: center; justify-content: center; padding: 8px; border-radius: 50%; color: var(--text-secondary); cursor: pointer; transition: background 0.2s ease; border: none; background: transparent; }
.icon-btn:hover { background: var(--surface-soft); color: var(--text-main); }
.admin-search { display: none; align-items: center; gap: 8px; padding: 7px 14px; border-radius: var(--radius); background: var(--surface-soft); border: 1px solid var(--border); color: var(--text-muted); font-size: 13px; }
.system-badge { display: none; align-items: center; gap: 8px; padding: 5px 12px; border-radius: var(--radius-full); background: rgba(22,163,74,0.1); border: 1px solid rgba(22,163,74,0.2); font-size: 12px; font-weight: 600; color: var(--success); }
.dot-ping { width: 8px; height: 8px; border-radius: 50%; background: var(--success); animation: pingDot 1.8s ease-out infinite; }
@keyframes pingDot { 0% { box-shadow: 0 0 0 0 rgba(22,163,74,0.4); } 100% { box-shadow: 0 0 0 6px rgba(22,163,74,0); } }
@media (min-width: 640px) { .admin-search { display: flex; } }
@media (min-width: 768px) { .system-badge { display: flex; } }

.admin-content { flex: 1; padding: 24px; }
@media (min-width: 1024px) { .admin-content { padding: 32px; } }

.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: var(--border); border-radius: 4px; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.15s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
@media (prefers-reduced-motion: reduce) { .dot-ping { animation: none !important; } }
</style>
