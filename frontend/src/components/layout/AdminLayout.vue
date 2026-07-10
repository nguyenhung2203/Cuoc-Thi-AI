<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { LayoutDashboard, Users, Building, Settings, LogOut, ChevronLeft, ChevronRight, Menu, Activity, ShieldCheck, Bell, Search, Sparkles } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()

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
  <div class="flex h-screen overflow-hidden bg-gray-50 dark:bg-slate-900 text-slate-800 dark:text-slate-200 font-sans transition-colors duration-300">
    
    <!-- Sidebar -->
    <aside :class="[isCollapsed ? 'w-16' : 'w-64', 'bg-white dark:bg-slate-800 border-r border-slate-200 dark:border-slate-700 flex flex-col transition-all duration-300 shadow-sm z-20 shrink-0']">
      <!-- Logo & Toggle -->
      <div class="p-4 border-b border-slate-100 dark:border-slate-700 flex items-center justify-between h-[64px] relative">
        <div v-if="!isCollapsed" class="cursor-pointer flex items-center gap-3 overflow-hidden" @click="router.push('/admin/dashboard')">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-indigo-600 to-purple-600 flex items-center justify-center text-white font-bold text-lg shadow-md shadow-indigo-500/20 shrink-0">
            A
          </div>
          <div class="flex flex-col">
            <span class="font-bold text-base whitespace-nowrap text-slate-800 dark:text-white tracking-tight">WeMake <span class="text-indigo-600 dark:text-indigo-400 font-extrabold">Admin</span></span>
            <span class="text-[10px] text-slate-400 font-medium tracking-wider uppercase">System Portal</span>
          </div>
        </div>

        <button v-if="!isCollapsed" @click="isCollapsed = true" class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-500 transition-colors shrink-0 ml-auto" title="Thu gọn menu">
          <ChevronLeft size="18" />
        </button>
        <button v-else @click="isCollapsed = false" class="mx-auto w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 flex items-center justify-center text-white shadow-md shadow-indigo-500/20 transition-all duration-300 transform hover:scale-105" title="Mở rộng menu Admin">
          <ChevronRight size="20" />
        </button>
      </div>

      <!-- Nav Items -->
      <nav class="flex-1 overflow-y-auto py-4 space-y-1.5 custom-scrollbar" :class="isCollapsed ? 'px-2' : 'px-3'">
        <div v-if="!isCollapsed" class="px-3 py-1 text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
          Quản trị viên
        </div>
        
        <router-link 
          v-for="item in navItems" 
          :key="item.path" 
          :to="item.path"
          class="flex items-center gap-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group relative"
          :class="[
            isCollapsed ? 'justify-center px-0' : 'px-3.5',
            isActive(item.path) 
              ? 'bg-gradient-to-r from-indigo-600 to-indigo-700 text-white shadow-md shadow-indigo-600/20 font-semibold' 
              : 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700/50 hover:text-indigo-600 dark:hover:text-indigo-400'
          ]"
          :title="isCollapsed ? item.name : ''"
        >
          <component :is="item.icon" :size="isCollapsed ? 20 : 19" class="shrink-0" :class="isActive(item.path) ? 'text-white' : 'text-slate-400 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors'" />
          <span v-if="!isCollapsed" class="whitespace-nowrap overflow-hidden text-ellipsis">{{ item.name }}</span>
        </router-link>
      </nav>

      <!-- Admin Profile & Logout -->
      <div class="border-t border-slate-200 dark:border-slate-700 bg-slate-50/60 dark:bg-slate-800/60" :class="isCollapsed ? 'p-2' : 'p-4'">
        <div class="flex items-center gap-3 mb-3 rounded-xl hover:bg-white dark:hover:bg-slate-700 transition-colors cursor-pointer border border-transparent hover:border-slate-200 dark:hover:border-slate-600" :class="isCollapsed ? 'justify-center p-1' : 'p-2'" :title="isCollapsed ? (authStore.user?.full_name || 'Super Admin') : ''">
          <div class="rounded-xl bg-gradient-to-br from-amber-500 to-orange-500 text-white flex items-center justify-center font-bold shadow-md shrink-0" :class="isCollapsed ? 'w-8 h-8 text-xs' : 'w-10 h-10 text-sm'">
            <ShieldCheck size="20" />
          </div>
          <div v-if="!isCollapsed" class="overflow-hidden flex-1">
            <div class="text-sm font-semibold whitespace-nowrap overflow-hidden text-ellipsis text-slate-800 dark:text-slate-200">
              {{ authStore.user?.full_name || 'Super Admin' }}
            </div>
            <div class="text-xs text-amber-600 dark:text-amber-400 font-medium whitespace-nowrap overflow-hidden text-ellipsis flex items-center gap-1">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span> Toàn quyền hệ thống
            </div>
          </div>
        </div>
        <button @click="handleLogout" :title="isCollapsed ? 'Đăng xuất' : ''" class="w-full flex items-center gap-2 py-2 text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-xl transition-colors" :class="isCollapsed ? 'justify-center px-0' : 'px-3'">
          <LogOut size="16" class="shrink-0" />
          <span v-if="!isCollapsed">Đăng xuất</span>
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col h-screen overflow-y-auto transition-colors duration-300 bg-gray-50 dark:bg-slate-900 custom-scrollbar">
      <!-- Top Header -->
      <header class="sticky top-0 h-[64px] flex justify-between items-center px-8 border-b border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-slate-800/80 backdrop-blur-md transition-colors duration-300 gap-4 shrink-0 z-40">
        <div class="flex items-center gap-4">
          <button @click="isCollapsed = !isCollapsed" class="md:hidden text-slate-500 hover:text-slate-800 dark:hover:text-white">
            <Menu size="22" />
          </button>
          
          <!-- Quick Search Placeholder / Title -->
          <div class="hidden sm:flex items-center gap-2 text-slate-400 text-sm bg-slate-100 dark:bg-slate-700/50 px-3 py-1.5 rounded-lg border border-slate-200/60 dark:border-slate-600/60">
            <Search size="15" />
            <span>Tìm kiếm nhà tuyển dụng, user, log... (Ctrl + K)</span>
          </div>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- Live System Badge -->
          <div class="hidden md:flex items-center gap-2 px-3 py-1 bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800/40 rounded-full">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-ping"></span>
            <span class="text-xs font-semibold text-emerald-700 dark:text-emerald-400">System Normal</span>
          </div>

          <!-- Notification Bell -->
          <div class="relative cursor-pointer" @click="showNotifications = !showNotifications">
            <div class="p-2 rounded-full hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors text-slate-500 dark:text-slate-400">
              <Bell size="20" />
            </div>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <div class="flex-1 p-6 lg:p-8">
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
.custom-scrollbar::-webkit-scrollbar {
  width: 5px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}
.dark .custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
