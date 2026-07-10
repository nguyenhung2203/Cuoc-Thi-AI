<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { LogOut, Home, Briefcase, Users, Calendar, BarChart2, BookOpen, Bot, Settings, Bell, Scale, FileText } from 'lucide-vue-next'
import { notificationService } from '../../services/notification.service'

const router = useRouter()
const route = useRoute()
const showNotifications = ref(false)
const notifications = ref([])
const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

onMounted(async () => {
  // Ensure light mode is default
  document.documentElement.classList.remove('dark')
  localStorage.setItem('theme', 'light')

  try {
    if (authStore.user) {
      const data = await notificationService.getNotifications()
      notifications.value = data.notifications || []
    }
  } catch (error) {
    console.error('Failed to load notifications:', error)
  }
})

const handleMarkAllAsRead = async (e) => {
  e.stopPropagation();
  try {
    if (authStore.user) {
      await notificationService.markAllAsRead()
      notifications.value = notifications.value.map(n => ({ ...n, is_read: true }))
    }
  } catch (error) {
    console.error('Failed to mark all as read:', error)
  }
}

const handleMarkAsRead = async (e, id) => {
  e.stopPropagation();
  try {
    if (authStore.user) {
      await notificationService.markAsRead(id)
      const notif = notifications.value.find(n => n.id === id)
      if (notif) notif.is_read = true
    }
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

const formatTimeAgo = (isoStr) => {
  if (!isoStr) return ''
  const diff = new Date() - new Date(isoStr)
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return 'Vừa xong'
  if (minutes < 60) return `${minutes} phút trước`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} giờ trước`
  return `${Math.floor(hours / 24)} ngày trước`
}


const handleLogout = async () => {
  await authStore.logout()
}

const adminMenu = [
  { path: '/admin/dashboard', name: 'Tổng quan', icon: Home },
  { path: '/admin/users', name: 'Người dùng', icon: Users },
  { path: '/admin/companies', name: 'Công ty', icon: Briefcase },
  { path: '/admin/prompts', name: 'Mẫu AI', icon: Bot },
  { path: '/admin/audit-logs', name: 'Nhật ký HT', icon: FileText }
]

const currentMenu = computed(() => adminMenu)
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-gray-50 dark:bg-slate-900 text-slate-800 dark:text-slate-200 font-sans transition-colors duration-300">
    <!-- Sidebar -->
    <aside class="w-64 bg-white dark:bg-slate-800 border-r border-slate-200 dark:border-slate-700 flex flex-col transition-colors duration-300 shadow-sm z-20">
      <div class="p-6 cursor-pointer border-b border-slate-100 dark:border-slate-700 flex items-center justify-center h-[88px]" @click="router.push('/dashboard')">
        <img src="/images/logo.png" alt="Logo" class="h-14 object-contain" />
      </div>
      
      <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-1 scrollbar-thin">
        <router-link 
          v-for="item in currentMenu" 
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/50 hover:text-indigo-600 dark:hover:text-indigo-400"
          active-class="bg-indigo-50 dark:bg-indigo-500/10 text-indigo-700 dark:text-indigo-400"
        >
          <component :is="item.icon" size="18" />
          {{ item.name }}
        </router-link>
      </nav>
      
      <div class="p-4 border-t border-slate-200 dark:border-slate-700 bg-slate-50/50 dark:bg-slate-800/50">
        <div class="flex items-center gap-3 mb-3 p-2 rounded-lg hover:bg-white dark:hover:bg-slate-700 transition-colors cursor-pointer border border-transparent hover:border-slate-200 dark:hover:border-slate-600">
          <div class="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-blue-500 text-white flex items-center justify-center font-bold shadow-md">
            {{ authStore.user?.full_name ? authStore.user.full_name.charAt(0).toUpperCase() : 'R' }}
          </div>
          <div class="overflow-hidden flex-1">
            <div class="text-sm font-semibold whitespace-nowrap overflow-hidden text-ellipsis text-slate-800 dark:text-slate-200">
              {{ authStore.user?.full_name || 'Admin User' }}
            </div>
            <div class="text-xs text-slate-500 dark:text-slate-400 whitespace-nowrap overflow-hidden text-ellipsis">
              System Administrator
            </div>
          </div>
        </div>
        <button @click="handleLogout" class="w-full flex items-center gap-2 px-3 py-2 text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-lg transition-colors">
          <LogOut size="18" />
          Đăng xuất
        </button>
      </div>
    </aside>
    
    <!-- Main Content -->
    <main class="flex-1 flex flex-col h-screen overflow-y-auto transition-colors duration-300 bg-slate-50 dark:bg-slate-900 scrollbar-thin">
      <header class="sticky top-0 h-[72px] flex justify-end items-center px-8 border-b border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-slate-800/80 backdrop-blur-md transition-colors duration-300 gap-4 shrink-0 z-40">

        <!-- Notifications -->
        <div class="relative cursor-pointer" @click="showNotifications = !showNotifications">
          <div class="p-2 rounded-full hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors text-slate-500 dark:text-slate-400">
            <Bell size="20" />
          </div>
          <div v-if="unreadCount > 0" class="absolute top-1 right-1 w-4 h-4 bg-red-500 rounded-full border-2 border-white dark:border-slate-800 flex items-center justify-center text-[9px] text-white font-bold">
            {{ unreadCount }}
          </div>
          
          <!-- Notifications Dropdown -->
          <div v-if="showNotifications" class="absolute top-full right-0 mt-2 w-80 sm:w-96 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-xl z-50 overflow-hidden" @click.stop>
            <div class="p-4 border-b border-slate-100 dark:border-slate-700 flex justify-between items-center bg-slate-50/50 dark:bg-slate-800/50">
              <span class="font-semibold text-slate-800 dark:text-slate-200">Thông báo</span>
              <span class="text-xs text-indigo-600 dark:text-indigo-400 font-medium cursor-pointer hover:underline" @click="handleMarkAllAsRead">Đánh dấu đã đọc</span>
            </div>
            
            <div class="max-h-[400px] overflow-y-auto">
              <div v-if="notifications.length === 0" class="p-8 text-center text-slate-500 dark:text-slate-400 text-sm">
                Chưa có thông báo nào.
              </div>
              <div v-for="n in notifications" :key="n.id" 
                   class="p-4 border-b border-slate-100 dark:border-slate-700 transition-colors relative hover:bg-slate-50 dark:hover:bg-slate-700/50" 
                   :class="{ 'bg-indigo-50/50 dark:bg-indigo-500/5': !n.is_read }">
                <div class="flex justify-between items-start mb-1">
                  <div class="text-sm font-semibold text-slate-800 dark:text-slate-200">{{ n.title }}</div>
                  <div v-if="!n.is_read" class="w-2 h-2 bg-indigo-500 rounded-full shrink-0 cursor-pointer" title="Đánh dấu đã đọc" @click="handleMarkAsRead($event, n.id)"></div>
                </div>
                <div class="text-xs text-slate-600 dark:text-slate-400 leading-relaxed mb-2">
                  {{ n.content || n.message }}
                </div>
                <div class="text-[11px] text-slate-400 dark:text-slate-500">{{ formatTimeAgo(n.created_at) }}</div>
              </div>
            </div>

            <div class="p-3 text-center text-indigo-600 dark:text-indigo-400 text-xs font-medium cursor-pointer bg-slate-50 hover:bg-slate-100 dark:bg-slate-800/80 dark:hover:bg-slate-700 transition-colors">
              Xem tất cả thông báo
            </div>
          </div>
        </div>
      </header>
      
      <!-- Page Content -->
      <div class="flex-1 p-6 lg:p-8">
        <div class="max-w-7xl mx-auto">
          <router-view />
        </div>
      </div>
    </main>
  </div>
</template>
