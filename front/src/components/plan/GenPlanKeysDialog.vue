<template>
  <v-dialog v-model="show" max-width="400" :persistent="true">
    <v-card class="rounded-lg flex flex-col h-[100%] overflow-auto">
      <v-card-title class="flex items-center">
        <v-icon color="primary" size="small" class="mr-2">
          mdi-plus
        </v-icon>
        {{ t('plans.generateKeysTitle') }}
      </v-card-title>
      <v-form class="flex-1 flex flex-col overflow-auto" v-model="isFormValid">
        <div class="pa-5 flex-1 overflow-auto ga-2 flex flex-col">
          <v-text-field
            variant="outlined"
            :label="t('plans.generateCount')"
            v-model.number="submitData!.size"
            color="primary"
            density="comfortable"
            type="number"
            :rules="[
              value => value >= 0 || t('validation.cannotBeNegative'),
              value => /^\d+$/.test(String(value)) || t('validation.cannotBeDecimal'),
              value => value <= 2000 || t('validation.maximum', { value: 2000 }),
            ]"
          />
        </div>
        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn @click="onCancel">
            {{ t('common.cancel') }}
          </v-btn>

          <v-btn @click.prevent="onOk" :loading="loading" color="primary" type="submit">
            {{ t('common.generate') }}
          </v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { t } from '@/i18n'
import * as planReq from '@/network/details/plan'
import { useGlobalSnackbar } from '@/stores/snackbar'
import { GenPlanKeyInfo } from '@/types'
import { ref } from 'vue'

const { showSnackbar } = useGlobalSnackbar()

const show = defineModel<boolean>()
const props = defineProps<{
  id?: string,
}>()
const loading = ref(false)
const submitData = ref<GenPlanKeyInfo>({} as GenPlanKeyInfo)
const isFormValid = ref<boolean>(false)

const onCancel = () => {
  show.value = false
  loading.value = false
}

const onOk = () => {
  loading.value = true
  if (!isFormValid.value) {
    loading.value = false
    return
  }
  const size = submitData.value.size
  planReq.generatePlanKeys(props.id!, size).then(res => {
    showSnackbar({
      text: t('common.generateResult', { result: res ? t('common.success') : t('common.operationFailed') }),
    }, !res)
    if (res) {
      show.value = false
    }
  }).finally(() => {
    loading.value = false
  })
}
</script>
