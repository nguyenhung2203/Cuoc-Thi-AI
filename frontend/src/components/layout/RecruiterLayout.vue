<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { LogOut, Home, Briefcase, Users, Calendar, BarChart2, BookOpen, Bot, Settings, FileText, CheckSquare, Award, User as UserIcon, Bell } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const showNotifications = ref(false)

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

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('role')
  router.push('/login')
}
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
            R
          </div>
          <div style="overflow: hidden">
            <div style="font-size: 14px; font-weight: 500; white-space: nowrap; text-overflow: ellipsis">
              Recruiter User
            </div>
            <div style="font-size: 12px; color: var(--text-muted)">
              HR Department
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
          <div style="position: absolute; top: -2px; right: -2px; width: 8px; height: 8px; background-color: var(--danger); border-radius: 50%; border: 2px solid var(--surface)"></div>
          
          <!-- Notifications Dropdown -->
          <div v-if="showNotifications" style="position: absolute; top: 100%; right: 0; margin-top: 16px; width: 340px; background-color: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1); z-index: 50; overflow: hidden">
            <div style="padding: 16px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center">
              <span style="font-weight: 600; color: var(--text-main)">Thông báo</span>
              <span style="font-size: 12px; color: var(--primary); font-weight: 500">Đánh dấu đã đọc</span>
            </div>
            
            <div style="padding: 16px; border-bottom: 1px solid var(--border); background-color: rgba(37, 99, 235, 0.05); transition: background-color 0.2s" class="hover-bg">
              <div class="text-body" style="font-weight: 600; margin-bottom: 6px; color: var(--text-main)">Có ứng viên mới 🚀</div>
              <div class="text-helper" style="color: var(--text-secondary); line-height: 1.5">
                <strong>Nguyễn Văn A</strong> vừa nộp hồ sơ ứng tuyển vào vị trí <strong>Backend Engineer</strong>. Điểm phù hợp AI đánh giá: 85%.
              </div>
              <div style="font-size: 11px; color: var(--text-muted); margin-top: 8px">5 phút trước</div>
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
