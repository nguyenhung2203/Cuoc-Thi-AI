import { createApp } from 'vue'
import './styles/global.css'
import App from './App.vue'
import router from './router'

import { authStore } from './stores/auth.store'

const app = createApp(App)

// Khởi tạo Auth Store (Lấy thông tin User nếu có token) trước khi load app
authStore.init().finally(() => {
  app.use(router)
  app.mount('#app')
})
