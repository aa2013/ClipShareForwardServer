<script setup lang="ts">
import { t } from '@/i18n'
</script>

<template>
  <div class="mt-5">
    <div class="h-[1px] bg-[#efefef]" />
    <footer class="text-center mt-5 text-[#3C3C43C6] font-weight-medium">
      {{ t('footer.copyright', { year: new Date().getFullYear() }) }}
    </footer>
  </div>
</template>
