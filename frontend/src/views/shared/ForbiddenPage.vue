<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ShieldAlert, ArrowLeft, LogIn, Lock, HelpCircle, UserX, AlertTriangle, Home, Sparkles, KeyRound } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()

const currentRole = computed(() => {
  return localStorage.getItem('user_role') || 'Chưa đăng nhập'
})

const userEmail = computed(() => {
  return localStorage.getItem('user_email') || 'N/A'
})

const attemptedPath = computed(() => {
  return route.query.attempted || route.redirectedFrom?.path || '/admin'
})

const reasonCode = computed(() => {
  return route.query.reason || 'unauthorized'
})

const warningDetails = computed(() => {
  switch (reasonCode.value) {
    case 'admin_required':
      return {
        title: 'Khu Vực Giới Hạn Quản Trị Viên (Admin Portal)',
        desc: 'Trang bạn đang cố gắng truy cập chỉ dành riêng cho Quản trị viên hệ thống (Admin). Tài khoản hiện tại không có đặc quyền để xem hoặc thao tác trên phân hệ này.',
        requiredRole: 'Quản trị viên (Admin)',
        actionText: 'Quay về Không gian của bạn',
        actionPath: currentRole.value === 'recruiter' ? '/dashboard' : '/home'
      }
    case 'recruiter_required':
      return {
        title: 'Khu Vực Dành Cho Nhà Tuyển Dụng (Recruiter Portal)',
        desc: 'Bạn đang đăng nhập với tư cách Ứng viên (Candidate). Các chức năng tạo tin tuyển dụng, quản lý ngân hàng câu hỏi, rubric AI và hồ sơ ứng viên chỉ dành cho Nhà tuyển dụng.',
        requiredRole: 'Nhà tuyển dụng (Recruiter)',
        actionText: 'Quay về Bảng điều khiển Ứng viên',
        actionPath: '/home'
      }
    case 'candidate_required':
      return {
        title: 'Tính Năng Dành Riêng Cho Ứng Viên (Candidate Only)',
        desc: 'Các chức năng phỏng vấn thử AI (Mock Interview), luyện tập trả lời câu hỏi và phòng phỏng vấn 1-1 AI được thiết kế tối ưu dành riêng cho tài khoản Ứng viên.',
        requiredRole: 'Ứng viên (Candidate)',
        actionText: 'Quay về Bảng điều khiển Nhà tuyển dụng',
        actionPath: '/dashboard'
      }
    case 'account_blocked':
      return {
        title: 'Tài Khoản Đang Bị Tạm Khóa (Suspended)',
        desc: 'Tài khoản của bạn đã bị Quản trị viên tạm dừng hoạt động do vi phạm quy tắc sử dụng hoặc yêu cầu bảo mật hệ thống. Vui lòng liên hệ bộ phận hỗ trợ.',
        requiredRole: 'Tài khoản hoạt động (Active)',
        actionText: 'Quay lại Trang Đăng nhập',
        actionPath: '/login'
      }
    case 'pending_approval':
      return {
        title: 'Tài Khoản Đang Chờ Xét Duyệt (Pending Approval)',
        desc: 'Giấy phép hoặc thông tin nhà tuyển dụng đang được Quản trị viên thẩm định. Các tính năng chuyên sâu sẽ tự động mở khóa ngay sau khi tài khoản được phê duyệt.',
        requiredRole: 'Tài khoản đã xác minh (Verified)',
        actionText: 'Quay lại Trang Tổng quan',
        actionPath: '/dashboard'
      }
    default:
      return {
        title: 'Quyền Truy Cập Bị Từ Chối (Access Denied / 403)',
        desc: 'Hệ thống đã từ chối yêu cầu truy cập của bạn do thiếu giấy phép xác thực hợp lệ hoặc quyền hạn hiện tại không tương thích với phân hệ này.',
        requiredRole: 'Cấp quyền hợp lệ (Valid Token)',
        actionText: 'Quay lại Trang Chủ',
        actionPath: currentRole.value === 'admin' ? '/admin/dashboard' : (currentRole.value === 'recruiter' ? '/dashboard' : '/home')
      }
  }
})

const goBack = () => {
  const target = warningDetails.value.actionPath || '/'
  router.push(target)
}

const switchAccount = () => {
  localStorage.removeItem('access_token')
  localStorage.removeItem('user_role')
  localStorage.removeItem('user_email')
  localStorage.removeItem('user_id')
  router.push('/login')
}
</script>

<template>
  <div class="h-screen w-screen overflow-hidden bg-slate-950 text-slate-100 flex flex-col items-center justify-center p-4 relative font-sans select-none">
    <!-- Futuristic Deep Space Background & Grid -->
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(244,63,94,0.14),transparent_60%)] pointer-events-none"></div>
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_80%_20%,rgba(168,85,247,0.10),transparent_50%)] pointer-events-none"></div>
    
    <!-- Cyber Grid pattern overlay -->
    <div class="absolute inset-0 bg-[linear-gradient(to_right,rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:3rem_3rem] pointer-events-none"></div>

    <!-- Giant background watermark -->
    <div class="absolute inset-0 flex items-center justify-center pointer-events-none overflow-hidden opacity-[0.025]">
      <span class="text-[16vw] font-black tracking-tighter text-rose-500 whitespace-nowrap">403 FORBIDDEN</span>
    </div>

    <!-- Top Floating Security Badge (Compact) -->
    <div class="absolute top-4 left-6 right-6 flex justify-between items-center z-20 pointer-events-none">
      <div class="flex items-center gap-2 px-3 py-1 bg-rose-500/10 border border-rose-500/30 rounded-full text-rose-400 font-mono text-[11px] font-bold tracking-wider uppercase backdrop-blur-md shadow-md">
        <span class="w-1.5 h-1.5 rounded-full bg-rose-500 animate-ping"></span>
        <span>Security Portal // Error 403</span>
      </div>
      <div class="hidden sm:flex items-center gap-2 text-[11px] font-mono text-slate-400 px-3 py-1 bg-slate-900/80 rounded-full border border-slate-800">
        <span>Session ID: {{ Math.random().toString(36).substring(2, 10).toUpperCase() }}</span>
      </div>
    </div>

    <!-- Main Holographic Card (Compact 1-Frame Layout) -->
    <div class="max-w-xl w-full bg-slate-900/95 backdrop-blur-2xl border border-rose-500/30 rounded-3xl p-6 sm:p-8 shadow-[0_0_80px_rgba(244,63,94,0.18)] relative z-10 space-y-5 transition-all max-h-[90vh] flex flex-col justify-between">
      
      <!-- Top Decorative Glow bar -->
      <div class="absolute -top-px left-12 right-12 h-1 bg-gradient-to-r from-transparent via-rose-500 to-transparent opacity-80"></div>

      <!-- Icon & Status Header (Compact) -->
      <div class="flex flex-col items-center text-center space-y-3">
        <div class="relative w-16 h-16 flex items-center justify-center">
          <!-- Pulsing Rings -->
          <div class="absolute inset-0 bg-rose-500/20 rounded-2xl transform rotate-6 animate-pulse"></div>
          <div class="absolute inset-0 bg-rose-500/10 rounded-2xl transform -rotate-6 animate-ping opacity-30"></div>
          
          <!-- Shield Box -->
          <div class="w-16 h-16 bg-gradient-to-tr from-rose-600 via-red-600 to-amber-600 rounded-2xl flex items-center justify-center shadow-lg shadow-rose-600/40 border-2 border-rose-400/50 relative z-10 transform transition-transform hover:scale-105">
            <ShieldAlert size="32" class="text-white animate-bounce" />
          </div>
        </div>

        <div class="space-y-1">
          <h1 class="text-2xl sm:text-3xl font-black tracking-tight bg-gradient-to-r from-rose-400 via-red-500 to-amber-400 bg-clip-text text-transparent uppercase">
            CẢNH BÁO: TRUY CẬP BỊ TỪ CHỐI
          </h1>
          <p class="text-rose-300 font-extrabold text-sm sm:text-base flex items-center justify-center gap-1.5">
            <Lock size="15" class="text-rose-400 shrink-0" /> {{ warningDetails.title }}
          </p>
        </div>
      </div>

      <!-- Security Details Box (Tight & Balanced) -->
      <div class="bg-slate-950/80 border border-slate-800/90 rounded-2xl p-4 sm:p-5 space-y-3.5 shadow-inner">
        <!-- Explanatory message -->
        <div class="border-l-3 border-rose-500 bg-rose-950/25 px-3.5 py-2.5 rounded-r-xl">
          <p class="text-slate-200 text-xs sm:text-sm leading-relaxed font-medium">
            {{ warningDetails.desc }}
          </p>
        </div>

        <!-- Role Comparison Pills (Compact Grid) -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
          <div class="bg-slate-900/90 border border-slate-800 p-3 rounded-xl space-y-1">
            <span class="text-[10px] font-mono uppercase text-slate-500 font-semibold block">Vai trò hiện tại (Current)</span>
            <div class="flex items-center justify-between gap-2">
              <span class="text-indigo-400 font-bold capitalize text-sm truncate">{{ currentRole }}</span>
              <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-indigo-500/20 text-indigo-300 border border-indigo-500/30 truncate max-w-[120px]" :title="userEmail">
                {{ userEmail }}
              </span>
            </div>
          </div>

          <div class="bg-slate-900/90 border border-slate-800 p-3 rounded-xl space-y-1">
            <span class="text-[10px] font-mono uppercase text-slate-500 font-semibold block">Yêu cầu quyền (Required)</span>
            <div class="flex items-center justify-between gap-2">
              <span class="text-amber-400 font-bold text-sm truncate">{{ warningDetails.requiredRole }}</span>
              <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-rose-500/20 text-rose-300 border border-rose-500/30 flex items-center gap-1 shrink-0">
                <KeyRound size="10" /> Restricted
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Help note (Compact) -->
      <div class="flex items-center justify-center gap-1.5 text-[11px] text-slate-400 text-center font-medium">
        <HelpCircle size="13" class="text-amber-400 shrink-0" />
        <span>Cần hỗ trợ? Liên hệ Quản trị viên qua <a href="mailto:admin@wemake.vn" class="text-amber-400 hover:underline font-bold">admin@wemake.vn</a>.</span>
      </div>

      <!-- Action Buttons (Tight row) -->
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3 pt-1">
        <button 
          @click="goBack"
          class="w-full sm:w-auto flex-1 py-3.5 px-6 bg-gradient-to-r from-rose-600 via-red-600 to-amber-600 hover:from-rose-500 hover:to-amber-500 text-white font-bold text-sm rounded-xl shadow-lg shadow-rose-600/25 hover:shadow-rose-600/40 transition-all duration-300 transform hover:-translate-y-0.5 active:translate-y-0 flex items-center justify-center gap-2 group"
        >
          <Home size="16" class="transition-transform group-hover:scale-110 shrink-0" />
          <span class="truncate">{{ warningDetails.actionText }}</span>
        </button>

        <button 
          @click="switchAccount"
          class="w-full sm:w-auto px-5 py-3.5 bg-slate-800/90 hover:bg-slate-700/90 text-slate-200 font-bold text-sm rounded-xl border border-slate-700 hover:border-slate-600 transition-all duration-300 flex items-center justify-center gap-2 shadow-md shrink-0"
        >
          <LogIn size="16" class="text-indigo-400 shrink-0" />
          <span>Đổi tài khoản khác</span>
        </button>
      </div>
    </div>
  </div>
</template>
