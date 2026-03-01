<script setup lang="ts">
import emptyPng from '@/assets/empty.png'
</script>

<template>
  <div class="flex flex-col text-center justify-center size-full">
    <img alt="empty-content"
      class="mx-auto"
      :src="emptyPng"
    />
    <span class="text-[#a3a3a3] text-2xl ">
      Empty Data
    </span>
  </div>
</template>

<style scoped>

</style>
