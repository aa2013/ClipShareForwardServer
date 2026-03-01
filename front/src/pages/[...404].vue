<script setup lang="ts">

import Error404 from "@/components/Error404.vue";
import {ref} from "vue";
import router from "@/router";

const seconds = ref(5)
const timer = ref(0)
const timeSub = () => {
  timer.value = setTimeout(() => {
    seconds.value -= 1;
    if (seconds.value <= 0) {
      router.push("/login")
    } else {
      timeSub()
    }
  }, 1000)
}
timeSub()
</script>

<template>
  <div>
    <Error404/>
    <div class="text-center">
      {{ seconds }} 秒后跳转登录...
    </div>
  </div>
</template>

<style scoped>

</style>
