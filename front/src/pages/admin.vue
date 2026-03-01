<template>
  <v-app id="inspire">
    <AppDrawer v-model="drawer" v-if="mobile"/>
    <v-app-bar elevation="0" class="border-b">
      <v-app-bar-nav-icon @click="drawer = !drawer" v-if="mobile"/>
      <v-app-bar-title>
        <AppLogoTitle/>
      </v-app-bar-title>
      <v-btn flat icon="mdi-github" @click="gotoGithub"/>

      <v-menu transition="scale-transition">
        <template v-slot:activator="{ props }">
          <v-btn flat :icon="currentThemeIcon" color="primary" v-bind="props"/>
        </template>
        <v-list class="w-[120px]">
          <v-list-item nav :active="currentTheme===item.mode" :key="item.mode" :value="item.mode"
                       v-for="(item) in themeModes" :color="currentTheme===item.mode?'primary':undefined"
                       @click="()=>setTheme(item.mode)">
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
    <AppDrawer v-model="drawer" v-if="!mobile"/>
    <v-main>
      <v-container class="flex flex-col h-full">
        <router-view class="flex-1"/>
        <Footer/>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import AppDrawer from "@/components/AppDrawer.vue";
import AppLogoTitle from "@/components/AppLogo.vue";
import {useDisplay} from 'vuetify'
import {ref, watchEffect} from "vue";
import Footer from "@/components/Footer.vue";
import {ThemeMode} from "@/types";
import {useLocalTheme} from "@/stores/theme";
import {storeToRefs} from "pinia";


const {mobile} = useDisplay()
const {currentTheme, currentThemeIcon} = storeToRefs(useLocalTheme())
const {setTheme} = useLocalTheme()
const drawer = ref<boolean>(!mobile.value)
const themeModes = ref<ThemeMode[]>([
  {mode: 'auto', name: '自动', icon: 'mdi-brightness-auto'},
  {mode: 'light', name: '亮色模式', icon: 'mdi-brightness-5'},
  {mode: 'dark', name: '暗色模式', icon: 'mdi-brightness-2'},
])
const gotoGithub = () => {
  window.open('https://github.com/aa2013/ClipShareForwardServerWeb', '_blank')
}
watchEffect(() => {
  drawer.value = !mobile.value
})
</script>
