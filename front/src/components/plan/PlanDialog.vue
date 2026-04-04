<template>
  <v-dialog v-model="show" max-width="400" :persistent="true">
    <v-card class="rounded-lg flex flex-col h-[100%] overflow-auto">
      <v-card-title class="flex items-center">
        <v-icon color="primary" size="small" class="mr-2">
          {{ isEdit ? 'mdi-square-edit-outline' : 'mdi-plus' }}
        </v-icon>
        {{ t('plans.addOrEditPlan', { action: isEdit ? t('common.edit') : t('common.add') }) }}
      </v-card-title>
      <v-form class="flex-1 flex flex-col overflow-auto" v-model="isFormValid">
        <div class="pa-5 flex-1 overflow-auto ga-2 flex flex-col">
          <v-text-field
            variant="outlined"
            :label="t('plans.planName')"
            v-model="submitData!.name"
            color="primary"
            density="comfortable"
            :rules="[v => !!v || t('validation.required')]"
          />
          <v-text-field
            variant="outlined"
            :label="t('plans.rateLimitLabel')"
            v-model.number="submitData!.rate"
            color="primary"
            density="comfortable"
            type="number"
            :rules="numericRules"
          />
          <v-text-field
            variant="outlined"
            :label="t('plans.validityDays')"
            v-model.number="submitData!.lifespan"
            color="primary"
            density="comfortable"
            type="number"
            :rules="numericRules"
          />
          <v-text-field
            variant="outlined"
            :label="t('plans.deviceLimitLabel')"
            v-model.number="submitData!.deviceLimit"
            color="primary"
            density="comfortable"
            type="number"
            :rules="numericRules"
          />
          <v-textarea variant="outlined" :label="t('plans.remark')" density="comfortable" color="primary" rows="2" v-model="submitData!.remark" />
        </div>
        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn @click="onCancelClick">
            {{ t('common.cancel') }}
          </v-btn>

          <v-btn @click.prevent="onOkClick" :loading="loading" color="primary" type="submit">
            {{ t('plans.submitAction', { action: isEdit ? t('common.edit') : t('common.add') }) }}
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
import { PlanType } from '@/types'
import { computed, ref, watch } from 'vue'

const { showSnackbar } = useGlobalSnackbar()

const show = defineModel<boolean>()
const loading = ref(false)
const props = defineProps<{
  data?: PlanType,
}>()
const emit = defineEmits(['onOk'])
const daySeconds = 24 * 60 * 60
const isEdit = computed(() => !!props.data)
const submitData = ref<PlanType>({} as PlanType)
const isFormValid = ref<boolean>(false)
const numericRules = [
  (value: number) => value >= 0 || t('validation.cannotBeNegative'),
  (value: number) => /^\d+$/.test(String(value)) || t('validation.cannotBeDecimal'),
]

watch(() => props.data, (newV) => {
  if (newV) {
    const lf = newV.lifespan ? newV.lifespan / daySeconds : undefined
    const rate = newV.rate ? newV.rate / 1024 : undefined
    submitData.value = { ...newV, lifespan: lf, rate }
    return
  }
  submitData.value = {} as PlanType
}, { immediate: true })

const onCancelClick = () => {
  show.value = false
  loading.value = false
}

const onOkClick = () => {
  loading.value = true
  if (!isFormValid.value) {
    loading.value = false
    return
  }
  let promise
  const data = { ...submitData.value! }
  data.lifespan = data.lifespan ? data.lifespan * daySeconds : undefined
  data.rate = data.rate ? data.rate * 1024 : undefined
  if (isEdit.value) {
    promise = planReq.editPlan(data)
  } else {
    promise = planReq.addPlan(data)
  }
  promise.then(res => {
    showSnackbar({
      text: isEdit.value
        ? t('common.updateResult', { result: res ? t('common.success') : t('common.operationFailed') })
        : t('common.addResult', { result: res ? t('common.success') : t('common.operationFailed') }),
    }, !res)
    if (res) {
      show.value = false
      emit('onOk')
    }
  }).finally(() => {
    loading.value = false
  })
}
</script>
