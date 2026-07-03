<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { LogOut, Home, Briefcase, Users, Calendar, BarChart2, BookOpen, Bot, Settings, Bell } from 'lucide-vue-next'
import { notificationService } from '../../services/notification.service'

const router = useRouter()
const route = useRoute()
const showNotifications = ref(false)
const notifications = ref([])
const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

onMounted(async () => {
  try {
    if (authStore.user) {
      const data = await notificationService.getNotifications()
      notifications.value = data
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

const recruiterMenu = [
  { path: '/dashboard', name: 'Tổng quan', icon: Home },
  { path: '/jobs', name: 'Việc làm', icon: Briefcase },
  { path: '/candidates', name: 'Ứng viên', icon: Users },
  { path: '/interviews', name: 'Lịch phỏng vấn', icon: Calendar },
  { path: '/reports', name: 'Báo cáo', icon: BarChart2 },
  { path: '/question-bank', name: 'Kho câu hỏi', icon: BookOpen },
  { path: '/templates', name: 'Mẫu AI', icon: Bot },
  { path: '/settings', name: 'Cài đặt', icon: Settings }
]

const currentMenu = computed(() => recruiterMenu)
</script>

<template>
  <div class="main-layout">
    <aside class="sidebar">
      <div class="sidebar-header" style="cursor: pointer;" @click="router.push('/dashboard')">
        <div class="sidebar-logo">Interview AI</div>
        <div class="sidebar-subtitle">Enterprise Tier</div>
      </div>
      
      <nav class="sidebar-nav">
        <router-link 
          v-for="item in currentMenu" 
          :key="item.path"
          :to="item.path"
          class="nav-item"
          active-class="active"
        >
          <component :is="item.icon" size="18" />
          {{ item.name }}
        </router-link>
      </nav>
      
      <div style="padding: 16px; border-top: 1px solid var(--border)">
        <div style="display: flex; align-items: center; gap: 12px; padding: 10px; margin-bottom: 12px">
          <div style="width: 32px; height: 32px; border-radius: 50%; background-color: var(--primary); color: white; display: flex; align-items: center; justify-content: center; font-weight: bold">
            {{ authStore.user?.full_name ? authStore.user.full_name.charAt(0).toUpperCase() : 'R' }}
          </div>
          <div style="overflow: hidden">
            <div style="font-size: 14px; font-weight: 500; white-space: nowrap; text-overflow: ellipsis">
              {{ authStore.user?.full_name || 'Recruiter User' }}
            </div>
            <div style="font-size: 12px; color: var(--text-muted); white-space: nowrap; text-overflow: ellipsis; overflow: hidden">
              {{ authStore.user?.companies?.[0]?.name || 'HR Department' }}
            </div>
          </div>
        </div>
        <button class="btn btn-ghost" @click="handleLogout" style="width: 100%; justify-content: flex-start; color: var(--danger)">
          <LogOut size="18" style="margin-right: 8px;" />
          Đăng xuất
        </button>
      </div>
    </aside>
    
    <main class="content-area">
      <header class="navbar" style="display: flex; justify-content: flex-end; padding: 16px 32px; border-bottom: 1px solid var(--border); background-color: var(--surface)">
        <div style="position: relative; cursor: pointer" @click="showNotifications = !showNotifications">
          <Bell size="20" color="var(--text-secondary)" style="transition: color 0.2s" />
          <div v-if="unreadCount > 0" style="position: absolute; top: -2px; right: -2px; width: 14px; height: 14px; background-color: var(--danger); border-radius: 50%; border: 2px solid var(--surface); display: flex; align-items: center; justify-content: center; font-size: 9px; color: white; font-weight: bold;">
            {{ unreadCount }}
          </div>
          
          <!-- Notifications Dropdown -->
          <div v-if="showNotifications" style="position: absolute; top: 100%; right: 0; margin-top: 16px; width: 340px; background-color: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1); z-index: 50; overflow: hidden" @click.stop>
            <div style="padding: 16px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center">
              <span style="font-weight: 600; color: var(--text-main)">Thông báo</span>
              <span style="font-size: 12px; color: var(--primary); font-weight: 500; cursor: pointer" @click="handleMarkAllAsRead">Đánh dấu đã đọc</span>
            </div>
            
            <div style="max-height: 400px; overflow-y: auto">
              <div v-if="notifications.length === 0" style="padding: 32px; text-align: center; color: var(--text-muted)">
                Chưa có thông báo nào.
              </div>
              <div v-for="n in notifications" :key="n.id" 
                   style="padding: 16px; border-bottom: 1px solid var(--border); transition: background-color 0.2s; position: relative" 
                   :style="{ backgroundColor: n.is_read ? 'transparent' : 'rgba(37, 99, 235, 0.05)' }"
                   class="hover-bg">
                <div style="display: flex; justify-content: space-between; align-items: flex-start">
                  <div class="text-body" style="font-weight: 600; margin-bottom: 6px; color: var(--text-main)">{{ n.title }}</div>
                  <div v-if="!n.is_read" style="width: 8px; height: 8px; background-color: var(--primary); border-radius: 50%; flex-shrink: 0; cursor: pointer" title="Đánh dấu đã đọc" @click="handleMarkAsRead($event, n.id)"></div>
                </div>
                <div class="text-helper" style="color: var(--text-secondary); line-height: 1.5">
                  {{ n.content || n.message }}
                </div>
                <div style="font-size: 11px; color: var(--text-muted); margin-top: 8px">{{ formatTimeAgo(n.created_at) }}</div>
              </div>
            </div>

            <div style="padding: 12px; text-align: center; color: var(--primary); font-size: 13px; font-weight: 500; cursor: pointer; background-color: var(--surface-soft)">
              Xem tất cả thông báo
            </div>
          </div>
        </div>
      </header>
      <div class="page-content">
        <router-view />
      </div>
    </main>
  </div>
</template>
