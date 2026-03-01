<template>
  <v-dialog
    v-model="show"
    max-width="400"
    :persistent="true"
  >
    <v-card class="rounded-lg flex flex-col h-[100%] overflow-auto">
      <v-card-title class="flex items-center">
        <v-icon color="primary" size="small" class="mr-2">
          {{ isEdit ? 'mdi-square-edit-outline' : 'mdi-plus' }}
        </v-icon>
        {{ isEdit ? '修改' : '新增' }}套餐
      </v-card-title>
      <v-form class="flex-1 flex flex-col overflow-auto" v-model="isFormValid">
        <div class="pa-5 flex-1 overflow-auto ga-2 flex flex-col">
          <v-text-field variant="outlined" label="套餐名称" v-model="submitData!.name"
                        color="primary" density="comfortable"
                        :rules="[v=>!!v||'Required.']"/>
          <v-text-field variant="outlined" label="限速速率" v-model.number="submitData!.rate"
                        color="primary" density="comfortable" type="number"
                        :rules="[
                            value=>value>=0||'Cannot be negative.',
                            value=>/^\d+$/.test(value)||'Cannot be decimal.'
                        ]"/>
          <v-text-field variant="outlined" label="有效期（天）" v-model.number="submitData!.lifespan"
                        color="primary" density="comfortable" type="number"
                        :rules="[
                            value=>value>=0||'Cannot be negative.',
                            value=>/^\d+$/.test(value)||'Cannot be decimal.'
                        ]"/>
          <v-text-field variant="outlined" label="同时使用设备数量" v-model.number="submitData!.deviceLimit"
                        color="primary" density="comfortable" type="number"
                        :rules="[
                            value=>value>=0||'Cannot be negative.',
                            value=>/^\d+$/.test(value)||'Cannot be decimal.'
                        ]"/>
          <v-textarea variant="outlined" label="备注" density="comfortable"
                      color="primary" rows="2" v-model="submitData!.remark"/>
        </div>
        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn @click="onCancelClick">
            取消
          </v-btn>

          <v-btn @click.prevent="onOkClick" :loading="loading" color="primary" type="submit">
            提交{{ isEdit ? '修改' : '新增' }}
          </v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import {PlanType} from "@/types";
import {computed, ref, watch} from "vue";
import * as planReq from '@/network/details/plan'
import {useGlobalSnackbar} from "@/stores/snackbar";

const {showSnackbar} = useGlobalSnackbar()

const show = defineModel<boolean>()
const loading = ref(false)
const props = defineProps<{
  data?: PlanType,
}>()
const emit = defineEmits(['onOk']);
const daySeconds = 24 * 60 * 60
const isEdit = computed(() => !!props.data)
const submitData = ref<PlanType>({} as PlanType)
const isFormValid = ref<boolean>(false)
watch(() => props.data, (newV) => {
  if (newV) {
    const lf = newV.lifespan ? newV.lifespan / daySeconds : undefined
    const rate = newV.rate ? newV.rate / 1024 : undefined
    submitData.value = {...newV, lifespan: lf, rate}
    return
  } else {
    submitData.value = {} as PlanType
  }
})
const onCancelClick = () => {
  show.value = false
  loading.value = false
}
const onOkClick = () => {
  // show.value = false
  loading.value = true
  if (!isFormValid.value) {
    loading.value = false
    return
  }
  let promise;
  const data = {...submitData.value!}
  data.lifespan = data.lifespan ? data.lifespan * daySeconds : undefined
  data.rate = data.rate ? data.rate * 1024 : undefined
  if (isEdit.value) {
    promise = planReq.editPlan(data)
  } else {
    promise = planReq.addPlan(data)
  }
  promise.then(res => {
    showSnackbar({
      text: `${isEdit.value ? '更新' : '添加'}${res ? '成功' : '失败'}`
    }, !res)
    if (res) {
      show.value = false
      emit('onOk',);
    }
  }).finally(() => {
    loading.value = false
  })
}
</script>
<style scoped>

</style>
