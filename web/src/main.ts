import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import Casdoor from 'casdoor-vue-sdk'
import { useConfigStore } from '@/stores/config'
const app = createApp(App)
const pinia = createPinia()

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(pinia)
app.use(ElementPlus)
app.use(router)

useConfigStore()
  .getConfig()
  .then((c) => {
    const casdoorConfig = {
      serverUrl: c.casdoor.serverUrl || window.location.origin,
      clientId: c.casdoor.clientId,
      organizationName: c.casdoor.organizationName || 'built-in',
      appName: c.casdoor.appName || 'app-built-in',
      redirectPath: c.casdoor.redirectPath || '/callback',
    }
    // casdoor
    app.use(Casdoor, casdoorConfig)

    app.mount('#app')
  })
