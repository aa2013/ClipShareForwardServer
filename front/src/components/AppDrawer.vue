<template>
  <v-navigation-drawer v-model="showDrawer">
    <div class="flex flex-1 h-[100%] flex-col pb-[4px] px-[8px]">
      <AppLogoTitle v-if="showAppLogoTitle" class="border-b-[1px]" />
      <v-list density="compact" class="flex-1" :opened="openedGroup">
        <template v-for="(item, i) in items" :key="item.value || i">
          <v-list-item
            nav
            :active="item.route === baseRoutePath"
            v-if="!item.children"
            :value="item.value"
            color="primary"
            @click="() => gotoPage(item)"
          >
            <template #prepend>
              <v-icon :icon="item.icon"></v-icon>
            </template>

            <v-list-item-title v-text="item.text"></v-list-item-title>
          </v-list-item>
          <v-list-group v-else :value="item.value" :key="item.text" nav>
            <template #activator="{ props }">
              <v-list-item
                nav
                v-bind="props"
                :title="item.text"
                color="primary"
                :active="item.children.map(c => c.route).includes(baseRoutePath)"
                :prepend-icon="item.icon"
              />
            </template>
            <v-list-item
              v-for="(child, childIndex) in item.children"
              :key="childIndex"
              :prepend-icon="child.icon"
              :title="child.text"
              :value="child.value"
              nav
              :active="child.route === baseRoutePath"
              color="primary"
              @click="() => gotoPage(child)"
            />
          </v-list-group>
        </template>
      </v-list>
      <v-divider opacity="1" class="my-[4px]" />
      <v-card elevation="0" v-ripple>
        <template #title>
          <div class="flex justify-between ">
            <div class="flex align-center">
              <v-avatar color="primary">{{ local.username.substring(0, 1) }}</v-avatar>
              <div class="ml-[8px]">{{ local.username }}</div>
            </div>

            <v-menu>
              <template #activator="{ props }">
                <v-btn icon="mdi-logout" elevation="0" v-bind="props" />
              </template>
              <div class="px-[16px] py-[8px] rounded-[8px] bg-[white] m-[8px]" :style="{ boxShadow: '0 0px 20px #0000001a' }">
                <div class="mb-[8px] flex items-center">
                  <v-icon icon="mdi-information-outline" size="small" class="mr-1" color="orange" />
                  {{ t('drawer.logoutConfirm') }}
                </div>
                <div class="flex justify-end">
                  <v-btn variant="text" class="mr-[4px]" size="small">
                    {{ t('common.cancel') }}
                  </v-btn>
                  <v-btn variant="text" @click="logout" color="primary" size="small">
                    {{ t('common.confirm') }}
                  </v-btn>
                </div>
              </div>
            </v-menu>
          </div>
        </template>
      </v-card>
    </div>
  </v-navigation-drawer>
</template>
<script setup lang="ts">
import AppLogoTitle from '@/components/AppLogo.vue'
import { t } from '@/i18n'
import * as userReq from '@/network/details/user'
import router from '@/router'
import { DrawerMenu } from '@/types'
import { local } from '@/utils/user'
import { computed, defineModel, ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'

const baseRoutePath = computed(() => {
  const path = route.matched[route.matched.length - 1]?.path
  const pathParts = path?.split('/') ?? []
  const lastPart = pathParts[pathParts.length - 1]
  if (lastPart.includes(':')) {
    return pathParts.slice(0, pathParts.length - 1).join('/')
  }
  return route.path
})

const showDrawer = defineModel<boolean>()
const { mobile } = useDisplay()
const showAppLogoTitle = ref<boolean>(mobile.value)
const route = useRoute()

const gotoPage = (menu: DrawerMenu) => {
  let path = menu.route!
  const params = menu.defaultParams
  if (params) {
    path += `/${params.join('/')}`
  }
  router.push(path)
}

const items = computed<DrawerMenu[]>(() => [
  { text: t('nav.connectionInfo'), value: 'connection', icon: 'mdi-connection', route: '/admin/connection' },
  {
    text: t('nav.planManagement'),
    value: 'plan',
    icon: 'mdi-bookshelf',
    children: [
      { text: t('nav.plans'), value: 'plans', icon: 'mdi-bookmark-multiple', route: '/admin/plans' },
      { text: t('nav.keyManagement'), value: 'keys', icon: 'mdi-cloud-key-outline', route: '/admin/planKeys', defaultParams: [0] },
      { text: t('nav.keyVerification'), value: 'verifyKey', icon: 'mdi-key-alert-outline', route: '/admin/verifyKey' },
    ],
  },
  { text: t('nav.settings'), value: 'setting', icon: 'mdi-cogs', route: '/admin/setting' },
  { text: t('nav.logs'), value: 'log', icon: 'mdi-flag', route: '/admin/log' },
])

const openedGroup = computed(() => {
  return items.value.filter(item => item.children)
    .filter(item => item.children!.map(child => child.route).includes(baseRoutePath.value))
    .map(item => item.value)
})

watchEffect(() => {
  showAppLogoTitle.value = mobile.value
})

const logout = () => {
  userReq.logout().finally(() => {
    local.clearUser()
    router.push('/login')
  })
}
</script>
