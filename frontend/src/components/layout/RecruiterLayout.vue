<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { useNotificationStore } from '../../stores/notification.store'
import { LogOut, Home, Briefcase, Users, Calendar, BarChart2, BookOpen, Bot, Settings, Bell, Scale, ChevronLeft, ChevronRight } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const notificationStore = useNotificationStore()
const showNotifications = ref(false)
const isCollapsed = ref(false)

const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)

watch(() => route.path, (newPath) => {
  if (newPath.includes('/recruiter-room')) {
    isCollapsed.value = true
  } else {
    isCollapsed.value = false
  }
}, { immediate: true })

onMounted(async () => {
  // Ensure light mode is default
  document.documentElement.classList.remove('dark')
  localStorage.setItem('theme', 'light')

  if (authStore.user) {
    await notificationStore.fetch()
  }
})

const handleMarkAllAsRead = async (e) => {
  e.stopPropagation();
  await notificationStore.markAllRead()
}

const handleMarkAsRead = async (e, id) => {
  e.stopPropagation();
  await notificationStore.markRead(id)
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

const recruiterMenu = [
  { path: '/dashboard', name: 'Tổng quan', icon: Home },
  { path: '/jobs', name: 'Việc làm', icon: Briefcase },
  { path: '/candidates', name: 'Ứng viên', icon: Users },
  { path: '/interviews', name: 'Lịch phỏng vấn', icon: Calendar },
  { path: '/reports', name: 'Báo cáo', icon: BarChart2 },
  { path: '/question-bank', name: 'Kho câu hỏi', icon: BookOpen },
  { path: '/rubrics', name: 'Tiêu chí (Rubric)', icon: Scale },
  { path: '/templates', name: 'Mẫu AI', icon: Bot },
  { path: '/settings', name: 'Cài đặt', icon: Settings }
]

const currentMenu = computed(() => recruiterMenu)
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-gray-50 dark:bg-slate-900 text-slate-800 dark:text-slate-200 font-sans transition-colors duration-300">
    <!-- Sidebar -->
    <aside :class="[isCollapsed ? 'w-14' : 'w-46', 'bg-white dark:bg-slate-800 border-r border-slate-200 dark:border-slate-700 flex flex-col transition-all duration-300 shadow-sm z-20 shrink-0']">
      <div class="p-3 border-b border-slate-100 dark:border-slate-700 flex items-center justify-between h-[64px] relative">
        <div v-if="!isCollapsed" class="cursor-pointer flex items-center overflow-hidden" @click="router.push('/dashboard')">
          <img src="/images/logo.png" alt="Logo" class="h-10 object-contain" />
        </div>
        <button @click="isCollapsed = !isCollapsed" class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-500 transition-colors mx-auto shrink-0" :class="{ 'ml-auto mr-0': !isCollapsed }" :title="isCollapsed ? 'Mở rộng menu' : 'Thu gọn menu'">
          <ChevronLeft v-if="!isCollapsed" size="18" />
          <ChevronRight v-else size="18" />
        </button>
      </div>
      
      <nav class="flex-1 overflow-y-auto py-3 space-y-1 scrollbar-thin" :class="isCollapsed ? 'px-1.5' : 'px-3'">
        <router-link 
          v-for="item in currentMenu" 
          :key="item.path"
          :to="item.path"
          :title="isCollapsed ? item.name : ''"
          class="flex items-center gap-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/50 hover:text-indigo-600 dark:hover:text-indigo-400"
          :class="isCollapsed ? 'justify-center px-0' : 'px-3'"
          active-class="bg-indigo-50 dark:bg-indigo-500/10 text-indigo-700 dark:text-indigo-400"
        >
          <component :is="item.icon" :size="isCollapsed ? 18 : 20" class="shrink-0" />
          <span v-if="!isCollapsed" class="whitespace-nowrap overflow-hidden text-ellipsis">{{ item.name }}</span>
        </router-link>
      </nav>
      
      <div class="border-t border-slate-200 dark:border-slate-700 bg-slate-50/50 dark:bg-slate-800/50" :class="isCollapsed ? 'p-2' : 'p-4'">
        <div class="flex items-center gap-3 mb-2 rounded-lg hover:bg-white dark:hover:bg-slate-700 transition-colors cursor-pointer border border-transparent hover:border-slate-200 dark:hover:border-slate-600" :class="isCollapsed ? 'justify-center p-1' : 'p-2'" :title="isCollapsed ? (authStore.user?.full_name || 'Recruiter User') : ''">
          <div class="rounded-full bg-gradient-to-br from-indigo-500 to-blue-500 text-white flex items-center justify-center font-bold shadow-md shrink-0" :class="isCollapsed ? 'w-8 h-8 text-xs' : 'w-10 h-10 text-sm'">
            {{ authStore.user?.full_name ? authStore.user.full_name.charAt(0).toUpperCase() : 'R' }}
          </div>
          <div v-if="!isCollapsed" class="overflow-hidden flex-1">
            <div class="text-sm font-semibold whitespace-nowrap overflow-hidden text-ellipsis text-slate-800 dark:text-slate-200">
              {{ authStore.user?.full_name || 'Recruiter User' }}
            </div>
            <div class="text-xs text-slate-500 dark:text-slate-400 whitespace-nowrap overflow-hidden text-ellipsis">
              {{ authStore.user?.companies?.[0]?.name || 'HR Department' }}
            </div>
          </div>
        </div>
        <button @click="handleLogout" :title="isCollapsed ? 'Đăng xuất' : ''" class="w-full flex items-center gap-2 py-1.5 text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-lg transition-colors" :class="isCollapsed ? 'justify-center px-0' : 'px-3'">
          <LogOut size="16" class="shrink-0" />
          <span v-if="!isCollapsed">Đăng xuất</span>
        </button>
      </div>
    </aside>
    
    <!-- Main Content -->
    <main class="flex-1 flex flex-col h-screen overflow-y-auto transition-colors duration-300 bg-slate-50 dark:bg-slate-900 scrollbar-thin">
      <header v-if="!route.path.includes('/recruiter-room')" class="sticky top-0 h-[72px] flex justify-end items-center px-8 border-b border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-slate-800/80 backdrop-blur-md transition-colors duration-300 gap-4 shrink-0 z-40">

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
      <div :class="route.path.includes('/recruiter-room') ? 'flex-1 p-0 overflow-hidden flex flex-col' : 'flex-1 p-6 lg:p-8'">
        <div :class="route.path.includes('/recruiter-room') ? 'w-full h-full flex-1 flex flex-col' : 'max-w-7xl mx-auto'">
          <router-view />
        </div>
      </div>
    </main>
  </div>
</template>
