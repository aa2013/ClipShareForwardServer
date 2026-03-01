/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import {createRouter, createWebHashHistory} from 'vue-router/auto'
import {routes} from 'vue-router/auto-routes'
import {local} from "@/utils/user";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    ...routes,
    {
      path: '/',
      redirect: '/admin/connection',
    },
  ],
})
const loginPath = '/login'
router.beforeEach((to, from, next) => {
  if (local.token) {
    next()
    return;
  } else {
    if (to.path !== loginPath) {
      next(loginPath)
      return
    } else {
      next()
    }
  }
})
// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.('Failed to fetch dynamically imported module')) {
    if (!localStorage.getItem('vuetify:dynamic-reload')) {
      console.log('Reloading page to fix dynamic import error')
      localStorage.setItem('vuetify:dynamic-reload', 'true')
      location.assign(to.fullPath)
    } else {
      console.error('Dynamic import error, reloading page did not fix it', err)
    }
  } else {
    console.error(err)
  }
})

router.isReady().then(() => {
  localStorage.removeItem('vuetify:dynamic-reload')
})

export default router
