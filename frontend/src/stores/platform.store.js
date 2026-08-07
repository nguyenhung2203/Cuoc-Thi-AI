import { defineStore } from 'pinia'
import { ref } from 'vue'

export const usePlatformStore = defineStore('platform', () => {
  const brandName = ref(localStorage.getItem('platform_brand_name') || 'ViệcLàm')
  const brandBadge = ref(localStorage.getItem('platform_brand_badge') || 'AI')
  const brandSlogan = ref(localStorage.getItem('platform_brand_slogan') || 'Nền tảng Phỏng vấn & Tuyển dụng Thông minh')
  const brandLogoUrl = ref(localStorage.getItem('platform_brand_logo_url') || '/images/logo.png')
  const supportEmail = ref(localStorage.getItem('platform_support_email') || 'support@vieclam.ai')
  const systemName = ref(localStorage.getItem('platform_system_name') || 'ViệcLàm AI Recruitment')

  // Cập nhật cấu hình ngay lập tức (Trigger khi Admin ấn Lưu thiết lập)
  const updateConfig = (config) => {
    if (config.brand_name !== undefined) {
      brandName.value = config.brand_name
      localStorage.setItem('platform_brand_name', config.brand_name)
    }
    if (config.brand_badge !== undefined) {
      brandBadge.value = config.brand_badge
      localStorage.setItem('platform_brand_badge', config.brand_badge)
    }
    if (config.brand_slogan !== undefined) {
      brandSlogan.value = config.brand_slogan
      localStorage.setItem('platform_brand_slogan', config.brand_slogan)
    }
    if (config.brand_logo_url !== undefined) {
      brandLogoUrl.value = config.brand_logo_url
      localStorage.setItem('platform_brand_logo_url', config.brand_logo_url)
    }
    if (config.support_email !== undefined) {
      supportEmail.value = config.support_email
      localStorage.setItem('platform_support_email', config.support_email)
    }
    if (config.system_name !== undefined) {
      systemName.value = config.system_name
      localStorage.setItem('platform_system_name', config.system_name)
    }

    // Cập nhật tiêu đề trình duyệt <title>
    const fullName = brandBadge.value ? `${brandName.value} ${brandBadge.value}` : brandName.value
    document.title = `${fullName} | ${brandSlogan.value}`
  }

  // Tự động set title khi init store
  const init = () => {
    const fullName = brandBadge.value ? `${brandName.value} ${brandBadge.value}` : brandName.value
    if (typeof document !== 'undefined') {
      document.title = `${fullName} | ${brandSlogan.value}`
    }
  }

  return {
    brandName,
    brandBadge,
    brandSlogan,
    brandLogoUrl,
    supportEmail,
    systemName,
    updateConfig,
    init
  }
})
