<template>
  <v-app id="inspire">
    <AppDrawer v-model="drawer" v-if="mobile" />
    <v-app-bar elevation="0" class="border-b">
      <v-app-bar-nav-icon @click="drawer = !drawer" v-if="mobile" />
      <v-app-bar-title>
        <AppLogoTitle />
      </v-app-bar-title>
      <v-btn flat icon="mdi-github" @click="gotoGithub" />
      <LanguageSelector />

      <v-menu transition="scale-transition">
        <template #activator="{ props }">
          <v-btn flat :icon="currentThemeIcon" color="primary" v-bind="props" />
        </template>
        <v-list class="w-[140px]">
          <v-list-item
            nav
            :active="currentTheme === item.mode"
            :key="item.mode"
            :value="item.mode"
            v-for="item in themeModes"
            :color="currentTheme === item.mode ? 'primary' : undefined"
            @click="() => setTheme(item.mode)"
          >
            <v-list-item-title>
              <v-icon class="mr-2">
                {{ item.icon }}
              </v-icon>
              {{ item.name }}
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>
    <AppDrawer v-model="drawer" v-if="!mobile" />
    <v-main>
      <v-container class="flex flex-col h-full">
        <router-view class="flex-1" />
        <Footer />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import AppDrawer from '@/components/AppDrawer.vue'
import AppLogoTitle from '@/components/AppLogo.vue'
import Footer from '@/components/Footer.vue'
import LanguageSelector from '@/components/LanguageSelector.vue'
import { t } from '@/i18n'
import { ThemeMode } from '@/types'
import { useLocalTheme } from '@/stores/theme'
import { storeToRefs } from 'pinia'
import { computed, ref, watchEffect } from 'vue'
import { useDisplay } from 'vuetify'

const { mobile } = useDisplay()
const { currentTheme, currentThemeIcon } = storeToRefs(useLocalTheme())
const { setTheme } = useLocalTheme()
const drawer = ref<boolean>(!mobile.value)

const themeModes = computed<ThemeMode[]>(() => [
  { mode: 'auto', name: t('theme.auto'), icon: 'mdi-brightness-auto' },
  { mode: 'light', name: t('theme.light'), icon: 'mdi-brightness-5' },
  { mode: 'dark', name: t('theme.dark'), icon: 'mdi-brightness-2' },
])

const gotoGithub = () => {
  window.open('https://github.com/aa2013/ClipShareForwardServerWeb', '_blank')
}

watchEffect(() => {
  drawer.value = !mobile.value
})
</script>
