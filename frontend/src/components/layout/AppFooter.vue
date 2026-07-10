<script setup>
import { useRouter } from 'vue-router'
import { Sparkles, Mail, Phone, MapPin, Github, Linkedin, Facebook, Globe, ShieldCheck, Heart } from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <footer class="bg-white dark:bg-slate-900 border-t border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 text-sm transition-colors duration-300">
    <div class="max-w-7xl mx-auto px-6 lg:px-8 py-12 lg:py-16">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-10 lg:gap-8">
        
        <!-- Brand Column (2 cols) -->
        <div class="lg:col-span-2 space-y-4">
          <div @click="router.push('/')" class="flex items-center gap-3 cursor-pointer group w-fit">
            <img src="/images/logo.png" alt="WeMake AI Logo" class="h-12 object-contain transform transition-transform group-hover:scale-105" />
            <div>
              <span class="text-lg font-extrabold bg-gradient-to-r from-indigo-600 via-blue-600 to-purple-600 bg-clip-text text-transparent block tracking-tight">WeMake AI</span>
              <span class="text-[11px] font-semibold text-slate-400 uppercase tracking-widest block -mt-1">Recruitment Platform</span>
            </div>
          </div>
          
          <p class="text-xs sm:text-sm text-slate-500 dark:text-slate-400 leading-relaxed pr-6">
            Nền tảng phỏng vấn & tuyển dụng thông minh tích hợp Trí tuệ Nhân tạo thế hệ mới. Tự động bóc tách CV, phỏng vấn trực tuyến và đánh giá ứng viên chính xác tuyệt đối.
          </p>

          <div class="flex items-center gap-3 pt-2">
            <a href="https://github.com" target="_blank" class="w-9 h-9 rounded-xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-600 dark:text-slate-300 hover:bg-indigo-600 hover:text-white dark:hover:bg-indigo-600 transition-all shadow-sm">
              <Github size="18" />
            </a>
            <a href="https://linkedin.com" target="_blank" class="w-9 h-9 rounded-xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-600 dark:text-slate-300 hover:bg-blue-600 hover:text-white dark:hover:bg-blue-600 transition-all shadow-sm">
              <Linkedin size="18" />
            </a>
            <a href="https://facebook.com" target="_blank" class="w-9 h-9 rounded-xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-600 dark:text-slate-300 hover:bg-blue-500 hover:text-white dark:hover:bg-blue-500 transition-all shadow-sm">
              <Facebook size="18" />
            </a>
            <a href="https://wemake.vn" target="_blank" class="w-9 h-9 rounded-xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-600 dark:text-slate-300 hover:bg-purple-600 hover:text-white dark:hover:bg-purple-600 transition-all shadow-sm">
              <Globe size="18" />
            </a>
          </div>
        </div>

        <!-- Solutions Column -->
        <div class="space-y-4">
          <h3 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-1.5">
            <Sparkles size="14" class="text-indigo-600 dark:text-indigo-400" />
            Giải pháp AI
          </h3>
          <ul class="space-y-2.5 text-xs font-medium">
            <li><a @click.prevent="router.push('/login')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Phỏng vấn Trực tuyến AI</a></li>
            <li><a @click.prevent="router.push('/login')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Bóc tách & Chấm điểm CV</a></li>
            <li><a @click.prevent="router.push('/login')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Ngân hàng Câu hỏi Chuẩn</a></li>
            <li><a @click.prevent="router.push('/login')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Thang điểm Rubric Tự động</a></li>
          </ul>
        </div>

        <!-- Candidates Column -->
        <div class="space-y-4">
          <h3 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider">
            Dành cho Ứng viên
          </h3>
          <ul class="space-y-2.5 text-xs font-medium">
            <li><a @click.prevent="router.push('/job-board')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Cơ hội Việc làm AI</a></li>
            <li><a @click.prevent="router.push('/login')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Luyện tập Phỏng vấn 1-1</a></li>
            <li><a @click.prevent="router.push('/login')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Lịch sử Chấm điểm AI</a></li>
            <li><a @click.prevent="router.push('/login')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Hồ sơ Năng lực Trực tuyến</a></li>
          </ul>
        </div>

        <!-- Contact & Security Column -->
        <div class="space-y-4">
          <h3 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-1.5">
            <ShieldCheck size="14" class="text-emerald-600 dark:text-emerald-400" />
            Bảo mật & Hỗ trợ
          </h3>
          <ul class="space-y-2.5 text-xs font-medium">
            <li class="flex items-center gap-2 text-slate-500 dark:text-slate-400">
              <Mail size="14" class="text-indigo-500 shrink-0" />
              <span>support@wemake.vn</span>
            </li>
            <li class="flex items-center gap-2 text-slate-500 dark:text-slate-400">
              <Phone size="14" class="text-indigo-500 shrink-0" />
              <span>+84 (0) 28 8888 2026</span>
            </li>
            <li class="flex items-start gap-2 text-slate-500 dark:text-slate-400">
              <MapPin size="14" class="text-indigo-500 shrink-0 mt-0.5" />
              <span>Khu Công nghệ Cao, TP. Thủ Đức, TP. Hồ Chí Minh</span>
            </li>
          </ul>
        </div>

      </div>

      <!-- Bottom Bar -->
      <div class="mt-12 pt-8 border-t border-slate-100 dark:border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs">
        <p class="text-slate-400 dark:text-slate-500 flex items-center gap-1">
          <span>© 2026 WeMake AI Platform. Made with</span>
          <Heart size="13" class="text-rose-500 fill-rose-500 inline" />
          <span>& AI Technology.</span>
        </p>
        <div class="flex items-center gap-6 font-semibold text-slate-500 dark:text-slate-400">
          <a @click.prevent="router.push('/')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Điều khoản dịch vụ</a>
          <a @click.prevent="router.push('/')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Chính sách bảo mật 2FA</a>
          <a @click.prevent="router.push('/')" href="#" class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">Trợ giúp</a>
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
</style>
