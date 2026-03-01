/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'

// Components
import App from './App.vue'
import '@/styles/global.less'
const pinia = createPinia()

// Composables
import { createApp } from 'vue'
import {createPinia} from "pinia";

const app = createApp(App)

registerPlugins(app)
app.use(pinia)
app.mount('#app')
