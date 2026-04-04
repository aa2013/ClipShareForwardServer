<template>
  <v-card class="pa-2 h-[100%] flex flex-col" flat>
    <v-card-title class="d-flex align-center">
      <v-icon icon="mdi-flag" class="mr-2" color="primary" />
      {{ t('logs.title') }}

      <v-spacer />

      <v-text-field
        v-model="search"
        density="compact"
        :label="t('common.search')"
        prepend-inner-icon="mdi-magnify"
        variant="solo-filled"
        flat
        hide-details
        single-line
      />
    </v-card-title>
    <v-divider />
    <div class="flex-1-0 overflow-auto h-0" ref="tableContainer">
      <v-data-table-virtual
        width="100%"
        :headers="tableHeaders"
        :items="logs"
        :no-data-text="t('common.noData')"
        :fixed-header="true"
        sticky
        :height="tableHeight || tableContainer?.offsetHeight"
        :search="search"
      >
        <template #item.level="{ item }">
          <v-chip
            :color="LogLevelColors[item.level as LogLevel] ?? 'grey'"
            :text="item.level"
            class="text-uppercase"
            size="small"
            label
          />
        </template>
        <template #bottom></template>
      </v-data-table-virtual>
    </div>
    <div class="flex flex-wrap justify-between items-center ga-4 pa-2">
      <div>
        {{ t('common.totalItems', { count: logs.length }) }}
      </div>
    </div>
  </v-card>
</template>
<script setup lang="ts">
import { t } from '@/i18n'
import * as logReq from '@/network/details/log'
import { Log, LogLevel } from '@/types'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const LogLevelColors: Record<LogLevel, string> = {
  info: 'primary',
  warn: 'orange',
  error: 'red',
}

const tableHeaders = computed<any[]>(() => [
  {
    title: t('common.time'),
    align: 'start',
    key: 'time',
    minWidth: '200px',
  },
  {
    title: t('logs.level'),
    align: 'start',
    key: 'level',
    minWidth: '115px',
  },
  {
    title: t('logs.content'),
    align: 'start',
    key: 'content',
  },
])

const logs = ref<Log[]>([])
const search = ref<string>()
let logTimer: number

const fetchLogs = () => {
  const logLen = logs.value.length
  const lastTime = logs.value[logLen - 1]?.time
  logReq.getLogs(lastTime).then(list => {
    logs.value.push(...list.map(item => {
      const logParts = item.log.split(' ')
      const level = logParts[0].replace(/[[\]]/g, '')
      return {
        time: item.time,
        level,
        content: logParts.slice(3).join(' '),
      } as Log
    }))
  })
}

const tableContainer = ref<HTMLElement>()
const tableHeight = ref<number>()

const onWindowResize = () => {
  tableHeight.value = tableContainer.value?.offsetHeight
}

onMounted(() => {
  fetchLogs()
  logTimer = setInterval(fetchLogs, 1000)
  window.addEventListener('resize', onWindowResize)
})

onUnmounted(() => {
  clearInterval(logTimer)
  window.removeEventListener('resize', onWindowResize)
})
</script>
