<script setup>
import { Bot, Sparkles } from 'lucide-vue-next'
import { usePlatformStore } from '../../stores/platform.store'

const platformStore = usePlatformStore()

defineProps({
  size: {
    type: String,
    default: 'md' // sm, md, lg
  },
  showSubtitle: {
    type: Boolean,
    default: false
  },
  subtitle: {
    type: String,
    default: ''
  },
  variant: {
    type: String,
    default: 'default' // default, light
  }
})
</script>

<template>
  <div class="app-logo flex items-center select-none cursor-pointer group" :class="[size === 'sm' ? 'gap-2.5' : size === 'lg' ? 'gap-4' : 'gap-3']">
    <!-- Icon Emblem / Custom Image Logo Badge -->
    <div 
      class="logo-badge relative rounded-2xl bg-gradient-to-br from-[var(--primary)] via-indigo-600 to-[var(--accent)] text-white flex items-center justify-center shadow-md shadow-blue-500/20 border border-white/20 transition-all duration-300 group-hover:scale-105 group-hover:shadow-lg group-hover:shadow-blue-500/30 shrink-0 overflow-hidden"
      :class="[
        size === 'sm' ? 'w-8 h-8 rounded-xl' : size === 'lg' ? 'w-14 h-14 rounded-2xl' : 'w-11 h-11 rounded-2xl',
        platformStore.brandLogoUrl ? 'bg-transparent border-none shadow-none' : ''
      ]"
    >
      <!-- TRƯỜNG HỢP 1: NẾU ADMIN CẤU HÌNH ẢNH LOGO RIÊNG (brandLogoUrl) -->
      <img 
        v-if="platformStore.brandLogoUrl" 
        :src="platformStore.brandLogoUrl" 
        alt="Brand Logo" 
        class="w-full h-full object-contain p-0.5 transition-transform duration-300 group-hover:scale-110" 
      />

      <!-- TRƯỜNG HỢP 2: NẾU DÙNG ICON VECTOR AI MẶC ĐỊNH -->
      <template v-else>
        <!-- Subtle internal glow and sparkles -->
        <Sparkles 
          class="absolute top-1 right-1 text-cyan-200 opacity-80 animate-pulse" 
          :size="size === 'sm' ? 10 : size === 'lg' ? 14 : 12" 
        />
        
        <Bot 
          class="relative z-10 transition-transform duration-300 group-hover:rotate-6" 
          :size="size === 'sm' ? 18 : size === 'lg' ? 30 : 24" 
        />

        <!-- Glossy sheen effect overlay -->
        <div class="absolute inset-0 bg-gradient-to-tr from-white/0 via-white/10 to-white/20 pointer-events-none"></div>
      </template>
    </div>

    <!-- Wordmark & Tagline -->
    <div class="logo-text flex flex-col justify-center text-left">
      <div 
        class="font-black tracking-tight leading-none flex items-center"
        :class="[
          size === 'sm' ? 'text-base gap-1' : size === 'lg' ? 'text-2xl gap-2' : 'text-xl gap-1.5',
          variant === 'light' ? 'text-white' : 'text-[var(--text-main)]'
        ]"
      >
        <span class="group-hover:text-[var(--primary)] transition-colors duration-200">{{ platformStore.brandName }}</span>
        <span 
          v-if="platformStore.brandBadge"
          class="rounded-lg bg-gradient-to-r from-[var(--primary)] to-[var(--accent)] text-white font-black uppercase tracking-widest shadow-sm inline-flex items-center justify-center"
          :class="[
            size === 'sm' ? 'px-1.5 py-0.5 text-[9px]' : size === 'lg' ? 'px-2.5 py-1 text-xs' : 'px-2 py-0.5 text-[10px]'
          ]"
        >
          {{ platformStore.brandBadge }}
        </span>
      </div>

      <span 
        v-if="showSubtitle || subtitle" 
        class="font-bold uppercase tracking-wider block transition-colors"
        :class="[
          size === 'sm' ? 'text-[9px] mt-0.5' : size === 'lg' ? 'text-xs mt-1.5' : 'text-[11px] mt-1',
          variant === 'light' ? 'text-blue-100' : 'text-[var(--text-muted)] group-hover:text-[var(--text-secondary)]'
        ]"
      >
        {{ subtitle || platformStore.brandSlogan }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.app-logo {
  font-family: var(--sans, 'Inter', -apple-system, sans-serif);
}
</style>
