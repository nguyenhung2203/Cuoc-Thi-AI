<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Home, Calendar, Bot, Award, User as UserIcon, FileText, Settings, Bell, LogOut, ChevronDown } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const showNotifications = ref(false)
const showProfileMenu = ref(false)

const candidateMenu = [
  { path: '/home', name: 'Tổng quan', icon: Home },
  { path: '/my-interviews', name: 'Phỏng vấn', icon: Calendar },
  { path: '/mock-setup', name: 'Luyện tập AI', icon: Bot },
  { path: '/mock-results', name: 'Kết quả', icon: Award }
]

const profileMenu = [
  { path: '/profile', name: 'Hồ sơ của tôi', icon: UserIcon },
  { path: '/cv', name: 'Quản lý CV', icon: FileText },
  { path: '/settings', name: 'Cài đặt', icon: Settings }
]

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('role')
  router.push('/login')
}
</script>

<template>
  <div class="candidate-layout">
    <!-- Top Navbar -->
    <header class="top-navbar">
      <div class="nav-container">
        <!-- Logo -->
        <div class="nav-brand" style="cursor: pointer;" @click="router.push('/home')">
          <div class="brand-logo">Interview AI</div>
          <div class="brand-subtitle">Candidate Portal</div>
        </div>

        <!-- Center Menu -->
        <nav class="nav-menu">
          <router-link 
            v-for="item in candidateMenu" 
            :key="item.path"
            :to="item.path"
            class="nav-link"
            active-class="active"
          >
            <component :is="item.icon" size="18" />
            {{ item.name }}
          </router-link>
        </nav>

        <!-- Right Actions -->
        <div class="nav-actions">
          <!-- Notification Bell -->
          <div style="position: relative; cursor: pointer; margin-right: 16px" @click="showNotifications = !showNotifications; showProfileMenu = false">
            <div style="padding: 8px; border-radius: 50%; background-color: var(--surface-soft); transition: background-color 0.2s" class="hover-circle">
              <Bell size="20" color="var(--text-secondary)" />
              <div style="position: absolute; top: 6px; right: 8px; width: 8px; height: 8px; background-color: var(--danger); border-radius: 50%; border: 2px solid var(--surface)"></div>
            </div>
            
            <!-- Notifications Dropdown -->
            <div v-if="showNotifications" class="dropdown-menu">
              <div style="padding: 16px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center">
                <span style="font-weight: 600; color: var(--text-main)">Thông báo</span>
                <span style="font-size: 12px; color: var(--primary); font-weight: 500">Đánh dấu đã đọc</span>
              </div>
              <div style="padding: 16px; border-bottom: 1px solid var(--border); background-color: rgba(37, 99, 235, 0.05); transition: background-color 0.2s" class="hover-bg">
                <div class="text-body" style="font-weight: 600; margin-bottom: 6px; color: var(--text-main)">Lịch phỏng vấn mới! 🎉</div>
                <div class="text-helper" style="color: var(--text-secondary); line-height: 1.5">
                  Nhà tuyển dụng vừa lên lịch phỏng vấn với bạn cho vị trí <strong>Frontend Developer</strong> vào 10:00 sáng ngày mai. Hãy kiểm tra mục "Phỏng vấn của tôi".
                </div>
                <div style="font-size: 11px; color: var(--text-muted); margin-top: 8px">Vừa xong</div>
              </div>
              <div style="padding: 12px; text-align: center; color: var(--primary); font-size: 13px; font-weight: 500; cursor: pointer; background-color: var(--surface-soft)">
                Xem tất cả thông báo
              </div>
            </div>
          </div>

          <!-- User Profile Dropdown -->
          <div style="position: relative; cursor: pointer" @click="showProfileMenu = !showProfileMenu; showNotifications = false">
            <div style="display: flex; align-items: center; gap: 12px; padding: 4px 8px; border-radius: 24px; border: 1px solid var(--border); transition: border-color 0.2s" class="hover-border">
              <div style="width: 32px; height: 32px; border-radius: 50%; background-color: var(--primary); color: white; display: flex; align-items: center; justify-content: center; font-weight: bold">
                C
              </div>
              <div style="display: flex; flex-direction: column">
                <span style="font-size: 13px; font-weight: 600; color: var(--text-main)">Candidate User</span>
              </div>
              <ChevronDown size="16" color="var(--text-muted)" style="margin-right: 4px" />
            </div>

            <div v-if="showProfileMenu" class="dropdown-menu" style="width: 240px">
              <div style="padding: 16px; border-bottom: 1px solid var(--border)">
                <div style="font-weight: 600; color: var(--text-main)">Candidate User</div>
                <div style="font-size: 13px; color: var(--text-muted); margin-top: 4px">candidate@example.com</div>
              </div>
              <div style="padding: 8px">
                <router-link 
                  v-for="item in profileMenu" 
                  :key="item.path"
                  :to="item.path"
                  class="dropdown-item"
                >
                  <component :is="item.icon" size="16" />
                  {{ item.name }}
                </router-link>
              </div>
              <div style="padding: 8px; border-top: 1px solid var(--border)">
                <button class="dropdown-item text-danger" @click="handleLogout">
                  <LogOut size="16" />
                  Đăng xuất
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <main class="page-content" style="max-width: 1200px; margin: 0 auto; padding: 32px">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.candidate-layout {
  min-height: 100vh;
  background-color: var(--background);
}

.top-navbar {
  background-color: var(--surface);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 40;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 32px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-brand {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.brand-logo {
  font-size: 20px;
  font-weight: 700;
  color: var(--primary);
  letter-spacing: -0.5px;
}

.brand-subtitle {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 600;
  margin-top: 2px;
}

.nav-menu {
  display: flex;
  gap: 8px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-weight: 500;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s;
}

.nav-link:hover {
  background-color: var(--surface-soft);
  color: var(--text-main);
}

.nav-link.active {
  background-color: rgba(37, 99, 235, 0.1);
  color: var(--primary);
  font-weight: 600;
}

.nav-actions {
  display: flex;
  align-items: center;
}

.hover-circle:hover {
  background-color: var(--border) !important;
}

.hover-border:hover {
  border-color: var(--primary) !important;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 12px;
  width: 340px;
  background-color: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
  z-index: 50;
  overflow: hidden;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  text-decoration: none;
  border: none;
  background: none;
  width: 100%;
  cursor: pointer;
  transition: background-color 0.2s;
}

.dropdown-item:hover {
  background-color: var(--surface-soft);
  color: var(--text-main);
}

.dropdown-item.text-danger {
  color: var(--danger);
}

.dropdown-item.text-danger:hover {
  background-color: rgba(220, 38, 38, 0.05);
}

.page-content {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
