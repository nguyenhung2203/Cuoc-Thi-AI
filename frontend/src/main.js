import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import './styles/global.css'
import App from './App.vue'
import router from './router'

import { authStore } from './stores/auth.store'
import { langStore } from './stores/lang.store'
import { usePlatformStore } from './stores/platform.store'

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
app.use(pinia)

// Khởi tạo Auth Store (Lấy thông tin User nếu có token) trước khi load app
authStore.init().finally(() => {
  app.use(router)
  app.mount('#app')
  
  // Khởi tạo MutationObserver dịch tự động toàn cục
  langStore.initAutoTranslator()

  // Khởi tạo cấu hình thương hiệu từ Platform Store
  const platformStore = usePlatformStore()
  platformStore.init()
})
