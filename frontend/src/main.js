import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import './style.css'

// 全局标题和favicon管理
const updatePageTitle = () => {
  const settingsData = localStorage.getItem('site_settings')
  if (settingsData) {
    try {
      const settings = JSON.parse(settingsData)
      if (settings.siteTitle) {
        document.title = settings.siteTitle
      }
      if (settings.siteLogo) {
        const favicon = document.getElementById('favicon')
        if (favicon) {
          favicon.href = settings.siteLogo
        }
      }
    } catch (error) {
      console.error('Failed to parse site settings:', error)
    }
  }
}

const app = createApp(App)
const pinia = createPinia()

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 应用启动时更新标题
updatePageTitle()

// 路由变化时更新标题
router.beforeEach((to, from, next) => {
  updatePageTitle()
  next()
})

app.use(pinia)
app.use(router)
app.use(ElementPlus)

app.mount('#app')
