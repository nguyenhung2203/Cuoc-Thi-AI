<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { useNotificationStore } from '../../stores/notification.store'
import { fileService } from '../../services/file.service'
import { authService } from '../../services/auth.service'
import { LogOut, Home, Briefcase, Users, Calendar, BarChart2, BookOpen, Bot, Settings, Bell, Scale, ChevronLeft, ChevronRight } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const notificationStore = useNotificationStore()
const showNotifications = ref(false)
const isCollapsed = ref(false)
const uploading = ref(false)

const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)

const handleUploadDocument = async (event) => {
  const file = event.target.files[0];
  if (!file) return;
  
  try {
    uploading.value = true;
    const response = await fileService.uploadFile(file, 'image');
    
    // Call Verify Document API
    await authService.verifyDocument(response.data.id);
    
    // Update local state
    if (authStore.user) {
      authStore.user.verification_file_id = response.data.id;
    }
    
    alert('Đã tải tài liệu lên thành công. Vui lòng chờ Admin phê duyệt.');
  } catch (error) {
    console.error('Upload failed', error);
    alert('Lỗi tải file: ' + (error.response?.data?.message || error.message));
  } finally {
    uploading.value = false;
  }
};

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
    notificationStore.startPolling()
  }
})

onUnmounted(() => {
  notificationStore.stopPolling()
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
  <div class="rc-shell">
    <!-- Verification Overlay -->
    <div v-if="authStore.user?.role === 'recruiter' && authStore.user?.status === 'pending'" class="rc-verify-overlay">
      <div class="rc-verify-card">
        <div class="rc-verify-icon"><Scale :size="32" /></div>
        <h2 class="rc-verify-title">Xác thực tài khoản</h2>
        <p class="rc-verify-desc">
          Tài khoản Nhà tuyển dụng của bạn đang chờ phê duyệt. Vui lòng tải lên Giấy phép kinh doanh (hoặc tài liệu xác minh doanh nghiệp) để chúng tôi xử lý.
        </p>

        <div v-if="authStore.user?.verification_file_id" class="rc-verify-ok">
          <span>Đã tải lên tài liệu xác minh. Đang chờ Admin duyệt.</span>
        </div>
        <div v-else class="space-y-4">
          <label class="rc-verify-drop">
            <input type="file" class="hidden" @change="handleUploadDocument" accept=".pdf,.png,.jpg,.jpeg" />
            <div class="flex flex-col items-center">
              <span class="rc-verify-drop-title">Nhấn để chọn file tải lên</span>
              <span class="text-helper mt-1">Hỗ trợ PDF, PNG, JPG (Tối đa 5MB)</span>
            </div>
          </label>
          <div v-if="uploading" style="color: var(--primary);">Đang tải lên...</div>
        </div>

        <div class="mt-8 flex justify-center">
          <button @click="authStore.logout()" class="rc-verify-logout">Đăng xuất</button>
        </div>
      </div>
    </div>

    <!-- Sidebar -->
    <aside class="rc-sidebar" :class="{ 'is-collapsed': isCollapsed }">
      <div class="rc-brand">
        <div v-if="!isCollapsed" class="cursor-pointer flex items-center overflow-hidden" @click="router.push('/dashboard')">
          <img src="/images/logo.png" alt="Logo" class="h-10 object-contain" />
        </div>
        <button @click="isCollapsed = !isCollapsed" class="collapse-btn" :class="{ 'ml-auto': !isCollapsed, 'mx-auto': isCollapsed }" :title="isCollapsed ? 'Mở rộng menu' : 'Thu gọn menu'">
          <ChevronLeft v-if="!isCollapsed" :size="18" />
          <ChevronRight v-else :size="18" />
        </button>
      </div>

      <nav class="rc-nav custom-scrollbar">
        <router-link
          v-for="item in currentMenu"
          :key="item.path"
          :to="item.path"
          :title="isCollapsed ? item.name : ''"
          class="rc-nav-item"
          :class="{ 'is-collapsed': isCollapsed }"
          active-class="is-active"
        >
          <component :is="item.icon" :size="20" class="shrink-0" />
          <span v-if="!isCollapsed" class="truncate">{{ item.name }}</span>
        </router-link>
      </nav>

      <div class="rc-foot" :class="{ 'is-collapsed': isCollapsed }">
        <div class="rc-profile" :class="{ 'is-collapsed': isCollapsed }" :title="isCollapsed ? (authStore.user?.full_name || 'Recruiter User') : ''">
          <div class="rc-avatar">
            {{ authStore.user?.full_name ? authStore.user.full_name.charAt(0).toUpperCase() : 'R' }}
          </div>
          <div v-if="!isCollapsed" class="overflow-hidden flex-1">
            <div class="rc-name truncate">{{ authStore.user?.full_name || 'Recruiter User' }}</div>
            <div class="rc-org truncate">{{ authStore.user?.companies?.[0]?.name || 'HR Department' }}</div>
          </div>
        </div>
        <button @click="handleLogout" :title="isCollapsed ? 'Đăng xuất' : ''" class="logout-btn" :class="{ 'is-collapsed': isCollapsed }">
          <LogOut :size="16" class="shrink-0" />
          <span v-if="!isCollapsed">Đăng xuất</span>
        </button>
      </div>
    </aside>
    
    <!-- Main Content -->
    <main class="rc-main custom-scrollbar">
      <header v-if="!route.path.includes('/recruiter-room')" class="rc-header">
        <!-- Notifications -->
        <div class="relative cursor-pointer" @click="showNotifications = !showNotifications">
          <div class="icon-btn"><Bell :size="20" /></div>
          <div v-if="unreadCount > 0" class="notif-badge">{{ unreadCount }}</div>

          <!-- Notifications Dropdown -->
          <div v-if="showNotifications" class="notif-dropdown" @click.stop>
            <div class="notif-head">
              <span class="notif-title">Thông báo</span>
              <span class="notif-mark" @click="handleMarkAllAsRead">Đánh dấu đã đọc</span>
            </div>

            <div class="notif-list">
              <div v-if="notifications.length === 0" class="notif-empty">Chưa có thông báo nào.</div>
              <div v-for="n in notifications" :key="n.id" class="notif-item" :class="{ 'is-unread': !n.is_read }">
                <div class="flex justify-between items-start mb-1">
                  <div class="notif-item-title">{{ n.title }}</div>
                  <div v-if="!n.is_read" class="notif-dot" title="Đánh dấu đã đọc" @click="handleMarkAsRead($event, n.id)"></div>
                </div>
                <div class="notif-item-body">{{ n.content || n.message }}</div>
                <div class="notif-item-time">{{ formatTimeAgo(n.created_at) }}</div>
              </div>
            </div>

            <div class="notif-foot">Xem tất cả thông báo</div>
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

<style scoped>
.rc-shell { display: flex; height: 100vh; overflow: hidden; background: var(--background); color: var(--text-main); font-family: var(--sans); }

/* Sidebar */
.rc-sidebar { width: 232px; background: var(--surface); border-right: 1px solid var(--border); display: flex; flex-direction: column; transition: width 0.25s ease; box-shadow: var(--shadow-sm); z-index: 20; flex-shrink: 0; }
.rc-sidebar.is-collapsed { width: 64px; }
.rc-brand { display: flex; align-items: center; justify-content: space-between; height: 64px; padding: 0 12px; border-bottom: 1px solid var(--border); }
.collapse-btn { display: flex; align-items: center; justify-content: center; padding: 7px; border: none; background: transparent; border-radius: var(--radius); color: var(--text-muted); cursor: pointer; transition: all 0.2s ease; flex-shrink: 0; }
.collapse-btn:hover { background: var(--surface-soft); color: var(--text-main); }

.rc-nav { flex: 1; overflow-y: auto; padding: 12px; display: flex; flex-direction: column; gap: 3px; }
.rc-nav-item { display: flex; align-items: center; gap: 12px; padding: 10px 12px; border-radius: var(--radius); font-size: 14px; font-weight: 500; color: var(--text-secondary); transition: all 0.18s ease; }
.rc-nav-item.is-collapsed { justify-content: center; padding: 10px 0; }
.rc-nav-item :deep(svg) { color: var(--text-muted); transition: color 0.18s ease; flex-shrink: 0; }
.rc-nav-item:hover { background: var(--surface-soft); color: var(--primary); }
.rc-nav-item:hover :deep(svg) { color: var(--primary); }
.rc-nav-item.is-active { background: var(--primary-light); color: var(--primary); font-weight: 600; }
.rc-nav-item.is-active :deep(svg) { color: var(--primary); }

.rc-foot { border-top: 1px solid var(--border); background: var(--surface-soft); padding: 16px; }
.rc-foot.is-collapsed { padding: 8px; }
.rc-profile { display: flex; align-items: center; gap: 12px; padding: 8px; border-radius: var(--radius); margin-bottom: 6px; cursor: pointer; border: 1px solid transparent; transition: all 0.2s ease; }
.rc-profile.is-collapsed { justify-content: center; padding: 4px; }
.rc-profile:hover { background: var(--surface); border-color: var(--border); }
.rc-avatar { width: 40px; height: 40px; border-radius: 50%; background: var(--primary); color: #fff; display: flex; align-items: center; justify-content: center; font-weight: 700; flex-shrink: 0; }
.rc-name { font-size: 14px; font-weight: 600; color: var(--text-main); }
.rc-org { font-size: 12px; color: var(--text-muted); }
.logout-btn { width: 100%; display: flex; align-items: center; gap: 8px; padding: 8px 12px; border: none; background: transparent; border-radius: var(--radius); font-size: 14px; font-weight: 500; color: var(--danger); cursor: pointer; transition: background 0.2s ease; }
.logout-btn.is-collapsed { justify-content: center; padding: 8px 0; }
.logout-btn:hover { background: rgba(220,38,38,0.08); }

/* Main + header */
.rc-main { flex: 1; display: flex; flex-direction: column; height: 100vh; overflow-y: auto; background: var(--background); }
.rc-header { position: sticky; top: 0; height: 68px; display: flex; align-items: center; justify-content: flex-end; gap: 16px; padding: 0 32px; border-bottom: 1px solid var(--border); background: rgba(255,255,255,0.85); backdrop-filter: blur(8px); z-index: 40; flex-shrink: 0; }
.icon-btn { display: flex; align-items: center; justify-content: center; padding: 8px; border-radius: 50%; color: var(--text-secondary); cursor: pointer; transition: background 0.2s ease; }
.icon-btn:hover { background: var(--surface-soft); color: var(--text-main); }
.notif-badge { position: absolute; top: 2px; right: 2px; min-width: 16px; height: 16px; padding: 0 4px; background: var(--danger); color: #fff; border: 2px solid var(--surface); border-radius: var(--radius-full); font-size: 9px; font-weight: 700; display: flex; align-items: center; justify-content: center; }

.rc-content { flex: 1; padding: 24px; }
@media (min-width: 1024px) { .rc-content { padding: 32px; } }

.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: var(--border); border-radius: 4px; }
</style>

<style scoped>
/* Notification dropdown */
.notif-dropdown { position: absolute; top: 100%; right: 0; margin-top: 8px; width: 360px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-lg); z-index: 50; overflow: hidden; }
.notif-head { display: flex; justify-content: space-between; align-items: center; padding: 16px; border-bottom: 1px solid var(--border); background: var(--surface-soft); }
.notif-title { font-weight: 600; color: var(--text-main); }
.notif-mark { font-size: 12px; color: var(--primary); font-weight: 500; cursor: pointer; }
.notif-mark:hover { text-decoration: underline; }
.notif-list { max-height: 400px; overflow-y: auto; }
.notif-empty { padding: 32px; text-align: center; color: var(--text-muted); font-size: 14px; }
.notif-item { padding: 16px; border-bottom: 1px solid var(--border); transition: background 0.2s ease; }
.notif-item:hover { background: var(--surface-soft); }
.notif-item.is-unread { background: rgba(37,99,235,0.05); }
.notif-item-title { font-size: 14px; font-weight: 600; color: var(--text-main); }
.notif-dot { width: 8px; height: 8px; background: var(--primary); border-radius: 50%; flex-shrink: 0; cursor: pointer; }
.notif-item-body { font-size: 13px; color: var(--text-secondary); line-height: 1.5; margin: 4px 0 8px; }
.notif-item-time { font-size: 11px; color: var(--text-muted); }
.notif-foot { padding: 12px; text-align: center; color: var(--primary); font-size: 13px; font-weight: 500; cursor: pointer; background: var(--surface-soft); }
.notif-foot:hover { background: var(--border); }

/* Verification overlay */
.rc-verify-overlay { position: fixed; inset: 0; z-index: 100; display: flex; align-items: center; justify-content: center; background: rgba(15,23,42,0.8); backdrop-filter: blur(4px); padding: 16px; }
.rc-verify-card { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 32px; max-width: 32rem; width: 100%; text-align: center; box-shadow: var(--shadow-lg); }
.rc-verify-icon { width: 64px; height: 64px; margin: 0 auto 22px; border-radius: 50%; background: var(--primary-light); color: var(--primary); display: flex; align-items: center; justify-content: center; }
.rc-verify-title { font-size: 22px; font-weight: 700; color: var(--text-main); margin-bottom: 8px; }
.rc-verify-desc { color: var(--text-secondary); margin-bottom: 22px; line-height: 1.6; }
.rc-verify-ok { padding: 16px; background: rgba(22,163,74,0.08); border: 1px solid rgba(22,163,74,0.2); border-radius: var(--radius); color: var(--success); }
.rc-verify-drop { display: block; width: 100%; border: 2px dashed var(--border); border-radius: var(--radius); padding: 32px; cursor: pointer; transition: border-color 0.2s ease; background: var(--surface-soft); }
.rc-verify-drop:hover { border-color: var(--primary); }
.rc-verify-drop-title { color: var(--text-secondary); font-weight: 500; }
.rc-verify-logout { background: none; border: none; color: var(--text-muted); cursor: pointer; transition: color 0.2s ease; }
.rc-verify-logout:hover { color: var(--text-main); }
</style>
