import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import './style.css'
import { router } from './router'
import { i18n, applyDocumentDirection, getInitialLocale } from './i18n'
import { useBrandStore } from './stores/brand'
import { useAuthStore } from './stores/auth'

applyDocumentDirection(getInitialLocale())

const app = createApp(App)
app.use(createPinia())
app.use(i18n)
app.use(router)

async function bootstrap() {
  const brandStore = useBrandStore()
  const authStore = useAuthStore()
  await Promise.all([brandStore.load(), authStore.restoreSession()])
  await router.isReady()
  app.mount('#app')
}

bootstrap()
