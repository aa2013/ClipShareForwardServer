<template>
  <BaseChart id="flow-usage-chart" :options="options"/>
</template>
<script setup lang="ts">

import BaseChart from "@/components/charts/BaseChart.vue";
import * as echarts from "echarts";
import {CallbackDataParams, TopLevelFormatterParams} from "echarts/types/dist/shared";
import {computed, defineModel} from "vue";
import {TrafficChartData} from "@/types";
import * as commonUtil from '@/utils/common'

const data = defineModel<TrafficChartData[]>({
  get(value) {
    if (!value) {
      return []
    }
    return value
  }
})
const options = computed<echarts.EChartsOption>(() => {
  const trafficData = data.value!.map(item => item.traffic)
  const maxVal = Math.ceil(Math.max(...trafficData) * 1.2)
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
            ${marker} ${seriesName}&nbsp;&nbsp;&nbsp;&nbsp;${commonUtil.numberToSizeStr(value)}
        `
        })
        return html
      }
    },
    legend: {
      data: ['当前连接使用量']
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
      name:'单位：'+commonUtil.numberToSizeStr(maxVal,maxVal).slice(-2),
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
        name: '当前连接使用量',
        type: 'line',
        data: trafficData
      },
    ]
  }
})
</script>

<style scoped>

</style>
