<script setup lang="ts">
import { t } from '@/i18n'
import { TrafficTrends } from '@/types'
import { defineModel } from 'vue'

const trafficTrends = defineModel<TrafficTrends>()
</script>

<template>
  <div class="flex">
    <v-btn
      :variant="trafficTrends === 'network' ? 'tonal' : 'text'"
      @click="trafficTrends = 'network'"
      :color="trafficTrends === 'network' ? 'primary' : ''"
      prepend-icon="mdi-speedometer"
      density="comfortable"
    >
      {{ t('trends.network') }}
    </v-btn>
    <v-btn
      :variant="trafficTrends === 'usage' ? 'tonal' : 'text'"
      @click="trafficTrends = 'usage'"
      :color="trafficTrends === 'usage' ? 'primary' : ''"
      prepend-icon="mdi-waves"
      class="ml-2"
      density="comfortable"
    >
      {{ t('trends.usage') }}
    </v-btn>
    <v-btn
      :variant="trafficTrends === 'connection' ? 'tonal' : 'text'"
      @click="trafficTrends = 'connection'"
      :color="trafficTrends === 'connection' ? 'primary' : ''"
      prepend-icon="mdi-connection"
      class="ml-2"
      density="comfortable"
    >
      {{ t('trends.connection') }}
    </v-btn>
  </div>
</template>
