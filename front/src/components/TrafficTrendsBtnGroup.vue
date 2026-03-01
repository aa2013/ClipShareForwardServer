<script setup lang="ts">
import {defineModel} from 'vue';
import {TrafficTrends} from "@/types";

const trafficTrends = defineModel<TrafficTrends>()

</script>

<template>
  <div class="flex">
    <v-btn :variant="trafficTrends=='network'?'tonal':'text'"
           @click="trafficTrends='network'"
           :color="trafficTrends=='network'?'primary':''"
           prepend-icon="mdi-speedometer"
           density="comfortable">
      网络
    </v-btn>
    <v-btn :variant="trafficTrends=='usage'?'tonal':'text'"
           @click="trafficTrends='usage'"
           :color="trafficTrends=='usage'?'primary':''"
           prepend-icon="mdi-waves"
           class="ml-2"
           density="comfortable">
      使用量
    </v-btn>
    <v-btn :variant="trafficTrends=='connection'?'tonal':'text'"
           @click="trafficTrends='connection'"
           :color="trafficTrends=='connection'?'primary':''"
           prepend-icon="mdi-connection"
           class="ml-2"
           density="comfortable">
      连接数
    </v-btn>
  </div>
</template>

<style scoped>

</style>
