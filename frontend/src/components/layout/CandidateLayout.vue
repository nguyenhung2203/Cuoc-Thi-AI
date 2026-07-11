<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { useNotificationStore } from '../../stores/notification.store'
import { langStore } from '../../stores/lang.store'
import AppFooter from './AppFooter.vue'
import AppLogo from '../common/AppLogo.vue'
import { Home, Calendar, Bot, Award, User as UserIcon, FileText, Settings, Bell, LogOut, ChevronDown, Briefcase, Globe, Sun, Moon, Check, Bookmark } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const showNotifications = ref(false)
const showProfileMenu = ref(false)
const showLangMenu = ref(false)
const notificationStore = useNotificationStore()

const currentTheme = ref(localStorage.getItem('app_theme') || 'light')
const currentLang = computed(() => langStore.lang)

const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)

const applyTheme = () => {
  if (currentTheme.value === 'dark') {
    document.documentElement.setAttribute('data-theme', 'dark')
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.setAttribute('data-theme', 'light')
    document.documentElement.classList.remove('dark')
  }
}

const toggleTheme = () => {
  currentTheme.value = currentTheme.value === 'light' ? 'dark' : 'light'
  localStorage.setItem('app_theme', currentTheme.value)
  applyTheme()
}

const setLanguage = (code) => {
  langStore.setLang(code)
  showLangMenu.value = false
}

const closeAllMenus = () => {
  showNotifications.value = false
  showProfileMenu.value = false
  showLangMenu.value = false
}

onMounted(async () => {
  applyTheme()
  window.addEventListener('click', closeAllMenus)
  if (authStore.user) {
    notificationStore.startPolling()
  }
})

onUnmounted(() => {
  window.removeEventListener('click', closeAllMenus)
  notificationStore.stopPolling()
})

const handleMarkAllAsRead = async (e) => {
  e.stopPropagation()
  await notificationStore.markAllRead()
}

const handleMarkAsRead = async (e, id) => {
  e.stopPropagation()
  await notificationStore.markRead(id)
}

const formatTimeAgo = (isoStr) => {
  if (!isoStr) return ''
  const diff = new Date() - new Date(isoStr)
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return langStore.t('nav', 'justNow')
  if (minutes < 60) return `${minutes} ${langStore.t('nav', 'minutesAgo')}`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} ${langStore.t('nav', 'hoursAgo')}`
  return `${Math.floor(hours / 24)} ${langStore.t('nav', 'daysAgo')}`
}

const candidateMenu = computed(() => [
  { path: '/', name: langStore.t('nav', 'about'), icon: Globe },
  { path: '/home', name: langStore.t('nav', 'overview'), icon: Home },
  { path: '/job-board', name: langStore.t('nav', 'jobs'), icon: Briefcase },
  { path: '/my-interviews', name: langStore.t('nav', 'interviews'), icon: Calendar },
  { path: '/mock-setup', name: langStore.t('nav', 'aiPractice'), icon: Bot },
  { path: '/mock-results', name: langStore.t('nav', 'results'), icon: Award }
])

const profileMenu = computed(() => [
  { path: '/profile?tab=profile', name: langStore.t('nav', 'myProfile'), icon: UserIcon },
  { path: '/my-applications', name: langStore.t('nav', 'appliedJobs'), icon: FileText },
  { path: '/saved-jobs', name: langStore.t('nav', 'savedJobs'), icon: Bookmark }
])

const handleLogout = async () => {
  await authStore.logout()
}
</script>

<template>
  <div class="candidate-layout">
    <!-- Top Navbar -->
    <header v-if="!route.path.includes('/candidate-room')" class="top-navbar">
      <div class="nav-container">
        <!-- Left Side: Logo & Menu -->
        <div class="nav-left">
          <!-- Logo -->
          <div class="nav-brand cursor-pointer" @click="router.push('/home')">
            <AppLogo size="md" />
          </div>

          <!-- Menu -->
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
        </div>

        <!-- Right Side: Actions -->
        <div class="nav-actions">
          <!-- Language Selector Dropdown -->
          <div style="position: relative; cursor: pointer; margin-right: 12px" @click.stop="showLangMenu = !showLangMenu; showNotifications = false; showProfileMenu = false">
            <div style="display: flex; align-items: center; gap: 6px; padding: 6px 12px; border-radius: 999px; background-color: var(--surface-soft); border: 1px solid var(--border); transition: all 0.2s; font-size: 13px; font-weight: 600; color: var(--text-main)" class="hover-border">
              <Globe size="15" color="var(--primary)" />
              <span>{{ currentLang === 'en' ? 'EN' : 'VN' }}</span>
              <ChevronDown size="14" color="var(--text-muted)" />
            </div>

            <div v-if="showLangMenu" class="dropdown-menu" style="width: 210px; right: 0;">
              <div style="padding: 10px 14px; border-bottom: 1px solid var(--border); font-size: 11px; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.5px;">
                Ngôn ngữ Giao diện
              </div>
              <div 
                @click.stop="setLanguage('vi')"
                style="padding: 10px 14px; display: flex; align-items: center; justify-content: space-between; font-size: 13px; font-weight: 500; color: var(--text-main); cursor: pointer; transition: background 0.15s;"
                class="hover-bg"
                :style="{ background: currentLang === 'vi' ? 'var(--primary-light)' : 'transparent', color: currentLang === 'vi' ? 'var(--primary)' : 'var(--text-main)' }"
              >
                <span>Tiếng Việt</span>
                <Check v-if="currentLang === 'vi'" size="16" color="var(--primary)" />
              </div>
              <div 
                @click.stop="setLanguage('en')"
                style="padding: 10px 14px; display: flex; align-items: center; justify-content: space-between; font-size: 13px; font-weight: 500; color: var(--text-main); cursor: pointer; transition: background 0.15s;"
                class="hover-bg"
                :style="{ background: currentLang === 'en' ? 'var(--primary-light)' : 'transparent', color: currentLang === 'en' ? 'var(--primary)' : 'var(--text-main)' }"
              >
                <span>English</span>
                <Check v-if="currentLang === 'en'" size="16" color="var(--primary)" />
              </div>
            </div>
          </div>

          <!-- Dark / Light Mode Toggle Button -->
          <div style="cursor: pointer; margin-right: 14px" @click.stop="toggleTheme" :title="currentTheme === 'dark' ? 'Chuyển sang chế độ Sáng (Light Mode)' : 'Chuyển sang chế độ Tối (Dark Mode)'">
            <div style="width: 36px; height: 36px; border-radius: 50%; background-color: var(--surface-soft); border: 1px solid var(--border); display: flex; align-items: center; justify-content: center; transition: all 0.2s" class="hover-circle">
              <Sun v-if="currentTheme === 'dark'" size="18" style="color: #FACC15" />
              <Moon v-else size="18" style="color: var(--text-secondary)" />
            </div>
          </div>

          <template v-if="authStore.isAuthenticated">
            <!-- Notification Bell -->
            <div style="position: relative; cursor: pointer; margin-right: 16px" @click.stop="showNotifications = !showNotifications; showProfileMenu = false; showLangMenu = false">
              <div style="padding: 8px; border-radius: 50%; background-color: var(--surface-soft); transition: background-color 0.2s" class="hover-circle">
                <Bell size="20" color="var(--text-secondary)" />
                <div v-if="unreadCount > 0" style="position: absolute; top: 0px; right: 0px; background-color: var(--danger); color: white; border-radius: 50%; border: 2px solid var(--surface); font-size: 10px; font-weight: bold; width: 16px; height: 16px; display: flex; align-items: center; justify-content: center;">
                  {{ unreadCount }}
                </div>
              </div>
              
              <!-- Notifications Dropdown -->
              <div v-if="showNotifications" class="dropdown-menu">
                <div style="padding: 16px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center">
                  <span style="font-weight: 600; color: var(--text-main)">{{ langStore.t('nav', 'notifications') }}</span>
                  <span style="font-size: 12px; color: var(--primary); font-weight: 500; cursor: pointer;" @click="handleMarkAllAsRead">{{ langStore.t('nav', 'markAllRead') }}</span>
                </div>
                
                <div style="max-height: 400px; overflow-y: auto;">
                  <div v-if="notifications.length === 0" style="padding: 24px; text-align: center; color: var(--text-muted); font-size: 14px;">
                    {{ langStore.t('nav', 'noNotifications') }}
                  </div>
                  <div v-for="n in notifications" :key="n.id" 
                       style="padding: 16px; border-bottom: 1px solid var(--border); transition: background-color 0.2s" 
                       :style="{ backgroundColor: n.is_read ? 'transparent' : 'rgba(37, 99, 235, 0.05)' }"
                       class="hover-bg">
                    <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 6px;">
                      <div class="text-body" style="font-weight: 600; color: var(--text-main)">{{ n.title }}</div>
                      <div v-if="!n.is_read" @click="handleMarkAsRead($event, n.id)" style="width: 8px; height: 8px; background-color: var(--primary); border-radius: 50%; cursor: pointer; flex-shrink: 0;" title="Đánh dấu đã đọc"></div>
                    </div>
                    <div class="text-helper" style="color: var(--text-secondary); line-height: 1.5">
                      {{ n.content || n.message }}
                    </div>
                    <div style="font-size: 11px; color: var(--text-muted); margin-top: 8px">{{ formatTimeAgo(n.created_at) }}</div>
                  </div>
                </div>

                <div style="padding: 12px; text-align: center; color: var(--primary); font-size: 13px; font-weight: 500; cursor: pointer; background-color: var(--surface-soft)">
                  {{ langStore.t('nav', 'viewAllNotifications') }}
                </div>
              </div>
            </div>

            <!-- User Profile Dropdown -->
            <div style="position: relative; cursor: pointer" @click.stop="showProfileMenu = !showProfileMenu; showNotifications = false; showLangMenu = false">
              <div style="display: flex; align-items: center; gap: 8px; padding: 4px 10px 4px 4px; border-radius: 24px; border: 1px solid var(--border); transition: border-color 0.2s" class="hover-border">
                <div style="width: 32px; height: 32px; flex-shrink: 0; border-radius: 50%; background-color: var(--primary); color: white; display: flex; align-items: center; justify-content: center; font-weight: bold">
                  {{ authStore.user?.full_name ? authStore.user.full_name[0].toUpperCase() : 'C' }}
                </div>
                <div style="display: flex; flex-direction: column; max-width: 130px;">
                  <span style="font-size: 12px; font-weight: 600; color: var(--text-main); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ authStore.user?.full_name || 'Ứng viên' }}</span>
                </div>
                <ChevronDown size="14" color="var(--text-muted)" />
              </div>

              <div v-if="showProfileMenu" class="dropdown-menu" style="width: 240px">
                <div style="padding: 16px; border-bottom: 1px solid var(--border)">
                  <div style="font-weight: 600; color: var(--text-main)">{{ authStore.user?.full_name || 'Ứng viên' }}</div>
                  <div style="font-size: 13px; color: var(--text-muted); margin-top: 4px">{{ authStore.user?.email || 'Chưa cập nhật' }}</div>
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
                    {{ langStore.t('nav', 'logout') }}
                  </button>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <button class="btn btn-primary" @click="router.push('/login')">Đăng nhập / Đăng ký</button>
          </template>
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <main class="page-content" :class="[(route.path.includes('candidate-room') || route.path.includes('mock-room')) ? 'room-fullscreen' : (route.path !== '/' ? 'container-bounded' : '')]">
      <router-view />
    </main>
    <AppFooter v-if="!route.path.includes('candidate-room') && !route.path.includes('mock-room')" />
  </div>
</template>

<style scoped>
.candidate-layout {
  min-height: 100vh;
  background-color: var(--background);
  overflow-x: clip;
}

.top-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 40;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  height: 72px;
  display: flex;
  align-items: center;
}

.page-content {
  padding-top: 72px;
}

.page-content.room-fullscreen {
  padding-top: 0 !important;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.container-bounded {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px;
}

.page-content.container-bounded {
  padding-top: calc(72px + 32px);
}

.nav-container {
  max-width: 1440px;
  margin: 0 auto;
  padding: 0 32px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 40px;
}

.nav-brand {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.brand-img {
  height: 48px;
  object-fit: contain;
}

.nav-menu {
  display: flex;
  align-items: center;
  gap: 16px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-weight: 500;
  font-size: 14px;
  text-decoration: none;
  white-space: nowrap;
  transition: all 0.2s;
}

.nav-link:hover {
  background-color: var(--surface-soft);
  color: var(--text-main);
}

.nav-link.active {
  background: var(--primary-light);
  color: var(--primary);
  font-weight: 600;
  box-shadow: inset 0 -2px 0 var(--primary);
}

.nav-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
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
