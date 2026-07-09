import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import './styles/global.css'
import App from './App.vue'
import router from './router'

import { authStore } from './stores/auth.store'

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
app.use(pinia)

// Khởi tạo Auth Store (Lấy thông tin User nếu có token) trước khi load app
authStore.init().finally(() => {
  app.use(router)
  app.mount('#app')
})
