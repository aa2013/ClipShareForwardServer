<template>
  <v-card :elevation="elevation" class="border duration-300 rounded-[12px] px-2"
          @mouseover="elevation=3" @mouseleave="elevation=0">
    <v-skeleton-loader :loading="!plan" type="article,actions">
      <div/>
    </v-skeleton-loader>
    <div v-if="plan" class="flex p-3 text-xl items-center justify-between">
      <div class="flex items-center">
        <v-icon color="primary" size="25">mdi-bookmark</v-icon>
        {{ plan?.name }}
      </div>
      <v-btn icon="mdi-arrow-right" color="primary" variant="text" @click="router.push(`/admin/planKeys/${plan.id}`)"/>
    </div>
    <div v-if="plan" class="relative">
      <v-card-text class="ga-3 flex flex-col">
        <div class="flex ga-2 items-center">
          <v-icon color="primary">mdi-cellphone-link</v-icon>
          <div>设备同时连接数:</div>
          <div v-if="plan?.deviceLimit">{{ plan?.deviceLimit }}</div>
          <v-chip v-else color="primary" label density="compact" size="small">
            <v-icon>mdi-all-inclusive</v-icon>
          </v-chip>
          <div v-if="plan?.deviceLimit">台</div>
        </div>
        <div class="flex ga-2 items-center">
          <v-icon color="orange">mdi-car-speed-limiter</v-icon>
          <div>限速:</div>
          <div v-if="plan?.rate">{{ plan?.rate / 1024 }}</div>
          <v-chip v-else color="primary" label density="compact" size="small">
            <v-icon>mdi-all-inclusive</v-icon>
          </v-chip>
          <div v-if="plan?.rate">KB/s</div>
        </div>
        <div class="flex ga-2 items-center">
          <v-icon color="primary">mdi-weather-cloudy-clock</v-icon>
          <div>有效期:</div>
          <div v-if="plan?.lifespan">{{ plan?.lifespan / (24 * 60 * 60) }}</div>
          <v-chip v-else color="primary" label density="compact" size="small">
            <v-icon>mdi-all-inclusive</v-icon>
          </v-chip>
          <div v-if="plan?.lifespan">天</div>
        </div>
        <div v-if="plan?.remark">
          <div>
            备注：
          </div>
          <div>
            {{ plan?.remark }}
          </div>
        </div>
      </v-card-text>
      <div class="absolute right-0 bottom-0 mb-1 mr-3"
           v-if="!readonly">
        <v-tooltip :text="fold?'展开':'收起'">
          <template v-slot:activator="{props}">
            <v-btn :icon="fold?'mdi-chevron-down':'mdi-chevron-up'" variant="text"
                   v-bind="props" color="primary"
                   @click="fold=!fold"
                   size="small"/>
          </template>
        </v-tooltip>
      </div>
    </div>
    <v-slide-y-reverse-transition v-if="!readonly">
      <div v-show="!fold">
        <v-divider opacity="1"/>
        <v-card-actions>
          <v-btn :color="plan?.enable?'primary':undefined" prepend-icon="mdi-creation-outline" @click="genKey"
                 :disable="!plan?.enable">
            生成密钥
          </v-btn>
          <v-spacer/>
          <v-btn prepend-icon="mdi-square-edit-outline"
                 @click="onEditClick" :disable="!plan?.enable"
                 :color="plan?.enable?'primary':undefined">
            修改
          </v-btn>
          <v-btn :color="plan?.enable?'red':'primary'"
                 @click="switchStatus" :loading="switching"
                 :prepend-icon="plan?.enable?'mdi-stop-circle-outline':'mdi-check-circle-outline'">
            点击{{ plan?.enable ? '停用' : '启用' }}
          </v-btn>
        </v-card-actions>
      </div>
    </v-slide-y-reverse-transition>
    <gen-plan-keys-dialog v-model="showGenKeyDialog" :id="plan?.id"/>
  </v-card>
</template>
<script setup lang="ts">
import {PlanType} from "@/types";
import {ref} from "vue";
import * as planReq from '@/network/details/plan'
import {useGlobalSnackbar} from "@/stores/snackbar";
import router from "@/router";
import GenPlanKeysDialog from "@/components/plan/GenPlanKeysDialog.vue";

const plan = defineModel<PlanType>()
defineProps<{
  readonly?: boolean
}>()
const fold = ref<boolean>(true)
const elevation = ref<number>(0)
const emit = defineEmits(['onEditClick']);
const switching = ref<boolean>(false)
const {showSnackbar} = useGlobalSnackbar()
const showGenKeyDialog = ref<boolean>(false)
const onEditClick = () => {
  if (!plan.value?.enable) return
  // 触发事件并传递数据
  emit('onEditClick',);
};
const switchStatus = () => {
  switching.value = true
  const status = plan.value!.enable
  planReq.updateStatus(plan.value!.id!, !status).then(res => {
    showSnackbar({
      text: `${status ? '停用' : '启用'}${res ? '成功' : '失败'}`
    }, !res)
    plan.value!.enable = !status
  }).finally(() => {
    switching.value = false
  })
}
const genKey = () => {
  if (!plan.value?.enable) return
  showGenKeyDialog.value = true
}
</script>

<style scoped>

</style>
