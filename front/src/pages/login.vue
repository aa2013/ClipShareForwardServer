<template>
  <v-app>
    <div class="size-full login-bg flex justify-center items-center">
      <div class="bg-drop"></div>
      <div
        class="position-relative w-[350px] rounded-[12px] p-[16px] flex flex-col"
        :style="{ background: theme.global.name.value === 'dark' ? '#0f1417' : 'white' }"
      >
        <div class="absolute right-2 top-2">
          <LanguageSelector :compact="false" />
        </div>

        <div class="h-[68px] flex justify-center items-center">
          <img :src="logo" alt="logo" width="48px" height="48px" />
          <div class="text-[transparent] hero-text font-bold text-2xl">
            {{ t('loginPage.title') }}
          </div>
        </div>

        <v-form ref="form" v-model="isFormValid" class="my-[16px]" @submit.prevent="login" :disabled="logging">
          <v-text-field
            v-model="formData.username"
            :rules="[v => !!v || t('loginPage.usernameRequired')]"
            :label="t('common.username')"
            required
            class="mb-2"
            color="primary"
            variant="outlined"
          />
          <v-text-field
            v-model="formData.password"
            :rules="[v => !!v || t('loginPage.passwordRequired')]"
            :label="t('common.password')"
            required
            color="primary"
            variant="outlined"
            type="password"
          />

          <v-btn
            class="mt-4"
            size="large"
            block
            type="submit"
            color="primary"
            :elevation="loginBtnElevation"
            @mouseover="loginBtnElevation = 4"
            @mouseleave="loginBtnElevation = 0"
            :loading="logging"
          >
            {{ t('loginPage.loginButton') }}
          </v-btn>
        </v-form>
      </div>
    </div>
  </v-app>
</template>
<script setup lang="ts">
import logo from '@/assets/logo.png'
import LanguageSelector from '@/components/LanguageSelector.vue'
import { t } from '@/i18n'
import router from '@/router'
import * as userReq from '@/network/details/user'
import { ref } from 'vue'
import { useTheme } from 'vuetify'

const theme = useTheme()
const isFormValid = ref(false)
const loginBtnElevation = ref(0)
const logging = ref(false)
const formData = ref({
  username: '',
  password: '',
})

const login = () => {
  if (!isFormValid.value) return
  logging.value = true
  userReq.login(formData.value.username, formData.value.password)
    .then((res) => {
      localStorage.setItem('user', JSON.stringify({
        token: res.token,
        username: formData.value.username,
      }))
      router.push('/admin/connection')
    }).finally(() => {
      logging.value = false
    })
}
</script>

<style scoped>
.login-bg {
  background: url('../assets/login-bg.jpeg');
}

.bg-drop {
  --blur: 10px;
  position: fixed;
  height: 100vh;
  width: 100vw;
  backdrop-filter: blur(var(--blur));
  -webkit-backdrop-filter: blur(var(--blur));
}
</style>
