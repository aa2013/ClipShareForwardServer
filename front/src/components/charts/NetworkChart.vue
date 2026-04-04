<template>
  <BaseChart id="network-chart" :options="options" />
</template>
<script setup lang="ts">
import BaseChart from '@/components/charts/BaseChart.vue'
import { t } from '@/i18n'
import { NetworkChartData } from '@/types'
import * as commonUtil from '@/utils/common'
import * as echarts from 'echarts'
import { CallbackDataParams, TopLevelFormatterParams } from 'echarts/types/dist/shared'
import { computed, defineModel } from 'vue'

const data = defineModel<NetworkChartData[]>({
  get(value) {
    if (!value) {
      return []
    }
    return value
  },
})

const options = computed<echarts.EChartsOption>(() => {
  const totalSpeedData = data.value!.map(item => item.dataSync + item.fileSync)
  const dataSyncSpeedData = data.value!.map(item => item.dataSync)
  const fileSyncSpeedData = data.value!.map(item => item.fileSync)
  const maxVal = Math.ceil(Math.max(...totalSpeedData) * 1.2)
  return {
    tooltip: {
      show: true,
      trigger: 'axis',
      formatter: function (p: TopLevelFormatterParams) {
        const params = p as (CallbackDataParams & { axisValueLabel: string })[]
        let html = params[0].axisValueLabel
        html += params.map(item => {
          const marker = item.marker
          const seriesName = item.seriesName
          const value = parseInt(item.value?.toString() ?? '0')
          return `
            <br>
            ${marker} ${seriesName}&nbsp;&nbsp;&nbsp;&nbsp;${commonUtil.numberToSizeStr(value)}/s
        `
        })
        return html
      },
    },
    legend: {
      data: [t('charts.totalSpeed'), t('charts.dataSync'), t('charts.fileSync')],
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
      max: maxVal,
      name: t('charts.unitPerSecond', { unit: commonUtil.numberToSizeStr(maxVal, maxVal).slice(-2) }),
      axisLabel: {
        formatter: function (value) {
          const label = commonUtil.numberToSizeStr(value, maxVal)
          const match = label.match(/\d+/)
          return `${parseInt(match![0])}`
        },
      },
    },
    series: [
      {
        name: t('charts.totalSpeed'),
        type: 'line',
        data: totalSpeedData,
      },
      {
        name: t('charts.dataSync'),
        type: 'line',
        data: dataSyncSpeedData,
      },
      {
        name: t('charts.fileSync'),
        type: 'line',
        data: fileSyncSpeedData,
      },
    ],
  }
})
</script>
