<template>
  <div>
    <v-card :color="!isSameLoginSettings ? 'primary' : undefined" variant="tonal" class="mx-auto border my-5">
      <v-card-title>
        <v-icon icon="mdi-account-key" class="mr-2" />
        <span class="font-weight-black">
          {{ t('settings.loginSettings') }}
        </span>
      </v-card-title>
      <v-divider opacity="1" />

      <v-card-text class="pt-4">
        <v-skeleton-loader :loading="!loginSettings" type="list-item-two-line,button">
          <v-form class="size-full" v-model="formValid.loginSettings" @submit.prevent="() => updateConfig(loginSettings, 'loginSettings')">
            <v-text-field
              :label="t('settings.loginExpiredSeconds')"
              prepend-icon="mdi-clock-time-eight-outline"
              variant="underlined"
              type="number"
              v-model.number="loginSettings!.loginExpiredSeconds"
              :rules="nonNegativeIntegerRules"
            />
            <v-btn
              type="submit"
              :color="isSameLoginSettings ? 'grey' : 'primary'"
              variant="flat"
              :disabled="isSameLoginSettings"
              :loading="submitLoading.loginSettings"
              prepend-icon="mdi-content-save"
              class="my-2"
            >
              {{ t('common.save') }}
            </v-btn>
          </v-form>
        </v-skeleton-loader>
      </v-card-text>
    </v-card>

    <v-card :color="!isSamePublicModeConfig ? 'primary' : undefined" variant="tonal" class="mx-auto border my-5">
      <v-card-title>
        <v-icon icon="mdi-account-key" class="mr-2" />
        <span class="font-weight-black">
          {{ t('settings.modeSelection') }}
        </span>
      </v-card-title>
      <v-divider opacity="1" />

      <v-card-text class="pt-4">
        <v-skeleton-loader :loading="!publicMode && typeof publicMode !== 'boolean'" type="list-item-two-line,button">
          <v-form class="size-full" v-model="formValid.publicMode" @submit.prevent="() => updateConfig({ publicMode }, 'publicMode')">
            <v-switch color="primary" hide-details inset v-model="publicMode">
              <template #prepend>
                <v-icon icon="mdi-car-speed-limiter" class="mr-4" />
                {{ t('common.publicMode') }}
              </template>
            </v-switch>
            <v-btn
              type="submit"
              :color="isSamePublicModeConfig ? 'grey' : 'primary'"
              variant="flat"
              :disabled="isSamePublicModeConfig"
              :loading="submitLoading.publicMode"
              prepend-icon="mdi-content-save"
              class="my-2"
            >
              {{ t('common.save') }}
            </v-btn>
          </v-form>
        </v-skeleton-loader>
      </v-card-text>
    </v-card>

    <v-card :color="!isSameFileTransferLimit ? 'primary' : undefined" variant="tonal" class="mx-auto border my-5">
      <v-card-title>
        <v-icon icon="mdi-file-arrow-left-right-outline" class="mr-2" />
        <span class="font-weight-black">
          {{ t('settings.fileTransferLimit') }}
        </span>
      </v-card-title>
      <v-divider opacity="1" />
      <v-card-text class="pt-4">
        <v-skeleton-loader :loading="!fileTransferLimit" type="text,list-item-two-line,button">
          <v-form class="size-full" v-model="formValid.fileTransferLimit" @submit.prevent="() => updateConfig(fileTransferLimit, 'fileTransferLimit')">
            <v-switch color="primary" hide-details inset v-model="fileTransferLimit!.fileTransferEnabled">
              <template #prepend>
                <v-icon icon="mdi-car-speed-limiter" class="mr-4" />
                {{ t('settings.allowFileTransfer') }}
              </template>
            </v-switch>
            <v-text-field
              v-model.number="fileTransferLimit!.fileTransferRateLimit"
              :label="t('settings.fileTransferRateLimit')"
              prepend-icon="mdi-speedometer"
              variant="underlined"
              type="number"
              :rules="nonNegativeIntegerRules"
            />
            <v-btn
              type="submit"
              :color="isSameFileTransferLimit ? 'grey' : 'primary'"
              variant="flat"
              :disabled="isSameFileTransferLimit"
              :loading="submitLoading.fileTransferLimit"
              prepend-icon="mdi-content-save"
              class="my-2"
            >
              {{ t('common.save') }}
            </v-btn>
          </v-form>
        </v-skeleton-loader>
      </v-card-text>
    </v-card>

    <v-card :color="!isSameUnlimitedDevices ? 'primary' : undefined" variant="tonal" class="mx-auto border">
      <v-card-title>
        <v-icon icon="mdi-cellphone-link" class="mr-2" />
        <span class="font-weight-black">
          {{ t('settings.whitelistDevices') }}
        </span>
      </v-card-title>
      <v-divider opacity="1" />

      <v-card-text class="pt-4">
        <v-skeleton-loader :loading="!unlimitedDevices" type="paragraph,paragraph,button">
          <v-form class="size-full" v-model="formValid.unlimitedDevices" @submit.prevent="() => updateConfig({ unlimitedDevices }, 'unlimitedDevices')">
            <v-list>
              <v-list-item v-for="(device, i) in unlimitedDevices" :key="i">
                <div>
                  <v-chip class="my-4 flex items-center" label color="primary">
                    <v-icon icon="mdi-music-accidental-sharp" />
                    {{ t('settings.deviceCard', { index: i + 1 }) }}
                  </v-chip>
                  <v-btn size="small" class="ml-2" color="#ff4081" variant="tonal" @click="() => onDeleteUnlimitedDevice(i)">
                    <v-icon>mdi-trash-can-outline</v-icon>
                  </v-btn>
                </div>
                <v-text-field
                  v-model="device.id"
                  :label="t('common.deviceId')"
                  prepend-icon="mdi-identifier"
                  variant="underlined"
                  :rules="[
                    value => !!value || t('validation.required'),
                    value => {
                      if (unlimitedDevices!.some((item, idx) => item.id === value && idx !== i)) {
                        return t('validation.deviceIdExists')
                      }
                      return true
                    },
                  ]"
                />
                <v-text-field
                  v-model="device.name"
                  :label="t('common.deviceName')"
                  prepend-icon="mdi-label-outline"
                  variant="underlined"
                  :rules="[
                    value => !!value || t('validation.required'),
                  ]"
                />
                <v-text-field
                  v-model="device.desc"
                  :label="t('common.deviceDescription')"
                  prepend-icon="mdi-information-outline"
                  variant="underlined"
                />
              </v-list-item>
              <v-btn variant="text" color="primary" prepend-icon="mdi-plus" @click="unlimitedDevices?.push({} as UnlimitedDevice)">
                {{ t('settings.addDevice') }}
              </v-btn>
            </v-list>
            <v-btn
              type="submit"
              :color="isSameUnlimitedDevices ? 'grey' : 'primary'"
              variant="flat"
              :disabled="isSameUnlimitedDevices"
              :loading="submitLoading.unlimitedDevices"
              prepend-icon="mdi-content-save"
              class="my-2"
            >
              {{ t('common.save') }}
            </v-btn>
          </v-form>
        </v-skeleton-loader>
      </v-card-text>
    </v-card>

    <v-card :color="!isSameLogConfig ? 'primary' : undefined" variant="tonal" class="mx-auto border my-5">
      <v-card-title>
        <v-icon icon="mdi-flag" class="mr-2" />
        <span class="font-weight-black">
          {{ t('settings.logs') }}
        </span>
      </v-card-title>
      <v-divider opacity="1" />
      <v-card-text class="pt-4">
        <v-skeleton-loader :loading="!logConfig" type="text,list-item-two-line,button">
          <v-form class="size-full" v-model="formValid.log" @submit.prevent="() => updateConfig({ log: logConfig }, 'log')">
            <v-text-field
              v-model.number="logConfig!.memoryBufferSize"
              :label="t('settings.logMemoryBufferSize')"
              prepend-icon="mdi-speedometer"
              variant="underlined"
              type="number"
              :rules="nonNegativeIntegerRules"
            />
            <v-btn
              type="submit"
              :color="isSameLogConfig ? 'grey' : 'primary'"
              variant="flat"
              :disabled="isSameLogConfig"
              :loading="submitLoading.log"
              prepend-icon="mdi-content-save"
              class="my-2"
            >
              {{ t('common.save') }}
            </v-btn>
          </v-form>
        </v-skeleton-loader>
      </v-card-text>
    </v-card>
  </div>
</template>
<script setup lang="ts">
import { t } from '@/i18n'
import * as configReq from '@/network/details/config'
import { useGlobalDialog } from '@/stores/dialog'
import { useGlobalSnackbar } from '@/stores/snackbar'
import { FileTransferLimit, LogConfig, LoginSettings, SysConfig, UnlimitedDevice } from '@/types'
import { computed, onMounted, ref } from 'vue'

const { showGlobalDialog } = useGlobalDialog()
const { showSnackbar } = useGlobalSnackbar()

const loginSettings = ref<LoginSettings>()
const fileTransferLimit = ref<FileTransferLimit>()
const logConfig = ref<LogConfig>()
const unlimitedDevices = ref<UnlimitedDevice[]>()
const originConfig = ref<SysConfig>()
const publicMode = ref<boolean>()

const nonNegativeIntegerRules = [
  (value: number) => !!value || t('validation.required'),
  (value: number) => value >= 0 || t('validation.cannotBeNegative'),
  (value: number) => /^\d+$/.test(String(value)) || t('validation.cannotBeDecimal'),
]

const formValid = ref<Record<string, boolean>>({
  loginSettings: false,
  fileTransferLimit: false,
  unlimitedDevices: false,
  log: false,
  publicMode: false,
})

const submitLoading = ref<Record<string, boolean>>({
  loginSettings: false,
  fileTransferLimit: false,
  unlimitedDevices: false,
  log: false,
  publicMode: false,
})

const isSameLoginSettings = computed(() => {
  if (!loginSettings.value) return true
  const originData = {
    loginExpiredSeconds: originConfig.value?.loginExpiredSeconds,
  }
  return JSON.stringify(loginSettings.value) === JSON.stringify(originData)
})

const isSameFileTransferLimit = computed(() => {
  if (!fileTransferLimit.value) return true
  const originData = {
    fileTransferEnabled: originConfig.value?.fileTransferEnabled,
    fileTransferRateLimit: originConfig.value?.fileTransferRateLimit,
  }
  return JSON.stringify(fileTransferLimit.value) === JSON.stringify(originData)
})

const isSameUnlimitedDevices = computed(() => {
  if (!fileTransferLimit.value) return true
  const originData = originConfig.value?.unlimitedDevices ?? []
  return JSON.stringify(unlimitedDevices.value) === JSON.stringify(originData)
})

const isSameLogConfig = computed(() => {
  if (!logConfig.value) return true
  const originData = originConfig.value?.log ?? {}
  return JSON.stringify(logConfig.value) === JSON.stringify(originData)
})

const isSamePublicModeConfig = computed(() => {
  if (!publicMode.value && typeof publicMode.value !== 'boolean') return true
  const originData = originConfig.value?.publicMode ?? false
  return (publicMode.value ?? false) === originData
})

const loadConfigs = () => {
  configReq.getConfigs().then(configs => {
    originConfig.value = configs
    fileTransferLimit.value = {
      fileTransferEnabled: configs.fileTransferEnabled,
      fileTransferRateLimit: configs.fileTransferRateLimit,
    }

    loginSettings.value = {
      loginExpiredSeconds: configs.loginExpiredSeconds,
    }
    unlimitedDevices.value = JSON.parse(JSON.stringify(configs.unlimitedDevices ?? []))
    logConfig.value = JSON.parse(JSON.stringify(configs.log))
    publicMode.value = configs.publicMode
  })
}

const updateConfig = (config: any, group: string) => {
  const isValid = formValid.value[group]
  if (!isValid) {
    return
  }
  submitLoading.value[group] = true
  configReq.updateConfig(config).then(res => {
    showSnackbar({
      text: t('common.saveResult', { result: res ? t('common.success') : t('common.operationFailed') }),
    }, !res)
    if (res) {
      Object.assign(originConfig.value!, config)
      originConfig.value = JSON.parse(JSON.stringify(originConfig.value!))
    }
  }).finally(() => {
    submitLoading.value[group] = false
  })
}

onMounted(() => {
  setTimeout(loadConfigs, 500)
})

const onDeleteUnlimitedDevice = (idx: number) => {
  const device = unlimitedDevices.value![idx]
  if (device.name !== '' && !device.name && device.id !== '' && !device.id) {
    unlimitedDevices.value = unlimitedDevices.value!.filter((_, i) => i !== idx)
    return
  }
  showGlobalDialog({
    persistent: true,
    title: t('settings.deleteWhitelistDeviceTitle'),
    msg: t('settings.deleteWhitelistDeviceMessage', {
      name: device.name ?? '<unknown>',
      id: device.id ?? '<unknown>',
    }),
    showCancelBtn: true,
    onOk() {
      unlimitedDevices.value = unlimitedDevices.value!.filter((_, i) => i !== idx)
    },
  })
}
</script>
