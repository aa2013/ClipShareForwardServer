<template>
  <v-menu transition="scale-transition">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :icon="compact ? 'mdi-translate' : undefined"
        variant="text"
        color="primary"
        :title="t('language.label')"
      >
        <template v-if="!compact">
          <v-icon start icon="mdi-translate" />
          {{ currentLocaleLabel }}
        </template>
      </v-btn>
    </template>

    <v-list min-width="160" density="compact">
      <v-list-item
        v-for="item in localeOptions"
        :key="item.value"
        :active="item.value === currentLocale"
        :color="item.value === currentLocale ? 'primary' : undefined"
        @click="setLocale(item.value)"
      >
        <v-list-item-title>{{ item.label }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { currentLocale, localeOptions, setLocale, t } from '@/i18n'

withDefaults(defineProps<{
  compact?: boolean
}>(), {
  compact: true,
})

const currentLocaleLabel = computed(() => {
  return localeOptions.find(item => item.value === currentLocale.value)?.label ?? currentLocale.value
})
</script>
