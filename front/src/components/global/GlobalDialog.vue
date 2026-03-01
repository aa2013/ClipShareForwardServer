<template>
  <v-dialog
    v-model="dialogShow"
    max-width="400"
    :persistent="globalDialogProps.persistent"
  >
    <v-card>
      <v-card-title class="flex items-center">
        <v-icon :color="globalDialogProps.iconColor" size="small"
                v-if="globalDialogProps.icon" class="mr-2">
          {{ globalDialogProps.icon }}
        </v-icon>
        {{ globalDialogProps.title }}
      </v-card-title>
      <v-card-text>
        <div v-html="globalDialogProps.msg"/>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn @click="onDialogClose" :color="globalDialogProps.cancelBtnColor"
               v-show="globalDialogProps.showCancelBtn">
          {{ globalDialogProps.cancelBtnText }}
        </v-btn>

        <v-btn @click="onDialogOk" :color="globalDialogProps.okBtnColor" :loading="!!isOkLoading">
          {{ globalDialogProps.okBtnText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import {useGlobalDialog} from "@/stores/dialog";
import {storeToRefs} from "pinia";

const {dialogShow, globalDialogProps, isOkLoading} = storeToRefs(useGlobalDialog())
const {onDialogClose, onDialogOk} = useGlobalDialog()
</script>
<style scoped>

</style>
