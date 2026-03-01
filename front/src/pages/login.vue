<template>
  <v-app>
    <div class="size-full login-bg flex justify-center items-center">
      <div class="bg-drop"></div>
      <div class="position-relative w-[350px] rounded-[12px] p-[16px] flex flex-col"
           :style="{'background':theme.global.name.value=='dark'?'#0f1417':'white'}"
      >
        <div class="h-[68px] flex justify-center items-center">
          <img :src="logo" alt="logo" width="48px" height="48px"/>
          <div class="text-[transparent] hero-text font-bold text-2xl">
            Forward Server
          </div>
        </div>

        <v-form ref="form" v-model="isFormValid" class="my-[16px]" @submit.prevent="login" :disabled="logging">
          <v-text-field
            v-model="formData.username"
            :rules="[v=>!!v || '请输入用户名']"
            label="用户名"
            required
            class="mb-2"
            color="primary"
            variant="outlined"
          />
          <v-text-field
            v-model="formData.password"
            :rules="[v=>!!v || '请输入密码']"
            label="密码"
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
            @mouseover="loginBtnElevation=4"
            @mouseleave="loginBtnElevation=0"
            :loading="logging"
          >
            登录
          </v-btn>
        </v-form>
      </div>
    </div>
  </v-app>
</template>
<script setup lang="ts">
import logo from "@/assets/logo.png";
import {ref} from "vue";
import {useTheme} from "vuetify";
import * as  userReq from '@/network/details/user'
import router from "@/router";

const theme = useTheme();
const isFormValid = ref(false)
const loginBtnElevation = ref(0)
const logging = ref(false)
const formData = ref({
  username: '',
  password: ''
})
const login = () => {
  if (!isFormValid.value) return;
  logging.value = true
  userReq.login(formData.value.username, formData.value.password)
    .then((res) => {
      localStorage.setItem('user', JSON.stringify({
        'token': res.token,
        'username': formData.value.username
      }))
      router.push('/admin/connection')
    }).finally(() => {
    logging.value = false
  })
}
</script>

<style scoped>
.login-bg {
  background: url("../assets/login-bg.jpeg");
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
