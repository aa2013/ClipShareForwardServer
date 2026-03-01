<template>
  <div class="w-[350px] mx-auto">
    <v-form @submit.prevent="verifyKey" v-model="isFormValid"
            class="my-5">
      <v-text-field v-model="inputKey" label="密钥"
                    color="primary" variant="outlined"
                    placeholder="请输入密钥"
                    :loading="loading"
                    :disabled="loading"
                    :rules="[
                      v=>!!v||'Required.'
                    ]"/>
      <v-btn prepend-icon="mdi-check" type="submit"
             color="primary" block :loading="loading">
        提交验证
      </v-btn>
    </v-form>
    <plan v-if="loading || planType" v-model="planType" readonly/>
  </div>
</template>
<script setup lang="ts">
import {ref} from "vue";
import * as planKeyReq from '@/network/details/planKeys'
import {PlanType} from "@/types";
import {useGlobalSnackbar} from "@/stores/snackbar";
import Plan from "@/components/plan/Plan.vue";

const {showSnackbar} = useGlobalSnackbar()

const inputKey = ref<string>("")
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
          text: "该密钥不存在",
        }, true)
      }
    })
    .finally(() => {
      loading.value = false
    })
}
</script>
<style scoped>

</style>
