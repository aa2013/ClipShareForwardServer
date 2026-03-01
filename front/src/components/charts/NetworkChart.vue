<template>
  <BaseChart id="network-chart" :options="options"/>
</template>
<script setup lang="ts">

import BaseChart from "@/components/charts/BaseChart.vue";
import * as echarts from "echarts";
import {CallbackDataParams, TopLevelFormatterParams} from "echarts/types/dist/shared";
import {computed, defineModel} from "vue";
import {NetworkChartData} from "@/types";
import * as commonUtil from '@/utils/common'

const data = defineModel<NetworkChartData[]>({
  get(value) {
    if (!value) {
      return []
    }
    return value
  }
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
          const marker = item.marker;
          const seriesName = item.seriesName;
          const value = parseInt(item.value?.toString() ?? '0');
          return `
            <br>
            ${marker} ${seriesName}&nbsp;&nbsp;&nbsp;&nbsp;${commonUtil.numberToSizeStr(value)}/s
        `
        })
        return html
      }
    },
    legend: {
      data: ['总速度', '数据同步', '文件同步']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.value!.map(item => item.time)
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: maxVal,
      name: '单位：' + commonUtil.numberToSizeStr(maxVal, maxVal).slice(-2)+"/s",
      axisLabel: {
        formatter: function (value) {
          const label = commonUtil.numberToSizeStr(value, maxVal)
          const match = label.match(/\d+/)
          return parseInt(match![0]) + ""
        }
      }
    },
    series: [
      {
        name: '总速度',
        type: 'line',
        data: totalSpeedData
      },
      {
        name: '数据同步',
        type: 'line',
        data: dataSyncSpeedData
      },
      {
        name: '文件同步',
        type: 'line',
        data: fileSyncSpeedData
      },
    ]
  }
})
</script>

<style scoped>

</style>
