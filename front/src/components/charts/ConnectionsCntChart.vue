<template>
  <BaseChart id="connection-cnt-chart" :options="options" />
</template>
<script setup lang="ts">
import { t } from '@/i18n'
import { ConnectionChartData } from '@/types'
import * as echarts from 'echarts'
import { computed, defineModel } from 'vue'
import BaseChart from '@/components/charts/BaseChart.vue'

const data = defineModel<ConnectionChartData[]>({
  get(value) {
    if (!value) {
      return []
    }
    return value
  },
})

const options = computed<echarts.EChartsOption>(() => {
  const totalCntData = data.value!.map(item => item.baseCnt + item.dataSyncCnt + item.fileSyncCnt)
  const baseCntData = data.value!.map(item => item.baseCnt)
  const dataSyncCntData = data.value!.map(item => item.dataSyncCnt)
  const fileSyncCntData = data.value!.map(item => item.fileSyncCnt)
  return {
    tooltip: {
      show: true,
      trigger: 'axis',
    },
    legend: {
      data: [
        t('charts.totalConnections'),
        t('charts.deviceConnections'),
        t('charts.dataSyncConnections'),
        t('charts.fileSyncConnections'),
      ],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.value!.map(item => item.time),
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: Math.ceil(Math.max(...totalCntData) * 1.2),
    },
    series: [
      {
        name: t('charts.totalConnections'),
        type: 'line',
        data: totalCntData,
      },
      {
        name: t('charts.deviceConnections'),
        type: 'line',
        data: baseCntData,
      },
      {
        name: t('charts.dataSyncConnections'),
        type: 'line',
        data: dataSyncCntData,
      },
      {
        name: t('charts.fileSyncConnections'),
        type: 'line',
        data: fileSyncCntData,
      },
    ],
  }
})
</script>
