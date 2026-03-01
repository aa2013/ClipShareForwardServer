<template>
  <v-snackbar
    v-model="snackbarShow"
    :timeout="globalSnackbarProps.timeout"
    :color="globalSnackbarProps.color"
    :location="globalSnackbarProps.location"
  >
    <div class="flex items-center">
      <v-icon class="mr-2" v-if="globalSnackbarProps.icon">{{globalSnackbarProps.icon}}</v-icon>
      {{ globalSnackbarProps.text }}
    </div>
    <template v-slot:actions>
    <v-btn
      v-if="globalSnackbarProps.showAction"
      :color="globalSnackbarProps.actionColor"
      variant="text"
      @click="onActionClick"
    >
      {{ globalSnackbarProps.actionText }}
    </v-btn>
  </template>
  </v-snackbar>
</template>
<script setup lang="ts">
import {storeToRefs} from "pinia";
import {useGlobalSnackbar} from "@/stores/snackbar";

const {snackbarShow, globalSnackbarProps} = storeToRefs(useGlobalSnackbar())
const {onActionClick} = useGlobalSnackbar()
</script>
<style scoped>

</style>
