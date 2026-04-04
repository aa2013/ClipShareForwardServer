<template>
  <div class="w-[350px] mx-auto">
    <v-form @submit.prevent="verifyKey" v-model="isFormValid" class="my-5">
      <v-text-field
        v-model="inputKey"
        :label="t('verifyKey.inputLabel')"
        color="primary"
        variant="outlined"
        :placeholder="t('verifyKey.inputPlaceholder')"
        :loading="loading"
        :disabled="loading"
        :rules="[
          v => !!v || t('validation.required'),
        ]"
      />
      <v-btn prepend-icon="mdi-check" type="submit" color="primary" block :loading="loading">
        {{ t('verifyKey.submit') }}
      </v-btn>
    </v-form>
    <Plan v-if="loading || planType" v-model="planType" readonly />
  </div>
</template>
<script setup lang="ts">
import { t } from '@/i18n'
import * as planKeyReq from '@/network/details/planKeys'
import { useGlobalSnackbar } from '@/stores/snackbar'
import { PlanType } from '@/types'
import { ref } from 'vue'
import Plan from '@/components/plan/Plan.vue'

const { showSnackbar } = useGlobalSnackbar()

const inputKey = ref<string>('')
const isFormValid = ref<boolean>(false)
const loading = ref<boolean>(false)
const planType = ref<PlanType>()

const verifyKey = () => {
  if (!isFormValid.value) return
  planType.value = undefined
  loading.value = true
  planKeyReq.verify(inputKey.value)
    .then(plan => {
      planType.value = plan
      if (!plan) {
        showSnackbar({
          text: t('verifyKey.notFound'),
        }, true)
      }
    })
    .finally(() => {
      loading.value = false
    })
}
</script>
