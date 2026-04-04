<script setup lang="ts">
import Error404 from '@/components/Error404.vue'
import { t } from '@/i18n'
import router from '@/router'
import { ref } from 'vue'

const seconds = ref(5)
const timer = ref(0)

const timeSub = () => {
  timer.value = setTimeout(() => {
    seconds.value -= 1
    if (seconds.value <= 0) {
      router.push('/login')
    } else {
      timeSub()
    }
  }, 1000)
}

timeSub()
</script>

<template>
  <div>
    <Error404 />
    <div class="text-center">
      {{ t('notFound.redirectToLogin', { seconds }) }}
    </div>
  </div>
</template>
